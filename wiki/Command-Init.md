# Comando goca init

El comando `goca init` es el punto de partida para crear cualquier proyecto con Clean Architecture. Genera la estructura base completa siguiendo las mejores prácticas y convenciones establecidas.

## 📋 Sintaxis

```bash
goca init <project-name> [flags]
```

## 🎯 Propósito

Crea la estructura base de un proyecto Go siguiendo los principios de Clean Architecture, incluyendo:

- 📁 Estructura de directorios organizada por capas
- 📄 Archivos de configuración esenciales
- 🔧 Setup inicial de dependencias
- 📝 Documentación base del proyecto
- 🔐 Configuración de autenticación (opcional)
- 🌐 Setup del servidor API (opcional)

## 🚩 Flags Disponibles

| Flag         | Tipo     | Requerido | Valor por Defecto | Descripción                                              |
| ------------ | -------- | --------- | ----------------- | -------------------------------------------------------- |
| `--module`   | `string` | ✅ **Sí**  | -                 | Nombre del módulo Go (ej: `github.com/usuario/proyecto`) |
| `--database` | `string` | ❌ No      | `postgres`        | Tipo de base de datos (`postgres`, `mysql`, `mongodb`)   |
| `--auth`     | `bool`   | ❌ No      | `false`           | Incluir sistema de autenticación JWT                     |
| `--api`      | `string` | ❌ No      | `rest`            | Tipo de API (`rest`, `graphql`, `grpc`)                  |

## 📖 Ejemplos de Uso

### Ejemplo Básico
```bash
goca init mi-proyecto --module github.com/usuario/mi-proyecto
```

### Proyecto con Autenticación
```bash
goca init ecommerce --module github.com/empresa/ecommerce --auth --database postgres --api rest
```

### Proyecto con gRPC
```bash
goca init microservicio --module github.com/empresa/microservicio --api grpc --database mongodb
```

### Proyecto Completo
```bash
goca init plataforma --module github.com/empresa/plataforma --auth --database mysql --api both
```

## 📂 Estructura Generada

Después de ejecutar `goca init`, obtendrás esta estructura:

```
mi-proyecto/
├── cmd/
│   └── server/
│       └── main.go              # Punto de entrada del servidor
├── internal/
│   ├── domain/                  # 🟡 Capa de Dominio
│   ├── usecase/                 # 🔴 Capa de Casos de Uso
│   ├── repository/              # 🔵 Capa de Infraestructura
│   └── handler/                 # 🟢 Capa de Adaptadores
│       ├── http/                # Handlers HTTP REST
│       ├── grpc/                # Handlers gRPC (si se selecciona)
│       └── middleware/          # Middlewares comunes
├── pkg/
│   ├── config/
│   │   ├── config.go            # Configuración de la aplicación
│   │   └── database.go          # Configuración de base de datos
│   ├── logger/
│   │   └── logger.go            # Sistema de logging
│   └── auth/                    # Sistema de autenticación (si se activa)
│       ├── jwt.go
│       ├── middleware.go
│       └── service.go
├── go.mod                       # Dependencias del módulo
├── go.sum                       # Checksums de dependencias
├── .gitignore                   # Archivos a ignorar en Git
└── README.md                    # Documentación del proyecto
```

## 🔧 Archivos Generados en Detalle

### `cmd/server/main.go`
```go
package main

import (
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/usuario/mi-proyecto/pkg/config"
    "github.com/usuario/mi-proyecto/pkg/logger"
)

func main() {
    // Cargar configuración
    cfg := config.Load()
    
    // Inicializar logger
    logger.Init(cfg.LogLevel)
    
    // Configurar router
    router := gin.Default()
    
    // Middleware global
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })
    
    // Iniciar servidor
    log.Printf("Servidor iniciado en puerto %s", cfg.Port)
    if err := router.Run(":" + cfg.Port); err != nil {
        log.Fatal("Error al iniciar servidor:", err)
    }
}
```

### `pkg/config/config.go`
```go
package config

import (
    "os"
)

type Config struct {
    Port        string
    Environment string
    LogLevel    string
    Database    DatabaseConfig
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    SSLMode  string
}

func Load() *Config {
    return &Config{
        Port:        getEnv("PORT", "8080"),
        Environment: getEnv("ENVIRONMENT", "development"),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", ""),
            Name:     getEnv("DB_NAME", "mi_proyecto"),
            SSLMode:  getEnv("DB_SSL_MODE", "disable"),
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### `pkg/logger/logger.go`
```go
package logger

import (
    "log/slog"
    "os"
)

var Logger *slog.Logger

func Init(level string) {
    var logLevel slog.Level
    
    switch level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    default:
        logLevel = slog.LevelInfo
    }
    
    Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: logLevel,
    }))
    
    slog.SetDefault(Logger)
}

func Info(msg string, args ...any) {
    Logger.Info(msg, args...)
}

func Error(msg string, args ...any) {
    Logger.Error(msg, args...)
}

func Debug(msg string, args ...any) {
    Logger.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
    Logger.Warn(msg, args...)
}
```

## 🔐 Sistema de Autenticación (--auth)

Cuando usas el flag `--auth`, se genera automáticamente:

### `pkg/auth/jwt.go`
```go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
    secretKey []byte
    issuer    string
}

