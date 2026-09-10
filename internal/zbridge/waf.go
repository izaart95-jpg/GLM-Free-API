package zbridge

import (
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
    "sync"
    "time"
)

// WAF circuit breaker — issue #41.
//
// Sustained agent traffic trips the Aliyun WAF in front of chat.z.ai: the
// egress IP (not the bridge, not the account) gets served the block page on
// POST /api/v2/chat/completions for a while. The block is IP-wide — a real
// browser on the same IP fails identically — and outlives any single request,
// so retrying immediately is guaranteed to fail while making the block worse.
//
// Once the block page is detected the breaker opens: requests fail fast with
// a clean, actionable error instead of burning a captcha (each one costs a
// harvested device token and ~10-15 s of signing work), chat deletions pause
// (so a block no longer generates pointless upstream traffic that keeps it
// alive), and a dedicated background prober re-checks the blocked endpoint
// after each cooldown — doubling the backoff while the block persists — and
// resumes traffic automatically the moment it lifts.

// ErrWAFBlock is returned (and surfaced to clients) while the breaker is open.
// The message is deliberately stable, short, and JSON-safe: agents must be
// able to parse it and decide to wait, not wade through a 3 KB HTML page.
var ErrWAFBlock = errors.New(
    "chat.z.ai has temporarily blocked this server's IP (Aliyun WAF). " +
        "All requests will fail until the block expires; retry after the suggested delay")

// Breaker states. There is no half-open state on purpose: probing is owned by
// a dedicated background goroutine (see wafProber), so request handlers only
// ever see closed (proceed) or open (fail fast).
const (
    breakerClosed = 0 // normal traffic
    breakerOpen   = 1 // blocked: fail fast; background prober will re-check
)

// Breaker timing knobs. The cooldown starts at wafBaseCooldown and doubles on
// every consecutive probe that still sees the block page, capped at
// wafMaxCooldown. Probing while blocked costs one cheap request (a bare,
// unauthenticated POST to the blocked endpoint — never a captcha, never a
// session), so a long block does not bleed device tokens.
const (
    wafBaseCooldown = 60 * time.Second
    wafMaxCooldown  = 30 * time.Minute
    // wafProbeEndpoint is the endpoint the Aliyun WAF actually blocks (issue
    // #41: "GET / returns 200, GET /api/models returns 200 — POST
    // /api/v2/chat/completions returns 405 WAF block"). A block is per-path,
    // so probing any other path reports "cleared" while completions still
    // fail. The probe therefore sends a BARE POST — empty body, no
    // authorization, no captcha — which costs nothing upstream and splits
    // the two outcomes cleanly:
    //
    //   still blocked -> 405 + the HTML block page (isWAFBlockResponse: true)
    //   block lifted  -> 403/4xx + JSON {"detail":"Not authenticated"}
    wafProbeEndpoint = "/api/v2/chat/completions"
    // wafBlockMarker is the unique substring of the Aliyun block page
    // ("Sorry, your request has been blocked as it may cause potential
    // threats to the server's security." render template). The block page
    // ships the same marker in every language.
    wafBlockMarker = "blocked as it may cause potential threats"
)

// wafBreaker is the process-wide breaker. Every upstream Z.AI call funnels
// through it.
var wafBreaker = &wafBreakerState{}

type wafBreakerState struct {
    mu         sync.Mutex
    state      int           // breakerClosed | breakerOpen
    openedAt   time.Time     // when the current open period began
    cooldown   time.Duration // current probe backoff (doubles while blocked)
    blockCount int           // consecutive block-page observations

    // probeGen invalidates outstanding prober goroutines: every transition
    // that changes the plan (re-open, close, manual reset) bumps it, and a
    // prober that wakes up holding a stale generation exits without probing.
    probeGen uint64
}

// wafProbeTransport issues the actual probe request: a bare POST to the
// blocked endpoint with no auth and an empty body. Package var so tests can
// swap it for a stub without standing up a real server.
var wafProbeTransport = func() (*http.Response, error) {
    req, err := http.NewRequest(http.MethodPost, BASE_URL+wafProbeEndpoint, strings.NewReader("{}"))
    if err != nil {
        return nil, err
    }
    req.Header.Set("User-Agent", zaiUserAgent)
    req.Header.Set("Accept", "application/json, text/plain, */*")
    req.Header.Set("Accept-Language", "en-US,en;q=0.9")
    req.Header.Set("Content-Type", "application/json")
    return zaiHTTPClient.Do(req)
}

// isWAFBlockResponse reports whether an upstream response is the Aliyun WAF
// block page. The page is served with status 405 (occasionally 403) and its
// body carries the marker string; the status alone is not enough because Z.AI
// itself returns 405 for genuinely unsupported methods on other paths.
func isWAFBlockResponse(statusCode int, body []byte) bool {
    if statusCode != http.StatusMethodNotAllowed && statusCode != http.StatusForbidden {
        return false
    }
    if len(body) == 0 {
        return false
    }
    // Cheap pre-filter: the block page is HTML; a JSON error is not.
    s := string(body)
    if !strings.Contains(s, "<") {
        return false
    }
    return strings.Contains(s, wafBlockMarker) ||
        strings.Contains(s, "errors.aliyun.com") ||
        strings.Contains(s, "您的访问被阻断")
}

