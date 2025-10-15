# Code Style Guide

This document defines the code style guidelines for the Goca project. Following these guidelines ensures consistency and maintainability across the codebase.

## Table of Contents

- [General Principles](#general-principles)
- [Go Style Guidelines](#go-style-guidelines)
- [File Organization](#file-organization)
- [Naming Conventions](#naming-conventions)
- [Comments and Documentation](#comments-and-documentation)
- [Error Handling](#error-handling)
- [Testing](#testing)
- [Templates](#templates)
- [Command Line Interface](#command-line-interface)

## General Principles

### Readability First

Code is read more often than it is written. Prioritize:
- Clear, descriptive names
- Simple, straightforward logic
- Appropriate comments where needed
- Consistent formatting

### Clean Architecture Compliance

All code must respect Clean Architecture principles:
- Dependencies point inward
- Domain layer remains pure
- Clear separation of concerns
- Interface-based design

### Go Idioms

Follow standard Go idioms and conventions:
- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Adhere to [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## Go Style Guidelines

### Formatting

Use `gofmt` to format all Go code automatically:

```bash
go fmt ./...
```

### Line Length

- Prefer lines under 100 characters
- Break long lines at logical points
- Use continuation indentation for wrapped lines

```go
// Good
func CreateEntity(
    name string,
    fields []Field,
    options EntityOptions,
) error {
    // implementation
}

// Avoid
func CreateEntity(name string, fields []Field, options EntityOptions, config Config, metadata Metadata) error {
    // implementation
}
```

### Imports

Group imports into three sections:
1. Standard library
2. External dependencies
3. Internal packages

```go
import (
    // Standard library
    "context"
    "fmt"
    "time"

    // External dependencies
    "github.com/spf13/cobra"
    "gorm.io/gorm"

    // Internal packages
    "github.com/sazardev/goca/internal/domain"
    "github.com/sazardev/goca/internal/usecase"
)
```

### Variable Declaration

Use short variable declarations when possible:

```go
// Good
name := "example"
count := 10

// Acceptable for zero values
var buffer bytes.Buffer
var wg sync.WaitGroup

// Required for multiple variables
var (
    width  int
    height int
)
```

## File Organization

### Package Structure

```
package-name/
├── package.go       # Main package file with package documentation
├── types.go         # Type definitions
├── interface.go     # Interface definitions
├── implementation.go # Implementation
└── package_test.go  # Tests
```

### File Headers

Include package documentation in the main package file:

```go
// Package entity provides functionality for generating domain entities
// following Clean Architecture principles. It handles entity creation,
// field validation, and code generation.
package entity
```

### Import Organization

Place imports immediately after package declaration:

```go
package entity

import (
    "context"
    "fmt"
)

// Type definitions follow
```

## Naming Conventions

### Packages

- Use lowercase, single-word names
- Avoid underscores or mixed caps
- Choose clear, descriptive names

```go
// Good
package entity
package repository
package usecase

// Avoid
package entityManager
package repo_impl
```

### Types

Use PascalCase for exported types, camelCase for unexported:

```go
// Exported
type EntityGenerator struct {
    name string
}

// Unexported
type fieldValidator struct {
    rules []Rule
}
```

### Functions and Methods

Use camelCase for unexported, PascalCase for exported:

```go
// Exported
func GenerateEntity(name string) error {
    return validateName(name)
}

// Unexported
func validateName(name string) error {
    // implementation
}
```

### Variables

Use descriptive names, avoid single letters except in short scopes:

```go
// Good
entityName := "User"
fieldCount := len(fields)

// Acceptable for short scopes
for i, field := range fields {
    // i is clear in context
}

// Avoid in larger scopes
n := "User"  // What does n represent?
c := 10      // What does c count?
```

### Constants

Use PascalCase for exported, camelCase for unexported:

```go
const (
    // Exported
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3

    // Unexported
    maxFieldLength = 255
    defaultPort    = 8080
)
```

### Interface Names

Use descriptive names, often ending with "er" for single-method interfaces:

```go
// Single-method interfaces
type Generator interface {
    Generate() error
}

type Validator interface {
    Validate() error
}

// Multi-method interfaces
type EntityRepository interface {
    Create(entity Entity) error
    FindByID(id string) (Entity, error)
    Update(entity Entity) error
    Delete(id string) error
}
```

## Comments and Documentation

### Package Documentation

Every package should have documentation:

```go
// Package entity provides functionality for generating domain entities
// in Clean Architecture projects. It supports field validation, custom
// types, and automatic code generation following best practices.
//
// Basic usage:
//
//     gen := entity.NewGenerator("User")
//     gen.AddField("name", "string")
//     err := gen.Generate()
package entity
```

### Function Documentation

Document all exported functions:

```go
// GenerateEntity creates a new domain entity with the specified name and fields.
// It validates the entity structure and generates the corresponding Go file in
// the internal/domain package.
//
// The function returns an error if:
//   - The entity name is invalid
//   - Field definitions contain errors
//   - File generation fails
//
// Example:
//
//     err := GenerateEntity("User", []Field{
//         {Name: "name", Type: "string"},
//         {Name: "email", Type: "string"},
//     })
func GenerateEntity(name string, fields []Field) error {
    // implementation
}
```

### Inline Comments

Use inline comments for complex logic:

```go
// Convert field type to Go type
goType := convertToGoType(field.Type)

// Special handling for time.Time fields
if goType == "time.Time" {
    // Add time import if not already present
    addImport("time")
}
```

### TODO Comments

Format TODO comments consistently:

```go
// TODO(username): Add support for custom validators
// TODO(username): Optimize performance for large entity sets
```

## Error Handling

### Error Creation

Provide context in error messages:

```go
// Good
if err != nil {
    return fmt.Errorf("failed to generate entity %s: %w", name, err)
}

// Avoid
if err != nil {
    return err  // Loses context
}
```

### Error Checking

Always check errors explicitly:

```go
// Good
data, err := readFile(path)
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}

// Never ignore errors
_ = file.Close()  // Avoid

// Better
if err := file.Close(); err != nil {
    log.Printf("warning: failed to close file: %v", err)
}
```

### Custom Error Types

Use custom error types for specific conditions:

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}
```

## Testing

### Test File Naming

Use `_test.go` suffix:

```go
// Implementation
entity.go

// Tests
entity_test.go
```

### Test Function Naming

Use descriptive names starting with `Test`:

```go
func TestGenerateEntity(t *testing.T) {}
func TestGenerateEntity_InvalidName(t *testing.T) {}
func TestGenerateEntity_WithCustomFields(t *testing.T) {}
```

### Table-Driven Tests

Use table-driven tests for multiple scenarios:

```go
func TestValidateName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {
            name:    "valid name",
            input:   "User",
            wantErr: false,
        },
        {
            name:    "empty name",
            input:   "",
            wantErr: true,
        },
        {
            name:    "invalid characters",
            input:   "User-123",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateName() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Test Helpers

Create helper functions for common test setup:

```go
func setupTestEnvironment(t *testing.T) (cleanup func()) {
    // Setup code
    
    return func() {
        // Cleanup code
    }
}

func TestSomething(t *testing.T) {
    cleanup := setupTestEnvironment(t)
    defer cleanup()
    
    // Test code
}
```

## Templates

### Template Organization

Organize templates logically:

```go
const (
    entityTemplate = `package domain

type {{ .Name }} struct {
    {{ range .Fields -}}
    {{ .Name }} {{ .Type }}
    {{ end -}}
}`
)
```

### Template Formatting

Format templates for readability:

```go
const handlerTemplate = `
package handler

import (
    "{{ .Module }}/internal/domain"
    "{{ .Module }}/internal/usecase"
)

type {{ .Name }}Handler struct {
    useCase usecase.{{ .Name }}UseCase
}
`
```

## Command Line Interface

### Command Naming

Use clear, verb-based command names:

```bash
goca init      # Initialize project
goca feature   # Generate feature
goca entity    # Generate entity
```

### Flag Naming

Use descriptive flag names:

```go
cmd.Flags().StringVar(&name, "name", "", "entity name")
cmd.Flags().StringSliceVar(&fields, "fields", nil, "entity fields")
cmd.Flags().BoolVar(&force, "force", false, "force overwrite")
```

### Command Documentation

Provide clear usage information:

```go
var entityCmd = &cobra.Command{
    Use:   "entity [name]",
    Short: "Generate a domain entity",
    Long: `Generate a domain entity following Clean Architecture principles.
    
The entity command creates a new entity in the internal/domain package
with proper structure, validation, and documentation.

Example:
  goca entity User --fields "name:string,email:string"`,
    RunE: runEntityCommand,
}
```

## Code Review Checklist

Before submitting code, verify:

- [ ] Code is formatted with `gofmt`
- [ ] All functions are documented
- [ ] Error handling is appropriate
- [ ] Tests are included and pass
- [ ] No unnecessary dependencies
- [ ] Follows naming conventions
- [ ] Clean Architecture principles respected
- [ ] Code is clear and maintainable

## Tools

Recommended tools for maintaining code quality:

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run tests
go test ./...

# Check test coverage
go test -cover ./...

# Lint (if golangci-lint installed)
golangci-lint run
```

## Resources

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Clean Code Go](https://github.com/Pungyeon/clean-go-article)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

---

Following these guidelines helps maintain a high-quality, consistent codebase that's easy to understand and maintain.
