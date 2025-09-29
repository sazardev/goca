package usecase

import (
	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/messages"
	"github.com/sazardev/goca/internal/repository"
)

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductUseCase {
	return &productService{repo: repo}
}

func (p *productService) CreateProduct(input CreateProductInput) (CreateProductOutput, error) {
	product := domain.Product{
		Name:        input.Name,
		Price:       input.Price,
		Description: input.Description,
	}

	if err := product.Validate(); err != nil {
		return CreateProductOutput{}, err
	}

	if err := p.repo.Save(&product); err != nil {
		return CreateProductOutput{}, err
	}

	return CreateProductOutput{
		Product: product,
		Message: messages.ProductCreatedSuccessfully,
	}, nil
}

func (p *productService) GetProduct(id int) (*domain.Product, error) {
	return p.repo.FindByID(id)
}

func (p *productService) UpdateProduct(id int, input UpdateProductInput) error {
	product, err := p.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Actualizar campos según tu entidad
	if input.Nombre != "" {
		product.Nombre = input.Nombre
	}
	if input.Email != "" {
		product.Email = input.Email
	}
	// Agregar más campos según necesites

	return p.repo.Update(product)
}

func (p *productService) DeleteProduct(id int) error {
	return p.repo.Delete(id)
}

func (p *productService) ListProducts() (ListProductOutput, error) {
	products, err := p.repo.FindAll()
	if err != nil {
		return ListProductOutput{}, err
	}

	return ListProductOutput{
		Products: products,
		Total:    len(products),
		Message:  messages.ProductsListedSuccessfully,
	}, nil
}
