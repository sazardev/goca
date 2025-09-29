package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sazardev/goca/cmd"
)

// TestConfigCompleteWorkflow tests the complete workflow of the config system
func TestConfigCompleteWorkflow(t *testing.T) {
	t.Run("ConfigSystemCompleteFlow", func(t *testing.T) {
		// Setup test environment
		tempDir := t.TempDir()
		projectName := "goca-config-test"
		projectPath := filepath.Join(tempDir, projectName)

		if err := os.MkdirAll(projectPath, 0755); err != nil {
			t.Fatalf("Failed to create project directory: %v", err)
		}

		// Test 1: Create configuration system integrated project
		manager := cmd.NewConfigManager()
		config := manager.CreateDefaultConfig(projectPath)

		// Customize configuration for testing
		config.Project.Name = projectName
		config.Project.Module = "github.com/test/" + projectName
		config.Database.Type = "postgres"
		config.Generation.Validation.Enabled = true
		config.Testing.Enabled = true

		// Save configuration
		configPath := filepath.Join(projectPath, ".goca.yaml")
		manager.SetConfig(config)
		if err := manager.SaveConfig(configPath); err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}

		// Test 2: Create basic project structure compatible with config
		if err := createBasicStructure(projectPath, config); err != nil {
			t.Fatalf("Failed to create basic structure: %v", err)
		}

		// Test 3: Validate the structure was created correctly
		expectedDirs := []string{
			config.Architecture.Layers.Domain.Directory,
			config.Architecture.Layers.UseCase.Directory,
			config.Architecture.Layers.Repository.Directory,
			config.Architecture.Layers.Handler.Directory,
		}

		for _, dir := range expectedDirs {
			fullPath := filepath.Join(projectPath, dir)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Errorf("Expected directory %s to exist", dir)
			}
		}

		// Test 4: Verify go.mod file was created correctly
		goModPath := filepath.Join(projectPath, "go.mod")
		if _, err := os.Stat(goModPath); os.IsNotExist(err) {
			t.Error("Expected go.mod file to exist")
		}
	})

	t.Run("ConfigValidationIntegration", func(t *testing.T) {
		// Create valid configuration using defaults
		tempDir := t.TempDir()
		testPath := filepath.Join(tempDir, "validation-test")

		manager := cmd.NewConfigManager()
		config := manager.CreateDefaultConfig(testPath) // Use default config which should be valid

		// Test validation passes
		if err := manager.ValidateConfig(config); err != nil {
			t.Logf("Configuration validation failed: %v", err)
			// Show specific errors for debugging
			errors := manager.GetErrors()
			for i, e := range errors {
				t.Logf("Error %d: Field=%s, Message=%s, Value=%s", i+1, e.Field, e.Message, e.Value)
			}
			t.Errorf("Default configuration should pass validation: %v", err)
		}

		errors := manager.GetErrors()
		if len(errors) > 0 {
			t.Errorf("Default configuration should have no errors, got: %d", len(errors))
		}

		// Test warnings detection
		warnings := manager.GetWarnings()
		// Some warnings might be expected (like missing version)
		if len(warnings) > 5 {
			t.Errorf("Too many warnings for basic config: %d", len(warnings))
		}
	})

	t.Run("EndToEndWorkflow", func(t *testing.T) {
		// Simulate complete workflow: init -> configure -> generate
		tempDir := t.TempDir()
		projectName := "e2e-test"
		projectPath := filepath.Join(tempDir, projectName)

		// Step 1: Initialize with config
		manager := cmd.NewConfigManager()
		config := manager.CreateDefaultConfig(projectPath)
		config.Project.Name = projectName
		config.Project.Module = "github.com/e2e/test"

		// Step 2: Save configuration
		configPath := filepath.Join(projectPath, ".goca.yaml")
		manager.SetConfig(config)
		if err := manager.SaveConfig(configPath); err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}

		// Step 3: Load configuration (simulate CLI loading)
		loader := cmd.NewConfigManager()
		if err := loader.LoadConfig(projectPath); err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		loadedConfig := loader.GetConfig()
		if loadedConfig.Project.Name != projectName {
			t.Errorf("Loaded project name mismatch: expected %s, got %s",
				projectName, loadedConfig.Project.Name)
		}

		// Step 4: Simulate flag override (like CLI would do)
		flags := map[string]interface{}{
			"database": "mysql",
		}
		loader.MergeWithFlags(flags)

		finalConfig := loader.GetConfig()
		if finalConfig.Database.Type != "mysql" {
			t.Errorf("Flag override failed: expected mysql, got %s", finalConfig.Database.Type)
		}

		// Step 5: Use config for template generation
		templateManager := cmd.NewTemplateManager(&finalConfig.Templates, projectPath)
		if err := templateManager.LoadTemplates(); err != nil {
			// Expected if no template directory exists
		}

		// Test template execution with config values
		templateData := map[string]interface{}{
			"ProjectName":  finalConfig.Project.Name,
			"DatabaseType": finalConfig.Database.Type,
		}

		simpleTemplate := "Project: {{.ProjectName}}, DB: {{.DatabaseType}}"
		result, err := templateManager.ExecuteTemplateString(simpleTemplate, templateData)
		if err != nil {
			t.Fatalf("Template execution failed: %v", err)
		}

		expected := "Project: e2e-test, DB: mysql"
		if result != expected {
			t.Errorf("Template result mismatch:\nExpected: %s\nGot: %s", expected, result)
		}
	})
}

// Helper function to create basic project structure
func createBasicStructure(projectPath string, config *cmd.GocaConfig) error {
	// Create go.mod file
	goModContent := "module " + config.Project.Module + "\n\ngo 1.21\n"
	goModPath := filepath.Join(projectPath, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return err
	}

	// Create basic directories based on config
	dirs := []string{
		config.Architecture.Layers.Domain.Directory,
		config.Architecture.Layers.UseCase.Directory,
		config.Architecture.Layers.Repository.Directory,
		config.Architecture.Layers.Handler.Directory,
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}

		// Create a dummy .gitkeep file to make the directory exist
		keepFile := filepath.Join(fullPath, ".gitkeep")
		if err := os.WriteFile(keepFile, []byte(""), 0644); err != nil {
			return err
		}
	}

	return nil
}
