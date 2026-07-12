# Goca Safety & Dependency Management Features

## Overview

Goca's safety and dependency management system protects your project from accidental overwrites and automates parts of `go.mod` maintenance. The `--dry-run`, `--force`, and `--backup` flags are wired through every file-generating command.

## Features

### 1. Dry-Run Mode (`--dry-run`)

Preview all changes before they are made to your project.

**Usage (any file-generating command):**
```bash
goca feature User --fields "name:string,email:string" --dry-run
goca entity Product --fields "Name:string,Price:float64" --dry-run
goca usecase Order --dry-run
goca init myapp --module github.com/acme/myapp --dry-run
```

**Output (styled terminal):**
```
DRY-RUN PREVIEW
┌──────────────────────────┬─────────┬────────┐
│ File                     │ Action  │ Size   │
├──────────────────────────┼─────────┼────────┤
│ internal/domain/user.go  │ create  │ 1234 B │
└──────────────────────────┴─────────┴────────┘
ℹ 1 files would be written
⚠ 1 conflicts detected:
  - internal/domain/user.go

Run without --dry-run to actually create files
Use --force to overwrite existing files
Use --backup to backup files before overwriting
```

In non-interactive/plain mode (`--no-color`, or when output isn't a TTY), the summary falls back to plain text:
```
DRY-RUN SUMMARY:
   Would create 1 files
```

### 2. File Conflict Detection

Automatically detects existing files and prevents accidental overwrites.

**Scenarios:**

#### Scenario A: File Exists (No Force)
```bash
goca feature User --fields "name:string"
```
```
file already exists: internal/domain/user.go (use --force to overwrite or --backup to backup first)
```

#### Scenario B: Force Overwrite
```bash
goca feature User --fields "name:string" --force
```
```
Overwriting: internal/domain/user.go
✓ Created: internal/domain/user.go
```

#### Scenario C: Backup Before Overwrite
```bash
goca feature User --fields "name:string" --force --backup
```
```
Backed up: internal/domain/user.go -> .goca-backup/user.go.backup
✓ Created: internal/domain/user.go
```

Note: backups are flattened into `.goca-backup/` by filename only (`.goca-backup/user.go.backup`), not mirrored under their original directory path — if two entities share a base filename in different directories, later backups can overwrite earlier ones.

### 3. Name Conflict Detection

Detects duplicate entity/feature names across the project.

**Example:**
```bash
# User feature already exists
goca feature User --fields "email:string"
```
```
Error: feature 'User' already exists in the project
```

**Existing Entities Detection:**
`NameConflictDetector.ScanExistingEntities()` scans `internal/domain/` for `.go` files (skipping `_test.go`, `errors.go`, `validations.go`, `common.go`) and registers each remaining basename — with any `_seeds` suffix stripped — as both an existing entity and an existing feature name, to catch collisions from either `goca entity` or `goca feature`.

### 4. Dependency Management

Adds required dependencies to `go.mod` when generating features, and can tidy the module afterward.

**Example:**
```bash
goca feature Auth --fields "username:string,password:string" --validation
```
```
✓ Added dependency: github.com/go-playground/validator/v10 v10.16.0
✓ Added dependency: github.com/golang-jwt/jwt/v5 v5.2.0

OPTIONAL DEPENDENCIES:
   The following dependencies might be useful for your feature:

   golang.org/x/crypto v0.17.0
      Reason: password hashing
      Install: go get golang.org/x/crypto@v0.17.0
```

`goca feature` runs `go mod tidy` automatically after adding its dependencies. Other generators that add a single driver dependency (e.g. `goca repository --database dynamodb`) add it to `go.mod` but do **not** run `go mod tidy` for you — run it yourself if you want the module fully tidied.

### 5. Version Compatibility Checking

`DependencyManager.CheckGoVersion(required)` compares the installed `go version` against a minimum and returns an error if it's too old — it does not print a success banner when the check passes; you'll only see output if the version check fails.

### 6. Optional Dependency Suggestions

Suggests dependencies based on feature characteristics (validation, auth, testing, gRPC, etc.) via `DependencyManager.SuggestDependencies()` / `PrintDependencySuggestions()`.

**Dependency Categories:**

#### Validation
```
github.com/go-playground/validator/v10
   Reason: struct validation for DTOs
```

#### Authentication
```
github.com/golang-jwt/jwt/v5
   Reason: JWT authentication

golang.org/x/crypto
   Reason: password hashing
```

#### Testing
```
github.com/stretchr/testify
   Reason: testing assertions and mocks
```

#### gRPC
```
google.golang.org/grpc
   Reason: gRPC protocol support

google.golang.org/protobuf
   Reason: Protocol Buffers
```

## Implementation Details

### Key Files

1. **`cmd/safety.go`**
   - `SafetyManager`: handles dry-run, force, and backup modes; default backup directory is `.goca-backup`
   - `NameConflictDetector`: scans for existing entities/features
   - File conflict detection logic

2. **`cmd/dependency_manager.go`**
   - `DependencyManager`: manages `go.mod` updates (`AddDependency`, `UpdateGoMod`)
   - Dependency suggestion system (`SuggestDependencies`, `GetRequiredDependenciesForFeature`, `GetRequiredDependenciesForDatabase`)
   - Go version compatibility checking (`CheckGoVersion`)

### Updated Files

All file-generating commands (`cmd/entity.go`, `cmd/usecase.go`, `cmd/repository.go`, `cmd/handler.go`, `cmd/di.go`, `cmd/messages.go`, `cmd/interfaces.go`, `cmd/mocks.go`, `cmd/init.go`, `cmd/integrate.go`, `cmd/feature.go`, `cmd/test-integration`) thread a `*SafetyManager` through their sub-generators and register `--dry-run`, `--force`, `--backup` flags.

## Command Flag Reference

| Command | `--dry-run` | `--force` | `--backup` |
|---|---|---|---|
| `goca entity` | ✓ | ✓ | ✓ |
| `goca usecase` | ✓ | ✓ | ✓ |
| `goca repository` | ✓ | ✓ | ✓ |
| `goca handler` | ✓ | ✓ | ✓ |
| `goca di` | ✓ | ✓ | ✓ |
| `goca messages` | ✓ | ✓ | ✓ |
| `goca interfaces` | ✓ | ✓ | ✓ |
| `goca mocks` | ✓ | ✓ | ✓ |
| `goca init` | ✓ | ✓ | ✓ |
| `goca integrate` | ✓ | ✓ | ✓ |
| `goca feature` | ✓ | ✓ | ✓ |
| `goca test-integration` | ✓ | ✓ | ✓ |

| Flag        | Type | Description                             |
| ----------- | ---- | --------------------------------------- |
| `--dry-run` | bool | Preview changes without creating files  |
| `--force`   | bool | Overwrite existing files without asking |
| `--backup`  | bool | Backup files before overwriting         |

### Safety Workflow

```bash
# 1. Preview changes first
goca feature Product --fields "name:string,price:float64" --dry-run

# 2. If satisfied, generate for real
goca feature Product --fields "name:string,price:float64"

# 3. If files exist and you want to update
goca feature Product --fields "name:string,price:float64" --force --backup
```

## Benefits

### For Developers
- **Safety**: preview changes before committing
- **Confidence**: know exactly what will be created
- **Fewer accidents**: automatic conflict detection
- **Easy recovery**: automatic backups
- **Less manual work**: dependencies added for you on `goca feature`

### For Teams
- **Consistency**: standardized dependency versions
- **Onboarding**: dependency suggestions help new developers

## Examples

### Example 1: Safe Feature Generation

```bash
# Step 1: Preview
goca feature Order --fields "customer_id:int,total:float64,status:string" --dry-run

# Step 2: Generate with backup
goca feature Order --fields "customer_id:int,total:float64,status:string" --backup

# Step 3: Dependencies auto-added and go.mod tidied
```

### Example 2: Update Existing Feature

```bash
# Backup and force update
goca feature User --fields "name:string,email:string,age:int,role:string" --force --backup

# Old files saved to .goca-backup/
# New files generated
```

### Example 3: Team Workflow

```bash
# Developer A: Preview changes
goca feature Payment --fields "amount:float64,method:string" --dry-run

# Share preview output in PR

# Developer B: Generate with exact same command
goca feature Payment --fields "amount:float64,method:string"
```

## Best Practices

### 1. Dry-Run First on Existing Projects
```bash
goca feature NewFeature --fields "..." --dry-run
goca feature NewFeature --fields "..."
```

### 2. Use Backup for Updates
```bash
goca feature ExistingFeature --fields "..." --force --backup
```

### 3. Review Dependency Suggestions
```bash
# After generation, review suggested dependencies
go get github.com/suggested/package@version
```

### 4. Handle Backups Deliberately
```bash
# Either ignore them...
echo ".goca-backup/" >> .gitignore

# ...or commit them for safety
git add .goca-backup/
git commit -m "Backup before updating User feature"
```

## Troubleshooting

### Issue: "file already exists"
**Solution:** Use `--force` and `--backup`:
```bash
goca feature User --fields "..." --force --backup
```

### Issue: "feature already exists"
**Solution:** Either:
1. Use a different name
2. Use `--force` to regenerate
3. Delete existing feature files first

### Issue: Dependency or version check fails
**Solution:** Run manually:
```bash
cd your-project
go mod tidy
go mod verify
```

### Issue: Dry-run shows many conflicts
**Solution:** This is expected if updating an existing feature. Use `--force --backup` to proceed safely.

## Contributing

These features are open for community contribution. See:
- `cmd/safety.go` - Safety manager implementation
- `cmd/dependency_manager.go` - Dependency management
- `cmd/safety_test.go`, `cmd/dependency_manager_test.go` - Tests

## Support

- Documentation: https://sazardev.github.io/goca
- Issues: https://github.com/sazardev/goca/issues
- Discussions: https://github.com/sazardev/goca/discussions
