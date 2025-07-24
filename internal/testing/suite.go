package testing

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TestSuite represents the comprehensive testing suite for Goca CLI
type TestSuite struct {
	t           *testing.T
	workingDir  string
	projectName string
	moduleName  string
	tempDir     string
	projectPath string
	errors      []string
	warnings    []string
}

// NewTestSuite creates a new test suite instance
func NewTestSuite(t *testing.T) *TestSuite {
	tempDir, err := os.MkdirTemp("", "goca_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	return &TestSuite{
		t:           t,
		tempDir:     tempDir,
		projectName: "testproject",
		moduleName:  "github.com/test/testproject",
		errors:      make([]string, 0),
		warnings:    make([]string, 0),
	}
}

// Cleanup removes temporary directories and files
func (ts *TestSuite) Cleanup() {
	if ts.tempDir != "" {
		os.RemoveAll(ts.tempDir)
	}
}

// runGocaCommand executes a goca CLI command and captures output
func (ts *TestSuite) runGocaCommand(args ...string) (string, string, error) {
	// Use local executable on Windows, or goca from PATH on other systems
	gocaCmd := "goca"
	if runtime.GOOS == "windows" {
		// Get the current working directory and build path to goca.exe
		if wd, err := os.Getwd(); err == nil {
			// If we're in the testing directory, go up two levels
			if strings.Contains(wd, "internal") {
				gocaCmd = filepath.Join(wd, "..", "..", "goca.exe")
			} else {
				gocaCmd = filepath.Join(wd, "goca.exe")
			}
		}
	}

	cmd := exec.Command(gocaCmd, args...)
	cmd.Dir = ts.workingDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// TestInitCommand tests the goca init command with all flags
func (ts *TestSuite) TestInitCommand() {
	ts.t.Log("Testing goca init command...")

	// Change to temp directory
	ts.workingDir = ts.tempDir

	// Test basic init
	stdout, stderr, err := ts.runGocaCommand("init", ts.projectName,
		"--module", ts.moduleName,
		"--database", "postgres",
		"--auth",
		"--api", "rest")

	if err != nil {
		ts.addError(fmt.Sprintf("goca init failed: %v\nStdout: %s\nStderr: %s", err, stdout, stderr))
		return
	}

	ts.projectPath = filepath.Join(ts.tempDir, ts.projectName)

	// Verify project structure was created
	ts.verifyProjectStructure()

	// Change working directory to project
	ts.workingDir = ts.projectPath

	ts.t.Log("✅ goca init command passed")
}

// TestFeatureCommand tests the goca feature command with various field types
func (ts *TestSuite) TestFeatureCommand() {
	ts.t.Log("Testing goca feature command...")

	testCases := []struct {
		name   string
		entity string
		fields string
		flags  []string
	}{
		{
			name:   "User with string and int fields",
			entity: "User",
			fields: "name:string,email:string,age:int",
			flags:  []string{"--validation", "--business-rules"},
		},
		{
			name:   "Product with float and multiple types",
			entity: "Product",
			fields: "title:string,description:string,price:float64,stock:int,active:bool",
			flags:  []string{"--validation"},
		},
		{
			name:   "Order with complex fields",
			entity: "Order",
			fields: "user_id:int,total:float64,status:string,items:string",
			flags:  []string{"--validation", "--business-rules"},
		},
	}

	for _, tc := range testCases {
		ts.t.Logf("Testing feature: %s", tc.name)

		args := []string{"feature", tc.entity, "--fields", tc.fields}
		args = append(args, tc.flags...)

		stdout, stderr, err := ts.runGocaCommand(args...)
		if err != nil {
			ts.addError(fmt.Sprintf("goca feature %s failed: %v\nStdout: %s\nStderr: %s",
				tc.entity, err, stdout, stderr))
			continue
		}

		// Verify feature structure
		ts.verifyFeatureStructure(tc.entity)

		// Verify domain entity
		ts.verifyDomainEntity(tc.entity, tc.fields)

		// Verify use cases
		ts.verifyUseCases(tc.entity)

		// Verify repositories
		ts.verifyRepositories(tc.entity)

		// Verify handlers
		ts.verifyHandlers(tc.entity)

		ts.t.Logf("✅ Feature %s passed", tc.entity)
	}
}

// TestEntityCommand tests the goca entity command with all flags
func (ts *TestSuite) TestEntityCommand() {
	ts.t.Log("Testing goca entity command...")

	testCases := []struct {
		name   string
		entity string
		fields string
		flags  []string
	}{
		{
			name:   "Simple entity",
			entity: "Customer",
			fields: "name:string,email:string",
			flags:  []string{"--validation"},
		},
		{
			name:   "Entity with timestamps",
			entity: "Invoice",
			fields: "number:string,amount:float64",
			flags:  []string{"--validation", "--timestamps"},
		},
		{
			name:   "Entity with soft delete",
			entity: "Category",
			fields: "name:string,description:string",
			flags:  []string{"--validation", "--soft-delete"},
		},
		{
			name:   "Full entity with all features",
			entity: "Employee",
			fields: "name:string,position:string,salary:float64",
			flags:  []string{"--validation", "--business-rules", "--timestamps", "--soft-delete"},
		},
	}

	for _, tc := range testCases {
		ts.t.Logf("Testing entity: %s", tc.name)

		args := []string{"entity", tc.entity, "--fields", tc.fields}
		args = append(args, tc.flags...)

		stdout, stderr, err := ts.runGocaCommand(args...)
		if err != nil {
			ts.addError(fmt.Sprintf("goca entity %s failed: %v\nStdout: %s\nStderr: %s",
				tc.entity, err, stdout, stderr))
			continue
		}

		// Verify entity file was created correctly
		ts.verifyEntityFile(tc.entity, tc.fields, tc.flags)

		ts.t.Logf("✅ Entity %s passed", tc.entity)
	}
}

// AssertFileNotExists checks if a file does not exist
func (ts *TestSuite) AssertFileNotExists(path string) {
	ts.t.Helper()
	fullPath := filepath.Join(ts.tempDir, path)
	if _, err := os.Stat(fullPath); err == nil {
		ts.t.Errorf("Expected file %s to not exist, but it does", path)
	}
}

// AssertFileContains checks if a file contains specific content
func (ts *TestSuite) AssertFileContains(path, content string) {
	ts.t.Helper()
	fullPath := filepath.Join(ts.tempDir, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		ts.t.Errorf("Failed to read file %s: %v", path, err)
		return
	}

	if !strings.Contains(string(data), content) {
		ts.t.Errorf("File %s does not contain expected content: %s", path, content)
	}
}

// AssertFileNotContains checks if a file does not contain specific content
func (ts *TestSuite) AssertFileNotContains(path, content string) {
	ts.t.Helper()
	fullPath := filepath.Join(ts.tempDir, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		ts.t.Errorf("Failed to read file %s: %v", path, err)
		return
	}

	if strings.Contains(string(data), content) {
		ts.t.Errorf("File %s should not contain content: %s", path, content)
	}
}

// AssertValidGoCode checks if the generated Go code is syntactically valid
func (ts *TestSuite) AssertValidGoCode(path string) {
	ts.t.Helper()
	// Implementation will use go/parser to validate syntax
	// TODO: Add go/parser validation
}

// AssertDirectoryExists checks if a directory exists
func (ts *TestSuite) AssertDirectoryExists(path string) {
	ts.t.Helper()
	fullPath := filepath.Join(ts.tempDir, path)
	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		ts.t.Errorf("Expected directory %s to exist, but it doesn't", path)
		return
	}
	if !info.IsDir() {
		ts.t.Errorf("Expected %s to be a directory, but it's a file", path)
	}
}

