# goca entity

Generate pure domain entities following Domain-Driven Design (DDD) principles.

## Syntax

```bash
goca entity <EntityName> [flags]
```

## Description

The `goca entity` command generates domain entities that represent core business concepts. These entities are pure, containing only business logic without external dependencies.

::: tip Domain Layer
Entities are the heart of Clean Architecture - they contain enterprise-wide business rules and are completely independent of frameworks, databases, or external systems.
:::

## Arguments

### `<EntityName>`

**Required.** The name of your entity (singular, PascalCase).

```bash
goca entity User
goca entity Product
goca entity Order
```

## Flags

### `--fields`

**Required.** Define the structure of your entity.

**Format:** `"field1:type1,field2:type2,..."`

**Supported Types:**
- `string` - Text data
- `int`, `int64` - Integer numbers
- `float64` - Decimal numbers
- `bool` - Boolean values
- `time.Time` - Timestamps
- `[]type` - Arrays/slices

```bash
goca entity Product --fields "name:string,price:float64,stock:int"
```

### `--validation`

Include domain-level validation methods.

```bash
goca entity User --fields "name:string,email:string,age:int" --validation
```

Generates:
- `Validate()` method with business rules
- Domain-specific error constants
- Input sanitization

### `--business-rules`

Generate business rule methods.

```bash
goca entity Order --fields "total:float64,status:string" --business-rules
```

Generates methods like:
- `CanBeCancelled() bool`
- `IsCompleted() bool`
- Business logic calculations

### `--timestamps`

Add automatic timestamp fields.

```bash
goca entity Product --fields "name:string,price:float64" --timestamps
```

Adds:
- `CreatedAt time.Time`
- `UpdatedAt time.Time`

### `--soft-delete`

Enable soft delete functionality.

```bash
goca entity User --fields "name:string,email:string" --soft-delete
```

Adds:
- `DeletedAt *time.Time`
- `IsDeleted() bool` method

### `--tests`

Generate unit tests for the entity (enabled by default).

```bash
goca entity User --fields "name:string,email:string,age:int" --validation --tests
```

**Generates:** `internal/domain/user_test.go`

Creates comprehensive test suite including:
- **Validation tests**: Table-driven tests for `Validate()` method
- **Initialization tests**: Verify field assignment
- **Edge case tests**: Test boundary conditions for each field
- **Type-specific tests**: String length, numeric ranges, email format, etc.

**Disable tests:**
```bash
goca entity User --fields "name:string,email:string" --tests=false
```

