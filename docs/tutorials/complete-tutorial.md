# Complete Tutorial

Step-by-step guide to building a complete API with Goca.

## What We'll Build

A **Task Management API** with:
- Users
- Projects and tasks
- Task assignments
- Due dates and priorities

## Prerequisites

- Go 1.21+
- PostgreSQL installed (or Docker, to run one in a container)
- Basic Go knowledge

## Step 1: Initialize Project

```bash
goca init task-manager --module github.com/yourusername/task-manager --database postgres
cd task-manager
go mod tidy
```

This creates the complete project structure with clean architecture: `cmd/server/main.go`, `internal/{domain,usecase,repository,handler}/`, `pkg/{config,logger}/`, `migrations/`, `.goca.yaml`, and an `.env` / `.env.example` pair with the database and server settings.

The generated HTTP layer uses [gorilla/mux](https://github.com/gorilla/mux), not gin/chi/echo — handlers have the signature `func (h *XHandler) Method(w http.ResponseWriter, r *http.Request)` and read path parameters with `mux.Vars(r)`.

## Step 2: Configure the Database

Goca does **not** generate a `config/config.yaml` file — configuration is environment-variable based, defined in `pkg/config/config.go` and populated from `.env`. Edit `.env`:

```bash
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=taskmanager
DB_SSL_MODE=disable
```

`pkg/config/config.go` reads these with `os.Getenv`, it does not parse the `.env` file itself, so export the variables into your shell before running the server:

```bash
set -a
source .env
set +a
```

If you don't have PostgreSQL installed, you can start one with Docker:

```bash
docker run -d --name task-manager-db \
  -e POSTGRES_PASSWORD=yourpassword \
  -e POSTGRES_DB=taskmanager \
  -p 5432:5432 postgres:16-alpine
```

## Step 3: Create User Feature

```bash
goca feature User --fields "name:string,email:string,password:string,role:string"
```

This generates and **wires up automatically** (domain, use case, repository, HTTP handler, routes registered in `cmd/server/main.go`, and DI container updated in `internal/di/container.go`):
- Domain entity (`internal/domain/user.go`)
- Repository (`internal/repository/postgres_user_repository.go`)
- Use case (`internal/usecase/user_service.go`, `internal/usecase/user_usecase.go`)
- HTTP handler (`internal/handler/http/user_handler.go`)
- Routes (`internal/handler/http/routes.go`, called from `cmd/server/main.go`)

Run `go mod tidy` any time a feature adds a new dependency.

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

At this point `go build ./...` should already succeed — every `goca feature` run updates `cmd/server/main.go` (route registration and the GORM auto-migration list) and `internal/di/container.go` for you.

> Note: field names with underscores (`owner_id`, `project_id`, `assigned_to`, `due_date`) become Go struct fields in `PascalCase` (`OwnerID`, `ProjectID`, ...), but their JSON tags on the generated DTOs (`internal/usecase/dto.go`) drop the underscore, e.g. `ownerid`, `projectid`, `assignedto`, `duedate`. Keep this in mind for the `curl` calls in Step 12 — the request bodies use the DTO's JSON tags, not the original `--fields` spelling.

## Step 6: Review the Generated Relationships

Goca generates independent entities; it does not add GORM association tags (`foreignKey`, preloaded structs) between them automatically. `internal/domain/project.go` and `internal/domain/task.go` currently look like this:

```go
// internal/domain/project.go
type Project struct {
    ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
    Name        string `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
    Description string `json:"description" gorm:"type:text" validate:"required"`
    OwnerID     int    `json:"owner_id" gorm:"type:integer;not null;default:0" validate:"required,gte=0"`
}
```

```go
// internal/domain/task.go
type Task struct {
    ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    Title       string    `json:"title" gorm:"type:varchar(255);not null" validate:"required"`
    Description string    `json:"description" gorm:"type:text" validate:"required"`
    ProjectID   int       `json:"project_id" gorm:"type:integer;not null;default:0" validate:"required,gte=0"`
    AssignedTo  int       `json:"assigned_to" gorm:"type:integer;not null;default:0" validate:"required,gte=0"`
    Priority    string    `json:"priority" gorm:"type:varchar(255)" validate:"required"`
    Status      string    `json:"status" gorm:"type:varchar(255)" validate:"required"`
    DueDate     time.Time `json:"due_date" gorm:"not null" validate:"required"`
}
```

The `ProjectID`/`OwnerID`/`AssignedTo` fields are plain `int`, matching the type you passed with `--fields`. That's enough for GORM's auto-migration to create the foreign-key columns; you can add `gorm:"foreignKey:..."` association fields yourself later if you want GORM to preload related rows, but it's optional and not required for the rest of this tutorial.

## Step 7: Run Migrations

The server runs GORM auto-migration on startup for every entity generated so far (see the `entities := []interface{}{...}` list near the bottom of `cmd/server/main.go`), so for local development you don't need a separate migration step — just start the server (Step 11) and it creates the tables.

Goca also generates SQL migration files under `migrations/` (`001_initial.up.sql` / `001_initial.down.sql`) if you prefer applying migrations explicitly with a tool like [golang-migrate](https://github.com/golang-migrate/migrate):

```bash
migrate -path migrations -database "postgresql://postgres:yourpassword@localhost:5432/taskmanager?sslmode=disable" up
```

## Step 8: Add Business Logic

Add `AssignTask` and `GetTasksByProject` to the `TaskUseCase` interface in `internal/usecase/task_usecase.go`:

```go
type TaskUseCase interface {
    CreateTask(input CreateTaskInput) (CreateTaskOutput, error)
    GetTask(id int) (*domain.Task, error)
    UpdateTask(id int, input UpdateTaskInput) error
    DeleteTask(id int) error
    ListTasks() (ListTaskOutput, error)
    AssignTask(taskID, userID int) error
    GetTasksByProject(projectID int) ([]domain.Task, error)
}
```

Implement them in `internal/usecase/task_service.go` (add `"fmt"` to the imports):

```go
func (t *taskService) AssignTask(taskID, userID int) error {
    task, err := t.repo.FindByID(taskID)
    if err != nil {
        return fmt.Errorf("task not found: %w", err)
    }

    task.AssignedTo = userID
    return t.repo.Update(task)
}

