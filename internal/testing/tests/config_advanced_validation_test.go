package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sazardev/goca/cmd"
)

// TestConfigAdvancedValidation tests comprehensive validation scenarios
func TestConfigAdvancedValidation(t *testing.T) {
	t.Run("InvalidDatabaseTypes", func(t *testing.T) {
		testCases := []struct {
			name       string
			dbType     string
			shouldFail bool
		}{
			{"ValidPostgres", "postgres", false},
			{"ValidMySQL", "mysql", false},
			{"ValidMongoDB", "mongodb", false},
			{"ValidSQLite", "sqlite", false},
			{"InvalidOracle", "oracle", true},
			{"InvalidMSSQL", "mssql", true},
			// EmptyType now gets default "postgres" applied, so it passes
			{"EmptyType", "", false},
			{"CaseSensitive", "POSTGRES", true},
			{"Typo", "postgress", true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tempDir := t.TempDir()
				configContent := `
project:
  name: "test-project"
  module: "github.com/test/project"
database:
  type: "` + tc.dbType + `"
  port: 5432
`
				configPath := filepath.Join(tempDir, ".goca.yaml")
				if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
					t.Fatalf("Failed to write config: %v", err)
				}

				manager := cmd.NewConfigManager()
				err := manager.LoadConfig(tempDir)

				if tc.shouldFail && err == nil {
					t.Errorf("Expected validation to fail for database type '%s' but it passed", tc.dbType)
				}

				if !tc.shouldFail && err != nil {
					t.Errorf("Expected validation to pass for database type '%s' but got error: %v", tc.dbType, err)
				}
			})
		}
	})

	t.Run("PortBoundaryValidation", func(t *testing.T) {
		testCases := []struct {
			name       string
			port       int
			shouldFail bool
		}{
			{"ValidPort", 5432, false},
			{"ValidMinPort", 1, false},
			{"ValidMaxPort", 65535, false},
			// Port 0 is valid for SQLite (doesn't use network ports)
			{"ValidZero", 0, false},
			{"InvalidNegative", -1, true},
			{"InvalidTooHigh", 65536, true},
			{"InvalidMaxInt", 99999, true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tempDir := t.TempDir()
				configContent := fmt.Sprintf(`
project:
  name: "test-project"
  module: "github.com/test/project"
database:
  type: "postgres"
  port: %d
`, tc.port)
				configPath := filepath.Join(tempDir, ".goca.yaml")
				if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
					t.Fatalf("Failed to write config: %v", err)
				}

				manager := cmd.NewConfigManager()
				err := manager.LoadConfig(tempDir)

				if tc.shouldFail && err == nil {
					t.Errorf("Expected validation to fail for port %d but it passed", tc.port)
				}

				if !tc.shouldFail && err != nil {
					t.Errorf("Expected validation to pass for port %d but got error: %v", tc.port, err)
				}
			})
		}
	})

	t.Run("NamingConventionValidation", func(t *testing.T) {
		validConventions := []string{"PascalCase", "camelCase", "snake_case", "kebab-case", "UPPER_CASE", "lowercase"}
		// Empty convention now gets default "PascalCase" applied
		invalidConventions := []string{"InvalidCase", "random", "PascalCaseTypo", "snake_Case"}

		for _, convention := range validConventions {
			t.Run("Valid_"+convention, func(t *testing.T) {
				tempDir := t.TempDir()
				configContent := `
project:
  name: "test-project"
  module: "github.com/test/project"
architecture:
  naming:
    entities: "` + convention + `"
    fields: "` + convention + `"
    files: "` + convention + `"
`
				configPath := filepath.Join(tempDir, ".goca.yaml")
				if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
					t.Fatalf("Failed to write config: %v", err)
				}

				manager := cmd.NewConfigManager()
				err := manager.LoadConfig(tempDir)

				if err != nil {
					t.Errorf("Expected validation to pass for convention '%s' but got error: %v", convention, err)
				}
			})
		}

		for _, convention := range invalidConventions {
			t.Run("Invalid_"+convention, func(t *testing.T) {
				tempDir := t.TempDir()
				configContent := `
project:
  name: "test-project"
  module: "github.com/test/project"
architecture:
  naming:
    entities: "` + convention + `"
`
				configPath := filepath.Join(tempDir, ".goca.yaml")
				if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
					t.Fatalf("Failed to write config: %v", err)
				}

				manager := cmd.NewConfigManager()
				err := manager.LoadConfig(tempDir)

				if err == nil {
					t.Errorf("Expected validation to fail for invalid convention '%s' but it passed", convention)
				}
			})
		}
	})

	t.Run("DITypeValidation", func(t *testing.T) {
		testCases := []struct {
			name       string
			diType     string
			shouldFail bool
		}{
			{"ValidManual", "manual", false},
			{"ValidWire", "wire", false},
			{"ValidFx", "fx", false},
			{"ValidDig", "dig", false},
			{"InvalidSpring", "spring", true},
			// Empty DI type now gets default "manual" applied
			{"ValidEmpty", "", false},
			{"InvalidCaseSensitive", "Manual", true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tempDir := t.TempDir()
				configContent := `
project:
  name: "test-project"
  module: "github.com/test/project"
architecture:
  di:
    type: "` + tc.diType + `"
`
				configPath := filepath.Join(tempDir, ".goca.yaml")
				if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
					t.Fatalf("Failed to write config: %v", err)
				}

				manager := cmd.NewConfigManager()
				err := manager.LoadConfig(tempDir)

				if tc.shouldFail && err == nil {
					t.Errorf("Expected validation to fail for DI type '%s' but it passed", tc.diType)
				}

				if !tc.shouldFail && err != nil {
					t.Errorf("Expected validation to pass for DI type '%s' but got error: %v", tc.diType, err)
				}
			})
		}
	})

	t.Run("CoverageThresholdValidation", func(t *testing.T) {
		testCases := []struct {
			name       string
			threshold  string
			shouldFail bool
		}{
			{"ValidZero", "0.0", false},
			{"ValidFifty", "50.0", false},
			{"ValidHundred", "100.0", false},
			{"InvalidNegative", "-10.0", true},
			{"InvalidOverHundred", "101.0", true},
			{"InvalidLarge", "999.9", true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tempDir := t.TempDir()
				configContent := `
project:
  name: "test-project"
  module: "github.com/test/project"
testing:
  framework: "testify"
  coverage:
    threshold: ` + tc.threshold + `
`
				configPath := filepath.Join(tempDir, ".goca.yaml")
				if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
					t.Fatalf("Failed to write config: %v", err)
				}

				manager := cmd.NewConfigManager()
				err := manager.LoadConfig(tempDir)

				if tc.shouldFail && err == nil {
					t.Errorf("Expected validation to fail for threshold %s but it passed", tc.threshold)
				}

				if !tc.shouldFail && err != nil {
					t.Errorf("Expected validation to pass for threshold %s but got error: %v", tc.threshold, err)
				}
			})
		}
	})

	t.Run("TestingFrameworkValidation", func(t *testing.T) {
		validFrameworks := []string{"testify", "ginkgo", "builtin"}
		// Empty framework now gets default "testify" applied
		invalidFrameworks := []string{"junit", "pytest", "mocha", "jest"}

		for _, framework := range validFrameworks {
			t.Run("Valid_"+framework, func(t *testing.T) {
				tempDir := t.TempDir()
				configContent := `
project:
  name: "test-project"
  module: "github.com/test/project"
testing:
  framework: "` + framework + `"
`
				configPath := filepath.Join(tempDir, ".goca.yaml")
				if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
					t.Fatalf("Failed to write config: %v", err)
				}

				manager := cmd.NewConfigManager()
				err := manager.LoadConfig(tempDir)

				if err != nil {
					t.Errorf("Expected validation to pass for framework '%s' but got error: %v", framework, err)
				}
			})
		}

		for _, framework := range invalidFrameworks {
			t.Run("Invalid_"+framework, func(t *testing.T) {
				tempDir := t.TempDir()
				configContent := `
project:
  name: "test-project"
  module: "github.com/test/project"
testing:
  framework: "` + framework + `"
`
				configPath := filepath.Join(tempDir, ".goca.yaml")
				if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
					t.Fatalf("Failed to write config: %v", err)
				}

				manager := cmd.NewConfigManager()
				err := manager.LoadConfig(tempDir)

				if err == nil {
					t.Errorf("Expected validation to fail for invalid framework '%s' but it passed", framework)
				}
			})
		}
	})

	t.Run("AuthTypeValidation", func(t *testing.T) {
		validAuthTypes := []string{"jwt", "oauth2", "session", "basic"}
		// Empty auth type with enabled: true now gets default "jwt" applied
		invalidAuthTypes := []string{"saml", "openid", "ldap", "kerberos"}

		for _, authType := range validAuthTypes {
			t.Run("Valid_"+authType, func(t *testing.T) {
				tempDir := t.TempDir()
				configContent := `
project:
  name: "test-project"
  module: "github.com/test/project"
features:
  auth:
    enabled: true
    type: "` + authType + `"
`
				configPath := filepath.Join(tempDir, ".goca.yaml")
				if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
					t.Fatalf("Failed to write config: %v", err)
				}

				manager := cmd.NewConfigManager()
				err := manager.LoadConfig(tempDir)

				if err != nil {
					t.Errorf("Expected validation to pass for auth type '%s' but got error: %v", authType, err)
				}
			})
		}

		for _, authType := range invalidAuthTypes {
			t.Run("Invalid_"+authType, func(t *testing.T) {
				tempDir := t.TempDir()
				configContent := `
project:
  name: "test-project"
  module: "github.com/test/project"
features:
  auth:
    enabled: true
    type: "` + authType + `"
`
				configPath := filepath.Join(tempDir, ".goca.yaml")
				if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
					t.Fatalf("Failed to write config: %v", err)
				}

				manager := cmd.NewConfigManager()
				err := manager.LoadConfig(tempDir)

				if err == nil {
					t.Errorf("Expected validation to fail for invalid auth type '%s' but it passed", authType)
				}
			})
		}
	})

	t.Run("CacheTypeValidation", func(t *testing.T) {
		validCacheTypes := []string{"redis", "memcached", "inmemory"}
		// Empty cache type with enabled: true now gets default "redis" applied
		invalidCacheTypes := []string{"hazelcast", "ehcache", "varnish"}

		for _, cacheType := range validCacheTypes {
			t.Run("Valid_"+cacheType, func(t *testing.T) {
				tempDir := t.TempDir()
				configContent := `
project:
  name: "test-project"
  module: "github.com/test/project"
features:
  cache:
    enabled: true
    type: "` + cacheType + `"
`
				configPath := filepath.Join(tempDir, ".goca.yaml")
				if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
					t.Fatalf("Failed to write config: %v", err)
				}

				manager := cmd.NewConfigManager()
				err := manager.LoadConfig(tempDir)

				if err != nil {
					t.Errorf("Expected validation to pass for cache type '%s' but got error: %v", cacheType, err)
				}
			})
		}

		for _, cacheType := range invalidCacheTypes {
			t.Run("Invalid_"+cacheType, func(t *testing.T) {
				tempDir := t.TempDir()
				configContent := `
project:
  name: "test-project"
  module: "github.com/test/project"
features:
  cache:
    enabled: true
    type: "` + cacheType + `"
`
				configPath := filepath.Join(tempDir, ".goca.yaml")
				if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
					t.Fatalf("Failed to write config: %v", err)
				}

				manager := cmd.NewConfigManager()
				err := manager.LoadConfig(tempDir)

				if err == nil {
					t.Errorf("Expected validation to fail for invalid cache type '%s' but it passed", cacheType)
				}
			})
		}
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		testCases := []struct {
			name    string
			content string
		}{
			{
				name: "MissingProjectName",
				content: `
project:
  module: "github.com/test/project"
`,
			},
			{
				name: "MissingProjectModule",
				content: `
project:
  name: "test-project"
`,
			},
			{
				name: "EmptyProjectName",
				content: `
project:
  name: ""
  module: "github.com/test/project"
`,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, ".goca.yaml")
				if err := os.WriteFile(configPath, []byte(tc.content), 0644); err != nil {
					t.Fatalf("Failed to write config: %v", err)
				}

				manager := cmd.NewConfigManager()
				err := manager.LoadConfig(tempDir)

				if err == nil {
					t.Errorf("Expected validation to fail for %s but it passed", tc.name)
				}
			})
		}
	})

	t.Run("ComplexNestedConfiguration", func(t *testing.T) {
		tempDir := t.TempDir()
		configContent := `
project:
  name: "complex-test-project"
  module: "github.com/test/complex"
  version: "2.5.0"
  license: "Apache-2.0"
  tags:
    - "microservices"
    - "api"
    - "clean-architecture"
  metadata:
    team: "backend"
    environment: "production"

architecture:
  layers:
    domain:
      enabled: true
      directory: "internal/domain"
      patterns: ["factory", "builder"]
    usecase:
      enabled: true
      directory: "internal/usecase"
      validations: ["business", "authorization"]
    repository:
      enabled: true
      directory: "internal/repository"
      templates: ["postgres", "redis"]
    handler:
      enabled: true
      directory: "internal/handler"
  patterns:
    - "repository"
    - "factory"
    - "decorator"
  di:
    type: "wire"
    auto_wire: true
    providers: ["database", "logger", "metrics"]
  naming:
    entities: "PascalCase"
    fields: "camelCase"
    files: "snake_case"
    packages: "lowercase"
    constants: "UPPER_CASE"

database:
  type: "postgres"
  port: 5432
  name: "complex_db"
  migrations:
    enabled: true
    auto_generate: true
    directory: "migrations"
    naming: "timestamp"
  connection:
    max_open: 50
    max_idle: 25
    ssl_mode: "require"
    timezone: "America/New_York"
  features:
    soft_delete: true
    timestamps: true
    uuid: true
    audit: true

generation:
  validation:
    enabled: true
    library: "validator"
    tags: ["required", "email", "min", "max", "oneof"]
  business_rules:
    enabled: true
    patterns: ["saga", "event-sourcing"]
    events: true
    guards: true
  documentation:
    swagger:
      enabled: true
      version: "3.0.0"
      title: "Complex API"
      host: "api.example.com"
      schemes: ["https"]
    comments:
      enabled: true
      language: "english"
      style: "godoc"
  style:
    gofmt: true
    goimports: true
    golint: true
    govet: true
    line_length: 120

testing:
  enabled: true
  framework: "testify"
  integration: true
  benchmarks: true
  coverage:
    enabled: true
    threshold: 85.0
    format: "html"
  mocks:
    enabled: true
    tool: "testify"

features:
  auth:
    enabled: true
    type: "jwt"
    rbac: true
    middleware: true
  cache:
    enabled: true
    type: "redis"
    ttl: "1h"
  logging:
    enabled: true
    level: "info"
    format: "json"
    structured: true
  security:
    https: true
    cors: true
    rate_limit: true

deploy:
  docker:
    enabled: true
    multistage: true
  kubernetes:
    enabled: true
    namespace: "production"
`
		configPath := filepath.Join(tempDir, ".goca.yaml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		manager := cmd.NewConfigManager()
		err := manager.LoadConfig(tempDir)

		if err != nil {
			t.Fatalf("Expected complex nested configuration to be valid but got error: %v", err)
		}

		config := manager.GetConfig()
		if config == nil {
			t.Fatal("Expected config to be loaded but got nil")
		}

		// Verify complex nested values
		if config.Project.Name != "complex-test-project" {
			t.Errorf("Expected project name 'complex-test-project' but got '%s'", config.Project.Name)
		}

		if config.Database.Connection.MaxOpen != 50 {
			t.Errorf("Expected max_open 50 but got %d", config.Database.Connection.MaxOpen)
		}

		if !config.Features.Auth.RBAC {
			t.Error("Expected RBAC to be enabled")
		}

		if config.Testing.Coverage.Threshold != 85.0 {
			t.Errorf("Expected coverage threshold 85.0 but got %.1f", config.Testing.Coverage.Threshold)
		}
	})
}
