# AGENTS.md — Goca Code Generator

Quick navigation for AI agents working with the Goca CLI codebase. All critical signals—exactly what would be missed without this file.

## 🏗 Repository Structure

**Core entry point:** `cmd/root.go` — Cobra CLI root with subcommands.

**Command implementations:** `cmd/*.go`  
- Layer commands: `entity.go`, `usecase.go`, `handler.go`, `repository.go`
- Composite: `feature.go` (all layers at once), `init.go` (scaffold project)
- Utilities: `di.go`, `interfaces.go`, `messages.go`, `middleware.go`
- Supporting: `config_manager.go`, `template_manager.go`, `mcp_server.go`, `mcp_tools.go`

**Generated code templates:** `cmd/templates.go`, `cmd/template_components.go`, `cmd/project_templates.go`  
These are critical—all generated Go code flows from these templates.

**Clean Architecture layers** (the generated output structure):
```
internal/
  domain/        # Pure entities (no external dependencies)
  usecase/       # Business logic with DTOs
  repository/    # Data persistence interfaces & implementations
  handler/       # HTTP/gRPC adapters
  di/            # Dependency injection container
```

**Testing infrastructure:** `internal/testing/`  
- `framework/` — test utilities and mocks
- `tests/` — integration test scenarios
- `comprehensive_test.go` — main integration suite
- `cli.go`, `validator.go`, `architecture.go` — test helpers

**Configuration:** `.goca.yaml` (project defaults), `go.mod` (Go 1.25.1), `Makefile` (dev workflow).

## 🚀 Developer Commands

### Essential Build & Test Workflow

**Full validation (must pass before any commit):**
```bash
make dev          # fmt → lint → test → build (recommended pre-commit check)
```

**Individual steps:**
```bash
make lint         # golangci-lint run
make test         # go test -v ./...
make test-coverage  # Coverage report → coverage.html
make build        # Build binary with version info
```

**CLI-specific tests (separate from unit tests):**
```bash
make test-cli-comprehensive   # All CLI scenarios
make test-cli-init            # goca init command only
make test-cli-feature         # goca feature command only
make test-cli-entity          # goca entity command only
make test-cli-fast            # Compilation-only tests (no codegen)
make test-all                 # Unit tests + comprehensive CLI tests
```

**Debug specific failing test:**
```bash
# Run single test with verbose output
go test -v ./internal/testing -run TestGocaCLIComprehensive

# Run subset of cmd tests
go test -v ./cmd -run TestNameOfFunction
```

**Install locally:**
```bash
make install      # go install . → $GOPATH/bin/goca
```

## 🔑 Critical Constraints

### Code Generation Quality Gates

**All generated code must pass (non-negotiable):**
1. `go build ./...` — zero build errors across all database backends (postgres, mysql, sqlite, etc.)
2. `go vet ./...` — zero vet warnings
3. `go fmt ./...` — formatted correctly
4. No unused imports

**Templates render correctly with edge cases:**
- Empty `Fields` slice (minimal entity)
- All database backends (postgres, mysql, sqlite, sqlserver, mongodb, redis, cassandra, dynamodb)
- Optional features toggled on/off (timestamps, soft_delete, uuid, audit, versioning)

**See CodegenAuditor mode in `.github/AGENTS.md` for full template audit workflow.**

### Dependency Architecture (Enforced in Generated Code)

**Direction of dependencies (must point inward):**
```
HTTP Handler → UseCase Interface → Repository Interface → Domain Entity
      ↓              ↓                       ↓                    ↓
    adapter      business logic        persistence           pure logic
```

**Violations to prevent:**
- ❌ `internal/handler` imports `internal/usecase` concrete types (must use interfaces)
- ❌ `internal/usecase` imports `internal/handler` (upward dependency)
- ❌ `internal/domain` imports anything from `internal/repository` or `internal/handler`
- ✅ All layers use dependency injection—no direct instantiation of concrete types outside DI

**See ArchitectGuard mode in `.github/AGENTS.md`.**

### Security Validation Required

**All user input validated via `CommandValidator`:**
- Name validation regex: `^[A-Za-z][A-Za-z0-9]*$`
- Path construction uses `filepath.Join` only
- Constructed paths verified to stay within `os.Getwd()`
- `exec.Command` uses argument arrays—never shell string concatenation

**No risky operations:**
- ✅ File creation with `0644` (files), `0755` (dirs)
- ❌ No deletion operations in any command
- ❌ No sensitive data (passwords, tokens) in logs or CLI output

**See SecurityAuditor mode in `.github/AGENTS.md`.**

## 📋 Common Tasks & Patterns

### Adding a New Command

1. Create `cmd/<command_name>.go`
2. Define command in `cmd/root.go` using Cobra pattern (see `cmd/entity.go` as template)
3. Add validation step using `CommandValidator` struct
4. Add integration test to `internal/testing/tests/`
5. Add test case to `comprehensive_test.go`
6. Update `GUIDE.md` with command documentation

