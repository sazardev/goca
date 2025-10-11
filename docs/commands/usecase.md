# goca usecase

Generate use cases (application services) with DTOs and business logic.

## Syntax

```bash
goca usecase <ServiceName> [flags]
```

## Description

Creates application layer services that orchestrate business workflows, coordinate repositories, and define clear input/output contracts through DTOs.

## Flags

### `--entity`

Associated domain entity.

```bash
goca usecase ProductService --entity Product
```

### `--operations`

CRUD operations to generate.

**Options:** `create`, `read`, `update`, `delete`, `list`

```bash
goca usecase UserService --entity User --operations "create,read,update,delete,list"
```

### `--dto-validation`

Include DTO validation tags.

```bash
goca usecase OrderService --entity Order --dto-validation
```

## Examples

### Basic Use Case

```bash
goca usecase ProductService --entity Product
```

### Complete CRUD

```bash
goca usecase UserService \
  --entity User \
  --operations "create,read,update,delete,list" \
  --dto-validation
```

## Generated Files

```
internal/usecase/
├── user_dto.go         # Input/Output DTOs
├── user_interfaces.go  # Service interface
└── user_service.go     # Service implementation
```

## Generated Code Example

```go
// user_dto.go
package usecase

type CreateUserInput struct {
    Name  string `json:"name" validate:"required,min=2"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"min=0,max=150"`
}

type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// user_interfaces.go
type UserService interface {
    CreateUser(ctx context.Context, input CreateUserInput) (*UserResponse, error)
    GetUser(ctx context.Context, id uint) (*UserResponse, error)
    UpdateUser(ctx context.Context, id uint, input UpdateUserInput) error
    DeleteUser(ctx context.Context, id uint) error
    ListUsers(ctx context.Context) ([]*UserResponse, error)
}

// user_service.go
type userService struct {
    userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) UserService {
    return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (*UserResponse, error) {
    // Validate input
    if err := input.Validate(); err != nil {
        return nil, err
    }
    
    // Create domain entity
    user := &domain.User{
        Name:  input.Name,
        Email: input.Email,
        Age:   input.Age,
    }
    
    // Validate domain rules
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    // Persist
    if err := s.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // Return DTO
    return &UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}
```

## Best Practices

### ✅ DO

- Define clear DTOs for each operation
- Validate input at use case boundary
- Coordinate multiple repositories
- Transform entities to DTOs

### ❌ DON'T

- Include HTTP/gRPC logic
- Write SQL queries
- Import framework packages
- Skip validation

## See Also

- [`goca entity`](/commands/entity) - Generate entities
- [`goca repository`](/commands/repository) - Generate repositories
- [`goca feature`](/commands/feature) - Generate complete feature
