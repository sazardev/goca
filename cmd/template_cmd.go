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
			fmt.Printf("❌ Error: Could not load configuration: %v\n", err)
			fmt.Println("💡 Make sure you're in a GOCA project directory or run 'goca init' first")
			os.Exit(1)
		}

		// Initialize template system
		if err := configIntegration.InitializeTemplateSystem(); err != nil {
			fmt.Printf("❌ Error initializing templates: %v\n", err)
			os.Exit(1)
		}

		// Generate enhanced documentation if templates are available
		if err := configIntegration.GenerateProjectDocumentation(); err != nil {
			fmt.Printf("⚠️  Warning: Could not generate documentation: %v\n", err)
		}

		fmt.Println()
		fmt.Println("🎉 Template system initialized successfully!")
		fmt.Println()
		fmt.Println("📋 Next steps:")
		fmt.Printf("   1. Edit templates in: %s\n", configIntegration.GetTemplateConfig().Directory)
		fmt.Println("   2. Use functions like {{pascal .EntityName}}, {{snake .EntityName}}")
		fmt.Println("   3. Generate features: goca feature Product --fields \"name:string\"")
		fmt.Println("   4. Your custom templates will be used automatically!")
	},
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Long:  `List all available custom templates in the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Checking for templates...")

		// Get current working directory
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("❌ Error getting working directory: %v\n", err)
			return
		}

		fmt.Printf("📁 Working directory: %s\n", wd)

		// Try simple config manager first
		configManager := NewConfigManager()
		if err := configManager.LoadConfig(wd); err != nil {
			fmt.Printf("❌ Error loading config: %v\n", err)
			return
		}

		config := configManager.GetConfig()
		if config == nil {
			fmt.Println("❌ Config is nil")
			return
		}

		fmt.Printf("✅ Config loaded. Templates dir: %s\n", config.Templates.Directory)

		// Try template manager
		templateDir := filepath.Join(wd, config.Templates.Directory)
		fmt.Printf("📂 Full template path: %s\n", templateDir)

		if _, err := os.Stat(templateDir); os.IsNotExist(err) {
			fmt.Println("� No templates directory found.")
			fmt.Println("💡 Run 'goca template init' to create templates")
			return
		}

		fmt.Println("✅ Template directory exists")
	},
}

func init() {
	templateManagementCmd.AddCommand(templateInitCmd)
	templateManagementCmd.AddCommand(templateListCmd)
}
