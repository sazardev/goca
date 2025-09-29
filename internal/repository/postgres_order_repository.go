package repository

import (
	"github.com/sazardev/goca/internal/domain"

	"gorm.io/gorm"
)

type postgresOrderRepository struct {
	db *gorm.DB
}

func NewPostgresOrderRepository(db *gorm.DB) OrderRepository {
	return &postgresOrderRepository{
		db: db,
	}
}

func (p *postgresOrderRepository) Save(order *domain.Order) error {
	result := p.db.Create(order)
	return result.Error
}

func (p *postgresOrderRepository) FindByID(id int) (*domain.Order, error) {
	order := &domain.Order{}
	result := p.db.First(order, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

func (p *postgresOrderRepository) Update(order *domain.Order) error {
	result := p.db.Save(order)
	return result.Error
}

func (p *postgresOrderRepository) Delete(id int) error {
	result := p.db.Delete(&domain.Order{}, id)
	return result.Error
}

func (p *postgresOrderRepository) FindAll() ([]domain.Order, error) {
	var orders []domain.Order
	result := p.db.Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func (p *postgresOrderRepository) FindByCustomer_id(customer_id int) (*domain.Order, error) {
	order := &domain.Order{}
	result := p.db.Where("customer_id = ?", customer_id).First(order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

func (p *postgresOrderRepository) FindByStatus(status string) (*domain.Order, error) {
	order := &domain.Order{}
	result := p.db.Where("status = ?", status).First(order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}
