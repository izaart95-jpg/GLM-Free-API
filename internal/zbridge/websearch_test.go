package zbridge

// Issue #42 regression tests: the shape of the upstream completions payload
// for web search / advanced search, and the agent-mode gate.
//
// Facts these tests pin (captured from the real chat.z.ai frontend):
//   - "web_search" inside features is ALWAYS false upstream; the only key
//     that toggles the built-in WebSearch is "auto_web_search".
//   - "Advanced Search" is an MCP server: when enabled, the web client adds
//     a top-level "mcp_servers": ["advanced-search"] field to the
//     completions payload, on the same level as captcha_verify_param.
//   - Agent mode always disables web search, even when the request asked
//     for it (the agent shim's tool contract and Z.AI's internal
//     tool-call loop cannot be mixed).

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// captureUpstreamRequest spins a mock Z.AI upstream that records the JSON
// body of every POST /api/v2/chat/completions and answers with a minimal
// SSE stream. It returns the server and a function to read the last
// captured request body.
func captureUpstreamRequest(t *testing.T) (*httptest.Server, func() map[string]interface{}) {
	t.Helper()

	var mu sync.Mutex
	var lastBody map[string]interface{}

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/chat/completions" {
			http.NotFound(w, r)
			return
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		var parsed map[string]interface{}
		if err := json.Unmarshal(raw, &parsed); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		mu.Lock()
		lastBody = parsed
		mu.Unlock()

		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		fmt.Fprintf(w, "data: {\"data\":{\"delta_content\":\"ok\"}}\n\n")
		fmt.Fprintf(w, "data: {\"data\":{\"phase\":\"done\"}}\n\n")
		fmt.Fprintf(w, "data: [DONE]\n\n")
		if flusher != nil {
			flusher.Flush()
		}
	}))

	read := func() map[string]interface{} {
		mu.Lock()
		defer mu.Unlock()
		if lastBody == nil {
			t.Fatal("upstream never received a completion request")
		}
		return lastBody
	}
	return upstream, read
}

// runCompletionThroughSendToZAI drives sendToZAI (the real feature
// resolution + payload build) against the given mock upstream and drains
// the result channel. The captcha machinery is bypassed via the
// captchaParamOverride test seam (the agent-mode cache only serves
// agent-mode requests).
func runCompletionThroughSendToZAI(t *testing.T, opts SendOptions) {
	t.Helper()

	SeedCaptchaParam("test-captcha-param")

	ch, err := sendToZAI("test prompt", opts)
	if err != nil {
		t.Fatalf("sendToZAI: %v", err)
	}
	deadline := time.After(10 * time.Second)
	for {
		select {
		case res, ok := <-ch:
			if !ok {
				return // stream finished
			}
			if res.Err != nil {
				t.Fatalf("upstream stream error: %v", res.Err)
			}
		case <-deadline:
			t.Fatal("timed out waiting for the upstream stream to finish")
		}
	}
}

// websearchTestEnv points the bridge at the mock upstream, bypasses guest
// auth, and bypasses the captcha machinery. It returns a restore function.
//
// Pacing is switched off for these tests (UPSTREAM_MIN_INTERVAL_MS=0, the
// same trick as tests/main_test.go): they talk to a local mock upstream
// and would otherwise pay the 200-500 ms inter-request gap on every call,
// slowing the package and disturbing the timing-sensitive pacing tests
// that run in parallel packages.
func websearchTestEnv(t *testing.T, upstream *httptest.Server) func() {
	t.Helper()

	t.Setenv("UPSTREAM_MIN_INTERVAL_MS", "0")

	oldBase := BASE_URL
	BASE_URL = upstream.URL
	restoreSession := OverrideSessionState("test-token", "test-user", true)
	oldCaptchaOverride := captchaParamOverride
	captchaParamOverride = "test-captcha-param"

	return func() {
		BASE_URL = oldBase
		captchaParamOverride = oldCaptchaOverride
		restoreSession()
	}
}

func websearchBoolPtr(b bool) *bool { return &b }

