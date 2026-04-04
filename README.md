# CP Website API

> [!IMPORTANT]
> I'm an individual developer, and it's impossible for me to write all the code myself. I use AI for development assistance. I apologize if this doesn't align with your philosophy, but you don't have to use/contribute to it.

HTTP API for managing “CP” entries, tags, likes, and nested comments. Built with [Echo](https://echo.labstack.com/), [Ent](https://entgo.io/), SQLite, and [ZITADEL](https://zitadel.com/) OAuth2 token introspection.

**Release status: beta** — expect breaking changes; run through the [deployment checklist](#deployment) before exposing the service publicly.

## Requirements

- Go version **as specified in `go.mod`** (currently 1.26+)
- A ZITADEL application configured for **JWT private key** (client assertion) introspection
- SQLite (default) — database file `database.db` in the working directory

## Quick start

1. Create `secret.json` (see [Configuration](#configuration)).
2. Run migrations / create schema as you normally do for Ent (e.g. Atlas or `Schema.Create` during first deploy).
3. Start the server:

```bash
go run .
```

By default the server listens on **`0.0.0.0:8000`**. Override with `LISTEN_ADDR` or `PORT` (see below).

## Configuration

### `secret.json` (required)

Place next to the binary or set the working directory so the file is found. Shape:

| Field        | Description                          |
|-------------|--------------------------------------|
| `type`      | Application type (as in ZITADEL)     |
| `keyId`     | Key ID (`kid`) for the JWT assertion |
| `key`       | RSA private key PEM (string)         |
| `appId`     | Application ID                       |
| `clientId`  | OAuth2 client ID                     |

Introspection URL and JWT audience default to `https://auth.pdnode.com` unless overridden by environment variables.

### Environment variables

| Variable | Description |
|----------|-------------|
| `LISTEN_ADDR` | Full bind address. **Takes precedence** over `PORT`. Examples: `:8000`, `0.0.0.0:8000`, `127.0.0.1:3000`. |
| `PORT` | Port only (e.g. `8080`). Binds to **`0.0.0.0:PORT`**. Ignored if `LISTEN_ADDR` is set. |
| `HTTP_BODY_LIMIT` | Max request body size (Echo format, default `512K`). |
| `TRUSTED_PROXY_CIDRS` | Comma-separated CIDRs of reverse proxies. If **unset**, client IP is taken from the connection only (no `X-Forwarded-For`). If **set**, XFF is parsed only for requests whose direct peer is in these ranges. |
| `RATE_LIMIT_RPS` | Average requests per second per client IP (default `30`). |
| `RATE_LIMIT_BURST` | Burst size (default `60`). |
| `ZITADEL_INTROSPECT_URL` | OAuth2 introspection endpoint URL. |
| `ZITADEL_AUDIENCE` | `aud` claim for the client JWT assertion. |
| `ZITADEL_HTTP_TIMEOUT` | HTTP client timeout calling ZITADEL (e.g. `15s`). |
| `ZITADEL_INTROSPECT_CACHE_TTL` | Introspection result cache TTL (e.g. `30s`). Use `0s` or invalid value to disable. |

CORS allowed origins are currently defined in code (`main.go`); adjust for your frontend origins before production.

---

## API usage

### Base URL

All examples use `http://localhost:8000`. Replace with your deployed host and scheme (`https://…`).

### Authentication

Every route under **`/cp`** requires:

```http
Authorization: Bearer <access_token>
```

The access token is validated via ZITADEL **introspection** using a client assertion derived from `secret.json`.

### Response shape

**Success** (HTTP 2xx):

```json
{
  "status": "ok",
  "data": { }
}
```

**Error**:

```json
{
  "status": "error",
  "msg": "Human-readable message"
}
```

Server errors (5xx) return a generic message; details are logged server-side. Responses may include **`X-Request-ID`** for tracing.

### Public endpoints (no bearer token)

#### `GET /`

Liveness — short JSON message that the API process is running.

#### `GET /health`

Health check — returns `{"message":"OK"}` (not wired to the database).

---

### CP resources (authenticated)

Base path: **`/cp`**

#### `GET /cp`

List all CPs with tags and per-item `like_count`.

```bash
curl -sS -H "Authorization: Bearer $TOKEN" http://localhost:8000/cp
```

#### `POST /cp`

Create a CP. **Body (JSON):**

| Field | Type | Rules |
|-------|------|--------|
| `name` | string | Required, min length 1, unique |
| `category` | string | Required, length 2–20 |
| `link` | string | Optional |
| `tag_names` | string[] | Required, at least one tag (after trim, non-empty names) |

```bash
curl -sS -X POST http://localhost:8000/cp \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Example CP","category":"fiction","tag_names":["tag1","tag2"],"link":"https://example.com"}'
```

#### `GET /cp/:id`

Single CP with tags, `like_count`, and `is_liked` for the current user.

```bash
curl -sS -H "Authorization: Bearer $TOKEN" http://localhost:8000/cp/1234567890
```

`:id` is a numeric snowflake-style ID.

#### `PUT /cp/:id`

Update a CP (same JSON body as `POST /cp`). Caller must be the owner or an admin.

```bash
curl -sS -X PUT http://localhost:8000/cp/1234567890 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated","category":"fiction","tag_names":["tag1"]}'
```

#### `DELETE /cp/:id`

Delete a CP. Caller must be the owner or an admin.

```bash
curl -sS -X DELETE -H "Authorization: Bearer $TOKEN" http://localhost:8000/cp/1234567890
```

---

### Likes and comments (authenticated)

#### `POST /cp/:id/like`

Toggle like for the current user on CP `:id`. Response `data` includes `liked` (boolean, state after the request).

```bash
curl -sS -X POST -H "Authorization: Bearer $TOKEN" http://localhost:8000/cp/1234567890/like
```

#### `POST /cp/:id/comment`

Create a comment (or reply). **Body (JSON):**

| Field | Type | Rules |
|-------|------|--------|
| `content` | string | Required, min length 1 |
| `parent_id` | number | Optional; must reference an existing comment that belongs to **this** CP |

```bash
curl -sS -X POST http://localhost:8000/cp/1234567890/comment \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"Great!"}'
```

#### `GET /cp/:id/comments`

List top-level comments for CP `:id` with nested `children` and authors.

```bash
curl -sS -H "Authorization: Bearer $TOKEN" http://localhost:8000/cp/1234567890/comments
```

---

## Deployment

### Build a binary

```bash
go build -o cp-website .
```

Run from a directory that contains `secret.json` and is writable for `database.db` (or adjust the Ent DSN in code for production databases). Because this project uses `go-sqlite3`, build on the same OS family you deploy to (or cross-compile with a matching toolchain and libc).

### Process environment

Example for a VPS or bare-metal host listening behind Nginx/Caddy on the same machine (`127.0.0.1` upstream):

```bash
export LISTEN_ADDR="127.0.0.1:8000"
export TRUSTED_PROXY_CIDRS="127.0.0.0/8,::1/128"
export ZITADEL_INTROSPECT_CACHE_TTL="30s"
./cp-website
```

If the reverse proxy is on another host, set **`TRUSTED_PROXY_CIDRS`** to the subnet(s) from which your app sees those connections (not the public Internet). If nothing sits in front of the app, omit it and the server uses the TCP peer IP only.

You can bind on all interfaces and set only a port:

```bash
export PORT=8000
./cp-website
```

### Reverse proxy

- Terminate **TLS** at Nginx/Caddy (recommended). Keep `Secure` middleware HSTS at `0` in the app unless you serve HTTPS directly from Go.
- At the edge, **do not forward** client-supplied `X-Forwarded-For` blindly; have the proxy set or append the real client IP.
- Match **`TRUSTED_PROXY_CIDRS`** to the addresses that connect **to this API** (usually loopback if the proxy is local).

### Operations checklist

- [ ] `secret.json` is **not** in version control; restrict filesystem permissions (`chmod 600` or equivalent).
- [ ] `TRUSTED_PROXY_CIDRS` matches how traffic actually reaches the app (or unset if no trusted proxy).
- [ ] Regular backups of `database.db` (and a restore drill).
- [ ] CORS `AllowOrigins` in `main.go` matches your real frontend origin(s) for beta.
- [ ] Logs go somewhere durable if you care about incidents (`journald`, files, or a log shipper); correlate with `X-Request-ID`.

---

## License

This project is not open source for the time being, but will be open source under the MIT license at an appropriate time.