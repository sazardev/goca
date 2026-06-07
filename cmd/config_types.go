package cmd

import (
	"time"
)

// GocaConfig represents the complete configuration structure for a Goca project.
type GocaConfig struct {
	// Core project configuration
	Project ProjectConfig `json:"project" yaml:"project"`

	// Architecture and patterns configuration
	Architecture ArchitectureConfig `json:"architecture" yaml:"architecture"`

	// Database configuration
	Database DatabaseConfig `json:"database" yaml:"database"`

	// Code generation preferences
	Generation GenerationConfig `json:"generation" yaml:"generation"`

	// Testing configuration
	Testing TestingConfig `json:"testing" yaml:"testing"`

	// Templates and customization
	Templates TemplateConfig `json:"templates" yaml:"templates"`

	// Features and plugins
	Features FeatureConfig `json:"features" yaml:"features"`

	// Deployment and infrastructure
	Deploy DeployConfig `json:"deploy" yaml:"deploy"`
}

// ProjectConfig contains basic project information.
type ProjectConfig struct {
	Name        string            `json:"name"        yaml:"name"`
	Module      string            `json:"module"      yaml:"module"`
	Description string            `json:"description" yaml:"description"`
	Version     string            `json:"version"     yaml:"version"`
	Author      string            `json:"author"      yaml:"author"`
	License     string            `json:"license"     yaml:"license"`
	Repository  string            `json:"repository"  yaml:"repository"`
	Tags        []string          `json:"tags"        yaml:"tags"`
	Metadata    map[string]string `json:"metadata"    yaml:"metadata"`
}

// ArchitectureConfig defines Clean Architecture preferences.
type ArchitectureConfig struct {
	// Layers configuration
	Layers LayersConfig `json:"layers" yaml:"layers"`

	// Patterns to apply
	Patterns []string `json:"patterns" yaml:"patterns"`

	// Dependency injection type
	DI DIConfig `json:"di" yaml:"di"`

	// Naming conventions
	Naming NamingConfig `json:"naming" yaml:"naming"`
}

// LayersConfig defines which layers to generate and their structure.
type LayersConfig struct {
	Domain     LayerConfig   `json:"domain"     yaml:"domain"`
	UseCase    LayerConfig   `json:"usecase"    yaml:"usecase"`
	Repository LayerConfig   `json:"repository" yaml:"repository"`
	Handler    LayerConfig   `json:"handler"    yaml:"handler"`
	Custom     []LayerConfig `json:"custom"     yaml:"custom"`
}

// LayerConfig defines configuration for a specific layer.
type LayerConfig struct {
	Enabled     bool              `json:"enabled"     yaml:"enabled"`
	Directory   string            `json:"directory"   yaml:"directory"`
	Patterns    []string          `json:"patterns"    yaml:"patterns"`
	Templates   []string          `json:"templates"   yaml:"templates"`
	Validations []string          `json:"validations" yaml:"validations"`
	Extensions  map[string]string `json:"extensions"  yaml:"extensions"`
}

// DIConfig defines dependency injection preferences.
type DIConfig struct {
	Type       string            `json:"type"       yaml:"type"` // manual, wire, fx, dig
	AutoWire   bool              `json:"auto_wire"  yaml:"auto_wire"`
	Providers  []string          `json:"providers"  yaml:"providers"`
	Modules    []string          `json:"modules"    yaml:"modules"`
	Extensions map[string]string `json:"extensions" yaml:"extensions"`
}

// NamingConfig defines naming conventions.
type NamingConfig struct {
	Entities  string `json:"entities"  yaml:"entities"`  // PascalCase, camelCase, snake_case
	Fields    string `json:"fields"    yaml:"fields"`    // PascalCase, camelCase, snake_case
	Files     string `json:"files"     yaml:"files"`     // snake_case, kebab-case, camelCase
	Packages  string `json:"packages"  yaml:"packages"`  // lowercase, snake_case
	Constants string `json:"constants" yaml:"constants"` // UPPER_CASE, PascalCase
	Variables string `json:"variables" yaml:"variables"` // camelCase, snake_case
	Functions string `json:"functions" yaml:"functions"` // camelCase, PascalCase
}

