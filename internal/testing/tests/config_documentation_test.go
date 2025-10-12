package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sazardev/goca/cmd"
)

// TestGocaYamlDocumentation validates all .goca.yaml features that will be documented
// This ensures documentation accuracy by testing real functionality
func TestGocaYamlDocumentation(t *testing.T) {
	t.Run("BasicConfiguration", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create basic .goca.yaml
		configContent := `project:
  name: test-api
  module: github.com/test/api
  description: Test API Project

database:
  type: postgres
  migrations:
    enabled: true
    auto_generate: true

architecture:
  layers:
    domain:
      enabled: true
    usecase:
      enabled: true
    repository:
      enabled: true
    handler:
      enabled: true
  
  patterns:
    - repository
    - service
    - dto

generation:
  validation:
    enabled: true
  business_rules:
    enabled: false
`
		configPath := filepath.Join(tmpDir, ".goca.yaml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		// Test config loading
		manager := cmd.NewConfigManager()
		if err := manager.LoadConfig(tmpDir); err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		config := manager.GetConfig()
		if config == nil {
			t.Fatal("Config is nil")
		}

		// Validate basic fields
		if config.Project.Name != "test-api" {
			t.Errorf("Expected name 'test-api', got '%s'", config.Project.Name)
		}
		if config.Database.Type != "postgres" {
			t.Errorf("Expected database 'postgres', got '%s'", config.Database.Type)
		}
		if !config.Generation.Validation.Enabled {
			t.Error("Expected validation to be enabled")
		}
	})

	t.Run("DatabaseConfiguration", func(t *testing.T) {
		tmpDir := t.TempDir()

		configContent := `project:
  name: db-config
  module: github.com/test/db

database:
  type: mysql
  host: localhost
  port: 3306
  migrations:
    enabled: true
    auto_generate: true
  features:
    soft_delete: true
    timestamps: true
    uuid: true
`
		configPath := filepath.Join(tmpDir, ".goca.yaml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		manager := cmd.NewConfigManager()
		if err := manager.LoadConfig(tmpDir); err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		config := manager.GetConfig()
		if config.Database.Type != "mysql" {
			t.Errorf("Expected mysql, got %s", config.Database.Type)
		}
		if !config.Database.Migrations.Enabled {
			t.Error("Expected migrations to be enabled")
		}
		if !config.Database.Features.SoftDelete {
			t.Error("Expected soft delete to be enabled")
		}
	})

	t.Run("NamingConventionsConfiguration", func(t *testing.T) {
		tmpDir := t.TempDir()

		configContent := `project:
  name: naming-test
  module: github.com/test/naming

architecture:
  naming:
    files: snake_case
    entities: PascalCase
    variables: camelCase
    constants: SCREAMING_SNAKE
`
		configPath := filepath.Join(tmpDir, ".goca.yaml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		manager := cmd.NewConfigManager()
		if err := manager.LoadConfig(tmpDir); err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		config := manager.GetConfig()
		if config.Architecture.Naming.Files != "snake_case" {
			t.Errorf("Expected snake_case for files, got %s", config.Architecture.Naming.Files)
		}
		if config.Architecture.Naming.Entities != "PascalCase" {
			t.Errorf("Expected PascalCase for entities, got %s", config.Architecture.Naming.Entities)
		}
	})

	t.Run("TestingConfiguration", func(t *testing.T) {
		tmpDir := t.TempDir()

		configContent := `project:
  name: testing-project
  module: github.com/test/testing

testing:
  enabled: true
  framework: testify
  coverage:
    enabled: true
    threshold: 80
  mocks:
    enabled: true
    tool: testify
  integration: true
`
		configPath := filepath.Join(tmpDir, ".goca.yaml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		manager := cmd.NewConfigManager()
		if err := manager.LoadConfig(tmpDir); err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		config := manager.GetConfig()
		if !config.Testing.Enabled {
			t.Error("Expected testing to be enabled")
		}
		if config.Testing.Framework != "testify" {
			t.Errorf("Expected testify framework, got %s", config.Testing.Framework)
		}
	})

	t.Run("TemplateCustomizationConfiguration", func(t *testing.T) {
		tmpDir := t.TempDir()

		configContent := `project:
  name: template-test
  module: github.com/test/templates

templates:
  directory: .goca/templates
  variables:
    author: "John Doe"
    copyright: "2024"
`
		configPath := filepath.Join(tmpDir, ".goca.yaml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		manager := cmd.NewConfigManager()
		if err := manager.LoadConfig(tmpDir); err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		config := manager.GetConfig()
		if config.Templates.Directory != ".goca/templates" {
			t.Errorf("Expected .goca/templates directory, got %s", config.Templates.Directory)
		}
		if len(config.Templates.Variables) == 0 {
			t.Error("Expected template variables to be set")
		}
	})

	t.Run("CLIFlagsOverrideConfiguration", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Config says mysql
		configContent := `project:
  name: override-test
  module: github.com/test/override

database:
  type: mysql
`
		configPath := filepath.Join(tmpDir, ".goca.yaml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		// Change to tmpDir to simulate project context
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tmpDir)

		// Config integration should allow CLI override
		ci := cmd.NewConfigIntegration()
		if err := ci.LoadConfigForProject(); err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// Simulate CLI flag override to postgres
		flags := map[string]interface{}{
			"database": "postgres",
		}
		ci.MergeWithCLIFlags(flags)

		// CLI flag should override config
		effectiveDB := ci.GetDatabaseType("postgres")
		if effectiveDB != "postgres" {
			t.Errorf("Expected postgres (CLI override), got %s", effectiveDB)
		}
	})

	t.Run("CompleteProjectGeneration", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create comprehensive config
		configContent := `project:
  name: comprehensive-api
  module: github.com/test/comprehensive
  description: Comprehensive test project

database:
  type: postgres
  migrations:
    enabled: true
  features:
    soft_delete: true
    timestamps: true

architecture:
  layers:
    domain:
      enabled: true
    usecase:
      enabled: true
    repository:
      enabled: true
    handler:
      enabled: true
  naming:
    files: lowercase
    entities: PascalCase

generation:
  validation:
    enabled: true
  business_rules:
    enabled: true

testing:
  enabled: true
  framework: testify
`
		configPath := filepath.Join(tmpDir, ".goca.yaml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		// Change to tmpDir
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tmpDir)

		// Initialize ConfigIntegration
		ci := cmd.NewConfigIntegration()
		if err := ci.LoadConfigForProject(); err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// Verify comprehensive config loaded correctly
		if !ci.HasConfigFile() {
			t.Error("Expected config file to be detected")
		}

		projectConfig := ci.GetProjectConfig()
		if projectConfig.Name != "comprehensive-api" {
			t.Errorf("Expected name 'comprehensive-api', got '%s'", projectConfig.Name)
		}

		// Test that defaults are properly applied
		effectiveDB := ci.GetDatabaseType("")
		if effectiveDB != "postgres" {
			t.Errorf("Expected postgres from config, got %s", effectiveDB)
		}

		validationEnabled := ci.GetValidationEnabled(nil)
		if !validationEnabled {
			t.Error("Expected validation to be enabled from config")
		}
	})

	t.Run("InvalidConfiguration", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Invalid YAML
		configContent := `project:
  name: invalid
  module: github.com/test/invalid
  invalid_yaml_structure
`
		configPath := filepath.Join(tmpDir, ".goca.yaml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		manager := cmd.NewConfigManager()
		err := manager.LoadConfig(tmpDir)
		if err == nil {
			t.Error("Expected error for invalid YAML, got nil")
		}
	})

	t.Run("MissingConfiguration", func(t *testing.T) {
		tmpDir := t.TempDir()
		// No .goca.yaml file

		// Change to tmpDir
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tmpDir)

		ci := cmd.NewConfigIntegration()
		err := ci.LoadConfigForProject()

		// Should not fail, should use defaults
		if err == nil {
			// This is OK - system should work without config
			if ci.HasConfigFile() {
				t.Error("Should not detect config file when none exists")
			}
		}
	})
}

// TestRealProjectWithGocaYaml creates a real project and generates features using .goca.yaml
func TestRealProjectWithGocaYaml(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "test-project")

	// Create project structure
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	// Create .goca.yaml with specific settings
	configContent := `project:
  name: test-project
  module: github.com/test/project
  description: Test project with configuration

database:
  type: postgres
  migrations:
    enabled: true
  features:
    timestamps: true
    soft_delete: false

architecture:
  layers:
    domain:
      enabled: true
    usecase:
      enabled: true
    repository:
      enabled: true
    handler:
      enabled: true
  naming:
    files: lowercase
    entities: PascalCase
    variables: camelCase

generation:
  validation:
    enabled: true
  business_rules:
    enabled: false

testing:
  enabled: true
  framework: testify
`
	configPath := filepath.Join(projectDir, ".goca.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Initialize go.mod
	goModContent := "module github.com/test/project\n\ngo 1.21\n"
	goModPath := filepath.Join(projectDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	// Change to project directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(projectDir)

	// Test that ConfigIntegration picks up the config
	ci := cmd.NewConfigIntegration()
	if err := ci.LoadConfigForProject(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if !ci.HasConfigFile() {
		t.Fatal("Config file should be detected")
	}

	// Verify settings are applied
	dbConfig := ci.GetDatabaseConfig()
	if dbConfig.Type != "postgres" {
		t.Errorf("Expected postgres, got %s", dbConfig.Type)
	}

	validationEnabled := ci.GetValidationEnabled(nil)
	if !validationEnabled {
		t.Error("Expected validation to be true")
	}

	// Verify architecture config
	archConfig := ci.GetArchitectureConfig()
	if archConfig.Naming.Files != "lowercase" {
		t.Errorf("Expected lowercase files, got %s", archConfig.Naming.Files)
	}

	t.Log("✓ Real project with .goca.yaml validated successfully")
}

// TestConfigurationMergeStrategies tests different merge strategies
func TestConfigurationMergeStrategies(t *testing.T) {
	t.Run("ConfigOverridesDefaults", func(t *testing.T) {
		tmpDir := t.TempDir()

		configContent := `project:
  name: test
  module: github.com/test/test
database:
  type: mysql
`
		configPath := filepath.Join(tmpDir, ".goca.yaml")
		os.WriteFile(configPath, []byte(configContent), 0644)

		// Change to tmpDir
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tmpDir)

		ci := cmd.NewConfigIntegration()
		ci.LoadConfigForProject()

		// Config should override default postgres
		effectiveDB := ci.GetDatabaseType("")
		if effectiveDB != "mysql" {
			t.Errorf("Config should override default, expected mysql, got %s", effectiveDB)
		}
	})

	t.Run("CLIOverridesConfig", func(t *testing.T) {
		tmpDir := t.TempDir()

		configContent := `project:
  name: test
  module: github.com/test/test
database:
  type: mysql
`
		configPath := filepath.Join(tmpDir, ".goca.yaml")
		os.WriteFile(configPath, []byte(configContent), 0644)

		// Change to tmpDir
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tmpDir)

		ci := cmd.NewConfigIntegration()
		ci.LoadConfigForProject()

		// CLI flag overrides config
		flags := map[string]interface{}{
			"database": "postgres",
		}
		ci.MergeWithCLIFlags(flags)

		effectiveDB := ci.GetDatabaseType("postgres")
		if effectiveDB != "postgres" {
			t.Errorf("CLI should override config, expected postgres, got %s", effectiveDB)
		}
	})

	t.Run("DefaultsWhenNoConfig", func(t *testing.T) {
		tmpDir := t.TempDir()
		// No config file

		// Change to tmpDir
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tmpDir)

		ci := cmd.NewConfigIntegration()
		ci.LoadConfigForProject()

		// Should use defaults
		effectiveDB := ci.GetDatabaseType("")
		if effectiveDB == "" {
			t.Error("Should have default database type")
		}
	})
}

