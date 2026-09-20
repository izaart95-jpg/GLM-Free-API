// session_pool.go
//
// Throwaway chat sessions + async session pool, ported from the
// DeepseekFreeAPI reference implementation (see reference/internal/dsproxy:
// pool.go for the pool, proxy.go for the garbage collector) and adapted to
// the Z.AI / chat.z.ai platform.
//
// ── Why this exists (context rot) ─────────────────────────────────────────
//
// OpenAI-compatible clients are stateless: they re-send the ENTIRE
// conversation (user + assistant turns) on every request. The bridge
// forwards that history to Z.AI inside a chat identified by `chat_id`, and
// every chat referenced by a completion materializes server-side under the
// bridge account (ZAI_TOKEN, or the guest identity) together with its full
// history. Two failure modes follow:
//
//   1. Accumulation — nothing ever deletes those chats, so the account fills
//      up with dead sessions, one per proxied request.
//   2. Context rot — if a server-side chat outlives a single request, Z.AI's
//      own stored history stacks on top of the history the client already
//      re-sent, the model sees duplicated/stale context, and the
//      conversation quality rots.
//
// The fix mirrors the reference project:
//
//   - Every stateless request runs on a THROWAWAY chat session. As soon as
//     the response has been fully written (or definitively failed) the chat
//     is deleted on Z.AI (DELETE /api/v1/chats/{chat_id}) so no session
//     history survives the request and nothing accumulates on the account.
//   - Async mode (default): a standing batch of ready sessions
//     (SESSION_POOL_SIZE, default 5) is kept pre-made at all times; requests
//     grab one instantly, and each consumed session is deleted upstream +
//     replaced the moment its response is fully processed, so the batch
//     refills itself for as long as the app runs.
//   - Sync mode (--sync-mode / SYNC_MODE=true): the legacy flow where every
//     request creates its own session first; the used session is still
//     garbage-collected afterwards.
//   - Graceful shutdown (CTRL+C / SIGTERM): in-flight requests are drained,
//     then every still-pooled session is deleted on Z.AI before exiting
//     (a second CTRL+C force-exits).
//
// ── Platform difference vs the reference ──────────────────────────────────
//
// DeepSeek creates sessions with an upstream API call (POST
// /chat_session/create) and deletes them in bulk (POST /chat_session/delete),
// so the reference pool pre-warms against the network. Z.AI chat IDs are
// client-generated UUIDs: a chat only materializes server-side when a
// completion first references it, and deletion is one DELETE per chat:
//
//     DELETE https://chat.z.ai/api/v1/chats/<chat_id>
//     authorization: Bearer <token>
//       -> true                                            (deleted)
//       -> {"detail":"We could not find what you're looking for :/"}
//                                                          (already gone)
//
// Therefore pool warmup here is instant and local (just mint UUIDs), and a
// session that is never consumed never touches the account at all. The
// retirement discipline — delete after use, then refill — is unchanged, and
// the "already gone" reply counts as a successful delete (idempotent GC).

package zbridge

import (
    "context"
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
    "sync"
    "sync/atomic"
    "time"
)

var (
    // ErrPoolClosing is returned by Acquire once Shutdown has begun.
    ErrPoolClosing = errors.New("session pool is shutting down")
    // ErrPoolTimeout is returned by Acquire when no pooled session became
    // available within the configured wait window.
    ErrPoolTimeout = errors.New("timed out waiting for a pooled session")
)

