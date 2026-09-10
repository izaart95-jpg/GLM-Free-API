package zbridge

import (
    "crypto/rand"
    "math/big"
    "net/http"
    "os"
    "strconv"
    "sync"
    "time"
)

// Upstream pacing.
//
// One proxied chat request costs several upstream calls — a captcha challenge,
// the completion itself, and the DELETE that retires the throwaway chat — and
// the session pool refills in the background on top of that. Bursts therefore
// hit chat.z.ai far harder than the request count suggests, and Aliyun's WAF
// answers a burst by blocking the source IP: every later request comes back as
// the 405 block page, for hours, no matter which account or session is used
// (issues #20, #41).
//
// Pacing every upstream call through ONE SHARED gate keeps a burst from ever
// forming. It is deliberately a floor on the gap between requests rather than
// a token bucket: a bucket lets a burst through as long as the average holds,
// which is exactly the shape that trips the WAF.
//
// The gate is shared across all upstream transports (the Z.AI client and the
// Aliyun captcha client). Issue #41 showed why: per-transport gates each keep
// their own "next slot" clock, so alternating calls between two transports
// interleave their slots and the real gap between consecutive upstream
// requests collapses toward zero even when each transport is individually
// paced — sustained agent traffic interleaves exactly this way.
//
// UPSTREAM_MIN_INTERVAL_MS pins the interval explicitly; 0 disables pacing
// entirely. Unset, the interval is drawn once per transport from a
// cryptographically random 200–500 ms — a fixed gap is itself a fingerprint,
// and crypto/rand (not math/rand) so the draw is not guessable from the
// process's seed.
const (
    minUpstreamIntervalMS = 200
    maxUpstreamIntervalMS = 500
)

// randomUpstreamInterval returns a cryptographically random duration in
// [200ms, 500ms]. Falls back to the midpoint if the entropy source is
// exhausted (crypto/rand never returns an error here, but a dead source
// must not take pacing down with it).
func randomUpstreamInterval() time.Duration {
    n, err := rand.Int(rand.Reader, big.NewInt(maxUpstreamIntervalMS-minUpstreamIntervalMS+1))
    if err != nil {
        return time.Duration((minUpstreamIntervalMS+maxUpstreamIntervalMS)/2) * time.Millisecond
    }
    return time.Duration(minUpstreamIntervalMS+n.Int64()) * time.Millisecond
}

func upstreamMinInterval() time.Duration {
    raw := os.Getenv("UPSTREAM_MIN_INTERVAL_MS")
    if raw == "" {
        return randomUpstreamInterval()
    }
    ms, err := strconv.Atoi(raw)
    if err != nil || ms < 0 {
        return randomUpstreamInterval()
    }
    return time.Duration(ms) * time.Millisecond
}

// pacingGate is the process-wide slot allocator every upstream transport
// funnels through. A single mutex-protected "next allowed start" makes the
// minimum gap hold across ALL concurrent upstream calls, no matter which
// transport issues them.
type pacingGate struct {
    gapOnce sync.Once
    gap     time.Duration
    gapFor  func() time.Duration

    mu   sync.Mutex
    next time.Time // earliest instant the next request may start
}

// sharedPacingGate is the one gate all upstream transports share.
var sharedPacingGate = &pacingGate{gapFor: upstreamMinInterval}

func (g *pacingGate) minGap() time.Duration {
    g.gapOnce.Do(func() {
        gapFor := g.gapFor
        if gapFor == nil {
            gapFor = upstreamMinInterval
        }
        g.gap = gapFor()
    })
    return g.gap
}

// claim reserves the next start slot and reports how long the caller must
// wait before issuing its request.
func (g *pacingGate) claim() time.Duration {
    minGap := g.minGap()
    if minGap <= 0 {
        return 0
    }
    g.mu.Lock()
    now := time.Now()
    slot := g.next
    if slot.Before(now) {
        slot = now
    }
    g.next = slot.Add(minGap)
    g.mu.Unlock()
    return time.Until(slot)
}

// pacedTransport spaces out the requests handed to base through the shared
// pacing gate. Callers block in RoundTrip until their slot is due, so pacing
// applies no matter which code path issues the call — and, because every
// upstream transport wraps itself with the SAME gate, no matter which
// transport fires either.
type pacedTransport struct {
    base http.RoundTripper
    gate *pacingGate
}

// newPacedTransport wraps base with the shared pacing gate.
func newPacedTransport(base http.RoundTripper) http.RoundTripper {
    return &pacedTransport{base: base, gate: sharedPacingGate}
}

func (t *pacedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    if wait := t.gate.claim(); wait > 0 {
        timer := time.NewTimer(wait)
        defer timer.Stop()
        select {
        case <-timer.C:
        case <-req.Context().Done():
            return nil, req.Context().Err()
        }
    }
    return t.base.RoundTrip(req)
}
