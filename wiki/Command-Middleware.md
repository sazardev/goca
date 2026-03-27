# goca middleware Command

The `goca middleware` command generates composable HTTP middleware for your application. Middleware can be generated individually or as a full set, and the handler layer auto-detects and chains them.

## 📋 Syntax

```bash
goca middleware <name> [flags]
```

## 🎯 Purpose

Creates middleware components for HTTP request processing:

- 🔒 **Auth** — JWT/Bearer token authentication
- 🌐 **CORS** — Cross-Origin Resource Sharing headers
- 📝 **Logging** — Structured request/response logging
- ⚡ **Rate Limit** — Token-bucket rate limiting per IP
- 🛡️ **Recovery** — Panic recovery with stack traces
- 🆔 **Request ID** — UUID request tracing header
- ⏱️ **Timeout** — Context-based request timeouts

## 🚩 Available Flags

| Flag       | Type     | Required  | Default Value | Description                                           |
| ---------- | -------- | --------- | ------------- | ----------------------------------------------------- |
| `--types`  | `string` | ❌ No      | `all`         | Comma-separated middleware types to generate           |
| `--dry-run`| `bool`   | ❌ No      | `false`       | Preview files without writing                         |
| `--force`  | `bool`   | ❌ No      | `false`       | Overwrite existing files                              |
| `--backup` | `bool`   | ❌ No      | `false`       | Backup existing files before overwriting              |

### Valid Middleware Types

`cors`, `logging`, `auth`, `rate-limit`, `recovery`, `request-id`, `timeout`

## 📖 Usage Examples

### Generate All Middleware
```bash
goca middleware myapp
```
Generates all 7 middleware types plus a `chain.go` for composing them.

### Generate Specific Types
```bash
goca middleware myapp --types "cors,logging,auth"
```

### Auth + Rate Limiting Only
```bash
goca middleware myapp --types "auth,rate-limit"
```

### Preview Before Generating
```bash
goca middleware myapp --types "cors,logging" --dry-run
```

## 📁 Generated Files

```
internal/
└── middleware/
    ├── chain.go        # Always generated — Middleware type + Chain()
    ├── cors.go         # CORS middleware
    ├── logging.go      # Logging middleware
    ├── auth.go         # Auth middleware
    ├── rate_limit.go   # Rate limit middleware
    ├── recovery.go     # Recovery middleware
    ├── request_id.go   # Request ID middleware
    └── timeout.go      # Timeout middleware
```

### chain.go

Defines the `Middleware` type alias and `Chain()` function:

```go
package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
    return func(next http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            next = middlewares[i](next)
        }
        return next
    }
}
```

### Handler Integration

When a middleware package exists, `goca handler` automatically detects it and adds middleware imports and chaining to generated handlers.

## 🔗 Related Commands

- [`goca handler`](Command-Handler.md) — Generate HTTP handlers (auto-chains middleware)
- [`goca feature`](Command-Feature.md) — Generate complete feature stack with `--middleware-types`
