# Comando goca interfaces

El comando `goca interfaces` genera únicamente las interfaces de contratos entre capas, útil para desarrollo dirigido por pruebas (TDD) y para definir contratos claros en Clean Architecture.

## 📋 Sintaxis

```bash
goca interfaces <entity> [flags]
```

## 🎯 Propósito

Crea interfaces de contratos para desarrollo TDD:

- 🔗 **Interfaces de casos de uso** para la capa de aplicación
- 📊 **Interfaces de repositorio** para persistencia
- 🟢 **Interfaces de handlers** para adaptadores
- 🧪 **Desarrollo TDD** con contratos primero
- 📝 **Documentación de APIs** internas

## 🚩 Flags Disponibles

| Flag           | Tipo   | Requerido | Valor por Defecto | Descripción                        |
| -------------- | ------ | --------- | ----------------- | ---------------------------------- |
| `--all`        | `bool` | ❌ No      | `false`           | Generar todas las interfaces       |
| `--usecase`    | `bool` | ❌ No      | `false`           | Generar interfaces de casos de uso |
| `--repository` | `bool` | ❌ No      | `false`           | Generar interfaces de repositorio  |
| `--handler`    | `bool` | ❌ No      | `false`           | Generar interfaces de handlers     |

## 📖 Ejemplos de Uso

### Todas las Interfaces
```bash
goca interfaces User --all
```

### Solo Interfaces de Casos de Uso
```bash
goca interfaces Product --usecase
```

### Solo Interfaces de Repositorio
```bash
goca interfaces Order --repository
```

### Solo Interfaces de Handlers
```bash
goca interfaces Customer --handler
```

### Combinación Específica
```bash
goca interfaces User --usecase --repository
```

## 📂 Archivos Generados

### Estructura de Archivos
```
internal/interfaces/
├── user_usecase.go        # Interfaces de casos de uso
├── user_repository.go     # Interfaces de repositorio
└── user_handler.go        # Interfaces de handlers
```

## 🔍 Código Generado en Detalle

### Interfaces de Casos de Uso: `internal/interfaces/user_usecase.go`

```go
package interfaces

import (
    "context"
    
    "github.com/usuario/proyecto/internal/usecase/dto"
)

//go:generate mockgen -source=user_usecase.go -destination=mocks/user_usecase_mock.go

// UserUseCase define los contratos para los casos de uso de usuario
type UserUseCase interface {
    // Operaciones CRUD básicas
    Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
    GetByID(ctx context.Context, id uint) (*dto.UserResponse, error)
    Update(ctx context.Context, id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, req dto.ListUsersRequest) (*dto.ListUsersResponse, error)
    
    // Operaciones de búsqueda
    Search(ctx context.Context, query string, req dto.ListUsersRequest) (*dto.ListUsersResponse, error)
    FindByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
    
    // Operaciones de negocio
    Activate(ctx context.Context, id uint) error
    Deactivate(ctx context.Context, id uint) error
    ChangePassword(ctx context.Context, id uint, req dto.ChangePasswordRequest) error
    
    // Operaciones de estadísticas
    GetUserStats(ctx context.Context, id uint) (*dto.UserStatsResponse, error)
    GetUsersCount(ctx context.Context) (int64, error)
    
    // Operaciones en lote
    CreateBatch(ctx context.Context, users []dto.CreateUserRequest) (*dto.BatchCreateResponse, error)
    UpdateBatch(ctx context.Context, updates []dto.BatchUpdateUserRequest) (*dto.BatchUpdateResponse, error)
    DeleteBatch(ctx context.Context, ids []uint) (*dto.BatchDeleteResponse, error)
}

// UserNotificationUseCase interface para notificaciones de usuario
type UserNotificationUseCase interface {
    SendWelcomeEmail(ctx context.Context, userID uint) error
    SendPasswordResetEmail(ctx context.Context, userID uint) error
    SendActivationEmail(ctx context.Context, userID uint) error
    NotifyUserUpdate(ctx context.Context, userID uint, changes map[string]interface{}) error
}

// UserValidationUseCase interface para validaciones avanzadas
type UserValidationUseCase interface {
    ValidateUserCreation(ctx context.Context, req dto.CreateUserRequest) error
    ValidateUserUpdate(ctx context.Context, id uint, req dto.UpdateUserRequest) error
    ValidateEmailUniqueness(ctx context.Context, email string, excludeID *uint) error
    ValidateUserPermissions(ctx context.Context, userID uint, action string) error
}

// UserAnalyticsUseCase interface para análisis de usuarios
type UserAnalyticsUseCase interface {
    GetUserActivity(ctx context.Context, userID uint, from, to time.Time) (*dto.UserActivityResponse, error)
    GetUserEngagement(ctx context.Context, userID uint) (*dto.UserEngagementResponse, error)
    GetUsersGrowth(ctx context.Context, period string) (*dto.UsersGrowthResponse, error)
    GetActiveUsersCount(ctx context.Context, period string) (int64, error)
}
```

