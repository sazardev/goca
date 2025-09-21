# goca interfaces Command

The `goca interfaces` command generates only the contract interfaces between layers, useful for test-driven development (TDD) and for defining clear contracts in Clean Architecture.

## 📋 Syntax

```bash
goca interfaces <entity> [flags]
```

## 🎯 Purpose

Creates contract interfaces for TDD development:

- 🔗 **Use case interfaces** for the application layer
- 📊 **Repository interfaces** for persistence
- 🟢 **Handler interfaces** for adapters
- 🧪 **TDD development** with contracts first
- 📝 **Internal API documentation**

## 🚩 Available Flags

| Flag           | Type   | Required | Default Value | Description                    |
| -------------- | ------ | -------- | ------------- | ------------------------------ |
| `--all`        | `bool` | ❌ No     | `false`       | Generate all interfaces        |
| `--usecase`    | `bool` | ❌ No     | `false`       | Generate use case interfaces   |
| `--repository` | `bool` | ❌ No     | `false`       | Generate repository interfaces |
| `--handler`    | `bool` | ❌ No     | `false`       | Generate handler interfaces    |

## 📖 Usage Examples

### All Interfaces
```bash
goca interfaces User --all
```

### Use Case Interfaces Only
```bash
goca interfaces Product --usecase
```

### Repository Interfaces Only
```bash
goca interfaces Order --repository
```

### Handler Interfaces Only
```bash
goca interfaces Customer --handler
```

### Specific Combination
```bash
goca interfaces User --usecase --repository
```

## 📂 Generated Files

### File Structure
```
internal/interfaces/
├── user_usecase.go        # Use case interfaces
├── user_repository.go     # Repository interfaces
└── user_handler.go        # Handler interfaces
```

## 🔍 Generated Code in Detail

### Use Case Interfaces: `internal/interfaces/user_usecase.go`

```go
package interfaces

import (
    "context"
    
    "github.com/usuario/proyecto/internal/usecase/dto"
)

//go:generate mockgen -source=user_usecase.go -destination=mocks/user_usecase_mock.go

// UserUseCase defines contracts for user use cases
type UserUseCase interface {
    // Basic CRUD operations
    Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
    GetByID(ctx context.Context, id uint) (*dto.UserResponse, error)
    Update(ctx context.Context, id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, req dto.ListUsersRequest) (*dto.ListUsersResponse, error)
    
    // Search operations
    Search(ctx context.Context, query string, req dto.ListUsersRequest) (*dto.ListUsersResponse, error)
    FindByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
    
    // Business operations
    Activate(ctx context.Context, id uint) error
    Deactivate(ctx context.Context, id uint) error
    ChangePassword(ctx context.Context, id uint, req dto.ChangePasswordRequest) error
    
    // Statistics operations
    GetUserStats(ctx context.Context, id uint) (*dto.UserStatsResponse, error)
    GetUsersCount(ctx context.Context) (int64, error)
    
    // Batch operations
    CreateBatch(ctx context.Context, users []dto.CreateUserRequest) (*dto.BatchCreateResponse, error)
    UpdateBatch(ctx context.Context, updates []dto.BatchUpdateUserRequest) (*dto.BatchUpdateResponse, error)
    DeleteBatch(ctx context.Context, ids []uint) (*dto.BatchDeleteResponse, error)
}

// UserNotificationUseCase interface for user notifications
type UserNotificationUseCase interface {
    SendWelcomeEmail(ctx context.Context, userID uint) error
    SendPasswordResetEmail(ctx context.Context, userID uint) error
    SendActivationEmail(ctx context.Context, userID uint) error
    NotifyUserUpdate(ctx context.Context, userID uint, changes map[string]interface{}) error
}

// UserValidationUseCase interface for advanced validations
type UserValidationUseCase interface {
    ValidateUserCreation(ctx context.Context, req dto.CreateUserRequest) error
    ValidateUserUpdate(ctx context.Context, id uint, req dto.UpdateUserRequest) error
    ValidateEmailUniqueness(ctx context.Context, email string, excludeID *uint) error
    ValidateUserPermissions(ctx context.Context, userID uint, action string) error
}

// UserAnalyticsUseCase interface for user analytics
type UserAnalyticsUseCase interface {
    GetUserActivity(ctx context.Context, userID uint, from, to time.Time) (*dto.UserActivityResponse, error)
    GetUserEngagement(ctx context.Context, userID uint) (*dto.UserEngagementResponse, error)
    GetUsersGrowth(ctx context.Context, period string) (*dto.UsersGrowthResponse, error)
    GetActiveUsersCount(ctx context.Context, period string) (int64, error)
}
```

