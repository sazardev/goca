package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for template_generator.go functions NOT in template_generator_test.go

func TestTemplateGenerator_GenerateFromTemplate_Entity(t *testing.T) {
	t.Parallel()
	gen := NewTemplateGenerator()
	data := &TemplateData{
		Entity: EntityData{
			Name:       "Product",
			NameLower:  "product",
			NamePlural: "Products",
			Package:    "domain",
		},
		Fields: []FieldData{
			{Name: "Name", Type: "string", JSONTag: `json:"name"`, IsRequired: true},
		},
		Module:   "myproject",
		Imports:  []string{"errors"},
		Features: FeatureFlags{Validation: true},
	}
	result, err := gen.GenerateFromTemplate("entity", data)
	require.NoError(t, err)
	assert.Contains(t, result, "Product")
}

func TestTemplateGenerator_GenerateFromTemplate_Unknown(t *testing.T) {
	t.Parallel()
	gen := NewTemplateGenerator()
	_, err := gen.GenerateFromTemplate("nonexistent", &TemplateData{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template not found")
}

func TestTemplateGenerator_PrepareTemplateData_Valid(t *testing.T) {
	t.Parallel()
	gen := NewTemplateGenerator()
	data, err := gen.PrepareTemplateData("Product", "Name:string,Price:float64", FeatureFlags{Validation: true})
	require.NoError(t, err)
	assert.Equal(t, "Product", data.Entity.Name)
	assert.Equal(t, "product", data.Entity.NameLower)
	assert.Equal(t, "Products", data.Entity.NamePlural)
	assert.GreaterOrEqual(t, len(data.Fields), 2) // May include auto-added ID field
	// Find the Name field and verify it's marked as required
	var foundName bool
	for _, f := range data.Fields {
		if f.Name == "Name" {
			foundName = true
			assert.True(t, f.IsRequired)
		}
	}
	assert.True(t, foundName)
}

func TestTemplateGenerator_PrepareTemplateData_EmptyFields(t *testing.T) {
	t.Parallel()
	gen := NewTemplateGenerator()
	_, err := gen.PrepareTemplateData("Product", "", FeatureFlags{})
	assert.Error(t, err)
}

func TestTemplateGenerator_GenerateFromTemplate_AllTemplates(t *testing.T) {
	t.Parallel()
	gen := NewTemplateGenerator()
	data := &TemplateData{
		Entity: EntityData{
			Name:       "Order",
			NameLower:  "order",
			NamePlural: "Orders",
			Package:    "domain",
		},
		Module: "myproject",
	}
	for _, name := range []string{"entity", "usecase", "repository", "handler"} {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result, err := gen.GenerateFromTemplate(name, data)
			require.NoError(t, err)
			assert.NotEmpty(t, result)
		})
	}
}

func TestTemplateGenerator_GenerateValidationBody(t *testing.T) {
	t.Parallel()
	gen := NewTemplateGenerator()
	result := gen.generateValidationBody(nil)
	assert.NotEmpty(t, result)
}

func TestTemplateGenerator_GenerateBusinessRuleMethods(t *testing.T) {
	t.Parallel()
	gen := NewTemplateGenerator()
	entity := EntityData{Name: "Product"}
	result := gen.generateBusinessRuleMethods(entity, nil)
	assert.Empty(t, result) // Stub returns empty
}
