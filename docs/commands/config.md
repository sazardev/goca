---
layout: doc
title: goca config
titleTemplate: Commands | Goca
description: Manage .goca.yaml configuration files — initialize, validate, show, and apply templates for consistent code generation.
---

# goca config

Manage `.goca.yaml` configuration files for your project.

## Syntax

```bash
goca config [subcommand] [flags]
```

## Description

The `goca config` command provides tools to initialize, display, validate, and manage your project's `.goca.yaml` configuration file. When run without a subcommand, it defaults to `show`.

::: tip Configuration First
Running `goca config init` before generating features ensures consistent code generation settings across your team.
:::

## Subcommands

### `show`

Display the current project configuration loaded from `.goca.yaml`.

```bash
goca config show
```

### `init`

Initialize a new `.goca.yaml` configuration file in the current directory with intelligent defaults based on your project structure.

```bash
goca config init [flags]
```

**Flags:**

| Flag | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| `--template` | string | `""` | Use predefined template (`web`, `api`, `microservice`, `full`) |
| `--force` | bool | `false` | Overwrite existing configuration file |
| `--database` | string | `""` | Database type (`postgres`, `mysql`, `sqlite`) |
| `--handlers` | strings | `[]` | Handler types (`http`, `grpc`, `cli`) |

::: warning `.goca.yaml` schema differs from `goca init`'s
`goca config init` and `goca init --config` write two different `.goca.yaml` shapes — `goca config init`'s output is missing the `templates`, `testing`, `features` and `deploy` sections that `goca init`'s does include. In particular, this means `goca template init` (which reads `templates.directory` from the config) silently does nothing useful against a `goca config init`-generated file — it reports success but creates no templates. If you plan to use `goca template`, generate your `.goca.yaml` via `goca init --config` instead of `goca config init`.
:::

### `validate`

Validate the current `.goca.yaml` configuration file for errors and warnings.

```bash
goca config validate
```

### `template`

Show available predefined configuration templates for different project types.

```bash
goca config template
```

## Usage Examples

### Initialize with defaults

```bash
goca config init
```

### Initialize from a template

```bash
goca config init --template api --database postgres
```

The `project.module` value written is read from your `go.mod`'s `module` declaration (falling back to a placeholder only if there's no `go.mod` yet).

### Validate existing config

```bash
goca config validate
```

### View current settings

```bash
goca config
# or
goca config show
```

## Configuration File Format

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

## Related Commands

- [`goca init`](/commands/init) — Initialize a new project (can generate `.goca.yaml`)
- [`goca doctor`](/commands/doctor) — Verify project health including config
- [`goca template`](/commands/template) — Manage custom code templates
