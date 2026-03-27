# goca ci Command

The `goca ci` command generates CI/CD pipeline configuration files for your project. Currently supports GitHub Actions with automatic Go version detection.

## 📋 Syntax

```bash
goca ci [flags]
```

## 🎯 Purpose

Creates CI/CD pipeline configuration files:

- 🔄 **Test workflow** — runs `go test ./...` on pushes and PRs
- 🏗️ **Build workflow** — compiles the binary and uploads artifacts
- 🚀 **Deploy workflow** (optional) — triggered on tags for release deployments
- 🐳 **Docker support** (optional) — adds Dockerfile build/push steps
- 📊 **Auto-detects Go version** from `go.mod`

## 🚩 Available Flags

| Flag             | Type     | Required | Default Value | Description                               |
| ---------------- | -------- | -------- | ------------- | ----------------------------------------- |
| `--provider`     | `string` | ❌ No     | `github`      | CI provider (`github`)                    |
| `--with-docker`  | `bool`   | ❌ No     | `false`       | Include Docker build/push steps           |
| `--with-deploy`  | `bool`   | ❌ No     | `false`       | Include deploy workflow for tagged releases |
| `--go-version`   | `string` | ❌ No     | auto-detected | Go version (default: read from `go.mod`) |
| `--dry-run`      | `bool`   | ❌ No     | `false`       | Preview files without writing             |
| `--force`        | `bool`   | ❌ No     | `false`       | Overwrite existing files                  |
| `--backup`       | `bool`   | ❌ No     | `false`       | Backup existing files before overwriting  |

## 📖 Usage Examples

### Basic CI Pipeline
```bash
goca ci
```
Generates `.github/workflows/test.yml` and `.github/workflows/build.yml`.

### Full Pipeline with Docker and Deploy
```bash
goca ci --with-docker --with-deploy
```
Generates test, build, and deploy workflows with Docker build/push steps.

### Specify Go Version
```bash
goca ci --go-version 1.22
```

### Preview Before Generating
```bash
goca ci --with-docker --dry-run
```

## 📁 Generated Files

```
.github/
└── workflows/
    ├── test.yml      # Always generated
    ├── build.yml     # Always generated
    └── deploy.yml    # Only with --with-deploy
```

### test.yml
Triggers on push to `main` and pull requests. Runs:
- `go mod download`
- `go vet ./...`
- `go test -race -coverprofile=coverage.txt ./...`

### build.yml
Triggers on push to `main`. Runs:
- `go build -v -o app ./...`
- Uploads binary as artifact
- Optionally builds and pushes Docker image (`--with-docker`)

### deploy.yml (optional)
Triggers on version tags (`v*`). Runs:
- Downloads build artifact
- Placeholder deploy step (customize for your environment)
- Optionally pushes Docker image with version tag

## 🔗 Related Commands

- [`goca init`](Command-Init.md) — Initialize project structure
- [`goca feature`](Command-Feature.md) — Generate complete feature stack
