package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateDefaultTemplate(t *testing.T) {
	t.Parallel()

	t.Run("with postgres", func(t *testing.T) {
		t.Parallel()
		result := generateDefaultTemplate("postgres", []string{"http"})
		assert.Contains(t, result, "project:")
		assert.Contains(t, result, "database:")
		assert.Contains(t, result, `type: "postgres"`)
		assert.Contains(t, result, "defaults:")
		assert.Contains(t, result, "architecture:")
		assert.Contains(t, result, "generation:")
		assert.Contains(t, result, "quality:")
		assert.Contains(t, result, "infrastructure:")
	})

	t.Run("with mysql", func(t *testing.T) {
		t.Parallel()
		result := generateDefaultTemplate("mysql", []string{"grpc"})
		assert.Contains(t, result, `type: "mysql"`)
		assert.Contains(t, result, "grpc")
	})

	t.Run("empty database defaults to postgres", func(t *testing.T) {
		t.Parallel()
		result := generateDefaultTemplate("", nil)
		assert.Contains(t, result, `type: "postgres"`)
		assert.Contains(t, result, "http")
	})

	t.Run("multiple handlers", func(t *testing.T) {
		t.Parallel()
		result := generateDefaultTemplate("postgres", []string{"http", "grpc"})
		assert.Contains(t, result, "http, grpc")
	})
}

func TestGenerateWebTemplate(t *testing.T) {
	t.Parallel()
	result := generateWebTemplate("postgres", []string{"http"})
	assert.Contains(t, result, "project:")
	assert.NotEmpty(t, result)
}

func TestGenerateAPITemplate(t *testing.T) {
	t.Parallel()
	result := generateAPITemplate("mysql", []string{"grpc"})
	assert.Contains(t, result, "project:")
	assert.NotEmpty(t, result)
}

func TestGenerateMicroserviceTemplate(t *testing.T) {
	t.Parallel()
	result := generateMicroserviceTemplate("postgres", []string{"http"})
	assert.Contains(t, result, "project:")
	assert.NotEmpty(t, result)
}

func TestGenerateFullTemplate(t *testing.T) {
	t.Parallel()
	result := generateFullTemplate("postgres", []string{"http", "grpc"})
	assert.Contains(t, result, "project:")
	assert.NotEmpty(t, result)
}

func TestValidateConfigStructure(t *testing.T) {
	t.Parallel()

	t.Run("empty config has errors", func(t *testing.T) {
		t.Parallel()
		errs := validateConfigStructure(map[string]interface{}{})
		assert.Len(t, errs, 2) // missing project and defaults
		assert.Contains(t, strings.Join(errs, " "), "project")
		assert.Contains(t, strings.Join(errs, " "), "defaults")
	})

	t.Run("valid config no errors", func(t *testing.T) {
		t.Parallel()
		config := map[string]interface{}{
			"project":  map[string]interface{}{"name": "test"},
			"defaults": map[string]interface{}{"database": "postgres"},
		}
		errs := validateConfigStructure(config)
		assert.Empty(t, errs)
	})

	t.Run("missing project only", func(t *testing.T) {
		t.Parallel()
		config := map[string]interface{}{
			"defaults": map[string]interface{}{"database": "postgres"},
		}
		errs := validateConfigStructure(config)
		assert.Len(t, errs, 1)
		assert.Contains(t, errs[0], "project")
	})

	t.Run("missing defaults only", func(t *testing.T) {
		t.Parallel()
		config := map[string]interface{}{
			"project": map[string]interface{}{"name": "test"},
		}
		errs := validateConfigStructure(config)
		assert.Len(t, errs, 1)
		assert.Contains(t, errs[0], "defaults")
	})

	t.Run("extra sections are ok", func(t *testing.T) {
		t.Parallel()
		config := map[string]interface{}{
			"project":  map[string]interface{}{"name": "test"},
			"defaults": map[string]interface{}{"database": "postgres"},
			"extras":   "something",
		}
		errs := validateConfigStructure(config)
		assert.Empty(t, errs)
	})
}
