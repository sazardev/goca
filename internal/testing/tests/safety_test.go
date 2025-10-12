package tests
package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sazardev/goca/cmd"
)

// TestSafetyManager tests the SafetyManager functionality
func TestSafetyManager(t *testing.T) {
	t.Run("DryRunMode", func(t *testing.T) {
		sm := cmd.NewSafetyManager(true, false, false)
		
		err := sm.WriteFile("test.go", "package test")
		if err != nil {
			t.Errorf("WriteFile in dry-run should not error: %v", err)
		}

		files := sm.GetCreatedFiles()
		if len(files) != 1 {
			t.Errorf("Expected 1 file in dry-run, got %d", len(files))
		}
	})

	t.Run("FileConflictDetection", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "existing.go")
		
		// Create existing file
		err := os.WriteFile(testFile, []byte("existing content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		sm := cmd.NewSafetyManager(false, false, false)
		err = sm.CheckFileConflict(testFile)
		
		if err == nil {
			t.Error("Expected conflict error for existing file")
		}
	})

	t.Run("ForceOverwrite", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "existing.go")
		
		// Create existing file
		err := os.WriteFile(testFile, []byte("existing content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		sm := cmd.NewSafetyManager(false, true, false)
		err = sm.CheckFileConflict(testFile)
		
		if err != nil {
			t.Errorf("Force mode should not error on conflict: %v", err)
		}
	})

	t.Run("BackupFile", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.go")
		originalContent := "original content"
		
		err := os.WriteFile(testFile, []byte(originalContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		sm := cmd.NewSafetyManager(false, true, true)
		sm.BackupDir = filepath.Join(tempDir, ".backup")
		
		err = sm.BackupFile(testFile)
		if err != nil {
			t.Errorf("Backup failed: %v", err)
		}

		backupFile := filepath.Join(sm.BackupDir, "test.go.backup")
		content, err := os.ReadFile(backupFile)
		if err != nil {
			t.Errorf("Could not read backup file: %v", err)
		}

		if string(content) != originalContent {
			t.Errorf("Backup content mismatch. Expected '%s', got '%s'", originalContent, string(content))
		}
	})
}

// TestNameConflictDetector tests the NameConflictDetector functionality
func TestNameConflictDetector(t *testing.T) {
	t.Run("ScanExistingEntities", func(t *testing.T) {
		tempDir := t.TempDir()
		domainDir := filepath.Join(tempDir, "internal", "domain")
		err := os.MkdirAll(domainDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create domain directory: %v", err)
		}

		// Create test entity files
		entities := []string{"user.go", "product.go", "order.go", "errors.go"}
		for _, entity := range entities {
			err := os.WriteFile(filepath.Join(domainDir, entity), []byte("package domain"), 0644)
			if err != nil {
				t.Fatalf("Failed to create entity file %s: %v", entity, err)
			}
		}

		detector := cmd.NewNameConflictDetector(tempDir)
		err = detector.ScanExistingEntities()
		if err != nil {
			t.Errorf("Scan failed: %v", err)
		}

		existingEntities := detector.GetExistingEntities()
		// Should find 3 entities (excluding errors.go)
		if len(existingEntities) != 3 {
			t.Errorf("Expected 3 entities, found %d: %v", len(existingEntities), existingEntities)
		}
	})

	t.Run("DetectNameConflict", func(t *testing.T) {
		tempDir := t.TempDir()
		domainDir := filepath.Join(tempDir, "internal", "domain")
		err := os.MkdirAll(domainDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create domain directory: %v", err)
		}

		// Create existing User entity
		err = os.WriteFile(filepath.Join(domainDir, "user.go"), []byte("package domain"), 0644)
		if err != nil {
			t.Fatalf("Failed to create user.go: %v", err)
		}

		detector := cmd.NewNameConflictDetector(tempDir)
		err = detector.ScanExistingEntities()
		if err != nil {
			t.Errorf("Scan failed: %v", err)
		}

		// Should detect conflict
		err = detector.CheckNameConflict("User")
		if err == nil {
			t.Error("Expected conflict for duplicate entity 'User'")
		}

		// Should not detect conflict for new entity
		err = detector.CheckNameConflict("Product")
		if err != nil {
			t.Errorf("Should not have conflict for new entity: %v", err)
		}
	})

	t.Run("CaseInsensitiveDetection", func(t *testing.T) {
		tempDir := t.TempDir()
		domainDir := filepath.Join(tempDir, "internal", "domain")
		err := os.MkdirAll(domainDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create domain directory: %v", err)
		}

		err = os.WriteFile(filepath.Join(domainDir, "user.go"), []byte("package domain"), 0644)
		if err != nil {
			t.Fatalf("Failed to create user.go: %v", err)
		}

		detector := cmd.NewNameConflictDetector(tempDir)
		err = detector.ScanExistingEntities()
		if err != nil {
			t.Errorf("Scan failed: %v", err)
		}

		// Should detect conflict regardless of case
		testCases := []string{"User", "user", "USER", "uSeR"}
		for _, testCase := range testCases {
			err = detector.CheckNameConflict(testCase)
			if err == nil {
				t.Errorf("Expected conflict for '%s'", testCase)
			}
		}
	})
}

// TestDependencyManager tests the DependencyManager functionality
func TestDependencyManager(t *testing.T) {
	t.Run("CommonDependencies", func(t *testing.T) {
		tempDir := t.TempDir()
		dm := cmd.NewDependencyManager(tempDir, false)
		
		deps := dm.CommonDependencies()
		
		// Check for known dependencies
		if _, ok := deps["validator"]; !ok {
			t.Error("Expected 'validator' in common dependencies")
		}
		if _, ok := deps["jwt"]; !ok {
			t.Error("Expected 'jwt' in common dependencies")
		}
		if _, ok := deps["grpc"]; !ok {
			t.Error("Expected 'grpc' in common dependencies")
		}
	})

	t.Run("SuggestDependencies", func(t *testing.T) {
		tempDir := t.TempDir()
		dm := cmd.NewDependencyManager(tempDir, false)
		
		suggestions := dm.SuggestDependencies([]string{"validation", "auth", "grpc"})
		
		if len(suggestions) == 0 {
			t.Error("Expected dependency suggestions")
		}

		// Check that suggestions include relevant modules
		hasValidator := false
		hasJWT := false
		hasGRPC := false
		
		for _, dep := range suggestions {
			if dep.Module == "github.com/go-playground/validator/v10" {
				hasValidator = true
			}
			if dep.Module == "github.com/golang-jwt/jwt/v5" {
				hasJWT = true
			}
			if dep.Module == "google.golang.org/grpc" {
				hasGRPC = true
			}
		}

		if !hasValidator {
			t.Error("Expected validator suggestion for 'validation' feature")
		}
		if !hasJWT {
			t.Error("Expected JWT suggestion for 'auth' feature")
		}
		if !hasGRPC {
			t.Error("Expected gRPC suggestion for 'grpc' feature")
		}
	})

	t.Run("GetRequiredDependencies", func(t *testing.T) {
		tempDir := t.TempDir()
		dm := cmd.NewDependencyManager(tempDir, false)
		
		required := dm.GetRequiredDependenciesForFeature("grpc", map[string]bool{
			"validation": true,
		})

		if len(required) == 0 {
			t.Error("Expected required dependencies")
		}

		// Should include gRPC, protobuf, and validator
		hasGRPC := false
		hasProtobuf := false
		hasValidator := false

		for _, dep := range required {
			if dep.Module == "google.golang.org/grpc" {
				hasGRPC = true
			}
			if dep.Module == "google.golang.org/protobuf" {
				hasProtobuf = true
			}
			if dep.Module == "github.com/go-playground/validator/v10" {
				hasValidator = true
			}
		}

		if !hasGRPC {
			t.Error("Expected gRPC as required dependency")
		}
		if !hasProtobuf {
			t.Error("Expected protobuf as required dependency")
		}
		if !hasValidator {
			t.Error("Expected validator as required dependency for validation option")
		}
	})

	t.Run("DryRunMode", func(t *testing.T) {
		tempDir := t.TempDir()
		dm := cmd.NewDependencyManager(tempDir, true)
		
		dep := cmd.Dependency{
			Module:  "github.com/test/module",
			Version: "v1.0.0",
			Type:    "optional",
			Reason:  "testing",
		}

		// Should not actually add dependency in dry-run
		err := dm.AddDependency(dep)
		if err != nil {
			t.Errorf("AddDependency in dry-run should not error: %v", err)
		}
	})
}
