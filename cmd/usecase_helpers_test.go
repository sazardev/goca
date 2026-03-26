package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOperations(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty returns defaults", "", []string{"create", "read"}},
		{"single op", "create", []string{"create"}},
		{"multiple ops", "create,update,delete", []string{"create", "update", "delete"}},
		{"with spaces", "create , update , delete", []string{"create", "update", "delete"}},
		{"all ops", "create,read,update,delete,list", []string{"create", "read", "update", "delete", "list"}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := parseOperations(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetValidationTag(t *testing.T) {
	t.Parallel()

	cases := []struct {
		fieldType string
		expected  string
	}{
		{"string", "required,min=1"},
		{"int", "required,min=1"},
		{"int64", "required,min=1"},
		{"uint", "required,min=1"},
		{"uint64", "required,min=1"},
		{"float64", "required,min=0"},
		{"float32", "required,min=0"},
		{"bool", ""},
		{"time.Time", "required"},
		{"custom", "required"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.fieldType, func(t *testing.T) {
			t.Parallel()
			result := getValidationTag(tc.fieldType)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestEnsureImportInDTOFile(t *testing.T) {
	t.Parallel()

	t.Run("import already exists", func(t *testing.T) {
		t.Parallel()
		content := "package usecase\n\nimport (\n\t\"errors\"\n)\n"
		result := ensureImportInDTOFile(content, "errors", "mymodule")
		assert.Equal(t, content, result)
	})

	t.Run("adds missing import", func(t *testing.T) {
		t.Parallel()
		content := "package usecase\n\nimport (\n\t\"fmt\"\n)\n"
		result := ensureImportInDTOFile(content, "strings", "mymodule")
		assert.Contains(t, result, "\"strings\"")
		assert.Contains(t, result, "\"fmt\"")
	})

	t.Run("no import block", func(t *testing.T) {
		t.Parallel()
		content := "package usecase\n\ntype Foo struct{}\n"
		result := ensureImportInDTOFile(content, "errors", "mymodule")
		// No import block, returns unchanged
		assert.Equal(t, content, result)
	})
}

func TestGenerateCreateDTO(t *testing.T) {
	t.Parallel()

	t.Run("without validation", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateCreateDTO(&b, "Product", false)
		output := b.String()
		assert.Contains(t, output, "type CreateProductInput struct")
		assert.Contains(t, output, "type CreateProductOutput struct")
		assert.Contains(t, output, "json:\"nombre\"")
		assert.NotContains(t, output, "validate:")
	})

	t.Run("with validation", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateCreateDTO(&b, "Product", true)
		output := b.String()
		assert.Contains(t, output, "type CreateProductInput struct")
		assert.Contains(t, output, "validate:\"required,min=2\"")
	})
}

func TestGenerateUpdateDTO(t *testing.T) {
	t.Parallel()

	t.Run("without validation", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateUpdateDTO(&b, "Product", false)
		output := b.String()
		assert.Contains(t, output, "type UpdateProductInput struct")
		assert.Contains(t, output, "json:\"nombre,omitempty\"")
		assert.NotContains(t, output, "validate:")
	})

	t.Run("with validation", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateUpdateDTO(&b, "Product", true)
		output := b.String()
		assert.Contains(t, output, "validate:\"omitempty,min=2\"")
	})
}

func TestGenerateListDTO(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	generateListDTO(&b, "Product")
	output := b.String()
	assert.Contains(t, output, "type ListProductOutput struct")
	assert.Contains(t, output, "[]domain.Product")
	assert.Contains(t, output, "json:\"products\"")
	assert.Contains(t, output, "json:\"total\"")
	assert.Contains(t, output, "json:\"message\"")
}