func (t *taskService) GetTasksByProject(projectID int) ([]domain.Task, error) {
    return t.repo.FindByProjectID(projectID)
}
```

Note that generated use case methods take plain `int` IDs and DTOs, not `context.Context` — the generated code doesn't thread a context through these layers, so keep the same signature style for consistency.

## Step 9: Upgrade the Generated Repository Methods

`goca feature` already generated `FindByProjectID` and `FindByAssignedTo` on `TaskRepository`, but as single-record lookups (`(*domain.Task, error)`) since it doesn't know they should return collections. Upgrade both to return a slice.

In `internal/repository/interfaces.go`:

```go
type TaskRepository interface {
    Save(task *domain.Task) error
    FindByID(id int) (*domain.Task, error)
    FindByTitle(title string) (*domain.Task, error)
    FindByDescription(description string) (*domain.Task, error)
    FindByProjectID(projectid int) ([]domain.Task, error)
    FindByAssignedTo(assignedto int) ([]domain.Task, error)
    FindByPriority(priority string) (*domain.Task, error)
    FindByStatus(status string) (*domain.Task, error)
    Update(task *domain.Task) error
    Delete(id int) error
    FindAll() ([]domain.Task, error)
}
```

In `internal/repository/postgres_task_repository.go`, replace the two matching methods (the repository struct is unexported — `postgresTaskRepository` — so keep the receiver as-is):

```go
// FindByProjectID returns every task that belongs to the given project.
func (p *postgresTaskRepository) FindByProjectID(projectid int) ([]domain.Task, error) {
    var tasks []domain.Task
    result := p.db.Where("project_id = ?", projectid).Find(&tasks)
    if result.Error != nil {
        return nil, result.Error
    }
    return tasks, nil
}

