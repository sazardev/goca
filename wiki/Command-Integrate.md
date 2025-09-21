# goca integrate Command

The `goca integrate` command is a specialized tool for **integrating existing features** that are not connected to the dependency injection container or `main.go`.

## 🎯 Purpose

Useful for:
- 📦 Projects migrated from previous Goca versions
- 🔧 Manually generated features that need integration
- 🛠️ Repairing incomplete or damaged integrations
- 🔄 Updating existing projects with new auto-integration

## 📋 Syntax

```bash
goca integrate [flags]
```

## 🚩 Available Flags

| Flag         | Type     | Required | Description                                       |
| ------------ | -------- | -------- | ------------------------------------------------- |
| `--all`      | `bool`   | ❌ No     | Automatically detect and integrate all features   |
| `--features` | `string` | ❌ No     | Specific features to integrate (`"User,Product"`) |

## 📖 Usage Examples

### Automatic Integration (Recommended)
```bash
# Automatically detect all features and integrate them
goca integrate --all
```

**Expected output:**
```
🔍 Detecting existing features...
📋 Detected features: User, Product, Order

🔄 Starting integration process...

1️⃣  Setting up DI container...
   📦 Creating DI container...
   ✅ User integrated in DI container
   ✅ Product integrated in DI container
   ✅ Order integrated in DI container

2️⃣  Updating main.go...
   📍 Updating main.go at: main.go
   🔧 Rewriting complete main.go...
   ✅ main.go created with 3 features

3️⃣  Verifying integration...
   ✅ DI container exists
   ✅ main.go integrated (main.go)
   ✅ User routes integrated
   ✅ Product routes integrated
   ✅ Order routes integrated

🎯 Perfect integration! Everything is ready.

🎉 Integration completed!
✅ All features are now:
   🔗 Connected in DI container
   🛣️  With routes registered in main.go
   ⚡ Ready to use
```

### Specific Integration
```bash
# Integrate only specific features
goca integrate --features "User,Product"
```

### Common Use Case: After Cloning a Project
```bash
# 1. Clone existing project
git clone https://github.com/user/my-goca-project.git
cd my-goca-project

# 2. Automatically integrate all features
goca integrate --all

# 3. Verify everything works
go mod tidy
go run main.go
```

## 🔍 Automatic Detection

The `goca integrate --all` command automatically detects features by looking for:

1. **Domain entities** in `internal/domain/*.go`
2. **HTTP handlers** in `internal/handler/http/*_handler.go`
3. **Use cases** in `internal/usecase/*.go`

### Files Ignored in Detection
- `errors.go`
- `validations.go`
- `common.go`
- `types.go`

## 🔧 Integration Process

### 1. DI Container
- ✅ Creates `internal/di/container.go` if it doesn't exist
- ✅ Adds fields for repositories, use cases and handlers
- ✅ Sets up setup methods and getters
- ✅ Detects and avoids duplicates

### 2. Main.go
- ✅ Detects `main.go` in multiple locations
- ✅ Adds necessary imports (`internal/di`)
- ✅ Sets up DI container
- ✅ Registers all HTTP routes
- ✅ Preserves existing configuration

### 3. Generated Routes
For each feature, the following routes are automatically registered:
```go
// User routes
userHandler := container.UserHandler()
router.HandleFunc("/api/v1/users", userHandler.CreateUser).Methods("POST")
router.HandleFunc("/api/v1/users/{id}", userHandler.GetUser).Methods("GET")
router.HandleFunc("/api/v1/users/{id}", userHandler.UpdateUser).Methods("PUT")
router.HandleFunc("/api/v1/users/{id}", userHandler.DeleteUser).Methods("DELETE")
router.HandleFunc("/api/v1/users", userHandler.ListUsers).Methods("GET")
```

## 🏗️ Expected Project Structure

For integration to work correctly, the project must have:

```
myproject/
├── go.mod
├── main.go                          # Will be updated/created
├── internal/
│   ├── domain/
│   │   ├── user.go                  # ← Detected as "User" feature
│   │   └── product.go               # ← Detected as "Product" feature
│   ├── usecase/
│   │   ├── user_usecase.go
│   │   └── product_usecase.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   └── product_repository.go
│   ├── handler/
│   │   └── http/
│   │       ├── user_handler.go      # ← Detected as "User" feature
│   │       └── product_handler.go   # ← Detected as "Product" feature
│   └── di/                          # ← Will be created/updated
│       └── container.go
└── pkg/
    ├── config/
    └── logger/
```

## ⚠️ Special Cases

### Main.go Not Found
If `main.go` is not found, the command will create a complete new one:
```
⚠️  main.go not found, creating new one...
✅ main.go created with 3 features
```

### Features Already Integrated
If some features are already integrated, they will be skipped:
```
✅ User is already in DI container
➕ Adding Product to DI container...
✅ Product integrated in DI container
```

### Integration Errors
If there are problems, manual instructions are shown:
```
⚠️  Could not update main.go: permission denied

📋 Manual integration instructions:
1. Add import in main.go:
   "myproject/internal/di"
2. Add in main(), after connecting the DB:
   container := di.NewContainer(db)
3. Add feature routes:
   userHandler := container.UserHandler()
   router.HandleFunc("/api/v1/users", userHandler.CreateUser).Methods("POST")
   ...
```

## 🔄 Complete Workflow

### Scenario: Migrate Existing Project

```bash
# 1. Check current structure
ls internal/domain/     # See existing features

# 2. Run automatic integration
goca integrate --all

# 3. Verify everything compiles
go mod tidy
go build

# 4. Test the server
go run main.go

# 5. Test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/users
```

### Scenario: Add Feature to Integrated Project

```bash
# 1. Generate new feature (auto-integrated)
goca feature Order --fields "user_id:int,total:float64,status:string"

# 2. Ready to go! No integrate needed
go run main.go
```

## 🤝 Comparison with `goca feature`

| Aspect               | `goca feature`     | `goca integrate`            |
| -------------------- | ------------------ | --------------------------- |
| **Purpose**          | Create new feature | Integrate existing features |
| **Generates code**   | ✅ Yes (all layers) | ❌ No (integration only)     |
| **Auto-integration** | ✅ Yes (automatic)  | ✅ Yes (repair/update)       |
| **Typical use**      | New development    | Migration/repair            |
| **Detection**        | Not applicable     | ✅ Yes (automatic)           |

## 💡 Tips and Best Practices

### ✅ Recommendations
- **Use `--all`** for automatic detection instead of specifying features manually
- **Run after cloning** existing Goca projects
- **Combine with `go mod tidy`** after integration
- **Backup `main.go`** before integration if you have custom code

### ⚠️ Precautions
- **Review `main.go`** after integration if you had special configurations
- **Verify routes** if you already had custom endpoints registered
- **Compile after** integration to verify everything works

### 🔄 Continuous Integration
```bash
# Script for CI/CD
#!/bin/bash
go mod download
goca integrate --all
go mod tidy
go test ./...
go build
```

## 🆘 Troubleshooting

### Problem: "No features found"
```bash
# Check structure
ls internal/domain/
ls internal/handler/http/

# If files exist but are not detected:
goca integrate --features "User,Product"  # Specify manually
```

### Problem: "main.go could not be updated"
```bash
# Check permissions
ls -la main.go

# If it's a Windows permission issue:
# Run terminal as administrator

# Alternative: manual integration
goca integrate --features "User" --dry-run  # See instructions
```

### Problem: "Duplicate routes"
```bash
# The command automatically detects existing routes and skips them
# If there are conflicts, review main.go manually
```

## 🔗 Related Commands

- [`goca feature`](Command-Feature.md) - Create new features (includes auto-integration)
- [`goca di`](Command-DI.md) - Generate DI container only
- [`goca init`](Command-Init.md) - Initialize new project

---

**Next**: [DI Command](Command-DI.md) | **Previous**: [Feature Command](Command-Feature.md) | **Index**: [Commands](README.md)
