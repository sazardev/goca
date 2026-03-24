---
mode: agent
description: Audit the code that Goca generates. Verify all templates produce compilable, vet-clean Go code for every database backend and handler type. Find and fix template bugs.
tools:
  - mcp_oraios_serena_get_symbols_overview
  - mcp_oraios_serena_find_symbol
  - mcp_oraios_serena_find_referencing_symbols
  - mcp_oraios_serena_replace_symbol_body
  - mcp_context7_resolve-library-id
  - mcp_context7_get-library-docs
  - run_in_terminal
  - read_file
  - replace_string_in_file
---

# Audit Generated Code

Systematically verify that every template in Goca produces valid, compile-clean Go code.

## Scope

`$SCOPE` — which templates to audit. Default: all. Options: `entity`, `usecase`, `repository`, `handler`, `di`, `init`, `interfaces`, `messages`, `mocks`.

## Step 1 — Inventory All Templates

Use Serena to list all template constants:

```
mcp_oraios_serena_get_symbols_overview("cmd/templates.go")
mcp_oraios_serena_get_symbols_overview("cmd/template_components.go")
mcp_oraios_serena_get_symbols_overview("cmd/project_templates.go")
```

## Step 2 — For Each Template, Verify

### Template Syntax Correctness

- Open/close `{{` `}}` pairs are balanced
- All `.Fields` range loops have correct field access: `{{.Name}}`, `{{.Type}}`, `{{.JSONTag}}`
- All feature flags checked correctly: `{{if .Features.Timestamps}}` not `{{if .Timestamps}}`
- Module path used as `{{.Module}}` not hardcoded
- Package declarations use `{{.Entity.Package}}` not hardcoded

### Import Completeness Check

For each template, verify imports match used packages:

- Uses `fmt` → `"fmt"` in imports
- Uses `time.Time` → `"time"` in imports
- Uses GORM → `"gorm.io/gorm"` in imports
- Uses context → `"context"` in imports
- No unused imports (causes `go build` failure)

### Edge Cases to Test

1. **Zero fields:** `TemplateData{Fields: []FieldData{}}` — should produce valid empty struct
2. **Single field:** One field of each basic type
3. **All feature flags false:** No timestamps, no soft delete, no validation
4. **All feature flags true:** All optional sections enabled
5. **Long entity names:** 50+ character names (within validation limits)

## Step 3 — Generate Sample Output

For each template variant, generate with a test entity:

```bash
mkdir -p /tmp/goca-audit
cd /tmp/goca-audit
go mod init audittest

# Test each command
goca entity Product --fields "Name:string,Price:float64" --validation --timestamps
goca usecase Product --fields "Name:string,Price:float64"
goca repository Product --database postgres
goca handler Product --handlers http
goca di

# Verify compilation
go build ./...
go vet ./...
```

## Step 4 — Database Backend Matrix

Test repository templates for ALL 8 databases:

```bash
for db in postgres mysql sqlite sqlserver mongodb redis cassandra dynamodb; do
    mkdir -p /tmp/goca-audit-$db
    cd /tmp/goca-audit-$db
    go mod init auditdb
    goca feature Product --database $db --fields "Name:string"
    go build ./... 2>&1 | grep -v "^$" && echo "FAIL: $db" || echo "PASS: $db"
done
```

## Step 5 — Handler Types Matrix

Test handler templates for all protocols:

```bash
for handler in http grpc cli worker; do
    goca handler Product --handlers $handler
    # verify generated file has correct package imports
done
```

## Step 6 — Verify No Go Vet Warnings

```bash
go vet ./... 2>&1
```

Common vet issues in generated code:

- `printf` format string mismatches (`%s` with non-string value)
- Unreachable code after `return`
- Shadowed loop variables (avoid `i := i` pattern in generated code)
- Struct field misalignment (informational in vet, but clean code preferred)

## Step 7 — Report Issues Found

For each issue:

1. Identify the template constant name and line
2. Identify the specific template syntax error
3. Use Serena `find_symbol` with `include_body=true` to read the template
4. Apply fix with `replace_symbol_body`
5. Re-run the compilation test for that template

## Acceptance Criteria

- [ ] All templates produce code that compiles with `go build ./...`
- [ ] All templates produce code that passes `go vet ./...`
- [ ] Zero unused imports in any generated file
- [ ] Zero undefined variables in any generated template output
- [ ] All 8 database backends generate valid repository code
- [ ] All handler type templates generate valid handler code
- [ ] Empty fields edge case produces valid (empty struct body) code
