package main

import (
    "bufio"
    "database/sql"
    "flag"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
    "time"

    _ "modernc.org/sqlite"

    playwright "github.com/playwright-community/playwright-go"
)

// ---------- Configuration ----------
const (
    maxTokens                = 1250
    defaultTokens            = 750
    sendWaitMs               = 7000
    maxRetries               = 3
    tokenCollectionTimeoutMs = 90000
    targetURL                = "https://chat.z.ai"
)

// ---------- Prompt user for token count ----------
func promptTokenCount() int {
    fmt.Printf("How many tokens to collect? [default: %d, max: %d] ", defaultTokens, maxTokens)
    reader := bufio.NewReader(os.Stdin)
    answer, _ := reader.ReadString('\n')
    answer = strings.TrimSpace(answer)

    n := defaultTokens
    if answer != "" {
        if parsed, err := strconv.Atoi(answer); err == nil && parsed > 0 {
            n = parsed
        }
    }
    if n > maxTokens {
        fmt.Printf("⚠️  Capping to max %d.\n", maxTokens)
        n = maxTokens
    }
    return n
}

// ---------- Selector wait result ----------
type selResult struct {
    loc playwright.Locator
    err error
}

// ---------- Single page collection attempt ----------
func tryCollect(page playwright.Page, total int) ([]string, error) {
    fmt.Printf("  🌐 Navigating to %s\n", targetURL)
    _, err := page.Goto(targetURL, playwright.PageGotoOptions{
        WaitUntil: playwright.WaitUntilStateDomcontentloaded,
        Timeout:   playwright.Float(60000),
    })
    if err != nil {
        return nil, fmt.Errorf("navigation failed: %w", err)
    }

    // ---------- Wait for model button + textarea in parallel ----------
    fmt.Println("  🔍 Locating UI elements in parallel...")

    modelBtnCh := make(chan selResult, 1)
    textareaCh := make(chan selResult, 1)

    go func() {
        loc, err := page.WaitForSelector("#model-selector-glm-4_7-button",
            playwright.PageWaitForSelectorOptions{Timeout: playwright.Float(15000)})
        modelBtnCh <- selResult{loc, err}
    }()

    go func() {
        loc, err := page.WaitForSelector("#chat-input",
            playwright.PageWaitForSelectorOptions{Timeout: playwright.Float(15000)})
        textareaCh <- selResult{loc, err}
    }()

    mbRes := <-modelBtnCh
    if mbRes.err != nil {
        return nil, fmt.Errorf("model button not found: %w", mbRes.err)
    }
    taRes := <-textareaCh
    if taRes.err != nil {
        return nil, fmt.Errorf("textarea not found: %w", taRes.err)
    }
    textarea := taRes.loc
    fmt.Println("  ✅ Model button & textarea found")

    // ---------- Fill textarea & click send ----------
    if err := textarea.Fill("__"); err != nil {
        return nil, fmt.Errorf("failed to fill textarea: %w", err)
    }
    fmt.Println(`  ✅ Textarea filled with "__"`)

    sendBtn, err := page.WaitForSelector("#send-message-button",
        playwright.PageWaitForSelectorOptions{Timeout: playwright.Float(5000)})
    if err != nil {
        return nil, fmt.Errorf("send button not found: %w", err)
    }
    if err := sendBtn.Click(); err != nil {
        return nil, fmt.Errorf("failed to click send: %w", err)
    }
    fmt.Println("  ✅ Send clicked")

    // ---------- Wait for token endpoint ----------
    fmt.Printf("  ⏳ Waiting %dms for token endpoint to initialize...\n", sendWaitMs)
    time.Sleep(time.Duration(sendWaitMs) * time.Millisecond)

    // ---------- Collect tokens with timeout ----------
    fmt.Println("  🚀 Collecting tokens...")
    t0 := time.Now()

    page.SetDefaultTimeout(float64(tokenCollectionTimeoutMs))

    jsExpr := `async (total) => {
        const out = new Array(total);
        for (let i = 0; i < total; i++) {
            const tok = window.z_um.getToken();
            out[i] = (tok && typeof tok.then === 'function') ? await tok : tok;
            if (i % 50 === 0) {
                await new Promise(r => setTimeout(r, 0));
            }
        }
        return out;
    }`

    result, err := page.Evaluate(jsExpr, total)
    if err != nil {
        return nil, fmt.Errorf("token collection failed: %w", err)
    }

    elapsed := time.Since(t0).Seconds()
    fmt.Printf("  ✅ Collected tokens in %.2fs\n", elapsed)

    // ---------- Convert result to []string ----------
    tokens := make([]string, 0, total)
    if arr, ok := result.([]interface{}); ok {
        for _, v := range arr {
            tokens = append(tokens, fmt.Sprintf("%v", v))
        }
    }

    return tokens, nil
}

