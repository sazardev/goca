# Project Structure

This page explains the directory and file organization that Goca generates, following Clean Architecture best practices in Go.

## 📁 Complete Structure

```
my-project/
├── cmd/                              # Application entry points
│   └── server/
│       └── main.go                   # Main HTTP server
├── internal/                         # Private application code
│   ├── domain/                       # 🟡 Domain Layer
│   │   ├── user.go                   # Business entities
│   │   ├── product.go
│   │   └── errors.go                 # Domain errors
│   ├── usecase/                      # 🔴 Use Case Layer
│   │   ├── dto/                      # Data Transfer Objects
│   │   │   ├── user_dto.go
│   │   │   └── product_dto.go
│   │   ├── interfaces/               # Contracts between layers
│   │   │   ├── user_interfaces.go
│   │   │   └── product_interfaces.go
│   │   ├── user_usecase.go           # Application services
│   │   └── product_usecase.go
│   ├── repository/                   # 🔵 Infrastructure Layer
│   │   ├── interfaces/               # Repository interfaces
│   │   │   ├── user_repository.go
│   │   │   └── product_repository.go
│   │   ├── postgres/                 # PostgreSQL implementations
│   │   │   ├── user_repository.go
│   │   │   └── product_repository.go
│   │   ├── mysql/                    # MySQL implementations
│   │   └── mongodb/                  # MongoDB implementations
│   ├── handler/                      # 🟢 Adapter Layer
│   │   ├── http/                     # HTTP REST handlers
│   │   │   ├── dto/                  # HTTP-specific DTOs
│   │   │   ├── user_handler.go
│   │   │   ├── user_routes.go
│   │   │   ├── product_handler.go
│   │   │   └── middleware/           # HTTP middlewares
│   │   ├── grpc/                     # gRPC handlers
│   │   │   ├── user.proto
│   │   │   ├── user_server.go
│   │   │   └── product_server.go
│   │   ├── cli/                      # CLI commands
│   │   │   ├── user_commands.go
│   │   │   └── product_commands.go
│   │   └── worker/                   # Background workers
│   │       ├── user_worker.go
│   │       └── order_worker.go
│   ├── infrastructure/               # Infrastructure configuration
│   │   ├── di/                       # Dependency injection
│   │   │   ├── container.go
│   │   │   ├── wire.go               # Wire.dev (optional)
│   │   │   └── wire_gen.go
│   │   ├── database/                 # DB configuration
│   │   │   ├── postgres.go
│   │   │   ├── mysql.go
│   │   │   └── migrations/
│   │   └── cache/                    # Cache configuration
│   │       ├── redis.go
│   │       └── memory.go
│   ├── messages/                     # Messages and constants
│   │   ├── errors.go                 # Error messages
│   │   ├── responses.go              # Response messages
│   │   └── constants.go              # System constants
│   └── constants/                    # Feature-specific constants
│       ├── user_constants.go
│       └── product_constants.go
├── pkg/                              # Reusable/public code
│   ├── config/                       # Application configuration
│   │   ├── config.go
│   │   └── database.go
│   ├── logger/                       # Logging system
│   │   ├── logger.go
│   │   └── interfaces.go
│   ├── auth/                         # Authentication system
│   │   ├── jwt.go
│   │   ├── middleware.go
│   │   └── service.go
│   ├── validator/                    # Reusable validations
│   │   ├── validator.go
│   │   └── custom_rules.go
│   ├── utils/                        # General utilities
│   │   ├── crypto.go
│   │   ├── time.go
│   │   └── strings.go
│   └── errors/                       # Global error handling
│       ├── errors.go
│       ├── codes.go
│       └── handler.go
├── api/                              # API documentation
│   ├── openapi/                      # OpenAPI specifications
│   │   ├── swagger.yaml
│   │   └── user.yaml
│   └── proto/                        # Protocol Buffers files
│       ├── user.proto
│       └── product.proto
├── web/                              # Static web files (optional)
│   ├── static/
│   └── templates/
├── docs/                             # Project documentation
│   ├── architecture.md
│   ├── api.md
│   └── deployment.md
├── scripts/                          # Automation scripts
│   ├── build.sh
│   ├── test.sh
│   └── migrate.sh
├── deployments/                      # Deployment configurations
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   ├── kubernetes/
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   └── terraform/
├── test/                             # Integration and E2E tests
│   ├── integration/
│   ├── e2e/
│   └── fixtures/
├── migrations/                       # Database migrations
│   ├── 001_initial_schema.sql
│   ├── 002_add_users_table.sql
│   └── 003_add_products_table.sql
├── .github/                          # GitHub configuration
│   └── workflows/
│       ├── ci.yml
│       └── release.yml
├── go.mod                            # Module dependencies
├── go.sum                            # Dependency checksums
├── .env.example                      # Environment variables example
├── .gitignore                        # Files ignored by Git
├── Makefile                          # Automation commands
├── README.md                         # Main documentation
└── CHANGELOG.md                      # Change history
```

