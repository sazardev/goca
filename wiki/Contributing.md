# Contribution Guide

Thank you for your interest in contributing to Goca! This guide will help you understand how you can participate in the project development.

## 🎯 Ways to Contribute

### 🐛 Report Bugs
- Use the [bug report template](https://github.com/sazardev/goca/issues/new?template=bug_report.md)
- Include version information (`goca version`)
- Provide steps to reproduce the problem
- Include code examples if relevant

### 💡 Suggest Features
- Use the [feature request template](https://github.com/sazardev/goca/issues/new?template=feature_request.md)
- Explain the use case and benefits
- Consider compatibility with Clean Architecture
- Discuss implementation in issues before coding

### 📖 Improve Documentation
- Fix typos
- Add examples and use cases
- Translate documentation
- Improve clarity of explanations

### 🔧 Contribute Code
- Implement new features
- Fix existing bugs
- Optimize performance
- Add tests

## 🚀 Development Environment Setup

### Prerequisites
- **Go 1.21+**
- **Git**
- **Make** (optional)

### Initial Setup
```bash
# 1. Fork the repository on GitHub
# 2. Clone your fork
git clone https://github.com/your-username/goca.git
cd goca

# 3. Add upstream remote
git remote add upstream https://github.com/sazardev/goca.git

# 4. Install dependencies
go mod tidy

# 5. Verify everything works
go build
./goca version
```

### Development Project Structure
```
goca/
├── cmd/                     # CLI commands
│   ├── di.go               # di command
│   ├── entity.go           # entity command
│   ├── feature.go          # feature command
│   ├── handler.go          # handler command
│   ├── init.go             # init command
│   ├── repository.go       # repository command
│   ├── usecase.go          # usecase command
│   ├── version.go          # version command
│   └── utils.go            # Common utilities
├── examples/               # Examples and demos
├── scripts/                # Automation scripts
├── wiki/                   # Wiki documentation
├── .github/workflows/      # CI/CD
├── go.mod
├── main.go
└── README.md
```

## 📝 Development Process

### 1. Create Branch
```bash
# Update main
git checkout main
git pull upstream main

# Create branch for feature/fix
git checkout -b feature/new-functionality
# or
git checkout -b fix/bug-description
```

### 2. Development
```bash
# Make changes
# Run tests
go test ./...

# Verify it compiles
go build

# Test manually
./goca help
```

### 3. Commit Guidelines
We follow [Conventional Commits](https://www.conventionalcommits.org/):

```bash
# Features
git commit -m "feat: add MongoDB support in repositories"

# Bug fixes
git commit -m "fix: correct email validation in entities"

# Documentation
git commit -m "docs: update examples in README"

# Tests
git commit -m "test: add tests for feature command"

# Refactoring
git commit -m "refactor: simplify DTO generation"
```

### 4. Push and Pull Request
```bash
# Push branch
git push origin feature/new-functionality

# Create Pull Request on GitHub
# Use the provided template
# Include detailed description
# Reference related issues
```

## 🧪 Testing

### Run Tests
```bash
# All tests
go test ./...

# Tests with coverage
go test -cover ./...

# Verbose tests
go test -v ./...

# Specific tests
go test ./cmd -run TestEntityGeneration
```

### Write Tests
```go
func TestGenerateEntity(t *testing.T) {
    tests := []struct {
        name     string
        entity   string
        fields   string
        expected string
    }{
        {
            name:     "basic entity",
            entity:   "User",
            fields:   "name:string,email:string",
            expected: "package domain",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := generateEntity(tt.entity, tt.fields, false, false, false, false)
            if !strings.Contains(result, tt.expected) {
                t.Errorf("Expected %s to contain %s", result, tt.expected)
            }
        })
    }
}
```

### Integration Tests
```bash
# Create test project
mkdir test-project
cd test-project

# Test init command
../goca init test --module github.com/test/test

# Verify structure
ls -la test/

# Test feature generation
../goca feature User --fields "name:string,email:string"

# Verify it compiles
cd test && go mod tidy && go build
```

## 📚 Add New Functionality

### 1. New Command
To add a new command (e.g. `goca migrate`):

```go
// cmd/migrate.go
package cmd

import (
    "github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Generate database migrations",
    Long:  `Long description of the command...`,
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(migrateCmd)
    
    // Flags
    migrateCmd.Flags().StringP("database", "d", "postgres", "Database type")
}
```

### 2. New Functionality in Existing Command
To add a flag or modify behavior:

```go
// In the existing command
func init() {
    // New flag
    featureCmd.Flags().BoolP("swagger", "s", false, "Generate Swagger documentation")
}

// In the main function
swagger, _ := cmd.Flags().GetBool("swagger")
if swagger {
    generateSwaggerDocs(featureName)
}
```

### 3. New Templates
To add support for new technologies:

```go
// cmd/repository.go
func generateRedisRepository(dir, entity string) {
    content := `package redis

import (
    "context"
    "github.com/go-redis/redis/v8"
)

type %sRepository struct {
    client *redis.Client
}

func New%sRepository(client *redis.Client) *%sRepository {
    return &%sRepository{
        client: client,
    }
}
`
    content = fmt.Sprintf(content, entity, entity, entity, entity)
    writeFile(filepath.Join(dir, strings.ToLower(entity)+"_repository.go"), content)
}
```

## 🎨 Code Standards

### Formatting
```bash
# Format code
go fmt ./...

# Organize imports
goimports -w .

# Linting
golangci-lint run
```

### Conventions
- **Public functions**: PascalCase with comments
- **Variables**: camelCase descriptive
- **Constants**: UPPER_SNAKE_CASE
- **Files**: snake_case.go
- **Packages**: lowercase, singular

### Comments
```go
// generateEntity creates a new domain entity with the specified fields.
// Parameters:
//   - entityName: entity name (e.g. "User")
//   - fields: comma-separated fields (e.g. "name:string,email:string")
//   - validation: whether to include automatic validations
//   - businessRules: whether to generate business rule methods
func generateEntity(entityName, fields string, validation, businessRules bool) string {
    // Implementation...
}
```

## 🚀 Release Process

### Versioning
We follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Incompatible API changes
- **MINOR**: Backwards compatible new functionality
- **PATCH**: Backwards compatible bug fixes

### Release Process
```bash
# 1. Update version.go
# cmd/version.go
var Version = "1.1.0"

# 2. Update CHANGELOG.md
# Add new section with changes

# 3. Commit and tag
git commit -m "release: v1.1.0"
git tag v1.1.0
git push origin main --tags

# 4. GitHub Actions automatically:
# - Runs tests
# - Compiles binaries
# - Creates GitHub release
# - Publishes to repositories
```

## 📖 Documentation

### Wiki
Documentation is in the `wiki/` directory:

```bash
# Edit documentation
vim wiki/Command-Entity.md

# Verify markdown
markdownlint wiki/*.md

# Preview locally
cd wiki && python -m http.server 8000
```

### README
- Keep examples updated
- Include common use cases
- Verify links work

### Code Comments
- Document public functions
- Explain complex algorithms
- Include usage examples

## 🤝 Community Guidelines

### Communication
- **Be respectful** and constructive
- **Help newcomers** with patience
- **Discuss ideas** before implementing
- **Give useful feedback** in code reviews

### Code Review
- **Review logic** and architecture
- **Verify tests** are included
- **Check documentation** is updated
- **Suggest improvements** constructively

### Issues and Discussions
- **Search for duplicates** before creating
- **Use appropriate templates**
- **Provide complete context**
- **Follow up** on conversations

## 🏆 Recognition

### Contributors
All contributors are recognized in:
- README.md
- Release notes
- Contributors page

### Types of Contribution
- 💻 **Code**: Feature implementation and fixes
- 📖 **Documentation**: Improvements in docs and examples
- 🐛 **Bug Reports**: Issue identification and reporting
- 💡 **Ideas**: Feature suggestions and discussions
- 🎨 **Design**: UX/UI and architecture
- 🔍 **Testing**: Writing and improving tests

## 📞 Contact

### Communication Channels
- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and discussions
- **Email**: sazardev@example.com (main maintainer)

### Expected Response Time
- **Issues**: 24-48 hours
- **Pull Requests**: 2-7 days
- **Discussions**: 1-3 days

## 📋 Checklist for Contributors

### Before Submitting PR
- [ ] Tests pass (`go test ./...`)
- [ ] Code formatted (`go fmt ./...`)
- [ ] Documentation updated
- [ ] CHANGELOG.md updated (for features)
- [ ] Commits follow conventions
- [ ] Branch is updated with main

### For Maintainers
- [ ] Complete code review
- [ ] Integration tests pass
- [ ] Documentation reviewed
- [ ] Breaking changes documented
- [ ] Release notes prepared

---

**Thank you for contributing to Goca! Your participation makes this project better for the entire community. 🙏**

**← [Troubleshooting](Troubleshooting) | [Development](Development) →**
