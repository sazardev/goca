# Tutorial Completo: E-commerce con Goca

En este tutorial crearemos un **sistema de e-commerce completo** usando Goca, desde la inicialización del proyecto hasta tener una API funcional con múltiples features.

## 🎯 Objetivo

Al finalizar este tutorial tendrás:

- ✅ Sistema de usuarios con autenticación
- ✅ Catálogo de productos con categorías  
- ✅ Sistema de órdenes con items
- ✅ API REST completa
- ✅ Base de datos PostgreSQL
- ✅ Inyección de dependencias configurada
- ✅ Estructura Clean Architecture

## 📋 Prerrequisitos

- **Go 1.21+** instalado
- **PostgreSQL** instalado y ejecutándose
- **Goca CLI** instalado (`go install github.com/sazardev/goca@latest`)
- **curl** o **Postman** para probar APIs

## 🚀 Paso 1: Inicializar el Proyecto

### Crear el proyecto base
```bash
# Crear directorio del proyecto
mkdir ecommerce-api
cd ecommerce-api

# Inicializar con Goca
goca init ecommerce-api \
  --module github.com/miempresa/ecommerce-api \
  --database postgres \
  --auth \
  --api rest

# Entrar al directorio generado
cd ecommerce-api

# Instalar dependencias
go mod tidy
```

### Verificar la estructura creada
```bash
tree
```

**Salida esperada:**
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

## 🔧 Paso 2: Configurar Base de Datos

### Crear base de datos PostgreSQL
```sql
-- Conectar a PostgreSQL
psql -U postgres

-- Crear base de datos
CREATE DATABASE ecommerce_db;

-- Crear usuario (opcional)
CREATE USER ecommerce_user WITH PASSWORD 'mypassword';
GRANT ALL PRIVILEGES ON DATABASE ecommerce_db TO ecommerce_user;

\q
```

### Configurar variables de entorno
```bash
# Crear archivo .env
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

## 👤 Paso 3: Crear Feature de Usuarios

### Generar feature completo de usuarios
```bash
goca feature User \
  --fields "name:string,email:string,password:string,role:string,phone:string" \
  --validation \
  --business-rules \
  --database postgres \
  --handlers http
```

### Verificar archivos generados
```bash
find internal/ -name "*user*" -type f
```

**Archivos creados:**
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

## 🛍️ Paso 4: Crear Feature de Productos

### Generar feature de productos
```bash
goca feature Product \
  --fields "name:string,description:string,price:float64,category:string,stock:int,sku:string,image_url:string" \
  --validation \
  --business-rules \
  --database postgres \
  --handlers http
```

### Generar mensajes para productos
```bash
goca messages Product --all
```

## 📦 Paso 5: Crear Feature de Órdenes

### Generar feature de órdenes
```bash
goca feature Order \
  --fields "user_id:int,total:float64,status:string,payment_method:string,shipping_address:string" \
  --validation \
  --business-rules \
  --database postgres \
  --handlers "http,worker"
```

### Generar feature de items de orden
```bash
goca feature OrderItem \
  --fields "order_id:int,product_id:int,quantity:int,price:float64" \
  --validation \
  --database postgres \
  --handlers http
```

## 🔌 Paso 6: Configurar Inyección de Dependencias

### Generar contenedor DI
```bash
goca di \
  --features "User,Product,Order,OrderItem" \
  --database postgres
```

### Verificar archivos DI generados
```bash
ls -la internal/infrastructure/di/
```

## 🗄️ Paso 7: Crear Tablas de Base de Datos

### Script SQL para crear tablas
```bash
# Crear archivo de migración
cat > migrations/001_initial_schema.sql << 'EOF'
-- Tabla de usuarios
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

-- Tabla de productos
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

-- Tabla de órdenes
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

-- Tabla de items de orden
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id),
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para mejor performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
EOF
```

### Ejecutar migración
```bash
# Crear directorio de migraciones
mkdir -p migrations

# Aplicar migración
psql -h localhost -U ecommerce_user -d ecommerce_db -f migrations/001_initial_schema.sql
```

## 🔧 Paso 8: Integrar Todo en Main

### Actualizar cmd/server/main.go
```go
package main

import (
    "database/sql"
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    _ "github.com/lib/pq"
    
    "github.com/miempresa/ecommerce-api/internal/infrastructure/di"
    userHTTP "github.com/miempresa/ecommerce-api/internal/handler/http"
    "github.com/miempresa/ecommerce-api/pkg/config"
    "github.com/miempresa/ecommerce-api/pkg/logger"
)

