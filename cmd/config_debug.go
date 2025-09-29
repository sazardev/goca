package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// configCmd command group for configuration management
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage GOCA configuration",
	Long: `Manage GOCA YAML configuration files. 

This command provides subcommands to initialize, validate, and manage
your .goca.yaml configuration files for consistent project setup.

Available subcommands:
  show      - Display current configuration  
  init      - Initialize new configuration file
  validate  - Validate configuration file
  template  - Manage template configuration

If no subcommand is provided, defaults to 'show'.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default to 'show' behavior when no subcommand is provided
		showCurrentConfig()
	},
}

// configShowCmd shows current configuration
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current project configuration loaded from .goca.yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		showCurrentConfig()
	},
}

// configInitCmd initializes a new configuration file
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long: `Initialize a new .goca.yaml configuration file in the current directory.

This creates a comprehensive configuration file with intelligent defaults
based on your project structure and specified options.`,
	Run: func(cmd *cobra.Command, args []string) {
		initializeConfig(cmd)
	},
}

// configValidateCmd validates the current configuration
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration file",
	Long:  `Validate the current .goca.yaml configuration file for errors and warnings.`,
	Run: func(cmd *cobra.Command, args []string) {
		validateConfiguration()
	},
}

// configTemplateCmd manages template configurations
var configTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage configuration templates",
	Long:  `Manage predefined configuration templates for different project types.`,
	Run: func(cmd *cobra.Command, args []string) {
		showTemplateOptions()
	},
}

func init() {
	// Add subcommands to config command
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configTemplateCmd)

	// Add flags for init command
	configInitCmd.Flags().String("template", "", "Use predefined template (web, api, microservice, full)")
	configInitCmd.Flags().Bool("force", false, "Overwrite existing file")
	configInitCmd.Flags().String("database", "", "Database type (postgres, mysql, sqlite)")
	configInitCmd.Flags().StringSlice("handlers", []string{}, "Handler types (http, grpc, cli)")

	// Add to root command
	rootCmd.AddCommand(configCmd)
}

// showCurrentConfig displays the current configuration
func showCurrentConfig() {
	fmt.Println("=== Current GOCA Configuration ===")

	// Try to load existing config
	configPath := ".goca.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("❌ Configuration file not found: %s\n", configPath)
		fmt.Println("💡 Run 'goca config init' to create a new one")
		return
	}

	// Read and display config
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("❌ Error reading configuration: %v\n", err)
		return
	}

	// Parse YAML to validate structure
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("❌ Invalid YAML file: %v\n", err)
		return
	}

	fmt.Printf("✅ File found: %s\n", configPath)
	fmt.Printf("📁 Directory: %s\n", getCurrentDir())
	fmt.Println("\n--- Content ---")
	fmt.Println(string(data))

	// Show validation status
	fmt.Println("\n--- Validation ---")
	validateConfigSilent(config)
}

// initializeConfig creates a new configuration file
func initializeConfig(cmd *cobra.Command) {
	template, _ := cmd.Flags().GetString("template")
	force, _ := cmd.Flags().GetBool("force")
	database, _ := cmd.Flags().GetString("database")
	handlers, _ := cmd.Flags().GetStringSlice("handlers")

	configPath := ".goca.yaml"

	// Check if file exists
	if _, err := os.Stat(configPath); err == nil && !force {
		fmt.Printf("❌ File %s already exists. Use --force to overwrite.\n", configPath)
		return
	}

	fmt.Println("🚀 Initializing GOCA configuration...")

	// Generate config content based on template
	var configContent string
	switch template {
	case "web":
		configContent = generateWebTemplate(database, handlers)
	case "api":
		configContent = generateAPITemplate(database, handlers)
	case "microservice":
		configContent = generateMicroserviceTemplate(database, handlers)
	case "full":
		configContent = generateFullTemplate(database, handlers)
	default:
		configContent = generateDefaultTemplate(database, handlers)
	}

	// Write config file
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		fmt.Printf("❌ Error writing configuration: %v\n", err)
		return
	}

	fmt.Printf("✅ Configuration file created: %s\n", configPath)
	if template != "" {
		fmt.Printf("📋 Template applied: %s\n", template)
	}
	fmt.Println("💡 Run 'goca config show' to view the configuration")
}

// validateConfiguration validates the current config file
func validateConfiguration() {
	fmt.Println("🔍 Validating configuration...")

	configPath := ".goca.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("❌ File not found: %s\n", configPath)
		return
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("❌ Error reading file: %v\n", err)
		return
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("❌ Invalid YAML: %v\n", err)
		return
	}

	// Validate structure
	errors := validateConfigStructure(config)
	if len(errors) == 0 {
		fmt.Println("✅ Configuration is valid")
	} else {
		fmt.Printf("❌ Found %d errors:\n", len(errors))
		for i, err := range errors {
			fmt.Printf("  %d. %s\n", i+1, err)
		}
	}
}

// showTemplateOptions shows available configuration templates
func showTemplateOptions() {
	fmt.Println("📋 Available Configuration Templates")
	fmt.Println("===================================")

	templates := map[string]string{
		"web":          "Complete web application with frontend and backend",
		"api":          "REST API with database",
		"microservice": "Microservice with multiple handlers",
		"full":         "Complete configuration with all features",
		"default":      "Basic minimal configuration",
	}

	for name, desc := range templates {
		fmt.Printf("• %s\n  %s\n\n", name, desc)
	}

	fmt.Println("Usage:")
	fmt.Println("  goca config init --template <name>")
	fmt.Println("  goca config init --template web --database postgres --handlers http,grpc")
}

// Helper functions for config generation
func generateDefaultTemplate(database string, handlers []string) string {
	db := database
	if db == "" {
		db = "postgres"
	}

	handlersStr := strings.Join(handlers, ", ")
	if handlersStr == "" {
		handlersStr = "http"
	}

	return fmt.Sprintf(`# GOCA Configuration File  
