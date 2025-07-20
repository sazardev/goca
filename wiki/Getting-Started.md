# Primeros Pasos con Goca

Esta guía te ayudará a crear tu primer proyecto con Goca en **menos de 5 minutos**. Al final tendrás una API funcional con Clean Architecture.

## 🎯 Lo que Construiremos

En esta guía crearemos:
- ✅ Un proyecto básico con estructura Clean Architecture
- ✅ Una entidad `User` completa
- ✅ API REST funcional con CRUD
- ✅ Base de datos PostgreSQL configurada

## ⏱️ Tiempo Estimado: 5 minutos

## 📋 Prerrequisitos

- ✅ **Go 1.21+** - [Descargar aquí](https://golang.org/dl/)
- ✅ **Goca instalado** - [Ver guía de instalación](Installation)
- ✅ **PostgreSQL** (opcional para este tutorial)

## 🚀 Paso 1: Crear el Proyecto (30 segundos)

```bash
# Crear y entrar en directorio
mkdir mi-primer-proyecto
cd mi-primer-proyecto

# Inicializar con Goca
goca init mi-api --module github.com/usuario/mi-api --database postgres

# Entrar al directorio generado
cd mi-api
```

**✅ Resultado:** Estructura de proyecto completa creada

## 👤 Paso 2: Crear Feature de Usuario (30 segundos)

```bash
# Generar feature completo de usuario
goca feature User --fields "name:string,email:string,age:int" --validation

# Ver lo que se generó
find internal/ -name "*user*" -type f
```

**✅ Resultado:** 8+ archivos generados con todas las capas de Clean Architecture

## 🔌 Paso 3: Configurar Dependencias (30 segundos)

```bash
# Generar inyección de dependencias
goca di --features "User" --database postgres

# Instalar dependencias Go
go mod tidy
```

**✅ Resultado:** Contenedor DI configurado y dependencias instaladas

## 🏃‍♂️ Paso 4: Ejecutar el Proyecto (30 segundos)

```bash
# Ejecutar el servidor
go run cmd/server/main.go
```

**✅ Resultado:** Servidor corriendo en http://localhost:8080

## 🧪 Paso 5: Probar la API (3 minutos)

### Health Check
```bash
curl http://localhost:8080/health
```

**Respuesta:**
```json
{"status": "ok"}
```

### Crear Usuario
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ana García",
    "email": "ana@example.com",
    "age": 28
  }'
```

**Respuesta:**
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

### Obtener Usuario
```bash
curl http://localhost:8080/api/v1/users/1
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
- `POST /api/v1/users` - Crear usuario
- `GET /api/v1/users` - Listar usuarios  
- `GET /api/v1/users/:id` - Obtener usuario
- `PUT /api/v1/users/:id` - Actualizar usuario
- `DELETE /api/v1/users/:id` - Eliminar usuario

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
