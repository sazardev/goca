# Comando goca handler

El comando `goca handler` crea adaptadores de entrega que manejan diferentes protocolos (HTTP, gRPC, CLI, Worker, SOAP) manteniendo la separación de capas y siguiendo Clean Architecture.

## 📋 Sintaxis

```bash
goca handler <entity> [flags]
```

## 🎯 Propósito

Crea handlers (adaptadores de entrada) para diferentes protocolos:

- 🟢 **HTTP REST** con Gin y documentación Swagger
- 🔷 **gRPC** con Protocol Buffers y servidores
- 💻 **CLI** para herramientas de línea de comandos
- ⚙️ **Worker** para tareas en background
- 🌐 **SOAP** para servicios web legacy
- 🛡️ **Middleware** y validaciones por protocolo

## 🚩 Flags Disponibles

| Flag           | Tipo     | Requerido | Valor por Defecto | Descripción                                               |
| -------------- | -------- | --------- | ----------------- | --------------------------------------------------------- |
| `--type`       | `string` | ❌ No      | `http`            | Tipo de handler (`http`, `grpc`, `cli`, `worker`, `soap`) |
| `--swagger`    | `bool`   | ❌ No      | `false`           | Generar documentación Swagger (solo HTTP)                 |
| `--middleware` | `bool`   | ❌ No      | `false`           | Incluir setup de middleware                               |
| `--validation` | `bool`   | ❌ No      | `false`           | Validación de entrada en handler                          |

## 📖 Ejemplos de Uso

### Handler HTTP REST
```bash
goca handler User --type http --swagger --middleware --validation
```

### Handler gRPC
```bash
goca handler Product --type grpc
```

### Handler CLI
```bash
goca handler Order --type cli
```

### Handler Worker
```bash
goca handler Notification --type worker
```

### Handler SOAP
```bash
goca handler Payment --type soap
```

## 📂 Archivos Generados por Tipo

### HTTP REST (`--type http`)
```
internal/handler/http/
├── user_handler.go     # Controladores HTTP
├── user_routes.go      # Definición de rutas
├── dto.go              # DTOs específicos para HTTP
└── swagger.yaml        # Documentación Swagger (si --swagger)
```

### gRPC (`--type grpc`)
```
internal/handler/grpc/
├── user.proto          # Definición Protocol Buffers
└── user_server.go      # Servidor gRPC
```

### CLI (`--type cli`)
```
internal/handler/cli/
└── user_commands.go    # Comandos CLI con Cobra
```

### Worker (`--type worker`)
```
internal/handler/worker/
└── user_worker.go      # Worker para tareas en background
```

### SOAP (`--type soap`)
```
internal/handler/soap/
└── user_client.go      # Cliente SOAP
```

## 🔍 Código Generado en Detalle

### HTTP Handler: `internal/handler/http/user_handler.go`

```go
package http

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/usuario/proyecto/internal/usecase"
    "github.com/usuario/proyecto/internal/usecase/dto"
)

// UserHandler maneja las peticiones HTTP para usuarios
type UserHandler struct {
    userUseCase usecase.UserUseCase
}

// NewUserHandler crea una nueva instancia del handler
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
    return &UserHandler{
        userUseCase: userUseCase,
    }
}

// @Summary Crear usuario
// @Description Crea un nuevo usuario en el sistema
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "Datos del usuario"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "Invalid request body",
            Message: err.Error(),
        })
        return
    }
    
    userResponse, err := h.userUseCase.Create(c.Request.Context(), req)
    if err != nil {
        statusCode := http.StatusInternalServerError
        if err == domain.ErrUserEmailAlreadyExists {
            statusCode = http.StatusConflict
        }
        
        c.JSON(statusCode, dto.ErrorResponse{
            Error:   "Failed to create user",
            Message: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusCreated, userResponse)
}

// @Summary Obtener usuario por ID
// @Description Obtiene un usuario específico por su ID
// @Tags users
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "Invalid user ID",
            Message: "User ID must be a valid number",
        })
        return
    }
    
    userResponse, err := h.userUseCase.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        statusCode := http.StatusInternalServerError
        if err == domain.ErrUserNotFound {
            statusCode = http.StatusNotFound
        }
        
        c.JSON(statusCode, dto.ErrorResponse{
            Error:   "Failed to get user",
            Message: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, userResponse)
}

// @Summary Listar usuarios
// @Description Obtiene una lista paginada de usuarios
// @Tags users
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Success 200 {object} dto.ListUsersResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
    page := 1
    limit := 10
    
    if pageStr := c.Query("page"); pageStr != "" {
        if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
            page = p
        }
    }
    
    if limitStr := c.Query("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
            limit = l
        }
    }
    
    req := dto.ListUsersRequest{
        Page:  page,
        Limit: limit,
    }
    
    usersResponse, err := h.userUseCase.List(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Error:   "Failed to list users",
            Message: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, usersResponse)
}
```

