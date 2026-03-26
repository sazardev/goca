package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCreateHandlerMethod(t *testing.T) {
	t.Parallel()

	t.Run("without validation", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateCreateHandlerMethod(&b, "Product", "ProductHandler", false)
		output := b.String()
		assert.Contains(t, output, "func (p *ProductHandler) CreateProduct(")
		assert.Contains(t, output, "CreateProductInput")
		assert.Contains(t, output, "json.NewDecoder")
		assert.Contains(t, output, "StatusCreated")
		assert.NotContains(t, output, "validator.New()")
	})

	t.Run("with validation", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateCreateHandlerMethod(&b, "Product", "ProductHandler", true)
		output := b.String()
		assert.Contains(t, output, "validator.New().Struct(input)")
		assert.Contains(t, output, "StatusUnprocessableEntity")
	})
}

func TestGenerateGetHandlerMethod(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	generateGetHandlerMethod(&b, "Product", "ProductHandler")
	output := b.String()
	assert.Contains(t, output, "func (p *ProductHandler) GetProduct(")
	assert.Contains(t, output, "mux.Vars(r)")
	assert.Contains(t, output, "strconv.Atoi")
	assert.Contains(t, output, "Invalid product ID")
	assert.Contains(t, output, "StatusNotFound")
}

func TestGenerateUpdateHandlerMethod(t *testing.T) {
	t.Parallel()

	t.Run("without validation", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateUpdateHandlerMethod(&b, "Product", "ProductHandler", false)
		output := b.String()
		assert.Contains(t, output, "func (p *ProductHandler) UpdateProduct(")
		assert.Contains(t, output, "UpdateProductInput")
		assert.Contains(t, output, "StatusNoContent")
		assert.NotContains(t, output, "validator.New()")
	})

	t.Run("with validation", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateUpdateHandlerMethod(&b, "Product", "ProductHandler", true)
		output := b.String()
		assert.Contains(t, output, "validator.New().Struct(input)")
	})
}

func TestGenerateDeleteHandlerMethod(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	generateDeleteHandlerMethod(&b, "Product", "ProductHandler")
	output := b.String()
	assert.Contains(t, output, "func (p *ProductHandler) DeleteProduct(")
	assert.Contains(t, output, "mux.Vars(r)")
	assert.Contains(t, output, "StatusNoContent")
}

func TestGenerateListHandlerMethod(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	generateListHandlerMethod(&b, "Product", "ProductHandler")
	output := b.String()
	assert.Contains(t, output, "func (p *ProductHandler) ListProducts(")
	assert.Contains(t, output, "ListProducts()")
	assert.Contains(t, output, "json.NewEncoder(w).Encode(output)")
}
