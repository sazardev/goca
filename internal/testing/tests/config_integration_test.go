package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sazardev/goca/cmd"
)

// TestConfigSystemIntegration prueba la integración completa del sistema de configuración
func TestConfigSystemIntegration(t *testing.T) {
	tempDir := t.TempDir()
	projectPath := filepath.Join(tempDir, "integration-test")
	err := os.MkdirAll(projectPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	t.Run("EndToEndConfigWorkflow", func(t *testing.T) {
		// Paso 1: Crear configuración por defecto
		manager := cmd.NewConfigManager()
		defaultConfig := manager.CreateDefaultConfig(projectPath)

		// Verificaciones básicas
		if defaultConfig == nil {
			t.Fatal("Expected default config to be created")
		}
		if defaultConfig.Project.Name != "integration-test" {
			t.Errorf("Expected project name 'integration-test', got '%s'", defaultConfig.Project.Name)
		}
		if defaultConfig.Database.Type != "postgres" {
			t.Errorf("Expected default database 'postgres', got '%s'", defaultConfig.Database.Type)
		}

		// Paso 2: Guardar configuración
		configPath := filepath.Join(projectPath, ".goca.yaml")
		manager.SetConfig(defaultConfig)
		if err := manager.SaveConfig(configPath); err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}

		// Verificar que el archivo existe
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("Config file should have been created")
		}

		// Paso 3: Cargar configuración desde archivo
		loader := cmd.NewConfigManager()
		if err := loader.LoadConfig(projectPath); err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		loadedConfig := loader.GetConfig()
		if loadedConfig == nil {
			t.Fatal("Expected loaded config to be non-nil")
		}
		if loadedConfig.Project.Name != defaultConfig.Project.Name {
			t.Errorf("Loaded config project name mismatch: expected '%s', got '%s'",
				defaultConfig.Project.Name, loadedConfig.Project.Name)
		}

		// Paso 4: Merge con flags de CLI
		flags := map[string]interface{}{
			"database": "mysql",
		}
		loader.MergeWithFlags(flags)

		mergedConfig := loader.GetConfig()
		if mergedConfig.Database.Type != "mysql" {
			t.Errorf("Expected merged database type 'mysql', got '%s'", mergedConfig.Database.Type)
		}

		// Paso 5: Validación de configuración
		validator := cmd.NewConfigManager()
		if err := validator.ValidateConfig(mergedConfig); err != nil {
			t.Errorf("Valid config should not have validation errors: %v", err)
		}
	})

	t.Run("ConfigValidationAndErrors", func(t *testing.T) {
		manager := cmd.NewConfigManager()

		// Crear configuración inválida
		invalidConfig := &cmd.GocaConfig{
			Project: cmd.ProjectConfig{
				Name:   "", // Inválido: nombre vacío
				Module: "", // Inválido: módulo vacío
			},
			Database: cmd.DatabaseConfig{
				Type: "invalid-db-type", // Inválido: tipo de BD no soportado
				Port: -1,                // Inválido: puerto negativo
			},
		}

		// Validar configuración inválida
		err := manager.ValidateConfig(invalidConfig)
		if err == nil {
			t.Error("Expected validation errors for invalid config")
		}

		errors := manager.GetErrors()
		if len(errors) == 0 {
			t.Error("Expected validation errors to be collected")
		}

		// Verificar campos específicos con errores
		foundErrors := make(map[string]bool)
		for _, e := range errors {
			foundErrors[e.Field] = true
		}

		expectedErrors := []string{"project.name", "project.module", "database.type", "database.port"}
		for _, expectedField := range expectedErrors {
			if !foundErrors[expectedField] {
				t.Errorf("Expected validation error for field: %s", expectedField)
			}
		}
	})

	t.Run("ConfigFileDiscovery", func(t *testing.T) {
		manager := cmd.NewConfigManager()

		// Test 1: No config file found
		emptyDir := filepath.Join(tempDir, "empty")
		os.MkdirAll(emptyDir, 0755)

		foundPath := manager.FindConfigFile(emptyDir)
		if foundPath != "" {
			t.Errorf("Expected empty string for no config file, got '%s'", foundPath)
		}

		// Test 2: Config file found
		configDir := filepath.Join(tempDir, "with-config")
		os.MkdirAll(configDir, 0755)

		configPath := filepath.Join(configDir, ".goca.yaml")
		testConfig := manager.CreateDefaultConfig(configDir)
		manager.SetConfig(testConfig)
		manager.SaveConfig(configPath)

		foundPath = manager.FindConfigFile(configDir)
		if foundPath != configPath {
			t.Errorf("Expected to find config at '%s', got '%s'", configPath, foundPath)
		}
	})

	t.Run("ConfigDefaults", func(t *testing.T) {
		manager := cmd.NewConfigManager()
		config := manager.CreateDefaultConfig(projectPath)

		// Test default postgres port
		if config.Database.Type != "postgres" || config.Database.Port != 5432 {
			t.Errorf("Expected postgres:5432, got %s:%d", config.Database.Type, config.Database.Port)
		}

		// Test MySQL defaults
		config.Database.Type = "mysql"
		config.Database.Port = 0 // Reset port
		manager.ApplyDefaults(config)
		if config.Database.Port != 3306 {
			t.Errorf("Expected MySQL port 3306, got %d", config.Database.Port)
		}

		// Test MongoDB defaults
		config.Database.Type = "mongodb"
		config.Database.Port = 0 // Reset port
		manager.ApplyDefaults(config)
		if config.Database.Port != 27017 {
			t.Errorf("Expected MongoDB port 27017, got %d", config.Database.Port)
		}

		// Test intelligent defaults (auth enables rate limiting)
		config.Features.Auth.Enabled = true
		config.Features.Security.RateLimit = false
		manager.ApplyDefaults(config)

		if !config.Features.Security.RateLimit {
			t.Error("Expected rate limiting to be auto-enabled when auth is enabled")
		}

		warnings := manager.GetWarnings()
		hasRateLimitWarning := false
		for _, w := range warnings {
			if w.Field == "features.security.rate_limit" {
				hasRateLimitWarning = true
				break
			}
		}
		if !hasRateLimitWarning {
			t.Error("Expected warning about auto-enabled rate limiting")
		}
	})
}

