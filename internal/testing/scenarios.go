package testing

import (
	"fmt"
	"path/filepath"
	"testing"
)

// TestScenario represents a comprehensive test scenario for the CLI
type TestScenario struct {
	Name        string
	Description string
	Setup       func(*TestSuite) error
	Execute     func(*TestSuite) error
	Validate    func(*TestSuite) []*TestError
	Cleanup     func(*TestSuite) error
}

// ComprehensiveTestRunner executes all test scenarios
type ComprehensiveTestRunner struct {
	scenarios []TestScenario
}

// NewComprehensiveTestRunner creates a new comprehensive test runner
func NewComprehensiveTestRunner() *ComprehensiveTestRunner {
	runner := &ComprehensiveTestRunner{}
	runner.registerAllScenarios()
	return runner
}

// RunAllTests executes all registered test scenarios
func (r *ComprehensiveTestRunner) RunAllTests(t *testing.T) {
	fmt.Printf("🚀 Starting Goca CLI Comprehensive Testing Suite\n")
	fmt.Printf("Total scenarios: %d\n\n", len(r.scenarios))

	var allErrors []*TestError
	passedScenarios := 0

	for i, scenario := range r.scenarios {
		fmt.Printf("📋 [%d/%d] Running: %s\n", i+1, len(r.scenarios), scenario.Name)
		fmt.Printf("   %s\n", scenario.Description)

		suite := NewTestSuite(t)
		errors := r.runScenario(scenario, suite)

		if len(errors) == 0 {
			fmt.Printf("   ✅ PASSED\n\n")
			passedScenarios++
		} else {
			fmt.Printf("   ❌ FAILED (%d errors)\n", len(errors))
			for _, err := range errors {
				fmt.Printf("      - %s\n", err.Error())
			}
			fmt.Printf("\n")
		}

		allErrors = append(allErrors, errors...)
		suite.Cleanup()
	}

	// Print final summary
	r.printFinalSummary(passedScenarios, allErrors)

	// Fail the test if there are critical errors
	summary := NewErrorSummary(allErrors)
	if summary.HasCriticalErrors() {
		t.Fatalf("Test suite failed with %d critical errors", summary.Critical)
	}
}

// runScenario executes a single test scenario
func (r *ComprehensiveTestRunner) runScenario(scenario TestScenario, suite *TestSuite) []*TestError {
	var errors []*TestError

	// Setup phase
	if scenario.Setup != nil {
		if err := scenario.Setup(suite); err != nil {
			return []*TestError{NewGenerationError("setup", err.Error())}
		}
	}

	// Execute phase
	if scenario.Execute != nil {
		if err := scenario.Execute(suite); err != nil {
			return []*TestError{NewGenerationError("execution", err.Error())}
		}
	}

	// Validate phase
	if scenario.Validate != nil {
		errors = scenario.Validate(suite)
	}

	// Cleanup phase
	if scenario.Cleanup != nil {
		if err := scenario.Cleanup(suite); err != nil {
			errors = append(errors, NewGenerationError("cleanup", err.Error()))
		}
	}

	return errors
}

// registerAllScenarios registers all test scenarios
func (r *ComprehensiveTestRunner) registerAllScenarios() {
	r.scenarios = []TestScenario{
		r.createBasicInitScenario(),
		r.createFullFeatureScenario(),
		r.createEntityOnlyScenario(),
		r.createUseCaseOnlyScenario(),
		r.createRepositoryOnlyScenario(),
		r.createHandlerOnlyScenario(),
		r.createDIScenario(),
		r.createInterfacesScenario(),
		r.createMessagesScenario(),
		r.createCustomModuleScenario(),
		r.createComplexFieldsScenario(),
		r.createArchitectureValidationScenario(),
		r.createErrorHandlingScenario(),
		r.createFlagValidationScenario(),
		r.createVersionAndHelpScenario(),
	}
}

// Scenario implementations

