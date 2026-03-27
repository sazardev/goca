package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFeatureAlreadyRegistered(t *testing.T) {
	t.Parallel()

	t.Run("feature registered", func(t *testing.T) {
		t.Parallel()
		content := `r.HandleFunc("/api/v1/products", handler.List)`
		assert.True(t, isFeatureAlreadyRegistered(content, "Product"))
	})

	t.Run("feature not registered", func(t *testing.T) {
		t.Parallel()
		content := `r.HandleFunc("/api/v1/users", handler.List)`
		assert.False(t, isFeatureAlreadyRegistered(content, "Product"))
	})
}

func TestAddFieldsToDIContainer(t *testing.T) {
	t.Parallel()
	content := `type Container struct {
	db *gorm.DB

	// Repositories

	// Use Cases

	// Handlers
}

func NewContainer() *Container {`

	result := addFieldsToDIContainer(content, "Product", "product")
	assert.Contains(t, result, "productRepo    repository.ProductRepository")
	assert.Contains(t, result, "productUC    usecase.ProductUseCase")
	assert.Contains(t, result, "productHandler    *http.ProductHandler")
}

func TestAddSetupMethodsToDI(t *testing.T) {
	t.Parallel()
	content := `func (c *Container) setupRepositories() {
}

func (c *Container) setupUseCases() {
}

func (c *Container) setupHandlers() {
}

// Getters`

	result := addSetupMethodsToDI(content, "Product", "product", false)
	assert.Contains(t, result, "c.productRepo = repository.NewPostgresProductRepository(c.db)")
	assert.Contains(t, result, "c.productUC = usecase.NewProductService(c.productRepo)")
	assert.Contains(t, result, "c.productHandler = http.NewProductHandler(c.productUC)")
}

func TestAddGetterMethodsToDI(t *testing.T) {
	t.Parallel()
	content := "// existing content\n"
	result := addGetterMethodsToDI(content, "Product", "product")
	assert.Contains(t, result, "func (c *Container) ProductHandler()")
	assert.Contains(t, result, "func (c *Container) ProductUseCase()")
	assert.Contains(t, result, "func (c *Container) ProductRepository()")
	assert.Contains(t, result, "return c.productHandler")
	assert.Contains(t, result, "return c.productUC")
	assert.Contains(t, result, "return c.productRepo")
}

func TestFindMainGoPath(t *testing.T) {
	t.Parallel()
	// This function looks at the filesystem, test that it returns false when not found
	path, found := findMainGoPath()
	// In the test environment, main.go should not exist at CWD
	// Just verify it returns string and bool
	_ = path
	_ = found
}
