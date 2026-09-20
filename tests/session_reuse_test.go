package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"zai-api/internal/zbridge"
)

// Sync sticky reuse: same chat_id until SESSION_REUSE_COUNT, then rotate.
func TestSyncStickyReusesUpToMax(t *testing.T) {
	defer zbridge.OverrideSessionReuseCount(3)()
	restore := zbridge.AttachSessionPool(nil, time.Second)
	defer restore()
	zbridge.ResetSyncStickySession()

	// Mock DELETE so rotation GC has somewhere to land.
	var mu sync.Mutex
	deleted := []string{}
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/api/v1/chats/") {
			mu.Lock()
			deleted = append(deleted, strings.TrimPrefix(r.URL.Path, "/api/v1/chats/"))
			mu.Unlock()
			w.WriteHeader(200)
			w.Write([]byte("true"))
			return
		}
		http.NotFound(w, r)
	}))
	defer upstream.Close()
	oldBase := zbridge.BASE_URL
	zbridge.BASE_URL = upstream.URL
	defer func() { zbridge.BASE_URL = oldBase }()
	defer zbridge.OverrideSessionState("test-token", "test-user", true)()

	id1, _, err := zbridge.AcquireStatelessSession(context.Background())
	if err != nil {
		t.Fatalf("acquire 1: %v", err)
	}
	zbridge.ReleaseStatelessSession(id1, false)

	id2, _, err := zbridge.AcquireStatelessSession(context.Background())
	if err != nil {
		t.Fatalf("acquire 2: %v", err)
	}
	if id2 != id1 {
		t.Fatalf("second acquire = %q, want sticky reuse %q", id2, id1)
	}
	zbridge.ReleaseStatelessSession(id2, false)

	// Still same sticky for the 3rd request.
	id3, _, err := zbridge.AcquireStatelessSession(context.Background())
	if err != nil {
		t.Fatalf("acquire 3: %v", err)
	}
	if id3 != id1 {
		t.Fatalf("third acquire = %q, want sticky %q", id3, id1)
	}
	zbridge.ReleaseStatelessSession(id3, false) // hits max=3 -> rotate + GC delete

	waitFor(t, "sticky rotation delete", 5*time.Second, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(deleted) == 1 && deleted[0] == id1
	})

	id4, _, err := zbridge.AcquireStatelessSession(context.Background())
	if err != nil {
		t.Fatalf("acquire 4: %v", err)
	}
	if id4 == id1 {
		t.Fatalf("after max uses sticky did not rotate (still %q)", id4)
	}
}

// Handler-level pool reuse: 3 completions share one chat_id, no DELETE until max.
func TestHTTPChatCompletionsReusesPooledSession(t *testing.T) {
	defer zbridge.OverrideSessionReuseCount(3)()

	var mu sync.Mutex
	chatIDs := []string{}
	deleted := []string{}
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v2/chat/completions":
			var body struct {
				ChatID string `json:"chat_id"`
			}
			json.NewDecoder(r.Body).Decode(&body)
			mu.Lock()
			chatIDs = append(chatIDs, body.ChatID)
			mu.Unlock()
			w.Header().Set("Content-Type", "text/event-stream")
			flusher, _ := w.(http.Flusher)
			fmt.Fprintf(w, "data: {\"data\":{\"delta_content\":\"hi\"}}\n\n")
			fmt.Fprintf(w, "data: {\"data\":{\"phase\":\"done\"}}\n\n")
			fmt.Fprintf(w, "data: [DONE]\n\n")
			if flusher != nil {
				flusher.Flush()
			}
		case r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/api/v1/chats/"):
			mu.Lock()
			deleted = append(deleted, strings.TrimPrefix(r.URL.Path, "/api/v1/chats/"))
			mu.Unlock()
			w.WriteHeader(200)
			w.Write([]byte("true"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer upstream.Close()
	oldBase := zbridge.BASE_URL
	zbridge.BASE_URL = oldBase
	zbridge.BASE_URL = upstream.URL
	defer func() { zbridge.BASE_URL = oldBase }()
	defer zbridge.OverrideSessionState("test-token", "test-user", true)()

	cfg := zbridge.GetConfig()
	oldAgent := cfg.AgentMode
	cfg.AgentMode = true
	defer func() { cfg.AgentMode = oldAgent }()
	zbridge.SeedCaptchaParam("test-captcha-param")
	zbridge.SeedCaptchaParam("test-captcha-param")
	zbridge.SeedCaptchaParam("test-captcha-param")
	defer zbridge.OverrideCaptchaParam("test-captcha-param")()
	// Ensure non-agent path also bypasses captcha (Seed only serves agent cache).
	// The Override above forces the value regardless of mode; restore after.

	pool := zbridge.NewSessionPool(zbridge.NewZAIChatBackend(), 1)
	pool.SetMaxUses(3)
	restore := zbridge.AttachSessionPool(pool, time.Second)
	defer func() {
		pool.Shutdown()
		restore()
	}()
	pool.Start()
	waitFor(t, "warmup", 3*time.Second, func() bool { return pool.Ready() == 1 })

	doReq := func() {
		t.Helper()
		body := `{"model":"glm-4.7","stream":false,"messages":[{"role":"user","content":"hi"}]}`
		req := httptest.NewRequest("POST", "/v1/chat/completions", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+cfg.Auth.Token)
		rec := httptest.NewRecorder()
		zbridge.NewHandler().ServeHTTP(rec, req)
		if rec.Code != 200 {
			t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
		}
	}

	doReq()
	doReq()

	mu.Lock()
	if len(chatIDs) != 2 {
		mu.Unlock()
		t.Fatalf("upstream saw %d completions, want 2", len(chatIDs))
	}
	if chatIDs[0] != chatIDs[1] {
		mu.Unlock()
		t.Fatalf("first two completions used different chats %v, want reuse", chatIDs)
	}
	if len(deleted) != 0 {
		mu.Unlock()
		t.Fatalf("deleted=%v after 2/3 uses, want none yet", deleted)
	}
	mu.Unlock()

	// Third use hits max -> DELETE + refill. Poll for it.
	doReq()
	waitFor(t, "delete after max uses", 5*time.Second, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(deleted) == 1 && deleted[0] == chatIDs[0] && pool.Ready() == 1
	})
}
