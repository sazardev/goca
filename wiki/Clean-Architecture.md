# Clean Architecture

This page explains in detail how Goca implements and enforces **Clean Architecture** principles from Uncle Bob (Robert C. Martin) in Go projects.

## 🎯 What is Clean Architecture?

Clean Architecture is an architectural pattern that organizes code in **concentric layers** where dependencies point toward the center of the system, ensuring:

- 🔒 **Framework independence**
- 🧪 **Complete testability**
- 🎨 **UI independence**
- 💾 **Database independence**
- 🌐 **External agent independence**

## 🏗️ The 4 Layers of Clean Architecture

### 🟡 1. Domain Layer (Entities)
**Location**: `internal/domain/`  
**Responsibility**: Core business logic and enterprise rules

#### ✅ What SHOULD be here:
- Business entities
- Fundamental business rules
- Domain validations
- Repository interfaces
- Domain-specific errors

#### ❌ What should NOT be here:
- External dependencies (databases, APIs)
- Presentation logic
- Implementation details
- External frameworks or libraries

#### 📄 Entity Example:
```go
package domain

import (
    "errors"
    "strings"
    "time"
)

// User represents a user in the system
type User struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Validate implements business rules for validating a user
func (u *User) Validate() error {
    if strings.TrimSpace(u.Name) == "" {
        return ErrUserNameRequired
    }
    
    if len(u.Name) < 2 {
        return ErrUserNameTooShort
    }
    
    if !u.isValidEmail() {
        return ErrUserEmailInvalid
    }
    
    return nil
}

// CanUpdateProfile verifies if the user can update their profile
func (u *User) CanUpdateProfile() bool {
    return u.ID > 0 && u.Name != ""
}

// isValidEmail validates email format (business rule)
func (u *User) isValidEmail() bool {
    return strings.Contains(u.Email, "@") && 
           strings.Contains(u.Email, ".") &&
           len(u.Email) > 5
}

// Domain errors
var (
    ErrUserNameRequired  = errors.New("user name is required")
    ErrUserNameTooShort  = errors.New("user name must be at least 2 characters")
    ErrUserEmailInvalid  = errors.New("user email format is invalid")
    ErrUserNotFound      = errors.New("user not found")
)
```

### 🔴 2. Use Cases Layer (Use Cases)
**Location**: `internal/usecase/`  
**Responsibility**: Application logic and orchestration

#### ✅ What SHOULD be here:
- DTOs (Data Transfer Objects)
- Use case interfaces
- Application services
- Input validations
- Repository coordination

#### ❌ What should NOT be here:
- Presentation logic
- Database details
- Web framework logic
- Specific infrastructure implementations

#### 📄 Use Case Example:
```go
package usecase

import (
    "context"
    "github.com/user/project/internal/domain"
)

// UserUseCase defines contracts for user use cases
type UserUseCase interface {
    Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
    GetByID(ctx context.Context, id uint) (*UserResponse, error)
    Update(ctx context.Context, id uint, req UpdateUserRequest) (*UserResponse, error)
    Delete(ctx context.Context, id uint) error
}

// UserRepository defines contracts for persistence
type UserRepository interface {
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id uint) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uint) error
}

// userUseCase implements application logic
type userUseCase struct {
    userRepo UserRepository
}

// NewUserUseCase creates a new use case instance
func NewUserUseCase(userRepo UserRepository) UserUseCase {
    return &userUseCase{
        userRepo: userRepo,
    }
}

// Create creates a new user
func (uc *userUseCase) Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
    // 1. Validate input DTO
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // 2. Create domain entity
    user := &domain.User{
        Name:  req.Name,
        Email: req.Email,
    }
    
    // 3. Validate business rules
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    // 4. Check application rules (unique email)
    existingUser, _ := uc.userRepo.FindByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, domain.ErrUserEmailAlreadyExists
    }
    
    // 5. Persist
    if err := uc.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // 6. Return response DTO
    return &UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }, nil
}
```

#### 📄 DTOs (Data Transfer Objects):
```go
package usecase

// CreateUserRequest DTO for creating user
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
}

// Validate validates the DTO
func (r *CreateUserRequest) Validate() error {
    if strings.TrimSpace(r.Name) == "" {
        return errors.New("name is required")
    }
    if len(r.Name) < 2 {
        return errors.New("name must be at least 2 characters")
    }
    if !strings.Contains(r.Email, "@") {
        return errors.New("invalid email format")
    }
    return nil
}

// UserResponse response DTO
type UserResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 🟢 3. Adapters Layer (Interface Adapters)
**Location**: `internal/handler/`  
**Responsibility**: Adapt input/output between protocols and use cases

#### ✅ What SHOULD be here:
- HTTP/gRPC/CLI handlers
- REST controllers
- Protocol adapters
- Protocol-specific DTOs
- Middlewares

#### ❌ What should NOT be here:
- Business logic
- Direct database access
- Business validations
- Enterprise rules

#### 📄 HTTP Handler Example:
```go
package http

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/user/project/internal/usecase"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
    userUseCase usecase.UserUseCase
}

