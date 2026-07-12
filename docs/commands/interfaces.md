---
layout: doc
title: goca interfaces
titleTemplate: Commands | Goca
description: Generate interface contracts for use-case and repository layers to enable Test-Driven Development and mock generation.
---

# goca interfaces

Generate interface contracts for Test-Driven Development.

## Syntax

```bash
goca interfaces <EntityName> [flags]
```

## Description

Generates only the interface contracts without implementations, perfect for TDD workflows where you define contracts first.

## Flags

### `--usecase`

Generate use case interfaces.

```bash
goca interfaces Product --usecase
```

### `--repository`

Generate repository interfaces.

```bash
goca interfaces User --repository
```

### `--all`

Generate all interfaces.

```bash
goca interfaces Order --all
```

## Examples

### Use Case Interfaces

```bash
goca interfaces Product --usecase
```

**Generates:** `internal/interfaces/product_usecase.go` (note: a separate `internal/interfaces` package — not `internal/usecase`)

```go
package interfaces

import "myproject/internal/domain"

// Product UseCase DTOs
type CreateProductInput interface {
    GetName() string
    Validate() error
}

type CreateProductOutput interface {
    GetProduct() domain.Product
    GetMessage() string
}
// ...similar Update/Delete/List DTO interfaces
```

No `context.Context` parameters — generated methods are synchronous.

### Repository Interfaces

```bash
goca interfaces User --repository
```

**Generates:** `internal/interfaces/user_repository.go` (also under `internal/interfaces`, not `internal/repository`)

```go
package interfaces

import (
    "myproject/internal/domain"
    "gorm.io/gorm"
)

type UserRepository interface {
    Save(user *domain.User) error
    FindByID(id int) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id int) error
    FindAll() ([]domain.User, error)
    Count() (int, error)
    Exists(id int) (bool, error)
    SaveBatch(users []domain.User) error
    DeleteBatch(ids []int) error
    SaveWithTx(tx *gorm.DB, user *domain.User) error
    UpdateWithTx(tx *gorm.DB, user *domain.User) error
    DeleteWithTx(tx *gorm.DB, id int) error
}
```

### All Interfaces

```bash
goca interfaces Order --all
```

## TDD Workflow

1. **Generate Interfaces:**
   ```bash
   goca interfaces Payment --all
   ```

2. **Write Tests:**
   ```go
   func TestPaymentService_CreatePayment(t *testing.T) {
       mockRepo := &MockPaymentRepository{}
       service := NewPaymentService(mockRepo)
       // ... test implementation
   }
   ```

3. **Implement:**
   Implement the actual service and repository.

## See Also

- [`goca usecase`](/commands/usecase) - Generate full use cases
- [`goca repository`](/commands/repository) - Generate repositories
