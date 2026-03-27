# goca ci

Generate CI/CD pipeline configuration for your project.

## Syntax

```bash
goca ci [flags]
```

## Description

The `goca ci` command generates Continuous Integration pipeline files tailored to your project. It reads `go.mod` to detect the Go version and `.goca.yaml` to detect the database driver, then produces provider-specific workflow files.

Currently supported providers:

- **github-actions** — GitHub Actions workflows (`.github/workflows/`)

::: tip Auto-Detection
`goca ci` automatically reads your `go.mod` for the Go version and `.goca.yaml` for the database driver. No manual configuration is needed in most cases.
:::

## Generated Workflows

### test.yml

Runs on every push and pull request:

- Checks out code
- Sets up Go (version from `go.mod`)
- Starts database service container (PostgreSQL or MySQL) when detected in `.goca.yaml`
- Runs `go vet ./...`
- Runs `go test -race -coverprofile=coverage.out ./...`
- Runs `go build ./...`

### build.yml

Runs on every push to `main` and on pull requests:

- Checks out code
- Sets up Go
- Builds the binary (`go build -o bin/app ./...`)
- Uploads `bin/` as a build artifact
- Optionally builds a Docker image (when `--with-docker` is used)

### deploy.yml (optional)

Generated when `--with-deploy` is provided. Runs on tag pushes matching `v*`:

- Checks out code
- Sets up Go
- Builds the binary
- Placeholder deploy step (customize to your target)

## Flags

### `--provider`

CI provider to generate configuration for.

```bash
goca ci --provider github-actions
```

**Default:** `github-actions`
**Supported values:** `github-actions`

### `--with-docker`

Include a Docker image build step in the build workflow.

```bash
goca ci --with-docker
```

### `--with-deploy`

Generate an additional deploy workflow triggered by version tags.

```bash
goca ci --with-deploy
```

### `--go-version`

Override the Go version used in CI. By default, Goca reads the version from your `go.mod` file.

```bash
goca ci --go-version 1.25
```

### `--dry-run`

Preview what files would be generated without writing anything to disk.

```bash
goca ci --dry-run
```

### `--force`

Overwrite existing workflow files without prompting.

```bash
goca ci --force
```

### `--backup`

Create backups of existing workflow files before overwriting.

```bash
goca ci --backup
```

## Usage Examples

### Basic CI generation

```bash
goca ci
```

Creates `.github/workflows/test.yml` and `.github/workflows/build.yml`.

### Full pipeline with Docker and deploy

```bash
goca ci --with-docker --with-deploy
```

Creates `test.yml`, `build.yml` (with Docker step), and `deploy.yml`.

### Preview before generating

```bash
goca ci --dry-run
```

**Example output:**

```
Goca CI — Pipeline Generation

  Provider:   github-actions
  Go version: 1.25
  Docker:     true
  Deploy:     true
  Database:   postgres

Dry-Run Summary
┌──────────────────────────────────────────────────┐
│  Action   Path                                    │
├──────────────────────────────────────────────────┤
│  CREATE   .github/workflows/test.yml              │
│  CREATE   .github/workflows/build.yml             │
│  CREATE   .github/workflows/deploy.yml            │
└──────────────────────────────────────────────────┘
```

### With explicit Go version

```bash
goca ci --go-version 1.24
```

### Database-Aware CI

When your `.goca.yaml` specifies a database driver, `goca ci` automatically adds the corresponding service container to `test.yml`:

**PostgreSQL:**

```yaml
services:
  postgres:
    image: postgres:16
    env:
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test
      POSTGRES_DB: testdb
    ports:
      - 5432:5432
    options: >-
      --health-cmd pg_isready
      --health-interval 10s
      --health-timeout 5s
      --health-retries 5
```

**MySQL:**

```yaml
services:
  mysql:
    image: mysql:8
    env:
      MYSQL_ROOT_PASSWORD: test
      MYSQL_DATABASE: testdb
    ports:
      - 3306:3306
    options: >-
      --health-cmd "mysqladmin ping"
      --health-interval 10s
      --health-timeout 5s
      --health-retries 5
```

Environment variables `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, and `DB_NAME` are set automatically for the test job.

## Integration with Other Commands

| Scenario | Related Command |
| -------- | --------------- |
| Initialize a project first | `goca init <project>` |
| Check project health | `goca doctor` |
| Generate integration tests | `goca test-integration` |

## See Also

- [Commands Overview](/goca/commands/)
- [goca init](/goca/commands/init)
- [goca doctor](/goca/commands/doctor)
- [goca test-integration](/goca/commands/test-integration)
