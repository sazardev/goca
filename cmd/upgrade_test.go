package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestBuildSectionRows(t *testing.T) {
	t.Parallel()

	t.Run("empty config all defaults", func(t *testing.T) {
		t.Parallel()
		rows := buildSectionRows(GocaConfig{})
		assert.Len(t, rows, 8)
		for _, row := range rows {
			assert.Len(t, row, 3)
			assert.Equal(t, "○ default", row[1])
		}
	})

	t.Run("project populated", func(t *testing.T) {
		t.Parallel()
		cfg := GocaConfig{}
		cfg.Project.Name = "myproject"
		rows := buildSectionRows(cfg)
		assert.Equal(t, "✓ set", rows[0][1])
		assert.Equal(t, "project", rows[0][0])
	})

	t.Run("database populated", func(t *testing.T) {
		t.Parallel()
		cfg := GocaConfig{}
		cfg.Database.Type = "postgres"
		rows := buildSectionRows(cfg)
		assert.Equal(t, "✓ set", rows[2][1])
		assert.Equal(t, "database", rows[2][0])
	})

	t.Run("generation validation enabled", func(t *testing.T) {
		t.Parallel()
		cfg := GocaConfig{}
		cfg.Generation.Validation.Enabled = true
		rows := buildSectionRows(cfg)
		assert.Equal(t, "✓ set", rows[3][1])
	})

	t.Run("templates populated", func(t *testing.T) {
		t.Parallel()
		cfg := GocaConfig{}
		cfg.Templates.Directory = "custom"
		rows := buildSectionRows(cfg)
		assert.Equal(t, "✓ set", rows[6][1])
	})

	t.Run("deploy docker enabled", func(t *testing.T) {
		t.Parallel()
		cfg := GocaConfig{}
		cfg.Deploy.Docker.Enabled = true
		rows := buildSectionRows(cfg)
		assert.Equal(t, "✓ set", rows[7][1])
	})
}

func TestInjectGocaVersion(t *testing.T) {
	t.Parallel()

	t.Run("inject into empty yaml", func(t *testing.T) {
		t.Parallel()
		yamlStr := `project:
  name: test`
		var doc yaml.Node
		require.NoError(t, yaml.Unmarshal([]byte(yamlStr), &doc))

		injectGocaVersion(&doc, "1.2.3")

		// Re-marshal and check
		out, err := yaml.Marshal(&doc)
		require.NoError(t, err)
		assert.Contains(t, string(out), "goca_version")
		assert.Contains(t, string(out), "1.2.3")
	})

	t.Run("update existing version", func(t *testing.T) {
		t.Parallel()
		yamlStr := `project:
  name: test
  metadata:
    goca_version: "1.0.0"`
		var doc yaml.Node
		require.NoError(t, yaml.Unmarshal([]byte(yamlStr), &doc))

		injectGocaVersion(&doc, "2.0.0")

		out, err := yaml.Marshal(&doc)
		require.NoError(t, err)
		assert.Contains(t, string(out), "2.0.0")
		assert.NotContains(t, string(out), "1.0.0")
	})

	t.Run("nil doc", func(t *testing.T) {
		t.Parallel()
		injectGocaVersion(nil, "1.0.0") // should not panic
	})

	t.Run("empty doc", func(t *testing.T) {
		t.Parallel()
		doc := &yaml.Node{}
		injectGocaVersion(doc, "1.0.0") // should not panic
	})
}

func TestFindMappingKey(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()
		node := &yaml.Node{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "name"},
				{Kind: yaml.ScalarNode, Value: "test"},
				{Kind: yaml.ScalarNode, Value: "version"},
				{Kind: yaml.ScalarNode, Value: "1.0"},
			},
		}
		assert.Equal(t, 0, findMappingKey(node, "name"))
		assert.Equal(t, 2, findMappingKey(node, "version"))
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		node := &yaml.Node{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "name"},
				{Kind: yaml.ScalarNode, Value: "test"},
			},
		}
		assert.Equal(t, -1, findMappingKey(node, "missing"))
	})

	t.Run("empty node", func(t *testing.T) {
		t.Parallel()
		node := &yaml.Node{Kind: yaml.MappingNode}
		assert.Equal(t, -1, findMappingKey(node, "name"))
	})
}

func TestHandleRegenerate(t *testing.T) {
	t.Parallel()

	t.Run("empty name returns error", func(t *testing.T) {
		t.Parallel()
		err := handleRegenerate("", false)
		assert.Error(t, err)
	})

	// handleRegenerate with non-empty name calls ui.* which requires
	// the global ui to be non-nil. We skip those since they are UI-dependent.
}
