# Goca Safety & Dependency Management Features

## Overview

Goca's safety and dependency management system protects your project from accidental overwrites and automates `go.mod` maintenance. As of **v1.18.7**, the `--dry-run`, `--force`, and `--backup` flags are fully wired through **all 12 file-generating commands**.

## New Features

### 1. Dry-Run Mode (`--dry-run`)

Preview all changes before they are made to your project.

**Usage (any file-generating command):**
```bash
goca feature User --fields "name:string,email:string" --dry-run
goca entity Product --fields "Name:string,Price:float64" --dry-run
goca usecase Order --dry-run
goca init myapp --module github.com/acme/myapp --dry-run
```

**Output:**
```
🔍 DRY-RUN MODE: Previewing changes without creating files

📝 [DRY-RUN] Would create: internal/domain/user.go (1234 bytes)
📝 [DRY-RUN] Would create: internal/usecase/user_service.go (2345 bytes)
📝 [DRY-RUN] Would create: internal/repository/postgres_user_repository.go (1567 bytes)
📝 [DRY-RUN] Would create: internal/handler/http/user_handler.go (2890 bytes)

📋 DRY-RUN SUMMARY:
   Would create 15 files
   ⚠️  2 conflicts detected:
      - internal/domain/user.go
      - internal/usecase/user_service.go

💡 Run without --dry-run to actually create files
   Use --force to overwrite existing files
   Use --backup to backup files before overwriting
```

### 2. File Conflict Detection

Automatically detects existing files and prevents accidental overwrites.

**Scenarios:**

#### Scenario A: File Exists (No Force)
```bash
goca feature User --fields "name:string"
```
```
❌ file already exists: internal/domain/user.go (use --force to overwrite or --backup to backup first)
```

#### Scenario B: Force Overwrite
```bash
goca feature User --fields "name:string" --force
```
```
⚠️  Overwriting: internal/domain/user.go
✅ Created: internal/domain/user.go
```

#### Scenario C: Backup Before Overwrite
```bash
goca feature User --fields "name:string" --force --backup
```
```
📦 Backed up: internal/domain/user.go -> .goca-backup/internal/domain/user.go.backup
✅ Created: internal/domain/user.go
```

### 3. Name Conflict Detection

Detects duplicate entity/feature names across the project.

**Example:**
```bash
# User feature already exists
goca feature User --fields "email:string"
```
```
❌ feature 'User' already exists in the project
💡 Use --force to generate anyway
```

**Existing Entities Detection:**
The system scans `internal/domain/` for existing entities and prevents duplicates.

### 4. Automatic go.mod Management

Automatically updates `go.mod` when generating features with dependencies.

**Features:**
- ✅ Adds required dependencies automatically
- ✅ Runs `go mod tidy` after generation
- ✅ Verifies dependency compatibility
- ✅ Suggests optional dependencies

**Example:**
```bash
goca feature Auth --fields "username:string,password:string" --validation
```
```
7️⃣  Managing dependencies...
✅ Added dependency: github.com/go-playground/validator/v10 v10.16.0
✅ Added dependency: github.com/golang-jwt/jwt/v5 v5.2.0

📦 Updating go.mod...
✅ Updated go.mod and go.sum

💡 OPTIONAL DEPENDENCIES:
   The following dependencies might be useful for your feature:

   📦 golang.org/x/crypto v0.17.0
      Reason: password hashing
      Install: go get golang.org/x/crypto@v0.17.0
```

### 5. Version Compatibility Checking

Verifies Go version and dependency compatibility.

**Features:**
- ✅ Checks minimum Go version (1.21+)
- ✅ Verifies dependency versions are compatible
- ✅ Warns about potential conflicts

**Example:**
```bash
goca init myproject --module github.com/user/myproject
```
```
✅ Go version check: go1.25.2 (compatible with go1.21+)
✅ All dependencies verified
```

### 6. Optional Dependency Suggestions

Intelligently suggests dependencies based on feature characteristics.

**Dependency Categories:**

#### Validation
```
📦 github.com/go-playground/validator/v10
   Reason: struct validation for DTOs
```

#### Authentication
```
📦 github.com/golang-jwt/jwt/v5
   Reason: JWT authentication
   
📦 golang.org/x/crypto
   Reason: password hashing
```

#### Testing
```
📦 github.com/stretchr/testify
   Reason: testing assertions and mocks
   
📦 github.com/golang/mock
   Reason: mock generation for testing
```

#### gRPC
```
📦 google.golang.org/grpc
   Reason: gRPC protocol support
   
📦 google.golang.org/protobuf
   Reason: Protocol Buffers
```

## Implementation Details

### New Files Created

1. **`cmd/safety.go`**
   - `SafetyManager`: Handles dry-run, force, and backup modes
   - `NameConflictDetector`: Scans for existing entities/features
   - File conflict detection logic
   - Backup system

2. **`cmd/dependency_manager.go`**
   - `DependencyManager`: Manages go.mod updates
   - Dependency suggestion system
   - Version compatibility checking
   - Automatic dependency installation

### Updated Files

