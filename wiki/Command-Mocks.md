# goca mocks

Generate [testify/mock](https://pkg.go.dev/github.com/stretchr/testify/mock) stubs for all interfaces in your project.

## Syntax

```bash
goca mocks [flags]
goca mocks <InterfaceName> [flags]
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--output` | `string` | `internal/mocks/` | Output directory |
| `--dry-run` | `bool` | `false` | Preview without writing |
| `--force` | `bool` | `false` | Overwrite existing files |

## Examples

```bash
# Generate mocks for all interfaces
goca mocks

# Generate a mock for a specific interface
goca mocks UserRepository

# Dry-run preview
goca mocks --dry-run
```

## Generated Files

```
internal/mocks/
└── mock_<interface_name>.go
```

## See Also

- [goca interfaces](Command-Interfaces) — generate the interfaces first
- [goca feature](Command-Feature) — full feature generation
- Full docs → https://sazardev.github.io/goca/commands/mocks
