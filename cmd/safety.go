package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SafetyManager handles file safety operations like conflict detection,
// backups, and dry-run mode
type SafetyManager struct {
	DryRun       bool
	Force        bool
	Backup       bool
	BackupDir    string
	conflicts    []string
	createdFiles []string
}

// NewSafetyManager creates a new safety manager instance
func NewSafetyManager(dryRun, force, backup bool) *SafetyManager {
	return &SafetyManager{
		DryRun:       dryRun,
		Force:        force,
		Backup:       backup,
		BackupDir:    ".goca-backup",
		conflicts:    make([]string, 0),
		createdFiles: make([]string, 0),
	}
}

// CheckFileConflict checks if a file exists and handles it according to flags
func (sm *SafetyManager) CheckFileConflict(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		// File exists
		if sm.DryRun {
			sm.conflicts = append(sm.conflicts, filePath)
			return fmt.Errorf("file already exists (dry-run): %s", filePath)
		}

		if !sm.Force {
			return fmt.Errorf("file already exists: %s (use --force to overwrite or --backup to backup first)", filePath)
		}

		if sm.Backup {
			if err := sm.BackupFile(filePath); err != nil {
				return fmt.Errorf("failed to backup file %s: %v", filePath, err)
			}
		}
	}
	return nil
}

// BackupFile creates a backup of an existing file
func (sm *SafetyManager) BackupFile(filePath string) error {
	// Create backup directory if it doesn't exist
	backupPath := filepath.Join(sm.BackupDir, filepath.Dir(filePath))
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return err
	}

	// Read original file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Write to backup location with timestamp
	backupFile := filepath.Join(sm.BackupDir, filePath+".backup")
	if err := os.WriteFile(backupFile, content, 0644); err != nil {
		return err
	}

	fmt.Printf("Backed up: %s -> %s\n", filePath, backupFile)
	return nil
}

// WriteFile writes a file with safety checks
func (sm *SafetyManager) WriteFile(filePath, content string) error {
	// Check for conflicts first
	if err := sm.CheckFileConflict(filePath); err != nil && !sm.Force {
		return err
	}

	if sm.DryRun {
		fmt.Printf("[DRY-RUN] Would create: %s (%d bytes)\n", filePath, len(content))
		sm.createdFiles = append(sm.createdFiles, filePath)
		return nil
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %v", filePath, err)
	}

	sm.createdFiles = append(sm.createdFiles, filePath)
	fmt.Printf("Created: %s\n", filePath)
	return nil
}

// GetConflicts returns list of file conflicts found
func (sm *SafetyManager) GetConflicts() []string {
	return sm.conflicts
}

// GetCreatedFiles returns list of files that would be/were created
func (sm *SafetyManager) GetCreatedFiles() []string {
	return sm.createdFiles
}

// PrintSummary prints a summary of the operation
func (sm *SafetyManager) PrintSummary() {
	if sm.DryRun {
		fmt.Println("\nDRY-RUN SUMMARY:")
		fmt.Printf("   Would create %d files\n", len(sm.createdFiles))
		if len(sm.conflicts) > 0 {
			fmt.Printf("   Warning: %d conflicts detected:\n", len(sm.conflicts))
			for _, conflict := range sm.conflicts {
				fmt.Printf("      - %s\n", conflict)
			}
		}
		fmt.Println("\nTip: Run without --dry-run to actually create files")
		fmt.Println("   Use --force to overwrite existing files")
		fmt.Println("   Use --backup to backup files before overwriting")
	} else {
		fmt.Printf("\nSuccessfully created %d files\n", len(sm.createdFiles))
	}
}

// NameConflictDetector detects naming conflicts in the project
type NameConflictDetector struct {
	projectRoot string
	entities    map[string]bool
	features    map[string]bool
}

// NewNameConflictDetector creates a new name conflict detector
func NewNameConflictDetector(projectRoot string) *NameConflictDetector {
	return &NameConflictDetector{
		projectRoot: projectRoot,
		entities:    make(map[string]bool),
		features:    make(map[string]bool),
	}
}

// ScanExistingEntities scans the project for existing entities
func (ncd *NameConflictDetector) ScanExistingEntities() error {
	domainPath := filepath.Join(ncd.projectRoot, "internal", "domain")

	if _, err := os.Stat(domainPath); os.IsNotExist(err) {
		return nil // No domain directory yet
	}

	entries, err := os.ReadDir(domainPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			// Extract entity name from filename
			name := strings.TrimSuffix(entry.Name(), ".go")
			name = strings.TrimSuffix(name, "_seeds")
			if name != "errors" {
				ncd.entities[strings.ToLower(name)] = true
				ncd.features[strings.ToLower(name)] = true
			}
		}
	}

	return nil
}

// CheckNameConflict checks if a name conflicts with existing entities/features
func (ncd *NameConflictDetector) CheckNameConflict(name string) error {
	normalizedName := strings.ToLower(name)

	if ncd.entities[normalizedName] {
		return fmt.Errorf("entity '%s' already exists in the project", name)
	}

	if ncd.features[normalizedName] {
		return fmt.Errorf("feature '%s' already exists in the project", name)
	}

	return nil
}

// GetExistingEntities returns list of existing entities
func (ncd *NameConflictDetector) GetExistingEntities() []string {
	entities := make([]string, 0, len(ncd.entities))
	for entity := range ncd.entities {
		entities = append(entities, entity)
	}
	return entities
}
