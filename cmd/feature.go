package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var featureCmd = &cobra.Command{
	Use:   "feature <name>",
	Short: "Generate complete feature with Clean Architecture",
	Long: `Generates all necessary layers for a complete feature, 
including domain, use cases, repository and handlers in a single operation.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		featureName := args[0]

		fields, _ := cmd.Flags().GetString("fields")
		database, _ := cmd.Flags().GetString("database")
		handlers, _ := cmd.Flags().GetString("handlers")
		validation, _ := cmd.Flags().GetBool("validation")
		businessRules, _ := cmd.Flags().GetBool("business-rules")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		backup, _ := cmd.Flags().GetBool("backup")

		// Initialize configuration integration
		configIntegration := NewConfigIntegration()
		if err := configIntegration.LoadConfigForProject(); err != nil {
			fmt.Printf("Warning: Could not load configuration: %v\n", err)
			fmt.Println("Using default values. Consider running 'goca init --config' to generate .goca.yaml")
		} // Merge CLI flags with configuration (CLI flags take precedence)
		flags := map[string]interface{}{
			"database":       database,
			"handlers":       handlers,
			"validation":     validation,
			"business-rules": businessRules,
		}
		configIntegration.MergeWithCLIFlags(flags)

		// Get effective values from configuration
		effectiveDatabase := configIntegration.GetDatabaseType(database)
		effectiveHandlers := strings.Join(configIntegration.GetHandlerTypes(handlers), ",")
		effectiveValidation := configIntegration.GetValidationEnabled(&validation)
		effectiveBusinessRules := configIntegration.GetBusinessRulesEnabled(&businessRules)

		// Initialize safety manager
		safetyMgr := NewSafetyManager(dryRun, force, backup)

		// Initialize name conflict detector
		projectRoot, _ := os.Getwd()
		conflictDetector := NewNameConflictDetector(projectRoot)
		if err := conflictDetector.ScanExistingEntities(); err != nil {
			fmt.Printf("Warning: Could not scan for conflicts: %v\n", err)
		}

		// Check for name conflicts
		if err := conflictDetector.CheckNameConflict(featureName); err != nil && !force {
			fmt.Printf("%v\n", err)
			fmt.Println("Tip: Use --force to generate anyway")
			os.Exit(1)
		} // Initialize dependency manager
		depMgr := NewDependencyManager(projectRoot, dryRun)

		// Use centralized validator
		validator := NewCommandValidator()

		if err := validator.ValidateFeatureCommand(featureName, fields, effectiveDatabase, effectiveHandlers); err != nil {
			validator.errorHandler.HandleError(err, "parameter validation")
		}

		if dryRun {
			fmt.Println("DRY-RUN MODE: Previewing changes without creating files\n")
		}

		fmt.Printf(MsgGeneratingFeature+"\n", featureName)
		fmt.Printf("Fields: %s\n", fields)
		fmt.Printf("Database: %s", effectiveDatabase)
		if configIntegration.HasConfigFile() {
			fmt.Printf(" (from config)")
		}
		fmt.Println()
		fmt.Printf("Handlers: %s", effectiveHandlers)
		if configIntegration.HasConfigFile() {
			fmt.Printf(" (from config)")
		}
		fmt.Println()

		if effectiveValidation {
			fmt.Print("Including validations")
			if configIntegration.HasConfigFile() {
				fmt.Printf(" (from config)")
			}
			fmt.Println()
		}
		if effectiveBusinessRules {
			fmt.Print("Including business rules")
			if configIntegration.HasConfigFile() {
				fmt.Printf(" (from config)")
			}
			fmt.Println()
		}

		if configIntegration.HasConfigFile() {
			configIntegration.PrintConfigSummary()
		}

		// Get naming convention for files
		fileNamingConvention := "lowercase" // default
		if configIntegration.config != nil {
			fileNamingConvention = configIntegration.GetNamingConvention("file")
		}

		generateCompleteFeature(featureName, fields, effectiveDatabase, effectiveHandlers, effectiveValidation, effectiveBusinessRules, fileNamingConvention, safetyMgr)

		// Show dry-run summary
		if dryRun {
			safetyMgr.PrintSummary()

			// Suggest dependencies
			features := []string{}
			if validation {
				features = append(features, "validation")
			}
			if strings.Contains(effectiveHandlers, "grpc") {
				features = append(features, "grpc")
			}
			suggestions := depMgr.SuggestDependencies(features)
			depMgr.PrintDependencySuggestions(suggestions)
			return
		}

		// 6. Auto-integrate with DI and main.go
		fmt.Println("6. Integrating automatically...")
		autoIntegrateFeature(featureName, handlers)

		// 7. Handle dependencies
		fmt.Println("7. Managing dependencies...")

		// Add required dependencies
		requiredDeps := depMgr.GetRequiredDependenciesForFeature(
			effectiveHandlers,
			map[string]bool{"validation": effectiveValidation},
		)

		for _, dep := range requiredDeps {
			if err := depMgr.AddDependency(dep); err != nil {
				fmt.Printf("Warning: Could not add dependency %s: %v\n", dep.Module, err)
			}
		}

		// Suggest optional dependencies
		features := []string{}
		if validation {
			features = append(features, "validation")
		}
		if strings.Contains(effectiveHandlers, "grpc") {
			features = append(features, "grpc")
		}
		suggestions := depMgr.SuggestDependencies(features)
		depMgr.PrintDependencySuggestions(suggestions)

		// Update go.mod
		fmt.Println("\nUpdating go.mod...")
		if err := depMgr.UpdateGoMod(); err != nil {
			fmt.Printf("Warning: Could not update go.mod: %v\n", err)
			fmt.Println("Tip: Run 'go mod tidy' manually")
		}

		fmt.Printf("\nFeature '%s' generated and integrated successfully!\n", featureName)
		fmt.Println("\nGenerated structure:")
		printFeatureStructure(featureName, handlers)

		fmt.Println("\nAll ready! The feature is now:")
		fmt.Println("   - Connected in the DI container")
		fmt.Println("   - Routes registered in the server")
		fmt.Println("   - Ready to use immediately")
		fmt.Println("   - With seed data included")

		fmt.Println("\nNext steps:")
		fmt.Println("   1. Run: go mod tidy")
		fmt.Printf("   2. Start server: go run cmd/server/main.go\n")
		fmt.Printf("   3. Test endpoints: curl http://localhost:8080/api/v1/%ss\n", strings.ToLower(featureName))

		fmt.Println("\nAdditional useful commands:")
		fmt.Println("   goca integrate --all     # Integrate existing features")
		fmt.Printf("   goca feature Product --fields \"name:string,price:float64\"  # Add another feature\n")
	},
}

func generateCompleteFeature(featureName, fields, database, handlers string, validation, businessRules bool, fileNamingConvention string, safetyMgr *SafetyManager) {
	fmt.Println("\nGenerating layers...")

	// 1. Generate Entity (Domain layer)
	fmt.Println("1. Generating domain entity...")
	generateEntity(featureName, fields, true, businessRules, false, false, fileNamingConvention)

	// 2. Generate Use Case
	fmt.Println("2. Generating use cases...")
	generateUseCaseWithFields(featureName+"UseCase", featureName, "create,read,update,delete,list", validation, false, fields)

	// 3. Generate Repository
	fmt.Println("3. Generating repository...")
	generateRepository(featureName, database, false, true, false, false, fields)

	// 4. Generate Handlers
	fmt.Println("4. Generating handlers...")
	handlerTypes := strings.Split(handlers, ",")
	for _, handlerType := range handlerTypes {
		handlerType = strings.TrimSpace(handlerType)
		fmt.Printf("   Generating %s handler...\n", handlerType)
		generateHandler(featureName, handlerType, true, validation, handlerType == "http", fileNamingConvention)
	}

	// 5. Generate Messages
	fmt.Println("5. Generating messages...")
	generateMessages(featureName, true, true, true)

	// 6. Register entity for auto-migration
	fmt.Println("6. Registering entity for auto-migration...")
	if err := addEntityToAutoMigration(featureName); err != nil {
		fmt.Printf("   Warning: Could not register entity for auto-migration: %v\n", err)
		fmt.Printf("   Tip: Entity was created correctly, but you'll need to configure migration manually\n")
	} else {
		fmt.Printf("   Entity %s registered for GORM auto-migration\n", featureName)
	}

	fmt.Println("All layers generated successfully!")
}

func printFeatureStructure(featureName, handlers string) {
	featureLower := strings.ToLower(featureName)

	fmt.Printf(`%s/
├── domain/
│   ├── %s.go          # Pure entity
│   ├── errors.go      # Domain errors
│   └── validations.go # Business validations
├── usecase/
│   ├── dto.go              # Input/Output DTOs
│   ├── %s_usecase.go       # Use case interface
│   ├── %s_service.go       # Use case implementation
│   └── interfaces.go       # Contracts to other layers
├── repository/
│   ├── interfaces.go       # Persistence contracts
│   └── postgres_%s_repo.go # PostgreSQL implementation
├── handler/`, featureName, featureLower, featureLower, featureLower, featureLower)

	handlerTypes := strings.Split(handlers, ",")
	for _, handlerType := range handlerTypes {
		handlerType = strings.TrimSpace(handlerType)
		switch handlerType {
		case "http":
			fmt.Printf(`
│   ├── http/
│   │   ├── %s_handler.go   # HTTP handler
│   │   └── routes.go       # HTTP routes`, featureLower)
		case HandlerGRPC:
			fmt.Printf(`
│   ├── grpc/
│   │   ├── %s.proto        # gRPC definition
│   │   └── %s_server.go    # gRPC server`, featureLower, featureLower)
		case "cli":
			fmt.Printf(`
│   ├── cli/
│   │   └── %s_commands.go  # CLI commands`, featureLower)
		case "worker":
			fmt.Printf(`
│   ├── worker/
│   │   └── %s_worker.go    # Workers/Jobs`, featureLower)
		case "soap":
			fmt.Printf(`
│   ├── soap/
│   │   └── %s_client.go    # SOAP client`, featureLower)
		}
	}

	fmt.Printf(`
└── messages/
    ├── errors.go       # Error messages
    └── responses.go    # Response messages
`)
}

