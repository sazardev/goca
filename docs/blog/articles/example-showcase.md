---
layout: doc
title: Example - Advanced Features Showcase
titleTemplate: Articles | Goca Blog
description: Demonstration of blog post capabilities including Mermaid diagrams, code blocks, and markdown features
---

<script setup>
import Badge from '../../.vitepress/theme/components/Badge.vue'
</script>

# Advanced Features Showcase

<div style="display: flex; gap: 0.5rem; margin-bottom: 1rem;">
<Badge type="tip">Example</Badge>
<Badge type="info">Tutorial</Badge>
</div>

This is an example article demonstrating the full capabilities of the Goca blog system, including Mermaid diagrams, syntax-highlighted code blocks, and advanced markdown features.

---

## Clean Architecture Flow

Here's how Goca implements Clean Architecture principles using a Mermaid diagram:

```mermaid
graph TD
    A[HTTP Request] --> B[Handler Layer]
    B --> C[Use Case Layer]
    C --> D[Repository Interface]
    D --> E[Repository Implementation]
    E --> F[Database]
    
    B -.->|DTO| C
    C -.->|Domain Entity| D
```

### Layer Dependencies

The dependency rule states that source code dependencies must point inward:

```mermaid
graph LR
    A[External Interfaces<br/>Frameworks & Drivers] --> B[Interface Adapters<br/>Controllers, Presenters]
    B --> C[Application Business Rules<br/>Use Cases]
    C --> D[Enterprise Business Rules<br/>Entities]
```

## Code Generation Workflow

### Entity Generation Process

```go
// Generate an entity with Goca
package main

import (
    "fmt"
    "time"
)

// User represents a user entity in the domain layer
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"type:varchar(255);not null"`
    Email     string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Validate performs business logic validation
func (u *User) Validate() error {
    if u.Name == "" {
        return fmt.Errorf("name cannot be empty")
    }
    if u.Email == "" {
        return fmt.Errorf("email cannot be empty")
    }
    return nil
}
```

### Command Execution Sequence

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Validator
    participant Generator
    participant FileSystem
    
    User->>CLI: goca feature User --fields "name:string,email:string"
    CLI->>Validator: Validate input
    Validator-->>CLI: Validation OK
    CLI->>Generator: Generate files
    Generator->>FileSystem: Create domain/user.go
    Generator->>FileSystem: Create usecase/user_service.go
    Generator->>FileSystem: Create repository/postgres_user_repository.go
    Generator->>FileSystem: Create handler/http/user_handler.go
    FileSystem-->>User: ✅ Feature generated successfully
```

## Database Support Matrix

Goca supports multiple databases with specific implementations:

| Database        | Type         | Primary Use Case | Status   |
| --------------- | ------------ | ---------------- | -------- |
| PostgreSQL      | SQL          | OLTP/General     | ✅ Stable |
| PostgreSQL JSON | SQL+Document | Semi-structured  | ✅ Stable |
| MySQL           | SQL          | Web Applications | ✅ Stable |
| MongoDB         | NoSQL        | Document Store   | ✅ Stable |
| SQLite          | SQL          | Embedded/Testing | ✅ Stable |
| SQL Server      | SQL          | Enterprise       | ✅ Stable |
| Elasticsearch   | Search       | Full-text Search | ✅ Stable |
| DynamoDB        | NoSQL        | Serverless AWS   | ✅ Stable |

### Database Selection Decision Tree

```mermaid
graph TD
    A[Choose Database] --> B{Data Structure?}
    B -->|Relational| C{Scale?}
    B -->|Document| D{Managed?}
    B -->|Search| E[Elasticsearch]
    
    C -->|Small| F[SQLite]
    C -->|Medium| G{Platform?}
    C -->|Large| H[PostgreSQL]
    
    G -->|Open Source| I[MySQL]
    G -->|Enterprise| J[SQL Server]
    
    D -->|Self-hosted| K[MongoDB]
    D -->|AWS| L[DynamoDB]
```

## Project Structure Visualization

```mermaid
graph TB
    subgraph "Clean Architecture Layers"
        A[internal/domain/<br/>Entities & Business Rules]
        B[internal/usecase/<br/>Application Services]
        C[internal/repository/<br/>Data Access]
        D[internal/handler/<br/>External Interfaces]
    end
    
    subgraph "Infrastructure"
        E[internal/di/<br/>Dependency Injection]
        F[cmd/<br/>Entry Points]
    end
    
    F --> E
    E --> D
    D --> B
    B --> A
    C --> A
```

## Testing Strategy

### Test Pyramid

