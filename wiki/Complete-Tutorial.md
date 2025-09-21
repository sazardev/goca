# Complete Tutorial: E-commerce with Goca

In this tutorial we'll create a **complete e-commerce system** using Goca, from project initialization to having a functional API with multiple features.

## 🎯 Objective

By the end of this tutorial you'll have:

- ✅ User system with authentication
- ✅ Product catalog with categories  
- ✅ Order system with items
- ✅ Complete REST API
- ✅ PostgreSQL database
- ✅ Configured dependency injection
- ✅ Clean Architecture structure

## 📋 Prerequisites

- **Go 1.21+** installed
- **PostgreSQL** installed and running
- **Goca CLI** installed (`go install github.com/sazardev/goca@latest`)
- **curl** or **Postman** to test APIs

## 🚀 Step 1: Initialize the Project

### Create the base project
```bash
# Create project directory
mkdir ecommerce-api
cd ecommerce-api

# Initialize with Goca
goca init ecommerce-api \
  --module github.com/mycompany/ecommerce-api \
  --database postgres \
  --auth \
  --api rest

# Enter the generated directory
cd ecommerce-api

# Install dependencies
go mod tidy
```

### Verify the created structure
```bash
tree
```

**Expected output:**
```
ecommerce-api/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   ├── usecase/
│   ├── repository/
│   └── handler/
├── pkg/
│   ├── config/
│   ├── logger/
│   └── auth/
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## 🔧 Step 2: Configure Database

### Create PostgreSQL database
```sql
-- Connect to PostgreSQL
psql -U postgres

-- Create database
CREATE DATABASE ecommerce_db;

-- Create user (optional)
CREATE USER ecommerce_user WITH PASSWORD 'mypassword';
GRANT ALL PRIVILEGES ON DATABASE ecommerce_db TO ecommerce_user;

\q
```

### Configure environment variables
```bash
# Create .env file
cat > .env << EOF
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=ecommerce_user
DB_PASSWORD=mypassword
DB_NAME=ecommerce_db
DB_SSL_MODE=disable

# Server
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=debug

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_ISSUER=ecommerce-api
EOF
```

## 👤 Step 3: Create User Feature

### Generate complete user feature
```bash
goca feature User \
  --fields "name:string,email:string,password:string,role:string,phone:string" \
  --validation \
  --business-rules \
  --database postgres \
  --handlers http
```

### Verify generated files
```bash
find internal/ -name "*user*" -type f
```

**Created files:**
```
internal/domain/user.go
internal/domain/errors.go
internal/usecase/dto/user_dto.go
internal/usecase/user_usecase.go
internal/usecase/interfaces/user_interfaces.go
internal/repository/interfaces/user_repository.go
internal/repository/postgres/user_repository.go
internal/handler/http/user_handler.go
internal/handler/http/user_routes.go
internal/handler/http/dto/user_dto.go
internal/messages/errors.go
internal/messages/responses.go
```

## 🛍️ Step 4: Create Product Feature

### Generate product feature
```bash
goca feature Product \
  --fields "name:string,description:string,price:float64,category:string,stock:int,sku:string,image_url:string" \
  --validation \
  --business-rules \
  --database postgres \
  --handlers http
```

### Generate messages for products
```bash
goca messages Product --all
```

## 📦 Step 5: Create Order Feature

### Generate order feature
```bash
goca feature Order \
  --fields "user_id:int,total:float64,status:string,payment_method:string,shipping_address:string" \
  --validation \
  --business-rules \
  --database postgres \
  --handlers "http,worker"
```

### Generate order item feature
```bash
goca feature OrderItem \
  --fields "order_id:int,product_id:int,quantity:int,price:float64" \
  --validation \
  --database postgres \
  --handlers http
```

## 🔌 Step 6: Configure Dependency Injection

### Generate DI container
```bash
goca di \
  --features "User,Product,Order,OrderItem" \
  --database postgres
```

### Verify generated DI files
```bash
ls -la internal/infrastructure/di/
```

## 🗄️ Step 7: Create Database Tables

### SQL script to create tables
```bash
# Create migration file
cat > migrations/001_initial_schema.sql << 'EOF'
-- User table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'customer',
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Product table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    category VARCHAR(100) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    sku VARCHAR(100) UNIQUE NOT NULL,
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Order table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    total DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    payment_method VARCHAR(50),
    shipping_address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Order items table
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id),
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
EOF
```

### Execute migration
```bash
# Create migrations directory
mkdir -p migrations

# Apply migration
psql -h localhost -U ecommerce_user -d ecommerce_db -f migrations/001_initial_schema.sql
```

## 🔧 Step 8: Integrate Everything in Main

### Update cmd/server/main.go
```go
package main

