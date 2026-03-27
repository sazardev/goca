package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSetupRepositories(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		database string
		prefix   string
	}{
		{"postgres", "postgres", "NewPostgres"},
		{"mysql", "mysql", "NewMySQL"},
		{"mongodb", "mongodb", "NewMongo"},
		{"default", "other", "NewPostgres"},
	}

	features := []string{"Product", "User"}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var b strings.Builder
			generateSetupRepositories(&b, features, tc.database, false)
			output := b.String()
			assert.Contains(t, output, "func (c *Container) setupRepositories()")
			assert.Contains(t, output, tc.prefix+"ProductRepository")
			assert.Contains(t, output, tc.prefix+"UserRepository")
		})
	}
}

func TestGenerateSetupUseCases(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	generateSetupUseCases(&b, []string{"Product", "User"})
	output := b.String()
	assert.Contains(t, output, "func (c *Container) setupUseCases()")
	assert.Contains(t, output, "usecase.NewProductService(c.productRepo)")
	assert.Contains(t, output, "usecase.NewUserService(c.userRepo)")
}

func TestGenerateSetupHandlers(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	generateSetupHandlers(&b, []string{"Product", "User"})
	output := b.String()
	assert.Contains(t, output, "func (c *Container) setupHandlers()")
	assert.Contains(t, output, "http.NewProductHandler(c.productUC)")
	assert.Contains(t, output, "http.NewUserHandler(c.userUC)")
}

func TestGenerateGetters(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	generateGetters(&b, []string{"Product"})
	output := b.String()
	assert.Contains(t, output, "func (c *Container) ProductHandler()")
	assert.Contains(t, output, "func (c *Container) ProductUseCase()")
	assert.Contains(t, output, "func (c *Container) ProductRepository()")
	assert.Contains(t, output, "return c.productHandler")
	assert.Contains(t, output, "return c.productUC")
	assert.Contains(t, output, "return c.productRepo")
}

func TestWriteWireHeader(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	writeWireHeader(&b)
	output := b.String()
	assert.Contains(t, output, "//go:build wireinject")
	assert.Contains(t, output, "package di")
}

func TestWriteWireImports(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	writeWireImports(&b, "mymodule")
	output := b.String()
	assert.Contains(t, output, "import (")
	assert.Contains(t, output, "database/sql")
	assert.Contains(t, output, "github.com/google/wire")
	assert.Contains(t, output, "mymodule/internal/repository")
	assert.Contains(t, output, "mymodule/internal/usecase")
}

func TestWriteWireSets(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	writeWireSets(&b, []string{"Product"}, "postgres")
	output := b.String()
	assert.Contains(t, output, "Wire sets")
	assert.Contains(t, output, "RepositorySet")
	assert.Contains(t, output, "UseCaseSet")
	assert.Contains(t, output, "HandlerSet")
	assert.Contains(t, output, "AllSet")
}

func TestWriteRepositorySet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		database string
		expected string
	}{
		{"postgres", "postgres", "NewPostgresProductRepository"},
		{"mysql", "mysql", "NewMySQLProductRepository"},
		{"mongodb", "mongodb", "NewMongoProductRepository"},
		{"default", "other", "NewPostgresProductRepository"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var b strings.Builder
			writeRepositorySet(&b, []string{"Product"}, tc.database)
			assert.Contains(t, b.String(), tc.expected)
		})
	}
}

func TestWriteUseCaseSet(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	writeUseCaseSet(&b, []string{"Product", "User"})
	output := b.String()
	assert.Contains(t, output, "usecase.NewProductService")
	assert.Contains(t, output, "usecase.NewUserService")
}

func TestWriteHandlerSet(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	writeHandlerSet(&b, []string{"Product"})
	assert.Contains(t, b.String(), "http.NewProductHandler")
}

func TestWriteAllSet(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	writeAllSet(&b)
	output := b.String()
	assert.Contains(t, output, "RepositorySet")
	assert.Contains(t, output, "UseCaseSet")
	assert.Contains(t, output, "HandlerSet")
}

func TestWriteWireFunctions(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	writeWireFunctions(&b, []string{"Product"})
	output := b.String()
	assert.Contains(t, output, "func InitializeProductHandler(db *sql.DB)")
	assert.Contains(t, output, "wire.Build(AllSet)")
	assert.Contains(t, output, "func InitializeContainer(db *sql.DB)")
}
