# Goca CLI — AI Agent Skill

## Description
Specialized skill for AI agents using **Goca** (Go Clean Architecture Code Generator). Provides complete command reference, workflow patterns, error tracing, and architecture rules based on the actual codebase.

**Activate when:** user mentions `goca`, Clean Architecture code generation, Go project scaffolding, or any goca subcommand.

---

## Clean Architecture — Layer Map

Dependencies point **inward** — each layer only knows the layer below it:

```
HTTP Handler → UseCase Interface → Repository Interface → Domain Entity
   adapter        business logic       persistence          pure logic
```

### Generated directories:

| Layer | Path | Depends On |
|-------|------|------------|
| Domain | `internal/domain/` | Nothing (pure Go) |
| UseCase | `internal/usecase/` | Domain + Repository interfaces |
| Repository | `internal/repository/` | Domain |
| Handler | `internal/handler/` | UseCase interfaces |
| DI | `internal/di/` | All layers (wires concretes) |
| Messages | `internal/messages/` | Nothing |
| Interfaces | `internal/interfaces/` | Nothing (contracts only) |
| Middleware | `internal/middleware/` | Nothing |

**Violations to prevent:**
- ❌ Handler imports concrete usecase types (must use interfaces)
- ❌ UseCase imports Handler or Repository concretes
- ❌ Domain imports anything outside `internal/domain/`

---

## Complete Command Reference

### Global flags (apply to all commands):
`--no-color`, `--no-interactive`, `-q` / `--quiet`, `-v` / `--verbose`

### Code Generation Commands

| Command | Args | Required Flags | When to Use |
|---------|------|----------------|-------------|
| `goca init <name>` | 0-1 | `--module` | **Scaffold new project** from scratch. Creates full Clean Architecture structure. |
| `goca feature <name>` | Exactly 1 | `--fields "field:type"` | **Generate ALL layers at once** (entity + usecase + handler + repository + messages + DI). Fastest path to a working feature. |
| `goca entity <name>` | Exactly 1 | `--fields "field:type"` | **Generate only the domain entity** with validation, business rules, seeds, tests. Use when you already have other layers or want fine-grained control. |
| `goca usecase <name>` | Exactly 1 | `--entity <name>` | **Generate business logic** (DTOs + service + interface). Requires existing entity. |
| `goca handler <entity>` | Exactly 1 | (none) | **Generate protocol adapter** (HTTP, gRPC, CLI, Worker, SOAP). Requires existing usecase. |
| `goca repository <entity>` | Exactly 1 | (none) | **Generate persistence layer** (interface + DB-specific impl). Choose from 8 databases. |
| `goca messages <entity>` | Exactly 1 | (none) | **Generate error + response constants** per feature. |
| `goca di` | None | `--features` | **Generate or update DI container** (manual or Google Wire). |
| `goca interfaces <entity>` | Exactly 1 | (none) | **Generate layer contracts only** — ideal for TDD approach. |
| `goca middleware <name>` | Exactly 1 | (none) | **Generate HTTP middleware** (cors, logging, auth, rate-limit, recovery, request-id, timeout). |

### Diagnostic & Utility Commands

| Command | Args | When to Use |
|---------|------|-------------|
| `goca doctor` | None | **Health check** — verifies go.mod, .goca.yaml, dirs, go build, go vet, DI container. Always run first when something is wrong. |
| `goca analyze` | None | **Deep self-analysis** — checks architecture, quality, security (OWASP), standards, tests, deps. Use after `doctor` passes. |
| `goca integrate` | None | **Auto-wire existing features** into DI + main.go. Use after manually creating entities/layers. |
| `goca upgrade` | None | **Check config version** against installed Goca. Use after updating goca binary. |
| `goca ci` | None | **Generate GitHub Actions** CI/CD pipeline. Use when setting up automation. |
| `goca template init\|list` | Varies | **Manage custom templates** — init creates template dir, list shows available. |
| `goca mcp-server` | None | **Start MCP server** for AI tools integration. Use `--print-config vscode\|claude\|cursor\|zed` to get client config. |
| `goca version` | None | **Show version** + build info. `--short` for just the number. |

---

## Common Workflows

### 1. Full feature from scratch (recommended)
```bash
goca init myproject --module github.com/user/myproject
cd myproject
goca feature User --fields "Name:string,Email:string,Age:int"
```
Generates: entity + usecase + handler + repository + messages + DI + routes + main.go update.

### 2. Modular build (fine-grained control)
```bash
goca init myproject --module github.com/user/myproject
goca entity Product --fields "Name:string,Price:float64" --validation --timestamps
goca usecase Product --entity Product --operations "create,read,update,delete,list"
goca repository Product --database postgres --transactions
goca handler Product --type http --validation --swagger
goca messages Product --all
goca di --features "Product" --database postgres
```

