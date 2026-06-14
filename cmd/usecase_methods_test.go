package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for usecase.go functions NOT already in usecase_helpers_test.go

func TestGenerateCreateMethod_Pure(t *testing.T) {
	t.Parallel()
	var sb strings.Builder
	generateCreateMethod(&sb, "productService", "Product")
	result := sb.String()
	assert.Contains(t, result, "func (p *productService) CreateProduct")
	assert.Contains(t, result, "CreateProductInput")
	assert.Contains(t, result, "CreateProductOutput")
	// Generated comments/identifiers must be English, not the old Spanish ones.
	assert.NotContains(t, result, "Nombre")
	assert.NotContains(t, result, "Edad")
}

func TestGenerateGetMethod_Pure(t *testing.T) {
	t.Parallel()
	var sb strings.Builder
	generateGetMethod(&sb, "productService", "Product")
	result := sb.String()
	assert.Contains(t, result, "func")
	assert.Contains(t, result, "productService")
	assert.Contains(t, result, "Product")
}

func TestGenerateDeleteMethod_Pure(t *testing.T) {
	t.Parallel()
	var sb strings.Builder
	generateDeleteMethod(&sb, "productService", "Product")
	result := sb.String()
	assert.Contains(t, result, "func (p *productService) DeleteProduct")
	assert.Contains(t, result, "repo.Delete")
}

func TestGenerateListMethod_Pure(t *testing.T) {
	t.Parallel()
	var sb strings.Builder
	generateListMethod(&sb, "productService", "Product")
	result := sb.String()
	assert.Contains(t, result, "func")
	assert.Contains(t, result, "productService")
	assert.Contains(t, result, "Product")
}

func TestGenerateUpdateMethodWithFields_Pure(t *testing.T) {
	t.Parallel()
	var sb strings.Builder
	generateUpdateMethodWithFields(&sb, "productService", "Product", "Name:string,Price:float64")
	result := sb.String()
	assert.Contains(t, result, "func (p *productService) UpdateProduct")
	assert.Contains(t, result, "Name")
	assert.Contains(t, result, "Price")
}

func TestGenerateCreateMethodWithFields_Pure(t *testing.T) {
	t.Parallel()
	var sb strings.Builder
	generateCreateMethodWithFields(&sb, "productService", "Product", "Name:string,Price:float64", true)
	result := sb.String()
	assert.Contains(t, result, "func (p *productService) CreateProduct")
	assert.Contains(t, result, "Name")
	// When DTO validation is enabled, the create method must call input.Validate().
	assert.Contains(t, result, "input.Validate()")
}

func TestGenerateCreateMethodWithFields_NoDTOValidate(t *testing.T) {
	t.Parallel()
	var sb strings.Builder
	generateCreateMethodWithFields(&sb, "productService", "Product", "Name:string,Price:float64", false)
	result := sb.String()
	assert.Contains(t, result, "func (p *productService) CreateProduct")
	assert.NotContains(t, result, "input.Validate()")
}

func TestGenerateAsyncCreateMethod(t *testing.T) {
	t.Parallel()
	var sb strings.Builder
	generateAsyncCreateMethod(&sb, "productService", "Product")
	result := sb.String()
	assert.Contains(t, result, "func (p *productService) CreateProductAsync(input CreateProductInput)")
	assert.Contains(t, result, "p.asyncChannel <- func() {")
	assert.Contains(t, result, "p.CreateProduct(input)")
}
