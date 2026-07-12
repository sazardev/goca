---
layout: doc
title: goca messages
titleTemplate: Commands | Goca
description: Generate typed error messages, response structs, and shared constants for a feature's communication layer.
---

# goca messages

Generate error messages, responses, and constants for features.

## Syntax

```bash
goca messages <EntityName> [flags]
```

## Description

Creates centralized message files for errors, success responses, and feature-specific constants.

## Flags

### `--errors`

Generate error messages.

```bash
goca messages User --errors
```

### `--responses`

Generate response messages.

```bash
goca messages Product --responses
```

### `--constants` / `-c`

Generate feature constants.

```bash
goca messages Order --constants
```

### `--all` / `-a`

Generate errors, responses and constants together.

```bash
goca messages Order --all
```

## Examples

### Error Messages

```bash
goca messages User --errors
```

**Generates/appends to:** `internal/messages/messages.go` — shared across every entity in the project (not a per-entity file); running this for multiple entities keeps adding to the same file rather than creating `user_errors.go`, `product_errors.go`, etc.

```go
package messages

import "errors"

var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidUserData   = errors.New("invalid user data")
    ErrUserUnauthorized  = errors.New("user unauthorized")
)
```

### Response Messages

```bash
goca messages Product --responses
```

**Generates/appends to:** `internal/messages/responses.go` (shared across every entity)

```go
package messages

const (
    ProductCreatedSuccess = "Product created successfully"
    ProductUpdatedSuccess = "Product updated successfully"
    ProductDeletedSuccess = "Product deleted successfully"
    ProductListSuccess    = "Products retrieved successfully"
)
```

### Constants

```bash
goca messages Order --constants
```

**Generates/appends to:** `internal/constants/constants.go` (note: a separate `internal/constants` package, not `internal/messages`; shared across every entity)

```go
package messages

const (
    OrderStatusPending    = "pending"
    OrderStatusProcessing = "processing"
    OrderStatusCompleted  = "completed"
    OrderStatusCancelled  = "cancelled"
    
    DefaultPageSize = 20
    MaxPageSize     = 100
)
```

## See Also

- [`goca entity`](/commands/entity) - Generate entities
- [`goca feature`](/commands/feature) - Generate complete feature
