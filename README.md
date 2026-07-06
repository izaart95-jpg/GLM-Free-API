# GLM Bridge — Z.AI Proxy API

An OpenAI-compatible API proxy for [chat.z.ai](https://chat.z.ai). Drop it in front of any OpenAI-compatible tool and start using Z.AI's GLM models without browser automation or complex setup.

---

## Features

- **OpenAI-compatible** — Works as a drop-in replacement for `/v1/chat/completions` and `/v1/models`
- **Pure HTTP** — No Playwright, no Selenium, no browser overhead
- **In-process captcha** — Aliyun CaptchaV3 verification handled entirely in-memory (no FIFO / named pipe)
- **Streaming + non-streaming** — Full SSE support with keep-alive ticks every 5s
- **Session management** — Per-client conversation threads via `X-Session-Id`, with 30-minute TTL (cleaned every 5 min)
- **Per-model feature resolution** — Features are resolved per-model from Z.AI server capabilities, with user overrides stored per-model. `image_generation` is **always forced to `false`**.
- **Token pool** — Device tokens stored in `tokens.sqlite`, consumed FIFO and removed after use (max 2 retries per request)
- **Live dashboard** — Status, features, and curl examples at `/`
- **Pure-Go SQLite** — Uses `modernc.org/sqlite` — no CGO required
- **HTTP/2 + pooled connections** — Optimised transport for both Aliyun and Z.AI endpoints

---

## Supported Models

Models are fetched live from Z.AI's `/api/models` (cached 5 min). The fallback list (used if Z.AI is unreachable) is:

| Model ID | Notes |
|---|---|
| `glm-5.2` | Flagship model, excels at coding and long-horizon tasks |
| `GLM-5.1` | Previous flagship model |
| `GLM-5-Turbo` | New model for chat, coding, and agentic tasks |
| `GLM-5v-Turbo` | Vision model with evolved intelligence |
| `glm-4.7` | Classic high-performance model |

> **Note:**
> - If you don't pass `model` in `/v1/chat/completions`, the server defaults to `glm-5`.
> - Z.AI's guest session (no `ZAI_TOKEN`) typically only allows `glm-4.7`. Use `glm-4.7` for tokenless testing.
> - `/models` (plural) returns `{ models: [...], currentModel: "glm-5.2" }` for clients that expect that shape.

---

## Getting `ZAI_TOKEN` (optional, but recommended)

`ZAI_TOKEN` is a Z.AI JWT. Setting it skips guest initialization and unlocks all models.

1. Go to **https://chat.z.ai** and log in.
2. Open browser **DevTools** (`F12` or `Ctrl+Shift+I`).
3. Navigate to **Application → Local Storage → https://chat.z.ai**.
4. Find the key named **`token`** and copy its value.
5. Export it before starting the server:

   ```bash
   # Linux / macOS
   export ZAI_TOKEN="paste-the-copied-jwt-here"

   # Windows PowerShell
   $env:ZAI_TOKEN="paste-the-copied-jwt-here"
   ```

   Or, in the DevTools **Console** tab, run:

   ```js
   localStorage.getItem('token')
   ```

   and copy the printed string.

---

## Getting Started

```bash
# 1. Clone the repo
git clone https://github.com/izaart95-jpg/GLM-Free-API/ zai-api
cd zai-api

# 2. Initialize the Go module
go mod init zai-api
go mod tidy

# 3. Generate the token database
go run init.go
# Recommended: build first for better performance and faster startup:
#   go build -o token-collector -ldflags="-s -w" init.go && ./token-collector

# 4. Start the server
go run main.go
# Recommended: build first for better performance and faster startup:
#   go build -o zai-api -ldflags="-s -w" main.go && ./zai-api
```

On startup, you'll see a banner with your dashboard URL and auth token. The Z.AI session is initialised asynchronously — if guest init fails, the first chat request will retry it.

---

## Configuration

### CLI Flags

| Flag | Default | Description |
|---|---|---|
| `--db-path` | `tokens.sqlite` | Path to the SQLite token database |
| `--verbose` | `false` | Enable verbose captcha/debug logging (`logError` / `logInfo` are silent unless this is set) |

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `3001` | HTTP server port |
| `HOST` | `0.0.0.0` | Bind address |
| `AUTH_TOKEN` | `Waguri` | Bearer token for client authentication |
| `TIMEOUT` | `300000` | Request timeout in milliseconds |
| `ZAI_TOKEN` | *(empty)* | Hardcoded Z.AI JWT — skips guest initialization |
| `LOG_LEVEL` | `debug` | Log level (`debug` dumps every Z.AI request/response, SSE lines, and headers) |
| `LOG_FORMAT` | `text` | Log format |

---

## API Reference

### OpenAI-Compatible

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/v1/chat/completions` | ✅ | Chat completions (streaming + non-streaming) |
| `GET`  | `/v1/models` | ✅ | OpenAI-style model list |
| `GET`  | `/models` | ✅ | Compact `{ models, currentModel }` shape |

#### `/v1/chat/completions` body

| Field | Type | Default | Notes |
|---|---|---|---|
| `model` | string | `glm-5` | Any model ID from `/v1/models` |
| `messages` | array | *(required)* | OpenAI-style message array |
| `stream` | bool | `true` | SSE stream when true |
| `webSearch` | bool | *(per-model)* | Enables `web_search` + `auto_web_search` for this request |
| `search` | bool | *(per-model)* | Alias for `webSearch` (used if `webSearch` is absent) |
| `deepThink` | bool | *(per-model)* | Enables `think` + `enable_thinking` for this request |

#### Request headers

| Header | Purpose |
|---|---|
| `Authorization: Bearer <AUTH_TOKEN>` | Required auth |
| `X-Session-Id: <id>` | Pin a conversation thread |
| `X-Fresh-Session: true` | Force a new session for this `X-Session-Id` |
| `Include-All-Features: true` | (Only for `POST /features`) Send all server capabilities to `/completions` |

### Management

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET`  | `/` | ❌ | HTML dashboard |
| `GET`  | `/status` | ❌ | Live session + feature status (JSON) |
| `GET`  | `/admin/health` | ❌ | Health check (`200` if initialised, else `503`) |
| `GET`  | `/admin/stats` | ❌ | Mode, active sessions, request count |
| `GET`  | `/admin/clients` | ❌ | Client list (always `[]` or one idle entry in direct mode) |
| `POST` | `/features` | ✅ | Per-model feature overrides (see below) |
| `GET`  | `/features` | ✅ | Inspect resolved features / stored states |
| `POST` | `/admin/session/clear` | ✅ | Clear all conversation histories |
| `POST` | `/admin/clients/<id>/clear` | ✅ | Clear history (clears all in direct mode) |
| `POST` | `/prompt` | ✅ | Non-OpenAI simple prompt endpoint (`{ prompt, model, ... }`) |
| `POST` | `/stop` | ✅ | Acknowledged no-op (returns `{ success: true }`) |
| `GET`  | `/inject.js` | ❌ | Returns `{"message":"Direct mode"}` |

---

## `/features` — Per-Model Feature Configuration

Features are resolved **per model** (not globally). The resolution logic is:

1. Start from the model's server capabilities.
2. If `Include-All-Features: true` header has been set for that model → include **all** capabilities.
   Otherwise include only `web_search`, `think`, `preview_mode` by default.
3. Apply stored user overrides (per-model).
4. **Always force `image_generation = false`** (overrides are ignored for this key).

### `GET /features`

- Without query: returns all per-model states:

  ```json
  { "states": { "glm-4.7": { "includeAll": false, "overrides": {...} } } }
  ```

- With `?model=glm-4.7`: returns the resolved feature map, the stored `includeAll` flag, stored `overrides`, and the model's raw `capabilities`:

  ```json
  {
    "model": "glm-4.7",
    "features": { "web_search": false, "think": false, "preview_mode": false, "image_generation": false },
    "includeAll": false,
    "overrides": {},
    "capabilities": { ... }
  }
  ```

### `POST /features`

Body **must** contain `model`. Any other key is treated as a feature override and is normalised to snake_case (e.g. `webSearch` → `web_search`, `deepThink` → `think`, `imageGen` → `image_generation`).

```bash
curl -X POST http://localhost:3001/features \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer Waguri" \
  -d '{"model":"glm-4.7","webSearch":true,"thinking":true}'
```

To enable **all** server capabilities for a model (e.g. for testing), send the header:

```bash
curl -X POST http://localhost:3001/features \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer Waguri" \
  -H "Include-All-Features: true" \
  -d '{"model":"glm-4.7"}'
```

Response:

```json
{
  "success": true,
  "model": "glm-4.7",
  "includeAll": false,
  "overrides": { "web_search": true, "think": true },
  "features": { "web_search": true, "think": true, "preview_mode": false, "image_generation": false }
}
```

> `imageGen` / `image_generation` is **always** `false` — sending it in the body has no effect.

---

## Examples

All examples use `glm-4.7` so they work **without** `ZAI_TOKEN`.

**Basic non-streaming request**

```bash
curl -X POST http://localhost:3001/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer Waguri" \
  -d '{
    "model": "glm-4.7",
    "stream": false,
    "messages": [{"role": "user", "content": "Hello, who are you?"}]
  }'
```

**Streaming (SSE)**

```bash
curl -N -X POST http://localhost:3001/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer Waguri" \
  -d '{
    "model": "glm-4.7",
    "stream": true,
    "messages": [{"role": "user", "content": "Write a haiku about Go."}]
  }'
```

**Web search + deep thinking**

```bash
curl -X POST http://localhost:3001/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer Waguri" \
  -d '{
    "model": "glm-4.7",
    "stream": true,
    "webSearch": true,
    "deepThink": true,
    "messages": [{"role": "user", "content": "Summarize today'\''s top AI news."}]
  }'
```

**Toggle per-model features**

```bash
curl -X POST http://localhost:3001/features \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer Waguri" \
  -d '{"model":"glm-4.7","webSearch":true,"thinking":true}'
```

**Python (OpenAI SDK)**

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:3001/v1",
    api_key="Waguri",
)

