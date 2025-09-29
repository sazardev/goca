package cmd

import (
	"time"
)

// GocaConfig represents the complete configuration structure for a Goca project
type GocaConfig struct {
	// Core project configuration
	Project ProjectConfig `yaml:"project" json:"project"`

	// Architecture and patterns configuration
	Architecture ArchitectureConfig `yaml:"architecture" json:"architecture"`

	// Database configuration
	Database DatabaseConfig `yaml:"database" json:"database"`

	// Code generation preferences
	Generation GenerationConfig `yaml:"generation" json:"generation"`

	// Testing configuration
	Testing TestingConfig `yaml:"testing" json:"testing"`

	// Templates and customization
	Templates TemplateConfig `yaml:"templates" json:"templates"`

	// Features and plugins
	Features FeatureConfig `yaml:"features" json:"features"`

	// Deployment and infrastructure
	Deploy DeployConfig `yaml:"deploy" json:"deploy"`
}

// ProjectConfig contains basic project information
type ProjectConfig struct {
	Name        string            `yaml:"name" json:"name"`
	Module      string            `yaml:"module" json:"module"`
	Description string            `yaml:"description" json:"description"`
	Version     string            `yaml:"version" json:"version"`
	Author      string            `yaml:"author" json:"author"`
	License     string            `yaml:"license" json:"license"`
	Repository  string            `yaml:"repository" json:"repository"`
	Tags        []string          `yaml:"tags" json:"tags"`
	Metadata    map[string]string `yaml:"metadata" json:"metadata"`
}

// ArchitectureConfig defines Clean Architecture preferences
type ArchitectureConfig struct {
	// Layers configuration
	Layers LayersConfig `yaml:"layers" json:"layers"`

	// Patterns to apply
	Patterns []string `yaml:"patterns" json:"patterns"`

	// Dependency injection type
	DI DIConfig `yaml:"di" json:"di"`

	// Naming conventions
	Naming NamingConfig `yaml:"naming" json:"naming"`
}

// LayersConfig defines which layers to generate and their structure
type LayersConfig struct {
	Domain     LayerConfig   `yaml:"domain" json:"domain"`
	UseCase    LayerConfig   `yaml:"usecase" json:"usecase"`
	Repository LayerConfig   `yaml:"repository" json:"repository"`
	Handler    LayerConfig   `yaml:"handler" json:"handler"`
	Custom     []LayerConfig `yaml:"custom" json:"custom"`
}

// LayerConfig defines configuration for a specific layer
type LayerConfig struct {
	Enabled     bool              `yaml:"enabled" json:"enabled"`
	Directory   string            `yaml:"directory" json:"directory"`
	Patterns    []string          `yaml:"patterns" json:"patterns"`
	Templates   []string          `yaml:"templates" json:"templates"`
	Validations []string          `yaml:"validations" json:"validations"`
	Extensions  map[string]string `yaml:"extensions" json:"extensions"`
}

// DIConfig defines dependency injection preferences
type DIConfig struct {
	Type       string            `yaml:"type" json:"type"` // manual, wire, fx, dig
	AutoWire   bool              `yaml:"auto_wire" json:"auto_wire"`
	Providers  []string          `yaml:"providers" json:"providers"`
	Modules    []string          `yaml:"modules" json:"modules"`
	Extensions map[string]string `yaml:"extensions" json:"extensions"`
}

// NamingConfig defines naming conventions
type NamingConfig struct {
	Entities  string `yaml:"entities" json:"entities"`   // PascalCase, camelCase, snake_case
	Fields    string `yaml:"fields" json:"fields"`       // PascalCase, camelCase, snake_case
	Files     string `yaml:"files" json:"files"`         // snake_case, kebab-case, camelCase
	Packages  string `yaml:"packages" json:"packages"`   // lowercase, snake_case
	Constants string `yaml:"constants" json:"constants"` // UPPER_CASE, PascalCase
	Variables string `yaml:"variables" json:"variables"` // camelCase, snake_case
	Functions string `yaml:"functions" json:"functions"` // camelCase, PascalCase
}

