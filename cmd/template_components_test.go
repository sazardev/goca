package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemplateBuilder(t *testing.T) {
	t.Parallel()
	tb := NewTemplateBuilder()
	require.NotNil(t, tb)
	assert.Empty(t, tb.components)
}

func TestTemplateBuilder_AddComponent(t *testing.T) {
	t.Parallel()
	tb := NewTemplateBuilder()
	component := TemplateComponent{Name: "test", Template: "hello", Required: true}

	result := tb.AddComponent(component)
	assert.Same(t, tb, result, "should return self for chaining")
	assert.Len(t, tb.components, 1)
	assert.Equal(t, "test", tb.components[0].Name)
}

func TestTemplateBuilder_AddComponentByName(t *testing.T) {
	t.Parallel()
	tb := NewTemplateBuilder()

	t.Run("existing component", func(t *testing.T) {
		tb2 := NewTemplateBuilder()
		tb2.AddComponentByName("header", EntityTemplateComponents)
		assert.Len(t, tb2.components, 1)
	})

	t.Run("nonexistent component", func(t *testing.T) {
		tb.AddComponentByName("nonexistent", EntityTemplateComponents)
		assert.Empty(t, tb.components)
	})
}

func TestTemplateBuilder_Build(t *testing.T) {
	t.Parallel()
	tb := NewTemplateBuilder()
	tb.AddComponent(TemplateComponent{Template: "AAA"})
	tb.AddComponent(TemplateComponent{Template: "BBB"})

	result := tb.Build()
	assert.Equal(t, "AAABBB", result)
}

func TestBuildTemplate(t *testing.T) {
	t.Parallel()

	t.Run("entity header+struct+close", func(t *testing.T) {
		t.Parallel()
		result := BuildTemplate([]string{"header", "struct", "structClose"}, EntityTemplateComponents)
		assert.Contains(t, result, "package domain")
		assert.Contains(t, result, "type {{.Entity.Name}} struct")
	})

	t.Run("empty components", func(t *testing.T) {
		t.Parallel()
		result := BuildTemplate([]string{}, EntityTemplateComponents)
		assert.Empty(t, result)
	})
}

func TestGetEntityTemplate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name             string
		timestamps       bool
		softDelete       bool
		validation       bool
		methods          bool
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name:          "minimal",
			shouldContain: []string{"package domain", "type {{.Entity.Name}} struct"},
		},
		{
			name:          "with timestamps",
			timestamps:    true,
			shouldContain: []string{"CreatedAt", "UpdatedAt"},
		},
		{
			name:          "with soft delete",
			softDelete:    true,
			shouldContain: []string{"DeletedAt"},
		},
		{
			name:          "with validation",
			validation:    true,
			shouldContain: []string{"Validate()"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := GetEntityTemplate(tc.timestamps, tc.softDelete, tc.validation, tc.methods)
			for _, s := range tc.shouldContain {
				assert.Contains(t, result, s)
			}
			for _, s := range tc.shouldNotContain {
				assert.NotContains(t, result, s)
			}
		})
	}
}

func TestGetUseCaseTemplate(t *testing.T) {
	t.Parallel()
	result := GetUseCaseTemplate()
	assert.Contains(t, result, "package usecase")
	assert.Contains(t, result, "UseCase interface")
	assert.Contains(t, result, "Service")
}

func TestValidateTemplate(t *testing.T) {
	t.Parallel()

	t.Run("valid template", func(t *testing.T) {
		t.Parallel()
		err := ValidateTemplate("Hello {{.Name}}")
		assert.NoError(t, err)
	})

	t.Run("invalid template", func(t *testing.T) {
		t.Parallel()
		err := ValidateTemplate("Hello {{.Name")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "template validation failed")
	})
}

func TestEntityTemplateComponents_AllRequired(t *testing.T) {
	t.Parallel()
	requiredNames := []string{"header", "struct", "structClose"}
	for _, name := range requiredNames {
		comp, ok := EntityTemplateComponents[name]
		assert.True(t, ok, "component %s should exist", name)
		assert.True(t, comp.Required, "component %s should be required", name)
	}
}

func TestUseCaseTemplateComponents_AllRequired(t *testing.T) {
	t.Parallel()
	requiredNames := []string{"header", "interface", "service", "dtos"}
	for _, name := range requiredNames {
		comp, ok := UseCaseTemplateComponents[name]
		assert.True(t, ok, "component %s should exist", name)
		assert.True(t, comp.Required, "component %s should be required", name)
	}
}
