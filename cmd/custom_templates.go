package cmd

import (
	"fmt"
	"regexp"
	"strings"
)

// templateFieldData mirrors the shape expected by the built-in customizable
// templates (Name/Type/JSONName/Validations/GormTags), adapted from the
// internal Field{Name, Type, Tag} representation used by the real generators.
type templateFieldData struct {
	Name        string
	Type        string
	JSONName    string
	Validations []string
	GormTags    []string
}

var tagPartRegexp = regexp.MustCompile(`(\w+):"([^"]*)"`)

// fieldToTemplateData parses a Field's raw struct tag into the pieces a
// custom template can reference individually.
func fieldToTemplateData(f Field) templateFieldData {
	data := templateFieldData{
		Name:     f.Name,
		Type:     f.Type,
		JSONName: strings.ToLower(f.Name),
	}

	raw := strings.Trim(f.Tag, "`")
	for _, m := range tagPartRegexp.FindAllStringSubmatch(raw, -1) {
		key, val := m[1], m[2]
		switch key {
		case "json":
			if parts := strings.Split(val, ","); len(parts) > 0 && parts[0] != "" {
				data.JSONName = parts[0]
			}
		case "validate":
			data.Validations = strings.Split(val, ",")
		case "gorm":
			data.GormTags = strings.Split(val, ";")
		}
	}

	return data
}

// nonSystemTemplateFields converts fields to template data, skipping the
// system-managed ones (ID/CreatedAt/UpdatedAt/DeletedAt) that templates
// already render explicitly via dedicated struct fields/feature flags.
func nonSystemTemplateFields(fields []Field) []templateFieldData {
	var out []templateFieldData
	for _, f := range fields {
		if isSystemField(f.Name) {
			continue
		}
		out = append(out, fieldToTemplateData(f))
	}
	return out
}

// loadTemplateManagerForCWD returns the project's TemplateManager if the
// current directory has a loaded .goca.yaml configuration, or nil otherwise
// (no config file, or the config carries no templates). A nil result means
// "no custom templates available, use the built-in generator" and callers
// must treat it as such rather than an error.
func loadTemplateManagerForCWD() *TemplateManager {
	ci := NewConfigIntegration()
	if err := ci.LoadConfigForProject(); err != nil {
		return nil
	}
	return ci.GetTemplateManager()
}

// buildEntityTemplateData builds the data map for the "domain/entity" custom
// template out of the same inputs the built-in generator receives. System
// fields (ID/CreatedAt/UpdatedAt/DeletedAt) are excluded from Fields since the
// template renders them itself via Features.*.
func buildEntityTemplateData(entityName string, fields []Field, validation, businessRules, timestamps, softDelete bool) map[string]interface{} {
	return map[string]interface{}{
		"EntityName":        entityName,
		"EntityDescription": fmt.Sprintf("a %s entity", strings.ToLower(entityName)),
		"Fields":            nonSystemTemplateFields(fields),
		"Features": map[string]interface{}{
			// The real generator has no UUID-primary-key option today; this is
			// always false so the template's ID field matches what's actually
			// generated (uint, auto-increment).
			"UUID":       false,
			"Timestamps": timestamps,
			"SoftDelete": softDelete,
		},
		"ValidationEnabled": validation,
		"BusinessRules":     businessRules,
	}
}

// buildDTOTemplateData builds the data map for the "usecase/dto" custom
// template, mirroring the Create/Update validation tags the built-in
// generator computes (dtoValidationTag/dtoUpdateValidationTag) so a
// custom-templated DTO file stays consistent with the rest of the project.
func buildDTOTemplateData(entity, fields string, validation bool) map[string]interface{} {
	fieldsList := parseFields(fields)

	var createFields, updateFields []templateFieldData
	for _, f := range fieldsList {
		if isSystemField(f.Name) {
			continue
		}
		cf := fieldToTemplateData(f)
		uf := fieldToTemplateData(f)
		if validation {
			cf.Validations = []string{dtoValidationTag(f)}
			if tag := dtoUpdateValidationTag(f); tag != "" {
				uf.Validations = []string{tag}
			}
		}
		createFields = append(createFields, cf)
		updateFields = append(updateFields, uf)
	}

	return map[string]interface{}{
		"EntityName":   entity,
		"CreateFields": createFields,
		"UpdateFields": updateFields,
		"Module":       getImportPath(getModuleName()),
	}
}

// buildHandlerTemplateData builds the data map for the "handler/http/handler"
// custom template. The handler file is self-contained (no cross-entity merge
// concerns), so it can always be swapped for a custom template when present.
func buildHandlerTemplateData(entity string) map[string]interface{} {
	return map[string]interface{}{
		"EntityName": entity,
		"Module":     getImportPath(getModuleName()),
	}
}

// buildRepositoryTemplateData builds the data map for the "repository/repo"
// custom template.
func buildRepositoryTemplateData(entity string) map[string]interface{} {
	return map[string]interface{}{
		"EntityName": entity,
		"Module":     getImportPath(getModuleName()),
	}
}

// renderCustomTemplate renders the named custom template with data if the
// project has defined one, reporting ok=false when there is none (or it
// fails to execute, in which case a warning is emitted and generation should
// fall back to the built-in generator instead of aborting).
func renderCustomTemplate(name string, data map[string]interface{}) (string, bool) {
	tm := loadTemplateManagerForCWD()
	if tm == nil || !tm.HasTemplate(name) {
		return "", false
	}

	content, err := tm.ExecuteTemplate(name, data)
	if err != nil {
		ui.Warning(fmt.Sprintf("Custom template %q failed, falling back to built-in generator: %v", name, err))
		return "", false
	}

	return content, true
}
