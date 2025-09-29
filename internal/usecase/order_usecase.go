package usecase

import "github.com/sazardev/goca/internal/domain"

type OrderUseCase interface {
	CreateOrder(input CreateOrderInput) (CreateOrderOutput, error)
	GetOrder(id int) (*domain.Order, error)
	UpdateOrder(id int, input UpdateOrderInput) error
	DeleteOrder(id int) error
	ListOrders() (ListOrderOutput, error)
}
