# Comando goca di

El comando `goca di` crea un contenedor de inyección de dependencias que conecta automáticamente todas las capas del sistema siguiendo los principios de Clean Architecture.

## 📋 Sintaxis

```bash
goca di [flags]
```

## 🎯 Propósito

Crea el sistema de inyección de dependencias para conectar todas las capas:

- 🔧 **Contenedor manual** con configuración explícita
- ⚡ **Google Wire** para generación automática
- 🔗 **Conexión de capas** respetando las dependencias
- 🗄️ **Configuración de DB** específica por tipo
- 📦 **Features modulares** para diferentes funcionalidades

## 🚩 Flags Disponibles

| Flag         | Tipo     | Requerido | Valor por Defecto | Descripción                                               |
| ------------ | -------- | --------- | ----------------- | --------------------------------------------------------- |
| `--features` | `string` | ✅ **Sí**  | -                 | Características del proyecto (`crud,auth,validation,etc`) |
| `--database` | `string` | ❌ No      | `postgres`        | Tipo de base de datos (`postgres`, `mysql`, `mongodb`)    |
| `--wire`     | `bool`   | ❌ No      | `false`           | Usar Google Wire para inyección de dependencias           |

## 📖 Ejemplos de Uso

### Contenedor Manual Básico
```bash
goca di --features "crud" --database postgres
```

### Con Google Wire
```bash
goca di --features "crud,auth,validation" --database postgres --wire
```

### Múltiples Features
```bash
goca di --features "crud,auth,validation,logging,metrics" --database mysql
```

### MongoDB con Wire
```bash
goca di --features "crud,auth" --database mongodb --wire
```

## 📂 Archivos Generados

### Contenedor Manual
```
internal/di/
└── container.go           # Contenedor manual de dependencias
```

### Con Google Wire
```
internal/di/
├── wire.go                # Definiciones Wire
├── wire_gen.go            # Código generado por Wire
└── container.go           # Wrapper del contenedor Wire
```

## 🔍 Código Generado en Detalle

### Contenedor Manual: `internal/di/container.go`

```go
package di

import (
    "database/sql"
    "log"
    
    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
    
    "github.com/usuario/proyecto/internal/usecase"
    "github.com/usuario/proyecto/internal/repository/postgres"
    "github.com/usuario/proyecto/internal/handler/http"
    "github.com/usuario/proyecto/pkg/config"
)

// Container contiene todas las dependencias del sistema
type Container struct {
    // Base de datos
    DB *sql.DB
    
    // Repositorios
    UserRepository interfaces.UserRepository
    
    // Casos de uso
    UserUseCase usecase.UserUseCase
    
    // Handlers
    UserHandler *http.UserHandler
    
    // Router
    Router *gin.Engine
}

// NewContainer crea y configura un nuevo contenedor de dependencias
func NewContainer() (*Container, error) {
    container := &Container{}
    
    // Configurar base de datos
    if err := container.setupDatabase(); err != nil {
        return nil, err
    }
    
    // Configurar repositorios
    container.setupRepositories()
    
    // Configurar casos de uso
    container.setupUseCases()
    
    // Configurar handlers
    container.setupHandlers()
    
    // Configurar router
    container.setupRouter()
    
    return container, nil
}

// setupDatabase configura la conexión a la base de datos
func (c *Container) setupDatabase() error {
    cfg := config.Load()
    
    dbUrl := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Name,
        cfg.Database.SSLMode,
    )
    
    db, err := sql.Open("postgres", dbUrl)
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }
    
    c.DB = db
    log.Println("Database connected successfully")
    
    return nil
}

// setupRepositories configura todos los repositorios
func (c *Container) setupRepositories() {
    c.UserRepository = postgres.NewUserRepository(c.DB)
    log.Println("Repositories configured")
}

// setupUseCases configura todos los casos de uso
func (c *Container) setupUseCases() {
    c.UserUseCase = usecase.NewUserService(c.UserRepository)
    log.Println("Use cases configured")
}

// setupHandlers configura todos los handlers
func (c *Container) setupHandlers() {
    c.UserHandler = http.NewUserHandler(c.UserUseCase)
    log.Println("Handlers configured")
}

// setupRouter configura el router con todas las rutas
func (c *Container) setupRouter() {
    router := gin.Default()
    
    // Middleware global
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
    router.Use(corsMiddleware())
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // API v1
    v1 := router.Group("/api/v1")
    {
        // Rutas de usuarios
        http.RegisterUserRoutes(v1, c.UserUseCase)
    }
    
    c.Router = router
    log.Println("Router configured with all routes")
}

// Close cierra todas las conexiones
func (c *Container) Close() error {
    if c.DB != nil {
        if err := c.DB.Close(); err != nil {
            return fmt.Errorf("failed to close database: %w", err)
        }
    }
    
    return nil
}

// corsMiddleware configura CORS
func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

### Google Wire: `internal/di/wire.go`

```go
//go:build wireinject
// +build wireinject

