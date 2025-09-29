package repository

import "github.com/sazardev/goca/internal/domain"

// Repository interfaces
type UserRepository interface {
	Save(user *domain.User) error
	FindByID(id int) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id int) error
	FindAll() ([]domain.User, error)
}

type ProductRepository interface {
	Save(product *domain.Product) error
	FindByID(id int) (*domain.Product, error)
	FindByName(name string) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id int) error
	FindAll() ([]domain.Product, error)
}

type OrderRepository interface {
	Save(order *domain.Order) error
	FindByID(id int) (*domain.Order, error)
	Update(order *domain.Order) error
	Delete(id int) error
	FindAll() ([]domain.Order, error)
}

type TestFeatureRepository interface {
	Save(testfeature *domain.TestFeature) error
	FindByID(id int) (*domain.TestFeature, error)
	FindByEmail(email string) (*domain.TestFeature, error)
	Update(testfeature *domain.TestFeature) error
	Delete(id int) error
	FindAll() ([]domain.TestFeature, error)
}
