# Clean Architecture

Learn how Goca implements and enforces **Clean Architecture** principles from Robert C. Martin (Uncle Bob) in your Go projects.

## What is Clean Architecture?

Clean Architecture is a software design philosophy that separates code into **layers** with clear responsibilities and dependencies that point inward toward the core business logic.

### Core Principles

1. **Independence of Frameworks** - Business logic doesn't depend on libraries
2. **Testability** - Business rules can be tested without UI, database, or external services
3. **Independence of UI** - Change UI without changing business logic
4. **Independence of Database** - Swap databases without affecting business rules
5. **Independence of External Systems** - Business logic knows nothing about the outside world

## The Dependency Rule

> **Source code dependencies must point only inward, toward higher-level policies.**

```
┌─────────────────────────────────────┐
│     External Interfaces & I/O       │  ← Frameworks, Drivers, APIs
└─────────────────┬───────────────────┘
                  │ depends on
┌─────────────────▼───────────────────┐
│    Interface Adapters (Handlers)    │  ← Controllers, Presenters
└─────────────────┬───────────────────┘
                  │ depends on
┌─────────────────▼───────────────────┐
│    Application Business Rules       │  ← Use Cases
└─────────────────┬───────────────────┘
                  │ depends on
┌─────────────────▼───────────────────┐
│   Enterprise Business Rules         │  ← Entities (Domain)
│         (No Dependencies)            │
└─────────────────────────────────────┘
```

**Inner layers never know about outer layers.**

## Goca's 4 Layers

### 🟡 Layer 1: Domain (Entities)

**Location**: `internal/domain/`

The innermost layer containing **enterprise-wide business rules**.

#### Responsibilities

- Define business entities
- Implement core business rules
- Define domain errors
- Declare repository interfaces
- Domain-specific validations

#### Example: User Entity

```go
package domain

import (
    "errors"
    "strings"
    "time"
)

// User represents a user entity in our system
type User struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Validate enforces business rules
func (u *User) Validate() error {
    if strings.TrimSpace(u.Name) == "" {
        return ErrUserNameRequired
    }
    
    if len(u.Name) < 2 {
        return ErrUserNameTooShort
    }
    
    if !u.IsValidEmail() {
        return ErrInvalidEmail
    }
    
    return nil
}

// IsAdmin checks if user has admin privileges
func (u *User) IsAdmin() bool {
    return u.Role == "admin"
}

// IsValidEmail validates email format (business rule)
func (u *User) IsValidEmail() bool {
    return strings.Contains(u.Email, "@") && len(u.Email) > 5
}

// Domain-specific errors
var (
    ErrUserNameRequired = errors.New("user name is required")
    ErrUserNameTooShort = errors.New("user name must be at least 2 characters")
    ErrInvalidEmail     = errors.New("invalid email format")
    ErrUserNotFound     = errors.New("user not found")
)
```

::: tip Domain Layer Rules
✅ **Do**: Pure business logic, no external dependencies  
❌ **Don't**: Import HTTP, database, or framework packages
:::

### 🔴 Layer 2: Use Cases (Application Logic)

**Location**: `internal/usecase/`

Contains **application-specific business rules**.

#### Responsibilities

- Define use case interfaces
- Implement application workflows
- Coordinate between repositories
- Define DTOs (Data Transfer Objects)
- Input validation

#### Example: User Use Case

```go
package usecase

import (
    "context"
    "myproject/internal/domain"
)

// UserUseCase defines user-related operations
type UserUseCase interface {
    CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
    GetUser(ctx context.Context, id uint) (*UserResponse, error)
    UpdateUser(ctx context.Context, id uint, req UpdateUserRequest) error
    DeleteUser(ctx context.Context, id uint) error
    ListUsers(ctx context.Context) ([]*UserResponse, error)
}

// CreateUserRequest - Input DTO
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2"`
    Email string `json:"email" validate:"required,email"`
    Role  string `json:"role" validate:"required,oneof=admin user"`
}

// UserResponse - Output DTO
type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

// userUseCase implements UserUseCase
type userUseCase struct {
    userRepo domain.UserRepository // Depends on interface!
}

func NewUserUseCase(userRepo domain.UserRepository) UserUseCase {
    return &userUseCase{userRepo: userRepo}
}

func (uc *userUseCase) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
    // 1. Validate input
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // 2. Create domain entity
    user := &domain.User{
        Name:  req.Name,
        Email: req.Email,
        Role:  req.Role,
    }
    
    // 3. Validate business rules
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    // 4. Persist through repository
    if err := uc.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // 5. Return DTO
    return &UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        Role:  user.Role,
    }, nil
}
```

::: tip Use Case Layer Rules
✅ **Do**: Application workflows, DTOs, coordinate repositories  
❌ **Don't**: HTTP/gRPC details, SQL queries, framework-specific code
:::

### 🔵 Layer 3: Repository (Infrastructure)

**Location**: `internal/repository/`

Implements **data access and external communication**.

#### Responsibilities

- Implement repository interfaces from domain
- Handle database operations
- Manage database connections
- Transform between DB models and domain entities

#### Example: PostgreSQL Repository

```go
package repository

