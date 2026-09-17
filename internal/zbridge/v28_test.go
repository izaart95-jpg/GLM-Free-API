// v28_test.go — V28 loader: hardware policy, download/verify/cache,
// canonical spec, and repair-prompt shape (§5).

package zbridge

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestV28CanonicalSpecPinned(t *testing.T) {
	spec := V28CanonicalSpec()
	for _, want := range []string{agentToolStart, agentToolEnd, agentCallSchema, `"name"`, `"arguments"`} {
		if !strings.Contains(spec, want) {
			t.Errorf("canonical spec missing %q:\n%s", want, spec)
		}
	}
}

func TestV28RepairPromptUsesFragmentOnly(t *testing.T) {
	frag := `<tool_call>{"tool":"bash"}</tool_call>`
	tools := `[{"name":"bash"}]`
	p := V28RepairPrompt(frag, tools)
	if !strings.Contains(p, frag) {
		t.Errorf("prompt missing fragment")
	}
	if !strings.Contains(p, agentToolStart) || !strings.Contains(p, agentToolEnd) {
		t.Errorf("prompt missing canonical spec")
	}
	// The prompt carries the fragment + tools, never a second smuggled copy
	// of surrounding prose: exactly one <malformed> section.
	if strings.Count(p, "<malformed>") != 1 {
		t.Errorf("prompt must carry exactly one <malformed> fragment section")
	}
	if !strings.Contains(p, "VERBATIM") || !strings.Contains(p, "never invent") {
		t.Errorf("prompt missing faithfulness rules:\n%s", p)
	}
}

func TestV28CPUGuardAbortsWithoutForce(t *testing.T) {
	restore := withUltraGate(true, "ultra", false)
	defer restore()
	hasGPU, _ := DetectGPU()
	if hasGPU {
		t.Skip("GPU present: CPU guard not exercisable on this host")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := EnsureV28ForUltra(ctx)
	if err == nil {
		t.Fatal("CPU-only ultra without --force-cpu must abort startup")
	}
	msg := strings.ToLower(err.Error())
	if !strings.Contains(msg, "force-cpu") {
		t.Errorf("abort must name --force-cpu (got %v)", err)
	}
	if !strings.Contains(msg, "not recommended") {
		t.Errorf("abort must warn CPU inference is NOT recommended (got %v)", err)
	}
}

func TestV28AssetsLocalOverrideNoNetwork(t *testing.T) {
	// A directory containing weights + config satisfies the bundle gate
	// without any download (air-gapped / CI path).
	lora := t.TempDir()
	if err := os.WriteFile(filepath.Join(lora, "adapter_model.safetensors"), []byte("fake-weights"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(lora, "adapter_config.json"), []byte(`{"r":16,"lora_alpha":32}`), 0o644); err != nil {
		t.Fatal(err)
	}
	base := t.TempDir()
	t.Setenv("V28_LORA_PATH", lora)
	t.Setenv("V28_BASE_PATH", base)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	gotBase, gotLora, err := EnsureV28Assets(ctx)
	if err != nil {
		t.Fatalf("local-override assets failed: %v", err)
	}
	if gotBase != base {
		t.Errorf("base = %q, want override %q", gotBase, base)
	}
	if gotLora != lora {
		t.Errorf("lora = %q, want override %q", gotLora, lora)
	}
}

func TestV28BundleVerificationRejectsIncomplete(t *testing.T) {
	empty := t.TempDir() // no weights/config
	if err := verifyLoRABundle(empty); err == nil {
		t.Error("empty dir must fail bundle verification")
	} else if !strings.Contains(err.Error(), "adapter_model.safetensors") {
		t.Errorf("error must name missing weights (got %v)", err)
	}
	badJSON := t.TempDir()
	_ = os.WriteFile(filepath.Join(badJSON, "adapter_model.safetensors"), []byte("w"), 0o644)
	_ = os.WriteFile(filepath.Join(badJSON, "adapter_config.json"), []byte("{not json"), 0o644)
	if err := verifyLoRABundle(badJSON); err == nil {
		t.Error("invalid config JSON must fail verification")
	}
}

func TestV28DetectGPUReturnsEvidence(t *testing.T) {
	ok, evidence := DetectGPU()
	if strings.TrimSpace(evidence) == "" {
		t.Error("DetectGPU must always return human-readable evidence")
	}
	t.Logf("gpu=%v evidence=%s", ok, evidence)
}
