# 📋 YAML Configuration System for GOCA

## 🎯 Overview

GOCA now includes a powerful YAML configuration system that allows you to centralize and reuse project configurations. The `.goca.yaml` file lets you define your configuration once and apply it consistently across all commands.

## 🚀 Configuration System Benefits

- **⚡ Productivity**: Configure once, use everywhere
- **📐 Consistency**: Maintain uniform patterns across all your projects  
- **🔄 Reusability**: Share configurations between teams
- **🛡️ Validation**: Automatic error detection in configuration
- **🎨 Customization**: Custom templates and code styles
- **🧩 Modularity**: Layer-based and feature-based configuration

## 📁 .goca.yaml File Structure

```yaml
# Project configuration
project:
  name: "my-project"
  module: "github.com/user/my-project"
  description: "Project description"
  version: "1.0.0"
  license: "MIT"
  tags: ["api", "clean-architecture", "go"]

# Architecture and layers
architecture:
  layers:
    domain:
      enabled: true
      path: "internal/domain"
    usecase:
      enabled: true 
      path: "internal/usecase"
    repository:
      enabled: true
      path: "internal/repository"
    handler:
      enabled: true
      path: "internal/handler"

# Database
database:
  type: "postgres"  # postgres, mysql, mongodb, sqlite
  migrations:
    enabled: true
    auto_generate: true

# Code generation
generation:
  validation:
    enabled: true
    tags: ["required", "min", "max", "email"]
  business_rules:
    enabled: false

# Testing
testing:
  enabled: true
  coverage_threshold: 80.0

# Custom templates
templates:
  directory: ".goca/templates"

# Project features
features:
  auth:
    enabled: false
  cors:
    enabled: true
  logging:
    enabled: true

# Deployment
deploy:
  docker:
    enabled: true
```

## 🛠️ Configuration System Commands

### Initialize with Configuration

```bash
# Create project with automatic .goca.yaml file
goca init my-project --config

# Create project with custom configuration
goca init my-project --database mysql --auth --config
```

### Configuration Management

```bash
# View current configuration
goca config show

# See available templates  
goca config template

# Initialize configuration system
goca config init
```

### Using Configuration in Commands

```bash
# Before: Specifying options every time
goca feature user --fields "name:string,email:string" --database postgres --validation
goca feature product --fields "name:string,price:float64" --database postgres --validation
goca feature order --fields "user_id:int,total:float64" --database postgres --validation

# After: Configuration loaded automatically
goca feature user --fields "name:string,email:string"
goca feature product --fields "name:string,price:float64" 
goca feature order --fields "user_id:int,total:float64"
```

## 📋 Complete Configuration Reference

### Project Section

```yaml
project:
  name: "ecommerce-api"                    # Project name
  module: "github.com/company/ecommerce"   # Go module name
  description: "E-commerce REST API"       # Project description
  version: "1.0.0"                         # Semantic version
  author: "Development Team"               # Author/Team name
  license: "MIT"                           # License type
  tags: ["api", "ecommerce", "go"]        # Project tags
```

### Architecture Section

```yaml
architecture:
  pattern: "clean_architecture"            # Architecture pattern
  layers:
    domain:
      enabled: true                        # Enable domain layer
      path: "internal/domain"              # Domain layer path
      pattern: ""                          # Specific pattern for layer
      testing: true                        # Generate tests for layer
    usecase:
      enabled: true
      path: "internal/usecase"
      pattern: ""
      testing: true
    repository:
      enabled: true  
      path: "internal/repository"
      pattern: ""
      testing: true
    handler:
      enabled: true
      path: "internal/handler"
      pattern: ""
      testing: true
  patterns: ["repository", "factory"]      # Design patterns to use
  di:                                      # Dependency Injection config
    type: "manual"                         # DI type: manual, wire, dig
    auto_wire: true                        # Auto-wire dependencies
    container_path: "internal/di"          # DI container path
    interfaces: true                       # Generate interfaces
  naming:                                  # Naming conventions
    entities: "PascalCase"                 # Entity naming
    fields: "PascalCase"                   # Field naming
    files: "snake_case"                    # File naming
    packages: "lowercase"                  # Package naming
    constants: "UPPER_CASE"                # Constant naming
    variables: "camelCase"                 # Variable naming
    functions: "PascalCase"                # Function naming
```

### Database Section

