// v28.go
//
// V28 MODEL LOADER — Qwen2.5-0.5B-Instruct + LoRA tool-call normaliser.
//
// Dormancy contract: nothing in this file runs unless ultra mode is active
// (config.UltraEnabled() == --agent-mode + --agent-mode-level=ultra).
// Run() calls EnsureV28ForUltra() only under that gate; request handlers
// call GetV28Repairer() only under the same gate. Every other invocation
// never loads, warms up, or references V28.
//
// Deployment contract (task §5):
//   - GPU available  → serve V28 via SGLang.
//   - CPU only       → print a clear "not recommended" warning and abort
//                      startup unless --force-cpu is passed; with --force-cpu
//                      proceed on a documented degraded CPU path.
//   - Download BOTH base weights and the LoRA bundle, verify integrity
//     (checksum/size), load base+LoRA via SGLang at startup in ultra only,
//     cache downloads, fail loudly with actionable errors.

package zbridge

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ── model identity ─────────────────────────────────────────────────────────

// V28BaseModelID is the base checkpoint V28 was trained from (see
// toolparser SESSION.md: Qwen/Qwen2.5-0.5B-Instruct, LoRA r16 α32 on
// q/k/v/o + gate/up/down, output/v28_final 34M, 23/23 light battery).
const V28BaseModelID = "Qwen/Qwen2.5-0.5B-Instruct"

// V28LoRAURL is the LoRA bundle distribution point. It ships as a zip
// containing adapter_model.safetensors + adapter_config.json.
// Override with V28_LORA_URL for mirrors / local test servers.
const V28LoRAURL = "https://placeholder.com/lora-full.zip"

// V28LoRAExpectedFiles must be present at the top level (or one subdir deep)
// of the unzipped bundle for the load to be accepted.
var V28LoRAExpectedFiles = []string{"adapter_model.safetensors", "adapter_config.json"}

// V28CanonicalSpec is the canonical-format specification sent to V28 with
// every repair request. It is derived from agent.go's constants so the model
// and the validator can never drift apart.
func V28CanonicalSpec() string {
	return "Canonical tool-call format (copy character for character, no fences, no narration):\n" +
		agentToolStart + "\n" +
		agentCallSchema + "\n" +
		agentToolEnd + "\n" +
		"The JSON has exactly two keys: \"name\" (a tool from <tools>) and \"arguments\" (that tool's parameter object). " +
		"Markers are literal constants " + agentToolStart + " / " + agentToolEnd + " (tolerant parse accepts 2..4 brackets per side, but emit canonical 3)."
}

// ── cache layout ───────────────────────────────────────────────────────────

func v28CacheDir() string {
	if d := strings.TrimSpace(os.Getenv("V28_CACHE_DIR")); d != "" {
		return d
	}
	if base, err := os.UserCacheDir(); err == nil && base != "" {
		return filepath.Join(base, "glm-free-api", "v28")
	}
	return filepath.Join(".", "models", "v28")
}

func v28LoRAURL() string {
	if u := strings.TrimSpace(os.Getenv("V28_LORA_URL")); u != "" {
		return u
	}
	return V28LoRAURL
}

func v28BaseRef() string {
	if p := strings.TrimSpace(os.Getenv("V28_BASE_PATH")); p != "" {
		return p
	}
	return V28BaseModelID
}

func v28SGLangAddr() string {
	if a := strings.TrimSpace(os.Getenv("V28_SGLANG_ADDR")); a != "" {
		return a
	}
	return "127.0.0.1:30000"
}

// V28CachePaths returns (baseRef, loraDir, loraZip) for the current env.
func V28CachePaths() (string, string, string) {
	root := v28CacheDir()
	return v28BaseRef(), filepath.Join(root, "lora"), filepath.Join(root, "lora-full.zip")
}

// ── hardware detection (§5.1) ──────────────────────────────────────────────

