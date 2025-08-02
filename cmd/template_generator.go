package cmd

import (
	"fmt"
	"strings"
	"text/template"
)

// TemplateData holds all data needed for template generation
type TemplateData struct {
	Entity      EntityData
	Fields      []FieldData
	Module      string
	Database    string
	Features    FeatureFlags
	Imports     []string
	Methods     []MethodData
	Validations []ValidationData
}

// EntityData holds entity-specific information
type EntityData struct {
	Name       string
	NameLower  string
	NamePlural string
	Package    string
}

// FieldData holds field-specific information
type FieldData struct {
	Name         string
	Type         string
	JSONTag      string
	GormTag      string
	ValidateTag  string
	IsRequired   bool
	IsUnique     bool
	IsSearchable bool
}

// FeatureFlags holds feature configuration
type FeatureFlags struct {
	Validation    bool
	BusinessRules bool
	Timestamps    bool
	SoftDelete    bool
	Cache         bool
	Transactions  bool
	Auth          bool
}

// MethodData holds method generation information
type MethodData struct {
	Name       string
	Params     []ParamData
	ReturnType string
	Body       string
}

// ParamData holds parameter information
type ParamData struct {
	Name string
	Type string
}

// ValidationData holds validation information
type ValidationData struct {
	Field    string
	Rule     string
	Message  string
	Priority int
}

// TemplateGenerator generates code from templates with dynamic data
type TemplateGenerator struct {
	fieldValidator *FieldValidator
}

// NewTemplateGenerator creates a new template generator
func NewTemplateGenerator() *TemplateGenerator {
	return &TemplateGenerator{
		fieldValidator: NewFieldValidator(),
	}
}

// GenerateFromTemplate generates code using template and data
func (g *TemplateGenerator) GenerateFromTemplate(templateName string, data *TemplateData) (string, error) {
	tmpl, err := g.getTemplate(templateName)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", err
	}

	return result.String(), nil
}

// PrepareTemplateData prepares template data from command parameters
func (g *TemplateGenerator) PrepareTemplateData(entityName, fields string, features FeatureFlags) (*TemplateData, error) {
	// Parse fields using the field validator
	fieldsList, err := g.fieldValidator.ParseFieldsWithValidation(fields)
	if err != nil {
		return nil, err
	}

	// Convert to FieldData
	var fieldData []FieldData
	for _, field := range fieldsList {
		fieldData = append(fieldData, FieldData{
			Name:         field.Name,
			Type:         field.Type,
			JSONTag:      fmt.Sprintf("json:\"%s\"", strings.ToLower(field.Name)),
			GormTag:      getGormTag(field.Name, field.Type),
			ValidateTag:  getValidationTag(field.Type),
			IsRequired:   isRequiredField(field.Name),
			IsUnique:     g.fieldValidator.isLikelyUniqueField(strings.ToLower(field.Name)),
			IsSearchable: isSearchableField(strings.ToLower(field.Name), field.Type),
		})
	}

	// Prepare entity data
	entityData := EntityData{
		Name:       entityName,
		NameLower:  strings.ToLower(entityName),
		NamePlural: makePlural(entityName),
		Package:    "domain",
	}

	// Generate imports based on features
	imports := g.generateImports(features, fieldData)

	// Generate methods based on fields and features
	methods := g.generateMethods(entityData, fieldData, features)

	// Generate validations
	validations := g.generateValidations(fieldData, features)

	return &TemplateData{
		Entity:      entityData,
		Fields:      fieldData,
		Module:      getModuleName(),
		Features:    features,
		Imports:     imports,
		Methods:     methods,
		Validations: validations,
	}, nil
}

// getTemplate returns the appropriate template for the given name
func (g *TemplateGenerator) getTemplate(templateName string) (*template.Template, error) {
	switch templateName {
	case "entity":
		return g.getEntityTemplate(), nil
	case "usecase":
		return g.getUseCaseTemplate(), nil
	case "repository":
		return g.getRepositoryTemplate(), nil
	case "handler":
		return g.getHandlerTemplate(), nil
	default:
		return nil, fmt.Errorf("template not found: %s", templateName)
	}
}

// Helper functions
func isRequiredField(fieldName string) bool {
	requiredFields := []string{"name", "email", "title"}
	for _, required := range requiredFields {
		if strings.EqualFold(fieldName, required) {
			return true
		}
	}
	return false
}

func makePlural(word string) string {
	if strings.HasSuffix(word, "y") {
		return strings.TrimSuffix(word, "y") + "ies"
	}
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "sh") {
		return word + "es"
	}
	return word + "s"
}

func (g *TemplateGenerator) generateImports(features FeatureFlags, fields []FieldData) []string {
	imports := []string{}

	if features.Validation {
		imports = append(imports, "errors")
	}

	if features.Timestamps || features.SoftDelete {
		imports = append(imports, "time")
	}

	// Check if any field requires specific imports
	for _, field := range fields {
		if field.Type == "time.Time" {
			imports = append(imports, "time")
		}
		if strings.Contains(field.Type, "[]byte") {
			imports = append(imports, "bytes")
		}
	}

	return removeDuplicates(imports)
}

func (g *TemplateGenerator) generateMethods(entity EntityData, fields []FieldData, features FeatureFlags) []MethodData {
	var methods []MethodData

	if features.Validation {
		methods = append(methods, MethodData{
			Name:       "Validate",
			Params:     []ParamData{},
			ReturnType: "error",
			Body:       g.generateValidationBody(fields),
		})
	}

	if features.BusinessRules {
		methods = append(methods, g.generateBusinessRuleMethods(entity, fields)...)
	}

	return methods
}

func (g *TemplateGenerator) generateValidations(fields []FieldData, features FeatureFlags) []ValidationData {
	var validations []ValidationData

	if !features.Validation {
		return validations
	}

	for _, field := range fields {
		if field.IsRequired {
			validations = append(validations, ValidationData{
				Field:    field.Name,
				Rule:     "required",
				Message:  fmt.Sprintf("%s es requerido", strings.ToLower(field.Name)),
				Priority: 1,
			})
		}
	}

	return validations
}

func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

// Template generation methods would be implemented here
func (g *TemplateGenerator) getEntityTemplate() *template.Template {
	return template.Must(template.New("entity").Parse(entityTemplate))
}

func (g *TemplateGenerator) getUseCaseTemplate() *template.Template {
	return template.Must(template.New("usecase").Parse(useCaseTemplate))
}

func (g *TemplateGenerator) getRepositoryTemplate() *template.Template {
	return template.Must(template.New("repository").Parse(repositoryTemplate))
}

func (g *TemplateGenerator) getHandlerTemplate() *template.Template {
	return template.Must(template.New("handler").Parse(handlerTemplate))
}

func (g *TemplateGenerator) generateValidationBody(fields []FieldData) string {
	// Implementation for validation body generation
	return "// Validation logic generated dynamically"
}

func (g *TemplateGenerator) generateBusinessRuleMethods(entity EntityData, fields []FieldData) []MethodData {
	// Implementation for business rule methods generation
	return []MethodData{}
}
