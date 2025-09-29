package usecase

import (
	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/repository"
)

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductUseCase {
	return &productService{repo: repo}
}

func (p *productService) Create(input CreateProductInput) (*CreateProductOutput, error) {
	product := domain.Product{
		Name:        input.Name,
		Price:       input.Price,
		Description: input.Description,
	}

	if err := product.Validate(); err != nil {
		return nil, err
	}

	if err := p.repo.Save(&product); err != nil {
		return nil, err
	}

	return &CreateProductOutput{
		Product: product,
		Message: "Product created successfully",
	}, nil
}

func (p *productService) GetByID(id uint) (*domain.Product, error) {
	return p.repo.FindByID(int(id))
}

func (p *productService) Update(id uint, input UpdateProductInput) (*domain.Product, error) {
	product, err := p.repo.FindByID(int(id))
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Price != nil {
		product.Price = *input.Price
	}
	if input.Description != nil {
		product.Description = *input.Description
	}

	err = p.repo.Update(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *productService) Delete(id uint) error {
	return p.repo.Delete(int(id))
}

func (p *productService) List() (*ListProductOutput, error) {
	products, err := p.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return &ListProductOutput{
		Products: products,
		Total:    len(products),
		Message:  "Products listed successfully",
	}, nil
}
