package usecase

import (
	"github.com/sazardev/goca/internal/domain"
)

// CreateTestFeatureRequest DTO para crear un nuevo testfeature
type CreateTestFeatureRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// CreateTestFeatureResponse DTO para la respuesta de creación
type CreateTestFeatureResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Age     int    `json:"age"`
	Message string `json:"message"`
}

type UpdateTestFeatureInput struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Age   *int    `json:"age,omitempty"`
}

type ListTestFeatureOutput struct {
	TestFeatures []domain.TestFeature `json:"testfeatures"`
	Total        int                  `json:"total"`
	Message      string               `json:"message"`
}

// CreateUserRequest DTO para crear un nuevo user
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// CreateUserResponse DTO para la respuesta de creación
type CreateUserResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Age     int    `json:"age"`
	Message string `json:"message"`
}

type UpdateUserInput struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Age   *int    `json:"age,omitempty"`
}

type ListUserOutput struct {
	Users   []domain.User `json:"users"`
	Total   int           `json:"total"`
	Message string        `json:"message"`
}

// CreateProductRequest DTO para crear un nuevo product
type CreateProductRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

// CreateProductResponse DTO para la respuesta de creación
type CreateProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Message     string  `json:"message"`
}

type UpdateProductInput struct {
	Name        *string  `json:"name,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Description *string  `json:"description,omitempty"`
}

type ListProductOutput struct {
	Products []domain.Product `json:"products"`
	Total    int              `json:"total"`
	Message  string           `json:"message"`
}

// CreateOrderRequest DTO para crear un nuevo order
type CreateOrderRequest struct {
	Customer_id int     `json:"customer_id"`
	Total       float64 `json:"total"`
	Status      string  `json:"status"`
}

// CreateOrderResponse DTO para la respuesta de creación
type CreateOrderResponse struct {
	ID          uint    `json:"id"`
	Customer_id int     `json:"customer_id"`
	Total       float64 `json:"total"`
	Status      string  `json:"status"`
	Message     string  `json:"message"`
}

type UpdateOrderInput struct {
	Customer_id *int     `json:"customer_id,omitempty"`
	Total       *float64 `json:"total,omitempty"`
	Status      *string  `json:"status,omitempty"`
}

type ListOrderOutput struct {
	Orders  []domain.Order `json:"orders"`
	Total   int            `json:"total"`
	Message string         `json:"message"`
}
