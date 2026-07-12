---
layout: doc
title: goca di
titleTemplate: Commands | Goca
description: Generate a dependency injection container that automatically wires all Clean Architecture layers.
---

# goca di

Generate dependency injection container for automatic wiring.

## Syntax

```bash
goca di [flags]
```

## Description

Creates a dependency injection container that automatically wires all features, repositories, use cases, and handlers together.

## Flags

### `--features` / `-f`

**Required.** Comma-separated list of features to wire.

```bash
goca di --features "User,Product,Order"
```

### `--database` / `-d`

Database type (`postgres`, `mysql`, `mongodb`). Defaults to the project's configured database (`.goca.yaml`), falling back to `postgres` if there is no config.

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

### `--wire` / `-w`

Generate a [Google Wire](https://github.com/google/wire)-based container (`internal/di/wire.go` + `wire_container.go`) instead of the default manual-wiring container. This only emits the `wire.Build(...)` provider-set annotations Wire's own code generator consumes — it does not require the `wire` binary to be installed to run `goca di --wire` itself, only to later run `wire` against the generated file.

```bash
goca di --features "User,Product" --wire
```

### `--dry-run`, `--force`, `--backup`

Standard safety flags: preview without writing, overwrite without asking, and back up an existing `container.go` before overwriting, respectively.

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

::: warning Illustrative only
The exact type/method names below have drifted from current output (e.g. the real use case type is `usecase.UserUseCase`/`NewUserService`, the container holds a `*gorm.DB` not `*sql.DB`, and there are also `Get<Entity>UseCase()`-style getters). The overall shape — repositories, use cases and handlers wired together in one `Container` — is accurate; generate a container yourself with `goca di` to see the current exact code.
:::

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