### Interfaces de Repositorio: `internal/interfaces/user_repository.go`

```go
package interfaces

import (
    "context"
    "time"
    
    "github.com/usuario/proyecto/internal/domain"
)

//go:generate mockgen -source=user_repository.go -destination=mocks/user_repository_mock.go

// UserRepository define los contratos para la persistencia de usuarios
type UserRepository interface {
    // Operaciones CRUD básicas
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id uint) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uint) error
    
    // Operaciones de consulta
    List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    Search(ctx context.Context, query string, limit, offset int) ([]*domain.User, int64, error)
    Exists(ctx context.Context, id uint) (bool, error)
    ExistsByEmail(ctx context.Context, email string) (bool, error)
    
    // Operaciones de filtrado
    FindByStatus(ctx context.Context, status string, limit, offset int) ([]*domain.User, int64, error)
    FindByDateRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*domain.User, int64, error)
    FindActive(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
    FindInactive(ctx context.Context, inactiveDays int, limit, offset int) ([]*domain.User, int64, error)
    
    // Operaciones de agregación
    Count(ctx context.Context) (int64, error)
    CountByStatus(ctx context.Context, status string) (int64, error)
    CountByDateRange(ctx context.Context, from, to time.Time) (int64, error)
    CountActive(ctx context.Context) (int64, error)
    
    // Operaciones en lote
    SaveBatch(ctx context.Context, users []*domain.User) error
    UpdateBatch(ctx context.Context, users []*domain.User) error
    DeleteBatch(ctx context.Context, ids []uint) error
    FindByIDs(ctx context.Context, ids []uint) ([]*domain.User, error)
    
    // Operaciones de transacción
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
    
    // Operaciones de caché
    ClearCache(ctx context.Context, id uint) error
    ClearAllCache(ctx context.Context) error
}

// UserAuditRepository interface para auditoría de usuarios
type UserAuditRepository interface {
    LogUserAction(ctx context.Context, userID uint, action string, details map[string]interface{}) error
    GetUserAuditLog(ctx context.Context, userID uint, limit, offset int) ([]*domain.UserAuditLog, error)
    GetAuditLogByAction(ctx context.Context, action string, limit, offset int) ([]*domain.UserAuditLog, error)
    GetAuditLogByDateRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*domain.UserAuditLog, error)
}

// UserSessionRepository interface para sesiones de usuario
type UserSessionRepository interface {
    CreateSession(ctx context.Context, session *domain.UserSession) error
    GetSession(ctx context.Context, token string) (*domain.UserSession, error)
    GetUserSessions(ctx context.Context, userID uint) ([]*domain.UserSession, error)
    UpdateSession(ctx context.Context, session *domain.UserSession) error
    DeleteSession(ctx context.Context, token string) error
    DeleteUserSessions(ctx context.Context, userID uint) error
    DeleteExpiredSessions(ctx context.Context) error
}

// UserStatsRepository interface para estadísticas de usuario
type UserStatsRepository interface {
    GetUserStats(ctx context.Context, userID uint) (*domain.UserStats, error)
    UpdateUserStats(ctx context.Context, stats *domain.UserStats) error
    GetGlobalStats(ctx context.Context) (*domain.GlobalUserStats, error)
    GetStatsHistory(ctx context.Context, userID uint, days int) ([]*domain.UserStatsHistory, error)
}
```

### Interfaces de Handlers: `internal/interfaces/user_handler.go`