// webSearchTrue_leavesWebSearchFalse pins the core issue #42 discovery:
// enabling web search must set auto_web_search=true and keep web_search
// at false — sending web_search=true upstream does nothing (the server
// always mirrors it back as false) and only confuses clients.
func TestWebSearchToggleSetsAutoWebSearchOnly(t *testing.T) {
	upstream, readBody := captureUpstreamRequest(t)
	defer upstream.Close()
	restore := websearchTestEnv(t, upstream)
	defer restore()

	on := true
	runCompletionThroughSendToZAI(t, SendOptions{
		Model:     "glm-4.7",
		ChatID:    "chat-websearch-1",
		WebSearch: &on,
	})

	body := readBody()
	features, _ := body["features"].(map[string]interface{})
	if v, _ := features["auto_web_search"].(bool); !v {
		t.Errorf("auto_web_search = %v, want true (the real WebSearch toggle)", features["auto_web_search"])
	}
	if v, _ := features["web_search"].(bool); v {
		t.Errorf("web_search = %v, want false (upstream always reports false; only auto_web_search toggles WebSearch)", features["web_search"])
	}
	if _, present := body["mcp_servers"]; present {
		t.Errorf("mcp_servers present on a plain web-search request: %v — Advanced Search must be requested explicitly", body["mcp_servers"])
	}
}

// webSearchFalse_disablesAutoWebSearch pins that an explicit webSearch
// false strips auto_web_search from the features payload.
func TestWebSearchFalseRemovesAutoWebSearch(t *testing.T) {
	upstream, readBody := captureUpstreamRequest(t)
	defer upstream.Close()
	restore := websearchTestEnv(t, upstream)
	defer restore()

	off := false
	runCompletionThroughSendToZAI(t, SendOptions{
		Model:     "glm-4.7",
		ChatID:    "chat-websearch-2",
		WebSearch: &off,
	})

	body := readBody()
	features, _ := body["features"].(map[string]interface{})
	if v, present := features["auto_web_search"]; present {
		t.Errorf("auto_web_search present (%v) on a webSearch=false request, want it removed", v)
	}
	if v, _ := features["web_search"].(bool); v {
		t.Errorf("web_search = %v, want false", features["web_search"])
	}
}

// defaultRequest_noSearchKeys: with neither flag set, the payload must not
// carry auto_web_search at all and web_search stays false.
func TestDefaultRequestCarriesNoSearchFeatures(t *testing.T) {
	upstream, readBody := captureUpstreamRequest(t)
	defer upstream.Close()
	restore := websearchTestEnv(t, upstream)
	defer restore()

	runCompletionThroughSendToZAI(t, SendOptions{
		Model:  "glm-4.7",
		ChatID: "chat-websearch-3",
	})

	body := readBody()
	features, _ := body["features"].(map[string]interface{})
	if v, present := features["auto_web_search"]; present {
		t.Errorf("auto_web_search present (%v) on a default request, want it removed", v)
	}
	if v, _ := features["web_search"].(bool); v {
		t.Errorf("web_search = %v, want false", features["web_search"])
	}
	if _, present := body["mcp_servers"]; present {
		t.Errorf("mcp_servers present on a default request: %v", body["mcp_servers"])
	}
}

