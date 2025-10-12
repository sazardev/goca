# goca repository

Generate repository interfaces and database implementations.

## Syntax

```bash
goca repository <EntityName> [flags]
```

## Description

Creates repository pattern implementations for data persistence, abstracting database operations behind clean interfaces.

## Flags

### `--database`

Database type. Default: `postgres`

**Options:** `postgres` | `mysql` | `mongodb` | `sqlite`

```bash
goca repository Product --database postgres
```

### `--interface-only`

Generate only the interface.

```bash
goca repository User --interface-only
```

### `--implementation`

Generate only the implementation.

```bash
goca repository User --implementation --database mysql
```

## Examples

### PostgreSQL Repository

```bash
goca repository User --database postgres
```

### MongoDB Repository

```bash
goca repository Product --database mongodb
```

### Interface Only

```bash
goca repository Order --interface-only
```

## Generated Files

```
internal/repository/
├── interfaces.go               # Repository interfaces
└── postgres_user_repository.go # Implementation
```

## Generated Code Example

```go
// interfaces.go
package repository

import (
    "context"
    "myproject/internal/domain"
)

type UserRepository interface {
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id uint) (*domain.User, error)
    FindAll(ctx context.Context) ([]*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uint) error
}

// postgres_user_repository.go
package repository

import (
    "context"
    "database/sql"
    "myproject/internal/domain"
)

type postgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
    return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Save(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (name, email, age, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
    
    return r.db.QueryRowContext(
        ctx, query,
        user.Name, user.Email, user.Age,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
    query := `
        SELECT id, name, email, age, created_at, updated_at
        FROM users
        WHERE id = $1
    `
    
    user := &domain.User{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID, &user.Name, &user.Email, &user.Age,
        &user.CreatedAt, &user.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, domain.ErrUserNotFound
    }
    
    return user, err
}

func (r *postgresUserRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
    query := `
        SELECT id, name, email, age, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []*domain.User
    for rows.Next() {
        user := &domain.User{}
        if err := rows.Scan(
            &user.ID, &user.Name, &user.Email, &user.Age,
            &user.CreatedAt, &user.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    return users, rows.Err()
}

func (r *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {
    query := `
        UPDATE users
        SET name = $1, email = $2, age = $3, updated_at = NOW()
        WHERE id = $4
    `
    
    _, err := r.db.ExecContext(ctx, query,
        user.Name, user.Email, user.Age, user.ID,
    )
    
    return err
}

func (r *postgresUserRepository) Delete(ctx context.Context, id uint) error {
    query := `DELETE FROM users WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}
```

## Database-Specific Features

### PostgreSQL
- RETURNING clause support
- JSON/JSONB columns
- Array types
- Full-text search

### MySQL
- AUTO_INCREMENT handling
- LIMIT/OFFSET pagination
- JSON column support

### MongoDB
- Document-based storage
- Flexible schema
- Aggregation pipelines
- Index management

### SQLite
- Embedded database
- File-based storage
- Great for testing
- Simplified queries

## Best Practices

###  DO

- Keep repositories simple (CRUD + specific queries)
- Return domain entities
- Handle database errors properly
- Use prepared statements
- Implement transactions when needed

###  DON'T

- Include business logic
- Return database-specific types
- Expose SQL to callers
- Skip error handling

## See Also

- [`goca entity`](/commands/entity) - Generate entities
- [`goca usecase`](/commands/usecase) - Generate use cases
- [`goca feature`](/commands/feature) - Generate complete feature