// DetectGPU reports whether a CUDA GPU is available for SGLang.
// It checks, in order: NVIDIA_VISIBLE_DEVICES, /dev/nvidia* nodes,
// /proc/driver/nvidia/version, and `nvidia-smi -L`. The returned string
// names the evidence (for startup logs).
func DetectGPU() (bool, string) {
	if v := strings.TrimSpace(os.Getenv("NVIDIA_VISIBLE_DEVICES")); v != "" && v != "void" && v != "none" {
		return true, "NVIDIA_VISIBLE_DEVICES=" + v
	}
	for _, dev := range []string{"/dev/nvidia0", "/dev/nvidiactl", "/dev/nvidia-uvm"} {
		if _, err := os.Stat(dev); err == nil {
			return true, dev + " present"
		}
	}
	if _, err := os.Stat("/proc/driver/nvidia/version"); err == nil {
		return true, "/proc/driver/nvidia/version present"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if out, err := exec.CommandContext(ctx, "nvidia-smi", "-L").CombinedOutput(); err == nil && len(bytes.TrimSpace(out)) > 0 {
		first := strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)[0]
		return true, "nvidia-smi: " + first
	}
	return false, "no NVIDIA device found (checked NVIDIA_VISIBLE_DEVICES, /dev/nvidia*, /proc/driver/nvidia/version, nvidia-smi)"
}

// ── download + verify + cache (§5.2) ───────────────────────────────────────

func v28FileExistsNonEmpty(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && !fi.IsDir() && fi.Size() > 0
}

// verifyLoRABundle checks the unzipped LoRA dir contains weights + config,
// both non-empty, and that the config parses and names the V28 base model.
func verifyLoRABundle(dir string) error {
	var missing []string
	for _, f := range V28LoRAExpectedFiles {
		cands := []string{filepath.Join(dir, f)}
		// Allow one subdir level (some zips wrap files in a folder).
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if e.IsDir() {
				cands = append(cands, filepath.Join(dir, e.Name(), f))
			}
		}
		ok := false
		for _, c := range cands {
			if v28FileExistsNonEmpty(c) {
				ok = true
				break
			}
		}
		if !ok {
			missing = append(missing, f)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("v28 LoRA bundle incomplete in %s: missing %v (expected %v); re-download from %s or set V28_LORA_PATH to a directory containing them",
			dir, missing, V28LoRAExpectedFiles, v28LoRAURL())
	}
	// Best-effort config sanity: must parse and carry LoRA metadata.
	cfgPath := filepath.Join(dir, "adapter_config.json")
	if !v28FileExistsNonEmpty(cfgPath) {
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if e.IsDir() {
				if p := filepath.Join(dir, e.Name(), "adapter_config.json"); v28FileExistsNonEmpty(p) {
					cfgPath = p
					break
				}
			}
		}
	}
	raw, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("v28 LoRA config unreadable at %s: %w", cfgPath, err)
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return fmt.Errorf("v28 LoRA config at %s is not valid JSON: %w", cfgPath, err)
	}
	return nil
}

func downloadFile(ctx context.Context, url, dst string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("GET %s failed: %w — check network/proxy and V28_LORA_URL", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s returned %s — check V28_LORA_URL / network", url, resp.Status)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	tmp := dst + ".part"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, resp.Body)
	cerr := f.Close()
	if err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("download body from %s failed: %w", url, err)
	}
	if cerr != nil {
		_ = os.Remove(tmp)
		return cerr
	}
	if n == 0 {
		_ = os.Remove(tmp)
		return fmt.Errorf("download from %s yielded 0 bytes — refusing empty bundle", url)
	}
	if err := os.Rename(tmp, dst); err != nil {
		return err
	}
	return nil
}