```mermaid
graph TB
    subgraph "Test Pyramid"
        A[E2E Tests<br/>Integration Tests]
        B[Service Tests<br/>Use Case Tests]
        C[Unit Tests<br/>Domain Logic Tests]
    end
    
    A --> B --> C
```

### Test Coverage by Layer

```bash
# Run tests with coverage
go test ./internal/domain/... -cover
# PASS coverage: 95.2% of statements

go test ./internal/usecase/... -cover
# PASS coverage: 89.7% of statements

go test ./internal/repository/... -cover
# PASS coverage: 78.4% of statements

go test ./internal/handler/... -cover
# PASS coverage: 82.1% of statements
```

## Best Practices

### Repository Pattern Implementation

```go
// Repository interface (domain layer)
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uint) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uint) error
    FindAll(ctx context.Context) ([]*User, error)
}

// PostgreSQL implementation (infrastructure layer)
type postgresUserRepository struct {
    db *gorm.DB
}

func (r *postgresUserRepository) Save(ctx context.Context, user *User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id uint) (*User, error) {
    var user User
    if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }
    return &user, nil
}
```

### Use Case Pattern

```go
type UserUseCase interface {
    CreateUser(ctx context.Context, input CreateUserInput) (*UserOutput, error)
    GetUser(ctx context.Context, id uint) (*UserOutput, error)
    UpdateUser(ctx context.Context, input UpdateUserInput) (*UserOutput, error)
    DeleteUser(ctx context.Context, id uint) error
    ListUsers(ctx context.Context) ([]*UserOutput, error)
}

type userService struct {
    repo UserRepository
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (*UserOutput, error) {
    // 1. Validate input
    if err := input.Validate(); err != nil {
        return nil, err
    }
    
    // 2. Create domain entity
    user := &User{
        Name:  input.Name,
        Email: input.Email,
    }
    
    // 3. Validate entity
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    // 4. Save to repository
    if err := s.repo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // 5. Return output DTO
    return &UserOutput{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
    }, nil
}
```

## Command Reference Quick Guide

```bash
# Initialize new project
goca init myproject --database postgres

# Generate complete feature
goca feature User --fields "name:string,email:string,age:int"

# Generate with testing support
goca feature Product --fields "name:string,price:float64" \
    --integration-tests --mocks

# Generate only entity
goca entity Order --fields "total:float64,status:string"

# Generate repository
goca repository User --database postgres

# Generate handler
goca handler User --type http

# Wire everything together
goca integrate --all

# Generate dependency injection
goca di
```

## State Machine Example

Here's how you might model an order state machine:

```mermaid
stateDiagram-v2
    [*] --> Pending
    Pending --> Processing: Payment Confirmed
    Pending --> Cancelled: User Cancelled
    
    Processing --> Shipped: Items Dispatched
    Processing --> Cancelled: Out of Stock
    
    Shipped --> Delivered: Customer Received
    Shipped --> Returned: Customer Rejected
    
    Delivered --> [*]
    Returned --> Refunded
    Refunded --> [*]
    Cancelled --> [*]
    
    note right of Processing
        Inventory reserved
        Payment captured
    end note
    
    note right of Shipped
        Tracking number generated
        Notification sent
    end note
```

## Performance Considerations

### Database Query Optimization

```sql
-- Before: N+1 query problem
SELECT * FROM users WHERE id = 1;
SELECT * FROM orders WHERE user_id = 1;
SELECT * FROM orders WHERE user_id = 2;
SELECT * FROM orders WHERE user_id = 3;

-- After: Single query with join
SELECT u.*, o.*
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE u.id IN (1, 2, 3);
```

```go
// Use GORM preloading to avoid N+1
var users []User
db.Preload("Orders").Find(&users)
```

## Deployment Architecture

```mermaid
graph LR
    A[Load Balancer] --> B[App Server 1]
    A --> C[App Server 2]
    A --> D[App Server 3]
    
    B --> E[(Primary DB)]
    C --> E
    D --> E
    
    E --> F[(Replica DB)]
    E --> G[(Replica DB)]
    
    B --> H[Redis Cache]
    C --> H
    D --> H
```

## Conclusion

This example demonstrates the powerful capabilities available in Goca blog posts:

- Mermaid diagrams for architecture visualization
- Syntax-highlighted code blocks
- Tables and structured data
- State machines and flow diagrams
- Sequence diagrams
- Best practices and patterns

Use these features to create comprehensive, professional documentation and blog posts for your Goca projects.

---

<div style="text-align: center; margin-top: 3rem; padding-top: 2rem; border-top: 1px solid var(--vp-c-divider);">

[Edit on GitHub](https://github.com/sazardev/goca/edit/master/docs/blog/articles/example-showcase.md) • [Report Issue](https://github.com/sazardev/goca/issues)

</div>
