# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

```bash
# Run the application
go run main.go

# Install/update dependencies
go mod tidy

# Generate Swagger documentation
swag init

# Run all tests
go test ./...

# Docker workflow (recommended)
make all          # down → swag-init → build → up
make docker-up    # start services
make docker-down  # stop services
```

## Architecture Overview

The project follows **Clean Architecture** with strict layer separation.
Dependencies always point inward: Infrastructure → Application → Domain.

```
internal/
├── domain/                   # Enterprise rules — no external imports allowed
│   ├── entities/             # User, Post, Video, Product, Like
│   └── repositories/         # Pure interfaces (UserRepository, VideoRepository …)
│
├── application/              # Use-case layer — depends only on domain
│   ├── dto/
│   │   ├── request/          # Incoming request DTOs (binding tags)
│   │   └── response/         # Outgoing response DTOs
│   └── usecases/
│       ├── auth/             # register, login, JWT helpers
│       ├── like/             # like/unlike posts, get liked posts
│       ├── post/             # legacy post CRUD
│       ├── user/             # get users
│       ├── video/            # create/get video, feed, like, unlike, view, stats
│       └── product/          # create/get product, get by video
│
├── infrastructure/           # Adapters — depends on application + domain
│   ├── database/
│   │   ├── memory/           # In-memory repo impls (dev/test only)
│   │   └── postgres/
│   │       ├── models/       # GORM models (UserModel, VideoModel, ProductModel …)
│   │       ├── connection.go # DB struct, NewDB(), Migrate()
│   │       ├── init.go       # InitRepositories() — reads env vars
│   │       └── *_repository_impl.go
│   ├── middleware/           # AuthMiddleware (JWT)
│   └── router/               # Gin router, route registration
│
└── interfaces/
    └── http/                 # Gin handlers — thin; delegate to use cases
```

## Entities

| Entity    | Key Fields                                                        |
|-----------|-------------------------------------------------------------------|
| User      | ID, Username, Name, Email, Password (bcrypt), CreatedAt          |
| Post      | ID, UserID, Title, Content, Likes                                 |
| Video     | ID, UserID, Title, Description, VideoURL, LikesCount, ViewsCount, CreatedAt |
| Product   | ID, VideoID, Name, Price, ImageURL                                |
| Like      | ID, UserID, PostID (also used for video likes), CreatedAt        |

## API Endpoints

### Public
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/register` | Register |
| POST | `/api/v1/auth/login` | Login → JWT |
| GET  | `/api/v1/feed?cursor=&limit=` | Cursor-paginated video feed |
| GET  | `/api/v1/videos/:id` | Get video |
| GET  | `/api/v1/videos/:id/stats` | Views, likes, engagement rate |
| GET  | `/api/v1/videos/:id/product` | Product linked to a video |
| GET  | `/api/v1/products/:id` | Get product |
| GET  | `/api/v1/posts` | Legacy posts |
| GET  | `/api/v1/users` | All users |

### Protected (Bearer JWT)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/videos` | Upload video |
| POST | `/api/v1/videos/:id/like` | Like video |
| POST | `/api/v1/videos/:id/unlike` | Unlike video |
| POST | `/api/v1/videos/:id/view` | Track view (dedup 1 h window) |
| POST | `/api/v1/products` | Create product |

## Pagination Strategy

The feed uses **cursor-based pagination** keyed on `created_at` (RFC3339).

```
GET /api/v1/feed?limit=20
→ { videos: [...], next_cursor: "2024-01-15T10:00:00Z" }

GET /api/v1/feed?cursor=2024-01-15T10:00:00Z&limit=20
→ { videos: [...], next_cursor: "2024-01-14T08:00:00Z" }
```

Why cursor instead of OFFSET:
- No row re-scan on deep pages (OFFSET n scans and discards n rows).
- Stable under concurrent inserts — new rows do not shift subsequent pages.
- The composite index on `created_at DESC` makes each page a single range scan.

## View Deduplication

Views are deduplicated with an **in-process time-window map** (1 h per user+video pair).

Trade-off: simple and zero-dependency, but resets on restart and does not work
across multiple instances. Production upgrade path: replace with a Redis
`SETEX "view:{userID}:{videoID}" 3600 1` check.

## Like Counters

`likes_count` / `views_count` are incremented/decremented with atomic SQL
expressions (`UPDATE … SET col = col + 1`) rather than read-modify-write,
preventing lost-update races without application-level locking.

## Database Design Decisions

- `videos.created_at` is indexed → efficient cursor-based feed queries.
- `products.video_id` has a unique index → one product per video enforced at DB level.
- `likes(user_id, post_id)` composite unique index → duplicate likes blocked at DB level.
- `likes.post_id` is also used for video IDs (shared table); separate tables are the
  clean long-term solution but share-table avoids schema churn during early development.

## Extending the Project

- **New entity**: add to `domain/entities/`, define interface in `domain/repositories/`,
  implement in `infrastructure/database/postgres/`, write use cases, handler, and register routes.
- **Redis caching**: wrap repository reads in a cache layer implementing the same interface.
- **Background view flush**: replace `IncrementViews` call with a Redis enqueue;
  a worker goroutine (or separate service) drains the queue in batches.
- **Read replica**: change `ReadDB` DSN in `NewDB()`; no application code changes needed.