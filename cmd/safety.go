package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DryRunEntry represents a file that would be created or modified in dry-run mode.
type DryRunEntry struct {
	Path   string
	Action string // "create" or "overwrite"
	Size   int
}

// backups, and dry-run mode.
type SafetyManager struct {
	DryRun       bool
	Force        bool
	Backup       bool
	BackupDir    string
	conflicts    []string
	createdFiles []string
	pendingFiles []DryRunEntry
}

// NewSafetyManager creates a new safety manager instance.
func NewSafetyManager(dryRun, force, backup bool) *SafetyManager {
	return &SafetyManager{
		DryRun:       dryRun,
		Force:        force,
		Backup:       backup,
		BackupDir:    ".goca-backup",
		conflicts:    make([]string, 0),
		createdFiles: make([]string, 0),
		pendingFiles: make([]DryRunEntry, 0),
	}
}

// CheckFileConflict checks if a file exists and handles it according to flags.
func (sm *SafetyManager) CheckFileConflict(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		// File exists
		if sm.DryRun {
			// In dry-run mode an existing file is not an error: it is simply
			// recorded as a conflict and previewed as a "would overwrite".
			sm.conflicts = append(sm.conflicts, filePath)
			return nil
		}

		if !sm.Force {
			return fmt.Errorf("file already exists: %s (use --force to overwrite or --backup to backup first)", filePath)
		}

		if sm.Backup {
			if err := sm.BackupFile(filePath); err != nil {
				return fmt.Errorf("failed to backup file %s: %w", filePath, err)
			}
		}
	}
	return nil
}

// BackupFile creates a backup of an existing file.
func (sm *SafetyManager) BackupFile(filePath string) error {
	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(sm.BackupDir, 0o755); err != nil {
		return err
	}

	// Read original file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Use only the base filename for backup to avoid path issues
	baseFileName := filepath.Base(filePath)
	backupFile := filepath.Join(sm.BackupDir, baseFileName+".backup")
	//#nosec G703 // backup path is derived from validated filePath
	if err := os.WriteFile(backupFile, content, 0o644); err != nil {
		return err
	}

	if ui != nil {
		ui.FileBackedUp(filePath, backupFile)
	} else {
		fmt.Printf("Backed up: %s -> %s\n", filePath, backupFile)
	}
	return nil
}

