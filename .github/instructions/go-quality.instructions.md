---
applyTo: "**/*.go"
---

# Go Code Quality Rules

## Microfunctions — Single Responsibility

Every function MUST do exactly ONE thing. Enforce these limits at all times:

- **Max lines per function:** 50 (excluding blank lines and comments)
- **Max cyclomatic complexity:** 10
- **Max parameters:** 5 — if more are needed, group into a struct

```go
// BAD — does too many things
func processEntity(name, fields, db, handler string, validation, timestamps bool) error { ... }

// GOOD — split into focused helpers
func buildEntityConfig(name string, flags EntityFlags) (EntityConfig, error) { ... }
func writeEntityFiles(cfg EntityConfig, sm *SafetyManager) error { ... }
```

## Error Handling

Always wrap errors with context. Never swallow or lose error information.

```go
// REQUIRED pattern
if err := os.MkdirAll(dir, 0755); err != nil {
    return fmt.Errorf("create directory %s: %w", dir, err)
}

// FORBIDDEN
os.MkdirAll(dir, 0755) // ignoring error
if err != nil { return err } // no context added
```

In `cmd/` commands, use `ErrorHandler` for fatal errors:

```go
if err := validator.ValidateEntityName(name); err != nil {
    validator.errorHandler.HandleError(err, "entity name validation")
    // HandleError exits — do not add return after
}
```

## No Named Returns

Named returns are forbidden. Always use explicit return values.

```go
// FORBIDDEN
func loadConfig() (cfg *Config, err error) { ... }

// CORRECT
func loadConfig() (*Config, error) { ... }
```

## No Global Mutable State

All state in `cmd/` must be held in struct receivers. Never use package-level `var` for mutable state.

```go
// FORBIDDEN
var currentConfig *Config

// CORRECT
type ConfigManager struct {
    config *Config
}
```

## Interfaces Over Concrete Types

Depend on interfaces everywhere. Concrete types are only used at initialization.

```go
// CORRECT — repository pattern
type ProductUseCase interface {
    Create(input CreateProductInput) (*domain.Product, error)
    GetByID(id uint) (*domain.Product, error)
    List() ([]domain.Product, error)
    Update(id uint, input UpdateProductInput) (*domain.Product, error)
    Delete(id uint) error
}
```

## Zero `interface{}` in New Code

Use generics or typed interfaces. The `any` / `interface{}` type is forbidden in new code.

```go
// FORBIDDEN
func processField(value interface{}) string { ... }

// CORRECT
func processField[T any](value T) string { ... }
// or use typed union via interface constraints
```

## Import Organization

Imports must be grouped in this order, separated by blank lines:

1. Standard library
2. External dependencies
3. Internal `github.com/sazardev/goca/...` packages

```go
import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "gopkg.in/yaml.v3"

    "github.com/sazardev/goca/internal/domain"
)
```

## Documentation Comments

- Every exported type, function, and method MUST have a Go doc comment
- Comments explain WHY, not WHAT (the code says what)
- Template strings in `templates.go` must have a comment describing the generated structure

## Naming Conventions

| Kind                   | Convention                  | Example                               |
| ---------------------- | --------------------------- | ------------------------------------- |
| Packages               | lowercase single word       | `cmd`, `domain`, `usecase`            |
| Types, Functions       | PascalCase                  | `SafetyManager`, `ValidateEntityName` |
| Variables, Fields      | camelCase                   | `entityName`, `fieldValidator`        |
| Constants (exported)   | PascalCase                  | `DefaultDatabase`                     |
| Constants (unexported) | camelCase                   | `defaultTimeout`                      |
| Test functions         | `TestFunctionName_Scenario` | `TestValidateEntityName_EmptyString`  |

## Generated Code Invariants

Any template string in `templates.go`, `template_components.go`, or `project_templates.go` MUST produce code that:

1. Compiles with `go build` — zero errors
2. Passes `go vet` — zero warnings
3. Has zero unused imports (use `// Intentionally imported for side effects` if needed)
4. Has correct package declarations matching the target directory
5. Uses the project module path from `TemplateData.Module`

Verify after template changes by running the relevant integration test.

## File Permissions

- Generated source files: `0644`
- Generated directories: `0755`
- Never use `0777` or `0666`