### HTTP Routes: `internal/handler/http/user_routes.go`

```go
package http

import (
    "github.com/gin-gonic/gin"
    "github.com/usuario/proyecto/internal/usecase"
)

// RegisterUserRoutes registra todas las rutas relacionadas con usuarios
func RegisterUserRoutes(router *gin.RouterGroup, userUseCase usecase.UserUseCase) {
    userHandler := NewUserHandler(userUseCase)
    
    users := router.Group("/users")
    {
        users.POST("", userHandler.Create)       // POST /api/v1/users
        users.GET("", userHandler.List)          // GET /api/v1/users
        users.GET("/:id", userHandler.GetByID)   // GET /api/v1/users/:id
        users.PUT("/:id", userHandler.Update)    // PUT /api/v1/users/:id
        users.DELETE("/:id", userHandler.Delete) // DELETE /api/v1/users/:id
    }
}

// RegisterUserRoutesWithMiddleware registra rutas con middleware específico
func RegisterUserRoutesWithMiddleware(router *gin.RouterGroup, userUseCase usecase.UserUseCase) {
    userHandler := NewUserHandler(userUseCase)
    
    users := router.Group("/users")
    
    // Middleware específico para usuarios
    users.Use(AuthMiddleware())
    users.Use(RateLimitMiddleware())
    users.Use(ValidationMiddleware())
    
    {
        users.POST("", userHandler.Create)
        users.GET("", userHandler.List)
        users.GET("/:id", userHandler.GetByID)
        users.PUT("/:id", userHandler.Update)
        users.DELETE("/:id", userHandler.Delete)
    }
}
```

### gRPC Server: `internal/handler/grpc/user_server.go`

```go
package grpc

import (
    "context"
    
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "github.com/usuario/proyecto/internal/usecase"
    "github.com/usuario/proyecto/internal/usecase/dto"
)

// UserServer implementa el servidor gRPC para usuarios
type UserServer struct {
    UnimplementedUserServiceServer
    userUseCase usecase.UserUseCase
}

// NewUserServer crea una nueva instancia del servidor gRPC
func NewUserServer(userUseCase usecase.UserUseCase) *UserServer {
    return &UserServer{
        userUseCase: userUseCase,
    }
}

// CreateUser crea un nuevo usuario via gRPC
func (s *UserServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
    createReq := dto.CreateUserRequest{
        Name:  req.GetName(),
        Email: req.GetEmail(),
    }
    
    userResponse, err := s.userUseCase.Create(ctx, createReq)
    if err != nil {
        if err == domain.ErrUserEmailAlreadyExists {
            return nil, status.Error(codes.AlreadyExists, "Email already exists")
        }
        return nil, status.Error(codes.Internal, err.Error())
    }
    
    return &UserResponse{
        Id:        uint32(userResponse.ID),
        Name:      userResponse.Name,
        Email:     userResponse.Email,
        CreatedAt: userResponse.CreatedAt.Unix(),
        UpdatedAt: userResponse.UpdatedAt.Unix(),
    }, nil
}

// GetUser obtiene un usuario por ID via gRPC
func (s *UserServer) GetUser(ctx context.Context, req *GetUserRequest) (*UserResponse, error) {
    userResponse, err := s.userUseCase.GetByID(ctx, uint(req.GetId()))
    if err != nil {
        if err == domain.ErrUserNotFound {
            return nil, status.Error(codes.NotFound, "User not found")
        }
        return nil, status.Error(codes.Internal, err.Error())
    }
    
    return &UserResponse{
        Id:        uint32(userResponse.ID),
        Name:      userResponse.Name,
        Email:     userResponse.Email,
        CreatedAt: userResponse.CreatedAt.Unix(),
        UpdatedAt: userResponse.UpdatedAt.Unix(),
    }, nil
}

// ListUsers lista usuarios via gRPC
func (s *UserServer) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
    listReq := dto.ListUsersRequest{
        Page:  int(req.GetPage()),
        Limit: int(req.GetLimit()),
    }
    
    usersResponse, err := s.userUseCase.List(ctx, listReq)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    
    users := make([]*UserResponse, len(usersResponse.Users))
    for i, user := range usersResponse.Users {
        users[i] = &UserResponse{
            Id:        uint32(user.ID),
            Name:      user.Name,
            Email:     user.Email,
            CreatedAt: user.CreatedAt.Unix(),
            UpdatedAt: user.UpdatedAt.Unix(),
        }
    }
    
    return &ListUsersResponse{
        Users:       users,
        Total:       usersResponse.Total,
        Page:        int32(usersResponse.Page),
        Limit:       int32(usersResponse.Limit),
        TotalPages:  int32(usersResponse.TotalPages),
        HasNextPage: usersResponse.HasNextPage,
        HasPrevPage: usersResponse.HasPrevPage,
    }, nil
}
```

