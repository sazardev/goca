# Getting Started with Goca

This guide will help you create your first project with Goca in **less than 5 minutes**. By the end, you'll have a functional API with Clean Architecture.

## 🎯 What We'll Build

In this guide we'll create:
- ✅ A basic project with Clean Architecture structure
- ✅ A complete `User` entity
- ✅ Functional REST API with CRUD
- ✅ PostgreSQL database configured

## ⏱️ Estimated Time: 5 minutes

## 📋 Prerequisites

- ✅ **Go 1.21+** - [Download here](https://golang.org/dl/)
- ✅ **Goca installed** - [See installation guide](Installation)
- ✅ **PostgreSQL** (optional for this tutorial)

## 🚀 Step 1: Create the Project (30 seconds)

```bash
# Create and enter directory
mkdir my-first-project
cd my-first-project

# Initialize with Goca
goca init my-api --module github.com/user/my-api --database postgres

# Enter generated directory
cd my-api
```

**✅ Result:** Complete project structure created

## 👤 Step 2: Create User Feature (30 seconds)

```bash
# Generate complete user feature
goca feature User --fields "name:string,email:string,age:int" --validation

# See what was generated
find internal/ -name "*user*" -type f
```

**✅ Result:** 8+ files generated with all Clean Architecture layers

## 🔌 Step 3: Configure Dependencies (30 seconds)

```bash
# Generate dependency injection
goca di --features "User" --database postgres

# Install Go dependencies
go mod tidy
```

**✅ Result:** DI container configured and dependencies installed

## 🏃‍♂️ Step 4: Run the Project (30 seconds)

```bash
# Run the server
go run cmd/server/main.go
```

**✅ Result:** Server running at http://localhost:8080

## 🧪 Step 5: Test the API (3 minutes)

### Health Check
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{"status": "ok"}
```

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ana García",
    "email": "ana@example.com",
    "age": 28
  }'
```

**Response:**
```json
{
  "id": 1,
  "name": "Ana García", 
  "email": "ana@example.com",
  "age": 28,
  "created_at": "2025-07-20T10:30:00Z",
  "updated_at": "2025-07-20T10:30:00Z"
}
```

### Get User
```bash
curl http://localhost:8080/api/v1/users/1
```

### List Users
```bash
curl http://localhost:8080/api/v1/users
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ana García Updated",
    "email": "ana.updated@example.com",
    "age": 29
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

### Listar Usuarios
```bash
curl http://localhost:8080/api/v1/users
```

### Actualizar Usuario
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ana García López",
    "age": 29
  }'
```

### Eliminar Usuario
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## 🎉 ¡Felicitaciones!

En solo 5 minutos has creado:

### ✅ Estructura Completa
```
mi-api/
├── cmd/server/main.go               # Servidor HTTP
├── internal/
│   ├── domain/user.go               # Entidad de dominio
│   ├── usecase/user_usecase.go      # Lógica de aplicación
│   ├── repository/postgres/         # Acceso a datos
│   ├── handler/http/                # API REST
│   └── infrastructure/di/           # Inyección de dependencias
└── pkg/
    ├── config/                      # Configuración
    └── logger/                      # Logging
```

### ✅ API REST Funcional
- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - List users  
- `GET /api/v1/users/:id` - Get user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### ✅ Clean Architecture
- **🟡 Dominio**: Entidad `User` con validaciones
- **🔴 Casos de Uso**: CRUD completo con DTOs
- **🔵 Repositorio**: Interface y implementación PostgreSQL
- **🟢 Handler**: API REST con manejo de errores

## 🔍 ¿Qué Acabas de Crear?

### Entidad de Dominio
```go
// internal/domain/user.go
type User struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Age       int       `json:"age"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
    if strings.TrimSpace(u.Name) == "" {
        return ErrUserNameRequired
    }
    if u.Age < 0 || u.Age > 150 {
        return ErrUserAgeInvalid
    }
    return nil
}
```

### Caso de Uso
```go
// internal/usecase/user_usecase.go
func (uc *userUseCase) Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
    // 1. Validar DTO
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // 2. Crear entidad
    user := &domain.User{
        Name:  req.Name,
        Email: req.Email,
        Age:   req.Age,
    }
    
    // 3. Validar entidad
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    // 4. Guardar
    if err := uc.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    return &UserResponse{...}, nil
}
```

### Handler HTTP
```go
// internal/handler/http/user_handler.go
func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Error: err.Error()})
        return
    }
    
    user, err := h.userUseCase.Create(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, ErrorResponse{Error: err.Error()})
        return
    }
    
    c.JSON(201, user)
}
```

## 🚀 Próximos Pasos

### 1. **Agregar Base de Datos Real**
```bash
# Configurar PostgreSQL
echo "DB_PASSWORD=mypassword" > .env

# Crear tabla
psql -c "CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(255), email VARCHAR(255), age INTEGER, created_at TIMESTAMP, updated_at TIMESTAMP);"
```

### 2. **Agregar Más Features**
```bash
# Crear feature de productos
goca feature Product --fields "name:string,price:float64,category:string" --validation

# Actualizar DI
goca di --features "User,Product" --database postgres
```

### 3. **Agregar Autenticación**
```bash
# Recrear proyecto con auth
goca init mi-api-segura --module github.com/usuario/mi-api-segura --auth
```

### 4. **Explorar Funcionalidades Avanzadas**
- [Tutorial Completo de E-commerce](Complete-Tutorial)
- [Comando goca feature](Command-Feature)
- [Clean Architecture](Clean-Architecture)

## 🐛 Troubleshooting

### Error: "command not found"
```bash
# Verificar instalación
goca version

# Si no está instalado
go install github.com/sazardev/goca@latest
```

### Error: "module not found" 
```bash
# Limpiar cache de módulos
go clean -modcache
go mod tidy
```

### Error: "port already in use"
```bash
# Cambiar puerto
echo "PORT=8081" > .env
```

### Error de Base de Datos
```bash
# Por defecto usa in-memory, PostgreSQL es opcional
# Para usar PostgreSQL, configurar .env primero
```

## 📚 Conceptos Clave Aprendidos

### 🏗️ Clean Architecture
- **Separación de capas** por responsabilidades
- **Dependencias apuntando hacia adentro**
- **Interfaces para desacoplar** implementaciones

### 🔄 Flujo de Datos
```
HTTP Request → Handler → UseCase → Domain → Repository → Database
                ↓
HTTP Response ← DTO ← Response ← Entity ← Data ← Query Result
```

### 🎯 Principios SOLID
- **SRP**: Cada archivo tiene una responsabilidad
- **OCP**: Extensible via interfaces
- **LSP**: Implementaciones intercambiables
- **ISP**: Interfaces específicas
- **DIP**: Depende de abstracciones

## 🎊 ¡Excelente Trabajo!

Has completado exitosamente tus primeros pasos con Goca. Ahora tienes:

- ✅ **Comprensión práctica** de Clean Architecture
- ✅ **API funcional** con CRUD completo
- ✅ **Código de calidad** siguiendo mejores prácticas
- ✅ **Base sólida** para proyectos más complejos

**¿Listo para el siguiente nivel?** Prueba el [Tutorial Completo](Complete-Tutorial) donde construiremos un sistema de e-commerce completo.

---

**← [Instalación](Installation) | [Tutorial Completo](Complete-Tutorial) →**