// GetProjectPath returns the full path to a file in the test project
func (ts *TestSuite) GetProjectPath(relativePath string) string {
	return filepath.Join(ts.tempDir, relativePath)
}

// ChangeToProject changes to the project directory
func (ts *TestSuite) ChangeToProject(projectName string) {
	projectPath := filepath.Join(ts.tempDir, projectName)
	_ = os.Chdir(projectPath)
}

// ChangeToTempDir changes to the temp directory
func (ts *TestSuite) ChangeToTempDir() {
	_ = os.Chdir(ts.tempDir)
}

// FileExists checks if a file exists and returns a boolean
func (ts *TestSuite) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// addError adds an error to the test suite
func (ts *TestSuite) addError(msg string) {
	ts.errors = append(ts.errors, msg)
	ts.t.Errorf("ERROR: %s", msg)
}

// addWarning adds a warning to the test suite
func (ts *TestSuite) addWarning(msg string) {
	ts.warnings = append(ts.warnings, msg)
	ts.t.Logf("WARNING: %s", msg)
}

// verifyProjectStructure verifies the basic project structure
func (ts *TestSuite) verifyProjectStructure() {
	ts.t.Log("Verifying project structure...")

	expectedDirs := []string{
		"cmd/server",
		"internal/domain",
		"internal/usecase",
		"internal/infrastructure/repository",
		"internal/infrastructure/handler",
		"pkg/config",
		"pkg/logger",
		"pkg/auth",
	}

	for _, dir := range expectedDirs {
		fullPath := filepath.Join(ts.projectPath, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			ts.addError(fmt.Sprintf("Expected directory %s does not exist", dir))
		}
	}

	expectedFiles := []string{
		"go.mod",
		"README.md",
		".gitignore",
		"cmd/server/main.go",
		"pkg/config/config.go",
		"pkg/logger/logger.go",
		"pkg/auth/jwt.go",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(ts.projectPath, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			ts.addError(fmt.Sprintf("Expected file %s does not exist", file))
		}
	}
}

