---
layout: doc
title: goca repository
titleTemplate: Commands | Goca
description: Generate repository interfaces and concrete database implementations for PostgreSQL, MySQL, SQLite, and MongoDB.
---

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

**Options:** 
- `postgres` - PostgreSQL (GORM)
- `postgres-json` - PostgreSQL with JSONB support
- `mysql` - MySQL — shares the `postgres` GORM implementation (`postgres_<entity>_repository.go`); only the dialector in `main.go` differs
- `mongodb` - MongoDB (native driver)
- `sqlite` - SQLite (embedded) — also shares the `postgres` GORM implementation, same as `mysql`
- `sqlserver` - SQL Server (GORM)
- `elasticsearch` - Elasticsearch (v8 client)
- `dynamodb` - DynamoDB (AWS SDK v2)

```bash
goca repository Product --database postgres
goca repository Config --database postgres-json
goca repository Article --database elasticsearch
```

### `--cache`

Generate a Redis cache decorator for the repository. Creates a `Cached<Entity>Repository` that wraps the database implementation with Redis caching.

```bash
goca repository Product --database postgres --cache
```

When enabled, this generates:
- `internal/repository/cached_product_repository.go` — cache decorator
- `internal/cache/redis.go` — Redis client factory

The cache decorator:
- **Caches** `FindByID` and `FindAll` results in Redis
- **Invalidates** cache on `Save`, `Update`, and `Delete`
- **Delegates** search methods directly to the underlying repository

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

### PostgreSQL JSON (Semi-structured Data)

```bash
goca repository Config --database postgres-json
```

### MongoDB Repository

```bash
goca repository Product --database mongodb
```

### Elasticsearch Full-Text Search

```bash
goca repository Article --database elasticsearch
```

### DynamoDB (AWS Serverless)

```bash
goca repository Order --database dynamodb
```

### SQL Server (Enterprise)

```bash
goca repository Employee --database sqlserver
```

### SQLite (Embedded)

```bash
goca repository Setting --database sqlite
```

### Interface Only

```bash
goca repository Order --interface-only
```

### With Redis Cache

```bash
goca repository Product --database postgres --cache
```

Generates a decorator pattern:

```
internal/repository/
├── interfaces.go                       # Repository interface
├── postgres_product_repository.go      # Database implementation
└── cached_product_repository.go        # Redis cache decorator
internal/cache/
└── redis.go                            # Redis client factory
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

import "myproject/internal/domain"

type UserRepository interface {
    Save(user *domain.User) error
    FindByID(id int) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id int) error
    FindAll() ([]domain.User, error)
}

// postgres_user_repository.go — used for postgres, mysql AND sqlite alike;
// all three run on GORM and share this one implementation (only the
// dialector passed to gorm.Open in main.go differs between them).
package repository

import (
    "myproject/internal/domain"
    "gorm.io/gorm"
)

type postgresUserRepository struct {
    db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) UserRepository {
    return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Save(user *domain.User) error {
    return r.db.Create(user).Error
}

func (r *postgresUserRepository) FindByID(id int) (*domain.User, error) {
    var user domain.User
    if err := r.db.First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *postgresUserRepository) Update(user *domain.User) error {
    return r.db.Save(user).Error
}

func (r *postgresUserRepository) Delete(id int) error {
    return r.db.Delete(&domain.User{}, id).Error
}

func (r *postgresUserRepository) FindAll() ([]domain.User, error) {
    var users []domain.User
    if err := r.db.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}
```

There is no `context.Context` parameter on any of these methods — the generated repository is synchronous.

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
