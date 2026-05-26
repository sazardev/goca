---
name: ddd-teaching-skill
description: >
  Enseña Domain-Driven Design y Clean Architecture de forma guiada usando Goca CLI.
  Actívar cuando el usuario pregunte "aprender DDD", "como usar Goca", "enseñame Clean Architecture",
  "que es domain", "como generar entity/usecase/repository/handler", o pida practicar DDD con Goca.
  Inactívar al completar el flujo o si el usuario cambia a tareas no relacionadas con aprendizaje.
---

# DDD Teaching — Aprendizaje Guiado con Goca

## Filosofía

DDD + Clean Architecture resuelven el problema de **código que muere lentamente**: lógica de negocio mezclada con frameworks, tests imposibles, cambios que cascada. La solución es la **Dependency Rule** — las dependencias apuntan hacia adentro.

Goca materializa esta arquitectura generando código en 4 capas:

```
Handler (adapter)  →  UseCase (app logic)  →  Repository (persistence)  →  Domain (pure)
     ↓                       ↓                        ↓                        ↓
   HTTP/gRPC           business logic              SQL/Redis               entities
```

Cada capa solo conoce la que está inmediatamente **dentro** de ella, y siempre a través de interfaces.

## Activación

Invocar este skill automáticamente cuando el usuario:
- Pregunte conceptos de DDD o Clean Architecture
- Quiera aprender a usar Goca paso a paso
- Pida "generar un feature completo" sin entender qué genera
- Muestre código con violaciones arquitectónicas (handler con GORM, usecase con http, etc.)

## Mapa Conceptual ↔ Comandos Goca

| Concepto DDD | Comando Goca | Archivo generado | Propósito |
|---|---|---|---|
| Entidad | `goca entity` | `internal/domain/product.go` | Pureza del dominio, invariantes |
| Value Object | `goca entity --validation` | `domain/product.go` — Validate() | Validación encapsulada |
| Repository interface | `goca repository --interface-only` | `repository/interfaces.go` | Contrato de persistencia |
| Repository impl | `goca repository -d postgres` | `repository/postgres_product.go` | GORM encapsulado |
| DTO / Application Service | `goca usecase` | `usecase/dto.go`, `service.go` | Separación capas |
| UseCase interface | `goca usecase` | `usecase/product_usecase.go` | Contrato para handlers |
| Handler / Adapter | `goca handler -t http` | `handler/http/product_handler.go` | Delivery |
| DI Container | `goca di` | `di/container.go` | Wiring |

---

## Flujo Guiado — 6 Pasos (con TDD)

Cada paso: **Concepto DDD** → **Comando Goca** → **Código generado** → **Test**

### Paso 1: Init — Scaffold del Proyecto

**Concepto:** Separación en capas limpias. El directorio `internal/` es la frontera — nada fuera de `internal/` puede importar lo de adentro, pero lo de adentro sí puede importar `pkg/`.

**Comando:**
```bash
goca init ecommerce --module github.com/myapp/ecommerce --database postgres
```

**Estructura generada:**
```
internal/
  domain/       ← Puro: sin imports externos, solo lógica de negocio
  usecase/      ← Solo importa domain + repository interfaces
  repository/   ← Implementa interfaces; conoce GORM, no el usecase
  handler/      ← Conoce interfaces de usecase, nunca repository directo
```

**Análisis:** Cada carpeta corresponde exactamente a un círculo de Clean Architecture. La regla de dependencia se aplica a nivel de import: `handler` → `usecase` → `repository` → `domain`. Prohibido saltar capas.

**Test:**
```bash
cd ecommerce && go build ./... && go vet ./...
# Debe pasar sin errores: es el esqueleto vacío pero correcto
```

---

### Paso 2: Entity — Domain Puro

**Concepto:** Una **entidad DDD** tiene identidad (ID) e **invariantes** — reglas que siempre deben cumplirse. Los **Value Objects** son inmutables y se comparan por valor. Ambos viven en `domain/` sin dependencias externas.

**Comando:**
```bash
goca entity Product --fields "name:string,price:float64,category:string" --validation --business-rules
```

