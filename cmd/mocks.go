package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	mockAll        bool
	mockRepository bool
	mockUseCase    bool
	mockHandler    bool
)

// mocksCmd represents the mocks command.
var mocksCmd = &cobra.Command{
	Use:   "mocks [entity]",
	Short: "Generate mock implementations for interfaces",
	Long: `Generate mock implementations for repository, use case, and handler interfaces
using testify/mock. These mocks are essential for unit testing and test-driven development.

Mock files are generated in internal/mocks/ directory with the following naming:
- mock_{entity}_repository.go
- mock_{entity}_usecase.go
- mock_{entity}_handler.go

Examples:
  goca mocks User
  goca mocks Product --repository
  goca mocks Order --all
  goca mocks Customer --usecase --repository`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entityName := args[0]

		// Validate entity name
		validator := NewCommandValidator()
		if err := validator.fieldValidator.ValidateEntityName(entityName); err != nil {
			validator.errorHandler.HandleError(err, "entity validation")
			return
		}

		// Check if we're in a Goca project
		if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
			ui.Error("Not in a Go module directory. Run 'go mod init' first.")
			return
		}

		// If no flags specified, generate all mocks
		if !mockRepository && !mockUseCase && !mockHandler && !mockAll {
			mockAll = true
		}

		// Initialize safety manager
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		backup, _ := cmd.Flags().GetBool("backup")
		sm := NewSafetyManager(dryRun, force, backup)

		if dryRun {
			ui.DryRun("Previewing changes without creating files")
		}

		// Generate mocks
		if err := generateMocks(entityName, mockAll, mockRepository, mockUseCase, mockHandler, sm); err != nil {
			ui.Error(fmt.Sprintf("Error generating mocks: %v", err))
			return
		}

		if dryRun {
			sm.PrintSummary()
			return
		}

		ui.Success(fmt.Sprintf("Mocks generated successfully for '%s'", entityName))
		ui.Dim("See internal/mocks/examples/ for unit test examples")
	},
}

func init() {
	rootCmd.AddCommand(mocksCmd)

	mocksCmd.Flags().BoolVar(&mockAll, "all", false, "Generate all mocks (repository, usecase, handler)")
	mocksCmd.Flags().BoolVar(&mockRepository, "repository", false, "Generate repository mock only")
	mocksCmd.Flags().BoolVar(&mockUseCase, "usecase", false, "Generate use case mock only")
	mocksCmd.Flags().BoolVar(&mockHandler, "handler", false, "Generate handler mock only")
	mocksCmd.Flags().Bool("dry-run", false, "Preview changes without creating files")
	mocksCmd.Flags().Bool("force", false, "Overwrite existing files without asking")
	mocksCmd.Flags().Bool("backup", false, "Backup existing files before overwriting")
}

// generateMocks generates mock files based on flags.
func generateMocks(entityName string, all, repository, usecase, handler bool, sm ...*SafetyManager) error {
	// Create mocks directory
	mocksDir := filepath.Join("internal", "mocks")
	if err := os.MkdirAll(mocksDir, 0o755); err != nil {
		return err
	}

	importPath := getImportPath(getModuleName())

	// Recover the entity's fields so the generated mocks include the same
	// per-field finder methods (FindBy<Field>) that the real repository
	// interface declares.
	var fields []Field
	if fs := readEntityFieldsString(entityName); fs != "" {
		fields = parseFields(fs)
	}

	// Generate repository mock
	if all || repository {
		mockFile := filepath.Join(mocksDir, fmt.Sprintf("mock_%s_repository.go", strings.ToLower(entityName)))
		content := fixGeneratedModulePath(generateRepositoryMock(entityName, fields), importPath)
		if err := writeFile(mockFile, content, sm...); err != nil {
			return err
		}
	}

	// Generate use case mock
	if all || usecase {
		mockFile := filepath.Join(mocksDir, fmt.Sprintf("mock_%s_usecase.go", strings.ToLower(entityName)))
		content := fixGeneratedModulePath(generateUseCaseMock(entityName), importPath)
		if err := writeFile(mockFile, content, sm...); err != nil {
			return err
		}
	}

	// Generate handler mock
	if all || handler {
		mockFile := filepath.Join(mocksDir, fmt.Sprintf("mock_%s_handler.go", strings.ToLower(entityName)))
		content := fixGeneratedModulePath(generateHandlerMock(entityName), importPath)
		if err := writeFile(mockFile, content, sm...); err != nil {
			return err
		}
	}

	// Generate usage examples
	if all || repository || usecase || handler {
		examplesDir := filepath.Join(mocksDir, "examples")
		if err := os.MkdirAll(examplesDir, 0o755); err != nil {
			return err
		}

		exampleFile := filepath.Join(examplesDir, fmt.Sprintf("%s_mock_examples_test.go", strings.ToLower(entityName)))
		exampleContent := fixGeneratedModulePath(generateMockUsageExamples(entityName), importPath)
		if err := writeFile(exampleFile, exampleContent, sm...); err != nil {
			return err
		}
	}

	return nil
}