package di

import (
    "database/sql"
    
    "github.com/google/wire"
    "github.com/gin-gonic/gin"
    
    "github.com/usuario/proyecto/internal/usecase"
    "github.com/usuario/proyecto/internal/repository/postgres"
    "github.com/usuario/proyecto/internal/handler/http"
    "github.com/usuario/proyecto/pkg/config"
)

// WireContainer contiene las dependencias inyectadas por Wire
type WireContainer struct {
    DB          *sql.DB
    UserHandler *http.UserHandler
    Router      *gin.Engine
}

// DatabaseSet conjunto de proveedores para base de datos
var DatabaseSet = wire.NewSet(
    provideDatabaseConnection,
)

// RepositorySet conjunto de proveedores para repositorios
var RepositorySet = wire.NewSet(
    postgres.NewUserRepository,
    wire.Bind(new(interfaces.UserRepository), new(*postgres.UserRepository)),
)

// UseCaseSet conjunto de proveedores para casos de uso
var UseCaseSet = wire.NewSet(
    usecase.NewUserService,
    wire.Bind(new(usecase.UserUseCase), new(*usecase.UserService)),
)

// HandlerSet conjunto de proveedores para handlers
var HandlerSet = wire.NewSet(
    http.NewUserHandler,
)

// RouterSet conjunto de proveedores para router
var RouterSet = wire.NewSet(
    provideRouter,
)

// InitializeContainer inicializa el contenedor con Wire
func InitializeContainer() (*WireContainer, error) {
    wire.Build(
        DatabaseSet,
        RepositorySet,
        UseCaseSet,
        HandlerSet,
        RouterSet,
        wire.Struct(new(WireContainer), "*"),
    )
    return &WireContainer{}, nil
}

// provideDatabaseConnection provee la conexión a la base de datos
func provideDatabaseConnection() (*sql.DB, error) {
    cfg := config.Load()
    
    dbUrl := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Name,
        cfg.Database.SSLMode,
    )
    
    db, err := sql.Open("postgres", dbUrl)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return db, nil
}

// provideRouter provee el router configurado
func provideRouter(userHandler *http.UserHandler) *gin.Engine {
    router := gin.Default()
    
    // Middleware global
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // API v1
    v1 := router.Group("/api/v1")
    {
        // Rutas de usuarios
        users := v1.Group("/users")
        {
            users.POST("", userHandler.Create)
            users.GET("", userHandler.List)
            users.GET("/:id", userHandler.GetByID)
            users.PUT("/:id", userHandler.Update)
            users.DELETE("/:id", userHandler.Delete)
        }
    }
    
    return router
}
```

### Wrapper Wire: `internal/di/wire_container.go`

```go
package di

import (
    "database/sql"
    "log"
    
    "github.com/gin-gonic/gin"
)

// Container wrapper para el contenedor Wire
type Container struct {
    wireContainer *WireContainer
}

// NewContainer crea un nuevo contenedor usando Wire
func NewContainer() (*Container, error) {
    wireContainer, err := InitializeContainer()
    if err != nil {
        return nil, err
    }
    
    log.Println("Wire container initialized successfully")
    
    return &Container{
        wireContainer: wireContainer,
    }, nil
}

// GetDB retorna la conexión a la base de datos
func (c *Container) GetDB() *sql.DB {
    return c.wireContainer.DB
}

// GetRouter retorna el router configurado
func (c *Container) GetRouter() *gin.Engine {
    return c.wireContainer.Router
}