**Código generado** (extraído de `cmd/templates.go` y `cmd/template_components.go`):
```go
// internal/domain/product.go
package domain

type Product struct {
    ID       int     `json:"id" gorm:"primaryKey"`
    Name     string  `json:"name" gorm:"type:varchar(255);not null"`
    Price    float64 `json:"price" gorm:"type:decimal(10,2);not null"`
    Category string  `json:"category" gorm:"type:varchar(100)"`
}

func (p *Product) Validate() error {
    if p.Name == "" {
        return errors.New("product name is required")
    }
    if p.Price <= 0 {
        return errors.New("product price must be positive")
    }
    return nil
}

func (p *Product) IsExpensive() bool {
    return p.Price > 1000
}
```

```go
// internal/domain/errors.go
package domain

import "errors"

var (
    ErrInvalidProductName  = errors.New("product name is required")
    ErrInvalidProductPrice = errors.New("product price must be positive")
)
```

**Análisis:**
- Struct `Product`: sin lógica externa. Los tags de campo son metadata de serialización, no comportamiento.
- `Validate()`: invariante expresado como **método del dominio**. Sin dependencias externas. Sin ORM. Sin HTTP.
- `IsExpensive()`: regla de negocio co-localizada con los datos. Cambia la regla? Cambias este método, no un service remoto.
- Errores como **sentinel values** (`var Err...`), no strings mágicos. El consumidor puede hacer `errors.Is(err, domain.ErrInvalidProductName)`.

**Test (TDD primero):**
```go
// internal/domain/product_test.go
func TestProduct_Validate(t *testing.T) {
    tests := []struct {
        name    string
        product Product
        wantErr bool
    }{
        {"valid product", Product{Name: "Laptop", Price: 999.99, Category: "Electronics"}, false},
        {"empty name", Product{Price: 999.99}, true},
        {"zero price", Product{Name: "Laptop", Price: 0}, true},
        {"negative price", Product{Name: "Laptop", Price: -1}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.product.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
            }
        })
    }
}

func TestProduct_IsExpensive(t *testing.T) {
    cheap := Product{Name: "Notebook", Price: 10}
    expensive := Product{Name: "MacBook", Price: 2500}

    if cheap.IsExpensive() {
        t.Error("expected cheap product to not be expensive")
    }
    if !expensive.IsExpensive() {
        t.Error("expected expensive product to be expensive")
    }
}
```

---

### Paso 3: UseCase — Lógica de Aplicación

**Concepto:** El **Application Service** orquesta la lógica de negocio. No contiene reglas de dominio (esas van en la entidad). Usa **DTOs** para desacoplar el mundo exterior del dominio. Depende de **interfaces de repositorio**, no de implementaciones concretas.

**Comando:**
```bash
goca usecase ProductService --entity Product --operations "create,read,update,delete,list"
```

**Código generado** (extraído de `cmd/template_components.go` + `cmd/templates.go`):
```go
// internal/usecase/product_usecase.go
package usecase

import (
    "github.com/myapp/ecommerce/internal/domain"
    "github.com/myapp/ecommerce/internal/repository"
)

type ProductUseCase interface {
    CreateProduct(input CreateProductInput) (*CreateProductOutput, error)
    GetProductByID(id int) (*ProductOutput, error)
    UpdateProduct(id int, input UpdateProductInput) error
    DeleteProduct(id int) error
    ListProducts() (*ListProductsOutput, error)
}

type productService struct {
    repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductUseCase {
    return &productService{repo: repo}
}

func (s *productService) CreateProduct(input CreateProductInput) (*CreateProductOutput, error) {
    product := domain.Product{
        Name:     input.Name,
        Price:    input.Price,
        Category: input.Category,
    }

    if err := product.Validate(); err != nil {
        return nil, err
    }

    if err := s.repo.Save(&product); err != nil {
        return nil, err
    }

    return &CreateProductOutput{
        Product: &product,
        Message: "Product created successfully",
    }, nil
}
```

