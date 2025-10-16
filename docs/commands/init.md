# goca init

Initialize a new Clean Architecture project with complete structure and configuration.

## Syntax

```bash
goca init <project-name> [flags]
```

## Description

The `goca init` command creates a production-ready Go project following Clean Architecture principles. It generates the complete directory structure, configuration files, and boilerplate code to get you started immediately.

::: tip Git Initialization
Projects are automatically initialized with Git, including an initial commit. This ensures your project is version-control ready from the start.
:::

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

**Options:**
- `postgres` - PostgreSQL (GORM)
- `postgres-json` - PostgreSQL with JSONB
- `mysql` - MySQL (GORM)
- `mongodb` - MongoDB (native driver)
- `sqlite` - SQLite (embedded)
- `sqlserver` - SQL Server
- `elasticsearch` - Elasticsearch (v8)
- `dynamodb` - DynamoDB (AWS)

```bash
goca init myproject --module github.com/user/myproject --database mysql
goca init config-server --module github.com/user/config --database postgres-json
goca init search-app --module github.com/user/search --database elasticsearch
```

See [Database Support](/features/database-support) for detailed comparison.

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

## Project Templates

GOCA provides predefined project templates that automatically configure your project with optimized settings for specific use cases. Templates generate a complete `.goca.yaml` configuration file tailored to your project type.

### Using Templates

#### `--template`

Initialize project with a predefined configuration template.

```bash
goca init myproject --module github.com/user/myproject --template rest-api
```

#### `--list-templates`

List all available templates with descriptions.

```bash
goca init --list-templates
```

### Available Templates

#### `minimal`

**Lightweight starter with essential features only**

Perfect for:
- Quick prototypes
- Learning Clean Architecture
- Minimal dependencies

```bash
goca init quick-start \
  --module github.com/user/quick-start \
  --template minimal
```

**Includes:**
- Basic project structure
- PostgreSQL database
- Essential layers (domain, usecase, repository, handler)
- Simple validation
- Testify for testing

#### `rest-api`

**Production-ready REST API with PostgreSQL, validation, and testing**

Perfect for:
- RESTful web services
- API backends
- Standard CRUD applications

```bash
goca init api-service \
  --module github.com/company/api-service \
  --template rest-api
```

**Includes:**
- Complete Clean Architecture layers
- PostgreSQL with migrations
- Input validation and sanitization
- Swagger/OpenAPI documentation
- Comprehensive testing with testify
- Test coverage (70% threshold)
- Integration tests
- Soft deletes and timestamps

#### `microservice`

**Microservice with gRPC, events, and comprehensive testing**

Perfect for:
- Distributed systems
- Event-driven architecture
- Service-oriented architecture

```bash
goca init user-service \
  --module github.com/company/user-service \
  --template microservice
```

**Includes:**
- UUID primary keys
- Audit logging
- Event-driven patterns
- Domain events support
- Specification pattern
- Advanced validation (validator library)
- High test coverage (80% threshold)
- Integration and benchmark tests
- Optimized for horizontal scaling

#### `monolith`

**Full-featured monolithic application with web interface**

Perfect for:
- Traditional web applications
- Internal tools
- Admin panels

```bash
goca init admin-panel \
  --module github.com/company/admin-panel \
  --template monolith
```

**Includes:**
- JWT authentication with RBAC
- Redis caching
- Structured logging (JSON)
- Health check endpoints
- Soft deletes and timestamps
- Audit trail
- Versioning support
- Markdown documentation
- Test fixtures and seeds
- Guards and authorization patterns

#### `enterprise`

**Enterprise-grade with all features, security, and monitoring**

Perfect for:
- Production applications
- Enterprise systems
- Mission-critical services

```bash
goca init enterprise-app \
  --module github.com/corp/enterprise-app \
  --template enterprise
```

**Includes:**
- **Security**: HTTPS, CORS, rate limiting, header security
- **Authentication**: JWT + OAuth2, RBAC
- **Caching**: Redis with multi-layer caching
- **Monitoring**: Prometheus metrics, distributed tracing, health checks, profiling
- **Documentation**: Swagger 3.0, Postman collections, comprehensive markdown
- **Testing**: 85% coverage threshold, mocks, integration, benchmarks, examples
- **Deployment**: Docker (multistage), Kubernetes (manifests, Helm), CI/CD (GitHub Actions)
- **Code Quality**: gofmt, goimports, golint, staticcheck
- **Database**: Advanced features (partitioning, connection pooling)

### Template Configuration

When you initialize a project with a template, GOCA:

1. Creates the standard project structure
2. Generates a `.goca.yaml` configuration file with template settings
3. All future feature generation uses these settings automatically
4. You can still override settings using CLI flags when needed

**Example workflow:**

```bash
# Initialize with template
goca init my-api --module github.com/user/my-api --template rest-api

# Navigate to project
cd my-api

# Generate features - automatically uses template configuration
goca feature Product --fields "name:string,price:float64,stock:int"
# ✓ Uses REST API settings from template
# ✓ Includes validation
# ✓ Generates Swagger docs
# ✓ Creates comprehensive tests

# Override specific settings if needed
goca feature Order --fields "total:float64" --database mysql
# ✓ Uses template settings except database
```

### Customizing Template Configuration

After initialization, you can customize the generated `.goca.yaml`:

```bash
# Initialize project
goca init my-project --module github.com/user/my-project --template rest-api

# Edit configuration
cd my-project
vim .goca.yaml  # Customize as needed

# Features use your customized settings
goca feature User --fields "name:string,email:string"
```

See [Configuration Guide](/guide/configuration) for detailed `.goca.yaml` documentation.

### Choosing the Right Template

| Template       | Best For                | Complexity    | Features               |
| -------------- | ----------------------- | ------------- | ---------------------- |
| `minimal`      | Learning, prototypes    | ⭐ Simple      | Essential only         |
| `rest-api`     | Web APIs, CRUD services | ⭐⭐ Standard   | Production-ready API   |
| `microservice` | Distributed systems     | ⭐⭐⭐ Advanced  | Events, gRPC, scaling  |
| `monolith`     | Web applications        | ⭐⭐⭐ Advanced  | Auth, caching, logging |
| `enterprise`   | Mission-critical apps   | ⭐⭐⭐⭐ Complete | Everything included    |

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

After initialization, follow these steps:

1. **Navigate to project:**
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

4. **Verify Git initialization:**
   ```bash
   git log --oneline  # See initial commit
   git status         # Check repository status
   ```

5. **Generate your first feature:**
   ```bash
   goca feature User --fields "name:string,email:string"
   ```

6. **Run the application:**
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
