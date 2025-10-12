package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sazardev/goca/cmd"
)

// TestProjectTemplates tests the project template system
func TestProjectTemplates(t *testing.T) {
	t.Run("GetProjectTemplates", func(t *testing.T) {
		templates := cmd.GetProjectTemplates()

		// Should have all predefined templates
		expectedTemplates := []string{"minimal", "rest-api", "microservice", "monolith", "enterprise"}

		for _, name := range expectedTemplates {
			if _, exists := templates[name]; !exists {
				t.Errorf("Expected template '%s' not found", name)
			}
		}

		// Each template should have required fields
		for name, template := range templates {
			if template.Name == "" {
				t.Errorf("Template '%s' has empty Name", name)
			}
			if template.Description == "" {
				t.Errorf("Template '%s' has empty Description", name)
			}
			if template.Config == "" {
				t.Errorf("Template '%s' has empty Config", name)
			}
		}
	})

	t.Run("ValidateTemplateName", func(t *testing.T) {
		// Valid templates
		validTemplates := []string{"minimal", "rest-api", "microservice", "monolith", "enterprise"}
		for _, name := range validTemplates {
			if !cmd.ValidateTemplateName(name) {
				t.Errorf("Template '%s' should be valid", name)
			}
		}

		// Invalid templates
		invalidTemplates := []string{"invalid", "unknown", "test"}
		for _, name := range invalidTemplates {
			if cmd.ValidateTemplateName(name) {
				t.Errorf("Template '%s' should be invalid", name)
			}
		}

		// Empty should be valid (means no template)
		if !cmd.ValidateTemplateName("") {
			t.Error("Empty template name should be valid")
		}
	})

	t.Run("GetTemplateConfig", func(t *testing.T) {
		// Test each template configuration
		templates := []string{"minimal", "rest-api", "microservice", "monolith", "enterprise"}

		for _, name := range templates {
			config, err := cmd.GetTemplateConfig(name)
			if err != nil {
				t.Errorf("Failed to get config for template '%s': %v", name, err)
			}

			if config == "" {
				t.Errorf("Config for template '%s' is empty", name)
			}

			// Config should be valid YAML with project section
			if !strings.Contains(config, "project:") {
				t.Errorf("Config for template '%s' missing 'project:' section", name)
			}
		}

		// Test invalid template
		_, err := cmd.GetTemplateConfig("invalid-template")
		if err == nil {
			t.Error("Expected error for invalid template, got nil")
		}
	})

	t.Run("GetTemplateNames", func(t *testing.T) {
		names := cmd.GetTemplateNames()

		if len(names) == 0 {
			t.Error("GetTemplateNames returned empty list")
		}

		// Should include all expected templates
		expectedNames := map[string]bool{
			"minimal":      false,
			"rest-api":     false,
			"microservice": false,
			"monolith":     false,
			"enterprise":   false,
		}

		for _, name := range names {
			if _, exists := expectedNames[name]; exists {
				expectedNames[name] = true
			}
		}

		for name, found := range expectedNames {
			if !found {
				t.Errorf("Expected template name '%s' not found in GetTemplateNames()", name)
			}
		}
	})
}

// TestTemplateConfigurations validates each template configuration
func TestTemplateConfigurations(t *testing.T) {
	templates := cmd.GetProjectTemplates()

	t.Run("MinimalTemplate", func(t *testing.T) {
		template := templates["minimal"]
		config := template.Config

		// Should have essential configuration only
		requiredSections := []string{"project:", "database:", "architecture:", "generation:", "testing:"}
		for _, section := range requiredSections {
			if !strings.Contains(config, section) {
				t.Errorf("Minimal template missing required section: %s", section)
			}
		}

		// Should NOT have complex features
		complexFeatures := []string{"features:", "deploy:", "monitoring:"}
		for _, feature := range complexFeatures {
			if strings.Contains(config, feature) {
				t.Errorf("Minimal template should not have complex feature: %s", feature)
			}
		}
	})

	t.Run("RestAPITemplate", func(t *testing.T) {
		template := templates["rest-api"]
		config := template.Config

		// Should have REST API specific configuration
		if !strings.Contains(config, "swagger:") {
			t.Error("REST API template should include Swagger configuration")
		}

		if !strings.Contains(config, "validation:") {
			t.Error("REST API template should include validation configuration")
		}

		if !strings.Contains(config, "testing:") {
			t.Error("REST API template should include testing configuration")
		}
	})

	t.Run("MicroserviceTemplate", func(t *testing.T) {
		template := templates["microservice"]
		config := template.Config

		// Should have microservice specific features
		microserviceFeatures := []string{
			"events: true",
			"uuid: true",
			"audit: true",
		}

		for _, feature := range microserviceFeatures {
			if !strings.Contains(config, feature) {
				t.Errorf("Microservice template should include: %s", feature)
			}
		}
	})

	t.Run("MonolithTemplate", func(t *testing.T) {
		template := templates["monolith"]
		config := template.Config

		// Should have monolith specific features
		monolithFeatures := []string{
			"features:",
			"auth:",
			"cache:",
			"logging:",
			"monitoring:",
		}

		for _, feature := range monolithFeatures {
			if !strings.Contains(config, feature) {
				t.Errorf("Monolith template should include: %s", feature)
			}
		}
	})

	t.Run("EnterpriseTemplate", func(t *testing.T) {
		template := templates["enterprise"]
		config := template.Config

		// Should have comprehensive enterprise features
		enterpriseFeatures := []string{
			"security:",
			"monitoring:",
			"deploy:",
			"docker:",
			"kubernetes:",
			"ci:",
		}

		for _, feature := range enterpriseFeatures {
			if !strings.Contains(config, feature) {
				t.Errorf("Enterprise template should include: %s", feature)
			}
		}

		// Should have high test coverage threshold
		if !strings.Contains(config, "threshold: 85") {
			t.Error("Enterprise template should have high test coverage threshold (85)")
		}
	})
}

