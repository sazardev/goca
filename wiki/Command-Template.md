# goca template

Manage Goca's custom code templates. List, add, edit, and reset the templates used by all generation commands.

## Syntax

```bash
goca template <subcommand> [flags]
```

## Subcommands

| Subcommand | Description |
|-----------|-------------|
| `list` | List all available templates |
| `show <name>` | Display the content of a template |
| `reset` | Reset all templates to built-in defaults |

## Examples

```bash
# List all templates
goca template list

# Show the entity template
goca template show entity

# Reset templates to defaults
goca template reset
```

## See Also

- Full docs → https://sazardev.github.io/goca/commands/template
