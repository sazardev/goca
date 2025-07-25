package testing

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ArchitectureValidator validates Clean Architecture patterns in generated code
type ArchitectureValidator struct {
	suite *TestSuite
}

// NewArchitectureValidator creates a new architecture validator
func NewArchitectureValidator(suite *TestSuite) *ArchitectureValidator {
	return &ArchitectureValidator{suite: suite}
}

// ValidateProjectStructure checks if project follows Clean Architecture structure
func (v *ArchitectureValidator) ValidateProjectStructure(projectDir string) []*TestError {
	var errors []*TestError

	expectedDirs := []string{
		"internal/domain",
		"internal/usecase",
		"internal/repository",
		"internal/handler",
		"internal/di",
		"internal/constants",
		"internal/messages",
		"pkg/messages",
		"pkg/interfaces",
	}

	for _, dir := range expectedDirs {
		fullPath := filepath.Join(projectDir, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			errors = append(errors, NewStructureError(dir, "directory not found"))
		}
	}

	return errors
}

// ValidateLayerSeparation ensures layers don't have forbidden dependencies
func (v *ArchitectureValidator) ValidateLayerSeparation(projectDir string) []*TestError {
	var errors []*TestError

	// Domain layer should not import infrastructure
	domainFiles := v.findGoFiles(filepath.Join(projectDir, "internal/domain"))
	for _, file := range domainFiles {
		if v.hasInfrastructureImports(file, projectDir) {
			errors = append(errors, NewDependencyError(file, "domain", "infrastructure",
				"Domain layer should not depend on infrastructure"))
		}
	}

	// UseCase layer should not import infrastructure except interfaces
	usecaseFiles := v.findGoFiles(filepath.Join(projectDir, "internal/usecase"))
	for _, file := range usecaseFiles {
		if v.hasForbiddenInfrastructureImports(file, projectDir) {
			errors = append(errors, NewDependencyError(file, "usecase", "infrastructure",
				"UseCase layer should only depend on domain and interfaces"))
		}
	}

	return errors
}

// ValidateEntityCompliance checks if entity follows domain rules
func (v *ArchitectureValidator) ValidateEntityCompliance(entityFile string) []*TestError {
	var errors []*TestError

	// Check if entity is in correct location
	if !strings.Contains(entityFile, "internal/domain") {
		errors = append(errors, NewLocationError(entityFile, "internal/domain",
			"Entity should be in domain layer"))
	}

	// Validate entity structure using code validator
	codeValidator := NewCodeValidator(v.suite)

	// Check package declaration
	if err := codeValidator.ValidatePackageDeclaration(entityFile, "domain"); err != nil {
		errors = append(errors, NewComplianceError(entityFile, "package", err.Error()))
	}

	// Check naming conventions
	if namingErrors := codeValidator.ValidateNamingConventions(entityFile); len(namingErrors) > 0 {
		for _, ne := range namingErrors {
			errors = append(errors, NewComplianceError(entityFile, "naming", ne.Error()))
		}
	}

	return errors
}

// ValidateUseCaseCompliance checks if usecase follows clean architecture rules
func (v *ArchitectureValidator) ValidateUseCaseCompliance(usecaseFile string) []*TestError {
	var errors []*TestError

	// Check if usecase is in correct location
	if !strings.Contains(usecaseFile, "internal/usecase") {
		errors = append(errors, NewLocationError(usecaseFile, "internal/usecase",
			"UseCase should be in usecase layer"))
	}

	// Validate usecase structure
	codeValidator := NewCodeValidator(v.suite)

	// Check package declaration
	if err := codeValidator.ValidatePackageDeclaration(usecaseFile, "usecase"); err != nil {
		errors = append(errors, NewComplianceError(usecaseFile, "package", err.Error()))
	}

	// Check for required imports (should import domain and interfaces)
	content, err := os.ReadFile(usecaseFile)
	if err == nil {
		contentStr := string(content)
		if !strings.Contains(contentStr, "/internal/domain") {
			errors = append(errors, NewComplianceError(usecaseFile, "imports",
				"UseCase should import domain entities"))
		}
		if !strings.Contains(contentStr, "/pkg/interfaces") {
			errors = append(errors, NewComplianceError(usecaseFile, "imports",
				"UseCase should import repository interfaces"))
		}
	}

	return errors
}

