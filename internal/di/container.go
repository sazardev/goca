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
	testperfRepo repository.TestPerfRepository

	// Use Cases
	testperfUC usecase.TestPerfUseCase

	// Handlers
	testPerfHandler *http.TestPerfHandler
}

func NewContainer(db *gorm.DB) *Container {
	c := &Container{db: db}
	c.setupRepositories()
	c.setupUseCases()
	c.setupHandlers()
	return c
}

func (c *Container) setupRepositories() {
	c.testperfRepo = repository.NewPostgresTestPerfRepository(c.db)
}

func (c *Container) setupUseCases() {
	c.testperfUC = usecase.NewTestPerfService(c.testperfRepo)
}

func (c *Container) setupHandlers() {
	c.testPerfHandler = http.NewTestPerfHandler(c.testperfUC)
}

// Getters
func (c *Container) TestPerfHandler() *http.TestPerfHandler {
	return c.testPerfHandler
}

func (c *Container) TestPerfUseCase() usecase.TestPerfUseCase {
	return c.testperfUC
}

func (c *Container) TestPerfRepository() repository.TestPerfRepository {
	return c.testperfRepo
}
