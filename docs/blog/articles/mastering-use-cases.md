---
layout: doc
title: Mastering Use Cases in Clean Architecture
titleTemplate: Articles | Goca Blog
description: A deep dive into use cases, application services, DTOs, and how they orchestrate business workflows while maintaining clean separation of concerns
tags:
  - Use Cases
  - Application Services
  - DTOs
  - Clean Architecture
  - Go
---

<script setup>
import Badge from '../../.vitepress/theme/components/Badge.vue'
</script>

# Mastering Use Cases in Clean Architecture

<div style="display: flex; gap: 0.5rem; margin-bottom: 1rem;">
<Badge type="info">Architecture</Badge>
<Badge type="tip">Application Layer</Badge>
</div>

Use cases represent application-specific business rules and orchestrate the flow of data between entities and external systems. Understanding use cases correctly is critical for building well-structured applications that adapt to changing requirements without compromising core business logic.

---

## What is a Use Case?

A use case is an application service that coordinates domain entities and infrastructure to fulfill a specific user or system goal. Use cases answer the question: "What can the application do?"

### Core Responsibilities

**Orchestration**: Use cases coordinate multiple domain entities, repositories, and external services to complete a workflow. They do not contain business rules; they apply them.

**Data Flow Control**: Use cases manage the flow of data between the UI layer and the domain layer, transforming external requests into domain operations and domain results into external responses.

**Application Logic**: Use cases implement application-specific rules that do not belong in the domain. These rules depend on the use case context, not on core business concepts.

**Transaction Management**: Use cases define transaction boundaries, ensuring that operations either complete fully or roll back entirely.

**Permission and Security**: Use cases enforce authorization rules, checking whether the requesting user can perform the operation.

## Use Case vs Controller vs Service

Many developers confuse use cases with controllers or generic services. This confusion leads to bloated classes and violated architectural boundaries.

### What a Use Case Is NOT

**Not a Controller**: Controllers are adapters that convert HTTP requests to use case calls. Controllers handle protocol concerns; use cases handle application logic.

**Not a Generic Service**: A use case serves a specific goal, not general utilities. Services like "EmailService" or "LoggerService" are infrastructure concerns, not use cases.

**Not a Transaction Script**: Use cases orchestrate domain entities. They do not implement business rules. Domain logic belongs in entities and value objects.

**Not a Facade**: Use cases are not simple pass-throughs to repositories. They add application-level coordination and workflow management.

### The Clear Distinction

```
HTTP Request
    ↓
Controller (Adapter - Outer Layer)
    ↓ Converts to DTO
Use Case (Application Layer)
    ↓ Orchestrates
Domain Entity (Domain Layer)
    ↓ Enforces rules
Repository Interface (Domain Layer)
    ↓ Implements
Repository Implementation (Infrastructure Layer)
    ↓ Persists
Database
```

Each layer has distinct responsibilities. Use cases sit between adapters and domain, orchestrating operations without implementing business rules or handling external protocols.

## The Application Layer

Use cases form the application layer in Clean Architecture, distinct from both the domain layer and the infrastructure layer.

### Application Layer Characteristics

**Depends on Domain**: Use cases depend on domain entities and interfaces. They call entity methods and use repository interfaces defined in the domain.

**Independent of Infrastructure**: Use cases do not import database drivers, HTTP libraries, or external service clients. They work with interfaces.

**Stateless by Design**: Use cases do not maintain state between calls. Each operation is independent.

**Transaction Boundaries**: Use cases define where transactions begin and end, ensuring data consistency.

### Why a Separate Layer?

The application layer exists because application logic and domain logic are different:

**Domain Logic**: "A user must have a valid email address" is domain logic. This rule exists regardless of how you access users.

**Application Logic**: "To create a user, check if the email exists, create the user, send a welcome email" is application logic. This workflow is specific to the user creation use case.

Separating these concerns allows you to:

- Change workflows without changing domain rules
- Test business rules without application context
- Reuse domain logic across different workflows
- Evolve the application independently of the domain

## Data Transfer Objects (DTOs)

DTOs are simple structures that carry data between layers without behavior. Use cases use DTOs to receive input and provide output.

