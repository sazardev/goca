package usecase

import (
	"github.com/sazardev/goca/internal/domain"
)

type CreateUserInput struct {
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,min=1"`
}

type CreateUserOutput struct {
	User    domain.User `json:"user"`
	Message string      `json:"message"`
}

type UpdateUserInput struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	Age   *int    `json:"age,omitempty" validate:"omitempty,min=1"`
}

type ListUserOutput struct {
	Users   []domain.User `json:"users"`
	Total   int           `json:"total"`
	Message string        `json:"message"`
}

type CreateProductInput struct {
	Name        string  `json:"name" validate:"required,min=2"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Description string  `json:"description" validate:"required,min=5"`
}

type CreateProductOutput struct {
	Product domain.Product `json:"product"`
	Message string         `json:"message"`
}

type UpdateProductInput struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=2"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	Description *string  `json:"description,omitempty" validate:"omitempty,min=5"`
}

type ListProductOutput struct {
	Products []domain.Product `json:"products"`
	Total    int              `json:"total"`
	Message  string           `json:"message"`
}

type CreateOrderInput struct {
	Customer_id int     `json:"customer_id" validate:"required,min=0"`
	Total       float64 `json:"total" validate:"required,min=0"`
	Status      string  `json:"status" validate:"required,min=2"`
}

type CreateOrderOutput struct {
	Order   domain.Order `json:"order"`
	Message string       `json:"message"`
}

type UpdateOrderInput struct {
	Customer_id *int     `json:"customer_id,omitempty" validate:"omitempty,min=0"`
	Total       *float64 `json:"total,omitempty" validate:"omitempty,min=0"`
	Status      *string  `json:"status,omitempty" validate:"omitempty,min=2"`
}

type ListOrderOutput struct {
	Orders  []domain.Order `json:"orders"`
	Total   int            `json:"total"`
	Message string         `json:"message"`
}

type CreateTestFeatureInput struct {
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,min=1"`
}

type CreateTestFeatureOutput struct {
	TestFeature domain.TestFeature `json:"testfeature"`
	Message     string             `json:"message"`
}

type UpdateTestFeatureInput struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	Age   *int    `json:"age,omitempty" validate:"omitempty,min=1"`
}

type ListTestFeatureOutput struct {
	TestFeatures []domain.TestFeature `json:"testfeatures"`
	Total        int                  `json:"total"`
	Message      string               `json:"message"`
}
