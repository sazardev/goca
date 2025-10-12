package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateManager manages custom templates for code generation
type TemplateManager struct {
	config    *TemplateConfig
	baseDir   string
	templates map[string]*template.Template
	variables map[string]string
	functions template.FuncMap
}

// NewTemplateManager creates a new template manager
func NewTemplateManager(config *TemplateConfig, projectPath string) *TemplateManager {
	baseDir := filepath.Join(projectPath, config.Directory)

	tm := &TemplateManager{
		config:    config,
		baseDir:   baseDir,
		templates: make(map[string]*template.Template),
		variables: make(map[string]string),
		functions: template.FuncMap{
			"title": func(s string) string {
				if len(s) == 0 {
					return s
				}
				return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
			},
			"lower":      strings.ToLower,
			"upper":      strings.ToUpper,
			"camel":      toCamelCase,
			"pascal":     toPascalCase,
			"snake":      toSnakeCase,
			"kebab":      toKebabCase,
			"plural":     toPlural,
			"singular":   toSingular,
			"join":       strings.Join,
			"split":      strings.Split,
			"contains":   strings.Contains,
			"hasPrefix":  strings.HasPrefix,
			"hasSuffix":  strings.HasSuffix,
			"trimSpace":  strings.TrimSpace,
			"replace":    strings.Replace,
			"replaceAll": strings.ReplaceAll,
			// Aliases for consistency with existing templates
			"toCamelCase":  toCamelCase,
			"toPascalCase": toPascalCase,
			"toSnakeCase":  toSnakeCase,
			"toKebabCase":  toKebabCase,
			"toPlural":     toPlural,
			"toSingular":   toSingular,
		},
	}

	// Copy configuration variables
	for k, v := range config.Variables {
		tm.variables[k] = v
	}

	return tm
}

// LoadTemplates loads all templates from the templates directory
func (tm *TemplateManager) LoadTemplates() error {
	if _, err := os.Stat(tm.baseDir); os.IsNotExist(err) {
		// Templates directory doesn't exist, create it with defaults
		return tm.createDefaultTemplates()
	}

	// Walk through templates directory
	return filepath.Walk(tm.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Only process .tmpl and .tpl files
		if !strings.HasSuffix(path, ".tmpl") && !strings.HasSuffix(path, ".tpl") {
			return nil
		}

		return tm.loadTemplate(path)
	})
}

// loadTemplate loads a single template file
func (tm *TemplateManager) loadTemplate(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", filePath, err)
	}

	// Get template name from relative path
	relPath, err := filepath.Rel(tm.baseDir, filePath)
	if err != nil {
		return fmt.Errorf("failed to get relative path for %s: %w", filePath, err)
	}

	// Remove extension for template name
	templateName := strings.TrimSuffix(relPath, filepath.Ext(relPath))
	templateName = strings.ReplaceAll(templateName, "\\", "/") // Use forward slashes

	// Create template with custom functions
	tmpl := template.New(templateName).Funcs(tm.functions)

	// Parse template content
	_, err = tmpl.Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	tm.templates[templateName] = tmpl
	return nil
}

// createDefaultTemplates creates default template structure
func (tm *TemplateManager) createDefaultTemplates() error {
	// Create templates directory
	if err := os.MkdirAll(tm.baseDir, 0755); err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}

	// Create default template directories
	dirs := []string{
		"domain",
		"usecase",
		"repository",
		"handler/http",
		"messages",
		"docs",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(tm.baseDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create template directory %s: %w", dir, err)
		}
	}

	// Create default templates
	return tm.createBuiltinTemplates()
}

