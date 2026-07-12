# Goca CLI - Guía Completa de Comandos

Esta guía contiene la documentación detallada de todos los comandos disponibles en Goca CLI, incluyendo sus flags, ejemplos de uso y código generado.

## 📋 Índice de Comandos

- [Goca CLI - Guía Completa de Comandos](#goca-cli---guía-completa-de-comandos)
  - [📋 Índice de Comandos](#-índice-de-comandos)
  - [Comandos Principales](#comandos-principales)
    - [`goca help`](#goca-help)
    - [`goca version`](#goca-version)
    - [`goca init`](#goca-init)
    - [`goca feature`](#goca-feature)
  - [Comandos por Capas](#comandos-por-capas)
    - [`goca entity`](#goca-entity)
    - [`goca usecase`](#goca-usecase)
    - [`goca handler`](#goca-handler)
    - [`goca repository`](#goca-repository)
  - [Comandos Auxiliares](#comandos-auxiliares)
    - [`goca messages`](#goca-messages)
    - [`goca di`](#goca-di)
    - [`goca interfaces`](#goca-interfaces)
  - [📝 Notas Adicionales](#-notas-adicionales)
    - [Convenciones de Nomenclatura](#convenciones-de-nomenclatura)
    - [Validaciones Automáticas](#validaciones-automáticas)
    - [Integración con Editores](#integración-con-editores)

---

## Comandos Principales

### `goca help`

**Título**: Muestra ayuda del CLI de Goca

**Descripción**: Proporciona información detallada sobre todos los comandos disponibles en Goca CLI, incluyendo ejemplos de uso y descripción de flags.

**Uso**:
```bash
goca help [command]
```

**Flags**:
Ninguno

**Ejemplos**:
```bash
# Mostrar ayuda general
goca help

# Mostrar ayuda específica de un comando
goca help entity
goca help usecase
goca help handler
```

**Salida esperada**:
```
Goca - Go Clean Architecture Code Generator

USAGE:
  goca [command]

AVAILABLE COMMANDS:
  help        Ayuda sobre cualquier comando
  version     Muestra la versión de Goca
  init        Inicializa un nuevo proyecto con Clean Architecture
  feature     Genera un feature completo con todas las capas
  entity      Genera entidades de dominio puras
  usecase     Genera casos de uso con DTOs
  handler     Genera handlers para diferentes protocolos
  repository  Genera repositorios con interfaces
  messages    Genera mensajes y constantes
  di          Genera contenedor de inyección de dependencias
  interfaces  Genera solo interfaces para TDD

Use "goca [command] --help" para más información sobre un comando.
```

---

### `goca version`

**Título**: Versión de Goca CLI

**Descripción**: Muestra la versión actual de Goca CLI junto con información de compilación.

**Uso**:
```bash
goca version
```

**Flags**:
- `--short` (opcional): Muestra solo el número de versión

**Ejemplos**:
```bash
# Versión completa
goca version

# Versión corta
goca version --short
```

**Salida esperada**:
```
Goca v1.0.0
Build: 2025-01-19T10:30:00Z
Go Version: go1.24.5
```

---

### `goca init`

**Título**: Inicializar proyecto Clean Architecture

**Descripción**: Crea la estructura base de un proyecto Go siguiendo los principios de Clean Architecture, incluyendo directorios, archivos de configuración y estructura de capas.

**Uso**:
```bash
goca init <project-name> [flags]
```

**Flags**:
- `--module string` (requerido): Nombre del módulo Go
- `--database string` (opcional): Tipo de base de datos (postgres, postgres-json, mysql, sqlite, mongodb, sqlserver, dynamodb, elasticsearch) (default: "sqlite")
- `--auth` (opcional): Incluir boilerplate de autenticación
- `--api string` (opcional): Tipo de API (rest, grpc, graphql) (default: "rest")

**Ejemplos**:
```bash
# Proyecto básico
goca init ecommerce --module github.com/mycompany/ecommerce

# Proyecto con autenticación y MongoDB
goca init blog --module github.com/myblog/api --database mongodb --auth

# Proyecto con gRPC
goca init microservice --module github.com/company/ms --api grpc
```

**Estructura generada**:
```
myproject/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   ├── usecase/
│   ├── repository/
│   └── handler/
├── pkg/
│   ├── config/
│   └── logger/
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

---

### `goca feature`

**Título**: Generar feature completo con Clean Architecture

**Descripción**: Genera todas las capas necesarias para un feature completo, incluyendo dominio, casos de uso, repositorio y handlers en una sola operación.

**Uso**:
```bash
goca feature <name> [flags]
```

**Flags**:
- `--fields string` (requerido): Campos de la entidad "field:type,field2:type"
- `--database string` (opcional): Tipo de base de datos (postgres, mysql, mongodb) (default: "postgres")
- `--handlers string` (opcional): Tipos de handlers "http,grpc,cli" (default: "http")
- `--validation` (opcional): Incluir validaciones en todas las capas
- `--business-rules` (opcional): Incluir métodos de reglas de negocio

**Ejemplos**:
```bash
# Feature básico
goca feature Product --fields "name:string,price:float64,category:string"

# Feature completo con validaciones
goca feature Employee --fields "name:string,email:string,role:string" --validation --business-rules

# Feature con múltiples handlers
goca feature Order --fields "total:float64,status:string" --handlers "http,grpc,cli"
```

**Archivos generados**:
```
product/
├── domain/
│   ├── product.go
│   ├── errors.go
│   └── validations.go
├── usecase/
│   ├── dto.go
│   ├── product_usecase.go
│   └── product_service.go
├── repository/
│   ├── interfaces.go
│   └── postgres_product_repo.go
├── handler/
│   └── http/
│       ├── product_handler.go
│       └── routes.go
└── messages/
    ├── errors.go
    └── responses.go
```

---

## Comandos por Capas

### `goca entity`

**Título**: Generar entidad de dominio pura

**Descripción**: Crea entidades de dominio siguiendo los principios DDD, sin dependencias externas y con validaciones de negocio.

**Uso**:
```bash
goca entity <name> [flags]
```

**Flags**:
- `--fields string` (requerido): Campos de la entidad "name:type,email:string"
- `--validation` (opcional): Agregar validaciones de dominio
- `--business-rules` (opcional): Incluir métodos de reglas de negocio
- `--timestamps` (opcional): Agregar campos created_at y updated_at
- `--soft-delete` (opcional): Agregar funcionalidad de soft delete

**Ejemplos**:
```bash
# Entidad básica
goca entity User --fields "name:string,email:string,age:int"

# Entidad con validaciones y reglas de negocio
goca entity Product --fields "name:string,price:float64" --validation --business-rules

# Entidad con timestamps
goca entity Article --fields "title:string,content:string" --timestamps
```

**Código generado**:
```go
// domain/user.go
package domain

import "errors"

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func (u *User) Validate() error {
    if u.Name == "" {
        return ErrInvalidUserName
    }
    if u.Email == "" {
        return ErrInvalidUserEmail
    }
    if u.Age < 0 {
        return ErrInvalidUserAge
    }
    return nil
}

func (u *User) IsAdult() bool {
    return u.Age >= 18
}

// domain/errors.go
package domain

import "errors"

var (
    ErrInvalidUserData  = errors.New("invalid user data")
    ErrInvalidUserName  = errors.New("user name is required")
    ErrInvalidUserEmail = errors.New("user email is required")
    ErrInvalidUserAge   = errors.New("user age must be positive")
)
```

---

### `goca usecase`

**Título**: Generar casos de uso con DTOs

**Descripción**: Crea servicios de aplicación con DTOs bien definidos, interfaces claras y lógica de negocio encapsulada.

**Uso**:
```bash
goca usecase <name> [flags]
```

**Flags**:
- `--entity string` (requerido): Entidad asociada al caso de uso
- `--operations string` (opcional): Operaciones CRUD "create,read,update,delete,list" (default: "create,read")
- `--dto-validation` (opcional): DTOs con validaciones específicas
- `--async` (opcional): Incluir operaciones asíncronas

**Ejemplos**:
```bash
# Caso de uso básico
goca usecase UserService --entity User

# Caso de uso CRUD completo
goca usecase ProductService --entity Product --operations "create,read,update,delete,list"

# Caso de uso con validaciones
goca usecase OrderService --entity Order --dto-validation
```

**Código generado**:
```go
// usecase/dto.go
package usecase

type CreateUserInput struct {
    Name  string `json:"name" validate:"required,min=2"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"required,min=0"`
}

type CreateUserOutput struct {
    User    domain.User `json:"user"`
    Message string      `json:"message"`
}

type UpdateUserInput struct {
    Name  string `json:"name,omitempty" validate:"omitempty,min=2"`
    Email string `json:"email,omitempty" validate:"omitempty,email"`
    Age   int    `json:"age,omitempty" validate:"omitempty,min=0"`
}

// usecase/interfaces.go
package usecase

import "myproject/domain"

type UserRepository interface {
    Save(user *domain.User) error
    FindByID(id int) (*domain.User, error)
    FindByEmail(email string) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id int) error
    FindAll() ([]domain.User, error)
}

// usecase/user_usecase.go
package usecase

import "myproject/domain"

type UserUseCase interface {
    CreateUser(input CreateUserInput) (CreateUserOutput, error)
    GetUser(id int) (*domain.User, error)
    UpdateUser(id int, input UpdateUserInput) error
    DeleteUser(id int) error
    ListUsers() ([]domain.User, error)
}

// usecase/user_service.go
package usecase

import (
    "myproject/domain"
    "myproject/messages"
)

type userService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) UserUseCase {
    return &userService{repo: repo}
}

func (s *userService) CreateUser(input CreateUserInput) (CreateUserOutput, error) {
    user := domain.User{
        Name:  input.Name,
        Email: input.Email,
        Age:   input.Age,
    }
    
    if err := user.Validate(); err != nil {
        return CreateUserOutput{}, err
    }
    
    if err := s.repo.Save(&user); err != nil {
        return CreateUserOutput{}, err
    }
    
    return CreateUserOutput{
        User:    user,
        Message: messages.UserCreatedSuccessfully,
    }, nil
}

func (s *userService) GetUser(id int) (*domain.User, error) {
    return s.repo.FindByID(id)
}
```

---

### `goca handler`

**Título**: Generar handlers para diferentes protocolos

**Descripción**: Crea adaptadores de entrega que manejan diferentes protocolos (HTTP, gRPC, CLI) manteniendo la separación de capas.

**Uso**:
```bash
goca handler <entity> [flags]
```

**Flags**:
- `--type string` (requerido): Tipo de handler (http, grpc, cli, worker, soap)
- `--middleware` (opcional): Incluir setup de middleware
- `--validation` (opcional): Validación de entrada en handler
- `--swagger` (opcional): Generar documentación Swagger (solo HTTP)

**Ejemplos**:
```bash
# Handler HTTP básico
goca handler User --type http

# Handler HTTP con middleware
goca handler Product --type http --middleware --validation

# Handler gRPC
goca handler Order --type grpc

# Handler CLI
goca handler Employee --type cli
```

**Código generado para HTTP**:
```go
// handler/http/user_handler.go
package http

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "github.com/gorilla/mux"
    "myproject/usecase"
)

type UserHandler struct {
    usecase usecase.UserUseCase
}

func NewUserHandler(uc usecase.UserUseCase) *UserHandler {
    return &UserHandler{usecase: uc}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var input usecase.CreateUserInput
    
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    output, err := h.usecase.CreateUser(input)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(output)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    user, err := h.usecase.GetUser(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// handler/http/routes.go
package http

import (
    "github.com/gorilla/mux"
    "myproject/usecase"
)

func SetupUserRoutes(router *mux.Router, uc usecase.UserUseCase) {
    handler := NewUserHandler(uc)
    
    router.HandleFunc("/users", handler.CreateUser).Methods("POST")
    router.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
    router.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
    router.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")
    router.HandleFunc("/users", handler.ListUsers).Methods("GET")
}
```

---

### `goca repository`

**Título**: Generar repositorios con interfaces

**Descripción**: Crea repositorios que implementan el patrón Repository con interfaces bien definidas e implementaciones específicas por base de datos.

**Uso**:
```bash
goca repository <entity> [flags]
```

**Flags**:
- `--database string` (requerido): Tipo de base de datos (postgres, mysql, mongodb)
- `--interface-only` (opcional): Solo generar interfaces
- `--implementation` (opcional): Solo generar implementación
- `--cache` (opcional): Incluir capa de caché
- `--transactions` (opcional): Incluir soporte para transacciones

**Ejemplos**:
```bash
# Repositorio básico
goca repository User --database postgres

# Solo interfaces (útil para TDD)
goca repository Product --interface-only

# Con caché y transacciones
goca repository Order --database postgres --cache --transactions
```

**Código generado**:
```go
// repository/interfaces.go
package repository

import "myproject/domain"

type UserRepository interface {
    Save(user *domain.User) error
    FindByID(id int) (*domain.User, error)
    FindByEmail(email string) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id int) error
    FindAll() ([]domain.User, error)
}

// repository/postgres_user_repository.go
package repository

import (
    "database/sql"
    "myproject/domain"
    
    _ "github.com/lib/pq"
)

type postgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
    return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Save(user *domain.User) error {
    query := `INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id`
    err := r.db.QueryRow(query, user.Name, user.Email, user.Age).Scan(&user.ID)
    return err
}

func (r *postgresUserRepository) FindByID(id int) (*domain.User, error) {
    user := &domain.User{}
    query := `SELECT id, name, email, age FROM users WHERE id = $1`
    err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Age)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *postgresUserRepository) FindByEmail(email string) (*domain.User, error) {
    user := &domain.User{}
    query := `SELECT id, name, email, age FROM users WHERE email = $1`
    err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Age)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

---

## Comandos Auxiliares

### `goca messages`

**Título**: Generar mensajes y constantes

**Descripción**: Crea archivos de mensajes de error, respuestas y constantes organizados por feature para mantener consistencia en la aplicación.

**Uso**:
```bash
goca messages <entity> [flags]
```

**Flags**:
- `--errors` (opcional): Generar mensajes de error
- `--responses` (opcional): Generar mensajes de respuesta
- `--constants` (opcional): Generar constantes del feature
- `--all` (opcional): Generar todos los tipos de mensajes

**Ejemplos**:
```bash
# Solo errores
goca messages User --errors

# Solo respuestas
goca messages Product --responses

# Todo
goca messages Order --all
```

**Código generado**:
```go
// messages/errors.go
package messages

const (
    // User errors
    ErrUserNotFound        = "user not found"
    ErrUserAlreadyExists   = "user already exists"
    ErrInvalidUserData     = "invalid user data"
    ErrUserEmailRequired   = "user email is required"
    ErrUserNameRequired    = "user name is required"
    ErrUserAgeInvalid      = "user age must be positive"
)

// messages/responses.go
package messages

const (
    // User success messages
    UserCreatedSuccessfully = "user created successfully"
    UserUpdatedSuccessfully = "user updated successfully"
    UserDeletedSuccessfully = "user deleted successfully"
    UserFoundSuccessfully   = "user found successfully"
    UsersListedSuccessfully = "users listed successfully"
)

// constants/constants.go
package constants

const (
    // User constants
    MinUserAge        = 0
    MaxUserAge        = 150
    MinUserNameLength = 2
    MaxUserNameLength = 100
    UserTableName     = "users"
)
```

---

### `goca di`

**Título**: Generar contenedor de inyección de dependencias

**Descripción**: Crea un contenedor de inyección de dependencias que conecta automáticamente todas las capas del sistema.

**Uso**:
```bash
goca di [flags]
```

**Flags**:
- `--features string` (requerido): Features a incluir "User,Product,Order"
- `--database string` (opcional): Tipo de base de datos (default: "postgres")
- `--wire` (opcional): Usar Google Wire para DI

**Ejemplos**:
```bash
# DI básico
goca di --features "User,Product"

# DI con Wire
goca di --features "User,Product,Order" --wire
```

**Código generado**:
```go
// infrastructure/di/container.go
package di

import (
    "database/sql"
    
    "myproject/repository"
    "myproject/usecase"
    "myproject/handler/http"
)

type Container struct {
    db *sql.DB
    
    // Repositories
    userRepo    repository.UserRepository
    productRepo repository.ProductRepository
    
    // Use Cases
    userUC    usecase.UserUseCase
    productUC usecase.ProductUseCase
    
    // Handlers
    userHandler    *http.UserHandler
    productHandler *http.ProductHandler
}

func NewContainer(db *sql.DB) *Container {
    c := &Container{db: db}
    c.setupRepositories()
    c.setupUseCases()
    c.setupHandlers()
    return c
}

func (c *Container) setupRepositories() {
    c.userRepo = repository.NewPostgresUserRepository(c.db)
    c.productRepo = repository.NewPostgresProductRepository(c.db)
}

func (c *Container) setupUseCases() {
    c.userUC = usecase.NewUserService(c.userRepo)
    c.productUC = usecase.NewProductService(c.productRepo)
}

func (c *Container) setupHandlers() {
    c.userHandler = http.NewUserHandler(c.userUC)
    c.productHandler = http.NewProductHandler(c.productUC)
}

// Getters
func (c *Container) UserHandler() *http.UserHandler {
    return c.userHandler
}

func (c *Container) ProductHandler() *http.ProductHandler {
    return c.productHandler
}
```

---

### `goca interfaces`

**Título**: Generar solo interfaces para TDD

**Descripción**: Genera únicamente las interfaces de contratos entre capas, útil para desarrollo dirigido por pruebas (TDD).

**Uso**:
```bash
goca interfaces <entity> [flags]
```

**Flags**:
- `--usecase` (opcional): Generar interfaces de casos de uso
- `--repository` (opcional): Generar interfaces de repositorio
- `--handler` (opcional): Generar interfaces de handlers
- `--all` (opcional): Generar todas las interfaces

**Ejemplos**:
```bash
# Solo interfaces de repositorio
goca interfaces User --repository

# Todas las interfaces
goca interfaces Product --all
```

**Código generado**:
```go
// interfaces/user_usecase.go
package interfaces

import "myproject/domain"

type UserUseCase interface {
    CreateUser(input CreateUserInput) (CreateUserOutput, error)
    GetUser(id int) (*domain.User, error)
    UpdateUser(id int, input UpdateUserInput) error
    DeleteUser(id int) error
    ListUsers() ([]domain.User, error)
}

// interfaces/user_repository.go
package interfaces

import "myproject/domain"

type UserRepository interface {
    Save(user *domain.User) error
    FindByID(id int) (*domain.User, error)
    FindByEmail(email string) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id int) error
    FindAll() ([]domain.User, error)
}
```

---

## 📝 Notas Adicionales

### Convenciones de Nomenclatura

- **Entidades**: PascalCase singular (User, Product, Order)
- **Archivos**: snake_case con sufijo descriptivo (user_repository.go, product_service.go)
- **Paquetes**: lowercase singular (domain, usecase, repository)

### Validaciones Automáticas

Todos los comandos incluyen validaciones para:
- Nombres de entidades válidos
- Tipos de campos soportados
- Dependencias entre capas
- Consistencia de nomenclatura

### Integración con Editores

Los comentarios `// filepath:` permiten integración directa con editores que soporten este formato para navegación automática de archivos.

---

**Nota**: Esta guía está basada en los principios de Clean Architecture definidos en el archivo `rules.md` del proyecto.