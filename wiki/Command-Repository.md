# Comando goca repository

El comando `goca repository` crea repositorios que implementan el patrón Repository con interfaces bien definidas e implementaciones específicas por base de datos siguiendo Clean Architecture.

## 📋 Sintaxis

```bash
goca repository <entity> [flags]
```

## 🎯 Propósito

Crea repositorios para manejar la persistencia de entidades:

- 🔵 **Abstracción de persistencia** sin acoplar el dominio
- 📊 **Implementaciones específicas** por base de datos
- 🔗 **Interfaces claras** para casos de uso
- 💾 **Transacciones** y manejo de errores
- ⚡ **Caché** opcional para optimización
- 🔍 **Queries optimizadas** por tecnología

## 🚩 Flags Disponibles

| Flag               | Tipo     | Requerido | Valor por Defecto | Descripción                                            |
| ------------------ | -------- | --------- | ----------------- | ------------------------------------------------------ |
| `--database`       | `string` | ❌ No      | -                 | Tipo de base de datos (`postgres`, `mysql`, `mongodb`) |
| `--interface-only` | `bool`   | ❌ No      | `false`           | Solo generar interfaces                                |
| `--implementation` | `bool`   | ❌ No      | `false`           | Solo generar implementación                            |
| `--transactions`   | `bool`   | ❌ No      | `false`           | Incluir soporte para transacciones                     |
| `--cache`          | `bool`   | ❌ No      | `false`           | Incluir capa de caché                                  |

## 📖 Ejemplos de Uso

### Repositorio Básico con PostgreSQL
```bash
goca repository User --database postgres
```

### Solo Generar Interfaces
```bash
goca repository Product --interface-only
```

### Con Transacciones y Caché
```bash
goca repository Order --database postgres --transactions --cache
```

### Diferentes Bases de Datos
```bash
# PostgreSQL
goca repository User --database postgres --transactions

# MySQL
goca repository Product --database mysql --cache

# MongoDB
goca repository Order --database mongodb
```

## 📂 Archivos Generados

### Estructura de Archivos
```
internal/repository/
├── interfaces/
│   └── user_repository.go      # Interface del repositorio
├── postgres/
│   └── user_repository.go      # Implementación PostgreSQL
├── mysql/
│   └── user_repository.go      # Implementación MySQL (si se especifica)
└── mongodb/
    └── user_repository.go      # Implementación MongoDB (si se especifica)
```

## 🔍 Código Generado en Detalle

### Interface: `internal/repository/interfaces/user_repository.go`

```go
package interfaces

import (
    "context"
    
    "github.com/usuario/proyecto/internal/domain"
)

//go:generate mockgen -source=user_repository.go -destination=mocks/user_repository_mock.go

// UserRepository define los contratos para la persistencia de usuarios
type UserRepository interface {
    // Operaciones básicas CRUD
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id uint) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uint) error
    
    // Operaciones de consulta
    List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
    Search(ctx context.Context, query string, limit, offset int) ([]*domain.User, int64, error)
    Exists(ctx context.Context, id uint) (bool, error)
    
    // Operaciones de transacción (si --transactions)
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
    
    // Operaciones de caché (si --cache)
    ClearCache(ctx context.Context, id uint) error
}
```

### Implementación PostgreSQL: `internal/repository/postgres/user_repository.go`

```go
package postgres

import (
    "context"
    "database/sql"
    "fmt"
    
    "github.com/lib/pq"
    "github.com/usuario/proyecto/internal/domain"
    "github.com/usuario/proyecto/internal/repository/interfaces"
)

// UserRepository implementa interfaces.UserRepository para PostgreSQL
type UserRepository struct {
    db *sql.DB
}

// NewUserRepository crea una nueva instancia del repositorio
func NewUserRepository(db *sql.DB) interfaces.UserRepository {
    return &UserRepository{
        db: db,
    }
}

// Save guarda un usuario en la base de datos
func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (name, email, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
    
    err := r.db.QueryRowContext(
        ctx,
        query,
        user.Name,
        user.Email,
        user.CreatedAt,
        user.UpdatedAt,
    ).Scan(&user.ID)
    
    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok {
            switch pqErr.Code {
            case "23505": // unique_violation
                return domain.ErrUserEmailAlreadyExists
            }
        }
        return fmt.Errorf("failed to save user: %w", err)
    }
    
    return nil
}

// FindByID busca un usuario por ID
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `
    
    user := &domain.User{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, domain.ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to find user by ID: %w", err)
    }
    
    return user, nil
}

