---
layout: doc
title: Mastering the Repository Pattern in Clean Architecture
titleTemplate: Articles | Goca Blog
description: A comprehensive guide to the Repository pattern, data access abstraction, and how Goca generates database-agnostic repositories that enforce Clean Architecture boundaries
tags:
  - Repository Pattern
  - Data Access
  - Clean Architecture
  - Infrastructure Layer
  - Go
---

<script setup>
import Badge from '../../.vitepress/theme/components/Badge.vue'
</script>

# Mastering the Repository Pattern in Clean Architecture

<div style="display: flex; gap: 0.5rem; margin-bottom: 1rem;">
<Badge type="info">Infrastructure</Badge>
<Badge type="tip">Data Access</Badge>
</div>

The Repository pattern is a critical abstraction that isolates domain logic from data access concerns. Understanding repositories correctly is essential for building applications that remain testable, maintainable, and independent of specific database technologies. When implemented properly, repositories enable you to change databases, add caching, or switch ORMs without touching business logic.

---

## What is the Repository Pattern?

A repository mediates between the domain layer and data mapping layers, acting as an in-memory collection of domain objects. Repositories provide a clean, domain-centric API for data access while hiding the complexities of database interactions, query construction, and persistence mechanisms.

### Core Responsibilities

**Abstraction**: Repositories abstract the details of data access. The domain layer works with repository interfaces that express intent ("find user by email") rather than implementation ("execute SQL query with WHERE clause").

**Collection Semantics**: Repositories present data as collections. You add entities to repositories, remove them, and query them using domain-meaningful methods. The fact that data persists to a database is an implementation detail.

**Domain-Centric API**: Repository methods use domain language and work with domain entities. A `UserRepository` has methods like `Save(user)`, not `Insert(tableName, columns, values)`.

**Testability**: Because repositories are interfaces, you can replace real implementations with in-memory fakes during testing. This allows unit testing of business logic without database setup.

**Database Independence**: Repositories isolate database-specific code. Changing from PostgreSQL to MongoDB requires changing only repository implementations, not domain or application logic.

## Repository vs DAO vs Data Mapper

Developers often confuse repositories with Data Access Objects (DAOs) or Data Mappers. These patterns serve different purposes and operate at different abstraction levels.

### What a Repository Is NOT

**Not a DAO**: DAOs provide CRUD operations on database tables. Repositories provide domain operations on entity collections. A DAO might have `insertUser(name, email)`. A repository has `Save(user *User)`.

**Not a Data Mapper**: Data Mappers convert between database rows and objects. Repositories use data mappers internally but provide higher-level operations that reflect business operations.

**Not a Generic Interface**: Repositories are not `IRepository<T>` with generic CRUD. Each repository has a domain-specific interface. `UserRepository` has methods that make sense for users; `OrderRepository` has methods that make sense for orders.

**Not Query Builders**: Repositories do not expose SQL or query DSLs. They provide intention-revealing methods. Instead of `repository.Query("SELECT * FROM users WHERE age > ?", 18)`, you have `repository.FindAdults()`.

### The Clear Distinction

```
Domain Layer (Business Logic)
    ↓ Uses interface
Repository Interface (Contract)
    ↓ Defined in domain
Infrastructure Layer (Data Access)
    ↓ Implements interface
Repository Implementation (Database-Specific)
    ↓ Uses ORM/Driver
Database
```

Repositories are defined in the domain layer as interfaces but implemented in the infrastructure layer with database-specific code. This inverts the dependency, making infrastructure depend on the domain rather than the reverse.

## The Infrastructure Layer

Repositories form the infrastructure layer in Clean Architecture. This layer contains all the concrete implementations of persistence, external services, and framework-specific code.

### Infrastructure Layer Characteristics

**Implements Domain Interfaces**: The infrastructure layer provides concrete implementations of repository interfaces defined in the domain. The domain depends on abstractions; infrastructure depends on the domain.

**Database-Specific Code**: Infrastructure contains ORM configurations, SQL queries, connection management, and database-specific optimizations. This code is hidden behind interfaces.

