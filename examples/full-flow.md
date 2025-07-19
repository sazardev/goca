Basándome en el código de Goca CLI que analicé, aquí tienes una guía completa para crear un proyecto desde cero usando tu paquete:

## 🚀 Crear Proyecto Completo desde Cero con Goca

### 1. **Instalación de Goca CLI**
```bash
# Instalar desde tu repositorio
go install github.com/sazardev/goca@latest

# Verificar instalación
goca version
```

### 2. **Inicializar Proyecto Base**
```bash
# Crear proyecto con Clean Architecture
goca init ecommerce --module github.com/tuusuario/ecommerce --database postgres --auth --api rest

# Entrar al directorio del proyecto
cd ecommerce

# Instalar dependencias
go mod tidy
```

**Esto crea la estructura base:**
```
ecommerce/
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
├── .gitignore
└── README.md
```

### 3. **Generar Features Completos**

#### Feature: Usuario
```bash
goca feature User --fields "name:string,email:string,password:string,role:string,age:int" --validation --business-rules --database postgres --handlers "http,cli"
```

#### Feature: Producto
```bash
goca feature Product --fields "name:string,description:string,price:float64,category:string,stock:int" --validation --business-rules --database postgres --handlers "http,grpc"
```

#### Feature: Orden
```bash
goca feature Order --fields "user_id:int,total:float64,status:string,items:string" --validation --business-rules --database postgres --handlers "http,worker"
```

### 4. **Generar Mensajes y Constantes**
```bash
# Generar mensajes para todos los features
goca messages User --all
goca messages Product --all
goca messages Order --all
```

### 5. **Generar Contenedor de Inyección de Dependencias**
```bash
goca di --features "User,Product,Order" --database postgres
```

### 6. **Estructura Final Generada**
```
ecommerce/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── order.go
│   │   └── errors.go
│   ├── usecase/
│   │   ├── dto.go
│   │   ├── user_usecase.go
│   │   ├── user_service.go
│   │   ├── product_usecase.go
│   │   ├── product_service.go
│   │   ├── order_usecase.go
│   │   ├── order_service.go
│   │   └── interfaces.go
│   ├── repository/
│   │   ├── interfaces.go
│   │   ├── postgres_user_repository.go
│   │   ├── postgres_product_repository.go
│   │   └── postgres_order_repository.go
│   ├── handler/
│   │   ├── http/
│   │   │   ├── user_handler.go
│   │   │   ├── product_handler.go
│   │   │   ├── order_handler.go
│   │   │   └── routes.go
│   │   ├── grpc/
│   │   │   ├── product.proto
│   │   │   └── product_server.go
│   │   ├── cli/
│   │   │   └── user_commands.go
│   │   └── worker/
│   │       └── order_worker.go
│   ├── infrastructure/
│   │   └── di/
│   │       └── container.go
│   └── messages/
│       ├── errors.go
│       └── responses.go
├── pkg/
│   ├── config/
│   ├── logger/
│   └── auth/
├── go.mod
└── README.md
```

### 7. **Configurar Base de Datos**
```sql
-- Crear tablas en PostgreSQL
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    age INTEGER CHECK (age >= 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    category VARCHAR(50) NOT NULL,
    stock INTEGER NOT NULL CHECK (stock >= 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    total DECIMAL(10,2) NOT NULL CHECK (total >= 0),
    status VARCHAR(50) NOT NULL,
    items TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);
```

### 8. **Configurar Variables de Entorno**
```bash
# .env
PORT=8080
DATABASE_URL=postgres://user:password@localhost/ecommerce?sslmode=disable
LOG_LEVEL=info
JWT_SECRET=your-super-secret-key
```

### 9. **Compilar y Ejecutar**
```bash
# Compilar
go build -o ecommerce cmd/server/main.go

# Ejecutar
./ecommerce

# O directamente
go run cmd/server/main.go
```

