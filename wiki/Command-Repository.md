# goca repository Command

The `goca repository` command creates repositories that implement the Repository pattern with well-defined interfaces and database-specific implementations following Clean Architecture.

## 📋 Syntax

```bash
goca repository <entity> [flags]
```

## 🎯 Purpose

Creates repositories to handle entity persistence:

- 🔵 **Persistence abstraction** without coupling the domain
- 📊 **Database-specific implementations** for each technology
- 🔗 **Clear interfaces** for use cases
- 💾 **Transactions** and error handling
- ⚡ **Optional cache** for optimization
- 🔍 **Optimized queries** per technology

## 🚩 Available Flags

| Flag               | Type     | Required | Default Value | Description                                    |
| ------------------ | -------- | -------- | ------------- | ---------------------------------------------- |
| `--database`       | `string` | ❌ No     | -             | Database type (`postgres`, `mysql`, `mongodb`) |
| `--interface-only` | `bool`   | ❌ No     | `false`       | Generate interfaces only                       |
| `--implementation` | `bool`   | ❌ No     | `false`       | Generate implementation only                   |
| `--transactions`   | `bool`   | ❌ No     | `false`       | Include transaction support                    |
| `--cache`          | `bool`   | ❌ No     | `false`       | Include cache layer                            |

## 📖 Usage Examples

### Basic Repository with PostgreSQL
```bash
goca repository User --database postgres
```

### Generate Interfaces Only
```bash
goca repository Product --interface-only
```

### With Transactions and Cache
```bash
goca repository Order --database postgres --transactions --cache
```

### Different Databases
```bash
# PostgreSQL
goca repository User --database postgres --transactions

# MySQL
goca repository Product --database mysql --cache

# MongoDB
goca repository Order --database mongodb
```

## 📂 Generated Files

### File Structure
```
internal/repository/
├── interfaces/
│   └── user_repository.go      # Repository interface
├── postgres/
│   └── user_repository.go      # PostgreSQL implementation
├── mysql/
│   └── user_repository.go      # MySQL implementation (if specified)
└── mongodb/
    └── user_repository.go      # MongoDB implementation (if specified)
```

## 🔍 Generated Code in Detail

### Interface: `internal/repository/interfaces/user_repository.go`

```go
package interfaces

import (
    "context"
    
    "github.com/usuario/proyecto/internal/domain"
)

//go:generate mockgen -source=user_repository.go -destination=mocks/user_repository_mock.go

// UserRepository defines contracts for user persistence
type UserRepository interface {
    // Basic CRUD operations
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id uint) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uint) error
    
    // Query operations
    List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
    Search(ctx context.Context, query string, limit, offset int) ([]*domain.User, int64, error)
    Exists(ctx context.Context, id uint) (bool, error)
    
    // Transaction operations (if --transactions)
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
    
    // Cache operations (if --cache)
    ClearCache(ctx context.Context, id uint) error
}
```

### PostgreSQL Implementation: `internal/repository/postgres/user_repository.go`

```go
package postgres

import (
    "context"
    "database/sql"
    "fmt"
    
    "github.com/lib/pq"
    "github.com/usuario/proyecto/internal/domain"
    "github.com/usuario/proyecto/internal/repository/interfaces"
)

// UserRepository implements interfaces.UserRepository for PostgreSQL
type UserRepository struct {
    db *sql.DB
}

// NewUserRepository creates a new repository instance
func NewUserRepository(db *sql.DB) interfaces.UserRepository {
    return &UserRepository{
        db: db,
    }
}

// Save stores a user in the database
func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (name, email, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
    
    err := r.db.QueryRowContext(
        ctx,
        query,
        user.Name,
        user.Email,
        user.CreatedAt,
        user.UpdatedAt,
    ).Scan(&user.ID)
    
    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok {
            switch pqErr.Code {
            case "23505": // unique_violation
                return domain.ErrUserEmailAlreadyExists
            }
        }
        return fmt.Errorf("failed to save user: %w", err)
    }
    
    return nil
}

// FindByID searches for a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `
    
    user := &domain.User{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, domain.ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to find user by ID: %w", err)
    }
    
    return user, nil
}