// createBuiltinTemplates creates built-in default templates
func (tm *TemplateManager) createBuiltinTemplates() error {
	templates := map[string]string{
		"domain/entity.tmpl": `package domain

import (
	"time"
{{- if .ValidationEnabled }}
	"github.com/go-playground/validator/v10"
{{- end }}
)

// {{.EntityName}} represents {{.EntityDescription}}
type {{.EntityName}} struct {
{{- if .Features.UUID }}
	ID   string ` + "`json:\"id\" gorm:\"type:uuid;primaryKey\"`" + `
{{- else }}
	ID   uint ` + "`json:\"id\" gorm:\"primaryKey\"`" + `
{{- end }}
{{- range .Fields }}
	{{.Name}} {{.Type}} ` + "`json:\"{{.JSONName}}\"{{if .Validations}} validate:\"{{join .Validations \",\"}}\"{{end}}{{if .GormTags}} gorm:\"{{join .GormTags \";\"}}\"{{end}}`" + `
{{- end }}
{{- if .Features.Timestamps }}
	CreatedAt time.Time ` + "`json:\"created_at\" gorm:\"autoCreateTime\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\" gorm:\"autoUpdateTime\"`" + `
{{- end }}
{{- if .Features.SoftDelete }}
	DeletedAt *time.Time ` + "`json:\"deleted_at,omitempty\" gorm:\"index\"`" + `
{{- end }}
}

// TableName returns the table name for {{.EntityName}}
func ({{lower (slice .EntityName 0 1)}}) TableName() string {
	return "{{snake .EntityName}}"
}

{{- if .ValidationEnabled }}

// Validate validates {{.EntityName}} fields
func ({{lower (slice .EntityName 0 1)}} *{{.EntityName}}) Validate() error {
	validate := validator.New()
	return validate.Struct({{lower (slice .EntityName 0 1)}})
}
{{- end }}

{{- if .BusinessRules }}

// Business Rules for {{.EntityName}}

// IsValid checks if the {{.EntityName}} is in a valid state
func ({{lower (slice .EntityName 0 1)}} *{{.EntityName}}) IsValid() bool {
	// Add business logic here
	return true
}
{{- end }}`,

		"usecase/dto.tmpl": `package usecase

import (
{{- if .Features.Timestamps }}
	"time"
{{- end }}
)

// Create{{.EntityName}}Request represents request to create {{.EntityName}}
type Create{{.EntityName}}Request struct {
{{- range .Fields }}
	{{.Name}} {{.Type}} ` + "`json:\"{{.JSONName}}\"{{if .Validations}} validate:\"{{join .Validations \",\"}}\"{{end}}`" + `
{{- end }}
}

// Update{{.EntityName}}Request represents request to update {{.EntityName}}
type Update{{.EntityName}}Request struct {
{{- range .Fields }}
	{{.Name}} *{{.Type}} ` + "`json:\"{{.JSONName}},omitempty\"`" + `
{{- end }}
}

// {{.EntityName}}Response represents {{.EntityName}} response
type {{.EntityName}}Response struct {
{{- if .Features.UUID }}
	ID   string ` + "`json:\"id\"`" + `
{{- else }}
	ID   uint ` + "`json:\"id\"`" + `
{{- end }}
{{- range .Fields }}
	{{.Name}} {{.Type}} ` + "`json:\"{{.JSONName}}\"`" + `
{{- end }}
{{- if .Features.Timestamps }}
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
{{- end }}
}

// List{{.EntityName}}Response represents paginated list response
type List{{.EntityName}}Response struct {
	Data       []{{.EntityName}}Response ` + "`json:\"data\"`" + `
	Total      int64                      ` + "`json:\"total\"`" + `
	Page       int                        ` + "`json:\"page\"`" + `
	PerPage    int                        ` + "`json:\"per_page\"`" + `
	TotalPages int                        ` + "`json:\"total_pages\"`" + `
}`,

		"handler/http/handler.tmpl": `package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"{{.Module}}/internal/usecase"
	"{{.Module}}/internal/messages"
	
	"github.com/gorilla/mux"
)

// {{.EntityName}}Handler handles HTTP requests for {{.EntityName}}
type {{.EntityName}}Handler struct {
	usecase usecase.{{.EntityName}}UseCase
}

// New{{.EntityName}}Handler creates a new {{.EntityName}} handler
func New{{.EntityName}}Handler(uc usecase.{{.EntityName}}UseCase) *{{.EntityName}}Handler {
	return &{{.EntityName}}Handler{
		usecase: uc,
	}
}

// Create handles POST /{{kebab (plural .EntityName)}}
func (h *{{.EntityName}}Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req usecase.Create{{.EntityName}}Request
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, messages.ErrInvalidJSON, http.StatusBadRequest)
		return
	}
	
	result, err := h.usecase.Create(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// GetByID handles GET /{{kebab (plural .EntityName)}}/:id
func (h *{{.EntityName}}Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, messages.ErrInvalidID, http.StatusBadRequest)
		return
	}
	
	result, err := h.usecase.GetByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Update handles PUT /{{kebab (plural .EntityName)}}/:id
func (h *{{.EntityName}}Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, messages.ErrInvalidID, http.StatusBadRequest)
		return
	}
	
	var req usecase.Update{{.EntityName}}Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, messages.ErrInvalidJSON, http.StatusBadRequest)
		return
	}
	
	result, err := h.usecase.Update(r.Context(), uint(id), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Delete handles DELETE /{{kebab (plural .EntityName)}}/:id
func (h *{{.EntityName}}Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, messages.ErrInvalidID, http.StatusBadRequest)
		return
	}
	
	if err := h.usecase.Delete(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// List handles GET /{{kebab (plural .EntityName)}}
func (h *{{.EntityName}}Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}
	
	result, err := h.usecase.List(r.Context(), page, perPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}`,

		"docs/README.tmpl": `# {{.ProjectName}}

{{.Description}}

## 📋 Configuración GOCA

Este proyecto fue generado usando GOCA CLI con configuración YAML personalizada.

### Archivo de configuración: .goca.yaml

This project uses centralized configuration in .goca.yaml for:

- **Architecture**: Layers, patterns, DI, naming conventions
- **Database**: Type ({{.DatabaseType}}), migrations, features
- **Generation**: Validation, business rules, documentation
- **Testing**: Framework, coverage, mocks
- **Templates**: Customizable in {{.TemplateDirectory}}

### Comandos disponibles

` + "```" + `bash
# Generar nuevas features usando configuración
goca feature Product --fields "name:string,price:float64"

# Los valores CLI sobrescriben la configuración
goca feature Order --fields "total:float64" --database mysql

# Generar documentación
goca docs generate

# Integrar features existentes
goca integrate --all
` + "```" + `

### Personalización de templates

Los templates se pueden personalizar en {{.TemplateDirectory}}:

` + "```" + `
{{.TemplateDirectory}}/
├── domain/
│   ├── entity.tmpl      # Template for entities
│   └── validations.tmpl # Template para validaciones
├── usecase/
│   ├── dto.tmpl         # Template para DTOs
│   └── service.tmpl     # Template para servicios
├── repository/
│   └── repo.tmpl        # Template para repositorios
├── handler/
│   └── http/
│       └── handler.tmpl # Template para handlers HTTP
└── docs/
    └── README.tmpl      # Este template
` + "```" + `

## Funciones de template disponibles

| Función | Descripción | Ejemplo |
|---------|-------------|---------|
| title | Primera letra mayúscula | {{title "hello"}} → "Hello" |
| pascal | PascalCase | {{pascal "user_name"}} → "UserName" |
| camel | camelCase | {{camel "user_name"}} → "userName" |
| snake | snake_case | {{snake "UserName"}} → "user_name" |
| kebab | kebab-case | {{kebab "UserName"}} → "user-name" |
| plural | Pluralización | {{plural "user"}} → "users" |
| singular | Singularización | {{singular "users"}} → "user" |

---

Generado por **GOCA CLI** v{{.Version}} 🚀`,
	}

	// Write default templates
	for name, content := range templates {
		filePath := filepath.Join(tm.baseDir, name)

		// Create directory if it doesn't exist
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create template directory: %w", err)
		}

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write template %s: %w", name, err)
		}
	}

	return nil
}