// DatabaseConfig contains database configuration
type DatabaseConfig struct {
	Type        string                `yaml:"type" json:"type"` // postgres, mysql, mongodb, sqlite
	Host        string                `yaml:"host" json:"host"`
	Port        int                   `yaml:"port" json:"port"`
	Name        string                `yaml:"name" json:"name"`
	Migrations  MigrationConfig       `yaml:"migrations" json:"migrations"`
	Connection  ConnectionConfig      `yaml:"connection" json:"connection"`
	Features    DatabaseFeatureConfig `yaml:"features" json:"features"`
	Extensions  []string              `yaml:"extensions" json:"extensions"`
	CustomTypes map[string]string     `yaml:"custom_types" json:"custom_types"`
}

// MigrationConfig defines migration preferences
type MigrationConfig struct {
	Enabled      bool     `yaml:"enabled" json:"enabled"`
	AutoGenerate bool     `yaml:"auto_generate" json:"auto_generate"`
	Directory    string   `yaml:"directory" json:"directory"`
	Naming       string   `yaml:"naming" json:"naming"`
	Versioning   string   `yaml:"versioning" json:"versioning"`
	Tools        []string `yaml:"tools" json:"tools"`
}

// ConnectionConfig defines database connection settings
type ConnectionConfig struct {
	MaxOpen     int           `yaml:"max_open" json:"max_open"`
	MaxIdle     int           `yaml:"max_idle" json:"max_idle"`
	MaxLifetime time.Duration `yaml:"max_lifetime" json:"max_lifetime"`
	SSLMode     string        `yaml:"ssl_mode" json:"ssl_mode"`
	Timezone    string        `yaml:"timezone" json:"timezone"`
	Charset     string        `yaml:"charset" json:"charset"`
	Collation   string        `yaml:"collation" json:"collation"`
}

// DatabaseFeatureConfig defines database-specific features
type DatabaseFeatureConfig struct {
	SoftDelete   bool     `yaml:"soft_delete" json:"soft_delete"`
	Timestamps   bool     `yaml:"timestamps" json:"timestamps"`
	UUID         bool     `yaml:"uuid" json:"uuid"`
	Audit        bool     `yaml:"audit" json:"audit"`
	Versioning   bool     `yaml:"versioning" json:"versioning"`
	Partitioning bool     `yaml:"partitioning" json:"partitioning"`
	Indexes      []string `yaml:"indexes" json:"indexes"`
	Constraints  []string `yaml:"constraints" json:"constraints"`
}

// GenerationConfig defines code generation preferences
type GenerationConfig struct {
	// Field validation preferences
	Validation ValidationConfig `yaml:"validation" json:"validation"`

	// Business rules generation
	BusinessRules BusinessRulesConfig `yaml:"business_rules" json:"business_rules"`

	// Documentation generation
	Documentation DocumentationConfig `yaml:"documentation" json:"documentation"`

	// Code style preferences
	Style StyleConfig `yaml:"style" json:"style"`

	// Import management
	Imports ImportConfig `yaml:"imports" json:"imports"`
}

// ValidationConfig defines validation generation preferences
type ValidationConfig struct {
	Enabled   bool     `yaml:"enabled" json:"enabled"`
	Library   string   `yaml:"library" json:"library"` // builtin, validator, ozzo-validation
	Tags      []string `yaml:"tags" json:"tags"`
	Custom    []string `yaml:"custom" json:"custom"`
	Sanitize  bool     `yaml:"sanitize" json:"sanitize"`
	Transform bool     `yaml:"transform" json:"transform"`
}

// BusinessRulesConfig defines business rules generation
type BusinessRulesConfig struct {
	Enabled   bool     `yaml:"enabled" json:"enabled"`
	Patterns  []string `yaml:"patterns" json:"patterns"`
	Templates []string `yaml:"templates" json:"templates"`
	Events    bool     `yaml:"events" json:"events"`
	Guards    bool     `yaml:"guards" json:"guards"`
}

