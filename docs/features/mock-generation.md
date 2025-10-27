# Mock Generation Guide

## Overview

Goca automatically generates mock implementations for your interfaces using `testify/mock`. These mocks are essential for unit testing, enabling you to test components in isolation with full control over dependencies.

## What are Mocks?

Mocks are test doubles that:
- **Replace real dependencies** in unit tests
- **Record method calls** for verification
- **Return configured values** for different scenarios
- **Verify behavior** through assertions
- **Enable isolated testing** of single components

## Quick Start

### Generate Mocks

#### For a New Feature
```bash
# Generate feature with mocks
goca feature User --fields "name:string,email:string" --mocks

# This creates:
# - internal/mocks/mock_user_repository.go
# - internal/mocks/mock_user_usecase.go
# - internal/mocks/mock_user_handler.go
# - internal/mocks/examples/user_mock_examples_test.go
```

#### For an Existing Feature
```bash
# Generate all mocks
goca mocks Product

# Generate specific mocks
goca mocks Order --repository --usecase
goca mocks Customer --repository
```

### Use Mocks in Tests

```go
package usecase_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "yourproject/internal/domain"
    "yourproject/internal/mocks"
    "yourproject/internal/usecase"
)

func TestCreateUser_Success(t *testing.T) {
    // Arrange
    mockRepo := mocks.NewMockUserRepository()
    service := usecase.NewUserService(mockRepo)

    input := usecase.CreateUserInput{
        Name:  "John Doe",
        Email: "john@example.com",
        Age:   30,
    }

    // Setup mock expectation
    mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil)

    // Act
    output, err := service.Create(input)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, output)
    mockRepo.AssertExpectations(t)
}
```

## Generated Mock Structure

### Repository Mock
**Location**: `internal/mocks/mock_{entity}_repository.go`

```go
package mocks

import (
    "github.com/stretchr/testify/mock"
    "yourproject/internal/domain"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Save(user *domain.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) FindByID(id int) (*domain.User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}
```

### Use Case Mock
**Location**: `internal/mocks/mock_{entity}_usecase.go`

```go
type MockUserUseCase struct {
    mock.Mock
}

func (m *MockUserUseCase) Create(input usecase.CreateUserInput) (*usecase.CreateUserOutput, error) {
    args := m.Called(input)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*usecase.CreateUserOutput), args.Error(1)
}
```

### Handler Mock
**Location**: `internal/mocks/mock_{entity}_handler.go`

```go
type MockUserHandler struct {
    mock.Mock
}

func (m *MockUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    m.Called(w, r)
}
```

## Mock Testing Patterns

### 1. Basic Setup and Assertions

```go
func TestBasicMockUsage(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository()
    
    // Setup expectation
    mockRepo.On("FindByID", 1).Return(&domain.User{ID: 1, Name: "John"}, nil)
    
    // Execute
    user, err := mockRepo.FindByID(1)
    
    // Verify
    assert.NoError(t, err)
    assert.Equal(t, "John", user.Name)
    mockRepo.AssertExpectations(t)
}
```

### 2. Error Scenarios

```go
func TestErrorHandling(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository()
    service := usecase.NewUserService(mockRepo)
    
    // Setup error expectation
    expectedErr := errors.New("user not found")
    mockRepo.On("FindByID", 999).Return(nil, expectedErr)
    
    // Execute
    user, err := service.GetByID(999)
    
    // Verify error handling
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Equal(t, expectedErr, err)
}
```

### 3. Argument Matchers

```go
func TestArgumentMatchers(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository()
    
    // Match any user with Age > 18
    mockRepo.On("Save", mock.MatchedBy(func(u *domain.User) bool {
        return u.Age > 18
    })).Return(nil)
    
    // This will match
    err := mockRepo.Save(&domain.User{Name: "John", Age: 25})
    assert.NoError(t, err)
    
    // This won't match and will cause test failure
    err = mockRepo.Save(&domain.User{Name: "Jane", Age: 15})
    assert.Error(t, err)
}
```

### 4. Multiple Calls

```go
func TestMultipleCalls(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository()
    
    // Setup for first call
    mockRepo.On("FindByID", 1).Return(&domain.User{ID: 1}, nil).Once()
    
    // Setup for second call with different return
    mockRepo.On("FindByID", 1).Return(&domain.User{ID: 1, Name: "Updated"}, nil).Once()
    
    // First call
    user1, _ := mockRepo.FindByID(1)
    assert.Empty(t, user1.Name)
    
    // Second call
    user2, _ := mockRepo.FindByID(1)
    assert.Equal(t, "Updated", user2.Name)
    
    mockRepo.AssertNumberOfCalls(t, "FindByID", 2)
}
```

### 5. Call Verification

```go
func TestCallVerification(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository()
    service := usecase.NewUserService(mockRepo)
    
    mockRepo.On("Delete", 1).Return(nil)
    
    // Execute
    err := service.Delete(1)
    
    // Verify
    assert.NoError(t, err)
    mockRepo.AssertCalled(t, "Delete", 1)
    mockRepo.AssertNumberOfCalls(t, "Delete", 1)
    mockRepo.AssertNotCalled(t, "Save")
}
```

