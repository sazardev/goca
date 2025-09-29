package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sazardev/goca/cmd"
	"github.com/sazardev/goca/internal/testing/framework"
)

// TestConfigManagerBasics tests basic ConfigManager functionality
func TestConfigManagerBasics(t *testing.T) {
	tc := framework.NewTestContext(t)
	defer tc.Cleanup()

	t.Run("CreateDefaultConfig", func(t *testing.T) {
		manager := cmd.NewConfigManager()

		// Test creating default config
		config := manager.CreateDefaultConfig("test-project")
		if config == nil {
			t.Fatal("Expected default config to be created")
		}

		// Validate default values
		if config.Project.Name != "test-project" {
			t.Errorf("Expected project name 'test-project', got '%s'", config.Project.Name)
		}

		if config.Database.Type != "postgres" {
			t.Errorf("Expected default database 'postgres', got '%s'", config.Database.Type)
		}

		if config.Database.Port != 5432 {
			t.Errorf("Expected default port 5432, got %d", config.Database.Port)
		}

		if !config.Generation.Validation.Enabled {
			t.Error("Expected validation to be enabled by default")
		}

		if !config.Testing.Enabled {
			t.Error("Expected testing to be enabled by default")
		}
	})

	t.Run("SaveAndLoadConfig", func(t *testing.T) {
		manager := cmd.NewConfigManager()

		// Create a test config
		config := manager.CreateDefaultConfig("save-load-test")
		config.Project.Description = "Test description"
		config.Database.Type = "mysql"
		config.Database.Port = 3306

		// Save config
		configPath := filepath.Join(tc.TempDir, ".goca.yaml")
		manager.SetConfig(config) // Set config internally first
		if err := manager.SaveConfig(configPath); err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}

		// Load config back
		manager2 := cmd.NewConfigManager()
		if err := manager2.LoadFromFile(configPath); err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		loadedConfig := manager2.GetConfig()
		if loadedConfig == nil {
			t.Fatal("Expected config to be loaded")
		}

		// Verify loaded values
		if loadedConfig.Project.Name != "save-load-test" {
			t.Errorf("Expected project name 'save-load-test', got '%s'", loadedConfig.Project.Name)
		}

		if loadedConfig.Project.Description != "Test description" {
			t.Errorf("Expected description 'Test description', got '%s'", loadedConfig.Project.Description)
		}

		if loadedConfig.Database.Type != "mysql" {
			t.Errorf("Expected database 'mysql', got '%s'", loadedConfig.Database.Type)
		}

		if loadedConfig.Database.Port != 3306 {
			t.Errorf("Expected port 3306, got %d", loadedConfig.Database.Port)
		}
	})
}

