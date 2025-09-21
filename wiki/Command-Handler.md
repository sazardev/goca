# goca handler Command

The `goca handler` command creates delivery adapters that handle different protocols (HTTP, gRPC, CLI, Worker, SOAP) while maintaining layer separation and following Clean Architecture.

## 📋 Syntax

```bash
goca handler <entity> [flags]
```

## 🎯 Purpose

Creates handlers (input adapters) for different protocols:

- 🟢 **HTTP REST** with Gin and Swagger documentation
- 🔷 **gRPC** with Protocol Buffers and servers
- 💻 **CLI** for command-line tools
- ⚙️ **Worker** for background tasks
- 🌐 **SOAP** for legacy web services
- 🛡️ **Middleware** and protocol-specific validations

## 🚩 Available Flags

| Flag           | Type     | Required | Default Value | Description                                            |
| -------------- | -------- | -------- | ------------- | ------------------------------------------------------ |
| `--type`       | `string` | ❌ No     | `http`        | Handler type (`http`, `grpc`, `cli`, `worker`, `soap`) |
| `--swagger`    | `bool`   | ❌ No     | `false`       | Generate Swagger documentation (HTTP only)             |
| `--middleware` | `bool`   | ❌ No     | `false`       | Include middleware setup                               |
| `--validation` | `bool`   | ❌ No     | `false`       | Input validation in handler                            |

## 📖 Usage Examples

### HTTP REST Handler
```bash
goca handler User --type http --swagger --middleware --validation
```

### gRPC Handler
```bash
goca handler Product --type grpc
```

### CLI Handler
```bash
goca handler Order --type cli
```

### Worker Handler
```bash
goca handler Notification --type worker
```

### SOAP Handler
```bash
goca handler Payment --type soap
```

## 📂 Generated Files by Type

### HTTP REST (`--type http`)
```
internal/handler/http/
├── user_handler.go     # HTTP controllers
├── user_routes.go      # Route definitions
├── dto.go              # HTTP-specific DTOs
└── swagger.yaml        # Swagger documentation (if --swagger)
```

### gRPC (`--type grpc`)
```
internal/handler/grpc/
├── user.proto          # Protocol Buffers definition
└── user_server.go      # gRPC server
```

### CLI (`--type cli`)
```
internal/handler/cli/
└── user_commands.go    # CLI commands with Cobra
```

### Worker (`--type worker`)
```
internal/handler/worker/
└── user_worker.go      # Worker for background tasks
```

### SOAP (`--type soap`)
```
internal/handler/soap/
└── user_client.go      # SOAP client
```

## 🔍 Generated Code in Detail

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

// UserHandler handles HTTP requests for users
type UserHandler struct {
    userUseCase usecase.UserUseCase
}

// NewUserHandler creates a new handler instance
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
    return &UserHandler{
        userUseCase: userUseCase,
    }
}

// @Summary Create user
// @Description Creates a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User data"
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

// @Summary Get user by ID
// @Description Gets a specific user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
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

// @Summary List users
// @Description Gets a paginated list of users
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
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

// RegisterUserRoutes registers all user-related routes
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

// RegisterUserRoutesWithMiddleware registers routes with specific middleware
func RegisterUserRoutesWithMiddleware(router *gin.RouterGroup, userUseCase usecase.UserUseCase) {
    userHandler := NewUserHandler(userUseCase)
    
    users := router.Group("/users")
    
    // User-specific middleware
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

// UserServer implements the gRPC server for users
type UserServer struct {
    UnimplementedUserServiceServer
    userUseCase usecase.UserUseCase
}

// NewUserServer creates a new gRPC server instance
func NewUserServer(userUseCase usecase.UserUseCase) *UserServer {
    return &UserServer{
        userUseCase: userUseCase,
    }
}

// CreateUser creates a new user via gRPC
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

// GetUser gets a user by ID via gRPC
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

// ListUsers lists users via gRPC
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

// createUserCmd command to create user
func (uc *UserCommands) createUserCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "create-user",
        Short: "Create a new user",
        Long:  "Creates a new user in the system with name and email",
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
    
    cmd.Flags().StringP("name", "n", "", "User name")
    cmd.Flags().StringP("email", "e", "", "User email")
    cmd.MarkFlagRequired("name")
    cmd.MarkFlagRequired("email")
    
    return cmd
}

// getUserCmd command to get user by ID
func (uc *UserCommands) getUserCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "get-user",
        Short: "Get user by ID",
        Long:  "Gets information for a specific user by their ID",
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
    
    cmd.Flags().StringP("id", "i", "", "User ID")
    cmd.MarkFlagRequired("id")
    
    return cmd
}

