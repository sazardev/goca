# Definitive Guide to Clean Architecture in Go

## 📌 Introduction

This document aims to provide a complete and comprehensive reference for the technical team for precise and effective implementation of Clean Architecture in Go projects. The standard described here is designed to be accessible and useful for both novice and experienced developers.

---

## 🔖 Fundamental Principles of Clean Architecture

### ✅ Best Practices:

* Orient dependencies exclusively towards the system core.
* Establish clear and explicit interfaces and contracts.
* Encapsulate business logic exclusively in internal layers.
* Clearly segregate responsibilities in different layers.
* Use dependency injection to optimize testability and modularity.

### 🚫 Bad Practices:

* Combine technical logic and business logic in the same layer.
* Allow entities to have direct dependency on external infrastructure.
* Maintain excessively large or poorly structured packages.

---

## 📂 Complete Project Structure (Feature: Employee)

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

## 🧩 Layers, Rules and Examples

### 🟡 Domain Layer (Entities)

* **DO:**

  * Keep entities pure.
  * Validate invariants and essential basic rules.
* **DON'T:**

  * Include external or technical dependencies.

### Advanced Example:

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

### 🔴 Use Case Layer (Application)

* **DO:**

  * Define unique and specific use cases.
  * Concrete business validations.
  * Well-structured DTOs.
* **DON'T:**

  * Incorporate technical logic or direct persistence.

### Complete Example:

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

### 🟢 Adapter Layer (Interface Adapters)

* **DO:**

  * Adapt external DTOs to internal ones.
  * Maintain specific handlers for each interface type.
* **DON'T:**

  * Add business logic or directly couple persistence.

### Complete HTTP Example:

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

### 🔵 Repository and Infrastructure Layer

* **DO:**

  * Implement defined interfaces.
  * Manage storage-specific details.
* **DON'T:**

  * Expose specific DB types in other layers.

### Advanced Repository Example:

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

## 🧠 Patterns and Anti-Patterns

### Patterns:

* **Repository Pattern**
* **CQRS**
* **Dependency Injection**

### Anti-Patterns:

* **Fat Controller**
* **God Object**
* **Anemic Domain Model**

---

## 🛠️ Essential Tips

* Implement unit tests for each layer.
* Ensure concise and well-documented interfaces.
* Handle errors explicitly with clear and specific messages.
* Maintain a robust system of messages and constants.

---

## 📚 External Resources

* [Clean Architecture in Go](https://github.com/bxcodec/go-clean-arch)
* [Clean Code by Robert C. Martin](https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882)
* [Golang Project Layout](https://github.com/golang-standards/project-layout)

---

This guide represents a complete and rigorous standard for developing applications in Go, aligned with the fundamental principles of Clean Architecture.