// DocumentationConfig defines documentation generation
type DocumentationConfig struct {
	Swagger  SwaggerConfig  `yaml:"swagger" json:"swagger"`
	Postman  PostmanConfig  `yaml:"postman" json:"postman"`
	Markdown MarkdownConfig `yaml:"markdown" json:"markdown"`
	Comments CommentsConfig `yaml:"comments" json:"comments"`
}

// SwaggerConfig defines Swagger/OpenAPI generation
type SwaggerConfig struct {
	Enabled     bool              `yaml:"enabled" json:"enabled"`
	Version     string            `yaml:"version" json:"version"`
	Output      string            `yaml:"output" json:"output"`
	Title       string            `yaml:"title" json:"title"`
	Description string            `yaml:"description" json:"description"`
	Host        string            `yaml:"host" json:"host"`
	BasePath    string            `yaml:"base_path" json:"base_path"`
	Schemes     []string          `yaml:"schemes" json:"schemes"`
	Tags        []SwaggerTag      `yaml:"tags" json:"tags"`
	Extensions  map[string]string `yaml:"extensions" json:"extensions"`
}

// SwaggerTag defines Swagger tag configuration
type SwaggerTag struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
}

// PostmanConfig defines Postman collection generation
type PostmanConfig struct {
	Enabled     bool   `yaml:"enabled" json:"enabled"`
	Output      string `yaml:"output" json:"output"`
	Environment bool   `yaml:"environment" json:"environment"`
	Tests       bool   `yaml:"tests" json:"tests"`
	Variables   bool   `yaml:"variables" json:"variables"`
}

// MarkdownConfig defines Markdown documentation generation
type MarkdownConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Output   string `yaml:"output" json:"output"`
	Template string `yaml:"template" json:"template"`
	TOC      bool   `yaml:"toc" json:"toc"`
	Examples bool   `yaml:"examples" json:"examples"`
}

// CommentsConfig defines code comments generation
type CommentsConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled"`
	Language   string `yaml:"language" json:"language"` // spanish, english
	Style      string `yaml:"style" json:"style"`       // godoc, standard
	Examples   bool   `yaml:"examples" json:"examples"`
	TODO       bool   `yaml:"todo" json:"todo"`
	Deprecated bool   `yaml:"deprecated" json:"deprecated"`
}

// StyleConfig defines code style preferences
type StyleConfig struct {
	Gofmt       bool     `yaml:"gofmt" json:"gofmt"`
	Goimports   bool     `yaml:"goimports" json:"goimports"`
	Golint      bool     `yaml:"golint" json:"golint"`
	Govet       bool     `yaml:"govet" json:"govet"`
	Staticcheck bool     `yaml:"staticcheck" json:"staticcheck"`
	Custom      []string `yaml:"custom" json:"custom"`
	LineLength  int      `yaml:"line_length" json:"line_length"`
	TabWidth    int      `yaml:"tab_width" json:"tab_width"`
}

// ImportConfig defines import management
type ImportConfig struct {
	GroupStandard   bool    `yaml:"group_standard" json:"group_standard"`
	GroupThirdParty bool    `yaml:"group_third_party" json:"group_third_party"`
	GroupLocal      bool    `yaml:"group_local" json:"group_local"`
	SortAlpha       bool    `yaml:"sort_alpha" json:"sort_alpha"`
	RemoveUnused    bool    `yaml:"remove_unused" json:"remove_unused"`
	Aliases         []Alias `yaml:"aliases" json:"aliases"`
}

// Alias defines import alias configuration
type Alias struct {
	Package string `yaml:"package" json:"package"`
	Alias   string `yaml:"alias" json:"alias"`
}

// TestingConfig defines testing generation preferences
type TestingConfig struct {
	Enabled     bool           `yaml:"enabled" json:"enabled"`
	Framework   string         `yaml:"framework" json:"framework"` // testify, ginkgo, builtin
	Coverage    CoverageConfig `yaml:"coverage" json:"coverage"`
	Mocks       MockConfig     `yaml:"mocks" json:"mocks"`
	Integration bool           `yaml:"integration" json:"integration"`
	Benchmarks  bool           `yaml:"benchmarks" json:"benchmarks"`
	Examples    bool           `yaml:"examples" json:"examples"`
	Fixtures    FixtureConfig  `yaml:"fixtures" json:"fixtures"`
}

