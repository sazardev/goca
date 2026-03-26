package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigIntegration(t *testing.T) {
	t.Parallel()
	ci := NewConfigIntegration()
	require.NotNil(t, ci)
	assert.NotNil(t, ci.manager)
}

func TestConfigIntegration_GetDatabaseType(t *testing.T) {
	t.Parallel()
	ci := NewConfigIntegration()

	t.Run("CLI flag takes precedence", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "mysql", ci.GetDatabaseType("mysql"))
	})

	t.Run("config value used when no flag", func(t *testing.T) {
		t.Parallel()
		ci2 := NewConfigIntegration()
		ci2.config = &GocaConfig{Database: DatabaseConfig{Type: "sqlite"}}
		assert.Equal(t, "sqlite", ci2.GetDatabaseType(""))
	})

	t.Run("default when no flag or config", func(t *testing.T) {
		t.Parallel()
		ci2 := NewConfigIntegration()
		assert.Equal(t, "postgres", ci2.GetDatabaseType(""))
	})
}

func TestConfigIntegration_GetHandlerTypes(t *testing.T) {
	t.Parallel()
	ci := NewConfigIntegration()

	t.Run("CLI flag", func(t *testing.T) {
		t.Parallel()
		handlers := ci.GetHandlerTypes("http,grpc")
		assert.Equal(t, []string{"http", "grpc"}, handlers)
	})

	t.Run("config value", func(t *testing.T) {
		t.Parallel()
		ci2 := NewConfigIntegration()
		ci2.config = &GocaConfig{}
		ci2.config.Architecture.Layers.Handler.Enabled = true
		handlers := ci2.GetHandlerTypes("")
		assert.Contains(t, handlers, "http")
	})

	t.Run("default", func(t *testing.T) {
		t.Parallel()
		ci2 := NewConfigIntegration()
		handlers := ci2.GetHandlerTypes("")
		assert.Equal(t, []string{"http"}, handlers)
	})
}

func TestConfigIntegration_GetValidationEnabled(t *testing.T) {
	t.Parallel()

	t.Run("CLI flag true", func(t *testing.T) {
		t.Parallel()
		ci := NewConfigIntegration()
		val := true
		assert.True(t, ci.GetValidationEnabled(&val))
	})

	t.Run("CLI flag false", func(t *testing.T) {
		t.Parallel()
		ci := NewConfigIntegration()
		val := false
		assert.False(t, ci.GetValidationEnabled(&val))
	})

	t.Run("config value", func(t *testing.T) {
		t.Parallel()
		ci := NewConfigIntegration()
		ci.config = &GocaConfig{}
		ci.config.Generation.Validation.Enabled = false
		assert.False(t, ci.GetValidationEnabled(nil))
	})

	t.Run("default", func(t *testing.T) {
		t.Parallel()
		ci := NewConfigIntegration()
		assert.True(t, ci.GetValidationEnabled(nil))
	})
}

func TestConfigIntegration_GetBusinessRulesEnabled(t *testing.T) {
	t.Parallel()

	t.Run("CLI flag", func(t *testing.T) {
		t.Parallel()
		ci := NewConfigIntegration()
		val := true
		assert.True(t, ci.GetBusinessRulesEnabled(&val))
	})

	t.Run("default", func(t *testing.T) {
		t.Parallel()
		ci := NewConfigIntegration()
		assert.False(t, ci.GetBusinessRulesEnabled(nil))
	})
}

func TestConfigIntegration_GetProjectConfig(t *testing.T) {
	t.Parallel()

	t.Run("with config", func(t *testing.T) {
		t.Parallel()
		ci := NewConfigIntegration()
		ci.config = &GocaConfig{}
		ci.config.Project.Name = "myproject"
		ci.config.Project.Module = "github.com/test/myproject"
		pc := ci.GetProjectConfig()
		assert.Equal(t, "myproject", pc.Name)
	})

	t.Run("without config", func(t *testing.T) {
		t.Parallel()
		ci := NewConfigIntegration()
		pc := ci.GetProjectConfig()
		assert.NotEmpty(t, pc.Name)
		assert.NotEmpty(t, pc.Module)
	})
}

