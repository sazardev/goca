# Configuration Guide

## Overview

GOCA supports centralized project configuration through the `.goca.yaml` file. This powerful feature allows you to define project-wide settings, conventions, and preferences that will be applied automatically to all code generation commands.

### Benefits

- **Consistency**: Maintain uniform settings across your entire project
- **Productivity**: Avoid repeating the same CLI flags for every command
- **Team Collaboration**: Share standardized configurations across your team
- **Flexibility**: Override configuration with CLI flags when needed
- **Documentation**: Configuration files serve as self-documenting project preferences

### When to Use Configuration

Use `.goca.yaml` when:
- Working on projects with consistent patterns and conventions
- Managing multiple features with the same database type
- Enforcing team-wide coding standards
- Automating CI/CD workflows
- Customizing code generation templates

## Configuration File Location

Create a `.goca.yaml` file in your project root directory (where your `go.mod` file is located):

```
my-project/
├── .goca.yaml          ← Configuration file
├── go.mod
├── go.sum
├── cmd/
├── internal/
└── ...
```

## Configuration Precedence

GOCA uses a three-tier configuration system with the following precedence (highest to lowest):

1. **CLI Flags**: Command-line arguments (highest priority)
2. **Configuration File**: Settings from `.goca.yaml`
3. **Defaults**: Built-in GOCA defaults (lowest priority)

This allows you to define common settings in your configuration file while still being able to override them for specific commands.

## Core Configuration Sections

### Project Configuration

Define basic project metadata and information:

```yaml
project:
  name: my-api
  module: github.com/mycompany/my-api
  description: RESTful API for customer management
  version: 1.0.0
  author: Development Team
  license: MIT
```

**Fields:**
- `name`: Project name
- `module`: Go module path
- `description`: Project description
- `version`: Project version
- `author`: Author or team name
- `license`: License type

### Database Configuration

Configure database settings and features:

```yaml
database:
  type: postgres
  host: localhost
  port: 5432
  migrations:
    enabled: true
    auto_generate: true
    directory: migrations
  features:
    soft_delete: true
    timestamps: true
    uuid: true
    audit: false
```

**Supported database types:**
- `postgres`: PostgreSQL
- `mysql`: MySQL/MariaDB
- `mongodb`: MongoDB

**Migration settings:**
- `enabled`: Enable/disable migrations
- `auto_generate`: Auto-generate migration files
- `directory`: Migration files directory

**Database features:**
- `soft_delete`: Add soft delete functionality to entities
- `timestamps`: Add created_at/updated_at fields
- `uuid`: Use UUID for primary keys
- `audit`: Enable audit logging

### Architecture Configuration

Define Clean Architecture layers and naming conventions:

```yaml
architecture:
  layers:
    domain:
      enabled: true
      directory: internal/domain
    usecase:
      enabled: true
      directory: internal/usecase
    repository:
      enabled: true
      directory: internal/repository
    handler:
      enabled: true
      directory: internal/handler
  
  patterns:
    - repository
    - service
    - dto
  
  naming:
    files: lowercase
    entities: PascalCase
    variables: camelCase
    functions: PascalCase
    constants: SCREAMING_SNAKE
```

**Layer configuration:**
- `enabled`: Enable/disable layer generation
- `directory`: Custom directory for layer

**Naming conventions:**
- `lowercase`: user_service.go
- `snake_case`: user_service.go
- `PascalCase`: UserService
- `camelCase`: userService
- `SCREAMING_SNAKE`: MAX_RETRIES

### Generation Configuration

Control code generation preferences:

```yaml
generation:
  validation:
    enabled: true
    library: builtin
    sanitize: true
  
  business_rules:
    enabled: true
    patterns:
      - validation
      - authorization
    events: true
  
  documentation:
    swagger:
      enabled: true
      version: "2.0"
      output: docs/swagger.yaml
    comments:
      enabled: true
      language: english
      style: godoc
```

**Validation options:**
- `enabled`: Enable field validation
- `library`: Validation library (builtin, validator, ozzo-validation)
- `sanitize`: Enable input sanitization

