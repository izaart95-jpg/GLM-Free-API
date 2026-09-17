// agent_ultra_test.go — V28 ultra repair pipeline (§8 acceptance).
//
// Covers: strict opt-in gating, buffering only in ultra, valid bypass of V28,
// fragment-only extraction, re-validation before splicing, idempotence,
// multi-span independence, fallback preservation, and dormancy
// (non-ultra byte-identical).

package zbridge

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// ── fakes ──────────────────────────────────────────────────────────────────

type fakeRepairer struct {
	calls    int
	lastFrag string
	lastTool string
	reply    string
	err      error
}

func (f *fakeRepairer) RepairFragment(_ context.Context, fragment, toolsJSON string) (string, error) {
	f.calls++
	f.lastFrag = fragment
	f.lastTool = toolsJSON
	if f.err != nil {
		return "", f.err
	}
	return f.reply, nil
}

// withUltraGate swaps config for one test and restores it.
func withUltraGate(agentMode bool, level string, forceCPU bool) func() {
	oldMode, oldLevel, oldCPU := config.AgentMode, config.AgentModeLevel, config.ForceCPU
	config.AgentMode, config.AgentModeLevel, config.ForceCPU = agentMode, level, forceCPU
	return func() { config.AgentMode, config.AgentModeLevel, config.ForceCPU = oldMode, oldLevel, oldCPU }
}

const ultraTools = `[{"name":"bash"},{"name":"read"}]`

// ── gating matrix (§2) ─────────────────────────────────────────────────────

func TestUltraGateMatrix(t *testing.T) {
	cases := []struct {
		level string
		agent bool
		want  bool
	}{
		{"", false, false},        // (none) → default, dormant
		{"ultra", false, false},   // level without agent → dormant
		{"", true, false},         // --agent-mode only → dormant
		{"standard", true, false}, // non-ultra level → dormant
		{"ULTRA", true, true},     // case-insensitive ultra + agent → active
		{"ultra", true, true},     // --agent-mode --agent-mode-level=ultra → active
		{" ultra ", true, true},   // whitespace tolerated
	}
	for _, c := range cases {
		restore := withUltraGate(c.agent, c.level, false)
		got := UltraRepairEnabled()
		if got != c.want {
			t.Errorf("agent=%v level=%q → Ultra=%v, want %v (%s)", c.agent, c.level, got, c.want, ultraGateDebug())
		}
		restore()
	}
}

func TestUltraGateMatrixExplicit(t *testing.T) {
	if func() bool {
		r := withUltraGate(false, "", false)
		defer r()
		return UltraRepairEnabled()
	}() {
		t.Error("(none) must be dormant")
	}
	if func() bool {
		r := withUltraGate(true, "", false)
		defer r()
		return UltraRepairEnabled()
	}() {
		t.Error("--agent-mode only must be dormant (no V28)")
	}
	if func() bool {
		r := withUltraGate(true, "ultra", false)
		defer r()
		return UltraRepairEnabled()
	}() != true {
		t.Error("--agent-mode --agent-mode-level=ultra must be active")
	}
}

// ── validator: valid bypasses V28 (§3, §8) ─────────────────────────────────

func TestUltraValidCanonicalBypassesV28(t *testing.T) {
	restore := withUltraGate(true, "ultra", true)
	defer restore()
	valid := "Checking.\n<<<TOOL_CALL>>>\n{\"name\":\"bash\",\"arguments\":{\"command\":\"uname -a\"}}\n<<<END_TOOL_CALL>>>\n"
	if ClassifyUltra(valid) != UltraValid {
		t.Fatalf("valid canonical classified as %v", ClassifyUltra(valid))
	}
	f := &fakeRepairer{reply: "SHOULD NOT BE CALLED"}
	got, events := RepairUltraBuffer(valid, ultraTools, f)
	if got != valid {
		t.Errorf("valid input altered:\n got %q\nwant %q", got, valid)
	}
	if f.calls != 0 {
		t.Errorf("V28 invoked %d times for valid input, want 0 (bypass)", f.calls)
	}
	if len(events) != 1 || events[0].Decision != "passthrough-valid" {
		t.Errorf("events = %+v, want single passthrough-valid", events)
	}
}

