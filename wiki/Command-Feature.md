# Comando goca feature

El comando `goca feature` es la herramienta más poderosa de Goca. Genera **todas las capas** de Clean Architecture para un feature completo de una sola vez, siguiendo las mejores prácticas y convenciones establecidas.

## 📋 Sintaxis

```bash
goca feature <name> [flags]
```

## 🎯 Propósito

Genera un feature completo con **todas las capas de Clean Architecture**:

- 🟡 **Dominio**: Entidad con validaciones y reglas de negocio
- 🔴 **Casos de Uso**: Servicios de aplicación con DTOs
- 🔵 **Repositorio**: Interfaz y implementación de persistencia
- 🟢 **Handlers**: Adaptadores para diferentes protocolos
- 📄 **Mensajes**: Constantes y mensajes de error/éxito

## 🚩 Flags Disponibles

| Flag               | Tipo     | Requerido | Valor por Defecto | Descripción                                                      |
| ------------------ | -------- | --------- | ----------------- | ---------------------------------------------------------------- |
| `--fields`         | `string` | ✅ **Sí**  | -                 | Campos de la entidad (`"name:string,email:string"`)              |
| `--database`       | `string` | ❌ No      | `postgres`        | Base de datos (`postgres`, `mysql`, `mongodb`)                   |
| `--handlers`       | `string` | ❌ No      | `http`            | Tipos de handlers (`http`, `grpc`, `cli`, `worker`, `http,grpc`) |
| `--validation`     | `bool`   | ❌ No      | `true`            | Incluir validaciones en entidad y DTOs                           |
| `--business-rules` | `bool`   | ❌ No      | `false`           | Generar métodos de reglas de negocio                             |

## 📖 Ejemplos de Uso

### Ejemplo Básico
```bash
goca feature User --fields "name:string,email:string"
```

### Feature con Validaciones
```bash
goca feature Product --fields "name:string,price:float64,category:string,stock:int" --validation --business-rules
```

### Feature Multi-Handler
```bash
goca feature Order --fields "user_id:int,total:float64,status:string" --handlers "http,grpc,worker"
```

### Feature Completo
```bash
goca feature Employee \
  --fields "name:string,email:string,department:string,salary:float64,hire_date:time.Time" \
  --validation \
  --business-rules \
  --database postgres \
  --handlers "http,grpc"
```

## 📂 Estructura Generada

Un feature completo genera esta estructura:

```
internal/
├── domain/
│   ├── user.go                  # Entidad de dominio
│   └── errors.go                # Errores específicos del dominio
├── usecase/
│   ├── dto/
│   │   └── user_dto.go          # DTOs para casos de uso
│   ├── interfaces/
│   │   └── user_interfaces.go   # Interfaces de contratos
│   └── user_usecase.go          # Implementación de casos de uso
├── repository/
│   ├── interfaces/
│   │   └── user_repository.go   # Interface del repositorio
│   └── postgres/
│       └── user_repository.go   # Implementación PostgreSQL
├── handler/
│   ├── http/
│   │   ├── user_handler.go      # Handler HTTP REST
│   │   ├── user_routes.go       # Definición de rutas
│   │   └── dto/
│   │       └── user_dto.go      # DTOs específicos para HTTP
│   ├── grpc/                    # (si se especifica)
│   │   ├── user.proto           # Definición Protocol Buffers
│   │   └── user_server.go       # Servidor gRPC
│   └── worker/                  # (si se especifica)
│       └── user_worker.go       # Worker para tareas en background
└── messages/
    ├── errors.go                # Mensajes de error
    └── responses.go             # Mensajes de respuesta
```

## 🔍 Análisis de Archivos Generados

### 🟡 Dominio: `internal/domain/user.go`

```go
package domain

import (
    "errors"
    "strings"
    "time"
)

// User representa la entidad de dominio del usuario
type User struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Validate valida la entidad User según las reglas de negocio
func (u *User) Validate() error {
    if strings.TrimSpace(u.Name) == "" {
        return ErrUserNameRequired
    }
    
    if len(u.Name) < 2 {
        return ErrUserNameTooShort
    }
    
    if len(u.Name) > 100 {
        return ErrUserNameTooLong
    }
    
    if strings.TrimSpace(u.Email) == "" {
        return ErrUserEmailRequired
    }
    
    if !isValidEmail(u.Email) {
        return ErrUserEmailInvalid
    }
    
    return nil
}

// isValidEmail valida si el email tiene formato correcto
func isValidEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// Business Rules (si --business-rules está activado)

// CanUpdateEmail verifica si el usuario puede actualizar su email
func (u *User) CanUpdateEmail() bool {
    return u.ID > 0
}

// IsEmailDomainAllowed verifica si el dominio del email está permitido
func (u *User) IsEmailDomainAllowed() bool {
    allowedDomains := []string{"gmail.com", "company.com", "example.com"}
    
    parts := strings.Split(u.Email, "@")
    if len(parts) != 2 {
        return false
    }
    
    domain := parts[1]
    for _, allowed := range allowedDomains {
        if domain == allowed {
            return true
        }
    }
    
    return false
}
```