**Business rules:**
- `enabled`: Generate business rules layer
- `patterns`: Patterns to apply
- `events`: Enable domain events

### Testing Configuration

Configure testing generation preferences:

```yaml
testing:
  enabled: true
  framework: testify
  coverage:
    enabled: true
    threshold: 80
  mocks:
    enabled: true
    tool: testify
    directory: internal/mocks
  integration: true
  benchmarks: true
```

**Testing options:**
- `enabled`: Enable test generation
- `framework`: Testing framework (testify, ginkgo, builtin)
- `coverage`: Code coverage settings
- `mocks`: Mock generation settings
- `integration`: Generate integration tests
- `benchmarks`: Generate benchmark tests

### Template Configuration

Customize code generation templates:

```yaml
templates:
  directory: .goca/templates
  variables:
    author: "Development Team"
    copyright: "2024 MyCompany Inc"
    license: "MIT"
  custom:
    entity:
      path: .goca/templates/entity.tmpl
      type: go-template
```

**Template settings:**
- `directory`: Custom templates directory
- `variables`: Template variables
- `custom`: Custom template definitions

## Configuration Examples

### Minimal Configuration

Simple configuration for small projects:

```yaml
project:
  name: minimal-api
  module: github.com/user/minimal-api
```

### Web Application

Configuration for a web application with HTTP handlers:

```yaml
project:
  name: web-app
  module: github.com/user/web-app

generation:
  validation:
    enabled: true
```

### Microservice

Configuration for a microservice with database and testing:

```yaml
project:
  name: user-service
  module: github.com/company/user-service

database:
  type: postgres
  migrations:
    enabled: true

architecture:
  patterns:
    - repository
    - service

testing:
  enabled: true
  framework: testify
  mocks:
    enabled: true
```

### Enterprise Application

Comprehensive configuration for enterprise applications:

```yaml
project:
  name: enterprise-api
  module: github.com/corp/enterprise-api
  description: Enterprise-grade API

database:
  type: postgres
  migrations:
    enabled: true
  features:
    soft_delete: true
    timestamps: true

generation:
  validation:
    enabled: true
  business_rules:
    enabled: true

testing:
  enabled: true
  framework: testify
  coverage:
    enabled: true
  mocks:
    enabled: true
  integration: true

templates:
  directory: .goca/templates

architecture:
  naming:
    files: lowercase
    entities: PascalCase
    variables: camelCase
```

## Using Configuration with Commands

### Initialize Project with Configuration

When you have a `.goca.yaml` file, GOCA will automatically use it:

```bash
# Configuration will be loaded automatically
goca feature Product --fields "name:string,price:float64,stock:int"
```

### Override Configuration

You can override configuration settings using CLI flags:

```bash
# Override database type from config
goca feature Order --database mysql
```

### View Effective Configuration

Check what configuration is being used:

```bash
goca config show
```

## Template Customization

### Using Custom Templates

1. Create a templates directory:

```bash
mkdir -p .goca/templates
```

2. Configure template directory:

```yaml
templates:
  directory: .goca/templates
```

3. Create custom templates in the directory following GOCA's template structure

4. GOCA will use your custom templates instead of built-in ones

### Template Variables

Define custom variables for templates:

```yaml
templates:
  variables:
    author: "Your Name"
    company: "Your Company"
    copyright: "2024"
    license: "Apache 2.0"
```

Access variables in templates:

```go
// Generated by GOCA CLI
// Author: {{.Author}}
// Copyright: {{.Copyright}} {{.Company}}
```

## Best Practices

### 1. Start Simple

Begin with minimal configuration and add settings as needed:

```yaml
project:
  name: my-project
  module: github.com/user/my-project

database:
  type: postgres
```

### 2. Document Your Choices

Add comments to explain configuration decisions:

```yaml
# Using PostgreSQL for advanced features
database:
  type: postgres
  
  # Enable soft deletes for audit trail
  features:
    soft_delete: true
    timestamps: true
```