// HasTemplate checks if a custom template exists
func (tm *TemplateManager) HasTemplate(name string) bool {
	_, exists := tm.templates[name]
	return exists
}

// ExecuteTemplate executes a template with given data
func (tm *TemplateManager) ExecuteTemplate(name string, data interface{}) (string, error) {
	tmpl, exists := tm.templates[name]
	if !exists {
		return "", fmt.Errorf("template %s not found", name)
	}

	// Merge template variables with data
	templateData := tm.enrichData(data)

	var result strings.Builder
	if err := tmpl.Execute(&result, templateData); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return result.String(), nil
}

// enrichData enriches template data with variables and functions
func (tm *TemplateManager) enrichData(data interface{}) interface{} {
	// If data is a map, merge with variables
	if dataMap, ok := data.(map[string]interface{}); ok {
		// Add template variables
		for k, v := range tm.variables {
			if _, exists := dataMap[k]; !exists {
				dataMap[k] = v
			}
		}

		// Add template configuration
		dataMap["TemplateDirectory"] = tm.config.Directory
		dataMap["TemplateVariables"] = tm.variables

		return dataMap
	}

	return data
}

// GetAvailableTemplates returns list of available templates
func (tm *TemplateManager) GetAvailableTemplates() []string {
	templates := make([]string, 0, len(tm.templates))
	for name := range tm.templates {
		templates = append(templates, name)
	}
	return templates
}

// Helper functions for template functions

func toCamelCase(s string) string {
	if s == "" {
		return ""
	}

	words := strings.FieldsFunc(s, func(c rune) bool {
		return c == '_' || c == '-' || c == ' '
	})

	if len(words) == 0 {
		return s
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		result += strings.Title(strings.ToLower(words[i]))
	}

	return result
}

func toPascalCase(s string) string {
	if s == "" {
		return ""
	}

	words := strings.FieldsFunc(s, func(c rune) bool {
		return c == '_' || c == '-' || c == ' '
	})

	var result string
	for _, word := range words {
		result += strings.Title(strings.ToLower(word))
	}

	return result
}

func toSnakeCase(s string) string {
	if s == "" {
		return ""
	}

	var result []rune
	for i, r := range s {
		if i > 0 && (r >= 'A' && r <= 'Z') {
			result = append(result, '_')
		}
		result = append(result, rune(strings.ToLower(string(r))[0]))
	}

	return string(result)
}

func toKebabCase(s string) string {
	return strings.ReplaceAll(toSnakeCase(s), "_", "-")
}

func toPlural(s string) string {
	if s == "" {
		return ""
	}

	// Simple pluralization rules
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") {
		return s + "es"
	}

	return s + "s"
}

func toSingular(s string) string {
	if s == "" {
		return ""
	}

	// Simple singularization rules
	if strings.HasSuffix(s, "ies") {
		return s[:len(s)-3] + "y"
	}
	if strings.HasSuffix(s, "es") {
		return s[:len(s)-2]
	}
	if strings.HasSuffix(s, "s") {
		return s[:len(s)-1]
	}

	return s
}

// ExecuteTemplateString executes a template from string content (useful for testing)
func (tm *TemplateManager) ExecuteTemplateString(templateContent string, data interface{}) (string, error) {
	tmpl, err := template.New("temp").Funcs(tm.functions).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	// Enrich data with additional context
	enrichedData := tm.enrichData(data)

	var buf strings.Builder
	if err := tmpl.Execute(&buf, enrichedData); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}