func TestConfigIntegration_GetArchitectureConfig(t *testing.T) {
	t.Parallel()

	t.Run("default", func(t *testing.T) {
		t.Parallel()
		ci := NewConfigIntegration()
		ac := ci.GetArchitectureConfig()
		assert.True(t, ac.Layers.Domain.Enabled)
		assert.Equal(t, "manual", ac.DI.Type)
	})
}

func TestConfigIntegration_GetDatabaseConfig(t *testing.T) {
	t.Parallel()

	t.Run("default", func(t *testing.T) {
		t.Parallel()
		ci := NewConfigIntegration()
		dc := ci.GetDatabaseConfig()
		assert.Equal(t, "postgres", dc.Type)
		assert.Equal(t, 5432, dc.Port)
	})
}

func TestConfigIntegration_GetGenerationConfig(t *testing.T) {
	t.Parallel()

	ci := NewConfigIntegration()
	gc := ci.GetGenerationConfig()
	assert.True(t, gc.Validation.Enabled)
	assert.Equal(t, "builtin", gc.Validation.Library)
}

func TestConfigIntegration_GetTestingConfig(t *testing.T) {
	t.Parallel()

	ci := NewConfigIntegration()
	tc := ci.GetTestingConfig()
	assert.True(t, tc.Enabled)
	assert.Equal(t, "testify", tc.Framework)
}

func TestConfigIntegration_GetTemplateConfig(t *testing.T) {
	t.Parallel()

	ci := NewConfigIntegration()
	tc := ci.GetTemplateConfig()
	assert.Equal(t, ".goca/templates", tc.Directory)
}

func TestConfigIntegration_GetFeatureConfig(t *testing.T) {
	t.Parallel()

	ci := NewConfigIntegration()
	fc := ci.GetFeatureConfig()
	assert.False(t, fc.Auth.Enabled)
	assert.True(t, fc.Logging.Enabled)
}

func TestConfigIntegration_HasConfigFile(t *testing.T) {
	t.Parallel()

	t.Run("no config", func(t *testing.T) {
		t.Parallel()
		ci := NewConfigIntegration()
		assert.False(t, ci.HasConfigFile())
	})

	t.Run("with config but no file path", func(t *testing.T) {
		t.Parallel()
		ci := NewConfigIntegration()
		ci.config = &GocaConfig{}
		assert.False(t, ci.HasConfigFile())
	})
}

func TestConfigIntegration_GetConfigPath(t *testing.T) {
	t.Parallel()

	ci := NewConfigIntegration()
	assert.Empty(t, ci.GetConfigPath())
}

func TestConfigIntegration_ValidateConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("nil manager", func(t *testing.T) {
		t.Parallel()
		ci := &ConfigIntegration{}
		assert.NoError(t, ci.ValidateConfiguration())
	})

	t.Run("no errors", func(t *testing.T) {
		t.Parallel()
		ci := NewConfigIntegration()
		assert.NoError(t, ci.ValidateConfiguration())
	})
}

func TestConfigIntegration_GetValidationErrors(t *testing.T) {
	t.Parallel()
	ci := &ConfigIntegration{}
	assert.Nil(t, ci.GetValidationErrors())
}

func TestConfigIntegration_GetValidationWarnings(t *testing.T) {
	t.Parallel()
	ci := &ConfigIntegration{}
	assert.Nil(t, ci.GetValidationWarnings())
}

func TestConfigIntegration_MergeWithCLIFlags(t *testing.T) {
	t.Parallel()
	ci := NewConfigIntegration()
	ci.config = ci.manager.CreateDefaultConfig("/tmp/test")
	ci.manager.SetConfig(ci.config)

	ci.MergeWithCLIFlags(map[string]interface{}{
		"database": "mysql",
	})

	assert.Equal(t, "mysql", ci.config.Database.Type)
}