### gRPC Proto: `internal/handler/grpc/user.proto`

```protobuf
syntax = "proto3";

package user;

option go_package = "./;grpc";

// Servicio de usuarios
service UserService {
    rpc CreateUser(CreateUserRequest) returns (UserResponse);
    rpc GetUser(GetUserRequest) returns (UserResponse);
    rpc UpdateUser(UpdateUserRequest) returns (UserResponse);
    rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}

// Request para crear usuario
message CreateUserRequest {
    string name = 1;
    string email = 2;
}

// Request para obtener usuario
message GetUserRequest {
    uint32 id = 1;
}

// Request para actualizar usuario
message UpdateUserRequest {
    uint32 id = 1;
    optional string name = 2;
    optional string email = 3;
}

// Request para eliminar usuario
message DeleteUserRequest {
    uint32 id = 1;
}

// Response para eliminar usuario
message DeleteUserResponse {
    bool success = 1;
    string message = 2;
}

// Request para listar usuarios
message ListUsersRequest {
    int32 page = 1;
    int32 limit = 2;
}

// Response de usuario
message UserResponse {
    uint32 id = 1;
    string name = 2;
    string email = 3;
    int64 created_at = 4;
    int64 updated_at = 5;
}

// Response para listar usuarios
message ListUsersResponse {
    repeated UserResponse users = 1;
    int64 total = 2;
    int32 page = 3;
    int32 limit = 4;
    int32 total_pages = 5;
    bool has_next_page = 6;
    bool has_prev_page = 7;
}
```

### CLI Commands: `internal/handler/cli/user_commands.go`