func (r *ComprehensiveTestRunner) createBasicInitScenario() TestScenario {
	return TestScenario{
		Name:        "Basic Project Initialization",
		Description: "Tests basic project initialization with default settings",
		Execute: func(suite *TestSuite) error {
			cli := NewCLIRunner(suite)
			cli.InitProject("testapp", "", "--force")
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			arch := NewArchitectureValidator(suite)
			code := NewCodeValidator(suite)

			var errors []*TestError

			// Validate project structure
			errors = append(errors, arch.ValidateProjectStructure(suite.tempDir)...)

			// Validate all generated Go files have valid syntax
			if syntaxErrors := code.ValidateAllGoFiles(suite.tempDir); len(syntaxErrors) > 0 {
				for _, err := range syntaxErrors {
					errors = append(errors, NewSyntaxError("", err.Error()))
				}
			}

			return errors
		},
	}
}

func (r *ComprehensiveTestRunner) createFullFeatureScenario() TestScenario {
	return TestScenario{
		Name:        "Full Feature Generation",
		Description: "Tests complete feature generation with all components",
		Execute: func(suite *TestSuite) error {
			cli := NewCLIRunner(suite)
			cli.InitProject("fullapp", "github.com/test/fullapp", "--force")
			cli.GenerateFeature("User", "ID:uint,Name:string,Email:string,Active:bool")
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			var errors []*TestError

			// Validate all expected files exist
			expectedFiles := []string{
				"internal/domain/user.go",
				"internal/usecase/user_usecase.go",
				"internal/repository/user_repository.go",
				"internal/handler/http/user_handler.go",
				"pkg/messages/user_messages.go",
				"pkg/interfaces/user_repository.go",
				"internal/di/container.go",
			}

			for _, file := range expectedFiles {
				fullPath := filepath.Join(suite.tempDir, file)
				if !suite.FileExists(fullPath) {
					errors = append(errors, NewFileError(file, "existence", "file not found"))
				}
			}

			// Validate architecture compliance
			arch := NewArchitectureValidator(suite)
			if entityFile := filepath.Join(suite.tempDir, "internal/domain/user.go"); suite.FileExists(entityFile) {
				errors = append(errors, arch.ValidateEntityCompliance(entityFile)...)
			}

			return errors
		},
	}
}

func (r *ComprehensiveTestRunner) createEntityOnlyScenario() TestScenario {
	return TestScenario{
		Name:        "Entity Only Generation",
		Description: "Tests generating only entity without other components",
		Execute: func(suite *TestSuite) error {
			cli := NewCLIRunner(suite)
			cli.InitProject("entityapp", "github.com/test/entityapp", "--force")
			cli.GenerateEntity("Product", "ID:uint,Name:string,Price:float64,CategoryID:uint")
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			var errors []*TestError

			entityFile := filepath.Join(suite.tempDir, "internal/domain/product.go")
			if !suite.FileExists(entityFile) {
				errors = append(errors, NewFileError(entityFile, "existence", "entity file not found"))
				return errors
			}

			// Validate entity structure
			code := NewCodeValidator(suite)
			expectedFields := []string{"ID", "Name", "Price", "CategoryID"}
			if err := code.ValidateStructFields(entityFile, "Product", expectedFields); err != nil {
				errors = append(errors, NewComplianceError(entityFile, "struct", err.Error()))
			}

			return errors
		},
	}
}

func (r *ComprehensiveTestRunner) createCustomModuleScenario() TestScenario {
	return TestScenario{
		Name:        "Custom Module Generation",
		Description: "Tests generation with custom module name",
		Execute: func(suite *TestSuite) error {
			cli := NewCLIRunner(suite)
			cli.InitProject("customapp", "example.com/my/custom/module", "--force")
			cli.GenerateFeature("Order", "ID:uint,CustomerID:uint,Total:float64")
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			var errors []*TestError

			// Check all files use the correct module
			code := NewCodeValidator(suite)
			files := []string{
				"internal/domain/order.go",
				"internal/usecase/order_usecase.go",
				"internal/repository/order_repository.go",
			}

			for _, file := range files {
				fullPath := filepath.Join(suite.tempDir, file)
				if suite.FileExists(fullPath) {
					if err := code.ValidateModuleReferences(fullPath, "example.com/my/custom/module"); err != nil {
						errors = append(errors, NewImportError(file, "module", err.Error()))
					}
				}
			}

			return errors
		},
	}
}