// generateRepositoryMock generates a mock that satisfies repository.<Entity>Repository.
// The generated finder methods (FindBy<Field>) mirror those produced for the real
// repository interface by generateSearchMethods, so the mock implements the
// interface exactly.
func generateRepositoryMock(entityName string, fields []Field) string {
	lowerEntity := strings.ToLower(entityName)

	var b strings.Builder
	b.WriteString("package mocks\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"github.com/stretchr/testify/mock\"\n")
	b.WriteString("\t\"github.com/sazardev/goca/internal/domain\"\n")
	b.WriteString(")\n\n")

	fmt.Fprintf(&b, "// Mock%sRepository is a mock implementation of repository.%sRepository\n", entityName, entityName)
	fmt.Fprintf(&b, "type Mock%sRepository struct {\n\tmock.Mock\n}\n\n", entityName)

	// Save
	fmt.Fprintf(&b, "// Save mocks the Save method\n")
	fmt.Fprintf(&b, "func (m *Mock%sRepository) Save(%s *domain.%s) error {\n", entityName, lowerEntity, entityName)
	fmt.Fprintf(&b, "\targs := m.Called(%s)\n\treturn args.Error(0)\n}\n\n", lowerEntity)

	// FindByID
	fmt.Fprintf(&b, "// FindByID mocks the FindByID method\n")
	fmt.Fprintf(&b, "func (m *Mock%sRepository) FindByID(id int) (*domain.%s, error) {\n", entityName, entityName)
	fmt.Fprintf(&b, "\targs := m.Called(id)\n\tif args.Get(0) == nil {\n\t\treturn nil, args.Error(1)\n\t}\n")
	fmt.Fprintf(&b, "\treturn args.Get(0).(*domain.%s), args.Error(1)\n}\n\n", entityName)

	// Per-field finders, matching generateSearchMethods.
	for _, method := range generateSearchMethods(fields, entityName) {
		paramName := strings.ToLower(method.FieldName)
		fmt.Fprintf(&b, "// %s mocks the %s method\n", method.MethodName, method.MethodName)
		fmt.Fprintf(&b, "func (m *Mock%sRepository) %s(%s %s) (*domain.%s, error) {\n",
			entityName, method.MethodName, paramName, method.FieldType, entityName)
		fmt.Fprintf(&b, "\targs := m.Called(%s)\n\tif args.Get(0) == nil {\n\t\treturn nil, args.Error(1)\n\t}\n", paramName)
		fmt.Fprintf(&b, "\treturn args.Get(0).(*domain.%s), args.Error(1)\n}\n\n", entityName)
	}

	// Update
	fmt.Fprintf(&b, "// Update mocks the Update method\n")
	fmt.Fprintf(&b, "func (m *Mock%sRepository) Update(%s *domain.%s) error {\n", entityName, lowerEntity, entityName)
	fmt.Fprintf(&b, "\targs := m.Called(%s)\n\treturn args.Error(0)\n}\n\n", lowerEntity)

	// Delete
	fmt.Fprintf(&b, "// Delete mocks the Delete method\n")
	fmt.Fprintf(&b, "func (m *Mock%sRepository) Delete(id int) error {\n", entityName)
	fmt.Fprintf(&b, "\targs := m.Called(id)\n\treturn args.Error(0)\n}\n\n")

	// FindAll
	fmt.Fprintf(&b, "// FindAll mocks the FindAll method\n")
	fmt.Fprintf(&b, "func (m *Mock%sRepository) FindAll() ([]domain.%s, error) {\n", entityName, entityName)
	fmt.Fprintf(&b, "\targs := m.Called()\n\tif args.Get(0) == nil {\n\t\treturn nil, args.Error(1)\n\t}\n")
	fmt.Fprintf(&b, "\treturn args.Get(0).([]domain.%s), args.Error(1)\n}\n\n", entityName)

	fmt.Fprintf(&b, "// NewMock%sRepository creates a new mock repository\n", entityName)
	fmt.Fprintf(&b, "func NewMock%sRepository() *Mock%sRepository {\n\treturn &Mock%sRepository{}\n}\n",
		entityName, entityName, entityName)

	return b.String()
}

