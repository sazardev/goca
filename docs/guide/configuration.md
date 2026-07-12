# Configuration Guide

## Overview

GOCA supports centralized project configuration through the `.goca.yaml` file. This allows you to define project-wide settings, conventions, and preferences that are applied automatically to code generation commands.

### Benefits

- **Consistency**: Maintain uniform settings across your entire project
- **Productivity**: Avoid repeating the same CLI flags for every command
- **Standardization**: Share a configuration across your team
- **Flexibility**: Override configuration with CLI flags when needed
- **Documentation**: Configuration files serve as self-documenting project preferences

### When to Use Configuration

Use `.goca.yaml` when:
- Working on projects with consistent patterns and conventions
- Managing multiple features with the same database type
- Enforcing team-wide naming conventions
- Customizing code generation templates

### Quick Start

The `.goca.yaml` file is generated automatically when you scaffold a project:

```bash
goca init my-api --module github.com/user/my-api
```

`goca init` writes `.goca.yaml` unless you pass `--config=false`. You can also regenerate it later, or inspect/validate an existing one:

```bash
goca config show      # print the current .goca.yaml
goca config validate  # validate it (usable as a CI gate)
```

::: warning `goca config init` writes a different, incompatible file
`goca config init` is a separate, older command that writes its own hand-authored template with sections (`quality`, `infrastructure`, `defaults`) that **do not match** the schema described below and are silently ignored by every other Goca command. Prefer `goca init --config` (the default) to generate `.goca.yaml` — this guide documents that schema. Treat `goca config init` as legacy/experimental until this duplication is resolved.
:::

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

The sections below reflect the real `.goca.yaml` schema (the `GocaConfig` struct). Not every field influences code generation yet — each section notes what is actually wired in versus parsed-and-reserved for future use.

### Project Configuration

```yaml
project:
  name: my-api
  module: github.com/mycompany/my-api
  description: RESTful API for customer management
  version: 1.0.0
  author: Development Team
  license: MIT
```

`name` and `module` are used as generation inputs; the rest is descriptive metadata.

### Database Configuration

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

**Supported database types (`database.type`):**
- `postgres`: PostgreSQL
- `postgres-json`: PostgreSQL with JSONB fields
- `mysql`: MySQL/MariaDB
- `sqlite`: SQLite (no server required — good default for local development)
- `sqlserver`: Microsoft SQL Server
- `mongodb`: MongoDB
- `dynamodb`: AWS DynamoDB
- `elasticsearch`: Elasticsearch

`mysql` and `sqlite` reuse the same GORM-based repository implementation as `postgres` — there is one generated file per database "family", not one per driver.

**Actually applied to generation:**
- `type` drives which repository/DI implementation is generated
- `features.soft_delete` and `features.timestamps` drive entity field generation

**Parsed but not yet consumed by generators:** `migrations`, `connection`, and the remaining `features` flags (`uuid`, `audit`, `versioning`, `partitioning`, `indexes`, `constraints`).

### Architecture Configuration

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

  naming:
    entities: PascalCase
    fields: camelCase
    files: lowercase
    packages: lowercase
    constants: SCREAMING_SNAKE
    variables: camelCase
    functions: PascalCase
```

**Layer configuration:**
- `enabled`: enable/disable that layer (only `layers.handler.enabled` is actually read — it controls whether HTTP handler code is generated)
- `directory`: custom directory for the layer (the key is `directory`, not `path`)

**Naming conventions** (all actively applied via `GetNamingConvention()`):
- `lowercase` / `snake_case`: `user_service.go`
- `PascalCase`: `UserService`
- `camelCase`: `userService`
- `SCREAMING_SNAKE`: `MAX_RETRIES`

There is also an `architecture.di` block (`type: manual|wire|fx|dig`, `auto_wire`, `providers`, `modules`) — the type is validated but not yet wired into the dependency-injection generator, which always emits manual wiring.

### Generation Configuration

```yaml
generation:
  validation:
    enabled: true
    library: builtin
    sanitize: true

  business_rules:
    enabled: true
    events: true
```

**Actually applied:** `validation.enabled` and `business_rules.enabled` control whether `goca feature` generates validation code and a business-rules layer. `validation.library` accepts `builtin`, `validator`, or `ozzo-validation`.

The section also supports `documentation` (swagger/postman/markdown/comment settings) and `style`/`imports` (formatting preferences) — these are parsed and validated but not yet consumed by any generator.

### Testing Configuration

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

`framework` accepts `testify`, `ginkgo`, or `builtin`; `coverage.threshold` must be 0-100. Note: this section is currently only read by `goca upgrade`'s diagnostics — no generator branches on it yet (test scaffolding commands like `goca feature --testing` are driven by CLI flags, not this config block).

### Template Configuration

```yaml
templates:
  directory: .goca/templates
  variables:
    author: "Development Team"
    copyright: "2024 MyCompany Inc"
    license: "MIT"
```

This section is genuinely functional: `directory` points Goca at your `.goca/templates` folder, and any `.tmpl` file placed there (e.g. `entity.tmpl`, `usecase.tmpl`) overrides the corresponding built-in template. See [Template Customization](#template-customization) below.

### Features Configuration

```yaml
features:
  auth:
    enabled: true
    type: jwt
  cache:
    enabled: true
    type: redis
