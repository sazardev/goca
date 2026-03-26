package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Initialize Clean Architecture project",
	Long: `Creates the base structure of a Go project following Clean Architecture principles, 
including directories, configuration files and layer structure.

Use --template to initialize with predefined configurations:
  goca init myproject --module github.com/user/myproject --template rest-api`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Allow 0 args if --list-templates is set
		listTemplates, _ := cmd.Flags().GetBool("list-templates")
		if listTemplates {
			return nil
		}
		// Allow 0 args in interactive mode (wizard will prompt for project name)
		noInteractive, _ := cmd.Root().PersistentFlags().GetBool("no-interactive")
		if noInteractive {
			return cobra.ExactArgs(1)(cmd, args)
		}
		return cobra.MaximumNArgs(1)(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		module, _ := cmd.Flags().GetString("module")
		database, _ := cmd.Flags().GetString("database")
		auth, _ := cmd.Flags().GetBool("auth")
		api, _ := cmd.Flags().GetString("api")
		config, _ := cmd.Flags().GetBool("config")
		template, _ := cmd.Flags().GetString("template")
		listTemplates, _ := cmd.Flags().GetBool("list-templates")

		// Handle --list-templates flag
		if listTemplates {
			ListAvailableTemplates()
			return
		}

		// Determine project name from args or interactive wizard
		projectName := ""
		if len(args) > 0 {
			projectName = args[0]
		}

		if projectName == "" || module == "" {
			if ui.IsInteractive() {
				// Launch interactive wizard; also collects projectName when not in args
				var err error
				projectName, module, database, api, auth, config, err = runInitWizard(projectName)
				if err != nil {
					ui.Error(fmt.Sprintf("Interactive setup cancelled: %v", err))
					os.Exit(1)
				}
			} else {
				if projectName == "" {
					ui.Error("project name is required")
					cmd.Usage()
					os.Exit(1)
				}
				ui.Error("--module flag is required (or run without --no-interactive for guided setup)")
				os.Exit(1)
			}
		}

		// Validate template if provided
		if template != "" {
			if !ValidateTemplateName(template) {
				ui.Error(fmt.Sprintf("invalid template '%s'", template))
				ui.Println("\nAvailable templates:")
				for _, name := range GetTemplateNames() {
					ui.Dim(fmt.Sprintf("  - %s", name))
				}
				ui.Dim("\nUse --list-templates for detailed descriptions")
				os.Exit(1)
			}
			// When using template, always generate config
			config = true
			ui.Info(fmt.Sprintf("Using template: %s", template))
		}

		ui.Header(fmt.Sprintf("Initializing project '%s' with module '%s'", projectName, module))
		ui.KeyValue("Database", database)
		ui.KeyValue("API", api)
		if auth {
			ui.Feature("Including authentication", false)
		}
		if config {
			ui.Feature("Generating YAML configuration", false)
		}

		// Create configuration integration
		configIntegration := NewConfigIntegration()

		// Merge CLI flags with configuration
		flags := map[string]interface{}{
			"database": database,
			"auth":     auth,
			"api":      api,
		}
		configIntegration.MergeWithCLIFlags(flags)

		stop := ui.Spinner(fmt.Sprintf("Creating project '%s'", projectName))
		// Initialize safety manager
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		backup, _ := cmd.Flags().GetBool("backup")
		sm := NewSafetyManager(dryRun, force, backup)

		if dryRun {
			ui.DryRun("Previewing changes without creating files")
		}

		createProjectStructure(projectName, module, database, auth, api, configIntegration, config, template, sm)
		stop()

		if dryRun {
			sm.PrintSummary()
			return
		}

		ui.Success(fmt.Sprintf("Project '%s' created successfully!", projectName))
		ui.KeyValue("Directory", fmt.Sprintf("./%s", projectName))

		if config || template != "" {
			ui.KeyValue("Configuration file", fmt.Sprintf("./%s/.goca.yaml", projectName))
		}

		nextSteps := []string{
			fmt.Sprintf("cd %s", projectName),
			"go mod tidy",
		}
		if config || template != "" {
			nextSteps = append(nextSteps, "Edit .goca.yaml to customize your project")
		}
		nextSteps = append(nextSteps, "goca feature User --fields \"name:string,email:string\"")
		if !(config || template != "") {
			nextSteps = append(nextSteps, fmt.Sprintf("Or use a template: goca init %s --module %s --template rest-api", projectName, module))
		}
		ui.NextSteps(nextSteps)
	},
}

