// agent_ultra.go
//
// ULTRA MODE — V28 intelligent malformed-tool repair layer.
//
// Strict opt-in: EVERY entry point in this file first checks
// UltraRepairEnabled() (--agent-mode + --agent-mode-level=ultra).
// When the gate is off the helpers return their input untouched and never
// touch V28, so default and plain --agent-mode behaviour stay byte-for-byte
// identical to upstream (see §6 non-regression).
//
// Ultra pipeline (§3–§4), active only under the gate:
//   1. Buffer the full upstream text (UltraBuffer) instead of forwarding
//     incrementally.
//  2. Validate against the canonical format from agent.go
//     (findAgentSpans + agentLooseParse). Valid canonical → passthrough.
//  3. Otherwise extract ONLY the malformed tool-call span(s) — never the
//     whole response — via intent heuristics, normalize structural noise,
//     variablize long values (varpass mirror), and send the fragment +
//     canonical spec to V28 for ONE canonical object back.
//  4. Re-validate V28 output with the same validator; splice on success,
//     else forward the original fragment with a warning (never drop intent).
//
// Repair is idempotent (valid canonical input is returned unchanged) and
// non-destructive (surrounding prose is preserved byte-for-byte; each
// malformed span is handled independently). Every repair event is logged
// with original / extracted / V28 output / decision.

package zbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

// ── gating ─────────────────────────────────────────────────────────────────