1. **All 12 generator commands** (`cmd/entity.go`, `cmd/usecase.go`, `cmd/repository.go`, `cmd/handler.go`, `cmd/di.go`, `cmd/messages.go`, `cmd/interfaces.go`, `cmd/mocks.go`, `cmd/init.go`, `cmd/integrate.go`, `cmd/feature.go`, `cmd/test_integration.go`)
   - `--dry-run`, `--force`, `--backup` flags registered on each
   - SafetyManager threaded through all sub-generators
   - `feature` and `integrate` forward SafetyManager to every generator they call

2. **`cmd/utils.go`** (v1.18.7)
   - `writeFile()` and `writeGoFile()` accept variadic `*SafetyManager` parameter
   - When provided, all writes route through `SafetyManager.WriteFile()`

## Command Flag Reference

### All Commands Support

As of v1.18.7, `--dry-run`, `--force`, and `--backup` are registered and fully functional on every file-generating command:

| Command | `--dry-run` | `--force` | `--backup` |
|---|---|---|---|
| `goca entity` | ✅ | ✅ | ✅ |
| `goca usecase` | ✅ | ✅ | ✅ |
| `goca repository` | ✅ | ✅ | ✅ |
| `goca handler` | ✅ | ✅ | ✅ |
| `goca di` | ✅ | ✅ | ✅ |
| `goca messages` | ✅ | ✅ | ✅ |
| `goca interfaces` | ✅ | ✅ | ✅ |
| `goca mocks` | ✅ | ✅ | ✅ |
| `goca init` | ✅ | ✅ | ✅ |
| `goca integrate` | ✅ | ✅ | ✅ |
| `goca feature` | ✅ | ✅ | ✅ |
| `goca test-integration` | ✅ | ✅ | ✅ |

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
✅ **Safety**: Preview changes before committing
✅ **Confidence**: Know exactly what will be created
✅ **No Accidents**: Automatic conflict detection
✅ **Easy Recovery**: Automatic backups
✅ **Less Manual Work**: Automatic dependency management

### For Teams
✅ **Consistency**: Standardized dependency versions
✅ **Documentation**: Clear what each feature requires
✅ **Onboarding**: Suggestions help new developers
✅ **Best Practices**: Automatic inclusion of common libraries

## Configuration

### .goca.yaml Support

```yaml
# Enable safety features by default
safety:
  dry_run_default: false
  backup_enabled: true
  conflict_detection: true

# Dependency management
dependencies:
  auto_update: true
  suggest_optional: true
  verify_versions: true
```

## Examples

### Example 1: Safe Feature Generation

```bash
# Step 1: Preview
goca feature Order --fields "customer_id:int,total:float64,status:string" --dry-run

# Step 2: Check for conflicts
# (automatically done)

# Step 3: Generate with backup
goca feature Order --fields "customer_id:int,total:float64,status:string" --backup

# Step 4: Dependencies auto-added
# go.mod updated automatically
```

### Example 2: Update Existing Feature

```bash
# Backup and force update
goca feature User --fields "name:string,email:string,age:int,role:string" --force --backup

# Old files saved to .goca-backup/
# New files generated
# Dependencies updated
```

### Example 3: Team Workflow

```bash
# Developer A: Preview changes
goca feature Payment --fields "amount:float64,method:string" --dry-run

# Share preview output in PR
# Team reviews

# Developer B: Generate with exact same command
goca feature Payment --fields "amount:float64,method:string"

# Consistent results across team
```

## Best Practices

### 1. Always Dry-Run First
```bash
# Good
goca feature NewFeature --fields "..." --dry-run
goca feature NewFeature --fields "..."

# Risky
goca feature NewFeature --fields "..."
```

### 2. Use Backup for Updates
```bash
# Safe
goca feature ExistingFeature --fields "..." --force --backup

# Risky
goca feature ExistingFeature --fields "..." --force
```

### 3. Review Dependency Suggestions
```bash
# After generation, review suggested dependencies
# Install only what you need
go get github.com/suggested/package@version
```

### 4. Commit Backups
```bash
# Add backups to .gitignore
echo ".goca-backup/" >> .gitignore

# Or commit them for safety
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

### Issue: "dependency verification failed"
**Solution:** Run manually:
```bash
cd your-project
go mod tidy
go mod verify
```

### Issue: Dry-run shows many conflicts
**Solution:** This is expected if updating an existing feature. Use `--force --backup` to proceed safely.

## Migration Guide

### Existing Projects

No changes needed! New features are opt-in via flags.

### Updating Commands

Old command:
```bash
goca feature User --fields "name:string"
```

New (safer):
```bash
# Preview first
goca feature User --fields "name:string" --dry-run

# Then generate
goca feature User --fields "name:string"
```

## Future Enhancements

Planned for v1.12.0:
- Interactive conflict resolution
- Merge tool for conflicting files
- Undo/rollback command
- Dependency version suggestions
- Security vulnerability scanning

## Contributing

These features are open for community contribution. See:
- `cmd/safety.go` - Safety manager implementation
- `cmd/dependency_manager.go` - Dependency management
- Tests in `internal/testing/tests/safety_test.go`

## Support

- 📚 Documentation: https://sazardev.github.io/goca
- 🐛 Issues: https://github.com/sazardev/goca/issues
- 💬 Discussions: https://github.com/sazardev/goca/discussions