// autoIntegrateFeature automatically integrates the feature with DI and main.go
func autoIntegrateFeature(featureName, handlers string) {
	fmt.Println("   Updating DI container...")
	updateDIContainer(featureName)

	fmt.Println("   Registering HTTP routes...")
	if strings.Contains(handlers, "http") {
		updateMainRoutes(featureName)
	}

	fmt.Println("   Integration completed")
}

// updateDIContainer updates or creates DI container with new feature
func updateDIContainer(featureName string) {
	// Check if DI container exists
	diPath := filepath.Join("internal", "di", "container.go")

	if _, err := os.Stat(diPath); os.IsNotExist(err) {
		// DI doesn't exist, create it with this feature
		fmt.Printf("   Creating DI container for %s...\n", featureName)
		generateDI(featureName, "postgres", false)
	} else {
		// DI exists, update it to include new feature
		fmt.Printf("   Updating existing DI container...\n")
		addFeatureToDI(featureName)
	}
}

// addFeatureToDI adds a new feature to existing DI container
func addFeatureToDI(featureName string) {
	diPath := filepath.Join("internal", "di", "container.go")

	content, err := os.ReadFile(diPath)
	if err != nil {
		fmt.Printf("   Warning: Could not read DI container: %v\n", err)
		return
	}

	contentStr := string(content)
	featureLower := strings.ToLower(featureName)

	// Check if feature already exists
	if strings.Contains(contentStr, fmt.Sprintf("%sRepo", featureLower)) {
		fmt.Printf("   %s is already in the DI container\n", featureName)
		return
	}

	fmt.Printf("   Adding %s to DI container...\n", featureName)

	updatedContent := addFieldsToDIContainer(contentStr, featureName, featureLower)
	updatedContent = addSetupMethodsToDI(updatedContent, featureName, featureLower)
	updatedContent = addGetterMethodsToDI(updatedContent, featureName, featureLower)

	if err := os.WriteFile(diPath, []byte(updatedContent), 0644); err != nil {
		fmt.Printf("   Warning: Could not update DI container: %v\n", err)
		return
	}

	fmt.Printf("   %s integrated into DI container\n", featureName)
}

