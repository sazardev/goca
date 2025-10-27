# Integration Testing Guide

## Overview

Goca's integration testing scaffolding automatically generates comprehensive tests that verify interactions between different layers of your Clean Architecture application.

## What are Integration Tests?

Integration tests verify that:
- **Use cases and repositories** work together correctly
- **Handlers and use cases** interact properly
- **Database operations** execute as expected
- **CRUD operations** function end-to-end
- **Transactions** roll back properly on errors

## Quick Start

### Generate Integration Tests

#### For a New Feature
```bash
# Generate feature with integration tests
goca feature User --fields "name:string,email:string" --integration-tests

# This creates:
# - internal/testing/integration/user_integration_test.go
# - internal/testing/integration/fixtures/user_fixtures.go
# - internal/testing/integration/helpers.go (if doesn't exist)
```

#### For an Existing Feature
```bash
# Generate just the integration tests
goca test-integration Product

# With test containers (recommended for CI/CD)
goca test-integration Order --container

# Without fixtures
goca test-integration Customer --fixtures=false
```

### Run Integration Tests

```bash
# Run all integration tests
go test ./internal/testing/integration -v

# Run specific test
go test ./internal/testing/integration -run TestUserIntegration -v

# With coverage
go test ./internal/testing/integration -cover -v
```

## Generated Test Structure

When you generate integration tests for a feature called `User`, Goca creates:

### 1. Main Integration Test File
**Location**: `internal/testing/integration/user_integration_test.go`

```go
package integration

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// TestUserIntegration tests the complete User feature integration
func TestUserIntegration(t *testing.T) {
    // Setup test database
    db := setupTestDatabase(t, "postgres")
    defer cleanupTestDatabase(t, db)

    // Initialize dependencies
    repo := repository.NewPostgresUserRepository(db)
    service := usecase.NewUserService(repo)

    ctx := context.Background()

    t.Run("CreateAndRetrieveUser", func(t *testing.T) {
        // Test CRUD operations
    })

    t.Run("UpdateUser", func(t *testing.T) {
        // Test update logic
    })

    t.Run("DeleteUser", func(t *testing.T) {
        // Test deletion
    })

    t.Run("ListUsers", func(t *testing.T) {
        // Test listing
    })

    t.Run("TransactionUser", func(t *testing.T) {
        // Test transaction rollback
    })
}
```

### 2. Test Fixtures
**Location**: `internal/testing/integration/fixtures/user_fixtures.go`

```go
package fixtures

import "github.com/sazardev/goca/internal/domain"

// NewUserFixture creates a User with default test values
func NewUserFixture() *domain.User {
    return &domain.User{
        Name:  "Test User",
        Email: "test@example.com",
    }
}

// NewUserFixtureWithCustomFields creates customized User
func NewUserFixtureWithCustomFields(fields map[string]interface{}) *domain.User {
    user := NewUserFixture()
    // Override fields as needed
    return user
}

// NewUserFixtureList creates multiple User fixtures
func NewUserFixtureList(count int) []*domain.User {
    fixtures := make([]*domain.User, count)
    for i := 0; i < count; i++ {
        fixtures[i] = NewUserFixture()
        // Vary data for each fixture
    }
    return fixtures
}
```

### 3. Database Helpers
**Location**: `internal/testing/integration/helpers.go`

```go
package integration

import (
    "database/sql"
    "testing"
    "gorm.io/gorm"
)

// setupTestDatabase initializes a test database connection
func setupTestDatabase(t *testing.T, dbType string) *gorm.DB {
    // Setup test database (in-memory or test container)
}

// cleanupTestDatabase cleans up after tests
func cleanupTestDatabase(t *testing.T, db *gorm.DB) {
    // Drop tables, close connection
}

// truncateTables truncates specified tables
func truncateTables(t *testing.T, db *gorm.DB, tables ...string) {
    // Clean data between tests
}

// beginTransaction starts a test transaction
func beginTransaction(t *testing.T, db *gorm.DB) *gorm.DB {
    // For testing rollback behavior
}

// seedTestData seeds database with test data
func seedTestData(t *testing.T, db *gorm.DB) {
    // Populate test database
}
```

## Test Database Setup

### Option 1: In-Memory Database (SQLite)
**Fastest for development, no external dependencies**

```go
// In helpers.go
func setupTestDatabase(t *testing.T, dbType string) *gorm.DB {
    switch dbType {
    case "sqlite":
        return setupSQLiteMemory(t)
    // ...
    }
}

func setupSQLiteMemory(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to open SQLite: %v", err)
    }
    
    // Run migrations
    db.AutoMigrate(&domain.User{}, &domain.Product{})
    
    return db
}
```

### Option 2: Test Database Server
**Most realistic, requires running database**

