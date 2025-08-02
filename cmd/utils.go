package cmd

import (
	"bufio"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
)

// getModuleName reads the module name from go.mod file
func getModuleName() string {
	goMod, err := os.Open("go.mod")
	if err != nil {
		return "myproject" // fallback
	}
	defer goMod.Close()

	scanner := bufio.NewScanner(goMod)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}
	return "myproject" // fallback
}

// writeFile creates a file with the given content, creating directories if needed
func writeFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creando directorio %s: %w", dir, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creando archivo %s: %w", path, err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error escribiendo archivo %s: %w", path, err)
	}

	return nil
}

// writeGoFile creates a Go file with auto-formatting
func writeGoFile(path, content string) error {
	// Format Go code if it's a .go file
	if strings.HasSuffix(path, ".go") {
		formatted, err := format.Source([]byte(content))
		if err != nil {
			fmt.Printf("⚠️  Advertencia: No se pudo formatear el código Go para %s: %v\n", path, err)
			// Continuar con código sin formatear
		} else {
			content = string(formatted)
		}
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creando directorio %s: %w", dir, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creando archivo %s: %w", path, err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error escribiendo archivo %s: %w", path, err)
	}

	return nil
}

// getImportPath determines whether to use the full module path or relative path for imports
func getImportPath(moduleName string) string {
	// Check if we're in a test environment with a GitHub-style fake module
	if strings.Contains(moduleName, "github.com/goca/testproject") {
		// For test environments, use the module name as-is since it has valid structure
		return moduleName
	}

	// For all projects (both local and remote), use the module name as defined in go.mod
	// Go modules handle internal imports automatically based on the module declaration
	return moduleName
}

// generateFieldsForCreate genera campos reales para estructuras Create basados en los fields
func generateFieldsForCreate(content *strings.Builder, fields []Field, validation bool) {
	for _, field := range fields {
		if field.Name != "ID" { // Excluir ID en estructuras de creación
			validationTag := ""
			if validation {
				validationTag = getValidationTag(field.Type)
			}

			jsonTag := fmt.Sprintf("json:\"%s\"", strings.ToLower(field.Name))
			if validationTag != "" {
				jsonTag += fmt.Sprintf(" validate:\"%s\"", validationTag)
			}

			content.WriteString(fmt.Sprintf("\t%s %s `%s`\n", field.Name, field.Type, jsonTag))
		}
	}
}

// generateFieldsForUpdate genera campos reales para estructuras Update basados en los fields
func generateFieldsForUpdate(content *strings.Builder, fields []Field, validation bool) {
	for _, field := range fields {
		if field.Name != "ID" { // Excluir ID en estructuras de actualización
			validationTag := ""
			if validation {
				validationTag = getValidationTagForUpdate(field.Type)
			}

			jsonTag := fmt.Sprintf("json:\"%s,omitempty\"", strings.ToLower(field.Name))
			if validationTag != "" {
				jsonTag += fmt.Sprintf(" validate:\"%s\"", validationTag)
			}

			content.WriteString(fmt.Sprintf("\t%s %s `%s`\n", field.Name, field.Type, jsonTag))
		}
	}
}

// generateFieldMapping genera el mapeo real de campos de input a entity
func generateFieldMapping(content *strings.Builder, fields []Field) {
	for _, field := range fields {
		if field.Name != "ID" {
			content.WriteString(fmt.Sprintf("\t\t%s: input.%s,\n", field.Name, field.Name))
		}
	}
}

// generateUpdateMapping genera el mapeo real para operaciones de actualización
func generateUpdateMapping(content *strings.Builder, fields []Field, entityVar string) {
	for _, field := range fields {
		if field.Name != "ID" {
			content.WriteString(fmt.Sprintf("\tif input.%s != \"\" {\n", field.Name))
			content.WriteString(fmt.Sprintf("\t\t%s.%s = input.%s\n", entityVar, field.Name, field.Name))
			content.WriteString("\t}\n")
		}
	}
}

// getValidationTagForUpdate retorna validaciones para campos de actualización (más permisivas)
func getValidationTagForUpdate(fieldType string) string {
	switch fieldType {
	case "string":
		return "omitempty,min=1"
	case "int", "int64", "uint", "uint64":
		return "omitempty,min=1"
	case "float64", "float32":
		return "omitempty,min=0"
	case "bool":
		return ""
	case "time.Time":
		return "omitempty"
	default:
		return "omitempty"
	}
}

// generateSQLFields genera los campos SQL para queries basados en los fields
func generateSQLFields(fields []Field, operation string) (string, string) {
	var fieldNames []string
	var placeholders []string

	for _, field := range fields {
		if field.Name != "ID" || operation == "SELECT" {
			columnName := strings.ToLower(field.Name)
			if field.Name == "ID" {
				columnName = "id"
			}
			fieldNames = append(fieldNames, columnName)
			if operation == "INSERT" && field.Name != "ID" {
				placeholders = append(placeholders, "?")
			}
		}
	}

	return strings.Join(fieldNames, ", "), strings.Join(placeholders, ", ")
}

// generateSQLAssignments genera asignaciones SQL para operaciones UPDATE
func generateSQLAssignments(fields []Field) string {
	var assignments []string
	for _, field := range fields {
		if field.Name != "ID" {
			columnName := strings.ToLower(field.Name)
			assignments = append(assignments, fmt.Sprintf("%s = ?", columnName))
		}
	}
	return strings.Join(assignments, ", ")
}

// writeGoFileOrPanic wrapper para compatibilidad - maneja errores imprimiendo y continuando
func writeGoFileOrPanic(path, content string) {
	if err := writeGoFile(path, content); err != nil {
		fmt.Printf("❌ Error escribiendo archivo %s: %v\n", path, err)
		// No hacer panic, solo imprimir el error y continuar
	}
}

// writeFileOrPanic wrapper para compatibilidad - maneja errores imprimiendo y continuando
func writeFileOrPanic(path, content string) {
	if err := writeFile(path, content); err != nil {
		fmt.Printf("❌ Error escribiendo archivo %s: %v\n", path, err)
		// No hacer panic, solo imprimir el error y continuar
	}
}