// NewUserHandler creates a new handler
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
    return &UserHandler{
        userUseCase: userUseCase,
    }
}

// Create handles POST /users
func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserHTTPRequest
    
    // 1. Parse HTTP input
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: "Invalid request format",
        })
        return
    }
    
    // 2. Convert to use case DTO
    useCaseReq := usecase.CreateUserRequest{
        Name:  req.Name,
        Email: req.Email,
    }
    
    // 3. Execute use case
    user, err := h.userUseCase.Create(c.Request.Context(), useCaseReq)
    if err != nil {
        status := http.StatusInternalServerError
        if err == domain.ErrUserEmailAlreadyExists {
            status = http.StatusConflict
        }
        
        c.JSON(status, ErrorResponse{
            Error: err.Error(),
        })
        return
    }
    
    // 4. Convert to HTTP response
    response := UserHTTPResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
    
    c.JSON(http.StatusCreated, response)
}
```

### 🔵 4. Infrastructure Layer (Frameworks & Drivers)
**Location**: `internal/repository/`, `pkg/`  
**Responsibility**: Technology-specific implementations

#### ✅ What SHOULD be here:
- Repository implementations
- Database connections
- HTTP clients
- Configuration
- Logging
- Caches

#### ❌ What should NOT be here:
- Business logic
- Enterprise rules
- Domain validations
- Use case DTOs

#### 📄 Repository Example:
```go
package postgres

import (
    "context"
    "database/sql"
    
    "github.com/user/project/internal/domain"
)

// userRepository implements the repository for PostgreSQL
type userRepository struct {
    db *sql.DB
}

// NewUserRepository creates a new repository
func NewUserRepository(db *sql.DB) domain.UserRepository {
    return &userRepository{
        db: db,
    }
}