import (
    "database/sql"
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    _ "github.com/lib/pq"
    
    "github.com/mycompany/ecommerce-api/internal/infrastructure/di"
    userHTTP "github.com/mycompany/ecommerce-api/internal/handler/http"
    "github.com/mycompany/ecommerce-api/pkg/config"
    "github.com/mycompany/ecommerce-api/pkg/logger"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Initialize logger
    logger.Init(cfg.LogLevel)
    
    // Connect to database
    db, err := sql.Open("postgres", buildDSN(cfg))
    if err != nil {
        log.Fatal("Error connecting to database:", err)
    }
    defer db.Close()
    
    // Verify connection
    if err := db.Ping(); err != nil {
        log.Fatal("Error pinging database:", err)
    }
    
    // Initialize DI container
    container := di.NewContainer(db)
    
    // Configure router
    router := gin.Default()
    
    // Middleware
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "ok",
            "service": "ecommerce-api",
        })
    })
    
    // API routes
    api := router.Group("/api/v1")
    
    // Register routes
    userHTTP.RegisterUserRoutes(api, container.GetUserUseCase())
    // productHTTP.RegisterProductRoutes(api, container.GetProductUseCase())
    // orderHTTP.RegisterOrderRoutes(api, container.GetOrderUseCase())
    
    // Start server
    log.Printf("🚀 Server starting on port %s", cfg.Port)
    log.Printf("📖 API Documentation: http://localhost:%s/api/v1", cfg.Port)
    
    if err := router.Run(":" + cfg.Port); err != nil {
        log.Fatal("Error starting server:", err)
    }
}

func buildDSN(cfg *config.Config) string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Name,
        cfg.Database.SSLMode,
    )
}
```

## 🏃‍♂️ Step 9: Run and Test

### Run the server
```bash
# Install missing dependencies
go mod tidy

# Run server
go run cmd/server/main.go
```

**Expected output:**
```
🚀 Server starting on port 8080
📖 API Documentation: http://localhost:8080/api/v1
[GIN-debug] Listening and serving HTTP on :8080
```

### Test Health Check
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "service": "ecommerce-api"
}
```

## 🧪 Step 10: Test APIs

### 1. Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "password": "password123",
    "role": "customer",
    "phone": "+1234567890"
  }'
```

**Expected response:**
```json
{
  "id": 1,
  "name": "Juan Pérez",
  "email": "juan@example.com",
  "role": "customer",
  "phone": "+1234567890",
  "created_at": "2025-07-20T10:30:00Z",
  "updated_at": "2025-07-20T10:30:00Z"
}
```

### 2. Get User by ID
```bash
curl http://localhost:8080/api/v1/users/1
```

### 3. List Users
```bash
curl "http://localhost:8080/api/v1/users?page=1&limit=10"
```

### 4. Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Carlos Pérez",
    "phone": "+1234567899"
  }'
```

### 5. Create Product
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro",
    "description": "Latest iPhone from Apple",
    "price": 999.99,
    "category": "smartphones",
    "stock": 50,
    "sku": "IPHONE15PRO",
    "image_url": "https://example.com/iphone15pro.jpg"
  }'
```

## 📈 Step 11: Advanced Features

### Add Authentication Middleware
```bash
# Add JWT middleware to protected routes
# (Specific code dependent on JWT implementation)
```

### Add Advanced Validations
```bash
# Validations are already included with --validation
# Review generated files for customizations
```

### Add Structured Logging
```bash
# Logging is already configured in pkg/logger/
# Customize according to needs
```

## 🐛 Troubleshooting

### Database Connection Error
```bash
# Verify PostgreSQL is running
pg_isready -h localhost -p 5432

# Verify credentials in .env
cat .env | grep DB_
```

### Module Not Found Error
```bash
# Clean module cache
go clean -modcache

# Reinstall dependencies
go mod tidy
```

### Port in Use Error
```bash
# Check what process is using port 8080
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Change port in .env
echo "PORT=8081" >> .env
```

## 🎉 Final Result

Congratulations! You now have:

### ✅ Complete System
- **REST API** functional at http://localhost:8080
- **4 complete features**: User, Product, Order, OrderItem
- **PostgreSQL database** configured
- **Clean Architecture** correctly implemented

### ✅ Professional Structure
```
ecommerce-api/
├── cmd/server/main.go                    # Entry point
├── internal/
│   ├── domain/                           # Entities and business rules
│   ├── usecase/                          # Application logic
│   ├── repository/                       # Data access
│   ├── handler/                          # HTTP adapters
│   ├── infrastructure/di/                # Dependency injection
│   └── messages/                         # Messages and constants
├── pkg/
│   ├── config/                           # Configuration
│   ├── logger/                           # Logging system
│   └── auth/                             # JWT authentication
└── migrations/                           # Database scripts
```

### ✅ Available APIs
- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - List users
- `GET /api/v1/users/:id` - Get user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- *(Similar for products, orders and order items)*

## 🚀 Next Steps

### Additional Features
1. **Complete JWT authentication**
2. **Authorization middleware** by roles
3. **Automatic Swagger documentation**
4. **Unit and integration tests**
5. **Docker** for deployment
6. **CI/CD** with GitHub Actions

### System Expansion
1. **Product categories system**
2. **Temporary shopping cart**
3. **Discount and coupon system**
4. **Email notifications**
5. **Admin dashboard**

---

**You have successfully completed your first complete project with Goca! 🎊**

**← [Getting Started](Getting-Started) | [REST API Example](Example-REST-API) →**