// TestConfigAndTemplateIntegration prueba la integración entre configuración y templates
func TestConfigAndTemplateIntegration(t *testing.T) {
	tempDir := t.TempDir()
	projectPath := filepath.Join(tempDir, "template-test")
	err := os.MkdirAll(projectPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	t.Run("TemplateManagerBasics", func(t *testing.T) {
		// Crear configuración con template personalizado
		manager := cmd.NewConfigManager()
		config := manager.CreateDefaultConfig(projectPath)

		config.Project.Name = "TestProject"
		config.Project.Module = "github.com/test/testproject"
		config.Templates.Variables["author"] = "Test Author"
		config.Templates.Variables["year"] = "2024"

		// Crear template manager
		templateManager := cmd.NewTemplateManager(&config.Templates, projectPath)
		if templateManager == nil {
			t.Fatal("Expected template manager to be created")
		}

		// Cargar templates y funciones helper
		if err := templateManager.LoadTemplates(); err != nil {
			t.Fatalf("Failed to load templates: %v", err)
		}

		// Test template string execution con variables
		templateString := `
Project: {{.ProjectName}}
Module: {{.Module}}
Author: {{.Author}}
Year: {{.Year}}
CamelCase: {{toCamelCase .ProjectName}}
SnakeCase: {{toSnakeCase .ProjectName}}
`

		templateData := map[string]interface{}{
			"ProjectName": config.Project.Name,
			"Module":      config.Project.Module,
			"Author":      config.Templates.Variables["author"],
			"Year":        config.Templates.Variables["year"],
		}

		result, err := templateManager.ExecuteTemplateString(templateString, templateData)
		if err != nil {
			t.Fatalf("Failed to execute template: %v", err)
		}

		// Verificar contenido generado
		expectedContents := []string{
			"Project: TestProject",
			"Module: github.com/test/testproject",
			"Author: Test Author",
			"Year: 2024",
			"CamelCase: testproject", // toCamelCase converts first letter to lowercase
			"SnakeCase: test_project",
		}

		for _, expected := range expectedContents {
			if !contains(result, expected) {
				t.Errorf("Expected result to contain '%s'\nFull result:\n%s", expected, result)
			}
		}
	})
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr ||
		len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findInString(s, substr))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
