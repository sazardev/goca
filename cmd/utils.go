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
	defer func() {
		if err := goMod.Close(); err != nil {
			fmt.Printf("Error closing go.mod file: %v\n", err)
		}
	}()

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
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file %s: %v\n", path, err)
		}
	}()

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
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creando archivo %s: %w", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file %s: %v\n", path, err)
		}
	}()

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

// generateSearchMethods genera métodos de búsqueda basados en los campos de la entidad
func generateSearchMethods(fields []Field, entity string) []SearchMethod {
	var methods []SearchMethod

	for _, field := range fields {
		if field.Name == "ID" {
			continue // ID ya tiene FindByID por defecto
		}

		// Generar métodos para campos que comúnmente se usan para búsquedas
		if isSearchableField(field.Name, field.Type) {
			method := SearchMethod{
				MethodName: fmt.Sprintf("FindBy%s", field.Name),
				FieldName:  field.Name,
				FieldType:  field.Type,
				ReturnType: fmt.Sprintf("(*domain.%s, error)", entity),
				IsUnique:   isUniqueField(field.Name),
			}
			methods = append(methods, method)
		}
	}

	return methods
}

// isSearchableField determina si un campo debería tener un método de búsqueda
func isSearchableField(fieldName, fieldType string) bool {
	// Tipos que no son apropiados para búsquedas
	if fieldType == "[]byte" || fieldType == "interface{}" {
		return false
	}

	fieldLower := strings.ToLower(fieldName)

	// Campos comunes que suelen usarse para búsquedas
	searchableFields := []string{
		"email", "username", "nombre", "name", "codigo", "code",
		"sku", "slug", "telefono", "phone", "documento", "dni",
		"cedula", "passport", "license", "titulo", "title",
	}

	for _, searchable := range searchableFields {
		if strings.Contains(fieldLower, searchable) {
			return true
		}
	}

	// Solo campos string, int y uint son buenos para búsquedas
	return fieldType == "string" || fieldType == "int" || fieldType == "uint"
}

// isUniqueField determina si un campo probablemente debería ser único
func isUniqueField(fieldName string) bool {
	fieldLower := strings.ToLower(fieldName)
	uniqueFields := []string{
		"email", "username", "codigo", "code", "sku", "slug",
		"documento", "dni", "cedula", "passport", "license",
	}

	for _, unique := range uniqueFields {
		if strings.Contains(fieldLower, unique) {
			return true
		}
	}

	return false
}

// SearchMethod representa un método de búsqueda generado dinámicamente
type SearchMethod struct {
	MethodName string // FindByEmail, FindByUsername, etc.
	FieldName  string // Email, Username, etc.
	FieldType  string // string, int, etc.
	ReturnType string // (*domain.User, error)
	IsUnique   bool   // true si debería retornar un solo resultado
}

// generateSearchMethodSignature genera la firma del método de búsqueda
func (sm SearchMethod) generateSearchMethodSignature() string {
	paramName := strings.ToLower(sm.FieldName)
	return fmt.Sprintf("\t%s(%s %s) %s", sm.MethodName, paramName, sm.FieldType, sm.ReturnType)
}

// generateSearchMethodImplementation genera la implementación del método de búsqueda
func (sm SearchMethod) generateSearchMethodImplementation(receiverName, receiverType, entity string) string {
	paramName := strings.ToLower(sm.FieldName)
	entityVar := strings.ToLower(entity)

	var implementation strings.Builder
	implementation.WriteString(fmt.Sprintf("func (%s *%s) %s(%s %s) %s {\n",
		receiverName, receiverType, sm.MethodName, paramName, sm.FieldType, sm.ReturnType))

	implementation.WriteString(fmt.Sprintf("\t%s := &domain.%s{}\n", entityVar, entity))
	implementation.WriteString(fmt.Sprintf("\tresult := %s.db.Where(\"%s = ?\", %s).First(%s)\n",
		receiverName, strings.ToLower(sm.FieldName), paramName, entityVar))
	implementation.WriteString("\tif result.Error != nil {\n")
	implementation.WriteString("\t\treturn nil, result.Error\n")
	implementation.WriteString("\t}\n")
	implementation.WriteString(fmt.Sprintf("\treturn %s, nil\n", entityVar))
	implementation.WriteString("}\n\n")

	return implementation.String()
}