const (
    // defaultPoolSize is the standing batch of pre-made ready sessions.
    defaultPoolSize = 5
    // defaultSessionReuse is how many requests one chat session serves
    // before it is deleted on Z.AI and replaced (SESSION_REUSE_COUNT).
    // 10 looks human (a short chat run), cuts DELETE churn 10-fold so the
    // Aliyun WAF stops flagging one-chat-per-message as bot behaviour,
    // and still rotates often enough that any future server-side retention
    // can never accumulate into context rot. 1 = legacy throwaway flow.
    defaultSessionReuse = 10
    // defaultPoolWait bounds how long a completion request waits for a
    // pooled session before creating one directly (SESSION_ACQUIRE_TIMEOUT).
    defaultPoolWait = 10 * time.Second
    // poolOpTimeout bounds one upstream delete call.
    poolOpTimeout = 30 * time.Second
    // poolCreateBackoffStart / poolCreateBackoffMax shape the retry delay
    // when session creation fails (kept for parity with the reference; on
    // Z.AI creation is local and does not fail unless the backend is
    // swapped for one that calls upstream).
    poolCreateBackoffStart = 1 * time.Second
    poolCreateBackoffMax   = 15 * time.Second
    // poolDrainWait bounds how long Shutdown waits for in-flight
    // retire/refill operations before reporting leftovers.
    poolDrainWait = 20 * time.Second
)

// SessionBackend is the slice of the Z.AI bridge the pool needs. Tests
// substitute a stub; the production backend is zaiSessionBackend.
type SessionBackend interface {
    CreateChatSession(ctx context.Context) (string, error)
    DeleteChatSession(ctx context.Context, sessionIDs ...string) error
}

// zaiSessionBackend implements SessionBackend against chat.z.ai.
type zaiSessionBackend struct{}

// NewZAIChatBackend returns the production SessionBackend: sessions are
// client-generated chat IDs, deleted via DELETE /api/v1/chats/{id}. Run
// attaches a pool built from it at startup; tests use it to exercise the
// real delete path against a mock upstream.
func NewZAIChatBackend() SessionBackend { return zaiSessionBackend{} }

// CreateChatSession mints one fresh chat ID. Z.AI chats are identified by
// client-generated UUIDs and only materialize server-side when a completion
// first references one, so there is nothing to do upstream here — which also
// means warmup never leaves unconsumed sessions on the account.
func (zaiSessionBackend) CreateChatSession(ctx context.Context) (string, error) {
    if err := ctx.Err(); err != nil {
        return "", err
    }
    return randomUUID(), nil
}

// DeleteChatSession deletes chats one by one (Z.AI has no bulk endpoint).
// Best-effort: every ID is attempted; the first error is returned after the
// rest have been tried.
func (zaiSessionBackend) DeleteChatSession(ctx context.Context, sessionIDs ...string) error {
    var firstErr error
    for _, id := range sessionIDs {
        if id == "" {
            continue
        }
        if err := DeleteZAIChat(ctx, id); err != nil && firstErr == nil {
            firstErr = err
        }
    }
    return firstErr
}

