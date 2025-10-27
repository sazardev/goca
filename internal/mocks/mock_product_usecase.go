package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/usecase"
)

// MockProductUseCase is a mock implementation of usecase.ProductUseCase
type MockProductUseCase struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockProductUseCase) Create(input usecase.CreateProductInput) (*usecase.CreateProductOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.CreateProductOutput), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockProductUseCase) GetByID(id uint) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

// Update mocks the Update method
func (m *MockProductUseCase) Update(id uint, input usecase.UpdateProductInput) (*domain.Product, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockProductUseCase) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// List mocks the List method
func (m *MockProductUseCase) List() (*usecase.ListProductOutput, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.ListProductOutput), args.Error(1)
}

// NewMockProductUseCase creates a new mock use case
func NewMockProductUseCase() *MockProductUseCase {
	return &MockProductUseCase{}
}
