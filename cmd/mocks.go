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

// mocksCmd represents the mocks command
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
			fmt.Println("❌ Error: Not in a Go module directory. Run 'go mod init' first.")
			return
		}

		// If no flags specified, generate all mocks
		if !mockRepository && !mockUseCase && !mockHandler && !mockAll {
			mockAll = true
		}

		// Generate mocks
		if err := generateMocks(entityName, mockAll, mockRepository, mockUseCase, mockHandler); err != nil {
			fmt.Printf("❌ Error generating mocks: %v\n", err)
			return
		}

		fmt.Printf("\n✅ Mocks generated successfully for '%s'\n", entityName)
		fmt.Println("\nGenerated files:")
		if mockAll || mockRepository {
			fmt.Printf("   - internal/mocks/mock_%s_repository.go\n", strings.ToLower(entityName))
		}
		if mockAll || mockUseCase {
			fmt.Printf("   - internal/mocks/mock_%s_usecase.go\n", strings.ToLower(entityName))
		}
		if mockAll || mockHandler {
			fmt.Printf("   - internal/mocks/mock_%s_handler.go\n", strings.ToLower(entityName))
		}
		fmt.Println("\nUsage example:")
		fmt.Printf("   See internal/mocks/examples/ for unit test examples\n")
	},
}

func init() {
	rootCmd.AddCommand(mocksCmd)

	mocksCmd.Flags().BoolVar(&mockAll, "all", false, "Generate all mocks (repository, usecase, handler)")
	mocksCmd.Flags().BoolVar(&mockRepository, "repository", false, "Generate repository mock only")
	mocksCmd.Flags().BoolVar(&mockUseCase, "usecase", false, "Generate use case mock only")
	mocksCmd.Flags().BoolVar(&mockHandler, "handler", false, "Generate handler mock only")
}

// generateMocks generates mock files based on flags
func generateMocks(entityName string, all, repository, usecase, handler bool) error {
	// Create mocks directory
	mocksDir := filepath.Join("internal", "mocks")
	if err := os.MkdirAll(mocksDir, 0755); err != nil {
		return err
	}

	// Generate repository mock
	if all || repository {
		mockFile := filepath.Join(mocksDir, fmt.Sprintf("mock_%s_repository.go", strings.ToLower(entityName)))
		content := generateRepositoryMock(entityName)
		if err := os.WriteFile(mockFile, []byte(content), 0644); err != nil {
			return err
		}
	}

	// Generate use case mock
	if all || usecase {
		mockFile := filepath.Join(mocksDir, fmt.Sprintf("mock_%s_usecase.go", strings.ToLower(entityName)))
		content := generateUseCaseMock(entityName)
		if err := os.WriteFile(mockFile, []byte(content), 0644); err != nil {
			return err
		}
	}

	// Generate handler mock
	if all || handler {
		mockFile := filepath.Join(mocksDir, fmt.Sprintf("mock_%s_handler.go", strings.ToLower(entityName)))
		content := generateHandlerMock(entityName)
		if err := os.WriteFile(mockFile, []byte(content), 0644); err != nil {
			return err
		}
	}

	// Generate usage examples
	if all || repository || usecase || handler {
		examplesDir := filepath.Join(mocksDir, "examples")
		if err := os.MkdirAll(examplesDir, 0755); err != nil {
			return err
		}

		exampleFile := filepath.Join(examplesDir, fmt.Sprintf("%s_mock_examples_test.go", strings.ToLower(entityName)))
		exampleContent := generateMockUsageExamples(entityName)
		if err := os.WriteFile(exampleFile, []byte(exampleContent), 0644); err != nil {
			return err
		}
	}

	return nil
}

