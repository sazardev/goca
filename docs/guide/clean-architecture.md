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
 **Do**: Pure business logic, no external dependencies  
 **Don't**: Import HTTP, database, or framework packages
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
    "myproject/internal/domain"
    "myproject/internal/repository"
)

// UserUseCase defines user-related operations
type UserUseCase interface {
    Create(input CreateUserInput) (*CreateUserOutput, error)
    GetByID(id uint) (*domain.User, error)
    Update(id uint, input UpdateUserInput) (*domain.User, error)
    Delete(id uint) error
    List() (*ListUserOutput, error)
}

// CreateUserInput - Input DTO
type CreateUserInput struct {
    Name  string `json:"name" validate:"required,min=2"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"required,min=1"`
}

// CreateUserOutput - Output DTO
type CreateUserOutput struct {
    User    domain.User `json:"user"`
    Message string      `json:"message"`
}

// UpdateUserInput - Input DTO
type UpdateUserInput struct {
    Name  *string `json:"name,omitempty" validate:"omitempty,min=2"`
    Email *string `json:"email,omitempty" validate:"omitempty,email"`
    Age   *int    `json:"age,omitempty" validate:"omitempty,min=1"`
}

// ListUserOutput - Output DTO
type ListUserOutput struct {
    Users   []domain.User `json:"users"`
    Total   int           `json:"total"`
    Message string        `json:"message"`
}

// userService implements UserUseCase
type userService struct {
    repo repository.UserRepository // Depends on interface!
}

func NewUserService(repo repository.UserRepository) UserUseCase {
    return &userService{repo: repo}
}

func (s *userService) Create(input CreateUserInput) (*CreateUserOutput, error) {
    // 1. Create domain entity
    user := &domain.User{
        Name:  input.Name,
        Email: input.Email,
    }
    
    // 2. Validate business rules
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    // 3. Persist through repository
    if err := s.repo.Save(user); err != nil {
        return nil, err
    }
    
    // 4. Return DTO
    return &CreateUserOutput{
        User:    *user,
        Message: "User created successfully",
    }, nil
}
```

::: tip Use Case Layer Rules
 **Do**: Application workflows, DTOs, coordinate repositories  
 **Don't**: HTTP/gRPC details, SQL queries, framework-specific code
:::

### 🔵 Layer 3: Repository (Infrastructure)

**Location**: `internal/repository/`

Implements **data access and external communication**.

#### Responsibilities

- Define repository interfaces
- Implement repository interfaces
- Handle database operations
- Manage database connections
- Transform between DB models and domain entities

#### Example: GORM Repository

```go
package repository

import (
    "myproject/internal/domain"

    "gorm.io/gorm"
)

// UserRepository defines data access interface
type UserRepository interface {
    Save(user *domain.User) error
    FindByID(id int) (*domain.User, error)
    FindByEmail(email string) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id int) error
    FindAll() ([]domain.User, error)
}

type postgresUserRepository struct {
    db *gorm.DB
}

// NewPostgresUserRepository creates a new repository
func NewPostgresUserRepository(db *gorm.DB) UserRepository {
    return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Save(user *domain.User) error {
    result := r.db.Create(user)
    return result.Error
}

func (r *postgresUserRepository) FindByID(id int) (*domain.User, error) {
    user := &domain.User{}
    result := r.db.First(user, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return user, nil
}

func (r *postgresUserRepository) Update(user *domain.User) error {
    result := r.db.Save(user)
    return result.Error
}

func (r *postgresUserRepository) Delete(id int) error {
    result := r.db.Delete(&domain.User{}, id)
    return result.Error
}

func (r *postgresUserRepository) FindAll() ([]domain.User, error) {
    var users []domain.User
    result := r.db.Find(&users)
    if result.Error != nil {
        return nil, result.Error
    }
    return users, nil
}
```

::: tip Repository Layer Rules
 **Do**: Define and implement repository interfaces, handle persistence  
 **Don't**: Business logic, validation rules
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
    var input usecase.CreateUserInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    // 2. Call use case
    output, err := h.userUseCase.Create(input)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    // 3. Send HTTP response
    respondJSON(w, http.StatusCreated, output)
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
    
    user, err := h.userUseCase.GetByID(uint(id))
    if err != nil {
        respondError(w, http.StatusNotFound, "User not found")
        return
    }
    
    respondJSON(w, http.StatusOK, user)
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
    respondJSON(w, status, map[string]string{"error": message})
}
```

::: tip Handler Layer Rules
 **Do**: Protocol-specific concerns, request/response formatting  
 **Don't**: Business logic, direct database access
:::

## Complete Data Flow

Here's how a request flows through all layers:

```
1. HTTP Request
   ↓
2. Handler parses request → CreateUserInput DTO
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
9. UseCase transforms to CreateUserOutput DTO
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
    useCase := usecase.NewUserService(mockRepo)
    
    input := usecase.CreateUserInput{
        Name:  "John Doe",
        Email: "john@example.com",
        Age:   30,
    }
    
    output, err := useCase.Create(input)
    assert.NoError(t, err)
    assert.Equal(t, "John Doe", output.User.Name)
    assert.Equal(t, "User created successfully", output.Message)
}
```

### 2. Flexibility

Swap implementations without touching business logic:

```go
// Switch from PostgreSQL to MongoDB
// Old: postgresRepo := repository.NewPostgresUserRepository(db)
// New: mongoRepo := repository.NewMongoUserRepository(client)
userService := usecase.NewUserService(postgresRepo) // Same interface!
```

### 3. Maintainability

Changes are localized to specific layers:

- UI change? → Only handler layer
- Database change? → Only repository layer
- Business rule change? → Only domain/usecase layer

## Common Mistakes to Avoid

###  Skip Layers

```go
// BAD: Handler directly accessing database
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    db.Exec("INSERT INTO users...") //  Skipping use case!
}
```

###  Wrong Dependencies

```go
// BAD: Domain depending on outer layer
package domain

import "net/http" //  Domain shouldn't know about HTTP!
```

###  Business Logic in Handlers

```go
// BAD: Validation in handler
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    if user.Name == "" { //  This belongs in domain/usecase!
        return errors.New("name required")
    }
}
```

## Learn More

-  [Project Structure](/guide/project-structure) - Directory organization
-  [Complete Tutorial](/tutorials/complete-tutorial) - Build a real app
-  [Best Practices](/guide/best-practices) - Tips and conventions

## Resources

- [Clean Architecture Book](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164) by Robert C. Martin
- [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - Uncle Bob's blog
