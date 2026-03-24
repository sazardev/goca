# goca upgrade

Upgrade your project's Goca configuration to the current binary version.

## Syntax

```bash
goca upgrade [flags]
```

## Description

The `goca upgrade` command reads your `.goca.yaml` configuration and compares it
with the schema supported by the installed Goca version:

- Reports which config sections are configured and which are at defaults
- Detects a version mismatch between the `goca_version` stored in your metadata
  and the installed binary version
- Optionally records the current version with `--update`
- Optionally re-runs code generation for a feature with `--regenerate`

::: info Note
`goca upgrade` does **not** modify your Go source files. It only manages the
`.goca.yaml` configuration. Use `goca feature <name> --force` to regenerate
source code for a specific feature.
:::

## Flags

### `--update`

Write the current Goca binary version into `.goca.yaml` under
`project.metadata.goca_version`.

```bash
goca upgrade --update
```

### `--regenerate <feature>`

Print instructions for regenerating boilerplate for a named feature.

```bash
goca upgrade --regenerate User
```

### `--dry-run`

Preview any changes that `--update` or `--regenerate` would make, without
writing anything to disk.

```bash
goca upgrade --update --dry-run
```

## Usage Examples

### Check if config is current

```bash
goca upgrade
```

**Example output (up to date):**

```
Goca Upgrade

Project: myproject
Module:  github.com/user/myproject

ℹ Installed Goca version : v1.18.2
ℹ Recorded goca_version  : v1.18.2 (up to date)

Config Section Status
┌─────────────────────────────────────────────────┐
│  Section       Status      Note                 │
├─────────────────────────────────────────────────┤
│  project       ✓ set       name, module         │
│  architecture  ○ default   layers, DI type      │
│  database      ○ default   type, host           │
│  generation    ○ default   validation, style    │
│  testing       ○ default   framework, mocks     │
│  features      ○ default   auth, cache          │
│  templates     ○ default   custom dir           │
│  deploy        ○ default   docker, kubernetes   │
└─────────────────────────────────────────────────┘

✓ Project configuration is up to date
```

### Record the installed version

```bash
goca upgrade --update
```

This writes `goca_version: v1.18.2` into `.goca.yaml` metadata and is
useful after upgrading Goca itself.

### Regenerate a feature's boilerplate

```bash
goca upgrade --regenerate User
```

Prints the exact `goca feature` command to run with `--force`.

### Preview --update in dry-run mode

```bash
goca upgrade --update --dry-run
```

Shows what would be written without touching the file.

## When to Run `goca upgrade`

- After running `go install github.com/sazardev/goca@latest`
- When `goca doctor` reports a version mismatch
- Before regenerating features to ensure config is consistent

## Related Commands

| Command | Purpose |
| ------- | ------- |
| `goca doctor` | Full project health check |
| `goca feature <name> --force` | Regenerate all boilerplate for a feature |
| `goca init` | Initialize a new `.goca.yaml` |
