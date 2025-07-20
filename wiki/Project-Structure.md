# Estructura de Proyecto

Esta página explica la organización de directorios y archivos que Goca genera, siguiendo las mejores prácticas de Clean Architecture en Go.

## 📁 Estructura Completa

```
mi-proyecto/
├── cmd/                              # Puntos de entrada de la aplicación
│   └── server/
│       └── main.go                   # Servidor HTTP principal
├── internal/                         # Código privado de la aplicación
│   ├── domain/                       # 🟡 Capa de Dominio
│   │   ├── user.go                   # Entidades de negocio
│   │   ├── product.go
│   │   └── errors.go                 # Errores del dominio
│   ├── usecase/                      # 🔴 Capa de Casos de Uso
│   │   ├── dto/                      # Data Transfer Objects
│   │   │   ├── user_dto.go
│   │   │   └── product_dto.go
│   │   ├── interfaces/               # Contratos entre capas
│   │   │   ├── user_interfaces.go
│   │   │   └── product_interfaces.go
│   │   ├── user_usecase.go           # Servicios de aplicación
│   │   └── product_usecase.go
│   ├── repository/                   # 🔵 Capa de Infraestructura
│   │   ├── interfaces/               # Interfaces de repositorios
│   │   │   ├── user_repository.go
│   │   │   └── product_repository.go
│   │   ├── postgres/                 # Implementaciones PostgreSQL
│   │   │   ├── user_repository.go
│   │   │   └── product_repository.go
│   │   ├── mysql/                    # Implementaciones MySQL
│   │   └── mongodb/                  # Implementaciones MongoDB
│   ├── handler/                      # 🟢 Capa de Adaptadores
│   │   ├── http/                     # Handlers HTTP REST
│   │   │   ├── dto/                  # DTOs específicos para HTTP
│   │   │   ├── user_handler.go
│   │   │   ├── user_routes.go
│   │   │   ├── product_handler.go
│   │   │   └── middleware/           # Middlewares HTTP
│   │   ├── grpc/                     # Handlers gRPC
│   │   │   ├── user.proto
│   │   │   ├── user_server.go
│   │   │   └── product_server.go
│   │   ├── cli/                      # Comandos CLI
│   │   │   ├── user_commands.go
│   │   │   └── product_commands.go
│   │   └── worker/                   # Workers en background
│   │       ├── user_worker.go
│   │       └── order_worker.go
│   ├── infrastructure/               # Configuración de infraestructura
│   │   ├── di/                       # Inyección de dependencias
│   │   │   ├── container.go
│   │   │   ├── wire.go               # Wire.dev (opcional)
│   │   │   └── wire_gen.go
│   │   ├── database/                 # Configuración de DB
│   │   │   ├── postgres.go
│   │   │   ├── mysql.go
│   │   │   └── migrations/
│   │   └── cache/                    # Configuración de cache
│   │       ├── redis.go
│   │       └── memory.go
│   ├── messages/                     # Mensajes y constantes
│   │   ├── errors.go                 # Mensajes de error
│   │   ├── responses.go              # Mensajes de respuesta
│   │   └── constants.go              # Constantes del sistema
│   └── constants/                    # Constantes específicas por feature
│       ├── user_constants.go
│       └── product_constants.go
├── pkg/                              # Código reutilizable/público
│   ├── config/                       # Configuración de la aplicación
│   │   ├── config.go
│   │   └── database.go
│   ├── logger/                       # Sistema de logging
│   │   ├── logger.go
│   │   └── interfaces.go
│   ├── auth/                         # Sistema de autenticación
│   │   ├── jwt.go
│   │   ├── middleware.go
│   │   └── service.go
│   ├── validator/                    # Validaciones reutilizables
│   │   ├── validator.go
│   │   └── custom_rules.go
│   ├── utils/                        # Utilidades generales
│   │   ├── crypto.go
│   │   ├── time.go
│   │   └── strings.go
│   └── errors/                       # Manejo global de errores
│       ├── errors.go
│       ├── codes.go
│       └── handler.go
├── api/                              # Documentación de APIs
│   ├── openapi/                      # Especificaciones OpenAPI
│   │   ├── swagger.yaml
│   │   └── user.yaml
│   └── proto/                        # Archivos Protocol Buffers
│       ├── user.proto
│       └── product.proto
├── web/                              # Archivos web estáticos (opcional)
│   ├── static/
│   └── templates/
├── docs/                             # Documentación del proyecto
│   ├── architecture.md
│   ├── api.md
│   └── deployment.md
├── scripts/                          # Scripts de automatización
│   ├── build.sh
│   ├── test.sh
│   └── migrate.sh
├── deployments/                      # Configuraciones de despliegue
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   ├── kubernetes/
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   └── terraform/
├── test/                             # Tests de integración y E2E
│   ├── integration/
│   ├── e2e/
│   └── fixtures/
├── migrations/                       # Migraciones de base de datos
│   ├── 001_initial_schema.sql
│   ├── 002_add_users_table.sql
│   └── 003_add_products_table.sql
├── .github/                          # Configuración de GitHub
│   └── workflows/
│       ├── ci.yml
│       └── release.yml
├── go.mod                            # Dependencias del módulo
├── go.sum                            # Checksums de dependencias
├── .env.example                      # Ejemplo de variables de entorno
├── .gitignore                        # Archivos ignorados por Git
├── Makefile                          # Comandos de automatización
├── README.md                         # Documentación principal
└── CHANGELOG.md                      # Historial de cambios
```

