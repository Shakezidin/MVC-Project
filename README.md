# Bank Server

Production-grade banking backend service built in Go, designed for clean architecture and future MCP (Model Context Protocol) integration where each API endpoint becomes an MCP tool.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.23 |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Router | Gorilla Mux |
| Auth | JWT (golang-jwt/jwt/v5) |
| Logging | Zap (structured JSON) |
| Password Hashing | bcrypt |
| DB Driver | pgx/v5 |
| API Docs | Swagger/OpenAPI |
| Containerization | Docker + Docker Compose |

## Architecture

The project follows **Hexagonal / Clean Architecture** with clear layer separation:

```
┌─────────────────────────────────────────────────────────┐
│                      HTTP Layer                         │
│  Handlers → thin, parse request, call service, respond  │
├─────────────────────────────────────────────────────────┤
│                   Middleware Chain                       │
│  RequestID → Recovery → Logging → Security → CORS →     │
│  Timeout → RateLimit → [JWT Auth on protected routes]   │
├─────────────────────────────────────────────────────────┤
│                   Service Layer                          │
│  Business logic, caching strategy, ownership validation │
├─────────────────────────────────────────────────────────┤
│                  Repository Layer                        │
│  Interfaces + PostgreSQL implementations (pgx)          │
├─────────────────────────────────────────────────────────┤
│              Infrastructure Layer                        │
│  PostgreSQL │ Redis │ JWT │ Config │ Logger             │
└─────────────────────────────────────────────────────────┘
```

### Folder Structure

```
bank-server/
├── cmd/
│   ├── server/main.go      # REST API server entrypoint
│   ├── seed/main.go        # Database seeder (bcrypt passwords)
│   └── mcp/main.go         # MCP server entrypoint
├── internal/
│   ├── auth/               # JWT + bcrypt
│   ├── cache/              # Redis client + cache helpers
│   ├── config/             # Environment configuration
│   ├── handler/            # HTTP handlers
│   ├── observability/      # Zap logger + GCP PubSub integration
│   ├── middleware/         # HTTP middleware chain
│   ├── model/              # Domain models
│   ├── repository/         # Data access interfaces + impl
│   ├── response/           # Standard API response envelope
│   ├── router/             # Route registration
│   ├── service/            # Business logic
│   ├── utils/              # Context helpers, masking, pagination
│   └── validator/          # Request validation
├── migrations/             # SQL up/down migrations
├── scripts/                # Migration + seed scripts
├── docker/                 # Dockerfile
├── mcp/                    # MCP server components
│   ├── client/             # Bank API client
│   ├── config/             # MCP config
│   ├── tools/              # MCP tool implementations
│   └── types/              # MCP types
├── tests/                  # Integration test structure
├── docs/                   # Swagger documentation
└── docker-compose.yml
```

## Quick Start

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Make (optional)

### 1. Clone and Configure

```bash
cd bank-server
cp .env.example .env
```

### 2. Start Infrastructure

```bash
docker compose up -d postgres redis
```

### 3. Run Migrations

```bash
make migrate
# or on Windows with Git Bash:
bash scripts/migrate.sh
```

### 4. Seed Data

```bash
make seed-go
```

This creates test users with password `password123`:

| Email | Password |
|-------|----------|
| john.doe@example.com | password123 |
| jane.smith@example.com | password123 |

### 5. Run the Server

```bash
make run
# or
go run ./cmd/server
```

### 6. Run the MCP Server

The project includes an MCP (Model Context Protocol) server that exposes tools for interacting with the bank API:

```bash
make mcp-run
# or
go run ./cmd/mcp
```

### Full Docker Setup

```bash
# Start all services (postgres, redis, app)
docker compose up -d

# Seed data (after postgres is ready)
docker compose --profile seed run --rm seed
```

## API Documentation

### Swagger UI

Open [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Authentication

**Login** to obtain a JWT token:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john.doe@example.com","password":"password123"}'
```

Response:

```json
{
  "success": true,
  "message": "login successful",
  "request_id": "req-...",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "token_type": "Bearer"
  }
}
```

Use the token in subsequent requests:

```bash
export TOKEN="<your-jwt-token>"
```

### Endpoints

#### Health & Observability

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/health` | No | Full health check (DB + Redis) |
| GET | `/ready` | No | Readiness probe |
| GET | `/live` | No | Liveness probe |

#### Accounts

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/accounts` | Yes | List user accounts |
| GET | `/api/v1/accounts/balances` | Yes | All account balances |
| GET | `/api/v1/accounts/{accountId}/balance` | Yes | Single account balance |

```bash
# List accounts
curl http://localhost:8080/api/v1/accounts \
  -H "Authorization: Bearer $TOKEN"

# All balances
curl http://localhost:8080/api/v1/accounts/balances \
  -H "Authorization: Bearer $TOKEN"

# Single account balance
curl http://localhost:8080/api/v1/accounts/c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33/balance \
  -H "Authorization: Bearer $TOKEN"
```

#### Beneficiaries

```bash
curl http://localhost:8080/api/v1/beneficiaries \
  -H "Authorization: Bearer $TOKEN"
```

#### Transfer Modes

```bash
curl http://localhost:8080/api/v1/transfer-modes \
  -H "Authorization: Bearer $TOKEN"