**Framework Dependencies**: Infrastructure can depend on GORM, database drivers, caching libraries, and external SDKs. These dependencies do not leak into the domain.

**Swappable Implementations**: You can have multiple repository implementations for the same interface: PostgreSQL for production, in-memory for testing, MongoDB for a specific feature.

### Why Separate Infrastructure?

The infrastructure layer exists because data access mechanisms change independently of business rules:

**Domain Logic**: "A user must have a unique email" is domain logic. It does not care how you check uniqueness.

**Infrastructure Logic**: "Execute SELECT COUNT(*) FROM users WHERE email = ? to check uniqueness" is infrastructure logic. It is specific to SQL databases.

Separating these concerns allows you to:

- Test domain logic without database setup
- Change databases without changing business rules
- Optimize queries without touching domain code
- Support multiple databases simultaneously

## Repository Interface Design

Well-designed repository interfaces express domain intent clearly while remaining independent of implementation details.

### Interface Location: Domain Layer

Repository interfaces belong in the domain layer, typically in a package like `internal/domain` or alongside entity definitions. This placement is critical:

**Domain Owns the Contract**: The domain defines what operations it needs. Infrastructure adapts to the domain's requirements, not vice versa.

**Dependency Inversion**: By placing interfaces in the domain, you invert the dependency. Infrastructure imports domain types, not the other way around.

**No Infrastructure Leakage**: Domain interfaces use only domain types. They do not reference database connections, ORM types, or SQL constructs.

### Method Design Principles

**Intention-Revealing Names**: Methods should express why you are querying, not how. `FindActiveUsers()` is better than `QueryUsersWithStatus("active")`.

**Domain Types Only**: Parameters and return types are domain entities and value objects, never database-specific types like `sql.Row` or `bson.Document`.

**Error Handling**: Return domain errors, not database errors. Instead of returning `sql.ErrNoRows`, return `ErrUserNotFound`.

**No Leaky Abstractions**: Methods should not expose pagination cursors, query builders, or transaction objects unless these concepts exist in the domain.

### Common Repository Methods

Every repository typically includes these fundamental operations:

```go
type UserRepository interface {
    // Save adds or updates a user
    Save(user *User) error
    
    // FindByID retrieves a user by unique identifier
    FindByID(id uint) (*User, error)
    
    // FindAll retrieves all users (use with caution in production)
    FindAll() ([]*User, error)
    
    // Update modifies an existing user
    Update(user *User) error
    
    // Delete removes a user
    Delete(id uint) error
}
```

Beyond these basics, add domain-specific query methods:

```go
type UserRepository interface {
    // ... basic methods ...
    
    // Domain-specific queries
    FindByEmail(email string) (*User, error)
    FindAdults() ([]*User, error)
    FindByLastLoginAfter(date time.Time) ([]*User, error)
    CountByStatus(status Status) (int, error)
}
```

## Repository Implementation Patterns

Repository implementations encapsulate all database-specific logic, translating domain operations into database operations.

### Basic PostgreSQL Implementation

```go
package repository

import (
    "github.com/yourorg/yourapp/internal/domain"
    "gorm.io/gorm"
)

type postgresUserRepository struct {
    db *gorm.DB
}

// NewPostgresUserRepository creates a PostgreSQL implementation
func NewPostgresUserRepository(db *gorm.DB) domain.UserRepository {
    return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Save(user *domain.User) error {
    return r.db.Create(user).Error
}

func (r *postgresUserRepository) FindByID(id uint) (*domain.User, error) {
    var user domain.User
    err := r.db.First(&user, id).Error
    if err == gorm.ErrRecordNotFound {
        return nil, domain.ErrUserNotFound
    }
    return &user, err
}

func (r *postgresUserRepository) FindByEmail(email string) (*domain.User, error) {
    var user domain.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err == gorm.ErrRecordNotFound {
        return nil, domain.ErrUserNotFound
    }
    return &user, err
}

func (r *postgresUserRepository) Update(user *domain.User) error {
    return r.db.Save(user).Error
}

func (r *postgresUserRepository) Delete(id uint) error {
    return r.db.Delete(&domain.User{}, id).Error
}

func (r *postgresUserRepository) FindAll() ([]*domain.User, error) {
    var users []*domain.User
    err := r.db.Find(&users).Error
    return users, err
}
```

