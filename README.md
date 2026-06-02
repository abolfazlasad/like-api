# Like API — Video Commerce & Social Feed

A RESTful backend for a video-based social commerce platform, built with **Go** and **Gin** as part of a 48-hour technical challenge. The goal was to design a scalable, production-minded system that supports video feeds, social interactions, and shoppable products.

---

## Table of Contents

- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Makefile Commands](#makefile-commands)
- [Tests](#tests)
- [API Endpoints](#api-endpoints)
- [Swagger Docs](#swagger-docs)
- [Environment Variables](#environment-variables)
- [Branch Info](#branch-info)
- [TODO](#todo)

---

## Architecture

The project follows **Clean Architecture** principles, organizing code into four distinct layers:

```
Domain → Application → Interface → Infrastructures
```

| Layer | Responsibility |
|-------|---------------|
| `domain/` | Core entities and repository interfaces — zero external dependencies |
| `application/` | Use cases, DTOs, and business logic |
| `interfaces/http/` | Thin Gin handlers that delegate to use cases |
| `infrastructure/` | Concrete implementations: PostgreSQL, in-memory, middleware |

### Why Clean Architecture?

**Swappable dependencies.** The database, caching layer, or any external adapter can be replaced without touching business logic. A clear example in this project: switching between the PostgreSQL and in-memory implementations requires only a single environment variable (`DB_TYPE`). No use case code changes at all.

**Testability.** Because business rules live in the `domain` and `application` layers — with no direct dependency on Gin, GORM, or any framework — they can be tested with pure Go and simple in-memory fakes. This led to clean, fast unit tests with no database setup required.

**Separation of concerns.** Handlers are intentionally thin: they parse the request, call a use case, and map the output to an HTTP response. All real logic lives one layer deeper.

### Database Design

The PostgreSQL layer uses **two separate connection pools**: one for writes and one for reads. Write operations (inserts, updates, atomic counters) go through `writeDB`. Read operations (queries, existence checks) go through `readDB`.

```
writeDB  →  INSERT, UPDATE, DELETE
readDB   →  SELECT, EXISTS checks, feed queries
```

This design is forward-looking: when the system needs to scale reads, a read replica can be introduced by simply pointing `readDB` at the replica — no application code changes required.

Both PostgreSQL and an in-memory implementation are available. Switching between them is purely configuration (`DB_TYPE=postgres` or `DB_TYPE=memory`).

### Scalability & Statelessness

All request handling is **fully stateless**. No session state is stored on the server. Authentication relies entirely on self-contained JWTs. Because the application is stateless, horizontal scaling is straightforward — multiple instances can run behind a load balancer without sticky sessions or coordination.

### Trade-offs

The view deduplication window is currently an in-memory map (works for a single instance). In a multi-instance deployment this should be replaced with a Redis `SETEX` call. Like and view counters are written synchronously to PostgreSQL; under high traffic these should be buffered in Redis and flushed asynchronously via a background worker. The architecture is designed to accommodate both changes without modifying use case code.

---

## Tech Stack

- **Go** — primary language
- **Gin** — HTTP framework
- **PostgreSQL** — main database (GORM)
- **JWT** — authentication
- **Swagger** — API documentation
- **Docker / Docker Compose** — runtime and testing environment

---

## Project Structure

```
.
├── api-tests/                    # Black-Box integration tests (separate Go module)
│   ├── admin/                    # GET /admin/users tests
│   ├── analytics/                # GET /videos/:id/stats tests
│   ├── auth/                     # Register + Login tests
│   ├── feed/                     # GET /feed tests
│   ├── helpers/                  # HTTP client, shared types, seed helpers
│   ├── interactions/             # Like, Unlike, View tests
│   ├── products/                 # Products CRUD tests
│   ├── videos/                   # Videos CRUD tests
│   ├── go.mod
│   ├── Makefile
│   └── README.md
├── internal/
│   ├── application/              # Use Case layer
│   │   ├── dto/                  # Request / Response DTOs
│   │   └── usecases/             # auth, user, video, product
│   ├── domain/                   # Core — no external dependencies
│   │   ├── entities/             # User, Video, Product, Like
│   │   └── repositories/         # Pure interfaces
│   ├── infrastructure/           # Adapters
│   │   ├── database/
│   │   │   ├── memory/           # In-Memory implementation (dev/unit tests)
│   │   │   └── postgres/         # PostgreSQL implementation (GORM)
│   │   ├── middleware/           # Auth, OptionalAuth, Admin
│   │   └── router/               # Gin route registration
│   └── interfaces/http/          # Gin handlers (thin, delegate to use cases)
├── docs/                         # Generated by swag init
├── main.go
├── Makefile
├── docker-compose.yml
├── docker-compose.override.yml
└── docker-compose.test.yml
```

---

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- [`swag`](https://github.com/swaggo/swag) — for generating Swagger docs

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

---

## Quick Start

### With Docker (recommended)

```bash
# Clone and full setup (down → swag-init → build → up)
git clone https://github.com/abolfazlasad/like-api.git
cd like-api
make all
```

The server will be available at `http://localhost:8080`.

### Local (without Docker)

```bash
go mod tidy
swag init
cp .env.example .env
go run main.go
```
Note: The default database mode is PostgreSQL (`DB_TYPE=postgres`).
Update the `.env` file with your PostgreSQL credentials (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`).
If you want to use in-memory mode instead, change `DB_TYPE=memory`.


---

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make all` | down → swag-init → build → up |
| `make docker-up` | Start services |
| `make docker-down` | Stop services |
| `make docker-build` | Build images |
| `make swag-init` | Generate Swagger docs |
| `make test` | Run internal unit tests (without Docker) |
| `make go-tidy` | Update dependencies via Docker |
| `make go-get PKG=...` | Install a new package via Docker |
| `make integration-test` | build → up → run Black-Box tests → down |
| `make integration-test-build` | Build integration test image only |
| `make integration-test-up` | Run integration tests (attach to output) |
| `make integration-test-down` | Clean up after tests |

---

## Tests

### Unit Tests (internal)

Unit tests exist for the domain and application layers:

```bash
# Run without Docker
go test ./...

# Run with richgo (colored output)
make test
# equivalent: richgo test -v -cover -race -short ./...
```

Internal test files:

| Path | Coverage |
|------|----------|
| `internal/domain/entities/*_test.go` | Entity validation |
| `internal/application/usecases/auth/*_test.go` | JWT, Register, Login logic |
| `internal/application/usecases/video/create_video_test.go` | Video creation |

The existing tests are sample implementations that showcase the testing structure and approach. They are database- and technology-agnostic, making them easy to extend and adapt. Memory database is used for mocking dependencies during tests.


### Integration / Black-Box Tests

Black-Box tests live in `api-tests/` and work by firing real HTTP requests at a running server.

**Run with Docker (recommended):**
The database is started only for testing purposes and will be removed automatically after the tests complete:

```bash
make integration-test
```

**Run locally (server must be running at `http://localhost:8080`):**
Test data created during execution remains on the server):

```bash
make all

cd api-tests

# All packages
go test -v -count=1 ./...

# Individual packages
go test -v ./auth/...
go test -v ./videos/...
go test -v ./feed/...
go test -v ./products/...
go test -v ./interactions/...
go test -v ./analytics/...
go test -v ./admin/...
```

**Environment variables for tests:**

```bash
# Custom base URL (default: http://localhost:8080)
BASE_URL=http://staging.example.com:8080 go test -v ./...

# Admin credentials (for admin tests)
ADMIN_EMAIL=myadmin@example.com ADMIN_PASSWORD=MyPass123! go test -v ./admin/...
```

**Scenarios covered:**

| Package | Scenarios |
|---------|-----------|
| `auth` | Successful registration, duplicate email, required fields, password/username length limits, correct/wrong login |
| `videos` | Create with/without auth, invalid token, title/URL/description validation, get by valid/invalid/nonexistent ID |
| `feed` | Anonymous, authenticated, custom limit, cursor pagination, invalid cursor, newest-first ordering, limit boundaries |
| `products` | Create success, no auth, video not found, conflict (one product per video), required fields, name boundaries |
| `interactions` | Like/unlike success, no auth, video not found, duplicate like, unlike without prior like, view deduplication |
| `analytics` | Zero stats on new video, after view, after like+view, engagement rate formula, no-views edge case, anonymous access |
| `admin` | With admin token, no token, regular user token, invalid token |

> Admin tests skip gracefully if no admin account is available.

---

## API Endpoints

### Public (no auth required)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/register` | Register → JWT (role: user) |
| POST | `/api/v1/auth/login` | Login → JWT |
| GET  | `/api/v1/videos/:id` | Get single video |
| GET  | `/api/v1/videos/:id/stats` | Views, likes, engagement rate |
| GET  | `/api/v1/products/:id` | Get single product |
| GET  | `/api/v1/videos/:id/product` | Product linked to a video |

### Optional Auth (JWT decoded when present, never rejected)

| Method | Path | Description |
|--------|------|-------------|
| GET  | `/api/v1/feed?cursor=&limit=` | Cursor-paginated video feed |
| POST | `/api/v1/videos/:id/view` | Track view (deduplication + user identification for future personalisation) |


### Protected (Bearer JWT — any role)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/videos` | Upload video |
| POST | `/api/v1/videos/:id/like` | Like video (409 if duplicate) |
| POST | `/api/v1/videos/:id/unlike` | Unlike video (404 if not liked) |
| POST | `/api/v1/products` | Create product for a video |

### Admin (Bearer JWT — role: admin)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/admin/users` | List all users |

The only addition beyond the original specification is `GET /api/v1/admin/users`, which lists all registered users and is restricted to the `admin` role.

---

## Swagger Docs

After starting the server:

```
http://localhost:8080/swagger/index.html
```

To regenerate docs after changing handler annotations:

```bash
swag init
make docker-build && make docker-up
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | `dev-secret-change-in-production` | JWT signing secret |
| `DB_TYPE`      | `postgres`  | Database type (`postgres` or `memory`) |
| `DB_HOST` | — (required for postgres) | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `appuser` | Database user |
| `DB_PASSWORD` | `secret` | Database password |
| `DB_NAME` | `video_commerce` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `ADMIN_EMAIL` | `admin@example.com` | Initial admin email |
| `ADMIN_PASSWORD` | `admin` | Initial admin password |
| `ADMIN_NAME` | `Admin` | Admin display name |
| `ADMIN_USERNAME` | `admin` | Admin username |

---

## Branch Info

New changes and Clean Architecture refactoring are in the `develop` branch:

```bash
git checkout develop
```

---

## TODO

The following improvements were identified during development but could not be completed within the 48-hour time constraint. The architecture is designed so each of these can be added incrementally without restructuring existing code.

- [ ] **Error sanitization** — Internal errors (database errors, GORM messages, etc.) are currently propagated directly to the HTTP response in some paths. These should be mapped to safe, generic user-facing messages so that internal implementation details are never leaked to the client.

- [ ] **Redis caching for read queries** — Frequently accessed data (video metadata, feed pages, product lookups) hits the read PostgreSQL replica on every request. A Redis cache layer in front of the read repository implementations would dramatically reduce SQL load and response latency for popular content.

- [ ] **Redis-backed like and view counters** — High-frequency write operations like likes and views currently go directly to PostgreSQL. A better approach is to absorb these in Redis (atomic `INCR`) and flush them to the main database asynchronously via a queue and background worker. This decouples hot write paths from the primary DB and makes the system far more resilient under traffic spikes.

- [ ] **Improved feed mechanism** — The current feed is ordered strictly by `created_at`. A production feed should move toward a **recommendation-based model** that accounts for engagement signals (likes, views, recency, user affinity) rather than pure chronological order. This also addresses the cursor stability issue when content is ranked dynamically.