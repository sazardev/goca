# goca middleware

Generate a dedicated middleware package for HTTP handlers.

## Syntax

```bash
goca middleware <name> [flags]
```

## Description

The `goca middleware` command generates an `internal/middleware/` package containing composable HTTP middleware functions. Each middleware follows the standard `func(http.Handler) http.Handler` signature and can be chained together using the generated `Chain()` helper.

::: tip Composable Design
Generated middleware is designed to work with any `net/http` compatible router (gorilla/mux, chi, standard library). Use `middleware.Chain()` to compose multiple middleware in a single call.
:::

## Supported Middleware Types

| Type | Function | Description |
| ---- | -------- | ----------- |
| `cors` | `CORS(cfg CORSConfig)` | Configurable CORS headers (origins, methods, headers) |
| `logging` | `Logging()` | Structured request logging with method, path, status, duration |
| `auth` | `Auth()` | JWT Bearer token validation with claims extraction into context |
| `rate-limit` | `RateLimit(cfg RateLimitConfig)` | Per-IP token bucket rate limiting |
| `recovery` | `Recovery(debugMode bool)` | Panic recovery returning JSON 500 with optional stack trace |
| `request-id` | `RequestID()` | Inject `X-Request-ID` into context and response headers |
| `timeout` | `Timeout(d time.Duration)` | Per-request context deadline |

## Generated Files

```
internal/middleware/
├── middleware.go      # Middleware type + Chain() helper (always)
├── cors.go            # CORS middleware
├── logging.go         # Request logging middleware
├── auth.go            # JWT authentication middleware
├── rate_limit.go      # Rate limiting middleware
├── recovery.go        # Panic recovery middleware
├── request_id.go      # Request ID middleware
└── timeout.go         # Request timeout middleware
```

Only the types specified by `--types` are generated. The `middleware.go` chain helper is always included.

## Flags

### `--types`

Comma-separated list of middleware types to generate.

```bash
goca middleware MyApp --types cors,logging,auth,recovery
```

**Default:** `cors,logging,recovery`
**Supported values:** `cors`, `logging`, `auth`, `rate-limit`, `recovery`, `request-id`, `timeout`

### `--dry-run`

Preview what files would be generated without writing anything to disk.

```bash
goca middleware MyApp --dry-run
```

### `--force`

Overwrite existing middleware files without prompting.

```bash
goca middleware MyApp --force
```

### `--backup`

Create backups of existing files before overwriting.

```bash
goca middleware MyApp --backup
```

## Usage Examples

### Default middleware (CORS + Logging + Recovery)

```bash
goca middleware MyApp
```

Generates `middleware.go`, `cors.go`, `logging.go`, and `recovery.go`.

### All middleware types

```bash
goca middleware MyApp --types cors,logging,auth,rate-limit,recovery,request-id,timeout
```

### Auth-focused setup

```bash
goca middleware MyApp --types auth,rate-limit,cors,logging
```

### Preview before generating

```bash
goca middleware MyApp --types cors,logging,auth --dry-run
```

## Using Generated Middleware

### Chaining middleware with a router

```go
package main

import (
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "myproject/internal/middleware"
)

func main() {
    r := mux.NewRouter()

    // Compose middleware
    stack := middleware.Chain(
        middleware.Recovery(true),
        middleware.Logging(),
        middleware.CORS(middleware.DefaultCORSConfig()),
        middleware.Timeout(30 * time.Second),
    )

    // Apply to router
    http.ListenAndServe(":8080", stack(r))
}
```

### Accessing JWT claims in a handler

```go
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.ClaimsFromContext(r.Context())
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }
    userID := claims["sub"].(string)
    // ...
}
```

### Accessing the request ID

```go
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    reqID := middleware.RequestIDFromContext(r.Context())
    log.Printf("[%s] creating resource", reqID)
    // ...
}
```

### Custom CORS configuration

```go
corsConfig := middleware.CORSConfig{
    AllowedOrigins: []string{"https://myapp.com", "https://admin.myapp.com"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders: []string{"Content-Type", "Authorization"},
}
middleware.CORS(corsConfig)
```

## External Dependencies

Some middleware types require external packages. Goca does **not** automatically run `go get` for these — add them to your project as needed:

| Type | Dependency |
| ---- | ---------- |
| `auth` | `github.com/golang-jwt/jwt/v5` |
| `rate-limit` | `golang.org/x/time` |
| `request-id` | `github.com/google/uuid` |

## Integration with Other Commands

| Scenario | Related Command |
| -------- | --------------- |
| Generate handlers that use middleware | `goca handler <entity> --middleware` |
| Generate a full feature with handlers | `goca feature <entity> --fields "..."` |
| Wire middleware into DI container | `goca di` |

## See Also

- [Commands Overview](/goca/commands/)
- [goca handler](/goca/commands/handler)
- [goca feature](/goca/commands/feature)
