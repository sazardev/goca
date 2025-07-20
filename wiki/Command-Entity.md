# Comando goca entity

El comando `goca entity` genera entidades de dominio puras siguiendo los principios de Domain-Driven Design (DDD), sin dependencias externas y con validaciones de negocio incorporadas.

## 📋 Sintaxis

```bash
goca entity <name> [flags]
```

## 🎯 Propósito

Crea entidades de dominio que representan los conceptos centrales del negocio:

- 🟡 **Entidad pura** sin dependencias externas
- ✅ **Validaciones de negocio** incorporadas
- 🔧 **Reglas de negocio** específicas del dominio
- ⏰ **Timestamps** automáticos (opcional)
- 🗑️ **Soft delete** (opcional)
- 📄 **Errores específicos** del dominio

## 🚩 Flags Disponibles

| Flag               | Tipo     | Requerido | Valor por Defecto | Descripción                                         |
| ------------------ | -------- | --------- | ----------------- | --------------------------------------------------- |
| `--fields`         | `string` | ✅ **Sí**  | -                 | Campos de la entidad (`"name:string,email:string"`) |
| `--validation`     | `bool`   | ❌ No      | `false`           | Incluir validaciones automáticas                    |
| `--business-rules` | `bool`   | ❌ No      | `false`           | Generar métodos de reglas de negocio                |
| `--timestamps`     | `bool`   | ❌ No      | `false`           | Agregar campos `created_at` y `updated_at`          |
| `--soft-delete`    | `bool`   | ❌ No      | `false`           | Agregar campo `deleted_at` para soft delete         |

## 📖 Ejemplos de Uso

### Entidad Básica
```bash
goca entity User --fields "name:string,email:string,age:int"
```

### Entidad con Validaciones
```bash
goca entity Product --fields "name:string,price:float64,category:string" --validation
```

### Entidad Completa
```bash
goca entity Order \
  --fields "user_id:int,total:float64,status:string" \
  --validation \
  --business-rules \
  --timestamps \
  --soft-delete
```

### Entidad con Tipos Complejos
```bash
goca entity Employee \
  --fields "name:string,email:string,salary:float64,hire_date:time.Time,department:string,is_active:bool" \
  --validation \
  --business-rules \
  --timestamps
```

## 📂 Archivos Generados

### Estructura de Archivos
```
internal/domain/
├── user.go              # Entidad principal
├── errors.go            # Errores específicos (si --validation)
└── validations.go       # Validaciones reutilizables (si --business-rules)
```

## 🔍 Código Generado en Detalle

### Entidad Básica: `internal/domain/user.go`

```go
package domain

import (
    "errors"
    "strings"
    "time"
)

// User representa un usuario en el sistema
type User struct {
    ID        uint      `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Email     string    `json:"email" db:"email"`
    Age       int       `json:"age" db:"age"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewUser crea una nueva instancia de User
func NewUser(name, email string, age int) *User {
    return &User{
        Name:      name,
        Email:     email,
        Age:       age,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

// TableName retorna el nombre de la tabla en la base de datos
func (User) TableName() string {
    return "users"
}
```

### Con Validaciones: `--validation`

```go
// Validate valida la entidad User según las reglas de negocio
func (u *User) Validate() error {
    if strings.TrimSpace(u.Name) == "" {
        return ErrUserNameRequired
    }
    
    if len(u.Name) < 2 {
        return ErrUserNameTooShort
    }
    
    if len(u.Name) > 100 {
        return ErrUserNameTooLong
    }
    
    if strings.TrimSpace(u.Email) == "" {
        return ErrUserEmailRequired
    }
    
    if !u.isValidEmail() {
        return ErrUserEmailInvalid
    }
    
    if u.Age < 0 || u.Age > 150 {
        return ErrUserAgeInvalid
    }
    
    return nil
}

// isValidEmail valida el formato del email
func (u *User) isValidEmail() bool {
    return strings.Contains(u.Email, "@") && 
           strings.Contains(u.Email, ".") &&
           len(u.Email) >= 5
}

// BeforeSave ejecuta validaciones antes de guardar
func (u *User) BeforeSave() error {
    u.UpdatedAt = time.Now()
    return u.Validate()
}

// BeforeCreate ejecuta validaciones antes de crear
func (u *User) BeforeCreate() error {
    now := time.Now()
    u.CreatedAt = now
    u.UpdatedAt = now
    return u.Validate()
}
```

### Con Reglas de Negocio: `--business-rules`

