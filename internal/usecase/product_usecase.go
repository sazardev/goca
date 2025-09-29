package usecase

import "github.com/sazardev/goca/internal/domain"

type ProductUseCase interface {
	CreateProduct(input CreateProductInput) (CreateProductOutput, error)
	GetProduct(id int) (*domain.Product, error)
	UpdateProduct(id int, input UpdateProductInput) error
	DeleteProduct(id int) error
	ListProducts() (ListProductOutput, error)
}