// FindByAssignedTo returns every task assigned to the given user.
func (p *postgresTaskRepository) FindByAssignedTo(assignedto int) ([]domain.Task, error) {
    var tasks []domain.Task
    result := p.db.Where("assigned_to = ?", assignedto).Find(&tasks)
    if result.Error != nil {
        return nil, result.Error
    }
    return tasks, nil
}
```

## Step 10: Add Custom Routes

Add an `AssignTask` handler and a `GetTasksByProject` handler to `internal/handler/http/task_handler.go`:

```go
type assignTaskInput struct {
    UserID int `json:"user_id"`
}

func (t *TaskHandler) AssignTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    var input assignTaskInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := t.usecase.AssignTask(id, input.UserID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (t *TaskHandler) GetTasksByProject(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid project ID", http.StatusBadRequest)
        return
    }

    tasks, err := t.usecase.GetTasksByProject(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}
```

Register both routes at the bottom of `SetupTaskRoutes` in `internal/handler/http/routes.go` (this file already exists — `goca feature` appends a `SetupXRoutes` function per entity and `cmd/server/main.go` already calls all of them, so you only need to add lines inside the existing function, not wire anything new into `main.go`):

```go
func SetupTaskRoutes(router *mux.Router, uc usecase.TaskUseCase) {
    handler := NewTaskHandler(uc)

    taskRouter := router.PathPrefix("/tasks").Subrouter()
    taskRouter.Use(corsMiddleware)
    taskRouter.Use(loggingMiddleware)

    taskRouter.HandleFunc("", handler.CreateTask).Methods("POST")
    taskRouter.HandleFunc("/{id}", handler.GetTask).Methods("GET")
    taskRouter.HandleFunc("/{id}", handler.UpdateTask).Methods("PUT")
    taskRouter.HandleFunc("/{id}", handler.DeleteTask).Methods("DELETE")
    taskRouter.HandleFunc("", handler.ListTasks).Methods("GET")
    taskRouter.HandleFunc("/{id}/assign", handler.AssignTask).Methods("POST")

    // Nested route: list the tasks that belong to a project
    router.HandleFunc("/projects/{id}/tasks", handler.GetTasksByProject).Methods("GET")
}
```

## Step 11: Run the Server

```bash
set -a
source .env
set +a
go run cmd/server/main.go
```

Output:
```
Starting application vdev (built: unknown)
Environment: development
Connecting to database at localhost:5432/taskmanager
Database connected successfully
Running GORM auto-migrations...
GORM auto-migrations completed successfully
Database schema is up to date
Server starting on port 8080
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

Recall from Step 5's note: the DTO's JSON tag for `owner_id` is `ownerid`.

```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Website Redesign",
    "description": "Modernize company website",
    "ownerid": 1
  }'
```

### Create a Task

Same rule applies: `project_id` -> `projectid`, `due_date` -> `duedate`.

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Design homepage mockup",
    "description": "Create Figma designs for new homepage",
    "projectid": 1,
    "priority": "high",
    "status": "pending",
    "duedate": "2027-02-01T00:00:00Z"
  }'
```

### Assign the Task

This hits the custom route added in Step 10, whose request body we defined ourselves (`user_id`, with the underscore):

```bash
curl -X POST http://localhost:8080/api/v1/tasks/1/assign \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1}'
```

### List Tasks for a Project

```bash
curl http://localhost:8080/api/v1/projects/1/tasks
```

All of the above were verified against a real running instance of the generated server backed by PostgreSQL.

## Step 13: Add Tests

`goca feature` already generates a validation test per entity (e.g. `internal/domain/task_test.go`), and `github.com/stretchr/testify` is already in `go.mod`. Add a use case test with a small in-memory fake repository — create `internal/usecase/task_service_test.go`:

```go
package usecase_test

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"

    "github.com/yourusername/task-manager/internal/domain"
    "github.com/yourusername/task-manager/internal/usecase"
)

