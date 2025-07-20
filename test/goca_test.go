package main

import (
	"testing"

	gocaTesting "github.com/cerberus-org/goca/internal/testing"
)

// TestGocaCLI runs the comprehensive test suite for Goca CLI
func TestGocaCLI(t *testing.T) {
	runner := gocaTesting.NewComprehensiveTestRunner()
	runner.RunAllTests(t)
}

// TestBasicProjectInit tests basic project initialization
func TestBasicProjectInit(t *testing.T) {
	suite := gocaTesting.NewTestSuite(t, "basic-init")
	defer suite.Cleanup()

	cli := gocaTesting.NewCLIRunner(suite)
	validator := gocaTesting.NewCodeValidator(suite)

	// Initialize project
	cli.InitProject("testapp", "", "--force")

	// Validate basic structure
	suite.AssertDirectoryExists("internal/domain")
	suite.AssertDirectoryExists("internal/usecase")
	suite.AssertDirectoryExists("internal/infrastructure")
	suite.AssertFileExists("go.mod")
	suite.AssertFileExists("main.go")

	// Validate go.mod syntax
	if err := validator.ValidateGoSyntax(suite.GetProjectPath("go.mod")); err != nil {
		t.Errorf("Invalid go.mod syntax: %v", err)
	}
}

// TestFeatureGeneration tests complete feature generation
func TestFeatureGeneration(t *testing.T) {
	suite := gocaTesting.NewTestSuite(t, "feature-gen")
	defer suite.Cleanup()

	cli := gocaTesting.NewCLIRunner(suite)
	arch := gocaTesting.NewArchitectureValidator(suite)

	// Initialize and generate feature
	cli.InitProject("featureapp", "github.com/test/feature", "--force")
	cli.GenerateFeature("User", "ID:uint,Name:string,Email:string")

	// Validate all files were created
	expectedFiles := []string{
		"internal/domain/user.go",
		"internal/usecase/user_usecase.go",
		"internal/infrastructure/repository/user_repository.go",
		"internal/infrastructure/handler/user_handler.go",
		"pkg/messages/user_messages.go",
		"pkg/interfaces/user_repository.go",
	}

	for _, file := range expectedFiles {
		suite.AssertFileExists(file)
	}

	// Validate architecture compliance
	errors := arch.ValidateProjectStructure(suite.TempDir)
	if len(errors) > 0 {
		t.Errorf("Architecture validation failed with %d errors", len(errors))
		for _, err := range errors {
			t.Logf("Error: %s", err.Error())
		}
	}
}

// TestEntityGeneration tests entity-only generation
func TestEntityGeneration(t *testing.T) {
	suite := gocaTesting.NewTestSuite(t, "entity-gen")
	defer suite.Cleanup()

	cli := gocaTesting.NewCLIRunner(suite)
	validator := gocaTesting.NewCodeValidator(suite)

	// Initialize and generate entity
	cli.InitProject("entityapp", "", "--force")
	cli.GenerateEntity("Product", "ID:uint,Name:string,Price:float64")

	// Validate entity file
	entityFile := suite.GetProjectPath("internal/domain/product.go")
	suite.AssertFileExists("internal/domain/product.go")

	// Validate syntax
	if err := validator.ValidateGoSyntax(entityFile); err != nil {
		t.Errorf("Entity syntax validation failed: %v", err)
	}

	// Validate struct fields
	expectedFields := []string{"ID", "Name", "Price"}
	if err := validator.ValidateStructFields(entityFile, "Product", expectedFields); err != nil {
		t.Errorf("Struct fields validation failed: %v", err)
	}
}

// TestCustomModule tests generation with custom module
func TestCustomModule(t *testing.T) {
	suite := gocaTesting.NewTestSuite(t, "custom-module")
	defer suite.Cleanup()

	cli := gocaTesting.NewCLIRunner(suite)
	validator := gocaTesting.NewCodeValidator(suite)

	customModule := "example.com/my/custom/module"

	// Initialize with custom module
	cli.InitProject("customapp", customModule, "--force")
	cli.GenerateEntity("Order", "ID:uint,Total:float64")

	// Validate module references
	entityFile := suite.GetProjectPath("internal/domain/order.go")
	if err := validator.ValidateModuleReferences(entityFile, customModule); err != nil {
		t.Errorf("Module reference validation failed: %v", err)
	}

	// Check go.mod contains correct module
	suite.AssertFileContains("go.mod", customModule)
}

// TestArchitectureCompliance tests Clean Architecture compliance
func TestArchitectureCompliance(t *testing.T) {
	suite := gocaTesting.NewTestSuite(t, "arch-compliance")
	defer suite.Cleanup()

	cli := gocaTesting.NewCLIRunner(suite)
	arch := gocaTesting.NewArchitectureValidator(suite)

	// Generate full feature
	cli.InitProject("archapp", "github.com/test/arch", "--force")
	cli.GenerateFeature("Payment", "ID:uint,Amount:float64,Status:string")

	// Validate layer separation
	errors := arch.ValidateLayerSeparation(suite.TempDir)
	if len(errors) > 0 {
		t.Errorf("Layer separation validation failed with %d errors", len(errors))
		for _, err := range errors {
			t.Logf("Error: %s", err.Error())
		}
	}

	// Validate entity compliance
	entityFile := suite.GetProjectPath("internal/domain/payment.go")
	if suite.FileExists(entityFile) {
		entityErrors := arch.ValidateEntityCompliance(entityFile)
		if len(entityErrors) > 0 {
			t.Errorf("Entity compliance validation failed with %d errors", len(entityErrors))
			for _, err := range entityErrors {
				t.Logf("Error: %s", err.Error())
			}
		}
	}
}

// TestCLICommands tests CLI command functionality
func TestCLICommands(t *testing.T) {
	suite := gocaTesting.NewTestSuite(t, "cli-commands")
	defer suite.Cleanup()

	cli := gocaTesting.NewCLIRunner(suite)

	// Test version command
	version := cli.Version()
	if version == "" {
		t.Error("Version command returned empty output")
	}

	// Test help command
	help := cli.Help()
	if help == "" {
		t.Error("Help command returned empty output")
	}

	// Test invalid command should fail
	if _, err := cli.Run("invalid-command"); err == nil {
		t.Error("Invalid command should have failed but succeeded")
	}
}

// TestErrorHandling tests CLI error handling
func TestErrorHandling(t *testing.T) {
	suite := gocaTesting.NewTestSuite(t, "error-handling")
	defer suite.Cleanup()

	cli := gocaTesting.NewCLIRunner(suite)

	// Test missing required flags
	if _, err := cli.Run("entity", "TestEntity"); err == nil {
		t.Error("Entity command without fields flag should fail")
	}

	if _, err := cli.Run("usecase", "TestUseCase"); err == nil {
		t.Error("UseCase command without entity flag should fail")
	}
}

// BenchmarkFeatureGeneration benchmarks feature generation performance
func BenchmarkFeatureGeneration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		suite := gocaTesting.NewTestSuite(&testing.T{}, "benchmark")
		cli := gocaTesting.NewCLIRunner(suite)

		cli.InitProject("benchapp", "", "--force")
		cli.GenerateFeature("BenchEntity", "ID:uint,Name:string,Value:int")

		suite.Cleanup()
	}
}
