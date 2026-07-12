# Best Practices

Guidelines and recommendations for building with Goca.

## Code Organization

### Keep Layers Independent

Each layer should only depend on interfaces, not concrete implementations.

** Good:**
```go
type UserService struct {
    repo repository.UserRepository // Interface
}
```

** Bad:**
```go
type UserService struct {
    repo *repository.PostgresUserRepository // Concrete
}
```

### Use Dependency Injection

Let the DI container wire dependencies:

```go
func NewUserService(repo repository.UserRepository) usecase.UserService {
    return &userService{repo: repo}
}
```

### Single Responsibility

Each struct should have one clear purpose:

```go
//  One responsibility: handle HTTP requests
type UserHandler struct {
    service usecase.UserService
}

//  Too many responsibilities
type UserHandler struct {
    db      *gorm.DB
    cache   *redis.Client
    mailer  *smtp.Client
}
```

## Error Handling

### Use Domain Errors

Define errors in the domain layer:

```go
// internal/domain/errors.go
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidEmail = errors.New("invalid email format")
)
```

### Wrap Errors with Context

```go
func (s *userService) GetUser(ctx context.Context, id uint) (*UserResponse, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    return toUserResponse(user), nil
}
```

### Handle Errors at HTTP Layer

Goca-generated handlers use the standard library `net/http` together with
[gorilla/mux](https://github.com/gorilla/mux) — not Gin, Chi, or Echo. Path
parameters come from `mux.Vars(r)`, not `c.Param(...)`:

```go
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    user, err := h.usecase.GetUser(id)
    if errors.Is(err, domain.ErrUserNotFound) {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    if err != nil {
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

## Testing

### Use Table-Driven Tests

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserInput
        wantErr bool
    }{
        {
            name: "valid user",
            input: CreateUserInput{
                Name:  "John Doe",
                Email: "john@example.com",
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            input: CreateUserInput{
                Name:  "Jane Doe",
                Email: "invalid-email",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Mock External Dependencies

```go
type MockUserRepository struct {
    SaveFunc    func(ctx context.Context, user *domain.User) error
    FindByIDFunc func(ctx context.Context, id uint) (*domain.User, error)
}

func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
    return m.SaveFunc(ctx, user)
}
```

## Database Operations

### Use Transactions

```go
func (s *orderService) CreateOrder(ctx context.Context, input CreateOrderInput) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Create order
        if err := tx.Create(&order).Error; err != nil {
            return err
        }
        
        // Create order items
        if err := tx.Create(&items).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

### Use Context

Always pass context for cancellation and timeouts:

```go
func (r *PostgresUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    return &user, err
}
```

## API Design

### Use DTOs

Don't expose domain entities directly:

```go
//  Use DTOs
type CreateUserInput struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### Consistent Response Format

```go
type SuccessResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
}

type ErrorResponse struct {
    Success bool   `json:"success"`
    Error   string `json:"error"`
}
```

## Security

### Validate Input

Goca doesn't rely on framework-level binding tags (there is no Gin/`binding`
struct tag support). Validation is explicit, in the domain layer, via a
`Validate()` method that the generated `validations.go` / entity file calls
before persisting:

```go
type User struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

func (u *User) Validate() error {
    if u.Name == "" {
        return ErrInvalidUserName
    }
    if len(u.Name) < 2 || len(u.Name) > 100 {
        return ErrInvalidUserNameLength
    }
    if u.Email == "" || !strings.Contains(u.Email, "@") {
        return ErrInvalidUserEmail
    }
    return nil
}
```

### Don't Log Sensitive Data

```go
//  Bad - logs password
log.Printf("Creating user: %+v", input)

//  Good - excludes sensitive fields
log.Printf("Creating user: name=%s, email=%s", input.Name, input.Email)
```

## Performance

### Use Pagination

```go
// Read pagination params from the query string (no gin "form" binding tags —
// mux handlers parse r.URL.Query() directly).
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }
    // ...call usecase with page, pageSize...
}

func (r *postgresUserRepository) FindAll(page, pageSize int) ([]domain.User, error) {
    var users []domain.User
    offset := (page - 1) * pageSize
    err := r.db.
        Offset(offset).
        Limit(pageSize).
        Find(&users).Error
    return users, err
}
```

### Use Indexes

```go
type User struct {
    ID    uint   `gorm:"primaryKey"`
    Email string `gorm:"unique;index"` // Add index for lookups
    Name  string `gorm:"index"`        // Index for searches
}
```

## See Also

- [Project Structure](/guide/project-structure) - Directory organization
- [Clean Architecture](/guide/clean-architecture) - Architecture principles
