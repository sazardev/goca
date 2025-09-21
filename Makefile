# Goca Makefile

.PHONY: help build test clean install release dev-setup lint fmt

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build
build: ## Build the CLI binary
	@echo "Building Goca CLI..."
	go build -o goca .

build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	GOOS=windows GOARCH=amd64 go build -o dist/goca-windows-amd64.exe
	GOOS=linux GOARCH=amd64 go build -o dist/goca-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o dist/goca-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o dist/goca-darwin-arm64
	@echo "Binaries built in dist/ directory"

# Test
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Development
dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	go mod tidy
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development environment ready!"

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

# CLI Testing
test-cli-comprehensive: ## Run CLI comprehensive tests
	@echo "Running comprehensive CLI tests..."
	go run internal/testing/test_runner.go -type=all -v

test-cli-init: ## Test only goca init command
	@echo "Testing goca init command..."
	go run internal/testing/test_runner.go -type=init -v

test-cli-feature: ## Test only goca feature command
	@echo "Testing goca feature command..."
	go run internal/testing/test_runner.go -type=feature -v

test-cli-entity: ## Test only goca entity command
	@echo "Testing goca entity command..."
	go run internal/testing/test_runner.go -type=entity -v

test-cli-quality: ## Test code quality of generated code
	@echo "Testing generated code quality..."
	go run internal/testing/test_runner.go -type=quality -v

test-cli-fast: ## Run fast CLI tests (compilation only)
	@echo "Running fast CLI tests..."
	go test ./internal/testing -run TestGocaInitCommand -v
	go test ./internal/testing -run TestCodeQuality -v

test-cli-benchmark: ## Run CLI performance benchmarks
	@echo "Running CLI benchmarks..."
	go test ./internal/testing -bench=. -benchmem -v

test-all: ## Run all tests (unit + CLI)
	@echo "Running all tests..."
	$(MAKE) test
	$(MAKE) test-cli-comprehensive

test-cli-basic: build ## Test basic CLI functionality
	@echo "Testing basic CLI functionality..."
	./goca version
	./goca help
	@echo "Basic CLI tests passed!"

# Installation
install: ## Install CLI to GOPATH/bin
	@echo "Installing Goca CLI..."
	go install .
	@echo "Goca CLI installed successfully!"

# Release
release: ## Create a new release (usage: make release VERSION=1.0.1)
ifndef VERSION
	@echo "Error: VERSION is required. Usage: make release VERSION=1.0.1"
	@exit 1
endif
	@echo "Creating release $(VERSION)..."
	./scripts/release.sh $(VERSION)

# Maintenance
clean: ## Clean build artifacts
	@echo "Cleaning up..."
	rm -f goca goca.exe
	rm -rf dist/
	rm -f coverage.out coverage.html
	go clean
	@echo "Clean completed!"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Quick development workflow
dev: fmt lint test build ## Format, lint, test, and build

# Documentation
docs: ## Generate documentation
	@echo "Documentation available at:"
	@echo "- README.md: Main documentation"
	@echo "- GUIDE.md: Complete command guide"  
	@echo "- rules.md: Clean Architecture rules"

# Version info
version: ## Show current version
	@grep "Version.*=" cmd/version.go | sed 's/.*"\(.*\)"/\1/'

# Git helpers  
status: ## Show git status and current version
	@echo "Current version: $$(make version)"
	@echo "Git status:"
	@git status --short

# Docker (if needed in future)
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t sazardev/goca:latest .

# Check if everything is ready for release
pre-release-check: fmt lint test test-cli ## Check if everything is ready for release
	@echo "✅ All checks passed! Ready for release."
	@echo "Current version: $$(make version)"
	@echo "To create a release, run: make release VERSION=x.y.z"
