# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

```bash
# Run the application (requires DB_HOST env var or use Docker)
go run main.go

# Install/update dependencies
go mod tidy

# Generate Swagger documentation (must run before docker-build)
swag init

# Run all tests
go test ./...

# Docker workflow (recommended)
make all          # down → swag-init → build → up
make docker-up    # start services
make docker-down  # stop services

# Install a new package
make go-get PKG=github.com/some/package
```

## Architecture Overview

The project follows **Clean Architecture** with strict layer separation.
Dependencies always point inward: Infrastructure → Application → Domain.

```
internal/
├── domain/                   # Enterprise rules — zero external imports
│   ├── entities/             # User, Video, Product, Like
│   └── repositories/         # Pure interfaces
│
├── application/              # Use-case layer — depends only on domain
│   ├── dto/
│   │   ├── request/          # Incoming DTOs with binding tags
│   │   └── response/         # Outgoing DTOs
│   └── usecases/
│       ├── auth/             # register, login, JWT helpers
│       ├── user/             # get users
│       ├── video/            # create/get, feed, like, unlike, view, stats
│       └── product/          # create/get product, get by video
│
├── infrastructure/           # Adapters — depends on application + domain
│   ├── database/
│   │   ├── memory/           # In-memory impls (dev / unit tests)
│   │   │   ├── init.go                   # InitRepositories(seedData bool)
│   │   │   ├── user_repository_impl.go
│   │   │   ├── like_repository_impl.go
│   │   │   ├── video_repository_impl.go
│   │   │   └── product_repository_impl.go
│   │   └── postgres/
│   │       ├── models/       # GORM models (UserModel, VideoModel, …)
│   │       ├── connection.go # DB struct, NewDB(), Migrate()
│   │       ├── init.go       # InitRepositories() — reads env vars
│   │       └── *_repository_impl.go
│   ├── middleware/           # AuthMiddleware (JWT)
│   └── router/               # Gin router + route registration
│
└── interfaces/
    └── http/                 # Gin handlers — thin, delegate to use cases
        ├── auth_handler.go
        ├── user_handler.go
        ├── video_handler.go
        └── product_handler.go
```

> **Note:** The legacy `post` domain (Post entity, PostRepository, post usecases,
> PostHandler, post routes, PostModel, post memory/postgres impls) has been
> **fully removed**. Only the `Like` entity remains and its `PostID` field is
> reused as `VideoID` for video likes (shared `likes` table pattern).

## Domain Entities

| Entity  | Key Fields |
|---------|-----------|
| User    | ID, Username, Name, Email, Password (bcrypt hash), CreatedAt |
| Video   | ID, UserID, Title, Description, VideoURL, LikesCount, ViewsCount, CreatedAt |
| Product | ID, VideoID, Name, Price, ImageURL |
| Like    | ID, UserID, PostID (stores VideoID), CreatedAt |

## Repository Interfaces

| Interface | Methods |
|-----------|---------|
| UserRepository | Create, FindAll, FindByID, FindByEmail, Exists, ExistsByEmail |
| VideoRepository | Create, FindByID, FindAll(cursor,limit), FindByUserID, Exists, IncrementLikes, DecrementLikes, IncrementViews |
| ProductRepository | Create, FindByID, FindByVideoID, Exists |
| LikeRepository | Create, Delete, FindByPostID, FindByUserID, Exists |

## API Endpoints

### Public (no auth)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login → JWT |
| GET  | `/api/v1/users` | List all users |
| GET  | `/api/v1/feed?cursor=&limit=` | Cursor-paginated video feed |
| GET  | `/api/v1/videos/:id` | Get single video |
| GET  | `/api/v1/videos/:id/stats` | Views, likes, engagement rate |
| GET  | `/api/v1/videos/:id/product` | Product linked to a video |
| GET  | `/api/v1/products/:id` | Get single product |