## 🏗️ Clean Architecture Layers

### 🟡 Domain Layer (`internal/domain/`)

**Purpose**: Contains core business logic and business rules.

**Typical files:**
```
domain/
├── user.go              # User entity with business methods
├── product.go           # Product entity with validations
├── order.go             # Order entity with business rules
├── errors.go            # Domain-specific errors
└── validations.go       # Reusable business validations
```

**Characteristics:**
- ✅ **No external dependencies**
- ✅ **Rich entities** with behavior
- ✅ **Business rules** encapsulated
- ✅ **Domain validations**
- ❌ **Must not know** infrastructure

### 🔴 Use Case Layer (`internal/usecase/`)

**Purpose**: Orchestrates application logic and coordinates between domain and infrastructure.

**Typical files:**
```
usecase/
├── dto/                          # Data Transfer Objects
│   ├── user_dto.go              # DTOs for user operations
│   └── common_dto.go            # Shared DTOs
├── interfaces/                   # Contracts between layers
│   ├── user_interfaces.go       # User UseCase interfaces
│   └── repositories.go          # Repository interfaces
├── user_usecase.go              # Use case implementation
├── product_usecase.go           # Product use cases
└── common_usecase.go            # Shared logic
```

**Characteristics:**
- ✅ **Orchestrates** workflows
- ✅ **DTOs** for data transfer
- ✅ **Interfaces** to decouple layers
- ✅ **Application validations**
- ❌ **Must not know** HTTP/DB details

### 🟢 Adapter Layer (`internal/handler/`)

**Purpose**: Adapts external interfaces (HTTP, gRPC, CLI) to internal use cases.

**Typical files:**
```
handler/
├── http/                         # HTTP adapters
│   ├── dto/                     # HTTP-specific DTOs
│   │   ├── user_http_dto.go     # HTTP Request/Response
│   │   └── error_dto.go         # Error responses
│   ├── user_handler.go          # Handler for user endpoints
│   ├── user_routes.go           # Route definitions
│   └── middleware/              # HTTP middlewares
│       ├── auth.go              # Authentication middleware
│       ├── cors.go              # CORS middleware
│       └── logging.go           # Logging middleware
├── grpc/                        # gRPC adapters
│   ├── user.proto              # Service definitions
│   ├── user_server.go          # Server implementation
│   └── interceptors/           # gRPC interceptors
├── cli/                        # CLI commands
│   ├── user_commands.go        # User commands
│   └── root.go                 # Root command
└── worker/                     # Background workers
    ├── user_worker.go          # User worker
    └── queue.go                # Queue configuration
```

**Characteristics:**
- ✅ **Adapts** external protocols
- ✅ **Protocol-specific DTOs**
- ✅ **Appropriate error handling**
- ✅ **Input validation**
- ❌ **Must not contain** business logic

### 🔵 Infrastructure Layer (`internal/repository/`, `pkg/`)

**Purpose**: Implements technical details like persistence, logging, configuration.

**Typical files:**
```
repository/
├── interfaces/                   # Persistence contracts
│   ├── user_repository.go       # Repository interface
│   └── transaction.go           # Transaction interface
├── postgres/                    # PostgreSQL implementation
│   ├── user_repository.go      # Specific repository
│   ├── migrations.go           # Migrations
│   └── connection.go           # Connection configuration
├── mysql/                      # MySQL implementation
├── mongodb/                    # MongoDB implementation
└── memory/                     # In-memory implementation (tests)
    └── user_repository.go
```

**Characteristics:**
- ✅ **Implements** domain interfaces
- ✅ **Technology-specific** details
- ✅ **Connection configuration**
- ✅ **Database migrations**
- ❌ **Must not expose** technical details

