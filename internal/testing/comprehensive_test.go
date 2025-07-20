package testing

import (
	"fmt"
	"testing"
)

// TestGocaCLIComprehensive runs the complete test suite for Goca CLI
func TestGocaCLIComprehensive(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Cleanup()

	suite.RunAllTests()
}

// TestGocaInitCommand tests only the goca init command
func TestGocaInitCommand(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Cleanup()

	suite.TestInitCommand()
	suite.TestCodeCompilation()
}

// TestGocaFeatureCommand tests only the goca feature command
func TestGocaFeatureCommand(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Cleanup()

	// First init a project
	suite.TestInitCommand()
	if len(suite.errors) > 0 {
		t.Fatalf("Init failed, cannot test feature command: %v", suite.errors)
	}

	// Then test features
	suite.TestFeatureCommand()
	suite.TestCodeCompilation()
}

// TestGocaEntityCommand tests only the goca entity command
func TestGocaEntityCommand(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Cleanup()

	// First init a project
	suite.TestInitCommand()
	if len(suite.errors) > 0 {
		t.Fatalf("Init failed, cannot test entity command: %v", suite.errors)
	}

	// Then test entities
	suite.TestEntityCommand()
	suite.TestCodeCompilation()
}

// TestCodeQuality tests the quality of generated code
func TestCodeQuality(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Cleanup()

	// Generate a full project
	suite.TestInitCommand()
	if len(suite.errors) > 0 {
		t.Fatalf("Init failed, cannot test code quality: %v", suite.errors)
	}

	suite.TestFeatureCommand()
	if len(suite.errors) > 0 {
		t.Fatalf("Feature generation failed, cannot test code quality: %v", suite.errors)
	}

	// Test code quality
	suite.TestCodeCompilation()
	suite.TestCodeLinting()
	suite.TestCodeFormatting()
}

// BenchmarkGocaInit benchmarks the goca init command performance
func BenchmarkGocaInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		suite := NewTestSuite(&testing.T{})
		suite.TestInitCommand()
		suite.Cleanup()
	}
}

// BenchmarkGocaFeature benchmarks the goca feature command performance
func BenchmarkGocaFeature(b *testing.B) {
	// Setup once
	suite := NewTestSuite(&testing.T{})
	suite.TestInitCommand()
	defer suite.Cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a unique entity name for each iteration
		entityName := fmt.Sprintf("TestEntity%d", i)
		stdout, stderr, err := suite.runGocaCommand("feature", entityName,
			"--fields", "name:string,value:int", "--validation")

		if err != nil {
			b.Errorf("Feature command failed: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
		}
	}
}
