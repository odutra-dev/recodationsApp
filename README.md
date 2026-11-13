# Recodations API (Golang Modular Monolith)

High-performance, modular monolith backend built in **Go (1.26)** using **Fiber v3**, following **Clean Architecture**, **SOLID**, **Repository Pattern**, **TDD**, **Testcontainers**, and advanced distributed systems patterns.

---

## 🏛️ Architectural Highlights & Implemented Patterns

1. **Modular Monolith**:
   - Strictly decoupled domain modules located under `internal/modules/` (`auth`, `users`, `rules`, `messages`, `inbox`, `audit`).
   - Shared cross-cutting infrastructure and domain building blocks under `internal/shared/` (`database`, `storage`, `cache`, `logger`, `middleware`, `presenter`).

2. **Result Pattern & Uniform Response Envelope**:
   - Every API endpoint returns the requested structured HTTP contract:
     ```json
     {
       "data": { ... },
       "message": "Operation completed successfully",
       "traceId": "c39fb678-...",
       "success": true,
       "timestamp": "2026-11-12T20:00:00Z",
       "detailsError": null,
       "path": "/api/v1/messages",
       "apiVersion": 1,
       "errors": null,
       "errorValids": null
     }
     ```

3. **Snowflake ID Generator**:
   - 64-bit distributed integer unique IDs generated using Twitter Snowflake (`pkg/snowflake`), eliminating auto-increment bottlenecks.

4. **Optimistic Locking (Concurrency Control)**:
   - Entities (`User`, `Message`, `InboxMessage`) feature a `version` column.
   - Updates are executed via `WHERE id = ? AND version = ?`, atomically incrementing `version`. Concurrency conflicts immediately yield an HTTP `409 Conflict` (`ErrOptimisticLockConflict`).

5. **State Transition Pattern**:
   - Finite State Machine validator (`internal/shared/domain/transition.go`) governing state evolutions.
   - Example for Messages: `DRAFT -> PUBLISHED -> ARCHIVED`. Illegal transitions reject with HTTP `422 Unprocessable Entity`.

6. **Inbox Pattern**:
   - Incoming events and transactional message dispatches are idempotently recorded in `inbox_messages` with an `event_id` idempotency key, status transitions (`PENDING -> PROCESSED` or `FAILED`), and retry count.

7. **Two-Level Cache (L1/L2)**:
   - **L1**: High-speed thread-safe in-memory cache with TTL.
   - **L2**: Distributed Redis cache.
   - Read strategy: Checks L1 -> Miss -> Checks L2 -> Hit -> Backfills L1 -> Returns.
   - Graceful Degradation: If Redis is unreachable, the system transparently serves traffic using L1 without interruption.

8. **Resilient Google Cloud Storage (`floci-gcp`)**:
   - Cloud Storage integration using the official `cloud.google.com/go/storage` SDK.
   - Local development and test environments emulated via **`floci-gcp`** on port `4588`.
   - **Circuit Breaker** (`sony/gobreaker`): Automatically trips when failure thresholds are breached to protect upstream dependencies.
   - **Exponential Backoff & Jitter Retry**: Handles transient network/5xx glitches.

9. **RBAC & Security**:
   - Granular permission rules in `rules` table (`messages:read`, `messages:write`, `messages:delete`, `users:manage`, `audit:read`).
   - N:N associations in `user_rules`.
   - Signed JWT authentication (HMAC-SHA256) and route-level authorization middleware (`RequireRule`).

10. **Audit Logging & Structured Logging**:
    - Structured logging via Go's native `log/slog` with request-scoped `TraceId`.
    - Mutating actions (POST, PUT, PATCH, DELETE) automatically persisted to the PostgreSQL `audit_logs` table.

11. **Interactive API Documentation (Swagger & Scalar)**:
    - OpenAPI specifications generated via `swaggo/swag`.
    - Classic Swagger JSON served at `/swagger/doc.json`.
    - Interactive **Scalar** UI served at `/docs`.

---

## 🚀 Quick Start

### 1. Prerequisites
- Go 1.26+
- Docker & Docker Compose

### 2. Run Infrastructure (PostgreSQL, floci-gcp, Redis)
```bash
docker compose up -d postgres floci-gcp redis
```

### 3. Run Application Locally
```bash
# Copy development environment variables
cp .env.example .env.dev

# Run with Go
go run ./cmd/api/main.go
```

The API will start at `http://localhost:3000`.
- **Swagger UI**: `http://localhost:3000/swagger`
- **Scalar Documentation**: `http://localhost:3000/docs`
- **Swagger JSON**: `http://localhost:3000/swagger/doc.json`
- **Health Checks**: `http://localhost:3000/healthz`, `http://localhost:3000/readyz`

---

## 🔑 Default Seeded Admin Credentials

On first boot, the system automatically migrates tables and seeds default rules and a root administrator:
- **Email**: `admin@recodations.com`
- **Password**: `Admin@123`

---

## 🧪 Testing

### Run All Unit Tests
```bash
go test -v ./...
```

### Run Testcontainers Integration Tests
```bash
go test -v ./test/...
```

### Generate Test Coverage
```bash
make test-coverage
```

---

## 📦 Project Layout

```
.
├── cmd/
│   └── api/
│       └── main.go               # Server initialization, dependency injection & routing
├── docs/                         # Generated Swagger / OpenAPI artifacts
├── internal/
│   ├── modules/
│   │   ├── audit/                # Audit trail domain, repository & controller
│   │   ├── auth/                 # Authentication, JWT issuance, login & register use cases
│   │   ├── inbox/                # Inbox pattern for idempotent event processing
│   │   ├── messages/             # Message domain with GCS upload, state transitions & optimistic lock
│   │   ├── rules/                # RBAC permission rules domain & repository
│   │   └── users/                # User management, password hashing & rules assignment
│   └── shared/
│       ├── domain/               # Result pattern, common errors & StateMachine transition pattern
│       ├── infrastructure/
│       │   ├── cache/            # L1 Memory + L2 Redis tiered caching
│       │   ├── database/         # PostgreSQL GORM connection, automigrations & seed
│       │   ├── logger/           # Slog structured logging with TraceID
│       │   ├── middleware/       # JWT auth, TraceID, RateLimit, Panic Recovery, Audit
│       │   └── storage/          # Resilient Google Cloud Storage (Circuit Breaker + Retry)
│       └── presenter/            # HTTP envelope ResponseHttp<T> & response helpers
├── legacy/
│   └── nestjs/                   # Legacy NestJS reference codebase
├── pkg/
│   └── snowflake/                # Distributed 64-bit Snowflake ID generator
├── test/
│   └── integration_test.go       # Testcontainers-go integration test suite
├── docker-compose.yml            # PostgreSQL, floci-gcp, Redis, API services
├── Dockerfile                    # Multi-stage production container build
├── Makefile                      # Dev tasks & shortcuts
└── README.md
```