```yaml
database:
  type: "postgres"                         # Database type
  host: "localhost"                        # Database host
  port: 5432                               # Database port
  name: "ecommerce_db"                     # Database name
  migrations:
    enabled: true                          # Enable migrations
    auto_generate: true                    # Auto-generate migrations
    directory: "migrations"                # Migration directory
    naming: "timestamp"                    # Naming convention
    versioning: "sequential"               # Version strategy
    tools: ["migrate", "sql-migrate"]      # Migration tools
  connection:
    max_open: 25                           # Max open connections
    max_idle: 5                            # Max idle connections
    max_lifetime: "5m"                     # Connection max lifetime
    ssl_mode: "disable"                    # SSL mode
    timezone: "UTC"                        # Database timezone
  features:
    soft_delete: true                      # Enable soft delete
    timestamps: true                       # Enable timestamps
    uuid: true                             # Use UUID for IDs
    audit: false                           # Enable audit trails
    versioning: false                      # Enable record versioning
    partitioning: false                    # Enable table partitioning
```

### Generation Section

```yaml
generation:
  validation:
    enabled: true                          # Enable validation
    tags: ["required", "min", "max"]       # Validation tags to use
    custom_rules: true                     # Allow custom rules
    error_handling: "detailed"             # Error handling style
    localization: "english"                # Language for messages
    strict_mode: false                     # Strict validation mode
  business_rules:
    enabled: true                          # Enable business rules
    directory: "internal/domain/rules"     # Business rules directory
    naming: "rule"                         # Naming convention
    testing: true                          # Generate tests
    documentation: true                    # Generate documentation
  docker:
    enabled: true                          # Generate Docker files
    compose: true                          # Generate docker-compose
    dockerfile: true                       # Generate Dockerfile
    multi_stage: true                      # Use multi-stage builds
    base_image: "alpine"                   # Base Docker image
  docs:
    swagger:
      enabled: true                        # Generate Swagger docs
      title: "Ecommerce API"               # API title
      version: "1.0.0"                     # API version
      description: "REST API for ecommerce" # API description
      host: "localhost:8080"               # API host
      base_path: "/api/v1"                 # API base path
      schemes: ["http", "https"]           # Supported schemes
    postman:
      enabled: true                        # Generate Postman collection
      output: "docs/postman"               # Output directory
      environment: true                    # Generate environment
      tests: true                          # Include tests
      variables: true                      # Include variables
    markdown:
      enabled: true                        # Generate Markdown docs
      output: "docs"                       # Output directory
      template: "default"                  # Documentation template
      toc: true                           # Generate table of contents
      examples: true                       # Include examples
  comments:
    enabled: true                          # Generate code comments
    language: "english"                    # Comment language
    style: "godoc"                         # Comment style
    examples: true                         # Include examples
    todo: true                            # Include TODO comments
    deprecated: true                       # Mark deprecated items
```

### Testing Section

```yaml
testing:
  enabled: true                            # Enable testing
  coverage_threshold: 80.0                 # Coverage threshold %
  benchmark: true                          # Generate benchmarks
  integration: true                        # Generate integration tests
  e2e: false                              # Generate E2E tests
  mocks:
    enabled: true                          # Generate mocks
    tool: "testify"                        # Mock generation tool
    directory: "mocks"                     # Mock directory
    suffix: "_mock"                        # Mock file suffix
    interfaces: []                         # Interfaces to mock
  fixtures:
    enabled: true                          # Generate test fixtures
    directory: "fixtures"                  # Fixture directory
    format: "json"                         # Fixture format
    seeds: true                           # Generate seed data
    factories: []                          # Factory patterns
  parallel: true                          # Enable parallel testing
  timeout: "5m"                           # Test timeout
```

### Features Section

```yaml
features:
  handlers: ["http", "grpc"]               # Handler types
  middleware: ["cors", "auth", "logging"]  # Middleware to include
  authentication:
    enabled: true                          # Enable authentication
    type: "jwt"                           # Auth type
    providers: ["local", "oauth"]         # Auth providers
  authorization:
    enabled: true                          # Enable authorization
    type: "rbac"                          # Authorization type
    roles: ["admin", "user"]              # Default roles
  validation: true                         # Enable validation
  soft_delete: true                        # Enable soft delete
  timestamps: true                         # Enable timestamps
  business_rules: false                    # Enable business rules
  caching:
    enabled: true                          # Enable caching
    type: "redis"                         # Cache type
    ttl: "1h"                            # Default TTL
  rate_limiting:
    enabled: true                          # Enable rate limiting
    requests_per_minute: 100              # Rate limit
  cors: true                              # Enable CORS
  compression: true                       # Enable compression
  metrics: true                           # Enable metrics
  tracing: false                          # Enable tracing
  logging:
    enabled: true                          # Enable logging
    level: "info"                         # Log level
    format: "json"                        # Log format
```

