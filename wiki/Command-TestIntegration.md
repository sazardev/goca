# goca test-integration

Generate a full integration test suite for a feature, including database helpers, fixtures, and test lifecycle management.

## Syntax

```bash
goca test-integration <EntityName> [flags]
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--database` | `string` | `postgres` | Target database (postgres, mysql, sqlite) |
| `--container` | `bool` | `false` | Generate testcontainers-go setup |
| `--dry-run` | `bool` | `false` | Preview without writing |
| `--force` | `bool` | `false` | Overwrite existing files |

## Examples

```bash
# Generate integration tests for a Product entity
goca test-integration Product

# With test containers
goca test-integration Product --container

# For MySQL
goca test-integration Order --database mysql
```

## Generated Files

```
internal/testing/tests/
└── <entity>_integration_test.go
└── helpers_test.go
└── fixtures_test.go
```

## See Also

- [goca feature](Command-Feature) — generate the feature first
- Full docs → https://sazardev.github.io/goca/commands/test-integration
