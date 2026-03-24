---
applyTo: "internal/**/*.go"
---

# Clean Architecture Rules for Goca Internal Packages

## Layer Dependency Direction

The dependency rule is strict and non-negotiable:

```
Handler → UseCase (interface) → Repository (interface) → Entity
```

- Handlers import `usecase` interfaces — NEVER `repository` or `domain` directly for business logic
- UseCases import `repository` interfaces — NEVER import handler packages
- Repositories import `domain` — NEVER import usecase or handler packages
- Domain (entities) imports NOTHING from internal packages

### Import Violations — Immediately Fix

```go
// FORBIDDEN in usecase/ package
import "github.com/sazardev/goca/internal/handler/http"

// FORBIDDEN in domain/ package
import "github.com/sazardev/goca/internal/usecase"

// FORBIDDEN in repository/ package
import "github.com/sazardev/goca/internal/handler/http"
```

## Entity (Domain) Layer Rules

Entities in `internal/domain/` MUST:

- Be pure Go structs with no framework imports (no GORM, no validator libs in actual business logic methods)
- GORM struct tags are allowed for ORM mapping (they are annotations, not dependencies)
- Have business validation methods that return domain errors
- Use domain-specific error variables (`var ErrProductNotFound = errors.New("product not found")`)

```go
// CORRECT entity
type Product struct {
    ID    uint    `json:"id" gorm:"primaryKey;autoIncrement"`
    Name  string  `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
    Price float64 `json:"price" gorm:"type:decimal(10,2)" validate:"required,gte=0"`
}

func (p *Product) Validate() error {
    if p.Name == "" {
        return ErrProductNameRequired
    }
    if p.Price < 0 {
        return ErrProductPriceNegative
    }
    return nil
}
```

## UseCase Layer Rules

UseCase files in `internal/usecase/` MUST:

- Define an interface (e.g., `ProductUseCase`) AND implement it (`productService` struct — unexported)
- Accept repository interfaces in the constructor — never concrete implementations
- Define DTOs (Input/Output types) in `dto.go` — never expose domain entities directly in HTTP responses
- Constructor uses `New<Entity>Service(repo repository.Interface) UseCase` pattern

```go
// CORRECT usecase pattern
type productService struct {
    repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductUseCase {
    return &productService{repo: repo}
}

func (s *productService) Create(input CreateProductInput) (*domain.Product, error) {
    product := input.toDomain()
    if err := product.Validate(); err != nil {
        return nil, fmt.Errorf("create product: %w", err)
    }
    return s.repo.Create(product)
}
```

## Repository Layer Rules

Repository files in `internal/repository/` define interfaces ONLY. Implementations go in a sub-package named after the database (e.g., `postgres/`, `sqlite/`):

```go
// internal/repository/interfaces.go — interface definition
type ProductRepository interface {
    Create(product *domain.Product) (*domain.Product, error)
    GetByID(id uint) (*domain.Product, error)
    List() ([]domain.Product, error)
    Update(product *domain.Product) (*domain.Product, error)
    Delete(id uint) error
}
```

## Handler Layer Rules

Handlers in `internal/handler/http/` MUST:

- Accept UseCase interfaces in the constructor
- Never contain business logic — delegate 100% to UseCase
- Handle HTTP concerns only: status codes, request parsing, response encoding
- Use `json.NewDecoder` for request bodies — never `ioutil.ReadAll` + Unmarshal

```go
// CORRECT handler
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
    var input usecase.CreateProductInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }
    product, err := h.usecase.Create(input)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(product)
}
```

## DI Container Rules

`internal/di/container.go` MUST:

- Wire ALL dependencies using constructor injection
- Accept `*gorm.DB` (or relevant DB connection) as the single infrastructure input
- Never instantiate concrete implementations in non-`di` packages
- Expose getters for each handler (used by `main.go`)

```go
type Container struct {
    productHandler *httphandler.ProductHandler
    // ...
}

func NewContainer(db *gorm.DB) *Container {
    productRepo := repository.NewProductGORMRepository(db)
    productUseCase := usecase.NewProductService(productRepo)
    productHandler := httphandler.NewProductHandler(productUseCase)
    return &Container{productHandler: productHandler}
}
```

## Testing in `internal/testing/`

- `test_framework.go` — base test utilities and shared helpers
- `suite.go` — test suite definitions (optional, for grouped integration tests)
- `validator.go` — helpers to assert generated code structure
- `tests/` — integration tests that invoke `goca` commands and verify output
- All integration tests use `t.TempDir()` — never write to the real workspace

## Seed Data Files (`*_seeds.go`)

Seed files are only for local development and testing. They MUST:

- Be in the `domain` package
- Not import any test packages in production builds (use build tags if needed)
- Return deterministic data — no random values without a fixed seed