func (r *ComprehensiveTestRunner) createArchitectureValidationScenario() TestScenario {
	return TestScenario{
		Name:        "Architecture Validation",
		Description: "Tests strict Clean Architecture compliance",
		Execute: func(suite *TestSuite) error {
			cli := NewCLIRunner(suite)
			cli.InitProject("archapp", "github.com/test/arch", "--force")
			cli.GenerateFeature("Payment", "ID:uint,Amount:float64,Status:string")
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			arch := NewArchitectureValidator(suite)

			var errors []*TestError

			// Validate layer separation
			errors = append(errors, arch.ValidateLayerSeparation(suite.tempDir)...)

			// Validate individual components
			files := map[string]func(string) []*TestError{
				"internal/domain/payment.go":                arch.ValidateEntityCompliance,
				"internal/usecase/payment_usecase.go":       arch.ValidateUseCaseCompliance,
				"internal/repository/payment_repository.go": arch.ValidateRepositoryCompliance,
				"internal/handler/http/payment_handler.go":  arch.ValidateHandlerCompliance,
			}

			for file, validator := range files {
				fullPath := filepath.Join(suite.tempDir, file)
				if suite.FileExists(fullPath) {
					errors = append(errors, validator(fullPath)...)
				}
			}

			return errors
		},
	}
}

func (r *ComprehensiveTestRunner) createErrorHandlingScenario() TestScenario {
	return TestScenario{
		Name:        "Error Handling",
		Description: "Tests CLI error handling and validation",
		Execute: func(suite *TestSuite) error {
			// This scenario intentionally tests error conditions
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			var errors []*TestError
			cli := NewCLIRunner(suite)

			// Test invalid commands should fail
			if _, err := cli.Run("invalid-command"); err == nil {
				errors = append(errors, NewFlagError("command", "error", "success", "invalid command should fail"))
			}

			// Test missing required flags
			if _, err := cli.Run("entity", "TestEntity"); err == nil {
				errors = append(errors, NewFlagError("fields", "error", "success", "missing fields flag should fail"))
			}

			return errors
		},
	}
}

func (r *ComprehensiveTestRunner) createVersionAndHelpScenario() TestScenario {
	return TestScenario{
		Name:        "Version and Help Commands",
		Description: "Tests version and help functionality",
		Execute: func(suite *TestSuite) error {
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			var errors []*TestError
			cli := NewCLIRunner(suite)

			// Test version command
			version := cli.Version()
			if version == "" {
				errors = append(errors, NewFlagError("version", "non-empty", "empty", "version should return output"))
			}

			// Test help command
			help := cli.Help()
			if help == "" {
				errors = append(errors, NewFlagError("help", "non-empty", "empty", "help should return output"))
			}

			return errors
		},
	}
}

// Helper scenarios for individual components
func (r *ComprehensiveTestRunner) createUseCaseOnlyScenario() TestScenario {
	return TestScenario{
		Name:        "UseCase Only Generation",
		Description: "Tests generating only usecase component",
		Execute: func(suite *TestSuite) error {
			cli := NewCLIRunner(suite)
			cli.InitProject("ucapp", "", "--force")
			cli.GenerateEntity("Invoice", "ID:uint,Number:string,Amount:float64")
			cli.GenerateUseCase("ProcessInvoice", "Invoice")
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			usecaseFile := filepath.Join(suite.tempDir, "internal/usecase/process_invoice_usecase.go")
			if !suite.FileExists(usecaseFile) {
				return []*TestError{NewFileError(usecaseFile, "existence", "usecase file not found")}
			}

			arch := NewArchitectureValidator(suite)
			return arch.ValidateUseCaseCompliance(usecaseFile)
		},
	}
}