### 3. Add feature to existing project
```bash
goca feature Order --fields "CustomerID:uint,Total:float64,Status:string"
# or manually:
goca entity Order --fields "..."
goca usecase Order --entity Order
goca repository Order --database postgres
goca handler Order --type http
goca messages Order --all
goca integrate --all  # auto-wires everything
```

### 4. TDD-first approach
```bash
goca interfaces User --all
# Now write tests against interfaces before implementation
goca entity User --fields "..."
goca usecase User --entity User
goca repository User --database postgres
goca handler User --type http
```

### 5. Add middleware
```bash
goca middleware api --types "cors,logging,recovery,request-id,rate-limit"
```

### 6. CI/CD setup
```bash
goca ci --provider github-actions --with-docker --with-deploy
```

### 7. Diagnose issues
```bash
goca doctor        # quick health check
goca analyze       # deep analysis (arch, security, quality, deps)
goca doctor --fix  # auto-create missing directories
```

### 8. After updating Goca
```bash
goca upgrade            # check config compatibility
goca upgrade --update   # write new version to .goca.yaml
```

---

## Field Syntax

All `--fields` flags use the same format:
```
--fields "FieldName:Type,FieldName2:Type2"
```

**Supported types (20):** `string`, `int`, `int8`-`int64`, `uint`, `uint8`-`uint64`, `uintptr`, `byte`, `rune`, `float32`, `float64`, `bool`, `time.Time`, `[]byte`, `interface{}`

**Naming rules:**
- Name regex: `^[A-Za-z][A-Za-z0-9]*$`
- PascalCase for fields (e.g., `Email`, `CreatedAt`)
- Length: 1-50 characters per field name

**Auto-generated fields:**
- `ID` (uint, primary key, auto-increment) — always added
- `CreatedAt`, `UpdatedAt` (time.Time) — with `--timestamps`
- `DeletedAt` (gorm.DeletedAt) — with `--soft-delete`

---

## Error Tracing & Debugging

### Common errors and solutions

| Error | Cause | Fix |
|-------|-------|-----|
| `--fields flag is required` | Missing required fields | Add `--fields "Name:type"` |
| `invalid entity name` | Name doesn't match regex | Use PascalCase alphanumeric only (`^[A-Za-z][A-Za-z0-9]*$`) |
| `invalid field format: ...` | Wrong field syntax | Use `Name:type` pairs separated by commas |
| `file already exists` | Conflict with existing file | Add `--force` to overwrite, or `--backup` to backup first |
| `entity already exists` | Duplicate entity name | Use different name or `--force` |
| `unknown database type` | Invalid `--database` | Use one of: postgres, postgres-json, mysql, mongodb, sqlite, sqlserver, elasticsearch, dynamodb |
| `unknown handler type` | Invalid `--type` | Use one of: http, grpc, cli, worker, soap |
| Go build fails after generation | Missing deps or imports | Run `go mod tidy` then `go build ./...` |

### Debug workflow
```bash
goca doctor              # 1. Check project health
goca doctor --fix        # 2. Auto-fix missing dirs
goca analyze             # 3. Deep analysis
goca analyze --quality   # or just quality checks
goca analyze --security  # or just security checks
goca analyze --output json  # machine-readable output
```

### Safety first (always available on generation commands)
```bash
goca feature User --fields "..." --dry-run    # preview without creating
goca feature User --fields "..." --force       # overwrite existing
goca feature User --fields "..." --backup      # backup before overwrite
goca feature User --fields "..." --dry-run --force --backup  # combine safely
```

---

## Flag Reference by Command

### `goca feature <name>` — Full feature generation
| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--fields` | `-f` | **required** | Entity fields |
| `--database` | `-d` | `postgres` | postgres, mysql, mongodb, sqlite, sqlserver, elasticsearch, dynamodb |
| `--handlers` | | `http` | http, grpc, cli, worker, soap |
| `--validation` | | `false` | Include validations |
| `--business-rules` | `-b` | `false` | Business rule methods |
| `--integration-tests` | | `false` | Generate integration tests |
| `--mocks` | | `false` | Generate testify mocks |
| `--cache` | `-c` | `false` | Redis cache decorator |
| `--dry-run` | | `false` | Preview only |
| `--force` | | `false` | Overwrite existing |
| `--backup` | | `false` | Backup before write |

### `goca entity <name>` — Domain entity
| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--fields` | `-f` | **required** | Entity fields |
| `--validation` | | `false` | Include Validate() method |
| `--business-rules` | `-b` | `false` | Business rules (IsAdult, IsExpensive, etc.) |
| `--timestamps` | `-t` | `false` | CreatedAt + UpdatedAt |
| `--soft-delete` | `-s` | `false` | DeletedAt field |
| `--tests` | | `true` | Unit tests |
| `--dry-run` / `--force` / `--backup` | | `false` | Safety flags |

