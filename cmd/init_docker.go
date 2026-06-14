package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func createMigrations(projectName string, sm ...*SafetyManager) {
	// In dry-run mode do not touch the filesystem; file writes below are routed
	// through the SafetyManager which records them instead.
	dryRun := len(sm) > 0 && sm[0] != nil && sm[0].DryRun

	// Create migrations directory
	migrationDir := filepath.Join(projectName, "migrations")
	if !dryRun {
		if err := os.MkdirAll(migrationDir, 0o755); err != nil {
			ui.Warning(fmt.Sprintf("Error creating migrations directory: %v", err))
			return
		}
	}

	// Create initial migration
	migrationContent := `-- Initial migration
-- This file contains the initial database schema

-- Enable UUID extension (PostgreSQL)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Example: Users table
-- Uncomment and modify based on your needs
/*
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
*/

-- Add your tables here
-- Remember to create corresponding down migration files for rollbacks
`

	if err := writeFile(filepath.Join(migrationDir, "001_initial.up.sql"), migrationContent, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error creating up migration: %v", err))
	}

	// Create down migration
	downMigrationContent := `-- Down migration for initial schema
-- This file should reverse changes made in 001_initial.up.sql

-- Example: Drop users table
-- DROP TABLE IF EXISTS users;

-- Add your down migration here
`

	if err := writeFile(filepath.Join(migrationDir, "001_initial.down.sql"), downMigrationContent, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error creating down migration: %v", err))
	}

	// Create README for migrations
	migrationReadme := `# Database Migrations

This directory contains database migration files.

## Structure
- *.up.sql - Migration files (apply changes)
- *.down.sql - Rollback files (reverse changes)

## Usage

### Using golang-migrate tool
` + "```bash" + `
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path ./migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" up

# Rollback last migration
migrate -path ./migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" down 1
` + "```" + `

### Manual execution
` + "```bash" + `
# Apply migration
psql -h localhost -U postgres -d your_db -f migrations/001_initial.up.sql

# Rollback migration
psql -h localhost -U postgres -d your_db -f migrations/001_initial.down.sql
` + "```" + `

## Creating new migrations
1. Create new files: 002_description.up.sql and 002_description.down.sql
2. Add your changes in the .up.sql file
3. Add the reverse changes in the .down.sql file
`

	if err := writeFile(filepath.Join(migrationDir, "README.md"), migrationReadme, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error creating migration README: %v", err))
	}
}

func createMakefile(projectName string, sm ...*SafetyManager) {
	makefileContent := fmt.Sprintf(`# Makefile for %s
.PHONY: help build run test clean docker-build docker-run deps lint migrate-up migrate-down

# Variables
APP_NAME := %s
DOCKER_IMAGE := %s:latest
MIGRATE_PATH := ./migrations
DATABASE_URL := postgres://postgres:@localhost/%s?sslmode=disable

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%%-20s\033[0m %%s\n", $$1, $$2}'

deps: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o bin/$(APP_NAME) cmd/server/main.go

run: ## Run the application
	go run cmd/server/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	golangci-lint run

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

# Database migrations
migrate-install: ## Install migrate tool
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up: ## Apply database migrations
	migrate -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" up

migrate-down: ## Rollback last migration
	migrate -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" down 1

migrate-force: ## Force migration version (use with caution)
	migrate -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" force $(VERSION)

# Docker commands
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run application in Docker
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

docker-compose-up: ## Start services with Docker Compose
	docker-compose up -d

docker-compose-down: ## Stop services with Docker Compose
	docker-compose down

# Development helpers
dev-db: ## Start development database
	docker run --name %s-postgres -e POSTGRES_DB=%s -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15

dev-db-stop: ## Stop development database
	docker stop %s-postgres && docker rm %s-postgres

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

mod-upgrade: ## Upgrade dependencies
	go get -u ./...
	go mod tidy

# Production helpers
build-prod: ## Build for production
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/$(APP_NAME) cmd/server/main.go

# Security
sec-scan: ## Run security scan
	gosec ./...

# API documentation
api-docs: ## Generate API documentation
	swag init -g cmd/server/main.go
`, projectName, projectName, projectName, projectName, projectName, projectName, projectName, projectName)

	if err := writeFile(filepath.Join(projectName, "Makefile"), makefileContent, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error creating Makefile: %v", err))
	}
}