// DeleteZAIChat deletes one chat from the Z.AI account so its history does
// not accumulate on the account:
//
//     DELETE {BASE_URL}/api/v1/chats/{chatID}   (authorization: Bearer <token>)
//       200 "true"                                            -> deleted
//       404 / {"detail":"We could not find what you're ..."}  -> already gone
//       401                                                   -> re-init + retry once
//
// Deletion is idempotent: an "already gone" reply is treated as success so
// the GC never gets stuck on a chat that was collected twice.
//
// Issue #41: while the WAF breaker is open (IP blocked), the DELETE would be
// answered with the block page too — a wasted upstream call per retired
// session that also keeps the block alive. Deletion is skipped instead: the
// chats are throwaway UUIDs that only materialize server-side when a
// completion references them, so a chat whose delete was skipped during a
// block is at worst a short-lived leftover that a later GC pass retires
// through the idempotent 404 "already gone" path.
func DeleteZAIChat(ctx context.Context, chatID string) error {
    if chatID == "" {
        return nil
    }

    // Issue #41: do not poke the blocked edge while the breaker is open.
    if err, _ := CheckWAF(); err != nil {
        return fmt.Errorf("chat delete skipped: %s", err.Error())
    }

    for attempt := 0; attempt < 2; attempt++ {
        session.mu.Lock()
        token := session.Token
        initialized := session.Initialized
        feVersion := session.FeVersion
        session.mu.Unlock()

        if token == "" || !initialized {
            // Own context: the triggering request may already be gone.
            if err := initializeSession(); err != nil {
                return fmt.Errorf("session init for chat delete: %s", err.Error())
            }
            continue // re-read the fresh token
        }

        urlStr := BASE_URL + "/api/v1/chats/" + chatID
        req, err := http.NewRequestWithContext(ctx, "DELETE", urlStr, nil)
        if err != nil {
            return fmt.Errorf("chat delete request build failed: %s", err.Error())
        }
        req.Header.Set("authorization", "Bearer "+token)
        req.Header.Set("User-Agent", zaiUserAgent)
        req.Header.Set("content-type", "application/json")
        req.Header.Set("x-fe-Version", feVersion)

        if config.Logging.Level == "debug" {
            log.Printf("[DEBUG] Z.AI chat delete: DELETE %s", urlStr)
        }

        resp, err := zaiHTTPClient.Do(req)
        if err != nil {
            return fmt.Errorf("Z.AI chat delete connection error: %s", err.Error())
        }
        body, _ := io.ReadAll(resp.Body)
        resp.Body.Close()
        text := strings.TrimSpace(string(body))

        if config.Logging.Level == "debug" {
            log.Printf("[DEBUG] Z.AI chat delete response: %d %s", resp.StatusCode, text)
        }

        switch {
        case resp.StatusCode == 401:
            // Token expired mid-flight: force re-init and retry once.
            session.mu.Lock()
            session.Initialized = false
            session.mu.Unlock()
            continue
        case resp.StatusCode == 404 || strings.Contains(text, "could not find"):
            return nil // already gone — nothing left to collect
        case resp.StatusCode >= 200 && resp.StatusCode < 300:
            return nil // "true" (or any 2xx) — deleted
        default:
            return fmt.Errorf("Z.AI chat delete failed: %d: %s", resp.StatusCode, text)
        }
    }
    return errors.New("chat delete: max retries exceeded")
}

// SessionPool holds the standing batch of ready stateless chat sessions.
//
// With the pool attached, every completion draws a pre-made session instead
// of minting one per request. Each session is REUSED up to maxUses times
// (SESSION_REUSE_COUNT, default 10) to look human — one chat carrying
// several turns instead of one chat per message, which is what the Aliyun
// WAF flags as bot behaviour — and only then is it deleted upstream and
// replaced, so DELETE churn drops N-fold while the batch stays at full
// strength. Verified e2e: reusing one chat_id for isolated single-turn
// payloads builds no server-side history, so reuse does not thread
// conversations.
type SessionPool struct {
    backend SessionBackend
    size    int

    ready chan string // buffered with size; members are reusable sessions

    stopOnce sync.Once
    stopCh   chan struct{}
    stopped  atomic.Bool

    wg sync.WaitGroup // outstanding create/delete operations

    // uses tracks completed requests per pooled session ID. Guarded by mu.
    mu      sync.Mutex
    uses    map[string]int
    maxUses int // <=0 means "read live from config.SessionReuseCount"
}

// NewSessionPool builds a pool that keeps size sessions ready. size < 1 is
// clamped to the default batch. Call Start to begin warmup.
func NewSessionPool(backend SessionBackend, size int) *SessionPool {
    if size < 1 {
        size = defaultPoolSize
    }
    maxUses := defaultSessionReuse
    // config is package-initialised before any pool is built in production;
    // honour it so SESSION_REUSE_COUNT / --session-reuse-count applies.
    // (Tests may override per-pool via SetMaxUses.)
    if config.SessionReuseCount >= 1 {
        maxUses = config.SessionReuseCount
    }
    return &SessionPool{
        backend: backend,
        size:    size,
        ready:   make(chan string, size),
        stopCh:  make(chan struct{}),
        uses:    make(map[string]int),
        maxUses: maxUses,
    }
}

// SetMaxUses overrides the reuse limit for this pool (tests + ops tooling).
// n < 1 restores "read live from config".
func (p *SessionPool) SetMaxUses(n int) {
    p.mu.Lock()
    p.maxUses = n
    p.mu.Unlock()
}

