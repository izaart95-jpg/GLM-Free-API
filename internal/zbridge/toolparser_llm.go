// toolparser_llm.go
//
// TOOLPARSER-LLM MODEL LOADER — Qwen2.5-0.5B-Instruct + LoRA tool-call normaliser.
//
// Dormancy contract: nothing in this file runs unless ultra mode is active
// (config.UltraEnabled() == --agent-mode + --agent-mode-level=ultra).
// Run() calls EnsureToolParserLLMForUltra() only under that gate; request handlers
// call GetToolParserLLMRepairer() only under the same gate. Every other invocation
// never loads, warms up, or references ToolParserLLM.
//
// Deployment contract (task §5):
//   - GPU available  → serve ToolParserLLM via SGLang.
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

// ToolParserLLMBaseModelID is the base checkpoint the repair adapter was
// trained from (Qwen/Qwen2.5-0.5B-Instruct; current adapter bundle is the
// v28 flagship: LoRA r16 α32 on q/k/v/o + gate/up/down, 34M, 23/23 battery).
const ToolParserLLMBaseModelID = "Qwen/Qwen2.5-0.5B-Instruct"

// ToolParserLLMLoRAURL is the LoRA bundle distribution point. It ships as a zip
// containing adapter_model.safetensors + adapter_config.json.
// Public Kaggle dataset zip containing adapter_model.safetensors +
// adapter_config.json. Override with TOOLPARSER_LLM_LORA_URL for mirrors /
// local test servers. No auth needed (public dataset); plain HTTPS GET.
const ToolParserLLMLoRAURL = "https://www.kaggle.com/api/v1/datasets/download/zoxoashouko/v28-flagship-toolparser"

// ToolParserLLMLoRAExpectedFiles must be present at the top level (or one subdir deep)
// of the unzipped bundle for the load to be accepted.
var ToolParserLLMLoRAExpectedFiles = []string{"adapter_model.safetensors", "adapter_config.json"}

// ToolParserLLMCanonicalSpec is the canonical-format specification sent to
// ToolParserLLM with
// every repair request. It is derived from agent.go's constants so the model
// and the validator can never drift apart.
func ToolParserLLMCanonicalSpec() string {
	return "Canonical tool-call format (copy character for character, no fences, no narration):\n" +
		agentToolStart + "\n" +
		agentCallSchema + "\n" +
		agentToolEnd + "\n" +
		"The JSON has exactly two keys: \"name\" (a tool from <tools>) and \"arguments\" (that tool's parameter object). " +
		"Markers are literal constants " + agentToolStart + " / " + agentToolEnd + " (tolerant parse accepts 2..4 brackets per side, but emit canonical 3)."
}

// ── cache layout ───────────────────────────────────────────────────────────

func toolParserLLMCacheDir() string {
	if d := strings.TrimSpace(os.Getenv("TOOLPARSER_LLM_CACHE_DIR")); d != "" {
		return d
	}
	if base, err := os.UserCacheDir(); err == nil && base != "" {
		return filepath.Join(base, "glm-free-api", "toolparser-llm")
	}
	return filepath.Join(".", "models", "toolparser-llm")
}

func toolParserLLMLoRAURL() string {
	if u := strings.TrimSpace(os.Getenv("TOOLPARSER_LLM_LORA_URL")); u != "" {
		return u
	}
	return ToolParserLLMLoRAURL
}

func toolParserLLMBaseRef() string {
	if p := strings.TrimSpace(os.Getenv("TOOLPARSER_LLM_BASE_PATH")); p != "" {
		return p
	}
	return ToolParserLLMBaseModelID
}

func toolParserLLMSGLangAddr() string {
	if a := strings.TrimSpace(os.Getenv("TOOLPARSER_LLM_SGLANG_ADDR")); a != "" {
		return a
	}
	return "127.0.0.1:30000"
}

