# GOCA YAML Configuration Structure Reference

**Version:** 1.0  
**Last Updated:** 2025-09-30

## 🎯 Purpose

This document provides the **authoritative reference** for `.goca.yaml` configuration structure. All YAML examples in documentation, tests, and user projects should follow these exact structures.

---

## ⚠️ Common Mistakes to Avoid

### ❌ INCORRECT Structures

```yaml
# WRONG: timestamps and soft_delete under generation
generation:
  timestamps:
    enabled: true
  soft_delete:
    enabled: true
```

```yaml
# WRONG: validation as simple boolean
generation:
  validation: true
```

### ✅ CORRECT Structures

```yaml
# CORRECT: timestamps and soft_delete under database.features
database:
  features:
    timestamps: true
    soft_delete: true
```

```yaml
# CORRECT: validation as nested object
generation:
  validation:
    enabled: true
    library: "validator"
```

---

## 📋 Complete Structure Reference

### 1. Project Configuration

```yaml
project:
  name: "my-project"              # Required: Project name
  module: "github.com/user/repo"  # Required: Go module path
  version: "1.0.0"                # Optional: Semantic version
  description: "Description"      # Optional: Project description
  author: "Author Name"           # Optional: Project author
  license: "MIT"                  # Optional: License type
  repository: "https://..."       # Optional: Repository URL
  tags: ["api", "rest"]           # Optional: Project tags
  metadata: {}                    # Optional: Custom metadata
```

**Required Fields:**
- `project.name`
- `project.module`

---

### 2. Architecture Configuration

```yaml
architecture:
  pattern: "clean_architecture"   # Architecture pattern
  
  layers:
    domain:
      enabled: true
      directory: "internal/domain"
    usecase:
      enabled: true
      directory: "internal/usecase"
    repository:
      enabled: true
      directory: "internal/repository"
    handler:
      enabled: true
      directory: "internal/handler"
  
  di:
    type: "manual"                # Options: manual, wire, fx, dig
    auto_wire: false
    providers: []
    modules: []
  
  naming:
    entities: "PascalCase"        # PascalCase, camelCase, snake_case
    fields: "camelCase"           # Field naming convention
    files: "snake_case"           # File naming: snake_case, kebab-case, camelCase
    packages: "lowercase"         # Package naming
    constants: "UPPER_CASE"       # Constant naming
    variables: "camelCase"        # Variable naming
    functions: "camelCase"        # Function naming
```

**Key Points:**
- `architecture.naming.files` controls generated filename convention
- Valid values: `PascalCase`, `camelCase`, `snake_case`, `kebab-case`, `UPPER_CASE`, `lowercase`

---

### 3. Database Configuration

```yaml
database:
  type: "postgres"                # Required: postgres, mysql, mongodb, sqlite
  host: "localhost"
  port: 5432
  name: "mydb"
  
  migrations:
    enabled: true
    auto_generate: true
    directory: "migrations"
    naming: "timestamp"
    versioning: "sequential"
    tools: ["migrate", "sql-migrate"]
  
  connection:
    max_open: 25
    max_idle: 5
    max_lifetime: "5m"
    ssl_mode: "disable"
    timezone: "UTC"
    charset: "utf8mb4"
    collation: "utf8mb4_unicode_ci"
  
  features:
    soft_delete: true             # ✅ CORRECT LOCATION
    timestamps: true              # ✅ CORRECT LOCATION
    uuid: false
    audit: false
    versioning: false
    partitioning: false
    full_text_search: false
    json_support: false
  
  extensions: ["uuid-ossp", "pgcrypto"]
  custom_types: {}
```

**Key Points:**
- ✅ **`database.features.timestamps`** - Adds CreatedAt/UpdatedAt fields
- ✅ **`database.features.soft_delete`** - Adds DeletedAt field and soft delete methods
- Valid `database.type` values: `postgres`, `mysql`, `mongodb`, `sqlite`

---

### 4. Generation Configuration

```yaml
generation:
  validation:
    enabled: true                 # ✅ CORRECT LOCATION
    library: "validator"          # Options: builtin, validator, ozzo-validation
    tags: ["required", "min", "max", "email"]
    custom: []
    sanitize: false
    transform: false
  
  business_rules:
    enabled: false
    patterns: ["domain-events", "saga"]
    templates: []
    events: false
    guards: false
  
  documentation:
    swagger:
      enabled: true
      version: "3.0"
      output: "docs/swagger.yaml"
      title: "API Documentation"
      description: "API Description"
      host: "localhost:8080"
      base_path: "/api/v1"
      schemes: ["http", "https"]
      tags: []
      extensions: {}
    
    postman:
      enabled: false
      output: "docs/postman"
      environment: true
      tests: true
      variables: true
    
    markdown:
      enabled: true
      output: "docs"
      format: "github"
      toc: true
      examples: true
      diagrams: false
    
    comments:
      enabled: true
      style: "godoc"
      detail_level: "detailed"
      examples: true
      links: false
  
  style:
    line_length: 100
    indent: "tabs"
    bracket_style: "k&r"
    import_grouping: true
    comment_wrapping: true
  
  imports:
    grouping: true
    order: ["stdlib", "external", "internal"]
    aliases: {}
    format: "goimports"
```

**Key Points:**
- ✅ **`generation.validation.enabled`** - Controls validation tag generation
- ❌ **NOT** `generation.timestamps` - timestamps belong under `database.features`

