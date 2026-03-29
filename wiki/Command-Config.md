# goca config Command

The `goca config` command manages `.goca.yaml` configuration files — initialize, display, validate, and use predefined templates for consistent code generation across your team.

## 📋 Syntax

```bash
goca config [subcommand] [flags]
```

Defaults to `show` when no subcommand is provided.

## 🎯 Purpose

Centralizes `.goca.yaml` lifecycle management:

- 📄 **Show** — display current loaded configuration
- 🔧 **Init** — scaffold a new `.goca.yaml` from defaults or templates
- ✅ **Validate** — check configuration for errors and warnings
- 📋 **Template** — list available predefined configuration templates

## 📝 Subcommands

### `show`

Display the current project configuration (default when no subcommand is given):

```bash
goca config
# or
goca config show
```

### `init`

Initialize a `.goca.yaml` in the current directory:

```bash
goca config init [flags]
```

| Flag | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| `--template` | string | `""` | Predefined template (`web`, `api`, `microservice`, `full`) |
| `--force` | bool | `false` | Overwrite existing config file |
| `--database` | string | `""` | Database type (`postgres`, `mysql`, `sqlite`) |
| `--handlers` | strings | `[]` | Handler types (`http`, `grpc`, `cli`) |

### `validate`

Validate the `.goca.yaml` configuration for errors and warnings:

```bash
goca config validate
```

### `template`

List available predefined configuration templates:

```bash
goca config template
```

## 📖 Usage Examples

### Initialize with defaults

```bash
goca config init
```

### Initialize from a template

```bash
goca config init --template api --database postgres
```

### Force overwrite existing config

```bash
goca config init --force
```

### Validate before generating features

```bash
goca config validate
```

## 📄 Configuration File Format

The `.goca.yaml` file controls code generation behavior:

```yaml
project:
  name: myproject
  module: github.com/user/myproject

database:
  type: postgres

handlers:
  - http

templates:
  directory: .goca/templates

naming:
  file_convention: snake_case
```

## 🔗 Related Commands

- [goca init](Command-Init) — Initialize a new project (can generate `.goca.yaml`)
- [goca doctor](Command-Doctor) — Verify project health including config validity
- [goca template](Command-Template) — Manage custom code templates