## 🏗️ Capas de Clean Architecture

### 🟡 Capa de Dominio (`internal/domain/`)

**Propósito**: Contiene la lógica de negocio central y las reglas empresariales.

**Archivos típicos:**
```
domain/
├── user.go              # Entidad User con métodos de negocio
├── product.go           # Entidad Product con validaciones
├── order.go             # Entidad Order con reglas de negocio
├── errors.go            # Errores específicos del dominio
└── validations.go       # Validaciones de negocio reutilizables
```

**Características:**
- ✅ **Sin dependencias externas**
- ✅ **Entidades ricas** con comportamiento
- ✅ **Reglas de negocio** encapsuladas
- ✅ **Validaciones** de dominio
- ❌ **No debe conocer** infraestructura

### 🔴 Capa de Casos de Uso (`internal/usecase/`)

**Propósito**: Orquesta la lógica de aplicación y coordina entre el dominio y la infraestructura.

**Archivos típicos:**
```
usecase/
├── dto/                          # Data Transfer Objects
│   ├── user_dto.go              # DTOs para operaciones de usuario
│   └── common_dto.go            # DTOs compartidos
├── interfaces/                   # Contratos entre capas
│   ├── user_interfaces.go       # Interfaces de User UseCase
│   └── repositories.go          # Interfaces de repositorios
├── user_usecase.go              # Implementación de casos de uso
├── product_usecase.go           # Casos de uso de productos
└── common_usecase.go            # Lógica compartida
```

**Características:**
- ✅ **Orquesta** flujos de trabajo
- ✅ **DTOs** para transferencia de datos
- ✅ **Interfaces** para desacoplar capas
- ✅ **Validaciones** de aplicación
- ❌ **No debe conocer** detalles de HTTP/DB

### 🟢 Capa de Adaptadores (`internal/handler/`)

**Propósito**: Adapta las interfaces externas (HTTP, gRPC, CLI) a los casos de uso internos.

**Archivos típicos:**
```
handler/
├── http/                         # Adaptadores HTTP
│   ├── dto/                     # DTOs específicos para HTTP
│   │   ├── user_http_dto.go     # Request/Response HTTP
│   │   └── error_dto.go         # Respuestas de error
│   ├── user_handler.go          # Handler para endpoints de usuario
│   ├── user_routes.go           # Definición de rutas
│   └── middleware/              # Middlewares HTTP
│       ├── auth.go              # Middleware de autenticación
│       ├── cors.go              # Middleware de CORS
│       └── logging.go           # Middleware de logging
├── grpc/                        # Adaptadores gRPC
│   ├── user.proto              # Definición de servicios
│   ├── user_server.go          # Implementación del servidor
│   └── interceptors/           # Interceptors gRPC
├── cli/                        # Comandos CLI
│   ├── user_commands.go        # Comandos de usuario
│   └── root.go                 # Comando raíz
└── worker/                     # Workers en background
    ├── user_worker.go          # Worker para usuarios
    └── queue.go                # Configuración de colas
```

**Características:**
- ✅ **Adapta** protocolos externos
- ✅ **DTOs específicos** por protocolo
- ✅ **Manejo de errores** apropiado
- ✅ **Validación** de entrada
- ❌ **No debe contener** lógica de negocio

### 🔵 Capa de Infraestructura (`internal/repository/`, `pkg/`)

**Propósito**: Implementa detalles técnicos como persistencia, logging, configuración.

**Archivos típicos:**
```
repository/
├── interfaces/                   # Contratos de persistencia
│   ├── user_repository.go       # Interface del repositorio
│   └── transaction.go           # Interface de transacciones
├── postgres/                    # Implementación PostgreSQL
│   ├── user_repository.go      # Repositorio específico
│   ├── migrations.go           # Migraciones
│   └── connection.go           # Configuración de conexión
├── mysql/                      # Implementación MySQL
├── mongodb/                    # Implementación MongoDB
└── memory/                     # Implementación en memoria (tests)
    └── user_repository.go
```

**Características:**
- ✅ **Implementa interfaces** del dominio
- ✅ **Detalles específicos** de tecnología
- ✅ **Configuración** de conexiones
- ✅ **Migraciones** de base de datos
- ❌ **No debe exponer** detalles técnicos

## 📦 Directorios Especiales

