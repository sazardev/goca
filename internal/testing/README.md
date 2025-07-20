# Goca CLI Testing Framework

Este directorio contiene el framework de testing comprensivo para Goca CLI que verifica que el código generado esté correcto, compile sin errores y respete todas las convenciones.

## 🎯 Objetivo

Como solicitaste: **"me gustaría que cuando generemos código esté todo bien hecho, detecte las ubicaciones, las ponga acorde, limpiamente y genere 0 errores o alertas el código"**

Este framework de testing garantiza que:
- ✅ **Cero errores de compilación**
- ✅ **Cero alertas de linting**
- ✅ **Estructura de archivos correcta**
- ✅ **Ubicaciones adecuadas de archivos**
- ✅ **Código limpio y bien formateado**
- ✅ **Validación de todas las flags y opciones**

## 📁 Estructura

```
internal/testing/
├── suite.go                 # Suite principal de testing
├── comprehensive_test.go     # Tests comprensivos
├── validator.go             # Validadores de código
├── test_runner.go           # Ejecutor de tests standalone
├── scenarios.go             # Escenarios de testing específicos
├── cli.go                   # Utilidades para testing CLI
├── errors.go                # Manejo de errores de testing
└── architecture.go          # Validación de arquitectura
```

## 🚀 Uso Rápido

### Desde Makefile (Recomendado)

```bash
# Ejecutar todos los tests comprensivos
make test-cli

# Test solo el comando init
make test-cli-init

# Test solo el comando feature 
make test-cli-feature

# Test solo calidad del código
make test-cli-quality

# Tests rápidos
make test-cli-fast

# Benchmarks de rendimiento
make test-cli-benchmark

# Todos los tests (unit + CLI)
make test-all
```

### Ejecución Manual

```bash
# Test comprensivo completo
go run internal/testing/test_runner.go -type=all -v

# Tests específicos
go run internal/testing/test_runner.go -type=init -v
go run internal/testing/test_runner.go -type=feature -v
go run internal/testing/test_runner.go -type=entity -v
go run internal/testing/test_runner.go -type=quality -v

# Sin limpiar archivos temporales (para debugging)
go run internal/testing/test_runner.go -type=all -v -cleanup=false
```

### Tests Unitarios Go

```bash
# Test específico
go test ./internal/testing -run TestGocaCLIComprehensive -v

# Test de compilación únicamente
go test ./internal/testing -run TestCodeQuality -v

# Benchmarks
go test ./internal/testing -bench=. -benchmem
```

## 🧪 Qué Se Testa

### 1. Comando `goca init`
- [x] Estructura de directorios correcta
- [x] Archivos base generados
- [x] go.mod con módulo correcto
- [x] Configuración inicial válida
- [x] Compilación exitosa

### 2. Comando `goca feature`
- [x] Entidades de dominio con validación
- [x] Use cases con interfaces correctas
- [x] Repositorios con implementación Postgres
- [x] Handlers HTTP con endpoints
- [x] Mensajes de error y respuesta
- [x] Integración completa sin errores

### 3. Comando `goca entity`
- [x] Campos con tipos correctos
- [x] Métodos de validación
- [x] Flags: `--validation`, `--timestamps`, `--soft-delete`, `--business-rules`
- [x] Sintaxis Go correcta

### 4. Comandos Adicionales
- [x] `goca usecase` - Use cases con operaciones CRUD
- [x] `goca repository` - Repositorios con diferentes bases de datos
- [x] `goca handler` - Handlers para HTTP, gRPC, CLI, Worker
- [x] `goca messages` - Mensajes de error, respuesta y constantes
- [x] `goca di` - Dependency injection con y sin Wire
- [x] `goca interfaces` - Interfaces para todas las capas

### 5. Calidad del Código
- [x] **Compilación**: `go build ./...` sin errores
- [x] **Linting**: `go vet ./...` sin alertas
- [x] **Formato**: `gofmt` aplicado correctamente
- [x] **Sintaxis**: Parser Go valida todos los archivos
- [x] **Imports**: Sin conflictos ni duplicados
- [x] **Convenciones**: Nombres exportados correctos

