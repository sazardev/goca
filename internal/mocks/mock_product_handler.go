package mocks

import (
	"net/http"
	"github.com/stretchr/testify/mock"
)

// MockProductHandler is a mock implementation of HTTP handler
type MockProductHandler struct {
	mock.Mock
}

// CreateProduct mocks the CreateProduct HTTP handler method
func (m *MockProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// GetProduct mocks the GetProduct HTTP handler method
func (m *MockProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// UpdateProduct mocks the UpdateProduct HTTP handler method
func (m *MockProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// DeleteProduct mocks the DeleteProduct HTTP handler method
func (m *MockProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// ListProducts mocks the ListProducts HTTP handler method
func (m *MockProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// NewMockProductHandler creates a new mock handler
func NewMockProductHandler() *MockProductHandler {
	return &MockProductHandler{}
}