// getMaxUses resolves the effective reuse limit: per-pool override wins,
// otherwise the live global config, otherwise the built-in default.
func (p *SessionPool) getMaxUses() int {
    p.mu.Lock()
    override := p.maxUses
    p.mu.Unlock()
    if override >= 1 {
        return override
    }
    if config.SessionReuseCount >= 1 {
        return config.SessionReuseCount
    }
    return defaultSessionReuse
}

// Uses returns completed-request count for one pooled session (observability).
func (p *SessionPool) Uses(id string) int {
    p.mu.Lock()
    defer p.mu.Unlock()
    return p.uses[id]
}

// Size reports the configured batch size.
func (p *SessionPool) Size() int { return p.size }

// Ready reports how many sessions are currently stocked.
func (p *SessionPool) Ready() int { return len(p.ready) }

// Start launches the warmup goroutines that pre-make the initial batch.
func (p *SessionPool) Start() {
    log.Printf("[Pool] warming up %d stateless session(s)...", p.size)
    for i := 0; i < p.size; i++ {
        p.wg.Add(1)
        go func() {
            defer p.wg.Done()
            p.fillSlot("warmup")
        }()
    }
}

// Acquire hands out one ready session, blocking until one is available, ctx
// is done, or wait elapses (wait <= 0 waits indefinitely). The caller MUST
// eventually call Release with the returned ID — even on error paths — so
// the used session is retired and the batch refilled.
func (p *SessionPool) Acquire(ctx context.Context, wait time.Duration) (string, error) {
    var timeout <-chan time.Time
    if wait > 0 {
        timer := time.NewTimer(wait)
        defer timer.Stop()
        timeout = timer.C
    }
    for {
        select {
        case id := <-p.ready:
            if config.Logging.Level == "debug" {
                log.Printf("[Pool] handed out session %s (%d/%d ready)", id, len(p.ready), p.size)
            }
            return id, nil
        case <-ctx.Done():
            return "", ctx.Err()
        case <-timeout:
            return "", ErrPoolTimeout
        case <-p.stopCh:
            return "", ErrPoolClosing
        }
    }
}

// Release retires a consumed session. It is called only after the response
// has been fully written and processed (or definitively failed), so the
// session is never yanked out from under an in-flight completion. The
// session is reused until it has served maxUses requests (SESSION_REUSE_COUNT):
// below the limit it is stocked straight back into the batch for the next
// request; once the limit is reached it is deleted upstream and a
// replacement is created to fill the gap. Both steps run in the background.
func (p *SessionPool) Release(sessionID string) {
    if sessionID == "" {
        return
    }
    p.wg.Add(1)
    go func() {
        defer p.wg.Done()
        max := p.getMaxUses()
        if max < 1 {
            max = 1
        }
        p.mu.Lock()
        p.uses[sessionID]++
        uses := p.uses[sessionID]
        p.mu.Unlock()

        if uses < max {
            if p.stopped.Load() {
                // Shutting down: retire only, don't rebuild the batch.
                p.deleteOne(sessionID, "shutdown")
                p.mu.Lock()
                delete(p.uses, sessionID)
                p.mu.Unlock()
                return
            }
            select {
            case p.ready <- sessionID:
                log.Printf("[Pool:reuse] session %s reused %d/%d (%d/%d ready)", sessionID, uses, max, len(p.ready), p.size)
            case <-p.stopCh:
                p.deleteOne(sessionID, "shutdown-race")
                p.mu.Lock()
                delete(p.uses, sessionID)
                p.mu.Unlock()
            }
            return
        }
        p.deleteOne(sessionID, "used")
        p.mu.Lock()
        delete(p.uses, sessionID)
        p.mu.Unlock()
        if p.stopped.Load() {
            return // shutting down: retire only, don't rebuild the batch
        }
        p.fillSlot("refill")
    }()
}

