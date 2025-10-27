# Unit Test Generation Feature for Goca

## Overview

Goca now automatically generates comprehensive unit tests for domain entities when using the `goca entity` command. This feature implements **Issue #4** from the project roadmap.

## Features

### Automatic Test Generation

When generating an entity with validation enabled, Goca creates a complete test suite that includes:

1. **Validation Tests** - Table-driven tests for the `Validate()` method
2. **Initialization Tests** - Verify proper field assignment
3. **Edge Case Tests** - Test boundary conditions for string fields
4. **Type-Specific Tests** - Numeric validation, email format validation

### Test Structure

Generated tests follow Go best practices:
- ✅ Table-driven test patterns
- ✅ Using `testify/assert` for readable assertions
- ✅ Comprehensive coverage of validation logic
- ✅ Clear test naming conventions
- ✅ Proper test organization

## Usage

### Basic Usage

```bash
# Generate entity with tests (enabled by default)
goca entity User --fields "name:string,email:string,age:int" --validation --tests
```

### Disable Tests

```bash
# Generate entity without tests
goca entity User --fields "name:string,email:string" --validation --tests=false
```

## Generated Test File Example

For an entity with fields `name:string,email:string,age:int`:

```go
package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUser_Validate tests the Validate method with various scenarios
func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg   string
	}{
		{
			name: "valid entity",
			user: User{
				Name: "John Doe",
				Email: "test@example.com",
				Age: 25,
			},
			wantErr: false,
		},
		{
			name: "invalid user - empty name",
			user: User{
				Name: "",
				Email: "test@example.com",
				Age: 25,
			},
			wantErr: true,
			errMsg: "name",
		},
		// More test cases...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUser_Initialization tests entity field initialization
func TestUser_Initialization(t *testing.T) {
	user := &User{
		Name: "John Doe",
		Email: "test@example.com",
		Age: 25,
	}

	assert.Equal(t, "John Doe", user.Name, "Name should be set correctly")
	assert.Equal(t, "test@example.com", user.Email, "Email should be set correctly")
	assert.Equal(t, 25, user.Age, "Age should be set correctly")
}

// TestUser_Name_EdgeCases tests edge cases for Name field
func TestUser_Name_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid value", value: "Valid Name", wantErr: false},
		{name: "empty string", value: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Name: tt.value,
				Email: "test@example.com",
				Age: 25,
			}

			err := user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Additional field-specific tests...
```

## Running Tests

```bash
# Navigate to domain directory
cd internal/domain

# Run all tests for User entity
go test -v -run TestUser

# Run specific test
go test -v -run TestUser_Validate

# Run with coverage
go test -cover -run TestUser
```

## Implementation Details

### File Structure

```
cmd/
├── entity.go              # Entity command with --tests flag
├── entity_test_generator.go  # Test generation logic
└── feature.go             # Updated to use new signature
```

### Key Functions

1. **`generateEntityTests()`** - Main test generation function
2. **`generateValidationTests()`** - Creates table-driven validation tests
3. **`generateConstructorTests()`** - Creates initialization tests
4. **`generateFieldTests()`** - Creates field-specific edge case tests

### Dependencies

- `github.com/stretchr/testify/assert` - Assertion library (auto-installed)
- Standard Go testing package

## Test Coverage

### Validation Tests

✅ Valid entity cases
✅ Invalid cases for each field
✅ Empty string validation
✅ Negative number validation
✅ Email format validation (for email fields)

### Edge Case Tests

For **string fields**:
- Empty string
- Valid values
- Email format (if field name contains "email")

For **numeric fields** (int, float64):
- Positive values
- Zero values
- Negative values

### Initialization Tests

✅ All fields properly assigned
✅ Field values match expectations

## Configuration Integration

The `--tests` flag respects Goca's configuration system:

```yaml
# .goca.yaml
generation:
  tests:
    enabled: true  # Default
```

CLI flags override configuration:
```bash
goca entity User --fields "name:string" --tests=false
```

## Benefits

1. **Time Savings** - No need to write boilerplate tests manually
2. **Consistency** - All entities have standardized test structure
3. **Coverage** - Comprehensive validation and edge case testing
4. **Best Practices** - Follows Go testing conventions
5. **Maintainability** - Clean, readable test code

## Roadmap Implementation

This feature completes **Issue #4: Unit Test Generation for Entities** from the v1.14.0 Testing Support milestone.

### Acceptance Criteria Met

✅ Test generation integrated into `goca entity` command
✅ Table-driven test pattern implementation
✅ Validation method testing
✅ Field initialization testing
✅ Edge case coverage
✅ testify/assert integration
✅ Documentation in VitePress
✅ Zero compilation errors
✅ Tests pass successfully

## Future Enhancements

Potential improvements for future versions:

1. **Mock generation** for repository interfaces
2. **Benchmark tests** for performance-critical entities
3. **Property-based testing** integration
4. **Test data builders** for complex entities
5. **Custom test templates** via configuration

## Example Workflow

```bash
# 1. Generate entity with tests
goca entity Product --fields "name:string,price:float64,stock:int" --validation --tests

# Output:
# ✓ Generated: internal/domain/product.go
# ✓ Generated: internal/domain/product_test.go
# ✓ Generated: internal/domain/errors.go
# ✓ Generated: internal/domain/product_seeds.go

# 2. Run tests
cd internal/domain
go test -v -run TestProduct

# 3. Check coverage
go test -cover -run TestProduct

# Output:
# PASS
# coverage: 85.7% of statements
# ok      yourproject/internal/domain    0.013s
```

## Documentation

Full documentation available at:
- **Command Reference**: `docs/commands/entity.md`
- **VitePress Web Docs**: https://goca.dev/commands/entity.html#tests

## Testing

All existing tests pass:
```bash
go test ./cmd/... -v
# PASS
# ok      github.com/sazardev/goca/cmd    0.015s
```

## Contribution

For improvements or bug reports:
1. Open an issue on GitHub
2. Reference Issue #4 for context
3. Follow CONTRIBUTING.md guidelines

---

**Version**: 1.14.0 (Q1 2026)
**Status**: ✅ Implemented and Tested
**Maintainer**: @sazardev
