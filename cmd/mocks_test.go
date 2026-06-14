package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRepositoryMock(t *testing.T) {
	t.Parallel()
	fields := parseFields("name:string,email:string,age:int")
	result := generateRepositoryMock("Product", fields)
	assert.Contains(t, result, "MockProductRepository")
	assert.Contains(t, result, "mock.Mock")
	assert.Contains(t, result, "func (m *MockProductRepository) Save(")
	assert.Contains(t, result, "func (m *MockProductRepository) FindByID(id int)")
	assert.Contains(t, result, "func (m *MockProductRepository) Update(")
	assert.Contains(t, result, "func (m *MockProductRepository) Delete(id int)")
	assert.Contains(t, result, "func (m *MockProductRepository) FindAll(")
	// Per-field finders matching the real repository interface.
	assert.Contains(t, result, "func (m *MockProductRepository) FindByName(name string) (*domain.Product, error)")
	assert.Contains(t, result, "func (m *MockProductRepository) FindByEmail(email string) (*domain.Product, error)")
	assert.NotContains(t, result, "TODO")
	assert.Contains(t, result, "NewMockProductRepository")
	assert.Contains(t, result, "domain.Product")
}

func TestGenerateUseCaseMock(t *testing.T) {
	t.Parallel()
	result := generateUseCaseMock("Product")
	assert.Contains(t, result, "MockProductUseCase")
	assert.Contains(t, result, "mock.Mock")
	// Method names and signatures must match usecase.ProductUseCase exactly.
	assert.Contains(t, result, "func (m *MockProductUseCase) CreateProduct(input usecase.CreateProductInput) (usecase.CreateProductOutput, error)")
	assert.Contains(t, result, "func (m *MockProductUseCase) GetProduct(id int) (*domain.Product, error)")
	assert.Contains(t, result, "func (m *MockProductUseCase) UpdateProduct(id int, input usecase.UpdateProductInput) error")
	assert.Contains(t, result, "func (m *MockProductUseCase) DeleteProduct(id int) error")
	assert.Contains(t, result, "func (m *MockProductUseCase) ListProducts() (usecase.ListProductOutput, error)")
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
	assert.Contains(t, result, "TestMockProductRepository_Usage")
	assert.Contains(t, result, "MockProductRepository")
	assert.Contains(t, result, "assert")
	// The example must use the real API, not the old nonexistent methods.
	assert.NotContains(t, result, "service.Create(")
	assert.NotContains(t, result, "output.Product.ID")
	// Compile-time interface assertions are included.
	assert.Contains(t, result, "_ repository.ProductRepository = (*mocks.MockProductRepository)(nil)")
}