```go
package cli

import (
    "context"
    "fmt"
    "os"
    "strconv"
    "time"
    
    "github.com/spf13/cobra"
    "github.com/usuario/proyecto/internal/usecase"
    "github.com/usuario/proyecto/internal/usecase/dto"
)

// UserCommands contiene todos los comandos CLI para usuarios
type UserCommands struct {
    userUseCase usecase.UserUseCase
}

// NewUserCommands crea una nueva instancia de comandos de usuario
func NewUserCommands(userUseCase usecase.UserUseCase) *UserCommands {
    return &UserCommands{
        userUseCase: userUseCase,
    }
}

// GetCommands retorna todos los comandos de usuario
func (uc *UserCommands) GetCommands() []*cobra.Command {
    return []*cobra.Command{
        uc.createUserCmd(),
        uc.getUserCmd(),
        uc.listUsersCmd(),
        uc.updateUserCmd(),
        uc.deleteUserCmd(),
    }
}

// createUserCmd comando para crear usuario
func (uc *UserCommands) createUserCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "create-user",
        Short: "Crear un nuevo usuario",
        Long:  "Crea un nuevo usuario en el sistema con nombre y email",
        RunE: func(cmd *cobra.Command, args []string) error {
            name, _ := cmd.Flags().GetString("name")
            email, _ := cmd.Flags().GetString("email")
            
            if name == "" || email == "" {
                return fmt.Errorf("name and email are required")
            }
            
            req := dto.CreateUserRequest{
                Name:  name,
                Email: email,
            }
            
            ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
            defer cancel()
            
            user, err := uc.userUseCase.Create(ctx, req)
            if err != nil {
                return fmt.Errorf("failed to create user: %w", err)
            }
            
            fmt.Printf("User created successfully:\n")
            fmt.Printf("ID: %d\n", user.ID)
            fmt.Printf("Name: %s\n", user.Name)
            fmt.Printf("Email: %s\n", user.Email)
            fmt.Printf("Created: %s\n", user.CreatedAt.Format(time.RFC3339))
            
            return nil
        },
    }
    
    cmd.Flags().StringP("name", "n", "", "Nombre del usuario")
    cmd.Flags().StringP("email", "e", "", "Email del usuario")
    cmd.MarkFlagRequired("name")
    cmd.MarkFlagRequired("email")
    
    return cmd
}

// getUserCmd comando para obtener usuario por ID
func (uc *UserCommands) getUserCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "get-user",
        Short: "Obtener usuario por ID",
        Long:  "Obtiene la información de un usuario específico por su ID",
        RunE: func(cmd *cobra.Command, args []string) error {
            idStr, _ := cmd.Flags().GetString("id")
            id, err := strconv.ParseUint(idStr, 10, 32)
            if err != nil {
                return fmt.Errorf("invalid user ID: %w", err)
            }
            
            ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
            defer cancel()
            
            user, err := uc.userUseCase.GetByID(ctx, uint(id))
            if err != nil {
                return fmt.Errorf("failed to get user: %w", err)
            }
            
            fmt.Printf("User information:\n")
            fmt.Printf("ID: %d\n", user.ID)
            fmt.Printf("Name: %s\n", user.Name)
            fmt.Printf("Email: %s\n", user.Email)
            fmt.Printf("Created: %s\n", user.CreatedAt.Format(time.RFC3339))
            fmt.Printf("Updated: %s\n", user.UpdatedAt.Format(time.RFC3339))
            
            return nil
        },
    }
    
    cmd.Flags().StringP("id", "i", "", "ID del usuario")
    cmd.MarkFlagRequired("id")
    
    return cmd
}

// listUsersCmd comando para listar usuarios
func (uc *UserCommands) listUsersCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "list-users",
        Short: "Listar usuarios",
        Long:  "Lista todos los usuarios del sistema con paginación",
        RunE: func(cmd *cobra.Command, args []string) error {
            page, _ := cmd.Flags().GetInt("page")
            limit, _ := cmd.Flags().GetInt("limit")
            
            if page < 1 {
                page = 1
            }
            if limit < 1 || limit > 100 {
                limit = 10
            }
            
            req := dto.ListUsersRequest{
                Page:  page,
                Limit: limit,
            }
            
            ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
            defer cancel()
            
            response, err := uc.userUseCase.List(ctx, req)
            if err != nil {
                return fmt.Errorf("failed to list users: %w", err)
            }
            
            fmt.Printf("Users (Page %d of %d, Total: %d):\n\n", 
                response.Page, response.TotalPages, response.Total)
            
            for _, user := range response.Users {
                fmt.Printf("ID: %d | Name: %s | Email: %s | Created: %s\n",
                    user.ID, user.Name, user.Email, 
                    user.CreatedAt.Format("2006-01-02 15:04:05"))
            }
            
            fmt.Printf("\n")
            if response.HasPrevPage {
                fmt.Printf("Previous: --page %d\n", response.Page-1)
            }
            if response.HasNextPage {
                fmt.Printf("Next: --page %d\n", response.Page+1)
            }
            
            return nil
        },
    }
    
    cmd.Flags().IntP("page", "p", 1, "Número de página")
    cmd.Flags().IntP("limit", "l", 10, "Elementos por página")
    
    return cmd
}
```

### Worker: `internal/handler/worker/user_worker.go`

