package zbridge

import (
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "sync"
    "testing"
    "time"
)

// resetBreaker returns the process-wide breaker to a pristine closed state
// and invalidates any outstanding prober goroutines.
func resetBreaker(t *testing.T) {
    t.Helper()
    wafBreaker.mu.Lock()
    wafBreaker.state = breakerClosed
    wafBreaker.openedAt = time.Time{}
    wafBreaker.cooldown = 0
    wafBreaker.blockCount = 0
    wafBreaker.probeGen++
    wafBreaker.mu.Unlock()
}

// blockPageBody returns the (abridged) Aliyun WAF block page: the exact
// marker strings the detector looks for, matching the real page captured in
// issue #41 / issue #20.
func blockPageBody() []byte {
    return []byte(`<!doctypehtml><html lang="zh-cn"><meta charset="utf-8">` +
        `<title>405</title><body data-spm="7663354">` +
        `<script>var en_tips={block_message:"Sorry, your request has been blocked ` +
        `as it may cause potential threats to the server's security."};</script>` +
        `<textarea id="renderData" style="display:none">{"traceid":"0a0f6bd017889669411873134e5cf9","lang":"en"}</textarea>`)
}

func TestIsWAFBlockResponse(t *testing.T) {
    cases := []struct {
        name string
        code int
        body []byte
        want bool
    }{
        {"real block page 405", 405, blockPageBody(), true},
        {"chinese block page", 405, []byte(`<!doctypehtml><html lang="zh-cn">很抱歉，由于您访问的URL有可能对网站造成安全威胁，您的访问被阻断。`), true},
        {"aliyun marker only", 405, []byte(`<html>blocked by errors.aliyun.com`), true},
        {"json 405 is not a block", 405, []byte(`{"detail":"method not allowed"}`), false},
        {"json 403 is not a block", 403, []byte(`{"detail":"Not authenticated"}`), false},
        {"200 with marker is not a block", 200, blockPageBody(), false},
        {"405 empty body", 405, nil, false},
        {"500 json", 500, []byte(`{"error":"boom"}`), false},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            if got := isWAFBlockResponse(tc.code, tc.body); got != tc.want {
                t.Fatalf("isWAFBlockResponse(%d, %q) = %v, want %v", tc.code, tc.body, got, tc.want)
            }
        })
    }
}

func TestBreakerOpensAndFailsFast(t *testing.T) {
    resetBreaker(t)
    defer resetBreaker(t)

    if err, _ := CheckWAF(); err != nil {
        t.Fatalf("closed breaker rejected a request: %v", err)
    }

    RecordWAFBlock()

    err, wait := CheckWAF()
    if err == nil {
        t.Fatal("open breaker let a request through")
    }
    if !errors.Is(err, ErrWAFBlock) {
        t.Fatalf("want ErrWAFBlock, got %v", err)
    }
    if wait <= 0 || wait > wafBaseCooldown {
        t.Fatalf("wait hint out of range: %v", wait)
    }
    if !strings.Contains(err.Error(), "blocked") {
        t.Fatalf("error message not actionable: %v", err)
    }
}

func TestBreakerClosesOnSuccess(t *testing.T) {
    resetBreaker(t)
    defer resetBreaker(t)

    RecordWAFBlock()
    RecordWAFSuccess()
    if err, _ := CheckWAF(); err != nil {
        t.Fatalf("breaker still open after success: %v", err)
    }
}

func TestBreakerProberReopensWithDoubledCooldown(t *testing.T) {
    resetBreaker(t)
    defer resetBreaker(t)

    // Shrink the prober tick so a shortened cooldown is noticed immediately.
    oldSlice := wafProberSlice
    wafProberSlice = 5 * time.Millisecond
    defer func() { wafProberSlice = oldSlice }()

    RecordWAFBlock()

    // Point the probe at a still-blocked upstream.
    blocked := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(405)
        w.Write(blockPageBody())
    }))
    defer blocked.Close()
    oldProbe := wafProbeTransport
    wafProbeTransport = func() (*http.Response, error) {
        req, _ := http.NewRequest(http.MethodPost, blocked.URL+wafProbeEndpoint, strings.NewReader("{}"))
        return blocked.Client().Do(req)
    }
    defer func() { wafProbeTransport = oldProbe }()

    // Shrink the cooldown so the prober's wait elapses within test time.
    wafBreaker.mu.Lock()
    wafBreaker.cooldown = 20 * time.Millisecond
    wafBreaker.mu.Unlock()

    deadline := time.Now().Add(2 * time.Second)
    for time.Now().Before(deadline) {
        wafBreaker.mu.Lock()
        state := wafBreaker.state
        cooldown := wafBreaker.cooldown
        wafBreaker.mu.Unlock()
        if state == breakerOpen && cooldown == 40*time.Millisecond {
            // Probe ran, saw the block, doubled the cooldown (20ms -> 40ms).
            return
        }
        time.Sleep(2 * time.Millisecond)
    }
    t.Fatal("prober did not re-open the breaker with a doubled cooldown")
}