resp = client.chat.completions.create(
    model="glm-4.7",
    messages=[{"role": "user", "content": "Hello!"}],
)
print(resp.choices[0].message.content)
```

---

## Session Persistence

Pass `X-Session-Id` to pin a conversation thread across requests. Use `X-Fresh-Session: true` to start a new one. Sessions expire after **30 minutes** of inactivity (reaper runs every 5 minutes).

History is only appended server-side when `persistHistory` is enabled for the model via `POST /features` (e.g. `{"model":"glm-4.7","persistHistory":true}`).

```bash
curl -X POST http://localhost:3001/v1/chat/completions \
  -H "Authorization: Bearer Waguri" \
  -H "X-Session-Id: my-thread-1" \
  -H "Content-Type: application/json" \
  -d '{"model":"glm-4.7","messages":[{"role":"user","content":"My name is Alice."}]}'
```

---

## How It Works

1. **Guest token** — On startup, the server calls Z.AI's `/api/v1/auths/guest` (and `/api/v1/auths/`) for a session JWT, or uses `ZAI_TOKEN` if provided. The frontend version (`prod-fe-x.y.z`) is scraped from the Z.AI homepage.
2. **Captcha** — For each request, an Aliyun `captcha_verify_param` is generated **in-memory** (no FIFO, no named pipe):
   - `InitCaptchaV3` → obtain `certifyId`
   - Generate `arg` via RC4-like permutation cipher (KSA + PRGA over a 64-byte state)
   - Build a `Track` JSON, compute `ali_hash` (custom 16-byte-state hash), zlib-compress, base64-encode, then `encrypt` (second RC4-like pass with a different key)
   - `VerifyCaptchaV3` with a pooled device token → receive `securityToken`
   - Base64-encode the final `{ certifyId, isSign, sceneId, securityToken }` payload
   - Tokens are consumed FIFO from `tokens.sqlite` and deleted after use (up to 2 retries)
3. **Signature** — HMAC-SHA256 over `(sortedPayload | promptBase64 | timestamp)` with a salted bucket key derived from `SALT_KEY` and `timestamp / 300000`.
4. **Streaming** — POST to `/api/v2/chat/completions` with `stream: true`, parse SSE chunks (`edit_content`, `delta_content`, `content`, or OpenAI-style `choices[0].delta.content`), and forward as OpenAI-formatted SSE. Inline errors (HTTP 200 with `data.error`) are detected and surfaced as `api_error`. On `401`, the session is re-initialised and the request retried once.

---

## Project Structure

```
zai-api/
├── main.go          # HTTP server, captcha generation, Z.AI bridge, OpenAI shim
├── init.go          # Seeds tokens.sqlite with device tokens
├── tokens.sqlite    # Generated token pool (consumed at runtime)
├── go.mod
└── README.md
```

---

## Notes

- Device tokens are **consumed and deleted** after use. Re-run `init.go` to replenish the pool. Each request tries up to 2 tokens.
- The default auth token (`Waguri`) is a placeholder — set `AUTH_TOKEN` in production.
- `ZAI_TOKEN` bypasses guest initialization entirely. Without it, Z.AI's guest session typically only permits `glm-4.7`.
- `LOG_LEVEL=debug` dumps every Z.AI request body, response status/headers, and SSE lines — useful for troubleshooting.
- `image_generation` is **always `false`** and cannot be enabled via `/features` or per-request overrides.
- The captcha step has a hard 90-second timeout; if it fails, the request returns `500`.
- `--verbose` controls only the captcha subsystem's `logInfo` / `logError` output. Standard `log.*` calls (Z.AI bridge, SSE debug) are gated by `LOG_LEVEL=debug`.

---

## License

Provided as-is for educational and interoperability purposes. Use responsibly and in accordance with Z.AI's terms of service.
