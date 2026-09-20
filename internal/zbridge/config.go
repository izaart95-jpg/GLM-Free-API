package zbridge

import (
    "os"
    "strconv"
    "strings"
    "time"
)

// ============================================================================
// CONFIGURATION
// ============================================================================

const (
    // Aliyun captcha credentials
    accessKey       = "LTAI5tSEBwYMwVKAQGpxmvTd"
    secretKey       = "YSKfst7GaVkXwZYvVihJsKF9r89koz"
    sceneID         = "didk33e0"
    maxTokenRetries = 5

    // Z.AI direct config
    SALT_KEY           = "key-@@@@)))()((9))-xxxx&&&%%%%%"
    DEFAULT_FE_VERSION = "prod-fe-1.1.93"
    zaiUserAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
)

// BASE_URL is a var (not const) only so tests can point the bridge at a
// mock upstream; the default value is the production endpoint.
var BASE_URL = "https://chat.z.ai"

// ---------- Config struct (Z.AI) ----------

type Config struct {
    Server struct {
        Port int
        Host string
    }
    Auth struct {
        Enabled bool
        Token   string
    }
    Timeouts struct {
        Default int
    }
    ZaiToken  string
    AgentMode bool
    // AgentModeVariant selects the agent-mode compatibility shim:
    //   "modern" (default) — XML-sectioned prompt shim ported from
    //                        DeepseekFreeAPI (see agent.go)
    //   "legacy"           — the original [ROLE: ...] rewrite shim
    AgentModeVariant string
    // AgentModeLevel gates the opt-in ToolParserLLM intelligent repair layer
    // (see agent_ultra.go / toolparser_llm.go). Only the exact value "ultra"
    // (case insensitive) enables buffering + repair, and only together with
    // AgentMode. Any other value (including the default "") keeps ToolParserLLM
    // completely dormant: no load, no warm-up, no reference.
    AgentModeLevel string
    // ForceCPU explicitly opts into the degraded CPU inference path for
    // ToolParserLLM
    // in ultra mode on machines without a GPU. Without it, ultra mode on a
    // CPU-only host aborts startup with an actionable error.
    ForceCPU bool
    Logging   struct {
        Level  string
        Format string
    }
    KnownModels []string
    // StreamHoldback is the number of runes kept pending at the tail of the
    // streamed content before it is forwarded to clients. Z.AI's stream is
    // edit-based (edit_content can backtrack and rewrite the tail), and an
    // append-only SSE client cannot take back text it already received.
    // Holding back a small window lets ordinary trailing backtracks be
    // absorbed invisibly. 0 disables the hold-back. See issue #23.
    StreamHoldback int
    // SyncMode disables the async session pool and reuses one sticky session
    // instead (up to SESSION_REUSE_COUNT requests before delete + rotate).
    // Used sessions are still deleted on Z.AI after their reuse limit
    // (reused sessions either way — see session_pool.go).
    SyncMode bool
    // SessionPoolSize is the standing batch of pre-made ready chat sessions
    // kept by the async session pool (SESSION_POOL_SIZE, default 5).
    SessionPoolSize int
    // SessionAcquireTimeout bounds, in seconds, how long a request waits for
    // a pooled session before creating one directly instead of stalling
    // (SESSION_ACQUIRE_TIMEOUT, default 10; 0 waits indefinitely).
    SessionAcquireTimeout int
    // SessionReuseCount bounds how many requests one chat session serves
    // before it is deleted on Z.AI and replaced (SESSION_REUSE_COUNT,
    // default 10). 1 restores the legacy throwaway-per-request flow.
    // Reusing a session makes traffic look human (one chat carrying several
    // turns instead of one chat per message) and cuts DELETE churn N-fold,
    // which keeps the Aliyun WAF from flagging the egress IP as a bot.
    // Safe: upstream keeps no server-side history for isolated single-turn
    // payloads (verified e2e: same chat_id + "remember cat" then recall
    // answers "no history"), so reuse does not thread conversations.
    SessionReuseCount int
}

