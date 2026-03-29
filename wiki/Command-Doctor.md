# goca doctor

Run automated health checks on your project to verify that the Clean Architecture structure is intact and the codebase is wired correctly.

## Syntax

```bash
goca doctor [flags]
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--fix` | `bool` | `false` | Auto-create missing Clean Architecture directories |

## Checks Performed

| Check | Description |
|-------|-------------|
| `go.mod` present | Module file is found |
| `.goca.yaml` present | Goca config file exists |
| Directory structure | `internal/domain`, `usecase`, `repository`, `handler` exist |
| Build | `go build ./...` passes |
| Vet | `go vet ./...` passes |

## Examples

```bash
# Run all health checks
goca doctor

# Auto-fix missing directories
goca doctor --fix
```

## See Also

- [goca upgrade](Command-Upgrade) — keep `.goca.yaml` up to date
- Full docs → https://sazardev.github.io/goca/commands/doctor