```go
package worker

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"
    
    "github.com/usuario/proyecto/internal/usecase"
    "github.com/usuario/proyecto/internal/usecase/dto"
)

// UserWorker maneja tareas en background relacionadas con usuarios
type UserWorker struct {
    userUseCase usecase.UserUseCase
}

// NewUserWorker crea una nueva instancia del worker
func NewUserWorker(userUseCase usecase.UserUseCase) *UserWorker {
    return &UserWorker{
        userUseCase: userUseCase,
    }
}

// UserTask representa una tarea de usuario
type UserTask struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
    UserID  uint        `json:"user_id,omitempty"`
}

// ProcessUserTask procesa una tarea de usuario
func (w *UserWorker) ProcessUserTask(ctx context.Context, taskData []byte) error {
    var task UserTask
    if err := json.Unmarshal(taskData, &task); err != nil {
        return fmt.Errorf("failed to unmarshal task: %w", err)
    }
    
    switch task.Type {
    case "create_user":
        return w.processCreateUser(ctx, task.Payload)
    case "send_welcome_email":
        return w.processSendWelcomeEmail(ctx, task.UserID)
    case "update_user_stats":
        return w.processUpdateUserStats(ctx, task.UserID)
    case "cleanup_inactive_users":
        return w.processCleanupInactiveUsers(ctx)
    default:
        return fmt.Errorf("unknown task type: %s", task.Type)
    }
}

// processCreateUser procesa la creación de un usuario
func (w *UserWorker) processCreateUser(ctx context.Context, payload interface{}) error {
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }
    
    var req dto.CreateUserRequest
    if err := json.Unmarshal(payloadBytes, &req); err != nil {
        return fmt.Errorf("failed to unmarshal create request: %w", err)
    }
    
    user, err := w.userUseCase.Create(ctx, req)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    log.Printf("User created successfully: ID=%d, Email=%s", user.ID, user.Email)
    
    // Programar tarea de email de bienvenida
    welcomeTask := UserTask{
        Type:   "send_welcome_email",
        UserID: user.ID,
    }
    
    if err := w.scheduleTask(welcomeTask, 5*time.Second); err != nil {
        log.Printf("Failed to schedule welcome email: %v", err)
    }
    
    return nil
}

// processSendWelcomeEmail envía email de bienvenida
func (w *UserWorker) processSendWelcomeEmail(ctx context.Context, userID uint) error {
    user, err := w.userUseCase.GetByID(ctx, userID)
    if err != nil {
        return fmt.Errorf("failed to get user: %w", err)
    }
    
    // Simular envío de email
    log.Printf("Sending welcome email to %s (%s)", user.Name, user.Email)
    time.Sleep(2 * time.Second) // Simular latencia de email
    log.Printf("Welcome email sent successfully to %s", user.Email)
    
    return nil
}

// processUpdateUserStats actualiza estadísticas del usuario
func (w *UserWorker) processUpdateUserStats(ctx context.Context, userID uint) error {
    log.Printf("Updating stats for user ID: %d", userID)
    
    // Aquí implementarías la lógica de actualización de estadísticas
    // Por ejemplo: calcular número de pedidos, última actividad, etc.
    
    time.Sleep(1 * time.Second) // Simular procesamiento
    log.Printf("Stats updated for user ID: %d", userID)
    
    return nil
}

// processCleanupInactiveUsers limpia usuarios inactivos
func (w *UserWorker) processCleanupInactiveUsers(ctx context.Context) error {
    log.Println("Starting cleanup of inactive users")
    
    // Aquí implementarías la lógica de limpieza
    // Por ejemplo: marcar como inactivos usuarios sin actividad en X días
    
    time.Sleep(5 * time.Second) // Simular procesamiento pesado
    log.Println("Inactive users cleanup completed")
    
    return nil
}

// scheduleTask programa una tarea para ejecutar después de un delay
func (w *UserWorker) scheduleTask(task UserTask, delay time.Duration) error {
    go func() {
        time.Sleep(delay)
        
        taskData, err := json.Marshal(task)
        if err != nil {
            log.Printf("Failed to marshal scheduled task: %v", err)
            return
        }
        
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        if err := w.ProcessUserTask(ctx, taskData); err != nil {
            log.Printf("Failed to process scheduled task: %v", err)
        }
    }()
    
    return nil
}

// StartWorker inicia el worker para procesar tareas
func (w *UserWorker) StartWorker(ctx context.Context) error {
    log.Println("Starting user worker...")
    
    // Aquí conectarías con tu sistema de colas (Redis, RabbitMQ, etc.)
    // Por ahora, simularemos con un ticker
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            log.Println("User worker stopping...")
            return ctx.Err()
        case <-ticker.C:
            // Procesar tareas pendientes
            log.Println("Checking for pending user tasks...")
            
            // Ejemplo de tarea programada
            cleanupTask := UserTask{
                Type: "cleanup_inactive_users",
            }
            
            taskData, _ := json.Marshal(cleanupTask)
            if err := w.ProcessUserTask(ctx, taskData); err != nil {
                log.Printf("Error processing cleanup task: %v", err)
            }
        }
    }
}
```

