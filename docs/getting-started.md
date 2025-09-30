# Getting Started

Get started with GOCA CLI in minutes and generate your first Clean Architecture project.

## Prerequisites

- **Go 1.21+** installed on your system
- Basic understanding of Go programming
- Familiarity with Clean Architecture concepts (helpful but not required)

## Installation

### Using Go Install (Recommended)

```bash
go install github.com/sazardev/goca@v2.0.0
```

### Verify Installation

```bash
goca version
# Output: GOCA CLI v2.0.0
```

## Your First Project

### 1. Initialize a New Project

```bash
goca init myapp --database postgres --handlers http
```

This creates a complete project structure:

```
myapp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   ├── usecase/
│   ├── repository/
│   └── handler/
├── pkg/
│   ├── config/
│   └── database/
├── migrations/
├── .goca.yaml
└── Makefile
```

### 2. Navigate to Your Project

```bash
cd myapp
```

### 3. Generate Your First Feature

```bash
goca feature user \
  --fields "name:string,email:string,age:int,active:bool" \
  --validation \
  --business-rules
```

This generates:

- ✅ **Domain Entity** (`internal/domain/user.go`)
  - Struct definition with all fields
  - Validation methods
  - Business rule methods

- ✅ **Use Case** (`internal/usecase/user_usecase.go`)
  - CRUD operations
  - Business logic
  - Error handling

- ✅ **Repository** (`internal/repository/user_repository.go`)
  - Database operations
  - GORM integration
  - Query methods

- ✅ **HTTP Handler** (`internal/handler/http/user_handler.go`)
  - REST endpoints
  - Request/response handling
  - JSON serialization

- ✅ **Dependency Injection** (auto-integrated)
- ✅ **Database Migration** (`migrations/`)

### 4. Build and Run

```bash
# Install dependencies
go mod tidy

# Run migrations
make migrate-up

# Start the server
make run
```

Your API is now running at `http://localhost:8080`!

## Test Your API

### Create a User

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "active": true
  }'
```

### Get All Users

```bash
curl http://localhost:8080/api/users
```

### Get User by ID

```bash
curl http://localhost:8080/api/users/1
```

### Update User

```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated",
    "email": "john.updated@example.com",
    "age": 31,
    "active": true
  }'
```

### Delete User

```bash
curl -X DELETE http://localhost:8080/api/users/1
```

## Using YAML Configuration

Create a `.goca.yaml` file for project-wide configuration:

```yaml
project:
  name: "myapp"
  module: "github.com/myorg/myapp"

database:
  type: "postgres"
  host: "localhost"
  port: 5432
  name: "myapp_db"

handlers:
  - "http"
  - "grpc"

features:
  validation: true
  business_rules: true
  soft_delete: true
  timestamps: true
```

Now you can generate features without flags:

```bash
goca feature product --fields "name:string,price:float64,stock:int"
```

GOCA reads settings from `.goca.yaml` automatically!

## Next Steps

### Learn More

- 📚 [Command Reference](/commands/) - All available commands
- ⚙️ [Configuration Guide](/configuration/yaml) - YAML configuration
- 🏗️ [Architecture](/guide/architecture) - Clean Architecture concepts

### Advanced Features

- 🔧 [Custom Validation](/guide/validation) - Add custom validators
- 📊 [Business Rules](/guide/business-rules) - Implement business logic
- 🗄️ [Multiple Databases](/guide/databases) - PostgreSQL, MySQL, SQLite
- 🌐 [gRPC Support](/guide/grpc) - Generate gRPC services

### Examples

- 📦 [E-commerce API](/examples/ecommerce) - Complete example
- 🔐 [Auth System](/examples/auth) - Authentication patterns
- 📝 [Blog API](/examples/blog) - Content management

## Common Questions

### How do I add a new field to an entity?

Regenerate the feature with updated fields:

```bash
goca feature user \
  --fields "name:string,email:string,age:int,active:bool,bio:string"
```

### Can I customize generated code?

Yes! All generated code is yours to modify. GOCA generates starting point code.

### What databases are supported?

- PostgreSQL (recommended)
- MySQL
- SQLite

### Can I use with existing projects?

Yes! GOCA can be used in existing projects. Just initialize in a new directory and copy generated code.

## Need Help?

- 📖 [Full Documentation](/)
- 💬 [GitHub Discussions](https://github.com/sazardev/goca/discussions)
- 🐛 [Report Issues](https://github.com/sazardev/goca/issues)
- 📧 [Contact Support](mailto:support@goca.dev)

---

Ready to build amazing Go applications? Let's go! 🚀