type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func NewJWTService(secretKey, issuer string) *JWTService {
    return &JWTService{
        secretKey: []byte(secretKey),
        issuer:    issuer,
    }
}

func (j *JWTService) GenerateToken(userID uint, email, role string) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    j.issuer,
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secretKey)
}

func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return j.secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, jwt.ErrInvalidKey
}
```

## 🌐 Configuración de API

### REST API (--api rest)
Genera handlers HTTP con Gin:
- Routing RESTful
- Middleware de CORS
- Validación de entrada
- Manejo de errores

### gRPC API (--api grpc)
Genera configuración para gRPC:
- Archivos `.proto` base
- Configuración del servidor gRPC
- Interceptors de logging y autenticación

### Ambos (--api both)
Configura tanto REST como gRPC en el mismo proyecto.

## 💾 Bases de Datos Soportadas

### PostgreSQL (--database postgres)
```go
// Configuración automática para PostgreSQL
Database: DatabaseConfig{
    Host:     getEnv("DB_HOST", "localhost"),
    Port:     getEnv("DB_PORT", "5432"),
    User:     getEnv("DB_USER", "postgres"),
    Password: getEnv("DB_PASSWORD", ""),
    Name:     getEnv("DB_NAME", "mi_proyecto"),
    SSLMode:  getEnv("DB_SSL_MODE", "disable"),
}
```

### MySQL (--database mysql)
```go
// Configuración automática para MySQL
Database: DatabaseConfig{
    Host:     getEnv("DB_HOST", "localhost"),
    Port:     getEnv("DB_PORT", "3306"),
    User:     getEnv("DB_USER", "root"),
    Password: getEnv("DB_PASSWORD", ""),
    Name:     getEnv("DB_NAME", "mi_proyecto"),
}
```

### MongoDB (--database mongodb)
```go
// Configuración automática para MongoDB
Database: DatabaseConfig{
    URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
    Database: getEnv("MONGO_DB", "mi_proyecto"),
}
```

## 📄 Documentación Generada

### README.md
Se genera automáticamente con:
- Descripción del proyecto
- Instrucciones de instalación
- Configuración de variables de entorno
- Comandos para ejecutar el proyecto
- Estructura del proyecto explicada
- Contribución y licencia

### .gitignore
Incluye patrones comunes para proyectos Go:
```gitignore
# Binarios
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output de go build
main

# Dependencias
vendor/

# Variables de entorno
.env
.env.local

# Logs
*.log

# Base de datos local
*.db
*.sqlite

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db
```

## 🔄 Flujo de Trabajo Después del Init

1. **Entrar al directorio:**
   ```bash
   cd mi-proyecto
   ```

2. **Instalar dependencias:**
   ```bash
   go mod tidy
   ```

3. **Configurar variables de entorno:**
   ```bash
   # Crear archivo .env (opcional)
   echo "DB_PASSWORD=mipassword" > .env
   ```

4. **Generar tu primer feature:**
   ```bash
   goca feature User --fields "name:string,email:string"
   ```

5. **Configurar inyección de dependencias:**
   ```bash
   goca di --features "User" --database postgres
   ```

6. **Ejecutar el proyecto:**
   ```bash
   go run cmd/server/main.go
   ```

## ⚠️ Consideraciones Importantes

### ✅ Buenas Prácticas
- **Usar módulos descriptivos:** `github.com/empresa/proyecto` en lugar de `test` o `app`
- **Configurar Git:** Inicializar repositorio después del init
- **Variables de entorno:** Nunca commitear `.env` con credenciales reales
- **Documentación:** Actualizar README.md con información específica del proyecto

### ❌ Errores Comunes
- **Directorio existente:** No puedes usar `init` en un directorio que ya contiene archivos
- **Nombre de módulo inválido:** Debe seguir las convenciones de Go modules
- **Permisos:** Asegúrate de tener permisos de escritura en el directorio

## 🚀 Ejemplos Completos

### Proyecto de E-commerce
```bash
goca init ecommerce \
  --module github.com/miempresa/ecommerce \
  --auth \
  --database postgres \
  --api rest

cd ecommerce
go mod tidy

# Generar features principales
goca feature User --fields "name:string,email:string,password:string" --validation
goca feature Product --fields "name:string,price:float64,category:string" --validation
goca feature Order --fields "user_id:int,total:float64,status:string" --validation

# Configurar DI
goca di --features "User,Product,Order" --database postgres
```

### Microservicio gRPC
```bash
goca init user-service \
  --module github.com/miempresa/user-service \
  --auth \
  --database mongodb \
  --api grpc

cd user-service
go mod tidy

goca feature User --fields "name:string,email:string" --validation
goca handler User --type grpc
```

## 📞 Soporte

Si tienes problemas con `goca init`:

- 🔍 Revisa que el directorio esté vacío
- 📝 Verifica que el nombre del módulo sea válido
- 🐛 Reporta issues en [GitHub](https://github.com/sazardev/goca/issues)

---

**← [Instalación](Installation) | [Comando goca feature](Command-Feature) →**
