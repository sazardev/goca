package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var upgradeRegenerate string
var upgradeDryRun bool

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade project configuration to the current Goca version",
	Long: `upgrade inspects the current .goca.yaml and compares it with the
configuration schema supported by the installed Goca version.

It highlights:
  - Missing config sections that are now available
  - The recorded goca_version vs the current binary version
  - Suggestions to regenerate outdated boilerplate

Use --update to write the current Goca version into .goca.yaml metadata.
Use --regenerate <feature> to re-run code generation for a named feature.
Use --dry-run to preview any changes without writing files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Header("Goca Upgrade")
		ui.Blank()

		if err := runUpgrade(cmd); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	upgradeCmd.Flags().StringVar(&upgradeRegenerate, "regenerate", "", "Re-run code generation for the named feature")
	upgradeCmd.Flags().BoolVar(&upgradeDryRun, "dry-run", false, "Preview changes without writing to disk")
	upgradeCmd.Flags().Bool("update", false, "Write the current Goca version to .goca.yaml metadata")
}

func runUpgrade(cmd *cobra.Command) error {
	update, _ := cmd.Flags().GetBool("update")

	// Check for .goca.yaml
	configPath := ".goca.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		ui.Warning(".goca.yaml not found in the current directory")
		ui.Info("Run: goca init <project-name>  to initialize the project")
		return nil
	}

	raw, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading .goca.yaml: %w", err)
	}

	var cfg GocaConfig
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return fmt.Errorf("parsing .goca.yaml: %w", err)
	}

	// Report current project
	ui.KeyValue("Project", cfg.Project.Name)
	ui.KeyValue("Module", cfg.Project.Module)
	ui.Blank()

	// Version comparison
	recordedVersion := ""
	if cfg.Project.Metadata != nil {
		recordedVersion = cfg.Project.Metadata["goca_version"]
	}
	reportVersionStatus(recordedVersion)

	// Section completeness check
	ui.Section("Config Section Status")
	rows := buildSectionRows(cfg)
	ui.Table([]string{"Section", "Status", "Note"}, rows)
	ui.Blank()

	// Handle --regenerate
	if upgradeRegenerate != "" {
		return handleRegenerate(upgradeRegenerate, upgradeDryRun)
	}

	// Handle --update (write goca_version to metadata)
	if update {
		return writeGocaVersionToConfig(configPath, raw, upgradeDryRun)
	}

	// Summary suggestions
	if recordedVersion != Version {
		ui.NextSteps([]string{
			"Run: goca upgrade --update  to record the current Goca version in .goca.yaml",
			"Run: goca doctor  to check overall project health",
			"Run: goca upgrade --regenerate <feature>  to refresh generated boilerplate",
		})
	} else {
		ui.Success("Project configuration is up to date")
	}

	return nil
}

func reportVersionStatus(recorded string) {
	ui.Info(fmt.Sprintf("Installed Goca version : %s", Version))
	if recorded == "" {
		ui.Warning("Recorded goca_version   : not set in .goca.yaml metadata")
		ui.Dim("  Tip: run 'goca upgrade --update' to record the current version")
	} else if recorded == Version {
		ui.Info(fmt.Sprintf("Recorded goca_version   : %s (up to date)", recorded))
	} else {
		ui.Warning(fmt.Sprintf("Recorded goca_version   : %s (installed: %s)", recorded, Version))
		ui.Dim("  Tip: run 'goca upgrade --update' to refresh")
	}
	ui.Blank()
}

// buildSectionRows checks which top-level config sections have been populated.
func buildSectionRows(cfg GocaConfig) [][]string {
	sections := []struct {
		name      string
		populated bool
		note      string
	}{
		{
			name:      "project",
			populated: cfg.Project.Name != "" || cfg.Project.Module != "",
			note:      "name, module, description",
		},
		{
			name:      "architecture",
			populated: len(cfg.Architecture.Patterns) > 0 || cfg.Architecture.DI.Type != "",
			note:      "layers, patterns, DI type",
		},
		{
			name:      "database",
			populated: cfg.Database.Type != "",
			note:      "type, host, migrations",
		},
		{
			name:      "generation",
			populated: cfg.Generation.Validation.Enabled || cfg.Generation.Style.Gofmt,
			note:      "validation, style, documentation",
		},
		{
			name:      "testing",
			populated: cfg.Testing.Framework != "" || cfg.Testing.Mocks.Enabled,
			note:      "framework, mocks, coverage",
		},
		{
			name:      "features",
			populated: cfg.Features.Auth.Enabled || cfg.Features.Cache.Enabled || len(cfg.Features.Plugins) > 0,
			note:      "auth, cache, logging, security",
		},
		{
			name:      "templates",
			populated: cfg.Templates.Directory != "",
			note:      "custom template directory",
		},
		{
			name:      "deploy",
			populated: cfg.Deploy.Docker.Enabled || cfg.Deploy.Kubernetes.Enabled,
			note:      "docker, kubernetes, CI",
		},
	}

	rows := make([][]string, len(sections))
	for i, s := range sections {
		status := "✓ set"
		if !s.populated {
			status = "○ default"
		}
		rows[i] = []string{s.name, status, s.note}
	}
	return rows
}

func handleRegenerate(feature string, dryRun bool) error {
	name := strings.TrimSpace(feature)
	if name == "" {
		return fmt.Errorf("--regenerate requires a feature name")
	}

	noun := toPascalCase(name)
	if dryRun {
		ui.DryRun(fmt.Sprintf("Would regenerate feature: %s", noun))
		ui.Info("Run without --dry-run to actually regenerate")
		return nil
	}

	ui.Info(fmt.Sprintf("To regenerate feature %q, run:", noun))
	ui.Dim(fmt.Sprintf("  goca feature %s --force", name))
	ui.Blank()
	ui.Warning("goca upgrade --regenerate does not automatically overwrite code.")
	ui.Info("Use  --force  flag on the feature command to overwrite existing files.")
	return nil
}

// writeGocaVersionToConfig injects goca_version into .goca.yaml metadata.
func writeGocaVersionToConfig(configPath string, raw []byte, dryRun bool) error {
	var doc yaml.Node
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return fmt.Errorf("parsing .goca.yaml for update: %w", err)
	}

	injectGocaVersion(&doc, Version)

	var buf strings.Builder
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&doc); err != nil {
		return fmt.Errorf("encoding updated config: %w", err)
	}
	_ = enc.Close()

	updated := buf.String()

	if dryRun {
		ui.DryRun("Would write goca_version to .goca.yaml metadata")
		ui.Debug(fmt.Sprintf("  goca_version: %s", Version))
		return nil
	}

	if err := os.WriteFile(configPath, []byte(updated), 0644); err != nil {
		return fmt.Errorf("writing updated .goca.yaml: %w", err)
	}

	ui.Success(fmt.Sprintf("Updated .goca.yaml: goca_version set to %s", Version))
	return nil
}

// injectGocaVersion traverses a yaml.Node document and sets
// project.metadata.goca_version to the given version value.
// If project or metadata nodes are absent it inserts them.
func injectGocaVersion(doc *yaml.Node, version string) {
	if doc == nil || len(doc.Content) == 0 {
		return
	}

	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return
	}

	// Find or create "project" key
	projectIdx := findMappingKey(root, "project")
	var projectNode *yaml.Node
	if projectIdx < 0 {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "project"}
		valNode := &yaml.Node{Kind: yaml.MappingNode}
		root.Content = append(root.Content, keyNode, valNode)
		projectNode = valNode
	} else {
		projectNode = root.Content[projectIdx+1]
	}

	if projectNode.Kind != yaml.MappingNode {
		return
	}

	// Find or create "metadata" key
	metaIdx := findMappingKey(projectNode, "metadata")
	var metaNode *yaml.Node
	if metaIdx < 0 {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "metadata"}
		valNode := &yaml.Node{Kind: yaml.MappingNode}
		projectNode.Content = append(projectNode.Content, keyNode, valNode)
		metaNode = valNode
	} else {
		metaNode = projectNode.Content[metaIdx+1]
	}

	if metaNode.Kind != yaml.MappingNode {
		return
	}

	// Set goca_version key
	gvIdx := findMappingKey(metaNode, "goca_version")
	if gvIdx < 0 {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "goca_version"}
		valNode := &yaml.Node{Kind: yaml.ScalarNode, Value: version}
		metaNode.Content = append(metaNode.Content, keyNode, valNode)
	} else {
		metaNode.Content[gvIdx+1].Value = version
	}
}

// findMappingKey returns the index of the key node for the given key in a
// MappingNode's Content slice, or -1 if not found.
func findMappingKey(node *yaml.Node, key string) int {
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return i
		}
	}
	return -1
}