// Close cierra todas las conexiones
func (c *Container) Close() error {
    if c.wireContainer.DB != nil {
        if err := c.wireContainer.DB.Close(); err != nil {
            return err
        }
    }
    
    return nil
}
```

## 💾 Configuración por Base de Datos

### PostgreSQL
```go
func setupPostgresDatabase() (*sql.DB, error) {
    cfg := config.Load()
    
    dbUrl := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Name,
        cfg.Database.SSLMode,
    )
    
    db, err := sql.Open("postgres", dbUrl)
    if err != nil {
        return nil, err
    }
    
    // Configuración específica de PostgreSQL
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    return db, nil
}
```

### MySQL
```go
func setupMySQLDatabase() (*sql.DB, error) {
    cfg := config.Load()
    
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.Name,
    )
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    
    // Configuración específica de MySQL
    db.SetMaxOpenConns(20)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(3 * time.Minute)
    
    return db, nil
}
```

### MongoDB
```go
func setupMongoDatabase() (*mongo.Database, error) {
    cfg := config.Load()
    
    clientOptions := options.Client().ApplyURI(cfg.Database.URI)
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        return nil, err
    }
    
    // Verificar conexión
    if err := client.Ping(context.Background(), nil); err != nil {
        return nil, err
    }
    
    db := client.Database(cfg.Database.Name)
    return db, nil
}
```

## 🎯 Features Configurables

### CRUD Feature
```go
// Con --features "crud"
func setupCRUDFeatures(container *Container) {
    // Repositorios CRUD básicos
    container.UserRepository = postgres.NewUserRepository(container.DB)
    
    // Casos de uso CRUD
    container.UserUseCase = usecase.NewUserService(container.UserRepository)
    
    // Handlers CRUD
    container.UserHandler = http.NewUserHandler(container.UserUseCase)
}
```

### Auth Feature
```go
// Con --features "crud,auth"
func setupAuthFeatures(container *Container) {
    // JWT Service
    container.JWTService = auth.NewJWTService(
        os.Getenv("JWT_SECRET"),
        "goca-app",
    )
    
    // Auth Middleware
    container.AuthMiddleware = auth.NewAuthMiddleware(container.JWTService)
    
    // Auth Use Case
    container.AuthUseCase = auth.NewAuthService(
        container.UserRepository,
        container.JWTService,
    )
    
    // Auth Handler
    container.AuthHandler = http.NewAuthHandler(container.AuthUseCase)
}
```

### Validation Feature
```go
// Con --features "crud,validation"
func setupValidationFeatures(container *Container) {
    // Validator
    container.Validator = validator.New()
    
    // Validation Middleware
    container.ValidationMiddleware = middleware.NewValidationMiddleware(container.Validator)
    
    // Enhanced Use Cases con validación
    container.UserUseCase = usecase.NewUserServiceWithValidation(
        container.UserRepository,
        container.Validator,
    )
}
```

### Logging Feature
```go
// Con --features "crud,logging"
func setupLoggingFeatures(container *Container) {
    // Logger
    container.Logger = logger.New()
    
    // Logging Middleware
    container.LoggingMiddleware = middleware.NewLoggingMiddleware(container.Logger)
    
    // Enhanced Use Cases con logging
    container.UserUseCase = usecase.NewUserServiceWithLogging(
        container.UserRepository,
        container.Logger,
    )
}
```

## 🔧 Uso del Contenedor

### En main.go
```go
package main

import (
    "log"
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/usuario/proyecto/internal/di"
)

func main() {
    // Inicializar contenedor
    container, err := di.NewContainer()
    if err != nil {
        log.Fatal("Failed to initialize container:", err)
    }
    defer container.Close()
    
    // Configurar servidor
    server := &http.Server{
        Addr:    ":8080",
        Handler: container.GetRouter(),
    }
    
    // Iniciar servidor en goroutine
    go func() {
        log.Println("Server starting on :8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Failed to start server:", err)
        }
    }()
    
    // Esperar señal de terminación
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    
    log.Println("Server exited")
}
```

### En Tests
```go
func setupTestContainer(t *testing.T) *di.Container {
    // Configurar base de datos de prueba
    testDB := setupTestDB(t)
    
    container := &di.Container{
        DB: testDB,
    }
    
    // Configurar dependencias para pruebas
    container.UserRepository = postgres.NewUserRepository(testDB)
    container.UserUseCase = usecase.NewUserService(container.UserRepository)
    container.UserHandler = http.NewUserHandler(container.UserUseCase)
    
    return container
}

func TestUserCreation(t *testing.T) {
    container := setupTestContainer(t)
    defer container.Close()
    
    // Usar el contenedor en las pruebas
    // ...
}
```

## 🏗️ Comandos Wire

### Generar código Wire
```bash
# Instalar Wire
go install github.com/google/wire/cmd/wire@latest

# Generar código
cd internal/di
wire
```

### Verificar dependencias
```bash
# Verificar que Wire puede resolver todas las dependencias
wire check
```

## ⚠️ Consideraciones Importantes

### ✅ Buenas Prácticas
- **Interfaces claras**: Definir interfaces para todas las dependencias
- **Singleton pattern**: Usar una sola instancia del contenedor
- **Graceful shutdown**: Cerrar recursos apropiadamente
- **Configuration centralized**: Centralizar toda la configuración

### ❌ Errores Comunes
- **Circular dependencies**: Evitar dependencias circulares
- **Memory leaks**: No cerrar recursos adecuadamente
- **Hard-coded values**: Usar configuración para todos los valores
- **Missing error handling**: Manejar errores de inicialización

### 🔄 Dependencias Recomendadas

Para usar Wire, agregar a `go.mod`:
```go
require (
    github.com/google/wire v0.5.0
)
```

Para contenedor manual:
```go
require (
    github.com/gin-gonic/gin v1.9.1
    github.com/lib/pq v1.10.9          // PostgreSQL
    github.com/go-sql-driver/mysql v1.7.1  // MySQL
    go.mongodb.org/mongo-driver v1.12.1     // MongoDB
)
```

---

**← [Comando goca handler](Command-Handler) | [Comando goca interfaces](Command-Interfaces) →**