```go
// Business Rules

// CanUpdateProfile verifica si el usuario puede actualizar su perfil
func (u *User) CanUpdateProfile() bool {
    return u.ID > 0 && u.Name != ""
}

// CanChangeEmail verifica si el usuario puede cambiar su email
func (u *User) CanChangeEmail() bool {
    return u.ID > 0
}

// IsAdult verifica si el usuario es mayor de edad
func (u *User) IsAdult() bool {
    return u.Age >= 18
}

// IsEmailDomainAllowed verifica si el dominio del email está permitido
func (u *User) IsEmailDomainAllowed() bool {
    allowedDomains := []string{
        "gmail.com", "yahoo.com", "hotmail.com", 
        "company.com", "organization.org",
    }
    
    parts := strings.Split(u.Email, "@")
    if len(parts) != 2 {
        return false
    }
    
    domain := strings.ToLower(parts[1])
    for _, allowed := range allowedDomains {
        if domain == allowed {
            return true
        }
    }
    
    return false
}

// GetDisplayName retorna el nombre para mostrar
func (u *User) GetDisplayName() string {
    if strings.TrimSpace(u.Name) == "" {
        return "Usuario Anónimo"
    }
    return u.Name
}

// GetInitials retorna las iniciales del nombre
func (u *User) GetInitials() string {
    parts := strings.Fields(u.Name)
    if len(parts) == 0 {
        return "?"
    }
    
    initials := ""
    for _, part := range parts {
        if len(part) > 0 {
            initials += strings.ToUpper(string(part[0]))
        }
    }
    
    return initials
}
```

### Con Timestamps: `--timestamps`

```go
type User struct {
    ID        uint       `json:"id" db:"id"`
    Name      string     `json:"name" db:"name"`
    Email     string     `json:"email" db:"email"`
    Age       int        `json:"age" db:"age"`
    CreatedAt time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// UpdateTimestamp actualiza el timestamp de modificación
func (u *User) UpdateTimestamp() {
    u.UpdatedAt = time.Now()
}

// SetCreationTime establece el tiempo de creación
func (u *User) SetCreationTime() {
    now := time.Now()
    u.CreatedAt = now
    u.UpdatedAt = now
}
```

### Con Soft Delete: `--soft-delete`

```go
type User struct {
    ID        uint       `json:"id" db:"id"`
    Name      string     `json:"name" db:"name"`
    Email     string     `json:"email" db:"email"`
    Age       int        `json:"age" db:"age"`
    CreatedAt time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// IsDeleted verifica si la entidad está marcada como eliminada
func (u *User) IsDeleted() bool {
    return u.DeletedAt != nil
}

// SoftDelete marca la entidad como eliminada
func (u *User) SoftDelete() {
    now := time.Now()
    u.DeletedAt = &now
    u.UpdatedAt = now
}

// Restore restaura una entidad eliminada
func (u *User) Restore() {
    u.DeletedAt = nil
    u.UpdatedAt = time.Now()
}

// CanBeDeleted verifica si la entidad puede ser eliminada
func (u *User) CanBeDeleted() bool {
    return u.ID > 0 && !u.IsDeleted()
}
```

### Errores del Dominio: `internal/domain/errors.go`

```go
package domain

import "errors"

// User domain errors
var (
    // Validation errors
    ErrUserNameRequired    = errors.New("user name is required")
    ErrUserNameTooShort    = errors.New("user name must be at least 2 characters")
    ErrUserNameTooLong     = errors.New("user name must be less than 100 characters")
    ErrUserEmailRequired   = errors.New("user email is required")
    ErrUserEmailInvalid    = errors.New("user email format is invalid")
    ErrUserAgeInvalid      = errors.New("user age must be between 0 and 150")
    
    // Business rule errors
    ErrUserNotFound        = errors.New("user not found")
    ErrUserAlreadyExists   = errors.New("user already exists")
    ErrUserEmailTaken      = errors.New("email address is already taken")
    ErrUserCannotUpdate    = errors.New("user cannot be updated")
    ErrUserCannotDelete    = errors.New("user cannot be deleted")
    ErrUserInactive        = errors.New("user account is inactive")
    ErrUserDeleted         = errors.New("user account has been deleted")
    
    // Authorization errors
    ErrUserUnauthorized    = errors.New("user is not authorized")
    ErrUserInvalidRole     = errors.New("user has invalid role")
    ErrUserPermissionDenied = errors.New("user permission denied")
)

// IsValidationError verifica si el error es de validación
func IsValidationError(err error) bool {
    validationErrors := []error{
        ErrUserNameRequired,
        ErrUserNameTooShort,
        ErrUserNameTooLong,
        ErrUserEmailRequired,
        ErrUserEmailInvalid,
        ErrUserAgeInvalid,
    }
    
    for _, validationErr := range validationErrors {
        if errors.Is(err, validationErr) {
            return true
        }
    }
    
    return false
}

// IsBusinessRuleError verifica si el error es de regla de negocio
func IsBusinessRuleError(err error) bool {
    businessErrors := []error{
        ErrUserNotFound,
        ErrUserAlreadyExists,
        ErrUserEmailTaken,
        ErrUserCannotUpdate,
        ErrUserCannotDelete,
        ErrUserInactive,
        ErrUserDeleted,
    }
    
    for _, businessErr := range businessErrors {
        if errors.Is(err, businessErr) {
            return true
        }
    }
    
    return false
}
```

