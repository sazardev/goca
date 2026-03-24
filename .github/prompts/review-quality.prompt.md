---
mode: agent
description: Full code quality review of a Goca file or package. Check microfunctions, error handling, security, test coverage, and Clean Architecture compliance.
tools:
  - mcp_oraios_serena_get_symbols_overview
  - mcp_oraios_serena_find_symbol
  - mcp_oraios_serena_find_referencing_symbols
  - mcp_context7_resolve-library-id
  - mcp_context7_get-library-docs
  - run_in_terminal
  - read_file
  - replace_string_in_file
---

# Code Quality Review

Perform a thorough quality review of the target file or package. Identify and fix all violations.

## Target

`$TARGET` — file path or package to review (e.g., `cmd/entity.go`, `cmd/`, `internal/usecase/`).

## Review Process — Use Serena Throughout

**Do NOT read entire files.** Use symbolic tools:

```
1. mcp_oraios_serena_get_symbols_overview(target)    → list all functions
2. mcp_oraios_serena_find_symbol(name, body=false)   → get signature
3. mcp_oraios_serena_find_symbol(name, body=true)    → read body only when needed
```

## Checklist 1: Microfunctions

For each function found:

- [ ] ≤ 50 lines (excluding blank lines and comments)
- [ ] Cyclomatic complexity ≤ 10 (count `if`, `for`, `switch`, `&&`, `||`)
- [ ] ≤ 5 parameters (more → group into struct)
- [ ] Does exactly ONE thing (its name describes it fully)

**Fix pattern for oversized functions:**

```go
// Before: 80-line function with 3 responsibilities
func generateEntity(name, fields, database string, ...) error {
    // 1. validate (20 lines)
    // 2. build template data (30 lines)
    // 3. write files (30 lines)
}

// After: split into focused helpers
func validateEntityInputs(name, fields string) error { ... }        // ≤ 20 lines
func buildEntityTemplateData(name, fields string) (TemplateData, error) { ... } // ≤ 30 lines
func writeEntityFiles(data TemplateData, sm *SafetyManager) error { ... }       // ≤ 30 lines
```

## Checklist 2: Error Handling

- [ ] All errors are wrapped with context: `fmt.Errorf("context: %w", err)`
- [ ] No silent error drops: `_ = someFunc()` is forbidden for error returns
- [ ] Fatal errors in commands use `errorHandler.HandleError()`, not `log.Fatal` or `panic`
- [ ] `os.Exit()` only called via `ErrorHandler` — never directly in command logic
- [ ] No named returns in any function

## Checklist 3: Security

- [ ] All user inputs pass through `CommandValidator` before use in paths/templates
- [ ] `filepath.Join` used for all path construction (no string concatenation)
- [ ] Constructed paths verified to stay within project root
- [ ] `exec.Command` uses argument array (no shell string concat)
- [ ] No `interface{}` / `any` in new code
- [ ] File permissions: `0644` for files, `0755` for dirs

## Checklist 4: Import Quality

- [ ] Imports grouped: stdlib | external | internal (blank lines between groups)
- [ ] No unused imports
- [ ] No dot imports (`import . "pkg"`)
- [ ] Correct import alias if needed (no ambiguity)

## Checklist 5: Naming

- [ ] Types/Functions: PascalCase
- [ ] Variables/Fields: camelCase
- [ ] Constants: PascalCase (exported) / camelCase (unexported)
- [ ] No stuttering (`config.ConfigManager` → `config.Manager`)
- [ ] Test functions: `TestFunctionName_Scenario`

## Checklist 6: Interface Usage

- [ ] Struct fields use interface types, not concrete types
- [ ] Constructor functions accept interface parameters
- [ ] No concrete type assertions except in DI container or test setup

## Checklist 7: Test Coverage

```bash
go test ./... -coverprofile=cover.out
go tool cover -func=cover.out | grep "$TARGET"
```

- [ ] ≥ 80% coverage for `cmd/` package
- [ ] Every exported function has at least one test
- [ ] Security edge cases tested

## Checklist 8: Clean Architecture (for `internal/`)

- [ ] No upward dependency violations
- [ ] Domain has zero internal imports
- [ ] UseCase only imports repository interfaces + domain
- [ ] Handler only imports usecase interfaces
- [ ] DI container is only place that creates concrete instances

## Output Format

For each violation found, report:

```
FILE: cmd/entity.go
FUNCTION: generateEntityFiles
VIOLATION: Microfunctions — function is 73 lines, exceeds 50-line limit
SEVERITY: Medium
FIX: Split into buildEntityData() and writeEntityFiles()
```

Then implement fixes immediately using Serena `replace_symbol_body`.

## Verification After Fixes

```bash
go build ./...    # must pass
go vet ./...      # must pass
go test ./...     # must pass
```