**Command structure (required):**
- Use `cobra.Command` with clear `Use`, `Short`, `Long` descriptions
- Validate all user inputs before code generation
- Return detailed errors with context (`fmt.Errorf("...: %w", err)`)
- Generate files to deterministic paths relative to `os.Getwd()`

### Modifying Code Generation Templates

**Critical: After ANY template change, run:**
```bash
make test-cli-comprehensive   # Verify generated code compiles on all backends
```

1. Find template in `cmd/templates.go` or `cmd/template_components.go`
2. Trace where it's used via `TemplateData` struct (defined in same file)
3. Verify all `{{if .Features.*}}` blocks handle both true/false
4. Test with sample project: `goca init test-proj && cd test-proj && go build ./...`
5. Commit only when `go vet` and `go build` pass

**Template data flows:**
- `entity.go` → reads flags → builds `TemplateData` → passes to template → writes `domain/<entity>.go`
- `feature.go` → orchestrates entity + usecase + handler + repository (all templates used together)

### Adding Tests

**Naming convention:** `TestFunctionName_Scenario`  
**Use table-driven tests** for functions with multiple input variants.  
**Requirements:**
- `t.TempDir()` for all filesystem operations
- `t.Parallel()` for state-independent tests
- Include edge cases: empty input, special characters, boundary values
- Integration tests must verify `go build ./...` and `go vet ./...` on generated output

```go
func TestValidateName(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        wantError bool
    }{
        {"valid name", "User", false},
        {"empty name", "", true},
        {"special chars", "User-123", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            err := ValidateName(tt.input)
            if (err != nil) != tt.wantError {
                t.Errorf("got error %v, want %v", err, tt.wantError)
            }
        })
    }
}
```

## 📚 Key Files by Role

| Role | Key Files | Purpose |
|------|-----------|---------|
| **Template author** | `cmd/templates.go`, `cmd/template_components.go` | All generated code structure |
| **Command developer** | `cmd/entity.go`, `cmd/feature.go` | CLI command implementations |
| **Test engineer** | `cmd/*_test.go`, `internal/testing/` | Unit & integration tests |
| **Architecture enforcer** | `internal/domain/`, `internal/repository/` | Clean Architecture layers |
| **Security auditor** | `cmd/command_validator.go`, `cmd/safety.go` | Input validation & file safety |

## 🔗 Specialized Agent Modes

Goca project defines custom agent modes in `.github/AGENTS.md` with specific workflows:

- **CodegenAuditor** — Template validation & generated code quality
- **TestEngineer** — Unit & integration test authoring
- **ArchitectGuard** — Layer boundary enforcement
- **DocsWriter** — VitePress docs & wiki maintenance
- **SecurityAuditor** — Security review checklist

Load these modes when your task aligns with their domain.

## 📊 CI/CD Pipeline

**Triggers on:** push to `main`/`master`, all PRs

**Key workflow:** `.github/workflows/test.yml`
```
go mod download → go mod tidy → go mod verify →
go build ./... → goca init/help/version (smoke tests) →
go test ./internal/testing -run TestGocaCLIComprehensive →
go test ./cmd/... (unit tests)
```

**Pre-release checks:**
```bash
make pre-release-check    # fmt → lint → test → basic CLI tests
make release-auto         # Auto-detects version bump, creates git tag
```

## ⚙️ Style & Conventions

**Code format:** Go standard (`gofmt`) with 100-char line preference.  
**Import groups:** stdlib → external deps → internal packages (see `STYLE_GUIDE.md`).  
**Naming:** PascalCase for exported, camelCase for unexported.  
**Docs:** Every exported function must have godoc comment.  
**Error handling:** Always wrap errors with context: `fmt.Errorf("operation failed: %w", err)`.

See `STYLE_GUIDE.md` for full conventions.

## 🔗 External References

- **Complete docs:** https://sazardev.github.io/goca
- **Quick start:** https://sazardev.github.io/goca/getting-started
- **Architecture guide:** https://sazardev.github.io/goca/guide/clean-architecture
- **Cobra CLI framework:** https://cobra.dev/
- **Go text/template:** https://golang.org/pkg/text/template/
- **testify (test assertions):** https://github.com/stretchr/testify

## 🚨 Gotchas & Pitfalls

1. **Template data must be fully populated** before rendering—missing fields → nil pointer dereference
2. **Generated imports must include all used packages** or `go vet` fails—templates must match actual generated code
3. **Path traversal bugs** if `filepath.Join` not used—always validate constructed paths stay in project root
4. **Test isolation** — use `t.TempDir()`, never write to actual filesystem during tests
5. **Integration tests are slow** — keep `make test-cli-fast` for quick iteration, use full suite before commit
6. **Version info requires git tags** — `go build` injects version from git describe; use `VERSION=dev make build` if no tags exist

---

**Last updated:** May 2026  
**Goca version:** Go 1.25.1 required  
**Go Clean Architecture enforced:** Domain layer pure, dependencies point inward, all layers interface-based.