// generateRepositoryMock generates a mock for repository interface
func generateRepositoryMock(entityName string) string {
	lowerEntity := strings.ToLower(entityName)

	return fmt.Sprintf(`package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/sazardev/goca/internal/domain"
)

// Mock%sRepository is a mock implementation of repository.%sRepository
type Mock%sRepository struct {
	mock.Mock
}

// Save mocks the Save method
func (m *Mock%sRepository) Save(%s *domain.%s) error {
	args := m.Called(%s)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *Mock%sRepository) FindByID(id int) (*domain.%s, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.%s), args.Error(1)
}

// Update mocks the Update method
func (m *Mock%sRepository) Update(%s *domain.%s) error {
	args := m.Called(%s)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *Mock%sRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// FindAll mocks the FindAll method
func (m *Mock%sRepository) FindAll() ([]domain.%s, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.%s), args.Error(1)
}

// NewMock%sRepository creates a new mock repository
func NewMock%sRepository() *Mock%sRepository {
	return &Mock%sRepository{}
}
`,
		entityName, entityName, entityName,
		entityName, lowerEntity, entityName, lowerEntity,
		entityName, entityName, entityName,
		entityName, lowerEntity, entityName, lowerEntity,
		entityName,
		entityName, entityName, entityName,
		entityName, entityName, entityName, entityName,
	)
}

// generateUseCaseMock generates a mock for use case interface
func generateUseCaseMock(entityName string) string {
	return fmt.Sprintf(`package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/usecase"
)

// Mock%sUseCase is a mock implementation of usecase.%sUseCase
type Mock%sUseCase struct {
	mock.Mock
}

// Create mocks the Create method
func (m *Mock%sUseCase) Create(input usecase.Create%sInput) (*usecase.Create%sOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.Create%sOutput), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *Mock%sUseCase) GetByID(id uint) (*domain.%s, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.%s), args.Error(1)
}

// Update mocks the Update method
func (m *Mock%sUseCase) Update(id uint, input usecase.Update%sInput) (*domain.%s, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.%s), args.Error(1)
}

// Delete mocks the Delete method
func (m *Mock%sUseCase) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// List mocks the List method
func (m *Mock%sUseCase) List() (*usecase.List%sOutput, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.List%sOutput), args.Error(1)
}

// NewMock%sUseCase creates a new mock use case
func NewMock%sUseCase() *Mock%sUseCase {
	return &Mock%sUseCase{}
}
`,
		entityName, entityName, entityName,
		entityName, entityName, entityName, entityName,
		entityName, entityName, entityName,
		entityName, entityName, entityName, entityName,
		entityName,
		entityName, entityName, entityName,
		entityName, entityName, entityName, entityName,
	)
}