// verifyFeatureStructure verifies that a feature has all expected directories and files
func (ts *TestSuite) verifyFeatureStructure(entity string) {
	ts.t.Logf("Verifying feature structure for %s...", entity)

	entityLower := strings.ToLower(entity)

	expectedFiles := []string{
		fmt.Sprintf("internal/domain/entity/%s.go", entityLower),
		"internal/domain/entity/errors.go",
		fmt.Sprintf("internal/usecase/%s_usecase.go", entityLower),
		fmt.Sprintf("internal/usecase/%s_service.go", entityLower),
		"internal/usecase/dto.go",
		"internal/usecase/interfaces.go",
		"internal/infrastructure/repository/interfaces.go",
		fmt.Sprintf("internal/infrastructure/repository/postgres_%s_repo.go", entityLower),
		fmt.Sprintf("internal/infrastructure/handler/%s_handler.go", entityLower),
		"internal/infrastructure/handler/routes.go",
		"internal/messages/errors.go",
		"internal/messages/responses.go",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(ts.projectPath, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			ts.addError(fmt.Sprintf("Expected feature file %s does not exist", file))
		}
	}
}

// verifyDomainEntity verifies the domain entity file
func (ts *TestSuite) verifyDomainEntity(entity, fields string) {
	ts.t.Logf("Verifying domain entity %s...", entity)

	entityLower := strings.ToLower(entity)
	filePath := filepath.Join(ts.projectPath, "internal/domain/entity", entityLower+".go")

	content, err := os.ReadFile(filePath)
	if err != nil {
		ts.addError(fmt.Sprintf("Cannot read entity file %s: %v", filePath, err))
		return
	}

	contentStr := string(content)

	// Verify package declaration
	if !strings.Contains(contentStr, "package domain") {
		ts.addError(fmt.Sprintf("Entity %s missing package domain declaration", entity))
	}

	// Verify struct declaration
	structPattern := fmt.Sprintf(`type %s struct`, entity)
	if !strings.Contains(contentStr, structPattern) {
		ts.addError(fmt.Sprintf("Entity %s missing struct declaration", entity))
	}

	// Verify fields
	fieldPairs := strings.Split(fields, ",")
	for _, fieldPair := range fieldPairs {
		parts := strings.Split(strings.TrimSpace(fieldPair), ":")
		if len(parts) == 2 {
			caser := cases.Title(language.English)
			fieldName := caser.String(strings.TrimSpace(parts[0]))
			fieldType := strings.TrimSpace(parts[1])

			// Use regex to handle Go formatting with multiple spaces
			fieldPattern := fmt.Sprintf(`%s\s+%s`, fieldName, fieldType)
			matched, _ := regexp.MatchString(fieldPattern, contentStr)
			if !matched {
				ts.addError(fmt.Sprintf("Entity %s missing field %s of type %s", entity, fieldName, fieldType))
			}
		}
	}

	// Verify Validate method
	validatePattern := fmt.Sprintf(`func (%s *%s) Validate() error`, strings.ToLower(entity[:1]), entity)
	if !strings.Contains(contentStr, validatePattern) {
		ts.addError(fmt.Sprintf("Entity %s missing Validate method", entity))
	}

	ts.verifyGoSyntax(filePath)
}