## 📦 Special Directories

### `cmd/` - Entry Points

```
cmd/
├── server/                      # Main HTTP server
│   └── main.go
├── migrate/                     # Migration tool
│   └── main.go
├── worker/                      # Background worker
│   └── main.go
└── cli/                        # CLI tool
    └── main.go
```

**Purpose**: Each subdirectory represents a different executable.

### `pkg/` - Reusable Code

```
pkg/
├── config/                      # Global configuration
├── logger/                      # Logging system
├── auth/                        # Authentication/authorization
├── validator/                   # Reusable validations
├── utils/                       # General utilities
└── errors/                      # Global error handling
```

**Purpose**: Code that can be imported by other projects.

### `api/` - API Documentation

```
api/
├── openapi/                     # OpenAPI/Swagger specifications
│   ├── swagger.yaml            # Main documentation
│   ├── user.yaml               # User endpoints
│   └── product.yaml            # Product endpoints
└── proto/                      # Protocol Buffers for gRPC
    ├── user.proto
    ├── product.proto
    └── common.proto
```

**Purpose**: API contracts and external documentation.

### `migrations/` - Database Schema

```
migrations/
├── 001_initial_schema.sql       # Initial schema
├── 002_add_users_table.sql      # Add users table
├── 003_add_products_table.sql   # Add products table
└── 004_add_indexes.sql          # Add indexes
```

**Purpose**: Database schema version control.

## 🔄 Dependency Flow

```
┌─────────────────────────────────────────┐
│                 cmd/                    │ ← Entry points
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            internal/handler/            │ ← 🟢 Adapters
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            internal/usecase/            │ ← 🔴 Use Cases
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            internal/domain/             │ ← 🟡 Domain
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│          internal/repository/           │ ← 🔵 Infrastructure
└─────────────────────────────────────────┘
```

**Fundamental Rule**: Dependencies always point inward.

## 📁 Naming Conventions

### Files
- **Entities**: `user.go`, `product.go`, `order.go`
- **DTOs**: `user_dto.go`, `create_user_request.go`
- **UseCase**: `user_usecase.go`, `product_service.go`
- **Repository**: `user_repository.go`, `postgres_user_repo.go`
- **Handler**: `user_handler.go`, `user_routes.go`
- **Interfaces**: `user_interfaces.go`, `repositories.go`

### Packages
- **Lowercase**: Always lowercase
- **Descriptive**: `usecase`, `repository`, `handler`
- **No hyphens**: `userservice` not `user-service`
- **Singular**: `user` not `users` (except when appropriate)

### Structures
```go
// Entities: PascalCase
type User struct {}
type OrderItem struct {}

// Interfaces: PascalCase + suffix
type UserRepository interface {}
type UserUseCase interface {}

// DTOs: PascalCase + purpose
type CreateUserRequest struct {}
type UserResponse struct {}
```

## 🎯 Benefits of this Structure

### ✅ Clear Separation of Responsibilities
- Each directory has a specific purpose
- Easy to locate related code
- Changes in one layer don't affect others

### ✅ Testability
- Interfaces allow easy mocks
- Unit tests per layer
- Separate integration tests

### ✅ Scalability
- Adding new features is predictable
- Consistent structure between features
- Easy onboarding for new developers

### ✅ Maintainability
- Organized and predictable code
- Safe refactoring by layers
- Explicit dependencies

### ✅ Flexibility
- Easy to change implementations
- Support multiple protocols
- Add new functionality without breaking existing

## 🛠️ Customization

### Add New Layer
```bash
# Create new events layer
mkdir -p internal/events
mkdir -p internal/events/handlers
mkdir -p internal/events/publishers
```

### Add New Protocol
```bash
# Add GraphQL support
mkdir -p internal/handler/graphql
mkdir -p internal/handler/graphql/resolvers
mkdir -p internal/handler/graphql/schemas
```

### Add New Database
```bash
# Add Redis support
mkdir -p internal/repository/redis
mkdir -p pkg/cache/redis
```

## 📚 Additional Resources

- [Clean Architecture Principles](Clean-Architecture)
- [Design Patterns Used](Design-Patterns)
- [Testing Guide](Testing-Guide)
- [Best Practices](Best-Practices)

---

**← [Clean Architecture](Clean-Architecture) | [Implemented Patterns](Design-Patterns) →**
