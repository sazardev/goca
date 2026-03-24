---
applyTo: "**/*_test.go"
---

# Testing Standards for Goca

## Table-Driven Tests — Mandatory for Multi-Variant Functions

Use table-driven tests for any function that accepts multiple input variants.

```go
func TestValidateEntityName(t *testing.T) {
    t.Parallel()

    cases := []struct {
        name      string
        input     string
        wantErr   bool
        errSubstr string
    }{
        {name: "valid PascalCase", input: "Product", wantErr: false},
        {name: "valid single char", input: "A", wantErr: false},
        {name: "empty string", input: "", wantErr: true, errSubstr: "empty"},
        {name: "path traversal dot", input: "../../etc", wantErr: true},
        {name: "slash injection", input: "foo/bar", wantErr: true},
        {name: "starts with digit", input: "1Product", wantErr: true},
    }

    for _, tc := range cases {
        tc := tc // capture range variable
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            validator := NewCommandValidator()
            err := validator.ValidateEntityName(tc.input)
            if tc.wantErr {
                require.Error(t, err)
                if tc.errSubstr != "" {
                    assert.Contains(t, err.Error(), tc.errSubstr)
                }
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

## Filesystem Tests — Always Use `t.TempDir()`

Never use hard-coded paths or `os.TempDir()` directly. Always use `t.TempDir()`:

```go
func TestSafetyManager_WriteFile(t *testing.T) {
    t.Parallel()

    dir := t.TempDir()
    sm := NewSafetyManager(false, false, false)
    filePath := filepath.Join(dir, "output.go")

    err := sm.WriteFile(filePath, "package test\n")
    require.NoError(t, err)

    content, err := os.ReadFile(filePath)
    require.NoError(t, err)
    assert.Equal(t, "package test\n", string(content))
}
```

## Integration Tests — Verify Generated Code Compiles

Every integration test in `internal/testing/tests/` MUST verify that generated code compiles:

```go
func TestFeatureGeneratesValidCode(t *testing.T) {
    projectDir := t.TempDir()
    // ... generate files into projectDir ...

    // Verify compilation
    cmd := exec.Command("go", "build", "./...")
    cmd.Dir = projectDir
    out, err := cmd.CombinedOutput()
    require.NoError(t, err, "generated code must compile:\n%s", string(out))

    // Verify zero vet warnings
    vetCmd := exec.Command("go", "vet", "./...")
    vetCmd.Dir = projectDir
    vetOut, err := vetCmd.CombinedOutput()
    require.NoError(t, err, "generated code must pass vet:\n%s", string(vetOut))
}
```

## Assertion Imports

Always import both `assert` and `require` from testify:

```go
import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

- Use `require` for fatal assertions (test cannot continue if this fails)
- Use `assert` for non-fatal checks (test continues to collect all failures)

## Test Naming Convention

Format: `TestFunctionName_Scenario`

```go
func TestValidateEntityName_EmptyString(t *testing.T) { ... }
func TestValidateEntityName_PathTraversalAttack(t *testing.T) { ... }
func TestSafetyManager_DryRunPreventsWrite(t *testing.T) { ... }
func TestTemplateGenerator_RendersAllFields(t *testing.T) { ... }
```

## Parallelism

Add `t.Parallel()` to every test that does NOT mutate global state or shared directories:

```go
func TestMyFunction(t *testing.T) {
    t.Parallel()
    // ...
}
```

For subtests, add both:

```go
for _, tc := range cases {
    tc := tc
    t.Run(tc.name, func(t *testing.T) {
        t.Parallel()
        // ...
    })
}
```

## Coverage Requirements

| Package                 | Minimum Coverage                     |
| ----------------------- | ------------------------------------ |
| `cmd/`                  | ≥ 80%                                |
| New command files       | 100% of happy path + top error paths |
| `template_generator.go` | 100% data-path                       |
| `field_validator.go`    | 100% of type variants                |

## Mocking Strategy

Use `testify/mock` for all interface mocks. Never write manual stubs.

```go
type MockSafetyManager struct {
    mock.Mock
}

func (m *MockSafetyManager) WriteFile(path, content string) error {
    args := m.Called(path, content)
    return args.Error(0)
}
```

## Test File Organization

Unit tests are co-located with source:

```
cmd/
  entity.go
  entity_test.go        ← test ValidateEntityName, ParseFields, WriteEntityFile
  field_validator.go
  field_validator_test.go
```

Integration tests in `internal/testing/tests/`:

```
internal/testing/tests/
  entity_test.go         ← test full entity generation end-to-end
  feature_test.go        ← test full feature generation + compilation
  safety_test.go         ← test SafetyManager behavior
```

## Security Test Requirements

Every input validation function MUST include security-focused test cases:

- Empty string injection
- Path traversal (`../`, `../../`)
- Null bytes
- Shell metacharacters (`;`, `|`, `` ` ``, `$`)
- Extremely long inputs (> 1000 chars)
- Unicode normalization exploits
