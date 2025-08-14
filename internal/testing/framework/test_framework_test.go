package framework

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFilePathResolution(t *testing.T) {
	// Create a test context
	tc := NewTestContext(t)
	defer tc.Cleanup()
	tc.CurrentTestName = "TestFilePathResolution"

	// Create a test directory structure
	testProjectDir := filepath.Join(tc.TempDir, "test-project")
	err := os.MkdirAll(filepath.Join(testProjectDir, "internal", "domain"), 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a test file
	userEntityPath := filepath.Join(testProjectDir, "internal", "domain", "user.go")
	err = os.WriteFile(userEntityPath, []byte("package domain\n\ntype User struct {\n\tName string\n\tEmail string\n}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test file existence checks
	testCases := []string{
		"test-project/internal/domain/user.go",
		"internal/domain/user.go",
	}

	for _, testPath := range testCases {
		t.Logf("Testing path: %s", testPath)
		if !tc.AssertFileExists(testPath) {
			t.Errorf("AssertFileExists failed for path: %s", testPath)
		}
	}

	// Test file content checks
	for _, testPath := range testCases {
		if !tc.AssertFileContains(testPath, "type User struct {") {
			t.Errorf("AssertFileContains failed for path: %s", testPath)
		}
	}

	// List all files for debugging
	tc.ListProjectFiles()

	// Summary
	tc.PrintTestSummary()
}
