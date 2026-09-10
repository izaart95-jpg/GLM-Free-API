package zbridge

import (
    "context"
    "net/http"
    "net/http/httptest"
    "sync"
    "testing"
    "time"
)

// recordingTransport notes when each request reached the underlying transport.
type recordingTransport struct {
    mu    sync.Mutex
    times []time.Time
}

func (r *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    r.mu.Lock()
    r.times = append(r.times, time.Now())
    r.mu.Unlock()
    return &http.Response{
        StatusCode: 200,
        Body:       http.NoBody,
        Request:    req,
    }, nil
}

func (r *recordingTransport) gaps() []time.Duration {
    r.mu.Lock()
    defer r.mu.Unlock()
    var out []time.Duration
    for i := 1; i < len(r.times); i++ {
        out = append(out, r.times[i].Sub(r.times[i-1]))
    }
    return out
}

func (r *recordingTransport) snapshot() []time.Time {
    r.mu.Lock()
    defer r.mu.Unlock()
    return append([]time.Time{}, r.times...)
}

func TestPacedTransportSpacesSequentialRequests(t *testing.T) {
    rec := &recordingTransport{}
    const gap = 40 * time.Millisecond
    rt := pacedWithGap(rec, gap)

    for i := 0; i < 4; i++ {
        req := httptest.NewRequest("GET", "https://example.invalid/", nil)
        if _, err := rt.RoundTrip(req); err != nil {
            t.Fatalf("request %d: %v", i, err)
        }
    }
    for i, got := range rec.gaps() {
        // Timers fire no earlier than requested; allow slack for scheduling.
        if got < gap-5*time.Millisecond {
            t.Fatalf("gap %d too small: %v, want >= %v", i, got, gap)
        }
    }
}

// A burst is the shape that trips the WAF, so concurrent callers must be
// spread out too, not merely delayed by the same amount.
func TestPacedTransportSpreadsBurst(t *testing.T) {
    rec := &recordingTransport{}
    const gap = 25 * time.Millisecond
    const n = 6
    rt := pacedWithGap(rec, gap)

    var wg sync.WaitGroup
    start := time.Now()
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            req := httptest.NewRequest("GET", "https://example.invalid/", nil)
            if _, err := rt.RoundTrip(req); err != nil {
                t.Errorf("burst request: %v", err)
            }
        }()
    }
    wg.Wait()

    if got := len(rec.times); got != n {
        t.Fatalf("want %d requests through, got %d", n, got)
    }
    // n requests at one per gap cannot finish faster than (n-1) gaps.
    if elapsed := time.Since(start); elapsed < time.Duration(n-1)*gap-10*time.Millisecond {
        t.Fatalf("burst was not paced: %d requests in %v", n, elapsed)
    }
}