### `cmd/` - Puntos de Entrada

```
cmd/
├── server/                      # Servidor HTTP principal
│   └── main.go
├── migrate/                     # Herramienta de migraciones
│   └── main.go
├── worker/                      # Worker en background
│   └── main.go
└── cli/                        # Herramienta CLI
    └── main.go
```

**Propósito**: Cada subdirectorio representa un ejecutable diferente.

### `pkg/` - Código Reutilizable

```
pkg/
├── config/                      # Configuración global
├── logger/                      # Sistema de logging
├── auth/                        # Autenticación/autorización
├── validator/                   # Validaciones reutilizables
├── utils/                       # Utilidades generales
└── errors/                      # Manejo global de errores
```

**Propósito**: Código que puede ser importado por otros proyectos.

### `api/` - Documentación de APIs

```
api/
├── openapi/                     # Especificaciones OpenAPI/Swagger
│   ├── swagger.yaml            # Documentación principal
│   ├── user.yaml               # Endpoints de usuario
│   └── product.yaml            # Endpoints de producto
└── proto/                      # Protocol Buffers para gRPC
    ├── user.proto
    ├── product.proto
    └── common.proto
```

**Propósito**: Contratos de API y documentación externa.

### `migrations/` - Esquema de Base de Datos

```
migrations/
├── 001_initial_schema.sql       # Esquema inicial
├── 002_add_users_table.sql      # Agregar tabla usuarios
├── 003_add_products_table.sql   # Agregar tabla productos
└── 004_add_indexes.sql          # Agregar índices
```

**Propósito**: Control de versiones del esquema de base de datos.

## 🔄 Flujo de Dependencias

```
┌─────────────────────────────────────────┐
│                 cmd/                    │ ← Puntos de entrada
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            internal/handler/            │ ← 🟢 Adaptadores
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            internal/usecase/            │ ← 🔴 Casos de Uso
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            internal/domain/             │ ← 🟡 Dominio
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│          internal/repository/           │ ← 🔵 Infraestructura
└─────────────────────────────────────────┘
```

**Regla Fundamental**: Las dependencias siempre apuntan hacia adentro.

## 📁 Convenciones de Nomenclatura

### Archivos
- **Entidades**: `user.go`, `product.go`, `order.go`
- **DTOs**: `user_dto.go`, `create_user_request.go`
- **UseCase**: `user_usecase.go`, `product_service.go`
- **Repository**: `user_repository.go`, `postgres_user_repo.go`
- **Handler**: `user_handler.go`, `user_routes.go`
- **Interfaces**: `user_interfaces.go`, `repositories.go`

### Packages
- **Lowercase**: Siempre en minúsculas
- **Descriptivos**: `usecase`, `repository`, `handler`
- **Sin guiones**: `userservice` no `user-service`
- **Singulares**: `user` no `users` (excepto cuando sea apropiado)

### Estructuras
```go
// Entidades: PascalCase
type User struct {}
type OrderItem struct {}

// Interfaces: PascalCase + sufijo
type UserRepository interface {}
type UserUseCase interface {}

// DTOs: PascalCase + propósito
type CreateUserRequest struct {}
type UserResponse struct {}
```

## 🎯 Beneficios de esta Estructura

### ✅ Separación Clara de Responsabilidades
- Cada directorio tiene un propósito específico
- Fácil localizar código relacionado
- Cambios en una capa no afectan otras

### ✅ Testabilidad
- Interfaces permiten mocks fáciles
- Tests unitarios por capa
- Tests de integración separados

### ✅ Escalabilidad
- Agregar nuevos features es predecible
- Estructura consistente entre features
- Fácil onboarding de nuevos desarrolladores

### ✅ Mantenibilidad
- Código organizado y predecible
- Refactoring seguro por capas
- Dependencias explícitas

### ✅ Flexibilidad
- Fácil cambiar implementaciones
- Soporte múltiples protocolos
- Agregar nuevas funcionalidades sin romper existentes

## 🛠️ Personalización

### Agregar Nueva Capa
```bash
# Crear nueva capa de eventos
mkdir -p internal/events
mkdir -p internal/events/handlers
mkdir -p internal/events/publishers
```

### Agregar Nuevo Protocolo
```bash
# Agregar soporte para GraphQL
mkdir -p internal/handler/graphql
mkdir -p internal/handler/graphql/resolvers
mkdir -p internal/handler/graphql/schemas
```

### Agregar Nueva Base de Datos
```bash
# Agregar soporte para Redis
mkdir -p internal/repository/redis
mkdir -p pkg/cache/redis
```

## 📚 Recursos Adicionales

- [Clean Architecture Principles](Clean-Architecture)
- [Design Patterns Used](Design-Patterns)
- [Testing Guide](Testing-Guide)
- [Best Practices](Best-Practices)

---

**← [Clean Architecture](Clean-Architecture) | [Patrones Implementados](Design-Patterns) →**
