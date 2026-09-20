// Entry point of the Z.AI bridge (package zbridge).
//
// Run() is what the thin root main.go calls: it parses the CLI flags, opens
// the token database, starts the captcha cache (agent mode), attaches the
// session pool, serves NewHandler() and blocks until CTRL+C / SIGTERM —
// then drains in-flight requests and clears every still-pooled chat session
// on Z.AI before exiting (ported from the DeepseekFreeAPI reference).
//
// NewHandler() is exported on its own so integration tests can drive the
// full HTTP surface (all routes + auth + CORS) without starting a listener.

package zbridge

import (
    "context"
    "errors"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

// ============================================================================
// ENTRY POINT
// ============================================================================

// NewHandler assembles the bridge's complete HTTP surface: every route with
// the auth and CORS middleware applied. Used by Run and by the blackbox
// integration tests in tests/.
func NewHandler() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/", dashboardHandler)
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/status", statusHandler)
    mux.HandleFunc("/v1/models", authMiddleware(modelsHandler))
    mux.HandleFunc("/models", authMiddleware(modelsHandler2))
    mux.HandleFunc("/v1/chat/completions", authMiddleware(chatCompletionsHandler))
    mux.HandleFunc("/v1/messages", authMiddleware(anthropicMessagesHandler))
    mux.HandleFunc("/features", authMiddleware(featuresHandler))
    mux.HandleFunc("/admin/stats", statsHandler)
    mux.HandleFunc("/admin/health", healthHandler)
    mux.HandleFunc("/admin/clients", clientsHandler)
    mux.HandleFunc("/inject.js", injectHandler)
    mux.HandleFunc("/stop", authMiddleware(stopHandler))
    // Hot-swap the active tokens.sqlite without restarting the server.
    // Auth-protected like every other mutating endpoint.
    mux.HandleFunc("/sqlite", authMiddleware(sqliteSwapHandler))

    return corsMiddleware(mux)
}