func TestPacedTransportHonoursContextCancellation(t *testing.T) {
    rec := &recordingTransport{}
    rt := pacedWithGap(rec, time.Second)

    // First request claims the slot immediately, the second must wait a second.
    if _, err := rt.RoundTrip(httptest.NewRequest("GET", "https://example.invalid/", nil)); err != nil {
        t.Fatalf("first request: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
    defer cancel()
    req := httptest.NewRequest("GET", "https://example.invalid/", nil).WithContext(ctx)
    if _, err := rt.RoundTrip(req); err == nil {
        t.Fatal("want the wait to abort with the request context, got nil error")
    }
    if got := len(rec.times); got != 1 {
        t.Fatalf("cancelled request still reached upstream: %d calls", got)
    }
}

// A zero interval must hand the request straight to the base transport.
func TestPacedTransportDisabled(t *testing.T) {
    rec := &recordingTransport{}
    rt := pacedWithGap(rec, 0)
    start := time.Now()
    for i := 0; i < 3; i++ {
        if _, err := rt.RoundTrip(httptest.NewRequest("GET", "https://example.invalid/", nil)); err != nil {
            t.Fatalf("request %d: %v", i, err)
        }
    }
    if elapsed := time.Since(start); elapsed > 20*time.Millisecond {
        t.Fatalf("pacing still applied with a zero interval: %v", elapsed)
    }
}

// The interval is read on first use, so a value set after construction still
// takes effect — that is what lets TestMain switch pacing off for the mock.
func TestPacedTransportReadsIntervalLazily(t *testing.T) {
    rec := &recordingTransport{}
    rt := newPacedTransport(rec)
    t.Setenv("UPSTREAM_MIN_INTERVAL_MS", "0")
    start := time.Now()
    for i := 0; i < 3; i++ {
        if _, err := rt.RoundTrip(httptest.NewRequest("GET", "https://example.invalid/", nil)); err != nil {
            t.Fatalf("request %d: %v", i, err)
        }
    }
    if elapsed := time.Since(start); elapsed > 20*time.Millisecond {
        t.Fatalf("environment set after construction was ignored: %v", elapsed)
    }
}

func TestUpstreamMinIntervalFromEnv(t *testing.T) {
    cases := map[string]time.Duration{
        "0":      0,
        "750":    750 * time.Millisecond,
        "100":    100 * time.Millisecond,
    }
    for raw, want := range cases {
        t.Setenv("UPSTREAM_MIN_INTERVAL_MS", raw)
        if got := upstreamMinInterval(); got != want {
            t.Fatalf("UPSTREAM_MIN_INTERVAL_MS=%q: got %v, want %v", raw, got, want)
        }
    }

    // Unset / unparsable / negative values fall back to a random draw in
    // [200ms, 500ms] — asserted below, not a single expected value.
    for _, raw := range []string{"", "nonsen", "-5"} {
        t.Setenv("UPSTREAM_MIN_INTERVAL_MS", raw)
        for i := 0; i < 50; i++ {
            got := upstreamMinInterval()
            if got < minUpstreamIntervalMS*time.Millisecond || got > maxUpstreamIntervalMS*time.Millisecond {
                t.Fatalf("UPSTREAM_MIN_INTERVAL_MS=%q: got %v, want within [%dms, %dms]", raw, got, minUpstreamIntervalMS, maxUpstreamIntervalMS)
            }
        }
    }
}

// pacedWithGap builds a transport with a fixed interval, bypassing the
// environment lookup.
func pacedWithGap(base http.RoundTripper, gap time.Duration) http.RoundTripper {
    return &pacedTransport{base: base, gate: &pacingGate{gapFor: func() time.Duration { return gap }}}
}

// The default draw must be crypto-random in [200ms, 500ms]. This cannot assert
// unpredictability (that is the point of crypto/rand), but it can assert the
// range and that repeated draws are not all identical — a pinned constant or a
// broken rand would fail the latter, since 200 draws of a uniform [200,500]
// draw colliding entirely is astronomically unlikely.
func TestUpstreamIntervalRandomDraw(t *testing.T) {
    t.Setenv("UPSTREAM_MIN_INTERVAL_MS", "")
    seen := make(map[time.Duration]bool)
    for i := 0; i < 200; i++ {
        got := upstreamMinInterval()
        if got < minUpstreamIntervalMS*time.Millisecond || got > maxUpstreamIntervalMS*time.Millisecond {
            t.Fatalf("draw %d out of range: %v", i, got)
        }
        seen[got] = true
    }
    if len(seen) < 2 {
        t.Fatalf("all %d draws identical — the interval is not random", len(seen))
    }
}

// Issue #41: with per-transport pacing, alternating calls between the Z.AI
// client and the Aliyun captcha client each consumed a slot from their OWN
// transport's clock, so the real gap between consecutive upstream requests
// collapsed toward zero even though both transports were individually paced.
// The shared gate must keep the minimum gap across transports.
func TestSharedGatePacesAcrossTransports(t *testing.T) {
    const gap = 30 * time.Millisecond
    zaiRec := &recordingTransport{}
    aliRec := &recordingTransport{}
    // Both transports deliberately share ONE gate, as production does.
    gate := &pacingGate{gapFor: func() time.Duration { return gap }}
    zaiRT := &pacedTransport{base: zaiRec, gate: gate}
    aliRT := &pacedTransport{base: aliRec, gate: gate}

    // Interleave exactly like the real request flow does: captcha call
    // (Aliyun), completion POST (Z.AI), chat delete (Z.AI), pool refill...
    for i := 0; i < 6; i++ {
        req := httptest.NewRequest("GET", "https://example.invalid/", nil)
        if _, err := aliRT.RoundTrip(req); err != nil {
            t.Fatalf("aliyun call %d: %v", i, err)
        }
        if _, err := zaiRT.RoundTrip(req); err != nil {
            t.Fatalf("zai call %d: %v", i, err)
        }
    }

    // Merge the recorded timestamps from BOTH transports and assert that the
    // minimum gap between any two consecutive upstream requests respects the
    // configured floor.
    all := append(append([]time.Time{}, zaiRec.snapshot()...), aliRec.snapshot()...)
    if len(all) != 12 {
        t.Fatalf("want 12 upstream calls, got %d", len(all))
    }
    sortTimes(all)
    for i := 1; i < len(all); i++ {
        got := all[i].Sub(all[i-1])
        if got < gap-8*time.Millisecond {
            t.Fatalf("cross-transport gap %d collapsed: %v, want >= %v (per-transport clocks interleaved)", i, got, gap)
        }
    }
}

func sortTimes(ts []time.Time) {
    for i := 1; i < len(ts); i++ {
        for j := i; j > 0 && ts[j].Before(ts[j-1]); j-- {
            ts[j], ts[j-1] = ts[j-1], ts[j]
        }
    }
}
