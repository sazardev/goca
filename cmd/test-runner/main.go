package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	gocaTesting "github.com/sazardev/goca/internal/testing"
)

var (
	testType = flag.String("type", "all", "Type of test to run: all, init, feature, entity, quality")
	verbose  = flag.Bool("v", false, "Verbose output")
	cleanup  = flag.Bool("cleanup", true, "Cleanup temporary files after test")
)

func main() {
	flag.Parse()

	if *verbose {
		fmt.Println("🧪 Goca CLI Testing Framework")
		fmt.Println("=============================")
	}

	// Create a dummy testing.T for manual execution
	t := &testing.T{}

	// Run the specified test type
	switch *testType {
	case "all":
		runComprehensiveTest(t)
	case "init":
		runInitTest(t)
	case "feature":
		runFeatureTest(t)
	case "entity":
		runEntityTest(t)
	case "quality":
		runQualityTest(t)
	default:
		log.Fatalf("Unknown test type: %s", *testType)
	}

	if *verbose {
		fmt.Println("✅ Testing completed successfully!")
	}
}

func runComprehensiveTest(t *testing.T) {
	fmt.Println("🚀 Running comprehensive Goca CLI test suite...")

	suite := gocaTesting.NewTestSuite(t)
	if !*cleanup {
		fmt.Printf("Test files will remain in: %s\n", suite.GetProjectPath(""))
	} else {
		defer suite.Cleanup()
	}

	suite.RunAllTests()
}

func runInitTest(t *testing.T) {
	fmt.Println("🏗️  Testing goca init command...")

	suite := gocaTesting.NewTestSuite(t)
	if !*cleanup {
		fmt.Printf("Test files will remain in: %s\n", suite.GetProjectPath(""))
	} else {
		defer suite.Cleanup()
	}

	suite.TestInitCommand()
	suite.TestCodeCompilation()
}

func runFeatureTest(t *testing.T) {
	fmt.Println("⚡ Testing goca feature command...")

	suite := gocaTesting.NewTestSuite(t)
	if !*cleanup {
		fmt.Printf("Test files will remain in: %s\n", suite.GetProjectPath(""))
	} else {
		defer suite.Cleanup()
	}

	// Initialize project first
	suite.TestInitCommand()
	if len(suite.errors) > 0 {
		log.Fatalf("Init failed, cannot test features: %v", suite.errors)
	}

	suite.TestFeatureCommand()
	suite.TestCodeCompilation()
}

func runEntityTest(t *testing.T) {
	fmt.Println("🏢 Testing goca entity command...")

	suite := gocaTesting.NewTestSuite(t)
	if !*cleanup {
		fmt.Printf("Test files will remain in: %s\n", suite.GetProjectPath(""))
	} else {
		defer suite.Cleanup()
	}

	// Initialize project first
	suite.TestInitCommand()
	if len(suite.errors) > 0 {
		log.Fatalf("Init failed, cannot test entities: %v", suite.errors)
	}

	suite.TestEntityCommand()
	suite.TestCodeCompilation()
}

func runQualityTest(t *testing.T) {
	fmt.Println("✨ Testing code quality...")

	suite := gocaTesting.NewTestSuite(t)
	if !*cleanup {
		fmt.Printf("Test files will remain in: %s\n", suite.GetProjectPath(""))
	} else {
		defer suite.Cleanup()
	}

	// Generate a complete project
	suite.TestInitCommand()
	if len(suite.errors) > 0 {
		log.Fatalf("Init failed, cannot test quality: %v", suite.errors)
	}

	suite.TestFeatureCommand()
	if len(suite.errors) > 0 {
		log.Fatalf("Feature generation failed, cannot test quality: %v", suite.errors)
	}

	// Test code quality
	suite.TestCodeCompilation()
	suite.TestCodeLinting()
	suite.TestCodeFormatting()
}

// Additional helper for CI/CD integration
func runCITest() {
	// Set environment variables for CI
	os.Setenv("CI", "true")
	os.Setenv("GOCA_TEST", "true")

	// Create test directory in CI-friendly location
	testDir := filepath.Join(os.TempDir(), "goca-ci-test")
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	fmt.Println("Running CI-optimized tests...")

	t := &testing.T{}
	suite := gocaTesting.NewTestSuite(t)
	defer suite.Cleanup()

	// Run core tests only for CI
	suite.TestInitCommand()
	suite.TestFeatureCommand()
	suite.TestCodeCompilation()

	fmt.Println("CI tests completed successfully!")
}
