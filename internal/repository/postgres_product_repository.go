package repository

import (
	"github.com/sazardev/goca/internal/domain"

	"gorm.io/gorm"
)

type postgresProductRepository struct {
	db *gorm.DB
}

func NewPostgresProductRepository(db *gorm.DB) ProductRepository {
	return &postgresProductRepository{
		db: db,
	}
}

func (p *postgresProductRepository) Save(product *domain.Product) error {
	result := p.db.Create(product)
	return result.Error
}

func (p *postgresProductRepository) FindByID(id int) (*domain.Product, error) {
	product := &domain.Product{}
	result := p.db.First(product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (p *postgresProductRepository) Update(product *domain.Product) error {
	result := p.db.Save(product)
	return result.Error
}

func (p *postgresProductRepository) Delete(id int) error {
	result := p.db.Delete(&domain.Product{}, id)
	return result.Error
}

func (p *postgresProductRepository) FindAll() ([]domain.Product, error) {
	var products []domain.Product
	result := p.db.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (p *postgresProductRepository) FindByName(name string) (*domain.Product, error) {
	product := &domain.Product{}
	result := p.db.Where("name = ?", name).First(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (p *postgresProductRepository) FindByDescription(description string) (*domain.Product, error) {
	product := &domain.Product{}
	result := p.db.Where("description = ?", description).First(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}
