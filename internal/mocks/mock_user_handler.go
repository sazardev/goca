package mocks

import (
	"net/http"
	"github.com/stretchr/testify/mock"
)

// MockUserHandler is a mock implementation of HTTP handler
type MockUserHandler struct {
	mock.Mock
}

// CreateUser mocks the CreateUser HTTP handler method
func (m *MockUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// GetUser mocks the GetUser HTTP handler method
func (m *MockUserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// UpdateUser mocks the UpdateUser HTTP handler method
func (m *MockUserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// DeleteUser mocks the DeleteUser HTTP handler method
func (m *MockUserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// ListUsers mocks the ListUsers HTTP handler method
func (m *MockUserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// NewMockUserHandler creates a new mock handler
func NewMockUserHandler() *MockUserHandler {
	return &MockUserHandler{}
}