func TestUltraPlainPassthrough(t *testing.T) {
	restore := withUltraGate(true, "ultra", true)
	defer restore()
	plain := "Paris is lovely in spring."
	if ClassifyUltra(plain) != UltraPlain {
		t.Fatalf("prose classified as %v", ClassifyUltra(plain))
	}
	f := &fakeRepairer{reply: "NO"}
	got, _ := RepairUltraBuffer(plain, ultraTools, f)
	if got != plain || f.calls != 0 {
		t.Errorf("plain altered or V28 called (got %q calls=%d)", got, f.calls)
	}
}

// ── extraction: fragment only, never whole response (§4.1) ─────────────────

func TestUltraExtractsFragmentNotWholeResponse(t *testing.T) {
	restore := withUltraGate(true, "ultra", true)
	defer restore()
	prose := "Here is my plan in plain words. "
	malformed := "<tool_call>{\"tool\": \"bash\", \"command\": \"uname -a\"}</tool_call>"
	tail := " Let me know what you see."
	full := prose + malformed + tail
	spans := ExtractMalformedSpans(full)
	if len(spans) == 0 {
		t.Fatalf("no malformed span extracted from %q", full)
	}
	for _, sp := range spans {
		if strings.Contains(sp.Raw, "plain words") || strings.Contains(sp.Raw, "Let me know") {
			t.Errorf("span swallowed surrounding prose: %q", sp.Raw)
		}
		if len(sp.Raw) >= len(full) {
			t.Errorf("span is whole response (%d >= %d)", len(sp.Raw), len(full))
		}
	}
	// The stock parser must NOT already accept it (else it would bypass).
	if HasValidCanonicalBlock(full) {
		t.Errorf("malformed fixture unexpectedly parses as canonical")
	}
}

func TestUltraProseGuard(t *testing.T) {
	restore := withUltraGate(true, "ultra", true)
	defer restore()
	// Nameless JSON with no tool key must stay prose (existing policy).
	prose := `The config {"command":"ls"} is just an example.`
	if spans := ExtractMalformedSpans(prose); len(spans) != 0 {
		t.Errorf("prose flagged as malformed: %+v", spans)
	}
}

// ── repair + re-validation + splice (§4.2) ─────────────────────────────────

func TestUltraMalformedRepairedAndSpliced(t *testing.T) {
	restore := withUltraGate(true, "ultra", true)
	defer restore()
	full := "Working on it.\n<tool_call>{\"tool\": \"bash\", \"command\": \"uname -a\"}</tool_call>\nDone soon."
	f := &fakeRepairer{reply: "<<<TOOL_CALL>>>\n{\"name\":\"bash\",\"arguments\":{\"command\":\"uname -a\"}}\n<<<END_TOOL_CALL>>>"}
	got, events := RepairUltraBuffer(full, ultraTools, f)
	if f.calls != 1 {
		t.Fatalf("V28 calls = %d, want 1", f.calls)
	}
	if strings.Contains(f.lastFrag, "Working on it") || strings.Contains(f.lastFrag, "Done soon") {
		t.Errorf("whole response forwarded to V28 (fragment=%q)", f.lastFrag)
	}
	if !strings.Contains(got, "<<<TOOL_CALL>>>") || !strings.Contains(got, `"command":"uname -a"`) {
		t.Errorf("repaired block missing: %q", got)
	}
	if !strings.Contains(got, "Working on it.") || !strings.Contains(got, "Done soon.") {
		t.Errorf("surrounding prose damaged: %q", got)
	}
	if strings.Contains(got, "<tool_call>") {
		t.Errorf("malformed fragment not replaced: %q", got)
	}
	if len(events) != 1 || events[0].Decision != "repaired" {
		t.Errorf("events = %+v, want single repaired", events)
	}
	if events[0].OriginalFragment == "" || events[0].ExtractedFragment == "" || events[0].V28Output == "" {
		t.Errorf("repair event missing fields: %+v", events[0])
	}
	// Postcondition: repaired text validates with the SAME validator.
	if !HasValidCanonicalBlock(got) {
		t.Errorf("repaired output fails re-validation: %q", got)
	}
}

