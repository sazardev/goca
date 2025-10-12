# Project Structure

Understanding the directory organization and file conventions.

## Overview

Goca follows Clean Architecture principles with clear separation of concerns:

```
myproject/
├── cmd/                    # Application entry points
│   └── server/
│       └── main.go
├── internal/               # Private application code
│   ├── domain/            # Business entities
│   ├── usecase/           # Business logic
│   ├── repository/        # Data access
│   ├── handler/           # Delivery mechanisms
│   ├── di/                # Dependency injection
│   └── messages/          # Errors and constants
├── pkg/                    # Public libraries
├── config/                 # Configuration files
├── migrations/             # Database migrations
├── scripts/                # Build and deployment scripts
└── docs/                   # Documentation
```

## Layer Responsibilities

### Domain (`internal/domain/`)

Contains business entities and core business rules.

**Files:**
- `{entity}.go` - Entity definition
- `{entity}_seeds.go` - Seed data
- `errors.go` - Domain-specific errors

**Example:**
```go
// user.go
package domain

type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string `gorm:"not null"`
    Email     string `gorm:"unique;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Use Case (`internal/usecase/`)

Business logic and application rules.

**Files:**
- `{entity}_service.go` - Business logic implementation
- `interfaces.go` - Service contracts
- `dto.go` - Data transfer objects

**Example:**
```go
// user_service.go
package usecase

type UserService interface {
    CreateUser(ctx context.Context, input CreateUserInput) (*UserResponse, error)
    GetUser(ctx context.Context, id uint) (*UserResponse, error)
}
```

### Repository (`internal/repository/`)

Data persistence layer.

**Files:**
- `postgres_{entity}_repository.go` - PostgreSQL implementation
- `interfaces.go` - Repository contracts

**Example:**
```go
// postgres_user_repository.go
package repository

type PostgresUserRepository struct {
    db *gorm.DB
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}
```

### Handler (`internal/handler/http/`)

HTTP delivery layer.

**Files:**
- `{entity}_handler.go` - HTTP handlers
- `routes.go` - Route registration
- `swagger.yaml` - API documentation

**Example:**
```go
// user_handler.go
package http

type UserHandler struct {
    service usecase.UserService
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    // Handle HTTP request
}
```

### Dependency Injection (`internal/di/`)

Wires all components together.

**Files:**
- `container.go` - DI container

**Example:**
```go
// container.go
package di

type Container struct {
    UserRepository repository.UserRepository
    UserService    usecase.UserService
    UserHandler    *http.UserHandler
}
```

## File Naming Conventions

| Pattern                           | Example                       | Purpose           |
| --------------------------------- | ----------------------------- | ----------------- |
| `{entity}.go`                     | `user.go`                     | Entity definition |
| `{entity}_service.go`             | `user_service.go`             | Business logic    |
| `postgres_{entity}_repository.go` | `postgres_user_repository.go` | Data access       |
| `{entity}_handler.go`             | `user_handler.go`             | HTTP handlers     |
| `{entity}_seeds.go`               | `user_seeds.go`               | Seed data         |
| `{entity}_test.go`                | `user_test.go`                | Unit tests        |

## Package Organization

### Internal vs Pkg

```
internal/     # Private to this application
pkg/          # Can be imported by other projects
```

### Feature Grouping

Each feature spans multiple layers:

```
User Feature:
├── internal/domain/user.go
├── internal/usecase/user_service.go
├── internal/repository/postgres_user_repository.go
└── internal/handler/http/user_handler.go
```

## Import Rules

Follow dependency direction:

```
Handler → Use Case → Repository → Domain
   ↓         ↓           ↓          ↓
 HTTP    Business    Database   Entities
```

** Allowed:**
```go
// Handler imports use case
import "myproject/internal/usecase"

// Use case imports repository interface
import "myproject/internal/repository"
```

** Not Allowed:**
```go
// Repository imports handler (wrong direction!)
import "myproject/internal/handler/http"
```

## Configuration Files

```
config/
├── config.yaml          # Application settings
├── .env.example         # Environment template
└── database.yaml        # Database configuration
```

## Migrations

```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_products_table.up.sql
└── 000002_create_products_table.down.sql
```

## Testing Structure

Mirror the source structure:

```
internal/
├── usecase/
│   ├── user_service.go
│   └── user_service_test.go
└── repository/
    ├── postgres_user_repository.go
    └── postgres_user_repository_test.go
```

## See Also

- [Clean Architecture](/guide/clean-architecture) - Architecture principles
- [Best Practices](/guide/best-practices) - Coding guidelines