// Shutdown gracefully clears the pool (the CTRL+C path): refills are
// stopped, every still-pooled session is deleted upstream so nothing is left
// behind on the account, and we briefly wait for in-flight retire/refill
// operations to finish. Sessions already checked out are deleted by their
// request's Release call.
func (p *SessionPool) Shutdown() {
    first := false
    p.stopOnce.Do(func() {
        first = true
        p.stopped.Store(true)
        close(p.stopCh)
    })
    if !first {
        return
    }

    // Collect whatever is still stocked and delete it. Z.AI deletion is one
    // call per chat, so walk the leftovers sequentially.
    var leftover []string
    for {
        select {
        case id := <-p.ready:
            leftover = append(leftover, id)
            continue
        default:
        }
        break
    }
    if len(leftover) > 0 {
        log.Printf("[Pool] clearing all remaining sessions (%d)...", len(leftover))
        ctx, cancel := context.WithTimeout(context.Background(), poolOpTimeout)
        err := p.backend.DeleteChatSession(ctx, leftover...)
        cancel()
        if err != nil {
            log.Printf("[Pool] warning: failed to clear %d session(s): %v", len(leftover), err)
        } else {
            log.Printf("[Pool] cleared %d pooled session(s): deleted %v", len(leftover), leftover)
        }
        p.mu.Lock()
        for _, id := range leftover {
            delete(p.uses, id)
        }
        p.mu.Unlock()
    } else {
        log.Printf("[Pool] clearing all sessions... none remaining")
    }

    // Wait (bounded) for outstanding creates/deletes to wind down.
    done := make(chan struct{})
    go func() {
        p.wg.Wait()
        close(done)
    }()
    select {
    case <-done:
        log.Printf("[Pool] all sessions accounted for")
    case <-time.After(poolDrainWait):
        log.Printf("[Pool] warning: some background session operations did not finish within %s", poolDrainWait)
    }
}

// fillSlot creates one session (retrying through transient failures) and
// stocks it, unless shutdown won the race. Runs synchronously; callers wrap
// it in a goroutine where concurrency is wanted.
func (p *SessionPool) fillSlot(reason string) {
    backoff := poolCreateBackoffStart
    loggedOnce := false
    for {
        if p.stopped.Load() {
            return
        }
        ctx, cancel := context.WithTimeout(context.Background(), poolOpTimeout)
        id, err := p.backend.CreateChatSession(ctx)
        cancel()
        if err != nil {
            if p.stopped.Load() {
                return
            }
            // Log the first failure loudly; repeats stay quiet so a
            // misconfigured backend doesn't spam the log.
            if !loggedOnce {
                log.Printf("[Pool:%s] session creation failed (%v); retrying...", reason, err)
                loggedOnce = true
            } else if config.Logging.Level == "debug" {
                log.Printf("[Pool:%s] session creation failed again (%v); retrying in %s", reason, err, backoff)
            }
            select {
            case <-time.After(backoff):
            case <-p.stopCh:
                return
            }
            backoff *= 2
            if backoff > poolCreateBackoffMax {
                backoff = poolCreateBackoffMax
            }
            continue
        }
        p.stock(id, reason)
        return
    }
}

// stock puts a freshly created session into the batch, or deletes it if
// shutdown raced in first (never stockpile sessions nobody will consume).
func (p *SessionPool) stock(id, reason string) {
    select {
    case p.ready <- id:
        log.Printf("[Pool:%s] session ready: %s (%d/%d)", reason, id, len(p.ready), p.size)
    case <-p.stopCh:
        p.deleteOne(id, "shutdown-race")
    }
}

// deleteOne deletes a single session upstream, best-effort.
func (p *SessionPool) deleteOne(id, reason string) {
    ctx, cancel := context.WithTimeout(context.Background(), poolOpTimeout)
    defer cancel()
    if err := p.backend.DeleteChatSession(ctx, id); err != nil {
        log.Printf("[Pool:%s] warning: failed to delete session %s: %v", reason, id, err)
        return
    }
    if config.Logging.Level == "debug" {
        log.Printf("[Pool:%s] deleted chat session: %s", reason, id)
    }
}

// ── Bridge glue ─────────────────────────────────────────────────────────────

