package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNamingConvention_NilConfig(t *testing.T) {
	t.Parallel()
	ci := &ConfigIntegration{}

	cases := []struct {
		element  string
		expected string
	}{
		{"entity", "PascalCase"},
		{"field", "PascalCase"},
		{"file", "snake_case"},
		{"package", "lowercase"},
		{"unknown", "PascalCase"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.element, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, ci.GetNamingConvention(tc.element))
		})
	}
}

func TestGetNamingConvention_WithConfig(t *testing.T) {
	t.Parallel()
	ci := &ConfigIntegration{
		config: &GocaConfig{},
	}
	ci.config.Architecture.Naming.Entities = "myentity"
	ci.config.Architecture.Naming.Fields = "myfield"
	ci.config.Architecture.Naming.Files = "myfile"
	ci.config.Architecture.Naming.Packages = "mypkg"
	ci.config.Architecture.Naming.Constants = "myconst"
	ci.config.Architecture.Naming.Variables = "myvar"
	ci.config.Architecture.Naming.Functions = "myfunc"

	assert.Equal(t, "myentity", ci.GetNamingConvention("entity"))
	assert.Equal(t, "myfield", ci.GetNamingConvention("field"))
	assert.Equal(t, "myfile", ci.GetNamingConvention("file"))
	assert.Equal(t, "mypkg", ci.GetNamingConvention("package"))
	assert.Equal(t, "myconst", ci.GetNamingConvention("constant"))
	assert.Equal(t, "myvar", ci.GetNamingConvention("variable"))
	assert.Equal(t, "myfunc", ci.GetNamingConvention("function"))
	assert.Equal(t, "myentity", ci.GetNamingConvention("unknown"))
}

func TestUpdateConfigAfterGeneration(t *testing.T) {
	t.Parallel()

	t.Run("nil config no panic", func(t *testing.T) {
		t.Parallel()
		ci := &ConfigIntegration{}
		ci.UpdateConfigAfterGeneration("entity", "Product")
	})

	t.Run("with config no panic", func(t *testing.T) {
		t.Parallel()
		ci := &ConfigIntegration{config: &GocaConfig{}}
		ci.UpdateConfigAfterGeneration("entity", "Product")
	})
}

func TestGetTemplateManager(t *testing.T) {
	t.Parallel()
	ci := &ConfigIntegration{}
	assert.Nil(t, ci.GetTemplateManager())
}

func TestHasCustomTemplate(t *testing.T) {
	t.Parallel()
	ci := &ConfigIntegration{}
	assert.False(t, ci.HasCustomTemplate("anything"))
}

func TestExecuteCustomTemplate(t *testing.T) {
	t.Parallel()
	ci := &ConfigIntegration{}
	_, err := ci.ExecuteCustomTemplate("test", nil)
	assert.Error(t, err)
}

func TestGetAvailableTemplates(t *testing.T) {
	t.Parallel()
	ci := &ConfigIntegration{}
	result := ci.GetAvailableTemplates()
	assert.Empty(t, result)
}

func TestInitializeTemplateSystem_NoConfig(t *testing.T) {
	t.Parallel()
	ci := &ConfigIntegration{}
	err := ci.InitializeTemplateSystem()
	assert.Error(t, err)
}

func TestGenerateProjectDocumentation(t *testing.T) {
	t.Parallel()
	ci := &ConfigIntegration{}
	err := ci.GenerateProjectDocumentation()
	assert.NoError(t, err) // nil template manager returns nil
}