// FindByEmail busca un usuario por email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE email = $1 AND deleted_at IS NULL
    `
    
    user := &domain.User{}
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // No error si no existe
        }
        return nil, fmt.Errorf("failed to find user by email: %w", err)
    }
    
    return user, nil
}

// Update actualiza un usuario existente
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
    query := `
        UPDATE users
        SET name = $2, email = $3, updated_at = $4
        WHERE id = $1 AND deleted_at IS NULL
    `
    
    result, err := r.db.ExecContext(
        ctx,
        query,
        user.ID,
        user.Name,
        user.Email,
        user.UpdatedAt,
    )
    
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return domain.ErrUserNotFound
    }
    
    return nil
}

// Delete elimina un usuario (soft delete)
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
    query := `
        UPDATE users
        SET deleted_at = NOW()
        WHERE id = $1 AND deleted_at IS NULL
    `
    
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return domain.ErrUserNotFound
    }
    
    return nil
}

// List obtiene una lista paginada de usuarios
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
    // Contar total
    countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
    var total int64
    err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count users: %w", err)
    }
    
    // Obtener usuarios paginados
    query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
    
    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to list users: %w", err)
    }
    defer rows.Close()
    
    var users []*domain.User
    for rows.Next() {
        user := &domain.User{}
        err := rows.Scan(
            &user.ID,
            &user.Name,
            &user.Email,
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan user: %w", err)
        }
        users = append(users, user)
    }
    
    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("failed to iterate users: %w", err)
    }
    
    return users, total, nil
}
```

## 💾 Implementaciones por Base de Datos

### PostgreSQL Features
- **JSONB** para campos complejos
- **Array types** para listas
- **UUID** como primary keys
- **Partial indexes** para soft delete
- **RETURNING** clause para obtener IDs

### MySQL Features
- **JSON** para campos complejos
- **Generated columns** para campos calculados
- **Multi-value indexes** para búsquedas
- **Foreign key constraints** con CASCADE

### MongoDB Features
- **Agregación pipeline** para queries complejas
- **Índices compuestos** para optimización
- **GridFS** para archivos grandes
- **Transacciones** en replica sets

## 🔄 Transacciones (--transactions)

Con `--transactions`, se agregan métodos para manejo de transacciones:

```go
// WithTransaction ejecuta operaciones dentro de una transacción
func (r *UserRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    // Crear contexto con transacción
    txCtx := context.WithValue(ctx, "tx", tx)
    
    if err := fn(txCtx); err != nil {
        if rollbackErr := tx.Rollback(); rollbackErr != nil {
            return fmt.Errorf("failed to rollback transaction: %v, original error: %w", rollbackErr, err)
        }
        return err
    }
    
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}

// SaveWithTx guarda usando transacción existente
func (r *UserRepository) SaveWithTx(ctx context.Context, user *domain.User) error {
    var executor interface {
        QueryRowContext(context.Context, string, ...interface{}) *sql.Row
    } = r.db
    
    if tx, ok := ctx.Value("tx").(*sql.Tx); ok {
        executor = tx
    }
    
    query := `
        INSERT INTO users (name, email, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
    
    return executor.QueryRowContext(
        ctx, query,
        user.Name, user.Email, user.CreatedAt, user.UpdatedAt,
    ).Scan(&user.ID)
}
```

## ⚡ Caché (--cache)

Con `--cache`, se integra una capa de caché:

```go
import (
    "encoding/json"
    "time"
    
    "github.com/go-redis/redis/v8"
)

type CachedUserRepository struct {
    repo  interfaces.UserRepository
    cache *redis.Client
    ttl   time.Duration
}

func NewCachedUserRepository(repo interfaces.UserRepository, cache *redis.Client) interfaces.UserRepository {
    return &CachedUserRepository{
        repo:  repo,
        cache: cache,
        ttl:   time.Hour,
    }
}

