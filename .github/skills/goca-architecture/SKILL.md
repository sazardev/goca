# Skill: Goca Architecture

**Domain:** Understanding the internal architecture of the Goca CLI code generator — how commands, templates, validators, and subsystems interact.

---

## How a Goca Command Works (End-to-End)

When a user runs `goca entity Product --fields "Name:string"`:

```
cobra.Command.Run()
    ↓
1. Flag parsing (cobra flags → local vars)
    ↓
2. ConfigIntegration.LoadConfigForProject() — reads .goca.yaml
   ConfigIntegration.MergeWithCLIFlags()    — CLI overrides config
    ↓
3. SafetyManager = NewSafetyManager(dryRun, force, backup)
    ↓
4. CommandValidator.ValidateEntityName(name) — rejects invalid names
   FieldValidator.ParseFields(fieldsStr)    — parses "Name:string,Price:float64"
    ↓
5. TemplateGenerator.BuildTemplateData()    — populates TemplateData struct
    ↓
6. TemplateGenerator.GenerateFromTemplate() — text/template rendering
    ↓
7. SafetyManager.WriteFile(path, content)   — writes with conflict checking
    ↓
8. DependencyManager.AddDependency()        — queues go.mod updates
   DependencyManager.RunGoGet()             — executes go get
    ↓
9. Print summary (files created, dry-run preview, etc.)
```

## Key Subsystem Details

### TemplateData Struct

The central data structure passed to all templates:

```go
type TemplateData struct {
    Entity   EntityData    // Name, NameLower, NamePlural, Package
    Fields   []FieldData   // Name, Type, JSONTag, GormTag, ValidateTag, flags
    Module   string        // go module path (e.g. "github.com/user/project")
    Database string        // "postgres", "mysql", "sqlite", etc.
    Features FeatureFlags  // Validation, BusinessRules, Timestamps, SoftDelete, etc.
    Imports  []string      // deduplicated import paths for the generated file
    Methods  []MethodData  // (for use case / handler templates)
}
```

All template variables MUST reference `TemplateData` fields — never hardcoded strings.

### SafetyManager Flow

```
WriteFile(path, content)
    → CheckFileConflict(path)
        → if file exists AND !Force → error
        → if file exists AND Backup → BackupFile(path)
        → if DryRun → record to createdFiles, return nil (no write)
    → os.MkdirAll(dir, 0755)
    → os.WriteFile(path, []byte(content), 0644)
    → append to createdFiles
```

### ConfigIntegration Priority

Config is resolved in this priority order (highest first):

1. CLI flags explicitly set by user (`--database postgres`)
2. `.goca.yaml` configuration file values
3. Built-in defaults (e.g., `database: sqlite`)

### FieldValidator Type Support

Supports all standard Go types plus complex ones:

- Basic: `string`, `int`, `int64`, `float64`, `bool`, `uint`, `byte`, `rune`
- Pointers: `*string`, `*User`
- Slices: `[]string`, `[]*User`, `[][]string`
- Arrays: `[10]string`
- Maps: `map[string]string`, `map[string]interface{}`
- Channels: `chan string`, `<-chan int`, `chan<- bool`
- Functions: `func()`, `func(string) error`
- Interfaces: `interface{}`, `io.Reader`
- Qualified: `time.Time`, `uuid.UUID`

Invalid map keys (not comparable): `map[[]string]int`, `map[func()]string`

## Database Backend Matrix

Goca supports 8 database backends. Repository templates vary by database:

| Database  | Driver Package                 | Connection Type | Generated Driver Import |
| --------- | ------------------------------ | --------------- | ----------------------- |
| postgres  | `gorm.io/driver/postgres`      | DSN string      | ✓                       |
| mysql     | `gorm.io/driver/mysql`         | DSN string      | ✓                       |
| sqlite    | `gorm.io/driver/sqlite`        | file path       | ✓                       |
| sqlserver | `gorm.io/driver/sqlserver`     | DSN string      | ✓                       |
| mongodb   | `go.mongodb.org/mongo-driver`  | URI + client    | ✓ (non-GORM)            |
| redis     | `github.com/redis/go-redis/v9` | options         | ✓ (non-GORM)            |
| cassandra | `github.com/gocql/gocql`       | cluster         | ✓ (non-GORM)            |
| dynamodb  | `github.com/aws/aws-sdk-go-v2` | config          | ✓ (non-GORM)            |

MongoDB, Redis, Cassandra, and DynamoDB use non-GORM templates — they have separate template functions.

## Handler Protocol Matrix

| Protocol | Template       | Generated File                       | Routes      |
| -------- | -------------- | ------------------------------------ | ----------- |
| `http`   | HTTP REST      | `handler/http/<entity>_handler.go`   | gorilla/mux |
| `grpc`   | gRPC service   | `handler/grpc/<entity>_handler.go`   | proto-based |
| `cli`    | CLI subcommand | `handler/cli/<entity>_handler.go`    | cobra       |
| `worker` | Background job | `handler/worker/<entity>_handler.go` | goroutine   |

## Template String Organization

```
cmd/templates.go              ← entity, usecase, repository, handler templates
cmd/template_components.go    ← reusable partial templates (field blocks, import blocks)
cmd/project_templates.go      ← goca init project structure templates (main.go, go.mod, etc.)
```

## Error Types (`cmd/errors.go`)

Typed errors for clean error handling:

- `ErrInvalidEntityName` — name fails regex validation
- `ErrInvalidFieldType` — field type not supported
- `ErrInvalidDatabase` — unknown database backend
- `ErrFileConflict` — file already exists (use --force)
- `ErrTemplateRender` — template execution failed
- `ErrDependencyInstall` — go get failed

## Integration Test Architecture

```
internal/testing/tests/
├── entity_test.go      Tests goca entity command end-to-end
├── feature_test.go     Tests goca feature command (all layers)
├── usecase_test.go     Tests goca usecase command
├── handler_test.go     Tests goca handler command
├── init_test.go        Tests goca init command (project scaffold)
└── safety_test.go      Tests SafetyManager, NameConflictDetector, DependencyManager
```

Each test:

1. Creates `t.TempDir()` as project root
2. Initializes a Go module with `go mod init <test-module>`
3. Runs the command targeting that directory
4. Verifies expected files exist
5. Runs `go build ./...` and `go vet ./...` to confirm validity
