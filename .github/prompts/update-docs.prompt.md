---
mode: agent
description: Update VitePress documentation and GitHub Wiki after adding or modifying a Goca command. Ensures docs match implementation exactly.
tools:
  - mcp_oraios_serena_find_symbol
  - mcp_oraios_serena_get_symbols_overview
  - mcp_context7_resolve-library-id
  - mcp_context7_get-library-docs
  - read_file
  - create_file
  - replace_string_in_file
---

# Update Documentation

Update all documentation pages for a changed or new Goca command.

## Target Command

`$COMMAND` — the command being documented (e.g., `entity`, `feature`, `repository`).

## Pre-work: Understand the Command

Before writing docs, read the implementation:

```
mcp_oraios_serena_get_symbols_overview("cmd/$COMMAND.go")
mcp_oraios_serena_find_symbol("${COMMAND}Cmd", include_body=false)
```

List all flags:

```
mcp_oraios_serena_find_symbol("init$COMMAND", include_body=true)  ← where flags are registered
```

## Step 1 — Create or Update `docs/commands/$COMMAND.md`

Read existing similar command for reference:

```
read_file("docs/commands/entity.md")  ← best reference
```

Required structure:

```markdown
---
layout: doc
title: goca $COMMAND
titleTemplate: Commands | Goca
description: <one-sentence description>
---

# goca $COMMAND

<opening paragraph>

## Syntax

\`\`\`bash
goca $COMMAND <name> [flags]
\`\`\`

## Description

<expanded description>

## Flags

### `--fields`

Comma-separated list of fields in `Name:Type` format.
...

## Examples

### Basic usage

\`\`\`bash
goca $COMMAND MyEntity
\`\`\`

### Complete example

\`\`\`bash
goca $COMMAND MyEntity --fields "Name:string,Price:float64" --database postgres
\`\`\`

## Generated Files

| File          | Location           | Description   |
| ------------- | ------------------ | ------------- |
| `<entity>.go` | `internal/domain/` | Domain entity |
| ...           | ...                | ...           |

## Generated Code Example

\`\`\`go
// Example of what gets generated
\`\`\`

## Related Commands

- [`goca feature`](/goca/commands/feature) — Generate all layers at once
```

**Critical:** All internal links MUST use `/goca/` prefix (base path).

## Step 2 — Update `docs/commands/index.md`

Add a row to the commands table (keep alphabetical order by command name):

```markdown
| [`goca $COMMAND`](/goca/commands/$COMMAND) | Short description | Generated components |
```

## Step 3 — Update VitePress Sidebar

Edit `docs/.vitepress/config.mts` to add the command to the Commands sidebar:

```typescript
{
  text: 'Commands',
  items: [
    // ... existing items ...
    { text: 'goca $COMMAND', link: '/commands/$COMMAND' },
  ]
}
```

Note: VitePress prepends the base path (`/goca/`) automatically to sidebar links — do NOT add `/goca/` prefix in sidebar config.

## Step 4 — Create Wiki Mirror `wiki/Command-$COMMAND.md`

The wiki uses emoji headers and a slightly different structure:

```markdown
# goca $COMMAND Command

<Same description as docs>

## 📋 Syntax

\`\`\`bash
goca $COMMAND <name> [flags]
\`\`\`

## 🎯 Purpose

<Purpose bullets>

## 🔧 Flags

...

## 📁 Generated Files

...

## 💡 Examples

...

## 🔗 Related Commands

...
```

## Step 5 — Update `wiki/Home.md` Navigation

Add entry to the appropriate section in the wiki home page navigation table.

## Step 6 — Verify All Links

```bash
# Check for broken internal links (must include /goca/ prefix)
grep -r "\](/commands/" docs/  # should find ZERO results (all should be /goca/commands/)
grep -r "\](/guide/" docs/     # should find ZERO results

# Verify frontmatter exists
head -10 docs/commands/$COMMAND.md  # should show layout, title, titleTemplate, description
```

## Acceptance Criteria

- [ ] `docs/commands/$COMMAND.md` exists with all 7 required sections
- [ ] Frontmatter has `layout`, `title`, `titleTemplate`, `description`
- [ ] All internal links include `/goca/` prefix
- [ ] `docs/commands/index.md` table updated
- [ ] `docs/.vitepress/config.mts` sidebar updated
- [ ] `wiki/Command-$COMMAND.md` created
- [ ] `wiki/Home.md` updated
- [ ] No broken link patterns (missing `/goca/` prefix)
