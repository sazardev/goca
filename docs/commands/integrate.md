# goca integrate

Integrate existing features with dependency injection and routing.

## Syntax

```bash
goca integrate [flags]
```

## Description

Automatically detects all features in your project and integrates them by updating the dependency injection container and registering routes.

::: tip Auto-Integration
The `goca feature` command now includes automatic integration. Use this command when you need to repair or update integration for manually created features.
:::

## Flags

### `--all`

Integrate all detected features.

```bash
goca integrate --all
```

### `--feature`

Integrate a specific feature.

```bash
goca integrate --feature User
```

### `--dry-run`

Show what would be integrated without making changes.

```bash
goca integrate --all --dry-run
```

## Examples

### Integrate All Features

```bash
goca integrate --all
```

**Output:**
```
 Scanning for features...
 Found: User
 Found: Product
 Found: Order

 Updating dependency injection...
 Updated internal/di/container.go

 Registering routes...
 Updated internal/handler/http/routes.go

 Integration complete! 3 features integrated.
```

### Integrate Specific Feature

```bash
goca integrate --feature Product
```

### Dry Run

```bash
goca integrate --all --dry-run
```

Shows what would be changed without actually modifying files.

## What Gets Integrated

1. **Dependency Injection**
   - Adds repositories to DI container
   - Wires use cases with dependencies
   - Registers handlers

2. **HTTP Routing**
   - Registers all HTTP routes
   - Sets up middleware
   - Configures path prefixes

3. **Database Migrations**
   - Creates migration files
   - Registers schema changes

## Use Cases

### After Manual Feature Creation

If you created features manually:

```bash
# You manually created files for Order feature
goca integrate --feature Order
```

### Project Repair

If DI or routes are out of sync:

```bash
goca integrate --all
```

### After Git Merge

After merging branches with new features:

```bash
goca integrate --all
```

## See Also

- [`goca feature`](/commands/feature) - Generate complete feature (auto-integrated)
- [`goca di`](/commands/di) - Generate DI container