### 🔴 Casos de Uso: `internal/usecase/user_usecase.go`

```go
package usecase

import (
    "context"
    
    "github.com/usuario/proyecto/internal/domain"
    "github.com/usuario/proyecto/internal/usecase/dto"
)

//go:generate mockgen -source=user_usecase.go -destination=mocks/user_usecase_mock.go

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
    FindAll(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
}

// userUseCase implementa UserUseCase
type userUseCase struct {
    userRepo UserRepository
}

// NewUserUseCase crea una nueva instancia de UserUseCase
func NewUserUseCase(userRepo UserRepository) UserUseCase {
    return &userUseCase{
        userRepo: userRepo,
    }
}

// Create crea un nuevo usuario
func (uc *userUseCase) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
    // Validar DTO
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // Verificar que el email no exista
    existingUser, _ := uc.userRepo.FindByEmail(ctx, req.Email)
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
    if err := uc.userRepo.Save(ctx, user); err != nil {
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
func (uc *userUseCase) GetByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
    user, err := uc.userRepo.FindByID(ctx, id)
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

// Update actualiza un usuario existente
func (uc *userUseCase) Update(ctx context.Context, id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
    // Validar DTO
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // Buscar usuario existente
    user, err := uc.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        return nil, domain.ErrUserNotFound
    }
    
    // Actualizar campos
    if req.Name != nil {
        user.Name = *req.Name
    }
    
    if req.Email != nil {
        // Verificar que el nuevo email no exista
        if *req.Email != user.Email {
            existingUser, _ := uc.userRepo.FindByEmail(ctx, *req.Email)
            if existingUser != nil && existingUser.ID != user.ID {
                return nil, domain.ErrUserEmailAlreadyExists
            }
        }
        user.Email = *req.Email
    }
    
    // Validar entidad actualizada
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    // Guardar cambios
    if err := uc.userRepo.Update(ctx, user); err != nil {
        return nil, err
    }
    
    return &dto.UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }, nil
}

// Delete elimina un usuario
func (uc *userUseCase) Delete(ctx context.Context, id uint) error {
    // Verificar que el usuario existe
    user, err := uc.userRepo.FindByID(ctx, id)
    if err != nil {
        return err
    }
    
    if user == nil {
        return domain.ErrUserNotFound
    }
    
    // Eliminar usuario
    return uc.userRepo.Delete(ctx, id)
}

// List obtiene una lista paginada de usuarios
func (uc *userUseCase) List(ctx context.Context, req dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
    // Validar parámetros de paginación
    if req.Page < 1 {
        req.Page = 1
    }
    if req.Limit < 1 || req.Limit > 100 {
        req.Limit = 10
    }
    
    offset := (req.Page - 1) * req.Limit
    
    // Obtener usuarios
    users, total, err := uc.userRepo.FindAll(ctx, req.Limit, offset)
    if err != nil {
        return nil, err
    }
    
    // Convertir a DTOs
    userResponses := make([]dto.UserResponse, len(users))
    for i, user := range users {
        userResponses[i] = dto.UserResponse{
            ID:        user.ID,
            Name:      user.Name,
            Email:     user.Email,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
        }
    }
    
    totalPages := (total + int64(req.Limit) - 1) / int64(req.Limit)
    
    return &dto.ListUsersResponse{
        Users:       userResponses,
        Total:       total,
        Page:        req.Page,
        Limit:       req.Limit,
        TotalPages:  int(totalPages),
        HasNextPage: req.Page < int(totalPages),
        HasPrevPage: req.Page > 1,
    }, nil
}
```

### 🔴 DTOs: `internal/usecase/dto/user_dto.go`

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
func (r *CreateUserRequest) Validate() error {
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
        return errors.New("email format is invalid")
    }
    
    return nil
}

// UpdateUserRequest DTO para actualizar usuario
type UpdateUserRequest struct {
    Name  *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
    Email *string `json:"email,omitempty" validate:"omitempty,email"`
}

