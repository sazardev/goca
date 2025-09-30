# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.0.0] - 2025-09-30

### Added
- **VitePress Documentation System** - Complete documentation site with beautiful UI
  - Interactive documentation at https://sazardev.github.io/goca
  - Full-text search capability
  - Dark mode support
  - Mobile-responsive design
  - Command reference, guides, and examples
  - GitHub Pages deployment workflow

### Changed
- **Major Version Bump** - Marking production-ready status with complete documentation
- **Enhanced Documentation** - 2,500+ lines of comprehensive documentation
  - BUG_FIXES_REPORT.md (18.5 KB)
  - BUG_FIXES_QUICK_REFERENCE.md (7.4 KB)
  - PRODUCTION_READY_ACHIEVEMENT_REPORT.md (12.5 KB)
  - DOCUMENTATION_INDEX.md (11.0 KB)
  - FINAL_DOCUMENTATION_SUMMARY.md
  - Complete VitePress site

### Documentation
- Created VitePress documentation system with npm package
- Added comprehensive release notes for v2.0.0
- Added getting started guide
- Added complete command reference
- Added configuration guides
- Total documentation: ~120 KB across all files

## [1.0.1] - 2025-09-30

### Fixed
- **[Bug #1]** Fixed missing GORM import in entity templates when soft-delete is enabled ([#BUG-001](docs/BUG_FIXES_REPORT.md))
  - Entities with `soft_delete: true` now automatically include `gorm.io/gorm` import
  - Fixes compilation error: `undefined: gorm.DeletedAt`
  - Modified: `cmd/entity.go` - Added conditional GORM import in `writeEntityHeader()`

- **[Bug #2]** Removed unused time import from seed file templates ([#BUG-002](docs/BUG_FIXES_REPORT.md))
  - Seed files no longer import unused `time` package
  - Eliminates compiler warning: `imported and not used: "time"`
  - Modified: `cmd/entity.go` - Simplified `writeSeedFileHeader()`

- **[Bug #3]** Fixed missing domain import when registering entities for auto-migration ([#BUG-003](docs/BUG_FIXES_REPORT.md))
  - `goca feature` now automatically adds domain package import to `main.go`
  - Fixes compilation error: `undefined: domain`
  - Added: `ensureDomainImport()` function (67 lines) - Intelligently manages imports
  - Added: `isEntityInMigrationList()` function (34 lines) - Checks entity registration correctly
  - Modified: `cmd/automigrate.go` (+101 lines total)

- **[Bug #4]** Fixed hardcoded postgres in MySQL configuration ([#BUG-004](docs/BUG_FIXES_REPORT.md))
  - `goca init --database mysql` now correctly writes `type: mysql` to `.goca.yaml`
  - Previously wrote `type: postgres` regardless of CLI flag
  - Modified: `cmd/init.go`, `cmd/config_integration.go`, `cmd/config_manager.go`
  - Database parameter now flows through complete config generation chain

- **[Bug #5]** Implemented kebab-case file naming convention ([#BUG-005](docs/BUG_FIXES_REPORT.md))
  - `files: kebab-case` config now generates correct filenames (e.g., `order-item.go`)
  - Previously fell back to lowercase (e.g., `orderitem.go`)
  - Modified: `cmd/entity.go` (1 location), `cmd/handler.go` (6 locations)
  - Supports kebab-case for: entities, HTTP handlers, gRPC files, CLI commands, workers, SOAP clients

### Testing
- Created 4 dedicated test projects to verify bug fixes:
  - `bug-fix-test`: Verified bugs #1 and #2
  - `domain-import-test`: Verified bug #3 (domain import auto-added)
  - `mysql-config-test`: Verified bug #4 (MySQL config correct)
  - `kebab-case-test`: Verified bug #5 (kebab-case naming works)
- All test projects compile successfully with zero errors
- Comprehensive testing report: [docs/BUG_FIXES_REPORT.md](docs/BUG_FIXES_REPORT.md)

### Documentation
- Added [BUG_FIXES_REPORT.md](docs/BUG_FIXES_REPORT.md) - Complete bug fix documentation
- Updated [SESSION_SUMMARY.md](../goca-test-projects/SESSION_SUMMARY.md) - Testing session summary
- Updated [EXTENDED_TESTING_REPORT.md](../goca-test-projects/EXTENDED_TESTING_REPORT.md) - Extended test results

### Impact
- **Zero manual fixes required** - Generated code compiles immediately
- **100% Production Ready** - All critical bugs resolved
- **Improved user experience** - No more compilation errors from generated code

## [1.0.0] - 2025-01-19

### Added
- Initial release of Goca CLI
- `goca init` - Initialize Clean Architecture projects with complete scaffolding
- `goca entity` - Generate domain entities with validation and business rules
- `goca usecase` - Generate use cases with DTOs and interfaces
- `goca handler` - Generate handlers for multiple protocols (HTTP, gRPC, CLI, Worker, SOAP)
- `goca repository` - Generate repositories with database-specific implementations (PostgreSQL, MySQL, MongoDB)
- `goca feature` - Generate complete features with all Clean Architecture layers
- `goca messages` - Generate consistent error messages and responses
- `goca di` - Generate dependency injection containers
- `goca interfaces` - Generate interfaces for TDD development
- `goca version` - Show detailed version information
- Automated release workflow with GitHub Actions
- Cross-platform binaries (Windows, Linux, macOS Intel, macOS Apple Silicon)
- Comprehensive documentation and guides

### Features
- **Complete Clean Architecture code generation** following Uncle Bob's principles
- **Multi-protocol handler support** (HTTP REST, gRPC, CLI commands, Background Workers, SOAP)
- **Database-agnostic repository patterns** with specific implementations
- **Dependency injection setup** with manual and Wire.dev support
- **Validation and business rules generation** at domain level
- **Comprehensive error handling** with structured messaging
- **DTOs and interface generation** for clean layer separation
- **Project scaffolding** with industry best practices
- **TDD support** with interface-first development
- **Automated testing** with GitHub Actions CI/CD

### Technical
- Go 1.21+ support
- Cobra CLI framework
- Clean Architecture enforcement
- Multi-platform builds
- Automated releases

[Unreleased]: https://github.com/sazardev/goca/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/sazardev/goca/releases/tag/v1.0.0
