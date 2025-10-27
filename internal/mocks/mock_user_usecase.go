package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/usecase"
)

// MockUserUseCase is a mock implementation of usecase.UserUseCase
type MockUserUseCase struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockUserUseCase) Create(input usecase.CreateUserInput) (*usecase.CreateUserOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.CreateUserOutput), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockUserUseCase) GetByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// Update mocks the Update method
func (m *MockUserUseCase) Update(id uint, input usecase.UpdateUserInput) (*domain.User, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockUserUseCase) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// List mocks the List method
func (m *MockUserUseCase) List() (*usecase.ListUserOutput, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.ListUserOutput), args.Error(1)
}

// NewMockUserUseCase creates a new mock use case
func NewMockUserUseCase() *MockUserUseCase {
	return &MockUserUseCase{}
}
