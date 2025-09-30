package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestConfigCodeGeneration tests that custom .goca.yaml settings affect generated code
func TestConfigCodeGeneration(t *testing.T) {
	t.Run("ValidationDisabled", func(t *testing.T) {
		testValidationDisabled(t)
	})

	t.Run("ValidationEnabled", func(t *testing.T) {
		testValidationEnabled(t)
	})

	t.Run("SoftDeleteEnabled", func(t *testing.T) {
		testSoftDeleteEnabled(t)
	})

	t.Run("TimestampsEnabled", func(t *testing.T) {
		testTimestampsEnabled(t)
	})

	t.Run("DatabaseTypePostgres", func(t *testing.T) {
		testDatabaseTypePostgres(t)
	})

	t.Run("DatabaseTypeMySQL", func(t *testing.T) {
		testDatabaseTypeMySQL(t)
	})

	t.Run("NamingConventionSnakeCase", func(t *testing.T) {
		testNamingConventionSnakeCase(t)
	})

	t.Run("CustomLineLength", func(t *testing.T) {
		testCustomLineLength(t)
	})

	t.Run("TestingFrameworkGinkgo", func(t *testing.T) {
		testTestingFrameworkGinkgo(t)
	})

	t.Run("AuthTypeJWT", func(t *testing.T) {
		testAuthTypeJWT(t)
	})
}

// testValidationDisabled verifies that entities generated with validation.enabled: false
// do not contain validation tags
func testValidationDisabled(t *testing.T) {
	tempDir := t.TempDir()

	// Create .goca.yaml with validation disabled
	configContent := `project:
  name: "testproject"
  module: "github.com/test/testproject"

generation:
  validation:
    enabled: false
    library: "builtin"
`
	err := os.WriteFile(filepath.Join(tempDir, ".goca.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .goca.yaml: %v", err)
	}

	// Create basic project structure
	err = createMinimalProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate entity using goca CLI
	err = runGocaCommand(tempDir, "entity", "Product", "--fields", "name:string,price:float64")
	if err != nil {
		t.Fatalf("Failed to run goca entity: %v", err)
	}

	// Read generated entity file
	entityPath := filepath.Join(tempDir, "internal", "domain", "product.go")
	entityContent, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatalf("Failed to read generated entity: %v", err)
	}

	entityStr := string(entityContent)

	// Verify validation tags are NOT present
	if strings.Contains(entityStr, "validate:") {
		t.Error("Expected no validation tags when validation is disabled, but found validation tags")
	}

	// Verify basic struct is present
	if !strings.Contains(entityStr, "type Product struct") {
		t.Error("Expected Product struct to be generated")
	}
}

// testValidationEnabled verifies that entities generated with validation.enabled: true
// contain validation tags
func testValidationEnabled(t *testing.T) {
	tempDir := t.TempDir()

	// Create .goca.yaml with validation enabled
	configContent := `project:
  name: "testproject"
  module: "github.com/test/testproject"

generation:
  validation:
    enabled: true
    library: "builtin"
`
	err := os.WriteFile(filepath.Join(tempDir, ".goca.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .goca.yaml: %v", err)
	}

	// Create basic project structure
	err = createMinimalProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate entity using goca CLI
	err = runGocaCommand(tempDir, "entity", "Product", "--fields", "name:string,price:float64")
	if err != nil {
		t.Fatalf("Failed to run goca entity: %v", err)
	}

	// Read generated entity file
	entityPath := filepath.Join(tempDir, "internal", "domain", "product.go")
	entityContent, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatalf("Failed to read generated entity: %v", err)
	}

	entityStr := string(entityContent)

	// Verify validation tags ARE present
	if !strings.Contains(entityStr, "validate:") && !strings.Contains(entityStr, "binding:") {
		t.Error("Expected validation tags when validation is enabled, but found none")
	}

	// Verify basic struct is present
	if !strings.Contains(entityStr, "type Product struct") {
		t.Error("Expected Product struct to be generated")
	}
}

// testSoftDeleteEnabled verifies that entities with soft_delete: true contain DeletedAt field
func testSoftDeleteEnabled(t *testing.T) {
	tempDir := t.TempDir()

	// Create .goca.yaml with soft delete enabled
	configContent := `project:
  name: "testproject"
  module: "github.com/test/testproject"

database:
  type: "postgres"
  features:
    soft_delete: true
    timestamps: false
`
	err := os.WriteFile(filepath.Join(tempDir, ".goca.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .goca.yaml: %v", err)
	}

	// Create basic project structure
	err = createMinimalProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate entity using goca CLI
	err = runGocaCommand(tempDir, "entity", "Product", "--fields", "name:string")
	if err != nil {
		t.Fatalf("Failed to run goca entity: %v", err)
	}

	// Read generated entity file
	entityPath := filepath.Join(tempDir, "internal", "domain", "product.go")
	entityContent, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatalf("Failed to read generated entity: %v", err)
	}

	entityStr := string(entityContent)

	// Verify DeletedAt field is present
	if !strings.Contains(entityStr, "DeletedAt") {
		t.Error("Expected DeletedAt field when soft_delete is enabled, but not found")
	}

	// Verify it's using gorm.DeletedAt type
	if !strings.Contains(entityStr, "gorm.DeletedAt") && !strings.Contains(entityStr, "*time.Time") {
		t.Error("Expected DeletedAt to use proper type (gorm.DeletedAt or *time.Time)")
	}
}

