# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
