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
│   ├── middleware/           # AuthMiddleware, OptionalAuthMiddleware, AdminMiddleware
│   └── router/               # Gin router + route registration
│
└── interfaces/
    └── http/                 # Gin handlers — thin, delegate to use cases
        ├── auth_handler.go
        ├── user_handler.go
        ├── video_handler.go
        └── product_handler.go
```

## Domain Entities

| Entity  | Key Fields |
|---------|-----------|
| User    | ID, Username, Name, Email, Password (bcrypt hash), **Role**, CreatedAt |
| Video   | ID, UserID, Title, Description, VideoURL, LikesCount, ViewsCount, CreatedAt |
| Product | ID, VideoID, Name, Price, ImageURL |
| Like    | ID, UserID, PostID (stores VideoID), CreatedAt |

## Role System

Two roles are defined in `domain/entities/user.go`:

| Role    | Constant        | Description |
|---------|-----------------|-------------|
| `user`  | `entities.RoleUser`  | Default role assigned on registration |
| `admin` | `entities.RoleAdmin` | Required for admin-only endpoints |

- The `role` field is stored in the `users` table (varchar 20, default `'user'`).
- Role is embedded in the JWT payload (`claims.Role`) so middleware can authorize
  without a DB lookup.
- To promote a user to admin, update the DB row directly or add an admin promotion
  endpoint behind another `AdminMiddleware`.

## Middleware

Three middleware functions live in `internal/infrastructure/middleware/auth.go`:

| Middleware | Behaviour |
|------------|-----------|
| `AuthMiddleware(secret)` | Rejects requests without a valid Bearer JWT (401). Sets `userID`, `username`, `role` in Gin context. |
| `OptionalAuthMiddleware(secret)` | Parses the JWT when present but never blocks. Sets `userID`/`username`/`role` to empty strings for anonymous callers. |
| `AdminMiddleware()` | Must follow `AuthMiddleware`. Rejects non-admin callers with 403. |

## Repository Interfaces

| Interface | Methods |
|-----------|---------|
| UserRepository | Create, FindAll, FindByID, FindByEmail, Exists, ExistsByEmail |
| VideoRepository | Create, FindByID, FindAll(cursor,limit), FindByUserID, Exists, IncrementLikes, DecrementLikes, IncrementViews |
| ProductRepository | Create, FindByID, FindByVideoID, Exists |
| LikeRepository | Create, Delete, FindByPostID, FindByUserID, Exists |

## API Endpoints

### Public (no auth required)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/register` | Register → JWT (role: user) |
| POST | `/api/v1/auth/login` | Login → JWT |
| GET  | `/api/v1/products/:id` | Get single product |

### Optional Auth (JWT decoded when present, never rejected)

| Method | Path | Description |
|--------|------|-------------|
| GET  | `/api/v1/feed?cursor=&limit=` | Cursor-paginated video feed |
| GET  | `/api/v1/videos/:id` | Get single video |
| GET  | `/api/v1/videos/:id/stats` | Views, likes, engagement rate |
| GET  | `/api/v1/videos/:id/product` | Product linked to a video |
| POST | `/api/v1/videos/:id/view` | Track view (auth → userID dedup; anon → IP dedup) |

### Protected (Bearer JWT required — any role)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/videos` | Upload video |
| POST | `/api/v1/videos/:id/like` | Like video (409 if duplicate) |
| POST | `/api/v1/videos/:id/unlike` | Unlike video (404 if not liked) |
| POST | `/api/v1/products` | Create product for a video |

### Admin (Bearer JWT required — role: admin)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/admin/users` | List all users |

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

## View Deduplication

An **in-process time-window map** prevents duplicate view counts within a 1-hour window.

```
key = viewerKey + ":" + videoID
window = 1 hour
```

- **Authenticated** viewer: `viewerKey = userID`
- **Anonymous** viewer: `viewerKey = "anon:" + clientIP`

Trade-off: resets on process restart; does not work across multiple instances.
Production path: replace with `SETEX "view:{viewerKey}:{videoID}" 3600 1` in Redis.

## Like System

- Duplicate-like prevention: `LikeRepository.Exists()` check before insert +
  `ON CONFLICT DO NOTHING` + unique index `(user_id, post_id)` at DB level.
- Counters updated atomically: `UPDATE videos SET likes_count = likes_count + 1`.
- Unlike decrements with `GREATEST(likes_count - 1, 0)` to prevent negative values.

## Database Design

| Table | Notable indexes |
|-------|----------------|
| users | `email` unique, `username` unique, `role` default `'user'` |
| videos | `created_at DESC` (cursor pagination), `user_id` |
| products | `video_id` unique (one product per video) |
| likes | composite unique `(user_id, post_id)`, `post_id` index |

## Swagger / API Docs

- Docs are generated by `swag init` and served at `GET /swagger/index.html`.
- Security scheme is defined in `main.go` via `@securityDefinitions.apikey BearerAuth`.
- All handlers use `godoc`-style Swagger annotations with `@Security BearerAuth` on
  protected endpoints and explicit `Authorization` header param docs on optional-auth endpoints.
- Run `swag init` after any change to handler annotations before rebuilding.

## Switching Between Memory and Postgres

**Postgres (default / Docker):**
```go
userRepo, likeRepo, videoRepo, productRepo := postgres.InitRepositories()
```

**Memory (local dev without Docker):**
```go
userRepo, likeRepo, videoRepo, productRepo := memory.InitRepositories(true)
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | `dev-secret-change-in-production` | JWT signing secret |
| `DB_HOST` | — (required for postgres) | Postgres host |
| `DB_PORT` | `5432` | Postgres port |
| `DB_USER` | `appuser` | Postgres user |
| `DB_PASSWORD` | `secret` | Postgres password |
| `DB_NAME` | `video_commerce` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |

## Extending the Project

- **New entity**: add to `domain/entities/`, define interface in `domain/repositories/`,
  implement in both `memory/` and `postgres/`, write use cases, handler, register routes,
  add to `postgres/connection.go` Migrate().
- **Admin promotion endpoint**: add `POST /api/v1/admin/users/:id/promote` behind
  `AuthMiddleware + AdminMiddleware`, update user role in DB.
- **Redis view dedup**: replace `globalViewTracker` in `track_view.go` with a Redis
  `SETEX` call; inject the Redis client via the use-case constructor.
- **Background view flush**: buffer view events in Redis; a worker goroutine drains
  the queue in batches and calls `IncrementViews` in bulk.
- **Read replica**: change `ReadDB` DSN in `NewDB()` — zero application code changes.
- **Rate limiting**: add a Gin middleware before protected routes using a token-bucket
  or sliding-window algorithm backed by Redis.

## Seed Data (Memory Mode)

When `initDefaultData=true`:
- 4 users (password for all: `12345678`):
  - `admin@example.com` — role: **admin**
  - `john@example.com` — role: user
  - `jane@example.com` — role: user
  - `bob@example.com` — role: user
- 3 videos across john, jane, bob
- 2 products linked to videos 1 and 2
- 2 pre-existing likes on video1
