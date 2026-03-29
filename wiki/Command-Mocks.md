# goca mocks

Generate [testify/mock](https://pkg.go.dev/github.com/stretchr/testify/mock) stubs for an entity's interfaces (repository, use case, handler).

## Syntax

```bash
goca mocks <EntityName> [flags]
```

`<EntityName>` is **required**. It is the entity name (e.g. `User`, `Product`), not an interface name.

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | `bool` | `false` | Generate all mocks (repository, use case, handler) |
| `--repository` | `bool` | `false` | Generate only the repository mock |
| `--usecase` | `bool` | `false` | Generate only the use-case mock |
| `--handler` | `bool` | `false` | Generate only the handler mock |
| `--dry-run` | `bool` | `false` | Preview without writing |
| `--force` | `bool` | `false` | Overwrite existing files |
| `--backup` | `bool` | `false` | Back up existing files before overwriting |

## Examples

```bash
# Generate all mocks for the User entity
goca mocks User --all

# Generate only the repository mock for Product
goca mocks Product --repository

# Generate repository + use-case mocks for Order
goca mocks Order --repository --usecase

# Dry-run preview
goca mocks User --all --dry-run
```

## Generated Files

```
internal/mocks/
└── mock_<entity>_repository.go
└── mock_<entity>_usecase.go
└── mock_<entity>_handler.go
```

## See Also

- [goca interfaces](Command-Interfaces) — generate the interfaces first
- [goca feature](Command-Feature) — full feature generation
- Full docs → https://sazardev.github.io/goca/commands/mocks
