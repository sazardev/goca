package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var templateManagementCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage custom templates",
	Long: `Manage custom templates for code generation. 
Initialize, list, and customize templates for personalized code generation.`,
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

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Long:  `List all available custom templates in the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.Info("Checking for templates...")

		// Get current working directory
		wd, err := os.Getwd()
		if err != nil {
			ui.Error(fmt.Sprintf("Error getting working directory: %v", err))
			return
		}

		ui.KeyValue("Working directory", wd)

		// Try simple config manager first
		configManager := NewConfigManager()
		if err := configManager.LoadConfig(wd); err != nil {
			ui.Error(fmt.Sprintf("Error loading config: %v", err))
			return
		}

		config := configManager.GetConfig()
		if config == nil {
			ui.Warning("Config is nil")
			return
		}

		ui.KeyValue("Templates dir", config.Templates.Directory)

		// Try template manager
		templateDir := filepath.Join(wd, config.Templates.Directory)
		ui.KeyValue("Full template path", templateDir)

		if _, err := os.Stat(templateDir); os.IsNotExist(err) {
			ui.Warning("No templates directory found.")
			ui.Dim("Tip: Run 'goca template init' to create templates")
			return
		}

		ui.Success("Template directory exists")
	},
}

func init() {
	templateManagementCmd.AddCommand(templateInitCmd)
	templateManagementCmd.AddCommand(templateListCmd)
}