// TestConfigValidation tests configuration validation
func TestConfigValidation(t *testing.T) {
	tc := framework.NewTestContext(t)
	defer tc.Cleanup()

	t.Run("ValidConfig", func(t *testing.T) {
		manager := cmd.NewConfigManager()
		config := manager.CreateDefaultConfig("valid-project")

		if err := manager.ValidateConfig(config); err != nil {
			t.Errorf("Expected valid config to pass validation, got error: %v", err)
		}

		errors := manager.GetErrors()
		if len(errors) > 0 {
			t.Errorf("Expected no validation errors, got %d errors", len(errors))
		}
	})

	t.Run("InvalidProjectName", func(t *testing.T) {
		manager := cmd.NewConfigManager()
		config := manager.CreateDefaultConfig("valid-project")
		config.Project.Name = "" // Invalid: empty name

		err := manager.ValidateConfig(config)
		if err == nil {
			t.Error("Expected validation error for empty project name")
		}

		errors := manager.GetErrors()
		if len(errors) == 0 {
			t.Error("Expected validation errors for empty project name")
		}

		// Check that error mentions project name
		found := false
		for _, e := range errors {
			if strings.Contains(e.Field, "project.name") {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error for project.name field")
		}
	})

	t.Run("InvalidDatabaseType", func(t *testing.T) {
		manager := cmd.NewConfigManager()
		config := manager.CreateDefaultConfig("valid-project")
		config.Database.Type = "invalid-db" // Invalid database type

		err := manager.ValidateConfig(config)
		if err == nil {
			t.Error("Expected validation error for invalid database type")
		}

		errors := manager.GetErrors()
		if len(errors) == 0 {
			t.Error("Expected validation errors for invalid database type")
		}
	})

	t.Run("InvalidPort", func(t *testing.T) {
		manager := cmd.NewConfigManager()
		config := manager.CreateDefaultConfig("valid-project")
		config.Database.Port = -1 // Invalid port

		err := manager.ValidateConfig(config)
		if err == nil {
			t.Error("Expected validation error for invalid port")
		}

		errors := manager.GetErrors()
		if len(errors) == 0 {
			t.Error("Expected validation errors for invalid port")
		}
	})

	t.Run("Warnings", func(t *testing.T) {
		manager := cmd.NewConfigManager()
		config := manager.CreateDefaultConfig("valid-project")
		config.Project.Version = "" // Should generate warning

		// Validation should pass but generate warnings
		if err := manager.ValidateConfig(config); err != nil {
			t.Errorf("Expected validation to pass with warnings, got error: %v", err)
		}

		warnings := manager.GetWarnings()
		if len(warnings) == 0 {
			t.Error("Expected warnings for empty version")
		}
	})
}

// TestConfigFileDiscovery tests configuration file discovery
func TestConfigFileDiscovery(t *testing.T) {
	tc := framework.NewTestContext(t)
	defer tc.Cleanup()

	t.Run("FindConfigFile", func(t *testing.T) {
		// Create test directory structure
		testDir := filepath.Join(tc.TempDir, "project")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Create .goca.yaml file
		configContent := `
project:
  name: test-discovery
  module: github.com/test/discovery
database:
  type: postgres
  port: 5432
`
		configPath := filepath.Join(testDir, ".goca.yaml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		// Test discovery
		manager := cmd.NewConfigManager()
		foundPath := manager.FindConfigFile(testDir)

		if foundPath == "" {
			t.Error("Expected to find config file")
		}

		if foundPath != configPath {
			t.Errorf("Expected to find config at '%s', got '%s'", configPath, foundPath)
		}
	})

	t.Run("ConfigFileNotFound", func(t *testing.T) {
		// Create empty directory
		testDir := filepath.Join(tc.TempDir, "empty-project")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		manager := cmd.NewConfigManager()
		foundPath := manager.FindConfigFile(testDir)

		if foundPath != "" {
			t.Errorf("Expected no config file to be found, got '%s'", foundPath)
		}
	})

	t.Run("MultipleConfigFiles", func(t *testing.T) {
		// Create test directory
		testDir := filepath.Join(tc.TempDir, "multi-config")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Create multiple config files (priority: .goca.yaml > .goca.yml > goca.yaml > goca.yml)
		configContent := `
project:
  name: multi-test
database:
  type: postgres
`

		// Create lower priority files first
		if err := os.WriteFile(filepath.Join(testDir, "goca.yml"), []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}
		if err := os.WriteFile(filepath.Join(testDir, "goca.yaml"), []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}
		if err := os.WriteFile(filepath.Join(testDir, ".goca.yml"), []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		// Create highest priority file
		highPriorityPath := filepath.Join(testDir, ".goca.yaml")
		if err := os.WriteFile(highPriorityPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		manager := cmd.NewConfigManager()
		foundPath := manager.FindConfigFile(testDir)

		if foundPath != highPriorityPath {
			t.Errorf("Expected highest priority config '%s', got '%s'", highPriorityPath, foundPath)
		}
	})
}

// TestConfigMerging tests configuration merging with CLI flags
func TestConfigMerging(t *testing.T) {
	manager := cmd.NewConfigManager()
	config := manager.CreateDefaultConfig("merge-test")

	// Original values
	if config.Database.Type != "postgres" {
		t.Errorf("Expected default database 'postgres', got '%s'", config.Database.Type)
	}

	// Set config and merge with CLI flags
	manager.SetConfig(config)
	flags := map[string]interface{}{
		"database": "mysql",
		"auth":     true,
	}

	manager.MergeWithFlags(flags)

	// Get merged config
	mergedConfig := manager.GetConfig()

	// Check merged values
	if mergedConfig.Database.Type != "mysql" {
		t.Errorf("Expected merged database 'mysql', got '%s'", mergedConfig.Database.Type)
	}

	if !mergedConfig.Features.Auth.Enabled {
		t.Error("Expected auth to be enabled after merging")
	}
}

// TestConfigDefaults tests intelligent default application
func TestConfigDefaults(t *testing.T) {
	manager := cmd.NewConfigManager()

	t.Run("PostgresDefaults", func(t *testing.T) {
		config := manager.CreateDefaultConfig("postgres-test")
		config.Database.Type = "postgres"
		config.Database.Port = 0 // Should get default

		manager.ApplyDefaults(config)

		if config.Database.Port != 5432 {
			t.Errorf("Expected postgres default port 5432, got %d", config.Database.Port)
		}
	})

	t.Run("MySQLDefaults", func(t *testing.T) {
		config := manager.CreateDefaultConfig("mysql-test")
		config.Database.Type = "mysql"
		config.Database.Port = 0 // Should get default

		manager.ApplyDefaults(config)

		if config.Database.Port != 3306 {
			t.Errorf("Expected mysql default port 3306, got %d", config.Database.Port)
		}
	})

	t.Run("AuthSecurityDefaults", func(t *testing.T) {
		config := manager.CreateDefaultConfig("auth-test")
		config.Features.Auth.Enabled = true
		config.Features.Security.RateLimit = false // Should be auto-enabled

		manager.ApplyDefaults(config)

		if !config.Features.Security.RateLimit {
			t.Error("Expected rate limit to be auto-enabled when auth is enabled")
		}

		warnings := manager.GetWarnings()
		found := false
		for _, w := range warnings {
			if strings.Contains(w.Field, "security.rate_limit") {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected warning about auto-enabling rate limit")
		}
	})
}
