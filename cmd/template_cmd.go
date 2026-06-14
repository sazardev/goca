package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var templateManagementCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage custom templates",
	Long: `Manage custom templates for code generation.
Initialize, list, show, and reset templates for personalized code generation.`,
}

var templateInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize custom templates",
	Long: `Initialize custom templates directory with default templates.
Creates .goca/templates/ with customizable templates for all layers.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize configuration integration
		configIntegration := NewConfigIntegration()
		if err := configIntegration.LoadConfigForProject(); err != nil {
			ui.Error(fmt.Sprintf("Could not load configuration: %v", err))
			ui.Dim("Tip: Make sure you're in a GOCA project directory or run 'goca init' first")
			os.Exit(1)
		}

		// Initialize template system
		if err := configIntegration.InitializeTemplateSystem(); err != nil {
			ui.Error(fmt.Sprintf("Error initializing templates: %v", err))
			os.Exit(1)
		}

		// Generate enhanced documentation if templates are available
		if err := configIntegration.GenerateProjectDocumentation(); err != nil {
			ui.Warning(fmt.Sprintf("Could not generate documentation: %v", err))
		}

		ui.Blank()
		ui.Success("Template system initialized successfully!")
		ui.NextSteps([]string{
			fmt.Sprintf("Edit templates in: %s", configIntegration.GetTemplateConfig().Directory),
			"Use functions like {{pascal .EntityName}}, {{snake .EntityName}}",
			"Generate features: goca feature Product --fields \"name:string\"",
			"Your custom templates will be used automatically!",
		})
	},
}

// resolveTemplateDir returns the absolute templates directory for the current project.
func resolveTemplateDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting working directory: %w", err)
	}

	configManager := NewConfigManager()
	if err := configManager.LoadConfig(wd); err != nil {
		return "", fmt.Errorf("error loading config: %w", err)
	}

	config := configManager.GetConfig()
	dir := ".goca/templates"
	if config != nil && config.Templates.Directory != "" {
		dir = config.Templates.Directory
	}

	return filepath.Join(wd, dir), nil
}

// collectTemplateFiles walks the template directory and returns the list of
// template names (relative path without extension) and their absolute paths.
func collectTemplateFiles(templateDir string) (map[string]string, error) {
	result := make(map[string]string)
	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".tmpl") && !strings.HasSuffix(path, ".tpl") {
			return nil
		}
		rel, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}
		name := strings.TrimSuffix(rel, filepath.Ext(rel))
		name = filepath.ToSlash(name)
		result[name] = path
		return nil
	})
	return result, err
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Long:  `List all available custom templates in the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		templateDir, err := resolveTemplateDir()
		if err != nil {
			ui.Error(err.Error())
			os.Exit(1)
		}

		if _, err := os.Stat(templateDir); os.IsNotExist(err) {
			ui.Warning("No templates directory found.")
			ui.Dim("Tip: Run 'goca template init' to create templates")
			return
		}

		files, err := collectTemplateFiles(templateDir)
		if err != nil {
			ui.Error(fmt.Sprintf("Error reading templates: %v", err))
			os.Exit(1)
		}

		if len(files) == 0 {
			ui.Warning("No templates found.")
			ui.Dim("Tip: Run 'goca template init' to create templates")
			return
		}

		names := make([]string, 0, len(files))
		for name := range files {
			names = append(names, name)
		}
		sort.Strings(names)

		if quietMode {
			// Plain, script-friendly output under --quiet.
			for _, name := range names {
				ui.Println(name)
			}
			return
		}

		ui.Header(fmt.Sprintf("Available templates (%d):", len(names)))
		rows := make([][]string, 0, len(names))
		for _, name := range names {
			rel, _ := filepath.Rel(templateDir, files[name])
			rows = append(rows, []string{name, filepath.ToSlash(rel)})
		}
		ui.Table([]string{"Template", "File"}, rows)
	},
}

var templateShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Display the content of a template",
	Long: `Display the raw content of a custom template.

The <name> is the template's relative path without extension,
e.g. "domain/entity" or "usecase/dto". Run 'goca template list' to see names.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := filepath.ToSlash(args[0])

		templateDir, err := resolveTemplateDir()
		if err != nil {
			ui.Error(err.Error())
			os.Exit(1)
		}

		if _, err := os.Stat(templateDir); os.IsNotExist(err) {
			ui.Error("No templates directory found.")
			ui.Dim("Tip: Run 'goca template init' to create templates")
			os.Exit(1)
		}

		files, err := collectTemplateFiles(templateDir)
		if err != nil {
			ui.Error(fmt.Sprintf("Error reading templates: %v", err))
			os.Exit(1)
		}

		// Accept either the bare name or a name including the extension.
		name = strings.TrimSuffix(name, ".tmpl")
		name = strings.TrimSuffix(name, ".tpl")

		path, ok := files[name]
		if !ok {
			ui.Error(fmt.Sprintf("Template %q not found", args[0]))
			ui.Dim("Tip: Run 'goca template list' to see available templates")
			os.Exit(1)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			ui.Error(fmt.Sprintf("Error reading template %q: %v", name, err))
			os.Exit(1)
		}

		// Print raw content (ungated) so it is usable in scripts.
		ui.Printf("%s", string(content))
		if len(content) > 0 && content[len(content)-1] != '\n' {
			ui.Println("")
		}
	},
}

var templateResetForce bool

var templateResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset all templates to built-in defaults",
	Long: `Reset all custom templates to the built-in defaults.

This removes the existing templates directory (after backing it up) and
regenerates the default templates. Because this is destructive, you must
pass --force to confirm.`,
	Run: func(cmd *cobra.Command, args []string) {
		templateDir, err := resolveTemplateDir()
		if err != nil {
			ui.Error(err.Error())
			os.Exit(1)
		}

		exists := true
		if _, statErr := os.Stat(templateDir); os.IsNotExist(statErr) {
			exists = false
		}

		if exists && !templateResetForce {
			ui.Warning(fmt.Sprintf("This will delete and regenerate all templates in: %s", templateDir))
			if ui != nil && ui.IsInteractive() {
				ui.Printf("Type 'yes' to continue: ")
				reader := bufio.NewReader(os.Stdin)
				line, _ := reader.ReadString('\n')
				if strings.TrimSpace(strings.ToLower(line)) != "yes" {
					ui.Info("Reset cancelled.")
					return
				}
			} else {
				ui.Error("Refusing to reset without confirmation. Re-run with --force.")
				os.Exit(1)
			}
		}

		// Back up the existing templates directory before removing it.
		if exists {
			backupDir := templateDir + ".bak"
			_ = os.RemoveAll(backupDir)
			if err := os.Rename(templateDir, backupDir); err != nil {
				ui.Error(fmt.Sprintf("Failed to back up existing templates: %v", err))
				os.Exit(1)
			}
			ui.Info(fmt.Sprintf("Backed up existing templates to: %s", backupDir))
		}

		// Regenerate defaults via the template manager (creates default tree).
		configIntegration := NewConfigIntegration()
		if err := configIntegration.LoadConfigForProject(); err != nil {
			ui.Error(fmt.Sprintf("Could not load configuration: %v", err))
			os.Exit(1)
		}
		if err := configIntegration.InitializeTemplateSystem(); err != nil {
			ui.Error(fmt.Sprintf("Error regenerating templates: %v", err))
			os.Exit(1)
		}

		ui.Success("Templates reset to built-in defaults.")
	},
}

func init() {
	templateResetCmd.Flags().BoolVar(&templateResetForce, "force", false, "Reset without confirmation")

	templateManagementCmd.AddCommand(templateInitCmd)
	templateManagementCmd.AddCommand(templateListCmd)
	templateManagementCmd.AddCommand(templateShowCmd)
	templateManagementCmd.AddCommand(templateResetCmd)
}