### Why DTOs?

**Layer Separation**: DTOs prevent external layers from depending on domain entities directly. Changing an entity does not break API contracts.

**Validation Boundary**: DTOs define what data the use case needs and validate it before processing.

**Security**: DTOs control what data external systems can provide or receive, preventing over-posting and data exposure.

**Versioning**: DTOs allow multiple API versions to coexist by mapping different external structures to the same domain entities.

### Input DTOs

Input DTOs represent the data a use case needs to perform an operation:

```go
type CreateUserInput struct {
    Name  string `json:"name" validate:"required,min=2"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"required,gte=0"`
}
```

Input DTOs include validation tags that define constraints:

- **required**: Field must be present
- **min/max**: String length or numeric range
- **email**: Valid email format
- **gte/lte**: Greater than or equal / less than or equal

These validations are input validations, not business rules. They ensure the data is well-formed before processing.

### Output DTOs

Output DTOs represent the data a use case returns:

```go
type CreateUserOutput struct {
    User    domain.User `json:"user"`
    Message string      `json:"message"`
}
```

Output DTOs can include:

- Domain entities for complete information
- Specific fields for minimal responses
- Metadata like messages or status codes
- Related entities for composite responses

### Update DTOs with Optional Fields

Update operations use optional fields to support partial updates:

```go
type UpdateUserInput struct {
    Name  *string `json:"name,omitempty" validate:"omitempty,min=2"`
    Email *string `json:"email,omitempty" validate:"omitempty,email"`
    Age   *int    `json:"age,omitempty" validate:"omitempty,gte=0"`
}
```

Pointer fields distinguish between "not provided" (nil) and "explicitly set to zero value" (non-nil pointer to zero value). This allows clients to update only specific fields without affecting others.

### List DTOs

List operations return collections with metadata:

```go
type ListUserOutput struct {
    Users   []domain.User `json:"users"`
    Total   int           `json:"total"`
    Message string        `json:"message"`
}
```

List DTOs can include pagination information, filters applied, and total counts.

## Use Case Implementation Patterns

Use cases follow consistent patterns regardless of the specific operation they perform.

### Basic Structure

Every use case implementation includes:

1. Dependency injection of required repositories
2. Input validation
3. Domain entity coordination
4. Business rule enforcement
5. Persistence through repositories
6. Output transformation

### Create Operation

The create operation instantiates a new domain entity, validates it, and persists it:

```go
func (s *userService) Create(input CreateUserInput) (*CreateUserOutput, error) {
    // 1. Create domain entity from input
    user := domain.User{
        Name:  input.Name,
        Email: input.Email,
        Age:   input.Age,
    }
    
    // 2. Validate business rules
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    // 3. Persist through repository
    if err := s.repo.Save(&user); err != nil {
        return nil, err
    }
    
    // 4. Return success output
    return &CreateUserOutput{
        User:    user,
        Message: "User created successfully",
    }, nil
}
```

This pattern ensures that:

- Input is transformed to domain entities
- Business rules are enforced before persistence
- Repository handles storage concerns
- Output is well-defined and structured

### Read Operation

The read operation retrieves an entity by identifier:

```go
func (s *userService) GetByID(id uint) (*domain.User, error) {
    return s.repo.FindByID(int(id))
}
```

Read operations are simple because they delegate directly to repositories. Complexity arises when reads require:

- Authorization checks
- Data enrichment from multiple sources
- Transformation to specific output formats

### Update Operation

The update operation retrieves an entity, modifies it, validates it, and persists changes:

```go
func (s *userService) Update(id uint, input UpdateUserInput) (*domain.User, error) {
    // 1. Retrieve existing entity
    user, err := s.repo.FindByID(int(id))
    if err != nil {
        return nil, err
    }
    
    // 2. Apply changes from input (only provided fields)
    if input.Name != nil {
        user.Name = *input.Name
    }
    if input.Email != nil {
        user.Email = *input.Email
    }
    if input.Age != nil {
        user.Age = *input.Age
    }
    
    // 3. Validate updated entity
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    // 4. Persist changes
    if err := s.repo.Update(user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

The update pattern:

- Retrieves current state
- Applies only provided changes
- Validates the result
- Persists atomically

### Delete Operation

The delete operation removes an entity:

```go
func (s *userService) Delete(id uint) error {
    return s.repo.Delete(int(id))
}
```

Delete operations can be:

- **Hard Delete**: Permanently removes the record
- **Soft Delete**: Marks the record as deleted without removing it

Soft deletes are preferable for audit trails and data recovery.

### List Operation

The list operation retrieves collections with optional filtering:

```go
func (s *userService) List() (*ListUserOutput, error) {
    users, err := s.repo.FindAll()
    if err != nil {
        return nil, err
    }
    
    return &ListUserOutput{
        Users:   users,
        Total:   len(users),
        Message: "Users listed successfully",
    }, nil
}
```

List operations often include:

- Pagination parameters
- Sort order specifications
- Filter conditions
- Total count calculation

## How Goca Generates Use Cases

Goca provides the `goca usecase` command to generate complete application services with DTOs and interfaces.

### Basic Use Case Generation

```bash
goca usecase UserService --entity User
```

This generates three files:

**dto.go**: Input and output DTOs

```go
package usecase

import (
    "github.com/yourorg/yourproject/internal/domain"
)

type CreateUserInput struct {
    Name  string `json:"name" validate:"required,min=2"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"required,min=1"`
}

type CreateUserOutput struct {
    User    domain.User `json:"user"`
    Message string      `json:"message"`
}

type UpdateUserInput struct {
    Name  *string `json:"name,omitempty" validate:"omitempty,min=2"`
    Email *string `json:"email,omitempty" validate:"omitempty,email"`
    Age   *int    `json:"age,omitempty" validate:"omitempty,min=1"`
}

type ListUserOutput struct {
    Users   []domain.User `json:"users"`
    Total   int           `json:"total"`
    Message string        `json:"message"`
}
```

**user_service.go**: Service implementation

```go
package usecase

import (
    "github.com/yourorg/yourproject/internal/domain"
    "github.com/yourorg/yourproject/internal/repository"
)

type userService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserUseCase {
    return &userService{repo: repo}
}

func (u *userService) Create(input CreateUserInput) (*CreateUserOutput, error) {
    user := domain.User{
        Name:  input.Name,
        Email: input.Email,
        Age:   input.Age,
    }
    
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    if err := u.repo.Save(&user); err != nil {
        return nil, err
    }
    
    return &CreateUserOutput{
        User:    user,
        Message: "User created successfully",
    }, nil
}

func (u *userService) GetByID(id uint) (*domain.User, error) {
    return u.repo.FindByID(int(id))
}