```

`auth.type` accepts `jwt`, `oauth2`, `session`, or `basic`; `cache.type` accepts `redis`, `memcached`, or `inmemory`. Beyond validation, this section is currently only exposed as data to custom templates and to `goca upgrade` diagnostics — enabling `features.auth` does not by itself generate authentication middleware (use `goca middleware` for that).

### Deploy Configuration

```yaml
deploy:
  docker:
    enabled: true
  kubernetes:
    enabled: false
  ci:
    provider: github-actions
```

Parsed and validated, but not yet applied to generation — Dockerfile/Kubernetes/CI file generation is currently driven by CLI flags and `--database`, not by this block. Treat it as reserved for a future release.

## Configuration Examples

### Minimal Configuration

```yaml
project:
  name: minimal-api
  module: github.com/user/minimal-api
```

### Web Application

```yaml
project:
  name: web-app
  module: github.com/user/web-app

generation:
  validation:
    enabled: true
```

### Microservice

```yaml
project:
  name: user-service
  module: github.com/company/user-service

database:
  type: postgres
  migrations:
    enabled: true

testing:
  enabled: true
  framework: testify
  mocks:
    enabled: true
```

### Full Configuration

```yaml
project:
  name: enterprise-api
  module: github.com/corp/enterprise-api
  description: Enterprise-grade API

database:
  type: postgres
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

### Generation Commands Load It Automatically

```bash
# Configuration will be loaded automatically
goca feature Product --fields "name:string,price:float64,stock:int"
```

### Override Configuration

```bash
# Override database type from config
goca feature Order --database mysql
```

### View and Validate Configuration

```bash
goca config show      # print the current .goca.yaml
goca config validate  # validate types, naming, thresholds, etc.
```

## Template Customization

### Using Custom Templates

1. Create a templates directory:

```bash
mkdir -p .goca/templates
```

2. Configure the template directory:

```yaml
templates:
  directory: .goca/templates
```

3. Add a template file named after the generator it should override, e.g. `.goca/templates/entity.tmpl`, `.goca/templates/usecase.tmpl`, `.goca/templates/handler.tmpl`, or `.goca/templates/repo.tmpl`, written as a Go `text/template`.

4. Goca uses your custom template instead of the built-in one for any generator whose `.tmpl` file is present; generators without a matching file fall back to the built-in template.

### Template Variables

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

### 1. Generate Configuration with `goca init`

Let `goca init` create your first `.goca.yaml` rather than hand-writing one — it fills in valid defaults for every section:

```bash
goca init my-api --module github.com/user/my-api
```

### 2. Start Simple

Begin with minimal configuration and add settings as needed:

```yaml
project:
  name: my-project
  module: github.com/user/my-project

database:
  type: postgres
```

### 3. Document Your Choices

Add comments to explain configuration decisions, especially if you've customized a template:

```yaml
# Using PostgreSQL for advanced features
database:
  type: postgres

  # Enable soft deletes for audit trail
  features:
    soft_delete: true
    timestamps: true
```

### 4. Version Control

Always commit `.goca.yaml` to version control so team members can use the same configuration:

```bash
git add .goca.yaml
git commit -m "Add GOCA configuration"
```

### 5. Team Standards

Use configuration to enforce team-wide naming standards:

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

### 6. Validate Before Committing

```bash
goca config validate
```

## Configuration Validation

`goca config validate` checks your configuration file when loading. Common validation errors:

### Invalid YAML Syntax

```yaml
# ERROR: Invalid indentation
project:
name: my-project  # Should be indented
```

### Invalid Values

```yaml
# ERROR: Unsupported database type
database:
  type: oracle  # Not supported
```

**Fix:** Use one of the supported values:

```yaml
database:
  type: postgres  # postgres, postgres-json, mysql, sqlite, sqlserver, mongodb, dynamodb, elasticsearch
```

## Troubleshooting

### Configuration Not Loading

**Problem:** GOCA doesn't seem to use your configuration file.

**Solution:**
1. Verify the file name is exactly `.goca.yaml`
2. Check the file is in the project root (where `go.mod` is)
3. Run `goca config validate` to check for syntax/value errors
4. Make sure you didn't generate it with `goca config init` instead of `goca init --config` (see the warning above — the two produce incompatible schemas)

### CLI Flags Not Overriding

**Problem:** CLI flags don't override configuration settings.

**Solution:**
- CLI flags have the highest precedence and should always override
- Verify you're using the correct flag name
- Check command output for the effective settings actually applied

### Template Customization Not Working

**Problem:** Custom templates are not being used.

**Solution:**
1. Verify `templates.directory` in the configuration
2. Check the template file name matches the generator (`entity.tmpl`, `usecase.tmpl`, `handler.tmpl`, `repo.tmpl`)
3. Ensure the template is valid Go `text/template` syntax

## Next Steps

- Learn about [Clean Architecture](./clean-architecture.md)
- Explore [Commands Reference](../commands/)
- Read [Best Practices](./best-practices.md)
- Check [Project Structure](./project-structure.md)

## Resources

- [YAML Specification](https://yaml.org/)
- [Go Templates](https://pkg.go.dev/text/template)
- [Clean Architecture Principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