import (
    "context"
    "database/sql"
    "myproject/internal/domain"
)

type postgresUserRepository struct {
    db *sql.DB
}

// NewPostgresUserRepository creates a new repository
func NewPostgresUserRepository(db *sql.DB) domain.UserRepository {
    return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Save(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (name, email, role, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
    
    err := r.db.QueryRowContext(
        ctx, query,
        user.Name, user.Email, user.Role,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    return err
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
    query := `
        SELECT id, name, email, role, created_at, updated_at
        FROM users WHERE id = $1
    `
    
    user := &domain.User{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID, &user.Name, &user.Email, &user.Role,
        &user.CreatedAt, &user.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, domain.ErrUserNotFound
    }
    
    return user, err
}

func (r *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {
    query := `
        UPDATE users 
        SET name = $1, email = $2, role = $3, updated_at = NOW()
        WHERE id = $4
    `
    
    _, err := r.db.ExecContext(ctx, query,
        user.Name, user.Email, user.Role, user.ID,
    )
    
    return err
}
```

::: tip Repository Layer Rules
✅ **Do**: Implement domain interfaces, handle persistence  
❌ **Don't**: Business logic, validation rules
:::

### 🟢 Layer 4: Handlers (Interface Adapters)

**Location**: `internal/handler/http/`

Adapts **external requests to use cases**.

#### Responsibilities

- Handle HTTP/gRPC/CLI requests
- Parse and validate input
- Call use cases
- Format responses
- Handle HTTP-specific concerns

#### Example: HTTP Handler

```go
package http

import (
    "encoding/json"
    "net/http"
    "strconv"
    "myproject/internal/usecase"
    "github.com/gorilla/mux"
)

type UserHandler struct {
    userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
    return &UserHandler{userUseCase: userUseCase}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 1. Parse HTTP request
    var req usecase.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    // 2. Call use case
    user, err := h.userUseCase.CreateUser(r.Context(), req)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    // 3. Send HTTP response
    respondJSON(w, http.StatusCreated, user)
}

// GetUser handles GET /users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    // Parse path parameter
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        respondError(w, http.StatusBadRequest, "Invalid user ID")
        return
    }
    
    user, err := h.userUseCase.GetUser(r.Context(), uint(id))
    if err != nil {
        respondError(w, http.StatusNotFound, "User not found")
        return
    }
    
    respondJSON(w, http.StatusOK, user)
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteStatus(status)
    json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
    respondJSON(w, status, map[string]string{"error": message})
}
```

::: tip Handler Layer Rules
✅ **Do**: Protocol-specific concerns, request/response formatting  
❌ **Don't**: Business logic, direct database access
:::

## Complete Data Flow

Here's how a request flows through all layers:

```
1. HTTP Request
   ↓
2. Handler parses request → CreateUserRequest DTO
   ↓
3. UseCase validates and applies business rules
   ↓
4. UseCase creates Domain Entity
   ↓
5. Entity validates its own business rules
   ↓
6. UseCase calls Repository interface
   ↓
7. Repository saves to database
   ↓
8. Repository returns Domain Entity
   ↓
9. UseCase transforms to UserResponse DTO
   ↓
10. Handler formats HTTP Response
```

## Benefits of This Architecture

### 1. Testability

Test each layer in isolation:

```go
// Test use case without HTTP or database
func TestCreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    useCase := usecase.NewUserUseCase(mockRepo)
    
    req := usecase.CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
        Role:  "user",
    }
    
    user, err := useCase.CreateUser(context.Background(), req)
    assert.NoError(t, err)
    assert.Equal(t, "John Doe", user.Name)
}
```

### 2. Flexibility

Swap implementations without touching business logic:

```go
// Switch from PostgreSQL to MongoDB
// Old: postgresRepo := repository.NewPostgresUserRepository(db)
// New: mongoRepo := repository.NewMongoUserRepository(client)
userUseCase := usecase.NewUserUseCase(mongoRepo) // Same interface!
```

### 3. Maintainability

Changes are localized to specific layers:

- UI change? → Only handler layer
- Database change? → Only repository layer
- Business rule change? → Only domain/usecase layer

## Common Mistakes to Avoid

### ❌ Skip Layers

```go
// BAD: Handler directly accessing database
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    db.Exec("INSERT INTO users...") // ❌ Skipping use case!
}
```

### ❌ Wrong Dependencies

```go
// BAD: Domain depending on outer layer
package domain

import "net/http" // ❌ Domain shouldn't know about HTTP!
```

### ❌ Business Logic in Handlers

```go
// BAD: Validation in handler
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    if user.Name == "" { // ❌ This belongs in domain/usecase!
        return errors.New("name required")
    }
}
```

## Learn More

- 📖 [Project Structure](/guide/project-structure) - Directory organization
- 🎓 [Complete Tutorial](/tutorials/complete-tutorial) - Build a real app
- 📚 [Best Practices](/guide/best-practices) - Tips and conventions

## Resources

- [Clean Architecture Book](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164) by Robert C. Martin
- [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - Uncle Bob's blog
