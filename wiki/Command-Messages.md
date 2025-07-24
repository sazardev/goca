# Comando goca messages

El comando `goca messages` crea archivos de mensajes de error, respuestas y constantes organizados por feature para mantener consistencia en la aplicación.

## 📋 Sintaxis

```bash
goca messages <entity> [flags]
```

## 🎯 Propósito

Crea sistemas de mensajes organizados y consistentes:

- ❌ **Mensajes de error** centralizados y tipados
- ✅ **Mensajes de respuesta** para diferentes operaciones
- 📄 **Constantes** del sistema organizadas por feature
- 🌍 **Internacionalización** preparada (i18n)
- 📊 **Códigos de estado** consistentes

## 🚩 Flags Disponibles

| Flag          | Tipo   | Requerido | Valor por Defecto | Descripción                         |
| ------------- | ------ | --------- | ----------------- | ----------------------------------- |
| `--all`       | `bool` | ❌ No      | `false`           | Generar todos los tipos de mensajes |
| `--errors`    | `bool` | ❌ No      | `false`           | Generar mensajes de error           |
| `--responses` | `bool` | ❌ No      | `false`           | Generar mensajes de respuesta       |
| `--constants` | `bool` | ❌ No      | `false`           | Generar constantes del feature      |

## 📖 Ejemplos de Uso

### Todos los Mensajes
```bash
goca messages User --all
```

### Solo Mensajes de Error
```bash
goca messages Product --errors
```

### Solo Mensajes de Respuesta
```bash
goca messages Order --responses
```

### Solo Constantes
```bash
goca messages Customer --constants
```

### Combinación Específica
```bash
goca messages User --errors --responses
```

## 📂 Archivos Generados

### Estructura de Archivos
```
internal/messages/
├── errors.go              # Mensajes de error
├── responses.go            # Mensajes de respuesta
└── constants.go            # Constantes del sistema
```

## 🔍 Código Generado en Detalle

### Mensajes de Error: `internal/messages/errors.go`