// FindByEmail searches for a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE email = $1 AND deleted_at IS NULL
    `
    
    user := &domain.User{}
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // No error if it doesn't exist
        }
        return nil, fmt.Errorf("failed to find user by email: %w", err)
    }
    
    return user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
    query := `
        UPDATE users
        SET name = $2, email = $3, updated_at = $4
        WHERE id = $1 AND deleted_at IS NULL
    `
    
    result, err := r.db.ExecContext(
        ctx,
        query,
        user.ID,
        user.Name,
        user.Email,
        user.UpdatedAt,
    )
    
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return domain.ErrUserNotFound
    }
    
    return nil
}

// Delete removes a user (soft delete)
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
    query := `
        UPDATE users
        SET deleted_at = NOW()
        WHERE id = $1 AND deleted_at IS NULL
    `
    
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return domain.ErrUserNotFound
    }
    
    return nil
}

// List gets a paginated list of users
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
    // Count total
    countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
    var total int64
    err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count users: %w", err)
    }
    
    // Get paginated users
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
    
    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to list users: %w", err)
    }
    defer rows.Close()
    
    var users []*domain.User
    for rows.Next() {
        user := &domain.User{}
        err := rows.Scan(
            &user.ID,
            &user.Name,
            &user.Email,
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan user: %w", err)
        }
        users = append(users, user)
    }
    
    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("failed to iterate users: %w", err)
    }
    
    return users, total, nil
}
```

## 💾 Database-Specific Implementations

### PostgreSQL Features
- **JSONB** for complex fields
- **Array types** for lists
- **UUID** as primary keys
- **Partial indexes** for soft delete
- **RETURNING** clause to get IDs

### MySQL Features
- **JSON** for complex fields
- **Generated columns** for calculated fields
- **Multi-value indexes** for searches
- **Foreign key constraints** with CASCADE

### MongoDB Features
- **Aggregation pipeline** for complex queries
- **Compound indexes** for optimization
- **GridFS** for large files
- **Transactions** in replica sets

## 🔄 Transactions (--transactions)

With `--transactions`, methods for transaction management are added:

```go
// WithTransaction executes operations within a transaction
func (r *UserRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    // Create context with transaction
    txCtx := context.WithValue(ctx, "tx", tx)
    
    if err := fn(txCtx); err != nil {
        if rollbackErr := tx.Rollback(); rollbackErr != nil {
            return fmt.Errorf("failed to rollback transaction: %v, original error: %w", rollbackErr, err)
        }
        return err
    }
    
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}

// SaveWithTx saves using existing transaction
func (r *UserRepository) SaveWithTx(ctx context.Context, user *domain.User) error {
    var executor interface {
        QueryRowContext(context.Context, string, ...interface{}) *sql.Row
    } = r.db
    
    if tx, ok := ctx.Value("tx").(*sql.Tx); ok {
        executor = tx
    }
    
    query := `
        INSERT INTO users (name, email, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
    
    return executor.QueryRowContext(
        ctx, query,
        user.Name, user.Email, user.CreatedAt, user.UpdatedAt,
    ).Scan(&user.ID)
}
```

## ⚡ Cache (--cache)

With `--cache`, a cache layer is integrated:

```go
import (
    "encoding/json"
    "time"
    
    "github.com/go-redis/redis/v8"
)

type CachedUserRepository struct {
    repo  interfaces.UserRepository
    cache *redis.Client
    ttl   time.Duration
}

func NewCachedUserRepository(repo interfaces.UserRepository, cache *redis.Client) interfaces.UserRepository {
    return &CachedUserRepository{
        repo:  repo,
        cache: cache,
        ttl:   time.Hour,
    }
}

