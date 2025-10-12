# goca init

Initialize a new Clean Architecture project with complete structure and configuration.

## Syntax

```bash
goca init <project-name> [flags]
```

## Description

The `goca init` command creates a production-ready Go project following Clean Architecture principles. It generates the complete directory structure, configuration files, and boilerplate code to get you started immediately.

## Arguments

### `<project-name>`

**Required.** The name of your project directory.

```bash
goca init my-api
```

This creates a directory named `my-api` with the full project structure.

## Flags

### `--module` (Required)

Go module name for your project.

```bash
--module github.com/username/projectname
```

**Example:**
```bash
goca init ecommerce --module github.com/sazardev/ecommerce
```

::: tip Module Naming Convention
Use your repository URL as the module name:
- GitHub: `github.com/username/repo`
- GitLab: `gitlab.com/username/repo`
- Custom: `example.com/project`
:::

### `--database`

Database system to use. Default: `postgres`

**Options:** `postgres` | `mysql` | `mongodb` | `sqlite`

```bash
goca init myproject --module github.com/user/myproject --database mysql
```

### `--auth`

Include JWT authentication system.

```bash
goca init myproject --module github.com/user/myproject --auth
```

Generates:
- JWT token generation and validation
- Authentication middleware
- User authentication endpoints
- Password hashing utilities

### `--api`

API type to generate. Default: `rest`

**Options:** `rest` | `grpc` | `graphql` | `both`

```bash
goca init myproject --module github.com/user/myproject --api grpc
```

## Examples

### Basic REST API

```bash
goca init blog-api \
  --module github.com/sazardev/blog-api \
  --database postgres
```

### E-commerce with Authentication

```bash
goca init ecommerce \
  --module github.com/company/ecommerce \
  --database postgres \
  --auth
```

### gRPC Microservice

```bash
goca init user-service \
  --module github.com/company/user-service \
  --database mongodb \
  --api grpc
```

### Full-Featured Application

```bash
goca init platform \
  --module github.com/startup/platform \
  --database postgres \
  --auth \
  --api both
```

## Generated Structure

```
myproject/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/                  # 🟡 Entities & business rules
│   │   └── errors.go
│   ├── usecase/                 # 🔴 Application logic
│   │   ├── dto.go
│   │   └── interfaces.go
│   ├── repository/              # 🔵 Data access
│   │   └── interfaces.go
│   └── handler/                 # 🟢 Input adapters
│       ├── http/
│       │   ├── routes.go
│       │   └── middleware.go
│       └── grpc/                # (if --api grpc)
├── pkg/
│   ├── config/
│   │   ├── config.go            # App configuration
│   │   └── database.go          # DB connection
│   ├── logger/
│   │   └── logger.go            # Structured logging
│   └── auth/                    # (if --auth)
│       ├── jwt.go
│       ├── middleware.go
│       └── password.go
├── migrations/                   # Database migrations
│   └── 001_initial.sql
├── .env.example                 # Environment variables template
├── .gitignore
├── go.mod
├── go.sum
├── Makefile                     # Common tasks
└── README.md
```

## Generated Files

### `cmd/server/main.go`

The application entry point with:
- Server initialization
- Database connection
- Route registration
- Graceful shutdown

```go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/user/myproject/pkg/config"
    "github.com/user/myproject/pkg/logger"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Initialize logger
    log := logger.New(cfg.LogLevel)
    
    // Connect to database
    db := config.ConnectDatabase(cfg)
    
    // Start server
    server := NewServer(cfg, db, log)
    
    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Info("Shutting down server...")
}
```

### `pkg/config/config.go`

Configuration management:

```go
package config

import "github.com/spf13/viper"

type Config struct {
    ServerPort   string
    DatabaseURL  string
    LogLevel     string
    JWTSecret    string // if --auth
}

func Load() *Config {
    viper.AutomaticEnv()
    // Load configuration
    return &Config{...}
}
```

### `.env.example`

Environment variables template:

```bash
SERVER_PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/dbname
LOG_LEVEL=info
JWT_SECRET=your-secret-key  # if --auth
```

### `Makefile`

Common development tasks:

```makefile
.PHONY: run build test

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test ./...

migrate-up:
	migrate -path migrations -database $(DATABASE_URL) up

migrate-down:
	migrate -path migrations -database $(DATABASE_URL) down
```

## Next Steps

After running `goca init`:

1. **Navigate to project directory:**
   ```bash
   cd myproject
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Configure environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

4. **Generate your first feature:**
   ```bash
   goca feature User --fields "name:string,email:string"
   ```

5. **Run the application:**
   ```bash
   make run
   # or
   go run cmd/server/main.go
   ```

## Tips

### Use Configuration Files

Create a `.goca.yaml` for reusable settings:

```yaml
module: github.com/company/projectname
database: postgres
auth: true
api: rest
```

Then simply run:
```bash
goca init myproject
```

### Customize Templates

After initialization, you can modify the generated code to fit your needs. The structure is designed to be a starting point, not a constraint.

### Version Control

Don't forget to initialize Git:

```bash
cd myproject
git init
git add .
git commit -m "Initial commit with Clean Architecture structure"
```

## Troubleshooting

### Module Name Errors

**Problem:** "invalid module name"

**Solution:** Ensure module name follows Go conventions:
```bash
#  Correct
--module github.com/user/project

#  Incorrect
--module my-project
--module Project Name
```

### Permission Denied

**Problem:** "permission denied creating directory"

**Solution:** Run with appropriate permissions or choose a directory you have write access to.

### Dependencies Not Found

**Problem:** Generated project can't find dependencies

**Solution:**
```bash
cd myproject
go mod tidy
go mod download
```

## See Also

- [`goca feature`](/commands/feature) - Generate complete features
- [`goca integrate`](/commands/integrate) - Wire features together
- [Getting Started Guide](/getting-started) - Complete walkthrough
- [Project Structure](/guide/project-structure) - Understand the layout

## Resources

- [GitHub Repository](https://github.com/sazardev/goca)
- [Example Projects](https://github.com/sazardev/goca-examples)
- [Report Issues](https://github.com/sazardev/goca/issues)