func (u *userService) Update(id uint, input UpdateUserInput) (*domain.User, error) {
    user, err := u.repo.FindByID(int(id))
    if err != nil {
        return nil, err
    }
    
    if input.Name != nil {
        user.Name = *input.Name
    }
    if input.Email != nil {
        user.Email = *input.Email
    }
    if input.Age != nil {
        user.Age = *input.Age
    }
    
    if err := u.repo.Update(user); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (u *userService) Delete(id uint) error {
    return u.repo.Delete(int(id))
}

func (u *userService) List() (*ListUserOutput, error) {
    users, err := u.repo.FindAll()
    if err != nil {
        return nil, err
    }
    
    return &ListUserOutput{
        Users:   users,
        Total:   len(users),
        Message: "Users listed successfully",
    }, nil
}
```

**interfaces.go**: Repository interface definition

The service depends on a repository interface:

```go
type UserRepository interface {
    Save(user *domain.User) error
    FindByID(id int) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id int) error
    FindAll() ([]domain.User, error)
}
```

This interface lives in the repository package but is used by the use case. The use case depends on the abstraction, not the implementation.

### Selecting Operations

Control which CRUD operations to generate:

```bash
goca usecase ProductService --entity Product --operations "create,read,update"
```

This generates only create, read, and update methods, omitting delete and list.

Available operations:

- **create**: Instantiate and persist new entities
- **read** or **get**: Retrieve entities by ID
- **update**: Modify existing entities
- **delete**: Remove entities
- **list**: Retrieve collections

### DTO Validation

Enable validation tags on DTOs:

```bash
goca usecase OrderService --entity Order --dto-validation
```

With validation enabled, DTOs include comprehensive validation rules:

```go
type CreateOrderInput struct {
    CustomerID int     `json:"customer_id" validate:"required,gt=0"`
    Total      float64 `json:"total" validate:"required,gte=0"`
    Status     string  `json:"status" validate:"required,oneof=pending confirmed shipped delivered"`
}
```

Validation rules ensure:

- Required fields are present
- Numeric values are within acceptable ranges
- Strings match expected patterns or enumerations
- Email addresses are valid
- Custom validation logic is applied

## Advanced Use Case Patterns

Beyond basic CRUD, use cases handle complex workflows and business processes.

### Transactional Use Cases

Some operations require multiple steps within a single transaction:

```go
func (s *orderService) CreateOrder(input CreateOrderInput) (*CreateOrderOutput, error) {
    // Begin transaction (pseudo-code, actual implementation depends on repository)
    
    // 1. Validate customer exists
    customer, err := s.customerRepo.FindByID(input.CustomerID)
    if err != nil {
        return nil, errors.New("customer not found")
    }
    
    // 2. Check product availability
    for _, item := range input.Items {
        product, err := s.productRepo.FindByID(item.ProductID)
        if err != nil {
            return nil, err
        }
        
        if product.Stock < item.Quantity {
            return nil, errors.New("insufficient stock")
        }
    }
    
    // 3. Create order
    order := &domain.Order{
        CustomerID: input.CustomerID,
        Items:      mapOrderItems(input.Items),
        Total:      calculateTotal(input.Items),
        Status:     "pending",
    }
    
    if err := order.Validate(); err != nil {
        return nil, err
    }
    
    // 4. Persist order
    if err := s.orderRepo.Save(order); err != nil {
        return nil, err
    }
    
    // 5. Update product stock
    for _, item := range order.Items {
        product, _ := s.productRepo.FindByID(item.ProductID)
        product.Stock -= item.Quantity
        s.productRepo.Update(product)
    }
    
    // Commit transaction
    
    return &CreateOrderOutput{
        Order:   *order,
        Message: "Order created successfully",
    }, nil
}
```

This use case:

- Validates dependencies (customer exists)
- Checks business constraints (sufficient stock)
- Creates the main entity (order)
- Updates related entities (product stock)
- Ensures atomicity through transactions

### Async Use Cases

Some operations can execute asynchronously to improve response times:

```bash
goca usecase NotificationService --entity Notification --operations "create" --async
```

Asynchronous use cases return immediately while processing continues in the background:

```go
func (s *notificationService) SendNotification(input SendNotificationInput) (*SendNotificationOutput, error) {
    // Validate input immediately
    if err := input.Validate(); err != nil {
        return nil, err
    }
    
    // Create notification record
    notification := &domain.Notification{
        UserID:  input.UserID,
        Message: input.Message,
        Status:  "queued",
    }
    
    if err := s.repo.Save(notification); err != nil {
        return nil, err
    }
    
    // Queue for async processing
    s.queue.Enqueue(notification.ID)
    
    // Return immediately
    return &SendNotificationOutput{
        NotificationID: notification.ID,
        Status:         "queued",
        Message:        "Notification queued successfully",
    }, nil
}
```

The actual sending happens asynchronously:

```go
func (s *notificationService) ProcessQueue() {
    for {
        notificationID := s.queue.Dequeue()
        
        notification, err := s.repo.FindByID(notificationID)
        if err != nil {
            continue
        }
        
        // Send notification via external service
        err = s.emailService.Send(notification.UserID, notification.Message)
        
        if err != nil {
            notification.Status = "failed"
        } else {
            notification.Status = "sent"
        }
        
        s.repo.Update(notification)
    }
}
```

Asynchronous use cases are appropriate for:

- Email sending
- File processing
- Report generation
- Third-party API calls
- Long-running computations

### Composite Use Cases

Some operations aggregate data from multiple sources:

```go
func (s *dashboardService) GetUserDashboard(userID uint) (*DashboardOutput, error) {
    // Retrieve user
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, err
    }
    
    // Retrieve user orders
    orders, err := s.orderRepo.FindByUserID(userID)
    if err != nil {
        return nil, err
    }
    
    // Calculate statistics
    totalSpent := calculateTotalSpent(orders)
    averageOrderValue := totalSpent / float64(len(orders))
    
    // Retrieve recent activity
    activity, err := s.activityRepo.FindRecentByUserID(userID, 10)
    if err != nil {
        return nil, err
    }
    
    return &DashboardOutput{
        User:              user,
        TotalOrders:       len(orders),
        TotalSpent:        totalSpent,
        AverageOrderValue: averageOrderValue,
        RecentActivity:    activity,
    }, nil
}
```

Composite use cases coordinate multiple repositories to build aggregate views.

## Testing Use Cases

Use cases are highly testable because they depend on interfaces, not concrete implementations.

### Unit Testing with Mocks

Test use cases by mocking repository dependencies:

```go
func TestUserService_Create(t *testing.T) {
    // Create mock repository
    mockRepo := new(MockUserRepository)
    
    // Setup expectations
    mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil)
    
    // Create service with mock
    service := NewUserService(mockRepo)
    
    // Execute use case
    input := CreateUserInput{
        Name:  "John Doe",
        Email: "john@example.com",
        Age:   30,
    }
    
    output, err := service.Create(input)
    
    // Assert results
    assert.NoError(t, err)
    assert.NotNil(t, output)
    assert.Equal(t, "John Doe", output.User.Name)
    assert.Equal(t, "User created successfully", output.Message)
    
    // Verify mock was called
    mockRepo.AssertExpectations(t)
}
```

Mock repositories allow you to:

- Test use case logic independently
- Simulate repository errors
- Verify correct method calls
- Control return values

### Testing Validation

Test that use cases enforce validation correctly:

```go
func TestUserService_Create_InvalidEmail(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)
    
    input := CreateUserInput{
        Name:  "John Doe",
        Email: "invalid-email",
        Age:   30,
    }
    
    output, err := service.Create(input)
    
    assert.Error(t, err)
    assert.Nil(t, output)
    assert.Contains(t, err.Error(), "email")
    
    // Repository should not be called
    mockRepo.AssertNotCalled(t, "Save")
}
```

### Testing Error Handling

Test that use cases handle repository errors gracefully:

```go
func TestUserService_Create_RepositoryError(t *testing.T) {
    mockRepo := new(MockUserRepository)
    
    // Simulate repository error
    mockRepo.On("Save", mock.Anything).Return(errors.New("database connection failed"))
    
    service := NewUserService(mockRepo)
    
    input := CreateUserInput{
        Name:  "John Doe",
        Email: "john@example.com",
        Age:   30,
    }
    
    output, err := service.Create(input)
    
    assert.Error(t, err)
    assert.Nil(t, output)
    assert.Equal(t, "database connection failed", err.Error())
}
```

## Integration with Other Layers

Use cases coordinate between domain, infrastructure, and adapter layers.

### Handler to Use Case

Handlers convert external requests to use case calls:

```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 1. Parse HTTP request
    var input usecase.CreateUserInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // 2. Call use case
    output, err := h.usecase.Create(input)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 3. Return HTTP response
    w.Header().Set("Content-Type", "application/json")
    w.WriteStatus(http.StatusCreated)
    json.NewEncoder(w).Encode(output)
}
```

The handler:

- Handles HTTP concerns (parsing, status codes, headers)
- Delegates business logic to the use case
- Transforms use case output to HTTP response

### Use Case to Repository

Use cases call repository methods through interfaces:

```go
type userService struct {
    repo repository.UserRepository // Interface, not implementation
}

