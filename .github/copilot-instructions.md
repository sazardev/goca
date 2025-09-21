# Goca CLI - AI Coding Agent Instructions

## Project Overview
Goca is a Go CLI tool that generates Clean Architecture projects and code. It's primarily written in Spanish for documentation but Go code uses English. The tool generates complete features following strict Clean Architecture layering.

## Architecture & Patterns

### Command Structure
- All CLI commands live in `cmd/` with one file per command (e.g., `feature.go`, `entity.go`)
- Each command uses Cobra framework and follows pattern: `<name>Cmd` variable + `init()` function
- Commands use centralized validation through `CommandValidator` in `cmd/command_validator.go`
- Error handling follows consistent pattern through `ErrorHandler` in `cmd/errors.go`

### Code Generation
- Template-based generation using Go text templates in `cmd/templates.go`
- All generated code follows Clean Architecture layers: `domain/`, `usecase/`, `repository/`, `handler/`
- Auto-integration with dependency injection happens in `generateCompleteFeature()` functions
- Field validation uses `FieldValidator` with specific rules for entity names, field types, etc.

### Key Components
- **Feature Generation**: Main workflow that creates complete CRUD with all layers
- **Template System**: Dynamic code generation with validation, timestamps, soft delete features
- **Validation Layer**: Multi-level validation (command params, field types, business rules)
- **Auto-Integration**: Automatic DI container updates and main.go modifications

## Development Workflows

### Building & Testing
```bash
make build          # Build CLI binary
make test-cli       # Run comprehensive CLI tests (recommended for changes)
make test-coverage  # Generate coverage reports
make lint           # Run linting (part of dev workflow)
```

### Testing Philosophy
- Tests live in `internal/testing/` with comprehensive validation framework
- **Critical**: All generated code must compile with zero errors/warnings
- Tests validate file structure, code quality, and architectural compliance
- Use `TestSuite` for integration testing that runs actual CLI commands

### Code Quality Standards
- Zero compilation errors policy - tests enforce this
- Spanish documentation/comments for user-facing content
- English for code identifiers and internal comments
- Consistent error messaging through centralized handlers

## Project-Specific Conventions

### Field Definition Pattern
Fields use format: `name:type,email:string,age:int` - this is parsed by `parseFields()` in multiple commands

### Flag Naming
- `--fields` for entity field definitions
- `--database` for database type (postgres, mysql, sqlite)
- `--handlers` for interface types (http, grpc, cli)
- `--validation` and `--business-rules` for feature flags

### File Generation Rules
- Always use `ensureDirectoryExists()` before file creation
- Template data structures in `generateCompleteFeature()` functions
- Auto-generated files include imports, validation, and proper Go formatting

### Integration Points
- DI container auto-updates happen in `autoIntegrateFeature()` functions
- Main.go modifications are automatic for new features
- Repository pattern implementation is consistent across all database types

## Anti-Patterns to Avoid
- Don't bypass the CommandValidator - all input validation goes through it
- Don't hardcode file paths - use filepath.Join() and constants
- Don't skip the comprehensive test suite when adding new commands
- Don't generate code without proper Clean Architecture layer separation

## Key Files for Understanding
- `cmd/feature.go` - Core feature generation logic and workflow
- `cmd/templates.go` - All code generation templates
- `internal/testing/suite.go` - Testing framework and validation rules
- `examples/full-flow.md` - Complete usage examples and generated structure

## When Adding New Commands
1. Create new file in `cmd/` following existing pattern
2. Add validation rules to `CommandValidator`
3. Create templates in `templates.go` if generating code
4. Add comprehensive tests in `internal/testing/`
5. Update auto-integration logic if needed
6. Test with `make test-cli` to ensure zero compilation errors