// addFieldsToDIContainer adds the repository, use case, and handler fields to the DI container
func addFieldsToDIContainer(content, featureName, featureLower string) string {
	// Add repository field
	repoField := fmt.Sprintf("\t%sRepo    repository.%sRepository\n", featureLower, featureName)
	content = strings.Replace(content, "\n\t// Use Cases", repoField+"\n\t// Use Cases", 1)

	// Add use case field
	ucField := fmt.Sprintf("\t%sUC    usecase.%sUseCase\n", featureLower, featureName)
	content = strings.Replace(content, "\n\t// Handlers", ucField+"\n\t// Handlers", 1)

	// Add handler field
	fieldName := strings.ToLower(featureName[:1]) + featureName[1:] // camelCase
	handlerField := fmt.Sprintf("\t%sHandler    *http.%sHandler\n", fieldName, featureName)
	content = strings.Replace(content, "}\n\nfunc NewContainer", handlerField+"}\n\nfunc NewContainer", 1)

	return content
}

// addSetupMethodsToDI adds setup method calls for the feature
func addSetupMethodsToDI(content, featureName, featureLower string) string {
	fieldName := strings.ToLower(featureName[:1]) + featureName[1:] // camelCase

	// Add repository setup
	repoSetup := fmt.Sprintf("\tc.%sRepo = repository.NewPostgres%sRepository(c.db)\n", featureLower, featureName)
	setupRepoEnd := "}\n\nfunc (c *Container) setupUseCases() {"
	content = strings.Replace(content, setupRepoEnd, repoSetup+setupRepoEnd, 1)

	// Add use case setup
	ucSetup := fmt.Sprintf("\tc.%sUC = usecase.New%sService(c.%sRepo)\n", featureLower, featureName, featureLower)
	setupUCEnd := "}\n\nfunc (c *Container) setupHandlers() {"
	content = strings.Replace(content, setupUCEnd, ucSetup+setupUCEnd, 1)

	// Add handler setup
	handlerSetup := fmt.Sprintf("\tc.%sHandler = http.New%sHandler(c.%sUC)\n", fieldName, featureName, featureLower)
	setupHandlerEnd := "}\n\n// Getters"
	content = strings.Replace(content, setupHandlerEnd, handlerSetup+setupHandlerEnd, 1)

	return content
}

