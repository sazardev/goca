---
layout: doc
title: goca test-integration
titleTemplate: Commands | Goca
description: Generate integration test scaffolding for a named entity, including test suite setup, database fixtures, and HTTP request helpers.
---

# goca test-integration

Generate integration test scaffolding for a named entity. The generated tests exercise the full request-to-database stack in a temporary test environment.

## Syntax

```bash
goca test-integration <EntityName> [flags]
```

## Description

`goca test-integration` creates an integration test file under `internal/testing/integration/` that:

- Spins up a test database via `setupTestDatabase`/`cleanupTestDatabase` helpers (in a shared `helpers.go`)
- Runs the full use-case and repository stack
- Uses a plain `func Test<Entity>Integration(t *testing.T)` with `t.Run` subtests — **not** `testify/suite` (no `Suite`/`SetupSuite`/`TearDownSuite`)

## Flags

### `--dry-run`

Preview the files that would be created without writing anything to disk.

```bash
goca test-integration Product --dry-run
```

### `--force`

Overwrite existing integration test files.

```bash
goca test-integration Product --force
```

### `--backup`

Back up existing test files to `.goca-backup/` before overwriting.

```bash
goca test-integration Product --backup
```

## Examples

### Generate integration tests for one entity

```bash
goca test-integration Order
```

### Preview without writing

```bash
goca test-integration Order --dry-run
```

## Generated Files

| File | Description |
| --- | --- |
| `internal/testing/integration/order_integration_test.go` | Integration test for the entity |
| `internal/testing/integration/fixtures/order_fixtures.go` | Test data fixtures |
| `internal/testing/integration/helpers.go` | Shared `setupTestDatabase`/`cleanupTestDatabase` helpers (written once, reused by every entity's test) |

## Generated Code Example

```go
package integration

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "myproject/internal/repository"
    "myproject/internal/usecase"
)

// TestOrderIntegration tests the complete Order feature integration
func TestOrderIntegration(t *testing.T) {
    db := setupTestDatabase(t, "postgres")
    defer cleanupTestDatabase(t, db)

    repo := repository.NewPostgresOrderRepository(db)
    service := usecase.NewOrderService(repo)

    t.Run("CreateAndRetrieveOrder", func(t *testing.T) {
        input := usecase.CreateOrderInput{
            // TODO: Add fields based on entity structure
        }
        output, err := service.CreateOrder(input)
        require.NoError(t, err)
        require.NotNil(t, output)
        assert.NotZero(t, output.ID)
    })
}
```

## Related Commands

- [`goca feature`](/commands/feature) — generate all layers including integration tests via `--integration-tests` flag
- [`goca mocks`](/commands/mocks) — generate unit test mocks for the same entity
