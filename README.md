# Goca - Go Clean Architecture Code Generator

[![Go Version](https://img.shields.io/badge/Go-1.24.5+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/sazardev/goca)

Goca es un potente generador de código CLI para Go que te ayuda a crear proyectos con Clean Architecture siguiendo las mejores prácticas. Genera código limpio y bien estructurado por capas, permitiéndote enfocarte en la lógica de negocio en lugar de tareas repetitivas de configuración.

## 🏗️ Filosofía de Clean Architecture

Cada feature generado por Goca sigue estrictamente los principios de Clean Architecture:

- **🟡 Capa de Dominio**: Entidades puras sin dependencias externas
- **🔴 Capa de Casos de Uso**: Lógica de aplicación con DTOs y validaciones de negocio
- **🟢 Capa de Adaptadores**: Interfaces HTTP, gRPC, CLI que adaptan entrada/salida
- **🔵 Capa de Infraestructura**: Repositorios que implementan persistencia de datos

### ✅ Garantías de Buenas Prácticas

- Dependencias orientadas hacia el núcleo del sistema
- Interfaces y contratos claros entre capas
- Lógica de negocio encapsulada en capas internas
- Responsabilidades claramente segregadas
- Inyección de dependencias para máxima testabilidad

### 🚫 Prevención de Malas Prácticas

- Evita mezclar lógica técnica con lógica de negocio
- Impide dependencias directas de entidades hacia infraestructura
- Genera paquetes bien estructurados y cohesivos

## 🧠 Principios y Anti-Patrones Implementados

### ✅ Patrones Aplicados
- **Repository Pattern**: Abstracción de persistencia de datos
- **Dependency Injection**: Inversión de control entre capas
- **CQRS**: Separación de comandos y consultas en casos de uso
- **Interface Segregation**: Contratos específicos por responsabilidad

### 🚫 Anti-Patrones Prevenidos
- **Fat Controller**: Lógica de negocio en handlers
- **God Object**: Entidades con demasiadas responsabilidades
- **Anemic Domain Model**: Entidades sin comportamiento
- **Direct Database Access**: Dependencias directas a infraestructura

## 🔍 Validación de Clean Architecture

Goca garantiza que cada archivo generado cumple con:

- **Regla de Dependencias**: Código interno nunca depende de código externo
- **Separación de Responsabilidades**: Cada capa tiene una única razón para cambiar
- **Principio de Inversión**: Detalles dependen de abstracciones
- **Interfaces Limpias**: Contratos claros entre capas

## 🚀 Features Principales

- **Generación por Capas**: Cada comando genera código específico para una capa de Clean Architecture
- **Feature Completo**: Un comando genera toda la estructura necesaria para un feature
- **Entidades de Dominio**: Genera entidades puras con validaciones de negocio
- **Casos de Uso**: Crea servicios de aplicación con DTOs bien definidos
- **Repositorios**: Genera interfaces y implementaciones siguiendo Repository Pattern
- **Handlers Multi-Protocolo**: Soporta HTTP, gRPC, CLI manteniendo separación de capas
- **Inyección de Dependencias**: Estructura preparada para DI desde el inicio

## 📦 Instalación

### Usando Go Install (Recomendado)
```bash
go install github.com/sazardev/goca@latest
```

### Descarga de Binarios
Descarga el binario directamente desde [GitHub Releases](https://github.com/sazardev/goca/releases):

**Windows:**
```bash
# Descargar goca-windows-amd64.exe desde releases
# Renombrar a goca.exe y agregar al PATH
```

**Linux:**
```bash
# Descargar y hacer ejecutable
wget https://github.com/sazardev/goca/releases/latest/download/goca-linux-amd64
chmod +x goca-linux-amd64
sudo mv goca-linux-amd64 /usr/local/bin/goca
```

**macOS:**
```bash
# Intel Macs
wget https://github.com/sazardev/goca/releases/latest/download/goca-darwin-amd64
chmod +x goca-darwin-amd64
sudo mv goca-darwin-amd64 /usr/local/bin/goca

# Apple Silicon Macs  
wget https://github.com/sazardev/goca/releases/latest/download/goca-darwin-arm64
chmod +x goca-darwin-arm64
sudo mv goca-darwin-arm64 /usr/local/bin/goca
```

### Desde Código Fuente
```bash
git clone https://github.com/sazardev/goca.git
cd goca
go build -o goca
```

## 🛠️ Inicio Rápido

### Inicializar Proyecto Clean Architecture
```bash
# Crear nuevo proyecto con estructura Clean Architecture
goca init myproject --module github.com/sazardev/myproject

# Navegar al proyecto
cd myproject

# Instalar dependencias
go mod tidy
```

### Generar Feature Completo
```bash
# Generar feature completo con todas las capas
goca feature Employee --fields "name:string,email:string,role:string" --database postgres
```

## 🔄 Flujo de Trabajo Recomendado

1. **Generar Dominio**: `goca entity Employee --fields "name:string,email:string"`
2. **Generar Casos de Uso**: `goca usecase EmployeeService --entity Employee`
3. **Generar Repositorio**: `goca repository Employee --database postgres`
4. **Generar Handlers**: `goca handler Employee --type http`
5. **Generar DI**: `goca di --features Employee`

## 📚 Comandos por Capas

### 🟡 Capa de Dominio

#### `goca entity`
Genera entidades de dominio puras siguiendo DDD.
```bash
goca entity <name> [flags]

# Flags:
--fields string     Campos de la entidad "name:type,email:string"
--validation       Agregar validaciones de dominio
--business-rules   Incluir métodos de reglas de negocio
```

**Ejemplo:**
```bash
goca entity Product --fields "name:string,price:float64,category:string" --validation --business-rules
```

**Código Generado:**
```go
// domain/product.go
package domain

type Product struct {
    ID       int
    Name     string
    Price    float64
    Category string
}

func (p *Product) Validate() error {
    if p.Name == "" || p.Price <= 0 {
        return ErrInvalidProductData
    }
    return nil
}

func (p *Product) IsExpensive() bool {
    return p.Price > 1000.0
}
```

### 🔴 Capa de Casos de Uso

#### `goca usecase`
Genera servicios de aplicación con DTOs y lógica de negocio.
```bash
goca usecase <name> [flags]

# Flags:
--entity string    Entidad asociada
--operations string Operaciones CRUD (create,read,update,delete,list)
--dto-validation   DTOs con validaciones específicas
```

**Ejemplo:**
```bash
goca usecase ProductService --entity Product --operations "create,read,update,delete,list" --dto-validation
```

**Código Generado:**
```go
// usecase/product_service.go
package usecase

import "myproject/domain"

type CreateProductInput struct {
    Name     string `validate:"required,min=3"`
    Price    float64 `validate:"required,gt=0"`
    Category string `validate:"required"`
}

type CreateProductOutput struct {
    Product domain.Product
    Message string
}

type ProductUseCase interface {
    CreateProduct(input CreateProductInput) (CreateProductOutput, error)
    GetProduct(id int) (domain.Product, error)
    UpdateProduct(id int, input UpdateProductInput) error
    DeleteProduct(id int) error
    ListProducts() ([]domain.Product, error)
}
```

### 🟢 Capa de Adaptadores (Handlers)

#### `goca handler`
Genera adaptadores de entrega para diferentes protocolos.
```bash
goca handler <entity> [flags]

# Flags:
--type string     Tipo de handler (http, grpc, cli, worker, soap)
--middleware      Incluir setup de middleware
--validation      Validación de entrada en handler
```

**Ejemplo HTTP:**
```bash
goca handler Product --type http --middleware --validation
```

**Código Generado:**
```go
// handler/http/product_handler.go
package http

import (
    "encoding/json"
    "net/http"
    "myproject/usecase"
)

type ProductHandler struct {
    usecase usecase.ProductUseCase
}

func NewProductHandler(uc usecase.ProductUseCase) *ProductHandler {
    return &ProductHandler{usecase: uc}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
    var input usecase.CreateProductInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    output, err := h.usecase.CreateProduct(input)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(output)
}
```

### 🔵 Capa de Infraestructura

#### `goca repository`
Genera repositorios con interfaces e implementaciones.
```bash
goca repository <entity> [flags]

# Flags:
--database string  Tipo de base de datos (postgres, mysql, mongodb)
--interface-only   Solo generar interfaces
--implementation   Solo generar implementación
```

**Ejemplo:**
```bash
goca repository Product --database postgres
```

**Código Generado:**
```go
// repository/interfaces/product_repository.go
package interfaces

import "myproject/domain"

type ProductRepository interface {
    Save(product *domain.Product) error
    FindByID(id int) (*domain.Product, error)
    FindAll() ([]domain.Product, error)
    Update(product *domain.Product) error
    Delete(id int) error
}

// repository/postgres/product_repository.go
package postgres

import (
    "database/sql"
    "myproject/domain"
    "myproject/repository/interfaces"
)

type postgresProductRepository struct {
    db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) interfaces.ProductRepository {
    return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) Save(product *domain.Product) error {
    query := `INSERT INTO products (name, price, category) VALUES ($1, $2, $3) RETURNING id`
    err := r.db.QueryRow(query, product.Name, product.Price, product.Category).Scan(&product.ID)
    return err
}
```

### Comandos Auxiliares

#### `goca messages`
Genera archivos de mensajes y constantes.
```bash
goca messages <entity> [flags]

# Flags:
--errors       Generar mensajes de error
--responses    Generar mensajes de respuesta
--constants    Generar constantes del feature
```

## 📁 Estructura Completa por Feature (Ejemplo: Employee)

```
employee/
├── domain/
│   ├── employee.go          # Entidad pura
│   ├── errors.go           # Errores de dominio
│   └── validations.go      # Validaciones de negocio
├── usecase/
│   ├── dto.go              # DTOs de entrada/salida
│   ├── employee_usecase.go # Interfaz de casos de uso
│   ├── employee_service.go # Implementación de casos de uso
│   └── interfaces.go       # Contratos hacia otras capas
├── repository/
│   ├── interfaces.go       # Contratos de persistencia
│   ├── postgres_employee_repo.go  # Implementación PostgreSQL
│   └── memory_employee_repo.go    # Implementación en memoria
├── handler/
│   ├── http/
│   │   ├── dto.go          # DTOs específicos de HTTP
│   │   └── handler.go      # Controlador HTTP
│   ├── grpc/
│   │   ├── employee.proto  # Definición gRPC
│   │   └── server.go       # Servidor gRPC
│   ├── cli/
│   │   └── commands.go     # Comandos CLI
│   ├── worker/
│   │   └── worker.go       # Workers/Jobs
│   └── soap/
│       └── soap_client.go  # Cliente SOAP
├── messages/
│   ├── errors.go           # Mensajes de error
│   └── responses.go        # Mensajes de respuesta
├── constants/
│   └── constants.go        # Constantes del feature
└── main.go                 # Punto de entrada
```

## 📋 Buenas vs Malas Prácticas por Capa

### 🟡 Dominio - Qué SÍ y NO hacer

#### ✅ Buenas Prácticas:
```go
// ✅ Entidad pura con validaciones de negocio
type Employee struct {
    ID    int
    Name  string
    Email string
    Role  string
}

func (e *Employee) Validate() error {
    if e.Name == "" || e.Email == "" {
        return ErrInvalidEmployeeData
    }
    return nil
}

func (e *Employee) IsManager() bool {
    return e.Role == "manager"
}
```

#### ❌ Malas Prácticas:
```go
// ❌ NUNCA: Dependencias de infraestructura en dominio
type Employee struct {
    ID int
    DB *sql.DB // ❌ Dependencia externa
}

// ❌ NUNCA: Lógica técnica en dominio
func (e *Employee) SaveToDatabase() error // ❌ Responsabilidad incorrecta

// ❌ NUNCA: Importar paquetes de capas externas
import "myproject/handler/http" // ❌ Violación de dependencias
```

### 🔴 Casos de Uso - Qué SÍ y NO hacer

#### ✅ Buenas Prácticas:
```go
// ✅ DTOs bien definidos
type CreateEmployeeInput struct {
    Name  string `validate:"required"`
    Email string `validate:"required,email"`
}

// ✅ Interfaces hacia otras capas
type EmployeeRepository interface {
    Save(*domain.Employee) error
}

// ✅ Lógica de aplicación pura
func (s *EmployeeService) CreateEmployee(input CreateEmployeeInput) error {
    emp := domain.Employee{Name: input.Name, Email: input.Email}
    if err := emp.Validate(); err != nil {
        return err
    }
    return s.repo.Save(&emp)
}
```

#### ❌ Malas Prácticas:
```go
// ❌ NUNCA: Dependencias de infraestructura directas
func (s *EmployeeService) CreateEmployee(db *sql.DB) error // ❌ Acoplamiento

// ❌ NUNCA: Lógica de presentación
func (s *EmployeeService) CreateEmployeeJSON() string // ❌ Responsabilidad incorrecta

// ❌ NUNCA: Detalles de implementación
func (s *EmployeeService) CreateEmployeeWithPostgres() error // ❌ Especificidad técnica
```

### 🟢 Adaptadores - Qué SÍ y NO hacer

#### ✅ Buenas Prácticas:
```go
// ✅ Solo transformación de datos
func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
    var httpInput HTTPCreateEmployeeInput
    json.NewDecoder(r.Body).Decode(&httpInput)
    
    usecaseInput := usecase.CreateEmployeeInput{
        Name:  httpInput.Name,
        Email: httpInput.Email,
    }
    
    err := h.usecase.CreateEmployee(usecaseInput)
    // Manejar respuesta HTTP
}
```

#### ❌ Malas Prácticas:
```go
// ❌ NUNCA: Lógica de negocio en handlers
func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
    // ❌ Validaciones de negocio aquí
    if employee.Salary < 0 {
        return errors.New("invalid salary")
    }
}

// ❌ NUNCA: Acceso directo a repositorios
func (h *EmployeeHandler) CreateEmployee(repo EmployeeRepository) // ❌ Saltarse casos de uso
```

### 🔵 Infraestructura - Qué SÍ y NO hacer

#### ✅ Buenas Prácticas:
```go
// ✅ Implementación específica de persistencia
func (r *PostgresEmployeeRepo) Save(emp *domain.Employee) error {
    query := "INSERT INTO employees (name, email) VALUES ($1, $2)"
    _, err := r.db.Exec(query, emp.Name, emp.Email)
    return err
}

// ✅ Implementar interfaces del dominio
func NewPostgresEmployeeRepo(db *sql.DB) domain.EmployeeRepository {
    return &PostgresEmployeeRepo{db: db}
}
```

#### ❌ Malas Prácticas:
```go
// ❌ NUNCA: Exponer tipos específicos de DB
func (r *PostgresEmployeeRepo) GetDB() *sql.DB // ❌ Detalle técnico expuesto

// ❌ NUNCA: Lógica de negocio en repositorios
func (r *PostgresEmployeeRepo) ValidateAndSave(emp *domain.Employee) error {
    if emp.Salary < 0 { // ❌ Validación de negocio aquí
        return errors.New("invalid salary")
    }
}
```

## 🔧 Comandos Avanzados

### Generar Feature Completo
```bash
# Genera todas las capas para un feature
goca feature <name> --fields "field:type,..." --database postgres --handlers "http,grpc,cli"
```

### Generar Solo Interfaces
```bash
# Útil para TDD - generar contratos primero
goca interfaces Product --usecase --repository
```

### Generar Inyección de Dependencias
```bash
# Genera contenedor DI para wiring automático
goca di --features "Product,User,Order"
```

## 🎯 Ventajas de Cada Capa

### 🟡 Dominio
- **Entidades puras** sin dependencias externas
- **Reglas de negocio** centralizadas y testables
- **Validaciones** específicas del dominio

### 🔴 Casos de Uso
- **DTOs específicos** para cada operación
- **Lógica de aplicación** bien definida
- **Interfaces claras** hacia otras capas

### 🟢 Adaptadores
- **Separación total** entre protocolos de entrada
- **Validaciones de entrada** específicas por protocolo
- **Transformación** de DTOs externos a internos

### 🔵 Infraestructura
- **Implementaciones intercambiables** de persistencia
- **Aislamiento** de detalles técnicos
- **Configuración** centralizada de recursos externos

## 🔄 Flujo de Dependencias

```
Handler → UseCase → Repository → Database
   ↓         ↓         ↓
 DTO ←→ Business ←→ Domain Entity
```

**Regla de Oro**: Las dependencias siempre apuntan hacia adentro, hacia el dominio.

## 🤝 Contribuir

1. Fork el repositorio
2. Crea tu rama de feature (`git checkout -b feature/clean-arch-enhancement`)
3. Commit tus cambios (`git commit -m 'Add enhanced clean architecture layer'`)
4. Push a la rama (`git push origin feature/clean-arch-enhancement`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## 🆘 Soporte

- 📧 Email: support@goca.dev
- 🐛 Issues: [GitHub Issues](https://github.com/sazardev/goca/issues)
- 📖 Documentación: [Documentación Completa](https://docs.goca.dev)

---

Hecho