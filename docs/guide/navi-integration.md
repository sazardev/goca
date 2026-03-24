---
layout: doc
title: Navi Integration
titleTemplate: Guide | Goca
description: Use goca with Navi CLI cheat sheets for fast command lookup and interactive snippet execution with auto-completion for all flags and options.
---

# Navi Integration

[Navi](https://github.com/denisidoro/navi) is an interactive command-line cheat sheet tool. Goca ships a `.cheat` file in the repository root that covers all commands, flags, field types, and common workflows.

---

## What You Get

The `goca.cheat` file provides **133 snippets** across 19 sections:

| Section | Examples |
|---|---|
| Info & discovery | `goca --help`, `goca version --short` |
| Doctor & upgrade | `goca doctor --fix`, `goca upgrade --update` |
| Project init | Wizard, all database types, API styles |
| Feature generation | Full-stack with fields, validation, soft-delete |
| Entity | With timestamps, business rules, all Go types |
| UseCase | With DTOs, pagination, validation |
| Repository | GORM, interfaces-only |
| Handler | HTTP, gRPC, CLI, Worker |
| DI container | Per-database wiring |
| Integrate | Wire existing features into DI + routing |
| Interfaces | Generate interface contracts |
| Messages | Error message constants |
| Mocks | testify/mock generation |
| Integration tests | Test scaffolding |
| Config management | Show/generate/validate `.goca.yaml` |
| Custom templates | Register and use custom templates |
| Global flags | `--dry-run`, `--force`, `--backup`, `--quiet`, `--verbose`, `--no-color` |
| Field type examples | All supported Go types with GORM tags |
| Workflows | Real end-to-end patterns |

Variables are presented as interactive selectors. For example, the `$database` variable shows a menu of `postgres`, `mysql`, `sqlite`, `mongodb`, `sqlserver`, `elasticsearch`, `dynamodb`, and `postgres-json`.

---

## Installation

### Option 1: `navi repo add` (recommended)

```bash
navi repo add sazardev/goca
```

> **Note:** Do not prefix the URL with `github.com/` — Navi prepends it automatically.

After adding, open navi and filter by typing `goca`:

```bash
navi
# type: goca
```

### Option 2: Symlink from a local clone

```bash
git clone https://github.com/sazardev/goca.git
ln -s $(pwd)/goca/goca.cheat ~/.local/share/navi/cheats/goca.cheat
```

### Option 3: Copy the file manually

```bash
curl -o ~/.local/share/navi/cheats/goca.cheat \
  https://raw.githubusercontent.com/sazardev/goca/master/goca.cheat
```

---

## Usage

### Browse all goca snippets

```bash
navi --query goca
```

### Run a specific snippet interactively

```bash
navi
# type: feature
# select: goca feature with all fields
# fill in: Entity, Field list, Database
```

### Use as a one-liner reference

```bash
navi --print --query "goca entity"
# outputs the command text without executing it
```

---

## Example Snippets

### Generate a full feature

```
% goca, go, clean-architecture, codegen

# Generate full Clean Architecture feature with all layers
goca feature <entity> --fields "<fields>" --database <database> --validation --timestamps
```

Variables are filled in interactively:
- `<entity>` — free text input (e.g. `Product`)
- `<fields>` — free text input (e.g. `Name:string,Price:float64,Stock:int`)
- `<database>` — selector: `postgres / mysql / sqlite / mongodb / sqlserver / elasticsearch / dynamodb / postgres-json`

### Safe generation workflow

```bash
# Preview first
goca feature Order --fields "UserID:uint,Total:float64" --dry-run

# Then apply with backup
goca feature Order --fields "UserID:uint,Total:float64" --backup
```

The cheat file includes both steps as separate snippets in the **WORKFLOWS** section.

---

## Cheat File Location

The file lives at the repository root: [`goca.cheat`](https://github.com/sazardev/goca/blob/master/goca.cheat)

It is plain text and can be edited or extended locally to add project-specific snippets.

---

## Related

- [Safety & Dependency Management](/features/safety-and-dependencies) — `--dry-run`, `--force`, `--backup` flags
- [Commands Overview](/commands/) — full command reference
- [Navi GitHub repository](https://github.com/denisidoro/navi) — Navi documentation and installation