// advancedSearch_addsMcpServersField pins the wire shape captured from the
// real chat.z.ai web client: Advanced Search adds
// "mcp_servers": ["advanced-search"] at the TOP LEVEL of the completions
// payload (same level as captcha_verify_param), not inside features.
func TestAdvancedSearchAddsMcpServersField(t *testing.T) {
	upstream, readBody := captureUpstreamRequest(t)
	defer upstream.Close()
	restore := websearchTestEnv(t, upstream)
	defer restore()

	on := true
	runCompletionThroughSendToZAI(t, SendOptions{
		Model:          "glm-4.7",
		ChatID:         "chat-websearch-4",
		WebSearch:      &on,
		AdvancedSearch: &on,
	})

	body := readBody()

	// Top level, same level as captcha_verify_param.
	mcpRaw, present := body["mcp_servers"]
	if !present {
		t.Fatalf("mcp_servers missing from the completions payload with Advanced Search on; body keys: %v", bodyKeys(body))
	}
	mcp, ok := mcpRaw.([]interface{})
	if !ok || len(mcp) != 1 {
		t.Fatalf("mcp_servers = %#v, want [\"advanced-search\"]", mcpRaw)
	}
	if mcp[0] != "advanced-search" {
		t.Fatalf("mcp_servers = %#v, want [\"advanced-search\"]", mcpRaw)
	}
	if _, present := body["captcha_verify_param"]; !present {
		t.Errorf("captcha_verify_param missing — mcp_servers must sit on the same payload level")
	}

	// Advanced Search rides on top of the base web-search toggle.
	features, _ := body["features"].(map[string]interface{})
	if v, _ := features["auto_web_search"].(bool); !v {
		t.Errorf("auto_web_search = %v, want true when Advanced Search is on (it runs behind the web-search mode)", features["auto_web_search"])
	}
	if v, _ := features["web_search"].(bool); v {
		t.Errorf("web_search = %v, want false", features["web_search"])
	}
	if _, present := features["mcp_servers"]; present {
		t.Errorf("mcp_servers also present inside features — it is a top-level payload field, not a feature")
	}
}

// advancedSearch_alone_stillEnablesWebSearch: the handlers auto-enable the
// base web-search toggle when Advanced Search is requested (same coupling
// as the web UI, where Advanced Search hides behind the globe). This pins
// the sendToZAI side of that contract: with both flags on (what the handler
// passes down), auto_web_search is set and mcp_servers is attached.
func TestAdvancedSearchAloneStillEnablesWebSearch(t *testing.T) {
	upstream, readBody := captureUpstreamRequest(t)
	defer upstream.Close()
	restore := websearchTestEnv(t, upstream)
	defer restore()

	on := true
	runCompletionThroughSendToZAI(t, SendOptions{
		Model:          "glm-4.7",
		ChatID:         "chat-websearch-5",
		WebSearch:      &on, // handler auto-enabled (advancedSearch implies it)
		AdvancedSearch: &on,
	})

	body := readBody()
	features, _ := body["features"].(map[string]interface{})
	if v, _ := features["auto_web_search"].(bool); !v {
		t.Errorf("auto_web_search = %v, want true (Advanced Search implies web search)", features["auto_web_search"])
	}
	if _, present := body["mcp_servers"]; !present {
		t.Errorf("mcp_servers missing with Advanced Search on")
	}
}

// advancedSearch_off_leavesNoMcpServers: an explicit advancedSearch=false
// must not attach mcp_servers even when webSearch stays on.
func TestAdvancedSearchFalseLeavesNoMcpServers(t *testing.T) {
	upstream, readBody := captureUpstreamRequest(t)
	defer upstream.Close()
	restore := websearchTestEnv(t, upstream)
	defer restore()

	on := true
	off := false
	runCompletionThroughSendToZAI(t, SendOptions{
		Model:          "glm-4.7",
		ChatID:         "chat-websearch-5b",
		WebSearch:      &on,
		AdvancedSearch: &off,
	})

	body := readBody()
	if _, present := body["mcp_servers"]; present {
		t.Errorf("mcp_servers present (%v) with Advanced Search explicitly off", body["mcp_servers"])
	}
	features, _ := body["features"].(map[string]interface{})
	if v, _ := features["auto_web_search"].(bool); !v {
		t.Errorf("auto_web_search = %v, want true (webSearch on, advancedSearch off)", features["auto_web_search"])
	}
}