func createDockerfiles(projectName, database string, sm ...*SafetyManager) {
	// Dockerfile
	dockerfileContent := fmt.Sprintf(`# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/%s cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/bin/%s .

# Copy migrations if they exist
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./%s"]
`, projectName, projectName, projectName)

	if err := writeFile(filepath.Join(projectName, "Dockerfile"), dockerfileContent, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error creating Dockerfile: %v", err))
	}

	// Docker Compose. SQLite is file-based, so emit only the app service with
	// no DB container (INIT-B6).
	var dockerComposeContent string
	if database == DBSQLite {
		dockerComposeContent = fmt.Sprintf(`version: '3.8'

services:
  %s:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_NAME=%s
    restart: unless-stopped
`, projectName, projectName)
	} else {
		port := getDatabasePort(database)
		dockerComposeContent = fmt.Sprintf(`version: '3.8'

services:
  %s:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=database
      - DB_PORT=%s
      - DB_USER=%s
      - DB_PASSWORD=password
      - DB_NAME=%s
    depends_on:
      database:
        condition: service_healthy
    restart: unless-stopped

  database:
    image: %s
    environment:%s
    ports:
      - "%s:%s"
    volumes:
      - db_data:%s
    healthcheck:
      test: %s
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  db_data:
`, projectName, port, getDatabaseUser(database), projectName, getDatabaseImage(database), getDatabaseEnvVars(database, projectName), port, port, getDatabaseVolumePath(database), getDatabaseHealthCheck(database))
	}

	if err := writeFile(filepath.Join(projectName, "docker-compose.yml"), dockerComposeContent, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error creating docker-compose.yml: %v", err))
	}

	// .dockerignore
	dockerignoreContent := `# Git
.git
.gitignore

# Documentation
README.md
*.md

# Environment files
.env
.env.local
.env.example

# Build artifacts
bin/
*.exe

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Logs
*.log

# Testing
coverage.out
coverage.html
`

	if err := writeFile(filepath.Join(projectName, ".dockerignore"), dockerignoreContent, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error creating .dockerignore: %v", err))
	}
}

func getDatabaseImage(database string) string {
	switch database {
	case DBMySQL:
		return "mysql:8.0"
	case DBMongoDB:
		return "mongo:7.0"
	case DBSQLServer:
		return "mcr.microsoft.com/mssql/server:2022-latest"
	case DBDynamoDB:
		return "amazon/dynamodb-local:latest"
	case DBElasticsearch:
		return "docker.elastic.co/elasticsearch/elasticsearch:8.10.1"
	default: // postgres, postgres-json
		return "postgres:15"
	}
}

func getDatabaseEnvVars(database, projectName string) string {
	switch database {
	case DBMySQL:
		return fmt.Sprintf("\n      - MYSQL_ROOT_PASSWORD=password\n      - MYSQL_DATABASE=%s", projectName)
	case DBMongoDB:
		return fmt.Sprintf("\n      - MONGO_INITDB_ROOT_USERNAME=admin\n      - MONGO_INITDB_ROOT_PASSWORD=password\n      - MONGO_INITDB_DATABASE=%s", projectName)
	case DBSQLServer:
		return "\n      - ACCEPT_EULA=Y\n      - MSSQL_SA_PASSWORD=Your_password123"
	case DBDynamoDB:
		return "\n      - AWS_ACCESS_KEY_ID=local\n      - AWS_SECRET_ACCESS_KEY=local"
	case DBElasticsearch:
		return "\n      - discovery.type=single-node\n      - xpack.security.enabled=false\n      - ES_JAVA_OPTS=-Xms512m -Xmx512m"
	default: // postgres, postgres-json
		return fmt.Sprintf("\n      - POSTGRES_USER=postgres\n      - POSTGRES_PASSWORD=password\n      - POSTGRES_DB=%s", projectName)
	}
}

func getDatabaseHealthCheck(database string) string {
	switch database {
	case DBMySQL:
		return `["CMD", "mysqladmin", "ping", "-h", "localhost"]`
	case DBMongoDB:
		// mongo:7.0 ships mongosh, not the removed legacy mongo shell (INIT-B7).
		return `["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]`
	case DBSQLServer:
		return `["CMD-SHELL", "/opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P Your_password123 -Q 'SELECT 1' || exit 1"]`
	case DBDynamoDB:
		return `["CMD-SHELL", "curl -f http://localhost:8000 || exit 1"]`
	case DBElasticsearch:
		return `["CMD-SHELL", "curl -f http://localhost:9200/_cluster/health || exit 1"]`
	default: // postgres, postgres-json
		return `["CMD-SHELL", "pg_isready -U postgres"]`
	}
}

// getDatabaseVolumePath returns the in-container data directory for the database
// image, so the named volume is mounted at the correct path (INIT-B5).
func getDatabaseVolumePath(database string) string {
	switch database {
	case DBMySQL:
		return "/var/lib/mysql"
	case DBMongoDB:
		return "/data/db"
	case DBSQLServer:
		return "/var/opt/mssql"
	case DBDynamoDB:
		return "/home/dynamodblocal/data"
	case DBElasticsearch:
		return "/usr/share/elasticsearch/data"
	default: // postgres, postgres-json
		return "/var/lib/postgresql/data"
	}
}
