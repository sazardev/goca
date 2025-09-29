package usecase

import (
	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/repository"
)

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderUseCase {
	return &orderService{repo: repo}
}

func (o *orderService) Create(input CreateOrderInput) (*CreateOrderOutput, error) {
	order := domain.Order{
		Customer_id: input.Customer_id,
		Total:       input.Total,
		Status:      input.Status,
	}

	if err := order.Validate(); err != nil {
		return nil, err
	}

	if err := o.repo.Save(&order); err != nil {
		return nil, err
	}

	return &CreateOrderOutput{
		Order:   order,
		Message: "Order created successfully",
	}, nil
}

func (o *orderService) GetByID(id uint) (*domain.Order, error) {
	return o.repo.FindByID(int(id))
}

func (o *orderService) Update(id uint, input UpdateOrderInput) (*domain.Order, error) {
	order, err := o.repo.FindByID(int(id))
	if err != nil {
		return nil, err
	}

	if input.Customer_id != nil {
		order.Customer_id = *input.Customer_id
	}
	if input.Total != nil {
		order.Total = *input.Total
	}
	if input.Status != nil {
		order.Status = *input.Status
	}

	err = o.repo.Update(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *orderService) Delete(id uint) error {
	return o.repo.Delete(int(id))
}

func (o *orderService) List() (*ListOrderOutput, error) {
	orders, err := o.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return &ListOrderOutput{
		Orders:  orders,
		Total:   len(orders),
		Message: "Orders listed successfully",
	}, nil
}