```go
package messages

import (
    "errors"
    "fmt"
)

// User Error Types
const (
    // Validation errors
    ErrTypeUserValidation = "USER_VALIDATION_ERROR"
    ErrTypeUserNotFound   = "USER_NOT_FOUND"
    ErrTypeUserExists     = "USER_ALREADY_EXISTS"
    
    // Business logic errors
    ErrTypeUserInactive   = "USER_INACTIVE"
    ErrTypeUserBlocked    = "USER_BLOCKED"
    ErrTypeUserPermission = "USER_PERMISSION_DENIED"
    
    // System errors
    ErrTypeUserDatabase   = "USER_DATABASE_ERROR"
    ErrTypeUserInternal   = "USER_INTERNAL_ERROR"
)

// User Domain Errors
var (
    // Validation errors
    ErrUserNameRequired    = NewUserError(ErrTypeUserValidation, "user name is required")
    ErrUserNameTooShort    = NewUserError(ErrTypeUserValidation, "user name must be at least 2 characters")
    ErrUserNameTooLong     = NewUserError(ErrTypeUserValidation, "user name must be less than 100 characters")
    ErrUserEmailRequired   = NewUserError(ErrTypeUserValidation, "user email is required")
    ErrUserEmailInvalid    = NewUserError(ErrTypeUserValidation, "user email format is invalid")
    ErrUserPasswordWeak    = NewUserError(ErrTypeUserValidation, "user password is too weak")
    
    // Existence errors
    ErrUserNotFound        = NewUserError(ErrTypeUserNotFound, "user not found")
    ErrUserEmailExists     = NewUserError(ErrTypeUserExists, "user with this email already exists")
    ErrUserUsernameExists  = NewUserError(ErrTypeUserExists, "user with this username already exists")
    
    // State errors
    ErrUserInactive        = NewUserError(ErrTypeUserInactive, "user account is inactive")
    ErrUserBlocked         = NewUserError(ErrTypeUserBlocked, "user account is blocked")
    ErrUserDeleted         = NewUserError(ErrTypeUserInactive, "user account has been deleted")
    
    // Permission errors
    ErrUserUnauthorized    = NewUserError(ErrTypeUserPermission, "user is not authorized")
    ErrUserPermissionDenied = NewUserError(ErrTypeUserPermission, "user permission denied")
    ErrUserInvalidRole     = NewUserError(ErrTypeUserPermission, "user has invalid role")
    
    // Operation errors
    ErrUserCannotUpdate    = NewUserError(ErrTypeUserValidation, "user cannot be updated")
    ErrUserCannotDelete    = NewUserError(ErrTypeUserValidation, "user cannot be deleted")
    ErrUserCannotActivate  = NewUserError(ErrTypeUserValidation, "user cannot be activated")
    
    // System errors
    ErrUserDatabaseSave    = NewUserError(ErrTypeUserDatabase, "failed to save user to database")
    ErrUserDatabaseUpdate  = NewUserError(ErrTypeUserDatabase, "failed to update user in database")
    ErrUserDatabaseDelete  = NewUserError(ErrTypeUserDatabase, "failed to delete user from database")
    ErrUserInternal        = NewUserError(ErrTypeUserInternal, "internal user service error")
)

// UserError representa un error específico de usuario
type UserError struct {
    Type    string `json:"type"`
    Message string `json:"message"`
    Code    string `json:"code,omitempty"`
    Field   string `json:"field,omitempty"`
    Value   string `json:"value,omitempty"`
}

// Error implementa la interface error
func (e *UserError) Error() string {
    if e.Field != "" {
        return fmt.Sprintf("%s: %s (field: %s)", e.Type, e.Message, e.Field)
    }
    return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// NewUserError crea un nuevo error de usuario
func NewUserError(errorType, message string) *UserError {
    return &UserError{
        Type:    errorType,
        Message: message,
    }
}

// NewUserValidationError crea un error de validación con campo específico
func NewUserValidationError(message, field, value string) *UserError {
    return &UserError{
        Type:    ErrTypeUserValidation,
        Message: message,
        Field:   field,
        Value:   value,
    }
}

// WithField agrega información de campo al error
func (e *UserError) WithField(field string) *UserError {
    e.Field = field
    return e
}

// WithValue agrega valor que causó el error
func (e *UserError) WithValue(value string) *UserError {
    e.Value = value
    return e
}

// WithCode agrega código específico al error
func (e *UserError) WithCode(code string) *UserError {
    e.Code = code
    return e
}

// IsValidationError verifica si el error es de validación
func IsValidationError(err error) bool {
    if userErr, ok := err.(*UserError); ok {
        return userErr.Type == ErrTypeUserValidation
    }
    return false
}

// IsNotFoundError verifica si el error es de entidad no encontrada
func IsNotFoundError(err error) bool {
    if userErr, ok := err.(*UserError); ok {
        return userErr.Type == ErrTypeUserNotFound
    }
    return false
}

// IsExistsError verifica si el error es de entidad ya existente
func IsExistsError(err error) bool {
    if userErr, ok := err.(*UserError); ok {
        return userErr.Type == ErrTypeUserExists
    }
    return false
}

// IsPermissionError verifica si el error es de permisos
func IsPermissionError(err error) bool {
    if userErr, ok := err.(*UserError); ok {
        return userErr.Type == ErrTypeUserPermission
    }
    return false
}

// IsDatabaseError verifica si el error es de base de datos
func IsDatabaseError(err error) bool {
    if userErr, ok := err.(*UserError); ok {
        return userErr.Type == ErrTypeUserDatabase
    }
    return false
}

// GetErrorType retorna el tipo de error
func GetErrorType(err error) string {
    if userErr, ok := err.(*UserError); ok {
        return userErr.Type
    }
    return "UNKNOWN_ERROR"
}

// WrapError envuelve un error con contexto de usuario
func WrapError(err error, message string) *UserError {
    return &UserError{
        Type:    ErrTypeUserInternal,
        Message: fmt.Sprintf("%s: %v", message, err),
    }
}

// MultiError maneja múltiples errores
type MultiError struct {
    Errors []*UserError `json:"errors"`
}

// Error implementa la interface error para múltiples errores
func (m *MultiError) Error() string {
    if len(m.Errors) == 0 {
        return "no errors"
    }
    
    if len(m.Errors) == 1 {
        return m.Errors[0].Error()
    }
    
    return fmt.Sprintf("multiple errors: %d validation errors", len(m.Errors))
}

// Add agrega un error a la colección
func (m *MultiError) Add(err *UserError) {
    m.Errors = append(m.Errors, err)
}

// HasErrors verifica si hay errores
func (m *MultiError) HasErrors() bool {
    return len(m.Errors) > 0
}

// NewMultiError crea una nueva instancia de múltiples errores
func NewMultiError() *MultiError {
    return &MultiError{
        Errors: make([]*UserError, 0),
    }
}
```