// RecordWAFBlock marks the IP as blocked and opens the breaker. Callers
// should stop issuing further upstream calls for this request and surface
// ErrWAFBlock instead. Opening the breaker also launches a background prober
// goroutine that re-checks the block after each cooldown (doubling on
// persistence) until the block clears — callers never own the probe.
func RecordWAFBlock() {
    wafBreaker.mu.Lock()
    wafBreaker.blockCount++
    if wafBreaker.state == breakerClosed {
        wafBreaker.state = breakerOpen
        wafBreaker.openedAt = time.Now()
        wafBreaker.cooldown = wafBaseCooldown
        wafBreaker.probeGen++
        gen := wafBreaker.probeGen
        wafBreaker.mu.Unlock()
        logWAFBlock("WAF block detected — opening circuit breaker (fail fast; probing with backoff)")
        go wafProber(gen)
        return
    }
    wafBreaker.mu.Unlock()
}

// wafProber waits out the current cooldown, then probes once. If the block
// persists it re-opens with a doubled cooldown (and schedules itself again);
// if the block has cleared it closes the breaker. A stale generation (the
// breaker was closed or re-opened by someone else meanwhile) exits without
// probing.
//
// The wait is tick-sliced (wafProberSlice at a time, re-reading the cooldown
// each slice) so a cooldown adjusted mid-wait — tests shrinking it, ops
// tooling lengthening it — applies promptly instead of at the next full
// cycle.
func wafProber(gen uint64) {
    for {
        if !wafSleepCooldown(gen) {
            return // breaker closed or superseded while waiting
        }
        if !wafClaimProbe(gen) {
            return // stale prober: a newer state transition owns the breaker
        }
        blocked := wafProbeOnce()
        wafBreaker.mu.Lock()
        if wafBreaker.probeGen != gen || wafBreaker.state == breakerClosed {
            // Closed or superseded while the probe was in flight; do nothing.
            wafBreaker.mu.Unlock()
            return
        }
        if blocked {
            wafBreaker.state = breakerOpen
            wafBreaker.openedAt = time.Now()
            wafBreaker.cooldown = minDuration(wafBreaker.cooldown*2, wafMaxCooldown)
            wafBreaker.mu.Unlock()
            logWAFBlock(fmt.Sprintf("WAF block persists — backing off %.0fs before next probe", wafBreaker.cooldown.Seconds()))
            continue
        }
        wafBreaker.state = breakerClosed
        wafBreaker.cooldown = 0
        wafBreaker.blockCount = 0
        wafBreaker.probeGen++
        wafBreaker.mu.Unlock()
        logWAFBlock("WAF block cleared — circuit breaker closed, traffic resumed")
        return
    }
}

// wafProberSlice is how often a waiting prober re-checks the breaker state.
// Cooldowns are minutes-long, so a one-second tick is negligible overhead
// while keeping a mid-wait cooldown change effective within a second. It is
// a var (not a const) purely so tests can shrink it.
var wafProberSlice = time.Second

// wafSleepCooldown waits out the CURRENT cooldown, sliced into
// wafProberSlice ticks with the deadline recomputed from live state every
// tick. It returns false when the breaker closed or was superseded while
// waiting (the prober must exit without probing).
func wafSleepCooldown(gen uint64) bool {
    for {
        wafBreaker.mu.Lock()
        if wafBreaker.probeGen != gen || wafBreaker.state == breakerClosed {
            wafBreaker.mu.Unlock()
            return false
        }
        remaining := wafBreaker.openedAt.Add(wafBreaker.cooldown).Sub(time.Now())
        wafBreaker.mu.Unlock()

        if remaining <= 0 {
            return true
        }
        if remaining > wafProberSlice {
            remaining = wafProberSlice
        }
        time.Sleep(remaining)
    }
}

// wafClaimProbe atomically verifies that gen is still current and the
// breaker is still open. It returns false when the generation is stale.
func wafClaimProbe(gen uint64) bool {
    wafBreaker.mu.Lock()
    defer wafBreaker.mu.Unlock()
    if wafBreaker.probeGen != gen || wafBreaker.state == breakerClosed {
        return false
    }
    return true
}

// wafProbeOnce issues the cheap probe request and reports whether the block
// is still in effect. A transport-level failure is treated as "still blocked"
// (conservative) but never panics.
func wafProbeOnce() bool {
    resp, err := wafProbeTransport()
    if err != nil {
        return true
    }
    defer resp.Body.Close()
    body := readProbeBody(resp)
    return isWAFBlockResponse(resp.StatusCode, body)
}