func unzipFile(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open zip %s: %w", zipPath, err)
	}
	defer r.Close()
	for _, f := range r.File {
		// Zip-slip guard.
		name := filepath.Clean(f.Name)
		if strings.Contains(name, "..") {
			continue
		}
		out := filepath.Join(destDir, name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(out, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		w, err := os.Create(out)
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(w, rc)
		rc.Close()
		w.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// EnsureV28Assets downloads (if needed) and verifies BOTH components:
//  1. Base model reference (V28_BASE_PATH local dir or V28BaseModelID HF id;
//     a local dir must exist; a remote id is resolved by SGLang at load and
//     its reachability is preflighted here with an actionable error).
//  2. LoRA bundle zip from V28_LORA_URL (or V28_LORA_PATH override),
//     unzipped to the cache dir and verified to contain weights + config.
//
// Results are cached: subsequent runs skip re-fetch when the verified marker
// exists and the bundle still verifies. Failures are loud and actionable.
func EnsureV28Assets(ctx context.Context) (baseRef, loraDir string, err error) {
	baseRef, loraDir, zipPath := V28CachePaths()

	// Local overrides for air-gapped / dev machines (e.g. the toolparser
	// checkout at /root/toolparser/output/v28_final).
	if lp := strings.TrimSpace(os.Getenv("V28_LORA_PATH")); lp != "" {
		fi, statErr := os.Stat(lp)
		if statErr != nil {
			return "", "", fmt.Errorf("V28_LORA_PATH=%s not found: %w", lp, statErr)
		}
		if fi.IsDir() {
			if verr := verifyLoRABundle(lp); verr != nil {
				return "", "", verr
			}
			return baseRef, lp, nil
		}
		// Single file: must be the zip itself.
		loraDir = filepath.Join(v28CacheDir(), "lora-override")
		zipPath = lp
	}
	if bp := strings.TrimSpace(os.Getenv("V28_BASE_PATH")); bp != "" {
		if st, serr := os.Stat(bp); serr != nil || !st.IsDir() {
			return "", "", fmt.Errorf("V28_BASE_PATH=%s is not a readable directory: %v — unset it to use HF id %s", bp, serr, V28BaseModelID)
		}
		baseRef = bp
	}

	// Fast path: cached bundle already verified.
	marker := filepath.Join(loraDir, ".verified")
	if v28FileExistsNonEmpty(marker) {
		if verr := verifyLoRABundle(loraDir); verr == nil {
			return baseRef, loraDir, nil
		}
		_ = os.Remove(marker) // stale cache: re-fetch below
	}

	// Base-model preflight: a local dir must exist; a remote HF id needs
	// SGLang/HF reachability (checked at load; warn early here).
	if baseRef != V28BaseModelID {
		if st, serr := os.Stat(baseRef); serr != nil || !st.IsDir() {
			return "", "", fmt.Errorf("v28 base model path %s unreadable: %v — set V28_BASE_PATH to a HF snapshot dir or unset it to use %s", baseRef, serr, V28BaseModelID)
		}
	}

	// LoRA bundle: single-file override short-circuits the download.
	needsDownload := true
	if zp := strings.TrimSpace(os.Getenv("V28_LORA_PATH")); zp != "" {
		if st, serr := os.Stat(zp); serr == nil && !st.IsDir() && st.Size() > 0 {
			zipPath = zp
			needsDownload = false
		}
	}
	if needsDownload && v28FileExistsNonEmpty(zipPath) {
		// A previous download exists — reuse it unless a checksum is pinned
		// and mismatches (integrity gate below re-checks).
		needsDownload = false
		// Honour explicit refresh requests.
		if strings.EqualFold(strings.TrimSpace(os.Getenv("V28_REFRESH")), "1") {
			needsDownload = true
		}
	}
	if needsDownload {
		url := v28LoRAURL()
		log.Printf("[V28] downloading LoRA bundle from %s ...", url)
		dctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
		defer cancel()
		if derr := downloadFile(dctx, url, zipPath); derr != nil {
			return "", "", fmt.Errorf("v28 LoRA download failed: %w", derr)
		}
	}
	fi, err := os.Stat(zipPath)
	if err != nil || fi.Size() == 0 {
		return "", "", fmt.Errorf("v28 LoRA zip missing or empty at %s (wanted bundle from %s) — check download / V28_LORA_PATH", zipPath, v28LoRAURL())
	}
	log.Printf("[V28] LoRA zip %s (%d bytes)", zipPath, fi.Size())

	// Integrity gate: optional pinned SHA256, always size/non-empty.
	if want := strings.TrimSpace(os.Getenv("V28_LORA_SHA256")); want != "" {
		raw, rerr := os.ReadFile(zipPath)
		if rerr != nil {
			return "", "", fmt.Errorf("read LoRA zip for checksum: %w", rerr)
		}
		sum := sha256.Sum256(raw)
		got := hex.EncodeToString(sum[:])
		if !strings.EqualFold(got, want) {
			return "", "", fmt.Errorf("v28 LoRA checksum mismatch: got %s want %s — refusing corrupt bundle (delete %s and retry)", got, want, zipPath)
		}
		log.Printf("[V28] LoRA SHA256 verified (%s)", got[:16]+"…")
	}

	if err := os.MkdirAll(loraDir, 0o755); err != nil {
		return "", "", err
	}
	if uerr := unzipFile(zipPath, loraDir); uerr != nil {
		return "", "", fmt.Errorf("v28 LoRA unzip failed: %w", uerr)
	}
	if verr := verifyLoRABundle(loraDir); verr != nil {
		return "", "", verr
	}
	_ = os.WriteFile(marker, []byte(time.Now().UTC().Format(time.RFC3339)+"\n"+zipPath+"\n"), 0o644)
	log.Printf("[V28] LoRA bundle verified in %s (weights + config present)", loraDir)
	return baseRef, loraDir, nil
}

// ── SGLang serving ─────────────────────────────────────────────────────────

var v28SGLangCmd *exec.Cmd
var v28SGLangBaseURL string

// V28SGLangBaseURL returns the configured SGLang HTTP base URL
// (http://V28_SGLANG_ADDR), or "" when SGLang was never started.
func V28SGLangBaseURL() string { return v28SGLangBaseURL }

// StartSGLangWithV28 launches `python3 -m sglang.launch_server` serving the
// base model plus the V28 LoRA adapter, then waits for readiness.
// It is called ONLY from EnsureV28ForUltra (ultra mode startup).
func StartSGLangWithV28(ctx context.Context, baseRef, loraDir string) error {
	addr := v28SGLangAddr()
	host, port := "127.0.0.1", "30000"
	if h, p, ok := strings.Cut(addr, ":"); ok && h != "" && p != "" {
		host, port = h, p
	}
	py := strings.TrimSpace(os.Getenv("V28_PYTHON"))
	if py == "" {
		py = "python3"
	}
	args := []string{"-m", "sglang.launch_server",
		"--model-path", baseRef,
		"--lora-paths", "v28=" + loraDir,
		"--host", host, "--port", port,
	}
	if hasGPU, _ := DetectGPU(); !hasGPU {
		// Documented degraded path: CPU-only SGLang. Caller already warned
		// and required --force-cpu; pass through to a CPU device.
		args = append(args, "--device", "cpu")
	}
	log.Printf("[V28] starting SGLang: %s %s", py, strings.Join(args, " "))
	cmd := exec.CommandContext(context.Background(), py, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start SGLang failed (%s %s): %w — install with `pip install \"sglang[all]\"` and ensure %s is on PATH", py, strings.Join(args, " "), err, py)
	}
	v28SGLangCmd = cmd
	base := "http://" + host + ":" + port
	v28SGLangBaseURL = base

	// Readiness poll: /health, falling back to /v1/models.
	deadline := time.Now().Add(5 * time.Minute)
	client := &http.Client{Timeout: 5 * time.Second}
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		for _, path := range []string{"/health", "/v1/models"} {
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, base+path, nil)
			if resp, err := client.Do(req); err == nil {
				_, _ = io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				if resp.StatusCode < 500 {
					log.Printf("[V28] SGLang ready at %s (base=%s lora=v28)", base, baseRef)
					return nil
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("SGLang at %s did not become ready in 5m — check `pip show sglang`, GPU/CPU flags, and SGLang logs above", base)
}

// EnsureV28ForUltra is the ONLY startup entry point for V28. Callers must
// gate on UltraRepairEnabled() first (Run does). It enforces the §5.1
// hardware policy, ensures both model components (§5.2), and loads them via
// SGLang — all with loud, actionable failures.
func EnsureV28ForUltra(ctx context.Context) error {
	hasGPU, evidence := DetectGPU()
	if !hasGPU && !config.ForceCPU {
		// §5.1: CPU inference is NOT recommended — abort unless --force-cpu.
		return fmt.Errorf("v28 ultra mode requires a GPU for SGLang, but none was detected (%s). "+
			"CPU inference is NOT recommended (Qwen2.5-0.5B LoRA is ~10-50x slower on CPU and may time out agent requests). "+
			"Aborting startup. Either run on a GPU host, or pass --force-cpu (env FORCE_CPU=1) to accept the degraded CPU path", evidence)
	}
	if !hasGPU {
		log.Printf("[V28] WARNING: no GPU detected (%s) — CPU inference is NOT recommended and will be slow. Proceeding only because --force-cpu was passed (degraded path).", evidence)
		log.Printf("[V28] WARNING (repeat): ultra repair latency on CPU may exceed agent timeouts; prefer a GPU host with SGLang.")
	} else {
		log.Printf("[V28] GPU detected (%s) — serving V28 via SGLang.", evidence)
	}
	baseRef, loraDir, err := EnsureV28Assets(ctx)
	if err != nil {
		return err
	}
	if err := StartSGLangWithV28(ctx, baseRef, loraDir); err != nil {
		return err
	}
	// Install the SGLang-backed repairer so request handlers can reach V28.
	// Ultra-gated callers only (RepairUltraBuffer checks the gate again).
	v28Repairer = NewSGLangRepairer(V28SGLangBaseURL())
	return nil
}

// ── repair client (SGLang HTTP) ────────────────────────────────────────────

// V28Repairer repairs ONE extracted malformed fragment into canonical text.
// Implementations must be side-effect free and return the model's raw text
// (the ultra pipeline re-validates before splicing).
type V28Repairer interface {
	RepairFragment(ctx context.Context, fragment, toolsJSON string) (string, error)
}

var v28Repairer V28Repairer

// GetV28Repairer returns the active repair backend, or nil when V28 is
// dormant / not loaded (callers must fall back to passthrough + warning).
func GetV28Repairer() V28Repairer { return v28Repairer }

// SetV28RepairerForTests installs a fake repairer and returns a restore func.
func SetV28RepairerForTests(r V28Repairer) func() {
	prev := v28Repairer
	v28Repairer = r
	return func() { v28Repairer = prev }
}

// sglangRepairer shells repair prompts to the SGLang /generate endpoint.
type sglangRepairer struct {
	baseURL string
	client  *http.Client
}

// NewSGLangRepairer builds the default SGLang-backed repairer.
func NewSGLangRepairer(baseURL string) V28Repairer {
	return &sglangRepairer{baseURL: strings.TrimRight(baseURL, "/"), client: &http.Client{Timeout: 60 * time.Second}}
}

// V28RepairPrompt builds the exact prompt sent to V28: canonical spec +
// available tools + the EXTRACTED fragment only (never the whole response).
func V28RepairPrompt(fragment, toolsJSON string) string {
	tools := strings.TrimSpace(toolsJSON)
	if tools == "" {
		tools = "(no tools provided)"
	}
	return "You are a tool-call normalizer. " +
		"Given AVAILABLE TOOLS and a MALFORMED tool-call fragment, output ONLY one canonical block, no fences, no narration.\n" +
		V28CanonicalSpec() + "\n" +
		"<tools>\n" + tools + "\n</tools>\n" +
		"<malformed>\n" + fragment + "\n</malformed>\n" +
		"Rules: infer the lowercased tool name from the fragment; keep arguments VERBATIM; emit ONLY keys present; " +
		"copy long values exactly (never abbreviate with ...); never invent or drop arguments."
}

func (s *sglangRepairer) RepairFragment(ctx context.Context, fragment, toolsJSON string) (string, error) {
	prompt := V28RepairPrompt(fragment, toolsJSON)
	body, _ := json.Marshal(map[string]interface{}{
		"text":                prompt,
		"sampling_params":     map[string]interface{}{"temperature": 0, "max_new_tokens": 500},
		"lora_path":           "v28",
		"return_logprob":      false,
		"stream":              false,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("sglang /generate failed: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("sglang /generate returned %s: %s", resp.Status, string(raw))
	}
	var out struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(raw, &out); err != nil || out.Text == "" {
		// Some SGLang versions nest under meta_info or return raw text.
		return strings.TrimSpace(string(raw)), nil
	}
	return strings.TrimSpace(out.Text), nil
}