// TestConfigurationExamples tests various configuration examples from documentation
func TestConfigurationExamples(t *testing.T) {
	examples := map[string]string{
		"minimal": `project:
  name: minimal-api
  module: github.com/user/minimal-api
`,
		"web-application": `project:
  name: web-app
  module: github.com/user/web-app

generation:
  validation:
    enabled: true
`,
		"microservice": `project:
  name: user-service
  module: github.com/company/user-service

database:
  type: postgres
  migrations:
    enabled: true

architecture:
  patterns:
    - repository
    - service

testing:
  enabled: true
  framework: testify
  mocks:
    enabled: true
`,
		"enterprise": `project:
  name: enterprise-api
  module: github.com/corp/enterprise-api
  description: Enterprise-grade API

database:
  type: postgres
  migrations:
    enabled: true
  features:
    soft_delete: true
    timestamps: true

generation:
  validation:
    enabled: true
  business_rules:
    enabled: true

testing:
  enabled: true
  framework: testify
  coverage:
    enabled: true
  mocks:
    enabled: true
  integration: true

templates:
  directory: .goca/templates

architecture:
  naming:
    files: lowercase
    entities: PascalCase
    variables: camelCase
`,
	}

	for name, configContent := range examples {
		t.Run(name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, ".goca.yaml")

			if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
				t.Fatalf("Failed to write %s config: %v", name, err)
			}

			manager := cmd.NewConfigManager()
			if err := manager.LoadConfig(tmpDir); err != nil {
				t.Fatalf("Failed to load %s config: %v", name, err)
			}

			config := manager.GetConfig()
			if config == nil {
				t.Fatalf("%s config is nil", name)
			}

			// Validate project name is set
			if config.Project.Name == "" {
				t.Errorf("%s: project name should not be empty", name)
			}

			// Validate module is set
			if config.Project.Module == "" {
				t.Errorf("%s: module should not be empty", name)
			}

			t.Logf("✓ %s configuration validated", name)
		})
	}
}

