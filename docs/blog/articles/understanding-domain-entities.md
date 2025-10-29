---
layout: doc
title: Understanding Domain Entities in Clean Architecture
titleTemplate: Articles | Goca Blog
description: A comprehensive guide to domain entities, their role in Clean Architecture, and how Goca generates production-ready entities following DDD principles
tags:
  - Domain Entities
  - Clean Architecture
  - DDD
  - Go
  - Goca
---

<script setup>
import Badge from '../../.vitepress/theme/components/Badge.vue'
</script>

# Understanding Domain Entities in Clean Architecture

<div style="display: flex; gap: 0.5rem; margin-bottom: 1rem;">
<Badge type="info">Architecture</Badge>
<Badge type="tip">Domain-Driven Design</Badge>
</div>

Domain entities are the heart of Clean Architecture. They represent the core business concepts and rules that define your application's purpose. Understanding entities correctly is fundamental to building maintainable, testable, and scalable software systems.

---

## What is a Domain Entity?

A domain entity is a representation of a business concept that has a unique identity and encapsulates business rules. In Clean Architecture, entities form the innermost layer, completely independent of external concerns like databases, frameworks, or UI.

### Core Characteristics

**Identity**: Each entity has a unique identifier that distinguishes it from other entities of the same type. Two entities with the same attributes but different identities are different entities.

**Business Logic**: Entities contain methods that enforce business rules and maintain invariants. They are not passive data structures but active participants in your domain model.

**Independence**: Entities have zero dependencies on external systems. They do not import HTTP libraries, database drivers, or framework code. This independence makes them portable, testable, and reusable.

**Validation**: Entities validate their own state, ensuring that business rules are never violated. Invalid states are impossible to represent.

## Entity vs Model: A Critical Distinction

Many developers confuse entities with database models or API models. This confusion leads to architectural problems and coupling.

### What an Entity Is NOT

**Not a Database Model**: Entities do not map directly to database tables. They represent business concepts, not storage structures. Database concerns belong to the infrastructure layer.

**Not an API Response**: Entities are not DTOs (Data Transfer Objects). API responses should be separate structures that adapt entities for external communication.

**Not Framework-Dependent**: Entities do not depend on ORMs, validation frameworks, or serialization libraries. These are implementation details.

### The Separation Principle

```
Domain Entity (Pure Business Logic)
        ↓
    Use Case (Application Logic)
        ↓
Repository Interface (Contract)
        ↓
Repository Implementation (Database Details)
        ↓
    Database Schema
```

This separation allows you to:

- Change databases without touching business logic
- Test business rules without database setup
- Evolve your domain model independently
- Swap ORMs or frameworks with minimal impact

## Domain-Driven Design Principles

Goca implements Domain-Driven Design (DDD) principles when generating entities, ensuring your code follows established best practices.

### Ubiquitous Language

Entities use the same terminology as your business domain. If your business talks about "Orders," "Customers," and "Products," your entities should use these exact terms.

### Aggregate Roots

Entities can serve as aggregate roots, controlling access to related objects and maintaining consistency boundaries.

### Value Objects vs Entities

Entities have identity; value objects do not. An email address is a value object. A user is an entity. Goca helps you model both correctly.

## How Goca Generates Entities

Goca provides the `goca entity` command to generate domain entities following Clean Architecture and DDD principles.

### Basic Entity Generation

```bash
goca entity User --fields "name:string,email:string,age:int"
```

This command generates a pure domain entity with no external dependencies:

```go
package domain

type User struct {
    ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
    Name  string `json:"name" gorm:"type:varchar(255);not null"`
    Email string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
    Age   int    `json:"age" gorm:"type:integer;not null;default:0"`
}
```

Notice that while GORM tags are present for infrastructure convenience, the entity itself remains a simple Go struct. The entity does not import GORM or any database package.

### Field Types and Conventions

Goca supports common field types that map to both Go types and database columns:

**String Fields**: `name:string`, `email:string`, `description:string`
- Generate `string` type
- Map to `varchar` or `text` columns
- Suitable for textual data

**Numeric Fields**: `age:int`, `price:float64`, `quantity:int64`
- Generate integer or floating-point types
- Map to appropriate numeric columns
- Support business calculations

