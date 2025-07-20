# Clean Architecture

Esta página explica en detalle cómo Goca implementa y hace cumplir los principios de **Clean Architecture** de Uncle Bob (Robert C. Martin) en proyectos Go.

## 🎯 ¿Qué es Clean Architecture?

Clean Architecture es un patrón arquitectónico que organiza el código en **capas concéntricas** donde las dependencias apuntan hacia el centro del sistema, garantizando:

- 🔒 **Independencia de frameworks**
- 🧪 **Testabilidad completa**
- 🎨 **Independencia de UI**
- 💾 **Independencia de base de datos**
- 🌐 **Independencia de agentes externos**

## 🏗️ Las 4 Capas de Clean Architecture

### 🟡 1. Capa de Dominio (Entities)
**Ubicación**: `internal/domain/`  
**Responsabilidad**: Lógica de negocio central y reglas empresariales

#### ✅ Lo que SÍ debe estar aquí:
- Entidades del negocio
- Reglas de negocio fundamentales
- Validaciones de dominio
- Interfaces de repositorios
- Errores específicos del dominio

#### ❌ Lo que NO debe estar aquí:
- Dependencias externas (bases de datos, APIs)
- Lógica de presentación
- Detalles de implementación
- Frameworks o librerías externas

#### 📄 Ejemplo de Entidad:
```go
package domain

import (
    "errors"
    "strings"
    "time"
)

// User representa un usuario en el sistema
type User struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Validate implementa las reglas de negocio para validar un usuario
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

// CanUpdateProfile verifica si el usuario puede actualizar su perfil
func (u *User) CanUpdateProfile() bool {
    return u.ID > 0 && u.Name != ""
}

// isValidEmail valida el formato del email (regla de negocio)
func (u *User) isValidEmail() bool {
    return strings.Contains(u.Email, "@") && 
           strings.Contains(u.Email, ".") &&
           len(u.Email) > 5
}

// Errores del dominio
var (
    ErrUserNameRequired  = errors.New("user name is required")
    ErrUserNameTooShort  = errors.New("user name must be at least 2 characters")
    ErrUserEmailInvalid  = errors.New("user email format is invalid")
    ErrUserNotFound      = errors.New("user not found")
)
```

### 🔴 2. Capa de Casos de Uso (Use Cases)
**Ubicación**: `internal/usecase/`  
**Responsabilidad**: Lógica de aplicación y orquestación

#### ✅ Lo que SÍ debe estar aquí:
- DTOs (Data Transfer Objects)
- Interfaces de casos de uso
- Servicios de aplicación
- Validaciones de entrada
- Coordinación entre repositorios

#### ❌ Lo que NO debe estar aquí:
- Lógica de presentación
- Detalles de base de datos
- Lógica de frameworks web
- Implementaciones específicas de infraestructura

#### 📄 Ejemplo de Caso de Uso:
```go
package usecase

import (
    "context"
    "github.com/usuario/proyecto/internal/domain"
)

// UserUseCase define los contratos para casos de uso de usuario
type UserUseCase interface {
    Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
    GetByID(ctx context.Context, id uint) (*UserResponse, error)
    Update(ctx context.Context, id uint, req UpdateUserRequest) (*UserResponse, error)
    Delete(ctx context.Context, id uint) error
}

// UserRepository define los contratos para persistencia
type UserRepository interface {
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id uint) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uint) error
}

// userUseCase implementa la lógica de aplicación
type userUseCase struct {
    userRepo UserRepository
}

// NewUserUseCase crea una nueva instancia del caso de uso
func NewUserUseCase(userRepo UserRepository) UserUseCase {
    return &userUseCase{
        userRepo: userRepo,
    }
}

// Create crea un nuevo usuario
func (uc *userUseCase) Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
    // 1. Validar DTO de entrada
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // 2. Crear entidad de dominio
    user := &domain.User{
        Name:  req.Name,
        Email: req.Email,
    }
    
    // 3. Validar reglas de negocio
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    // 4. Verificar reglas de aplicación (email único)
    existingUser, _ := uc.userRepo.FindByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, domain.ErrUserEmailAlreadyExists
    }
    
    // 5. Persistir
    if err := uc.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // 6. Retornar DTO de respuesta
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

// CreateUserRequest DTO para crear usuario
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
}

// Validate valida el DTO
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

// UserResponse DTO de respuesta
type UserResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 🟢 3. Capa de Adaptadores (Interface Adapters)
**Ubicación**: `internal/handler/`  
**Responsabilidad**: Adaptar entrada/salida entre protocolos y casos de uso

#### ✅ Lo que SÍ debe estar aquí:
- Handlers HTTP/gRPC/CLI
- Controladores REST
- Adaptadores de protocolos
- DTOs específicos por protocolo
- Middlewares

#### ❌ Lo que NO debe estar aquí:
- Lógica de negocio
- Acceso directo a base de datos
- Validaciones de negocio
- Reglas empresariales

#### 📄 Ejemplo de Handler HTTP:
```go
package http

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/usuario/proyecto/internal/usecase"
)