var (
    // sessionPool holds the standing batch of reusable chat sessions.
    // nil when running in sync mode (--sync-mode / SYNC_MODE=true); the
    // sync flow reuses one sticky session up to SESSION_REUSE_COUNT too.
    sessionPool *SessionPool
    // poolWait bounds how long a request waits for a pooled session before
    // creating one directly (0 waits forever). See SESSION_ACQUIRE_TIMEOUT.
    poolWait = defaultPoolWait

    // syncSticky is the single reused chat ID for sync mode (--sync-mode):
    // one chat carries up to SESSION_REUSE_COUNT requests (human-like)
    // instead of one chat per message (bot-like), then rotates. Guarded by
    // syncMu. Verified e2e: reuse builds no server-side history.
    syncMu   sync.Mutex
    syncID   string
    syncUses int
)

// syncMaxUses resolves the effective sync reuse limit.
func syncMaxUses() int {
    if config.SessionReuseCount >= 1 {
        return config.SessionReuseCount
    }
    return defaultSessionReuse
}

// ResetSyncSticky clears the sync-mode sticky session (tests + pool attach).
// The caller owns deleting the old ID if it still needs GC; here we just
// forget it so the next Acquire mints fresh. Production rotation deletes
// via gcSessions before clearing.
func ResetSyncSticky() {
    syncMu.Lock()
    syncID = ""
    syncUses = 0
    syncMu.Unlock()
}

// ShutdownSyncSession deletes the sticky sync session, if any (shutdown path
// so a reused chat is never left behind on the account).
func ShutdownSyncSession() {
    syncMu.Lock()
    id := syncID
    syncID = ""
    syncUses = 0
    syncMu.Unlock()
    if id == "" {
        return
    }
    ctx, cancel := context.WithTimeout(context.Background(), poolOpTimeout)
    defer cancel()
    if err := (zaiSessionBackend{}).DeleteChatSession(ctx, id); err != nil {
        log.Printf("[Sync] warning: failed to delete sticky session %s: %v", id, err)
        return
    }
    log.Printf("[Sync] cleared sticky session: %s", id)
}

// AttachSessionPool swaps the async session pool (and its acquire wait
// window) used by stateless requests and returns a function that restores
// the previous attachment. Passing nil detaches the pool, i.e. switches to
// the sync (sticky-reuse) flow. Run uses it once at startup; the
// blackbox tests in tests/ use it to exercise both flows.
func AttachSessionPool(p *SessionPool, wait time.Duration) func() {
    oldPool, oldWait := sessionPool, poolWait
    sessionPool, poolWait = p, wait
    // Avoid cross-mode pollution in tests: a stale sticky ID must not leak
    // into a pooled run or vice versa.
    ResetSyncSticky()
    return func() {
        sessionPool, poolWait = oldPool, oldWait
        ResetSyncSticky()
    }
}

// AcquireStatelessSession returns a chat ID for one stateless request.
// Async mode takes a pre-made session from the standing batch (each pooled
// session serves up to SESSION_REUSE_COUNT requests before rotating); if a
// burst exhausts the batch the request waits up to poolWait and then creates
// a session directly instead of stalling indefinitely. Sync mode reuses one
// sticky session up to SESSION_REUSE_COUNT requests before rotating.
//
// The second return value reports whether the session is pool-owned (retired
// through pool.Release) or not (retired through ReleaseStatelessSession's
// sync-reuse / on-demand path).
func AcquireStatelessSession(ctx context.Context) (chatID string, pooled bool, err error) {
    if sessionPool == nil {
        syncMu.Lock()
        if syncID == "" {
            syncID = randomUUID()
            syncUses = 0
            log.Printf("[Sync] sticky session: %s (new)", syncID)
        } else {
            log.Printf("[Sync] sticky session: %s (reuse %d/%d)", syncID, syncUses, syncMaxUses())
        }
        id := syncID
        syncMu.Unlock()
        return id, false, nil
    }
    id, acqErr := sessionPool.Acquire(ctx, poolWait)
    switch {
    case acqErr == nil:
        log.Printf("[Pool] stateless session: %s (%d/%d ready, uses %d/%d)", id, sessionPool.Ready(), sessionPool.Size(), sessionPool.Uses(id), sessionPool.getMaxUses())
        return id, true, nil
    case errors.Is(acqErr, ErrPoolTimeout):
        chatID = randomUUID()
        log.Printf("[Pool] busy — stateless session created on demand: %s", chatID)
        return chatID, false, nil
    default: // ErrPoolClosing, or the request's client went away
        if errors.Is(acqErr, ErrPoolClosing) {
            return "", false, errors.New("server is shutting down")
        }
        return "", false, acqErr
    }
}