func (s *userService) Create(input CreateUserInput) (*CreateUserOutput, error) {
    user := domain.User{
        Name:  input.Name,
        Email: input.Email,
        Age:   input.Age,
    }
    
    // Use case calls repository interface
    if err := s.repo.Save(&user); err != nil {
        return nil, err
    }
    
    return &CreateUserOutput{User: user}, nil
}
```

The repository implementation lives in the infrastructure layer:

```go
type postgresUserRepository struct {
    db *gorm.DB
}

func (r *postgresUserRepository) Save(user *domain.User) error {
    return r.db.Create(user).Error
}
```

This separation allows:

- Swapping database implementations
- Testing use cases without databases
- Changing persistence strategies independently

### Dependency Injection

Use cases receive dependencies through constructors:

```go
func NewUserService(
    repo repository.UserRepository,
    emailService EmailService,
    logger Logger,
) UserUseCase {
    return &userService{
        repo:         repo,
        emailService: emailService,
        logger:       logger,
    }
}
```

Dependencies are injected at the composition root, typically in a DI container:

```go
// Composition root
userRepo := repository.NewPostgresUserRepository(db)
emailService := email.NewSMTPService(config)
logger := log.NewStdLogger()

userService := usecase.NewUserService(userRepo, emailService, logger)
userHandler := handler.NewUserHandler(userService)
```

## Best Practices for Use Cases

Follow these practices to maintain clean, maintainable use cases.

**Keep Use Cases Thin**: Use cases orchestrate; they do not implement business rules. Business logic belongs in domain entities.

**One Use Case, One Goal**: Each use case serves a specific goal. "Create user" is one use case. "Create user and send email" might be one or two, depending on cohesion.

**Use DTOs for All External Data**: Never pass domain entities directly to or from external layers. DTOs provide a stable contract.

**Validate at Boundaries**: Validate input at the use case boundary. Do not assume data is valid.

**Return Errors, Don't Panic**: Use cases return errors for exceptional conditions. They do not panic or crash.

**Keep Dependencies Minimal**: Use cases should depend only on repositories and essential services. Avoid excessive dependencies.

**Write Comprehensive Tests**: Test use cases thoroughly with mocked dependencies. Use cases are the easiest layer to test.

**Document Complex Workflows**: Use cases with multiple steps should be documented clearly, explaining the workflow and error handling.

## Common Mistakes to Avoid

**Business Logic in Use Cases**: Do not implement business rules in use cases. Use cases apply rules defined in entities.

**Direct Database Access**: Use cases should not import database drivers or execute SQL. They call repository methods.

**Mixing Concerns**: Use cases should not handle HTTP parsing, logging details, or UI concerns. They orchestrate business operations.

**Returning Domain Entities Directly**: Always use DTOs for external communication. Domain entities are internal structures.

**Ignoring Errors**: Handle repository errors appropriately. Log them, wrap them, or transform them, but do not ignore them.

**Tight Coupling**: Use cases depending on concrete implementations cannot be tested or swapped easily. Depend on interfaces.

## Generating Complete Features

While `goca usecase` generates use cases, `goca feature` generates complete features including entities, use cases, repositories, and handlers:

```bash
goca feature User --fields "name:string,email:string,age:int"
```

This creates:

- Domain entity with validation
- Use case with all CRUD operations
- Repository interface and implementation
- HTTP handler
- DTOs for all operations
- Dependency injection wiring

All layers work together following Clean Architecture principles, with use cases at the center coordinating workflows.

## Conclusion

Use cases are the application layer in Clean Architecture, orchestrating domain entities and infrastructure services to fulfill user goals. They coordinate workflows without implementing business rules, maintain clear boundaries through DTOs, and depend on abstractions rather than implementations.

Understanding use cases correctly is essential for building maintainable applications. They are not controllers, not generic services, and not transaction scripts. They are focused coordinators that apply domain rules in application-specific contexts.

Goca generates production-ready use cases with comprehensive DTOs, following established patterns and best practices. By using Goca's use case generation and understanding these principles, you create systems that are:

- Easy to test with mocked dependencies
- Simple to modify as requirements change
- Clear in expressing application workflows
- Maintainable over long periods
- Adaptable to new platforms and interfaces

Start with clear use case boundaries. Orchestrate domain logic. Let adapters handle external concerns. Build applications that last.

## Further Reading

- Application layer patterns in the [guide section](/guide/clean-architecture)
- Complete command reference for [`goca usecase`](/commands/usecase)
- Full feature generation with [`goca feature`](/commands/feature)
- Understanding domain entities in [our previous article](/blog/articles/understanding-domain-entities)
- Repository pattern implementation examples
- Dependency injection patterns and best practices