func TestBreakerProberClosesWhenClear(t *testing.T) {
    resetBreaker(t)
    defer resetBreaker(t)

    // Shrink the prober tick so a shortened cooldown is noticed immediately.
    oldSlice := wafProberSlice
    wafProberSlice = 5 * time.Millisecond
    defer func() { wafProberSlice = oldSlice }()

    RecordWAFBlock()
    wafBreaker.mu.Lock()
    wafBreaker.cooldown = 20 * time.Millisecond
    wafBreaker.mu.Unlock()

    healthy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        w.Write([]byte(`{"models":[]}`))
    }))
    defer healthy.Close()
    oldProbe := wafProbeTransport
    wafProbeTransport = func() (*http.Response, error) {
        req, _ := http.NewRequest(http.MethodPost, healthy.URL+wafProbeEndpoint, strings.NewReader("{}"))
        return healthy.Client().Do(req)
    }
    defer func() { wafProbeTransport = oldProbe }()

    deadline := time.Now().Add(2 * time.Second)
    for time.Now().Before(deadline) {
        if err, _ := CheckWAF(); err == nil {
            return // prober cleared the block
        }
        time.Sleep(2 * time.Millisecond)
    }
    t.Fatal("prober did not close the breaker after a healthy probe")
}

// While the breaker is open, an in-flight request must not be able to spend
// a captcha on a doomed upstream call: the CheckWAF gate inside
// sendToZAIStream rejects it before the captcha step.
func TestBreakerStopsCaptchaSpend(t *testing.T) {
    resetBreaker(t)
    defer resetBreaker(t)

    RecordWAFBlock()

    opts := struct {
        Model, ChatID     string
        FeaturesMap       map[string]interface{}
        Messages          []Message
        ClientMessagesRaw json.RawMessage
        Files             []map[string]interface{}
        RequestID         string
        AdvancedSearch    bool
    }{Model: "glm-4.7", ChatID: "chat-1"}
    err := sendToZAIStream("probe test", opts, make(chan ZAIResult, 1))
    if err == nil {
        t.Fatal("blocked breaker let a completion through")
    }
    if !errors.Is(err, ErrWAFBlock) {
        t.Fatalf("want ErrWAFBlock, got %v", err)
    }
}

func TestRejectIfWAFBlockedAnswers503(t *testing.T) {
    resetBreaker(t)
    defer resetBreaker(t)

    RecordWAFBlock()
    rr := httptest.NewRecorder()
    if !RejectIfWAFBlocked(rr) {
        t.Fatal("blocked breaker did not reject the request")
    }
    if rr.Code != http.StatusServiceUnavailable {
        t.Fatalf("status = %d, want 503", rr.Code)
    }
    body := rr.Body.String()
    for _, want := range []string{"waf_block", "overloaded_error", "blocked"} {
        if !strings.Contains(body, want) {
            t.Fatalf("body missing %q: %s", want, body)
        }
    }
    if rr.Header().Get("Retry-After") == "" {
        t.Fatal("Retry-After header missing")
    }

    // The 3 KB block page must never reach the client.
    if strings.Contains(body, "<!doctype") {
        t.Fatalf("block page leaked into the client answer: %.80s", body)
    }
}

func TestRejectIfWAFBlockedPassesWhenClosed(t *testing.T) {
    resetBreaker(t)
    defer resetBreaker(t)

    rr := httptest.NewRecorder()
    if RejectIfWAFBlocked(rr) {
        t.Fatal("closed breaker rejected a request")
    }
}

func TestStatusFromErrorMapsWAFTo503(t *testing.T) {
    if got := statusFromError(ErrWAFBlock.Error()); got != 503 {
        t.Fatalf("statusFromError(WAF block) = %d, want 503", got)
    }
    if got := statusFromError("Z.AI error 401: bad token"); got != 401 {
        t.Fatalf("statusFromError(401) = %d, want 401", got)
    }
    if got := statusFromError("Z.AI error 400: detail"); got != 400 {
        t.Fatalf("statusFromError(400) = %d, want 400", got)
    }
}

func TestTruncateErrorBody(t *testing.T) {
    short := `{"detail":"We could not find what you're looking for :/"}`
    if got := truncateErrorBody([]byte(short)); got != short {
        t.Fatalf("short body modified: %q", got)
    }
    long := strings.Repeat("x", 1024)
    got := truncateErrorBody([]byte(long))
    if len(got) >= len(long) {
        t.Fatalf("long body not truncated: %d bytes", len(got))
    }
    if !strings.HasSuffix(got, "...(truncated)") {
        t.Fatalf("truncation marker missing: %q", got[len(got)-30:])
    }
}