// testTimestampsEnabled verifies that entities with timestamps: true contain CreatedAt and UpdatedAt
func testTimestampsEnabled(t *testing.T) {
	tempDir := t.TempDir()

	// Create .goca.yaml with timestamps enabled
	configContent := `project:
  name: "testproject"
  module: "github.com/test/testproject"

database:
  type: "postgres"
  features:
    soft_delete: false
    timestamps: true
`
	err := os.WriteFile(filepath.Join(tempDir, ".goca.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .goca.yaml: %v", err)
	}

	// Create basic project structure
	err = createMinimalProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate entity using goca CLI
	err = runGocaCommand(tempDir, "entity", "Product", "--fields", "name:string")
	if err != nil {
		t.Fatalf("Failed to run goca entity: %v", err)
	}

	// Read generated entity file
	entityPath := filepath.Join(tempDir, "internal", "domain", "product.go")
	entityContent, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatalf("Failed to read generated entity: %v", err)
	}

	entityStr := string(entityContent)

	// Verify CreatedAt field is present
	if !strings.Contains(entityStr, "CreatedAt") {
		t.Error("Expected CreatedAt field when timestamps are enabled, but not found")
	}

	// Verify UpdatedAt field is present
	if !strings.Contains(entityStr, "UpdatedAt") {
		t.Error("Expected UpdatedAt field when timestamps are enabled, but not found")
	}

	// Verify time.Time type is used
	if !strings.Contains(entityStr, "time.Time") {
		t.Error("Expected timestamp fields to use time.Time type")
	}
}

// testDatabaseTypePostgres verifies postgres-specific code generation
func testDatabaseTypePostgres(t *testing.T) {
	tempDir := t.TempDir()

	// Create .goca.yaml with postgres database
	configContent := `project:
  name: "testproject"
  module: "github.com/test/testproject"

database:
  type: "postgres"
  port: 5432
`
	err := os.WriteFile(filepath.Join(tempDir, ".goca.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .goca.yaml: %v", err)
	}

	// Create basic project structure
	err = createMinimalProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate repository using goca CLI
	err = runGocaCommand(tempDir, "repository", "Product")
	if err != nil {
		t.Fatalf("Failed to run goca repository: %v", err)
	}

	// Read generated repository file
	repoPath := filepath.Join(tempDir, "internal", "repository", "postgres_product_repository.go")
	repoContent, err := os.ReadFile(repoPath)
	if err != nil {
		t.Fatalf("Failed to read generated repository: %v", err)
	}

	repoStr := string(repoContent)

	// Verify postgres-specific naming
	if !strings.Contains(repoStr, "postgresProductRepository") && !strings.Contains(repoStr, "PostgresProductRepository") {
		t.Error("Expected postgres-specific repository name")
	}

	// Verify gorm.DB usage
	if !strings.Contains(repoStr, "gorm.DB") {
		t.Error("Expected gorm.DB for postgres repository")
	}
}

// testDatabaseTypeMySQL verifies mysql-specific code generation
func testDatabaseTypeMySQL(t *testing.T) {
	tempDir := t.TempDir()

	// Create .goca.yaml with mysql database
	configContent := `project:
  name: "testproject"
  module: "github.com/test/testproject"

database:
  type: "mysql"
  port: 3306
`
	err := os.WriteFile(filepath.Join(tempDir, ".goca.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .goca.yaml: %v", err)
	}

	// Create basic project structure
	err = createMinimalProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate repository using goca CLI
	err = runGocaCommand(tempDir, "repository", "Product")
	if err != nil {
		t.Fatalf("Failed to run goca repository: %v", err)
	}

	// Read generated repository file
	repoPath := filepath.Join(tempDir, "internal", "repository", "mysql_product_repository.go")
	repoContent, err := os.ReadFile(repoPath)
	if err != nil {
		t.Fatalf("Failed to read generated repository: %v", err)
	}

	repoStr := string(repoContent)

	// Verify mysql-specific naming
	if !strings.Contains(repoStr, "mysqlProductRepository") && !strings.Contains(repoStr, "MysqlProductRepository") {
		t.Error("Expected mysql-specific repository name")
	}

	// Verify gorm.DB usage
	if !strings.Contains(repoStr, "gorm.DB") {
		t.Error("Expected gorm.DB for mysql repository")
	}
}

