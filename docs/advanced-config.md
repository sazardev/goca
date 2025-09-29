# Advanced GOCA Configuration Commands

## Introduction

GOCA includes a complete configuration management system through the `goca config` command with various subcommands that facilitate working with `.goca.yaml` files.

## Available Subcommands

### `goca config show`

Displays the current project configuration in a detailed and readable manner.

```bash
# Show current configuration
goca config show
```

**Typical output:**
```
=== Current GOCA Configuration ===
✅ File found: .goca.yaml
📁 Directory: /path/to/project

--- Content ---
project:
  name: "my-project"
  module: "github.com/user/my-project"
  version: "1.0.0"
  description: "Project generated with GOCA CLI"

defaults:
  database: postgres
  handlers: [http, grpc]
  validation: true
  # ... rest of configuration

--- Validation ---
✅ Valid structure
```

**Use cases:**
- Verify current configuration before generating code
- Debug configuration issues
- Document project setup

### `goca config init`

Initializes a new configuration file with predefined templates.

```bash
# Basic configuration
goca config init

# With specific template
goca config init --template web --database postgres --handlers http,grpc

# Overwrite existing file
goca config init --template api --force
```

**Available options:**

| Flag | Description | Values | Example |
|------|-------------|--------|---------|
| `--template` | Predefined template | `web`, `api`, `microservice`, `full` | `--template web` |
| `--database` | Database type | `postgres`, `mysql`, `sqlite` | `--database postgres` |
| `--handlers` | Handler types | `http`, `grpc`, `cli` | `--handlers http,grpc` |
| `--force` | Overwrite existing | `true`, `false` | `--force` |

### `goca config validate`

Validates the structure and content of the current configuration file.

```bash
# Validate configuration
goca config validate
```

**Types of validation:**
- ✅ **YAML syntax**: Valid file structure
- ✅ **Required fields**: All mandatory sections present
- ✅ **Data types**: Correct values for each field
- ⚠️ **Warnings**: Suboptimal but valid configurations

### `goca config template`

Shows information about available templates and their characteristics.

```bash
# List available templates
goca config template
```

**Output:**
```
📋 Available Configuration Templates
===================================
• web
  Complete web application with frontend and backend

• api
  REST API with database

• microservice
  Microservice with multiple handlers

• full
  Complete configuration with all features

• default
  Basic minimal configuration

Usage:
  goca config init --template <name>
  goca config init --template web --database postgres --handlers http,grpc
```

## Available Templates

### 1. Template "web"
- **Use**: Complete web applications
- **Includes**: Frontend + Backend + Database
- **Handlers**: HTTP + WebSocket
- **Features**: Authentication, sessions, static files

### 2. Template "api"
- **Use**: Pure REST APIs
- **Includes**: Endpoints + Database + Documentation
- **Handlers**: HTTP + JSON
- **Features**: Swagger, validation, middleware

### 3. Template "microservice"
- **Use**: Distributed services
- **Includes**: Multiple communication interfaces
- **Handlers**: HTTP + gRPC + Message Queue
- **Features**: Metrics, tracing, health checks

### 4. Template "full"
- **Use**: Complex enterprise projects
- **Includes**: All available features
- **Handlers**: HTTP + gRPC + CLI + WebSocket
- **Features**: Everything included with enterprise configuration

### 5. Template "default"
- **Use**: Simple or custom projects
- **Includes**: Minimal functional configuration
- **Handlers**: Basic HTTP
- **Features**: Essential configuration only

## Common Workflows

### Starting a new project

```bash
# 1. Create project directory
mkdir my-new-project
cd my-new-project

# 2. Initialize configuration
goca config init --template api --database postgres --handlers http

# 3. Verify configuration
goca config show

# 4. Initialize complete project
goca init --config

# 5. Generate first feature
goca feature User --fields "name:string,email:string,age:int"
```

### Migrating existing project

```bash
# 1. In existing project directory
cd my-existing-project

# 2. Create configuration based on current structure
goca config init --template default

# 3. Manually adjust configuration if needed
# (edit .goca.yaml)

# 4. Validate configuration
goca config validate

# 5. Continue with normal development
goca feature Product --fields "name:string,price:float64"
```

### Team development

```bash
# 1. Clone team repository
git clone https://github.com/company/project.git
cd project

# 2. Verify project configuration
goca config show

# 3. If no configuration exists, create team standard
goca config init --template full --database postgres --handlers http,grpc

# 4. Commit configuration to repository
git add .goca.yaml
git commit -m "Add GOCA configuration"

# 5. Team can use same configuration
goca feature NewFeature --fields "field:string"
```

## Integration with Existing Commands

### Order of precedence

GOCA commands follow this precedence for configuration:

1. **CLI Flags** (highest priority)
2. **.goca.yaml file** (if exists)
3. **Default values** (lowest priority)

### Integration examples

```bash
# Use configuration from file
goca feature Product --fields "name:string,price:float64"

# Override database from configuration
goca feature Product --fields "name:string,price:float64" --database mysql

# Override handlers from configuration  
goca feature Product --fields "name:string,price:float64" --handlers grpc
```

## Advanced Customization

### Template variables

Templates support variables that are automatically substituted:

- `{{.ProjectName}}`: Project name (current directory)
- `{{.ModuleName}}`: Inferred Go module name
- `{{.DatabaseType}}`: Specified database type
- `{{.Handlers}}`: List of specified handlers

### Custom templates

Although not implemented in this version, the structure allows for custom templates:

```yaml
# .goca.yaml
templates:
  custom:
    my-template:
      path: ".goca/templates/my-template.yaml"
      variables:
        custom_var: "value"
```

## Troubleshooting

### Error: "Configuration validation failed"

```bash
# Check specific errors
goca config validate

# If there are structure errors, regenerate configuration
goca config init --force
```

### Error: "Could not load configuration"

```bash
# Verify file exists
ls -la .goca.yaml

# Check YAML syntax
goca config validate

# If it doesn't exist, create a new one
goca config init --template default
```

### Warning: "Using default values"

This warning appears when:
- There's no `.goca.yaml` file in the project
- The file exists but has validation errors
- The file is not compatible with expected structures

**Solution:** Create or fix the configuration file:
```bash
goca config init --template default
```

## Tips and Best Practices

### 1. Use appropriate templates
- **Simple project**: `default`
- **REST API**: `api`  
- **Web application**: `web`
- **Microservice**: `microservice`
- **Enterprise project**: `full`

### 2. Commit configuration to repository
```bash
git add .goca.yaml
git commit -m "Add project GOCA configuration"
```

### 3. Validate before generating code
```bash
goca config validate
goca feature MyFeature --fields "name:string"
```

### 4. Document team configuration
Create a `README-GOCA.md` explaining:
- Chosen template and why
- Project-specific configuration
- Team workflows
- Naming and structure standards

## Quick Reference Commands

```bash
# View current configuration
goca config show

# Create new configuration
goca config init --template <type>

# Validate configuration
goca config validate

# View available templates
goca config template

# Complete example
goca config init --template api --database postgres --handlers http,grpc --force
```

This configuration system makes GOCA much more flexible and easy to use in real projects, especially in development teams where configuration consistency is crucial.