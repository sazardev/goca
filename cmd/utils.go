package cmd

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
)

// readEntityFieldsString reconstructs the "name:type,..." field specification of
// an already-generated entity by parsing internal/domain/<entity>.go. System
// fields (ID, CreatedAt, ...) are skipped. It returns an empty string when the
// file cannot be read or parsed, so callers can fall back to their defaults.
func readEntityFieldsString(entity string) string {
	filename := filepath.Join("internal", "domain", strings.ToLower(entity)+".go")
	src, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return ""
	}

	var parts []string
	ast.Inspect(file, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok || ts.Name.Name != entity {
			return true
		}
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}
		for _, f := range st.Fields.List {
			for _, nm := range f.Names {
				if isSystemField(nm.Name) {
					continue
				}
				name := strings.ToLower(nm.Name[:1]) + nm.Name[1:]
				parts = append(parts, name+":"+types.ExprString(f.Type))
			}
		}
		return false
	})
	return strings.Join(parts, ",")
}

// getModuleName reads the module name from go.mod file.
func getModuleName() string {
	goMod, err := os.Open("go.mod")
	if err != nil {
		return "myproject" // fallback
	}
	defer func() {
		if err := goMod.Close(); err != nil {
			if ui != nil {
				ui.Error(fmt.Sprintf("Error closing go.mod file: %v", err))
			} else {
				fmt.Printf("Error closing go.mod file: %v\n", err)
			}
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

// writeFile creates a file with the given content, creating directories if needed.
// An optional SafetyManager can be passed to enable dry-run, force, and backup support.
func writeFile(path, content string, sm ...*SafetyManager) error {
	if len(sm) > 0 && sm[0] != nil {
		return sm[0].WriteFile(path, content)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			if ui != nil {
				ui.Error(fmt.Sprintf("Error closing file %s: %v", path, err))
			} else {
				fmt.Printf("Error closing file %s: %v\n", path, err)
			}
		}
	}()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing file %s: %w", path, err)
	}

	return nil
}

// writeGoFile creates a Go file with auto-formatting.
// An optional SafetyManager can be passed to enable dry-run, force, and backup support.
func writeGoFile(path, content string, sm ...*SafetyManager) error {
	// Format Go code if it's a .go file
	if strings.HasSuffix(path, ".go") {
		formatted, err := format.Source([]byte(content))
		if err != nil {
			if ui != nil {
				ui.Warning(fmt.Sprintf("Could not format Go code for %s: %v", path, err))
			} else {
				fmt.Printf("Warning: Could not format Go code for %s: %v\n", path, err)
			}
			// Continue with unformatted code
		} else {
			content = string(formatted)
		}
	}

	// Route through SafetyManager when available
	if len(sm) > 0 && sm[0] != nil {
		return sm[0].WriteFile(path, content)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			if ui != nil {
				ui.Error(fmt.Sprintf("Error closing file %s: %v", path, err))
			} else {
				fmt.Printf("Error closing file %s: %v\n", path, err)
			}
		}
	}()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing file %s: %w", path, err)
	}

	return nil
}

// writeGoFileMerged writes a Go file whose content was rebuilt from the existing
// file (shared, merge-aware files like errors.go, dto.go, messages.go and the
// layer interface files). It overwrites without requiring --force so that a
// second feature can extend these shared files. Falls back to a plain overwrite
// when no SafetyManager is provided.
func writeGoFileMerged(path, content string, sm ...*SafetyManager) error {
	if strings.HasSuffix(path, ".go") {
		if formatted, err := format.Source([]byte(content)); err == nil {
			content = string(formatted)
		} else if ui != nil {
			ui.Warning(fmt.Sprintf("Could not format Go code for %s: %v", path, err))
		}
	}

	if len(sm) > 0 && sm[0] != nil {
		return sm[0].WriteMergedFile(path, content)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}
	//#nosec G703 // path derived from validated entity/project names
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("error writing file %s: %w", path, err)
	}
	return nil
}

// getImportPath determines whether to use the full module path or relative path for imports.
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

// generateSearchMethods generates search methods based on entity fields.
func generateSearchMethods(fields []Field, entity string) []SearchMethod {
	var methods []SearchMethod

	for _, field := range fields {
		if field.Name == "ID" {
			continue // ID already has FindByID by default
		}

		// Generate methods for fields commonly used for searches
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

// isSearchableField determines if a field should have a search method.
func isSearchableField(fieldName, fieldType string) bool {
	// Types that are not suitable for searches
	if fieldType == "[]byte" || fieldType == "interface{}" {
		return false
	}

	fieldLower := strings.ToLower(fieldName)

	// Common fields typically used for searches
	searchableFields := []string{
		"email", "username", "name", "code",
		"sku", "slug", "phone", "document", "dni",
		"passport", "license", "title",
	}

	for _, searchable := range searchableFields {
		if strings.Contains(fieldLower, searchable) {
			return true
		}
	}

	// Only string, int and uint fields are good for searches
	return fieldType == "string" || fieldType == "int" || fieldType == "uint"
}

// isUniqueField determines if a field should likely be unique.
func isUniqueField(fieldName string) bool {
	fieldLower := strings.ToLower(fieldName)
	uniqueFields := []string{
		"email", "username", "code", "sku", "slug",
		"document", "dni", "passport", "license",
	}

	for _, unique := range uniqueFields {
		if strings.Contains(fieldLower, unique) {
			return true
		}
	}

	return false
}

// SearchMethod represents a dynamically generated search method.
type SearchMethod struct {
	MethodName string // FindByEmail, FindByUsername, etc.
	FieldName  string // Email, Username, etc.
	FieldType  string // string, int, etc.
	ReturnType string // (*domain.User, error)
	IsUnique   bool   // true if it should return a single result
}

// generateSearchMethodSignature generates the search method signature.
func (sm SearchMethod) generateSearchMethodSignature() string {
	paramName := strings.ToLower(sm.FieldName)
	return fmt.Sprintf("\t%s(%s %s) %s", sm.MethodName, paramName, sm.FieldType, sm.ReturnType)
}

// generateSearchMethodImplementation generates the search method implementation.
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