### 3. Version Control

Always commit `.goca.yaml` to version control:

```bash
git add .goca.yaml
git commit -m "Add GOCA configuration"
```

### 4. Team Standards

Use configuration to enforce team-wide standards:

```yaml
architecture:
  naming:
    files: lowercase      # Consistent file naming
    entities: PascalCase  # Go naming conventions
    variables: camelCase  # Go naming conventions

testing:
  enabled: true           # Always generate tests
  framework: testify      # Standard testing framework
```

### 5. Environment Separation

For environment-specific settings, use separate configuration files:

```
.goca.yaml           # Base configuration
.goca.dev.yaml       # Development overrides
.goca.prod.yaml      # Production overrides
```

### 6. Validate Configuration

Test your configuration before committing:

```bash
# Generate a test feature to verify settings
goca feature TestEntity --fields "name:string"

# Review generated code
# If correct, delete test entity and commit config
```

## Configuration Validation

GOCA validates your configuration file when loading. Common validation errors:

### Invalid YAML Syntax

```yaml
# ERROR: Invalid indentation
project:
name: my-project  # Should be indented
```

### Required Fields Missing

```yaml
# ERROR: Project section required
database:
  type: postgres
```

**Fix:** Add required project information:

```yaml
project:
  name: my-project
  module: github.com/user/my-project

database:
  type: postgres
```

### Invalid Values

```yaml
# ERROR: Unsupported database type
database:
  type: oracle  # Not supported
```

**Fix:** Use supported values:

```yaml
database:
  type: postgres  # postgres, mysql, mongodb
```

## Troubleshooting

### Configuration Not Loading

**Problem:** GOCA doesn't seem to use your configuration file.

**Solution:**
1. Verify file name is exactly `.goca.yaml`
2. Check file is in project root (where go.mod is)
3. Verify YAML syntax is valid
4. Check for validation errors in output

### CLI Flags Not Overriding

**Problem:** CLI flags don't override configuration settings.

**Solution:**
- CLI flags have highest precedence and should always override
- Verify you're using the correct flag name
- Check command output for applied settings

### Template Customization Not Working

**Problem:** Custom templates are not being used.

**Solution:**
1. Verify template directory path in configuration
2. Check template files follow GOCA's naming conventions
3. Ensure template syntax is valid Go templates

## Advanced Configuration

### Conditional Generation

Control what gets generated based on project type:

```yaml
# API-only project
architecture:
  layers:
    domain:
      enabled: true
    usecase:
      enabled: true
    repository:
      enabled: true
    handler:
      enabled: true  # Enable HTTP handlers
```

### Custom Directory Structure

Override default directory layout:

```yaml
architecture:
  layers:
    domain:
      enabled: true
      directory: pkg/domain      # Custom path
    repository:
      enabled: true
      directory: pkg/persistence # Custom path
```

### Multiple Pattern Support

Apply multiple architectural patterns:

```yaml
architecture:
  patterns:
    - repository    # Repository pattern
    - service       # Service layer
    - dto           # Data Transfer Objects
    - specification # Specification pattern
```

## Migration Guide

### From CLI-Only to Configuration

If you're currently using only CLI flags, migrate to configuration:

1. **Identify repeated flags:**

```bash
# You run this often:
goca feature User --database postgres --validation
goca feature Order --database postgres --validation
```

2. **Create configuration:**

```yaml
project:
  name: my-project
  module: github.com/user/my-project

database:
  type: postgres

generation:
  validation:
    enabled: true
```

3. **Simplify commands:**

```bash
# Now you can run:
goca feature User
goca feature Order
```

## Next Steps

- Learn about [Clean Architecture](./clean-architecture.md)
- Explore [Commands Reference](../commands/)
- Read [Best Practices](./best-practices.md)
- Check [Project Structure](./project-structure.md)

## Resources

- [YAML Specification](https://yaml.org/)
- [Go Templates](https://pkg.go.dev/text/template)
- [Clean Architecture Principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