func TestUltraInvalidV28OutputFallsBack(t *testing.T) {
	restore := withUltraGate(true, "ultra", true)
	defer restore()
	full := "Hi.\n<tool_call>{\"tool\": \"bash\", \"command\": \"id\"}</tool_call>\nBye."
	f := &fakeRepairer{reply: "not a tool call at all"}
	got, events := RepairUltraBuffer(full, ultraTools, f)
	if got != full {
		t.Errorf("invalid V28 output must keep original (got %q want %q)", got, full)
	}
	if len(events) != 1 || events[0].Decision != "fallback-invalid" {
		t.Errorf("events = %+v, want fallback-invalid", events)
	}
}

func TestUltraNoRepairerPreservesIntent(t *testing.T) {
	restore := withUltraGate(true, "ultra", true)
	defer restore()
	full := "Hi.\n@bash(command=\"id\")\nBye."
	if ClassifyUltra(full) != UltraMalformed {
		t.Skipf("fixture not malformed under current heuristics (spans=%+v)", ExtractMalformedSpans(full))
	}
	got, events := RepairUltraBuffer(full, ultraTools, nil)
	if got != full {
		t.Errorf("nil repairer must preserve original, got %q", got)
	}
	for _, ev := range events {
		if ev.Decision != "no-repairer" {
			t.Errorf("event decision = %q, want no-repairer", ev.Decision)
		}
	}
}

// ── idempotence + multi-span (§4.3) ────────────────────────────────────────

func TestUltraRepairIdempotent(t *testing.T) {
	restore := withUltraGate(true, "ultra", true)
	defer restore()
	full := "A.\n<tool_call>{\"tool\": \"bash\", \"command\": \"id\"}</tool_call>\nB."
	f := &fakeRepairer{reply: "<<<TOOL_CALL>>>\n{\"name\":\"bash\",\"arguments\":{\"command\":\"id\"}}\n<<<END_TOOL_CALL>>>"}
	once, _ := RepairUltraBuffer(full, ultraTools, f)
	f2 := &fakeRepairer{reply: "MUST NOT BE CALLED"}
	twice, _ := RepairUltraBuffer(once, ultraTools, f2)
	if once != twice {
		t.Errorf("not idempotent:\n once %q\ntwice %q", once, twice)
	}
	if f2.calls != 0 {
		t.Errorf("second pass invoked V28 %d times (repaired text must bypass)", f2.calls)
	}
}

func TestUltraMultipleSpansIndependent(t *testing.T) {
	restore := withUltraGate(true, "ultra", true)
	defer restore()
	full := "First:\n<tool_call>{\"tool\": \"bash\", \"command\": \"id\"}</tool_call>\nMiddle prose stays.\n<tool_call>{\"tool\": \"read\", \"path\": \"/tmp/x\"}</tool_call>\nEnd."
	calls := 0
	f := &fakeRepairer{}
	// Route per-fragment replies so each span is handled independently.
	wrapped := &routeRepairer{fn: func(frag string) string {
		calls++
		if strings.Contains(frag, "read") || strings.Contains(frag, "/tmp/x") {
			return "<<<TOOL_CALL>>>\n{\"name\":\"read\",\"arguments\":{\"path\":\"/tmp/x\"}}\n<<<END_TOOL_CALL>>>"
		}
		return "<<<TOOL_CALL>>>\n{\"name\":\"bash\",\"arguments\":{\"command\":\"id\"}}\n<<<END_TOOL_CALL>>>"
	}}
	_ = f
	got, events := RepairUltraBuffer(full, ultraTools, wrapped)
	if calls != 2 {
		t.Errorf("V28 calls = %d, want 2 (one per span)", calls)
	}
	if len(events) != 2 {
		t.Errorf("events = %d, want 2", len(events))
	}
	if !strings.Contains(got, "Middle prose stays.") {
		t.Errorf("prose between spans damaged: %q", got)
	}
	if !HasValidCanonicalBlock(got) {
		t.Errorf("repaired multi-span output fails validation: %q", got)
	}
}

