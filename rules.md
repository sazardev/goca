# Guía Definitiva de Clean Architecture en Go

## 📌 Introducción

Este documento tiene como propósito ofrecer una referencia completa y exhaustiva al equipo técnico para la implementación precisa y efectiva de Clean Architecture en proyectos desarrollados en Go. El estándar aquí descrito está diseñado para ser accesible y útil tanto para desarrolladores novatos como experimentados.

---

## 🔖 Principios Fundamentales de Clean Architecture

### ✅ Buenas Prácticas:

* Orientar dependencias exclusivamente hacia el núcleo del sistema.
* Establecer interfaces y contratos claros y explícitos.
* Encapsular la lógica de negocio exclusivamente en capas internas.
* Segregar claramente responsabilidades en diferentes capas.
* Utilizar la inyección de dependencias para optimizar la testabilidad y modularidad.

### 🚫 Malas Prácticas:

* Combinar lógica técnica y lógica de negocio en una misma capa.
* Permitir que entidades tengan dependencia directa de infraestructura externa.
* Mantener paquetes excesivamente grandes o mal estructurados.

---

## 📂 Estructura Completa del Proyecto (Feature: Employee)

```
employee/
├── domain/
│   ├── employee.go
│   ├── errors.go
│   └── validations.go
├── usecase/
│   ├── dto.go
│   ├── employee_usecase.go
│   ├── employee_service.go
│   └── interfaces.go
├── repository/
│   ├── interfaces.go
│   ├── postgres_employee_repo.go
│   └── memory_employee_repo.go
├── handler/
│   ├── http/
│   │   ├── dto.go
│   │   └── handler.go
│   ├── grpc/
│   │   ├── employee.proto
│   │   └── server.go
│   ├── cli/
│   │   └── commands.go
│   ├── worker/
│   │   └── worker.go
│   └── soap/
│       └── soap_client.go
├── messages/
│   ├── errors.go
│   └── responses.go
├── constants/
│   └── constants.go
└── main.go
```

---

## 🧩 Capas, Reglas y Ejemplos

### 🟡 Capa de Dominio (Entities)

* **SÍ:**

  * Mantener entidades puras.
  * Validar invariantes y reglas básicas esenciales.
* **NO:**

  * Incluir dependencias externas o técnicas.

### Ejemplo Avanzado:

```go
package domain

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
```

### 🔴 Capa de Casos de Uso (Application)

* **SÍ:**

  * Definir casos de uso únicos y específicos.
  * Validaciones concretas de negocio.
  * DTOs bien estructurados.
* **NO:**

  * Incorporar lógica técnica o de persistencia directa.

### Ejemplo Completo:

```go
package usecase

import "employee/domain"

type CreateEmployeeInput struct {
	Name  string
	Email string
	Role  string
}

type CreateEmployeeOutput struct {
	Employee domain.Employee
	Message  string
}

type EmployeeUseCase interface {
	CreateEmployee(input CreateEmployeeInput) (CreateEmployeeOutput, error)
}

func (svc *EmployeeService) CreateEmployee(input CreateEmployeeInput) (CreateEmployeeOutput, error) {
	emp := domain.Employee{Name: input.Name, Email: input.Email, Role: input.Role}
	if err := emp.Validate(); err != nil {
		return CreateEmployeeOutput{}, err
	}
	err := svc.repo.Save(&emp)
	if err != nil {
		return CreateEmployeeOutput{}, err
	}
	return CreateEmployeeOutput{Employee: emp, Message: SuccessEmployeeCreated}, nil
}
```

### 🟢 Capa de Adaptadores (Interface Adapters)

* **SÍ:**

  * Adaptar DTOs externos a internos.
  * Mantener manejadores específicos para cada tipo de interfaz.
* **NO:**

  * Añadir lógica de negocio o acoplar directamente persistencia.

### Ejemplo HTTP Completo:

```go
package http

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateEmployeeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.usecase.CreateEmployee(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(output)
}
```

### 🔵 Capa de Repositorio e Infraestructura

* **SÍ:**

  * Implementar interfaces definidas.
  * Gestionar detalles específicos del almacenamiento.
* **NO:**

  * Exponer tipos específicos de DB en otras capas.

### Ejemplo Repositorio Avanzado:

```go
package repository

func (r *PostgresRepo) Save(emp *domain.Employee) error {
	query := "INSERT INTO employees (name, email, role) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRow(query, emp.Name, emp.Email, emp.Role).Scan(&emp.ID)
	if err != nil {
		return ErrDatabaseInsertion
	}
	return nil
}
```

---

## 🧠 Patrones y Anti-Patrones

### Patrones:

* **Repository Pattern**
* **CQRS**
* **Dependency Injection**

### Anti-Patrones:

* **Fat Controller**
* **God Object**
* **Anemic Domain Model**

---

## 🛠️ Consejos Esenciales

* Implementa pruebas unitarias para cada capa.
* Asegura interfaces concisas y bien documentadas.
* Maneja explícitamente errores claros y específicos.
* Mantiene un sistema robusto de mensajes y constantes.

---

## 📚 Recursos Externos

* [Clean Architecture en Go](https://github.com/bxcodec/go-clean-arch)
* [Clean Code by Robert C. Martin](https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882)
* [Golang Project Layout](https://github.com/golang-standards/project-layout)

---

Esta guía representa un estándar completo y riguroso para el desarrollo de aplicaciones en Go, alineado con los principios fundamentales de Clean Architecture.