### Mensajes de Respuesta: `internal/messages/responses.go`

```go
package messages

// User Success Messages
const (
    // CRUD operation messages
    MsgUserCreated     = "User created successfully"
    MsgUserUpdated     = "User updated successfully"
    MsgUserDeleted     = "User deleted successfully"
    MsgUserRetrieved   = "User retrieved successfully"
    MsgUsersListed     = "Users listed successfully"
    
    // State change messages
    MsgUserActivated   = "User activated successfully"
    MsgUserDeactivated = "User deactivated successfully"
    MsgUserBlocked     = "User blocked successfully"
    MsgUserUnblocked   = "User unblocked successfully"
    
    // Authentication messages
    MsgUserLoggedIn    = "User logged in successfully"
    MsgUserLoggedOut   = "User logged out successfully"
    MsgPasswordChanged = "Password changed successfully"
    MsgPasswordReset   = "Password reset email sent"
    
    // Email/Notification messages
    MsgWelcomeEmailSent     = "Welcome email sent successfully"
    MsgVerificationEmailSent = "Verification email sent successfully"
    MsgNotificationSent     = "Notification sent successfully"
    
    // Bulk operation messages
    MsgUsersBatchCreated = "Users batch created successfully"
    MsgUsersBatchUpdated = "Users batch updated successfully"
    MsgUsersBatchDeleted = "Users batch deleted successfully"
    
    // Export/Import messages
    MsgUsersExported = "Users exported successfully"
    MsgUsersImported = "Users imported successfully"
    
    // Validation messages
    MsgUserValidated = "User data validated successfully"
    MsgEmailVerified = "Email verified successfully"
    MsgPhoneVerified = "Phone number verified successfully"
)

// User Response Templates
var (
    // Success templates
    SuccessUserCreated = ResponseTemplate{
        Message: MsgUserCreated,
        Code:    "USER_CREATED",
        Status:  "success",
    }
    
    SuccessUserUpdated = ResponseTemplate{
        Message: MsgUserUpdated,
        Code:    "USER_UPDATED",
        Status:  "success",
    }
    
    SuccessUserDeleted = ResponseTemplate{
        Message: MsgUserDeleted,
        Code:    "USER_DELETED",
        Status:  "success",
    }
    
    SuccessUserRetrieved = ResponseTemplate{
        Message: MsgUserRetrieved,
        Code:    "USER_RETRIEVED",
        Status:  "success",
    }
    
    SuccessUsersListed = ResponseTemplate{
        Message: MsgUsersListed,
        Code:    "USERS_LISTED",
        Status:  "success",
    }
    
    // State change templates
    SuccessUserActivated = ResponseTemplate{
        Message: MsgUserActivated,
        Code:    "USER_ACTIVATED",
        Status:  "success",
    }
    
    SuccessUserDeactivated = ResponseTemplate{
        Message: MsgUserDeactivated,
        Code:    "USER_DEACTIVATED",
        Status:  "success",
    }
    
    // Authentication templates
    SuccessUserLoggedIn = ResponseTemplate{
        Message: MsgUserLoggedIn,
        Code:    "USER_LOGGED_IN",
        Status:  "success",
    }
    
    SuccessPasswordChanged = ResponseTemplate{
        Message: MsgPasswordChanged,
        Code:    "PASSWORD_CHANGED",
        Status:  "success",
    }
)

// ResponseTemplate estructura base para respuestas
type ResponseTemplate struct {
    Message    string                 `json:"message"`
    Code       string                 `json:"code"`
    Status     string                 `json:"status"`
    Data       interface{}            `json:"data,omitempty"`
    Meta       map[string]interface{} `json:"meta,omitempty"`
    Timestamp  string                 `json:"timestamp,omitempty"`
}

// SuccessResponse crea una respuesta de éxito
func SuccessResponse(template ResponseTemplate, data interface{}) ResponseTemplate {
    response := template
    response.Data = data
    response.Timestamp = getCurrentTimestamp()
    return response
}

// SuccessResponseWithMeta crea una respuesta de éxito con metadatos
func SuccessResponseWithMeta(template ResponseTemplate, data interface{}, meta map[string]interface{}) ResponseTemplate {
    response := template
    response.Data = data
    response.Meta = meta
    response.Timestamp = getCurrentTimestamp()
    return response
}

// PaginatedResponse crea una respuesta paginada
func PaginatedResponse(template ResponseTemplate, data interface{}, page, limit int, total int64) ResponseTemplate {
    totalPages := (total + int64(limit) - 1) / int64(limit)
    
    meta := map[string]interface{}{
        "page":         page,
        "limit":        limit,
        "total":        total,
        "total_pages":  totalPages,
        "has_next":     page < int(totalPages),
        "has_prev":     page > 1,
    }
    
    return SuccessResponseWithMeta(template, data, meta)
}

// BatchResponse crea una respuesta para operaciones en lote
func BatchResponse(template ResponseTemplate, successful, failed int, errors []string) ResponseTemplate {
    meta := map[string]interface{}{
        "successful": successful,
        "failed":     failed,
        "total":      successful + failed,
    }
    
    if len(errors) > 0 {
        meta["errors"] = errors
    }
    
    return SuccessResponseWithMeta(template, nil, meta)
}

// ErrorResponse crea una respuesta de error
func ErrorResponse(err error) ResponseTemplate {
    if userErr, ok := err.(*UserError); ok {
        return ResponseTemplate{
            Message:   userErr.Message,
            Code:      userErr.Type,
            Status:    "error",
            Timestamp: getCurrentTimestamp(),
            Meta: map[string]interface{}{
                "field": userErr.Field,
                "value": userErr.Value,
            },
        }
    }
    
    return ResponseTemplate{
        Message:   err.Error(),
        Code:      "INTERNAL_ERROR",
        Status:    "error",
        Timestamp: getCurrentTimestamp(),
    }
}

// ValidationErrorResponse crea una respuesta para errores de validación
func ValidationErrorResponse(multiErr *MultiError) ResponseTemplate {
    errors := make([]map[string]interface{}, len(multiErr.Errors))
    
    for i, err := range multiErr.Errors {
        errors[i] = map[string]interface{}{
            "type":    err.Type,
            "message": err.Message,
            "field":   err.Field,
            "value":   err.Value,
            "code":    err.Code,
        }
    }
    
    return ResponseTemplate{
        Message:   "Validation failed",
        Code:      "VALIDATION_ERROR",
        Status:    "error",
        Data:      errors,
        Timestamp: getCurrentTimestamp(),
        Meta: map[string]interface{}{
            "error_count": len(multiErr.Errors),
        },
    }
}

// getCurrentTimestamp retorna el timestamp actual en formato ISO
func getCurrentTimestamp() string {
    return time.Now().UTC().Format(time.RFC3339)
}

// NotificationMessages mensajes para notificaciones
var NotificationMessages = map[string]string{
    "welcome_email":       "Welcome to our platform! Please verify your email address.",
    "password_reset":      "Password reset requested. Click the link to reset your password.",
    "account_activation":  "Your account has been activated successfully.",
    "account_deactivation": "Your account has been deactivated.",
    "profile_updated":     "Your profile has been updated successfully.",
    "login_notification":  "New login detected from a new device.",
    "security_alert":      "Suspicious activity detected on your account.",
}

// GetNotificationMessage obtiene un mensaje de notificación
func GetNotificationMessage(key string) string {
    if message, exists := NotificationMessages[key]; exists {
        return message
    }
    return "Notification message"
}

// EmailTemplates plantillas para emails
var EmailTemplates = map[string]EmailTemplate{
    "welcome": {
        Subject: "Welcome to {{.AppName}}!",
        Body:    "Hi {{.UserName}}, welcome to our platform. Please verify your email by clicking the link below.",
    },
    "password_reset": {
        Subject: "Password Reset Request",
        Body:    "Hi {{.UserName}}, you requested a password reset. Click the link below to reset your password.",
    },
    "account_verification": {
        Subject: "Verify Your Account",
        Body:    "Hi {{.UserName}}, please verify your account by clicking the verification link below.",
    },
}

// EmailTemplate estructura para plantillas de email
type EmailTemplate struct {
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

// GetEmailTemplate obtiene una plantilla de email
func GetEmailTemplate(templateName string) (EmailTemplate, bool) {
    template, exists := EmailTemplates[templateName]
    return template, exists
}
```