// ReleaseStatelessSession retires a used stateless chat session. It must be
// called only after the response has been fully written and processed (or
// definitively failed). Pooled sessions go back through pool.Release (reuse
// until SESSION_REUSE_COUNT, then delete + refill). Sync-mode sticky sessions
// are kept until they hit SESSION_REUSE_COUNT, then deleted and rotated.
// On-demand overflow sessions (pool busy) are still deleted immediately via GC.
func ReleaseStatelessSession(chatID string, pooled bool) {
    if chatID == "" {
        return
    }
    if pooled && sessionPool != nil {
        sessionPool.Release(chatID)
        return
    }
    if sessionPool != nil {
        // On-demand overflow (pool busy): single-use, delete immediately.
        gcSessions("stateless", chatID)
        return
    }
    // Sync mode: sticky reuse.
    syncMu.Lock()
    if chatID != syncID {
        // Stale ID (rotated concurrently): delete directly, don't touch counters.
        stale := chatID
        syncMu.Unlock()
        gcSessions("sync-stale", stale)
        return
    }
    syncUses++
    uses, max := syncUses, syncMaxUses()
    if uses < max {
        syncMu.Unlock()
        log.Printf("[Sync] session %s reused %d/%d", chatID, uses, max)
        return
    }
    // Limit reached: rotate. Clear first so the next Acquire mints fresh
    // even while the DELETE below is still in flight.
    syncID = ""
    syncUses = 0
    syncMu.Unlock()
    log.Printf("[Sync] session %s reached max uses %d, rotating", chatID, max)
    gcSessions("sync-reuse-max", chatID)
}

// gcSessions asynchronously deletes used-up chat sessions on Z.AI so their
// history doesn't accumulate on the account (port of the reference's GC):
// sessions are temporary by design and are removed right after their last
// use, in a background goroutine so response latency is unaffected. Deletion
// is best-effort — failures are logged and otherwise ignored.
func gcSessions(reason string, sessionIDs ...string) {
    ids := make([]string, 0, len(sessionIDs))
    for _, id := range sessionIDs {
        if id != "" {
            ids = append(ids, id)
        }
    }
    if len(ids) == 0 {
        return
    }
    go func() {
        // Own context: the triggering request may already be gone by now.
        ctx, cancel := context.WithTimeout(context.Background(), poolOpTimeout)
        defer cancel()
        backend := zaiSessionBackend{}
        if err := backend.DeleteChatSession(ctx, ids...); err != nil {
            log.Printf("[GC:%s] failed to delete chat session(s) %v: %v", reason, ids, err)
            return
        }
        log.Printf("[GC:%s] deleted chat session(s): %v", reason, ids)
    }()
}

// sessionPoolStatus reports the session-lifecycle mode for /status.
func sessionPoolStatus() map[string]interface{} {
    max := config.SessionReuseCount
    if max < 1 {
        max = defaultSessionReuse
    }
    if sessionPool == nil {
        syncMu.Lock()
        uses, cur := syncUses, syncID
        syncMu.Unlock()
        st := map[string]interface{}{
            "mode":       "sync",
            "throwaway":  false,
            "reuse":      true,
            "maxUses":    max,
            "gc_enabled": true,
        }
        if cur != "" {
            st["currentUses"] = uses
        }
        return st
    }
    return map[string]interface{}{
        "mode":       "async",
        "throwaway":  false,
        "reuse":      true,
        "maxUses":    sessionPool.getMaxUses(),
        "gc_enabled": true,
        "size":       sessionPool.Size(),
        "ready":      sessionPool.Ready(),
    }
}
