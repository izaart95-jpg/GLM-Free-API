# GLM Bridge — Z.AI Proxy API

An OpenAI API proxy for [chat.z.ai](https://chat.z.ai), letting you use Z.AI's models through any OpenAI-compatible tool.

---

## ✨ Features

- **OpenAI-Compatible API** — Drop-in replacement for `/v1/chat/completions` and `/v1/models`.
- **Direct HTTP Mode** — No browser automation, no Playwright, no Selenium. Pure Go HTTP.
- **In-Memory Captcha** — Aliyun CaptchaV3 verification computed in-process (no FIFO/named pipes).
- **Streaming (SSE) + Non-Streaming** — Full support for `stream: true` with keep-alive ticks.
- **Session Management** — Per-client sessions via `X-Session-Id` header with 30-minute TTL.
- **Feature Toggles** — Web search, deep thinking, image generation, preview mode, history persistence.
- **Token Pool** — Device tokens stored in `tokens.sqlite`, consumed FIFO and removed after use.
- **Live Dashboard** — Status, features, and curl examples served at `/`.
- **Pure-Go SQLite** — Uses `modernc.org/sqlite` (no CGO required).

---

## 📦 Supported Models

| Model ID         | Notes                       |
|------------------|-----------------------------|
| `glm-4.7`        |                             |
| `glm-5`          | Default                     |
| `GLM-5-Turbo`    |                             |
| `GLM-5v-Turbo`   | Vision variant              |
| `GLM-5.1`        | Latest                      |

---

## 🚀 Usage

```bash
# 1. Clone the repository
git clone https://github.com/izaart95-jpg/GLM-Free-API/ zai-api
cd zai-api

# 2. Initialize the Go module
go mod init zai-api

# 3. Pull dependencies
go mod tidy

# 4. Generate the tokens.sqlite database
go run init.go

# 5. After tokens.sqlite is generated, start the bridge
go run main.go

# 6. Enjoy 🎉
```

Once running, you'll see the startup banner with your dashboard URL and auth token.

### Building a Binary

```bash
go build -trimpath -ldflags="-s -w" -gcflags="all=-l=4" -o zai-bridge .
./zai-bridge --db-path tokens.sqlite --verbose
```

---

## ⚙️ Configuration

All settings are controlled via environment variables or CLI flags.

### CLI Flags

| Flag        | Default          | Description                          |
|-------------|------------------|--------------------------------------|
| `--db-path` | `tokens.sqlite`  | Path to the SQLite token database    |
| `--verbose` | `false`          | Enable verbose captcha/debug logging |

### Environment Variables

| Variable      | Default     | Description                                      |
|---------------|-------------|--------------------------------------------------|
| `PORT`        | `3001`      | HTTP server port                                 |
| `HOST`        | `0.0.0.0`   | HTTP server bind address                         |
| `AUTH_TOKEN`  | `Waguri`    | Bearer token clients must supply                 |
| `TIMEOUT`     | `300000`    | Default request timeout (ms)                     |
| `ZAI_TOKEN`   | *(empty)*   | Hardcoded Z.AI JWT (skips guest init)            |
| `LOG_LEVEL`   | `debug`     | Log level (`debug` enables Z.AI request dumps)   |
| `LOG_FORMAT`  | `text`      | Log format                                       |

---

## 🔌 API Endpoints

### OpenAI-Compatible

| Method | Path                    | Description                                  |
|--------|-------------------------|----------------------------------------------|
| `POST` | `/v1/chat/completions`  | Chat completion (streaming + non-streaming)  |
| `GET`  | `/v1/models`            | List available models                        |

### Management

| Method | Path                     | Description                                      |
|--------|--------------------------|--------------------------------------------------|
| `POST` | `/features`              | Toggle `webSearch`, `thinking`, `imageGen`, etc. |
| `POST` | `/admin/session/clear`   | Clear all conversation histories                 |
| `GET`  | `/status`                | Live session + feature status (JSON)             |
| `GET`  | `/admin/health`          | Health check (`200` / `503`)                     |
| `GET`  | `/`                      | HTML dashboard                                   |

---

## 🧪 Examples

### Non-Streaming

```bash
curl -X POST http://localhost:3001/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer Waguri" \
  -d '{
    "model": "glm-5",
    "stream": false,
    "messages": [
      {"role": "user", "content": "Hello, who are you?"}
    ]
  }'
```

### Streaming (SSE)

```bash
curl -N -X POST http://localhost:3001/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer Waguri" \
  -d '{
    "model": "glm-5",
    "stream": true,
    "messages": [
      {"role": "user", "content": "Write a haiku about Go."}
    ]
  }'
```

### Enable Web Search + Deep Thinking

```bash
curl -X POST http://localhost:3001/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer Waguri" \
  -d '{
    "model": "GLM-5.1",
    "stream": true,
    "webSearch": true,
    "deepThink": true,
    "messages": [
      {"role": "user", "content": "Summarize today'\''s top AI news."}
    ]
  }'
```

### Toggle Global Features

```bash
curl -X POST http://localhost:3001/features \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer Waguri" \
  -d '{"thinking": true, "webSearch": true, "imageGen": false}'
```

### Use with the OpenAI SDK (Python)

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:3001/v1",
    api_key="Waguri",
)