```go
func setupTestDatabase(t *testing.T, dbType string) *gorm.DB {
    dsn := "host=localhost user=test password=test dbname=goca_test port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to connect: %v", err)
    }
    
    // Run migrations
    db.AutoMigrate(&domain.User{})
    
    return db
}
```

### Option 3: Test Containers (Recommended for CI/CD)
**Isolated, automatic cleanup, CI-friendly**

```go
import (
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDatabase(t *testing.T, dbType string) *gorm.DB {
    ctx := context.Background()
    
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "test",
            "POSTGRES_PASSWORD": "test",
            "POSTGRES_DB":       "testdb",
        },
        WaitingFor: wait.ForLog("database system is ready"),
    }
    
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatalf("failed to start container: %v", err)
    }
    
    // Get connection details and connect
    // ...
    
    return db
}
```

## Best Practices

### 1. Test Isolation
Each test should be independent:

```go
t.Run("CreateUser", func(t *testing.T) {
    // Start fresh transaction or truncate tables
    tx := beginTransaction(t, db)
    defer rollbackTransaction(t, tx)
    
    // Run test with tx
    repo := repository.NewPostgresUserRepository(tx)
    // ...
})
```

### 2. Use Fixtures
Create reusable test data:

```go
func TestUserService(t *testing.T) {
    user := fixtures.NewUserFixture()
    // Customize for this test
    user.Email = "specific@test.com"
    
    err := repo.Save(ctx, user)
    require.NoError(t, err)
}
```

### 3. Test Edge Cases
Test both happy path and error scenarios:

```go
t.Run("CreateDuplicateUser", func(t *testing.T) {
    user := fixtures.NewUserFixture()
    
    // First creation should succeed
    _, err := service.CreateUser(ctx, input)
    require.NoError(t, err)
    
    // Second with same email should fail
    _, err = service.CreateUser(ctx, input)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "already exists")
})
```

### 4. Verify Database State
Check that database reflects expected state:

```go
t.Run("DeleteUser", func(t *testing.T) {
    user := createTestUser(t, service)
    
    err := service.DeleteUser(ctx, user.ID)
    require.NoError(t, err)
    
    // Verify user is deleted
    _, err = repo.FindByID(ctx, user.ID)
    assert.Error(t, err)
    assert.Equal(t, domain.ErrUserNotFound, err)
})
```

### 5. Test Transactions
Verify proper transaction handling:

```go
t.Run("TransactionRollback", func(t *testing.T) {
    initialCount := countUsers(t, repo)
    
    // Create invalid data that should rollback
    input := usecase.CreateUserInput{
        Email: "", // Invalid - should fail validation
    }
    
    _, err := service.CreateUser(ctx, input)
    assert.Error(t, err)
    
    // Verify no user was created
    finalCount := countUsers(t, repo)
    assert.Equal(t, initialCount, finalCount)
})
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: goca_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run integration tests
        run: go test ./internal/testing/integration -v
        env:
          DATABASE_URL: postgres://test:test@localhost:5432/goca_test?sslmode=disable
```

## Troubleshooting

### Tests Hang or Timeout
- Check database connection settings
- Ensure database is running and accessible
- Verify no resource leaks (unclosed connections)

### Flaky Tests
- Ensure proper test isolation
- Use transactions for test data
- Avoid hardcoded timing assumptions

### Slow Tests
- Use in-memory database for development
- Consider parallel test execution
- Optimize database queries in tests

## Advanced Topics

### Parallel Test Execution
```go
func TestUserIntegration(t *testing.T) {
    t.Parallel() // Run in parallel with other tests
    
    db := setupTestDatabase(t, "postgres")
    defer cleanupTestDatabase(t, db)
    
    // Each parallel test needs isolated database
}
```

### Custom Assertions
```go
func assertUserEquals(t *testing.T, expected, actual *domain.User) {
    t.Helper()
    assert.Equal(t, expected.Name, actual.Name)
    assert.Equal(t, expected.Email, actual.Email)
    // ... more assertions
}
```

### Test Data Builders
```go
type UserBuilder struct {
    user *domain.User
}

func NewUserBuilder() *UserBuilder {
    return &UserBuilder{
        user: &domain.User{
            Name:  "Default User",
            Email: "default@test.com",
        },
    }
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
    b.user.Name = name
    return b
}

func (b *UserBuilder) Build() *domain.User {
    return b.user
}

// Usage in tests
user := NewUserBuilder().
    WithName("Custom Name").
    Build()
```

## Related Commands

- `goca feature` - Generate complete feature with `--integration-tests` flag
- `goca test-integration` - Generate only integration tests for existing feature
- `goca init` - Initialize new project with testing support

## Next Steps

- See [Complete Tutorial](../tutorials/complete-tutorial.md) for hands-on examples
- Check [Database Support](./database-support.md) for database-specific configuration
- Review [Best Practices](../guide/best-practices.md) for testing strategies