**Boolean Fields**: `is_active:bool`, `verified:bool`
- Generate `bool` type
- Map to boolean columns
- Represent binary states

**Temporal Fields**: `birth_date:time.Time`
- Generate `time.Time` type
- Handle date and time data
- Work with the standard library

### Adding Validation

Business rules are enforced through validation methods:

```bash
goca entity User --fields "name:string,email:string,age:int" --validation
```

This generates a `Validate()` method and domain-specific errors:

```go
package domain

import (
    "time"
    
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
    Name      string         `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
    Email     string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,email"`
    Age       int            `json:"age" gorm:"type:integer;not null;default:0" validate:"required,gte=0"`
    CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (u *User) Validate() error {
    if u.Name == "" {
        return ErrInvalidUserName
    }
    if u.Email == "" {
        return ErrInvalidUserEmail
    }
    if u.Age < 0 {
        return ErrInvalidUserAge
    }
    return nil
}
```

The validation method ensures that no invalid user can exist in your system. This is domain logic, not input validation. Input validation happens in the use case or handler layer.

### Domain Errors

Goca generates a separate `errors.go` file containing domain-specific errors:

```go
package domain

import "errors"

var (
    ErrInvalidUserName  = errors.New("invalid user name")
    ErrInvalidUserEmail = errors.New("invalid user email")
    ErrInvalidUserAge   = errors.New("invalid user age")
    ErrUserNotFound     = errors.New("user not found")
)
```

These errors are part of your domain model. They communicate business rule violations clearly and can be handled appropriately by outer layers.

### Business Rules

Beyond validation, entities can contain business logic:

```bash
goca entity Order --fields "customer_id:int,total:float64,status:string" --business-rules
```

This generates methods that implement domain logic:

```go
func (o *Order) Validate() error {
    if o.Customer_id < 0 {
        return ErrInvalidOrderCustomer_id
    }
    if o.Total < 0 {
        return ErrInvalidOrderTotal
    }
    if o.Status == "" {
        return ErrInvalidOrderStatus
    }
    return nil
}
```

You can extend these with additional business methods:

```go
func (o *Order) CanBeCancelled() bool {
    return o.Status == "pending" || o.Status == "confirmed"
}

func (o *Order) Apply(discount float64) error {
    if discount < 0 || discount > 1 {
        return errors.New("discount must be between 0 and 1")
    }
    o.Total = o.Total * (1 - discount)
    return nil
}

func (o *Order) IsCompleted() bool {
    return o.Status == "delivered"
}
```

These methods encapsulate business knowledge. They answer domain questions and enforce domain rules.

### Timestamps and Soft Deletes

Entities often need audit trails and soft delete functionality:

```bash
goca entity Product --fields "name:string,price:float64,stock:int" \
  --timestamps \
  --soft-delete
```

This generates:

```go
type Product struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"type:varchar(255);not null"`
    Price       float64        `json:"price" gorm:"type:decimal(10,2);not null;default:0"`
    Stock       int            `json:"stock" gorm:"type:integer;not null;default:0"`
    CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (p *Product) SoftDelete() {
    p.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
}

func (p *Product) IsDeleted() bool {
    return p.DeletedAt.Valid
}
```

Soft deletes preserve data while marking it as inactive. The `DeletedAt` field enables this pattern without permanently removing records.

## Complete Entity Example

Let's examine a complete entity generated by Goca:

```bash
goca entity Product --fields "name:string,description:string,price:float64,stock:int,is_active:bool" \
  --validation \
  --business-rules \
  --timestamps \
  --soft-delete
```

This generates a production-ready entity:

```go
package domain

import (
    "time"
    
    "gorm.io/gorm"
)

type Product struct {
    ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
    Name        string         `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
    Description string         `json:"description" gorm:"type:text"`
    Price       float64        `json:"price" gorm:"type:decimal(10,2);not null;default:0" validate:"required,gte=0"`
    Stock       int            `json:"stock" gorm:"type:integer;not null;default:0" validate:"required,gte=0"`
    IsActive    bool           `json:"is_active" gorm:"type:boolean;not null;default:false"`
    CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (p *Product) Validate() error {
    if p.Name == "" {
        return ErrInvalidProductName
    }
    if p.Price < 0 {
        return ErrInvalidProductPrice
    }
    if p.Stock < 0 {
        return ErrInvalidProductStock
    }
    return nil
}

