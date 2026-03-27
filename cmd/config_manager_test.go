package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigManager(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()
	require.NotNil(t, cm)
	assert.Nil(t, cm.config)
	assert.Empty(t, cm.GetErrors())
	assert.Empty(t, cm.GetWarnings())
}

func TestConfigManager_CreateDefaultConfig(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()
	config := cm.CreateDefaultConfig(filepath.Join(t.TempDir(), "test-project"))

	require.NotNil(t, config)
	assert.Equal(t, "test-project", config.Project.Name)
	assert.Contains(t, config.Project.Module, "test-project")
	assert.Equal(t, "postgres", config.Database.Type)
	assert.Equal(t, 5432, config.Database.Port)
	assert.True(t, config.Architecture.Layers.Domain.Enabled)
	assert.True(t, config.Testing.Enabled)
}

func TestConfigManager_CreateDefaultConfig_DotPath(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()
	config := cm.CreateDefaultConfig(".")

	require.NotNil(t, config)
	assert.Equal(t, "goca-project", config.Project.Name)
}

func TestConfigManager_ApplyDefaults(t *testing.T) {
	t.Parallel()

	t.Run("empty config gets postgres defaults", func(t *testing.T) {
		t.Parallel()
		config := &GocaConfig{}
		cm2 := NewConfigManager()
		cm2.ApplyDefaults(config)
		assert.Equal(t, "postgres", config.Database.Type)
		assert.Equal(t, 5432, config.Database.Port)
	})

	t.Run("mysql gets port 3306", func(t *testing.T) {
		t.Parallel()
		config := &GocaConfig{}
		config.Database.Type = "mysql"
		cm2 := NewConfigManager()
		cm2.ApplyDefaults(config)
		assert.Equal(t, 3306, config.Database.Port)
	})

	t.Run("mongodb gets port 27017", func(t *testing.T) {
		t.Parallel()
		config := &GocaConfig{}
		config.Database.Type = "mongodb"
		cm2 := NewConfigManager()
		cm2.ApplyDefaults(config)
		assert.Equal(t, 27017, config.Database.Port)
	})

	t.Run("naming defaults", func(t *testing.T) {
		t.Parallel()
		config := &GocaConfig{}
		cm2 := NewConfigManager()
		cm2.ApplyDefaults(config)
		assert.Equal(t, "PascalCase", config.Architecture.Naming.Entities)
		assert.Equal(t, "camelCase", config.Architecture.Naming.Fields)
		assert.Equal(t, "snake_case", config.Architecture.Naming.Files)
		assert.Equal(t, "lowercase", config.Architecture.Naming.Packages)
	})

	t.Run("DI default", func(t *testing.T) {
		t.Parallel()
		config := &GocaConfig{}
		cm2 := NewConfigManager()
		cm2.ApplyDefaults(config)
		assert.Equal(t, "manual", config.Architecture.DI.Type)
	})

	t.Run("validation library default", func(t *testing.T) {
		t.Parallel()
		config := &GocaConfig{}
		cm2 := NewConfigManager()
		cm2.ApplyDefaults(config)
		assert.Equal(t, "builtin", config.Generation.Validation.Library)
	})

	t.Run("testing framework default", func(t *testing.T) {
		t.Parallel()
		config := &GocaConfig{}
		cm2 := NewConfigManager()
		cm2.ApplyDefaults(config)
		assert.Equal(t, "testify", config.Testing.Framework)
	})

	t.Run("connection defaults", func(t *testing.T) {
		t.Parallel()
		config := &GocaConfig{}
		cm2 := NewConfigManager()
		cm2.ApplyDefaults(config)
		assert.Equal(t, 25, config.Database.Connection.MaxOpen)
	})
}

func TestConfigManager_ValidateConfig(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()
	config := cm.CreateDefaultConfig(t.TempDir())

	err := cm.ValidateConfig(config)
	assert.NoError(t, err)
	assert.Empty(t, cm.GetErrors())
}

func TestConfigManager_ValidateConfig_Invalid(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()
	config := &GocaConfig{}
	// Missing project name and module
	config.Database.Type = "invalid-db"
	config.Database.Port = -1
	config.Architecture.Naming.Entities = "INVALID"
	config.Architecture.DI.Type = "invalid"
	config.Generation.Validation.Library = "invalid"
	config.Testing.Framework = "invalid"

	err := cm.ValidateConfig(config)
	require.Error(t, err)
	errors := cm.GetErrors()
	assert.NotEmpty(t, errors)
}