### SOAP Client: `internal/handler/soap/user_client.go`

```go
package soap

import (
    "bytes"
    "context"
    "encoding/xml"
    "fmt"
    "io"
    "net/http"
    "time"
    
    "github.com/usuario/proyecto/internal/usecase"
    "github.com/usuario/proyecto/internal/usecase/dto"
)

// UserSOAPClient cliente SOAP para usuarios
type UserSOAPClient struct {
    endpoint    string
    httpClient  *http.Client
    userUseCase usecase.UserUseCase
}

// NewUserSOAPClient crea un nuevo cliente SOAP
func NewUserSOAPClient(endpoint string, userUseCase usecase.UserUseCase) *UserSOAPClient {
    return &UserSOAPClient{
        endpoint: endpoint,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        userUseCase: userUseCase,
    }
}

// SOAPEnvelope estructura del sobre SOAP
type SOAPEnvelope struct {
    XMLName xml.Name    `xml:"soap:Envelope"`
    XMLNS   string      `xml:"xmlns:soap,attr"`
    Body    interface{} `xml:"soap:Body"`
}

// SOAPFault estructura de errores SOAP
type SOAPFault struct {
    XMLName xml.Name `xml:"soap:Fault"`
    Code    string   `xml:"faultcode"`
    String  string   `xml:"faultstring"`
    Detail  string   `xml:"detail,omitempty"`
}

// CreateUserRequest request SOAP para crear usuario
type CreateUserRequest struct {
    XMLName xml.Name `xml:"CreateUserRequest"`
    Name    string   `xml:"name"`
    Email   string   `xml:"email"`
}

// CreateUserResponse response SOAP para crear usuario
type CreateUserResponse struct {
    XMLName   xml.Name `xml:"CreateUserResponse"`
    ID        uint     `xml:"id"`
    Name      string   `xml:"name"`
    Email     string   `xml:"email"`
    CreatedAt string   `xml:"createdAt"`
}

// GetUserRequest request SOAP para obtener usuario
type GetUserRequest struct {
    XMLName xml.Name `xml:"GetUserRequest"`
    ID      uint     `xml:"id"`
}

// GetUserResponse response SOAP para obtener usuario
type GetUserResponse struct {
    XMLName   xml.Name `xml:"GetUserResponse"`
    ID        uint     `xml:"id"`
    Name      string   `xml:"name"`
    Email     string   `xml:"email"`
    CreatedAt string   `xml:"createdAt"`
    UpdatedAt string   `xml:"updatedAt"`
}

// CreateUser crea un usuario via SOAP
func (c *UserSOAPClient) CreateUser(ctx context.Context, name, email string) (*CreateUserResponse, error) {
    // Crear caso de uso request
    req := dto.CreateUserRequest{
        Name:  name,
        Email: email,
    }
    
    // Procesar con caso de uso
    userResponse, err := c.userUseCase.Create(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    // Convertir a response SOAP
    soapResponse := &CreateUserResponse{
        ID:        userResponse.ID,
        Name:      userResponse.Name,
        Email:     userResponse.Email,
        CreatedAt: userResponse.CreatedAt.Format(time.RFC3339),
    }
    
    return soapResponse, nil
}

// GetUser obtiene un usuario via SOAP
func (c *UserSOAPClient) GetUser(ctx context.Context, id uint) (*GetUserResponse, error) {
    // Procesar con caso de uso
    userResponse, err := c.userUseCase.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    // Convertir a response SOAP
    soapResponse := &GetUserResponse{
        ID:        userResponse.ID,
        Name:      userResponse.Name,
        Email:     userResponse.Email,
        CreatedAt: userResponse.CreatedAt.Format(time.RFC3339),
        UpdatedAt: userResponse.UpdatedAt.Format(time.RFC3339),
    }
    
    return soapResponse, nil
}

// sendSOAPRequest envía una petición SOAP
func (c *UserSOAPClient) sendSOAPRequest(ctx context.Context, action string, body interface{}) ([]byte, error) {
    envelope := SOAPEnvelope{
        XMLNS: "http://schemas.xmlsoap.org/soap/envelope/",
        Body:  body,
    }
    
    xmlData, err := xml.Marshal(envelope)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal SOAP request: %w", err)
    }
    
    req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(xmlData))
    if err != nil {
        return nil, fmt.Errorf("failed to create HTTP request: %w", err)
    }
    
    req.Header.Set("Content-Type", "text/xml; charset=utf-8")
    req.Header.Set("SOAPAction", action)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send SOAP request: %w", err)
    }
    defer resp.Body.Close()
    
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read SOAP response: %w", err)
    }
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("SOAP request failed with status %d: %s", resp.StatusCode, string(respBody))
    }
    
    return respBody, nil
}

// HandleSOAPRequest maneja peticiones SOAP entrantes
func (c *UserSOAPClient) HandleSOAPRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    
    var envelope SOAPEnvelope
    if err := xml.Unmarshal(body, &envelope); err != nil {
        http.Error(w, "Invalid SOAP envelope", http.StatusBadRequest)
        return
    }
    
    // Procesar según el tipo de petición
    // Aquí determinarías qué operación realizar basándote en el contenido del body
    
    w.Header().Set("Content-Type", "text/xml; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    
    // Retornar respuesta SOAP apropiada
    response := SOAPEnvelope{
        XMLNS: "http://schemas.xmlsoap.org/soap/envelope/",
        Body:  "<!-- Response body would go here -->",
    }
    
    xml.NewEncoder(w).Encode(response)
}
```