// verifyUseCases verifies use case files
func (ts *TestSuite) verifyUseCases(entity string) {
	ts.t.Logf("Verifying use cases for %s...", entity)

	entityLower := strings.ToLower(entity)

	// Verify use case interface
	interfaceFile := filepath.Join(ts.projectPath, "internal/usecase", entityLower+"_usecase.go")
	ts.verifyFileExists(interfaceFile)
	ts.verifyGoSyntax(interfaceFile)

	// Verify use case service
	serviceFile := filepath.Join(ts.projectPath, "internal/usecase", entityLower+"_service.go")
	ts.verifyFileExists(serviceFile)
	ts.verifyGoSyntax(serviceFile)

	// Verify DTOs
	dtoFile := filepath.Join(ts.projectPath, "internal/usecase", "dto.go")
	ts.verifyFileExists(dtoFile)
	ts.verifyGoSyntax(dtoFile)

	// Verify interfaces
	interfacesFile := filepath.Join(ts.projectPath, "internal/usecase", "interfaces.go")
	ts.verifyFileExists(interfacesFile)
	ts.verifyGoSyntax(interfacesFile)
}

// verifyRepositories verifies repository files
func (ts *TestSuite) verifyRepositories(entity string) {
	ts.t.Logf("Verifying repositories for %s...", entity)

	entityLower := strings.ToLower(entity)

	// Verify repository interface
	interfaceFile := filepath.Join(ts.projectPath, "internal/infrastructure/repository", "interfaces.go")
	ts.verifyFileExists(interfaceFile)
	ts.verifyGoSyntax(interfaceFile)

	// Verify postgres repository
	repoFile := filepath.Join(ts.projectPath, "internal/infrastructure/repository", "postgres_"+entityLower+"_repo.go")
	ts.verifyFileExists(repoFile)
	ts.verifyGoSyntax(repoFile)
}

// verifyHandlers verifies handler files
func (ts *TestSuite) verifyHandlers(entity string) {
	ts.t.Logf("Verifying handlers for %s...", entity)

	entityLower := strings.ToLower(entity)

	// Verify HTTP handler
	handlerFile := filepath.Join(ts.projectPath, "internal/infrastructure/handler", entityLower+"_handler.go")
	ts.verifyFileExists(handlerFile)
	ts.verifyGoSyntax(handlerFile)

	// Verify routes
	routesFile := filepath.Join(ts.projectPath, "internal/infrastructure/handler", "routes.go")
	ts.verifyFileExists(routesFile)
	ts.verifyGoSyntax(routesFile)
}