// generateUseCaseMock generates a mock that satisfies usecase.<Entity>UseCase
// exactly: Create<Entity>, Get<Entity>, Update<Entity>, Delete<Entity> and
// List<Entity>s, with the same signatures the real interface declares.
func generateUseCaseMock(entityName string) string {
	var b strings.Builder
	b.WriteString("package mocks\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"github.com/stretchr/testify/mock\"\n")
	b.WriteString("\t\"github.com/sazardev/goca/internal/domain\"\n")
	b.WriteString("\t\"github.com/sazardev/goca/internal/usecase\"\n")
	b.WriteString(")\n\n")

	fmt.Fprintf(&b, "// Mock%sUseCase is a mock implementation of usecase.%sUseCase\n", entityName, entityName)
	fmt.Fprintf(&b, "type Mock%sUseCase struct {\n\tmock.Mock\n}\n\n", entityName)

	// Create<Entity>(input Create<Entity>Input) (Create<Entity>Output, error)
	fmt.Fprintf(&b, "// Create%s mocks the Create%s method\n", entityName, entityName)
	fmt.Fprintf(&b, "func (m *Mock%sUseCase) Create%s(input usecase.Create%sInput) (usecase.Create%sOutput, error) {\n",
		entityName, entityName, entityName, entityName)
	fmt.Fprintf(&b, "\targs := m.Called(input)\n")
	fmt.Fprintf(&b, "\treturn args.Get(0).(usecase.Create%sOutput), args.Error(1)\n}\n\n", entityName)

	// Get<Entity>(id int) (*domain.<Entity>, error)
	fmt.Fprintf(&b, "// Get%s mocks the Get%s method\n", entityName, entityName)
	fmt.Fprintf(&b, "func (m *Mock%sUseCase) Get%s(id int) (*domain.%s, error) {\n", entityName, entityName, entityName)
	fmt.Fprintf(&b, "\targs := m.Called(id)\n\tif args.Get(0) == nil {\n\t\treturn nil, args.Error(1)\n\t}\n")
	fmt.Fprintf(&b, "\treturn args.Get(0).(*domain.%s), args.Error(1)\n}\n\n", entityName)

	// Update<Entity>(id int, input Update<Entity>Input) error
	fmt.Fprintf(&b, "// Update%s mocks the Update%s method\n", entityName, entityName)
	fmt.Fprintf(&b, "func (m *Mock%sUseCase) Update%s(id int, input usecase.Update%sInput) error {\n",
		entityName, entityName, entityName)
	fmt.Fprintf(&b, "\targs := m.Called(id, input)\n\treturn args.Error(0)\n}\n\n")

	// Delete<Entity>(id int) error
	fmt.Fprintf(&b, "// Delete%s mocks the Delete%s method\n", entityName, entityName)
	fmt.Fprintf(&b, "func (m *Mock%sUseCase) Delete%s(id int) error {\n", entityName, entityName)
	fmt.Fprintf(&b, "\targs := m.Called(id)\n\treturn args.Error(0)\n}\n\n")

	// List<Entity>s() (List<Entity>Output, error)
	fmt.Fprintf(&b, "// List%ss mocks the List%ss method\n", entityName, entityName)
	fmt.Fprintf(&b, "func (m *Mock%sUseCase) List%ss() (usecase.List%sOutput, error) {\n",
		entityName, entityName, entityName)
	fmt.Fprintf(&b, "\targs := m.Called()\n")
	fmt.Fprintf(&b, "\treturn args.Get(0).(usecase.List%sOutput), args.Error(1)\n}\n\n", entityName)

	fmt.Fprintf(&b, "// NewMock%sUseCase creates a new mock use case\n", entityName)
	fmt.Fprintf(&b, "func NewMock%sUseCase() *Mock%sUseCase {\n\treturn &Mock%sUseCase{}\n}\n",
		entityName, entityName, entityName)

	return b.String()
}