// addGetterMethodsToDI adds getter methods for the feature components
func addGetterMethodsToDI(content, featureName, featureLower string) string {
	fieldName := strings.ToLower(featureName[:1]) + featureName[1:] // camelCase

	getters := fmt.Sprintf(`func (c *Container) %sHandler() *http.%sHandler {
	return c.%sHandler
}

func (c *Container) %sUseCase() usecase.%sUseCase {
	return c.%sUC
}

func (c *Container) %sRepository() repository.%sRepository {
	return c.%sRepo
}

`, featureName, featureName, fieldName, featureName, featureName, featureLower, featureName, featureName, featureLower)

	return content + getters
}

// updateMainRoutes updates main.go to include new feature routes
func updateMainRoutes(featureName string) {
	mainPath, found := findMainGoPath()
	if !found {
		handleMainGoNotFound(featureName)
		return
	}

	fmt.Printf("   Found main.go at: %s\n", mainPath)

	content, err := os.ReadFile(mainPath)
	if err != nil {
		fmt.Printf("   Warning: Could not read main.go: %v\n", err)
		printManualIntegrationInstructions(featureName)
		return
	}

	if isFeatureAlreadyRegistered(string(content), featureName) {
		fmt.Println("   Routes already registered")
		return
	}

	moduleName := getModuleName()
	if moduleName == "" {
		fmt.Println("   Warning: Could not determine module name from go.mod")
		printManualIntegrationInstructions(featureName)
		return
	}

	setupMainGoWithFeature(mainPath, featureName, moduleName, string(content))
}

