package cmd

import (
	"fmt"
)

// ProjectTemplate defines a predefined project configuration
type ProjectTemplate struct {
	Name        string
	Description string
	Config      string // YAML configuration
	Features    []string
}

// GetProjectTemplates returns all available project templates
func GetProjectTemplates() map[string]ProjectTemplate {
	return map[string]ProjectTemplate{
		"rest-api": {
			Name:        "REST API",
			Description: "Production-ready REST API with PostgreSQL, validation, and testing",
			Config: `# REST API Project Configuration
# Optimized for building RESTful web services

project:
  description: REST API built with Clean Architecture

database:
  type: postgres
  migrations:
    enabled: true
    auto_generate: true
  features:
    soft_delete: true
    timestamps: true
    uuid: false

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
  
  naming:
    files: lowercase
    entities: PascalCase
    variables: camelCase
    functions: PascalCase

generation:
  validation:
    enabled: true
    library: builtin
    sanitize: true
  
  business_rules:
    enabled: true
    patterns:
      - validation
  
  documentation:
    swagger:
      enabled: true
      version: "2.0"
      title: "API Documentation"
    comments:
      enabled: true
      language: english
      style: godoc

testing:
  enabled: true
  framework: testify
  coverage:
    enabled: true
    threshold: 70
  mocks:
    enabled: true
    tool: testify
  integration: true
  benchmarks: false
`,
			Features: []string{},
		},
		"microservice": {
			Name:        "Microservice",
			Description: "Microservice with gRPC, events, and comprehensive testing",
			Config: `# Microservice Project Configuration
# Optimized for distributed systems and event-driven architecture

project:
  description: Microservice built with Clean Architecture

database:
  type: postgres
  migrations:
    enabled: true
    auto_generate: true
  features:
    soft_delete: false
    timestamps: true
    uuid: true
    audit: true

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
    - specification
  
  naming:
    files: lowercase
    entities: PascalCase
    variables: camelCase
    functions: PascalCase

generation:
  validation:
    enabled: true
    library: validator
    sanitize: true
  
  business_rules:
    enabled: true
    patterns:
      - validation
      - authorization
    events: true
  
  documentation:
    swagger:
      enabled: true
    comments:
      enabled: true
      language: english

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
  benchmarks: true
`,
			Features: []string{},
		},
		"monolith": {
			Name:        "Monolith",
			Description: "Full-featured monolithic application with web interface",
			Config: `# Monolith Project Configuration
# Optimized for traditional web applications

project:
  description: Monolithic application built with Clean Architecture

database:
  type: postgres
  migrations:
    enabled: true
    auto_generate: true
  features:
    soft_delete: true
    timestamps: true
    uuid: false
    audit: true
    versioning: true

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
  
  naming:
    files: lowercase
    entities: PascalCase
    variables: camelCase
    functions: PascalCase

generation:
  validation:
    enabled: true
    library: builtin
    sanitize: true
    transform: true
  
  business_rules:
    enabled: true
    patterns:
      - validation
      - authorization
    guards: true
  
  documentation:
    swagger:
      enabled: true
      version: "2.0"
    markdown:
      enabled: true
      toc: true
    comments:
      enabled: true
      language: english

testing:
  enabled: true
  framework: testify
  coverage:
    enabled: true
    threshold: 75
  mocks:
    enabled: true
  integration: true
  benchmarks: false
  fixtures:
    enabled: true
    seeds: true

features:
  auth:
    enabled: true
    type: jwt
    rbac: true
  
  cache:
    enabled: true
    type: redis
  
  logging:
    enabled: true
    level: info
    format: json
    structured: true
  
  monitoring:
    enabled: true
    metrics: true
    health_check: true
`,
			Features: []string{},
		},
		"minimal": {
			Name:        "Minimal",
			Description: "Lightweight starter with essential features only",
			Config: `# Minimal Project Configuration
# Essential features only for quick prototyping

project:
  description: Minimal Clean Architecture project

database:
  type: postgres

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

testing:
  enabled: true
  framework: testify
`,
			Features: []string{},
		},
		"enterprise": {
			Name:        "Enterprise",
			Description: "Enterprise-grade with all features, security, and monitoring",
			Config: `# Enterprise Project Configuration
# Production-ready with comprehensive features

project:
  description: Enterprise application built with Clean Architecture
  version: 1.0.0

database:
  type: postgres
  migrations:
    enabled: true
    auto_generate: true
    versioning: timestamp
  connection:
    max_open: 100
    max_idle: 10
  features:
    soft_delete: true
    timestamps: true
    uuid: true
    audit: true
    versioning: true
    partitioning: false

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
    - specification
  
  naming:
    files: lowercase
    entities: PascalCase
    variables: camelCase
    functions: PascalCase
    constants: SCREAMING_SNAKE

generation:
  validation:
    enabled: true
    library: validator
    sanitize: true
    transform: true
  
  business_rules:
    enabled: true
    patterns:
      - validation
      - authorization
    events: true
    guards: true
  
  documentation:
    swagger:
      enabled: true
      version: "3.0"
      title: "Enterprise API"
    postman:
      enabled: true
      environment: true
      tests: true
    markdown:
      enabled: true
      toc: true
      examples: true
    comments:
      enabled: true
      language: english
      style: godoc
      examples: true

  style:
    gofmt: true
    goimports: true
    golint: true
    staticcheck: true

testing:
  enabled: true
  framework: testify
  coverage:
    enabled: true
    threshold: 85
    format: html
  mocks:
    enabled: true
    tool: testify
  integration: true
  benchmarks: true
  examples: true
  fixtures:
    enabled: true
    seeds: true
    factories: 
      - user
      - admin

templates:
  directory: .goca/templates

features:
  auth:
    enabled: true
    type: jwt
    providers:
      - local
      - oauth2
    rbac: true
    middleware: true
  
  cache:
    enabled: true
    type: redis
    ttl: 1h
    layers:
      - query
      - entity
  
  logging:
    enabled: true
    level: info
    format: json
    output:
      - stdout
      - file
    structured: true
    tracing: true
  
  monitoring:
    enabled: true
    metrics: true
    tracing: true
    health_check: true
    profiling: true
    tools:
      - prometheus
  
  security:
    https: true
    cors: true
    rate_limit: true
    validation: true
    sanitization: true
    headers:
      - X-Content-Type-Options
      - X-Frame-Options
      - X-XSS-Protection

deploy:
  docker:
    enabled: true
    multistage: true
    compose: true
  
  kubernetes:
    enabled: true
    manifests: k8s
    helm: true
    ingress: true
    config_maps: true
    secrets: true
  
  ci:
    enabled: true
    provider: github-actions
    workflows:
      - test
      - build
      - deploy
    tests: true
    build: true
    deploy: true
`,
			Features: []string{},
		},
	}
}

