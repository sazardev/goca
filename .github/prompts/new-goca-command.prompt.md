---
mode: agent
description: Implement a new goca CLI command with all required layers — cobra command, tests, docs, and CHANGELOG entry.
tools:
  - mcp_oraios_serena_get_symbols_overview
  - mcp_oraios_serena_find_symbol
  - mcp_oraios_serena_find_referencing_symbols
  - mcp_oraios_serena_replace_symbol_body
  - mcp_oraios_serena_insert_after_symbol
  - mcp_context7_resolve-library-id
  - mcp_context7_get-library-docs
  - run_in_terminal
  - read_file
  - create_file
  - replace_string_in_file
---

# New Goca Command

Implement a new `goca` CLI command following the complete checklist. Never skip a step.

## Required Information

Before starting, confirm:

1. **Command name:** `$COMMAND_NAME` (e.g., `scaffold`, `model`)
2. **What it generates:** `$WHAT_IT_GENERATES`
3. **Required flags:** `$FLAGS` (e.g., `--fields`, `--database`)

## Implementation Steps

### Step 1 — Explore existing command for reference

Use Serena to read a similar command's structure:

```
mcp_oraios_serena_get_symbols_overview("cmd/entity.go")
mcp_oraios_serena_find_symbol("entityCmd", include_body=true)
```

### Step 2 — Implement `cmd/$COMMAND_NAME.go`

Follow the standard command structure:

1. `cobra.Command` with `Use`, `Short`, `Long`, `Args`
2. Parse flags
3. Load `ConfigIntegration`
4. Merge CLI flags with config
5. Initialize `SafetyManager(dryRun, force, backup)`
6. Validate inputs with `CommandValidator`
7. Build `TemplateData`
8. Write files via `safetyMgr.WriteFile()`
9. Add dependencies via `DependencyManager`
10. Print summary

### Step 3 — Register in `cmd/root.go`

```go
rootCmd.AddCommand($COMMAND_NAMECmd)
```

### Step 4 — Write unit tests `cmd/$COMMAND_NAME_test.go`

Table-driven tests covering:

- Valid inputs (happy path)
- Empty name (error)
- Invalid characters in name (error)
- Path traversal attempt (error)
- Each flag combination

### Step 5 — Write integration test `internal/testing/tests/$COMMAND_NAME_test.go`

Must verify:

1. Files are created at expected paths
2. Generated code compiles: `go build ./...`
3. Generated code passes vet: `go vet ./...`
4. Generated files contain expected signatures

### Step 6 — Create docs `docs/commands/$COMMAND_NAME.md`

Include all 7 required sections:

1. Title + description
2. Syntax block
3. Description paragraph
4. Flags section
5. Examples (basic + with options)
6. Generated Files table
7. Related Commands

### Step 7 — Update documentation index

- Add row to `docs/commands/index.md` table
- Add entry to `docs/.vitepress/config.mts` sidebar

### Step 8 — Update CHANGELOG.md

Under `[Unreleased]` → `### Added`:

```markdown
- `goca $COMMAND_NAME` — [brief description]
```

### Step 9 — Compile and vet verification

```bash
go build ./...
go vet ./...
go test ./cmd/... -run TestValidate -v
```

All must pass with zero errors, zero warnings.
