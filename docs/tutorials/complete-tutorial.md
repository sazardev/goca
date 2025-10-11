# Complete Tutorial

Step-by-step guide to building a complete API with Goca.

## What We'll Build

A **Task Management API** with:
- User authentication
- Projects and tasks
- Task assignments
- Due dates and priorities

## Prerequisites

- Go 1.21+
- PostgreSQL installed
- Basic Go knowledge

## Step 1: Initialize Project

```bash
goca init task-manager --module github.com/yourusername/task-manager --database postgres
cd task-manager
```

This creates the complete project structure with clean architecture and PostgreSQL configuration.

## Step 2: Configure Database

Edit `config/config.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: yourpassword
  dbname: taskmanager
  sslmode: disable

server:
  port: 8080
  environment: development
```

## Step 3: Create User Feature

```bash
goca feature User --fields "name:string,email:string,password:string,role:string"
```

This generates:
- Domain entity (`internal/domain/user.go`)
- Repository (`internal/repository/postgres_user_repository.go`)
- Use case (`internal/usecase/user_service.go`)
- HTTP handler (`internal/handler/http/user_handler.go`)
- Routes (automatically registered)

## Step 4: Create Project Feature

```bash
goca feature Project --fields "name:string,description:string,owner_id:int"
```

A project belongs to a user (owner).

## Step 5: Create Task Feature

```bash
goca feature Task --fields "title:string,description:string,project_id:int,assigned_to:int,priority:string,status:string,due_date:time.Time"
```

Tasks belong to projects and can be assigned to users.

## Step 6: Add Relationships

Edit `internal/domain/project.go`:

```go
type Project struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"not null"`
    Description string    `gorm:"type:text"`
    OwnerID     uint      `gorm:"not null"`
    Owner       User      `gorm:"foreignKey:OwnerID"`
    Tasks       []Task    `gorm:"foreignKey:ProjectID"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

Edit `internal/domain/task.go`:

```go
type Task struct {
    ID          uint      `gorm:"primaryKey"`
    Title       string    `gorm:"not null"`
    Description string    `gorm:"type:text"`
    ProjectID   uint      `gorm:"not null"`
    Project     Project   `gorm:"foreignKey:ProjectID"`
    AssignedTo  uint
    Assignee    User      `gorm:"foreignKey:AssignedTo"`
    Priority    string    `gorm:"default:'medium'"`
    Status      string    `gorm:"default:'pending'"`
    DueDate     time.Time
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

## Step 7: Run Migrations

```bash
# Auto-migrate database
go run cmd/server/main.go migrate
```

Or use manual migrations:

```bash
# Create migration
migrate create -ext sql -dir migrations -seq create_tables

# Edit migration files, then run:
migrate -path migrations -database "postgresql://postgres:password@localhost:5432/taskmanager?sslmode=disable" up
```

## Step 8: Add Business Logic

Edit `internal/usecase/task_service.go`:

```go
func (s *taskService) AssignTask(ctx context.Context, taskID, userID uint) error {
    task, err := s.repo.FindByID(ctx, taskID)
    if err != nil {
        return fmt.Errorf("task not found: %w", err)
    }
    
    // Verify user exists
    // Add your validation logic here
    
    task.AssignedTo = userID
    return s.repo.Update(ctx, task)
}

func (s *taskService) GetTasksByProject(ctx context.Context, projectID uint) ([]*TaskResponse, error) {
    tasks, err := s.repo.FindByProjectID(ctx, projectID)
    if err != nil {
        return nil, err
    }
    
    var responses []*TaskResponse
    for _, task := range tasks {
        responses = append(responses, toTaskResponse(task))
    }
    return responses, nil
}
```

## Step 9: Add Custom Repository Methods

Edit `internal/repository/postgres_task_repository.go`:

```go
func (r *PostgresTaskRepository) FindByProjectID(ctx context.Context, projectID uint) ([]*domain.Task, error) {
    var tasks []*domain.Task
    err := r.db.WithContext(ctx).
        Where("project_id = ?", projectID).
        Preload("Assignee").
        Find(&tasks).Error
    return tasks, err
}

func (r *PostgresTaskRepository) FindByAssignee(ctx context.Context, userID uint) ([]*domain.Task, error) {
    var tasks []*domain.Task
    err := r.db.WithContext(ctx).
        Where("assigned_to = ?", userID).
        Preload("Project").
        Find(&tasks).Error
    return tasks, err
}
```

## Step 10: Add Custom Routes

Edit `internal/handler/http/routes.go`:

```go
func SetupRoutes(router *gin.Engine, container *di.Container) {
    api := router.Group("/api/v1")
    
    // Users
    users := api.Group("/users")
    {
        users.POST("", container.UserHandler.CreateUser)
        users.GET("/:id", container.UserHandler.GetUser)
        users.GET("", container.UserHandler.ListUsers)
    }
    
    // Projects
    projects := api.Group("/projects")
    {
        projects.POST("", container.ProjectHandler.CreateProject)
        projects.GET("/:id", container.ProjectHandler.GetProject)
        projects.GET("", container.ProjectHandler.ListProjects)
        projects.GET("/:id/tasks", container.TaskHandler.GetTasksByProject)
    }
    
    // Tasks
    tasks := api.Group("/tasks")
    {
        tasks.POST("", container.TaskHandler.CreateTask)
        tasks.GET("/:id", container.TaskHandler.GetTask)
        tasks.PUT("/:id", container.TaskHandler.UpdateTask)
        tasks.DELETE("/:id", container.TaskHandler.DeleteTask)
        tasks.POST("/:id/assign", container.TaskHandler.AssignTask)
    }
}
```

## Step 11: Run the Server

```bash
go run cmd/server/main.go
```

Output:
```
🚀 Server starting on :8080
✓ Database connected
✓ Routes registered
✓ Server running
```

## Step 12: Test the API

### Create a User

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepass123",
    "role": "developer"
  }'
```

### Create a Project

```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Website Redesign",
    "description": "Modernize company website",
    "owner_id": 1
  }'
```

### Create a Task

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Design homepage mockup",
    "description": "Create Figma designs for new homepage",
    "project_id": 1,
    "priority": "high",
    "due_date": "2025-02-01T00:00:00Z"
  }'
```

### Assign Task

```bash
curl -X POST http://localhost:8080/api/v1/tasks/1/assign \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1}'
```

## Step 13: Add Tests

Create `internal/usecase/task_service_test.go`:

```go
package usecase_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestTaskService_CreateTask(t *testing.T) {
    mockRepo := new(MockTaskRepository)
    service := NewTaskService(mockRepo)
    
    input := CreateTaskInput{
        Title:       "Test Task",
        Description: "Test Description",
        ProjectID:   1,
    }
    
    mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
    
    result, err := service.CreateTask(context.Background(), input)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertExpectations(t)
}
```

Run tests:
```bash
go test ./...
```

## Next Steps

- [Adding Features](/tutorials/adding-features) - Extend your project
- [Best Practices](/guide/best-practices) - Code quality guidelines
- [Deployment Guide](https://github.com/sazardev/goca/wiki) - Production deployment

## Congratulations! 🎉

You've built a complete REST API with:
- ✅ Clean Architecture
- ✅ CRUD operations
- ✅ Relationships
- ✅ Business logic
- ✅ Database integration

Keep exploring and building! 🚀