// UserHandler maneja peticiones HTTP para usuarios
type UserHandler struct {
    userUseCase usecase.UserUseCase
}

// NewUserHandler crea un nuevo handler
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
    return &UserHandler{
        userUseCase: userUseCase,
    }
}

// Create maneja POST /users
func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserHTTPRequest
    
    // 1. Parsear entrada HTTP
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: "Invalid request format",
        })
        return
    }
    
    // 2. Convertir a DTO de caso de uso
    useCaseReq := usecase.CreateUserRequest{
        Name:  req.Name,
        Email: req.Email,
    }
    
    // 3. Ejecutar caso de uso
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
    
    // 4. Convertir a respuesta HTTP
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

### 🔵 4. Capa de Infraestructura (Frameworks & Drivers)
**Ubicación**: `internal/repository/`, `pkg/`  
**Responsabilidad**: Implementaciones específicas de tecnología

#### ✅ Lo que SÍ debe estar aquí:
- Implementaciones de repositorios
- Conexiones a bases de datos
- Clientes HTTP
- Configuración
- Logging
- Caches

#### ❌ Lo que NO debe estar aquí:
- Lógica de negocio
- Reglas empresariales
- Validaciones de dominio
- DTOs de casos de uso

#### 📄 Ejemplo de Repositorio:
```go
package postgres

import (
    "context"
    "database/sql"
    
    "github.com/usuario/proyecto/internal/domain"
)

// userRepository implementa el repositorio para PostgreSQL
type userRepository struct {
    db *sql.DB
}

// NewUserRepository crea un nuevo repositorio
func NewUserRepository(db *sql.DB) domain.UserRepository {
    return &userRepository{
        db: db,
    }
}

// Save implementa la persistencia específica de PostgreSQL
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

// FindByID busca un usuario por ID
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

## 🔄 Flujo de Dependencias

```
🌐 HTTP Request
     ↓
🟢 Handler (convierte HTTP → DTO)
     ↓
🔴 UseCase (ejecuta lógica de aplicación)
     ↓
🟡 Domain (valida reglas de negocio)
     ↓
🔵 Repository (persiste en base de datos)
```

### Regla de Dependencias:
> **Las dependencias SIEMPRE apuntan hacia adentro**

- 🟢 Handler depende de 🔴 UseCase
- 🔴 UseCase depende de 🟡 Domain
- 🔵 Repository implementa interfaces de 🟡 Domain
- 🟡 Domain NO depende de nada externo

## 🎭 Principios SOLID Aplicados

### 🔵 Single Responsibility Principle (SRP)
Cada clase tiene una sola razón para cambiar:

```go
// ✅ BIEN - Una responsabilidad
type UserValidator struct{}
func (v *UserValidator) Validate(user *User) error { /* ... */ }

// ✅ BIEN - Una responsabilidad
type UserRepository struct{}
func (r *UserRepository) Save(user *User) error { /* ... */ }

// ❌ MAL - Múltiples responsabilidades
type UserService struct{}
func (s *UserService) ValidateAndSave(user *User) error {
    // Validación + Persistencia = 2 responsabilidades
}
```

### 🔓 Open/Closed Principle (OCP)
Abierto para extensión, cerrado para modificación:

```go
// Interface estable
type NotificationSender interface {
    Send(message string) error
}

// Implementaciones extensibles
type EmailSender struct{} // Nueva implementación
type SMSSender struct{}   // Nueva implementación
type SlackSender struct{} // Nueva implementación

// UseCase cerrado para modificación
type UserUseCase struct {
    notifier NotificationSender // Usa interface
}
```

### 🔄 Liskov Substitution Principle (LSP)
Las implementaciones deben ser intercambiables:

```go
// Cualquier implementación de UserRepository
// debe comportarse igual desde el punto de vista del UseCase
type PostgreSQLUserRepo struct{}
type MySQLUserRepo struct{}
type MongoUserRepo struct{}