// DatabaseConfig contains database configuration.
type DatabaseConfig struct {
	Type        string                `json:"type"         yaml:"type"` // postgres, mysql, mongodb, sqlite
	Host        string                `json:"host"         yaml:"host"`
	Port        int                   `json:"port"         yaml:"port"`
	Name        string                `json:"name"         yaml:"name"`
	Migrations  MigrationConfig       `json:"migrations"   yaml:"migrations"`
	Connection  ConnectionConfig      `json:"connection"   yaml:"connection"`
	Features    DatabaseFeatureConfig `json:"features"     yaml:"features"`
	Extensions  []string              `json:"extensions"   yaml:"extensions"`
	CustomTypes map[string]string     `json:"custom_types" yaml:"custom_types"`
}

// MigrationConfig defines migration preferences.
type MigrationConfig struct {
	Enabled      bool     `json:"enabled"       yaml:"enabled"`
	AutoGenerate bool     `json:"auto_generate" yaml:"auto_generate"`
	Directory    string   `json:"directory"     yaml:"directory"`
	Naming       string   `json:"naming"        yaml:"naming"`
	Versioning   string   `json:"versioning"    yaml:"versioning"`
	Tools        []string `json:"tools"         yaml:"tools"`
}

// ConnectionConfig defines database connection settings.
type ConnectionConfig struct {
	MaxOpen     int           `json:"max_open"     yaml:"max_open"`
	MaxIdle     int           `json:"max_idle"     yaml:"max_idle"`
	MaxLifetime time.Duration `json:"max_lifetime" yaml:"max_lifetime"`
	SSLMode     string        `json:"ssl_mode"     yaml:"ssl_mode"`
	Timezone    string        `json:"timezone"     yaml:"timezone"`
	Charset     string        `json:"charset"      yaml:"charset"`
	Collation   string        `json:"collation"    yaml:"collation"`
}

// DatabaseFeatureConfig defines database-specific features.
type DatabaseFeatureConfig struct {
	SoftDelete   bool     `json:"soft_delete"  yaml:"soft_delete"`
	Timestamps   bool     `json:"timestamps"   yaml:"timestamps"`
	UUID         bool     `json:"uuid"         yaml:"uuid"`
	Audit        bool     `json:"audit"        yaml:"audit"`
	Versioning   bool     `json:"versioning"   yaml:"versioning"`
	Partitioning bool     `json:"partitioning" yaml:"partitioning"`
	Indexes      []string `json:"indexes"      yaml:"indexes"`
	Constraints  []string `json:"constraints"  yaml:"constraints"`
}

// GenerationConfig defines code generation preferences.
type GenerationConfig struct {
	// Field validation preferences
	Validation ValidationConfig `json:"validation" yaml:"validation"`

	// Business rules generation
	BusinessRules BusinessRulesConfig `json:"business_rules" yaml:"business_rules"`

	// Documentation generation
	Documentation DocumentationConfig `json:"documentation" yaml:"documentation"`

	// Code style preferences
	Style StyleConfig `json:"style" yaml:"style"`

	// Import management
	Imports ImportConfig `json:"imports" yaml:"imports"`
}

// ValidationConfig defines validation generation preferences.
type ValidationConfig struct {
	Enabled   bool     `json:"enabled"   yaml:"enabled"`
	Library   string   `json:"library"   yaml:"library"` // builtin, validator, ozzo-validation
	Tags      []string `json:"tags"      yaml:"tags"`
	Custom    []string `json:"custom"    yaml:"custom"`
	Sanitize  bool     `json:"sanitize"  yaml:"sanitize"`
	Transform bool     `json:"transform" yaml:"transform"`
}

// BusinessRulesConfig defines business rules generation.
type BusinessRulesConfig struct {
	Enabled   bool     `json:"enabled"   yaml:"enabled"`
	Patterns  []string `json:"patterns"  yaml:"patterns"`
	Templates []string `json:"templates" yaml:"templates"`
	Events    bool     `json:"events"    yaml:"events"`
	Guards    bool     `json:"guards"    yaml:"guards"`
}

