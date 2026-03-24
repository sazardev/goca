---
applyTo: "docs/**"
---

# VitePress Documentation Rules for Goca

## Base Path Requirement

The VitePress site is deployed with `base: '/goca/'`. This affects ALL links:

```typescript
// docs/.vitepress/config.mts
export default defineConfig({
  base: "/goca/",
  // ...
});
```

### Internal Link Rules

```markdown
<!-- CORRECT — includes /goca/ base prefix -->

[Entity Command](/goca/commands/entity)
[Installation Guide](/goca/guide/installation)

<!-- FORBIDDEN — missing base prefix, causes 404 -->

[Entity Command](/commands/entity)
[Installation Guide](../guide/installation)
```

## Frontmatter Requirements

### Command Documentation Pages (`docs/commands/*.md`)

```markdown
---
layout: doc
title: goca <command>
titleTemplate: Commands | Goca
description: <One sentence describing what this command generates>
---
```

### Blog Articles (`docs/blog/articles/*.md`)

```markdown
---
layout: doc
title: <Article Title>
titleTemplate: Articles | Goca Blog
description: <SEO-friendly description, 120-160 chars>
---

<script setup>
import Badge from '../../.vitepress/theme/components/Badge.vue'
</script>
```

### Release Notes (`docs/blog/releases/v*.md`)

```markdown
---
layout: doc
title: v<semver> Release Notes
titleTemplate: Releases | Goca Blog
description: <Brief summary of key changes>
---

<script setup>
import Badge from '../../.vitepress/theme/components/Badge.vue'
</script>
```

## Command Documentation Structure

Every `docs/commands/<command>.md` MUST include these sections in order:

```markdown
# goca <command>

<One paragraph description of what gets generated and why it's useful.>

## Syntax

\`\`\`bash
goca <command> <required-arg> [flags]
\`\`\`

## Description

<Expanded description covering the Clean Architecture layer this addresses.>

## Flags

### `--flag-name`

<What it does, default value, example values.>

## Examples

### Basic usage

\`\`\`bash
goca <command> MyEntity
\`\`\`

### With options

\`\`\`bash
goca <command> MyEntity --fields "Name:string,Price:float64" --database postgres
\`\`\`

## Generated Files

| File                          | Description   |
| ----------------------------- | ------------- |
| `internal/domain/<entity>.go` | Domain entity |
| ...                           | ...           |

## Generated Code Example

<Show a realistic generated output snippet>

## Related Commands

- [`goca feature`](/goca/commands/feature) — Generate all layers at once
- [`goca entity`](/goca/commands/entity) — Domain entity only
```

## Commands Index Table (`docs/commands/index.md`)

Every new command MUST be added to the table in `docs/commands/index.md`:

```markdown
| Command                                  | Description        | Generates                              |
| ---------------------------------------- | ------------------ | -------------------------------------- |
| [`goca init`](/goca/commands/init)       | Project scaffold   | Full project structure                 |
| [`goca feature`](/goca/commands/feature) | All layers at once | Entity + UseCase + Repo + Handler + DI |
| ...                                      | ...                | ...                                    |
```

## Code Block Language Tags

Always specify language for syntax highlighting:

- Go code: ` ```go `
- Shell commands: ` ```bash `
- YAML: ` ```yaml `
- File trees: ` ```text ` or ` ```  ` (no tag)

## Blog Article Guidelines

- Articles go in `docs/blog/articles/`
- Filename: `kebab-case-title.md`
- Must have OG image if the site supports it (reference `public/og-images/`)
- Import `Badge` component for feature labels

## Wiki Mirroring

The `wiki/` directory mirrors `docs/` for the GitHub Wiki. When updating docs:

1. Update `docs/commands/<command>.md`
2. Mirror key changes to `wiki/Command-<Command>.md`
3. Update both `docs/commands/index.md` and `wiki/Home.md` navigation

## VitePress Component Usage

```vue
<!-- Badge for feature labels -->
<Badge type="new" text="New in v1.x" />
<Badge type="tip" text="Recommended" />
<Badge type="warning" text="Breaking change" />

<!-- FeatureCard for the home page -->
<FeatureCard
  icon="mdi:lightning-bolt"
  title="Feature Name"
  description="What it does"
/>
```

## Navigation Updates

When adding a new section, update `docs/.vitepress/config.mts`:

- `nav` array for top navigation
- `sidebar` object for the relevant section
- Use relative paths (VitePress prepends the base)
