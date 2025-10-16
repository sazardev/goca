# goca init Command

The `goca init` command is the starting point for creating any project with Clean Architecture. It generates the complete base structure following established best practices and conventions.

## 📋 Syntax

```bash
goca init <project-name> [flags]
```

## 🎯 Purpose

Creates the base structure of a Go project following Clean Architecture principles, including:

- 📁 Directory structure organized by layers
- 📄 Essential configuration files
- 🔧 Initial dependency setup
- 📝 Project base documentation
- 🔐 Authentication configuration (optional)
- 🌐 API server setup (optional)

## 🚩 Available Flags

| Flag         | Type     | Required  | Default Value | Description                                                                                                         |
| ------------ | -------- | --------- | ------------- | ------------------------------------------------------------------------------------------------------------------- |
| `--module`   | `string` | ✅ **Yes** | -             | Go module name (e.g: `github.com/user/project`)                                                                     |
| `--database` | `string` | ❌ No      | `postgres`    | Database type (`postgres`, `postgres-json`, `mysql`, `mongodb`, `sqlite`, `sqlserver`, `elasticsearch`, `dynamodb`) |
| `--auth`     | `bool`   | ❌ No      | `false`       | Include JWT authentication system                                                                                   |
| `--api`      | `string` | ❌ No      | `rest`        | API type (`rest`, `graphql`, `grpc`)                                                                                |

## 📖 Usage Examples

### Basic Example
```bash
goca init my-project --module github.com/user/my-project
```

### Project with Authentication
```bash
goca init ecommerce --module github.com/company/ecommerce --auth --database postgres --api rest
```

### Project with gRPC
```bash
goca init microservice --module github.com/company/microservice --api grpc --database mongodb
```

### Complete Project
```bash
goca init platform --module github.com/company/platform --auth --database mysql --api both
```

## 📂 Generated Structure

After running `goca init`, you'll get this structure:

```
my-project/
├── cmd/
│   └── server/
│       └── main.go              # Server entry point
├── internal/
│   ├── domain/                  # 🟡 Domain Layer
│   ├── usecase/                 # 🔴 Use Cases Layer
│   ├── repository/              # 🔵 Infrastructure Layer
│   └── handler/                 # 🟢 Adapters Layer
│       ├── http/                # HTTP REST handlers
│       ├── grpc/                # gRPC handlers (if selected)
│       └── middleware/          # Common middlewares
├── pkg/
│   ├── config/
│   │   ├── config.go            # Application configuration
│   │   └── database.go          # Database configuration
│   ├── logger/
│   │   └── logger.go            # Logging system
│   └── auth/                    # Authentication system (if enabled)
│       ├── jwt.go
│       ├── middleware.go
│       └── service.go
├── go.mod                       # Module dependencies
├── go.sum                       # Dependency checksums
├── .gitignore                   # Files to ignore in Git
└── README.md                    # Project documentation
```

## 🔧 Generated Files in Detail

### `cmd/server/main.go`
```go
package main

import (
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/user/my-project/pkg/config"
    "github.com/user/my-project/pkg/logger"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Initialize logger
    logger.Init(cfg.LogLevel)
    
    // Configure router
    router := gin.Default()
    
    // Global middleware
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })
    
    // Start server
    log.Printf("Server started on port %s", cfg.Port)
    if err := router.Run(":" + cfg.Port); err != nil {
        log.Fatal("Error starting server:", err)
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
            Name:     getEnv("DB_NAME", "my_project"),
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

## 🔐 Authentication System (--auth)

When you use the `--auth` flag, it automatically generates:

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

## 🌐 API Configuration

### REST API (--api rest)
Generates HTTP handlers with Gin:
- RESTful routing
- CORS middleware
- Input validation
- Error handling

### gRPC API (--api grpc)
Generates gRPC configuration:
- Base `.proto` files
- gRPC server configuration
- Logging and authentication interceptors

### Both (--api both)
Configures both REST and gRPC in the same project.

## 💾 Supported Databases

### PostgreSQL (--database postgres)
```go
// Automatic configuration for PostgreSQL
Database: DatabaseConfig{
    Host:     getEnv("DB_HOST", "localhost"),
    Port:     getEnv("DB_PORT", "5432"),
    User:     getEnv("DB_USER", "postgres"),
    Password: getEnv("DB_PASSWORD", ""),
    Name:     getEnv("DB_NAME", "my_project"),
    SSLMode:  getEnv("DB_SSL_MODE", "disable"),
}
```

### MySQL (--database mysql)
```go
// Automatic configuration for MySQL
Database: DatabaseConfig{
    Host:     getEnv("DB_HOST", "localhost"),
    Port:     getEnv("DB_PORT", "3306"),
    User:     getEnv("DB_USER", "root"),
    Password: getEnv("DB_PASSWORD", ""),
    Name:     getEnv("DB_NAME", "my_project"),
}
```

### MongoDB (--database mongodb)
```go
// Automatic configuration for MongoDB
Database: DatabaseConfig{
    URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
    Database: getEnv("MONGO_DB", "my_project"),
}
```

## 📄 Generated Documentation

### README.md
Automatically generated with:
- Project description
- Installation instructions
- Environment variable configuration
- Commands to run the project
- Project structure explained
- Contribution and license

### .gitignore
Includes common patterns for Go projects:
```gitignore
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output from go build
main

# Dependencies
vendor/

# Environment variables
.env
.env.local

# Logs
*.log

# Local database
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

## 🔄 Workflow After Init

1. **Enter the directory:**
   ```bash
   cd my-project
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Configure environment variables:**
   ```bash
   # Create .env file (optional)
   echo "DB_PASSWORD=mypassword" > .env
   ```

4. **Generate your first feature:**
   ```bash
   goca feature User --fields "name:string,email:string"
   ```

5. **Configure dependency injection:**
   ```bash
   goca di --features "User" --database postgres
   ```

6. **Run the project:**
   ```bash
   go run cmd/server/main.go
   ```

## ⚠️ Important Considerations

### ✅ Best Practices
- **Use descriptive modules:** `github.com/company/project` instead of `test` or `app`
- **Configure Git:** Initialize repository after init
- **Environment variables:** Never commit `.env` with real credentials
- **Documentation:** Update README.md with project-specific information

### ❌ Common Errors
- **Existing directory:** You cannot use `init` in a directory that already contains files
- **Invalid module name:** Must follow Go modules conventions
- **Permissions:** Make sure you have write permissions in the directory

## 🚀 Complete Examples

### E-commerce Project
```bash
goca init ecommerce \
  --module github.com/mycompany/ecommerce \
  --auth \
  --database postgres \
  --api rest

cd ecommerce
go mod tidy

# Generate main features
goca feature User --fields "name:string,email:string,password:string" --validation
goca feature Product --fields "name:string,price:float64,category:string" --validation
goca feature Order --fields "user_id:int,total:float64,status:string" --validation

# Configure DI
goca di --features "User,Product,Order" --database postgres
```

### gRPC Microservice
```bash
goca init user-service \
  --module github.com/mycompany/user-service \
  --auth \
  --database mongodb \
  --api grpc

cd user-service
go mod tidy

goca feature User --fields "name:string,email:string" --validation
goca handler User --type grpc
```

## 📞 Support

If you have problems with `goca init`:

- 🔍 Check that the directory is empty
- 📝 Verify that the module name is valid
- 🐛 Report issues on [GitHub](https://github.com/sazardev/goca/issues)

---

**← [Installation](Installation) | [goca feature Command](Command-Feature) →**