// FindByID with cache
func (r *CachedUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
    cacheKey := fmt.Sprintf("user:id:%d", id)
    
    // Try to get from cache
    cached, err := r.cache.Get(ctx, cacheKey).Result()
    if err == nil {
        var user domain.User
        if err := json.Unmarshal([]byte(cached), &user); err == nil {
            return &user, nil
        }
    }
    
    // If not in cache, get from DB
    user, err := r.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Save to cache
    if user != nil {
        if userData, err := json.Marshal(user); err == nil {
            r.cache.Set(ctx, cacheKey, userData, r.ttl)
        }
    }
    
    return user, nil
}

// ClearCache clears user cache
func (r *CachedUserRepository) ClearCache(ctx context.Context, id uint) error {
    patterns := []string{
        fmt.Sprintf("user:id:%d", id),
        "user:list:*",
        "user:search:*",
    }
    
    for _, pattern := range patterns {
        keys, err := r.cache.Keys(ctx, pattern).Result()
        if err != nil {
            continue
        }
        
        if len(keys) > 0 {
            r.cache.Del(ctx, keys...)
        }
    }
    
    return nil
}
```

## 🔍 Advanced Queries

### Full-Text Search
```go
// Search searches users by text
func (r *UserRepository) Search(ctx context.Context, query string, limit, offset int) ([]*domain.User, int64, error) {
    searchQuery := `
        SELECT id, name, email, created_at, updated_at,
               ts_rank_cd(search_vector, plainto_tsquery($1)) as rank
        FROM users
        WHERE search_vector @@ plainto_tsquery($1)
          AND deleted_at IS NULL
        ORDER BY rank DESC, created_at DESC
        LIMIT $2 OFFSET $3
    `
    
    rows, err := r.db.QueryContext(ctx, searchQuery, query, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to search users: %w", err)
    }
    defer rows.Close()
    
    var users []*domain.User
    for rows.Next() {
        user := &domain.User{}
        var rank float64
        
        err := rows.Scan(
            &user.ID, &user.Name, &user.Email,
            &user.CreatedAt, &user.UpdatedAt, &rank,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan search result: %w", err)
        }
        
        users = append(users, user)
    }
    
    // Count results
    countQuery := `
        SELECT COUNT(*)
        FROM users
        WHERE search_vector @@ plainto_tsquery($1) AND deleted_at IS NULL
    `
    
    var total int64
    err = r.db.QueryRowContext(ctx, countQuery, query).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count search results: %w", err)
    }
    
    return users, total, nil
}
```

## 🧪 Testing

Generated repositories include interfaces for easy testing:

```go
func TestUserRepository_Save(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := postgres.NewUserRepository(db)
    
    user := &domain.User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    err := repo.Save(context.Background(), user)
    
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}

func TestUserRepository_WithTransaction(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := postgres.NewUserRepository(db)
    
    err := repo.WithTransaction(context.Background(), func(ctx context.Context) error {
        user1 := &domain.User{Name: "User 1", Email: "user1@test.com"}
        user2 := &domain.User{Name: "User 2", Email: "user2@test.com"}
        
        if err := repo.SaveWithTx(ctx, user1); err != nil {
            return err
        }
        
        return repo.SaveWithTx(ctx, user2)
    })
    
    assert.NoError(t, err)
}
```

## ⚠️ Important Considerations

### ✅ Best Practices
- **Context propagation**: Use context.Context in all methods
- **Error wrapping**: Wrap errors with contextual information
- **Prepared statements**: Prevent SQL injection
- **Connection pooling**: Configure pools appropriately

### ❌ Common Errors
- **Not using transactions**: For operations requiring consistency
- **Ignoring errors**: Always handle errors appropriately
- **N+1 queries**: Optimize with JOINs or eager loading
- **Not cleaning resources**: Close rows, statements, etc.

---

**← [goca usecase Command](Command-UseCase) | [goca handler Command](Command-Handler) →**
