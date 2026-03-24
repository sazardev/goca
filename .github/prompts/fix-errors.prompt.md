---
mode: agent
description: Fix all go build errors, go vet warnings, and test failures across the Goca codebase. Runs diagnostics first, then applies targeted fixes.
tools:
  - mcp_oraios_serena_get_symbols_overview
  - mcp_oraios_serena_find_symbol
  - mcp_oraios_serena_find_referencing_symbols
  - mcp_oraios_serena_replace_symbol_body
  - run_in_terminal
  - read_file
  - replace_string_in_file
---

# Fix Build Errors and Test Failures

Diagnose and fix all `go build`, `go vet`, and `go test` failures in the Goca codebase.

## Step 1 — Run Full Diagnostics

```bash
# Build errors
go build ./... 2>&1

# Vet warnings
go vet ./... 2>&1

# Test failures (with race detector)
go test ./... -race -timeout 120s 2>&1
```

Collect ALL output before starting any fixes.

## Step 2 — Categorize Errors

For each error, categorize as:

- **Build error** (compilation fails, must fix first)
- **Vet warning** (bad code pattern, fix before build errors resolved)
- **Test failure** (assertion fails, fix after build/vet pass)
- **Race condition** (concurrent state access, fix last)

## Step 3 — Fix Build Errors (Priority 1)

Use Serena to understand the failing code before touching it:

```
mcp_oraios_serena_find_symbol("<failing function>", include_body=true)
```

Common Goca build errors and their fixes:

**Undefined template variable:**

```go
// Error: template: X:10: "Entity" is not a field of struct type TemplateData
// Fix: check TemplateData struct, use correct field path
// Command: mcp_oraios_serena_find_symbol("TemplateData", include_body=true)
```

**Unused import:**

```go
// Error: "fmt" imported and not used
// Fix: remove import or add usage
// Use: mcp_oraios_serena_find_symbol("imports in file", include_body=true)
```

**Type mismatch:**

```go
// Error: cannot use x (type string) as type FeatureFlags
// Fix: trace the TemplateData construction
```

## Step 4 — Fix Vet Warnings (Priority 2)

Common vet issues in Goca:

```bash
# printf format mismatch
go vet ./... 2>&1 | grep "printf"

# Shadowed error variable
go vet ./... 2>&1 | grep "shadow"
```

## Step 5 — Fix Test Failures (Priority 3)

For each failing test:

1. Read test function with Serena `find_symbol`
2. Read the function under test with Serena `find_symbol`
3. Identify root cause (implementation bug vs. test bug)
4. Fix the implementation (prefer) or update test if requirements changed

**Critical:** Do not remove tests to make them pass. Fix the underlying code.

## Step 6 — Fix Race Conditions (Priority 4)

```bash
go test ./... -race -count=1 2>&1 | grep "DATA RACE"
```

Common causes in Goca:

- Shared `SafetyManager` state in parallel tests → use separate instance per test
- Global cobra command flags read concurrently → use `cmd.Flags()` local copies

## Step 7 — Verify Everything Passes

```bash
go build ./...           # must exit 0
go vet ./...             # must exit 0
go test ./... -race      # must exit 0, all tests pass
```

## Step 8 — Verify Generated Code Still Works

If template files were changed, run a quick generation test:

```bash
mkdir /tmp/goca-verify && cd /tmp/goca-verify
go mod init verify
goca entity TestEntity --fields "Name:string" --validation
go build ./...
go vet ./...
```

## Rules for Fixing

- NEVER remove a test to make it pass
- NEVER use `//nolint` comments without explanation
- NEVER use `_ =` to silence an error that should be handled
- ALWAYS wrap errors with context when fixing error handling
- ALWAYS use `t.TempDir()` if adding filesystem operations to tests
