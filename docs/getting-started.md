# Getting Started

This guide will help you create your first project with Goca in **less than 5 minutes**. By the end, you'll have a functional REST API following Clean Architecture principles.

## What We'll Build

In this guide we'll create:

- A project with complete Clean Architecture structure
- A `User` entity with domain validations
- Full CRUD REST API with all layers
- Dependency injection configured
- Repository pattern with PostgreSQL

::: tip Estimated Time
**5 minutes** from zero to running API
:::

## Prerequisites

Before starting, make sure you have:

- **Go 1.21+** - [Download here](https://golang.org/dl/)
- **Goca installed** - See [Installation Guide](/guide/installation)
- **PostgreSQL** (optional for this tutorial)

## Step 1: Create the Project

```bash
# Create and navigate to your project directory
mkdir my-first-api
cd my-first-api

# Initialize with Goca
goca init my-api --module github.com/yourusername/my-api --database postgres

# Navigate to generated directory
cd my-api
```

::: details What just happened?
Goca created a complete project structure with:
- `internal/` - All your business logic layers
- `cmd/` - Application entry points
- `pkg/` - Shared packages
- Configuration files for database, HTTP server, and more
:::

## Step 2: Generate Your First Feature

```bash
# Generate complete User feature with all layers
goca feature User --fields "name:string,email:string,age:int"

# See what was generated
ls internal/domain/
ls internal/usecase/
ls internal/repository/
ls internal/handler/http/
```

::: tip What Gets Generated?
This single command creates:
- **Domain Entity**: `user.go` with validations
- **Use Cases**: Service interfaces and DTOs
- **Repository**: Interface and PostgreSQL implementation
- **HTTP Handler**: REST endpoints for CRUD
- **Dependency Injection**: Automatic wiring
- **Routes**: Automatically registered
:::

## Step 3: Install Dependencies

```bash
# Download and install Go dependencies
go mod tidy
```

## Step 4: Run Your API

```bash
# Start the server
go run cmd/server/main.go
```

You should see:
```
→ Server starting on :8080
→ Database connected
→ Routes registered
```

## Step 5: Test Your API

Now let's interact with our new API!

### Health Check

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2025-10-11T10:30:00Z"
}
```

### Create a User

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 28
  }'
```

**Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 28,
  "created_at": "2025-10-11T10:30:00Z"
}
```

### Get User by ID

```bash
curl http://localhost:8080/api/v1/users/1
```

### List All Users

```bash
curl http://localhost:8080/api/v1/users
```

**Response:**
```json
{
  "users": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "age": 28
    }
  ],
  "total": 1
}
```

### Update User

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "john.smith@example.com",
    "age": 29
  }'
```

### Delete User

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Understanding the Architecture

Let's see how your code is organized:

```
my-api/
├── internal/
│   ├── domain/           # 🟡 Business entities
│   │   └── user.go
│   ├── usecase/          # 🔴 Application logic
│   │   ├── dto.go
│   │   └── user_service.go
│   ├── repository/       # 🔵 Data persistence
│   │   └── postgres_user_repository.go
│   └── handler/          # 🟢 Input adapters
│       └── http/
│           └── user_handler.go
└── cmd/
    └── server/
        └── main.go       # Entry point
```

### The Clean Architecture Layers

1. **🟡 Domain** - Pure business logic, no dependencies
2. **🔴 Use Cases** - Application rules and workflows
3. **🔵 Repository** - Data access abstraction
4. **🟢 Handlers** - External interface adapters

::: info Dependency Rule
Dependencies always point inward:
```
Handler → UseCase → Repository → Domain
```
The domain never knows about outer layers!
:::

## Next Steps

Congratulations! You've created your first Clean Architecture API with Goca.

Here's what you can do next:

- [Learn Clean Architecture Concepts](/guide/clean-architecture)
- [Add More Features](/tutorials/adding-features)
- [Complete Tutorial](/tutorials/complete-tutorial)
- [Explore All Commands](/commands/)

## Common Issues

### Port Already in Use

If you see "address already in use":

```bash
# Find and kill the process using port 8080
lsof -ti:8080 | xargs kill -9
```

### Database Connection Failed

If using PostgreSQL and connection fails:

1. Make sure PostgreSQL is running
2. Update connection string in `.env` or config file
3. Create the database: `createdb my-api`

### Module Not Found

```bash
# Re-initialize Go modules
go mod init github.com/yourusername/my-api
go mod tidy
```

## Need Help?

- [GitHub Issues](https://github.com/sazardev/goca/issues)
- [Discussions](https://github.com/sazardev/goca/discussions)
- [Full Documentation](/guide/introduction)
