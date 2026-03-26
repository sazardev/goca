package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRepositoryMock(t *testing.T) {
	t.Parallel()
	result := generateRepositoryMock("Product")
	assert.Contains(t, result, "MockProductRepository")
	assert.Contains(t, result, "mock.Mock")
	assert.Contains(t, result, "func (m *MockProductRepository) Save(")
	assert.Contains(t, result, "func (m *MockProductRepository) FindByID(")
	assert.Contains(t, result, "func (m *MockProductRepository) Update(")
	assert.Contains(t, result, "func (m *MockProductRepository) Delete(")
	assert.Contains(t, result, "func (m *MockProductRepository) FindAll(")
	assert.Contains(t, result, "NewMockProductRepository")
	assert.Contains(t, result, "domain.Product")
}

func TestGenerateUseCaseMock(t *testing.T) {
	t.Parallel()
	result := generateUseCaseMock("Product")
	assert.Contains(t, result, "MockProductUseCase")
	assert.Contains(t, result, "mock.Mock")
	assert.Contains(t, result, "func (m *MockProductUseCase) Create(")
	assert.Contains(t, result, "func (m *MockProductUseCase) GetByID(")
	assert.Contains(t, result, "func (m *MockProductUseCase) Update(")
	assert.Contains(t, result, "func (m *MockProductUseCase) Delete(")
	assert.Contains(t, result, "func (m *MockProductUseCase) List(")
	assert.Contains(t, result, "NewMockProductUseCase")
	assert.Contains(t, result, "usecase.CreateProductInput")
}

func TestGenerateHandlerMock(t *testing.T) {
	t.Parallel()
	result := generateHandlerMock("Product")
	assert.Contains(t, result, "MockProductHandler")
	assert.Contains(t, result, "mock.Mock")
	assert.Contains(t, result, "func (m *MockProductHandler) CreateProduct(")
	assert.Contains(t, result, "func (m *MockProductHandler) GetProduct(")
	assert.Contains(t, result, "func (m *MockProductHandler) UpdateProduct(")
	assert.Contains(t, result, "func (m *MockProductHandler) DeleteProduct(")
	assert.Contains(t, result, "func (m *MockProductHandler) ListProducts(")
	assert.Contains(t, result, "NewMockProductHandler")
}

func TestGenerateMockUsageExamples(t *testing.T) {
	t.Parallel()
	result := generateMockUsageExamples("Product")
	assert.Contains(t, result, "TestCreate")
	assert.Contains(t, result, "MockProductRepository")
	assert.Contains(t, result, "assert")
}