// ---------- Collect a batch with retries (fresh page each attempt) ----------
func collectBatch(browser playwright.Browser, total int) ([]string, error) {
    var lastErr error

    for attempt := 1; attempt <= maxRetries; attempt++ {
        fmt.Printf("  🔄 Attempt %d of %d\n", attempt, maxRetries)

        page, err := browser.NewPage()
        if err != nil {
            lastErr = err
            continue
        }

        tokens, err := tryCollect(page, total)
        _ = page.Close() // close old page before next attempt

        if err == nil {
            return tokens, nil
        }

        lastErr = err
        fmt.Printf("  ❌ Attempt %d failed: %v\n", attempt, err)
        if attempt < maxRetries {
            fmt.Println("  ♻️  Retrying with a fresh page load...")
        }
    }

    return nil, lastErr
}

// ---------- Build SQLite database ----------
func buildSQLite(tokens []string, filePath string) error {
    os.Remove(filePath) // start fresh

    db, err := sql.Open("sqlite", filePath)
    if err != nil {
        return err
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        return err
    }

    if _, err := db.Exec("CREATE TABLE tokens (id INTEGER PRIMARY KEY, token TEXT);"); err != nil {
        return err
    }

    tx, err := db.Begin()
    if err != nil {
        return err
    }

    stmt, err := tx.Prepare("INSERT INTO tokens (id, token) VALUES (?, ?);")
    if err != nil {
        _ = tx.Rollback()
        return err
    }
    defer stmt.Close()

    for i, token := range tokens {
        if _, err := stmt.Exec(i, token); err != nil {
            _ = tx.Rollback()
            return err
        }
    }

    return tx.Commit()
}

// ---------- Main ----------
func main() {
    // --conc=N  (0 < N < 6, default 3)
    concFlag := flag.Int("conc", 3, "number of collection batches (1-5)")
    flag.Parse()

    conc := *concFlag
    if conc < 1 || conc > 5 {
        fmt.Fprintf(os.Stderr, "❌ --conc must be between 1 and 5 (got %d)\n", conc)
        os.Exit(1)
    }

    total := promptTokenCount()
    fmt.Printf("\n🎯 Collecting %d tokens per batch, %d batches total\n", total, conc)

    // ---------- Launch Playwright ----------
    fmt.Println("📦 Installing Playwright browsers (if needed)...")
    if err := playwright.Install(); err != nil {
        log.Fatalf("Failed to install Playwright: %v", err)
    }

    pw, err := playwright.Run()
    if err != nil {
        log.Fatalf("Failed to start Playwright: %v", err)
    }
    defer pw.Stop()

    browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
        Headless: playwright.Bool(true),
    })
    if err != nil {
        log.Fatalf("Failed to launch browser: %v", err)
    }
    defer browser.Close()

    // ---------- Batch loop: N starts at 1, runs until N == conc ----------
    var allTokens []string

    for n := 1; n <= conc; n++ {
        fmt.Printf("\n📦 Batch %d of %d\n", n, conc)

        tokens, err := collectBatch(browser, total)
        if err != nil {
            fmt.Printf("❌ Batch %d failed after all retries: %v\n", n, err)
            continue
        }

        // Merge with previously collected tokens
        allTokens = append(allTokens, tokens...)
        fmt.Printf("✅ Batch %d: collected %d tokens (running total: %d)\n",
            n, len(tokens), len(allTokens))

        // N == limit → build final database
        if n == conc {
            fmt.Println("\n📋 Reached batch limit — building final database...")
        }
    }

    if len(allTokens) == 0 {
        fmt.Fprintln(os.Stderr, "🚫 No tokens collected. Aborting.")
        os.Exit(1)
    }

    // ---------- Build & save SQLite ----------
    fmt.Printf("\n🗄️  Building SQLite database with %d total tokens...\n", len(allTokens))
    if err := buildSQLite(allTokens, "tokens.sqlite"); err != nil {
        fmt.Fprintf(os.Stderr, "❌ Failed to build database: %v\n", err)
        os.Exit(1)
    }

    if info, err := os.Stat("tokens.sqlite"); err == nil {
        fmt.Printf("✅ Saved: tokens.sqlite (%.1f KB)\n", float64(info.Size())/1024)
    }

    fmt.Println("\n🎉 Script finished successfully.")
}