### Constantes del Sistema: `internal/messages/constants.go`

```go
package messages

import "time"

// User Status Constants
const (
    UserStatusActive     = "active"
    UserStatusInactive   = "inactive"
    UserStatusPending    = "pending"
    UserStatusBlocked    = "blocked"
    UserStatusDeleted    = "deleted"
    UserStatusSuspended  = "suspended"
)

// User Role Constants
const (
    UserRoleAdmin     = "admin"
    UserRoleModerator = "moderator"
    UserRoleUser      = "user"
    UserRoleGuest     = "guest"
    UserRoleManager   = "manager"
    UserRoleSupport   = "support"
)

// User Permission Constants
const (
    PermissionUserCreate = "user:create"
    PermissionUserRead   = "user:read"
    PermissionUserUpdate = "user:update"
    PermissionUserDelete = "user:delete"
    PermissionUserList   = "user:list"
    PermissionUserExport = "user:export"
    PermissionUserImport = "user:import"
    PermissionUserBlock  = "user:block"
    PermissionUserActivate = "user:activate"
)

// User Validation Constants
const (
    UserNameMinLength     = 2
    UserNameMaxLength     = 100
    UserEmailMaxLength    = 255
    UserPasswordMinLength = 8
    UserPasswordMaxLength = 128
    UserBioMaxLength      = 500
    UserPhoneMaxLength    = 20
)

// User Default Values
const (
    DefaultUserStatus = UserStatusPending
    DefaultUserRole   = UserRoleUser
    DefaultPageLimit  = 10
    DefaultMaxLimit   = 100
)

// User Session Constants
const (
    SessionDuration       = 24 * time.Hour
    RefreshTokenDuration  = 7 * 24 * time.Hour
    PasswordResetDuration = 1 * time.Hour
    EmailTokenDuration    = 24 * time.Hour
)

// User Cache Constants
const (
    UserCachePrefix    = "user:"
    UserCacheDuration  = 30 * time.Minute
    UserListCachePrefix = "user:list:"
    UserStatsCachePrefix = "user:stats:"
)

// User Event Constants
const (
    EventUserCreated     = "user.created"
    EventUserUpdated     = "user.updated"
    EventUserDeleted     = "user.deleted"
    EventUserActivated   = "user.activated"
    EventUserDeactivated = "user.deactivated"
    EventUserBlocked     = "user.blocked"
    EventUserLoggedIn    = "user.logged_in"
    EventUserLoggedOut   = "user.logged_out"
    EventPasswordChanged = "user.password_changed"
)

// User Queue Constants
const (
    QueueUserEmail        = "user.email"
    QueueUserNotification = "user.notification"
    QueueUserExport       = "user.export"
    QueueUserImport       = "user.import"
    QueueUserCleanup      = "user.cleanup"
    QueueUserAnalytics    = "user.analytics"
)

// User File Constants
const (
    UserAvatarMaxSize     = 5 * 1024 * 1024 // 5MB
    UserExportMaxRecords  = 10000
    UserImportMaxRecords  = 1000
    UserImportBatchSize   = 100
    AllowedImageTypes     = "jpg,jpeg,png,gif"
    AllowedDocumentTypes  = "pdf,doc,docx"
)

// User API Constants
const (
    APIVersionV1          = "v1"
    APIVersionV2          = "v2"
    APIRateLimitPerMinute = 60
    APIRateLimitPerHour   = 1000
    APITimeout            = 30 * time.Second
)

// User Database Constants
const (
    UserTableName        = "users"
    UserSessionTableName = "user_sessions"
    UserAuditTableName   = "user_audit_logs"
    UserStatsTableName   = "user_stats"
    
    // Index names
    UserEmailIndex    = "idx_users_email"
    UserUsernameIndex = "idx_users_username"
    UserStatusIndex   = "idx_users_status"
    UserCreatedAtIndex = "idx_users_created_at"
)

// User HTTP Status Constants
const (
    StatusUserCreated   = 201
    StatusUserUpdated   = 200
    StatusUserDeleted   = 204
    StatusUserNotFound  = 404
    StatusUserConflict  = 409
    StatusUserForbidden = 403
)

// User Validation Patterns
const (
    EmailPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    PhonePattern    = `^\+?[1-9]\d{1,14}$`
    UsernamePattern = `^[a-zA-Z0-9_]{3,20}$`
    PasswordPattern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`
)

