package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/sazardev/goca/internal/domain"
)

// MockProductRepository is a mock implementation of repository.ProductRepository
type MockProductRepository struct {
	mock.Mock
}

// Save mocks the Save method
func (m *MockProductRepository) Save(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockProductRepository) FindByID(id int) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

// Update mocks the Update method
func (m *MockProductRepository) Update(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockProductRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// FindAll mocks the FindAll method
func (m *MockProductRepository) FindAll() ([]domain.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Product), args.Error(1)
}

// NewMockProductRepository creates a new mock repository
func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{}
}