// TestTemplateYAMLValidity tests that each template generates valid YAML
func TestTemplateYAMLValidity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping template YAML validation in short mode")
	}

	templates := cmd.GetProjectTemplates()

	for name, template := range templates {
		t.Run(name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Write template config to file
			configPath := filepath.Join(tmpDir, ".goca.yaml")

			// Replace placeholders
			configContent := strings.ReplaceAll(template.Config, "project:", "project:\n  name: test-project\n  module: github.com/test/project")

			err := os.WriteFile(configPath, []byte(configContent), 0644)
			if err != nil {
				t.Fatalf("Failed to write config file: %v", err)
			}

			// Try to load with ConfigManager
			manager := cmd.NewConfigManager()
			err = manager.LoadConfig(tmpDir)
			if err != nil {
				t.Errorf("Template '%s' generated invalid YAML: %v", name, err)
			}

			config := manager.GetConfig()
			if config == nil {
				t.Errorf("Template '%s' loaded but config is nil", name)
			}
		})
	}
}

// TestTemplateIntegrationWithInit tests template integration with init command
func TestTemplateIntegrationWithInit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping template integration tests in short mode")
	}

	testCases := []struct {
		name     string
		template string
	}{
		{"Minimal", "minimal"},
		{"RestAPI", "rest-api"},
		{"Microservice", "microservice"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			projectDir := filepath.Join(tmpDir, "test-project")

			// Simulate project initialization with template
			configContent, err := cmd.GetTemplateConfig(tc.template)
			if err != nil {
				t.Fatalf("Failed to get template config: %v", err)
			}

			// Create project directory
			if err := os.MkdirAll(projectDir, 0755); err != nil {
				t.Fatalf("Failed to create project dir: %v", err)
			}

			// Write config file
			configPath := filepath.Join(projectDir, ".goca.yaml")
			configWithProject := strings.ReplaceAll(configContent, "project:", "project:\n  name: test-project\n  module: github.com/test/project")

			if err := os.WriteFile(configPath, []byte(configWithProject), 0644); err != nil {
				t.Fatalf("Failed to write config: %v", err)
			}

			// Verify config file exists
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				t.Error("Config file was not created")
			}

			// Verify config can be loaded
			manager := cmd.NewConfigManager()
			if err := manager.LoadConfig(projectDir); err != nil {
				t.Errorf("Failed to load generated config: %v", err)
			}

			config := manager.GetConfig()
			if config == nil {
				t.Error("Loaded config is nil")
			}

			t.Logf("Successfully initialized project with template '%s'", tc.template)
		})
	}
}

// TestTemplateConfigurationOverrides tests that templates work with CLI flag overrides
func TestTemplateConfigurationOverrides(t *testing.T) {
	tmpDir := t.TempDir()

	// Create project with rest-api template
	configContent, err := cmd.GetTemplateConfig("rest-api")
	if err != nil {
		t.Fatalf("Failed to get template: %v", err)
	}

	// Write config
	configPath := filepath.Join(tmpDir, ".goca.yaml")
	configWithProject := strings.ReplaceAll(configContent, "project:", "project:\n  name: test-api\n  module: github.com/test/api")

	if err := os.WriteFile(configPath, []byte(configWithProject), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Change to tmpDir
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Load config
	ci := cmd.NewConfigIntegration()
	if err := ci.LoadConfigForProject(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Config should have rest-api settings
	if !ci.HasConfigFile() {
		t.Error("Config file should be detected")
	}

	// Test CLI override
	flags := map[string]interface{}{
		"database": "mysql",
	}
	ci.MergeWithCLIFlags(flags)

	// CLI should override template config
	effectiveDB := ci.GetDatabaseType("mysql")
	if effectiveDB != "mysql" {
		t.Errorf("Expected mysql (CLI override), got %s", effectiveDB)
	}

	t.Log("Template configuration with CLI overrides works correctly")
}