### 6. Testing with Dependencies

```go
func TestServiceWithMultipleDependencies(t *testing.T) {
    mockUserRepo := mocks.NewMockUserRepository()
    mockEmailService := mocks.NewMockEmailService()
    service := usecase.NewUserService(mockUserRepo, mockEmailService)
    
    input := usecase.CreateUserInput{Name: "John", Email: "john@example.com"}
    
    // Setup expectations
    mockUserRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil)
    mockEmailService.On("SendWelcomeEmail", "john@example.com").Return(nil)
    
    // Execute
    _, err := service.Create(input)
    
    // Verify both dependencies were called
    assert.NoError(t, err)
    mockUserRepo.AssertExpectations(t)
    mockEmailService.AssertExpectations(t)
}
```

## Best Practices

### 1. **One Mock Per Test**
Create fresh mocks for each test to avoid interference:
```go
func TestCase1(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository() // Fresh instance
    // ... test logic
}

func TestCase2(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository() // Fresh instance
    // ... test logic
}
```

### 2. **Clear Setup-Execute-Verify**
Structure tests with clear sections:
```go
func TestExample(t *testing.T) {
    // Setup (Arrange)
    mockRepo := mocks.NewMockUserRepository()
    mockRepo.On("FindByID", 1).Return(&domain.User{}, nil)
    
    // Execute (Act)
    result, err := service.GetByID(1)
    
    // Verify (Assert)
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### 3. **Test Error Paths**
Always test both success and failure scenarios:
```go
func TestCreate_Success(t *testing.T) { /* happy path */ }
func TestCreate_ValidationError(t *testing.T) { /* validation fails */ }
func TestCreate_RepositoryError(t *testing.T) { /* DB error */ }
```

### 4. **Use Meaningful Test Names**
```go
// Good
func TestCreateUser_WhenEmailAlreadyExists_ReturnsError(t *testing.T)

// Bad
func TestCreate(t *testing.T)
```

### 5. **Verify Call Counts**
Be explicit about expected call counts:
```go
mockRepo.On("Save", mock.Anything).Return(nil).Once()
// or
mockRepo.AssertNumberOfCalls(t, "Save", 1)
```

## testify/mock API Reference

### Setting Expectations

```go
// Basic expectation
mock.On("MethodName", arg1, arg2).Return(returnValue1, returnValue2)

// Match any value
mock.On("Save", mock.Anything).Return(nil)

// Match specific type
mock.On("Save", mock.AnythingOfType("*domain.User")).Return(nil)

// Custom matcher
mock.On("Save", mock.MatchedBy(func(u *domain.User) bool {
    return u.Age > 18
})).Return(nil)

// Call count expectations
mock.On("Method").Return(nil).Once()
mock.On("Method").Return(nil).Twice()
mock.On("Method").Return(nil).Times(3)
```

### Assertions

```go
// Verify all expectations met
mock.AssertExpectations(t)

// Verify specific call
mock.AssertCalled(t, "MethodName", arg1, arg2)
mock.AssertNotCalled(t, "MethodName")

// Verify call count
mock.AssertNumberOfCalls(t, "MethodName", 3)
```

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: Unit Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run unit tests with mocks
        run: go test -v ./internal/usecase/...
      
      - name: Check test coverage
        run: go test -cover ./...
```

## Troubleshooting

### Mock Not Called
**Problem**: `AssertExpectations` fails with "EXPECTED CALL"

**Solution**: Verify method is actually called with exact arguments:
```go
// Debug by checking calls
fmt.Println(mockRepo.Calls)
```

### Type Assertion Panic
**Problem**: `panic: interface conversion`

**Solution**: Ensure mock returns correct type:
```go
// Wrong
mock.On("FindByID", 1).Return(&domain.Product{}, nil) // Wrong type

// Correct
mock.On("FindByID", 1).Return(&domain.User{}, nil) // Correct type
```

### Argument Mismatch
**Problem**: Mock not matching arguments

**Solution**: Use `mock.Anything` or `mock.AnythingOfType`:
```go
mockRepo.On("Save", mock.Anything).Return(nil)
```

## Example Test Suite

See `internal/mocks/examples/{entity}_mock_examples_test.go` for complete examples including:
- Basic CRUD operations
- Error scenarios
- Multiple repository calls
- Handler testing with mocked use cases
- Argument matchers
- Call verification

## Commands Reference

```bash
# Generate all mocks for an entity
goca mocks User

# Generate specific mock types
goca mocks Product --repository
goca mocks Order --usecase
goca mocks Customer --handler

# Generate mocks with feature
goca feature User --fields "name:string,email:string" --mocks

# List generated mocks
ls internal/mocks/
```

## Related Documentation

- [Integration Testing Guide](./integration-testing.md)
- [Unit Testing Best Practices](../guide/best-practices.md)
- [testify/mock Documentation](https://pkg.go.dev/github.com/stretchr/testify/mock)
- [Clean Architecture Testing Strategies](../guide/clean-architecture.md#testing)

## Next Steps

- Generate mocks for your features: `goca mocks [Entity]`
- Write unit tests using generated mocks
- Combine with integration tests for comprehensive coverage
- Set up CI/CD pipeline with automated testing