// agentMode_forcesWebSearchOff pins the gate from the task: with agent mode
// active, web search (and Advanced Search) are always disabled, even when
// the request explicitly asked for them. Z.AI serves WebSearch through an
// internal tool-call loop; mixing that with the agent shim's own tool
// contract makes the model emit internal markers (retrieve/open_url) as if
// they were the caller's tools.
func TestAgentModeForcesWebSearchOff(t *testing.T) {
	upstream, readBody := captureUpstreamRequest(t)
	defer upstream.Close()
	restore := websearchTestEnv(t, upstream)
	defer restore()

	oldAgentMode := config.AgentMode
	config.AgentMode = true
	defer func() { config.AgentMode = oldAgentMode }()

	on := true
	runCompletionThroughSendToZAI(t, SendOptions{
		Model:          "glm-4.7",
		ChatID:         "chat-websearch-6",
		WebSearch:      &on,
		AdvancedSearch: &on,
	})

	body := readBody()
	features, _ := body["features"].(map[string]interface{})
	if v, present := features["auto_web_search"]; present {
		t.Errorf("auto_web_search present (%v) with agent mode on, want it stripped", v)
	}
	if v, _ := features["web_search"].(bool); v {
		t.Errorf("web_search = %v with agent mode on, want false", features["web_search"])
	}
	if _, present := body["mcp_servers"]; present {
		t.Errorf("mcp_servers present (%v) with agent mode on, want it stripped", body["mcp_servers"])
	}
}

// web_search_ignoredAsStoredFeature: even if a stored per-model override
// (or server capability copied by IncludeAll) sneaks web_search=true into
// the resolved feature map, the payload build must pin it back to false —
// only auto_web_search controls WebSearch upstream.
func TestWebSearchFeatureKeyAlwaysPinnedFalse(t *testing.T) {
	upstream, readBody := captureUpstreamRequest(t)
	defer upstream.Close()
	restore := websearchTestEnv(t, upstream)
	defer restore()

	// Inject a stored override as if a user had POSTed it to /features.
	modelFeatureStatesMu.Lock()
	prev, existed := modelFeatureStates["glm-4.7"]
	modelFeatureStates["glm-4.7"] = &ModelFeatureState{
		IncludeAll: false,
		Overrides: map[string]interface{}{
			"web_search":      true, // must be pinned back to false
			"auto_web_search": true,
		},
	}
	modelFeatureStatesMu.Unlock()
	defer func() {
		modelFeatureStatesMu.Lock()
		if existed {
			modelFeatureStates["glm-4.7"] = prev
		} else {
			delete(modelFeatureStates, "glm-4.7")
		}
		modelFeatureStatesMu.Unlock()
	}()

	runCompletionThroughSendToZAI(t, SendOptions{
		Model:  "glm-4.7",
		ChatID: "chat-websearch-7",
	})

	body := readBody()
	features, _ := body["features"].(map[string]interface{})
	if v, _ := features["web_search"].(bool); v {
		t.Errorf("web_search = %v despite the stored override, want false (pinned by the payload build)", features["web_search"])
	}
	if v, _ := features["auto_web_search"].(bool); !v {
		t.Errorf("auto_web_search = %v, want true (stored override is respected)", features["auto_web_search"])
	}
}

func bodyKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// agentMode_disabledRequest keeps working: the gate must not break the
// normal no-search path when agent mode is on.
func TestAgentModePlainRequestStillClean(t *testing.T) {
	upstream, readBody := captureUpstreamRequest(t)
	defer upstream.Close()
	restore := websearchTestEnv(t, upstream)
	defer restore()

	oldAgentMode := config.AgentMode
	config.AgentMode = true
	defer func() { config.AgentMode = oldAgentMode }()

	runCompletionThroughSendToZAI(t, SendOptions{
		Model:  "glm-4.7",
		ChatID: "chat-websearch-8",
	})

	body := readBody()
	features, _ := body["features"].(map[string]interface{})
	if v, present := features["auto_web_search"]; present {
		t.Errorf("auto_web_search present (%v) on an agent-mode request without search, want stripped", v)
	}
	if v, _ := features["web_search"].(bool); v {
		t.Errorf("web_search = %v, want false", features["web_search"])
	}
	if _, present := body["mcp_servers"]; present {
		t.Errorf("mcp_servers present on a request without Advanced Search: %v", body["mcp_servers"])
	}
}