Notice how the implementation:
- Returns domain errors, not GORM errors
- Uses domain types in signatures
- Hides all GORM-specific code
- Implements the domain interface

### MongoDB Implementation

For NoSQL databases, the implementation differs dramatically, but the interface remains the same:

```go
package repository

import (
    "context"
    "github.com/yourorg/yourapp/internal/domain"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
    collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) domain.UserRepository {
    return &mongoUserRepository{
        collection: db.Collection("users"),
    }
}

func (r *mongoUserRepository) Save(user *domain.User) error {
    ctx := context.TODO()
    _, err := r.collection.InsertOne(ctx, user)
    return err
}

func (r *mongoUserRepository) FindByID(id uint) (*domain.User, error) {
    ctx := context.TODO()
    var user domain.User
    
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    if err == mongo.ErrNoDocuments {
        return nil, domain.ErrUserNotFound
    }
    return &user, err
}

func (r *mongoUserRepository) FindByEmail(email string) (*domain.User, error) {
    ctx := context.TODO()
    var user domain.User
    
    err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err == mongo.ErrNoDocuments {
        return nil, domain.ErrUserNotFound
    }
    return &user, err
}
```

The key insight: both implementations satisfy the same `UserRepository` interface. Application and domain code work with either database without modification.

### In-Memory Implementation for Testing

For unit tests, create an in-memory fake that implements the repository interface:

```go
package repository

import (
    "sync"
    "github.com/yourorg/yourapp/internal/domain"
)

type inMemoryUserRepository struct {
    users  map[uint]*domain.User
    nextID uint
    mu     sync.RWMutex
}

func NewInMemoryUserRepository() domain.UserRepository {
    return &inMemoryUserRepository{
        users:  make(map[uint]*domain.User),
        nextID: 1,
    }
}

func (r *inMemoryUserRepository) Save(user *domain.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if user.ID == 0 {
        user.ID = r.nextID
        r.nextID++
    }
    
    r.users[user.ID] = user
    return nil
}

func (r *inMemoryUserRepository) FindByID(id uint) (*domain.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    user, exists := r.users[id]
    if !exists {
        return nil, domain.ErrUserNotFound
    }
    return user, nil
}

func (r *inMemoryUserRepository) FindByEmail(email string) (*domain.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    for _, user := range r.users {
        if user.Email == email {
            return user, nil
        }
    }
    return nil, domain.ErrUserNotFound
}
```

This in-memory implementation enables fast, isolated unit tests without database dependencies.

## How Goca Generates Repositories

Goca's `goca repository` command generates both interfaces and implementations following Clean Architecture principles.

### Basic Repository Generation

```bash
goca repository User --database postgres
```

This generates:

**1. Repository Interface** (`internal/repository/interfaces.go`):
```go
package repository

import "yourapp/internal/domain"

type UserRepository interface {
    Save(user *domain.User) error
    FindByID(id uint) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id uint) error
    FindAll() ([]*domain.User, error)
}
```

**2. PostgreSQL Implementation** (`internal/repository/postgres_user_repository.go`):
```go
package repository

import (
    "yourapp/internal/domain"
    "gorm.io/gorm"
)

type postgresUserRepository struct {
    db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) UserRepository {
    return &postgresUserRepository{db: db}
}

// ... CRUD implementations ...
```

### Database-Specific Implementations

Goca supports multiple databases, generating appropriate implementations for each:

**PostgreSQL with GORM**:
```bash
goca repository User --database postgres
# Generates GORM-based implementation with SQL transactions
```

**MongoDB**:
```bash
goca repository User --database mongodb
# Generates MongoDB driver implementation with BSON
```