// listUsersCmd command to list users
func (uc *UserCommands) listUsersCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "list-users",
        Short: "List users",
        Long:  "Lists all users in the system with pagination",
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
    
    cmd.Flags().IntP("page", "p", 1, "Page number")
    cmd.Flags().IntP("limit", "l", 10, "Items per page")
    
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

// UserWorker handles background tasks related to users
type UserWorker struct {
    userUseCase usecase.UserUseCase
}

// NewUserWorker creates a new worker instance
func NewUserWorker(userUseCase usecase.UserUseCase) *UserWorker {
    return &UserWorker{
        userUseCase: userUseCase,
    }
}

// UserTask represents a user task
type UserTask struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
    UserID  uint        `json:"user_id,omitempty"`
}

// ProcessUserTask processes a user task
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

// processCreateUser processes user creation
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
    
    // Schedule welcome email task
    welcomeTask := UserTask{
        Type:   "send_welcome_email",
        UserID: user.ID,
    }
    
    if err := w.scheduleTask(welcomeTask, 5*time.Second); err != nil {
        log.Printf("Failed to schedule welcome email: %v", err)
    }
    
    return nil
}

// processSendWelcomeEmail sends welcome email
func (w *UserWorker) processSendWelcomeEmail(ctx context.Context, userID uint) error {
    user, err := w.userUseCase.GetByID(ctx, userID)
    if err != nil {
        return fmt.Errorf("failed to get user: %w", err)
    }
    
    // Simulate email sending
    log.Printf("Sending welcome email to %s (%s)", user.Name, user.Email)
    time.Sleep(2 * time.Second) // Simulate email latency
    log.Printf("Welcome email sent successfully to %s", user.Email)
    
    return nil
}

// processUpdateUserStats updates user statistics
func (w *UserWorker) processUpdateUserStats(ctx context.Context, userID uint) error {
    log.Printf("Updating stats for user ID: %d", userID)
    
    // Here you would implement the statistics update logic
    // For example: calculate number of orders, last activity, etc.
    
    time.Sleep(1 * time.Second) // Simulate processing
    log.Printf("Stats updated for user ID: %d", userID)
    
    return nil
}

// processCleanupInactiveUsers cleans up inactive users
func (w *UserWorker) processCleanupInactiveUsers(ctx context.Context) error {
    log.Println("Starting cleanup of inactive users")
    
    // Here you would implement the cleanup logic
    // For example: mark as inactive users without activity in X days
    
    time.Sleep(5 * time.Second) // Simulate heavy processing
    log.Println("Inactive users cleanup completed")
    
    return nil
}

// scheduleTask schedules a task to run after a delay
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

// StartWorker starts the worker to process tasks
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
    
    // Process according to request type
    // Here you would determine which operation to perform based on the body content
    
    w.Header().Set("Content-Type", "text/xml; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    
    // Return appropriate SOAP response
    response := SOAPEnvelope{
        XMLNS: "http://schemas.xmlsoap.org/soap/envelope/",
        Body:  "<!-- Response body would go here -->",
    }
    
    xml.NewEncoder(w).Encode(response)
}
```

## 🛡️ Middleware (--middleware)

With `--middleware`, specific middlewares are generated:

```go
// AuthMiddleware authentication middleware
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        // Validate JWT token
        // ... validation logic
        
        c.Next()
    }
}

// ValidationMiddleware validation middleware
func ValidationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Specific validations before reaching the handler
        c.Next()
    }
}

// RateLimitMiddleware rate limiting middleware
func RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Rate limiting logic
        c.Next()
    }
}
```

## 📄 Swagger Documentation (--swagger)

With `--swagger`, automatic documentation is generated:

```yaml
swagger: "2.0"
info:
  title: User API
  description: API for user management
  version: 1.0.0
host: localhost:8080
basePath: /api/v1
schemes:
  - http
  - https

paths:
  /users:
    post:
      summary: Create user
      description: Creates a new user in the system
      tags:
        - users
      consumes:
        - application/json
      produces:
        - application/json
      parameters:
        - in: body
          name: user
          description: User data
          required: true
          schema:
            $ref: '#/definitions/CreateUserRequest'
      responses:
        201:
          description: User created successfully
          schema:
            $ref: '#/definitions/UserResponse'
        400:
          description: Invalid data
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

## ⚠️ Important Considerations

### ✅ Best Practices
- **Protocol separation**: Each handler handles only its protocol
- **Protocol-specific DTOs**: Different DTOs per protocol if necessary
- **Error handling**: Consistent error handling per protocol
- **Context propagation**: Use context.Context in all operations

### ❌ Common Errors
- **Business logic in handlers**: Should be in use cases
- **Direct dependencies**: Don't access repositories directly
- **Exposing internal errors**: Map errors appropriately
- **Not validating input**: Always validate input data

---

**← [goca repository Command](Command-Repository) | [goca di Command](Command-DI) →**
