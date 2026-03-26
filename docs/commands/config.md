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

- [`goca init`](/goca/commands/init) — Initialize a new project (can generate `.goca.yaml`)
- [`goca doctor`](/goca/commands/doctor) — Verify project health including config
- [`goca template`](/goca/commands/template) — Manage custom code templates
