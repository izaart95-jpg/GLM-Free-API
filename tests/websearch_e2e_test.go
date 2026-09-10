// Blackbox end-to-end regression tests for issue #42 through the REAL HTTP
// stack, driven via the exported zbridge.NewHandler() surface:
//
//   NewHandler -> authMiddleware -> chatCompletionsHandler -> sendToZAI ->
//   sendToZAIStream -> upstream POST body on the wire.
//
// These pin the externally visible contract of the search flags:
//   - webSearch=true  -> features.auto_web_search=true, features.web_search=false
//   - advancedSearch=true -> top-level mcp_servers:["advanced-search"]
//   - agent mode on -> both search toggles are force-disabled server-side
//     (the request still succeeds — it just never asks Z.AI to search).
//
// The upstream is a local httptest mock; the captcha machinery is bypassed
// via the captchaParamOverride seam (whitebox tests own it; here it is set
// through the exported test surface indirectly by keeping agent mode on for
// the gate test and using the cache for the others).

package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"zai-api/internal/zbridge"
)

// captureBodyUpstream records the last /api/v2/chat/completions JSON body.
type capturedBody struct {
	mu   sync.Mutex
	body map[string]interface{}
}

func (c *capturedBody) record(raw []byte) error {
	var parsed map[string]interface{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return err
	}
	c.mu.Lock()
	c.body = parsed
	c.mu.Unlock()
	return nil
}

func (c *capturedBody) get(t *testing.T) map[string]interface{} {
	t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.body == nil {
		t.Fatal("upstream never received a completion request")
	}
	return c.body
}

func newCaptureUpstream(t *testing.T) (*httptest.Server, *capturedBody) {
	t.Helper()
	cap := &capturedBody{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/chat/completions" {
			http.NotFound(w, r)
			return
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := cap.record(raw); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		fmt.Fprintf(w, "data: {\"data\":{\"delta_content\":\"ok\"}}\n\n")
		fmt.Fprintf(w, "data: {\"data\":{\"phase\":\"done\"}}\n\n")
		fmt.Fprintf(w, "data: [DONE]\n\n")
		if flusher != nil {
			flusher.Flush()
		}
	}))
	t.Cleanup(srv.Close)
	return srv, cap
}