### Protected (Bearer JWT required)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/videos` | Upload video |
| POST | `/api/v1/videos/:id/like` | Like video (409 if duplicate) |
| POST | `/api/v1/videos/:id/unlike` | Unlike video (404 if not liked) |
| POST | `/api/v1/videos/:id/view` | Track view (dedup 1 h window) |
| POST | `/api/v1/products` | Create product for a video |

## Pagination Strategy — Cursor-Based

The feed uses `created_at` (RFC3339) as the cursor. Items are returned newest-first.

```
GET /api/v1/feed?limit=20
→ { "videos": [...20 items...], "next_cursor": "2024-01-15T10:00:00Z" }

GET /api/v1/feed?cursor=2024-01-15T10:00:00Z&limit=20
→ { "videos": [...next 20...], "next_cursor": "2024-01-14T08:00:00Z" }

# Empty next_cursor = no more pages
```

**Why not OFFSET?**
- OFFSET scans and discards N rows on every deep page — O(N) cost.
- Cursor uses a range predicate on an indexed column — O(log N) cost.
- Stable under concurrent inserts: new rows don't shift subsequent pages.

**Memory implementation** sorts the in-memory slice by `created_at DESC` and
applies the same exclusive-bound filter for API parity with the postgres impl.

## View Deduplication

An **in-process time-window map** (`sync.Mutex` + `map[string]viewRecord`) prevents
the same user from inflating view counts within a 1-hour window.

```
key = userID + ":" + videoID
window = 1 hour
```

Trade-off: resets on process restart; does not work across multiple instances.
Production path: replace with `SETEX "view:{userID}:{videoID}" 3600 1` in Redis.

## Like System

- Duplicate-like prevention: `LikeRepository.Exists()` check before insert +
  `ON CONFLICT DO NOTHING` + unique index `(user_id, post_id)` at DB level.
- Counters updated atomically: `UPDATE videos SET likes_count = likes_count + 1`
  (no read-modify-write race).
- Unlike decrements with `GREATEST(likes_count - 1, 0)` to prevent negative values.

## Database Design

| Table | Notable indexes |
|-------|----------------|
| users | `email` unique, `username` unique |
| videos | `created_at DESC` (cursor pagination), `user_id` |
| products | `video_id` unique (one product per video) |
| likes | composite unique `(user_id, post_id)`, `post_id` index |

## Switching Between Memory and Postgres

**Postgres (default / Docker):**
```go
// main.go
userRepo, likeRepo, videoRepo, productRepo := postgres.InitRepositories()
```

**Memory (local dev without Docker):**
```go
// main.go
userRepo, likeRepo, videoRepo, productRepo := memory.InitRepositories(true)
```

Set `DB_HOST` env var for postgres; leave unset to get a clear panic message.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | `dev-secret-change-in-production` | JWT signing secret |
| `DB_HOST` | — (required) | Postgres host |
| `DB_PORT` | `5432` | Postgres port |
| `DB_USER` | `appuser` | Postgres user |
| `DB_PASSWORD` | `secret` | Postgres password |
| `DB_NAME` | `video_commerce` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |

## Extending the Project

- **New entity**: add to `domain/entities/`, define interface in `domain/repositories/`,
  implement in both `memory/` and `postgres/`, write use cases, handler, register routes,
  add to `postgres/connection.go` Migrate().
- **Redis view dedup**: replace `globalViewTracker` in `track_view.go` with a Redis
  `SETEX` call; inject the Redis client via the use-case constructor.
- **Background view flush**: buffer view events in Redis; a worker goroutine drains
  the queue in batches and calls `IncrementViews` in bulk.
- **Read replica**: change `ReadDB` DSN in `NewDB()` — zero application code changes.
- **Rate limiting**: add a Gin middleware before protected routes using a token-bucket
  or sliding-window algorithm backed by Redis.

## Seed Data (Memory Mode)

When `initDefaultData=true`:
- 3 users: `john@example.com`, `jane@example.com`, `bob@example.com` (password: `12345678`)
- 3 videos across those users
- 2 products linked to videos 1 and 2
- 2 pre-existing likes on video1
