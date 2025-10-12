package cmd

import (
	"fmt"
	"text/template"
)

// TemplateComponent represents a reusable template component
type TemplateComponent struct {
	Name     string
	Template string
	Required bool
}

// EntityTemplateComponents contains all entity template components
var EntityTemplateComponents = map[string]TemplateComponent{
	"header": {
		Name: "header",
		Template: `package domain

{{range .Imports}}import "{{.}}"
{{end}}

`,
		Required: true,
	},
	"struct": {
		Name: "struct",
		Template: `type {{.Entity.Name}} struct {
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`" + `{{.JSONTag}} {{.GormTag}}` + "`" + `
{{end}}`,
		Required: true,
	},
	"timestamps": {
		Name: "timestamps",
		Template: `{{if .Features.Timestamps}}	CreatedAt time.Time ` + "`" + `json:"created_at" gorm:"autoCreateTime"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at" gorm:"autoUpdateTime"` + "`" + `
{{end}}`,
		Required: false,
	},
	"softDelete": {
		Name: "softDelete",
		Template: `{{if .Features.SoftDelete}}	DeletedAt gorm.DeletedAt ` + "`" + `json:"deleted_at,omitempty" gorm:"index"` + "`" + `
{{end}}`,
		Required: false,
	},
	"structClose": {
		Name: "structClose",
		Template: `}

`,
		Required: true,
	},
	"validation": {
		Name: "validation",
		Template: `{{if .Features.Validation}}
func ({{.Entity.NameLower}} *{{.Entity.Name}}) Validate() error {
{{range .Validations}}	if {{.Field}} == "" {
		return errors.New("{{.Message}}")
	}
{{end}}	return nil
}
{{end}}

`,
		Required: false,
	},
	"methods": {
		Name: "methods",
		Template: `{{range .Methods}}
func ({{$.Entity.NameLower}} *{{$.Entity.Name}}) {{.Name}}({{range .Params}}{{.Name}} {{.Type}}, {{end}}) {{.ReturnType}} {
	{{.Body}}
}
{{end}}`,
		Required: false,
	},
}

// UseCaseTemplateComponents contains all use case template components
var UseCaseTemplateComponents = map[string]TemplateComponent{
	"header": {
		Name: "header",
		Template: `package usecase

import (
	"{{.Module}}/internal/domain"
	"{{.Module}}/internal/repository"
)

`,
		Required: true,
	},
	"interface": {
		Name: "interface",
		Template: `type {{.Entity.Name}}UseCase interface {
	Create{{.Entity.Name}}(input Create{{.Entity.Name}}Input) (*Create{{.Entity.Name}}Output, error)
	Get{{.Entity.Name}}ByID(id int) (*{{.Entity.Name}}Output, error)
	Update{{.Entity.Name}}(id int, input Update{{.Entity.Name}}Input) error
	Delete{{.Entity.Name}}(id int) error
	List{{.Entity.Name}}s() (*List{{.Entity.Name}}sOutput, error)
}

`,
		Required: true,
	},
	"service": {
		Name: "service",
		Template: `type {{.Entity.NameLower}}Service struct {
	repo repository.{{.Entity.Name}}Repository
}

func New{{.Entity.Name}}Service(repo repository.{{.Entity.Name}}Repository) {{.Entity.Name}}UseCase {
	return &{{.Entity.NameLower}}Service{repo: repo}
}

`,
		Required: true,
	},
	"dtos": {
		Name: "dtos",
		Template: `type Create{{.Entity.Name}}Input struct {
{{range .Fields}}{{if ne .Name "ID"}}	{{.Name}} {{.Type}} ` + "`" + `json:"{{.JSONTag}}"{{if .ValidateTag}} validate:"{{.ValidateTag}}"{{end}}` + "`" + `
{{end}}{{end}}}

type Create{{.Entity.Name}}Output struct {
	{{.Entity.Name}} *domain.{{.Entity.Name}} ` + "`" + `json:"{{.Entity.NameLower}}"` + "`" + `
	Message string ` + "`" + `json:"message"` + "`" + `
}

type Update{{.Entity.Name}}Input struct {
{{range .Fields}}{{if ne .Name "ID"}}	{{.Name}} *{{.Type}} ` + "`" + `json:"{{.JSONTag}},omitempty"` + "`" + `
{{end}}{{end}}}

type {{.Entity.Name}}Output struct {
	{{.Entity.Name}} *domain.{{.Entity.Name}} ` + "`" + `json:"{{.Entity.NameLower}}"` + "`" + `
}

type List{{.Entity.Name}}sOutput struct {
	{{.Entity.NamePlural}} []domain.{{.Entity.Name}} ` + "`" + `json:"{{.Entity.NameLower}}s"` + "`" + `
	Total int ` + "`" + `json:"total"` + "`" + `
}
`,
		Required: true,
	},
}

// TemplateBuilder builds templates from components
type TemplateBuilder struct {
	components []TemplateComponent
}

// NewTemplateBuilder creates a new template builder
func NewTemplateBuilder() *TemplateBuilder {
	return &TemplateBuilder{
		components: make([]TemplateComponent, 0),
	}
}

// AddComponent adds a component to the template builder
func (tb *TemplateBuilder) AddComponent(component TemplateComponent) *TemplateBuilder {
	tb.components = append(tb.components, component)
	return tb
}

// AddComponentByName adds a component by name from the component map
func (tb *TemplateBuilder) AddComponentByName(name string, componentMap map[string]TemplateComponent) *TemplateBuilder {
	if component, exists := componentMap[name]; exists {
		tb.components = append(tb.components, component)
	}
	return tb
}

// Build builds the final template string
func (tb *TemplateBuilder) Build() string {
	var result string
	for _, component := range tb.components {
		result += component.Template
	}
	return result
}

// BuildTemplate creates a complete template from component names
func BuildTemplate(componentNames []string, componentMap map[string]TemplateComponent) string {
	builder := NewTemplateBuilder()
	for _, name := range componentNames {
		builder.AddComponentByName(name, componentMap)
	}
	return builder.Build()
}

// GetEntityTemplate builds entity template with specific components
func GetEntityTemplate(withTimestamps, withSoftDelete, withValidation, withMethods bool) string {
	components := []string{"header", "struct"}

	if withTimestamps {
		components = append(components, "timestamps")
	}

	if withSoftDelete {
		components = append(components, "softDelete")
	}

	components = append(components, "structClose")

	if withValidation {
		components = append(components, "validation")
	}

	if withMethods {
		components = append(components, "methods")
	}

	return BuildTemplate(components, EntityTemplateComponents)
}

// GetUseCaseTemplate builds use case template with all components
func GetUseCaseTemplate() string {
	components := []string{"header", "interface", "service", "dtos"}
	return BuildTemplate(components, UseCaseTemplateComponents)
}

// ValidateTemplate validates a template string
func ValidateTemplate(templateStr string) error {
	_, err := template.New("test").Parse(templateStr)
	if err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}
	return nil
}
