# Comando goca integrate

El comando `goca integrate` es una herramienta especializada para **integrar features existentes** que no están conectados con el contenedor de inyección de dependencias o `main.go`.

## 🎯 Propósito

Útil para:
- 📦 Proyectos migrados de versiones anteriores de Goca
- 🔧 Features generados manualmente que necesitan integración
- 🛠️ Reparar integraciones incompletas o dañadas
- 🔄 Actualizar proyectos existentes con la nueva auto-integración

## 📋 Sintaxis

```bash
goca integrate [flags]
```

## 🚩 Flags Disponibles

| Flag         | Tipo     | Requerido | Descripción                                          |
| ------------ | -------- | --------- | ---------------------------------------------------- |
| `--all`      | `bool`   | ❌ No      | Detecta e integra automáticamente todos los features |
| `--features` | `string` | ❌ No      | Features específicos a integrar (`"User,Product"`)   |

## 📖 Ejemplos de Uso

### Integración Automática (Recomendado)
```bash
# Detecta automáticamente todos los features y los integra
goca integrate --all
```

**Salida esperada:**
```
🔍 Detectando features existentes...
📋 Features detectados: User, Product, Order

🔄 Iniciando proceso de integración...

1️⃣  Configurando contenedor DI...
   📦 Creando contenedor DI...
   ✅ User integrado en el contenedor DI
   ✅ Product integrado en el contenedor DI
   ✅ Order integrado en el contenedor DI

2️⃣  Actualizando main.go...
   📍 Actualizando main.go en: main.go
   🔧 Reescribiendo main.go completo...
   ✅ main.go creado con 3 features

3️⃣  Verificando integración...
   ✅ Contenedor DI existe
   ✅ main.go integrado (main.go)
   ✅ User routes integradas
   ✅ Product routes integradas
   ✅ Order routes integradas

🎯 ¡Integración perfecta! Todo está listo.

🎉 ¡Integración completada!
✅ Todos los features están ahora:
   🔗 Conectados en el contenedor DI
   🛣️  Con rutas registradas en main.go
   ⚡ Listos para usar
```

### Integración Específica
```bash
# Integra solo features específicos
goca integrate --features "User,Product"
```

### Caso de Uso Común: Después de Clonar un Proyecto
```bash
# 1. Clonar proyecto existente
git clone https://github.com/user/my-goca-project.git
cd my-goca-project

# 2. Integrar automáticamente todos los features
goca integrate --all

# 3. Verificar que todo funciona
go mod tidy
go run main.go
```

## 🔍 Detección Automática

El comando `goca integrate --all` detecta features automáticamente buscando:

1. **Entidades de dominio** en `internal/domain/*.go`
2. **Handlers HTTP** en `internal/handler/http/*_handler.go`
3. **Casos de uso** en `internal/usecase/*.go`

### Archivos Ignorados en Detección
- `errors.go`
- `validations.go`
- `common.go`
- `types.go`

## 🔧 Proceso de Integración

### 1. Contenedor DI
- ✅ Crea `internal/di/container.go` si no existe
- ✅ Agrega campos para repositories, use cases y handlers
- ✅ Configura métodos de setup y getters
- ✅ Detecta y evita duplicados

### 2. Main.go
- ✅ Detecta `main.go` en múltiples ubicaciones
- ✅ Agrega imports necesarios (`internal/di`)
- ✅ Configura contenedor DI
- ✅ Registra todas las rutas HTTP
- ✅ Preserva configuración existente

### 3. Rutas Generadas
Para cada feature se registran automáticamente:
```go
// User routes
userHandler := container.UserHandler()
router.HandleFunc("/api/v1/users", userHandler.CreateUser).Methods("POST")
router.HandleFunc("/api/v1/users/{id}", userHandler.GetUser).Methods("GET")
router.HandleFunc("/api/v1/users/{id}", userHandler.UpdateUser).Methods("PUT")
router.HandleFunc("/api/v1/users/{id}", userHandler.DeleteUser).Methods("DELETE")
router.HandleFunc("/api/v1/users", userHandler.ListUsers).Methods("GET")
```

## 🏗️ Estructura Esperada del Proyecto

Para que la integración funcione correctamente, el proyecto debe tener:

```
myproject/
├── go.mod
├── main.go                          # Será actualizado/creado
├── internal/
│   ├── domain/
│   │   ├── user.go                  # ← Detectado como feature "User"
│   │   └── product.go               # ← Detectado como feature "Product"
│   ├── usecase/
│   │   ├── user_usecase.go
│   │   └── product_usecase.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   └── product_repository.go
│   ├── handler/
│   │   └── http/
│   │       ├── user_handler.go      # ← Detectado como feature "User"
│   │       └── product_handler.go   # ← Detectado como feature "Product"
│   └── di/                          # ← Será creado/actualizado
│       └── container.go
└── pkg/
    ├── config/
    └── logger/
```

## ⚠️ Casos Especiales

### Main.go No Encontrado
Si no se encuentra `main.go`, el comando creará uno nuevo completo:
```
⚠️  main.go no encontrado, creando nuevo...
✅ main.go creado con 3 features
```

### Features Ya Integrados
Si algunos features ya están integrados, se saltarán:
```
✅ User ya está en el contenedor DI
➕ Agregando Product al contenedor DI...
✅ Product integrado en el contenedor DI
```

### Errores de Integración
Si hay problemas, se muestran instrucciones manuales:
```
⚠️  No se pudo actualizar main.go: permission denied

📋 Instrucciones de integración manual:
1. Agregar import en main.go:
   "myproject/internal/di"
2. Agregar en main(), después de conectar la DB:
   container := di.NewContainer(db)
3. Agregar las rutas del feature:
   userHandler := container.UserHandler()
   router.HandleFunc("/api/v1/users", userHandler.CreateUser).Methods("POST")
   ...
```

## 🔄 Workflow Completo

### Escenario: Migrar Proyecto Existente

```bash
# 1. Verificar structure actual
ls internal/domain/     # Ver features existentes

# 2. Ejecutar integración automática
goca integrate --all

# 3. Verificar que todo compile
go mod tidy
go build

# 4. Probar el servidor
go run main.go

# 5. Probar endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/users
```

### Escenario: Agregar Feature a Proyecto Integrado

```bash
# 1. Generar nuevo feature (auto-integrado)
goca feature Order --fields "user_id:int,total:float64,status:string"

# 2. ¡Ya está listo! No necesitas integrate
go run main.go
```

## 🤝 Comparación con `goca feature`

| Aspecto              | `goca feature`         | `goca integrate`                |
| -------------------- | ---------------------- | ------------------------------- |
| **Propósito**        | Crear nuevo feature    | Integrar features existentes    |
| **Genera código**    | ✅ Sí (todas las capas) | ❌ No (solo integración)         |
| **Auto-integración** | ✅ Sí (automática)      | ✅ Sí (reparación/actualización) |
| **Uso típico**       | Desarrollo nuevo       | Migración/reparación            |
| **Detección**        | No aplica              | ✅ Sí (automática)               |

## 💡 Consejos y Mejores Prácticas

### ✅ Recomendaciones
- **Usa `--all`** para detección automática en lugar de especificar features manualmente
- **Ejecuta después de clonar** proyectos existentes de Goca
- **Combina con `go mod tidy`** después de la integración
- **Haz backup** de `main.go` antes de integrar si tienes código custom

### ⚠️ Precauciones
- **Revisa `main.go`** después de la integración si tenías configuraciones especiales
- **Verifica rutas** si ya tenías endpoints custom registrados
- **Compila después** de la integración para verificar que todo funciona

### 🔄 Integración Continua
```bash
# Script para CI/CD
#!/bin/bash
go mod download
goca integrate --all
go mod tidy
go test ./...
go build
```

## 🆘 Solución de Problemas

### Problema: "No se encontraron features"
```bash
# Verificar estructura
ls internal/domain/
ls internal/handler/http/

# Si los archivos existen pero no se detectan:
goca integrate --features "User,Product"  # Especificar manualmente
```

### Problema: "main.go no se pudo actualizar"
```bash
# Verificar permisos
ls -la main.go

# Si es problema de permisos en Windows:
# Ejecutar terminal como administrador

# Alternativa: integración manual
goca integrate --features "User" --dry-run  # Ver instrucciones
```

### Problema: "Rutas duplicadas"
```bash
# El comando detecta rutas existentes y las salta automáticamente
# Si hay conflictos, revisar main.go manualmente
```

## 🔗 Comandos Relacionados

- [`goca feature`](Command-Feature.md) - Crear nuevos features (incluye auto-integración)
- [`goca di`](Command-DI.md) - Generar solo contenedor DI
- [`goca init`](Command-Init.md) - Inicializar nuevo proyecto

---

**Próximo**: [Comando DI](Command-DI.md) | **Anterior**: [Comando Feature](Command-Feature.md) | **Índice**: [Comandos](README.md)