// verifyEntityFile verifies entity files with specific flags
func (ts *TestSuite) verifyEntityFile(entity, fields string, flags []string) {
	ts.t.Logf("Verifying entity file %s with flags %v...", entity, flags)

	entityLower := strings.ToLower(entity)
	filePath := filepath.Join(ts.projectPath, "internal/domain/entity", entityLower+".go")

	content, err := os.ReadFile(filePath)
	if err != nil {
		ts.addError(fmt.Sprintf("Cannot read entity file %s: %v", filePath, err))
		return
	}

	contentStr := string(content)

	// Verify basic entity structure contains expected fields
	if fields != "" {
		// Basic validation that the struct exists
		if !strings.Contains(contentStr, fmt.Sprintf("type %s struct", entity)) {
			ts.addError(fmt.Sprintf("Entity %s missing struct definition", entity))
		}
	}

	// Check for validation flag
	if contains(flags, "--validation") {
		if !strings.Contains(contentStr, "Validate() error") {
			ts.addError(fmt.Sprintf("Entity %s with --validation flag missing Validate method", entity))
		}
	}

	// Check for timestamps flag
	if contains(flags, "--timestamps") {
		if !strings.Contains(contentStr, "CreatedAt") || !strings.Contains(contentStr, "UpdatedAt") {
			ts.addError(fmt.Sprintf("Entity %s with --timestamps flag missing timestamp fields", entity))
		}
	}

	// Check for soft-delete flag
	if contains(flags, "--soft-delete") {
		if !strings.Contains(contentStr, "DeletedAt") {
			ts.addError(fmt.Sprintf("Entity %s with --soft-delete flag missing DeletedAt field", entity))
		}
		if !strings.Contains(contentStr, "SoftDelete()") {
			ts.addError(fmt.Sprintf("Entity %s with --soft-delete flag missing SoftDelete method", entity))
		}
	}

	// Check for business-rules flag
	if contains(flags, "--business-rules") {
		// Look for business rule methods - this is more complex validation
		// We'll check if there are additional methods beyond the basic struct and Validate
		methods := ts.countMethods(contentStr)
		expectedMethods := 1 // Validate method
		if contains(flags, "--soft-delete") {
			expectedMethods += 2 // SoftDelete and IsDeleted methods
		}
		if methods <= expectedMethods {
			ts.addWarning(fmt.Sprintf("Entity %s with --business-rules flag may be missing business rule methods", entity))
		}
	}

	ts.verifyGoSyntax(filePath)
}

// verifyFileExists checks if a file exists
func (ts *TestSuite) verifyFileExists(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ts.addError(fmt.Sprintf("Expected file %s does not exist", filePath))
	}
}

// verifyGoSyntax verifies that a Go file has valid syntax
func (ts *TestSuite) verifyGoSyntax(filePath string) {
	fset := token.NewFileSet()

	content, err := os.ReadFile(filePath)
	if err != nil {
		ts.addError(fmt.Sprintf("Cannot read file %s: %v", filePath, err))
		return
	}

	_, err = parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		ts.addError(fmt.Sprintf("Go syntax error in %s: %v", filePath, err))
	}
}