```go
package interfaces

import (
    "context"
    
    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"
)

//go:generate mockgen -source=user_handler.go -destination=mocks/user_handler_mock.go

// UserHTTPHandler define los contratos para handlers HTTP
type UserHTTPHandler interface {
    // Operaciones CRUD REST
    Create(c *gin.Context)
    GetByID(c *gin.Context)
    Update(c *gin.Context)
    Delete(c *gin.Context)
    List(c *gin.Context)
    
    // Operaciones de búsqueda
    Search(c *gin.Context)
    FindByEmail(c *gin.Context)
    
    // Operaciones de negocio
    Activate(c *gin.Context)
    Deactivate(c *gin.Context)
    ChangePassword(c *gin.Context)
    
    // Operaciones de estadísticas
    GetStats(c *gin.Context)
    GetActivity(c *gin.Context)
    
    // Operaciones en lote
    CreateBatch(c *gin.Context)
    UpdateBatch(c *gin.Context)
    DeleteBatch(c *gin.Context)
    
    // Operaciones de archivos
    UploadAvatar(c *gin.Context)
    DownloadData(c *gin.Context)
    ImportUsers(c *gin.Context)
    ExportUsers(c *gin.Context)
}

// UserGRPCHandler define los contratos para handlers gRPC
type UserGRPCHandler interface {
    // Operaciones CRUD gRPC
    CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
    GetUser(ctx context.Context, req *GetUserRequest) (*UserResponse, error)
    UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UserResponse, error)
    DeleteUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error)
    ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error)
    
    // Operaciones de streaming
    StreamUsers(req *StreamUsersRequest, stream UserService_StreamUsersServer) error
    BulkCreateUsers(stream UserService_BulkCreateUsersServer) error
    
    // Operaciones de negocio
    ActivateUser(ctx context.Context, req *ActivateUserRequest) (*ActivateUserResponse, error)
    ValidateUser(ctx context.Context, req *ValidateUserRequest) (*ValidateUserResponse, error)
    
    // Operaciones de estadísticas
    GetUserStats(ctx context.Context, req *GetUserStatsRequest) (*UserStatsResponse, error)
    GetUsersMetrics(ctx context.Context, req *GetUsersMetricsRequest) (*UsersMetricsResponse, error)
}

// UserCLIHandler define los contratos para handlers CLI
type UserCLIHandler interface {
    // Comandos CRUD
    CreateUserCommand() CLICommand
    GetUserCommand() CLICommand
    UpdateUserCommand() CLICommand
    DeleteUserCommand() CLICommand
    ListUsersCommand() CLICommand
    
    // Comandos de administración
    ActivateUserCommand() CLICommand
    DeactivateUserCommand() CLICommand
    ResetPasswordCommand() CLICommand
    
    // Comandos de importación/exportación
    ImportUsersCommand() CLICommand
    ExportUsersCommand() CLICommand
    
    // Comandos de estadísticas
    UserStatsCommand() CLICommand
    UsersReportCommand() CLICommand
    
    // Comandos de mantenimiento
    CleanupUsersCommand() CLICommand
    ValidateUsersCommand() CLICommand
}

// UserWorkerHandler define los contratos para workers
type UserWorkerHandler interface {
    // Procesamiento de tareas
    ProcessUserTask(ctx context.Context, taskData []byte) error
    
    // Tareas específicas
    ProcessWelcomeEmail(ctx context.Context, userID uint) error
    ProcessPasswordReset(ctx context.Context, userID uint) error
    ProcessUserActivation(ctx context.Context, userID uint) error
    ProcessUserDeactivation(ctx context.Context, userID uint) error
    
    // Tareas en lote
    ProcessBatchUserCreation(ctx context.Context, userData []byte) error
    ProcessBatchUserUpdate(ctx context.Context, userData []byte) error
    ProcessBatchUserDeletion(ctx context.Context, userIDs []uint) error
    
    // Tareas de mantenimiento
    ProcessInactiveUsersCleanup(ctx context.Context) error
    ProcessUserStatsUpdate(ctx context.Context) error
    ProcessUserDataExport(ctx context.Context, exportID string) error
    
    // Control de workers
    StartWorker(ctx context.Context) error
    StopWorker(ctx context.Context) error
    GetWorkerStatus() WorkerStatus
}

// UserSOAPHandler define los contratos para servicios SOAP
type UserSOAPHandler interface {
    // Operaciones SOAP
    CreateUser(ctx context.Context, req *SOAPCreateUserRequest) (*SOAPUserResponse, error)
    GetUser(ctx context.Context, req *SOAPGetUserRequest) (*SOAPUserResponse, error)
    UpdateUser(ctx context.Context, req *SOAPUpdateUserRequest) (*SOAPUserResponse, error)
    DeleteUser(ctx context.Context, req *SOAPDeleteUserRequest) (*SOAPDeleteUserResponse, error)
    ListUsers(ctx context.Context, req *SOAPListUsersRequest) (*SOAPListUsersResponse, error)
    
    // Operaciones de validación SOAP
    ValidateUserData(ctx context.Context, req *SOAPValidateUserRequest) (*SOAPValidationResponse, error)
    
    // Manejo de peticiones SOAP
    HandleSOAPRequest(w http.ResponseWriter, r *http.Request)
    ProcessSOAPEnvelope(envelope *SOAPEnvelope) (*SOAPEnvelope, error)
}

// CLICommand define la estructura de un comando CLI
type CLICommand interface {
    GetName() string
    GetDescription() string
    GetUsage() string
    Execute(args []string) error
    GetFlags() []CLIFlag
}

// CLIFlag define la estructura de un flag CLI
type CLIFlag struct {
    Name        string
    ShortName   string
    Description string
    Required    bool
    Default     interface{}
}

// WorkerStatus define el estado de un worker
type WorkerStatus struct {
    IsRunning     bool
    TasksProcessed int64
    Errors        int64
    LastActivity  time.Time
    Uptime        time.Duration
}
```

