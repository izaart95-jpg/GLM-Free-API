# GLM Bridge — Z.AI Proxy API

OpenAI- **and Anthropic-compatible** API proxy for [chat.z.ai](https://chat.z.ai). Pure HTTP, no browser at runtime.

**Hosted:** `https://glm.cognix.sryze.cc/v1` — auth `Waguri` (Bearer or `x-api-key`), agent mode, supports `glm-5.2`/`glm-5.3-flash`. OpenAI endpoint recommended; Anthropic works but is primitive.

---

## Features

- **Dual protocol** — OpenAI `/v1/chat/completions` + Anthropic `/v1/messages` on one server; streaming (SSE, 5s keep-alive) + non-streaming
- **Pure HTTP** — no Playwright/Selenium; HTTP/2 + pooled conns; pure-Go SQLite (no CGO)
- **Vision** — images on both endpoints (OpenAI `image_url`, Anthropic `image`; URL or base64). Downloaded/decoded, uploaded to Z.AI file endpoint, attached as `files` array like official web client. **10 images / 50MB each**. Requires `ZAI_TOKEN` (guest → `401`)
- **Throwaway sessions (context-rot guard)** — every chat deleted on Z.AI (`DELETE /api/v1/chats/{id}`) when response completes; no history outlives a request
- **Async session pool (default)** — pre-made batch (`SESSION_POOL_SIZE`, default 5); requests grab instantly, consumed sessions deleted + replaced. `--sync-mode` restores legacy flow
- **Agent mode** — translates OpenAI tools/roles into user-only prompt Z.AI accepts, rewrites `<<<TOOL_CALL>>>` back to native `tool_calls`/`tool_use`. Modern XML shim (default, `<already_called>` anti-loop) + legacy `[ROLE: ...]` (opt-in)
- **Per-model feature resolution** — from Z.AI capabilities + overrides; `image_generation` always `false`; `reasoning_effort` (`high`/`max`) forwarded only if supported (forces `enable_thinking=true`)
- **Live model list** — from Z.AI `/api/models` (cached 5min, static fallback) with OpenAI-style `architecture`/modality
- **Token pool** — device tokens minted by `cmd/token-collector` into `tokens.sqlite`, consumed FIFO, deleted after use (up to 5 retries/captcha)
- **Device tokens minted in Go** — one brief browser visit registers device with Aliyun FeiLin; session reused to mint any number of tokens in-process (AES-128-CBC + MD5 seal, `VerifyCaptchaV3` → **T001**), verifiable **≥35min after browser closes** — no browser/Playwright/per-token page
- **Background captcha cache** — agent mode pre-generates captcha params async (2 cached, 75s TTL, pauses after 3min idle)
- **Graceful shutdown** — CTRL+C/SIGTERM drains in-flight (10s), deletes remaining pooled sessions before exit (second CTRL+C force-exits)

---

## Supported Models

Fetched live from `/api/models` (cached 5min). Server keeps newest down to `glm-4.7` (inclusive). Fallback list:

| Model ID | Notes |
|---|---|
| `glm-5.3-flash` | Lightweight flagship, premium quality, instant |
| `glm-5.3` | Flagship, coding + long-horizon |
| `glm-5.2` | Previous flagship |
| `GLM-5.1` | Older flagship |
| `GLM-5-Turbo` | Chat, coding, agentic |
| `GLM-5v-Turbo` | Vision model |
| `glm-4.7` | Classic high-performance |

Each `/v1/models` entry carries OpenAI-style `architecture` from capabilities (`info.meta.capabilities`); vision models (`"vision": true`) advertise image input:

```json
{
  "id": "GLM-5v-Turbo", "object": "model", "created": 1774521032,
  "owned_by": "z-ai", "display_name": "GLM-5V-Turbo",
  "description": "Vision model with evolved intelligence",
  "architecture": {
    "modality": "text+image->text",
    "input_modalities": ["text", "image"],
    "output_modalities": ["text"]
  }
}
```

Text-only: `"modality": "text->text"` / `"input_modalities": ["text"]`; `output_modalities` always `["text"]`.

> **Notes:**
> - No `model` in body → defaults to `glm-5`.
> - Guest (no `ZAI_TOKEN`) typically only allows `glm-5.3-flash` + `glm-4.7` — use `glm-4.7` for tokenless testing.
> - `/models` (plural) returns `{ models: [...], currentModel: "glm-5.2" }`.

---

## Getting `ZAI_TOKEN` (optional, recommended)

Z.AI JWT. Setting it skips guest init, unlocks all models + vision uploads.

1. Log in at **https://chat.z.ai**, DevTools (`F12`) → **Application → Local Storage → https://chat.z.ai** → copy key **`token`** (or `localStorage.getItem('token')`).
2. Export:
   ```bash
   export ZAI_TOKEN="paste-jwt-here"        # Linux/macOS
   $env:ZAI_TOKEN="paste-jwt-here"          # Windows PowerShell
   ```

---

## Getting Started

```bash
git clone https://github.com/izaart95-jpg/GLM-Free-API/ zai-api && cd zai-api
go mod init zai-api && go mod tidy          # repo ships without go.mod by design
npx playwright install-deps                 # optional: only for token collector
go run ./cmd/token-collector                # needs display (headless: xvfb-run -a)
go run .
```

Better perf: `go build -o zai-api -trimpath -gcflags="all=-l=4" -ldflags="-s -w" . && ./zai-api` (same for collector).

Startup banner shows health URL, endpoints, auth. Z.AI session initializes async — guest init failure retries on first chat.

---

## Configuration

### CLI Flags

| Flag | Default | Description |
|---|---|---|
| `--db-path` | `tokens.sqlite` | SQLite token DB path |
| `--verbose` | `false` | Verbose captcha/debug logging |
| `--agent-mode` | `false` | Enable agent mode (tools/role translation; starts captcha cache) |
| `--agent-mode-variant` | `modern` | `modern` (recommended) or `legacy` |
| `--agent-mode-level` | *(empty)* | empty (stock parser) or `ultra` (buffer + ToolParserLLM repair; requires `--agent-mode`) |
| `--force-cpu` | `false` | Allow ToolParserLLM ultra on CPU-only (degraded; GPU+SGLang required otherwise) |
| `--sync-mode` | `false` | Synchronous: one sticky chat vs pre-warmed pool (still deleted + rotated) |
| `--session-reuse-count` | `10` | Requests per chat before delete + replace (`1` = throwaway per request) |

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `3001` | HTTP port |
| `HOST` | `0.0.0.0` | Bind address |
| `AUTH_TOKEN` | `Waguri` | Bearer / `x-api-key` for client auth |
| `TIMEOUT` | `300000` | Request timeout (ms) |
| `ZAI_TOKEN` | *(empty)* | Hardcoded Z.AI JWT — skips guest init |
| `AGENT_MODE` | `false` | `1`/`true`/`yes`/`on`/`modern` → modern; `legacy` → legacy |
| `AGENT_MODE_VARIANT` | `modern` | Override variant (precedence over `AGENT_MODE`) |
| `AGENT_MODE_LEVEL` | *(empty)* | `ultra` enables ToolParserLLM repair (requires `AGENT_MODE`) |
| `FORCE_CPU` / `TOOLPARSER_LLM_FORCE_CPU` | `false` | Accept degraded CPU path for ultra on GPU-less hosts |
| `TOOLPARSER_LLM_CACHE_DIR` | OS cache / `./models/toolparser-llm` | Base-model + LoRA cache |
| `TOOLPARSER_LLM_LORA_URL` | public Kaggle zip | Override LoRA bundle download |
| `TOOLPARSER_LLM_LORA_PATH` | *(empty)* | Local LoRA dir (`adapter_model.safetensors` + `adapter_config.json`) or zip (air-gapped) |
| `TOOLPARSER_LLM_BASE_PATH` | *(empty)* | Local base-model snapshot (unset = resolve `Qwen/Qwen2.5-0.5B-Instruct` via SGLang/HF) |
| `TOOLPARSER_LLM_SGLANG_ADDR` | `127.0.0.1:30000` | SGLang address |
| `TOOLPARSER_LLM_LORA_SHA256` | *(empty)* | Pinned SHA256 for LoRA zip (mismatch aborts) |
| `TOOLPARSER_LLM_REFRESH` | *(empty)* | `1` = force re-download LoRA despite cache |
| `TOOLPARSER_LLM_PYTHON` | `python3` | Python for SGLang |
| `LOG_LEVEL` | `debug` | `debug` dumps every Z.AI req/resp, SSE lines, headers |
| `LOG_FORMAT` | `text` | Log format |
| `STREAM_HOLDBACK` | `24` | Runes held at live stream tail to absorb `edit_content` backtracks (`0` disables; issue #23) |
| `SYNC_MODE` | `false` | Sync session flow (one sticky reused chat, no pool) |
| `SESSION_POOL_SIZE` | `5` | Standing ready-session batch |
| `SESSION_ACQUIRE_TIMEOUT` | `10` | Seconds to wait for pooled session before creating directly (`0` = forever) |
| `SESSION_REUSE_COUNT` | `10` | Requests per session before delete + replace (`1` = throwaway; reuse looks human, avoids WAF bot flags) |
| `UPSTREAM_MIN_INTERVAL_MS` | *(random 200–500)* | Min gap between upstream requests — paces bursts to avoid Aliyun WAF IP blocks (issues #20, #41). One **shared gate across all transports** (Z.AI + Aliyun captcha) so alternating calls can't collapse the gap. Random default is a fingerprint-avoidance measure; `0` disables |
| `TOKEN_COUNT_TTL_MS` | `5000` | `/health` token count cache TTL (`<=0` keeps default) |

---

## Session Lifecycle — Reused Sessions & Pool

OpenAI clients are **stateless** — they re-send the whole conversation every request; bridge forwards it inside a `chat_id` materialized server-side. Leaving chats behind fills account with dead sessions; one-chat-per-message looks like bot behaviour to Aliyun WAF (IP block). So each session is **reused up to `SESSION_REUSE_COUNT` (default 10)** — one chat carrying several turns, like a human — then deleted via `DELETE /api/v1/chats/{id}` + replaced. Deletion is idempotent (`{"detail":"We could not find what you're looking for :/"}` = success). Reuse is safe: verified e2e that isolated single-turn payloads build no server-side history on a shared `chat_id`.

- Clients never see `chat_id`; state lives entirely in their own `messages` array — reuse/rotation invisible.

**Async (default):** startup pre-makes `SESSION_POOL_SIZE` (5). Requests acquire instantly; burst exhausts batch → waiters give up after `SESSION_ACQUIRE_TIMEOUT` (10s) and create directly. Each pooled session restocked until `SESSION_REUSE_COUNT`, then deleted + replaced. Z.AI chat IDs are client-generated UUIDs (chat materializes on first completion), so warmup is instant/local and unconsumed sessions never touch the account.

**Sync:** `--sync-mode` / `SYNC_MODE=true` — one sticky chat reused up to `SESSION_REUSE_COUNT`, then deleted + rotated; no pre-warming.

**Graceful shutdown:** CTRL+C/SIGTERM stops accepting, drains in-flight (10s), prints `clearing all remaining sessions...`, deletes every still-pooled session (or sticky sync session) on Z.AI, exits. Second CTRL+C force-exits.

**Observability:** `GET /status` → `{ "mode": "async", "reuse": true, "maxUses": 10, "gc_enabled": true, "size": 5, "ready": 5 }`.

---

## API Reference

### OpenAI-Compatible

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/v1/chat/completions` | ✅ | Chat completions (stream + non-stream) |
| `GET`  | `/v1/models` | ✅ | OpenAI-style model list |
| `GET`  | `/models` | ✅ | Compact `{ models, currentModel }` |

### Anthropic-Compatible

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/v1/messages` | ✅ | Messages API (stream + non-stream, tool use, thinking) |

Anthropic clients auth via `x-api-key` (same as `AUTH_TOKEN`); `anthropic-version` allowed through CORS.

#### `/v1/chat/completions` body

| Field | Type | Default | Notes |
|---|---|---|---|
| `model` | string | `glm-5` | Any ID from `/v1/models` |
| `messages` | array | *(required)* | OpenAI-style; `content` string or parts with `text`/`image_url` |
| `stream` | bool | `true` | SSE when true |
| `reasoning` | bool | *(per-model)* | Enables `enable_thinking` |
| `thinking` | object | *(per-model)* | `{"type":"enabled"}`/`{"type":"disabled"}` → `enable_thinking` |
| `reasoning_effort` | string | *(empty)* | `"high"`/`"max"` — forwarded only if supported; forces `enable_thinking=true` |
| `tools` | array | *(empty)* | OpenAI-style (requires agent mode) |
| `webSearch` / `search` | bool | *(per-model)* | Toggle WebSearch (sets `auto_web_search`) |
| `advancedSearch` | bool | `false` | Advanced Search (MCP variant) — adds `"mcp_servers":["advanced-search"]`; implies `webSearch=true` |

#### Image input (vision)

Both endpoints accept images. Bridge extracts every image, uploads to Z.AI (`POST /api/v1/files/`), attaches returned file objects as top-level `files` array — same payload as official web client. Text forwarded as usual; pixels ride in `files`.

- **OpenAI** — `image_url` parts (URL or base64 data URL):
  ```json
  {"role":"user","content":[
    {"type":"text","text":"What is in this image?"},
    {"type":"image_url","image_url":{"url":"https://example.com/photo.jpg"}},
    {"type":"image_url","image_url":{"url":"data:image/png;base64,iVBORw..."}}
  ]}
  ```
- **Anthropic** — `image` blocks (`base64` or `url` source):
  ```json
  {"role":"user","content":[
    {"type":"text","text":"What is in this image?"},
    {"type":"image","source":{"type":"base64","media_type":"image/png","data":"iVBORw..."}},
    {"type":"image","source":{"type":"url","url":"https://example.com/photo.jpg"}}
  ]}
  ```

**Limits:** max **10 images / 50MB each** per request — exceeding returns `400` before upload. Held in memory only (no temp files). Use vision-capable model (`architecture.input_modalities`); images on non-vision log warning + forward.

> ⚠️ **Requires real account (`ZAI_TOKEN`).** `/api/v1/files/` upload only accepts logged-in token. Guest: listing + text completions work, but uploads → `401` → request fails `400 {"message":"image processing failed: file upload unauthorized (401)"}`. Images uploaded per request (never cached/deduped), not deleted afterwards.

#### `/v1/messages` body

Standard Anthropic fields: `model`, `messages`, `system`, `max_tokens`, `temperature`, `top_p`, `stop_sequences`, `stream`, `thinking`, `tools`, `reasoning_effort`. `tool_use`/`tool_result` translate to/from OpenAI internally (tool calls require agent mode). Image blocks → OpenAI `image_url` parts via shared vision pipeline.

#### Request headers

| Header | Purpose |
|---|---|
| `Authorization: Bearer <AUTH_TOKEN>` | OpenAI-style auth |
| `x-api-key: <AUTH_TOKEN>` | Anthropic-style auth |
| `Include-All-Features: true` | (Only `POST /features`) send all server capabilities to `/completions` |

### Management

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET`  | `/` | ❌ | Redirects to `/health` |
| `GET`  | `/health` | ❌ | `200` if initialized else `503` — carries real-time `tokenCount` |
| `GET`  | `/status` | ❌ | Live session status (connected, userName, userId, feVersion, features, mode, sessionPool, waf) |
| `GET`  | `/admin/health` | ❌ | Same as `/health` |
| `GET`  | `/admin/stats` | ❌ | Mode, totalClients, totalRequests |
| `GET`  | `/admin/clients` | ❌ | Client list |
| `POST` | `/features` | ✅ | Per-model feature overrides |
| `GET`  | `/features` | ✅ | Inspect resolved features / stored states |
| `POST` | `/sqlite` | ✅ | Hot-swap active token DB without restart |
| `POST` | `/stop` | ✅ | Acknowledged no-op |

---

## `/health` — Real-Time Token Count

```json
{ "healthy": true, "mode": "direct", "tokenCount": 1500 }
```

- Served from TTL cache (`TOKEN_COUNT_TTL_MS`, default **5000**ms) so frequent probes don't turn into per-probe `SELECT COUNT(*)`.
- `-1` = unknown: no DB attached or query failed.
- DB swap invalidates cache immediately — next probe reflects new file.

## WAF Circuit Breaker (issue #41)

Sustained agent traffic can trip Aliyun WAF: **egress IP** gets block page on `POST /api/v2/chat/completions` even though `GET /` and `GET /api/models` keep working. Block is IP-wide, outlives any single request.

End-to-end handling:
- **Shared pacing gate** — every upstream transport funnels through *one* process-wide slot allocator; minimum gap holds regardless of transport. Previously independent pacing could collapse real gap to zero on interleaved calls.
- **Block detection** — block page (405/403 + Aliyun marker HTML) recognized, distinguished from genuine Z.AI JSON errors.
- **Fail fast** — while open, requests rejected `503` + `Retry-After` **before** pooled session drawn, vision upload, or **captcha burning a device token**. Chat deletions skipped too, so block generates no pointless upstream traffic.
- **Background recovery** — prober re-checks blocked endpoint after cooldown, **doubling backoff** (60s → 30min cap) while block persists; resumes automatically when lifted. Probe targets `POST /api/v2/chat/completions` (bare, unauthenticated, no captcha) — the endpoint WAF actually blocks.
- **Clean client errors** — small structured JSON never raw 3KB HTML:

```json
{ "error": { "type": "overloaded_error", "code": "waf_block",
    "message": "chat.z.ai has temporarily blocked this server's IP (Aliyun WAF). ...",
    "retryIn": "35s" } }
```

`GET /status` exposes live state:
```json
"waf": { "blocked": true, "state": 1, "retryIn": "43s", "consecutiveBlocks": 3 }
```

## `/sqlite` — Hot-Swap Token DB

Point running bridge at fresh `tokens.sqlite` **without restarting** (restart drops in-flight + pool):

```bash
curl -X POST http://localhost:3001/sqlite \
  -H "Authorization: Bearer $AUTH_TOKEN" -H "Content-Type: application/json" \
  -d '{"db_path": "/new/path/tokens.sqlite"}'
```

Response:
```json
{ "success": true, "message": "database swapped", "db_path": "/new/path/tokens.sqlite",
  "swapped_in": "35ms", "token_count": 1500 }
```

**Validation before touching live.** Rejected `400` + precise error if: doesn't exist, is directory, unreadable, not valid SQLite, or lacks expected `tokens` (id, token, batch) schema. Rejected candidate never disturbs active DB.

**Graceful by construction:**
1. Candidate validated + opened via private connection first — live swap is a single pointer exchange (tens of ms end-to-end incl. open + ping).
2. In-flight requests keep captured handle, finish normally against old file; retired handle closed after in-flight count drains.
3. Requests during swap **throttled** (wait microseconds for lock) then see new DB — nothing denied/failed mid-swap.
4. Concurrent `POST /sqlite` serialised; second gets "swap already in progress — retry shortly".

Repeated swaps to same path supported (collector may have rewritten file in place).

---

## `/features` — Per-Model Feature Configuration

Resolves **per model**: server capabilities → if `IncludeAll`, all except `reasoning_effort` (support flag, not feature) → stored overrides → `enable_thinking` defaults `true` unless overridden. `think` never included (only `enable_thinking` reaches request), `reasoning_effort` never stored (per-request only), **`image_generation` always forced `false`**.

**`GET /features`** — no query = all per-model states; `?model=<id>` = resolved `features`, stored `includeAll`/`overrides`, raw `capabilities`:
```bash
curl -H "Authorization: Bearer Waguri" http://localhost:3001/features
curl -H "Authorization: Bearer Waguri" "http://localhost:3001/features?model=glm-4.7"
```

**`POST /features`** — body must contain `model`; other keys are overrides (camelCase → snake_case):
```bash
curl -X POST http://localhost:3001/features \
  -H "Content-Type: application/json" -H "Authorization: Bearer Waguri" \
  -d '{"model":"glm-4.7","enable_thinking":false}'
```

| Special key | Behaviour |
|---|---|
| `reasoning` (bool) | stored as `enable_thinking` |
| `thinking` (bool or `{"type":"enabled"}`) | stored as `enable_thinking` |
| `image_generation` | ignored — always `false` |
| `think` | ignored — use `enable_thinking`/`reasoning`/`thinking` |
| `reasoning_effort` | not stored — per-request only |

To enable **all** capabilities for a model (testing), add header `Include-All-Features: true` to POST.

---

## Agent Mode

Z.AI's unofficial `/api/v2/chat/completions` only accepts `role="user"` — system/assistant/tool roles + OpenAI `tools`/`tool_calls` cause `INTERNAL_ERROR`. With agent mode (`--agent-mode` / `AGENT_MODE=true`), bridge folds roles & tools into user-only prompt and converts `<<<TOOL_CALL>>>` output back to native `tool_calls` deltas (OpenAI) / `tool_use` blocks (Anthropic) with `finish_reason="tool_calls"` / `stop_reason="tool_use"`. **Required for tool-calling** — without it, `tools` ignored. Hosted instance already runs it.

**Modern shim (default, recommended)** — folds conversation + tool contract into **one XML-sectioned user message**: `<system>` (output contract), `<tools>` (definitions), filled-in example from first tool, `<history_summary>` (older tool exchanges summarized — anti context-rot), `<recent>` (last 6 exchanges verbatim in `<tool_exchange>` groups), `<already_called>` (deduped calls already made — anti-loop, issue #31), `<current_task>`, `<output_rules>` (contract repeated verbatim at tail — recency). `<system>` contract deliberately strict + short (issue #44): literal block shown on own lines in head + tail, parser described as literal scanner recognizing only exact byte strings, derailment shapes (`TOOL_CALL_BLOCK`, renamed markers, fenced blocks) explicitly named discarded — paraphrasing = tool never runs. Parser keeps own tolerance (markers 2–4 angle brackets, ` ```json ` fences stripped, flat/alternate payloads normalized). Fallback parsers (Zhipu native tag dialect, envelope-less bare invocations, stream-end re-extraction) were **removed**: strict literal contract makes model emit exactly the block, so default parser alone suffices.

**Streamed tool calls** — with `stream: true`, calls stream live like OpenAI/Anthropic: once a block's `{"name": …, "arguments": {` arrived, client gets header delta (`index`, `id`, `type`, `function.name`), then incremental `function.arguments` fragments (`input_json_delta` on Anthropic) while model still writing JSON, closing with `finish_reason="tool_calls"` / `stop_reason="tool_use"`. No silent buffer window between markers. Non-streamable shapes (arguments as JSON-encoded string, flat payloads) emit one complete buffered call; fragments rune-safe so markers/multi-byte chars split across SSE chunks never garble or leak.

**Legacy shim (opt-in)** — `--agent-mode-variant=legacy` / `AGENT_MODE_VARIANT=legacy`: prepends system-prefix user message, rewrites every non-user message as `[ROLE: <role>] ...`, renders tools into strict `<<<TOOL_CALL>>> {"name":...,"arguments":{...}} <<<END_TOOL_CALL>>>` contract.

**Ultra repair (opt-in beta, issue #44)** — `--agent-mode --agent-mode-level=ultra`: when long-horizon sessions rot output and model emits paraphrased envelope (`TOOL_CALL_BLOCK`, `<tool_call>`, fenced/flat payloads) instead of canonical block, small Text-to-Tool-Call model (ToolParserLLM: Qwen2.5-0.5B-Instruct + LoRA adapter) repairs it. Pipeline: buffer stream → validate canonical → valid blocks pass untouched (model never consulted) → each malformed span isolated (prose excluded), de-noised, VAR-patched for long values (spans ≥150 chars, or ≥50 with newline, become `VARn` sentinels, rehydrated verbatim after), sent to model with canonical spec → output re-validated before splicing; on failure original fragment forwarded with warning, never silently dropped. Repair idempotent, per-span independent; every event logged (original/extracted/model output/decision). Dormant otherwise: without *both* flags nothing loaded/warmed/referenced; default + plain `--agent-mode` byte-for-byte unchanged.

Model serving: at ultra startup downloads base-model reference + LoRA bundle (public Kaggle `zoxoashouko/v28-flagship-toolparser`: `adapter_model.safetensors` + `adapter_config.json`), verifies integrity (non-empty zip, optional pinned SHA256, weights+config present + valid JSON), caches under `TOOLPARSER_LLM_CACHE_DIR`, loads base+LoRA via SGLang. GPU hosts serve automatically; CPU-only abort with actionable error unless `--force-cpu` / `FORCE_CPU=1` (degraded — slow, may exceed agent timeouts). Failures loud: bad checksum, incomplete bundle, or SGLang never ready all abort startup with what to check. `TOOLPARSER_LLM_LORA_PATH` / `TOOLPARSER_LLM_BASE_PATH` allow fully local/air-gapped.

Enabling agent mode also starts background captcha cache (2 pre-generated params, 75s TTL, pauses after 3min idle).

---

## Web Search & Advanced Search (issue #42)

Two separate things behind globe icon on chat.z.ai; bridge maps both to request flags:

**WebSearch** (`webSearch: true`, alias `search`) — model's built-in search. Internally sets **`auto_web_search: true`** in upstream `features`. Only key that actually toggles it: sibling `web_search` key is **always `false` upstream** — sending `true` there does nothing (server mirrors it back), which is why a request that "looked" enabled could silently do nothing. Bridge now always pins `features.web_search = false` and never sends `true` regardless of stored overrides.

**Advanced Search** (`advancedSearch: true`) — not a feature flag but an **MCP server**. On real web, toggling makes frontend add top-level field:

```json
{
  "model": "glm-5.3", "chat_id": "...", "messages": [...],
  "captcha_verify_param": "...",
  "mcp_servers": ["advanced-search"],
  "features": { "auto_web_search": true, ... }
}
```

`mcp_servers` sits at same payload level as `captcha_verify_param` (not inside `features`). When present, Z.AI runs searches through advanced-search MCP tools — internal `tool_call` events (`retrieve`, `open_url`, …) in SSE are that loop; answered server-side by Z.AI, unrelated to agent-mode tool calling.

Bridge reproduces exactly: `advancedSearch: true` adds `"mcp_servers": ["advanced-search"]` + auto-enables base web-search (same coupling as web UI, Advanced Search hides behind globe). Citations come inline as `【turn0search0】` markers + `[ref_id=...]` blocks model quotes from — `metadata.browser.search_result` on internal tool responses carries full source list (title, url, snippet, favicon).

**Agent-mode gate** — when agent mode on, both search toggles **force-disabled server-side**, even if requested. Z.AI serves WebSearch through own internal tool-call loop; mixing with agent shim's tool contract makes model emit internal markers (`retrieve`, `open_url`) as if they were caller's tools, or refuse to search. Request still succeeds — runs without server-side search. If you need both, run tool-calling requests without search flags and let client's own tools fetch pages.

Both flags accepted on Anthropic `/v1/messages` (`webSearch`, `search`, `advancedSearch`) and map to same upstream payload.

---

## Examples

> Self-hosted examples use `glm-4.7` (works **without** `ZAI_TOKEN`). Hosted: base URL `https://glm.cognix.sryze.cc/v1`, model `glm-5.2`/`glm-5.3-flash`, same auth.

```bash
# OpenAI — non-streaming
curl -X POST http://localhost:3001/v1/chat/completions \
  -H "Content-Type: application/json" -H "Authorization: Bearer Waguri" \
  -d '{"model":"glm-4.7","stream":false,"messages":[{"role":"user","content":"Hello, who are you?"}]}'

# OpenAI — streaming with deep thinking
curl -N -X POST http://localhost:3001/v1/chat/completions \
  -H "Content-Type: application/json" -H "Authorization: Bearer Waguri" \
  -d '{"model":"glm-4.7","stream":true,"reasoning":true,"messages":[{"role":"user","content":"Summarize today'\''s top AI news."}]}'

# Anthropic — streaming with thinking (primitive on hosted; OpenAI recommended)
curl -N -X POST http://localhost:3001/v1/messages \
  -H "Content-Type: application/json" -H "x-api-key: Waguri" -H "anthropic-version: 2023-06-01" \
  -d '{"model":"glm-4.7","max_tokens":4096,"stream":true,"thinking":{"type":"enabled"},"messages":[{"role":"user","content":"Explain quantum entanglement."}]}'
```

```python
# OpenAI SDK (hosted: base_url="https://glm.cognix.sryze.cc/v1", model="glm-5.3-flash"/"glm-5.2")
from openai import OpenAI
client = OpenAI(base_url="http://localhost:3001/v1", api_key="Waguri")
resp = client.chat.completions.create(model="glm-4.7", messages=[{"role": "user", "content": "Hello!"}])
print(resp.choices[0].message.content)

# Anthropic SDK
from anthropic import Anthropic
client = Anthropic(base_url="http://localhost:3001", api_key="Waguri")
resp = client.messages.create(model="glm-4.7", max_tokens=1024, messages=[{"role": "user", "content": "Hello!"}])
print(resp.content[0].text)
```

---

## How It Works

1. **Session init** — startup calls Z.AI `/api/v1/auths/guest` (+ `/api/v1/auths/`) for JWT, or uses `ZAI_TOKEN`; frontend version (`prod-fe-x.y.z`) scraped from homepage.
2. **Captcha** — per request, Aliyun `captcha_verify_param` generated in-memory: `InitCaptchaV3` → `certifyId`; `arg` via RC4-like permutation (KSA + PRGA over 64-byte state); `Track` JSON + custom `ali_hash` (16-byte-state hash), zlib-compress, base64, second RC4-like `encrypt` pass (different key); `VerifyCaptchaV3` with pooled device token → `securityToken`; final `{certifyId, isSign, sceneId, securityToken}` base64. Tokens FIFO from `tokens.sqlite`, deleted after use (up to 5 retries). Hard 90s timeout → `500`.
3. **Signature** — HMAC-SHA256 over `(sortedPayload | promptBase64 | timestamp)` with salted bucket key from `SALT_KEY` and `timestamp / 300000`.
4. **Streaming** — POST `/api/v2/chat/completions` with `stream:true`; parses SSE (`edit_content` + `edit_index`, `delta_content`, `content`, OpenAI-style deltas) re-emits as OpenAI/Anthropic SSE. Semantics mirror official `prod-fe` bundle: `edit_content` replaces from `edit_index` (**UTF-16 code-unit offset**, missing = 0 = full replace), `content` = full replace, `delta_content` = append. Deltas cut on rune boundaries with `STREAM_HOLDBACK` tail so backtracks never surface as garble (issue #23). Inline errors (HTTP 200 + `data.error`) surfaced; on `401` session re-initialises + retries once.
5. **Session GC** — each chat reused up to `SESSION_REUSE_COUNT` then deleted in background (async restocks pool; sync rotates sticky chat). Shutdown clears all pooled/sticky sessions. See [Session Lifecycle](#session-lifecycle--reused-sessions--pool).

---

## Token Collection (`cmd/token-collector`)

Seeds `tokens.sqlite` with device tokens, Bubble Tea TUI (progress bar, live logs, spinner).

**How it works (rewritten 2026-10-21).** Browser used **once, briefly**:
1. Opens `chat.z.ai`, boots Aliyun-SAF captcha SDK registering device with FeiLin (`Log1` → `InitCaptchaV3` → `Log2` → `Log3`) → **trusted device session**;
2. Collector grabs `Log1` response (`DeviceConfig`) + **one** genuine device token (decrypted payload = template);
3. Browser closed;
4. Every token then **minted in-process** (`cmd/token-collector/feilin.go`): refresh per-session fields (42/72/87) + per-token field 74 + trailing 93/94 probe-log entries, AES-128-CBC-seal with session key, assemble 5-field token.

Browser needed **once per session, not once per token** — one session mints arbitrary number, stays usable after browser gone (measured: 18/18 verifies over 35min no browser).

> **Browser must run headful.** Aliyun risk layer judges *client* issuing device service's `Log1`: headless rejected exactly like pure HTTP client, tokens won't verify. `Log1` request bytes byte-identical to SDK; fingerprint inside token irrelevant — token minted from foreign machine's payload still verifies. See `.assets/reports/aliyun-captcha-report.md` §13 + `.assets/reports/Fielin-report.md` §28.

```bash
go build -ldflags="-s -w" -trimpath -o token-collector ./cmd/token-collector
# modern CPUs, fully static: CGO_ENABLED=0 GOAMD64=v3 go build -ldflags="-s -w" -trimpath -o token-collector ./cmd/token-collector
./token-collector                        # interactive TUI prompt
./token-collector --tokens 750          # mint specific number
xvfb-run -a ./token-collector           # headless server: virtual display enough
```

| Flag | Default | Description |
|---|---|---|
| `--tokens` | `0` | Tokens to mint (0 = prompt; default 850, max 1500) |
| `--headed` | `true` | Deprecated/ignored — always headful (headless rejected) |
| `--session-max-age` | `15m` | Re-harvest session after this long (server-side lifetime) |
| `--export-session` | *(empty)* | Write session + payload template as JSON (for `.assets/scripts/live_validate.py --session-mint`) |
| `--no-tui` | `false` | Plain text output |

Implementation: pure-Go SQLite (WAL, 64MB page cache, 256MB mmap); one browser visit per session, tokens minted in Go; TUI ring buffer (1000 lines, scroll ↑/↓ g/G). Pure computation → nothing to batch/parallelise: one harvested session yields whole pool; only network-bound step (harvest) retried up to 3 times.

Cross-checks: `go test ./cmd/token-collector/` replays archived golden vectors (`DeviceConfig` decrypt + byte-exact token re-seal); `.assets/scripts/live_validate.py --session-mint N --session <collector export>` re-mints + re-verifies collector's own session against live service.

---

## Project Structure

Bridge core in `internal/zbridge` (mirrors DeepseekFreeAPI layout); `main.go` thin entry; blackbox tests in `tests/`, whitebox with package.

```
zai-api/
├── main.go                        # Thin entry -> zbridge.Run()
├── internal/zbridge/
│   ├── run.go                     # Run(), NewHandler(), graceful shutdown
│   ├── config.go                  # Config + env/flag loading
│   ├── types.go                   # Shared types + global state
│   ├── features.go                # Per-model feature resolution
│   ├── util.go                    # Logging, HTTP clients, cookie jar
│   ├── db.go                      # Hot-swappable token DB + count cache
│   ├── captcha.go                 # Aliyun captcha machinery + cache
│   ├── zai.go                     # Signature, session init, streaming, SSE parse
│   ├── format.go                  # OpenAI response/error formatting
│   ├── anthropic.go               # Anthropic endpoint + translation
│   ├── middleware.go              # CORS, auth, JSON, misc handlers
│   ├── models.go                  # Live model list + capabilities
│   ├── vision.go                  # Image extract/download/upload/files payload
│   ├── handlers.go                # /v1/chat/completions + admin handlers
│   ├── agent.go                   # Modern agent shim (XML-sectioned)
│   ├── agent_legacy.go            # Legacy [ROLE: ...] shim
│   ├── session_pool.go            # Throwaway sessions + async pool + GC
│   ├── ratelimit.go               # Shared upstream pacing gate (#20, #41)
│   ├── waf.go                     # Aliyun WAF circuit breaker (#41)
│   ├── waf_test.go                # Whitebox: breaker + probe
│   ├── testhooks.go               # Exported seams for tests/
│   ├── agent_test.go              # Whitebox: modern shim
│   ├── db_test.go                 # Whitebox: DB holder, count TTL, swap
│   ├── sse_garble_test.go         # Whitebox: SSE parser (#23)
│   ├── websearch_test.go          # Whitebox: search payload + agent gate (#42)
│   └── vision_test.go             # Whitebox: url->file-id rewrite
├── cmd/token-collector/           # Standalone: seeds tokens.sqlite (TUI)
├── tests/                         # Blackbox tests
│   ├── session_pool_test.go       # Pool mechanics, delete client, GC
│   ├── integration_test.go        # E2E HTTP garble regression (#23)
│   ├── toolcall_stream_test.go    # E2E streamed tool calls (agent)
│   ├── db_hotswap_test.go         # E2E /health tokenCount + /sqlite
│   ├── websearch_e2e_test.go      # E2E search flags -> payload (#42)
│   └── vision_test.go             # Vision e2e (upload + files), models, limits
└── README.md
```

Run suite: `go test ./...`

---

## Notes

- Default auth token (`Waguri`) is placeholder — set `AUTH_TOKEN` in production.
- Device tokens **consumed + deleted** after use — re-run collector to replenish (each captcha tries up to 5 tokens).
- `--verbose` controls only captcha subsystem logging; standard bridge/SSE gated by `LOG_LEVEL=debug`.
- Each request gets fresh `chat_id` — no server-side persistence; chat deleted right after response completes.

---

## License

MIT License. Use responsibly and in accordance with Z.AI's terms of service.