// driveChatCompletion posts a /v1/chat/completions body through the real
// handler and returns the response recorder.
func driveChatCompletion(t *testing.T, cfg *zbridge.Config, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest("POST", "/v1/chat/completions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Auth.Token)
	rec := httptest.NewRecorder()
	zbridge.NewHandler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	return rec
}

// endToEndSearchEnv points the bridge at the mock, pre-seeds the session
// (skipping guest auth), and bypasses the captcha machinery for BOTH the
// agent and non-agent request paths (OverrideCaptchaParam). Returns the
// config so callers can toggle agent mode and restore it.
func endToEndSearchEnv(t *testing.T, upstream *httptest.Server) *zbridge.Config {
	t.Helper()

	oldBase := zbridge.BASE_URL
	zbridge.BASE_URL = upstream.URL
	t.Cleanup(func() { zbridge.BASE_URL = oldBase })

	restore := zbridge.OverrideSessionState("test-token", "test-user", true)
	t.Cleanup(restore)

	restoreCaptcha := zbridge.OverrideCaptchaParam("test-captcha-param")
	t.Cleanup(restoreCaptcha)

	cfg := zbridge.GetConfig()
	oldAgentMode := cfg.AgentMode
	t.Cleanup(func() { cfg.AgentMode = oldAgentMode })

	return cfg
}

// HTTPEndToEnd: webSearch=true must reach the upstream as
// auto_web_search=true with web_search pinned false (issue #42: only
// auto_web_search controls the WebSearch feature; web_search upstream is
// always false).
func TestHTTPEndToEndWebSearchPayload(t *testing.T) {
	upstream, cap := newCaptureUpstream(t)
	cfg := endToEndSearchEnv(t, upstream)
	cfg.AgentMode = false

	body := `{"model":"glm-4.7","stream":false,"webSearch":true,"messages":[{"role":"user","content":"hi"}]}`
	driveChatCompletion(t, cfg, body)

	got := cap.get(t)
	features, _ := got["features"].(map[string]interface{})
	if v, _ := features["auto_web_search"].(bool); !v {
		t.Errorf("auto_web_search = %v, want true (webSearch flag must toggle it)", features["auto_web_search"])
	}
	if v, _ := features["web_search"].(bool); v {
		t.Errorf("web_search = %v, want false (upstream always reports false; only auto_web_search toggles WebSearch)", features["web_search"])
	}
	if _, present := got["mcp_servers"]; present {
		t.Errorf("mcp_servers present (%v) on a plain web-search request — Advanced Search is separate and opt-in", got["mcp_servers"])
	}
}

// HTTPEndToEnd: advancedSearch=true must add the top-level
// mcp_servers:["advanced-search"] field captured from the real chat.z.ai
// frontend (same payload level as captcha_verify_param), and auto-enable
// the base web-search toggle.
func TestHTTPEndToEndAdvancedSearchPayload(t *testing.T) {
	upstream, cap := newCaptureUpstream(t)
	cfg := endToEndSearchEnv(t, upstream)
	cfg.AgentMode = false

	body := `{"model":"glm-4.7","stream":false,"advancedSearch":true,"messages":[{"role":"user","content":"hi"}]}`
	driveChatCompletion(t, cfg, body)

	got := cap.get(t)
	mcpRaw, present := got["mcp_servers"]
	if !present {
		t.Fatalf("mcp_servers missing from the upstream payload with advancedSearch=true; body keys: %v", keysOf(got))
	}
	mcp, ok := mcpRaw.([]interface{})
	if !ok || len(mcp) != 1 || mcp[0] != "advanced-search" {
		t.Fatalf("mcp_servers = %#v, want [\"advanced-search\"]", mcpRaw)
	}
	if _, present := got["captcha_verify_param"]; !present {
		t.Errorf("captcha_verify_param missing — mcp_servers must sit at the same payload level")
	}

	// The base toggle is auto-enabled with it (same coupling as the web UI).
	features, _ := got["features"].(map[string]interface{})
	if v, _ := features["auto_web_search"].(bool); !v {
		t.Errorf("auto_web_search = %v, want true (advancedSearch implies web search)", features["auto_web_search"])
	}
	if v, _ := features["web_search"].(bool); v {
		t.Errorf("web_search = %v, want false", features["web_search"])
	}
}

// HTTPEndToEnd: with agent mode on, the search flags are ignored — the
// upstream payload carries no auto_web_search and no mcp_servers even
// though the request asked for both. The request itself must still succeed
// (the client just doesn't get server-side search).
func TestHTTPEndToEndAgentModeGatesWebSearch(t *testing.T) {
	upstream, cap := newCaptureUpstream(t)
	cfg := endToEndSearchEnv(t, upstream)
	cfg.AgentMode = true

	body := `{"model":"glm-4.7","stream":false,"webSearch":true,"advancedSearch":true,"messages":[{"role":"user","content":"hi"}]}`
	driveChatCompletion(t, cfg, body)

	got := cap.get(t)
	features, _ := got["features"].(map[string]interface{})
	if v, present := features["auto_web_search"]; present {
		t.Errorf("auto_web_search present (%v) with agent mode on, want it stripped", v)
	}
	if v, _ := features["web_search"].(bool); v {
		t.Errorf("web_search = %v with agent mode on, want false", features["web_search"])
	}
	if _, present := got["mcp_servers"]; present {
		t.Errorf("mcp_servers present (%v) with agent mode on, want it stripped", got["mcp_servers"])
	}
}

// HTTPEndToEnd: the "search" alias must behave exactly like "webSearch".
func TestHTTPEndToEndSearchAliasPayload(t *testing.T) {
	upstream, cap := newCaptureUpstream(t)
	cfg := endToEndSearchEnv(t, upstream)
	cfg.AgentMode = false

	body := `{"model":"glm-4.7","stream":false,"search":true,"messages":[{"role":"user","content":"hi"}]}`
	driveChatCompletion(t, cfg, body)

	got := cap.get(t)
	features, _ := got["features"].(map[string]interface{})
	if v, _ := features["auto_web_search"].(bool); !v {
		t.Errorf("auto_web_search = %v, want true (the search alias must toggle it too)", features["auto_web_search"])
	}
	if v, _ := features["web_search"].(bool); v {
		t.Errorf("web_search = %v, want false", features["web_search"])
	}
}

// HTTPEndToEnd: the Anthropic /v1/messages endpoint accepts the same
// non-standard flags and they must land in the same upstream payload shape.
func TestHTTPEndToEndAnthropicSearchFlags(t *testing.T) {
	upstream, cap := newCaptureUpstream(t)
	cfg := endToEndSearchEnv(t, upstream)
	cfg.AgentMode = false

	body := `{"model":"glm-4.7","stream":false,"max_tokens":64,"webSearch":true,"advancedSearch":true,"messages":[{"role":"user","content":"hi"}]}`
	req := httptest.NewRequest("POST", "/v1/messages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", cfg.Auth.Token)
	req.Header.Set("anthropic-version", "2023-06-01")
	rec := httptest.NewRecorder()
	zbridge.NewHandler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	got := cap.get(t)
	mcpRaw, present := got["mcp_servers"]
	if !present {
		t.Fatalf("mcp_servers missing from the upstream payload for the Anthropic request; body keys: %v", keysOf(got))
	}
	mcp, ok := mcpRaw.([]interface{})
	if !ok || len(mcp) != 1 || mcp[0] != "advanced-search" {
		t.Fatalf("mcp_servers = %#v, want [\"advanced-search\"]", mcpRaw)
	}
	features, _ := got["features"].(map[string]interface{})
	if v, _ := features["auto_web_search"].(bool); !v {
		t.Errorf("auto_web_search = %v, want true", features["auto_web_search"])
	}
	if v, _ := features["web_search"].(bool); v {
		t.Errorf("web_search = %v, want false", features["web_search"])
	}
}

func keysOf(m map[string]interface{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