### Repository Interfaces: `internal/interfaces/user_repository.go`

```go
package interfaces

import (
    "context"
    "time"
    
    "github.com/usuario/proyecto/internal/domain"
)

//go:generate mockgen -source=user_repository.go -destination=mocks/user_repository_mock.go

// UserRepository defines contracts for user persistence
type UserRepository interface {
    // Basic CRUD operations
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id uint) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uint) error
    
    // Query operations
    List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    Search(ctx context.Context, query string, limit, offset int) ([]*domain.User, int64, error)
    Exists(ctx context.Context, id uint) (bool, error)
    ExistsByEmail(ctx context.Context, email string) (bool, error)
    
    // Filtering operations
    FindByStatus(ctx context.Context, status string, limit, offset int) ([]*domain.User, int64, error)
    FindByDateRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*domain.User, int64, error)
    FindActive(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
    FindInactive(ctx context.Context, inactiveDays int, limit, offset int) ([]*domain.User, int64, error)
    
    // Aggregation operations
    Count(ctx context.Context) (int64, error)
    CountByStatus(ctx context.Context, status string) (int64, error)
    CountByDateRange(ctx context.Context, from, to time.Time) (int64, error)
    CountActive(ctx context.Context) (int64, error)
    
    // Batch operations
    SaveBatch(ctx context.Context, users []*domain.User) error
    UpdateBatch(ctx context.Context, users []*domain.User) error
    DeleteBatch(ctx context.Context, ids []uint) error
    FindByIDs(ctx context.Context, ids []uint) ([]*domain.User, error)
    
    // Transaction operations
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
    
    // Cache operations
    ClearCache(ctx context.Context, id uint) error
    ClearAllCache(ctx context.Context) error
}

// UserAuditRepository interface for user auditing
type UserAuditRepository interface {
    LogUserAction(ctx context.Context, userID uint, action string, details map[string]interface{}) error
    GetUserAuditLog(ctx context.Context, userID uint, limit, offset int) ([]*domain.UserAuditLog, error)
    GetAuditLogByAction(ctx context.Context, action string, limit, offset int) ([]*domain.UserAuditLog, error)
    GetAuditLogByDateRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*domain.UserAuditLog, error)
}

// UserSessionRepository interface for user sessions
type UserSessionRepository interface {
    CreateSession(ctx context.Context, session *domain.UserSession) error
    GetSession(ctx context.Context, token string) (*domain.UserSession, error)
    GetUserSessions(ctx context.Context, userID uint) ([]*domain.UserSession, error)
    UpdateSession(ctx context.Context, session *domain.UserSession) error
    DeleteSession(ctx context.Context, token string) error
    DeleteUserSessions(ctx context.Context, userID uint) error
    DeleteExpiredSessions(ctx context.Context) error
}

// UserStatsRepository interface for user statistics
type UserStatsRepository interface {
    GetUserStats(ctx context.Context, userID uint) (*domain.UserStats, error)
    UpdateUserStats(ctx context.Context, stats *domain.UserStats) error
    GetGlobalStats(ctx context.Context) (*domain.GlobalUserStats, error)
    GetStatsHistory(ctx context.Context, userID uint, days int) ([]*domain.UserStatsHistory, error)
}
```

### Handler Interfaces: `internal/interfaces/user_handler.go`