---

### 5. Testing Configuration

```yaml
testing:
  framework: "testify"            # Options: testify, ginkgo, builtin
  
  coverage:
    enabled: true
    threshold: 80
    exclude: ["*_test.go", "mocks/*"]
    report_format: "html"
  
  mocking:
    enabled: true
    tool: "mockery"
    output: "mocks"
    interfaces: ["*Repository", "*UseCase"]
  
  integration:
    enabled: false
    database: true
    api: true
    containers: false
  
  benchmarking:
    enabled: false
    duration: "10s"
    count: 3
  
  fixtures:
    enabled: true
    directory: "testdata"
    format: "json"
```

---

### 6. Features Configuration

```yaml
features:
  authentication:
    enabled: false
    type: "jwt"                   # Options: jwt, oauth2, session, basic
    providers: ["local"]
    token_expiry: "24h"
    refresh_enabled: false
  
  authorization:
    enabled: false
    type: "rbac"
    roles: []
    permissions: []
  
  caching:
    enabled: false
    type: "redis"                 # Options: redis, memcached, inmemory
    ttl: "1h"
    max_entries: 1000
  
  logging:
    enabled: true
    level: "info"
    format: "json"
    output: "stdout"
    file: "logs/app.log"
    rotation: true
  
  monitoring:
    enabled: false
    prometheus: false
    jaeger: false
    sentry: false
  
  api:
    versioning: true
    rate_limiting: false
    cors: true
    compression: true
```

---

### 7. Deploy Configuration

```yaml
deploy:
  docker:
    enabled: true
    multi_stage: true
    base_image: "golang:1.21-alpine"
    expose_port: 8080
  
  kubernetes:
    enabled: false
    namespace: "default"
    replicas: 3
    resources:
      limits:
        cpu: "500m"
        memory: "512Mi"
      requests:
        cpu: "250m"
        memory: "256Mi"
  
  cloud:
    provider: ""
    region: ""
    environment: "development"
```

---

## 🔍 How to Verify Your YAML Structure

### Method 1: Use `goca config show`

```bash
cd your-project
goca config show
```

This will display your current configuration and highlight any errors.

### Method 2: Use `goca config validate`

```bash
goca config validate
```

This validates your `.goca.yaml` without making any changes.

---

## 📝 Quick Reference Card

### Timestamps & Soft Delete
```yaml
database:
  features:
    timestamps: true    # CreatedAt, UpdatedAt
    soft_delete: true   # DeletedAt, SoftDelete(), IsDeleted()
```

### Validation
```yaml
generation:
  validation:
    enabled: true
    library: "validator"
```

### File Naming Convention
```yaml
architecture:
  naming:
    files: "snake_case"  # product.go → product.go
                         # ProductCategory → product_category.go
```

### Database Type
```yaml
database:
  type: "postgres"  # postgres, mysql, mongodb, sqlite
```

### Handler Types
```yaml
# CLI flag (not in YAML):
goca handler Product --type http    # HTTP handler
goca handler Product --type grpc    # gRPC handler
goca handler Product --type cli     # CLI handler
goca handler Product --type worker  # Worker handler
goca handler Product --type soap    # SOAP handler
```

---

## 🧪 Testing Your Configuration

### Example Test Project

```bash
# Create test directory
mkdir test-config
cd test-config

# Create .goca.yaml with correct structure
cat > .goca.yaml << 'EOF'
project:
  name: test-project
  module: github.com/test/project

database:
  type: postgres
  features:
    timestamps: true
    soft_delete: true

generation:
  validation:
    enabled: true
    library: validator

architecture:
  naming:
    files: snake_case
EOF

# Generate entity
goca entity ProductCategory --fields "name:string,price:float64"

# Verify generated file
ls internal/domain/product_category.go  # Should exist with snake_case name
cat internal/domain/product_category.go  # Should have validate tags, timestamps, soft delete
```

**Expected Output:**
- File: `product_category.go` (snake_case naming)
- Contains: `validate:"required"` tags
- Contains: `CreatedAt`, `UpdatedAt` fields
- Contains: `DeletedAt`, `SoftDelete()`, `IsDeleted()`

---

## 🎯 Summary

### ✅ DO Use These Structures

- `database.features.timestamps: true`
- `database.features.soft_delete: true`
- `generation.validation.enabled: true`
- `architecture.naming.files: "snake_case"`

### ❌ DON'T Use These Structures

- ~~`generation.timestamps.enabled: true`~~
- ~~`generation.soft_delete.enabled: true`~~
- ~~`generation.validation: true`~~ (use nested object)
- ~~`naming.file_convention: "snake_case"`~~ (use `architecture.naming.files`)

---

## 📚 Related Documentation

- [Configuration System Guide](configuration-system.md)
- [Migration Guide](migration-guide.md)
- [Complete Tutorial](Complete-Tutorial.md)
- [Configuration Integration Status](CONFIGURATION_INTEGRATION_COMPLETE.md)

---

## 🔗 Code Reference

The authoritative structure is defined in:
- **File:** `cmd/config_types.go`
- **Type:** `GocaConfig` struct

Always refer to this file for the ground truth of the configuration structure.

---

**Last Verified:** 2025-09-30  
**Status:** ✅ Production Ready  
**Test Coverage:** 100% (74/74 tests passing)
