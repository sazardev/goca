package cmd

import "fmt"

// CITemplateData holds the values fed into CI workflow templates.
type CITemplateData struct {
	ProjectName string
	Module      string
	GoVersion   string
	Database    string
	WithDocker  bool
	WithDeploy  bool
}

// generateTestWorkflow returns a GitHub Actions test workflow YAML.
func generateTestWorkflow(data CITemplateData) string {
	svc := ""
	envBlock := ""
	if data.Database == "postgres" || data.Database == "postgres-json" {
		svc = `
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: testdb
        ports:
          - 5432:5432
        options: >-
          --health-cmd "pg_isready -U postgres"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5`
		envBlock = `
      env:
        DATABASE_URL: postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable`
	} else if data.Database == "mysql" {
		svc = `
    services:
      mysql:
        image: mysql:8
        env:
          MYSQL_ROOT_PASSWORD: root
          MYSQL_DATABASE: testdb
        ports:
          - 3306:3306
        options: >-
          --health-cmd "mysqladmin ping -h 127.0.0.1"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5`
		envBlock = `
      env:
        DATABASE_URL: root:root@tcp(127.0.0.1:3306)/testdb`
	}

	return fmt.Sprintf(`name: Test

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest%s

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '%s'

      - name: Download dependencies
        run: go mod download

      - name: Build
        run: go build ./...

      - name: Vet
        run: go vet ./...

      - name: Test
        run: go test -race -coverprofile=coverage.out -covermode=atomic ./...%s

      - name: Upload coverage
        if: github.event_name == 'push'
        uses: actions/upload-artifact@v4
        with:
          name: coverage
          path: coverage.out
`, svc, data.GoVersion, envBlock)
}

// generateBuildWorkflow returns a GitHub Actions build workflow YAML.
func generateBuildWorkflow(data CITemplateData) string {
	dockerStep := ""
	if data.WithDocker {
		dockerStep = fmt.Sprintf(`
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build Docker image
        run: |
          docker build -t %s:${{ github.sha }} .
          docker tag %s:${{ github.sha }} %s:latest
`, data.ProjectName, data.ProjectName, data.ProjectName)
	}

	return fmt.Sprintf(`name: Build

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '%s'

      - name: Download dependencies
        run: go mod download

      - name: Build binary
        run: |
          CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/%s ./cmd/server
%s`, data.GoVersion, data.ProjectName, dockerStep)
}

// generateDeployWorkflow returns a GitHub Actions deploy workflow YAML
// triggered by tag pushes.
func generateDeployWorkflow(data CITemplateData) string {
	return fmt.Sprintf(`name: Deploy

on:
  push:
    tags:
      - 'v*'

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '%s'

      - name: Download dependencies
        run: go mod download

      - name: Build release binary
        run: |
          VERSION=${{ github.ref_name }}
          CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
            -ldflags="-s -w -X main.Version=${VERSION}" \
            -o bin/%s ./cmd/server

      - name: Upload release artifact
        uses: actions/upload-artifact@v4
        with:
          name: %s-${{ github.ref_name }}
          path: bin/%s
`, data.GoVersion, data.ProjectName, data.ProjectName, data.ProjectName)
}