// Todas implementan la misma interface
type UserRepository interface {
    Save(user *User) error
    FindByID(id uint) (*User, error)
}
```

### 🎯 Interface Segregation Principle (ISP)
Interfaces específicas y cohesivas:

```go
// ✅ BIEN - Interfaces específicas
type UserReader interface {
    FindByID(id uint) (*User, error)
}

type UserWriter interface {
    Save(user *User) error
}

// ❌ MAL - Interface demasiado grande
type UserRepository interface {
    Save(user *User) error
    FindByID(id uint) (*User, error)
    FindByEmail(email string) (*User, error)
    Update(user *User) error
    Delete(id uint) error
    FindAll() ([]*User, error)
    Count() (int, error)
    // ... muchos más métodos
}
```

### ⬇️ Dependency Inversion Principle (DIP)
Depender de abstracciones, no de concreciones:

```go
// ✅ BIEN - Depende de interface (abstracción)
type UserUseCase struct {
    userRepo UserRepository // Interface
}

// ❌ MAL - Depende de implementación concreta
type UserUseCase struct {
    userRepo *PostgreSQLUserRepository // Implementación específica
}
```

## 🧪 Testabilidad

Clean Architecture facilita enormemente el testing:

### Unit Tests para Dominio
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

### Unit Tests para UseCase con Mocks
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

### Integration Tests para Repository
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

## 🔒 Beneficios de Clean Architecture

### 1. **Independencia de Frameworks**
```go
// Puedes cambiar de Gin a Echo sin afectar la lógica de negocio
// internal/handler/http/ ← Solo esta capa cambia
// internal/usecase/     ← Sin cambios
// internal/domain/      ← Sin cambios
```

### 2. **Independencia de Base de Datos**
```go
// Puedes cambiar de PostgreSQL a MongoDB
// internal/repository/postgres/ → internal/repository/mongo/
// internal/usecase/             ← Sin cambios (usa interfaces)
// internal/domain/              ← Sin cambios
```

### 3. **Independencia de UI**
```go
// Puedes agregar gRPC sin afectar REST
// internal/handler/http/  ← Existente
// internal/handler/grpc/  ← Nuevo
// internal/usecase/       ← Sin cambios
// internal/domain/        ← Sin cambios
```

### 4. **Testabilidad Completa**
- **Unit tests** para entidades de dominio
- **Unit tests** para casos de uso (con mocks)
- **Integration tests** para repositorios
- **End-to-end tests** para handlers

### 5. **Mantenibilidad**
- Cambios en una capa no afectan otras
- Código predecible y bien organizado
- Fácil agregar nuevas funcionalidades
- Refactoring seguro

## 🚫 Anti-Patrones Que Goca Previene

### ❌ Fat Controller
```go
// MAL - Toda la lógica en el handler
func (h *UserHandler) Create(c *gin.Context) {
    // Parsing
    var req CreateUserRequest
    c.ShouldBindJSON(&req)
    
    // Validación
    if req.Name == "" { /* ... */ }
    
    // Lógica de negocio
    if len(req.Name) < 2 { /* ... */ }
    
    // Base de datos
    db.Query("INSERT INTO users...")
    
    // Respuesta
    c.JSON(200, user)
}
```

```go
// ✅ BIEN - Handler delegando responsabilidades
func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserHTTPRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Error: err.Error()})
        return
    }
    
    // Delegar al caso de uso
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
// MAL - Entidad sin comportamiento
type User struct {
    ID    uint
    Name  string
    Email string
}

// Lógica en el servicio
func (s *UserService) ValidateUser(user User) error {
    if user.Name == "" {
        return errors.New("name required")
    }
    // ...
}
```

```go
// ✅ BIEN - Entidad rica con comportamiento
type User struct {
    ID    uint
    Name  string
    Email string
}

// Comportamiento en la entidad
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
// MAL - Una clase hace todo
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
// ✅ BIEN - Responsabilidades separadas
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

## 📊 Métricas de Calidad

### Complejidad por Capa
- **Dominio**: Alta complejidad de negocio, baja complejidad técnica
- **UseCase**: Media complejidad de orquestación
- **Handlers**: Baja complejidad, solo adaptación
- **Repository**: Baja complejidad, solo persistencia

### Acoplamiento
- **Bajo acoplamiento** entre capas (solo interfaces)
- **Alto acoplamiento** dentro de cada capa (cohesión)

### Testabilidad
- **100% testeable** sin dependencias externas
- **Mocks fáciles** por uso de interfaces
- **Tests rápidos** sin I/O en unit tests

---

**← [Estructura de Proyecto](Project-Structure) | [Patrones Implementados](Design-Patterns) →**