resp = client.chat.completions.create(
    model="glm-5",
    messages=[{"role": "user", "content": "Hello!"}],
)
print(resp.choices[0].message.content)
```

---

## 🧵 Session Persistence

Pass an `X-Session-Id` header to pin a conversation thread across requests. Send `X-Fresh-Session: true` to force a new chat. Sessions auto-expire after 30 minutes of inactivity.

```bash
curl -X POST http://localhost:3001/v1/chat/completions \
  -H "Authorization: Bearer Waguri" \
  -H "X-Session-Id: my-thread-1" \
  -H "Content-Type: application/json" \
  -d '{"model":"glm-5","messages":[{"role":"user","content":"My name is Alice."}]}'
```

---

## 🔐 How It Works

1. **Guest Token** — On startup, the server calls Z.AI's `/api/v1/auths/guest` to obtain a session JWT (or uses `ZAI_TOKEN` if provided).
2. **Captcha Param** — For each chat request, the server generates an Aliyun `captcha_verify_param` entirely in-memory:
   - `InitCaptchaV3` → get `certifyId`
   - Generate `arg` via a RC4-like permutation cipher
   - Compute `ali_hash`, zlib-compress, base64-encode, then `encrypt`
   - `VerifyCaptchaV3` with a pooled device token → receive `securityToken`
   - Base64-encode the final payload
3. **Signature** — HMAC-SHA256 over `(sortedPayload | promptBase64 | timestamp)` using a salted bucket key.
4. **Stream** — POST to `/api/v2/chat/completions` with `stream: true`, parse SSE chunks (`edit_content`, `delta_content`, `content`), and forward as OpenAI-formatted SSE.

---

## 📁 Project Structure

```
zai-api/
├── main.go            # HTTP server, captcha gen, Z.AI bridge, OpenAI shim
├── init.go            # Seed tokens.sqlite with device tokens
├── tokens.sqlite      # Generated token pool (consumed at runtime)
├── go.mod
└── README.md
```

---

## ⚠️ Notes

- Device tokens in `tokens.sqlite` are **consumed and deleted** after use. Re-run `init.go` to replenish.
- The default auth token (`Waguri`) is insecure — override with `AUTH_TOKEN`.
- If `ZAI_TOKEN` is set, guest initialization is skipped and the JWT is used directly.
- `LOG_LEVEL=debug` will dump every Z.AI request/response — useful for troubleshooting.

---

## 📜 License

Provided as-is for educational and interoperability purposes. Use responsibly and in accordance with Z.AI's terms of service.
