// cmd/token-collector/main.go
//
// Standalone binary that seeds tokens.sqlite with device tokens for the
// Z.AI bridge.
//
// HOW IT WORKS (2026-10-21 rewrite)
//
//	The browser is used exactly ONCE, briefly.  It boots the Aliyun-SAF captcha
//	SDK on chat.z.ai so the FeiLin device service issues this device a
//	*trusted* session (Log1) and the device uploads its collector fingerprint
//	(Log2/Log3).  The visit also yields the genuine collector payload.
//
//	Every device token is then minted in-process: the payload is refreshed
//	(per-session timestamps + the trailing probe-log entries), AES-128-CBC
//	sealed with the session key and assembled into the 5-field token.
//
//	Minting is pure computation, so there is nothing to batch, parallelise or
//	retry per token: one harvested session yields as many tokens as asked for.
//
//	Why the browser cannot be dropped: the risk layer trusts a device session
//	only when the Log1 that created it came from a client it accepts as a real
//	browser (the request bytes are byte-identical to the SDK's - it is the
//	*client* that is judged) and once the device registered its fingerprint for
//	that session.  See .assets/reports/aliyun-captcha-report.md.
//
// Build commands (pick one):
//
//   Linux/macOS (modern CPU, fully static, stripped):
//     CGO_ENABLED=0 GOAMD64=v3 go build -ldflags="-s -w" -trimpath -o token-collector ./cmd/token-collector
//
//   Windows PowerShell:
//     $env:CGO_ENABLED=0; $env:GOAMD64="v3"; go build -ldflags="-s -w" -trimpath -o token-collector.exe ./cmd/token-collector
//
//   Portable fallback (any CPU / any OS):
//     go build -ldflags="-s -w" -trimpath -o token-collector ./cmd/token-collector
//
// Usage:
//   ./token-collector                  # interactive prompt
//   ./token-collector --tokens 750
//   ./token-collector --session-max-age 30m   # reuse the session longer

package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mxschmitt/playwright-go"
	_ "modernc.org/sqlite" // pure-Go SQLite, no CGO needed
)

// ---------- Configuration ----------
const (
	MaxTokens     = 1500
	DefaultTokens = 850
	URL           = "https://chat.z.ai"
)

// ---------- Flags ----------
var (
	tokensFlag = flag.Int("tokens", 0, "number of device tokens to mint (0 = prompt)")
	headedFlag = flag.Bool("headed", true,
		"deprecated: the browser must run headful (headless device sessions are "+
			"rejected by the Aliyun risk layer), so the flag is ignored")
	noTUIFlag         = flag.Bool("no-tui", false, "disable TUI, use plain text output")
	exportSessionFlag = flag.String("export-session", "",
		"write the harvested device session + genuine payload template to this JSON file "+
			"(consumable by .assets/scripts/live_validate.py --session-mint)")
	sessionMaxAgeFlag = flag.Duration("session-max-age", 15*time.Minute,
		"how long a harvested device session may be reused before the browser is used again "+
			"(the device session has a server side lifetime; live traffic measured >= 35 min)")
)

// ---------- TUI (Bubble Tea) ----------
var spinnerChars = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// logCapture stores log lines in a ring buffer for the TUI to display.
type logCapture struct {
	mu     sync.Mutex
	lines  []string
	maxLen int
}

func (lc *logCapture) addLine(line string) {
	lc.mu.Lock()
	lc.lines = append(lc.lines, line)
	if len(lc.lines) > lc.maxLen {
		lc.lines = lc.lines[len(lc.lines)-lc.maxLen:]
	}
	lc.mu.Unlock()
}

func (lc *logCapture) Lines() []string {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	out := make([]string, len(lc.lines))
	copy(out, lc.lines)
	return out
}

// Global TUI state — lock-free atomics for reads from the TUI goroutine.
var (
	tuiLogCapture      = &logCapture{maxLen: 1000}
	tuiStatus          atomic.Value // string
	tuiTokensCollected atomic.Int64
	tuiStartTime       atomic.Value // time.Time
	tuiDone            atomic.Bool
	tuiErr             atomic.Value // error
)