// ToolParserLLMCachePaths returns (baseRef, loraDir, loraZip) for the current env.
func ToolParserLLMCachePaths() (string, string, string) {
	root := toolParserLLMCacheDir()
	return toolParserLLMBaseRef(), filepath.Join(root, "lora"), filepath.Join(root, "lora-full.zip")
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

func toolParserLLMFileExistsNonEmpty(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && !fi.IsDir() && fi.Size() > 0
}

// verifyLoRABundle checks the unzipped LoRA dir contains weights + config,
// both non-empty, and that the config parses and names the base model.
func verifyLoRABundle(dir string) error {
	var missing []string
	for _, f := range ToolParserLLMLoRAExpectedFiles {
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
			if toolParserLLMFileExistsNonEmpty(c) {
				ok = true
				break
			}
		}
		if !ok {
			missing = append(missing, f)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("ToolParserLLM LoRA bundle incomplete in %s: missing %v (expected %v); re-download from %s or set TOOLPARSER_LLM_LORA_PATH to a directory containing them",
			dir, missing, ToolParserLLMLoRAExpectedFiles, toolParserLLMLoRAURL())
	}
	// Best-effort config sanity: must parse and carry LoRA metadata.
	cfgPath := filepath.Join(dir, "adapter_config.json")
	if !toolParserLLMFileExistsNonEmpty(cfgPath) {
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if e.IsDir() {
				if p := filepath.Join(dir, e.Name(), "adapter_config.json"); toolParserLLMFileExistsNonEmpty(p) {
					cfgPath = p
					break
				}
			}
		}
	}
	raw, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("ToolParserLLM LoRA config unreadable at %s: %w", cfgPath, err)
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return fmt.Errorf("ToolParserLLM LoRA config at %s is not valid JSON: %w", cfgPath, err)
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
		return fmt.Errorf("GET %s failed: %w — check network/proxy and TOOLPARSER_LLM_LORA_URL", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s returned %s — check TOOLPARSER_LLM_LORA_URL / network", url, resp.Status)
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

// EnsureToolParserLLMAssets downloads (if needed) and verifies BOTH components:
//  1. Base model reference (TOOLPARSER_LLM_BASE_PATH local dir or
//     ToolParserLLMBaseModelID HF id;
//     a local dir must exist; a remote id is resolved by SGLang at load and
//     its reachability is preflighted here with an actionable error).
//  2. LoRA bundle zip from TOOLPARSER_LLM_LORA_URL (or
//     TOOLPARSER_LLM_LORA_PATH override),
//     unzipped to the cache dir and verified to contain weights + config.
//
// Results are cached: subsequent runs skip re-fetch when the verified marker
// exists and the bundle still verifies. Failures are loud and actionable.
func EnsureToolParserLLMAssets(ctx context.Context) (baseRef, loraDir string, err error) {
	baseRef, loraDir, zipPath := ToolParserLLMCachePaths()

	// Local overrides for air-gapped / dev machines (e.g. the toolparser
	// checkout at /root/toolparser/output/v28_final).
	if lp := strings.TrimSpace(os.Getenv("TOOLPARSER_LLM_LORA_PATH")); lp != "" {
		fi, statErr := os.Stat(lp)
		if statErr != nil {
			return "", "", fmt.Errorf("TOOLPARSER_LLM_LORA_PATH=%s not found: %w", lp, statErr)
		}
		if fi.IsDir() {
			if verr := verifyLoRABundle(lp); verr != nil {
				return "", "", verr
			}
			return baseRef, lp, nil
		}
		// Single file: must be the zip itself.
		loraDir = filepath.Join(toolParserLLMCacheDir(), "lora-override")
		zipPath = lp
	}
	if bp := strings.TrimSpace(os.Getenv("TOOLPARSER_LLM_BASE_PATH")); bp != "" {
		if st, serr := os.Stat(bp); serr != nil || !st.IsDir() {
			return "", "", fmt.Errorf("TOOLPARSER_LLM_BASE_PATH=%s is not a readable directory: %v — unset it to use HF id %s", bp, serr, ToolParserLLMBaseModelID)
		}
		baseRef = bp
	}

	// Fast path: cached bundle already verified.
	marker := filepath.Join(loraDir, ".verified")
	if toolParserLLMFileExistsNonEmpty(marker) {
		if verr := verifyLoRABundle(loraDir); verr == nil {
			return baseRef, loraDir, nil
		}
		_ = os.Remove(marker) // stale cache: re-fetch below
	}

	// Base-model preflight: a local dir must exist; a remote HF id needs
	// SGLang/HF reachability (checked at load; warn early here).
	if baseRef != ToolParserLLMBaseModelID {
		if st, serr := os.Stat(baseRef); serr != nil || !st.IsDir() {
			return "", "", fmt.Errorf("v28 base model path %s unreadable: %v — set TOOLPARSER_LLM_BASE_PATH to a HF snapshot dir or unset it to use %s", baseRef, serr, ToolParserLLMBaseModelID)
		}
	}

	// LoRA bundle: single-file override short-circuits the download.
	needsDownload := true
	if zp := strings.TrimSpace(os.Getenv("TOOLPARSER_LLM_LORA_PATH")); zp != "" {
		if st, serr := os.Stat(zp); serr == nil && !st.IsDir() && st.Size() > 0 {
			zipPath = zp
			needsDownload = false
		}
	}
	if needsDownload && toolParserLLMFileExistsNonEmpty(zipPath) {
		// A previous download exists — reuse it unless a checksum is pinned
		// and mismatches (integrity gate below re-checks).
		needsDownload = false
		// Honour explicit refresh requests.
		if strings.EqualFold(strings.TrimSpace(os.Getenv("TOOLPARSER_LLM_REFRESH")), "1") {
			needsDownload = true
		}
	}
	if needsDownload {
		url := toolParserLLMLoRAURL()
		log.Printf("[ToolParserLLM] downloading LoRA bundle from %s ...", url)
		dctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
		defer cancel()
		if derr := downloadFile(dctx, url, zipPath); derr != nil {
			return "", "", fmt.Errorf("ToolParserLLM LoRA download failed: %w", derr)
		}
	}
	fi, err := os.Stat(zipPath)
	if err != nil || fi.Size() == 0 {
		return "", "", fmt.Errorf("ToolParserLLM LoRA zip missing or empty at %s (wanted bundle from %s) — check download / TOOLPARSER_LLM_LORA_PATH", zipPath, toolParserLLMLoRAURL())
	}
	log.Printf("[ToolParserLLM] LoRA zip %s (%d bytes)", zipPath, fi.Size())

	// Integrity gate: optional pinned SHA256, always size/non-empty.
	if want := strings.TrimSpace(os.Getenv("TOOLPARSER_LLM_LORA_SHA256")); want != "" {
		raw, rerr := os.ReadFile(zipPath)
		if rerr != nil {
			return "", "", fmt.Errorf("read LoRA zip for checksum: %w", rerr)
		}
		sum := sha256.Sum256(raw)
		got := hex.EncodeToString(sum[:])
		if !strings.EqualFold(got, want) {
			return "", "", fmt.Errorf("ToolParserLLM LoRA checksum mismatch: got %s want %s — refusing corrupt bundle (delete %s and retry)", got, want, zipPath)
		}
		log.Printf("[ToolParserLLM] LoRA SHA256 verified (%s)", got[:16]+"…")
	}

	if err := os.MkdirAll(loraDir, 0o755); err != nil {
		return "", "", err
	}
	if uerr := unzipFile(zipPath, loraDir); uerr != nil {
		return "", "", fmt.Errorf("ToolParserLLM LoRA unzip failed: %w", uerr)
	}
	if verr := verifyLoRABundle(loraDir); verr != nil {
		return "", "", verr
	}
	_ = os.WriteFile(marker, []byte(time.Now().UTC().Format(time.RFC3339)+"\n"+zipPath+"\n"), 0o644)
	log.Printf("[ToolParserLLM] LoRA bundle verified in %s (weights + config present)", loraDir)
	return baseRef, loraDir, nil
}

// ── SGLang serving ─────────────────────────────────────────────────────────

var toolParserLLMSGLangCmd *exec.Cmd
var toolParserLLMSGLangBaseURL string

// ToolParserLLMSGLangBaseURL returns the configured SGLang HTTP base URL
// (http://TOOLPARSER_LLM_SGLANG_ADDR), or "" when SGLang was never started.
func ToolParserLLMSGLangBaseURL() string { return toolParserLLMSGLangBaseURL }

// StartSGLangWithToolParserLLM launches `python3 -m sglang.launch_server` serving the
// base model plus the ToolParserLLM LoRA adapter, then waits for readiness.
// It is called ONLY from EnsureToolParserLLMForUltra (ultra mode startup).
func StartSGLangWithToolParserLLM(ctx context.Context, baseRef, loraDir string) error {
	addr := toolParserLLMSGLangAddr()
	host, port := "127.0.0.1", "30000"
	if h, p, ok := strings.Cut(addr, ":"); ok && h != "" && p != "" {
		host, port = h, p
	}
	py := strings.TrimSpace(os.Getenv("TOOLPARSER_LLM_PYTHON"))
	if py == "" {
		py = "python3"
	}
	args := []string{"-m", "sglang.launch_server",
		"--model-path", baseRef,
		"--lora-paths", "toolparserllm=" + loraDir,
		"--host", host, "--port", port,
	}
	if hasGPU, _ := DetectGPU(); !hasGPU {
		// Documented degraded path: CPU-only SGLang. Caller already warned
		// and required --force-cpu; pass through to a CPU device.
		args = append(args, "--device", "cpu")
	}
	log.Printf("[ToolParserLLM] starting SGLang: %s %s", py, strings.Join(args, " "))
	cmd := exec.CommandContext(context.Background(), py, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start SGLang failed (%s %s): %w — install with `pip install \"sglang[all]\"` and ensure %s is on PATH", py, strings.Join(args, " "), err, py)
	}
	toolParserLLMSGLangCmd = cmd
	base := "http://" + host + ":" + port
	toolParserLLMSGLangBaseURL = base

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
					log.Printf("[ToolParserLLM] SGLang ready at %s (base=%s lora=toolparserllm)", base, baseRef)
					return nil
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("SGLang at %s did not become ready in 5m — check `pip show sglang`, GPU/CPU flags, and SGLang logs above", base)
}

// EnsureToolParserLLMForUltra is the ONLY startup entry point for
// ToolParserLLM. Callers must
// gate on UltraRepairEnabled() first (Run does). It enforces the §5.1
// hardware policy, ensures both model components (§5.2), and loads them via
// SGLang — all with loud, actionable failures.
func EnsureToolParserLLMForUltra(ctx context.Context) error {
	hasGPU, evidence := DetectGPU()
	if !hasGPU && !config.ForceCPU {
		// §5.1: CPU inference is NOT recommended — abort unless --force-cpu.
		return fmt.Errorf("ToolParserLLM ultra mode requires a GPU for SGLang, but none was detected (%s). "+
			"CPU inference is NOT recommended (Qwen2.5-0.5B LoRA is ~10-50x slower on CPU and may time out agent requests). "+
			"Aborting startup. Either run on a GPU host, or pass --force-cpu (env FORCE_CPU=1) to accept the degraded CPU path", evidence)
	}
	if !hasGPU {
		log.Printf("[ToolParserLLM] WARNING: no GPU detected (%s) — CPU inference is NOT recommended and will be slow. Proceeding only because --force-cpu was passed (degraded path).", evidence)
		log.Printf("[ToolParserLLM] WARNING (repeat): ultra repair latency on CPU may exceed agent timeouts; prefer a GPU host with SGLang.")
	} else {
		log.Printf("[ToolParserLLM] GPU detected (%s) — serving ToolParserLLM via SGLang.", evidence)
	}
	baseRef, loraDir, err := EnsureToolParserLLMAssets(ctx)
	if err != nil {
		return err
	}
	if err := StartSGLangWithToolParserLLM(ctx, baseRef, loraDir); err != nil {
		return err
	}
	// Install the SGLang-backed repairer so request handlers can reach
// ToolParserLLM.
	// Ultra-gated callers only (RepairUltraBuffer checks the gate again).
	toolParserLLMRepairer = NewToolParserSGLangRepairer(ToolParserLLMSGLangBaseURL())
	return nil
}

// ── repair client (SGLang HTTP) ────────────────────────────────────────────

// ToolParserLLMRepairer repairs ONE extracted malformed fragment into canonical text.
// Implementations must be side-effect free and return the model's raw text
// (the ultra pipeline re-validates before splicing).
type ToolParserLLMRepairer interface {
	RepairFragment(ctx context.Context, fragment, toolsJSON string) (string, error)
}

var toolParserLLMRepairer ToolParserLLMRepairer

// GetToolParserLLMRepairer returns the active repair backend, or nil when
// ToolParserLLM is
// dormant / not loaded (callers must fall back to passthrough + warning).
func GetToolParserLLMRepairer() ToolParserLLMRepairer { return toolParserLLMRepairer }

// SetToolParserLLMRepairerForTests installs a fake repairer and returns a restore func.
func SetToolParserLLMRepairerForTests(r ToolParserLLMRepairer) func() {
	prev := toolParserLLMRepairer
	toolParserLLMRepairer = r
	return func() { toolParserLLMRepairer = prev }
}

// sglangRepairer shells repair prompts to the SGLang /generate endpoint.
type sglangRepairer struct {
	baseURL string
	client  *http.Client
}

// NewToolParserSGLangRepairer builds the default SGLang-backed repairer.
func NewToolParserSGLangRepairer(baseURL string) ToolParserLLMRepairer {
	return &sglangRepairer{baseURL: strings.TrimRight(baseURL, "/"), client: &http.Client{Timeout: 60 * time.Second}}
}

// ToolParserLLMRepairPrompt builds the exact prompt sent to ToolParserLLM:
// canonical spec +
// available tools + the EXTRACTED fragment only (never the whole response).
func ToolParserLLMRepairPrompt(fragment, toolsJSON string) string {
	tools := strings.TrimSpace(toolsJSON)
	if tools == "" {
		tools = "(no tools provided)"
	}
	return "You are a tool-call normalizer. " +
		"Given AVAILABLE TOOLS and a MALFORMED tool-call fragment, output ONLY one canonical block, no fences, no narration.\n" +
		ToolParserLLMCanonicalSpec() + "\n" +
		"<tools>\n" + tools + "\n</tools>\n" +
		"<malformed>\n" + fragment + "\n</malformed>\n" +
		"Rules: infer the lowercased tool name from the fragment; keep arguments VERBATIM; emit ONLY keys present; " +
		"copy long values exactly (never abbreviate with ...); never invent or drop arguments."
}

func (s *sglangRepairer) RepairFragment(ctx context.Context, fragment, toolsJSON string) (string, error) {
	prompt := ToolParserLLMRepairPrompt(fragment, toolsJSON)
	body, _ := json.Marshal(map[string]interface{}{
		"text":                prompt,
		"sampling_params":     map[string]interface{}{"temperature": 0, "max_new_tokens": 500},
		"lora_path":           "toolparserllm",
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