# Configuration file for projects generated with GOCA CLI

# Project configuration
project:
  name: "%s"
  module: "%s"
  version: "1.0.0"
  description: "Project generated with GOCA CLI"
  author: ""
  license: "MIT"

# Architecture configuration
architecture:
  pattern: "clean_architecture"
  layers:
    domain:
      enabled: true
      path: "internal/domain"
    usecase:
      enabled: true 
      path: "internal/usecase"
    repository:
      enabled: true
      path: "internal/repository"
    handler:
      enabled: true
      path: "internal/handler"

# Database configuration
database:
  type: "%s"
  host: "localhost"
  port: 5432
  name: "%s_db"
  migrations:
    enabled: true
    auto_generate: true
    directory: "migrations"
    naming: "timestamp"
    versioning: "sequential"
    tools: ["migrate", "sql-migrate"]
  connection:
    max_open: 25
    max_idle: 5
    max_lifetime: "5m"
    ssl_mode: "disable"
    timezone: "UTC"
  features:
    soft_delete: true
    timestamps: true
    uuid: false
    audit: false
    versioning: false
    partitioning: false
    indexes: []
    constraints: []

# Code generation configuration  
generation:
  validation:
    enabled: true
    tags: ["required", "min", "max", "email", "url"]
    custom_rules: true
    error_handling: "detailed"
    localization: "english"
  business_rules:
    enabled: false
    directory: "internal/domain/rules"
    naming: "rule"
    testing: true
  docker:
    enabled: true
    compose: true
    dockerfile: true
    multi_stage: true
  docs:
    swagger:
      enabled: true
      title: "%s API"
      version: "1.0.0"
      description: "API generada con GOCA CLI"
      host: "localhost:8080"
      base_path: "/api/v1"
      schemes: ["http", "https"]
      tags: []
    postman:
      enabled: true
      output: "docs/postman"
      environment: true
      tests: true
    markdown:
      enabled: true
      output: "docs"
      template: "default"
      toc: true
      examples: true
  comments:
    enabled: true
    language: "english"
    style: "godoc"
    examples: true
    todo: true
    deprecated: true

# Quality configuration
quality:
  style:
    gofmt: true
    goimports: true
    golint: true
    govet: true
    staticcheck: true
    line_length: 120
    tab_width: 4
  testing:
    enabled: true
    coverage_threshold: 80.0
    benchmark: true
    integration: true
    e2e: false
    mocks: true
  security:
    enabled: true
    scanner: "gosec"
    rules: ["all"]
    exclude: []

# Infrastructure configuration
infrastructure:
  logging:
    enabled: true
    level: "info"
    format: "structured"
    output: ["stdout", "file"]
    structured: true
    tracing: false
  monitoring:
    enabled: false
    metrics: false
    tracing: false
    health_check: true
    profiling: false
    tools: []
  cache:
    enabled: false
    type: "redis"
    ttl: "1h"
    max_size: "100MB"
  message_queue:
    enabled: false
    type: "rabbitmq"
    exchanges: []
    queues: []
  deployment:
    type: "docker"
    registry: ""
    namespace: ""
    replicas: 1
    resources:
      cpu: "100m"
      memory: "128Mi"

# Default configuration for CLI commands
defaults:
  database: "%s"
  handlers: [%s]
  validation: true
  soft_delete: true
  timestamps: true
  business_rules: false
  docker: true
  docs: true
  testing: true
`, getCurrentProjectName(), getCurrentModuleName(), db, getCurrentProjectName(), getCurrentProjectName(), db, handlersStr)
}

func generateWebTemplate(database string, handlers []string) string {
	// Similar implementation for web template
	return generateDefaultTemplate(database, handlers) // Simplified for now
}

func generateAPITemplate(database string, handlers []string) string {
	// Similar implementation for API template
	return generateDefaultTemplate(database, handlers) // Simplified for now
}

func generateMicroserviceTemplate(database string, handlers []string) string {
	// Similar implementation for microservice template
	return generateDefaultTemplate(database, handlers) // Simplified for now
}

func generateFullTemplate(database string, handlers []string) string {
	// Similar implementation for full template
	return generateDefaultTemplate(database, handlers) // Simplified for now
}

// Helper validation functions
func validateConfigSilent(config map[string]interface{}) {
	errors := validateConfigStructure(config)
	if len(errors) == 0 {
		fmt.Println("✅ Valid structure")
	} else {
		fmt.Printf("⚠️  %d warnings found\n", len(errors))
	}
}

func validateConfigStructure(config map[string]interface{}) []string {
	var errors []string

	// Check required sections
	requiredSections := []string{"project", "defaults"}
	for _, section := range requiredSections {
		if _, exists := config[section]; !exists {
			errors = append(errors, fmt.Sprintf("Required section '%s' missing", section))
		}
	}

	return errors
}

func getCurrentDir() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	return "unknown"
}

func getCurrentProjectName() string {
	dir := getCurrentDir()
	return filepath.Base(dir)
}

func getCurrentModuleName() string {
	projectName := getCurrentProjectName()
	return fmt.Sprintf("github.com/usuario/%s", projectName)
}
