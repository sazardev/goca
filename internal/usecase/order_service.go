package usecase

import (
	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/messages"
	"github.com/sazardev/goca/internal/repository"
)

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderUseCase {
	return &orderService{repo: repo}
}

func (o *orderService) CreateOrder(input CreateOrderInput) (CreateOrderOutput, error) {
	order := domain.Order{
		Customer_id: input.Customer_id,
		Total:       input.Total,
		Status:      input.Status,
	}

	if err := order.Validate(); err != nil {
		return CreateOrderOutput{}, err
	}

	if err := o.repo.Save(&order); err != nil {
		return CreateOrderOutput{}, err
	}

	return CreateOrderOutput{
		Order:   order,
		Message: messages.OrderCreatedSuccessfully,
	}, nil
}

func (o *orderService) GetOrder(id int) (*domain.Order, error) {
	return o.repo.FindByID(id)
}

func (o *orderService) UpdateOrder(id int, input UpdateOrderInput) error {
	order, err := o.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Actualizar campos según tu entidad
	if input.Nombre != "" {
		order.Nombre = input.Nombre
	}
	if input.Email != "" {
		order.Email = input.Email
	}
	// Agregar más campos según necesites

	return o.repo.Update(order)
}

func (o *orderService) DeleteOrder(id int) error {
	return o.repo.Delete(id)
}

func (o *orderService) ListOrders() (ListOrderOutput, error) {
	orders, err := o.repo.FindAll()
	if err != nil {
		return ListOrderOutput{}, err
	}

	return ListOrderOutput{
		Orders:  orders,
		Total:   len(orders),
		Message: messages.OrdersListedSuccessfully,
	}, nil
}