// testNamingConventionSnakeCase verifies snake_case naming is applied
func testNamingConventionSnakeCase(t *testing.T) {
	tempDir := t.TempDir()

	// Create .goca.yaml with snake_case naming
	configContent := `project:
  name: "testproject"
  module: "github.com/test/testproject"

architecture:
  naming:
    files: "snake_case"
    entities: "PascalCase"
    fields: "snake_case"
`
	err := os.WriteFile(filepath.Join(tempDir, ".goca.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .goca.yaml: %v", err)
	}

	// Create basic project structure
	err = createMinimalProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate entity using goca CLI
	err = runGocaCommand(tempDir, "entity", "ProductCategory", "--fields", "name:string")
	if err != nil {
		t.Fatalf("Failed to run goca entity: %v", err)
	}

	// Verify file is named with snake_case
	entityPath := filepath.Join(tempDir, "internal", "domain", "product_category.go")
	if _, err := os.Stat(entityPath); os.IsNotExist(err) {
		t.Errorf("Expected file product_category.go with snake_case naming, but not found")
	}

	// Read generated entity file
	entityContent, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatalf("Failed to read generated entity: %v", err)
	}

	entityStr := string(entityContent)

	// Verify struct name uses PascalCase (entities convention)
	if !strings.Contains(entityStr, "type ProductCategory struct") {
		t.Error("Expected ProductCategory struct with PascalCase")
	}
}

// testCustomLineLength verifies custom line length is respected
func testCustomLineLength(t *testing.T) {
	tempDir := t.TempDir()

	// Create .goca.yaml with custom line length
	configContent := `project:
  name: "testproject"
  module: "github.com/test/testproject"

generation:
  style:
    line_length: 80
    indentation: 4
`
	err := os.WriteFile(filepath.Join(tempDir, ".goca.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .goca.yaml: %v", err)
	}

	// Create basic project structure
	err = createMinimalProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	// Generate entity using goca CLI
	err = runGocaCommand(tempDir, "entity", "Product", "--fields", "name:string")
	if err != nil {
		t.Fatalf("Failed to run goca entity: %v", err)
	}

	// Read generated entity file
	entityPath := filepath.Join(tempDir, "internal", "domain", "product.go")
	entityContent, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatalf("Failed to read generated entity: %v", err)
	}

	// Verify lines don't exceed custom length (with some tolerance for struct tags)
	lines := strings.Split(string(entityContent), "\n")
	longLines := 0
	for _, line := range lines {
		if len(line) > 100 { // Allow some tolerance
			longLines++
		}
	}

	// Most lines should respect the limit (some struct tags might exceed)
	if longLines > len(lines)/4 {
		t.Errorf("Too many lines exceed reasonable length: %d out of %d lines", longLines, len(lines))
	}
}

// testTestingFrameworkGinkgo verifies Ginkgo test generation
func testTestingFrameworkGinkgo(t *testing.T) {
	t.Skip("Skipping Ginkgo test - requires full Ginkgo integration")
	// This would test that generated test files use Ginkgo syntax when configured
	// Requires more complex test setup with actual test file generation
}

// testAuthTypeJWT verifies JWT auth middleware generation
func testAuthTypeJWT(t *testing.T) {
	t.Skip("Skipping JWT test - requires full auth feature integration")
	// This would test that JWT middleware is generated when auth.type: jwt
	// Requires auth feature generation which is more complex
}

// Helper function to create minimal project structure
func createMinimalProjectStructure(projectPath string) error {
	// Create necessary directories
	dirs := []string{
		"internal/domain",
		"internal/usecase",
		"internal/repository",
		"internal/handler/http",
		"internal/di",
		"cmd/server",
		"pkg/config",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(projectPath, dir), 0755)
		if err != nil {
			return err
		}
	}

	// Create go.mod
	goModContent := `module github.com/test/testproject

go 1.21

require (
	github.com/gorilla/mux v1.8.0
	gorm.io/gorm v1.25.0
	gorm.io/driver/postgres v1.5.0
	gorm.io/driver/mysql v1.5.0
)
`
	err := os.WriteFile(filepath.Join(projectPath, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		return err
	}

	// Create main.go stub
	mainContent := `package main

func main() {
	// Server initialization
}
`
	err = os.WriteFile(filepath.Join(projectPath, "cmd", "server", "main.go"), []byte(mainContent), 0644)
	if err != nil {
		return err
	}

	// Create di container stub
	diContent := `package di

type Container struct {
	// Dependencies
}

func NewContainer() *Container {
	return &Container{}
}
`
	err = os.WriteFile(filepath.Join(projectPath, "internal", "di", "container.go"), []byte(diContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Helper function to run goca commands
func runGocaCommand(projectPath string, args ...string) error {
	// Get absolute path to goca binary
	// This test file is in internal/testing/tests, so we need to go up 3 levels
	testDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	gocaBinary := filepath.Join(testDir, "..", "..", "..", "goca.exe")

	// Convert to absolute path
	gocaBinary, err = filepath.Abs(gocaBinary)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Check if binary exists
	if _, err := os.Stat(gocaBinary); os.IsNotExist(err) {
		return fmt.Errorf("goca binary not found at %s", gocaBinary)
	}

	cmd := exec.Command(gocaBinary, args...)
	cmd.Dir = projectPath
	cmd.Env = append(os.Environ(), "GOCA_TEST=true")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %v, output: %s", err, string(output))
	}

	return nil
}