::: tip Testing Best Practices
Generated tests use [testify/assert](https://github.com/stretchr/testify) for readable assertions and follow table-driven test patterns recommended by the Go community.
:::

## Examples

### Basic Entity

```bash
goca entity User --fields "name:string,email:string,age:int"
```

**Generates:** `internal/domain/user.go`

```go
package domain

import "time"

type User struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Age       int       `json:"age"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(name, email string, age int) *User {
    return &User{
        Name:  name,
        Email: email,
        Age:   age,
    }
}
```

### Entity with Validation

```bash
goca entity Product \
  --fields "name:string,price:float64,stock:int" \
  --validation
```

**Generates:**

```go
package domain

import (
    "errors"
    "strings"
)

type Product struct {
    ID    uint    `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
    Stock int     `json:"stock"`
}

func (p *Product) Validate() error {
    if strings.TrimSpace(p.Name) == "" {
        return ErrProductNameRequired
    }
    
    if p.Price < 0 {
        return ErrInvalidPrice
    }
    
    if p.Stock < 0 {
        return ErrInvalidStock
    }
    
    return nil
}

var (
    ErrProductNameRequired = errors.New("product name is required")
    ErrInvalidPrice        = errors.New("price cannot be negative")
    ErrInvalidStock        = errors.New("stock cannot be negative")
)
```

### Complete Entity

```bash
goca entity Order \
  --fields "customerID:int,total:float64,status:string" \
  --validation \
  --business-rules \
  --timestamps \
  --soft-delete
```

**Generates:**

```go
package domain

import (
    "errors"
    "time"
)

type Order struct {
    ID         uint       `json:"id"`
    CustomerID int        `json:"customer_id"`
    Total      float64    `json:"total"`
    Status     string     `json:"status"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
    DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

func (o *Order) Validate() error {
    if o.CustomerID <= 0 {
        return ErrInvalidCustomerID
    }
    
    if o.Total < 0 {
        return ErrInvalidTotal
    }
    
    if !o.IsValidStatus() {
        return ErrInvalidStatus
    }
    
    return nil
}

// Business Rules
func (o *Order) IsValidStatus() bool {
    validStatuses := []string{"pending", "processing", "completed", "cancelled"}
    for _, status := range validStatuses {
        if o.Status == status {
            return true
        }
    }
    return false
}

func (o *Order) CanBeCancelled() bool {
    return o.Status == "pending" || o.Status == "processing"
}

func (o *Order) IsCompleted() bool {
    return o.Status == "completed"
}

func (o *Order) IsDeleted() bool {
    return o.DeletedAt != nil
}

var (
    ErrInvalidCustomerID = errors.New("invalid customer ID")
    ErrInvalidTotal      = errors.New("total cannot be negative")
    ErrInvalidStatus     = errors.New("invalid order status")
)
```

### Entity with Unit Tests

```bash
goca entity User \
  --fields "name:string,email:string,age:int" \
  --validation \
  --tests
```

**Generates:** `internal/domain/user_test.go`

```go
package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUser_Validate tests the Validate method with various scenarios
func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg   string
	}{
		{
			name: "valid entity",
			user: User{
				Name: "John Doe",
				Email: "test@example.com",
				Age: 25,
			},
			wantErr: false,
		},
		{
			name: "invalid user - empty name",
			user: User{
				Name: "",
				Email: "test@example.com",
				Age: 25,
			},
			wantErr: true,
			errMsg: "name",
		},
		{
			name: "invalid user - negative age",
			user: User{
				Name: "John Doe",
				Email: "test@example.com",
				Age: -1,
			},
			wantErr: true,
			errMsg: "age",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUser_Initialization tests entity field initialization
func TestUser_Initialization(t *testing.T) {
	user := &User{
		Name: "John Doe",
		Email: "test@example.com",
		Age: 25,
	}

	assert.Equal(t, "John Doe", user.Name, "Name should be set correctly")
	assert.Equal(t, "test@example.com", user.Email, "Email should be set correctly")
	assert.Equal(t, 25, user.Age, "Age should be set correctly")
}

// Additional field-specific tests for edge cases...
```

**Run tests:**
```bash
cd internal/domain
go test -v -run TestUser
```

**Expected output:**
```
=== RUN   TestUser_Validate
=== RUN   TestUser_Validate/valid_entity
=== RUN   TestUser_Validate/invalid_user_-_empty_name
=== RUN   TestUser_Validate/invalid_user_-_negative_age
--- PASS: TestUser_Validate (0.00s)
=== RUN   TestUser_Initialization
--- PASS: TestUser_Initialization (0.00s)
PASS
ok      yourproject/internal/domain    0.013s
```

## Generated Structure

```
internal/
└── domain/
    ├── order.go           # Main entity
    ├── order_test.go      # Unit tests (if --tests)
    ├── order_seeds.go     # Seed data
    └── errors.go          # Domain errors (if --validation)
```

::: info Note
The `--tests` flag is enabled by default. All generated tests use table-driven patterns and testify/assert for clean, maintainable test code.
:::

## Best Practices

###  DO

- Keep entities pure (no external dependencies)
- Include business validations
- Add meaningful business rule methods
- Use value objects for complex types
- Document business logic

**Example:**
```go
//  Good: Pure domain logic
func (u *User) IsAdult() bool {
    return u.Age >= 18
}

func (u *User) CanPlaceOrder() bool {
    return u.IsActive && u.EmailVerified
}
```

###  DON'T

- Import database packages
- Import HTTP frameworks
- Include infrastructure logic
- Add persistence methods

**Example:**
```go
//  Bad: Infrastructure dependency
import "database/sql"

func (u *User) Save(db *sql.DB) error {
    // Wrong layer!
}

//  Bad: Framework dependency
import "github.com/gin-gonic/gin"

func (u *User) ToJSON(c *gin.Context) {
    // Wrong responsibility!
}
```

## Integration

After generating an entity, you typically:

1. **Generate Use Cases:**
   ```bash
   goca usecase OrderService --entity Order
   ```

2. **Generate Repository:**
   ```bash
   goca repository Order --database postgres
   ```

3. **Generate Handler:**
   ```bash
   goca handler Order --type http
   ```

Or use the complete feature command:

```bash
goca feature Order --fields "customerID:int,total:float64,status:string"
```

## Field Type Reference

| Type        | Description    | Example                 |
| ----------- | -------------- | ----------------------- |
| `string`    | Text data      | `"name:string"`         |
| `int`       | Integer        | `"age:int"`             |
| `int64`     | Large integer  | `"userID:int64"`        |
| `float64`   | Decimal number | `"price:float64"`       |
| `bool`      | Boolean        | `"isActive:bool"`       |
| `time.Time` | Timestamp      | `"birthDate:time.Time"` |
| `[]string`  | String array   | `"tags:[]string"`       |
| `[]int`     | Integer array  | `"scores:[]int"`        |

## Tips

### Naming Conventions

```bash
#  Good
goca entity User --fields "firstName:string,lastName:string"
goca entity Product --fields "productName:string,unitPrice:float64"

#  Avoid
goca entity user --fields "first_name:string"  # Use PascalCase
goca entity PRODUCT --fields "PRICE:float64"    # Too loud
```

### Complex Domains

For complex domains, start simple and add complexity:

```bash
# Step 1: Basic structure
goca entity Invoice --fields "number:string,amount:float64"

# Step 2: Add more fields manually
# Edit internal/domain/invoice.go to add:
# - items []InvoiceItem
# - taxes []Tax
# - metadata map[string]interface{}
```

## Troubleshooting

### Entity Already Exists

**Problem:** "entity already exists"

**Solution:**
```bash
# Remove existing entity
rm internal/domain/order.go
rm internal/domain/order_*.go

# Regenerate
goca entity Order --fields "..."
```

### Invalid Field Type

**Problem:** "unsupported field type"

**Solution:** Use basic types first, then manually add custom types:
```go
// After generation, manually add:
type User struct {
    // ... generated fields ...
    Address Address // Custom type
    Tags    []Tag   // Custom slice
}
```

## See Also

- [`goca feature`](/commands/feature) - Generate complete feature
- [`goca usecase`](/commands/usecase) - Generate use cases
- [`goca repository`](/commands/repository) - Generate repositories
- [Clean Architecture Guide](/guide/clean-architecture) - Architecture principles
- [Domain Layer](/guide/clean-architecture#layer-1-domain-entities) - Layer details

## Next Steps

After creating your entity:

1. Add business logic methods manually if needed
2. Generate use cases that use this entity
3. Create repositories for persistence
4. Build handlers for external access

Or use the shortcut:

```bash
goca feature YourEntity --fields "..." # Generates everything!
```