### `goca usecase <name>` — Business logic
| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--entity` | `-e` | **required** | Associated entity name |
| `--operations` | `-o` | `create,read` | create, read, update, delete, list |
| `--dto-validation` | `-d` | `false` | DTOs with validation |
| `--async` | `-a` | `false` | Async operations |

### `goca handler <entity>` — Protocol adapter
| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--type` | `-t` | `http` | http, grpc, cli, worker, soap |
| `--validation` | | `false` | Input validation |
| `--swagger` | `-s` | `false` | Swagger docs (HTTP only) |

### `goca repository <entity>` — Persistence
| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--database` | `-d` | `""` | postgres, mysql, mongodb, sqlite, sqlserver, elasticsearch, dynamodb |
| `--interface-only` | `-i` | `false` | Interface only |
| `--cache` | `-c` | `false` | Cache layer |
| `--transactions` | `-t` | `false` | Transaction support |

### `goca di` — Dependency injection
| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--features` | `-f` | **required** | Comma-separated feature names |
| `--database` | `-d` | `postgres` | Database type |
| `--wire` | `-w` | `false` | Use Google Wire |

### `goca middleware <name>` — HTTP middleware
| Flag | Default | Description |
|------|---------|-------------|
| `--types` | `cors,logging,recovery` | cors, logging, auth, rate-limit, recovery, request-id, timeout |

---

## Template System

### Built-in templates (in `cmd/templates.go`):
- `entityTemplate` — Entity struct + validation + methods
- `useCaseTemplate` — Interface + service + DTOs
- `repositoryTemplate` — Interface + DB implementation
- `handlerTemplate` — HTTP handler with CRUD

### Custom templates:
```bash
goca template init   # creates .goca/templates/ with defaults
goca template list   # shows available custom templates
```
Custom `.tmpl` / `.tpl` files override built-in templates.

### Template functions available:
`title`, `lower`, `upper`, `camel` / `toCamelCase`, `pascal` / `toPascalCase`, `snake` / `toSnakeCase`, `kebab` / `toKebabCase`, `plural` / `toPlural`, `singular` / `toSingular`, `join`, `split`, `contains`, `hasPrefix`, `hasSuffix`, `trimSpace`, `replace`, `replaceAll`

---

## MCP Integration (AI Tool Calling)

Goca can run as an MCP server for AI assistants:
```bash
goca mcp-server --print-config vscode   # VS Code config snippet
goca mcp-server --print-config claude   # Claude Desktop config
goca mcp-server --print-config cursor   # Cursor config
goca mcp-server --print-config zed      # Zed config
```

### Available MCP tools (5 core + 11 utility):

**Core generation tools (auto-adds `--no-interactive`):**
- `goca_feature` — name, fields, database, validation, business_rules, handlers, integration_tests, mocks, dry_run, force
- `goca_entity` — name, fields, validation, business_rules, timestamps, soft_delete, tests, dry_run, force
- `goca_usecase` — name, fields, validation, dry_run, force
- `goca_repository` — name, database, dry_run, force
- `goca_handler` — name, type, validation, dry_run, force

**Utility tools:**
`goca_di`, `goca_integrate`, `goca_interfaces`, `goca_messages`, `goca_mocks`, `goca_init`, `goca_doctor`, `goca_upgrade`, `goca_ci`, `goca_middleware`, `goca_analyze`

---

## Quick Reference Card

### Most common command patterns
```bash
# New project
goca init <name> -m <module>

# Quick full feature
goca feature <Name> -f "field:type,field2:type2" -d <db>

# Add to existing
goca entity <Name> -f "..."
goca integrate --all

# Validate
goca doctor && goca analyze

# Safety
<command> --dry-run        # preview
<command> --dry-run --force --backup  # safe overwrite preview
```

### Naming conventions
- Entities: PascalCase (`User`, `ProductOrder`)
- Files: snake_case (`user.go`, `product_order.go`)
- Packages: lowercase (`domain`, `usecase`)

### Database types available
`postgres`, `postgres-json`, `mysql`, `mongodb`, `sqlite`, `sqlserver`, `elasticsearch`, `dynamodb`

### Handler types available
`http`, `grpc`, `cli`, `worker`, `soap`

### Operations for usecase
`create`, `read`, `update`, `delete`, `list`