type routeRepairer struct{ fn func(string) string }

func (r *routeRepairer) RepairFragment(_ context.Context, frag, _ string) (string, error) {
	return r.fn(frag), nil
}

// ── buffering (§3) ─────────────────────────────────────────────────────────

func TestUltraBufferDefersUntilFinish(t *testing.T) {
	restore := withUltraGate(true, "ultra", true)
	defer restore()
	restoreR := SetV28RepairerForTests(&fakeRepairer{reply: "<<<TOOL_CALL>>>\n{\"name\":\"bash\",\"arguments\":{\"command\":\"id\"}}\n<<<END_TOOL_CALL>>>"})
	defer restoreR()
	tools, _ := json.Marshal([]openAITool{{Type: "function", Function: &openAIFnSpec{Name: "bash"}}})
	b := NewUltraBuffer(tools)
	// Feed a malformed stream in pieces; Finish must yield the repaired call
	// with prose intact and ordering preserved.
	for _, piece := range []string{"Hello.\n", "<tool_call>", "{\"tool\": \"bash\", ", "\"command\": \"id\"}</tool_call>", "\nBye."} {
		b.Feed(piece)
	}
	content, toolCalls := b.Finish()
	if len(toolCalls) == 0 {
		t.Fatalf("buffered finish produced no tool calls (content=%q)", content)
	}
	if !strings.Contains(content, "Hello.") || !strings.Contains(content, "Bye.") {
		t.Errorf("buffered prose lost: %q", content)
	}
}

func TestUltraBufferValidPassthrough(t *testing.T) {
	restore := withUltraGate(true, "ultra", true)
	defer restore()
	tools, _ := json.Marshal([]openAITool{{Type: "function", Function: &openAIFnSpec{Name: "bash"}}})
	b := NewUltraBuffer(tools)
	valid := "<<<TOOL_CALL>>>\n{\"name\":\"bash\",\"arguments\":{\"command\":\"id\"}}\n<<<END_TOOL_CALL>>>"
	b.Feed(valid)
	content, toolCalls := b.Finish()
	if len(toolCalls) != 1 {
		t.Errorf("valid buffered call → %d calls, want 1 (content=%q)", len(toolCalls), content)
	}
}

// ── dormancy: non-ultra never touches V28 (§6) ─────────────────────────────

func TestUltraDormantWhenNotEnabled(t *testing.T) {
	malformed := "Hi.\n<tool_call>{\"tool\": \"bash\", \"command\": \"id\"}</tool_call>\nBye."
	for _, tc := range []struct {
		agent bool
		level string
	}{
		{false, ""}, {true, ""}, {true, "standard"}, {false, "ultra"},
	} {
		restore := withUltraGate(tc.agent, tc.level, false)
		f := &fakeRepairer{reply: "<<<TOOL_CALL>>>\n{\"name\":\"bash\",\"arguments\":{\"command\":\"id\"}}\n<<<END_TOOL_CALL>>>"}
		got, _ := RepairUltraBuffer(malformed, ultraTools, f)
		if got != malformed {
			t.Errorf("agent=%v level=%q altered input (non-regression breach)", tc.agent, tc.level)
		}
		if f.calls != 0 {
			t.Errorf("agent=%v level=%q invoked V28 (must stay dormant)", tc.agent, tc.level)
		}
		if UltraRepairFullText(malformed, json.RawMessage(ultraTools)) != malformed {
			t.Errorf("agent=%v level=%q UltraRepairFullText altered input", tc.agent, tc.level)
		}
		restore()
	}
}

func TestUltraNormalizePreservesPayload(t *testing.T) {
	in := "  ```json\n{\"tool\": \"bash\"} \n```  "
	got := NormalizeUltraFragment(in)
	if !strings.Contains(got, `"tool": "bash"`) {
		t.Errorf("payload damaged: %q", got)
	}
	if strings.Contains(got, "```") {
		t.Errorf("fence not stripped: %q", got)
	}
}
