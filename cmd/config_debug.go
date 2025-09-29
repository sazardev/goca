package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show configuration information",
	Long:  `Display current project configuration and status.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Checking GOCA configuration...")

		// Get current working directory
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("❌ Error getting working directory: %v\n", err)
			return
		}
		fmt.Printf("📁 Current directory: %s\n", wd)

		// Try to find config file
		configManager := NewConfigManager()
		configPath := configManager.findConfigFile(wd)

		if configPath == "" {
			fmt.Println("⚠️  No .goca.yaml configuration file found")
			fmt.Println("💡 Run 'goca init --config' to generate one")
			return
		}

		fmt.Printf("📋 Found configuration: %s\n", configPath)

		// Try to load config
		if err := configManager.loadFromFile(configPath); err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			return
		}

		config := configManager.GetConfig()
		if config == nil {
			fmt.Println("❌ Configuration is nil")
			return
		}

		fmt.Println("✅ Configuration loaded successfully!")
		fmt.Printf("   📦 Project: %s\n", config.Project.Name)
		fmt.Printf("   🔗 Module: %s\n", config.Project.Module)
		fmt.Printf("   🗄️  Database: %s\n", config.Database.Type)
		fmt.Printf("   📝 Templates: %s\n", config.Templates.Directory)
	},
}

func init() {
	// We'll add this to root manually for testing
}
