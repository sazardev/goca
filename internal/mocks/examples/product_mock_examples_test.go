package examples

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
func TestCreateProduct_WithMockRepository(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockProductRepository()
	service := usecase.NewProductService(mockRepo)

	input := usecase.CreateProductInput{
		// Add your input fields here
	}

	expectedProduct := &domain.Product{
		ID: 1,
		// Add your expected fields here
	}

	// Setup mock expectation
	mockRepo.On("Save", mock.AnythingOfType("*domain.Product")).Return(nil).Run(func(args mock.Arguments) {
		product := args.Get(0).(*domain.Product)
		product.ID = 1 // Simulate auto-increment
	})

	// Act
	output, err := service.Create(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, expectedProduct.ID, output.Product.ID)

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockRepo.AssertCalled(t, "Save", mock.AnythingOfType("*domain.Product"))
}

// Example: Testing error scenarios
func TestGetProductByID_NotFound_WithMockRepository(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockProductRepository()
	service := usecase.NewProductService(mockRepo)

	expectedErr := errors.New("product not found")
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
func TestUpdateProduct_WithMockRepository(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockProductRepository()
	service := usecase.NewProductService(mockRepo)

	existingProduct := &domain.Product{
		ID: 1,
		// Add fields
	}

	updateInput := usecase.UpdateProductInput{
		// Add update fields
	}

	// Setup mock expectations (FindByID then Update)
	mockRepo.On("FindByID", 1).Return(existingProduct, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Product")).Return(nil)

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
func TestCreateProductHandler_WithMockUseCase(t *testing.T) {
	// Arrange
	mockUC := mocks.NewMockProductUseCase()
	// handler := http.NewProductHandler(mockUC)

	expectedOutput := &usecase.CreateProductOutput{
		Product: domain.Product{ID: 1},
		Message: "Product created successfully",
	}

	mockUC.On("Create", mock.AnythingOfType("usecase.CreateProductInput")).Return(expectedOutput, nil)

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
func TestSaveproduct_WithArgumentMatchers(t *testing.T) {
	mockRepo := mocks.NewMockProductRepository()

	// Match any product with specific field value
	mockRepo.On("Save", mock.MatchedBy(func(Product *domain.Product) bool {
		return Product.ID > 0
	})).Return(nil)

	// Test with matching condition
	validProduct := &domain.Product{ID: 1}
	err := mockRepo.Save(validProduct)
	assert.NoError(t, err)

	// Test with non-matching condition
	invalidProduct := &domain.Product{ID: 0}
	err = mockRepo.Save(invalidProduct)
	assert.Error(t, err) // Will fail because matcher doesn't match

	mockRepo.AssertExpectations(t)
}

// Example: Testing method call verification
func TestDelete%!s(MISSING)_CallVerification(t *testing.T) {
	mockRepo := mocks.NewMock%!s(MISSING)Repository()
	service := usecase.New%!s(MISSING)Service(mockRepo)

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