### 6. Validaciones Específicas
- [x] Entidades con campos requeridos
- [x] Interfaces con métodos esperados
- [x] Handlers con endpoints correctos
- [x] Repositorios con operaciones CRUD
- [x] Use cases con DTOs apropiados
- [x] Manejo de errores consistente

## 📊 Reportes

El framework genera reportes detallados:

```
📊 Test Suite Results:
✅ Generated project: /tmp/goca_test_xxx/testproject
⚠️  Warnings: 0
❌ Errors: 0

🎉 All tests passed successfully!
```

### Tipos de Resultados:
- **✅ PASS**: Test completado exitosamente
- **⚠️ WARNING**: Problemas menores detectados
- **❌ ERROR**: Errores críticos que fallan el test

## 🔧 Configuración

### Variables de Entorno

```bash
# Habilitar modo CI
export CI=true

# Habilitar testing de Goca
export GOCA_TEST=true

# Directorio temporal personalizado
export GOCA_TEST_DIR=/custom/test/dir
```

### Flags de Test Runner

```bash
-type=all|init|feature|entity|quality  # Tipo de test
-v                                      # Verbose output
-cleanup=true|false                     # Limpiar archivos temp
```

## 🎯 Escenarios de Testing

### Escenario 1: Proyecto Básico
```bash
goca init myproject --module=github.com/user/myproject --database=postgres --auth --api=rest
```

### Escenario 2: Feature Completa
```bash
goca feature User --fields="name:string,email:string,age:int" --validation --business-rules
```

### Escenario 3: Entidad con Todo
```bash
goca entity Product --fields="name:string,price:float64" --validation --timestamps --soft-delete --business-rules
```

### Escenario 4: Diferentes Bases de Datos
```bash
goca repository User --database=postgres --cache --transactions
goca repository Product --database=mysql --cache
goca repository Order --database=mongodb --transactions
```

## 🚨 Debugging

### Ver Archivos Generados
```bash
# No limpiar archivos temporales
go run internal/testing/test_runner.go -type=init -v -cleanup=false
```

### Ejecutar Test Específico
```bash
# Solo compilación
go test ./internal/testing -run TestCodeCompilation -v

# Solo estructura de proyecto
go test ./internal/testing -run TestInitCommand -v
```

### Logs Detallados
```bash
# Máximo verbosity
go run internal/testing/test_runner.go -type=all -v 2>&1 | tee test.log
```

## 🎉 Ejemplo de Éxito

Cuando todo funciona correctamente:

```
🚀 Starting comprehensive Goca CLI test suite...
Testing goca init command...
✅ goca init command passed
Testing goca feature command...
Testing feature: User with string and int fields
✅ Feature User passed
Testing feature: Product with float and multiple types  
✅ Feature Product passed
Testing feature: Order with complex fields
✅ Feature Order passed
Testing goca entity command...
Testing entity: Simple entity
✅ Entity Customer passed
Testing entity: Entity with timestamps
✅ Entity Invoice passed
Testing entity: Entity with soft delete
✅ Entity Category passed
Testing entity: Full entity with all features
✅ Entity Employee passed
Testing code compilation...
✅ Code compilation passed
Testing code linting...
✅ Code linting passed
Testing code formatting...
✅ Code formatting passed

📊 Test Suite Results:
✅ Generated project: /tmp/goca_test_xxx/testproject
⚠️  Warnings: 0
❌ Errors: 0

🎉 All tests passed successfully!
```

## 🔄 Integración Continua

### GitHub Actions
```yaml
- name: Test Goca CLI
  run: |
    make build
    make test-cli
```

### Pre-commit Hook
```bash
#!/bin/bash
make test-cli-fast
```

Este framework garantiza que **"cuando generemos código esté todo bien hecho, detecte las ubicaciones, las ponga acorde, limpiamente y genere 0 errores o alertas el código"** como solicitaste. 🎯