// CoverageConfig defines test coverage preferences
type CoverageConfig struct {
	Enabled   bool     `yaml:"enabled" json:"enabled"`
	Threshold float64  `yaml:"threshold" json:"threshold"`
	Output    string   `yaml:"output" json:"output"`
	Format    string   `yaml:"format" json:"format"`
	Exclude   []string `yaml:"exclude" json:"exclude"`
}

// MockConfig defines mock generation preferences
type MockConfig struct {
	Enabled    bool     `yaml:"enabled" json:"enabled"`
	Tool       string   `yaml:"tool" json:"tool"` // gomock, testify, counterfeiter
	Directory  string   `yaml:"directory" json:"directory"`
	Suffix     string   `yaml:"suffix" json:"suffix"`
	Interfaces []string `yaml:"interfaces" json:"interfaces"`
}

// FixtureConfig defines test fixture preferences
type FixtureConfig struct {
	Enabled   bool     `yaml:"enabled" json:"enabled"`
	Directory string   `yaml:"directory" json:"directory"`
	Format    string   `yaml:"format" json:"format"` // json, yaml, sql
	Seeds     bool     `yaml:"seeds" json:"seeds"`
	Factories []string `yaml:"factories" json:"factories"`
}

// TemplateConfig defines custom templates
type TemplateConfig struct {
	Directory string              `yaml:"directory" json:"directory"`
	Custom    map[string]Template `yaml:"custom" json:"custom"`
	Overrides map[string]Template `yaml:"overrides" json:"overrides"`
	Variables map[string]string   `yaml:"variables" json:"variables"`
}

// Template defines a custom template configuration
type Template struct {
	Path       string            `yaml:"path" json:"path"`
	Type       string            `yaml:"type" json:"type"`
	Variables  map[string]string `yaml:"variables" json:"variables"`
	Conditions []string          `yaml:"conditions" json:"conditions"`
}

// FeatureConfig defines feature flags and plugins
type FeatureConfig struct {
	// Authentication and authorization
	Auth AuthConfig `yaml:"auth" json:"auth"`

	// Caching configuration
	Cache CacheConfig `yaml:"cache" json:"cache"`

	// Logging configuration
	Logging LoggingConfig `yaml:"logging" json:"logging"`

	// Monitoring configuration
	Monitoring MonitoringConfig `yaml:"monitoring" json:"monitoring"`

	// Security features
	Security SecurityConfig `yaml:"security" json:"security"`

	// Plugins and extensions
	Plugins []PluginConfig `yaml:"plugins" json:"plugins"`
}

// AuthConfig defines authentication configuration
type AuthConfig struct {
	Enabled    bool     `yaml:"enabled" json:"enabled"`
	Type       string   `yaml:"type" json:"type"` // jwt, oauth2, session, basic
	Providers  []string `yaml:"providers" json:"providers"`
	RBAC       bool     `yaml:"rbac" json:"rbac"`
	Middleware bool     `yaml:"middleware" json:"middleware"`
}

// CacheConfig defines caching configuration
type CacheConfig struct {
	Enabled  bool     `yaml:"enabled" json:"enabled"`
	Type     string   `yaml:"type" json:"type"` // redis, memcached, inmemory
	TTL      string   `yaml:"ttl" json:"ttl"`
	Layers   []string `yaml:"layers" json:"layers"`
	Patterns []string `yaml:"patterns" json:"patterns"`
}

// LoggingConfig defines logging configuration
type LoggingConfig struct {
	Enabled    bool     `yaml:"enabled" json:"enabled"`
	Level      string   `yaml:"level" json:"level"`
	Format     string   `yaml:"format" json:"format"` // json, text, structured
	Output     []string `yaml:"output" json:"output"` // stdout, file, syslog
	Structured bool     `yaml:"structured" json:"structured"`
	Tracing    bool     `yaml:"tracing" json:"tracing"`
}