// init: initialise TUI default state
func init() {
	tuiStatus.Store("Initializing...")
	tuiStartTime.Store(time.Now())
}

func tuiSetStatus(s string) { tuiStatus.Store(s) }

// Bubble Tea messages
type tickMsg time.Time
type doneMsg struct{ err error }

func tuiTick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

// tuiModel is the Bubble Tea model for the token collector TUI.
type tuiModel struct {
	width     int
	height    int
	logOffset int // lines scrolled up from bottom
	done      bool
	err       error
}

func (m tuiModel) Init() tea.Cmd { return tuiTick() }

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.logOffset++
		case "down", "j":
			if m.logOffset > 0 {
				m.logOffset--
			}
		case "g":
			m.logOffset = len(tuiLogCapture.Lines())
		case "G":
			m.logOffset = 0
		}
	case tickMsg:
		if tuiDone.Load() {
			m.done = true
			if v := tuiErr.Load(); v != nil {
				m.err = v.(error)
			}
			return m, tea.Quit
		}
		return m, tuiTick()
	case doneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	}
	return m, nil
}

func (m tuiModel) View() string {
	if m.height < 10 || m.width < 40 {
		return "Terminal too small (min 40x10). Resize or press q to quit."
	}

	// Color styles
	stTitle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("213"))
	stLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	stAcc := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	stWarn := lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)
	stDim := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	// Gather state from atomics (lock-free)
	status := "Initializing..."
	if v := tuiStatus.Load(); v != nil {
		status = v.(string)
	}
	tc := tuiTokensCollected.Load()

	var st time.Time
	if v := tuiStartTime.Load(); v != nil {
		st = v.(time.Time)
	}
	elapsed := time.Since(st).Round(time.Second)

	// --- Header / status block ---
	var hdr strings.Builder
	hdr.WriteString(stTitle.Render("🔑 Token Collector"))
	hdr.WriteByte('\n')

	if m.done {
		if m.err != nil {
			hdr.WriteString(stLabel.Render("Status: ") + stWarn.Render(fmt.Sprintf("ERROR: %v", m.err)))
		} else {
			hdr.WriteString(stLabel.Render("Status: ") + stAcc.Render("✅ COMPLETE"))
		}
	} else {
		sp := spinnerChars[int(time.Now().UnixMilli()/200)%len(spinnerChars)]
		hdr.WriteString(stLabel.Render("Status: ") + fmt.Sprintf("%s %s", sp, stAcc.Render(status)))
	}
	hdr.WriteByte('\n')

	// Stats line
	stats := fmt.Sprintf("%s %s",
		stLabel.Render("Tokens:"),
		stAcc.Render(fmt.Sprintf("%d", tc)))
	stats += fmt.Sprintf("  %s %s", stLabel.Render("Elapsed:"), stAcc.Render(elapsed.String()))
	hdr.WriteString(stats)
	hdr.WriteByte('\n')

	headerStr := hdr.String()
	headerLines := strings.Count(headerStr, "\n")

	// --- Log pane ---
	logH := m.height - headerLines - 4 // -4: sep + log header + sep + footer
	if logH < 2 {
		logH = 2
	}

	logs := tuiLogCapture.Lines()
	start := len(logs) - logH - m.logOffset
	if start < 0 {
		start = 0
	}
	end := start + logH
	if end > len(logs) {
		end = len(logs)
	}

	truncSt := lipgloss.NewStyle().MaxWidth(m.width)
	var logLines []string
	if len(logs) == 0 {
		logLines = []string{stDim.Render("(waiting for output...)")}
	} else if start < end {
		for _, l := range logs[start:end] {
			logLines = append(logLines, truncSt.Render(l))
		}
	} else {
		logLines = []string{stDim.Render("(no more logs)")}
	}
	logStr := strings.Join(logLines, "\n")

	// Separator and footer
	sep := stDim.Render(strings.Repeat("─", m.width))
	footer := stDim.Render(" ↑/↓ scroll  •  q quit")

	return headerStr + sep + "\n" + stLabel.Render("📋 Logs") + "\n" + logStr + "\n" + sep + "\n" + footer
}

