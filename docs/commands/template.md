---
layout: doc
title: goca template
titleTemplate: Commands | Goca
description: Manage custom code generation templates — initialize and list templates for personalized code output.
---

# goca template

Manage custom templates for personalized code generation.

## Syntax

```bash
goca template [subcommand]
```

## Description

The `goca template` command provides tools to initialize and manage custom templates for code generation. Custom templates allow you to override Goca's built-in templates with project-specific patterns.

::: tip Customization
Custom templates let you enforce your team's coding style across all generated code while keeping the benefits of automated generation.
:::

## Subcommands

### `init`

Initialize the custom templates directory with default templates that you can customize.

```bash
goca template init
```

Creates a `.goca/templates/` directory with editable templates for all layers:
- Entity templates
- Use case templates
- Repository templates
- Handler templates
- DI container templates

After initialization, Goca will automatically use your custom templates when generating code.

### `list`

List all available custom templates in the current project.

```bash
goca template list
```

Displays the templates directory path and confirms whether templates exist.

## Usage Examples

### Set up custom templates

```bash
# Initialize templates
goca template init

# Edit to match your style
# Templates are in .goca/templates/
```

### Check available templates

```bash
goca template list
```

## Template Functions

Custom templates have access to these template functions:

| Function | Description | Example |
| -------- | ----------- | ------- |
| `{{pascal .EntityName}}` | PascalCase | `UserProfile` |
| `{{snake .EntityName}}` | snake_case | `user_profile` |
| `{{camel .EntityName}}` | camelCase | `userProfile` |
| `{{lower .EntityName}}` | lowercase | `userprofile` |
| `{{plural .EntityName}}` | Pluralized | `Users` |

## Template Data

Templates receive a `TemplateData` struct with these fields:

| Field | Type | Description |
| ----- | ---- | ----------- |
| `EntityName` | string | Entity name (PascalCase) |
| `Module` | string | Go module path |
| `Database` | string | Database type |
| `Fields` | []Field | Parsed field definitions |
| `HasValidation` | bool | Whether validation is enabled |

## Related Commands

- [`goca config`](/goca/commands/config) — Manage `.goca.yaml` configuration
- [`goca feature`](/goca/commands/feature) — Generate features using templates
- [`goca init`](/goca/commands/init) — Initialize a new project
