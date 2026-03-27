# goca di

Generate dependency injection container for automatic wiring.

## Syntax

```bash
goca di [flags]
```

## Description

Creates a dependency injection container that automatically wires all features, repositories, use cases, and handlers together.

## Flags

### `--features`

Comma-separated list of features to wire.

```bash
goca di --features "User,Product,Order"
```

### `--database`

Database type. Default: `postgres`

```bash
goca di --features "User,Product" --database postgres
```

### `--cache` / `-c`

Wire Redis cache decorators for all repositories. When enabled, each repository is wrapped with a `Cached<Entity>Repository` that provides Redis-backed read caching.

```bash
goca di --features "User,Product" --cache
```

The generated container accepts a `*redis.Client` in addition to `*gorm.DB`:

```go
container := di.NewContainer(db, redisClient)
```

## Examples

### Wire All Features

```bash
goca di --features "User,Product,Order,Payment"
```

### PostgreSQL with Authentication

```bash
goca di --features "User,Auth" --database postgres
```

## Generated Code

```go
// internal/di/container.go
package di

import (
    "database/sql"
    "myproject/internal/handler/http"
    "myproject/internal/repository"
    "myproject/internal/usecase"
)

type Container struct {
    // Repositories
    UserRepository    repository.UserRepository
    ProductRepository repository.ProductRepository
    
    // Use Cases
    UserService    usecase.UserService
    ProductService usecase.ProductService
    
    // Handlers
    UserHandler    *http.UserHandler
    ProductHandler *http.ProductHandler
}

func NewContainer(db *sql.DB) *Container {
    // Initialize repositories
    userRepo := repository.NewPostgresUserRepository(db)
    productRepo := repository.NewPostgresProductRepository(db)
    
    // Initialize use cases
    userService := usecase.NewUserService(userRepo)
    productService := usecase.NewProductService(productRepo)
    
    // Initialize handlers
    userHandler := http.NewUserHandler(userService)
    productHandler := http.NewProductHandler(productService)
    
    return &Container{
        UserRepository:    userRepo,
        ProductRepository: productRepo,
        UserService:       userService,
        ProductService:    productService,
        UserHandler:       userHandler,
        ProductHandler:    productHandler,
    }
}
```

## See Also

- [`goca integrate`](/commands/integrate) - Integrate existing features
- [`goca feature`](/commands/feature) - Generate complete feature