## 🔧 Tipos de Campos Soportados

### Tipos Básicos de Go
```bash
--fields "name:string,age:int,salary:float64,active:bool"
```

### Tipos de Tiempo
```bash
--fields "created_at:time.Time,birth_date:time.Time"
```

### Tipos Opcionales (Punteros)
```bash
--fields "nickname:*string,last_login:*time.Time"
```

### Slices y Arrays
```bash
--fields "tags:[]string,scores:[]int,metadata:[]byte"
```

### Tipos Personalizados
```bash
--fields "status:UserStatus,role:Role,priority:Priority"
```

### Tipos Numéricos Específicos
```bash
--fields "count:int32,amount:int64,ratio:float32,percentage:float64"
```

## 🎯 Casos de Uso Comunes

### Entidad de Usuario
```bash
goca entity User \
  --fields "name:string,email:string,password:string,role:string,last_login:*time.Time" \
  --validation \
  --business-rules \
  --timestamps
```

### Entidad de Producto
```bash
goca entity Product \
  --fields "name:string,description:string,price:float64,category:string,stock:int,sku:string" \
  --validation \
  --business-rules \
  --timestamps \
  --soft-delete
```

### Entidad de Orden
```bash
goca entity Order \
  --fields "user_id:int,total:float64,status:string,payment_method:string" \
  --validation \
  --business-rules \
  --timestamps
```

### Entidad de Configuración
```bash
goca entity Settings \
  --fields "key:string,value:string,type:string,description:string" \
  --validation \
  --timestamps
```

## ✅ Validaciones Automáticas

Cuando usas `--validation`, se generan automáticamente:

### Para Strings
- **Required**: No vacío después de trim
- **Length**: Mínimo y máximo de caracteres
- **Format**: Email, URL, etc.

### Para Números
- **Range**: Valores mínimos y máximos
- **Positive**: Solo números positivos
- **Non-zero**: No puede ser cero

### Para Fechas
- **Not zero**: No puede ser fecha cero
- **Range**: Entre fechas específicas
- **Future/Past**: Solo fechas futuras o pasadas

### Para Booleanos
- **Explicit**: Debe ser explícitamente true/false

## 🏗️ Reglas de Negocio Automáticas

Con `--business-rules`, se generan métodos como:

### Métodos de Estado
```go
func (e *Entity) IsActive() bool
func (e *Entity) IsValid() bool
func (e *Entity) IsComplete() bool
```

### Métodos de Capacidad
```go
func (e *Entity) CanBeUpdated() bool
func (e *Entity) CanBeDeleted() bool
func (e *Entity) CanPerformAction() bool
```

### Métodos de Transformación
```go
func (e *Entity) GetDisplayName() string
func (e *Entity) GetSummary() string
func (e *Entity) GetInitials() string
```

## 🧪 Testing

Las entidades generadas son fáciles de testear:

```go
func TestUser_Validate(t *testing.T) {
    tests := []struct {
        name    string
        user    domain.User
        wantErr bool
    }{
        {
            name: "valid user",
            user: domain.User{
                Name:  "John Doe",
                Email: "john@example.com",
                Age:   25,
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            user: domain.User{
                Name:  "John Doe",
                Email: "invalid-email",
                Age:   25,
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## 🔄 Integración con Otros Comandos

### Después de Crear Entidad
```bash
# 1. Crear entidad
goca entity User --fields "name:string,email:string" --validation

# 2. Crear caso de uso
goca usecase UserUseCase --entity User --operations "create,read,update,delete"

# 3. Crear repositorio
goca repository User --database postgres

# 4. Crear handler
goca handler User --type http
```

## ⚠️ Consideraciones Importantes

### ✅ Buenas Prácticas
- **Nombres descriptivos**: `User`, `Product`, `Order`
- **Campos específicos**: `email` mejor que `contact`
- **Validaciones consistentes**: Usar siempre `--validation`
- **Reglas de negocio**: Activar para dominios complejos

### ❌ Errores Comunes
- **Entidades anémicas**: Sin comportamiento ni validaciones
- **Dependencias externas**: No agregar imports de DB o HTTP
- **Demasiados campos**: Dividir en entidades más pequeñas
- **Nombres genéricos**: `Data`, `Info`, `Item`

---

**← [Comando goca feature](Command-Feature) | [Comando goca usecase](Command-UseCase) →**