// Save implements PostgreSQL-specific persistence
func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (name, email, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
    
    err := r.db.QueryRowContext(ctx, query, user.Name, user.Email).Scan(
        &user.ID,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    return err
}

// FindByID searches for a user by ID
func (r *userRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE id = $1
    `
    
    user := &domain.User{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, domain.ErrUserNotFound
        }
        return nil, err
    }
    
    return user, nil
}
```

## 🔄 Dependency Flow

```
🌐 HTTP Request
     ↓
🟢 Handler (converts HTTP → DTO)
     ↓
🔴 UseCase (executes application logic)
     ↓
🟡 Domain (validates business rules)
     ↓
🔵 Repository (persists to database)
```

### Dependency Rule:
> **Dependencies ALWAYS point inward**

- 🟢 Handler depends on 🔴 UseCase
- 🔴 UseCase depends on 🟡 Domain
- 🔵 Repository implements 🟡 Domain interfaces
- 🟡 Domain does NOT depend on anything external

## 🎭 SOLID Principles Applied

### 🔵 Single Responsibility Principle (SRP)
Each class has a single reason to change:

```go
// ✅ GOOD - One responsibility
type UserValidator struct{}
func (v *UserValidator) Validate(user *User) error { /* ... */ }

// ✅ GOOD - One responsibility
type UserRepository struct{}
func (r *UserRepository) Save(user *User) error { /* ... */ }

// ❌ BAD - Multiple responsibilities
type UserService struct{}
func (s *UserService) ValidateAndSave(user *User) error {
    // Validation + Persistence = 2 responsibilities
}
```

### 🔓 Open/Closed Principle (OCP)
Open for extension, closed for modification:

```go
// Stable interface
type NotificationSender interface {
    Send(message string) error
}

// Extensible implementations
type EmailSender struct{} // New implementation
type SMSSender struct{}   // New implementation
type SlackSender struct{} // New implementation

// UseCase closed for modification
type UserUseCase struct {
    notifier NotificationSender // Uses interface
}
```

### 🔄 Liskov Substitution Principle (LSP)
Implementations must be interchangeable:

```go
// Any implementation of UserRepository
// must behave the same from the UseCase perspective
type PostgreSQLUserRepo struct{}
type MySQLUserRepo struct{}
type MongoUserRepo struct{}

// All implement the same interface
type UserRepository interface {
    Save(user *User) error
    FindByID(id uint) (*User, error)
}
```

### 🎯 Interface Segregation Principle (ISP)
Specific and cohesive interfaces:

```go
// ✅ GOOD - Specific interfaces
type UserReader interface {
    FindByID(id uint) (*User, error)
}

type UserWriter interface {
    Save(user *User) error
}

// ❌ BAD - Interface too large
type UserRepository interface {
    Save(user *User) error
    FindByID(id uint) (*User, error)
    FindByEmail(email string) (*User, error)
    Update(user *User) error
    Delete(id uint) error
    FindAll() ([]*User, error)
    Count() (int, error)
    // ... many more methods
}
```

### ⬇️ Dependency Inversion Principle (DIP)
Depend on abstractions, not concretions:

```go
// ✅ GOOD - Depends on interface (abstraction)
type UserUseCase struct {
    userRepo UserRepository // Interface
}

// ❌ BAD - Depends on concrete implementation
type UserUseCase struct {
    userRepo *PostgreSQLUserRepository // Specific implementation
}
```

## 🧪 Testability

Clean Architecture greatly facilitates testing:

### Unit Tests for Domain
```go
func TestUser_Validate(t *testing.T) {
    tests := []struct {
        name    string
        user    domain.User
        wantErr bool
    }{
        {
            name: "valid user",
            user: domain.User{Name: "John", Email: "john@example.com"},
            wantErr: false,
        },
        {
            name: "invalid email",
            user: domain.User{Name: "John", Email: "invalid"},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Unit Tests for UseCase with Mocks
```go
func TestUserUseCase_Create(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    useCase := NewUserUseCase(mockRepo)
    
    req := CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Mock expectations
    mockRepo.On("FindByEmail", "john@example.com").Return(nil, nil)
    mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil)
    
    // Act
    result, err := useCase.Create(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "John Doe", result.Name)
    mockRepo.AssertExpectations(t)
}
```

### Integration Tests for Repository
```go
func TestUserRepository_Save(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := NewUserRepository(db)
    
    user := &domain.User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Act
    err := repo.Save(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    
    // Verify in database
    found, err := repo.FindByID(context.Background(), user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Name, found.Name)
}
```

## 🔒 Benefits of Clean Architecture

### 1. **Framework Independence**
```go
// You can change from Gin to Echo without affecting business logic
// internal/handler/http/ ← Only this layer changes
// internal/usecase/     ← No changes
// internal/domain/      ← No changes
```

### 2. **Database Independence**
```go
// You can change from PostgreSQL to MongoDB
// internal/repository/postgres/ → internal/repository/mongo/
// internal/usecase/             ← No changes (uses interfaces)
// internal/domain/              ← No changes
```

### 3. **UI Independence**
```go
// You can add gRPC without affecting REST
// internal/handler/http/  ← Existing
// internal/handler/grpc/  ← New
// internal/usecase/       ← No changes
// internal/domain/        ← No changes
```

### 4. **Complete Testability**
- **Unit tests** for domain entities
- **Unit tests** for use cases (with mocks)
- **Integration tests** for repositories
- **End-to-end tests** for handlers

### 5. **Maintainability**
- Changes in one layer don't affect others
- Predictable and well-organized code
- Easy to add new functionalities
- Safe refactoring

## 🚫 Anti-Patterns That Goca Prevents

### ❌ Fat Controller
```go
// BAD - All logic in the handler
func (h *UserHandler) Create(c *gin.Context) {
    // Parsing
    var req CreateUserRequest
    c.ShouldBindJSON(&req)
    
    // Validation
    if req.Name == "" { /* ... */ }
    
    // Business logic
    if len(req.Name) < 2 { /* ... */ }
    
    // Database
    db.Query("INSERT INTO users...")
    
    // Response
    c.JSON(200, user)
}
```

```go
// ✅ GOOD - Handler delegating responsibilities
func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserHTTPRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Error: err.Error()})
        return
    }
    
    // Delegate to use case
    useCaseReq := usecase.CreateUserRequest{
        Name:  req.Name,
        Email: req.Email,
    }
    
    user, err := h.userUseCase.Create(c.Request.Context(), useCaseReq)
    if err != nil {
        c.JSON(500, ErrorResponse{Error: err.Error()})
        return
    }
    
    c.JSON(201, UserHTTPResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    })
}
```

### ❌ Anemic Domain Model
```go
// BAD - Entity without behavior
type User struct {
    ID    uint
    Name  string
    Email string
}

// Logic in the service
func (s *UserService) ValidateUser(user User) error {
    if user.Name == "" {
        return errors.New("name required")
    }
    // ...
}
```

```go
// ✅ GOOD - Rich entity with behavior
type User struct {
    ID    uint
    Name  string
    Email string
}

// Behavior in the entity
func (u *User) Validate() error {
    if u.Name == "" {
        return ErrUserNameRequired
    }
    return nil
}

func (u *User) CanUpdateProfile() bool {
    return u.ID > 0
}
```

### ❌ God Object
```go
// BAD - One class does everything
type UserManager struct {
    db     *sql.DB
    logger *log.Logger
    cache  *redis.Client
}

func (um *UserManager) CreateUser(data string) error {
    // Parse JSON
    // Validate data
    // Check business rules
    // Save to database
    // Update cache
    // Send email
    // Log action
    // Return response
}
```

```go
// ✅ GOOD - Separate responsibilities
type UserUseCase struct {
    userRepo UserRepository
}

type EmailService struct {
    sender EmailSender
}

type UserHandler struct {
    userUseCase UserUseCase
}
```

## 📊 Quality Metrics

### Complexity by Layer
- **Domain**: High business complexity, low technical complexity
- **UseCase**: Medium orchestration complexity
- **Handlers**: Low complexity, only adaptation
- **Repository**: Low complexity, only persistence

### Coupling
- **Low coupling** between layers (only interfaces)
- **High coupling** within each layer (cohesion)

### Testability
- **100% testable** without external dependencies
- **Easy mocks** by using interfaces
- **Fast tests** without I/O in unit tests

---

**← [Project Structure](Project-Structure) | [Implemented Patterns](Design-Patterns) →**