// Validate valida el DTO de actualización
func (r *UpdateUserRequest) Validate() error {
    if r.Name != nil {
        if strings.TrimSpace(*r.Name) == "" {
            return errors.New("name cannot be empty")
        }
        
        if len(*r.Name) < 2 {
            return errors.New("name must be at least 2 characters")
        }
        
        if len(*r.Name) > 100 {
            return errors.New("name must be less than 100 characters")
        }
    }
    
    if r.Email != nil {
        if strings.TrimSpace(*r.Email) == "" {
            return errors.New("email cannot be empty")
        }
        
        if !strings.Contains(*r.Email, "@") {
            return errors.New("email format is invalid")
        }
    }
    
    return nil
}

// UserResponse DTO de respuesta para usuario
type UserResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// ListUsersRequest DTO para listar usuarios
type ListUsersRequest struct {
    Page  int `json:"page" validate:"min=1"`
    Limit int `json:"limit" validate:"min=1,max=100"`
}

// ListUsersResponse DTO de respuesta para lista de usuarios
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

### 🔵 Repositorio: `internal/repository/postgres/user_repository.go`

```go
package postgres

import (
    "context"
    "database/sql"
    "time"
    
    "github.com/usuario/proyecto/internal/domain"
)

// userRepository implementa el repositorio de usuarios para PostgreSQL
type userRepository struct {
    db *sql.DB
}

// NewUserRepository crea una nueva instancia del repositorio
func NewUserRepository(db *sql.DB) domain.UserRepository {
    return &userRepository{
        db: db,
    }
}

// Save guarda un nuevo usuario en la base de datos
func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (name, email, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
    
    now := time.Now()
    user.CreatedAt = now
    user.UpdatedAt = now
    
    err := r.db.QueryRowContext(ctx, query, user.Name, user.Email, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
    if err != nil {
        return err
    }
    
    return nil
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
            return nil, nil
        }
        return nil, err
    }
    
    return user, nil
}

// FindByEmail busca un usuario por email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE email = $1
    `
    
    user := &domain.User{}
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    
    return user, nil
}

// Update actualiza un usuario existente
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
    query := `
        UPDATE users
        SET name = $2, email = $3, updated_at = $4
        WHERE id = $1
    `
    
    user.UpdatedAt = time.Now()
    
    _, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.UpdatedAt)
    return err
}