// GetTemplateNames returns list of available template names
func GetTemplateNames() []string {
	templates := GetProjectTemplates()
	names := make([]string, 0, len(templates))
	for name := range templates {
		names = append(names, name)
	}
	return names
}

// ValidateTemplateName checks if a template name is valid
func ValidateTemplateName(name string) bool {
	if name == "" {
		return true // Empty means no template
	}
	templates := GetProjectTemplates()
	_, exists := templates[name]
	return exists
}

// GetTemplateConfig returns the configuration for a template
func GetTemplateConfig(name string) (string, error) {
	templates := GetProjectTemplates()
	template, exists := templates[name]
	if !exists {
		return "", fmt.Errorf("template '%s' not found", name)
	}
	return template.Config, nil
}

// ListAvailableTemplates prints all available templates with descriptions
func ListAvailableTemplates() {
	templates := GetProjectTemplates()
	fmt.Println("\nAvailable project templates:")
	fmt.Println()

	// Define order for consistent output
	order := []string{"minimal", "rest-api", "microservice", "monolith", "enterprise"}

	for _, name := range order {
		if template, exists := templates[name]; exists {
			fmt.Printf("  %s\n", name)
			fmt.Printf("    %s\n", template.Description)
			fmt.Println()
		}
	}

	fmt.Println("Usage:")
	fmt.Println("  goca init myproject --module github.com/user/myproject --template rest-api")
	fmt.Println()
}