```go
// internal/usecase/dto.go
type CreateProductInput struct {
    Name     string  `json:"name"`
    Price    float64 `json:"price"`
    Category string  `json:"category"`
}

type CreateProductOutput struct {
    Product *domain.Product `json:"product"`
    Message string          `json:"message"`
}

type UpdateProductInput struct {
    Name     *string  `json:"name,omitempty"`
    Price    *float64 `json:"price,omitempty"`
    Category *string  `json:"category,omitempty"`
}

type ProductOutput struct {
    Product *domain.Product `json:"product"`
}

type ListProductsOutput struct {
    Products []domain.Product `json:"products"`
    Total    int              `json:"total"`
}
```

```go
// internal/usecase/interfaces.go
package usecase

import "github.com/myapp/ecommerce/internal/domain"

type ProductRepository interface {
    Save(product *domain.Product) error
    FindByID(id int) (*domain.Product, error)
    Update(product *domain.Product) error
    Delete(id int) error
    FindAll() ([]domain.Product, error)
}
```

**Análisis:**
- `ProductUseCase` es **interfaz pública**. El handler programa contra esta interfaz, no contra el struct concreto.
- `productService` es **privado** (minúscula). Solo se exporta el constructor `NewProductService(repo)`. Nadie puede acoplar al tipo concreto.
- **Constructor injection**: el repo se inyecta. El service no crea su propio repo, no llama a `sql.Open()`, no sabe si es Postgres o Mock.
- `CreateProductInput` vs `domain.Product`: el input puede tener validaciones distintas al domain. El Update usa **punteros** (`*string`) para distinguir "no enviado" de "enviado vacío".
- La interfaz `ProductRepository` en `usecase/interfaces.go` es la **misma** que en `repository/interfaces.go`. El usecase la necesita para su firma.

**Test con Mock (TDD — RED antes de implementar):**
```go
// internal/usecase/product_service_test.go
func TestCreateProduct_Success(t *testing.T) {
    mockRepo := new(MockProductRepository)
    mockRepo.On("Save", mock.AnythingOfType("*domain.Product")).Return(nil)

    svc := NewProductService(mockRepo)
    input := CreateProductInput{Name: "Laptop", Price: 999.99, Category: "Electronics"}

    output, err := svc.CreateProduct(input)

    assert.NoError(t, err)
    assert.NotNil(t, output.Product)
    assert.Equal(t, "Laptop", output.Product.Name)
    assert.Equal(t, "Product created successfully", output.Message)
    mockRepo.AssertExpectations(t)
}

func TestCreateProduct_InvalidData(t *testing.T) {
    mockRepo := new(MockProductRepository)

    svc := NewProductService(mockRepo)
    input := CreateProductInput{Name: "", Price: 0}

    _, err := svc.CreateProduct(input)

    assert.Error(t, err)
    mockRepo.AssertNotCalled(t, "Save")
}
```

---

### Paso 4: Repository — Persistencia

**Concepto:** El **Repository Pattern** abstrae el almacenamiento detrás de una interfaz. El dominio y el usecase conocen la interfaz, no la implementación. GORM, SQL, Redis, MongoDB — todo queda encapsulado detrás de esta interfaz.

**Comando:**
```bash
goca repository Product --database postgres --transactions
```

**Código generado** (extraído de `cmd/templates.go` + `cmd/repository.go`):
```go
// internal/repository/interfaces.go
package repository

import "github.com/myapp/ecommerce/internal/domain"

type ProductRepository interface {
    Save(product *domain.Product) error
    FindByID(id int) (*domain.Product, error)
    Update(product *domain.Product) error
    Delete(id int) error
    FindAll() ([]domain.Product, error)
}
```

```go
// internal/repository/postgres_product_repository.go
package repository

import (
    "gorm.io/gorm"
    "github.com/myapp/ecommerce/internal/domain"
)

type postgresProductRepository struct {
    db *gorm.DB
}

func NewPostgresProductRepository(db *gorm.DB) ProductRepository {
    return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) Save(product *domain.Product) error {
    return r.db.Create(product).Error
}

func (r *postgresProductRepository) FindByID(id int) (*domain.Product, error) {
    var product domain.Product
    if err := r.db.First(&product, id).Error; err != nil {
        return nil, err
    }
    return &product, nil
}

func (r *postgresProductRepository) Update(product *domain.Product) error {
    return r.db.Save(product).Error
}

func (r *postgresProductRepository) Delete(id int) error {
    return r.db.Delete(&domain.Product{}, id).Error
}

func (r *postgresProductRepository) FindAll() ([]domain.Product, error) {
    var products []domain.Product
    if err := r.db.Find(&products).Error; err != nil {
        return nil, err
    }
    return products, nil
}
```