## 🧪 Generación de Mocks

Las interfaces incluyen directivas para `go generate`:

```go
//go:generate mockgen -source=user_usecase.go -destination=mocks/user_usecase_mock.go
```

### Comandos para generar mocks:
```bash
# Instalar mockgen
go install github.com/golang/mock/mockgen@latest

# Generar todos los mocks
cd internal/interfaces
go generate ./...
```

### Uso de mocks en tests:
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

## 🔄 Desarrollo TDD

### Flujo TDD Recomendado:

1. **Generar interfaces**:
```bash
goca interfaces User --all
```

2. **Escribir tests con mocks**:
```go
func TestUserUseCase_Create(t *testing.T) {
    // Test usando la interface
}
```

3. **Implementar casos de uso**:
```bash
goca usecase UserService --entity User
```

4. **Implementar repositorios**:
```bash
goca repository User --database postgres
```

5. **Implementar handlers**:
```bash
goca handler User --type http
```

## 📝 Documentación de Contratos

Las interfaces sirven como documentación viva:

```go
// UserUseCase define todos los contratos para casos de uso de usuario.
// Esta interface establece el comportamiento esperado sin revelar
// detalles de implementación, permitiendo flexibilidad y testabilidad.
type UserUseCase interface {
    // Create crea un nuevo usuario en el sistema.
    // Valida los datos de entrada y retorna el usuario creado.
    // Retorna error si el email ya existe o los datos son inválidos.
    Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
    
    // GetByID obtiene un usuario por su ID único.
    // Retorna error si el usuario no existe o está marcado como eliminado.
    GetByID(ctx context.Context, id uint) (*dto.UserResponse, error)
}
```

## ⚠️ Consideraciones Importantes

### ✅ Buenas Prácticas
- **Interfaces pequeñas**: Seguir el principio de segregación de interfaces
- **Context first**: Siempre usar context.Context como primer parámetro
- **Error handling**: Retornar errores descriptivos
- **Documentation**: Documentar todos los métodos de interface

### ❌ Errores Comunes
- **Interfaces demasiado grandes**: Dividir en interfaces específicas
- **Dependencias concretas**: Las interfaces no deben depender de implementaciones
- **Mixing concerns**: Separar responsabilidades en diferentes interfaces
- **Missing context**: Siempre propagar el contexto

### 🔄 Integración con Herramientas

#### GoMock
```bash
# Instalar
go install github.com/golang/mock/mockgen@latest

# Generar mocks
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