func (p *Product) SoftDelete() {
    p.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
}

func (p *Product) IsDeleted() bool {
    return p.DeletedAt.Valid
}
```

You can extend this with additional business methods:

```go
func (p *Product) IsAvailable() bool {
    return p.IsActive && p.Stock > 0 && !p.IsDeleted()
}

func (p *Product) Restock(quantity int) error {
    if quantity <= 0 {
        return errors.New("restock quantity must be positive")
    }
    p.Stock += quantity
    p.UpdatedAt = time.Now()
    return nil
}

func (p *Product) Sell(quantity int) error {
    if quantity <= 0 {
        return errors.New("sell quantity must be positive")
    }
    if p.Stock < quantity {
        return errors.New("insufficient stock")
    }
    p.Stock -= quantity
    p.UpdatedAt = time.Now()
    return nil
}

func (p *Product) ApplyDiscount(percentage float64) error {
    if percentage < 0 || percentage > 100 {
        return errors.New("discount percentage must be between 0 and 100")
    }
    p.Price = p.Price * (1 - percentage/100)
    p.UpdatedAt = time.Now()
    return nil
}
```

These methods capture business logic that belongs in the domain layer. They make the entity more than a data structure; they make it a behavior-rich business object.

## Testing Domain Entities

Domain entities are easy to test because they have no external dependencies. You can test business logic in isolation:

```go
package domain_test

import (
    "testing"
    
    "yourproject/internal/domain"
)