// Delete elimina un usuario por ID
func (r *userRepository) Delete(ctx context.Context, id uint) error {
    query := `DELETE FROM users WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}

// FindAll obtiene una lista paginada de usuarios
func (r *userRepository) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
    // Contar total de registros
    var total int64
    countQuery := `SELECT COUNT(*) FROM users`
    err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
    if err != nil {
        return nil, 0, err
    }
    
    // Obtener usuarios paginados
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
    
    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()
    
    var users []*domain.User
    for rows.Next() {
        user := &domain.User{}
        err := rows.Scan(
            &user.ID,
            &user.Name,
            &user.Email,
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            return nil, 0, err
        }
        users = append(users, user)
    }
    
    if err = rows.Err(); err != nil {
        return nil, 0, err
    }
    
    return users, total, nil
}
```

### 🟢 Handler HTTP: `internal/handler/http/user_handler.go`

```go
package http

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/usuario/proyecto/internal/usecase"
    "github.com/usuario/proyecto/internal/usecase/dto"
    "github.com/usuario/proyecto/internal/handler/http/dto"
)

// UserHandler maneja las peticiones HTTP para usuarios
type UserHandler struct {
    userUseCase usecase.UserUseCase
}

// NewUserHandler crea una nueva instancia del handler
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
    return &UserHandler{
        userUseCase: userUseCase,
    }
}

// Create maneja la creación de usuarios
// @Summary Crear usuario
// @Description Crea un nuevo usuario en el sistema
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "Datos del usuario"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
    var req dto.CreateUserRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "Invalid request format",
            Message: err.Error(),
        })
        return
    }
    
    // Convertir a DTO de caso de uso
    useCaseReq := dto.CreateUserRequest{
        Name:  req.Name,
        Email: req.Email,
    }
    
    userResponse, err := h.userUseCase.Create(c.Request.Context(), useCaseReq)
    if err != nil {
        status := http.StatusInternalServerError
        if err == domain.ErrUserEmailAlreadyExists {
            status = http.StatusConflict
        }
        
        c.JSON(status, dto.ErrorResponse{
            Error:   "Failed to create user",
            Message: err.Error(),
        })
        return
    }
    
    // Convertir a DTO de respuesta HTTP
    response := dto.UserResponse{
        ID:        userResponse.ID,
        Name:      userResponse.Name,
        Email:     userResponse.Email,
        CreatedAt: userResponse.CreatedAt,
        UpdatedAt: userResponse.UpdatedAt,
    }
    
    c.JSON(http.StatusCreated, response)
}

// GetByID maneja la obtención de un usuario por ID
// @Summary Obtener usuario por ID
// @Description Obtiene un usuario específico por su ID
// @Tags users
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "Invalid user ID",
            Message: "User ID must be a valid number",
        })
        return
    }
    
    userResponse, err := h.userUseCase.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        status := http.StatusInternalServerError
        if err == domain.ErrUserNotFound {
            status = http.StatusNotFound
        }
        
        c.JSON(status, dto.ErrorResponse{
            Error:   "Failed to get user",
            Message: err.Error(),
        })
        return
    }
    
    response := dto.UserResponse{
        ID:        userResponse.ID,
        Name:      userResponse.Name,
        Email:     userResponse.Email,
        CreatedAt: userResponse.CreatedAt,
        UpdatedAt: userResponse.UpdatedAt,
    }
    
    c.JSON(http.StatusOK, response)
}

// Update maneja la actualización de usuarios
// @Summary Actualizar usuario
// @Description Actualiza un usuario existente
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param user body dto.UpdateUserRequest true "Datos a actualizar"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "Invalid user ID",
            Message: "User ID must be a valid number",
        })
        return
    }
    
    var req dto.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "Invalid request format",
            Message: err.Error(),
        })
        return
    }
    
    // Convertir a DTO de caso de uso
    useCaseReq := dto.UpdateUserRequest{
        Name:  req.Name,
        Email: req.Email,
    }
    
    userResponse, err := h.userUseCase.Update(c.Request.Context(), uint(id), useCaseReq)
    if err != nil {
        status := http.StatusInternalServerError
        switch err {
        case domain.ErrUserNotFound:
            status = http.StatusNotFound
        case domain.ErrUserEmailAlreadyExists:
            status = http.StatusConflict
        }
        
        c.JSON(status, dto.ErrorResponse{
            Error:   "Failed to update user",
            Message: err.Error(),
        })
        return
    }
    
    response := dto.UserResponse{
        ID:        userResponse.ID,
        Name:      userResponse.Name,
        Email:     userResponse.Email,
        CreatedAt: userResponse.CreatedAt,
        UpdatedAt: userResponse.UpdatedAt,
    }
    
    c.JSON(http.StatusOK, response)
}

// Delete maneja la eliminación de usuarios
// @Summary Eliminar usuario
// @Description Elimina un usuario del sistema
// @Tags users
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "Invalid user ID",
            Message: "User ID must be a valid number",
        })
        return
    }
    
    err = h.userUseCase.Delete(c.Request.Context(), uint(id))
    if err != nil {
        status := http.StatusInternalServerError
        if err == domain.ErrUserNotFound {
            status = http.StatusNotFound
        }
        
        c.JSON(status, dto.ErrorResponse{
            Error:   "Failed to delete user",
            Message: err.Error(),
        })
        return
    }
    
    c.Status(http.StatusNoContent)
}

// List maneja la obtención de una lista paginada de usuarios
// @Summary Listar usuarios
// @Description Obtiene una lista paginada de usuarios
// @Tags users
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Límite por página" default(10)
// @Success 200 {object} dto.ListUsersResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
    page := 1
    limit := 10
    
    if pageStr := c.Query("page"); pageStr != "" {
        if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
            page = p
        }
    }
    
    if limitStr := c.Query("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
            limit = l
        }
    }
    
    req := dto.ListUsersRequest{
        Page:  page,
        Limit: limit,
    }
    
    usersResponse, err := h.userUseCase.List(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Error:   "Failed to list users",
            Message: err.Error(),
        })
        return
    }
    
    // Convertir usuarios a DTOs HTTP
    httpUsers := make([]dto.UserResponse, len(usersResponse.Users))
    for i, user := range usersResponse.Users {
        httpUsers[i] = dto.UserResponse{
            ID:        user.ID,
            Name:      user.Name,
            Email:     user.Email,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
        }
    }
    
    response := dto.ListUsersResponse{
        Users:       httpUsers,
        Total:       usersResponse.Total,
        Page:        usersResponse.Page,
        Limit:       usersResponse.Limit,
        TotalPages:  usersResponse.TotalPages,
        HasNextPage: usersResponse.HasNextPage,
        HasPrevPage: usersResponse.HasPrevPage,
    }
    
    c.JSON(http.StatusOK, response)
}
```

### 🟢 Rutas HTTP: `internal/handler/http/user_routes.go`

```go
package http

import (
    "github.com/gin-gonic/gin"
    "github.com/usuario/proyecto/internal/usecase"
)

// RegisterUserRoutes registra todas las rutas relacionadas con usuarios
func RegisterUserRoutes(router *gin.RouterGroup, userUseCase usecase.UserUseCase) {
    userHandler := NewUserHandler(userUseCase)
    
    users := router.Group("/users")
    {
        users.POST("", userHandler.Create)       // POST /api/v1/users
        users.GET("", userHandler.List)          // GET /api/v1/users
        users.GET("/:id", userHandler.GetByID)   // GET /api/v1/users/:id
        users.PUT("/:id", userHandler.Update)    // PUT /api/v1/users/:id
        users.DELETE("/:id", userHandler.Delete) // DELETE /api/v1/users/:id
    }
}
```

## 🔍 Tipos de Campos Soportados

Goca soporta múltiples tipos de campos para las entidades:

### Tipos Básicos
```bash
--fields "name:string,age:int,salary:float64,active:bool"
```

### Tipos de Tiempo
```bash
--fields "created_at:time.Time,birth_date:time.Time"
```

### Tipos Opcionales (Punteros)
```bash
--fields "nickname:*string,last_login:*time.Time"
```

### Tipos Personalizados
```bash
--fields "status:UserStatus,role:Role"
```

### Arrays y Slices
```bash
--fields "tags:[]string,scores:[]int"
```

## 🎛️ Configuraciones Avanzadas

### Multi-Handler con Configuraciones Específicas
```bash
goca feature Order \
  --fields "user_id:int,items:[]string,total:float64,status:string" \
  --handlers "http,grpc,worker" \
  --validation \
  --business-rules \
  --database postgres
```

Esto genera:
- **HTTP Handler**: API REST completa
- **gRPC Handler**: Servicio gRPC con `.proto`
- **Worker Handler**: Para procesamiento en background

### Validaciones Personalizadas
Con `--validation`, se generan validaciones:
- **En Dominio**: Reglas de negocio esenciales
- **En DTOs**: Validaciones de formato e integridad
- **En Handlers**: Validaciones de entrada específicas por protocolo

### Reglas de Negocio
Con `--business-rules`, se generan métodos adicionales:
```go
// Ejemplos de métodos generados automáticamente
func (u *User) CanUpdateProfile() bool
func (u *User) IsEmailDomainAllowed() bool
func (p *Product) CanDiscount(percentage float64) bool
func (o *Order) CanBeCanceled() bool
```

## 🔄 Flujo de Trabajo Completo

### 1. Generar Feature
```bash
goca feature Product --fields "name:string,price:float64,category:string" --validation
```

### 2. Configurar Base de Datos
```sql
-- SQL generado automáticamente comentado en el código
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    category VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 3. Configurar Inyección de Dependencias
```bash
goca di --features "Product" --database postgres
```

### 4. Integrar en Main
```go
// El código de integración se genera automáticamente
// en los comentarios de los archivos
```

## ⚡ Optimizaciones y Mejores Prácticas

### ✅ Recomendaciones
- **Campos descriptivos**: Usa nombres claros y específicos
- **Validaciones consistentes**: Siempre usar `--validation` en producción
- **Reglas de negocio**: Activar `--business-rules` para dominios complejos
- **Multi-handler inteligente**: Solo generar handlers que realmente necesites

### 🚀 Performance
- **DTOs optimizados**: Campos opcionales con punteros para actualizaciones
- **Paginación automática**: Incluida en todos los endpoints de listado
- **Queries eficientes**: Repositorios optimizados por base de datos
- **Validación temprana**: En múltiples capas para fallar rápido

### 🔒 Seguridad
- **Validación de entrada**: En todos los puntos de entrada
- **SQL Injection**: Prevención automática con prepared statements
- **Type Safety**: Tipado fuerte en toda la aplicación
- **Error Handling**: Manejo consistente de errores sin exposición de detalles

## 📊 Casos de Uso Reales

### E-commerce
```bash
# Entidades principales
goca feature User --fields "name:string,email:string,password:string" --validation
goca feature Product --fields "name:string,price:float64,stock:int" --validation
goca feature Order --fields "user_id:int,total:float64,status:string" --validation

# Configurar todo junto
goca di --features "User,Product,Order" --database postgres
```

### Sistema de Blog
```bash
goca feature Author --fields "name:string,email:string,bio:string" --validation
goca feature Post --fields "title:string,content:string,author_id:int" --validation
goca feature Comment --fields "content:string,post_id:int,author_id:int" --validation
```

### API de Microservicio
```bash
goca feature Customer --fields "name:string,email:string" --handlers "grpc" --validation
```

---

**← [Comando goca init](Command-Init) | [Comando goca entity](Command-Entity) →**
