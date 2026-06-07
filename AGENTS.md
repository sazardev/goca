# AGENTS.md — Goca Code Generator

Quick navigation for AI agents. Every line answers "would an agent miss this?"

## 📁 Repo Structure

**Entry point:** `doc.go` — `package main` with `func main()`, calls `cmd.Execute()`.
No `main.go` exists; `doc.go` serves both doc + entry point roles.

**CLI commands:** `cmd/*.go` — Cobra subcommands.
- Layer generation: `entity.go`, `usecase.go`, `handler.go`, `repository.go`
- Composite: `feature.go`, `init.go`, `integrate.go`
- Utilities: `di.go`, `interfaces.go`, `messages.go`, `middleware.go`
- Supporting: `config_manager.go`, `template_manager.go`, `mcp_server.go`, `analyze.go`, `ci.go`, `doctor.go`, `upgrade.go`
- Templates: `templates.go`, `template_components.go`, `project_templates.go`

**Goca's own demo packages** (NOT generated output — reference implementations):
`internal/domain/`, `internal/usecase/`, `internal/handler/`, `internal/repository/`, `internal/di/`, `internal/messages/`, `internal/interfaces/`, `internal/mocks/`, `internal/constants/`

**Testing:** `internal/testing/` — `comprehensive_test.go` + `suite.go` + `tests/` + `framework/`.

**External rules (per-file-type):** `.github/instructions/` — 6 instruction files (testing, security, clean-architecture, commands, quality, docs-vitepress).
**Workflow prompts:** `.github/prompts/` — 6 prompt templates (new-command, audit-codegen, add-tests, etc.).
**Agent definition:** `.github/agents/goca-forge.agent.md`.
**Reusable skills:** `.github/skills/` — 3 skills (codegen-testing, goca-architecture, mcp-tools).
**OpenCode plugin:** `.opencode/package.json` — `@opencode-ai/plugin: 1.15.13`.
**Config:** `.goca.yaml` (project defaults), `.goreleaser.yml` (Homebrew + multi-arch release).

## 🔧 Developer Commands

```bash
make dev          # fmt → lint → test → build (pre-commit check)
make lint         # golangci-lint run
make test         # go test -v ./...
make test-coverage  # → coverage.html
make build        # Build binary with version injection
make install      # go install . → $GOPATH/bin/goca
```

**CLI integration tests (`go run internal/testing/test_runner.go`):**
```bash
make test-cli-comprehensive   # All CLI scenarios
make test-cli-init            # goca init only
make test-cli-feature         # goca feature only
make test-cli-entity          # goca entity only
make test-cli-fast            # Compilation-only (quick iteration)
make test-cli-quality         # Code quality of generated output
make test-all                 # Unit + comprehensive CLI tests
```

**Run a specific test:**
```bash
go test -v ./internal/testing -run TestGocaInitCommand
go test -v ./cmd -run TestValidateEntityName
```

**Version injection:** `go build -ldflags "-X github.com/sazardev/goca/cmd.Version=$(VERSION) ..."`
Requires git tags. Use `VERSION=dev make build` if no tags exist.

## 🔑 Codegen Quality Gates (non-negotiable)

Every template change must pass for all 8 database backends (postgres, mysql, sqlite, sqlserver, mongodb, redis, cassandra, dynamodb):
1. `go build ./...` — zero errors
2. `go vet ./...` — zero warnings
3. `go fmt ./...` — properly formatted
4. No unused imports
5. Template renders with empty `Fields` slice (edge case)
6. Optional features toggle on/off: timestamps, soft_delete, uuid, audit, versioning

After ANY template change: `make test-cli-comprehensive`.

**Template data:** Build `TemplateData` completely before passing to `template.Execute`. Missing fields → nil pointer dereference at runtime. All template strings are static constants — never constructed from user input.

## 🏛️ Architecture Rules (generated code)

Dependency direction (enforced, points inward):
```
HTTP Handler → UseCase Interface → Repository Interface → Domain Entity
```
- `internal/handler` imports usecase interfaces only (never concrete types)
- `internal/usecase` imports repository interfaces (never handler)
- `internal/domain` imports nothing from `internal/`
- DI container is the only place that wires concrete → interface

All constructors accept interfaces, not concrete types. Business logic lives in usecases, never handlers.

## 🛡️ Security Validation

- Names match `^[A-Za-z][A-Za-z0-9]*$` — validated via `CommandValidator`
- Paths: `filepath.Join` with validated components only; verify within `os.Getwd()`
- Shell: `exec.Command` with argument arrays — never string concatenation
- File permissions: source `0644`, dirs `0755`
- No deletion operations in any command
- No sensitive data (DB passwords) logged or surfaced in CLI output
- All file writes go through `SafetyManager.WriteFile()` — never `os.WriteFile` directly in command code

See `.github/instructions/security.instructions.md` for full OWASP checklist.

## 🚫 Regla Absoluta: No Bypass de Hooks

NUNCA uses `--no-verify`, `LEFTHOOK=0`, ni ningún otro mecanismo para saltar los hooks de pre-commit o pre-push.

Si un hook falla:
1. Lee el error completo — identifica QUÉ herramienta falló (golangci-lint, nilaway, gosec, staticcheck, tests, build, docs)
2. Investiga la causa raíz del error en el código o configuración
3. Corrige el problema real — no lo escondas ni lo silencies
4. Si después de investigar no sabes cómo resolverlo, PREGUNTA al usuario explicando el error y mostrando las opciones

Esta regla es estricta e inquebrantable. Ignorarla viola la confianza del proyecto.

## ⚡ Key Files by Task

| Task | Files |
|------|-------|
| Template author | `cmd/templates.go`, `cmd/template_components.go`, `cmd/project_templates.go` |
| Command developer | `cmd/entity.go`, `cmd/feature.go` (pattern reference) |
| Test engineer | `cmd/*_test.go`, `internal/testing/` |
| Architecture | `.github/instructions/clean-architecture.instructions.md` |
| Security auditor | `cmd/command_validator.go`, `cmd/safety.go` |
| MCP integration | `cmd/mcp_server.go`, `cmd/mcp_tools*.go`, `cmd/mcp_resources.go` |

## 🚨 Gotchas

1. **Template data must be fully populated** before rendering — missing fields → nil pointer dereference
2. **Generated imports must match used packages** — mismatch breaks `go vet`
3. **Test isolation:** always use `t.TempDir()`, never real filesystem
4. **Integration tests are slow** — `make test-cli-fast` for quick feedback
5. **Version info requires git tags** — `VERSION=dev make build` if no tags
6. **`pre-release-check`** runs: `fmt → lint → test → test-cli-comprehensive`
7. **Config loading fails silently for non-init projects** — always handle config load warnings
8. **MCP server mode** — `goca mcp-server` starts a stdio MCP server exposing Goca as tools for AI assistants

---

**Go 1.25.1 required.** Module: `github.com/sazardev/goca`.
