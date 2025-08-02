package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// addEntityToAutoMigration adds a domain entity to the auto-migration list in main.go
func addEntityToAutoMigration(entity string) error {
	// Try multiple possible locations for main.go
	possiblePaths := []string{
		filepath.Join("cmd", "server", "main.go"),
		"main.go",
		filepath.Join("cmd", "main.go"),
	}

	var mainPath string
	var found bool

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			mainPath = path
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("main.go no encontrado en ninguna ubicación esperada")
	}

	// Read existing content
	content, err := os.ReadFile(mainPath)
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}

	contentStr := string(content)

	// Check if entity is already added to auto-migration
	entityReference := fmt.Sprintf("&domain.%s{}", entity)
	if strings.Contains(contentStr, entityReference) {
		// Entity already exists in auto-migration
		return nil
	}

	// Look for the comment marker where we should add entities
	commentPattern := "// Add domain entities here as they are created"

	if strings.Contains(contentStr, commentPattern) {
		// Add the entity after the comment
		replacement := fmt.Sprintf("%s\n\t\t%s,", commentPattern, entityReference)
		contentStr = strings.Replace(contentStr, commentPattern, replacement, 1)
	} else {
		// Look for the entities slice pattern as fallback
		entitiesPattern := "entities := []interface{}{"

		if strings.Contains(contentStr, entitiesPattern) {
			// Find the closing brace of the entities slice
			startIdx := strings.Index(contentStr, entitiesPattern)
			if startIdx != -1 {
				// Find the closing brace
				searchStart := startIdx + len(entitiesPattern)
				braceCount := 1
				i := searchStart
				for i < len(contentStr) && braceCount > 0 {
					if contentStr[i] == '{' {
						braceCount++
					} else if contentStr[i] == '}' {
						braceCount--
					}
					i++
				}

				if braceCount == 0 {
					// Insert before the closing brace
					beforeClosing := contentStr[:i-1]
					afterClosing := contentStr[i-1:]
					contentStr = beforeClosing + fmt.Sprintf("\n\t\t%s,", entityReference) + afterClosing
				}
			}
		} else {
			// If neither pattern is found, return an error
			return fmt.Errorf("could not find auto-migration entities slice in main.go")
		}
	}

	// Write updated content
	if err := os.WriteFile(mainPath, []byte(contentStr), 0644); err != nil {
		return fmt.Errorf("failed to update main.go: %w", err)
	}

	return nil
} // generateAutoMigrateFunction generates a function that will be called to auto-migrate entities
func generateAutoMigrateFunction(entities []string) string {
	var entitiesReferences strings.Builder

	for _, entity := range entities {
		entitiesReferences.WriteString(fmt.Sprintf("\t\t&domain.%s{},\n", entity))
	}

	return fmt.Sprintf(`func runAutoMigrations(database *gorm.DB) error {
	if database == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	log.Println("🔄 Running GORM auto-migrations...")
	
	// Create a slice of all domain entities to migrate
	entities := []interface{}{
%s	}
	
	// Run auto-migration for all entities
	for _, entity := range entities {
		if err := database.AutoMigrate(entity); err != nil {
			return fmt.Errorf("failed to auto-migrate entity %%T: %%w", entity, err)
		}
		log.Printf("✅ Auto-migrated entity: %%T", entity)
	}
	
	log.Println("✅ GORM auto-migrations completed successfully")
	return nil
}`, entitiesReferences.String())
}
