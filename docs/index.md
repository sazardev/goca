---
layout: home

hero:
  name: GOCA CLI
  text: Clean Architecture Code Generator
  tagline: Generate production-ready Go projects with Clean Architecture in seconds
  image:
    src: /logo.svg
    alt: GOCA CLI
  actions:
    - theme: brand
      text: Get Started
      link: /getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/sazardev/goca

features:
  - icon: 🚀
    title: Production Ready
    details: Generate complete, tested, and production-ready Go code following Clean Architecture principles.
  
  - icon: ⚡
    title: Fast Development
    details: Create full features with CRUD operations, validation, and business rules in seconds.
  
  - icon: 🎯
    title: Clean Architecture
    details: Enforces proper layer separation with domain, usecase, repository, and handler layers.
  
  - icon: 🔧
    title: Flexible Configuration
    details: YAML-based configuration for project-wide settings, database types, and code generation options.
  
  - icon: 🗄️
    title: Multiple Databases
    details: Support for PostgreSQL, MySQL, SQLite with automatic repository implementation.
  
  - icon: 🌐
    title: Multiple Handlers
    details: Generate HTTP REST, gRPC, and CLI handlers for your features automatically.
  
  - icon: ✅
    title: Built-in Validation
    details: Automatic field validation with customizable rules and error handling.
  
  - icon: 📦
    title: Zero Config
    details: Works out of the box with sensible defaults, customize only what you need.
  
  - icon: 🔄
    title: Auto Integration
    details: Automatically integrates new features into dependency injection and main.go.
---

## Quick Start

```bash
# Install GOCA
go install github.com/sazardev/goca@v2.0.0

# Initialize a new project
goca init myproject --database postgres --handlers http

# Generate a complete feature
cd myproject
goca feature user --fields "name:string,email:string,age:int" --validation --business-rules
```

## What's New in v2.0.0

🎉 **Major Release** - Complete documentation system with VitePress

- ✅ **VitePress Documentation** - Beautiful, searchable documentation site
- ✅ **5 Critical Bugs Fixed** - GORM imports, time imports, domain imports, MySQL config, kebab-case
- ✅ **Production Ready** - 100% tested with comprehensive test suite
- ✅ **Complete Documentation** - 2,500+ lines of documentation
- ✅ **Enhanced YAML Config** - Full .goca.yaml support with validation

[Read the full release notes →](/releases/v2.0.0)

## Features

### Generate Complete Features

```bash
goca feature product \
  --fields "name:string,price:float64,stock:int,category:string" \
  --validation \
  --business-rules
```

Creates:
- ✅ Domain entities with validation
- ✅ Use cases with business logic
- ✅ Repository implementation
- ✅ HTTP/gRPC/CLI handlers
- ✅ Dependency injection setup
- ✅ Database migrations

### YAML Configuration

```yaml
# .goca.yaml
project:
  name: "ecommerce-api"
  module: "github.com/myorg/ecommerce"
  
database:
  type: "postgres"
  
handlers:
  - "http"
  - "grpc"
  
features:
  validation: true
  business_rules: true
  soft_delete: true
  timestamps: true
```

### Type-Safe Code Generation

All generated code is:
- ✅ **Type-safe** - Full Go type checking
- ✅ **Tested** - Compiles with zero errors
- ✅ **Documented** - Inline comments and documentation
- ✅ **Idiomatic** - Follows Go best practices

## Architecture

GOCA follows Clean Architecture principles with clear layer separation:

```
project/
├── internal/
│   ├── domain/          # Entities and business rules
│   ├── usecase/         # Business logic
│   ├── repository/      # Data access
│   └── handler/         # Interface adapters
├── cmd/
│   └── server/          # Application entry point
├── pkg/
│   ├── config/          # Configuration
│   └── database/        # Database setup
└── migrations/          # Database migrations
```

## Why GOCA?

- **Save Time** - Generate complete features in seconds instead of hours
- **Best Practices** - Enforces Clean Architecture and Go idioms
- **Production Ready** - All code is tested and verified
- **Flexible** - Customize through YAML or command flags
- **Zero Dependencies** - Single binary, no runtime dependencies

## Community

- [GitHub Discussions](https://github.com/sazardev/goca/discussions)
- [Report Issues](https://github.com/sazardev/goca/issues)
- [Contributing Guide](https://github.com/sazardev/goca/blob/master/CONTRIBUTING.md)

## License

GOCA is released under the [MIT License](https://github.com/sazardev/goca/blob/master/LICENSE).
