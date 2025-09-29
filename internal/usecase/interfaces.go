package usecase

import "github.com/sazardev/goca/internal/domain"

// UseCase interfaces
type UserUseCase interface {
	Create(input CreateUserInput) (*CreateUserOutput, error)
	GetByID(id uint) (*domain.User, error)
	Update(id uint, input UpdateUserInput) (*domain.User, error)
	Delete(id uint) error
	List() (*ListUserOutput, error)
}

type ProductUseCase interface {
	Create(input CreateProductInput) (*CreateProductOutput, error)
	GetByID(id uint) (*domain.Product, error)
	Update(id uint, input UpdateProductInput) (*domain.Product, error)
	Delete(id uint) error
	List() (*ListProductOutput, error)
}

type OrderUseCase interface {
	Create(input CreateOrderInput) (*CreateOrderOutput, error)
	GetByID(id uint) (*domain.Order, error)
	Update(id uint, input UpdateOrderInput) (*domain.Order, error)
	Delete(id uint) error
	List() (*ListOrderOutput, error)
}

type TestFeatureUseCase interface {
	Create(input CreateTestFeatureInput) (*CreateTestFeatureOutput, error)
	GetByID(id uint) (*domain.TestFeature, error)
	Update(id uint, input UpdateTestFeatureInput) (*domain.TestFeature, error)
	Delete(id uint) error
	List() (*ListTestFeatureOutput, error)
}
