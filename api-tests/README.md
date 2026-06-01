# API Black-Box Tests

Black-Box integration tests for the Video Commerce & Social Feed API — written with Go test.

---

## Structure

```
api-tests/
├── helpers/
│   └── helpers.go          # HTTP client, shared types, seed helpers
├── auth/
│   └── auth_test.go        # Register + Login scenarios
├── videos/
│   └── videos_test.go      # POST /videos + GET /videos/:id
├── feed/
│   └── feed_test.go        # GET /feed (pagination, ordering, limits)
├── products/
│   └── products_test.go    # POST /products + GET /products/:id + GET /videos/:id/product
├── interactions/
│   └── interactions_test.go # Like, Unlike, View (dedup, auth, not-found)
├── analytics/
│   └── analytics_test.go   # GET /videos/:id/stats (engagement rate, edge cases)
├── admin/
│   └── admin_test.go       # GET /admin/users (RBAC scenarios)
├── go.mod
└── README.md
```

---

## Prerequisites

- Go 1.21+
- The API server must be running (default: `http://localhost:8080`)

---

## Running the Tests

### With Docker (recommended)

From the root of the main project:

```bash
# Full cycle: build + run + clean up
make integration-test

# Or step by step:
make integration-test-build   # build the test image
make integration-test-up      # run and attach to output
make integration-test-down    # clean up
```

### Locally (server already running)

```bash
# run server
make all

cd api-tests

# All packages
go test -v -count=1 ./...

# Or using the internal Makefile
make all
```

### Individual packages

```bash
cd api-tests

make auth          # go test -v ./auth/...
make videos        # go test -v ./videos/...
make feed          # go test -v ./feed/...
make products      # go test -v ./products/...
make interactions  # go test -v ./interactions/...
make analytics     # go test -v ./analytics/...
make admin         # go test -v ./admin/...
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `BASE_URL` | `http://localhost:8080` | API server address |
| `ADMIN_EMAIL` | `admin@example.com` | Admin account email |
| `ADMIN_PASSWORD` | `AdminPass123!` | Admin account password |

```bash
# Example with custom values
BASE_URL=http://staging.example.com:8080 \
ADMIN_EMAIL=myadmin@example.com \
ADMIN_PASSWORD=MyPass123! \
go test -v -count=1 ./...
```

---

## Notes

**Admin tests:** If no admin account is available, tests skip gracefully via `t.Skip()`. Receiving a 401/403 confirms the endpoint enforces access control correctly.

**Test isolation:** Each test creates its own unique users and videos using nanosecond timestamps, preventing interference between parallel test runs.

**`-count=1` flag:** Disables Go's test result cache, which is important for integration tests that depend on live server state.

---

## Scenarios Covered

| Package | Scenarios |
|---------|-----------|
| `auth` | Successful registration, duplicate email, required fields, password min length (8), username min/max length (3–50), correct/wrong login, unregistered email |
| `videos` | Create with/without auth, invalid token, required title/URL, title > 255 chars, description > 2000 chars, get by valid/nonexistent/malformed ID |
| `feed` | Anonymous, authenticated, custom limit, cursor pagination, invalid cursor, newest-first ordering, limit boundaries (0, 101, -1, `abc` → 400) |
| `products` | Create success, no auth, video not found, conflict (one product per video), required fields (name, price, image_url, video_id), name boundaries (empty → 400, 255 chars) |
| `interactions` | Like/unlike success, no auth, video not found, duplicate like (409), unlike without prior like (404), anonymous view counted, authenticated view counted, deduplication within time window |
| `analytics` | Zero stats on new video, after one view, after like+view, engagement rate = (likes/views)×100, like-without-view edge case (rate must be 0), anonymous access, authenticated access, video not found |
| `admin` | With admin token (user list + response shape), no token (401), regular user token (403), invalid token (401) |