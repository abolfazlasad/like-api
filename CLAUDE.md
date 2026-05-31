# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Running the Application
```bash
# Run the application
go run main.go

# Install/update dependencies
go mod tidy

# Generate Swagger documentation
swag init
```

### Testing
The project currently doesn't have test files. To add tests:
```bash
# Create test files following Go naming convention (*_test.go)
# Run tests
go test ./...
```

## Code Architecture

This project follows Clean Architecture principles with separation of concerns across layers:

### Layer Structure
```
internal/
├── application/          # Use cases and business logic
│   ├── dto/             # Data Transfer Objects (request/response)
│   └── usecases/        # Business use cases organized by entity (like, post, user)
├── domain/              # Enterprise business logic
│   ├── entities/        # Core business objects (User, Post, Like)
│   └── repositories/    # Interfaces for data access
├── infrastructure/      # External concerns
│   ├── database/        # Database implementations (currently memory-based)
│   └── router/          # HTTP router setup
└── interfaces/          # Interface adapters
    └── http/            # HTTP handlers (Gin controllers)
```

### Key Components
- **Entities** (`internal/domain/entities/`): Core domain objects representing business concepts
  - `User`: ID, Username, Name, Email
  - `Post`: ID, UserID, Content, Likes count
  - `Like`: ID, UserID, PostID, CreatedAt

- **Repositories** (`internal/domain/repositories/`): Interfaces defining data access contracts
  - `UserRepository`: Exists, FindByID, Save
  - `PostRepository`: Exists, FindByID, FindAll, FindByUserID, Update
  - `LikeRepository`: Exists, Create, FindByUserIDAndPostID, FindByPostID, FindByUserID

- **Use Cases** (`/internal/application/usecases/`): Application-specific business rules
  - `LikePostUseCase`: Like/unlike posts with validation
  - `GetUserLikedPostsUseCase`: Retrieve posts liked by a user
  - `GetPostLikesUseCase`: Get like count for a post
  - User/Post CRUD operations

- **HTTP Handlers** (`/internal/interfaces/http/`): Gin-based HTTP controllers
  - `UserHandler`: Handle user-related endpoints
  - `PostHandler`: Handle post-related endpoints
  - `LikeHandler`: Handle like/unlike endpoints

- **Router** (`/internal/infrastructure/router/`): Gin router setup
  - Configures routes and middleware
  - Sets up Swagger documentation routes

- **Database Implementations** (`/internal/infrastructure/database/memory/`): In-memory repository implementations
  - Thread-safe implementations using mutexes
  - Pre-populated with sample data on initialization

### Data Flow
1. HTTP request arrives at Gin handler (`internal/interfaces/http/`)
2. Handler validates input and delegates to appropriate use case (`internal/application/usecases/`)
3. Use case executes business logic using repository interfaces (`internal/domain/repositories/`)
4. Repository implementations handle data operations (`internal/infrastructure/database/memory/`)
5. Use case returns result to handler
6. Handler formats JSON response

### Extending the Application
- To add new features: Create new use cases in `internal/application/usecases/[feature]/`
- To add new entities: Define structs in `internal/domain/entities/` and corresponding repository interfaces
- To change persistence: Implement repository interfaces in `internal/infrastructure/database/` (e.g., for PostgreSQL)
- To modify API endpoints: Update handlers in `internal/interfaces/http/` and routes in `internal/infrastructure/router/router.go`
- To add validation: Enhance use case input DTOs or add middleware in router

## Project Conventions
- Uses Go modules for dependency management (`go.mod`, `go.sum`)
- Follows standard Go project layout with `internal/` package for private code
- Leverages Gin web framework for HTTP handling
- Uses Swaggo for API documentation generation (docs generated in `docs/` directory)
- Repository pattern for data access abstraction
- Dependency injection for use cases and handlers
- UUIDs for entity identifiers (using `github.com/google/uuid`)
- Error handling through return values in use cases