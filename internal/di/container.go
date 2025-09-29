package di

import (
	"gorm.io/gorm"

	"github.com/sazardev/goca/internal/handler/http"
	"github.com/sazardev/goca/internal/repository"
	"github.com/sazardev/goca/internal/usecase"
)

type Container struct {
	db *gorm.DB

	// Repositories
	testfeatureRepo repository.TestFeatureRepository
	userRepo    repository.UserRepository
	productRepo    repository.ProductRepository
	orderRepo    repository.OrderRepository

	// Use Cases
	testfeatureUC usecase.TestFeatureUseCase
	userUC    usecase.UserUseCase
	productUC    usecase.ProductUseCase
	orderUC    usecase.OrderUseCase

	// Handlers
	testFeatureHandler *http.TestFeatureHandler
	userHandler    *http.UserHandler
	productHandler    *http.ProductHandler
	orderHandler    *http.OrderHandler
}

func NewContainer(db *gorm.DB) *Container {
	c := &Container{db: db}
	c.setupRepositories()
	c.setupUseCases()
	c.setupHandlers()
	return c
}

func (c *Container) setupRepositories() {
	c.testfeatureRepo = repository.NewPostgresTestFeatureRepository(c.db)
	c.userRepo = repository.NewPostgresUserRepository(c.db)
	c.productRepo = repository.NewPostgresProductRepository(c.db)
	c.orderRepo = repository.NewPostgresOrderRepository(c.db)
}

func (c *Container) setupUseCases() {
	c.testfeatureUC = usecase.NewTestFeatureService(c.testfeatureRepo)
	c.userUC = usecase.NewUserService(c.userRepo)
	c.productUC = usecase.NewProductService(c.productRepo)
	c.orderUC = usecase.NewOrderService(c.orderRepo)
}

func (c *Container) setupHandlers() {
	c.testFeatureHandler = http.NewTestFeatureHandler(c.testfeatureUC)
	c.userHandler = http.NewUserHandler(c.userUC)
	c.productHandler = http.NewProductHandler(c.productUC)
	c.orderHandler = http.NewOrderHandler(c.orderUC)
}

// Getters
func (c *Container) TestFeatureHandler() *http.TestFeatureHandler {
	return c.testFeatureHandler
}

func (c *Container) TestFeatureUseCase() usecase.TestFeatureUseCase {
	return c.testfeatureUC
}

func (c *Container) TestFeatureRepository() repository.TestFeatureRepository {
	return c.testfeatureRepo
}
func (c *Container) UserHandler() *http.UserHandler {
	return c.userHandler
}

func (c *Container) UserUseCase() usecase.UserUseCase {
	return c.userUC
}

func (c *Container) UserRepository() repository.UserRepository {
	return c.userRepo
}

func (c *Container) ProductHandler() *http.ProductHandler {
	return c.productHandler
}

func (c *Container) ProductUseCase() usecase.ProductUseCase {
	return c.productUC
}

func (c *Container) ProductRepository() repository.ProductRepository {
	return c.productRepo
}

func (c *Container) OrderHandler() *http.OrderHandler {
	return c.orderHandler
}

func (c *Container) OrderUseCase() usecase.OrderUseCase {
	return c.orderUC
}

func (c *Container) OrderRepository() repository.OrderRepository {
	return c.orderRepo
}