// GetValidUserStatuses retorna todos los estados válidos de usuario
func GetValidUserStatuses() []string {
    return []string{
        UserStatusActive,
        UserStatusInactive,
        UserStatusPending,
        UserStatusBlocked,
        UserStatusDeleted,
        UserStatusSuspended,
    }
}

// GetValidUserRoles retorna todos los roles válidos de usuario
func GetValidUserRoles() []string {
    return []string{
        UserRoleAdmin,
        UserRoleModerator,
        UserRoleUser,
        UserRoleGuest,
        UserRoleManager,
        UserRoleSupport,
    }
}

// GetUserPermissions retorna todos los permisos de usuario
func GetUserPermissions() []string {
    return []string{
        PermissionUserCreate,
        PermissionUserRead,
        PermissionUserUpdate,
        PermissionUserDelete,
        PermissionUserList,
        PermissionUserExport,
        PermissionUserImport,
        PermissionUserBlock,
        PermissionUserActivate,
    }
}

// IsValidUserStatus verifica si el estado de usuario es válido
func IsValidUserStatus(status string) bool {
    validStatuses := GetValidUserStatuses()
    for _, validStatus := range validStatuses {
        if status == validStatus {
            return true
        }
    }
    return false
}

// IsValidUserRole verifica si el rol de usuario es válido
func IsValidUserRole(role string) bool {
    validRoles := GetValidUserRoles()
    for _, validRole := range validRoles {
        if role == validRole {
            return true
        }
    }
    return false
}

