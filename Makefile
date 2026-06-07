# Goca Makefile

.PHONY: help build test clean install release dev-setup lint fmt version \
        hooks-install hooks-update hooks-uninstall pre-commit pre-push

# Variables
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "dev")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -X github.com/sazardev/goca/cmd.Version=$(VERSION) -X github.com/sazardev/goca/cmd.BuildTime=$(BUILD_TIME) -X github.com/sazardev/goca/cmd.GitCommit=$(GIT_COMMIT)

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build
build: ## Build the CLI binary with version info
	@echo "Building Goca CLI v$(VERSION)..."
	go build -ldflags "$(LDFLAGS)" -o goca .

build-all: ## Build for all platforms with version info
	@echo "Building for all platforms..."
	@mkdir -p dist
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/goca-windows-amd64.exe
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/goca-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/goca-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/goca-darwin-arm64
	@echo "Binaries built in dist/ directory"
	@ls -la dist/

version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"

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
dev-setup: ## Setup development environment (all tools + hooks)
	@echo "Setting up development environment..."
	go mod tidy
	@echo "Installing dev tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/evilmartians/lefthook@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install go.uber.org/nilaway/cmd/nilaway@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "Installing docs dependencies..."
	cd docs && npm ci 2>/dev/null || echo "npm not available — skip docs deps"
	@echo "Installing git hooks..."
	lefthook install 2>/dev/null || echo "lefthook not in PATH — run: lefthook install"
	@echo "Development environment ready!"

fmt: ## Format code with gofumpt + goimports
	@echo "Formatting code with gofumpt..."
	gofumpt -w .
	goimports -w .

lint: ## Run linter (ultra-aggressive: 80+ linters, zero tolerance)
	@echo "Running linter (ultra-aggressive)..."
	golangci-lint run --max-issues-per-linter=0 --max-same-issues=0
	@echo "Running staticcheck..."
	staticcheck ./...
	@echo "Running gosec..."
	gosec -quiet -confidence=medium -severity=medium ./...
	@echo "Running nilaway..."
	nilaway ./...

# Release
release-patch: ## Create a patch release (x.y.Z)
	@echo "Creating patch release..."
	./scripts/release.sh patch

release-minor: ## Create a minor release (x.Y.0)
	@echo "Creating minor release..."
	./scripts/release.sh minor

release-major: ## Create a major release (X.0.0)
	@echo "Creating major release..."
	./scripts/release.sh major

release-auto: ## Auto-detect release type from commits
	@echo "Auto-detecting release type..."
	./scripts/release.sh auto

release: release-auto ## Alias for release-auto

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
dev: fmt lint test build ## Format (gofumpt), lint (80+ linters), test, build

# Documentation
docs: ## Generate documentation
	@echo "Documentation available at:"
	@echo "- README.md: Main documentation"
	@echo "- GUIDE.md: Complete command guide"  
	@echo "- rules.md: Clean Architecture rules"

# Version info  
show-version: ## Show current version from git tags
	@echo "Current git version: $(VERSION)"

# Git helpers  
status: ## Show git status and current version
	@echo "Current version: $(VERSION)"
	@echo "Git status:"
	@git status --short

# Docker (if needed in future)
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t sazardev/goca:latest .

# Check if everything is ready for release
pre-release-check: fmt lint test test-cli-comprehensive ## Check if everything is ready for release
	@echo "✅ All checks passed! Ready for release."
	@echo "Current version: $(VERSION)"
	@echo "To create a release, run: make release [patch|minor|major|auto]"

# ── Hooks ─────────────────────────────────────────

# OS detection for install script
ifeq ($(OS),Windows_NT)
    HOOK_INSTALL_SCRIPT = powershell -ExecutionPolicy Bypass -File scripts/install-hooks.ps1
else
    HOOK_INSTALL_SCRIPT = bash scripts/install-hooks.sh
endif

hooks-install: ## Install pre-commit + pre-push hooks (Linux/macOS/Windows)
	@echo "🔧 Installing hooks..."
	$(HOOK_INSTALL_SCRIPT)

hooks-update: ## Update hooks after lefthook.yml changes
	@echo "🔄 Updating hooks..."
	lefthook install -f

hooks-uninstall: ## Remove all git hooks
	@echo "🗑️  Removing hooks..."
	lefthook uninstall

pre-commit: ## Run pre-commit checks manually
	@echo "🔍 Running pre-commit checks..."
	lefthook run pre-commit

pre-push: ## Run pre-push checks manually
	@echo "🔍 Running pre-push checks..."
	lefthook run pre-push
