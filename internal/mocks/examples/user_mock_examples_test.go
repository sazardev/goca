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
func TestCreateUser_WithMockRepository(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockUserRepository()
	service := usecase.NewUserService(mockRepo)

	input := usecase.CreateUserInput{
		// Add your input fields here
	}

	expectedUser := &domain.User{
		ID: 1,
		// Add your expected fields here
	}

	// Setup mock expectation
	mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(0).(*domain.User)
		user.ID = 1 // Simulate auto-increment
	})

	// Act
	output, err := service.Create(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, expectedUser.ID, output.User.ID)

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockRepo.AssertCalled(t, "Save", mock.AnythingOfType("*domain.User"))
}

// Example: Testing error scenarios
func TestGetUserByID_NotFound_WithMockRepository(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockUserRepository()
	service := usecase.NewUserService(mockRepo)

	expectedErr := errors.New("user not found")
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
func TestUpdateUser_WithMockRepository(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockUserRepository()
	service := usecase.NewUserService(mockRepo)

	existingUser := &domain.User{
		ID: 1,
		// Add fields
	}

	updateInput := usecase.UpdateUserInput{
		// Add update fields
	}

	// Setup mock expectations (FindByID then Update)
	mockRepo.On("FindByID", 1).Return(existingUser, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

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
func TestCreateUserHandler_WithMockUseCase(t *testing.T) {
	// Arrange
	mockUC := mocks.NewMockUserUseCase()
	// handler := http.NewUserHandler(mockUC)

	expectedOutput := &usecase.CreateUserOutput{
		User: domain.User{ID: 1},
		Message: "User created successfully",
	}

	mockUC.On("Create", mock.AnythingOfType("usecase.CreateUserInput")).Return(expectedOutput, nil)

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
func TestSaveuser_WithArgumentMatchers(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()

	// Match any user with specific field value
	mockRepo.On("Save", mock.MatchedBy(func(user *domain.User) bool {
		return user.ID > 0
	})).Return(nil)

	// Test with matching condition
	validUser := &domain.User{ID: 1}
	err := mockRepo.Save(validUser)
	assert.NoError(t, err)

	// Test with non-matching condition
	invalidUser := &domain.User{ID: 0}
	err = mockRepo.Save(invalidUser)
	assert.Error(t, err) // Will fail because matcher doesn't match

	mockRepo.AssertExpectations(t)
}

// Example: Testing method call verification
func TestDeleteUser_CallVerification(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	service := usecase.NewUserService(mockRepo)

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