// ---------- End TUI ----------

// ---------- Prompt user for integer ----------
func promptInt(reader *bufio.Reader, prompt string, def, max int) int {
	fmt.Print(prompt)
	line, err := reader.ReadString('\n')
	if err != nil {
		return def
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return def
	}
	n, err := strconv.Atoi(line)
	if err != nil || n <= 0 {
		fmt.Printf("⚠️  Invalid input, using default %d.\n", def)
		return def
	}
	if n > max {
		fmt.Printf("⚠️  Capping to max %d.\n", max)
		return max
	}
	return n
}

// ---------- Persistent SQLite store with tuned PRAGMAs ----------
// Opening/closing the DB per batch (as the original did) forces a full
// fsync and reparse of the schema every time. Keeping one connection
// open with WAL mode + large cache eliminates that overhead entirely.
type tokenStore struct {
	db   *sql.DB
	stmt *sql.Stmt
	mu   sync.Mutex // serialise writes (SQLite is single-writer)
}

func openTokenStore(dbPath string) (*tokenStore, error) {
	dsn := "file:" + filepath.ToSlash(dbPath) +
		"?_pragma=busy_timeout(10000)" +
		"&_pragma=journal_mode(WAL)" +
		"&_pragma=synchronous(NORMAL)" +
		"&_pragma=cache_size(-65536)" + // 64 MB page cache
		"&_pragma=temp_store(MEMORY)" +
		"&_pragma=mmap_size(268435456)" + // 256 MB mmap
		"&_pragma=wal_autocheckpoint(1000)"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	// SQLite serialises writes internally — one connection avoids
	// "database is locked" errors and avoids pool overhead.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	// Force the connection open so PRAGMAs take effect immediately.
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	// The bridge consumes the rows FIFO ("SELECT token FROM tokens ORDER BY id
	// LIMIT 1") and deletes them after use; no other column is read.
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS tokens (
        id    INTEGER PRIMARY KEY AUTOINCREMENT,
        token TEXT    NOT NULL
    )`); err != nil {
		db.Close()
		return nil, err
	}

	stmt, err := db.Prepare(`INSERT INTO tokens (token) VALUES (?)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return &tokenStore{db: db, stmt: stmt}, nil
}

