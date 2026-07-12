package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDependencyManager(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	dm := NewDependencyManager(root, true)
	assert.NotNil(t, dm)
	assert.Equal(t, root, dm.projectRoot)
	assert.True(t, dm.dryRun)
	assert.Contains(t, dm.goModPath, "go.mod")
}

func TestCommonDependencies(t *testing.T) {
	t.Parallel()
	dm := NewDependencyManager(t.TempDir(), false)
	deps := dm.CommonDependencies()

	assert.NotEmpty(t, deps)
	assert.Contains(t, deps, "validator")
	assert.Contains(t, deps, "jwt")
	assert.Contains(t, deps, "cors")
	assert.Contains(t, deps, "uuid")
	assert.Contains(t, deps, "bcrypt")
	assert.Contains(t, deps, "testify")
	assert.Contains(t, deps, "grpc")
	assert.Contains(t, deps, "protobuf")

	// Verify structure
	v := deps["validator"]
	assert.NotEmpty(t, v.Module)
	assert.NotEmpty(t, v.Version)
	assert.NotEmpty(t, v.Type)
	assert.NotEmpty(t, v.Reason)
}

func TestSuggestDependencies(t *testing.T) {
	t.Parallel()
	dm := NewDependencyManager(t.TempDir(), false)

	t.Run("validation feature", func(t *testing.T) {
		t.Parallel()
		deps := dm.SuggestDependencies([]string{"validation"})
		assert.Len(t, deps, 1)
		assert.Contains(t, deps[0].Module, "validator")
	})

	t.Run("auth feature", func(t *testing.T) {
		t.Parallel()
		deps := dm.SuggestDependencies([]string{"auth"})
		assert.Len(t, deps, 2)
	})

	t.Run("grpc feature", func(t *testing.T) {
		t.Parallel()
		deps := dm.SuggestDependencies([]string{"grpc"})
		assert.Len(t, deps, 2)
	})

	t.Run("testing feature", func(t *testing.T) {
		t.Parallel()
		deps := dm.SuggestDependencies([]string{"testing"})
		assert.Len(t, deps, 2)
	})

	t.Run("uuid feature", func(t *testing.T) {
		t.Parallel()
		deps := dm.SuggestDependencies([]string{"uuid"})
		assert.Len(t, deps, 1)
	})

	t.Run("unknown feature", func(t *testing.T) {
		t.Parallel()
		deps := dm.SuggestDependencies([]string{"unknown"})
		assert.Empty(t, deps)
	})

	t.Run("multiple features", func(t *testing.T) {
		t.Parallel()
		deps := dm.SuggestDependencies([]string{"validation", "auth", "uuid"})
		assert.Len(t, deps, 4) // validator + jwt + bcrypt + uuid
	})

	t.Run("empty features", func(t *testing.T) {
		t.Parallel()
		deps := dm.SuggestDependencies(nil)
		assert.Empty(t, deps)
	})
}

func TestGetRequiredDependenciesForFeature(t *testing.T) {
	t.Parallel()
	dm := NewDependencyManager(t.TempDir(), false)

	t.Run("grpc", func(t *testing.T) {
		t.Parallel()
		deps := dm.GetRequiredDependenciesForFeature("grpc", nil)
		assert.Len(t, deps, 2)
	})

	t.Run("auth", func(t *testing.T) {
		t.Parallel()
		deps := dm.GetRequiredDependenciesForFeature("auth", nil)
		assert.Len(t, deps, 2)
	})

	t.Run("with validation option", func(t *testing.T) {
		t.Parallel()
		deps := dm.GetRequiredDependenciesForFeature("auth", map[string]bool{"validation": true})
		assert.Len(t, deps, 3) // jwt + bcrypt + validator
	})

	t.Run("unknown type no options", func(t *testing.T) {
		t.Parallel()
		deps := dm.GetRequiredDependenciesForFeature("unknown", nil)
		assert.Empty(t, deps)
	})
}

// Regression test: goca init bakes the driver dependency into go.mod for a
// project's initial database, but a standalone `goca feature`/`goca
// repository --database X` targeting a *different* database (most visibly
// dynamodb/elasticsearch, whose packages are never a transitive dependency
// of anything else already present) previously added nothing at all.
func TestGetRequiredDependenciesForDatabase(t *testing.T) {
	t.Parallel()
	dm := NewDependencyManager(t.TempDir(), false)

	cases := []struct {
		database    string
		wantModules []string
	}{
		{DBPostgres, []string{"gorm.io/gorm", "gorm.io/driver/postgres"}},
		{DBPostgresJSON, []string{"gorm.io/gorm", "gorm.io/driver/postgres"}},
		{DBMySQL, []string{"gorm.io/gorm", "gorm.io/driver/mysql"}},
		{DBSQLite, []string{"gorm.io/gorm", "gorm.io/driver/sqlite"}},
		{DBSQLServer, []string{"gorm.io/gorm", "gorm.io/driver/sqlserver"}},
		{DBMongoDB, []string{"go.mongodb.org/mongo-driver"}},
		{DBDynamoDB, []string{"github.com/aws/aws-sdk-go-v2", "github.com/aws/aws-sdk-go-v2/service/dynamodb", "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"}},
		{DBElasticsearch, []string{"github.com/elastic/go-elasticsearch/v8"}},
	}

	for _, tc := range cases {
		t.Run(tc.database, func(t *testing.T) {
			t.Parallel()
			deps := dm.GetRequiredDependenciesForDatabase(tc.database)
			require.Len(t, deps, len(tc.wantModules))
			for i, module := range tc.wantModules {
				assert.Equal(t, module, deps[i].Module)
			}
		})
	}

	t.Run("unknown database", func(t *testing.T) {
		t.Parallel()
		assert.Empty(t, dm.GetRequiredDependenciesForDatabase("unknown"))
	})
}

func TestIsVersionCompatible(t *testing.T) {
	t.Parallel()
	dm := NewDependencyManager(t.TempDir(), false)

	cases := []struct {
		name     string
		current  string
		required string
		expected bool
	}{
		{"same version", "1.21", "1.21", true},
		{"higher minor", "1.22", "1.21", true},
		{"lower minor", "1.20", "1.21", false},
		{"higher major", "2.0", "1.21", true},
		{"lower major", "0.9", "1.0", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, dm.isVersionCompatible(tc.current, tc.required))
		})
	}
}

func TestDependencyManager_DryRun_UpdateGoMod(t *testing.T) {
	t.Parallel()
	dm := NewDependencyManager(t.TempDir(), true)
	err := dm.UpdateGoMod()
	require.NoError(t, err) // dry run should succeed without real go.mod
}

func TestDependencyManager_PrintDependencySuggestions(t *testing.T) {
	t.Parallel()
	dm := NewDependencyManager(t.TempDir(), false)

	t.Run("empty suggestions", func(t *testing.T) {
		t.Parallel()
		// Should not panic
		dm.PrintDependencySuggestions(nil)
	})

	t.Run("with suggestions", func(t *testing.T) {
		t.Parallel()
		suggestions := []Dependency{
			{Module: "test/module", Version: "v1.0.0", Type: "optional", Reason: "testing"},
		}
		// Should not panic
		dm.PrintDependencySuggestions(suggestions)
	})
}