### 10. **Testing del API**
```bash
# Health check
curl http://localhost:8080/health

# Crear usuario
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Juan Pérez","email":"juan@example.com","password":"secret123","role":"customer","age":30}'

# Crear producto
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"Gaming laptop","price":1500.00,"category":"electronics","stock":10}'

# Listar usuarios
curl http://localhost:8080/users

# Listar productos
curl http://localhost:8080/products
```

### 11. **Comandos CLI Generados**
```bash
# Si generaste handlers CLI, puedes usar:
./ecommerce user create --name "Juan" --email "juan@example.com" --role "admin"
./ecommerce user get 1
./ecommerce product list
```

### 12. **Flujo Completo de Desarrollo**

1. **Inicializar**: `goca init`
2. **Generar features**: `goca feature`
3. **Configurar DI**: `goca di`
4. **Configurar DB**: Scripts SQL
5. **Configurar env**: Variables de entorno
6. **Ejecutar**: `go run`
7. **Probar**: APIs y comandos

### 13. **Comandos Adicionales por Necesidad**

Si necesitas componentes específicos después:

```bash
# Solo interfaces (para TDD)
goca interfaces Payment --all

# Solo repositorio
goca repository Payment --database postgres

# Solo handler específico
goca handler Payment --type grpc --validation

# Mensajes específicos
goca messages Payment --errors --responses
```

### 🎯 **Resultado Final**

Tendrás un **proyecto completo de e-commerce** con:
- ✅ **Clean Architecture** implementada correctamente
- ✅ **APIs REST** para todas las operaciones
- ✅ **Validaciones** en todas las capas
- ✅ **Base de datos PostgreSQL** configurada
- ✅ **Inyección de dependencias** lista
- ✅ **Autenticación JWT** incluida
- ✅ **Logging** estructurado
- ✅ **Configuración** por ambiente
- ✅ **Estructura escalable** para nuevos features

¡Todo esto en **menos de 10 comandos** gracias a tu CLI Goca! 🚀¡Perfecto! Te voy a mostrar cómo crear un proyecto completo desde cero usando tu paquete Goca. Aquí tienes una guía paso a paso:

## 🚀 **Crear un Proyecto Completo con Goca**

### **1. Instalar Goca CLI**
```bash
go install github.com/sazardev/goca@latest
```

### **2. Crear e Inicializar el Proyecto**
```bash
# Crear directorio del proyecto
mkdir mi-ecommerce
cd mi-ecommerce

# Inicializar con Clean Architecture
goca init mi-ecommerce --module github.com/miusuario/mi-ecommerce --database postgres --auth --api rest
```

### **3. Navegar al Proyecto e Instalar Dependencias**
```bash
cd mi-ecommerce
go mod tidy
```

### **4. Generar Features Completos**

**Feature de Usuarios:**
```bash
goca feature User --fields "name:string,email:string,password:string,role:string,active:bool" --database postgres --validation --business-rules
```

**Feature de Productos:**
```bash
goca feature Product --fields "name:string,description:string,price:float64,stock:int,category:string,sku:string" --database postgres --validation --business-rules
```

**Feature de Pedidos:**
```bash
goca feature Order --fields "user_id:int,total:float64,status:string,items:string" --database postgres --validation
```

**Feature de Categorías:**
```bash
goca feature Category --fields "name:string,description:string,parent_id:int" --database postgres --validation
```

### **5. Generar Componentes Específicos Adicionales**

**Handlers adicionales para diferentes protocolos:**
```bash
# Handler gRPC para productos
goca handler Product --type grpc --validation

# Handler CLI para administración
goca handler User --type cli --validation

# Worker para procesamiento asíncrono de pedidos
goca handler Order --type worker --validation
```

**Interfaces para TDD:**
```bash
goca interfaces User --all
goca interfaces Product --all
goca interfaces Order --all
```

**Mensajes y constantes:**
```bash
goca messages User --all
goca messages Product --all
goca messages Order --all
```