// generateHandlerMock generates a mock for HTTP handler interface
func generateHandlerMock(entityName string) string {
	return fmt.Sprintf(`package mocks

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
	"github.com/sazardev/goca/internal/usecase"
)

// Example: Testing use case with mocked repository
func TestCreate%s_WithMockRepository(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMock%sRepository()
	service := usecase.New%sService(mockRepo)

	input := usecase.Create%sInput{
		// Add your input fields here
	}

	expected%s := &domain.%s{
		ID: 1,
		// Add your expected fields here
	}

	// Setup mock expectation
	mockRepo.On("Save", mock.AnythingOfType("*domain.%s")).Return(nil).Run(func(args mock.Arguments) {
		%s := args.Get(0).(*domain.%s)
		%s.ID = 1 // Simulate auto-increment
	})

	// Act
	output, err := service.Create(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, expected%s.ID, output.%s.ID)

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockRepo.AssertCalled(t, "Save", mock.AnythingOfType("*domain.%s"))
}

// Example: Testing error scenarios
func TestGet%sByID_NotFound_WithMockRepository(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMock%sRepository()
	service := usecase.New%sService(mockRepo)

	expectedErr := errors.New("%s not found")
	mockRepo.On("FindByID", 999).Return(nil, expectedErr)

	// Act
	result, err := service.GetByID(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// Example: Testing use case with multiple repository calls
func TestUpdate%s_WithMockRepository(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMock%sRepository()
	service := usecase.New%sService(mockRepo)

	existing%s := &domain.%s{
		ID: 1,
		// Add fields
	}

	updateInput := usecase.Update%sInput{
		// Add update fields
	}

	// Setup mock expectations (FindByID then Update)
	mockRepo.On("FindByID", 1).Return(existing%s, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.%s")).Return(nil)

	// Act
	result, err := service.Update(1, updateInput)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify call order and expectations
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByID", 1)
	mockRepo.AssertNumberOfCalls(t, "Update", 1)
}

// Example: Testing handler with mocked use case
func TestCreate%sHandler_WithMockUseCase(t *testing.T) {
	// Arrange
	mockUC := mocks.NewMock%sUseCase()
	// handler := http.New%sHandler(mockUC)

	expectedOutput := &usecase.Create%sOutput{
		%s: domain.%s{ID: 1},
		Message: "%s created successfully",
	}

	mockUC.On("Create", mock.AnythingOfType("usecase.Create%sInput")).Return(expectedOutput, nil)

	// Act
	// Create HTTP request and response recorder
	// Call handler method
	// result, err := mockUC.Create(input)

	// Assert
	// assert.NoError(t, err)
	// assert.Equal(t, http.StatusCreated, recorder.Code)

	// Verify expectations
	mockUC.AssertExpectations(t)
}

// Example: Testing argument matchers
func TestSave%s_WithArgumentMatchers(t *testing.T) {
	mockRepo := mocks.NewMock%sRepository()

	// Match any %s with specific field value
	mockRepo.On("Save", mock.MatchedBy(func(%s *domain.%s) bool {
		return %s.ID > 0
	})).Return(nil)

	// Test with matching condition
	valid%s := &domain.%s{ID: 1}
	err := mockRepo.Save(valid%s)
	assert.NoError(t, err)

	// Test with non-matching condition
	invalid%s := &domain.%s{ID: 0}
	err = mockRepo.Save(invalid%s)
	assert.Error(t, err) // Will fail because matcher doesn't match

	mockRepo.AssertExpectations(t)
}

// Example: Testing method call verification
func TestDelete%s_CallVerification(t *testing.T) {
	mockRepo := mocks.NewMock%sRepository()
	service := usecase.New%sService(mockRepo)

	mockRepo.On("Delete", 1).Return(nil)

	// Act
	err := service.Delete(1)

	// Assert
	assert.NoError(t, err)

	// Verify Delete was called exactly once with argument 1
	mockRepo.AssertCalled(t, "Delete", 1)
	mockRepo.AssertNumberOfCalls(t, "Delete", 1)
	mockRepo.AssertNotCalled(t, "FindByID")
}
`,
		entityName,                                   // TestCreate%s_WithMockRepository
		entityName, entityName, entityName,           // NewMock%sRepository, New%sService, Create%sInput
		entityName, entityName,                       // expected%s, domain.%s
		entityName,                                   // *domain.%s
		lowerEntity, entityName, lowerEntity,         // %s := args.Get(0), domain.%s, %s.ID = 1
		entityName, entityName,                       // expected%s.ID, output.%s.ID
		entityName,                                   // *domain.%s
		entityName,                                   // TestGet%sByID_NotFound
		entityName, entityName,                       // NewMock%sRepository, New%sService
		lowerEntity,                                  // %s not found
		entityName,                                   // TestUpdate%s_WithMockRepository
		entityName, entityName,                       // NewMock%sRepository, New%sService
		entityName, entityName,                       // existing%s, domain.%s
		entityName,                                   // Update%sInput
		entityName, entityName,                       // existing%s, *domain.%s
		entityName,                                   // TestCreate%sHandler
		entityName, entityName, entityName,           // NewMock%sUseCase, New%sHandler, Create%sOutput
		entityName, entityName,                       // %s: domain.%s
		entityName,                                   // %s created successfully
		entityName,                                   // Create%sInput
		lowerEntity,                                  // TestSave%s_WithArgumentMatchers
		entityName,                                   // NewMock%sRepository
		lowerEntity,                                  // Match any %s with
		lowerEntity, entityName, lowerEntity,         // func(%s *domain.%s) bool, %s.ID > 0
		entityName, entityName, entityName,           // valid%s, domain.%s, valid%s
		entityName, entityName, entityName,           // invalid%s, domain.%s, invalid%s
		entityName,                                   // TestDelete%s_CallVerification
		entityName, entityName,                       // NewMock%sRepository, New%sService
	)
}