// ValidateRepositoryCompliance checks if repository follows clean architecture rules
func (v *ArchitectureValidator) ValidateRepositoryCompliance(repoFile string) []*TestError {
	var errors []*TestError

	// Check if repository is in correct location
	if !strings.Contains(repoFile, "internal/repository") {
		errors = append(errors, NewLocationError(repoFile, "internal/repository",
			"Repository should be in repository layer"))
	}

	// Validate repository structure
	codeValidator := NewCodeValidator(v.suite)

	// Check package declaration
	if err := codeValidator.ValidatePackageDeclaration(repoFile, "repository"); err != nil {
		errors = append(errors, NewComplianceError(repoFile, "package", err.Error()))
	}

	// Check for interface implementation
	content, err := os.ReadFile(repoFile)
	if err == nil {
		contentStr := string(content)
		if !strings.Contains(contentStr, "/pkg/interfaces") {
			errors = append(errors, NewComplianceError(repoFile, "imports",
				"Repository should import and implement interfaces"))
		}
		if !strings.Contains(contentStr, "/internal/domain") {
			errors = append(errors, NewComplianceError(repoFile, "imports",
				"Repository should import domain entities"))
		}
	}

	return errors
}

// ValidateHandlerCompliance checks if handler follows clean architecture rules
func (v *ArchitectureValidator) ValidateHandlerCompliance(handlerFile string) []*TestError {
	var errors []*TestError

	// Check if handler is in correct location
	if !strings.Contains(handlerFile, "internal/handler") {
		errors = append(errors, NewLocationError(handlerFile, "internal/handler",
			"Handler should be in handler layer"))
	}

	// Validate handler structure
	codeValidator := NewCodeValidator(v.suite)

	// Check package declaration
	if err := codeValidator.ValidatePackageDeclaration(handlerFile, "handler"); err != nil {
		errors = append(errors, NewComplianceError(handlerFile, "package", err.Error()))
	}

	// Check for proper dependencies
	content, err := os.ReadFile(handlerFile)
	if err == nil {
		contentStr := string(content)
		if !strings.Contains(contentStr, "/internal/usecase") {
			errors = append(errors, NewComplianceError(handlerFile, "imports",
				"Handler should import and use usecases"))
		}
		if !strings.Contains(contentStr, "/pkg/messages") {
			errors = append(errors, NewComplianceError(handlerFile, "imports",
				"Handler should import message types"))
		}
	}

	return errors
}

// ValidateDependencyInjection checks if DI configuration is correct
func (v *ArchitectureValidator) ValidateDependencyInjection(diFile string) []*TestError {
	var errors []*TestError

	// Check if DI is in correct location
	if !strings.Contains(diFile, "internal/infrastructure/di") {
		errors = append(errors, NewLocationError(diFile, "internal/infrastructure/di",
			"DI should be in infrastructure layer"))
	}

	// Check if DI wires all components correctly
	content, err := os.ReadFile(diFile)
	if err == nil {
		contentStr := string(content)

		// Should import all layers
		requiredImports := []string{
			"/internal/usecase",
			"/internal/infrastructure/repository",
			"/internal/infrastructure/handler",
		}

		for _, imp := range requiredImports {
			if !strings.Contains(contentStr, imp) {
				errors = append(errors, NewComplianceError(diFile, "imports",
					"DI should import "+imp))
			}
		}
	}

	return errors
}

// Helper methods

func (v *ArchitectureValidator) findGoFiles(dir string) []string {
	var files []string

	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})

	return files
}

func (v *ArchitectureValidator) hasInfrastructureImports(file, _ string) bool {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return false
	}

	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, "\"")
		if strings.Contains(importPath, "/internal/infrastructure") {
			return true
		}
	}
	return false
}

func (v *ArchitectureValidator) hasForbiddenInfrastructureImports(file, _ string) bool {
	content, err := os.ReadFile(file)
	if err != nil {
		return false
	}

	contentStr := string(content)
	// UseCase can import interfaces but not concrete infrastructure
	return strings.Contains(contentStr, "/internal/infrastructure/repository") ||
		strings.Contains(contentStr, "/internal/infrastructure/handler")
}