// generateHandlerMock generates a mock for HTTP handler interface.
func generateHandlerMock(entityName string) string {
	return fmt.Sprintf(
		`package mocks

import (
	"net/http"
	"github.com/stretchr/testify/mock"
)

// Mock%sHandler is a mock implementation of HTTP handler
type Mock%sHandler struct {
	mock.Mock
}

// Create%s mocks the Create%s HTTP handler method
func (m *Mock%sHandler) Create%s(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// Get%s mocks the Get%s HTTP handler method
func (m *Mock%sHandler) Get%s(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// Update%s mocks the Update%s HTTP handler method
func (m *Mock%sHandler) Update%s(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// Delete%s mocks the Delete%s HTTP handler method
func (m *Mock%sHandler) Delete%s(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// List%ss mocks the List%ss HTTP handler method
func (m *Mock%sHandler) List%ss(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// NewMock%sHandler creates a new mock handler
func NewMock%sHandler() *Mock%sHandler {
	return &Mock%sHandler{}
}
`,
		entityName, entityName,
		entityName, entityName, entityName, entityName,
		entityName, entityName, entityName, entityName,
		entityName, entityName, entityName, entityName,
		entityName, entityName, entityName, entityName,
		entityName, entityName, entityName, entityName,
		entityName, entityName, entityName, entityName,
	)
} // generateMockUsageExamples generates example test files showing how to use mocks

func generateMockUsageExamples(entityName string) string {
	lowerEntity := strings.ToLower(entityName)

	return fmt.Sprintf(`package examples

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/mocks"
	"github.com/sazardev/goca/internal/repository"
	"github.com/sazardev/goca/internal/usecase"
)

// Compile-time assertions that the generated mocks satisfy the real interfaces.
var (
	_ repository.%[1]sRepository = (*mocks.Mock%[1]sRepository)(nil)
	_ usecase.%[1]sUseCase       = (*mocks.Mock%[1]sUseCase)(nil)
)

// Example: stubbing repository methods and verifying expectations.
func TestMock%[1]sRepository_Usage(t *testing.T) {
	mockRepo := mocks.NewMock%[1]sRepository()

	expected := &domain.%[1]s{ID: 1}
	mockRepo.On("FindByID", 1).Return(expected, nil)
	mockRepo.On("Save", mock.AnythingOfType("*domain.%[1]s")).Return(nil)

	got, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	err = mockRepo.Save(&domain.%[1]s{})
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// Example: stubbing an error return.
func TestMock%[1]sRepository_NotFound(t *testing.T) {
	mockRepo := mocks.NewMock%[1]sRepository()

	expectedErr := errors.New("%[2]s not found")
	mockRepo.On("FindByID", 999).Return(nil, expectedErr)

	got, err := mockRepo.FindByID(999)
	assert.Nil(t, got)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}

// Example: stubbing use-case methods.
func TestMock%[1]sUseCase_Usage(t *testing.T) {
	mockUC := mocks.NewMock%[1]sUseCase()

	mockUC.On("Get%[1]s", 1).Return(&domain.%[1]s{ID: 1}, nil)
	mockUC.On("Delete%[1]s", 1).Return(nil)

	got, err := mockUC.Get%[1]s(1)
	assert.NoError(t, err)
	assert.NotNil(t, got)

	err = mockUC.Delete%[1]s(1)
	assert.NoError(t, err)

	mockUC.AssertExpectations(t)
}
`, entityName, lowerEntity)
}