### Templates Section

```yaml
templates:
  enabled: true                           # Enable custom templates
  directory: ".goca/templates"            # Templates directory
  custom_templates:                       # Custom template definitions
    entity_custom:
      path: "entity.tmpl"
      type: "entity"
      variables:
        author: "Team"
  overrides:                             # Override built-in templates
    entity:
      path: "custom_entity.tmpl"
```

### Quality Section

```yaml
quality:
  style:
    gofmt: true                          # Use gofmt
    goimports: true                      # Use goimports
    golint: true                         # Use golint
    govet: true                          # Use go vet
    staticcheck: true                    # Use staticcheck
    custom: ["golangci-lint"]            # Custom linters
    line_length: 120                     # Max line length
    tab_width: 4                         # Tab width
  security:
    enabled: true                        # Enable security checks
    scanner: "gosec"                     # Security scanner
    rules: ["all"]                       # Security rules
    exclude: []                          # Rules to exclude
```

### Infrastructure Section

```yaml
infrastructure:
  logging:
    enabled: true                        # Enable logging
    level: "info"                        # Log level
    format: "structured"                 # Log format
    output: ["stdout", "file"]           # Log outputs
    structured: true                     # Structured logging
    tracing: false                      # Enable trace logging
  monitoring:
    enabled: true                        # Enable monitoring
    metrics: true                        # Enable metrics
    tracing: true                        # Enable distributed tracing
    health_check: true                   # Enable health checks
    profiling: true                      # Enable profiling
    tools: ["prometheus", "jaeger"]      # Monitoring tools
  cache:
    enabled: true                        # Enable caching
    type: "redis"                        # Cache type
    ttl: "1h"                           # Default TTL
    max_size: "100MB"                    # Max cache size
  message_queue:
    enabled: true                        # Enable message queue
    type: "rabbitmq"                     # Queue type
    exchanges: ["events"]                # Exchanges to create
    queues: ["notifications"]            # Queues to create
  deployment:
    type: "kubernetes"                   # Deployment type
    registry: "docker.io/company"        # Docker registry
    namespace: "production"              # K8s namespace
    replicas: 3                          # Replica count
    resources:
      cpu: "500m"                        # CPU limit
      memory: "512Mi"                    # Memory limit
```

## 🎯 Usage Scenarios

### Scenario 1: Simple API Project

```yaml
project:
  name: "simple-api"
  module: "github.com/company/simple-api"
  description: "Simple REST API"

database:
  type: "postgres"

generation:
  validation:
    enabled: true
  docker:
    enabled: true

features:
  handlers: ["http"]
  cors: true
  logging: true
```

**Commands:**
```bash
goca init simple-api --config
goca feature user --fields "name:string,email:string"
goca feature product --fields "name:string,price:float64"
```

### Scenario 2: Microservice with gRPC

```yaml
project:
  name: "user-service"
  module: "github.com/company/user-service"
  description: "User management microservice"

database:
  type: "postgres"
  features:
    soft_delete: true
    timestamps: true

generation:
  validation:
    enabled: true
  docs:
    swagger:
      enabled: true

features:
  handlers: ["http", "grpc"]
  authentication:
    enabled: true
    type: "jwt"
  caching:
    enabled: true
    type: "redis"

infrastructure:
  monitoring:
    enabled: true
    metrics: true
  logging:
    enabled: true
    level: "debug"
```

### Scenario 3: Enterprise Application