```go
package interfaces

import (
    "context"
    
    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"
)

//go:generate mockgen -source=user_handler.go -destination=mocks/user_handler_mock.go

// UserHTTPHandler defines contracts for HTTP handlers
type UserHTTPHandler interface {
    // REST CRUD operations
    Create(c *gin.Context)
    GetByID(c *gin.Context)
    Update(c *gin.Context)
    Delete(c *gin.Context)
    List(c *gin.Context)
    
    // Search operations
    Search(c *gin.Context)
    FindByEmail(c *gin.Context)
    
    // Business operations
    Activate(c *gin.Context)
    Deactivate(c *gin.Context)
    ChangePassword(c *gin.Context)
    
    // Statistics operations
    GetStats(c *gin.Context)
    GetActivity(c *gin.Context)
    
    // Batch operations
    CreateBatch(c *gin.Context)
    UpdateBatch(c *gin.Context)
    DeleteBatch(c *gin.Context)
    
    // File operations
    UploadAvatar(c *gin.Context)
    DownloadData(c *gin.Context)
    ImportUsers(c *gin.Context)
    ExportUsers(c *gin.Context)
}

// UserGRPCHandler defines contracts for gRPC handlers
type UserGRPCHandler interface {
    // gRPC CRUD operations
    CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
    GetUser(ctx context.Context, req *GetUserRequest) (*UserResponse, error)
    UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UserResponse, error)
    DeleteUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error)
    ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error)
    
    // Streaming operations
    StreamUsers(req *StreamUsersRequest, stream UserService_StreamUsersServer) error
    BulkCreateUsers(stream UserService_BulkCreateUsersServer) error
    
    // Business operations
    ActivateUser(ctx context.Context, req *ActivateUserRequest) (*ActivateUserResponse, error)
    ValidateUser(ctx context.Context, req *ValidateUserRequest) (*ValidateUserResponse, error)
    
    // Statistics operations
    GetUserStats(ctx context.Context, req *GetUserStatsRequest) (*UserStatsResponse, error)
    GetUsersMetrics(ctx context.Context, req *GetUsersMetricsRequest) (*UsersMetricsResponse, error)
}

// UserCLIHandler defines contracts for CLI handlers
type UserCLIHandler interface {
    // CRUD commands
    CreateUserCommand() CLICommand
    GetUserCommand() CLICommand
    UpdateUserCommand() CLICommand
    DeleteUserCommand() CLICommand
    ListUsersCommand() CLICommand
    
    // Administration commands
    ActivateUserCommand() CLICommand
    DeactivateUserCommand() CLICommand
    ResetPasswordCommand() CLICommand
    
    // Import/export commands
    ImportUsersCommand() CLICommand
    ExportUsersCommand() CLICommand
    
    // Statistics commands
    UserStatsCommand() CLICommand
    UsersReportCommand() CLICommand
    
    // Maintenance commands
    CleanupUsersCommand() CLICommand
    ValidateUsersCommand() CLICommand
}

// UserWorkerHandler defines contracts for workers
type UserWorkerHandler interface {
    // Task processing
    ProcessUserTask(ctx context.Context, taskData []byte) error
    
    // Specific tasks
    ProcessWelcomeEmail(ctx context.Context, userID uint) error
    ProcessPasswordReset(ctx context.Context, userID uint) error
    ProcessUserActivation(ctx context.Context, userID uint) error
    ProcessUserDeactivation(ctx context.Context, userID uint) error
    
    // Batch tasks
    ProcessBatchUserCreation(ctx context.Context, userData []byte) error
    ProcessBatchUserUpdate(ctx context.Context, userData []byte) error
    ProcessBatchUserDeletion(ctx context.Context, userIDs []uint) error
    
    // Maintenance tasks
    ProcessInactiveUsersCleanup(ctx context.Context) error
    ProcessUserStatsUpdate(ctx context.Context) error
    ProcessUserDataExport(ctx context.Context, exportID string) error
    
    // Worker control
    StartWorker(ctx context.Context) error
    StopWorker(ctx context.Context) error
    GetWorkerStatus() WorkerStatus
}

// UserSOAPHandler defines contracts for SOAP services
type UserSOAPHandler interface {
    // SOAP operations
    CreateUser(ctx context.Context, req *SOAPCreateUserRequest) (*SOAPUserResponse, error)
    GetUser(ctx context.Context, req *SOAPGetUserRequest) (*SOAPUserResponse, error)
    UpdateUser(ctx context.Context, req *SOAPUpdateUserRequest) (*SOAPUserResponse, error)
    DeleteUser(ctx context.Context, req *SOAPDeleteUserRequest) (*SOAPDeleteUserResponse, error)
    ListUsers(ctx context.Context, req *SOAPListUsersRequest) (*SOAPListUsersResponse, error)
    
    // SOAP validation operations
    ValidateUserData(ctx context.Context, req *SOAPValidateUserRequest) (*SOAPValidationResponse, error)
    
    // SOAP request handling
    HandleSOAPRequest(w http.ResponseWriter, r *http.Request)
    ProcessSOAPEnvelope(envelope *SOAPEnvelope) (*SOAPEnvelope, error)
}

// CLICommand defines the structure of a CLI command
type CLICommand interface {
    GetName() string
    GetDescription() string
    GetUsage() string
    Execute(args []string) error
    GetFlags() []CLIFlag
}

// CLIFlag defines the structure of a CLI flag
type CLIFlag struct {
    Name        string
    ShortName   string
    Description string
    Required    bool
    Default     interface{}
}

// WorkerStatus defines the status of a worker
type WorkerStatus struct {
    IsRunning     bool
    TasksProcessed int64
    Errors        int64
    LastActivity  time.Time
    Uptime        time.Duration
}
```