**PostgreSQL with JSONB**:
```bash
goca repository User --database postgres-json
# Generates JSONB-specific queries for nested documents
```

**MySQL**:
```bash
goca repository User --database mysql
# Generates MySQL-specific GORM implementation
```

**DynamoDB**:
```bash
goca repository User --database dynamodb
# Generates AWS SDK v2 implementation with attribute mapping
```

### Custom Query Methods

Goca auto-generates query methods based on entity fields:

```bash
goca repository User --fields "name:string,email:string,age:int,status:string"
```

Generates these additional methods:

```go
type UserRepository interface {
    // Basic CRUD
    Save(user *domain.User) error
    FindByID(id uint) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id uint) error
    FindAll() ([]*domain.User, error)
    
    // Field-based queries (auto-generated)
    FindByName(name string) (*domain.User, error)
    FindByEmail(email string) (*domain.User, error)
    FindByAge(age int) (*domain.User, error)
    FindByStatus(status string) (*domain.User, error)
}
```

### Interface-Only Generation

For Test-Driven Development (TDD), generate only the interface first:

```bash
goca repository User --interface-only
```

This creates the contract without implementation, allowing you to:
1. Write use cases against the interface
2. Create mock implementations for tests
3. Implement the real repository later

### Implementation-Only Generation

If you already have the interface but need a new database implementation:

```bash
goca repository User --implementation --database mongodb
```

This generates only the MongoDB implementation without modifying the interface.

## Advanced Repository Patterns

Beyond basic CRUD, repositories support advanced patterns for complex data access scenarios.

### Specification Pattern

Use specifications to encapsulate query criteria:

```go
type UserSpecification interface {
    IsSatisfiedBy(user *domain.User) bool
    ToSQL() (string, []interface{})
}

type UserRepository interface {
    // ... basic methods ...
    FindBySpec(spec UserSpecification) ([]*domain.User, error)
}

// Usage
activeAdults := NewAndSpecification(
    NewAgeGreaterThanSpec(18),
    NewStatusEqualsSpec("active"),
)
users, err := repo.FindBySpec(activeAdults)
```

### Unit of Work Pattern

Coordinate multiple repository operations in a transaction:

```go
type UnitOfWork interface {
    Users() UserRepository
    Orders() OrderRepository
    
    Begin() error
    Commit() error
    Rollback() error
}

// Usage in use case
func (s *orderService) CreateOrder(userID uint, items []Item) error {
    uow := s.unitOfWork
    
    if err := uow.Begin(); err != nil {
        return err
    }
    defer uow.Rollback()
    
    user, err := uow.Users().FindByID(userID)
    if err != nil {
        return err
    }
    
    order := domain.NewOrder(user, items)
    if err := uow.Orders().Save(order); err != nil {
        return err
    }
    
    return uow.Commit()
}
```

### Caching Layer

Add caching transparently using the decorator pattern:

```go
type cachedUserRepository struct {
    repository UserRepository
    cache      Cache
}

func NewCachedUserRepository(repo UserRepository, cache Cache) UserRepository {
    return &cachedUserRepository{
        repository: repo,
        cache:      cache,
    }
}

func (r *cachedUserRepository) FindByID(id uint) (*domain.User, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("user:%d", id)
    if cached, found := r.cache.Get(cacheKey); found {
        return cached.(*domain.User), nil
    }
    
    // Cache miss: query database
    user, err := r.repository.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    r.cache.Set(cacheKey, user, 5*time.Minute)
    return user, nil
}
```

The use case layer does not know caching exists. The cached repository satisfies the same interface as the uncached version.

### Pagination Support

For large datasets, repositories should support pagination:

