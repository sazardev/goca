# Contributing to Goca

Thank you for your interest in contributing to Goca. This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Community](#community)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to sazardev@gmail.com.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a new branch for your contribution
4. Make your changes
5. Push to your fork and submit a pull request

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include as many details as possible:

- Use a clear and descriptive title
- Describe the exact steps to reproduce the problem
- Provide specific examples to demonstrate the steps
- Describe the behavior you observed and what you expected to see
- Include screenshots if applicable
- Specify your environment (OS, Go version, Goca version)

Use our [bug report template](.github/ISSUE_TEMPLATE/bug_report.md) when creating issues.

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- Use a clear and descriptive title
- Provide a detailed description of the proposed feature
- Explain why this enhancement would be useful
- List examples of how the feature would be used
- Consider if it aligns with the project's Clean Architecture philosophy

Use our [feature request template](.github/ISSUE_TEMPLATE/feature_request.md) when creating suggestions.

### Pull Requests

- Fill in the required template
- Follow the coding standards outlined below
- Include appropriate test coverage
- Update documentation as needed
- End all files with a newline
- Avoid platform-dependent code

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Setup Steps

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/goca.git
cd goca

# Add upstream remote
git remote add upstream https://github.com/sazardev/goca.git

# Install dependencies
go mod download

# Build the project
make build

# Run tests
make test
```

### Project Structure

```
goca/
├── cmd/                  # CLI commands implementation
├── internal/            # Internal packages
│   ├── constants/       # Constants and configuration
│   ├── domain/          # Domain entities and business logic
│   ├── handler/         # Handler implementations
│   ├── interfaces/      # Interface definitions
│   ├── messages/        # Message templates
│   ├── repository/      # Repository implementations
│   ├── testing/         # Testing utilities
│   └── usecase/         # Use case implementations
├── docs/                # Documentation
└── wiki/                # Wiki content
```

## Pull Request Process

1. **Branch Naming**: Use descriptive branch names
   - `feature/add-new-command`
   - `fix/template-generation-bug`
   - `docs/update-contributing-guide`

2. **Commit Messages**: Follow conventional commit format
   ```
   type(scope): subject
   
   body
   
   footer
   ```
   
   Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
   
   Example:
   ```
   feat(entity): add support for custom validation tags
   
   - Add ValidationTag field to EntityField struct
   - Update template to include custom validations
   - Add tests for validation tag generation
   
   Closes #123
   ```

3. **Code Review**: 
   - Address review comments promptly
   - Be open to feedback and suggestions
   - Keep discussions professional and constructive

4. **Update Documentation**:
   - Update README.md if adding new features
   - Add or update relevant documentation in `/docs`
   - Update command help text if modifying commands

5. **Testing**:
   - Add tests for new functionality
   - Ensure all tests pass before submitting
   - Maintain or improve code coverage

6. **Merge Requirements**:
   - All tests must pass
   - Code review approval from maintainers
   - No merge conflicts with base branch
   - Documentation is updated

## Coding Standards

### Go Style Guide

Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

### Formatting

- Use `gofmt` to format all Go code
- Run `go vet` to catch common mistakes
- Use `golangci-lint` for comprehensive linting

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Lint code (if golangci-lint is installed)
golangci-lint run
```

### Naming Conventions

- Use camelCase for variable and function names
- Use PascalCase for exported functions and types
- Use descriptive names that clearly indicate purpose
- Avoid abbreviations unless widely understood

### Error Handling

- Always handle errors explicitly
- Provide context in error messages
- Use custom error types for specific error conditions
- Return errors rather than panicking in library code

```go
// Good
if err != nil {
    return fmt.Errorf("failed to generate entity: %w", err)
}

// Avoid
if err != nil {
    panic(err)
}
```

### Comments

- Write clear, concise comments for exported functions and types
- Use complete sentences in comments
- Begin comments with the name of the element being described
- Document complex logic or non-obvious decisions

```go
// GenerateEntity creates a new entity file based on the provided configuration.
// It validates the entity name and fields before generating the file.
func GenerateEntity(config EntityConfig) error {
    // Implementation
}
```

## Testing Guidelines

### Test Coverage

- Aim for at least 70% code coverage
- Focus on testing critical paths and edge cases
- Write table-driven tests for functions with multiple scenarios

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
        {
            name:    "invalid input",
            input:   invalidInput,
            want:    OutputType{},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("FunctionName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Documentation

### Code Documentation

- Document all exported types, functions, and methods
- Use godoc-compatible comments
- Include examples for complex functionality

### User Documentation

- Update user-facing documentation in `/docs`
- Update command help text for CLI changes
- Include examples and use cases
- Keep documentation in sync with code changes

### Documentation Standards

- Use clear, concise language
- Include code examples where appropriate
- Use proper Markdown formatting
- Verify all links work correctly

## Community

### Communication Channels

- **GitHub Issues**: Bug reports and feature requests
- **Pull Requests**: Code contributions and discussions
- **Email**: sazardev@gmail.com for security issues or private matters

### Getting Help

- Check existing documentation and issues first
- Provide clear and detailed information when asking for help
- Be patient and respectful with community members

### Recognition

Contributors will be recognized in the project's release notes and documentation. Significant contributions may result in being added to the maintainers team.

## License

By contributing to Goca, you agree that your contributions will be licensed under the MIT License.

## Questions?

If you have questions about contributing, please open an issue with the question label or contact the maintainers directly.

Thank you for contributing to Goca!