// WriteFile writes a file with safety checks.
func (sm *SafetyManager) WriteFile(filePath, content string) error {
	// Check for conflicts first
	if err := sm.CheckFileConflict(filePath); err != nil && !sm.Force {
		return err
	}

	if sm.DryRun {
		action := "create"
		if _, statErr := os.Stat(filePath); statErr == nil {
			action = "overwrite"
		}
		entry := DryRunEntry{Path: filePath, Action: action, Size: len(content)}
		sm.pendingFiles = append(sm.pendingFiles, entry)
		sm.createdFiles = append(sm.createdFiles, filePath)
		return nil
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	//#nosec G703 // path is validated by SafetyManager.WriteFile
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	sm.createdFiles = append(sm.createdFiles, filePath)
	if ui != nil {
		ui.FileCreated(filePath)
	} else {
		fmt.Printf("Created: %s\n", filePath)
	}
	return nil
}

// WriteMergedFile writes a file that the caller has intentionally rebuilt from
// the existing content (shared files such as errors.go, dto.go, messages.go and
// the layer interface files). Because the previous content is already merged
// into content, it overwrites without requiring --force, while still honoring
// dry-run and backup.
func (sm *SafetyManager) WriteMergedFile(filePath, content string) error {
	if sm.DryRun {
		action := "create"
		if _, statErr := os.Stat(filePath); statErr == nil {
			action = "overwrite"
		}
		sm.pendingFiles = append(sm.pendingFiles, DryRunEntry{Path: filePath, Action: action, Size: len(content)})
		sm.createdFiles = append(sm.createdFiles, filePath)
		return nil
	}

	if sm.Backup {
		if _, err := os.Stat(filePath); err == nil {
			if err := sm.BackupFile(filePath); err != nil {
				return fmt.Errorf("failed to backup file %s: %w", filePath, err)
			}
		}
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	//#nosec G703 // path is validated by callers
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	sm.createdFiles = append(sm.createdFiles, filePath)
	if ui != nil {
		ui.FileCreated(filePath)
	} else {
		fmt.Printf("Created: %s\n", filePath)
	}
	return nil
}

// GetPendingFiles returns the dry-run entries (files that would be created/overwritten).
func (sm *SafetyManager) GetPendingFiles() []DryRunEntry {
	return sm.pendingFiles
}

// GetConflicts returns list of file conflicts found.
func (sm *SafetyManager) GetConflicts() []string {
	return sm.conflicts
}

// GetCreatedFiles returns list of files that would be/were created.
func (sm *SafetyManager) GetCreatedFiles() []string {
	return sm.createdFiles
}

// PrintSummary prints a summary of the operation.
func (sm *SafetyManager) PrintSummary() {
	if ui != nil {
		sm.printSummaryStyled()
	} else {
		sm.printSummaryPlain()
	}
}

func (sm *SafetyManager) printSummaryStyled() {
	ui.Blank()
	if sm.DryRun {
		ui.Section("DRY-RUN PREVIEW")
		if len(sm.pendingFiles) > 0 {
			rows := make([][]string, len(sm.pendingFiles))
			for i, e := range sm.pendingFiles {
				rows[i] = []string{e.Path, e.Action, fmt.Sprintf("%d B", e.Size)}
			}
			ui.Table([]string{"File", "Action", "Size"}, rows)
		}
		ui.Info(fmt.Sprintf("%d files would be written", len(sm.pendingFiles)))
		if len(sm.conflicts) > 0 {
			ui.Warning(fmt.Sprintf("%d conflicts detected:", len(sm.conflicts)))
			for _, conflict := range sm.conflicts {
				ui.Dim("  - " + conflict)
			}
		}
		ui.Blank()
		ui.NextSteps([]string{
			"Run without --dry-run to actually create files",
			"Use --force to overwrite existing files",
			"Use --backup to backup files before overwriting",
		})
	} else {
		ui.Success(fmt.Sprintf("Successfully created %d files", len(sm.createdFiles)))
		if len(sm.createdFiles) > 0 {
			rows := make([][]string, len(sm.createdFiles))
			for i, f := range sm.createdFiles {
				rows[i] = []string{f, "created"}
			}
			ui.Table([]string{"File", "Status"}, rows)
		}
	}
}

func (sm *SafetyManager) printSummaryPlain() {
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

// NameConflictDetector detects naming conflicts in the project.
type NameConflictDetector struct {
	projectRoot string
	entities    map[string]bool
	features    map[string]bool
}

// NewNameConflictDetector creates a new name conflict detector.
func NewNameConflictDetector(projectRoot string) *NameConflictDetector {
	return &NameConflictDetector{
		projectRoot: projectRoot,
		entities:    make(map[string]bool),
		features:    make(map[string]bool),
	}
}

// ScanExistingEntities scans the project for existing entities.
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
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
			// Extract entity name from filename
			name := strings.TrimSuffix(entry.Name(), ".go")
			name = strings.TrimSuffix(name, "_seeds")
			if name != "errors" && name != "validations" && name != "common" {
				ncd.entities[strings.ToLower(name)] = true
				ncd.features[strings.ToLower(name)] = true
			}
		}
	}

	return nil
}

// CheckNameConflict checks if a name conflicts with existing entities/features.
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

// GetExistingEntities returns list of existing entities.
func (ncd *NameConflictDetector) GetExistingEntities() []string {
	entities := make([]string, 0, len(ncd.entities))
	for entity := range ncd.entities {
		entities = append(entities, entity)
	}
	return entities
}
