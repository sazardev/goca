# goca handler

Generate input adapters for different protocols (HTTP, gRPC, CLI, etc.).

## Syntax

```bash
goca handler <EntityName> [flags]
```

## Description

Creates handlers that adapt external requests to use case calls, handling protocol-specific concerns while keeping business logic isolated.

## Flags

### `--type`

Handler type. Default: `http`

**Options:** `http` | `grpc` | `cli` | `worker` | `soap`

```bash
goca handler Product --type http
```

### `--middleware`

Include middleware setup.

```bash
goca handler User --type http --middleware
```

### `--validation`

Add request validation.

```bash
goca handler Order --type http --validation
```

## Examples

### HTTP REST Handler

```bash
goca handler User --type http
```

**Generates:** `internal/handler/http/user_handler.go`

```go
package http

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "github.com/gorilla/mux"
    "myproject/internal/usecase"
)

type UserHandler struct {
    userService usecase.UserService
}

func NewUserHandler(userService usecase.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/users", h.CreateUser).Methods(http.MethodPost)
    r.HandleFunc("/users/{id}", h.GetUser).Methods(http.MethodGet)
    r.HandleFunc("/users/{id}", h.UpdateUser).Methods(http.MethodPut)
    r.HandleFunc("/users/{id}", h.DeleteUser).Methods(http.MethodDelete)
    r.HandleFunc("/users", h.ListUsers).Methods(http.MethodGet)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req usecase.CreateUserInput
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    user, err := h.userService.CreateUser(r.Context(), req)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    respondJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        respondError(w, http.StatusBadRequest, "Invalid user ID")
        return
    }
    
    user, err := h.userService.GetUser(r.Context(), uint(id))
    if err != nil {
        respondError(w, http.StatusNotFound, "User not found")
        return
    }
    
    respondJSON(w, http.StatusOK, user)
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
    respondJSON(w, status, map[string]string{"error": message})
}
```

### gRPC Handler

```bash
goca handler Product --type grpc
```

**Generates:** `internal/handler/grpc/product_handler.go` + `.proto` file

### CLI Handler

```bash
goca handler Task --type cli
```

**Generates:** CLI commands with Cobra

### Worker Handler

```bash
goca handler Email --type worker
```

**Generates:** Background job handlers

## Handler Types Comparison

| Type       | Use Case                        | Generated                  |
| ---------- | ------------------------------- | -------------------------- |
| **http**   | REST APIs, Web services         | HTTP handlers with routing |
| **grpc**   | Microservices, High performance | gRPC server + proto files  |
| **cli**    | Command-line tools              | Cobra commands             |
| **worker** | Background jobs, Async tasks    | Job handlers               |
| **soap**   | Legacy systems integration      | SOAP client wrappers       |

## Best Practices

### ✅ DO

- Handle protocol-specific concerns only
- Transform requests to use case DTOs
- Format responses appropriately
- Add proper error handling
- Use middleware for cross-cutting concerns

### ❌ DON'T

- Include business logic
- Access repositories directly
- Skip input validation
- Return domain entities directly

## See Also

- [`goca usecase`](/commands/usecase) - Generate use cases
- [`goca feature`](/commands/feature) - Generate complete feature
- [Handler Layer](/guide/clean-architecture#layer-4-handlers-interface-adapters)