func createProjectStructure(projectName, module, database string, auth bool, api string, configIntegration *ConfigIntegration, generateConfig bool, template string, sm ...*SafetyManager) {
	// Create main directories
	dirs := []string{
		filepath.Join(projectName, "cmd", "server"),
		filepath.Join(projectName, "internal", "domain"),
		filepath.Join(projectName, "internal", "usecase"),
		filepath.Join(projectName, "internal", "repository"),
		filepath.Join(projectName, "internal", "handler"),
		filepath.Join(projectName, "pkg", "config"),
		filepath.Join(projectName, "pkg", "logger"),
	}

	if auth {
		dirs = append(dirs, filepath.Join(projectName, "pkg", "auth"))
	}

	for _, dir := range dirs {
		_ = os.MkdirAll(dir, 0755)
	}

	// Create go.mod
	createGoMod(projectName, module, database, auth, sm...)

	// Create main.go
	createMainGo(projectName, module, database, sm...)

	// Create .gitignore
	createGitignore(projectName, sm...)

	// Create README.md
	createReadme(projectName, module, sm...)

	// Create config
	createConfig(projectName, module, database, sm...)

	// Create environment files
	createEnvFiles(projectName, database, sm...)

	// Create migrations
	createMigrations(projectName, sm...)

	// Create Makefile and Docker files
	createMakefile(projectName, sm...)
	createDockerfiles(projectName, database, sm...)

	// Create logger
	createLogger(projectName, module, sm...)

	if auth {
		createAuth(projectName, module, sm...)
	}

	// Generate .goca.yaml configuration file if requested or template is used
	if generateConfig && configIntegration != nil {
		configPath := filepath.Join(projectName, ".goca.yaml")

		// Use template configuration if specified
		if template != "" {
			templateConfig, err := GetTemplateConfig(template)
			if err != nil {
				ui.Warning(fmt.Sprintf("Failed to get template config: %v", err))
			} else {
				// Replace placeholders in template
				configContent := strings.ReplaceAll(templateConfig, "project:", fmt.Sprintf("project:\n  name: %s\n  module: %s", projectName, module))

				if err := writeFile(configPath, configContent, sm...); err != nil {
					ui.Warning(fmt.Sprintf("Failed to write config file: %v", err))
				} else {
					ui.FileCreated(fmt.Sprintf("Generated configuration file from template '%s': %s", template, configPath))
				}
			}
		} else {
			// Generate standard configuration
			if err := configIntegration.GenerateConfigFile(projectName, projectName, module, database); err != nil {
				ui.Warning(fmt.Sprintf("Failed to generate config file: %v", err))
			} else {
				ui.FileCreated(fmt.Sprintf("Generated configuration file: %s", configPath))
			}
		}
	}

	// Download dependencies after creating go.mod
	if err := downloadDependencies(projectName); err != nil {
		ui.Warning(fmt.Sprintf("Failed to download dependencies: %v", err))
		ui.Dim("Tip: Run 'go mod download' manually in the project directory")
	}

	// Initialize Git repository
	ui.Info("Initializing Git repository...")
	if err := initializeGitRepository(projectName); err != nil {
		ui.Warning(fmt.Sprintf("Failed to initialize Git repository: %v", err))
		ui.Dim("Tip: You can initialize Git manually with 'git init' in the project directory")
	} else {
		ui.Success("Git repository initialized with initial commit")
	}
}

func init() {
	initCmd.Flags().StringP("module", "m", "", "Go module name (e.g: github.com/user/project)")
	initCmd.Flags().StringP("database", "d", "sqlite", "Database type (postgres, mysql, sqlite, mongodb, sqlserver, dynamodb, elasticsearch)")
	initCmd.Flags().StringP("api", "a", "rest", "API type (rest, graphql, grpc)")
	initCmd.Flags().Bool("auth", false, "Include authentication system")
	initCmd.Flags().Bool("config", true, "Generate .goca.yaml configuration file")
	initCmd.Flags().StringP("template", "t", "", "Use predefined template (minimal, rest-api, microservice, monolith, enterprise)")
	initCmd.Flags().Bool("list-templates", false, "List available project templates")
	initCmd.Flags().Bool("dry-run", false, "Preview changes without creating files")
	initCmd.Flags().Bool("force", false, "Overwrite existing files without asking")
	initCmd.Flags().Bool("backup", false, "Backup existing files before overwriting")
}