func loadConfig() *Config {
    c := &Config{}
    c.Server.Port = 3001
    c.Server.Host = "0.0.0.0"
    c.Auth.Enabled = true
    c.Auth.Token = "Waguri"
    c.Timeouts.Default = 300000
    c.ZaiToken = ""
    c.AgentMode = false
    c.AgentModeVariant = "modern"
    c.AgentModeLevel = ""
    c.ForceCPU = false
    c.Logging.Level = "debug"
    c.Logging.Format = "text"
    c.KnownModels = []string{"GLM-5.1", "GLM-5"}
    c.StreamHoldback = 24
    c.SyncMode = false
    c.SessionPoolSize = defaultPoolSize
    c.SessionAcquireTimeout = int(defaultPoolWait / time.Second)
    c.SessionReuseCount = defaultSessionReuse

    if p := os.Getenv("PORT"); p != "" {
        if n, err := strconv.Atoi(p); err == nil {
            c.Server.Port = n
        }
    }
    if h := os.Getenv("HOST"); h != "" {
        c.Server.Host = h
    }
    if t := os.Getenv("AUTH_TOKEN"); t != "" {
        c.Auth.Token = t
    }
    if t := os.Getenv("TIMEOUT"); t != "" {
        if n, err := strconv.Atoi(t); err == nil {
            c.Timeouts.Default = n
        }
    }
    if t := os.Getenv("ZAI_TOKEN"); t != "" {
        c.ZaiToken = t
    }
    if am := os.Getenv("AGENT_MODE"); am != "" {
        switch strings.ToLower(am) {
        case "1", "true", "yes", "on", "modern":
            c.AgentMode = true
        case "legacy":
            // Explicit opt-in to the old [ROLE: ...] rewrite shim.
            c.AgentMode = true
            c.AgentModeVariant = "legacy"
        case "0", "false", "no", "off":
            c.AgentMode = false
        }
    }
    // AGENT_MODE_VARIANT overrides the shim variant independently of the
    // AGENT_MODE on/off switch: "modern" (default) or "legacy".
    if v := os.Getenv("AGENT_MODE_VARIANT"); v != "" {
        switch strings.ToLower(v) {
        case "legacy":
            c.AgentModeVariant = "legacy"
        case "modern":
            c.AgentModeVariant = "modern"
        }
    }
    // AGENT_MODE_LEVEL gates the ToolParserLLM repair layer. Only "ultra"
    // enables it (together with AGENT_MODE). Any other value keeps it dormant.
    if v := os.Getenv("AGENT_MODE_LEVEL"); v != "" {
        c.AgentModeLevel = strings.ToLower(strings.TrimSpace(v))
    }
    // FORCE_CPU (or TOOLPARSER_LLM_FORCE_CPU alias) opts into degraded CPU
    // inference for ToolParserLLM ultra mode. Without it, ultra on CPU-only
    // aborts at startup.
    if v := os.Getenv("FORCE_CPU"); v != "" {
        switch strings.ToLower(strings.TrimSpace(v)) {
        case "1", "true", "yes", "on":
            c.ForceCPU = true
        case "0", "false", "no", "off":
            c.ForceCPU = false
        }
    }
    if v := os.Getenv("TOOLPARSER_LLM_FORCE_CPU"); v != "" {
        switch strings.ToLower(strings.TrimSpace(v)) {
        case "1", "true", "yes", "on":
            c.ForceCPU = true
        case "0", "false", "no", "off":
            c.ForceCPU = false
        }
    }
    if l := os.Getenv("LOG_LEVEL"); l != "" {
        c.Logging.Level = l
    }
    if f := os.Getenv("LOG_FORMAT"); f != "" {
        c.Logging.Format = f
    }
    if h := os.Getenv("STREAM_HOLDBACK"); h != "" {
        if n, err := strconv.Atoi(h); err == nil && n >= 0 {
            c.StreamHoldback = n
        }
    }
    // SYNC_MODE restores the legacy synchronous session flow (one chat
    // created per request). Used sessions are still deleted after use.
    if sm := os.Getenv("SYNC_MODE"); sm != "" {
        switch strings.ToLower(sm) {
        case "1", "true", "yes", "on":
            c.SyncMode = true
        case "0", "false", "no", "off":
            c.SyncMode = false
        }
    }
    if ps := os.Getenv("SESSION_POOL_SIZE"); ps != "" {
        if n, err := strconv.Atoi(ps); err == nil && n >= 1 {
            c.SessionPoolSize = n
        }
    }
    if at := os.Getenv("SESSION_ACQUIRE_TIMEOUT"); at != "" {
        if n, err := strconv.Atoi(at); err == nil && n >= 0 {
            c.SessionAcquireTimeout = n
        }
    }
    if rc := os.Getenv("SESSION_REUSE_COUNT"); rc != "" {
        if n, err := strconv.Atoi(rc); err == nil && n >= 1 {
            c.SessionReuseCount = n
        }
    }
    return c
}

var config = loadConfig()

// agentModern reports whether the modern agent-mode shim (XML-sectioned
// prompt, tolerant marker/payload parsing — see agent.go) is active.
func (c *Config) agentModern() bool {
    return c.AgentMode && !strings.EqualFold(c.AgentModeVariant, "legacy")
}

// agentLegacy reports whether the legacy agent-mode shim ([ROLE: ...]
// message rewriting — see transformMessagesForAgent) is active.
func (c *Config) agentLegacy() bool {
    return c.AgentMode && strings.EqualFold(c.AgentModeVariant, "legacy")
}

// UltraEnabled reports whether the ToolParserLLM intelligent repair layer
// is active. Strict opt-in: BOTH --agent-mode AND --agent-mode-level=ultra
// are required. Every other combination keeps ToolParserLLM completely
// dormant (no load, no warm-up, no reference). See agent_ultra.go /
// toolparser_llm.go.
func (c *Config) UltraEnabled() bool {
    return c.AgentMode && strings.EqualFold(strings.TrimSpace(c.AgentModeLevel), "ultra")
}

// UltraRepairEnabled is the package-level gate used by request handlers.
// It mirrors Config.UltraEnabled on the live config so call sites stay terse
// while the dormancy contract stays in one place.
func UltraRepairEnabled() bool {
    return config.UltraEnabled()
}