// GetDefaultPaginationValues retorna valores por defecto para paginación
func GetDefaultPaginationValues() (int, int) {
    return 1, DefaultPageLimit // page, limit
}

// GetMaxAllowedLimit retorna el límite máximo permitido
func GetMaxAllowedLimit() int {
    return DefaultMaxLimit
}

// UserConfig configuración específica de usuario
type UserConfig struct {
    NameMinLength        int           `json:"name_min_length"`
    NameMaxLength        int           `json:"name_max_length"`
    EmailMaxLength       int           `json:"email_max_length"`
    PasswordMinLength    int           `json:"password_min_length"`
    DefaultStatus        string        `json:"default_status"`
    DefaultRole          string        `json:"default_role"`
    SessionDuration      time.Duration `json:"session_duration"`
    CacheDuration        time.Duration `json:"cache_duration"`
    MaxFileSize          int64         `json:"max_file_size"`
    AllowedImageTypes    []string      `json:"allowed_image_types"`
    RateLimitPerMinute   int           `json:"rate_limit_per_minute"`
}

// GetDefaultUserConfig retorna la configuración por defecto
func GetDefaultUserConfig() UserConfig {
    return UserConfig{
        NameMinLength:        UserNameMinLength,
        NameMaxLength:        UserNameMaxLength,
        EmailMaxLength:       UserEmailMaxLength,
        PasswordMinLength:    UserPasswordMinLength,
        DefaultStatus:        DefaultUserStatus,
        DefaultRole:          DefaultUserRole,
        SessionDuration:      SessionDuration,
        CacheDuration:        UserCacheDuration,
        MaxFileSize:          UserAvatarMaxSize,
        AllowedImageTypes:    []string{"jpg", "jpeg", "png", "gif"},
        RateLimitPerMinute:   APIRateLimitPerMinute,
    }
}

