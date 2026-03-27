# Commands Overview

Goca provides a comprehensive set of commands to generate Clean Architecture components and manage your project structure.

## Command Categories

### Project Initialization
- [`goca init`](/commands/init) - Initialize a new Clean Architecture project

### Complete Features
- [`goca feature`](/commands/feature) - Generate a complete feature with all layers
- [`goca integrate`](/commands/integrate) - Integrate existing features with DI and routing

### Layer-Specific Generation

#### Domain Layer
- [`goca entity`](/commands/entity) - Generate domain entities

#### Application Layer
- [`goca usecase`](/commands/usecase) - Generate use cases and DTOs
- [`goca interfaces`](/commands/interfaces) - Generate interface contracts

#### Infrastructure Layer
- [`goca repository`](/commands/repository) - Generate repositories

#### Adapter Layer
- [`goca handler`](/commands/handler) - Generate handlers (HTTP, gRPC, CLI, etc.)
- [`goca middleware`](/commands/middleware) - Generate composable HTTP middleware package

### Configuration & Templates
- [`goca config`](/commands/config) - Manage `.goca.yaml` configuration files
- [`goca template`](/commands/template) - Manage custom code generation templates

### Utilities
- [`goca di`](/commands/di) - Generate dependency injection container
- [`goca messages`](/commands/messages) - Generate error messages and constants
- [`goca mocks`](/commands/mocks) - Generate testify/mock mocks for all interfaces
- [`goca doctor`](/commands/doctor) - Check project health and Clean Architecture structure
- [`goca analyze`](/commands/analyze) - Deep self-analysis: architecture, security, quality, standards, tests, dependencies
- [`goca ci`](/commands/ci) - Generate CI/CD pipeline configuration (GitHub Actions)
- [`goca upgrade`](/commands/upgrade) - Upgrade project configuration to current Goca version
- [`goca version`](/commands/version) - Display version information

### Testing
- [`goca test-integration`](/commands/test-integration) - Generate integration test scaffolding

## Quick Reference

| Command                   | Purpose                          | Auto-Integration |
| ------------------------- | -------------------------------- | ---------------- |
| `goca init`               | Create new project               |  Complete setup |
| `goca feature`            | Generate full feature            |  Automatic      |
| `goca integrate`          | Wire existing features           |  Automatic      |
| `goca entity`             | Create entities only             |  Manual         |
| `goca usecase`            | Create use cases only            |  Manual         |
| `goca repository`         | Create repositories only         |  Manual         |
| `goca handler`            | Create handlers only             |  Manual         |
| `goca middleware`         | Generate HTTP middleware package  |  Manual         |
| `goca di`                 | Generate DI container            |  Manual         |
| `goca interfaces`         | Generate interface contracts     |  Manual         |
| `goca messages`           | Generate error message constants |  Manual         |
| `goca mocks`              | Generate testify/mock mocks      |  Manual         |
| `goca test-integration`   | Generate integration test files  |  Manual         |
| `goca config`             | Manage project configuration     |  —              |
| `goca template`           | Manage custom templates          |  —              |
| `goca ci`                 | Generate CI/CD pipelines         |  —              |
| `goca doctor`             | Project health checks            |  —              |
| `goca analyze`            | Deep project self-analysis       |  —              |
| `goca upgrade`            | Upgrade config/metadata          |  —              |

## Common Workflows

### Workflow 1: Quick Start (Recommended)

```bash
# Initialize project
goca init myproject --module github.com/user/myproject

# Generate complete features
goca feature User --fields "name:string,email:string"
goca feature Product --fields "name:string,price:float64"

# Everything is integrated automatically!
go run cmd/server/main.go
```

### Workflow 2: Layer-by-Layer

```bash
# Generate entity
goca entity Order --fields "customer:string,total:float64"

# Generate use case
goca usecase OrderService --entity Order

# Generate repository
goca repository Order --database postgres

# Generate handler
goca handler Order --type http

# Wire everything together
goca integrate --all
```

### Workflow 3: Add to Existing Project

```bash
# Generate new feature
goca feature Payment --fields "amount:float64,method:string"

# Automatically integrated with existing features
```

## Global Flags

All commands support these flags:

```bash
--help, -h          Show help for command
--verbose, -v       Enable verbose output (includes debug details)
--quiet, -q         Suppress all output except errors and success messages
--dry-run           Show what would be generated without creating files
--no-color          Disable colored output
--no-interactive    Disable interactive prompts
```

## Examples by Use Case

### Building a REST API

```bash
goca init ecommerce-api --module github.com/user/ecommerce
cd ecommerce-api

goca feature Product --fields "name:string,price:float64,stock:int"
goca feature Order --fields "customer:string,total:float64,status:string"
goca feature User --fields "name:string,email:string,role:string"
```

### Building a Microservice

```bash
goca init payment-service --module github.com/user/payment
cd payment-service

goca feature Payment --fields "amount:float64,currency:string,status:string"
goca handler Payment --type grpc
goca handler Payment --type http
```

### Building a CLI Tool

```bash
goca init data-processor --module github.com/user/processor
cd data-processor

goca feature DataProcessor --fields "input:string,output:string"
goca handler DataProcessor --type cli
```

## Next Steps

Choose a command to learn more:

### Essential Commands
-  [goca init](/commands/init) - Start here
-  [goca feature](/commands/feature) - Fastest way to add features
-  [goca integrate](/commands/integrate) - Wire everything together

### Detailed Generation
- 🟡 [goca entity](/commands/entity) - Domain layer
- 🔴 [goca usecase](/commands/usecase) - Application layer
- 🔵 [goca repository](/commands/repository) - Infrastructure layer
- 🟢 [goca handler](/commands/handler) - Adapter layer

### Utilities
-  [goca di](/commands/di) - Dependency injection
- 📝 [goca messages](/commands/messages) - Messages and constants
-  [goca version](/commands/version) - Version info