// fakeTaskRepository is a minimal in-memory implementation of
// repository.TaskRepository used to test the use case layer in isolation.
type fakeTaskRepository struct {
    tasks  map[int]*domain.Task
    nextID uint
}

func newFakeTaskRepository() *fakeTaskRepository {
    return &fakeTaskRepository{tasks: make(map[int]*domain.Task)}
}

func (f *fakeTaskRepository) Save(task *domain.Task) error {
    f.nextID++
    task.ID = f.nextID
    f.tasks[int(task.ID)] = task
    return nil
}

func (f *fakeTaskRepository) FindByID(id int) (*domain.Task, error) {
    task, ok := f.tasks[id]
    if !ok {
        return nil, domain.ErrInvalidTaskData
    }
    return task, nil
}

func (f *fakeTaskRepository) FindByTitle(title string) (*domain.Task, error)   { return nil, nil }
func (f *fakeTaskRepository) FindByDescription(d string) (*domain.Task, error) { return nil, nil }
func (f *fakeTaskRepository) FindByPriority(p string) (*domain.Task, error)    { return nil, nil }
func (f *fakeTaskRepository) FindByStatus(s string) (*domain.Task, error)     { return nil, nil }

func (f *fakeTaskRepository) FindByProjectID(projectID int) ([]domain.Task, error) {
    var result []domain.Task
    for _, t := range f.tasks {
        if t.ProjectID == projectID {
            result = append(result, *t)
        }
    }
    return result, nil
}

func (f *fakeTaskRepository) FindByAssignedTo(userID int) ([]domain.Task, error) {
    var result []domain.Task
    for _, t := range f.tasks {
        if t.AssignedTo == userID {
            result = append(result, *t)
        }
    }
    return result, nil
}

func (f *fakeTaskRepository) Update(task *domain.Task) error {
    f.tasks[int(task.ID)] = task
    return nil
}

func (f *fakeTaskRepository) Delete(id int) error {
    delete(f.tasks, id)
    return nil
}

func (f *fakeTaskRepository) FindAll() ([]domain.Task, error) {
    var result []domain.Task
    for _, t := range f.tasks {
        result = append(result, *t)
    }
    return result, nil
}

func TestTaskService_CreateTask(t *testing.T) {
    repo := newFakeTaskRepository()
    service := usecase.NewTaskService(repo)

    input := usecase.CreateTaskInput{
        Title:       "Design homepage mockup",
        Description: "Create Figma designs for new homepage",
        ProjectID:   1,
        Priority:    "high",
        Status:      "pending",
        DueDate:     time.Now().Add(24 * time.Hour),
    }

    output, err := service.CreateTask(input)

    assert.NoError(t, err)
    assert.NotZero(t, output.ID)
    assert.Equal(t, input.Title, output.Title)
}

func TestTaskService_AssignTask(t *testing.T) {
    repo := newFakeTaskRepository()
    service := usecase.NewTaskService(repo)

    created, err := service.CreateTask(usecase.CreateTaskInput{
        Title:       "Design homepage mockup",
        Description: "Create Figma designs",
        ProjectID:   1,
        Priority:    "high",
        Status:      "pending",
        DueDate:     time.Now(),
    })
    assert.NoError(t, err)

    err = service.AssignTask(int(created.ID), 42)
    assert.NoError(t, err)

    task, err := service.GetTask(int(created.ID))
    assert.NoError(t, err)
    assert.Equal(t, 42, task.AssignedTo)
}
```

Run tests:
```bash
go test ./...
```

(If you'd rather use generated mocks instead of a hand-written fake, run `goca mocks Task` and use `--integration-tests` on `goca feature` for scaffolded integration test suites.)

## Next Steps

- [Adding Features](/tutorials/adding-features) - Extend your project
- [Best Practices](/guide/best-practices) - Code quality guidelines
- [Deployment Guide](https://github.com/sazardev/goca/wiki) - Production deployment

## Congratulations!

You've built a complete REST API with:
- Clean Architecture
- CRUD operations
- A custom cross-entity endpoint and business logic method
- Tests
- Database integration

Keep exploring and building!