// ultraGateDebug reports the current gate state for logs/tests.
func ultraGateDebug() string {
	return "agentMode=" + boolStr(config.AgentMode) + " level=" + config.AgentModeLevel
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// ── canonical validator (§3.2, reuses agent.go) ─────────────────────────────

// UltraClassification is the validator verdict on buffered text.
type UltraClassification int

const (
	// UltraPlain: no tool intent at all — passthrough as content.
	UltraPlain UltraClassification = iota
	// UltraValid: ≥1 valid canonical block and no malformed span — bypass V28.
	UltraValid
	// UltraMalformed: ≥1 malformed tool-call span — enter repair path.
	UltraMalformed
)

// HasValidCanonicalBlock reports whether text holds ≥1 parseable canonical
// block (the same tolerant parse the runtime uses to emit tool_calls).
func HasValidCanonicalBlock(text string) bool {
	return len(ParseAgentToolCalls(NormalizeAgentFences(text))) > 0
}

// ClassifyUltra validates buffered content per §3.2.
func ClassifyUltra(text string) UltraClassification {
	if HasValidCanonicalBlock(text) {
		// A response can carry BOTH a valid block and a malformed sibling
		// (multi-call responses). Those still need the repair path for the
		// malformed span — valid blocks are left untouched by the splicer.
		if len(ExtractMalformedSpans(text)) > 0 {
			return UltraMalformed
		}
		return UltraValid
	}
	if len(ExtractMalformedSpans(text)) > 0 {
		return UltraMalformed
	}
	return UltraPlain
}

// ── malformed-span extraction (§4.1) ───────────────────────────────────────

// UltraSpan is one isolated malformed tool-call region: [Start,End) byte
// offsets into the original buffered text plus the normalized fragment sent
// toward V28.
type UltraSpan struct {
	Start, End int
	Fragment   string // normalized, variablized skeleton actually sent to V28
	Raw        string // original bytes at [Start,End)
	VarMap     map[string]string
}

// Intent heuristics: each pattern names a malformed-tool shape observed in
// the wild or in V28 training (v28 GENERALIZE families + issue #44
// derailments). Canonical <<<TOOL_CALL>>> blocks are matched separately and
// explicitly EXCLUDED below, so valid calls always bypass V28 (§8).
var ultraIntentPatterns = []*regexp.Regexp{
	// Renamed markers / envelopes.
	regexp.MustCompile(`(?i)TOOL_CALL_BLOCK`),
	regexp.MustCompile(`(?is)<\s*/?\s*tool_call\s*[^>]*>.*?<\s*/\s*tool_call\s*>`),
	regexp.MustCompile(`(?is)<\s*tool_args\s*>.*?<\s*/\s*tool_args\s*>`),
	regexp.MustCompile(`(?i)<\s*invoke\b[^>]*>.*?<\s*/\s*invoke\s*>`),
	regexp.MustCompile(`DSML`),
	// V28 GENERALIZE wrapper families (train): xml_upper, sexpr, at_call,
	// shell_export, hash_flag, double_brace, json_toolkey, ini_section,
	// md_tool — plus held-out pipe_kv / csv_row / angle_json.
	regexp.MustCompile(`(?is)<TOOL><NAME>.*?</ARGS></TOOL>`),
	regexp.MustCompile(`(?s)\(tool\s+\w+.*?\)`),
	regexp.MustCompile(`@\w+\([^)]*=\s*("[^"]*"|\{[^\n]{0,1000})`),
	regexp.MustCompile(`(?m)^\s*TOOL=\w+.*`),
	regexp.MustCompile(`(?m)^\s*#TOOL\s+\w+.*`),
	regexp.MustCompile(`(?s)\{\{\w+\s*:.*?\}\}`),
	regexp.MustCompile(`(?s)\{"tool"\s*:.*?\}`),
	regexp.MustCompile(`(?m)^\s*\[[\w-]+\]\s*$`),
	regexp.MustCompile("(?is)```tool.*?```"),
	regexp.MustCompile(`(?s)<json>\s*\{.*?\}\s*</json>`),
	regexp.MustCompile(`(?m)^\s*\w+\s*::.*?=.*`),
	regexp.MustCompile(`(?m)^\s*tool,args:.*`),
	// Bare JSON-ish payloads naming a tool outside canonical markers:
	// {"tool": "bash", ...} / {"name": "bash", "command": ...} fragments.
	// Kept deliberately narrow (requires a tool-ish key + a brace) so prose
	// like `{"command":"ls"}` with no tool key is NOT flagged (see test).
	regexp.MustCompile(`(?s)\{[^{}]{0,1000}"(tool|tool_name|function|function_name)"\s*:[^{}]{0,1000}\}`),
	// Pseudo-calls: get_weather(city="Tokyo") / bash(command="...").
	regexp.MustCompile(`(?s)\b[a-zA-Z_][a-zA-Z0-9_]*\s*\(\s*[a-zA-Z_][a-zA-Z0-9_]*\s*=\s*("[^"]*"|'[^']*'|\{[^\n]{0,1000})`),
}

// ultraProtectedSpans returns the byte ranges of valid canonical blocks that
// must never be treated as malformed (bypass-V28 guarantee).
func ultraProtectedSpans(text string) [][2]int {
	norm := NormalizeAgentFences(text)
	// Map spans on the normalized text back conservatively: since fence
	// stripping only removes bytes, any span found on raw text via
	// findAgentSpans is authoritative for splicing. Use raw spans whose body
	// parses.
	var out [][2]int
	for _, sp := range findAgentSpans(text) {
		if _, _, ok := agentLooseParse(text[sp.bodyStart:sp.bodyEnd]); ok {
			out = append(out, [2]int{sp.start, sp.end})
		}
	}
	_ = norm
	return out
}

func overlapsProtected(start, end int, prot [][2]int) bool {
	for _, p := range prot {
		if start < p[1] && end > p[0] {
			return true
		}
	}
	return false
}

// NormalizeUltraFragment removes obvious structural noise BEFORE V28 (§4.1):
// surrounding whitespace, fences adjacent to the fragment, stray trailing
// commas/semicolons, and fullwidth lookalikes. The semantic payload
// (tool name + arguments) is preserved verbatim — only the skeleton is fixed.
func NormalizeUltraFragment(s string) string {
	t := strings.TrimSpace(s)
	// Strip fences hugging the fragment (never ordinary code elsewhere).
	t = strings.TrimSpace(agentFenceLead.ReplaceAllString(t, ""))
	t = strings.TrimSpace(agentFenceTail.ReplaceAllString(t, ""))
	// Fullwidth fold (mirror varpass.ascii_fold for ；→; etc.).
	t = strings.Map(func(r rune) rune {
		if r >= 0xFF01 && r <= 0xFF5E {
			return r - 0xFEE0
		}
		if r == 0x3000 {
			return ' '
		}
		return r
	}, t)
	t = strings.Trim(t, " \t\n\r;,")
	return strings.TrimSpace(t)
}

// ExtractMalformedSpans isolates malformed tool-call regions (§4.1).
// It returns spans sorted by offset, non-overlapping, excluding surrounding
// prose and excluding valid canonical blocks. Pure prose yields nil.
//
// The regex match itself IS the span — no line expansion — so surrounding
// prose is never swallowed (fragment-only guarantee). Patterns that need
// line context (TOOL=, #TOOL, ::, csv) already carry `.*` to end-of-line.
func ExtractMalformedSpans(text string) []UltraSpan {
	prot := ultraProtectedSpans(text)
	var cands []candSpan
	for _, re := range ultraIntentPatterns {
		for _, m := range re.FindAllStringIndex(text, -1) {
			s, e := m[0], m[1]
			if overlapsProtected(s, e, prot) {
				continue
			}
			cands = append(cands, candSpan{s, e})
		}
	}
	if len(cands) == 0 {
		return nil
	}
	// Merge overlaps, keep outermost.
	merged := mergeCands(cands)
	var out []UltraSpan
	for _, m := range merged {
		if overlapsProtected(m.s, m.e, prot) {
			continue
		}
		raw := text[m.s:m.e]
		norm := NormalizeUltraFragment(raw)
		if strings.TrimSpace(norm) == "" {
			continue
		}
		// Drop candidates with no tool intent (no tool-ish key, no
		// paren-call, no marker word) — pure prose protection.
		if !ultraHasIntent(norm) {
			continue
		}
		skel, vmap := ultraVariablize(norm)
		out = append(out, UltraSpan{Start: m.s, End: m.e, Fragment: skel, Raw: raw, VarMap: vmap})
	}
	return out
}

type candSpan struct{ s, e int }

func mergeCands(in []candSpan) []candSpan {
	// Insertion sort by start (n is tiny).
	for i := 1; i < len(in); i++ {
		for j := i; j > 0 && in[j].s < in[j-1].s; j-- {
			in[j], in[j-1] = in[j-1], in[j]
		}
	}
	var out []candSpan
	for _, c := range in {
		if len(out) == 0 || c.s > out[len(out)-1].e {
			out = append(out, c)
			continue
		}
		if c.e > out[len(out)-1].e {
			out[len(out)-1].e = c.e
		}
	}
	return out
}

// ultraHasIntent is the prose guard: a span is tool-intent only if it carries
// a tool-ish key, a call shape, or a marker word.
// NOTE: built with a double-quoted string (not backticks) because the
// pattern itself contains ``` fence literals.
var ultraIntentProbe = regexp.MustCompile("(?i)(tool|invoke|dsml|TOOL_CALL|<tool|@\\w+\\(|\\(tool\\s|\\{\\{|\\[[\\w-]+\\]|```tool|<json>|::|tool,args:|\"[a-z_]*name\"\\s*:|bash|read|write|edit|grep|get_weather|AskUser)")

func ultraHasIntent(s string) bool {
	return ultraIntentProbe.MatchString(s)
}

// ── varpass mirror (LONG sentinels) ────────────────────────────────────────
// V28 was trained with varpass delexicalization: long values (≥150 chars, or
// ≥50 with a newline) are replaced by PUA sentinels before the model and
// substituted back after. Mirroring that here keeps fragments small and the
// payload exact. Only the two highest-value grammars are ported (fenced
// blocks + long quoted literals); the rest ride through verbatim.

const ultraLongThreshold = 150

func ultraSent(i int) string { return "VAR" + itoa(i) + "" }

func itoa(i int) string { return fmt.Sprintf("%d", i) }

var (
	ultraFenceRe = regexp.MustCompile("(?s)```[a-zA-Z]*\n(.*?)```")
	ultraLongQRe = regexp.MustCompile(`"((?:[^"\\]|\\.){200,})"`)
)

func ultraLongEnough(body string) bool {
	return len(body) >= ultraLongThreshold || (len(body) >= 50 && strings.Contains(body, "\n"))
}

// ultraVariablize replaces long spans with collision-proof PUA sentinels.
// Returns (skeleton, varmap). Deterministic; longest-first via offset order.
func ultraVariablize(text string) (string, map[string]string) {
	type span struct{ s, e int }
	var spans []span
	for _, re := range []*regexp.Regexp{ultraFenceRe, ultraLongQRe} {
		for _, m := range re.FindAllStringSubmatchIndex(text, -1) {
			if len(m) < 4 || m[2] < 0 {
				continue
			}
			body := text[m[2]:m[3]]
			if ultraLongEnough(body) {
				spans = append(spans, span{m[2], m[3]})
			}
		}
	}
	if len(spans) == 0 {
		return text, nil
	}
	for i := 1; i < len(spans); i++ {
		for j := i; j > 0 && spans[j].s < spans[j-1].s; j-- {
			spans[j], spans[j-1] = spans[j-1], spans[j]
		}
	}
	var out strings.Builder
	varmap := map[string]string{}
	pos := 0
	last := -1
	for _, sp := range spans {
		if sp.s < last {
			continue
		}
		last = sp.e
		out.WriteString(text[pos:sp.s])
		sent := ultraSent(len(varmap))
		varmap[sent] = text[sp.s:sp.e]
		out.WriteString(sent)
		pos = sp.e
	}
	out.WriteString(text[pos:])
	return out.String(), varmap
}

// ultraRehydrate substitutes sentinels back with exact original values.
func ultraRehydrate(s string, varmap map[string]string) string {
	for sent, val := range varmap {
		s = strings.ReplaceAll(s, sent, val)
	}
	return s
}

// ── repair events + logging (§4.3) ─────────────────────────────────────────

type UltraRepairEvent struct {
	OriginalFragment  string
	ExtractedFragment string
	V28Output         string
	Decision          string // "repaired" | "passthrough-valid" | "fallback-passthrough" | "fallback-invalid" | "no-repairer" | "plain"
}

func logUltraEvent(ev UltraRepairEvent) {
	const max = 500
	tr := func(s string) string {
		s = strings.ReplaceAll(s, "\n", "\\n")
		if len(s) > max {
			return s[:max] + "…"
		}
		return s
	}
	log.Printf("[UltraRepair] decision=%s original=%q extracted=%q v28=%q",
		ev.Decision, tr(ev.OriginalFragment), tr(ev.ExtractedFragment), tr(ev.V28Output))
}

// ── core repair (§4.2) ─────────────────────────────────────────────────────

// RepairUltraBuffer is the full §3–§4 pipeline on buffered text.
// toolsJSON is the request's tools array (for the repair prompt context).
// It is idempotent: valid canonical / plain input returns byte-identical.
// Malformed spans are each repaired independently; V28 output is re-validated
// before splicing, else the original fragment is kept with a warning log.
// Never drops tool intent silently. Dormant (gate off) → input unchanged.
func RepairUltraBuffer(fullText, toolsJSON string, repairer V28Repairer) (string, []UltraRepairEvent) {
	if !UltraRepairEnabled() {
		return fullText, nil
	}
	class := ClassifyUltra(fullText)
	if class != UltraMalformed {
		return fullText, []UltraRepairEvent{{OriginalFragment: "", ExtractedFragment: "", V28Output: "", Decision: map[bool]string{true: "passthrough-valid", false: "plain"}[class == UltraValid]}}
	}
	spans := ExtractMalformedSpans(fullText)
	if len(spans) == 0 {
		return fullText, nil
	}
	if repairer == nil {
		repairer = GetV28Repairer()
	}
	out := fullText
	// Splice back-to-front so earlier offsets stay valid.
	var events []UltraRepairEvent
	for i := len(spans) - 1; i >= 0; i-- {
		sp := spans[i]
		// Recompute offsets against the evolving `out`: spans were computed
		// on the original; since we splice back-to-front, [Start,End) is
		// still valid for not-yet-spliced earlier regions. Later splices
		// only touched bytes after End.
		original := sp.Raw
		fragment := sp.Fragment
		if repairer == nil {
			ev := UltraRepairEvent{OriginalFragment: original, ExtractedFragment: fragment, V28Output: "", Decision: "no-repairer"}
			logUltraEvent(ev)
			events = append(events, ev)
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		v28out, err := repairer.RepairFragment(ctx, fragment, toolsJSON)
		cancel()
		if err != nil {
			ev := UltraRepairEvent{OriginalFragment: original, ExtractedFragment: fragment, V28Output: "ERROR: " + err.Error(), Decision: "fallback-passthrough"}
			log.Printf("[UltraRepair] WARNING: V28 repair failed (%v); forwarding original malformed fragment (intent preserved)", err)
			logUltraEvent(ev)
			events = append(events, ev)
			continue
		}
		rehyd := ultraRehydrate(strings.TrimSpace(v28out), sp.VarMap)
		// Re-validate with the SAME validator (§4.2).
		calls := ParseAgentToolCalls(rehyd)
		if len(calls) == 0 {
			ev := UltraRepairEvent{OriginalFragment: original, ExtractedFragment: fragment, V28Output: v28out, Decision: "fallback-invalid"}
			log.Printf("[UltraRepair] WARNING: V28 output failed re-validation; forwarding original malformed fragment (intent preserved)")
			logUltraEvent(ev)
			events = append(events, ev)
			continue
		}
		// Splice the FIRST repaired canonical block in place of the
		// malformed fragment (non-destructive to surrounding prose).
		block := firstCanonicalBlock(rehyd)
		if block == "" {
			block = rehyd
		}
		out = out[:sp.Start] + block + out[sp.End:]
		ev := UltraRepairEvent{OriginalFragment: original, ExtractedFragment: fragment, V28Output: v28out, Decision: "repaired"}
		logUltraEvent(ev)
		events = append(events, ev)
	}
	// Reverse events back to document order.
	for i, j := 0, len(events)-1; i < j; i, j = i+1, j-1 {
		events[i], events[j] = events[j], events[i]
	}
	return out, events
}

func firstCanonicalBlock(text string) string {
	for _, sp := range findAgentSpans(text) {
		body := text[sp.bodyStart:sp.bodyEnd]
		if _, _, ok := agentLooseParse(body); ok {
			return text[sp.start:sp.end]
		}
	}
	return ""
}

// UltraRepairFullText is the non-streaming convenience wrapper: repairs the
// finished assistant text before agentExtractToolCalls. Dormant → input.
func UltraRepairFullText(fullContent string, toolsRaw json.RawMessage) string {
	if !UltraRepairEnabled() {
		return fullContent
	}
	repaired, _ := RepairUltraBuffer(fullContent, string(toolsRaw), nil)
	return repaired
}

// ── streaming buffer (§3.1, §3.4) ──────────────────────────────────────────

// UltraBuffer accumulates upstream chunks in ultra mode instead of forwarding
// them immediately. At Finish it runs the validator + repair pipeline once,
// then replays the (possibly repaired) text through a fresh stock
// interceptor so chunk ordering, call-index sequencing, and the final
// terminator behave exactly like the non-ultra path from that point on.
type UltraBuffer struct {
	chunks   []string
	full     strings.Builder
	toolsRaw json.RawMessage
}

// NewUltraBuffer starts a per-request ultra buffer. toolsRaw is captured for
// repair-prompt context. Dormant callers must not construct one (handlers
// check UltraRepairEnabled() first).
func NewUltraBuffer(toolsRaw json.RawMessage) *UltraBuffer {
	return &UltraBuffer{toolsRaw: toolsRaw}
}

// Feed appends one upstream delta. Nothing is forwarded yet (buffering).
func (b *UltraBuffer) Feed(chunk string) {
	if chunk == "" {
		return
	}
	b.chunks = append(b.chunks, chunk)
	b.full.WriteString(chunk)
}

// ResetTo replaces the buffered content with text (used when an upstream
// edit_content rewrite rewinds already-buffered bytes — issue #23).
func (b *UltraBuffer) ResetTo(text string) {
	b.chunks = append(b.chunks[:0], text)
	b.full.Reset()
	b.full.WriteString(text)
}

// Finish runs validation + optional repair, then replays through the stock
// interceptor to preserve streaming emission semantics (content deltas +
// tool-call deltas in order, rune-safe, fence-tolerant).
func (b *UltraBuffer) Finish() (content string, toolCalls []map[string]interface{}) {
	text := b.full.String()
	if repaired, _ := RepairUltraBuffer(text, string(b.toolsRaw), nil); repaired != text {
		text = repaired
	}
	in := &AgentStreamInterceptor{}
	var cb strings.Builder
	// Replay in one pass: the interceptor's hold-back logic is a function
	// of chunk splits, and a single Feed+Finish preserves its invariants
	// (markers never leak, fences swallowed, indices sequential).
	parsed := in.Feed(text)
	cb.WriteString(parsed.Content)
	toolCalls = append(toolCalls, parsed.ToolCalls...)
	fin := in.Finish()
	cb.WriteString(fin.Content)
	toolCalls = append(toolCalls, fin.ToolCalls...)
	return cb.String(), toolCalls
}

// FullText returns the buffered raw text (pre-repair), for logging/tests.
func (b *UltraBuffer) FullText() string { return b.full.String() }