## 🧪 Mock Generation

Interfaces include directives for `go generate`:

```go
//go:generate mockgen -source=user_usecase.go -destination=mocks/user_usecase_mock.go
```

### Commands to generate mocks:
```bash
# Install mockgen
go install github.com/golang/mock/mockgen@latest

# Generate all mocks
cd internal/interfaces
go generate ./...
```

### Using mocks in tests:
```go
func TestUserService_Create(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockUserRepository(ctrl)
    service := usecase.NewUserService(mockRepo)
    
    req := dto.CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    mockRepo.EXPECT().
        FindByEmail(gomock.Any(), req.Email).
        Return(nil, nil)
    
    mockRepo.EXPECT().
        Save(gomock.Any(), gomock.Any()).
        DoAndReturn(func(ctx context.Context, user *domain.User) error {
            user.ID = 1
            return nil
        })
    
    result, err := service.Create(context.Background(), req)
    
    assert.NoError(t, err)
    assert.Equal(t, uint(1), result.ID)
    assert.Equal(t, req.Name, result.Name)
    assert.Equal(t, req.Email, result.Email)
}
```

## 🔄 TDD Development

### Recommended TDD Flow:

1. **Generate interfaces**:
```bash
goca interfaces User --all
```

2. **Write tests with mocks**:
```go
func TestUserUseCase_Create(t *testing.T) {
    // Test using the interface
}
```

3. **Implement use cases**:
```bash
goca usecase UserService --entity User
```

4. **Implement repositories**:
```bash
goca repository User --database postgres
```

5. **Implement handlers**:
```bash
goca handler User --type http
```

## 📝 Contract Documentation

Interfaces serve as living documentation:

```go
// UserUseCase defines all contracts for user use cases.
// This interface establishes expected behavior without revealing
// implementation details, allowing flexibility and testability.
type UserUseCase interface {
    // Create creates a new user in the system.
    // Validates input data and returns the created user.
    // Returns error if email already exists or data is invalid.
    Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
    
    // GetByID gets a user by their unique ID.
    // Returns error if user doesn't exist or is marked as deleted.
    GetByID(ctx context.Context, id uint) (*dto.UserResponse, error)
}
```

## ⚠️ Important Considerations

### ✅ Best Practices
- **Small interfaces**: Follow interface segregation principle
- **Context first**: Always use context.Context as first parameter
- **Error handling**: Return descriptive errors
- **Documentation**: Document all interface methods

### ❌ Common Mistakes
- **Interfaces too large**: Divide into specific interfaces
- **Concrete dependencies**: Interfaces should not depend on implementations
- **Mixing concerns**: Separate responsibilities into different interfaces
- **Missing context**: Always propagate context

### 🔄 Tool Integration

#### GoMock
```bash
# Install
go install github.com/golang/mock/mockgen@latest

# Generate mocks
mockgen -source=interfaces/user_usecase.go -destination=mocks/user_usecase_mock.go
```

#### Testify
```go
import (
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/assert"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}
```

---

**← [Comando goca di](Command-DI) | [Comando goca messages](Command-Messages) →**