```yaml
project:
  name: "enterprise-app"
  module: "github.com/company/enterprise-app"
  description: "Full-featured enterprise application"

database:
  type: "postgres"
  features:
    soft_delete: true
    timestamps: true
    uuid: true
    audit: true

generation:
  validation:
    enabled: true
    strict_mode: true
  business_rules:
    enabled: true
  docs:
    swagger:
      enabled: true
    postman:
      enabled: true
    markdown:
      enabled: true

testing:
  enabled: true
  coverage_threshold: 90.0
  integration: true
  e2e: true

features:
  handlers: ["http", "grpc", "graphql"]
  authentication:
    enabled: true
  authorization:
    enabled: true
  rate_limiting:
    enabled: true
  caching:
    enabled: true

infrastructure:
  monitoring:
    enabled: true
    metrics: true
    tracing: true
  logging:
    enabled: true
    level: "info"
  message_queue:
    enabled: true
  deployment:
    type: "kubernetes"
```

## 🚦 Configuration Priority

GOCA follows this priority order when loading configuration:

1. **CLI Flags** (highest priority)
2. **Environment Variables** 
3. **`.goca.yaml` File**
4. **Default Values** (lowest priority)

### Example:

```bash
# .goca.yaml specifies database: "postgres"
# CLI flag overrides it
goca feature user --fields "name:string" --database mysql

# Result: Uses MySQL despite .goca.yaml specifying Postgres
```

## 🔧 Advanced Configuration

### Environment-Specific Configuration

```yaml
# .goca.yaml (base)
project:
  name: "my-app"

database:
  type: "postgres"
  host: "localhost"

# .goca.production.yaml
database:
  host: "prod-db.company.com"
  connection:
    max_open: 100
```

```bash
# Load production config
GOCA_ENV=production goca feature user --fields "name:string"
```

### Template Variables

```yaml
templates:
  variables:
    author: "Development Team"
    company: "ACME Corp"
    year: "2024"
```

**In templates:**
```go
// Code generated by GOCA for {{.company}}
// Author: {{.author}}
// Year: {{.year}}
```

### Conditional Configuration

```yaml
generation:
  validation:
    enabled: true
    rules:
      - condition: "field.Type == 'email'"
        tags: ["email", "required"]
      - condition: "field.Type == 'password'"
        tags: ["min=8", "required"]
```

## 🐛 Troubleshooting

### Common Issues

#### 1. Configuration Not Loading
```bash
❌ Warning: Could not load configuration
```

**Solutions:**
- Ensure `.goca.yaml` exists in project root
- Check YAML syntax with `goca config validate`
- Verify file permissions

#### 2. Invalid YAML Structure
```bash
❌ Invalid YAML file: yaml: unmarshal errors
```

**Solutions:**
- Use `goca config validate` to see specific errors
- Check indentation (use spaces, not tabs)
- Validate YAML syntax online

#### 3. Missing Required Sections
```bash
❌ Required section 'project' missing
```

**Solution:**
```bash
goca config init --force  # Regenerate config
```

## 🎨 Best Practices

### 1. Project Organization
- Keep `.goca.yaml` in project root
- Use consistent naming across projects
- Document custom configurations

### 2. Team Workflow
```bash
# 1. Team lead creates project configuration
goca config init --template enterprise

# 2. Commit configuration to repository
git add .goca.yaml
git commit -m "Add GOCA configuration"

# 3. Team members use shared configuration
git clone project-repo
cd project-repo
goca feature new-feature --fields "data:string"
```

### 3. Configuration Management
- Use templates for similar projects
- Override specific values per environment
- Keep sensitive data in environment variables

### 4. Validation
```bash
# Always validate before sharing
goca config validate

# Test configuration works
goca feature test-feature --fields "name:string"
```

## 🔄 Migration from CLI-Only

### Before (CLI-Only)
```bash
goca init ecommerce-api --database postgres
goca feature user --fields "name:string,email:string" --database postgres --validation
goca feature product --fields "name:string,price:float64" --database postgres --validation
```

### After (Configuration-Based)
```bash
# 1. Create configuration
goca config init --template api --database postgres

# 2. Generate features (options loaded from config)
goca feature user --fields "name:string,email:string"
goca feature product --fields "name:string,price:float64"
```

### Benefits Achieved
- **50% less typing** in commands
- **100% consistency** across features
- **Team-wide standardization**
- **Version-controlled configuration**

## 📚 Next Steps

1. **Initialize Configuration**: `goca config init --template api`
2. **Validate Setup**: `goca config validate`
3. **Generate Features**: `goca feature user --fields "name:string,email:string"`
4. **Customize Templates**: Edit `.goca/templates/` files
5. **Share with Team**: Commit `.goca.yaml` to repository

The YAML configuration system makes GOCA much more powerful and easier to use! 🚀