// RecordWAFSuccess notes that upstream is reachable again. Any non-block
// response (even a normal API error like 401/400) proves the IP is not
// blocked, so it closes the breaker and resets the backoff. Outstanding
// probers notice via the generation bump.
func RecordWAFSuccess() {
    wafBreaker.mu.Lock()
    if wafBreaker.state != breakerClosed || wafBreaker.blockCount != 0 {
        wafBreaker.state = breakerClosed
        wafBreaker.cooldown = 0
        wafBreaker.blockCount = 0
        wafBreaker.probeGen++
    }
    wafBreaker.mu.Unlock()
}

// CheckWAF returns ErrWAFBlock if the breaker believes the IP is currently
// blocked. Probing is owned by the background prober, so this call never
// blocks and never issues network traffic: callers either proceed or fail
// fast.
//
// The returned wait hint (>= 0) is the remaining cooldown the caller can
// pass to clients as a Retry-After suggestion.
func CheckWAF() (error, time.Duration) {
    wafBreaker.mu.Lock()
    defer wafBreaker.mu.Unlock()

    if wafBreaker.state == breakerClosed {
        return nil, 0
    }
    elapsed := time.Since(wafBreaker.openedAt)
    remaining := wafBreaker.cooldown - elapsed
    if remaining < 0 {
        remaining = 0
    }
    return fmt.Errorf("%w (retry in ~%s)", ErrWAFBlock, remaining.Round(time.Second)), remaining
}

// WAFStatus reports the breaker state for /status.
func WAFStatus() map[string]interface{} {
    wafBreaker.mu.Lock()
    defer wafBreaker.mu.Unlock()
    st := map[string]interface{}{
        "blocked": wafBreaker.state != breakerClosed,
        "state":   wafBreaker.state,
    }
    if wafBreaker.state != breakerClosed {
        st["retryIn"] = time.Until(wafBreaker.openedAt.Add(wafBreaker.cooldown)).Round(time.Second).String()
        st["consecutiveBlocks"] = wafBreaker.blockCount
    }
    return st
}

// RejectIfWAFBlocked is the handler-level gate (issue #41). While the breaker
// is open it rejects the request immediately — before a pooled session is
// acquired, before vision upload, and above all before any captcha is
// computed — with 503 + Retry-After so well-behaved clients back off on
// their own. It returns true when the request has been answered.
func RejectIfWAFBlocked(w http.ResponseWriter) bool {
    err, retryIn := CheckWAF()
    if err == nil {
        return false
    }
    writeWAFResponse(w, err, retryIn)
    return true
}

// writeWAFResponse emits the uniform 503 + Retry-After answer for a blocked
// server. The body is a small structured JSON both OpenAI- and Anthropic-
// style clients can parse — never the raw block page.
func writeWAFResponse(w http.ResponseWriter, err error, retryIn time.Duration) {
    w.Header().Set("Retry-After", strconvMaxInt64(int64(retryIn.Seconds()+0.5)))
    writeJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
        "error": map[string]interface{}{
            "type":    "overloaded_error",
            "code":    "waf_block",
            "message": err.Error(),
            "retryIn": retryIn.Round(time.Second).String(),
        },
    })
}

func strconvMaxInt64(v int64) string {
    return fmt.Sprintf("%d", v)
}

// wafRetryHint reports the current cooldown so error messages can carry an
// actionable "retry in" suggestion.
func wafRetryHint() time.Duration {
    wafBreaker.mu.Lock()
    defer wafBreaker.mu.Unlock()
    if wafBreaker.state == breakerClosed {
        return wafBaseCooldown
    }
    elapsed := time.Since(wafBreaker.openedAt)
    remaining := wafBreaker.cooldown - elapsed
    if remaining < 0 {
        remaining = 0
    }
    return remaining
}

// truncateErrorBody keeps non-WAF upstream error bodies short in client-
// facing messages: a full HTML page helps nobody, and agents choke on it
// (issue #41). JSON bodies pass through untouched up to a hard cap.
func truncateErrorBody(body []byte) string {
    s := strings.TrimSpace(string(body))
    if len(s) <= 512 {
        return s
    }
    return s[:512] + "...(truncated)"
}

func minDuration(a, b time.Duration) time.Duration {
    if a < b {
        return a
    }
    return b
}

// logWAFBlock prints a single condensed line on state transitions only, so a
// long block does not spam the log with the raw 3 KB block page.
func logWAFBlock(msg string) {
    // Uses the standard log path (not the verbose-gated logInfo) on purpose:
    // an IP-level block is an operational event operators must see.
    log.Printf("[WAF] %s", msg)
}

// readProbeBody reads up to 8 KB of the probe response — enough to carry the
// marker text, small enough to never matter.
func readProbeBody(resp *http.Response) []byte {
    if resp.Body == nil {
        return nil
    }
    buf := make([]byte, 8192)
    n, _ := io.ReadFull(resp.Body, buf)
    if n <= 0 {
        // Body shorter than the buffer (or empty): io.ReadFull returns
        // io.ErrUnexpectedEOF, which is fine here.
    }
    if n > 8192 {
        n = 8192
    }
    return buf[:n]
}
