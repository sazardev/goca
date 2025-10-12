# Changelog

All notable changes to Goca CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.11.0] - 2025-01-12

### 🎉 Major Features Added

#### Safety Features
- **Dry-Run Mode** (`--dry-run`): Preview all file changes before generation
  - Shows exactly what files would be created/modified
  - Displays dependency suggestions
  - Zero risk file operations
  - Perfect for CI/CD pipelines and code reviews

- **File Conflict Detection**: Automatic detection of existing files
  - Prevents accidental overwrites by default
  - Clear error messages with resolution options
  - Integrates with `--force` and `--backup` flags

- **Name Conflict Detection**: Scans project for duplicate entities
  - Case-insensitive duplicate detection
  - Scans `internal/domain/` for existing features
  - Prevents confusion with similarly named entities

- **Backup System** (`--backup`): Automatic file backups
  - Creates timestamped backups in `.goca-backup/`
  - Works with `--force` flag for safe overwrites
  - Organized directory structure preserves original paths

- **Force Overwrite** (`--force`): Override safety protections
  - Explicitly overwrite existing files
  - Combines with `--backup` for safe updates
  - Useful for intentional feature updates

#### Dependency Management
- **Automatic go.mod Updates**: Auto-manages project dependencies
  - Runs `go get` for required dependencies
  - Executes `go mod tidy` automatically
  - Validates dependency versions

- **Version Compatibility Checking**: Go version verification
  - Requires Go 1.21+ for generated projects
  - Clear error messages for incompatible versions
  - Prevents runtime issues with older Go versions

- **Smart Dependency Suggestions**: Context-aware recommendations
  - Suggests optional dependencies based on feature type
  - Provides installation commands
  - Explains why each dependency is useful
  - Common suggestions:
    - `validator/v10` for struct validation
    - `jwt/v5` for authentication features
    - `grpc` for gRPC handlers
    - `testify` for testing
    - `swagger` for API documentation

### 🔧 Files Added
- `cmd/safety.go`: Core safety infrastructure (SafetyManager, NameConflictDetector)
- `cmd/dependency_manager.go`: Dependency management system
- `internal/testing/tests/safety_test.go`: Comprehensive test suite
- `docs/features/safety-and-dependencies.md`: Complete documentation
- `SAFETY_FEATURES_IMPLEMENTATION.md`: Implementation summary

### 🔄 Files Modified
- `cmd/feature.go`: Integrated safety and dependency features
  - Added `--dry-run`, `--force`, `--backup` flags
  - Integrated SafetyManager for file operations
  - Integrated DependencyManager for go.mod updates
  - Added name conflict checking before generation
  - Enhanced output with progress indicators

- `README.md`: Updated with v1.11.0 features section
- `docs/commands/init.md`: Added Git initialization documentation

### 🐛 Bug Fixes
- None (pure feature release)

### 📚 Documentation
- Added comprehensive safety features guide
- Updated README with new features showcase
- Created implementation summary document
- Enhanced init command documentation

### 🧪 Testing
- Added unit tests for SafetyManager
- Added unit tests for NameConflictDetector
- Added unit tests for DependencyManager
- All tests cover dry-run, force, and backup scenarios

### ⚠️ Breaking Changes
None. All features are opt-in via flags.

### 📦 Migration Guide
No migration needed. New flags are optional:
- Default behavior unchanged
- `--dry-run` is purely additive
- `--force` and `--backup` only active when specified

### 🎯 What's Next (v1.12.0)
- Interactive conflict resolution
- Merge tool for conflicting files
- Undo/rollback command
- History tracking
- Plugin system for custom safety checks
- Configuration file for default flags

---

## [1.10.10] - 2025-01-10

### ✨ Features Added
- **Git Initialization**: `goca init` now automatically initializes Git repository
  - Runs `git init` after project creation
  - Creates initial commit: "Initial commit - Goca Clean Architecture project"
  - Adds all generated files to initial commit
  - Gracefully handles git unavailability

### 🔄 Files Modified
- `cmd/init.go`: Added `initializeGitRepository()` function
- `docs/commands/init.md`: Documented Git initialization feature

### 🐛 Bug Fixes
- None

---

## [1.10.9] - 2025-01-08

### 🐛 Bug Fixes
- Fixed DI container generation for multi-feature projects
- Fixed route registration order in HTTP handlers
- Improved error messages in validation logic

---

## [1.10.8] - 2025-01-05

### ✨ Features Added
- Added `goca integrate` command for automatic feature integration
- Added `--all` flag to integrate all detected features
- Auto-detection of unintegrated features

### 🔄 Files Modified
- `cmd/integrate.go`: New command implementation
- `README.md`: Updated with integrate command

---

## [1.10.0] - 2024-12-20

### 🎉 Major Release
- Initial stable release with Clean Architecture support
- Complete feature generation (entity, usecase, repository, handler)
- Multi-protocol handler support (HTTP, gRPC, CLI)
- Dependency injection container generation
- VitePress documentation site

---

## [1.0.0] - 2024-11-15

### 🎉 Initial Release
- Basic entity generation
- Repository pattern support
- HTTP handler generation
- Clean Architecture structure

---

## Legend

- 🎉 Major features
- ✨ New features
- 🔧 New files
- 🔄 Modified files
- 🐛 Bug fixes
- 📚 Documentation
- 🧪 Testing
- ⚠️ Breaking changes
- 📦 Migration
- 🎯 Future plans
