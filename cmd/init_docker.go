package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func createMigrations(projectName string, sm ...*SafetyManager) {
	// Create migrations directory
	migrationDir := filepath.Join(projectName, "migrations")
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		ui.Warning(fmt.Sprintf("Error creating migrations directory: %v", err))
		return
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
bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path ./migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" up

# Rollback last migration
migrate -path ./migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" down 1


### Manual execution
bash
# Apply migration
psql -h localhost -U postgres -d your_db -f migrations/001_initial.up.sql

# Rollback migration  
psql -h localhost -U postgres -d your_db -f migrations/001_initial.down.sql


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

	// Docker Compose
	dockerComposeContent := fmt.Sprintf(`version: '3.8'

services:
  %s:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=database
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
      - db_data:/var/lib/postgresql/data
    healthcheck:
      test: %s
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  db_data:
`, projectName, getDatabaseUser(database), projectName, getDatabaseImage(database), getDatabaseEnvVars(database, projectName), getDatabasePort(database), getDatabasePort(database), getDatabaseHealthCheck(database))

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
	case "mysql":
		return "mysql:8.0"
	case "mongodb":
		return "mongo:7.0"
	default:
		return "postgres:15"
	}
}

func getDatabaseEnvVars(database, projectName string) string {
	switch database {
	case "mysql":
		return fmt.Sprintf("\n      - MYSQL_ROOT_PASSWORD=password\n      - MYSQL_DATABASE=%s", projectName)
	case "mongodb":
		return fmt.Sprintf("\n      - MONGO_INITDB_ROOT_USERNAME=admin\n      - MONGO_INITDB_ROOT_PASSWORD=password\n      - MONGO_INITDB_DATABASE=%s", projectName)
	default:
		return fmt.Sprintf("\n      - POSTGRES_USER=postgres\n      - POSTGRES_PASSWORD=password\n      - POSTGRES_DB=%s", projectName)
	}
}

func getDatabaseHealthCheck(database string) string {
	switch database {
	case "mysql":
		return `["CMD", "mysqladmin", "ping", "-h", "localhost"]`
	case "mongodb":
		return `["CMD", "mongo", "--eval", "db.adminCommand('ping')"]`
	default:
		return `["CMD-SHELL", "pg_isready -U postgres"]`
	}
}