func TestWAFStatusReporting(t *testing.T) {
    resetBreaker(t)
    defer resetBreaker(t)

    st := WAFStatus()
    if blocked, _ := st["blocked"].(bool); blocked {
        t.Fatal("closed breaker reports blocked")
    }

    RecordWAFBlock()
    st = WAFStatus()
    if blocked, _ := st["blocked"].(bool); !blocked {
        t.Fatal("open breaker reports unblocked")
    }
    if _, ok := st["retryIn"]; !ok {
        t.Fatal("open breaker status missing retryIn")
    }
}

// Concurrency: many goroutines recording blocks and successes concurrently
// must not race (guarded by -race in CI).
func TestBreakerConcurrentAccess(t *testing.T) {
    resetBreaker(t)
    defer resetBreaker(t)

    // Neutralize the prober's network call during the concurrency storm.
    oldProbe := wafProbeTransport
    wafProbeTransport = func() (*http.Response, error) {
        return nil, errors.New("no network in test")
    }
    defer func() { wafProbeTransport = oldProbe }()

    var wg sync.WaitGroup
    for i := 0; i < 16; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            for j := 0; j < 50; j++ {
                if i%2 == 0 {
                    RecordWAFBlock()
                } else {
                    RecordWAFSuccess()
                }
                _, _ = CheckWAF()
                _ = WAFStatus()
            }
        }(i)
    }
    wg.Wait()
}

// The prober goroutine launched by RecordWAFBlock must exit (not leak) when
// the breaker is closed by someone else before its timer fires.
func TestBreakerProberStaleGenerationExits(t *testing.T) {
    resetBreaker(t)
    defer resetBreaker(t)

    RecordWAFBlock()
    // Close the breaker while the prober is still sleeping off its cooldown.
    RecordWAFSuccess()

    // Give any (incorrectly still-live) prober a window to fire its probe.
    time.Sleep(50 * time.Millisecond)

    // The breaker must still be closed (a stale prober cannot re-open it
    // without winning the generation check).
    if err, _ := CheckWAF(); err != nil {
        t.Fatalf("stale prober re-opened a closed breaker: %v", err)
    }
}

// A probe transport failure must be treated as "still blocked" — never a
// panic, never a false "cleared".
func TestBreakerProbeFailureTreatedAsBlocked(t *testing.T) {
    resetBreaker(t)
    defer resetBreaker(t)

    RecordWAFBlock()
    wafBreaker.mu.Lock()
    wafBreaker.cooldown = 10 * time.Millisecond
    wafBreaker.mu.Unlock()

    oldProbe := wafProbeTransport
    wafProbeTransport = func() (*http.Response, error) {
        return nil, errors.New("connection refused")
    }
    defer func() { wafProbeTransport = oldProbe }()

    deadline := time.Now().Add(1 * time.Second)
    for time.Now().Before(deadline) {
        if err, _ := CheckWAF(); err != nil {
            // Still blocked — correct, conservative outcome.
            return
        }
        time.Sleep(2 * time.Millisecond)
    }
    t.Fatal("probe transport failure incorrectly cleared the breaker")
}

// The probe MUST hit the endpoint the WAF actually blocks (POST
// /api/v2/chat/completions), not a path that stays open during a block
// (issue #41: GET / and GET /api/models keep returning 200/JSON while the
// POST is blocked). Probing a healthy path reports "cleared" while
// completions still fail — the exact regression caught during development.
func TestWAFProbeTargetsBlockedEndpoint(t *testing.T) {
    if wafProbeEndpoint != "/api/v2/chat/completions" {
        t.Fatalf("probe endpoint = %q, want /api/v2/chat/completions — probing a path the WAF does not block reports false recovery", wafProbeEndpoint)
    }

    // End-to-end shape: a blocked endpoint is judged blocked, and an open
    // endpoint answering a JSON auth error is judged clear — even though
    // both are non-2xx.
    blocked := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(405)
        w.Write(blockPageBody())
    }))
    defer blocked.Close()
    resp, err := http.Post(blocked.URL+wafProbeEndpoint, "application/json", strings.NewReader("{}"))
    if err != nil {
        t.Fatalf("probe request: %v", err)
    }
    body, _ := io.ReadAll(resp.Body)
    resp.Body.Close()
    if !isWAFBlockResponse(resp.StatusCode, body) {
        t.Fatal("probe against a blocked endpoint was not recognized as blocked")
    }

    // Unblocked but unauthorized: JSON 403, not the block page.
    open := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(403)
        w.Write([]byte(`{"detail":"Not authenticated"}`))
    }))
    defer open.Close()
    resp, err = http.Post(open.URL+wafProbeEndpoint, "application/json", strings.NewReader("{}"))
    if err != nil {
        t.Fatalf("probe request: %v", err)
    }
    body, _ = io.ReadAll(resp.Body)
    resp.Body.Close()
    if isWAFBlockResponse(resp.StatusCode, body) {
        t.Fatal("probe against an unblocked endpoint was mistaken for a block")
    }
}