// MonitoringConfig defines monitoring and observability
type MonitoringConfig struct {
	Enabled     bool     `yaml:"enabled" json:"enabled"`
	Metrics     bool     `yaml:"metrics" json:"metrics"`
	Tracing     bool     `yaml:"tracing" json:"tracing"`
	HealthCheck bool     `yaml:"health_check" json:"health_check"`
	Profiling   bool     `yaml:"profiling" json:"profiling"`
	Tools       []string `yaml:"tools" json:"tools"` // prometheus, jaeger, datadog
}

// SecurityConfig defines security features
type SecurityConfig struct {
	HTTPS        bool     `yaml:"https" json:"https"`
	CORS         bool     `yaml:"cors" json:"cors"`
	RateLimit    bool     `yaml:"rate_limit" json:"rate_limit"`
	Validation   bool     `yaml:"validation" json:"validation"`
	Sanitization bool     `yaml:"sanitization" json:"sanitization"`
	Headers      []string `yaml:"headers" json:"headers"`
	Middleware   []string `yaml:"middleware" json:"middleware"`
}

// PluginConfig defines plugin configuration
type PluginConfig struct {
	Name     string            `yaml:"name" json:"name"`
	Version  string            `yaml:"version" json:"version"`
	Enabled  bool              `yaml:"enabled" json:"enabled"`
	Config   map[string]string `yaml:"config" json:"config"`
	Priority int               `yaml:"priority" json:"priority"`
}

// DeployConfig defines deployment and infrastructure
type DeployConfig struct {
	Docker      DockerConfig     `yaml:"docker" json:"docker"`
	Kubernetes  KubernetesConfig `yaml:"kubernetes" json:"kubernetes"`
	CI          CIConfig         `yaml:"ci" json:"ci"`
	Environment []EnvConfig      `yaml:"environments" json:"environments"`
}

// DockerConfig defines Docker configuration
type DockerConfig struct {
	Enabled    bool              `yaml:"enabled" json:"enabled"`
	Dockerfile string            `yaml:"dockerfile" json:"dockerfile"`
	Image      string            `yaml:"image" json:"image"`
	Registry   string            `yaml:"registry" json:"registry"`
	Compose    bool              `yaml:"compose" json:"compose"`
	Multistage bool              `yaml:"multistage" json:"multistage"`
	Labels     map[string]string `yaml:"labels" json:"labels"`
}

// KubernetesConfig defines Kubernetes configuration
type KubernetesConfig struct {
	Enabled    bool              `yaml:"enabled" json:"enabled"`
	Namespace  string            `yaml:"namespace" json:"namespace"`
	Manifests  string            `yaml:"manifests" json:"manifests"`
	Helm       bool              `yaml:"helm" json:"helm"`
	Ingress    bool              `yaml:"ingress" json:"ingress"`
	ConfigMaps bool              `yaml:"config_maps" json:"config_maps"`
	Secrets    bool              `yaml:"secrets" json:"secrets"`
	Labels     map[string]string `yaml:"labels" json:"labels"`
}

// CIConfig defines CI/CD configuration
type CIConfig struct {
	Enabled   bool     `yaml:"enabled" json:"enabled"`
	Provider  string   `yaml:"provider" json:"provider"` // github-actions, gitlab-ci, jenkins
	Workflows []string `yaml:"workflows" json:"workflows"`
	Tests     bool     `yaml:"tests" json:"tests"`
	Build     bool     `yaml:"build" json:"build"`
	Deploy    bool     `yaml:"deploy" json:"deploy"`
}

// EnvConfig defines environment-specific configuration
type EnvConfig struct {
	Name      string            `yaml:"name" json:"name"`
	Default   bool              `yaml:"default" json:"default"`
	Variables map[string]string `yaml:"variables" json:"variables"`
	Overrides GocaConfig        `yaml:"overrides" json:"overrides"`
}
