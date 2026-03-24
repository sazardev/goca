# goca doctor

Run health checks on your project to verify Clean Architecture structure and code quality.

## Syntax

```bash
goca doctor [flags]
```

## Description

The `goca doctor` command inspects the current directory and runs a series of automated health checks. It reports the status of each check in a table with actionable suggestions.

Use `goca doctor` regularly to catch structural drift, missing configuration, and build issues early.

::: tip Project Health
Running `goca doctor` after adding new features or switching branches helps ensure your Clean Architecture structure stays consistent.
:::

## Checks Performed

| Check | What It Verifies |
| ----- | ---------------- |
| `go.mod` | Module declaration is present |
| `.goca.yaml` | Configuration file exists and is non-empty |
| Clean Architecture dirs | `internal/domain`, `internal/usecase`, `internal/repository`, `internal/handler` exist |
| `go build ./...` | Project compiles without errors |
| `go vet ./...` | No static analysis warnings |
| DI container | `internal/di` directory exists |

### Status Icons

- `✓` — Check passed
- `⚠` — Warning (non-fatal, suggestion available)
- `✗` — Check failed (error, exits with non-zero code)

## Flags

### `--fix`

Automatically create missing Clean Architecture directories when the structure check fails.

```bash
goca doctor --fix
```

## Usage Examples

### Basic health check

```bash
goca doctor
```

**Example output:**

```
Goca Doctor — Project Health Check

┌────────────────────────────────────────────────────────────────────────────┐
│     Check                   Details                  Suggestion            │
├────────────────────────────────────────────────────────────────────────────┤
│  ✓  go.mod                  go.mod present            —                    │
│  ✓  .goca.yaml              .goca.yaml present        —                    │
│  ⚠  Clean Architecture dirs 2 dirs missing            goca doctor --fix    │
│  ✓  go build ./...          Compiles without errors   —                    │
│  ✓  go vet ./...            No warnings               —                    │
│  ⚠  DI container            No DI directory found     goca di              │
└────────────────────────────────────────────────────────────────────────────┘

ℹ Results: 4 passed, 2 warnings, 0 failed
```

### Auto-fix missing directories

```bash
goca doctor --fix
```

Creates any missing `internal/` layer directories automatically.

### In CI pipelines

```bash
# Fail the build if any check fails
goca doctor || exit 1
```

`goca doctor` exits with code `1` when any check returns `✗`.

## Integration with Other Commands

After `goca doctor` reports issues:

| Issue | Fix Command |
| ----- | ----------- |
| Missing directories | `goca doctor --fix` |
| No `.goca.yaml` | `goca init <project-name>` |
| Missing DI container | `goca di` |
| Build errors | Fix code, then `go build ./...` |

## Exit Codes

| Code | Meaning |
| ---- | ------- |
| `0` | All checks passed (warnings are allowed) |
| `1` | One or more checks failed |