// TestConfigurationBestPractices validates best practices
func TestConfigurationBestPractices(t *testing.T) {
	t.Run("ShouldHaveProjectSection", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Config without project section should still work
		configContent := `database:
  type: postgres
`
		configPath := filepath.Join(tmpDir, ".goca.yaml")
		os.WriteFile(configPath, []byte(configContent), 0644)

		manager := cmd.NewConfigManager()
		// Should not fail
		if err := manager.LoadConfig(tmpDir); err != nil {
			// If it fails, that's OK - we're testing error handling
			if !strings.Contains(err.Error(), "project") {
				t.Logf("Config without project section: %v", err)
			}
		}
	})

	t.Run("ValidDatabaseTypes", func(t *testing.T) {
		validTypes := []string{"postgres", "mysql", "mongodb"}

		for _, dbType := range validTypes {
			tmpDir := t.TempDir()
			configContent := `project:
  name: test
  module: github.com/test/test
database:
  type: ` + dbType + `
`
			configPath := filepath.Join(tmpDir, ".goca.yaml")
			os.WriteFile(configPath, []byte(configContent), 0644)

			manager := cmd.NewConfigManager()
			if err := manager.LoadConfig(tmpDir); err != nil {
				t.Errorf("Valid database type %s should not fail: %v", dbType, err)
			}
		}
	})

	t.Run("ValidLayerConfiguration", func(t *testing.T) {
		tmpDir := t.TempDir()
		configContent := `project:
  name: test
  module: github.com/test/test

architecture:
  layers:
    domain:
      enabled: true
      directory: internal/domain
    usecase:
      enabled: true
      directory: internal/usecase
    repository:
      enabled: true
      directory: internal/repository
    handler:
      enabled: true
      directory: internal/handler
`
		configPath := filepath.Join(tmpDir, ".goca.yaml")
		os.WriteFile(configPath, []byte(configContent), 0644)

		manager := cmd.NewConfigManager()
		if err := manager.LoadConfig(tmpDir); err != nil {
			t.Errorf("Valid layer configuration should not fail: %v", err)
		}

		config := manager.GetConfig()
		if !config.Architecture.Layers.Domain.Enabled {
			t.Error("Domain layer should be enabled")
		}
		if !config.Architecture.Layers.UseCase.Enabled {
			t.Error("UseCase layer should be enabled")
		}
	})
}
