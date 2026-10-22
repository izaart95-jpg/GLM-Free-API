# GLM Bridge — Z.AI Proxy API

OpenAI- & Anthropic-compatible API proxy for chat.z.ai. Pure HTTP, no browser automation at runtime.

**Hosted:** `https://glm.cognix.sryze.cc/v1` — auth `Waguri` (Bearer or `x-api-key`), agent mode, supports `glm-5.2`/`glm-5.3-flash`. OpenAI endpoint recommended.

## Features
- Dual protocol: OpenAI `/v1/chat/completions` + Anthropic `/v1/messages`; streaming (SSE, 5s keep-alive) + non-streaming
- Pure HTTP (HTTP/2 pooled), pure-Go SQLite (no CGO)
- Vision: images on both endpoints (URL/base64), uploaded to Z.AI file endpoint; 10 images / 50MB each; requires `ZAI_TOKEN`
- Throwaway sessions: chat deleted on Z.AI when response completes
- Async session pool (default 5 pre-made sessions; `--sync-mode` for legacy per-request flow)
- Agent mode: translates OpenAI tools/roles into user-only prompt, rewrites `<<<TOOL_CALL>>>` back to native `tool_calls`/`tool_use`. Modern XML shim (default) + legacy `[ROLE: ...]` shim
- Per-model feature resolution; `image_generation` always `false`; `reasoning_effort` (`high`/`max`) forwarded only if supported
- Live model list from Z.AI `/api/models` (cached 5 min)
- Token pool from `cmd/token-collector` (device tokens minted in Go, ≥35 min after browser closes)
- Background captcha cache (2 cached, 75s TTL)
- Graceful shutdown: drains in-flight (10s), deletes pooled sessions
- WAF circuit breaker (issue #41): shared pacing gate, block detection, fail-fast 503, background recovery with doubling backoff

## Supported Models
Fetched live (cached 5 min), keeps newest down to `glm-4.7`. Fallback list: `glm-5.3-flash`, `glm-5.3`, `glm-5.2`, `GLM-5.1`, `GLM-5-Turbo`, `GLM-5v-Turbo` (vision), `glm-4.7`.

Each `/v1/models` entry has OpenAI-style `architecture` (modality, input/output modalities). No `model` in body → defaults `glm-5`. Guest sessions typically only allow `glm-5.3-flash` and `glm-4.7`. `/models` returns `{ models, currentModel }`.

## Getting `ZAI_TOKEN` (recommended)
Log in at chat.z.ai → DevTools → Application → Local Storage → copy `token` key. Export as `ZAI_TOKEN`. Unlocks all models + vision uploads.

## Getting Started
```bash
git clone https://github.com/izaart95-jpg/GLM-Free-API/ zai-api && cd zai-api
go mod init zai-api && go mod tidy
npx playwright install-deps   # optional, for token collector
go run ./cmd/token-collector  # needs display (xvfb-run -a on headless)
go run .
```
Build: `go build -o zai-api -trimpath -gcflags="all=-l=4" -ldflags="-s -w" . && ./zai-api`

## Configuration

### CLI Flags
| Flag | Default | Description |
|---|---|---|
| `--db-path` | `tokens.sqlite` | SQLite token DB path |
| `--verbose` | `false` | Captcha debug logging |
| `--agent-mode` | `false` | Enable agent mode (tools/roles; starts captcha cache) |
| `--agent-mode-variant` | `modern` | `modern` or `legacy` shim |
| `--agent-mode-level` | *(empty)* | `ultra` (ToolParserLLM repair; requires agent mode) |
| `--force-cpu` | `false` | Allow CPU-only ultra repair (degraded) |
| `--sync-mode` | `false` | Reuse one sticky chat instead of pool |
| `--session-reuse-count` | `10` | Requests per chat before delete+replace |

### Environment Variables
| Variable | Default | Description |
|---|---|---|
| `PORT` | `3001` | HTTP port |
| `HOST` | `0.0.0.0` | Bind address |
| `AUTH_TOKEN` | `Waguri` | Client auth token |
| `TIMEOUT` | `300000` | Request timeout (ms) |
| `ZAI_TOKEN` | *(empty)* | Hardcoded Z.AI JWT |
| `AGENT_MODE` | `false` | Enable agent mode (`modern`/`legacy`) |
| `AGENT_MODE_VARIANT` | `modern` | Shim variant override |
| `AGENT_MODE_LEVEL` | *(empty)* | `ultra` repair level |
| `FORCE_CPU` | `false` | Accept degraded CPU path |
| `TOOLPARSER_LLM_*` | — | Cache dir, LoRA URL/path/SHA256, base path, SGLang addr, Python, refresh |
| `LOG_LEVEL` | `debug` | `debug` dumps all Z.AI req/resp |
| `LOG_FORMAT` | `text` | Log format |
| `STREAM_HOLDBACK` | `24` | Runes held back to absorb `edit_content` backtracks |
| `SYNC_MODE` | `false` | Synchronous session flow |
| `SESSION_POOL_SIZE` | `5` | Standing batch of sessions |
| `SESSION_ACQUIRE_TIMEOUT` | `10` | Seconds to wait for pooled session |
| `SESSION_REUSE_COUNT` | `10` | Requests per session |
| `UPSTREAM_MIN_INTERVAL_MS` | random 200–500 | Min gap between upstream requests (shared gate) |
| `TOKEN_COUNT_TTL_MS` | `5000` | `/health` token count cache TTL |

## Session Lifecycle
Clients are stateless (re-send full conversation). Each session reused up to `SESSION_REUSE_COUNT` (default 10) requests — one chat carrying several turns, like a human — then deleted via `DELETE /api/v1/chats/{chat_id}` and replaced. Deletion idempotent (`{"detail":"..."}` counts as success). `chat_id` never visible to clients.

**Async (default):** pre-makes 5 sessions; requests acquire instantly; burst exhausts → waiters give up after 10s and create directly. Each session restocked until reuse count hit, then deleted+replaced. Warmup instant/local (chat IDs are client-generated UUIDs).

**Sync:** one sticky chat reused, then deleted+rotated; no pre-warming.

**Shutdown:** stop accepting, drain 10s, delete all pooled/sticky sessions, exit. Second CTRL+C force-exits.

**Observability:** `GET /status` → `{ "mode", "reuse", "maxUses", "gc_enabled", "size", "ready", "waf" }`.

## API Reference

### OpenAI
| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/v1/chat/completions` | ✅ | Chat (streaming + non-streaming) |
| `GET` | `/v1/models` | ✅ | OpenAI-style model list |
| `GET` | `/models` | ✅ | `{ models, currentModel }` |

### Anthropic
| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/v1/messages` | ✅ | Messages API (streaming, tools, thinking) |

Anthropic auth via `x-api-key` (same as `AUTH_TOKEN`); `anthropic-version` allowed through CORS.

### `/v1/chat/completions` body
| Field | Type | Default | Notes |
|---|---|---|---|
| `model` | string | `glm-5` | Any from `/v1/models` |
| `messages` | array | required | String or parts with `text`/`image_url` |
| `stream` | bool | `true` | SSE when true |
| `reasoning` | bool | per-model | Enables `enable_thinking` |
| `thinking` | object | per-model | `{"type":"enabled"}`/`{"type":"disabled"}` |
| `reasoning_effort` | string | empty | `high`/`max`; forwarded only if supported; forces `enable_thinking=true` |
| `tools` | array | empty | OpenAI tools (requires agent mode) |
| `webSearch`/`search` | bool | per-model | Sets `auto_web_search` |
| `advancedSearch` | bool | `false` | Adds `"mcp_servers":["advanced-search"]`; implies `webSearch=true` |

### Image input
Both endpoints accept images; bridge uploads to Z.AI (`POST /api/v1/files/`) and attaches as top-level `files` array.
- OpenAI: `image_url` parts (URL or base64 data URL)
- Anthropic: `image` blocks (`base64` or `url` source)

Limits: max 10 images / 50MB each per request → `400` before upload. Held in memory only. Requires `ZAI_TOKEN` (guest uploads → `401`, request fails `400 {"message":"image processing failed: file upload unauthorized (401)"}`). Images uploaded per request, not cached/deleted.

### `/v1/messages` body
Standard Anthropic fields: `model`, `messages`, `system`, `max_tokens`, `temperature`, `top_p`, `stop_sequences`, `stream`, `thinking`, `tools`, `reasoning_effort`. `tool_use`/`tool_result` translate to/from OpenAI internally (requires agent mode). Images use shared vision pipeline.

### Request headers
| Header | Purpose |
|---|---|
| `Authorization: Bearer <AUTH_TOKEN>` | OpenAI auth |
| `x-api-key: <AUTH_TOKEN>` | Anthropic auth |
| `Include-All-Features: true` | (Only `POST /features`) send all capabilities |

### Management
| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/` | ❌ | Redirects to `/health` |
| `GET` | `/health` | ❌ | `200`/`503` + `tokenCount` |
| `GET` | `/status` | ❌ | Session status + WAF state |
| `GET` | `/admin/health` | ❌ | Same as `/health` |
| `GET` | `/admin/stats` | ❌ | Mode, totals |
| `GET` | `/admin/clients` | ❌ | Client list |
| `POST` | `/features` | ✅ | Per-model overrides |
| `GET` | `/features` | ✅ | Inspect features |
| `POST` | `/sqlite` | ✅ | Hot-swap token DB |
| `POST` | `/stop` | ✅ | No-op |

## `/health` — Token Count
Returns `{ "healthy": true, "mode": "direct", "tokenCount": 1500 }`. TTL cache (`TOKEN_COUNT_TTL_MS`, default 5000ms). `-1` = unknown (no DB / query failed). DB swap invalidates cache immediately.

## WAF Circuit Breaker (issue #41)
Aliyun WAF can block egress IP on `POST /api/v2/chat/completions` while `GET /` and `/api/models` keep working. Block is IP-wide, outlives requests.

- **Shared pacing gate** — one process-wide slot allocator across all upstream transports (Z.AI + Aliyun captcha)
- **Block detection** — Aliyun marker HTML (405/403) distinguished from genuine JSON errors
- **Fail fast** — 503 + `Retry-After` before session draw, vision upload, or captcha token burn; deletions skipped
- **Background recovery** — prober re-checks blocked endpoint after cooldown, doubling backoff (60s → 30min cap), resumes automatically. Probes `POST /api/v2/chat/completions` bare/unauthenticated
- **Clean errors** — `{"type":"overloaded_error","code":"waf_block","retryIn":"35s"}` with `Retry-After`

`/status` exposes `"waf": { "blocked", "state", "retryIn", "consecutiveBlocks" }`.

## `/sqlite` — Hot-Swap Token DB
```bash
curl -X POST http://localhost:3001/sqlite -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" -d '{"db_path": "/new/path/tokens.sqlite"}'
```
Response: `{ "success": true, "message", "db_path", "swapped_in", "token_count" }`.

Validation before live touch (400 on: missing, directory, unreadable, invalid SQLite, wrong schema). Swap: private connection first → single pointer exchange; in-flight keep old handle (closed after drain); requests during swap throttled; concurrent swaps serialised.

## `/features` — Per-Model Config
Resolve: server capabilities → `IncludeAll` (all except `reasoning_effort`) → stored overrides → `enable_thinking` defaults `true`. `think` never included, `reasoning_effort` never stored, `image_generation` always `false`.

`GET /features` (all) or `?model=<id>` (resolved + raw). `POST /features` body needs `model`; other keys are overrides (camelCase auto-converts).

Special keys: `reasoning` → `enable_thinking`; `thinking` → `enable_thinking`; `image_generation` ignored; `think` ignored; `reasoning_effort` not stored. Header `Include-All-Features: true` enables all capabilities.

## Agent Mode
Z.AI endpoint accepts only `role="user"` — system/assistant/tool roles and `tools` cause `INTERNAL_ERROR`. Agent mode folds roles/tools into user-only prompt, converts `<<<TOOL_CALL>>>` output back to native calls with `finish_reason="tool_calls"` / `stop_reason="tool_use"`. **Required for tool calling.**

**Modern shim (default):** XML-sectioned user message: `<system>` (output contract), `<tools>`, example call, `<history_summary>` (anti context-rot), `<recent>` (last 6 exchanges), `<already_called>` (anti-loop), `<current_task>`, `<output_rules>` (recency). Strict literal contract; parser tolerates 2–4 angle brackets, strips ` ```json ` fences, normalizes flat/alternate payloads. Fallback parsers removed.

**Streamed tool calls:** stream live like OpenAI/Anthropic — header delta on `{"name":…,"arguments":{`, then incremental fragments (`input_json_delta` on Anthropic), closing with `finish_reason`/`stop_reason`. Non-streamable shapes emit one buffered call; rune-safe fragments.

**Legacy shim (opt-in):** `--agent-mode-variant=legacy` — system-prefix user message, `[ROLE: <role>] ...` rewriting, strict `<<<TOOL_CALL>>> ... <<<END_TOOL_CALL>>>` contract.

**Ultra repair (opt-in beta, issue #44):** `--agent-mode --agent-mode-level=ultra` — small Text-to-Tool-Call model (Qwen2.5-0.5B + LoRA) repairs paraphrased envelopes. Pipeline: buffer → validate (valid pass untouched) → isolate malformed spans → de-noise → VAR-patch long values → model repair → re-validate → splice; on failure forward original with warning. Idempotent, per-span independent. Dormant without both flags.

Model serving: downloads base + LoRA bundle (`zoxoashouko/v28-flagship-toolparser`), verifies integrity, caches, loads via SGLang. GPU auto; CPU aborts unless `--force-cpu`. `TOOLPARSER_LLM_LORA_PATH`/`TOOLPARSER_LLM_BASE_PATH` for air-gapped.

Enabling agent mode starts background captcha cache (2 params, 75s TTL, pauses after 3min idle).

## Web Search & Advanced Search (issue #42)
**WebSearch** (`webSearch: true` / `search`) — sets `auto_web_search: true` upstream. `web_search` key always `false` upstream (bridge pins it).

**Advanced Search** (`advancedSearch: true`) — not a feature flag but an MCP server: adds `"mcp_servers": ["advanced-search"]` to completions payload (same level as `captcha_verify_param`), auto-enables base web search. Citations inline as `【turn0search0】` + `[ref_id=...]`; full source list in `metadata.browser.search_result`.

**Agent-mode gate:** both search toggles force-disabled server-side when agent mode on (mixing causes model to emit internal markers as caller tools or refuse). Request still succeeds.

Both accepted on `/v1/messages` too.

## Examples
```bash
# OpenAI non-streaming
curl -X POST http://localhost:3001/v1/chat/completions -H "Content-Type: application/json" \
  -H "Authorization: Bearer Waguri" \
  -d '{"model":"glm-4.7","stream":false,"messages":[{"role":"user","content":"Hello, who are you?"}]}'

# OpenAI streaming + thinking
curl -N -X POST http://localhost:3001/v1/chat/completions -H "Content-Type: application/json" \
  -H "Authorization: Bearer Waguri" \
  -d '{"model":"glm-4.7","stream":true,"reasoning":true,"messages":[{"role":"user","content":"Summarize AI news."}]}'

# Anthropic streaming + thinking
curl -N -X POST http://localhost:3001/v1/messages -H "Content-Type: application/json" \
  -H "x-api-key: Waguri" -H "anthropic-version: 2023-06-01" \
  -d '{"model":"glm-4.7","max_tokens":4096,"stream":true,"thinking":{"type":"enabled"},"messages":[{"role":"user","content":"Explain entanglement."}]}'
```
```python
from openai import OpenAI
client = OpenAI(base_url="http://localhost:3001/v1", api_key="Waguri")
resp = client.chat.completions.create(model="glm-4.7", messages=[{"role":"user","content":"Hello!"}])
print(resp.choices[0].message.content)

from anthropic import Anthropic
client = Anthropic(base_url="http://localhost:3001", api_key="Waguri")
resp = client.messages.create(model="glm-4.7", max_tokens=1024, messages=[{"role":"user","content":"Hello!"}])
print(resp.content[0].text)
```

## How It Works
1. **Session init** — `/api/v1/auths/guest` + `/api/v1/auths/` for JWT, or `ZAI_TOKEN`; frontend version scraped from homepage.
2. **Captcha** — per request generate `captcha_verify_param`: `InitCaptchaV3` → `certifyId`; `arg` via RC4-like cipher; `Track` JSON + `ali_hash`, zlib-compress, base64, second RC4 pass; `VerifyCaptchaV3` with pooled token → `securityToken`; final base64. Tokens FIFO from `tokens.sqlite`, deleted after use (5 retries). 90s timeout → `500`.
3. **Signature** — HMAC-SHA256 over `(sortedPayload | promptBase64 | timestamp)` with salted bucket key.
4. **Streaming** — POST `/api/v2/chat/completions` with `stream:true`; parses SSE (`edit_content` + `edit_index` (UTF-16 offset, missing=0=full replace), `delta_content`, `content`, OpenAI-style deltas) re-emitted as OpenAI/Anthropic SSE. Deltas cut on rune boundaries with `STREAM_HOLDBACK` tail. Inline errors surfaced; `401` re-inits and retries once.
5. **Session GC** — reuse up to `SESSION_REUSE_COUNT`, then delete (async restocks, sync rotates). Shutdown clears all.

## Token Collection (`cmd/token-collector`)
Seeds `tokens.sqlite` with device tokens; Bubble Tea TUI.

Browser used once per session: opens chat.z.ai, boots Aliyun-SAF SDK (registers device via `Log1`→`InitCaptchaV3`→`Log2`→`Log3`), grabs `DeviceConfig` + one genuine token as template, closes. Each token then minted in-process (refresh per-session fields 42/72/87, per-token 74 + 93/94 probe logs, AES-128-CBC seal, assemble 5-field token). One session mints arbitrary tokens, stays usable ≥35 min after browser closes (measured 18/18 over 35 min).

**Browser must run headful** — headless rejected by risk layer (tokens won't verify). `Log1` bytes byte-identical to SDK; fingerprint irrelevant (foreign-machine payload still verifies).

```bash
go build -ldflags="-s -w" -trimpath -o token-collector ./cmd/token-collector
./token-collector --tokens 750
xvfb-run -a ./token-collector   # headless server
```
Flags: `--tokens` (0=prompt, default 850, max 1500), `--headed` (deprecated/ignored), `--session-max-age` (15m), `--export-session`, `--no-tui`.

Implementation: pure-Go SQLite (WAL, 64MB page cache, 256MB mmap); one browser visit per session; TUI ring buffer (1000 lines). Minting pure computation — one session yields whole pool; only harvest retried up to 3x.

Cross-checks: `go test ./cmd/token-collector/` replays golden vectors; `.assets/scripts/live_validate.py --session-mint N --session <export>` re-mints/re-verifies live.

## Project Structure
Core in `internal/zbridge`; `main.go` thin entry; blackbox tests in `tests/`, whitebox with packages.
```
zai-api/
├── main.go
├── internal/zbridge/
│   ├── run.go, config.go, types.go, features.go, util.go, db.go
│   ├── captcha.go, zai.go, format.go, anthropic.go, middleware.go
│   ├── models.go, vision.go, handlers.go, agent.go, agent_legacy.go
│   ├── session_pool.go, ratelimit.go, waf.go, testhooks.go
│   └── *_test.go (waf, agent, db, sse_garble, websearch, vision)
├── cmd/token-collector/
├── tests/
│   ├── session_pool_test.go, integration_test.go, toolcall_stream_test.go
│   ├── db_hotswap_test.go, websearch_e2e_test.go, vision_test.go
└── README.md
```
Run: `go test ./...`

## Notes
- Default `Waguri` is placeholder — set `AUTH_TOKEN` in production.
- Device tokens consumed+deleted after use — re-run collector to replenish (5 tries per captcha).
- `--verbose` only captcha logging; bridge/SSE logging via `LOG_LEVEL=debug`.
- Each request gets fresh `chat_id`; chat deleted after response completes.

## License
MIT. Use responsibly per Z.AI ToS.