```go
type Page struct {
    Items      []*User
    TotalItems int
    Page       int
    PageSize   int
    TotalPages int
}

type UserRepository interface {
    // ... basic methods ...
    FindAllPaginated(page, pageSize int) (*Page, error)
    FindByStatusPaginated(status string, page, pageSize int) (*Page, error)
}

// Implementation
func (r *postgresUserRepository) FindAllPaginated(page, pageSize int) (*Page, error) {
    var users []*domain.User
    var total int64
    
    offset := (page - 1) * pageSize
    
    // Count total
    r.db.Model(&domain.User{}).Count(&total)
    
    // Fetch page
    err := r.db.Limit(pageSize).Offset(offset).Find(&users).Error
    
    return &Page{
        Items:      users,
        TotalItems: int(total),
        Page:       page,
        PageSize:   pageSize,
        TotalPages: int(total)/pageSize + 1,
    }, err
}
```

### Soft Delete Support

Implement soft deletes while maintaining a clean repository interface:

```go
func (r *postgresUserRepository) Delete(id uint) error {
    // Soft delete: set deleted_at timestamp
    return r.db.Model(&domain.User{}).
        Where("id = ?", id).
        Update("deleted_at", time.Now()).Error
}

func (r *postgresUserRepository) FindByID(id uint) (*domain.User, error) {
    var user domain.User
    // Automatically exclude soft-deleted records
    err := r.db.Where("deleted_at IS NULL").First(&user, id).Error
    if err == gorm.ErrRecordNotFound {
        return nil, domain.ErrUserNotFound
    }
    return &user, err
}

// Add method for finding deleted records if needed
func (r *postgresUserRepository) FindDeletedByID(id uint) (*domain.User, error) {
    var user domain.User
    err := r.db.Unscoped().First(&user, id).Error
    if err == gorm.ErrRecordNotFound {
        return nil, domain.ErrUserNotFound
    }
    return &user, err
}
```

## Testing Repository Implementations

Repositories require different testing strategies than domain logic because they involve external dependencies.

### Unit Testing with Mocks

For use case testing, use mock repositories:

```go
type MockUserRepository struct {
    SaveFunc    func(user *domain.User) error
    FindByIDFunc func(id uint) (*domain.User, error)
}

func (m *MockUserRepository) Save(user *domain.User) error {
    if m.SaveFunc != nil {
        return m.SaveFunc(user)
    }
    return nil
}

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
    if m.FindByIDFunc != nil {
        return m.FindByIDFunc(id)
    }
    return nil, domain.ErrUserNotFound
}

// Test use case
func TestCreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{
        FindByEmailFunc: func(email string) (*domain.User, error) {
            return nil, domain.ErrUserNotFound // Email available
        },
        SaveFunc: func(user *domain.User) error {
            user.ID = 1
            return nil
        },
    }
    
    service := NewUserService(mockRepo)
    
    user, err := service.CreateUser(CreateUserInput{
        Name:  "John",
        Email: "john@example.com",
    })
    
    assert.NoError(t, err)
    assert.Equal(t, uint(1), user.ID)
}
```

### Integration Testing with Real Database

For repository implementation testing, use a real database:

```go
func TestPostgresUserRepository_Save(t *testing.T) {
    // Setup test database
    db := setupTestDatabase(t)
    defer cleanupTestDatabase(t, db)
    
    repo := NewPostgresUserRepository(db)
    
    // Create test user
    user := &domain.User{
        Name:  "John Doe",
        Email: "john@example.com",
        Age:   30,
    }
    
    // Test Save
    err := repo.Save(user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    
    // Verify in database
    found, err := repo.FindByID(user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Name, found.Name)
    assert.Equal(t, user.Email, found.Email)
}

func setupTestDatabase(t *testing.T) *gorm.DB {
    db, err := gorm.Open(postgres.Open("postgres://test:test@localhost/testdb"), &gorm.Config{})
    require.NoError(t, err)
    
    // Run migrations
    db.AutoMigrate(&domain.User{})
    
    return db
}

func cleanupTestDatabase(t *testing.T, db *gorm.DB) {
    db.Exec("TRUNCATE TABLE users CASCADE")
}
```

### Testing with Docker

For isolated integration tests, use Docker containers:

