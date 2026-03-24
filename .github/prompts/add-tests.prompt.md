---
mode: agent
description: Add comprehensive tests to a Goca command or internal package. Covers unit tests, integration tests, security tests, and compilation verification.
tools:
  - mcp_oraios_serena_get_symbols_overview
  - mcp_oraios_serena_find_symbol
  - mcp_oraios_serena_find_referencing_symbols
  - mcp_context7_resolve-library-id
  - mcp_context7_get-library-docs
  - run_in_terminal
  - read_file
  - create_file
  - replace_string_in_file
---

# Add Tests

Add comprehensive tests for the target file or function. Cover all quality gates.

## Target

`$TARGET` — the file, function, or package to test.

## Pre-work: Understand the Code

1. Use Serena `get_symbols_overview` on the target file
2. Use `find_symbol` with `include_body=false` to list all exported functions
3. For each function to test, use `find_symbol` with `include_body=true` to understand logic
4. Use `find_referencing_symbols` to understand how it's called in practice

**Do NOT read the entire file** — use symbolic tools to get only what's needed.

## Unit Test Requirements

### File location

Co-located with source: `cmd/<target>_test.go` or `internal/<pkg>/<target>_test.go`

### Structure template

```go
package cmd

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFunctionName_HappyPath(t *testing.T) {
    t.Parallel()
    // ...
}

func TestFunctionName_TableDriven(t *testing.T) {
    t.Parallel()
    cases := []struct {
        name      string
        input     string
        wantErr   bool
        wantValue string
    }{
        // populate cases
    }
    for _, tc := range cases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            // test body
        })
    }
}
```

### Required test case categories

For every function with user input:

1. **Happy path** — valid typical input
2. **Edge cases** — boundary values, empty collections
3. **Error cases** — expected error conditions
4. **Security cases:**
   - Empty string `""`
   - Path traversal `"../etc/passwd"`, `"../../secret"`
   - Shell metacharacters `"; rm -rf /"`, `"| cat /etc/passwd"`
   - Null bytes `"valid\x00.go"`
   - Very long string (> 500 chars)

## Filesystem Test Requirements

```go
// ALWAYS use t.TempDir() — never os.TempDir() or hardcoded paths
func TestFileOperation(t *testing.T) {
    t.Parallel()
    dir := t.TempDir() // auto-cleaned after test
    // use dir for all file operations
}
```

## Integration Test Requirements (for command tests)

Location: `internal/testing/tests/<command>_test.go`

Must include:

```go
func TestCommandGeneratesValidCode(t *testing.T) {
    dir := t.TempDir()
    // Initialize a minimal Go module in dir
    // Run the command targeting dir
    // Verify files exist at expected paths
    // Verify compilation:
    buildCmd := exec.Command("go", "build", "./...")
    buildCmd.Dir = dir
    out, err := buildCmd.CombinedOutput()
    require.NoError(t, err, "generated code must compile: %s", string(out))
    // Verify vet:
    vetCmd := exec.Command("go", "vet", "./...")
    vetCmd.Dir = dir
    vetOut, err := vetCmd.CombinedOutput()
    require.NoError(t, err, "generated code must pass vet: %s", string(vetOut))
}
```

## Coverage Verification

After writing tests, run coverage:

```bash
go test ./cmd/... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep -E "(total|$TARGET)"
```

Target: ≥ 80% for `cmd/` package. Fix gaps found.

## Run Tests

```bash
# Unit tests
go test ./cmd/... -v -run TestFunctionName

# All tests
go test ./... -race

# With coverage
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
```

All must pass. No race conditions.