func TestUser_Validate(t *testing.T) {
    tests := []struct {
        name    string
        user    domain.User
        wantErr bool
    }{
        {
            name: "valid user",
            user: domain.User{
                Name:  "John Doe",
                Email: "john@example.com",
                Age:   30,
            },
            wantErr: false,
        },
        {
            name: "empty name",
            user: domain.User{
                Name:  "",
                Email: "john@example.com",
                Age:   30,
            },
            wantErr: true,
        },
        {
            name: "negative age",
            user: domain.User{
                Name:  "John Doe",
                Email: "john@example.com",
                Age:   -5,
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestProduct_Sell(t *testing.T) {
    product := domain.Product{
        Name:  "Test Product",
        Price: 100.0,
        Stock: 10,
    }
    
    // Valid sale
    err := product.Sell(5)
    if err != nil {
        t.Errorf("Sell(5) failed: %v", err)
    }
    if product.Stock != 5 {
        t.Errorf("Expected stock 5, got %d", product.Stock)
    }
    
    // Insufficient stock
    err = product.Sell(10)
    if err == nil {
        t.Error("Sell(10) should fail with insufficient stock")
    }
    
    // Negative quantity
    err = product.Sell(-1)
    if err == nil {
        t.Error("Sell(-1) should fail with negative quantity")
    }
}
```

These tests run instantly because they do not touch databases, networks, or files. They verify business logic in isolation.

## Best Practices for Domain Entities

Based on Clean Architecture and DDD principles, follow these best practices:

**Keep Entities Pure**: Do not import framework or infrastructure code. Entities should compile without external dependencies beyond the standard library.

**Encapsulate State**: Use methods to modify entity state. Avoid exposing fields directly if business rules govern their modification.

**Express Business Rules**: Write methods that answer business questions and enforce business constraints. Make implicit knowledge explicit.

**Use Value Objects**: For concepts without identity, create value objects. An email address, money amount, or date range should be a value object, not part of an entity.

**Avoid Anemic Domain Models**: Entities with only getters and setters are anemic. Add behavior. Rich domain models contain both data and behavior.

**Design for Invariants**: Entities should always be in a valid state. Constructor functions and validation methods enforce this.

**Use Domain Language**: Name entities, fields, and methods using terms from your business domain. Code should read like business documentation.

## Integration with Other Layers

Entities work with other Clean Architecture layers through well-defined interfaces.

### Use Case Layer

Use cases orchestrate entities to fulfill application requirements:

```go
package usecase

type CreateProductInput struct {
    Name        string  `json:"name" validate:"required"`
    Description string  `json:"description"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Stock       int     `json:"stock" validate:"required,gte=0"`
}

type productService struct {
    repo repository.ProductRepository
}

func (s *productService) Create(input CreateProductInput) (*domain.Product, error) {
    // Create entity
    product := &domain.Product{
        Name:        input.Name,
        Description: input.Description,
        Price:       input.Price,
        Stock:       input.Stock,
        IsActive:    true,
    }
    
    // Validate business rules
    if err := product.Validate(); err != nil {
        return nil, err
    }
    
    // Persist through repository
    if err := s.repo.Save(product); err != nil {
        return nil, err
    }
    
    return product, nil
}
```

The use case depends on the entity, not the other way around. This maintains the dependency rule.

### Repository Layer

Repositories provide persistence for entities through interfaces defined in the domain:

```go
package repository

type ProductRepository interface {
    Save(product *domain.Product) error
    FindByID(id uint) (*domain.Product, error)
    Update(product *domain.Product) error
    Delete(id uint) error
    FindAll() ([]domain.Product, error)
}
```

The repository interface lives in the domain package, but implementations live in the infrastructure layer:

```go
package repository

import (
    "yourproject/internal/domain"
    "gorm.io/gorm"
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

func (r *postgresProductRepository) FindByID(id uint) (*domain.Product, error) {
    var product domain.Product
    err := r.db.First(&product, id).Error
    if err != nil {
        return nil, err
    }
    return &product, nil
}
```

This separation allows you to swap database implementations without changing business logic.

## Advanced Entity Patterns

### Aggregate Roots

Entities can serve as aggregate roots, controlling access to related entities:

```go
type Order struct {
    ID         uint
    CustomerID uint
    Items      []OrderItem
    Total      float64
    Status     string
}

type OrderItem struct {
    ID        uint
    ProductID uint
    Quantity  int
    Price     float64
}

func (o *Order) AddItem(productID uint, quantity int, price float64) error {
    if quantity <= 0 {
        return errors.New("quantity must be positive")
    }
    
    item := OrderItem{
        ProductID: productID,
        Quantity:  quantity,
        Price:     price,
    }
    
    o.Items = append(o.Items, item)
    o.calculateTotal()
    return nil
}

func (o *Order) calculateTotal() {
    total := 0.0
    for _, item := range o.Items {
        total += float64(item.Quantity) * item.Price
    }
    o.Total = total
}
```

The `Order` aggregate root controls `OrderItem` access, maintaining consistency.

### Factories

Complex entity creation can use factory patterns:

```go
func NewUser(name, email string, age int) (*User, error) {
    user := &User{
        Name:      name,
        Email:     email,
        Age:       age,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

Factories ensure entities are always created in valid states.

## Generating Complete Features

While `goca entity` generates entities, `goca feature` generates complete features including entities, use cases, repositories, and handlers:

```bash
goca feature Product --fields "name:string,price:float64,stock:int"
```

This generates:

- Domain entity (`internal/domain/product.go`)
- Use case interface and implementation (`internal/usecase/product_service.go`)
- Repository interface (`internal/repository/interfaces.go`)
- Repository implementation (`internal/repository/postgres_product_repository.go`)
- HTTP handler (`internal/handler/http/product_handler.go`)
- DTOs (`internal/usecase/dto.go`)
- Error definitions (`internal/domain/errors.go`)
- Seed data (`internal/domain/product_seeds.go`)

All layers work together following Clean Architecture principles, with the entity at the core.

## Conclusion

Domain entities are the foundation of Clean Architecture. They represent your business concepts and enforce business rules without coupling to external systems. Goca generates production-ready entities following DDD principles, giving you a solid starting point for building maintainable applications.

Understanding entities correctly is essential for successful software architecture. They are not database models, not API responses, and not framework-dependent structures. They are pure business logic, testable in isolation, and independent of implementation details.

By using Goca's entity generation commands and following Clean Architecture principles, you create systems that are:

- Easy to test without external dependencies
- Simple to modify as business requirements change
- Portable across different frameworks and technologies
- Clear in expressing business intent
- Maintainable over long periods

Start with entities. Build your business logic correctly. Let outer layers adapt to your domain, not the other way around.

## Further Reading

- Clean Architecture documentation in the [guide section](/guide/clean-architecture)
- Complete command reference for [`goca entity`](/commands/entity)
- Full feature generation with [`goca feature`](/commands/feature)
- Domain-Driven Design principles and patterns
- Repository pattern implementation examples