// Run starts the bridge server and blocks until a fatal error or a
// termination signal. Called from the root package's main().
func Run() {
    flag.StringVar(&dbPath, "db-path", "tokens.sqlite", "Path to SQLite database")
    flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
    flag.BoolVar(&config.AgentMode, "agent-mode", config.AgentMode, "Enable agent mode: translate tools & roles for Z.AI compatibility (modern shim by default)")
    flag.StringVar(&config.AgentModeVariant, "agent-mode-variant", config.AgentModeVariant, "Agent mode shim variant: modern (default, XML-sectioned prompt) or legacy ([ROLE: ...] rewrite)")
    flag.StringVar(&config.AgentModeLevel, "agent-mode-level", config.AgentModeLevel, "Agent mode repair level: empty (default, stock parser only) or ultra (buffer + ToolParserLLM malformed-tool repair; requires --agent-mode)")
    flag.BoolVar(&config.ForceCPU, "force-cpu", config.ForceCPU, "Allow ToolParserLLM ultra repair on CPU-only hosts (degraded path; CPU inference is NOT recommended, GPU+SGLang is required otherwise)")
    flag.BoolVar(&config.SyncMode, "sync-mode", config.SyncMode, "Legacy synchronous session flow: reuse one sticky chat up to SESSION_REUSE_COUNT instead of drawing from the pre-warmed session pool")
    flag.IntVar(&config.SessionReuseCount, "session-reuse-count", config.SessionReuseCount, "How many requests one chat session serves before it is deleted on Z.AI and replaced (SESSION_REUSE_COUNT, default 10; 1 = throwaway per request)")
    flag.Parse()

    if _, err := os.Stat(dbPath); err != nil {
        log.Println("Captcha db not found! Please run the token collector first (cmd/token-collector)")
        os.Exit(1)
    }

    logInfo("Starting with db-path='" + dbPath + "' verbose=true")

    if err := initDB(); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to open database: %v\n", err)
        os.Exit(1)
    }
    defer closeDB()

    gRunning.Store(true)

    if config.AgentMode {
        go captchaCache.Run()
        logInfo("Agent mode: Captcha background cache started")
        if config.agentModern() {
            logInfo("Agent mode variant: MODERN (XML-sectioned prompt shim, tolerant marker/payload parsing)")
        } else {
            logInfo("Agent mode variant: LEGACY ([ROLE: ...] message rewrite shim)")
        }
    }

    // ── ToolParserLLM ultra repair (strict opt-in: --agent-mode + --agent-mode-level=ultra) ──
    // Dormant otherwise: no load, no warm-up, no reference to ToolParserLLM.
    if config.UltraEnabled() {
        log.Printf("[ToolParserLLM] ultra repair enabled (%s) — ensuring base+LoRA assets and SGLang backend", ultraGateDebug())
        if err := EnsureToolParserLLMForUltra(context.Background()); err != nil {
            fmt.Fprintf(os.Stderr, "ToolParserLLM ultra startup failed: %v\n", err)
            os.Exit(1)
        }
        log.Printf("[ToolParserLLM] ultra repair ready (SGLang at %s)", ToolParserLLMSGLangBaseURL())
    }

    handler := NewHandler()

    addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)

    tokenPadded := fmt.Sprintf("%-44s", config.Auth.Token)
    fmt.Printf(`
╔═══════════════════════════════════════════════════════════════╗
║           Z.AI Direct Bridge Server Started                   ║
╠═══════════════════════════════════════════════════════════════╣
║  Mode:          DIRECT HTTP (no browser needed)               ║
║  Captcha IPC:   IN-MEMORY (no FIFO / named pipe)             ║
║  Health:        http://localhost:%d/health               ║
╠═══════════════════════════════════════════════════════════════╣
║  OpenAI API:    http://localhost:%d/v1/chat/completions
║  Anthropic API: http://localhost:%d/v1/messages  ║
╠═══════════════════════════════════════════════════════════════╣
║  Auth Token:    %s║
╚═══════════════════════════════════════════════════════════════╝
`, config.Server.Port, config.Server.Port, config.Server.Port, tokenPadded)

    go func() {
        if err := initializeSession(); err != nil {
            log.Println("[Startup] Session init deferred — will retry on first request.")
        }
        // Warm up model cache
        fetchModelsFromZAI()
    }()

    // ── Session lifecycle ─────────────────────────────────────────────────
    // Every stateless request runs on a REUSED chat session (up to
    // SESSION_REUSE_COUNT requests per chat) instead of one chat per
    // message: one-chat-per-message is flagged by the Aliyun WAF as bot
    // behaviour and gets the egress IP blocked. Reuse is safe — verified
    // e2e that isolated single-turn payloads build no server-side history
    // on a shared chat_id — and each session is still deleted on Z.AI once
    // it hits its reuse limit, so the account never accumulates dead chats.
    // Async keeps a standing batch (SESSION_POOL_SIZE); --sync-mode reuses
    // one sticky session instead of the pool.
    if config.SyncMode {
        log.Printf("[Startup] Session mode: SYNC (--sync-mode: one sticky chat reused up to %dx, deleted on Z.AI + rotated after)", config.SessionReuseCount)
    } else {
        poolWait = time.Duration(config.SessionAcquireTimeout) * time.Second
        if config.SessionAcquireTimeout <= 0 {
            poolWait = 0 // 0 => wait indefinitely for a pooled session
        }
        sessionPool = NewSessionPool(NewZAIChatBackend(), config.SessionPoolSize)
        log.Printf("[Startup] Session mode: ASYNC (pre-made chat batch x%d, each reused up to %dx then deleted on Z.AI + refilled)", sessionPool.Size(), sessionPool.getMaxUses())
        log.Printf("[Startup]               SESSION_POOL_SIZE=%d SESSION_ACQUIRE_TIMEOUT=%ds SESSION_REUSE_COUNT=%d", sessionPool.Size(), config.SessionAcquireTimeout, sessionPool.getMaxUses())
    }

    srv := &http.Server{
        Addr:    addr,
        Handler: handler,
    }

    // Start serving before blocking on signals.
    serveErr := make(chan error, 1)
    go func() {
        serveErr <- srv.ListenAndServe()
    }()

    // Warm the standing session batch in the background; requests are served
    // meanwhile (they simply queue on Acquire until sessions appear).
    if sessionPool != nil {
        sessionPool.Start()
    }

    // ── Graceful shutdown (ported from the DeepseekFreeAPI reference) ─────
    // CTRL+C / SIGTERM stops accepting new connections, lets in-flight
    // responses finish (10s drain deadline), then deletes every still-pooled
    // chat session on Z.AI so nothing is left behind on the account. A
    // second CTRL+C force-exits immediately (default handling is re-armed).
    ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stopSignal()

    select {
    case err := <-serveErr:
        if err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Fatal(err)
        }
    case <-ctx.Done():
        stopSignal()
        log.Println("[Shutdown] Graceful shutdown requested — draining connections and clearing all chat sessions...")

        drainCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        if err := srv.Shutdown(drainCtx); err != nil {
            log.Printf("[Shutdown] drain deadline hit (%v); closing remaining connections", err)
            _ = srv.Close()
        }
        cancel()

        // Clear any sessions still pooled (or the sticky sync session) so
        // nothing is left behind on the Z.AI account (checked-out ones are
        // deleted by their own Release).
        if sessionPool != nil {
            sessionPool.Shutdown()
        } else {
            ShutdownSyncSession()
        }
        log.Println("[Shutdown] All chat sessions cleared. Goodbye.")
    }
}
