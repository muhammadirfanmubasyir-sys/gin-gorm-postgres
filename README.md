# Gin GORM PostgreSQL

A production-ready REST API for user management built with Gin, GORM, and PostgreSQL.

## Tech Stack

- **Go 1.25.6**
- **Gin** v1.11.0 — HTTP framework
- **GORM** v1.31.1 — ORM
- **PostgreSQL 16** — Database
- **SQLite** (in-memory) — Test database
- **Testify** — Assertions

## Project Structure

```
.
├── main.go                 # Entry point, graceful shutdown, health check
├── config/
│   └── db.go               # Database connection, pooling, migration
├── controller/
│   └── user.go             # User CRUD handlers
├── dto/
│   ├── user.go             # Request/response DTOs
│   └── response.go         # Standard success/error response wrappers
├── middleware/
│   ├── cors.go             # CORS configuration
│   ├── logger.go           # Request logging
│   ├── ratelimit.go        # Per-IP rate limiting
│   ├── requestid.go        # UUID request ID
│   └── security.go         # Security headers
├── models/
│   └── user.go             # User model, password hashing
├── routes/
│   └── user.go             # Route registration
├── testutil/
│   ├── testutil.go         # Test helpers (DB, router, seed)
│   └── testutil_test.go    # Test helper tests
├── docker-compose.yaml     # PostgreSQL service
├── .env.example            # Environment variable template
└── go.mod
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose (for PostgreSQL)

### Setup

1. Clone the repository

```bash
git clone https://github.com/muhammadirfanmubasyir-sys/gin-gorm-postgres.git
cd gin-gorm-postgres
```

2. Copy environment file

```bash
cp .env.example .env
```

3. Start PostgreSQL

```bash
docker compose up -d
```

4. Run the application

```bash
go run main.go
```

The server starts on `http://localhost:8080`.

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `DB_HOST` | Yes | — | PostgreSQL host |
| `DB_PORT` | Yes | — | PostgreSQL port |
| `DB_USER` | Yes | — | PostgreSQL user |
| `DB_PASSWORD` | Yes | — | PostgreSQL password |
| `DB_NAME` | Yes | — | Database name |
| `DB_SSLMODE` | Yes | — | SSL mode (`disable`, `require`, etc.) |
| `SERVER_ADDR` | No | `:8080` | Server listen address |
| `DB_MAX_OPEN_CONNS` | No | `25` | Max open database connections |
| `DB_MAX_IDLE_CONNS` | No | `10` | Max idle database connections |
| `DB_CONN_MAX_LIFETIME` | No | `300` | Connection max lifetime in seconds |
| `CORS_ALLOWED_ORIGINS` | No | `*` | Comma-separated allowed origins |
| `RATE_LIMIT_PER_MINUTE` | No | `100` | Requests per minute per IP |

## API Endpoints

Base URL: `/api/v1`

### Health Check

```
GET /health
```

**Response `200 OK`**
```json
{
  "status": "ok",
  "timestamp": "2026-01-01T00:00:00Z"
}
```

### List Users

```
GET /api/v1/users
```

**Response `200 OK`**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Alice",
      "email": "alice@example.com"
    }
  ],
  "timestamp": "2026-01-01T00:00:00Z"
}
```

### Get User

```
GET /api/v1/users/:id
```

**Response `200 OK`**
```json
{
  "data": {
    "id": 1,
    "name": "Alice",
    "email": "alice@example.com"
  },
  "timestamp": "2026-01-01T00:00:00Z"
}
```

**Response `404 Not Found`**
```json
{
  "error": "user not found",
  "code": "USER_NOT_FOUND",
  "timestamp": "2026-01-01T00:00:00Z"
}
```

### Create User

```
POST /api/v1/users
Content-Type: application/json
```

**Request Body**
```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "password": "secret123"
}
```

**Validation Rules**
- `name` — required, cannot be empty
- `email` — required, valid email format
- `password` — required, minimum 6 characters

**Response `201 Created`**
```json
{
  "data": {
    "id": 1,
    "name": "Alice",
    "email": "alice@example.com"
  },
  "timestamp": "2026-01-01T00:00:00Z"
}
```

**Response `400 Bad Request`**
```json
{
  "error": "name is required",
  "code": "VALIDATION_ERROR",
  "timestamp": "2026-01-01T00:00:00Z"
}
```

### Update User

```
PUT /api/v1/users/:id
Content-Type: application/json
```

**Request Body**
```json
{
  "name": "Alice Updated",
  "email": "alice.new@example.com",
  "password": "newsecret123"
}
```

**Response `200 OK`**
```json
{
  "data": {
    "id": 1,
    "name": "Alice Updated",
    "email": "alice.new@example.com"
  },
  "timestamp": "2026-01-01T00:00:00Z"
}
```

### Delete User

```
DELETE /api/v1/users/:id
```

**Response `200 OK`**
```json
{
  "data": {
    "message": "user deleted successfully"
  },
  "timestamp": "2026-01-01T00:00:00Z"
}
```

## Error Codes

| Code | HTTP Status | Description |
|---|---|---|
| `VALIDATION_ERROR` | 400 | Input validation failed |
| `INVALID_REQUEST` | 400 | Malformed JSON body |
| `USER_NOT_FOUND` | 404 | User with given ID does not exist |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

## Middleware

Applied in order on every request:

1. **Request ID** — Injects `X-Request-ID` header (UUID)
2. **CORS** — Configurable cross-origin resource sharing
3. **Security Headers** — `X-Content-Type-Options`, `X-Frame-Options`, `X-XSS-Protection`, `Referrer-Policy`, `Cache-Control`
4. **Logger** — Request/response logging
5. **Rate Limit** — Per-IP sliding window limiter

## Testing

```bash
go test ./... -v
```

Tests use an in-memory SQLite database — no external dependencies needed.

```
controller  — 19 tests (CRUD flow, validation, error cases)
models      — 1 test  (table name)
routes      — 1 test  (endpoint registration)
testutil    — 14 tests (DB setup, seeding, isolation)
```

## Production Features

- Graceful shutdown on SIGINT/SIGTERM
- Configurable database connection pooling
- Request context propagation for cancellation
- Rate limiting with configurable limits
- Security headers
- Request ID tracking
- Health check endpoint
- Input validation with clear error messages
- Password hashing with bcrypt