func (r *ComprehensiveTestRunner) createRepositoryOnlyScenario() TestScenario {
	return TestScenario{
		Name:        "Repository Only Generation",
		Description: "Tests generating only repository component",
		Execute: func(suite *TestSuite) error {
			cli := NewCLIRunner(suite)
			cli.InitProject("repoapp", "", "--force")
			cli.GenerateEntity("Customer", "ID:uint,Name:string,Email:string")
			cli.GenerateRepository("Customer")
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			repoFile := filepath.Join(suite.tempDir, "internal/repository/customer_repository.go")
			if !suite.FileExists(repoFile) {
				return []*TestError{NewFileError(repoFile, "existence", "repository file not found")}
			}

			arch := NewArchitectureValidator(suite)
			return arch.ValidateRepositoryCompliance(repoFile)
		},
	}
}

func (r *ComprehensiveTestRunner) createHandlerOnlyScenario() TestScenario {
	return TestScenario{
		Name:        "Handler Only Generation",
		Description: "Tests generating only handler component",
		Execute: func(suite *TestSuite) error {
			cli := NewCLIRunner(suite)
			cli.InitProject("handlerapp", "", "--force")
			cli.GenerateEntity("Account", "ID:uint,Balance:float64,OwnerID:uint")
			cli.GenerateHandler("Account")
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			handlerFile := filepath.Join(suite.tempDir, "internal/handler/http/account_handler.go")
			if !suite.FileExists(handlerFile) {
				return []*TestError{NewFileError(handlerFile, "existence", "handler file not found")}
			}

			arch := NewArchitectureValidator(suite)
			return arch.ValidateHandlerCompliance(handlerFile)
		},
	}
}

func (r *ComprehensiveTestRunner) createDIScenario() TestScenario {
	return TestScenario{
		Name:        "Dependency Injection Generation",
		Description: "Tests DI container generation",
		Execute: func(suite *TestSuite) error {
			cli := NewCLIRunner(suite)
			cli.InitProject("diapp", "", "--force")
			cli.GenerateFeature("Transaction", "ID:uint,Amount:float64,Type:string")
			cli.GenerateDI("Transaction")
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			diFile := filepath.Join(suite.tempDir, "internal/infrastructure/di/container.go")
			if !suite.FileExists(diFile) {
				return []*TestError{NewFileError(diFile, "existence", "DI file not found")}
			}

			arch := NewArchitectureValidator(suite)
			return arch.ValidateDependencyInjection(diFile)
		},
	}
}

func (r *ComprehensiveTestRunner) createInterfacesScenario() TestScenario {
	return TestScenario{
		Name:        "Interfaces Generation",
		Description: "Tests interface generation",
		Execute: func(suite *TestSuite) error {
			cli := NewCLIRunner(suite)
			cli.InitProject("intapp", "", "--force")
			cli.GenerateEntity("Category", "ID:uint,Name:string,Description:string")
			cli.GenerateInterfaces("Category")
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			interfaceFile := filepath.Join(suite.tempDir, "pkg/interfaces/category_repository.go")
			if !suite.FileExists(interfaceFile) {
				return []*TestError{NewFileError(interfaceFile, "existence", "interface file not found")}
			}

			code := NewCodeValidator(suite)
			expectedMethods := []string{"Create", "GetByID", "Update", "Delete", "List"}
			if err := code.ValidateInterfaceMethods(interfaceFile, "CategoryRepository", expectedMethods); err != nil {
				return []*TestError{NewComplianceError(interfaceFile, "interface", err.Error())}
			}

			return nil
		},
	}
}