func main() {
    // Cargar configuración
    cfg := config.Load()
    
    // Inicializar logger
    logger.Init(cfg.LogLevel)
    
    // Conectar a base de datos
    db, err := sql.Open("postgres", buildDSN(cfg))
    if err != nil {
        log.Fatal("Error connecting to database:", err)
    }
    defer db.Close()
    
    // Verificar conexión
    if err := db.Ping(); err != nil {
        log.Fatal("Error pinging database:", err)
    }
    
    // Inicializar contenedor DI
    container := di.NewContainer(db)
    
    // Configurar router
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
    
    // Registrar rutas
    userHTTP.RegisterUserRoutes(api, container.GetUserUseCase())
    // productHTTP.RegisterProductRoutes(api, container.GetProductUseCase())
    // orderHTTP.RegisterOrderRoutes(api, container.GetOrderUseCase())
    
    // Iniciar servidor
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

## 🏃‍♂️ Paso 9: Ejecutar y Probar

### Ejecutar el servidor
```bash
# Instalar dependencias faltantes
go mod tidy

# Ejecutar servidor
go run cmd/server/main.go
```

**Salida esperada:**
```
🚀 Server starting on port 8080
📖 API Documentation: http://localhost:8080/api/v1
[GIN-debug] Listening and serving HTTP on :8080
```

### Probar Health Check
```bash
curl http://localhost:8080/health
```

**Respuesta:**
```json
{
  "status": "ok",
  "service": "ecommerce-api"
}
```

## 🧪 Paso 10: Probar APIs

### 1. Crear Usuario
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

**Respuesta esperada:**
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

### 2. Obtener Usuario por ID
```bash
curl http://localhost:8080/api/v1/users/1
```

### 3. Listar Usuarios
```bash
curl "http://localhost:8080/api/v1/users?page=1&limit=10"
```

### 4. Actualizar Usuario
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Carlos Pérez",
    "phone": "+1234567899"
  }'
```

### 5. Crear Producto
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro",
    "description": "Último iPhone de Apple",
    "price": 999.99,
    "category": "smartphones",
    "stock": 50,
    "sku": "IPHONE15PRO",
    "image_url": "https://example.com/iphone15pro.jpg"
  }'
```

## 📈 Paso 11: Funcionalidades Avanzadas

### Agregar Middleware de Autenticación
```bash
# Agregar middleware JWT a las rutas protegidas
# (Código específico dependiente de la implementación JWT)
```

### Agregar Validaciones Avanzadas
```bash
# Las validaciones ya están incluidas con --validation
# Revisar archivos generados para personalizaciones
```

### Agregar Logging Estructurado
```bash
# El logging ya está configurado en pkg/logger/
# Personalizar según necesidades
```

## 🐛 Troubleshooting

### Error de Conexión a Base de Datos
```bash
# Verificar que PostgreSQL esté ejecutándose
pg_isready -h localhost -p 5432

# Verificar credenciales en .env
cat .env | grep DB_
```

### Error de Módulo No Encontrado
```bash
# Limpiar cache de módulos
go clean -modcache

# Reinstalar dependencias
go mod tidy
```

### Error de Puerto en Uso
```bash
# Verificar qué proceso usa el puerto 8080
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Cambiar puerto en .env
echo "PORT=8081" >> .env
```

## 🎉 Resultado Final

¡Felicitaciones! Ahora tienes:

### ✅ Sistema Completo
- **API REST** funcional en http://localhost:8080
- **4 features** completos: User, Product, Order, OrderItem
- **Base de datos PostgreSQL** configurada
- **Clean Architecture** implementada correctamente

### ✅ Estructura Profesional
```
ecommerce-api/
├── cmd/server/main.go                    # Punto de entrada
├── internal/
│   ├── domain/                           # Entidades y reglas de negocio
│   ├── usecase/                          # Lógica de aplicación
│   ├── repository/                       # Acceso a datos
│   ├── handler/                          # Adaptadores HTTP
│   ├── infrastructure/di/                # Inyección de dependencias
│   └── messages/                         # Mensajes y constantes
├── pkg/
│   ├── config/                           # Configuración
│   ├── logger/                           # Sistema de logging
│   └── auth/                             # Autenticación JWT
└── migrations/                           # Scripts de base de datos
```

### ✅ APIs Disponibles
- `POST /api/v1/users` - Crear usuario
- `GET /api/v1/users` - Listar usuarios
- `GET /api/v1/users/:id` - Obtener usuario
- `PUT /api/v1/users/:id` - Actualizar usuario
- `DELETE /api/v1/users/:id` - Eliminar usuario
- *(Similar para productos, órdenes y order items)*

## 🚀 Próximos Pasos

### Funcionalidades Adicionales
1. **Autenticación JWT** completa
2. **Middleware de autorización** por roles
3. **Documentación Swagger** automática
4. **Tests unitarios** e integración
5. **Docker** para deployment
6. **CI/CD** con GitHub Actions

### Expansión del Sistema
1. **Sistema de categorías** de productos
2. **Carrito de compras** temporal
3. **Sistema de descuentos** y cupones
4. **Notificaciones** por email
5. **Dashboard** de administración

---

**¡Has completado con éxito tu primer proyecto completo con Goca! 🎊**

**← [Primeros Pasos](Getting-Started) | [Ejemplo REST API](Example-REST-API) →**
