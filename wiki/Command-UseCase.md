# goca usecase Command

The `goca usecase` command generates application services with well-defined DTOs, clear interfaces and encapsulated business logic following Clean Architecture principles.

## 📋 Syntax

```bash
goca usecase <name> [flags]
```

## 🎯 Purpose

Creates use cases (application services) that coordinate between domain and infrastructure:

- 🔴 **Application logic** without external dependencies
- 📄 **Specific DTOs** for input and output
- 🔗 **Clear interfaces** for repositories
- ⚡ **Asynchronous operations** (optional)
- ✅ **DTO validations** (optional)
- 🔄 **Configurable CRUD operations**

## 🚩 Available Flags

| Flag               | Type     | Required  | Default Value | Description                                        |
| ------------------ | -------- | --------- | ------------- | -------------------------------------------------- |
| `--entity`         | `string` | ✅ **Yes** | -             | Entity associated with the use case                |
| `--operations`     | `string` | ❌ No      | `create,read` | CRUD operations (`create,read,update,delete,list`) |
| `--dto-validation` | `bool`   | ❌ No      | `false`       | DTOs with specific validations                     |
| `--async`          | `bool`   | ❌ No      | `false`       | Include asynchronous operations                    |

## 📖 Usage Examples

### Basic Use Case
```bash
goca usecase UserService --entity User
```

### With All CRUD Operations
```bash
goca usecase ProductService --entity Product --operations "create,read,update,delete,list"
```

### With DTO Validations
```bash
goca usecase OrderService --entity Order --operations "create,read,update" --dto-validation
```

### With Asynchronous Operations
```bash
goca usecase NotificationService --entity Notification --operations "create,list" --async --dto-validation
```

## 📂 Archivos Generados

### Estructura de Archivos
```
internal/usecase/
├── user_service.go          # Implementación del caso de uso
├── interfaces/
│   └── user_interfaces.go   # Interfaces de contratos
└── dto/
    └── user_dto.go          # DTOs de entrada y salida
```

## 🔍 Código Generado en Detalle

### Caso de Uso Principal: `internal/usecase/user_service.go`

```go
package usecase

import (
    "context"
    
    "github.com/usuario/proyecto/internal/domain"
    "github.com/usuario/proyecto/internal/usecase/dto"
    "github.com/usuario/proyecto/internal/usecase/interfaces"
)

// UserService implementa la lógica de aplicación para usuarios
type UserService struct {
    userRepo interfaces.UserRepository
}

// NewUserService crea una nueva instancia del servicio
func NewUserService(userRepo interfaces.UserRepository) *UserService {
    return &UserService{
        userRepo: userRepo,
    }
}

// Create crea un nuevo usuario
func (s *UserService) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
    // Validar DTO
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // Verificar que el email no exista
    existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, domain.ErrUserEmailAlreadyExists
    }
    
    // Crear entidad de dominio
    user := &domain.User{
        Name:  req.Name,
        Email: req.Email,
    }
    
    // Validar entidad
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    // Guardar en repositorio
    if err := s.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // Retornar DTO de respuesta
    return &dto.UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }, nil
}

// GetByID obtiene un usuario por ID
func (s *UserService) GetByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        return nil, domain.ErrUserNotFound
    }
    
    return &dto.UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }, nil
}
```

### DTOs: `internal/usecase/dto/user_dto.go`

```go
package dto

import (
    "errors"
    "strings"
    "time"
)

// CreateUserRequest DTO para crear usuario
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
}

// Validate valida el DTO de creación
func (r CreateUserRequest) Validate() error {
    if strings.TrimSpace(r.Name) == "" {
        return errors.New("name is required")
    }
    
    if len(r.Name) < 2 {
        return errors.New("name must be at least 2 characters")
    }
    
    if len(r.Name) > 100 {
        return errors.New("name must be less than 100 characters")
    }
    
    if strings.TrimSpace(r.Email) == "" {
        return errors.New("email is required")
    }
    
    if !strings.Contains(r.Email, "@") {
        return errors.New("invalid email format")
    }
    
    return nil
}

// UpdateUserRequest DTO para actualizar usuario
type UpdateUserRequest struct {
    Name  *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
    Email *string `json:"email,omitempty" validate:"omitempty,email"`
}

// UserResponse DTO de respuesta de usuario
type UserResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// ListUsersRequest DTO para listar usuarios
type ListUsersRequest struct {
    Page   int    `json:"page" validate:"min=1"`
    Limit  int    `json:"limit" validate:"min=1,max=100"`
    Search string `json:"search,omitempty"`
}

// ListUsersResponse DTO de respuesta de listado
type ListUsersResponse struct {
    Users       []UserResponse `json:"users"`
    Total       int64          `json:"total"`
    Page        int            `json:"page"`
    Limit       int            `json:"limit"`
    TotalPages  int            `json:"total_pages"`
    HasNextPage bool           `json:"has_next_page"`
    HasPrevPage bool           `json:"has_prev_page"`
}
```

