package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// addEntityToAutoMigration adds a domain entity to the auto-migration list in main.go
func addEntityToAutoMigration(entity string) error {
	mainPath, err := findMainGoFile()
	if err != nil {
		return err
	}

	content, err := os.ReadFile(mainPath)
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}

	contentStr := string(content)
	entityReference := fmt.Sprintf("&domain.%s{}", entity)

	// Check if entity is already added to auto-migration (not in a comment)
	// Look for the entity reference followed by a comma and not preceded by //
	if isEntityInMigrationList(contentStr, entityReference) {
		return nil
	}

	// 1. Add domain import if not present
	updatedContent, err := ensureDomainImport(contentStr)
	if err != nil {
		return fmt.Errorf("failed to add domain import: %w", err)
	}

	// 2. Add entity to migration slice
	updatedContent, err = addEntityToMigrationSlice(updatedContent, entityReference)
	if err != nil {
		return err
	}

	if err := os.WriteFile(mainPath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to update main.go: %w", err)
	}

	return nil
}

// findMainGoFile locates the main.go file in possible locations
func findMainGoFile() (string, error) {
	possiblePaths := []string{
		filepath.Join("cmd", "server", "main.go"),
		"main.go",
		filepath.Join("cmd", "main.go"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("main.go not found in any expected location")
}

// addEntityToMigrationSlice adds an entity reference to the migration entities slice
func addEntityToMigrationSlice(content, entityReference string) (string, error) {
	commentPattern := "// Add domain entities here as they are created"

	if strings.Contains(content, commentPattern) {
		replacement := fmt.Sprintf("%s\n\t\t%s,", commentPattern, entityReference)
		return strings.Replace(content, commentPattern, replacement, 1), nil
	}

	return addEntityToEntitiesSlice(content, entityReference)
}

// addEntityToEntitiesSlice adds an entity to the entities slice as fallback
func addEntityToEntitiesSlice(content, entityReference string) (string, error) {
	entitiesPattern := "entities := []interface{}{"

	if !strings.Contains(content, entitiesPattern) {
		return "", fmt.Errorf("could not find auto-migration entities slice in main.go")
	}

	startIdx := strings.Index(content, entitiesPattern)
	if startIdx == -1 {
		return "", fmt.Errorf("could not find entities slice start")
	}

	closingIdx := findSliceClosingBrace(content, startIdx+len(entitiesPattern))
	if closingIdx == -1 {
		return "", fmt.Errorf("could not find entities slice closing brace")
	}

	beforeClosing := content[:closingIdx]
	afterClosing := content[closingIdx:]
	return beforeClosing + fmt.Sprintf("\n\t\t%s,", entityReference) + afterClosing, nil
}

// findSliceClosingBrace finds the closing brace of a slice starting from the given position
func findSliceClosingBrace(content string, startPos int) int {
	braceCount := 1
	for i := startPos; i < len(content) && braceCount > 0; i++ {
		switch content[i] {
		case '{':
			braceCount++
		case '}':
			braceCount--
		}
		if braceCount == 0 {
			return i
		}
	}
	return -1
}

// isEntityInMigrationList checks if an entity reference exists in the migration list (not in comments)
func isEntityInMigrationList(content, entityReference string) bool {
	// Find the entities slice
	entitiesPattern := "entities := []interface{}{"
	startIdx := strings.Index(content, entitiesPattern)
	if startIdx == -1 {
		return false
	}

	// Find the closing brace
	closingIdx := findSliceClosingBrace(content, startIdx+len(entitiesPattern))
	if closingIdx == -1 {
		return false
	}

	// Get the slice content
	sliceContent := content[startIdx+len(entitiesPattern) : closingIdx]

	// Split by lines and check each line
	lines := strings.Split(sliceContent, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}
		// Check if this line contains the entity reference
		if strings.Contains(line, entityReference) {
			return true
		}
	}

	return false
}

// ensureDomainImport ensures that the internal/domain package is imported in main.go
func ensureDomainImport(content string) (string, error) {
	// Check if domain is already imported
	if strings.Contains(content, "/internal/domain\"") {
		return content, nil
	}

	// Get module name from go.mod
	moduleName := getModuleName()
	if moduleName == "" {
		return "", fmt.Errorf("could not determine module name from go.mod")
	}

	domainImport := fmt.Sprintf("\"%s/internal/domain\"", moduleName)

	// Find import block
	importStart := strings.Index(content, "import (")
	if importStart == -1 {
		// No import block, add one
		packageEnd := strings.Index(content, "\n\n")
		if packageEnd == -1 {
			return "", fmt.Errorf("could not find appropriate place to add import")
		}
		newImport := fmt.Sprintf("\n\nimport (\n\t%s\n)", domainImport)
		return content[:packageEnd] + newImport + content[packageEnd:], nil
	}

	// Find the end of import block
	importEnd := strings.Index(content[importStart:], ")")
	if importEnd == -1 {
		return "", fmt.Errorf("could not find end of import block")
	}
	importEnd += importStart

	// Add domain import before closing parenthesis
	beforeClose := content[:importEnd]
	afterClose := content[importEnd:]

	// Add newline and tab for proper formatting
	return beforeClose + fmt.Sprintf("\n\t%s", domainImport) + afterClose, nil
}