// DocumentationConfig defines documentation generation.
type DocumentationConfig struct {
	Swagger  SwaggerConfig  `json:"swagger"  yaml:"swagger"`
	Postman  PostmanConfig  `json:"postman"  yaml:"postman"`
	Markdown MarkdownConfig `json:"markdown" yaml:"markdown"`
	Comments CommentsConfig `json:"comments" yaml:"comments"`
}

// SwaggerConfig defines Swagger/OpenAPI generation.
type SwaggerConfig struct {
	Enabled     bool              `json:"enabled"     yaml:"enabled"`
	Version     string            `json:"version"     yaml:"version"`
	Output      string            `json:"output"      yaml:"output"`
	Title       string            `json:"title"       yaml:"title"`
	Description string            `json:"description" yaml:"description"`
	Host        string            `json:"host"        yaml:"host"`
	BasePath    string            `json:"base_path"   yaml:"base_path"`
	Schemes     []string          `json:"schemes"     yaml:"schemes"`
	Tags        []SwaggerTag      `json:"tags"        yaml:"tags"`
	Extensions  map[string]string `json:"extensions"  yaml:"extensions"`
}

// SwaggerTag defines Swagger tag configuration.
type SwaggerTag struct {
	Name        string `json:"name"        yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

// PostmanConfig defines Postman collection generation.
type PostmanConfig struct {
	Enabled     bool   `json:"enabled"     yaml:"enabled"`
	Output      string `json:"output"      yaml:"output"`
	Environment bool   `json:"environment" yaml:"environment"`
	Tests       bool   `json:"tests"       yaml:"tests"`
	Variables   bool   `json:"variables"   yaml:"variables"`
}

// MarkdownConfig defines Markdown documentation generation.
type MarkdownConfig struct {
	Enabled  bool   `json:"enabled"  yaml:"enabled"`
	Output   string `json:"output"   yaml:"output"`
	Template string `json:"template" yaml:"template"`
	TOC      bool   `json:"toc"      yaml:"toc"`
	Examples bool   `json:"examples" yaml:"examples"`
}

// CommentsConfig defines code comments generation.
type CommentsConfig struct {
	Enabled    bool   `json:"enabled"    yaml:"enabled"`
	Language   string `json:"language"   yaml:"language"` // spanish, english
	Style      string `json:"style"      yaml:"style"`    // godoc, standard
	Examples   bool   `json:"examples"   yaml:"examples"`
	TODO       bool   `json:"todo"       yaml:"todo"`
	Deprecated bool   `json:"deprecated" yaml:"deprecated"`
}

// StyleConfig defines code style preferences.
type StyleConfig struct {
	Gofmt       bool     `json:"gofmt"       yaml:"gofmt"`
	Goimports   bool     `json:"goimports"   yaml:"goimports"`
	Golint      bool     `json:"golint"      yaml:"golint"`
	Govet       bool     `json:"govet"       yaml:"govet"`
	Staticcheck bool     `json:"staticcheck" yaml:"staticcheck"`
	Custom      []string `json:"custom"      yaml:"custom"`
	LineLength  int      `json:"line_length" yaml:"line_length"`
	TabWidth    int      `json:"tab_width"   yaml:"tab_width"`
}

// ImportConfig defines import management.
type ImportConfig struct {
	GroupStandard   bool    `json:"group_standard"    yaml:"group_standard"`
	GroupThirdParty bool    `json:"group_third_party" yaml:"group_third_party"`
	GroupLocal      bool    `json:"group_local"       yaml:"group_local"`
	SortAlpha       bool    `json:"sort_alpha"        yaml:"sort_alpha"`
	RemoveUnused    bool    `json:"remove_unused"     yaml:"remove_unused"`
	Aliases         []Alias `json:"aliases"           yaml:"aliases"`
}

// Alias defines import alias configuration.
type Alias struct {
	Package string `json:"package" yaml:"package"`
	Alias   string `json:"alias"   yaml:"alias"`
}

// TestingConfig defines testing generation preferences.
type TestingConfig struct {
	Enabled     bool           `json:"enabled"     yaml:"enabled"`
	Framework   string         `json:"framework"   yaml:"framework"` // testify, ginkgo, builtin
	Coverage    CoverageConfig `json:"coverage"    yaml:"coverage"`
	Mocks       MockConfig     `json:"mocks"       yaml:"mocks"`
	Integration bool           `json:"integration" yaml:"integration"`
	Benchmarks  bool           `json:"benchmarks"  yaml:"benchmarks"`
	Examples    bool           `json:"examples"    yaml:"examples"`
	Fixtures    FixtureConfig  `json:"fixtures"    yaml:"fixtures"`
}

// CoverageConfig defines test coverage preferences.
type CoverageConfig struct {
	Enabled   bool     `json:"enabled"   yaml:"enabled"`
	Threshold float64  `json:"threshold" yaml:"threshold"`
	Output    string   `json:"output"    yaml:"output"`
	Format    string   `json:"format"    yaml:"format"`
	Exclude   []string `json:"exclude"   yaml:"exclude"`
}

// MockConfig defines mock generation preferences.
type MockConfig struct {
	Enabled    bool     `json:"enabled"    yaml:"enabled"`
	Tool       string   `json:"tool"       yaml:"tool"` // gomock, testify, counterfeiter
	Directory  string   `json:"directory"  yaml:"directory"`
	Suffix     string   `json:"suffix"     yaml:"suffix"`
	Interfaces []string `json:"interfaces" yaml:"interfaces"`
}

// FixtureConfig defines test fixture preferences.
type FixtureConfig struct {
	Enabled   bool     `json:"enabled"   yaml:"enabled"`
	Directory string   `json:"directory" yaml:"directory"`
	Format    string   `json:"format"    yaml:"format"` // json, yaml, sql
	Seeds     bool     `json:"seeds"     yaml:"seeds"`
	Factories []string `json:"factories" yaml:"factories"`
}

// TemplateConfig defines custom templates.
type TemplateConfig struct {
	Directory string              `json:"directory" yaml:"directory"`
	Custom    map[string]Template `json:"custom"    yaml:"custom"`
	Overrides map[string]Template `json:"overrides" yaml:"overrides"`
	Variables map[string]string   `json:"variables" yaml:"variables"`
}

// Template defines a custom template configuration.
type Template struct {
	Path       string            `json:"path"       yaml:"path"`
	Type       string            `json:"type"       yaml:"type"`
	Variables  map[string]string `json:"variables"  yaml:"variables"`
	Conditions []string          `json:"conditions" yaml:"conditions"`
}

// FeatureConfig defines feature flags and plugins.
type FeatureConfig struct {
	// Authentication and authorization
	Auth AuthConfig `json:"auth" yaml:"auth"`

	// Caching configuration
	Cache CacheConfig `json:"cache" yaml:"cache"`

	// Logging configuration
	Logging LoggingConfig `json:"logging" yaml:"logging"`

	// Monitoring configuration
	Monitoring MonitoringConfig `json:"monitoring" yaml:"monitoring"`

	// Security features
	Security SecurityConfig `json:"security" yaml:"security"`

	// Plugins and extensions
	Plugins []PluginConfig `json:"plugins" yaml:"plugins"`
}

// AuthConfig defines authentication configuration.
type AuthConfig struct {
	Enabled    bool     `json:"enabled"    yaml:"enabled"`
	Type       string   `json:"type"       yaml:"type"` // jwt, oauth2, session, basic
	Providers  []string `json:"providers"  yaml:"providers"`
	RBAC       bool     `json:"rbac"       yaml:"rbac"`
	Middleware bool     `json:"middleware" yaml:"middleware"`
}

// CacheConfig defines caching configuration.
type CacheConfig struct {
	Enabled  bool     `json:"enabled"  yaml:"enabled"`
	Type     string   `json:"type"     yaml:"type"` // redis, memcached, inmemory
	TTL      string   `json:"ttl"      yaml:"ttl"`
	Layers   []string `json:"layers"   yaml:"layers"`
	Patterns []string `json:"patterns" yaml:"patterns"`
}

// LoggingConfig defines logging configuration.
type LoggingConfig struct {
	Enabled    bool     `json:"enabled"    yaml:"enabled"`
	Level      string   `json:"level"      yaml:"level"`
	Format     string   `json:"format"     yaml:"format"` // json, text, structured
	Output     []string `json:"output"     yaml:"output"` // stdout, file, syslog
	Structured bool     `json:"structured" yaml:"structured"`
	Tracing    bool     `json:"tracing"    yaml:"tracing"`
}

// MonitoringConfig defines monitoring and observability.
type MonitoringConfig struct {
	Enabled     bool     `json:"enabled"      yaml:"enabled"`
	Metrics     bool     `json:"metrics"      yaml:"metrics"`
	Tracing     bool     `json:"tracing"      yaml:"tracing"`
	HealthCheck bool     `json:"health_check" yaml:"health_check"`
	Profiling   bool     `json:"profiling"    yaml:"profiling"`
	Tools       []string `json:"tools"        yaml:"tools"` // prometheus, jaeger, datadog
}

// SecurityConfig defines security features.
type SecurityConfig struct {
	HTTPS        bool     `json:"https"        yaml:"https"`
	CORS         bool     `json:"cors"         yaml:"cors"`
	RateLimit    bool     `json:"rate_limit"   yaml:"rate_limit"`
	Validation   bool     `json:"validation"   yaml:"validation"`
	Sanitization bool     `json:"sanitization" yaml:"sanitization"`
	Headers      []string `json:"headers"      yaml:"headers"`
	Middleware   []string `json:"middleware"   yaml:"middleware"`
}

// PluginConfig defines plugin configuration.
type PluginConfig struct {
	Name     string            `json:"name"     yaml:"name"`
	Version  string            `json:"version"  yaml:"version"`
	Enabled  bool              `json:"enabled"  yaml:"enabled"`
	Config   map[string]string `json:"config"   yaml:"config"`
	Priority int               `json:"priority" yaml:"priority"`
}

// DeployConfig defines deployment and infrastructure.
type DeployConfig struct {
	Docker      DockerConfig     `json:"docker"       yaml:"docker"`
	Kubernetes  KubernetesConfig `json:"kubernetes"   yaml:"kubernetes"`
	CI          CIConfig         `json:"ci"           yaml:"ci"`
	Environment []EnvConfig      `json:"environments" yaml:"environments"`
}

// DockerConfig defines Docker configuration.
type DockerConfig struct {
	Enabled    bool              `json:"enabled"    yaml:"enabled"`
	Dockerfile string            `json:"dockerfile" yaml:"dockerfile"`
	Image      string            `json:"image"      yaml:"image"`
	Registry   string            `json:"registry"   yaml:"registry"`
	Compose    bool              `json:"compose"    yaml:"compose"`
	Multistage bool              `json:"multistage" yaml:"multistage"`
	Labels     map[string]string `json:"labels"     yaml:"labels"`
}

// KubernetesConfig defines Kubernetes configuration.
type KubernetesConfig struct {
	Enabled    bool              `json:"enabled"     yaml:"enabled"`
	Namespace  string            `json:"namespace"   yaml:"namespace"`
	Manifests  string            `json:"manifests"   yaml:"manifests"`
	Helm       bool              `json:"helm"        yaml:"helm"`
	Ingress    bool              `json:"ingress"     yaml:"ingress"`
	ConfigMaps bool              `json:"config_maps" yaml:"config_maps"`
	Secrets    bool              `json:"secrets"     yaml:"secrets"`
	Labels     map[string]string `json:"labels"      yaml:"labels"`
}

// CIConfig defines CI/CD configuration.
type CIConfig struct {
	Enabled   bool     `json:"enabled"   yaml:"enabled"`
	Provider  string   `json:"provider"  yaml:"provider"` // github-actions, gitlab-ci, jenkins
	Workflows []string `json:"workflows" yaml:"workflows"`
	Tests     bool     `json:"tests"     yaml:"tests"`
	Build     bool     `json:"build"     yaml:"build"`
	Deploy    bool     `json:"deploy"    yaml:"deploy"`
}

// EnvConfig defines environment-specific configuration.
type EnvConfig struct {
	Name      string            `json:"name"      yaml:"name"`
	Default   bool              `json:"default"   yaml:"default"`
	Variables map[string]string `json:"variables" yaml:"variables"`
	Overrides GocaConfig        `json:"overrides" yaml:"overrides"`
}
