package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProjectTemplates(t *testing.T) {
	t.Parallel()
	templates := GetProjectTemplates()
	assert.NotEmpty(t, templates)
	assert.Contains(t, templates, "minimal")
	assert.Contains(t, templates, "rest-api")
}

func TestGetTemplateNames(t *testing.T) {
	t.Parallel()
	names := GetTemplateNames()
	assert.NotEmpty(t, names)
	assert.Contains(t, names, "minimal")
}

func TestValidateTemplateName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty is valid", "", true},
		{"minimal", "minimal", true},
		{"rest-api", "rest-api", true},
		{"unknown", "nonexistent", false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, ValidateTemplateName(tc.input))
		})
	}
}

func TestGetTemplateConfig(t *testing.T) {
	t.Parallel()

	t.Run("valid template", func(t *testing.T) {
		t.Parallel()
		config, err := GetTemplateConfig("minimal")
		require.NoError(t, err)
		assert.NotEmpty(t, config)
	})

	t.Run("invalid template", func(t *testing.T) {
		t.Parallel()
		_, err := GetTemplateConfig("nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
