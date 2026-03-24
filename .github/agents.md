# Goca — Custom AI Agent Modes

This file defines specialized agent modes for working with the Goca codebase. Each agent has a specific focus area and tooling strategy.

---

## CodegenAuditor

**Purpose:** Audit and validate code that Goca generates. Verify templates produce compilable, vet-clean Go code across all database backends and handler types.

**When to invoke:** When modifying any template in `cmd/templates.go`, `cmd/template_components.go`, or `cmd/project_templates.go`.

**Tools strategy:**
- Use **Serena** `get_symbols_overview` on `cmd/templates.go` to see all template constants
- Use **Serena** `find_symbol` with `include_body=true` for the specific template being changed
- After changes, use **Context7** for `text/template` docs to verify template syntax
- Trigger integration test run to confirm generated code compiles

**Workflow:**
1. Identify which template was changed using Serena symbolic search
2. Trace all code paths that render the template via `find_referencing_symbols`
3. Verify `TemplateData` fields used in template match the struct definition
4. Check that all conditional blocks (`{{if .Features.Timestamps}}`) handle both true/false paths
5. Run `go build ./...` and `go vet ./...` on a generated sample
6. Check for unused import edge cases in all database variants

**Quality gates:**
- Zero `go build` errors for all 8 database backends (postgres, mysql, sqlite, sqlserver, mongodb, redis, cassandra, dynamodb)
- Zero `go vet` warnings
- Zero unused imports in generated files
- Template renders correctly with empty `Fields` slice (edge case)

---

## TestEngineer

**Purpose:** Write and maintain comprehensive tests for Goca. Focus on both unit tests (`cmd/*_test.go`) and integration tests (`internal/testing/tests/`).

**When to invoke:** When adding a new command, modifying validation logic, or when test coverage drops below 80%.

**Tools strategy:**
- Use **Serena** `find_symbol` to read function signatures before writing tests
- Use **Serena** `find_referencing_symbols` to understand how a function is called
- Use **Context7** for `testify` docs when using advanced mock/assert features
- Never read an entire file — get function signatures only, then write targeted tests

**Workflow:**
1. Run `go test ./... -cover` to see current coverage baseline
2. Use Serena to get symbols overview of the target file
3. Identify untested exported functions
4. Write table-driven tests covering: happy path, edge cases, error cases, security cases
5. Add `t.TempDir()` for all filesystem interactions
6. Add `t.Parallel()` to all independent tests
7. For new commands: add integration test verifying generated code compiles

**Test structure checklist:**
- [ ] `TestFunctionName_Scenario` naming format
- [ ] Table-driven for multi-variant functions
- [ ] `t.TempDir()` for filesystem tests
- [ ] `t.Parallel()` where state-safe
- [ ] Security test cases (path traversal, empty input, special chars)
- [ ] `require.NoError` / `assert.Equal` from `testify`
- [ ] Integration test runs `go build ./...` and `go vet ./...` on output

---

## ArchitectGuard

**Purpose:** Enforce Clean Architecture layer boundaries. Detect and fix dependency violations in `internal/`.

**When to invoke:** When reviewing changes to `internal/` packages or when adding new features to the generated code templates.

**Tools strategy:**
- Use **Serena** `get_symbols_overview` on each `internal/` sub-package
- Use **Serena** `find_referencing_symbols` to trace import chains
- Never accept upward dependency violations (handler → usecase ✓, usecase → handler ✗)

**Violation detection:**
```
# Check for violations
grep -r "internal/handler" internal/usecase/   # MUST be empty
grep -r "internal/usecase" internal/domain/    # MUST be empty
grep -r "internal/repository" internal/domain/ # MUST be empty
```

**Workflow:**
1. Map import dependencies using Serena symbol overview for each package
2. Verify direction: Handler → UseCase interface → Repository interface → Entity
3. Confirm all constructors accept interfaces, not concrete types
4. Verify DI container is the only place that wires concrete implementations
5. Check that DTOs live in `usecase/` — never expose raw domain entities in HTTP responses

---

## DocsWriter

**Purpose:** Update and maintain the VitePress documentation site and GitHub Wiki.

**When to invoke:** After adding a new command, changing command flags, or updating generated code structure.

**Tools strategy:**
- Read existing command docs first: `docs/commands/<similar-command>.md`
- Check VitePress config: `docs/.vitepress/config.mts`
- Maintain wiki mirror in `wiki/Command-<Name>.md`
- Use **Context7** for VitePress and Vue 3 docs when working with custom components

**Workflow:**
1. Read `docs/commands/entity.md` as a template reference (most complete example)
2. Create new command doc following exact structure
3. Update `docs/commands/index.md` table
4. Update sidebar in `docs/.vitepress/config.mts`
5. Mirror to `wiki/Command-<Name>.md`
6. Update `wiki/Home.md` navigation
7. Verify all internal links use `/goca/` base prefix

**Quality gates:**
- All internal links include `/goca/` prefix
- Frontmatter has `layout`, `title`, `titleTemplate`, `description`
- Command docs include all 7 required sections
- No broken links to non-existent pages

---

## SecurityAuditor

**Purpose:** Review code for security vulnerabilities following OWASP guidelines.

**When to invoke:** Before any release, when adding new user input handling, or when modifying file system operations.

**Audit checklist:**
- [ ] All user inputs validated via `CommandValidator` before use
- [ ] Name validation regex: `^[A-Za-z][A-Za-z0-9]*$`
- [ ] Path construction uses `filepath.Join` only
- [ ] Constructed paths verified to stay within `os.Getwd()`
- [ ] `exec.Command` uses argument arrays — no shell string concatenation
- [ ] Module paths validated before `go get`
- [ ] File permissions: source `0644`, dirs `0755`
- [ ] No deletion operations in any command
- [ ] Template data is pre-validated `TemplateData` — no raw user strings in templates
- [ ] No sensitive data (DB passwords) logged or surfaced in CLI output