func (r *ComprehensiveTestRunner) createMessagesScenario() TestScenario {
	return TestScenario{
		Name:        "Messages Generation",
		Description: "Tests message types generation",
		Execute: func(suite *TestSuite) error {
			cli := NewCLIRunner(suite)
			cli.InitProject("msgapp", "", "--force")
			cli.GenerateEntity("Notification", "ID:uint,Title:string,Content:string")
			cli.GenerateMessages("Notification")
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			msgFile := filepath.Join(suite.tempDir, "pkg/messages/notification_messages.go")
			if !suite.FileExists(msgFile) {
				return []*TestError{NewFileError(msgFile, "existence", "messages file not found")}
			}

			code := NewCodeValidator(suite)
			expectedStructs := []string{"CreateNotificationRequest", "UpdateNotificationRequest", "NotificationResponse"}

			var errors []*TestError
			for _, structName := range expectedStructs {
				if err := code.ValidateStructFields(msgFile, structName, []string{}); err != nil {
					errors = append(errors, NewComplianceError(msgFile, "message struct",
						fmt.Sprintf("struct %s: %s", structName, err.Error())))
				}
			}

			return errors
		},
	}
}

func (r *ComprehensiveTestRunner) createComplexFieldsScenario() TestScenario {
	return TestScenario{
		Name:        "Complex Fields Generation",
		Description: "Tests generation with complex field types",
		Execute: func(suite *TestSuite) error {
			cli := NewCLIRunner(suite)
			cli.InitProject("complexapp", "", "--force")
			cli.GenerateEntity("Document", "ID:uint,Title:string,Content:[]byte,Tags:[]string,Metadata:map[string]interface{},CreatedAt:time.Time")
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			entityFile := filepath.Join(suite.tempDir, "internal/domain/document.go")
			if !suite.FileExists(entityFile) {
				return []*TestError{NewFileError(entityFile, "existence", "entity file not found")}
			}

			code := NewCodeValidator(suite)
			expectedFields := []string{"ID", "Title", "Content", "Tags", "Metadata", "CreatedAt"}
			if err := code.ValidateStructFields(entityFile, "Document", expectedFields); err != nil {
				return []*TestError{NewComplianceError(entityFile, "complex fields", err.Error())}
			}

			return nil
		},
	}
}

func (r *ComprehensiveTestRunner) createFlagValidationScenario() TestScenario {
	return TestScenario{
		Name:        "Flag Validation",
		Description: "Tests CLI flag handling and validation",
		Execute: func(suite *TestSuite) error {
			return nil
		},
		Validate: func(suite *TestSuite) []*TestError {
			var errors []*TestError
			cli := NewCLIRunner(suite)

			// Test force flag
			cli.InitProject("flagtest", "", "--force")
			if !suite.FileExists(filepath.Join(suite.tempDir, "go.mod")) {
				errors = append(errors, NewFlagError("force", "project created", "not created", "force flag should create project"))
			}

			return errors
		},
	}
}

// printFinalSummary prints the final test summary
func (r *ComprehensiveTestRunner) printFinalSummary(passedScenarios int, allErrors []*TestError) {
	fmt.Printf("🎯 GOCA CLI TESTING SUMMARY\n")
	fmt.Printf("==========================\n")
	fmt.Printf("Total Scenarios: %d\n", len(r.scenarios))
	fmt.Printf("Passed: %d\n", passedScenarios)
	fmt.Printf("Failed: %d\n", len(r.scenarios)-passedScenarios)

	summary := NewErrorSummary(allErrors)
	fmt.Printf("Error Summary: %s\n", summary.String())

	if len(allErrors) > 0 {
		fmt.Printf("\n📊 ERROR BREAKDOWN:\n")
		for errorType, count := range summary.ByType {
			fmt.Printf("  %s: %d\n", errorType, count)
		}

		if len(summary.ByFile) > 0 {
			fmt.Printf("\n📁 ERRORS BY FILE:\n")
			for file, count := range summary.ByFile {
				fmt.Printf("  %s: %d\n", file, count)
			}
		}
	}

	if summary.HasCriticalErrors() {
		fmt.Printf("\n❌ TEST SUITE FAILED - Critical errors found!\n")
	} else if len(allErrors) == 0 {
		fmt.Printf("\n✅ ALL TESTS PASSED - Goca CLI is working perfectly!\n")
	} else {
		fmt.Printf("\n⚠️  Tests completed with warnings - Review non-critical issues\n")
	}
}
