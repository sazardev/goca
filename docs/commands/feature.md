# goca feature

Generate a complete feature with all Clean Architecture layers in a single command.

## Syntax

```bash
goca feature <FeatureName> [flags]
```

## Description

The `goca feature` command is the **fastest way** to add new functionality to your project. It generates all layers (domain, use case, repository, handler) and automatically integrates them with dependency injection and routing.

::: tip Recommended Workflow
This is the recommended command for most use cases. It saves time and ensures all layers are properly connected.
:::

## Arguments

### `<FeatureName>`

**Required.** The name of your feature (singular, PascalCase).

```bash
goca feature User
goca feature Product
goca feature Order
```

## Flags

### `--fields`

Define the structure of your entity.

**Format:** `"field1:type1,field2:type2,..."`

**Supported Types:**
- `string` - Text data
- `int`, `int64` - Integer numbers
- `float64` - Decimal numbers
- `bool` - Boolean values
- `time.Time` - Timestamps
- `[]type` - Arrays/slices

```bash
goca feature Product --fields "name:string,price:float64,inStock:bool"
```

### `--validation`

Add domain-level validation rules.

```bash
goca feature User --fields "name:string,email:string,age:int" --validation
```

Generates validation methods like:
- Email format validation
- Required field checks
- Range validations
- Custom business rules

### `--database`

Specify database type for repository. Default: `postgres`

**Options:** `postgres` | `mysql` | `mongodb` | `sqlite`

```bash
goca feature Order --fields "total:float64" --database mysql
```

### `--handlers`

Generate multiple handler types.

**Options:** `http` | `grpc` | `cli` | `worker` | `soap`

```bash
goca feature Payment --fields "amount:float64" --handlers "http,grpc"
```

## Examples

### Basic Feature

```bash
goca feature User --fields "name:string,email:string,age:int"
```

**Generates:**
```
internal/
├── domain/
│   ├── user.go              # Entity with business rules
│   └── user_errors.go       # Domain-specific errors
├── usecase/
│   ├── user_dto.go          # Input/Output DTOs
│   ├── user_interfaces.go   # Use case contracts
│   └── user_service.go      # Business logic implementation
├── repository/
│   ├── user_repository.go   # Repository interface
│   └── postgres_user_repository.go  # PostgreSQL implementation
└── handler/
    └── http/
        └── user_handler.go  # HTTP REST endpoints
```

### Feature with Validation

```bash
goca feature Product \
  --fields "name:string,price:float64,stock:int,category:string" \
  --validation
```

Generates entity with:
```go
func (p *Product) Validate() error {
    if p.Name == "" {
        return ErrProductNameRequired
    }
    if p.Price < 0 {
        return ErrInvalidPrice
    }
    if p.Stock < 0 {
        return ErrInvalidStock
    }
    return nil
}
```

### E-commerce Order Feature

```bash
goca feature Order \
  --fields "customerID:int,items:[]OrderItem,total:float64,status:string,createdAt:time.Time" \
  --validation \
  --database postgres
```

### Multi-Protocol Service

```bash
goca feature Payment \
  --fields "amount:float64,currency:string,method:string,status:string" \
  --handlers "http,grpc,worker" \
  --validation
```

### Complex Domain Model

```bash
goca feature Invoice \
  --fields "number:string,customerID:int,items:[]InvoiceItem,subtotal:float64,tax:float64,total:float64,dueDate:time.Time,status:string" \
  --validation \
  --database postgres
```

## What Gets Generated

### 1. Domain Layer (`internal/domain/`)

**Entity:** `<feature>.go`
```go
package domain

type User struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Age       int       `json:"age"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error { /* ... */ }
```

**Errors:** `<feature>_errors.go`
```go
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUserNameRequired  = errors.New("name is required")
    ErrInvalidEmail      = errors.New("invalid email")
)
```

### 2. Use Case Layer (`internal/usecase/`)

**DTOs:** `<feature>_dto.go`
```go
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"min=0,max=150"`
}

type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}
```

**Service:** `<feature>_service.go`
```go
type UserService interface {
    Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
    GetByID(ctx context.Context, id uint) (*UserResponse, error)
    Update(ctx context.Context, id uint, req UpdateUserRequest) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context) ([]*UserResponse, error)
}
```

### 3. Repository Layer (`internal/repository/`)

**Interface:** Repository contract in domain
**Implementation:** `postgres_<feature>_repository.go`
```go
type postgresUserRepository struct {
    db *sql.DB
}

func (r *postgresUserRepository) Save(ctx context.Context, user *domain.User) error {
    // Implementation
}
```

### 4. Handler Layer (`internal/handler/http/`)

**Handler:** `<feature>_handler.go`
```go
type UserHandler struct {
    service usecase.UserService
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    // HTTP handling
}
```

**Routes:** Automatically registered
```go
// POST   /api/v1/users
// GET    /api/v1/users/:id
// PUT    /api/v1/users/:id
// DELETE /api/v1/users/:id
// GET    /api/v1/users
```

## Automatic Integration

After generating a feature, Goca automatically:

 Updates dependency injection container  
 Registers HTTP routes  
 Adds database migrations  
 Configures repository connections  
 Wires all dependencies

**You can immediately test your new feature!**

```bash
# Generate feature
goca feature User --fields "name:string,email:string"

# Run server (feature is already integrated!)
go run cmd/server/main.go

# Test the API
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

## After Generation

### 1. Review Generated Code

Check the generated files and customize as needed:
```bash
# View generated files
find internal -name "*user*"
```

### 2. Add Business Logic

Enhance the domain entity with business rules:
```go
// internal/domain/user.go
func (u *User) CanDelete() bool {
    return u.Status != "active"
}

func (u *User) IsAdmin() bool {
    return u.Role == "admin"
}
```

### 3. Run Tests

```bash
go test ./internal/...
```

### 4. Run Application

```bash
go run cmd/server/main.go
```

## Tips

### Start Simple

Begin with basic fields, then add complexity:
```bash
# Start
goca feature Product --fields "name:string,price:float64"

# Later, manually add relationships and complex logic
```

### Use Consistent Naming

- Feature names: **Singular**, **PascalCase** (User, Product, Order)
- Fields: **camelCase** (firstName, productName)

### Field Naming Conventions

```bash
#  Good
--fields "firstName:string,lastName:string,emailAddress:string"

#  Avoid
--fields "first_name:string,Last-Name:string,EMAIL:string"
```

### Complex Types

For complex types, generate basic structure first, then modify:
```bash
goca feature Order --fields "customerID:int,total:float64"

# Then manually add:
# - items []OrderItem
# - metadata map[string]interface{}
# - custom types
```

## Troubleshooting

### Feature Already Exists

**Problem:** "feature already exists"

**Solution:** Choose a different name or delete existing files:
```bash
rm -rf internal/domain/user*
rm -rf internal/usecase/user*
# ... then regenerate
```

### Invalid Field Type

**Problem:** "unsupported field type"

**Solution:** Use supported types or add manually after generation.

### Integration Failed

**Problem:** Routes not working

**Solution:** Run integration command:
```bash
goca integrate --feature User
```

## See Also

- [`goca init`](/commands/init) - Initialize project
- [`goca integrate`](/commands/integrate) - Manual integration
- [`goca entity`](/commands/entity) - Generate entity only
- [Complete Tutorial](/tutorials/complete-tutorial) - Step-by-step guide

## Next Steps

-  Learn about [individual layer commands](/commands/)
-  Follow the [Complete Tutorial](/tutorials/complete-tutorial)
-  Understand [Clean Architecture](/guide/clean-architecture)
