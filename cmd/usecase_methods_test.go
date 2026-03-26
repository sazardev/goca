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
	generateCreateMethodWithFields(&sb, "productService", "Product", "Name:string,Price:float64")
	result := sb.String()
	assert.Contains(t, result, "func (p *productService) CreateProduct")
	assert.Contains(t, result, "Name")
}
