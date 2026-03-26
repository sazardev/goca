package examples

import (
	"errors"
	"testing"

	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/mocks"
	"github.com/sazardev/goca/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Example: Testing use case with mocked repository
func TestCreateProduct_WithMockRepository(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockProductRepository()
	service := usecase.NewProductService(mockRepo)

	input := usecase.CreateProductInput{
		Name:        "Test Product",
		Price:       29.99,
		Description: "A test product for unit testing",
	}

	expectedProduct := &domain.Product{
		ID:          1,
		Name:        "Test Product",
		Price:       29.99,
		Description: "A test product for unit testing",
	}

	// Setup mock expectation
	mockRepo.On("Save", mock.AnythingOfType("*domain.Product")).Return(nil).Run(func(args mock.Arguments) {
		product := args.Get(0).(*domain.Product)
		product.ID = 1 // Simulate auto-increment
	})

	// Act
	output, err := service.Create(input)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, output)
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

	newName := "Updated Product"
	newPrice := 39.99
	existingProduct := &domain.Product{
		ID:          1,
		Name:        "Original Product",
		Price:       29.99,
		Description: "Original description",
	}

	updateInput := usecase.UpdateProductInput{
		Name:  &newName,
		Price: &newPrice,
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

	input := usecase.CreateProductInput{
		Name:        "Test Product",
		Price:       19.99,
		Description: "A test product",
	}

	expectedOutput := &usecase.CreateProductOutput{
		Product: domain.Product{ID: 1, Name: "Test Product", Price: 19.99, Description: "A test product"},
		Message: "Product created successfully",
	}

	mockUC.On("Create", mock.AnythingOfType("usecase.CreateProductInput")).Return(expectedOutput, nil)

	// Act — simulate what a handler would do
	result, err := mockUC.Create(input)

	// Assert
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, uint(1), result.Product.ID)
	assert.Equal(t, "Product created successfully", result.Message)

	// Verify expectations
	mockUC.AssertExpectations(t)
}

// Example: Testing argument matchers
func TestSaveproduct_WithArgumentMatchers(t *testing.T) {
	mockRepo := mocks.NewMockProductRepository()

	// Match any product with ID > 0
	mockRepo.On("Save", mock.MatchedBy(func(product *domain.Product) bool {
		return product.ID > 0
	})).Return(nil)

	// Test with matching condition — ID > 0
	validProduct := &domain.Product{ID: 1, Name: "Valid", Price: 10, Description: "Valid product"}
	err := mockRepo.Save(validProduct)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// Example: Testing method call verification
func TestDeleteProduct_CallVerification(t *testing.T) {
	mockRepo := mocks.NewMockProductRepository()
	service := usecase.NewProductService(mockRepo)

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
