package usecase

import "github.com/sazardev/goca/internal/domain"

type OrderRepository interface {
	Save(order *domain.Order) error
	FindByID(id int) (*domain.Order, error)
	FindByEmail(email string) (*domain.Order, error)
	Update(order *domain.Order) error
	Delete(id int) error
	FindAll() ([]domain.Order, error)
}
