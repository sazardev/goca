# goca template

Manage Goca's custom code templates. List, add, edit, and reset the templates used by all generation commands.

## Syntax

```bash
goca template <subcommand> [flags]
```

## Subcommands

| Subcommand | Description |
|-----------|-------------|
| `init` | Initialize the `.goca/templates` directory with the built-in defaults |
| `list` | List all available templates |
| `show <name>` | Display the content of a template |
| `reset` | Reset all templates to built-in defaults |

## Examples

```bash
# List all templates
goca template list

# Show the entity template (name is the path relative to the templates dir, without extension)
goca template show domain/entity

# Reset templates to defaults
goca template reset
```

## See Also

- Full docs → https://sazardev.github.io/goca/commands/template