func (ts *tokenStore) add(tokens []string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	tx, err := ts.db.Begin()
	if err != nil {
		return err
	}
	// Rollback is a no-op after Commit, so this is safe.
	defer tx.Rollback()

	// Bind the prepared statement to this transaction.
	txStmt := tx.Stmt(ts.stmt)
	defer txStmt.Close()

	for _, tok := range tokens {
		if _, err := txStmt.Exec(tok); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (ts *tokenStore) close() {
	ts.stmt.Close()
	ts.db.Close()
}

// ---------- Open the one-off harvest page ----------
// The harvest is a single visit: open chat.z.ai (so the page origin, referrer
// and TLS session look like a real visit) and let the captcha SDK register the
// device.  There is no fingerprint spoofing and no tracker filtering — the
// risk layer judges the client, not the UA, and the page is loaded once.
func newHarvestPage(browser playwright.Browser) (playwright.BrowserContext, playwright.Page, error) {
	ctx, err := browser.NewContext(playwright.BrowserNewContextOptions{})
	if err != nil {
		return nil, nil, err
	}
	page, err := ctx.NewPage()
	if err != nil {
		_ = ctx.Close()
		return nil, nil, err
	}
	return ctx, page, nil
}

// ---------- Mint the tokens ----------
//
// The browser is consulted only for the session itself (harvest.go); every
// token after that is produced in-process.  The harvest is the single flaky
// step, so it is the only thing worth retrying.
func mintDeviceTokens(minter *sessionMinter, total int, maxAge time.Duration) ([]string, error) {
	var err error
	for attempt := 1; attempt <= MaxHarvestAttempts; attempt++ {
		if err = minter.ensure(maxAge); err == nil {
			break
		}
		fmt.Printf("❌ Harvest attempt %d/%d failed: %v\n", attempt, MaxHarvestAttempts, err)
		if attempt < MaxHarvestAttempts {
			fmt.Println("♻️  Retrying the one-off browser visit...")
			minter.reset()
		}
	}
	if err != nil {
		return nil, err
	}

	tuiSetStatus(fmt.Sprintf("Minting %d device tokens...", total))
	fmt.Printf("⚙️  Minting %d device tokens from the harvested session (no further browser use)...\n", total)
	t0 := time.Now()
	toks, err := minter.mint(total)
	if err != nil {
		return nil, err
	}
	fmt.Printf("✅ Minted %d unique device tokens in %.2fs\n", len(toks), time.Since(t0).Seconds())
	return toks, nil
}

// ---------- Browser launch ----------
//
// The device session is only trusted when the risk layer accepts the CLIENT,
// and it accepts the client by looking at the browser itself - not at the
// request bytes (those are byte-identical to what we can send from Go) and
// not at the UA string.  Measured live:
//
//	headful Chromium            -> session trusted   -> tokens verify (T001)
//	headless Chromium (new)     -> session UNTRUSTED -> tokens fail (F001)
//	any pure HTTP client        -> session UNTRUSTED -> tokens fail (F001)
//
// So the browser always runs HEADED.  --enable-automation is stripped from
// Playwright's default args so the automation infobar flag never reaches the
// engine.
func launchBrowser(pw *playwright.Playwright) (playwright.Browser, error) {
	return pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless:          playwright.Bool(false),
		Args:              []string{"--disable-blink-features=AutomationControlled"},
		IgnoreDefaultArgs: []string{"--enable-automation"},
	})
}

// ---------- Core run logic ----------
func run(tokenCount int) error {
	// Install Playwright browsers (best-effort)
	tuiSetStatus("Installing Playwright...")
	fmt.Println("⏳ Ensuring Playwright Chromium browser is installed...")
	if err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"chromium"},
	}); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  playwright install: %v (continuing anyway)\n", err)
	}

	tuiSetStatus("Launching browser...")
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("playwright run: %w", err)
	}
	defer pw.Stop()

	browser, err := launchBrowser(pw)
	if err != nil {
		return fmt.Errorf("browser launch: %w", err)
	}
	defer browser.Close()

	// Database path — start fresh. Also nuke WAL/SHM sidecar files.
	dbPath := filepath.Join(".", "tokens.sqlite")
	_ = os.Remove(dbPath)
	_ = os.Remove(dbPath + "-wal")
	_ = os.Remove(dbPath + "-shm")

	// Open DB once and keep it — avoids per-batch open/close/fsync.
	ts, err := openTokenStore(dbPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer ts.close()

	// ONE browser visit registers the device; every token is then minted from
	// that session.  The session is re-harvested only when it ages past
	// --session-max-age.
	minter := &sessionMinter{browser: browser}

	tokens, err := mintDeviceTokens(minter, tokenCount, *sessionMaxAgeFlag)
	if err != nil {
		return err
	}

	if err := ts.add(tokens); err != nil {
		return fmt.Errorf("database insert: %w", err)
	}
	tuiTokensCollected.Store(int64(len(tokens)))

	if *exportSessionFlag != "" {
		if err := minter.exportSession(*exportSessionFlag); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  session export failed: %v\n", err)
		} else {
			fmt.Printf("  🧾 Session exported to %s\n", *exportSessionFlag)
		}
	}

	// ---------- Final summary ----------
	fmt.Printf("\n══════════════════════════════════════════\n")
	fmt.Printf("  ✅ ALL TOKENS MINTED\n")
	fmt.Printf("  📦 %d device tokens\n", len(tokens))
	fmt.Printf("  🌐 %s\n", minter.stats())
	if info, err := os.Stat(dbPath); err == nil {
		fmt.Printf("  💾 %s (%.1f KB)\n", dbPath, float64(info.Size())/1024.0)
	}
	fmt.Printf("══════════════════════════════════════════\n")

	return nil
}