```go
func TestWithDocker(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Start PostgreSQL container
    ctx := context.Background()
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_PASSWORD": "test",
            "POSTGRES_DB":       "testdb",
        },
        WaitingFor: wait.ForLog("database system is ready"),
    }
    
    postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    require.NoError(t, err)
    defer postgres.Terminate(ctx)
    
    // Get connection string
    host, _ := postgres.Host(ctx)
    port, _ := postgres.MappedPort(ctx, "5432")
    dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb", host, port.Port())
    
    // Run tests
    db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    repo := NewPostgresUserRepository(db)
    
    // ... test repository operations ...
}
```

## Repository Anti-Patterns

Understanding what not to do is as important as knowing best practices.

### Anti-Pattern: Generic Repository

**Problem**:
```go
type Repository[T any] interface {
    Save(entity T) error
    FindByID(id int) (T, error)
    Update(entity T) error
    Delete(id int) error
}

type UserRepository = Repository[User]
type OrderRepository = Repository[Order]
```

**Why It's Bad**: Generic repositories force all entities to have the same operations. `FindByEmail` makes sense for users, not orders. You lose domain-specific expressiveness.

**Solution**: Create specific interfaces for each entity with domain-meaningful methods.

### Anti-Pattern: Leaky Abstraction

**Problem**:
```go
type UserRepository interface {
    Query(sql string, args ...interface{}) ([]*User, error)
    GetDB() *sql.DB
}
```

**Why It's Bad**: This exposes database implementation details. Consumers must write SQL. You cannot swap databases.

**Solution**: Provide intention-revealing methods that hide database details.

### Anti-Pattern: Business Logic in Repository

**Problem**:
```go
func (r *postgresUserRepository) CreateAdminUser(name, email string) (*User, error) {
    user := &User{
        Name:  name,
        Email: email,
        Role:  "admin",
    }
    
    // Business logic: validate admin email domain
    if !strings.HasSuffix(email, "@company.com") {
        return nil, errors.New("admin must use company email")
    }
    
    return user, r.db.Create(user).Error
}
```

**Why It's Bad**: Business rules belong in the domain or use cases, not repositories. Repositories are for data access only.

**Solution**: Move validation to entity or use case. Repository only persists valid entities.

### Anti-Pattern: Repository Returning DTOs

**Problem**:
```go
type UserRepository interface {
    FindByID(id int) (*UserDTO, error)
}
```

**Why It's Bad**: Repositories work with domain entities, not DTOs. DTOs are for external communication, not internal operations.

**Solution**: Return domain entities. Use cases convert entities to DTOs.

## Best Practices Summary

### Do This

✅ **Define interfaces in domain layer**: Keep contracts with domain code  
✅ **Use domain types in signatures**: Parameters and returns are entities  
✅ **Name methods by intent**: `FindActiveUsers()`, not `QueryUsers()`  
✅ **Return domain errors**: `ErrUserNotFound`, not `sql.ErrNoRows`  
✅ **Keep implementations simple**: One responsibility per repository  
✅ **Test with real databases**: Integration tests verify SQL correctness  
✅ **Use in-memory fakes for unit tests**: Fast, isolated tests  

### Avoid This

❌ **Don't expose database types**: No `*sql.DB`, `*gorm.DB` in interfaces  
❌ **Don't add business logic**: Repositories persist, they don't validate  
❌ **Don't use generic interfaces**: Each entity gets specific methods  
❌ **Don't return DTOs**: Repositories work with domain entities  
❌ **Don't couple to ORM**: Use interfaces, not concrete ORM types  

## Conclusion

The Repository pattern is a cornerstone of Clean Architecture, providing a boundary between domain logic and data access concerns. When implemented correctly, repositories enable:

- **Database Independence**: Change databases without changing business logic
- **Testability**: Unit test without database setup using mocks
- **Maintainability**: Data access code isolated in one layer
- **Flexibility**: Support multiple databases simultaneously
- **Clean Domain**: Domain remains pure, focused on business rules

Goca's `goca repository` command generates repositories that follow these principles, creating both interfaces and database-specific implementations that maintain architectural boundaries while providing practical, production-ready data access code.

Master repositories, and you master one of the most critical patterns in Clean Architecture.