### Interfaces: `internal/usecase/interfaces/user_interfaces.go`

```go
package interfaces

import (
    "context"
    
    "github.com/usuario/proyecto/internal/domain"
    "github.com/usuario/proyecto/internal/usecase/dto"
)

//go:generate mockgen -source=user_interfaces.go -destination=mocks/user_interfaces_mock.go

// UserUseCase define los contratos para los casos de uso de usuario
type UserUseCase interface {
    Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
    GetByID(ctx context.Context, id uint) (*dto.UserResponse, error)
    Update(ctx context.Context, id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, req dto.ListUsersRequest) (*dto.ListUsersResponse, error)
}

// UserRepository define los contratos para el repositorio de usuario
type UserRepository interface {
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id uint) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
}
```

## 🔧 Operaciones Disponibles

### Operaciones CRUD Básicas
```bash
--operations "create,read,update,delete,list"
```

- **create**: Crear nueva entidad
- **read**: Leer entidad por ID
- **update**: Actualizar entidad existente
- **delete**: Eliminar entidad
- **list**: Listar entidades con paginación

### Operaciones Personalizadas
Se pueden agregar operaciones específicas del dominio:

```go
// Con --operations "create,read,activate,deactivate"

// Activate activa un usuario
func (s *UserService) Activate(ctx context.Context, id uint) error {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return err
    }
    
    if user == nil {
        return domain.ErrUserNotFound
    }
    
    user.Activate()
    return s.userRepo.Update(ctx, user)
}
```

## ⚡ Operaciones Asíncronas (--async)

Con el flag `--async`, se generan métodos adicionales:

```go
// CreateAsync crea un usuario de forma asíncrona
func (s *UserService) CreateAsync(ctx context.Context, req dto.CreateUserRequest) error {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Error en CreateAsync: %v", r)
            }
        }()
        
        _, err := s.Create(context.Background(), req)
        if err != nil {
            log.Printf("Error creando usuario async: %v", err)
        }
    }()
    
    return nil
}

// ProcessBatch procesa múltiples usuarios en lotes
func (s *UserService) ProcessBatch(ctx context.Context, requests []dto.CreateUserRequest) error {
    for _, req := range requests {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            if err := s.CreateAsync(ctx, req); err != nil {
                return err
            }
        }
    }
    
    return nil
}
```

## ✅ Validaciones de DTOs (--dto-validation)

Con `--dto-validation`, los DTOs incluyen validaciones detalladas:

```go
// ValidateForUpdate valida campos para actualización
func (r UpdateUserRequest) ValidateForUpdate() error {
    if r.Name != nil {
        if strings.TrimSpace(*r.Name) == "" {
            return errors.New("name cannot be empty")
        }
        if len(*r.Name) < 2 || len(*r.Name) > 100 {
            return errors.New("name must be between 2 and 100 characters")
        }
    }
    
    if r.Email != nil {
        if !isValidEmail(*r.Email) {
            return errors.New("invalid email format")
        }
    }
    
    return nil
}

// ValidateForCreation valida campos requeridos para creación
func (r CreateUserRequest) ValidateForCreation() error {
    if err := r.Validate(); err != nil {
        return err
    }
    
    // Validaciones adicionales específicas para creación
    if strings.Contains(r.Email, "+") {
        return errors.New("email aliases not allowed")
    }
    
    return nil
}
```

## 🔄 Integración con Otros Comandos

### Flujo Completo de Desarrollo
```bash
# 1. Crear entidad
goca entity User --fields "name:string,email:string" --validation

# 2. Crear caso de uso
goca usecase UserService --entity User --operations "create,read,update,delete,list" --dto-validation

# 3. Crear repositorio
goca repository User --database postgres --transactions

# 4. Crear handler
goca handler User --type http --swagger --validation
```

## 🧪 Testing

Los casos de uso generados son fáciles de testear:

```go
func TestUserService_Create(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository(ctrl)
    service := NewUserService(mockRepo)
    
    req := dto.CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    mockRepo.EXPECT().
        FindByEmail(gomock.Any(), req.Email).
        Return(nil, nil)
    
    mockRepo.EXPECT().
        Save(gomock.Any(), gomock.Any()).
        Return(nil)
    
    result, err := service.Create(context.Background(), req)
    
    assert.NoError(t, err)
    assert.Equal(t, req.Name, result.Name)
    assert.Equal(t, req.Email, result.Email)
}
```

## ⚠️ Important Considerations

### ✅ Best Practices
- **Specific DTOs**: Different DTOs for input and output
- **Early validations**: Validate in DTOs before domain
- **Clear interfaces**: Define explicit contracts
- **Context propagation**: Use context.Context in all methods

### ❌ Common Errors
- **Business logic in use cases**: Should be in domain
- **Direct dependencies**: Don't access DB/HTTP directly
- **Anemic DTOs**: Without validations or behavior
- **Mixing concerns**: Mixing different responsibilities

---

**← [goca entity Command](Command-Entity) | [goca repository Command](Command-Repository) →**