// Environment-specific constants
const (
    EnvDevelopment = "development"
    EnvTesting     = "testing"
    EnvStaging     = "staging"
    EnvProduction  = "production"
)

// Feature Flags
const (
    FeatureUserRegistration = "user_registration"
    FeatureEmailVerification = "email_verification"
    FeaturePhoneVerification = "phone_verification"
    FeatureTwoFactorAuth    = "two_factor_auth"
    FeatureSocialLogin      = "social_login"
    FeatureUserImport       = "user_import"
    FeatureUserExport       = "user_export"
    FeatureUserAnalytics    = "user_analytics"
)

// GetFeatureFlags retorna las feature flags disponibles
func GetFeatureFlags() []string {
    return []string{
        FeatureUserRegistration,
        FeatureEmailVerification,
        FeaturePhoneVerification,
        FeatureTwoFactorAuth,
        FeatureSocialLogin,
        FeatureUserImport,
        FeatureUserExport,
        FeatureUserAnalytics,
    }
}
```

## 🌍 Internacionalización (i18n)

Los mensajes están preparados para internacionalización:

```go
// i18n.go
package messages

import "fmt"

// MessageKey tipo para claves de mensajes
type MessageKey string

// Claves de mensajes para i18n
const (
    KeyUserCreated     MessageKey = "user.created"
    KeyUserUpdated     MessageKey = "user.updated"
    KeyUserDeleted     MessageKey = "user.deleted"
    KeyUserNotFound    MessageKey = "user.not_found"
    KeyUserEmailExists MessageKey = "user.email_exists"
)

// MessageTranslations traducciones de mensajes
var MessageTranslations = map[string]map[MessageKey]string{
    "en": {
        KeyUserCreated:     "User created successfully",
        KeyUserUpdated:     "User updated successfully", 
        KeyUserDeleted:     "User deleted successfully",
        KeyUserNotFound:    "User not found",
        KeyUserEmailExists: "User with this email already exists",
    },
    "es": {
        KeyUserCreated:     "Usuario creado exitosamente",
        KeyUserUpdated:     "Usuario actualizado exitosamente",
        KeyUserDeleted:     "Usuario eliminado exitosamente", 
        KeyUserNotFound:    "Usuario no encontrado",
        KeyUserEmailExists: "Ya existe un usuario con este email",
    },
    "fr": {
        KeyUserCreated:     "Utilisateur créé avec succès",
        KeyUserUpdated:     "Utilisateur mis à jour avec succès",
        KeyUserDeleted:     "Utilisateur supprimé avec succès",
        KeyUserNotFound:    "Utilisateur non trouvé", 
        KeyUserEmailExists: "Un utilisateur avec cet email existe déjà",
    },
}