func TestConfigManager_ValidateProject(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()
		cm2 := NewConfigManager()
		config := &GocaConfig{}
		cm2.ValidateConfig(config)
		errors := cm2.GetErrors()
		found := false
		for _, e := range errors {
			if e.Field == "project.name" {
				found = true
				break
			}
		}
		assert.True(t, found, "should have project.name error")
	})

	t.Run("empty version triggers warning", func(t *testing.T) {
		t.Parallel()
		cm2 := NewConfigManager()
		config := cm.CreateDefaultConfig(t.TempDir())
		config.Project.Version = ""
		cm2.ValidateConfig(config)
		warnings := cm2.GetWarnings()
		found := false
		for _, w := range warnings {
			if w.Field == "project.version" {
				found = true
				break
			}
		}
		assert.True(t, found, "should have project.version warning")
	})
}

func TestConfigManager_FindConfigFile(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()

	t.Run("no config file", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		result := cm.FindConfigFile(dir)
		assert.Empty(t, result)
	})

	t.Run("finds .goca.yaml", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		configPath := filepath.Join(dir, ".goca.yaml")
		require.NoError(t, os.WriteFile(configPath, []byte("project:\n  name: test\n  module: test\n"), 0644))

		result := cm.FindConfigFile(dir)
		assert.Equal(t, configPath, result)
	})

	t.Run("finds .goca.yml", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		configPath := filepath.Join(dir, ".goca.yml")
		require.NoError(t, os.WriteFile(configPath, []byte("project:\n  name: test\n"), 0644))

		result := cm.FindConfigFile(dir)
		assert.Equal(t, configPath, result)
	})
}

func TestConfigManager_MergeWithFlags(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()
	config := cm.CreateDefaultConfig(t.TempDir())
	cm.SetConfig(config)

	cm.MergeWithFlags(map[string]interface{}{
		"database":   "mysql",
		"validation": true,
		"auth":       true,
	})

	assert.Equal(t, "mysql", cm.GetConfig().Database.Type)
	assert.True(t, cm.GetConfig().Generation.Validation.Enabled)
	assert.True(t, cm.GetConfig().Features.Auth.Enabled)
}

func TestConfigManager_MergeWithFlags_NilConfig(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()
	// Should not panic
	cm.MergeWithFlags(map[string]interface{}{"database": "mysql"})
}

func TestConfigManager_SaveAndLoadConfig(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	// Create and save
	cm := NewConfigManager()
	config := cm.CreateDefaultConfig(dir)
	cm.SetConfig(config)

	configPath := filepath.Join(dir, ".goca.yaml")
	err := cm.SaveConfig(configPath)
	require.NoError(t, err)

	// Load
	cm2 := NewConfigManager()
	err = cm2.LoadFromFile(configPath)
	require.NoError(t, err)

	loaded := cm2.GetConfig()
	require.NotNil(t, loaded)
	assert.Equal(t, config.Project.Name, loaded.Project.Name)
	assert.Equal(t, config.Database.Type, loaded.Database.Type)
}

func TestConfigManager_SaveConfig_NilConfig(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()
	err := cm.SaveConfig(filepath.Join(t.TempDir(), "test.yaml"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no configuration loaded")
}

func TestConfigManager_LoadConfig_NoFile(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()
	dir := t.TempDir()
	err := cm.LoadConfig(dir)
	require.NoError(t, err)
	// Should have default config
	assert.NotNil(t, cm.GetConfig())
}

func TestConfigManager_LoadFromFile_InvalidYAML(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	configPath := filepath.Join(dir, "bad.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(":::invalid yaml:::"), 0644))

	cm := NewConfigManager()
	err := cm.LoadFromFile(configPath)
	require.Error(t, err)
}

func TestConfigManager_Contains(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()
	assert.True(t, cm.contains([]string{"a", "b", "c"}, "b"))
	assert.False(t, cm.contains([]string{"a", "b", "c"}, "d"))
	assert.False(t, cm.contains([]string{}, "a"))
}

func TestConfigManager_IsValidNaming(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()
	valid := []string{"PascalCase", "camelCase", "snake_case"}
	assert.True(t, cm.isValidNaming("PascalCase", valid))
	assert.False(t, cm.isValidNaming("INVALID", valid))
}

func TestConfigManager_CountEnabledLayers(t *testing.T) {
	t.Parallel()
	cm := NewConfigManager()
	config := cm.CreateDefaultConfig(t.TempDir())
	cm.SetConfig(config)

	assert.Equal(t, 4, cm.countEnabledLayers())

	config.Architecture.Layers.Handler.Enabled = false
	assert.Equal(t, 3, cm.countEnabledLayers())
}

func TestConfigManager_GenerateDefaultConfig(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cm := NewConfigManager()

	err := cm.GenerateDefaultConfig(dir, "myproject", "github.com/test/myproject", "mysql")
	require.NoError(t, err)

	configPath := filepath.Join(dir, ".goca.yaml")
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	config := cm.GetConfig()
	assert.Equal(t, "myproject", config.Project.Name)
	assert.Equal(t, "github.com/test/myproject", config.Project.Module)
	assert.Equal(t, "mysql", config.Database.Type)
}