### **6. Configurar Inyección de Dependencias**
```bash
goca di --features "User,Product,Order,Category" --database postgres --wire
```

### **7. Estructura Final del Proyecto**

Después de ejecutar todos los comandos, tendrás esta estructura:

```
mi-ecommerce/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── order.go
│   │   └── category.go
│   ├── usecase/
│   │   ├── dto.go
│   │   ├── interfaces.go
│   │   ├── user_service.go
│   │   ├── product_service.go
│   │   ├── order_service.go
│   │   └── category_service.go
│   ├── repository/
│   │   ├── interfaces.go
│   │   ├── postgres_user_repository.go
│   │   ├── postgres_product_repository.go
│   │   ├── postgres_order_repository.go
│   │   └── postgres_category_repository.go
│   ├── handler/
│   │   ├── http/
│   │   │   ├── user_handler.go
│   │   │   ├── product_handler.go
│   │   │   ├── order_handler.go
│   │   │   ├── category_handler.go
│   │   │   └── routes.go
│   │   ├── grpc/
│   │   │   ├── product.proto
│   │   │   └── product_server.go
│   │   ├── cli/
│   │   │   └── user_commands.go
│   │   └── worker/
│   │       └── order_worker.go
│   ├── messages/
│   │   ├── errors.go
│   │   └── responses.go
│   ├── constants/
│   │   └── constants.go
│   └── infrastructure/
│       └── di/
│           ├── container.go
│           └── wire.go
├── pkg/
│   ├── config/
│   │   └── config.go
│   └── logger/
│       └── logger.go
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

### **8. Comandos para Desarrollo**

**Compilar y ejecutar:**
```bash
go build -o mi-ecommerce cmd/server/main.go
./mi-ecommerce
```

**Agregar más features durante el desarrollo:**
```bash
# Feature de inventario
goca feature Inventory --fields "product_id:int,quantity:int,location:string" --database postgres

# Feature de notificaciones
goca feature Notification --fields "user_id:int,message:string,type:string,read:bool" --database postgres
```

### **9. Comandos de Testing**

```bash
# Generar solo interfaces para TDD
goca interfaces Inventory --usecase --repository

# Generar mensajes específicos
goca messages Inventory --errors --responses

# Actualizar DI con nuevos features
goca di --features "User,Product,Order,Category,Inventory,Notification" --database postgres
```

### **10. Ejemplo de Flujo Completo para un Nuevo Feature**

```bash
# 1. Generar feature completo
goca feature Review --fields "user_id:int,product_id:int,rating:int,comment:string" --database postgres --validation

# 2. Generar handler específico si necesitas más protocolos
goca handler Review --type grpc --validation

# 3. Generar mensajes
goca messages Review --all

# 4. Actualizar DI
goca di --features "User,Product,Order,Category,Review" --database postgres

# 5. Compilar y probar
go mod tidy
go build ./...
go test ./...
```

### **💡 Tips Avanzados:**

**Para proyectos con microservicios:**
```bash
# Crear servicio específico para cada dominio
mkdir user-service && cd user-service
goca init user-service --module github.com/miusuario/user-service --database postgres --api grpc
goca feature User --fields "name:string,email:string" --database postgres
```

**Para proyectos con múltiples bases de datos:**
```bash
# Usuarios en PostgreSQL
goca feature User --fields "name:string,email:string" --database postgres

# Logs en MongoDB
goca feature AuditLog --fields "action:string,user_id:string,timestamp:time.Time" --database mongodb
```

**Para agregar autenticación:**
```bash
goca feature Auth --fields "user_id:int,token:string,expires_at:time.Time" --database postgres --validation
goca handler Auth --type http --middleware --validation
```

¡Con estos comandos tendrás un proyecto completo con Clean Architecture, múltiples features, diferentes handlers y toda la estructura necesaria para un ecommerce funcional! 🎉