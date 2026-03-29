# goca upgrade

Check for new Goca releases and update the `.goca.yaml` configuration metadata to match the current binary version.

## Syntax

```bash
goca upgrade [flags]
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--update` | `bool` | `false` | Write the current binary version to `.goca.yaml` |

## Examples

```bash
# Check if .goca.yaml is up-to-date
goca upgrade

# Update .goca.yaml metadata to current binary version
goca upgrade --update
```

## See Also

- [goca doctor](Command-Doctor) — run project health checks
- Full docs → https://sazardev.github.io/goca/commands/upgrade
