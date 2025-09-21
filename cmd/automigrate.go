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

	// Check if entity is already added to auto-migration
	if strings.Contains(contentStr, entityReference) {
		return nil
	}

	updatedContent, err := addEntityToMigrationSlice(contentStr, entityReference)
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

	return "", fmt.Errorf("main.go no encontrado en ninguna ubicación esperada")
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