```

### Response Format

**Success:**

```json
{
  "success": true,
  "message": "accounts fetched successfully",
  "request_id": "req-abc-123",
  "data": {}
}
```

**Error:**

```json
{
  "success": false,
  "message": "something went wrong",
  "request_id": "req-abc-123",
  "error": {
    "code": "INTERNAL_SERVER_ERROR"
  }
}
```

## Middleware Flow

Every request passes through this middleware chain (in order):

1. **Request ID** — generates/propagates `X-Request-ID` in context, logs, and responses
2. **Recovery** — catches panics, returns 500, logs stack trace
3. **Logging** — structured request completion logs (method, path, status, latency, user_id)
4. **Secure Headers** — X-Content-Type-Options, HSTS, CSP, etc.
5. **CORS** — cross-origin headers
6. **Timeout** — enforces configurable request deadline
7. **Rate Limit** — per-IP token bucket (configurable RPS/burst)
8. **JWT Auth** — on protected `/api/v1/*` routes (except login)

## Redis Caching

| Resource | Cache Key | TTL | Invalidation |
|----------|-----------|-----|--------------|
| Account list | `accounts:user:{userId}` | 5 min | `InvalidateAccountList()` |
| Beneficiary list | `beneficiaries:user:{userId}` | 5 min | `InvalidateBeneficiaryList()` |
| Transfer modes | `transfer_modes:all` | 1 hour | `InvalidateTransferModes()` |

Cache-aside pattern: on miss, fetch from PostgreSQL and populate cache. On Redis failure, services fall back to DB transparently.

## Configuration

All configuration via environment variables (see `.env.example`):

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | 8080 | HTTP port |
| `JWT_SECRET` | (required) | Min 32 characters |
| `DB_HOST` | localhost | PostgreSQL host |
| `REDIS_URL` | redis://localhost:6379/0 | Redis connection |
| `LOG_LEVEL` | info | debug, info, warn, error |
| `REQUEST_TIMEOUT` | 30s | Per-request timeout |
| `RATE_LIMIT_RPS` | 100 | Rate limit requests/sec |
| `GCP_PROJECT_ID` | tempBankLogs | GCP project ID for PubSub logging (optional) |
| `GCP_PUBSUB_TOPIC_ID` | bankServerLogs | GCP PubSub topic ID for logging (optional) |

## Make Commands

```bash
make build          # Build binary
make run            # Build and run
make test           # Run all tests
make test-unit      # Unit tests only
make test-integration  # Integration tests (requires running server)
make migrate        # Run SQL migrations
make seed-go        # Seed with bcrypt-hashed passwords
make docker-up      # Start Docker services
make docker-down    # Stop Docker services
make swagger        # Regenerate Swagger docs
make tidy           # go mod tidy
make fmt            # Format code with gofmt
make deps           # Install dependencies
make mcp-build      # Build MCP server
make mcp-run        # Build and run MCP server
make mcp-tidy       # Tidy MCP server modules
make mcp-clean      # Clean MCP build files
```

## Testing

```bash
# Unit tests
go test -v -short ./internal/...

# Integration tests (server must be running)
go test -v -tags=integration ./tests/integration/...
```

Example unit tests are provided in:
- `internal/auth/` — JWT and password hashing
- `internal/validator/` — request validation
- `internal/utils/` — account number masking
- `internal/service/` — auth service with mock repository

## Database Schema

| Table | Description |
|-------|-------------|
| `users` | User accounts with bcrypt password hashes |
| `bank_accounts` | Bank accounts linked to users |
| `balances` | Account balances (1:1 with bank_accounts) |
| `beneficiaries` | Saved transfer beneficiaries |
| `transfer_modes` | UPI, NEFT, RTGS, IMPS |

All tables use UUID primary keys with `created_at` / `updated_at` timestamps, foreign keys, and indexes.

## Security

- JWT with configurable expiry
- bcrypt password hashing (cost 12)
- Parameterized SQL queries (SQL injection protection)
- Request timeouts and rate limiting
- Panic recovery middleware
- Secure HTTP headers
- Environment-based secrets (never hardcoded)
- Account ownership validation on balance queries
- Masked account numbers in API responses

## Graceful Shutdown

The server listens for `SIGINT` / `SIGTERM`, stops accepting new connections, drains in-flight requests within `SERVER_SHUTDOWN_TIMEOUT` (default 30s), then closes DB and Redis connections.

## Future MCP Integration

This service is designed to be wrapped by an MCP server where each REST endpoint maps to an MCP tool:

| API Endpoint | Future MCP Tool |
|--------------|-----------------|
| `GET /api/v1/accounts` | `list_accounts` |
| `GET /api/v1/accounts/balances` | `get_all_balances` |
| `GET /api/v1/accounts/{id}/balance` | `get_account_balance` |
| `GET /api/v1/beneficiaries` | `list_beneficiaries` |
| `GET /api/v1/transfer-modes` | `list_transfer_modes` |
| `POST /api/v1/auth/login` | `authenticate` |

The clean handler/service separation and consistent response envelope make this mapping straightforward — each MCP tool calls the corresponding service method with the JWT-propagated user context.

## License

MIT
