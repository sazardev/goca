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

**Generates:** `internal/usecase/product_interfaces.go`

```go
package usecase

import "context"

type ProductService interface {
    CreateProduct(ctx context.Context, input CreateProductInput) (*ProductResponse, error)
    GetProduct(ctx context.Context, id uint) (*ProductResponse, error)
    UpdateProduct(ctx context.Context, id uint, input UpdateProductInput) error
    DeleteProduct(ctx context.Context, id uint) error
    ListProducts(ctx context.Context) ([]*ProductResponse, error)
}
```

### Repository Interfaces

```bash
goca interfaces User --repository
```

**Generates:** `internal/repository/user_interfaces.go`

```go
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