// FindByID con caché
func (r *CachedUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
    cacheKey := fmt.Sprintf("user:id:%d", id)
    
    // Intentar obtener del caché
    cached, err := r.cache.Get(ctx, cacheKey).Result()
    if err == nil {
        var user domain.User
        if err := json.Unmarshal([]byte(cached), &user); err == nil {
            return &user, nil
        }
    }
    
    // Si no está en caché, obtener de DB
    user, err := r.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Guardar en caché
    if user != nil {
        if userData, err := json.Marshal(user); err == nil {
            r.cache.Set(ctx, cacheKey, userData, r.ttl)
        }
    }
    
    return user, nil
}

// ClearCache limpia el caché de un usuario
func (r *CachedUserRepository) ClearCache(ctx context.Context, id uint) error {
    patterns := []string{
        fmt.Sprintf("user:id:%d", id),
        "user:list:*",
        "user:search:*",
    }
    
    for _, pattern := range patterns {
        keys, err := r.cache.Keys(ctx, pattern).Result()
        if err != nil {
            continue
        }
        
        if len(keys) > 0 {
            r.cache.Del(ctx, keys...)
        }
    }
    
    return nil
}
```

## 🔍 Queries Avanzadas

### Búsqueda con Full-Text Search
```go
// Search busca usuarios por texto
func (r *UserRepository) Search(ctx context.Context, query string, limit, offset int) ([]*domain.User, int64, error) {
    searchQuery := `
        SELECT id, name, email, created_at, updated_at,
               ts_rank_cd(search_vector, plainto_tsquery($1)) as rank
        FROM users
        WHERE search_vector @@ plainto_tsquery($1)
          AND deleted_at IS NULL
        ORDER BY rank DESC, created_at DESC
        LIMIT $2 OFFSET $3
    `
    
    rows, err := r.db.QueryContext(ctx, searchQuery, query, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to search users: %w", err)
    }
    defer rows.Close()
    
    var users []*domain.User
    for rows.Next() {
        user := &domain.User{}
        var rank float64
        
        err := rows.Scan(
            &user.ID, &user.Name, &user.Email,
            &user.CreatedAt, &user.UpdatedAt, &rank,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan search result: %w", err)
        }
        
        users = append(users, user)
    }
    
    // Contar resultados
    countQuery := `
        SELECT COUNT(*)
        FROM users
        WHERE search_vector @@ plainto_tsquery($1) AND deleted_at IS NULL
    `
    
    var total int64
    err = r.db.QueryRowContext(ctx, countQuery, query).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count search results: %w", err)
    }
    
    return users, total, nil
}
```

## 🧪 Testing

Los repositorios generados incluyen interfaces para fácil testing:

```go
func TestUserRepository_Save(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := postgres.NewUserRepository(db)
    
    user := &domain.User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    err := repo.Save(context.Background(), user)
    
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}

func TestUserRepository_WithTransaction(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := postgres.NewUserRepository(db)
    
    err := repo.WithTransaction(context.Background(), func(ctx context.Context) error {
        user1 := &domain.User{Name: "User 1", Email: "user1@test.com"}
        user2 := &domain.User{Name: "User 2", Email: "user2@test.com"}
        
        if err := repo.SaveWithTx(ctx, user1); err != nil {
            return err
        }
        
        return repo.SaveWithTx(ctx, user2)
    })
    
    assert.NoError(t, err)
}
```

## ⚠️ Consideraciones Importantes

### ✅ Buenas Prácticas
- **Context propagation**: Usar context.Context en todos los métodos
- **Error wrapping**: Envolver errores con información contextual
- **Prepared statements**: Prevenir SQL injection
- **Connection pooling**: Configurar pools adecuadamente

### ❌ Errores Comunes
- **No usar transacciones**: Para operaciones que requieren consistencia
- **Ignorar errores**: Siempre manejar errores apropiadamente
- **Queries N+1**: Optimizar con JOINs o eager loading
- **No limpiar recursos**: Cerrar rows, statements, etc.

---

**← [Comando goca usecase](Command-UseCase) | [Comando goca handler](Command-Handler) →**