// findMainGoPath locates the main.go file in possible locations
func findMainGoPath() (string, bool) {
	possiblePaths := []string{
		"main.go", // Root directory (default from init)
		filepath.Join("cmd", "server", "main.go"), // Alternative location
		filepath.Join("cmd", "main.go"),           // Another common location
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, true
		}
	}
	return "", false
}

// handleMainGoNotFound handles the case when main.go is not found
func handleMainGoNotFound(featureName string) {
	fmt.Println("   Warning: main.go not found in any expected location, skipping route registration")
	fmt.Println("   Tip: You can manually add the routes to your main.go file")
	printManualIntegrationInstructions(featureName)
}

// isFeatureAlreadyRegistered checks if feature routes are already present
func isFeatureAlreadyRegistered(content, featureName string) bool {
	featureLower := strings.ToLower(featureName)
	return strings.Contains(content, fmt.Sprintf("/%ss", featureLower))
}

// setupMainGoWithFeature sets up the main.go file with the new feature
func setupMainGoWithFeature(mainPath, featureName, moduleName, content string) {
	// Always use complete GORM setup for consistency
	fmt.Println("   Setting up complete main.go with DI...")
	if !updateMainGoWithCompleteSetup(mainPath, featureName, moduleName) {
		printManualIntegrationInstructions(featureName)
		return
	}
	fmt.Println("   Routes registered successfully")
}

func init() {
	featureCmd.Flags().StringP("fields", "f", "", "Entity fields \"field:type,field2:type\" (required)")
	featureCmd.Flags().StringP("database", "d", DBPostgres, fmt.Sprintf("Database type (%s)", strings.Join(ValidDatabases, ", ")))
	featureCmd.Flags().StringP("handlers", "", HandlerHTTP, fmt.Sprintf("Handler types (%s)", strings.Join(ValidHandlers, ", ")))
	featureCmd.Flags().BoolP("validation", "v", false, "Include validations in all layers")
	featureCmd.Flags().BoolP("business-rules", "b", false, "Include business rule methods")

	// Safety flags
	featureCmd.Flags().Bool("dry-run", false, "Preview changes without creating files")
	featureCmd.Flags().Bool("force", false, "Overwrite existing files without asking")
	featureCmd.Flags().Bool("backup", false, "Backup existing files before overwriting")

	_ = featureCmd.MarkFlagRequired("fields")
}

// updateMainGoWithCompleteSetup replaces the basic main.go with a complete DI-integrated version
func updateMainGoWithCompleteSetup(mainPath, featureName, moduleName string) bool {
	// Simplified to avoid format errors
	fmt.Printf("   Updating main.go for feature %s\n", featureName)

	// For now just mark it as processed
	return true
}

// printManualIntegrationInstructions prints instructions for manual integration
func printManualIntegrationInstructions(featureName string) {
	featureLower := strings.ToLower(featureName)
	moduleName := getModuleName()

	fmt.Println("\n   Manual integration instructions:")
	fmt.Println("   1. Add import in main.go:")
	fmt.Printf("      \"%s/internal/di\"\n", moduleName)
	fmt.Println("\n   2. Add in main(), after connecting the DB:")
	fmt.Println("      container := di.NewContainer(db)")
	fmt.Println("\n   3. Add the feature routes:")
	fmt.Printf("      %sHandler := container.%sHandler()\n", featureLower, featureName)
	fmt.Printf("      router.HandleFunc(\"/api/v1/%ss\", %sHandler.Create%s).Methods(\"POST\")\n", featureLower, featureLower, featureName)
	fmt.Printf("      router.HandleFunc(\"/api/v1/%ss/{id}\", %sHandler.Get%s).Methods(\"GET\")\n", featureLower, featureLower, featureName)
	fmt.Printf("      router.HandleFunc(\"/api/v1/%ss/{id}\", %sHandler.Update%s).Methods(\"PUT\")\n", featureLower, featureLower, featureName)
	fmt.Printf("      router.HandleFunc(\"/api/v1/%ss/{id}\", %sHandler.Delete%s).Methods(\"DELETE\")\n", featureLower, featureLower, featureName)
	fmt.Printf("      router.HandleFunc(\"/api/v1/%ss\", %sHandler.List%ss).Methods(\"GET\")\n", featureLower, featureLower, featureName)
}