## 🛡️ Middleware (--middleware)

Con `--middleware`, se generan middlewares específicos:

```go
// AuthMiddleware middleware de autenticación
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        // Validar token JWT
        // ... lógica de validación
        
        c.Next()
    }
}

// ValidationMiddleware middleware de validación
func ValidationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Validaciones específicas antes de llegar al handler
        c.Next()
    }
}

// RateLimitMiddleware middleware de rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Lógica de rate limiting
        c.Next()
    }
}
```

## 📄 Swagger Documentation (--swagger)

Con `--swagger`, se genera documentación automática:

```yaml
swagger: "2.0"
info:
  title: User API
  description: API para gestión de usuarios
  version: 1.0.0
host: localhost:8080
basePath: /api/v1
schemes:
  - http
  - https

paths:
  /users:
    post:
      summary: Crear usuario
      description: Crea un nuevo usuario en el sistema
      tags:
        - users
      consumes:
        - application/json
      produces:
        - application/json
      parameters:
        - in: body
          name: user
          description: Datos del usuario
          required: true
          schema:
            $ref: '#/definitions/CreateUserRequest'
      responses:
        201:
          description: Usuario creado exitosamente
          schema:
            $ref: '#/definitions/UserResponse'
        400:
          description: Datos inválidos
          schema:
            $ref: '#/definitions/ErrorResponse'

definitions:
  CreateUserRequest:
    type: object
    required:
      - name
      - email
    properties:
      name:
        type: string
        example: "John Doe"
      email:
        type: string
        example: "john@example.com"
        
  UserResponse:
    type: object
    properties:
      id:
        type: integer
        example: 1
      name:
        type: string
        example: "John Doe"
      email:
        type: string
        example: "john@example.com"
```

## ⚠️ Consideraciones Importantes

### ✅ Buenas Prácticas
- **Separación de protocolos**: Cada handler maneja solo su protocolo
- **DTOs específicos**: Diferentes DTOs por protocolo si es necesario
- **Error handling**: Manejo consistente de errores por protocolo
- **Context propagation**: Usar context.Context en todas las operaciones

### ❌ Errores Comunes
- **Lógica de negocio en handlers**: Debe estar en casos de uso
- **Dependencias directas**: No acceder a repositorios directamente
- **Exposición de errores internos**: Mapear errores apropiadamente
- **No validar entrada**: Siempre validar datos de entrada

---

**← [Comando goca repository](Command-Repository) | [Comando goca di](Command-DI) →**