// GetMessage obtiene un mensaje traducido
func GetMessage(key MessageKey, lang string) string {
    if translations, exists := MessageTranslations[lang]; exists {
        if message, exists := translations[key]; exists {
            return message
        }
    }
    
    // Fallback a inglés
    if translations, exists := MessageTranslations["en"]; exists {
        if message, exists := translations[key]; exists {
            return message
        }
    }
    
    return string(key)
}

// GetMessageWithParams obtiene un mensaje con parámetros
func GetMessageWithParams(key MessageKey, lang string, params ...interface{}) string {
    message := GetMessage(key, lang)
    if len(params) > 0 {
        return fmt.Sprintf(message, params...)
    }
    return message
}
```

## 📊 Uso en Diferentes Capas

### En Handlers HTTP
```go
func (h *UserHandler) Create(c *gin.Context) {
    user, err := h.userUseCase.Create(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }
    
    response := SuccessResponse(SuccessUserCreated, user)
    c.JSON(http.StatusCreated, response)
}
```

### En Casos de Uso
```go
func (s *UserService) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
    if req.Name == "" {
        return nil, ErrUserNameRequired
    }
    
    existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, ErrUserEmailExists
    }
    
    // ... lógica de creación
    
    return userResponse, nil
}
```

### En Validaciones
```go
func ValidateUser(user *domain.User) *MultiError {
    multiErr := NewMultiError()
    
    if user.Name == "" {
        multiErr.Add(ErrUserNameRequired.WithField("name"))
    }
    
    if len(user.Name) < UserNameMinLength {
        multiErr.Add(ErrUserNameTooShort.WithField("name").WithValue(user.Name))
    }
    
    if !isValidEmail(user.Email) {
        multiErr.Add(ErrUserEmailInvalid.WithField("email").WithValue(user.Email))
    }
    
    if multiErr.HasErrors() {
        return multiErr
    }
    
    return nil
}
```

## ⚠️ Consideraciones Importantes

### ✅ Buenas Prácticas
- **Mensajes consistentes**: Usar templates para mantener consistencia
- **Códigos de error claros**: Usar códigos descriptivos y únicos
- **Internacionalización**: Preparar para múltiples idiomas
- **Metadatos útiles**: Incluir información contextual en errores

### ❌ Errores Comunes
- **Mensajes hardcodeados**: Centralizar todos los mensajes
- **Información sensible**: No exponer datos internos en errores
- **Mensajes genéricos**: Ser específico sobre el problema
- **Inconsistencia**: Mantener formato consistente

### 🔄 Integración Recomendada

#### Logger
```go
import "github.com/sirupsen/logrus"

func LogUserError(err error, userID uint) {
    if userErr, ok := err.(*UserError); ok {
        logrus.WithFields(logrus.Fields{
            "error_type": userErr.Type,
            "user_id":    userID,
            "field":      userErr.Field,
            "value":      userErr.Value,
        }).Error(userErr.Message)
    }
}
```

#### Metrics
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    userErrorCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "user_errors_total",
            Help: "Total number of user errors",
        },
        []string{"error_type"},
    )
)

func TrackUserError(err error) {
    if userErr, ok := err.(*UserError); ok {
        userErrorCounter.WithLabelValues(userErr.Type).Inc()
    }
}
```

---

**← [Comando goca interfaces](Command-Interfaces) | [Comando goca version](Command-Version) →**