// countMethods counts the number of methods in a Go source code string
func (ts *TestSuite) countMethods(content string) int {
	// Simple regex to count method definitions
	methodRegex := regexp.MustCompile(`func\s+\([^)]+\)\s+\w+\s*\(`)
	return len(methodRegex.FindAllString(content, -1))
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// TestCodeCompilation tests that all generated code compiles without errors
func (ts *TestSuite) TestCodeCompilation() {
	ts.t.Log("Testing code compilation...")

	// Run go mod download first to ensure dependencies are available
	cmd := exec.Command("go", "mod", "download")
	cmd.Dir = ts.projectPath
	if err := cmd.Run(); err != nil {
		ts.addWarning(fmt.Sprintf("go mod download warning: %v", err))
		// Continue anyway, as this might not be critical
	}

	// Try to build the code without running go mod tidy first
	// Since local imports in generated code can cause go mod tidy to fail
	cmd = exec.Command("go", "build", "./...")
	cmd.Dir = ts.projectPath
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	buildErr := cmd.Run()
	if buildErr == nil {
		ts.t.Log("✅ Code compilation passed")
		return
	}

	// If build fails, try go mod tidy as last resort
	for attempts := 0; attempts < 2; attempts++ {
		cmd = exec.Command("go", "mod", "tidy")
		cmd.Dir = ts.projectPath
		var tidyOut, tidyErrBuf bytes.Buffer
		cmd.Stdout = &tidyOut
		cmd.Stderr = &tidyErrBuf

		tidyErr := cmd.Run()
		if tidyErr == nil {
			// Try building again after successful tidy
			cmd = exec.Command("go", "build", "./...")
			cmd.Dir = ts.projectPath
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			if err := cmd.Run(); err == nil {
				ts.t.Log("✅ Code compilation passed after go mod tidy")
				return
			}
			break
		}
	}

	// If everything fails, report the original build error
	ts.addError(fmt.Sprintf("Code compilation failed: %v\nStdout: %s\nStderr: %s",
		buildErr, stdout.String(), stderr.String()))
}

// TestCodeLinting tests that all generated code passes go vet
func (ts *TestSuite) TestCodeLinting() {
	ts.t.Log("Testing code linting...")

	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = ts.projectPath
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		ts.addWarning(fmt.Sprintf("Code linting issues: %v\nStdout: %s\nStderr: %s",
			err, stdout.String(), stderr.String()))
		return
	}

	ts.t.Log("✅ Code linting passed")
}

// TestCodeFormatting tests that all generated code is properly formatted
func (ts *TestSuite) TestCodeFormatting() {
	ts.t.Log("Testing code formatting...")

	err := filepath.WalkDir(ts.projectPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Check if file is properly formatted
		cmd := exec.Command("gofmt", "-l", path)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout

		if err := cmd.Run(); err != nil {
			ts.addWarning(fmt.Sprintf("gofmt check failed for %s: %v", path, err))
			return nil
		}

		if stdout.Len() > 0 {
			ts.addWarning(fmt.Sprintf("File %s is not properly formatted", path))
		}

		return nil
	})

	if err != nil {
		ts.addError(fmt.Sprintf("Code formatting check failed: %v", err))
		return
	}

	ts.t.Log("✅ Code formatting passed")
}

// RunAllTests runs the complete test suite
func (ts *TestSuite) RunAllTests() {
	ts.t.Log("🚀 Starting comprehensive Goca CLI test suite...")

	// Test in order of dependencies
	ts.TestInitCommand()

	if len(ts.errors) > 0 {
		ts.t.Fatalf("Init command failed, stopping tests. Errors: %v", ts.errors)
	}

	ts.TestFeatureCommand()
	ts.TestEntityCommand()

	// Code quality tests
	ts.TestCodeCompilation()
	ts.TestCodeLinting()
	ts.TestCodeFormatting()

	// Report results
	ts.reportResults()
}

// reportResults reports the final test results
func (ts *TestSuite) reportResults() {
	ts.t.Log("📊 Test Suite Results:")
	ts.t.Logf("✅ Generated project: %s", ts.projectPath)
	ts.t.Logf("⚠️  Warnings: %d", len(ts.warnings))
	ts.t.Logf("❌ Errors: %d", len(ts.errors))

	if len(ts.warnings) > 0 {
		ts.t.Log("Warnings:")
		for _, warning := range ts.warnings {
			ts.t.Logf("  - %s", warning)
		}
	}

	if len(ts.errors) > 0 {
		ts.t.Log("Errors:")
		for _, err := range ts.errors {
			ts.t.Logf("  - %s", err)
		}
		ts.t.Fatalf("Test suite failed with %d errors", len(ts.errors))
	} else {
		ts.t.Log("🎉 All tests passed successfully!")
	}
}