**Análisis:**
- Interfaz en `repository/interfaces.go` define el **contrato**. El usecase la consume. El handler ni la conoce.
- `postgresProductRepository` es **privado** — nadie fuera del package la instancia excepto el DI container.
- El constructor `NewPostgresProductRepository(db *gorm.DB) ProductRepository` devuelve la **interfaz**, no el struct.
- GORM está **completamente encapsulado** en este package. Si migras a MongoDB, cambias este archivo, no tocas domain ni usecase.
- `--transactions` agrega `SaveWithTx(tx *gorm.DB, ...)` para Unit of Work.

---

### Paso 5: Handler — Delivery Adapter

**Concepto:** El **Adapter** convierte requests externos (HTTP, gRPC, CLI) en llamadas a usecase. No contiene lógica de negocio. Solo serialización, routing, delegación.

**Comando:**
```bash
goca handler Product --type http --validation
```

**Código generado** (extraído de `cmd/templates.go`):
```go
// internal/handler/http/product_handler.go
package http

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/myapp/ecommerce/internal/usecase"
)

type ProductHandler struct {
    usecase usecase.ProductUseCase
}

func NewProductHandler(uc usecase.ProductUseCase) *ProductHandler {
    return &ProductHandler{usecase: uc}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
    var input usecase.CreateProductInput

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    output, err := h.usecase.CreateProduct(input)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(output)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    output, err := h.usecase.GetProductByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(output)
}
```

```go
// internal/handler/http/routes.go
package http

import (
    "github.com/gorilla/mux"
    "github.com/myapp/ecommerce/internal/usecase"
)

func SetupProductRoutes(router *mux.Router, uc usecase.ProductUseCase) {
    handler := NewProductHandler(uc)

    router.HandleFunc("/products", handler.CreateProduct).Methods("POST")
    router.HandleFunc("/products/{id}", handler.GetProduct).Methods("GET")
    router.HandleFunc("/products/{id}", handler.UpdateProduct).Methods("PUT")
    router.HandleFunc("/products/{id}", handler.DeleteProduct).Methods("DELETE")
    router.HandleFunc("/products", handler.ListProducts).Methods("GET")
}
```

**Análisis:**
- Handler importa `usecase` **solo**. No importa `repository`, no importa `gorm`. ✅
- `ProductHandler` recibe `usecase.ProductUseCase` (interfaz), no el service concreto.
- No hay lógica de negocio aquí. Solo: parse request → call usecase → serialize response.
- Si cambias HTTP → gRPC, el handler cambia pero `usecase` y `domain` no se tocan.
- `Routes` se inyecta el usecase y construye el handler — el router no necesita saber cómo construirlo.

**Test con httptest:**
```go
func TestCreateProductHandler(t *testing.T) {
    mockUC := new(MockProductUseCase)
    mockUC.On("CreateProduct", mock.Anything).Return(&usecase.CreateProductOutput{
        Product: &domain.Product{Name: "Laptop", Price: 999.99},
        Message: "Product created successfully",
    }, nil)

    handler := NewProductHandler(mockUC)
    body := `{"name":"Laptop","price":999.99,"category":"Electronics"}`
    req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    handler.CreateProduct(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
    assert.Contains(t, w.Body.String(), "Laptop")
    mockUC.AssertExpectations(t)
}
```

---

### Paso 6: DI — Wiring Completo

**Concepto:** El **Composition Root** es el único lugar donde se instancian implementaciones concretas. Construye el grafo de objetos completo. Todas las demás capas reciben sus dependencias ya construidas.