// ---------- Main ----------
func main() {
	flag.Parse()

	if !*headedFlag {
		fmt.Fprintln(os.Stderr, "⚠️  --headed=false is ignored: the Aliyun risk layer rejects")
		fmt.Fprintln(os.Stderr, "    headless device sessions, so the tokens would not verify.")
	}

	reader := bufio.NewReader(os.Stdin)

	// ---------- Prompt for token count ----------
	tokenCount := *tokensFlag
	if tokenCount <= 0 {
		tokenCount = promptInt(reader,
			fmt.Sprintf("How many tokens to mint? [default: %d, max: %d] ", DefaultTokens, MaxTokens),
			DefaultTokens, MaxTokens)
	} else if tokenCount > MaxTokens {
		fmt.Printf("⚠️  Capping tokens to max %d.\n", MaxTokens)
		tokenCount = MaxTokens
	}

	// ---------- TUI setup (before the plan so the plan shows in the TUI logs) ----------
	useTUI := !*noTUIFlag
	var origStdout, origStderr *os.File
	var pipeWriter *os.File

	if useTUI {
		origStdout = os.Stdout
		origStderr = os.Stderr

		r, w, perr := os.Pipe()
		if perr != nil {
			fmt.Fprintf(os.Stderr, "pipe error: %v\n", perr)
			os.Exit(1)
		}
		pipeWriter = w
		os.Stdout = w
		os.Stderr = w

		// Goroutine: read piped output → logCapture ring buffer.
		// Non-blocking; uses 1MB scanner buffer for long lines.
		go func() {
			scanner := bufio.NewScanner(r)
			scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
			for scanner.Scan() {
				tuiLogCapture.addLine(scanner.Text())
			}
		}()

		tuiStartTime.Store(time.Now())
		tuiStatus.Store("Starting...")
	}

	// ---------- Plan ----------
	fmt.Printf("\n🎯 Plan: mint %d device tokens", tokenCount)
	fmt.Println()

	// ---------- Run ----------
	if !useTUI {
		// Plain text mode — no TUI
		if err := run(tokenCount); err != nil {
			fmt.Fprintf(os.Stderr, "\n🚫 Fatal error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("\n🎉 Script finished successfully.")
		return
	}

	// ---------- TUI mode ----------
	// tea.WithOutput(origStdout) sends TUI rendering to the real terminal
	// while fmt.Println goes to the pipe → logCapture.
	p := tea.NewProgram(tuiModel{},
		tea.WithAltScreen(),
		tea.WithOutput(origStdout),
	)

	go func() {
		err := run(tokenCount)
		if err == nil {
			tuiSetStatus("Complete!")
		}
		tuiDone.Store(true)
		if err != nil {
			tuiErr.Store(err)
		}
		pipeWriter.Close() // EOF → scanner goroutine exits
		p.Send(doneMsg{err: err})
	}()

	if _, err := p.Run(); err != nil {
		os.Stdout = origStdout
		os.Stderr = origStderr
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}

	os.Stdout = origStdout
	os.Stderr = origStderr

	if v := tuiErr.Load(); v != nil {
		fmt.Fprintf(os.Stderr, "\n🚫 Fatal error: %v\n", v.(error))
		os.Exit(1)
	}

	if !tuiDone.Load() {
		fmt.Println("\n⏹️  Interrupted by user.")
		os.Exit(0)
	}

	fmt.Println("\n🎉 Script finished successfully.")
}