**Comando:**
```bash
goca di --features "Product" --database postgres
```

**Código generado** (extraído de `cmd/di.go`):
```go
// internal/di/container.go
package di

import (
    "gorm.io/gorm"

    "github.com/myapp/ecommerce/internal/repository"
    "github.com/myapp/ecommerce/internal/usecase"
    "github.com/myapp/ecommerce/internal/handler/http"
)

type Container struct {
    db *gorm.DB

    productRepo    repository.ProductRepository
    productUC      usecase.ProductUseCase
    productHandler *http.ProductHandler
}

func NewContainer(db *gorm.DB) *Container {
    c := &Container{db: db}
    c.setupRepositories()
    c.setupUseCases()
    c.setupHandlers()
    return c
}

func (c *Container) setupRepositories() {
    c.productRepo = repository.NewPostgresProductRepository(c.db)
}

func (c *Container) setupUseCases() {
    c.productUC = usecase.NewProductService(c.productRepo)
}

func (c *Container) setupHandlers() {
    c.productHandler = http.NewProductHandler(c.productUC)
}

func (c *Container) ProductHandler() *http.ProductHandler {
    return c.productHandler
}
```

**Análisis:**
- `NewContainer(db)` recibe la conexión a BD. El container construye TODO el grafo.
- **Único lugar** donde se llama a `NewPostgresProductRepository`, `NewProductService`, `NewProductHandler`.
- Si quieres cambiar Postgres → MySQL, cambias UNA línea en `setupRepositories`.
- Los getters (`ProductHandler()`) exponen solo lo que `main.go` necesita.

---

## Feature Completo (atajo)

```bash
goca feature Product --fields "name:string,price:float64,category:string" --database postgres --validation --handlers http
```

Este comando ejecuta Pasos 2-6 en una sola operación. Equivale a:

```
1. goca entity Product --fields "name:string,price:float64,category:string" --validation --business-rules
2. goca usecase ProductService --entity Product --operations "create,read,update,delete,list"
3. goca repository Product --database postgres
4. goca handler Product --type http --validation
5. goca messages Product --all
6. goca di --features "Product" --database postgres
```

Además integra el handler en `main.go` y registra migraciones automáticas.

---

## Anti-Patterns (qué NO hacer)

❌ **Handler importando repository directamente**
```go
type ProductHandler struct {
    repo *gorm.DB  // MAL: handler conoce la BD y GORM!
}
```
✅ Handler solo conoce `usecase.ProductUseCase`

❌ **UseCase con GORM o SQL**
```go
type productService struct {
    db *gorm.DB  // MAL: lógica de negocio conoce GORM
}
```
✅ UseCase recibe `repository.ProductRepository` (interfaz abstracta)

❌ **Entidad anémica (solo getters/setters, sin comportamiento)**
```go
type Product struct {
    Name  string  // MAL: sin Validate(), sin reglas de negocio
    Price float64 // solo es una struct de datos, no una entidad DDD
}
```
✅ Domain tiene métodos: `Validate()`, `IsExpensive()`, etc.

❌ **DTOs expuestos desde handler** — usa `usecase.CreateProductInput`, no crees DTOs en handler.

❌ **init() y variables globales** — usa constructor injection siempre.

---

## Límites del Skill

Este skill enseña DDD + Clean Architecture a través de Goca. Para otros aspectos, delegar:

| Necesitas | Recurso |
|---|---|
| Benchmarks, fuzzing, test coverage avanzado | golang-testing skill |
| Table-driven tests, helpers, golden files | golang-testing skill |
| sync.Pool, zero allocation, performance | golang-performance skill |
| Interface design, functional options, patterns | golang-patterns skill |
| Validar dependencias entre capas archivo por archivo | ArchitectGuard agent mode en `.github/AGENTS.md` |
| Verificar que templates generan código compilable | CodegenAuditor agent mode en `.github/AGENTS.md` |
| Guía completa de comandos Goca | `GUIDE.md` |

Cuando el usuario cambie a una tarea no relacionada con aprendizaje DDD (ej: "arregla este bug", "agrega esta feature"), desactivar este skill y continuar como agente normal.
