package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var repositoryCmd = &cobra.Command{
	Use:   "repository <entity>",
	Short: "Generate repositories with interfaces",
	Long: `Creates repositories that implement the Repository pattern with 
well-defined interfaces and database-specific implementations.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entity := args[0]

		database, _ := cmd.Flags().GetString(DatabaseFlag)
		interfaceOnly, _ := cmd.Flags().GetBool(InterfaceOnlyFlag)
		implementation, _ := cmd.Flags().GetBool(ImplementationFlag)
		cache, _ := cmd.Flags().GetBool(CacheFlag)
		transactions, _ := cmd.Flags().GetBool(TransactionsFlag)
		fields, _ := cmd.Flags().GetString("fields")

		// Initialize config integration
		configIntegration := NewConfigIntegration()
		if err := configIntegration.LoadConfigForProject(); err != nil {
			ui.Warning(fmt.Sprintf("Could not load configuration: %v", err))
			ui.Dim("Using default values. Consider running 'goca init --config' to generate .goca.yaml")
		} // Merge CLI flags with configuration (only explicitly changed flags)
		flags := map[string]interface{}{}
		if cmd.Flags().Changed("database") {
			flags["database"] = database
		}
		if cmd.Flags().Changed("cache") {
			flags["cache"] = cache
		}
		if cmd.Flags().Changed("transactions") {
			flags["transactions"] = transactions
		}
		if len(flags) > 0 {
			configIntegration.MergeWithCLIFlags(flags)
		}

		// Get effective values from configuration
		effectiveDatabase := database
		if !cmd.Flags().Changed("database") && configIntegration.config != nil {
			effectiveDatabase = configIntegration.config.Database.Type
		}

		// Validate with the robust validator
		validator := NewFieldValidator()

		if err := validator.ValidateEntityName(entity); err != nil {
			ui.Error(fmt.Sprintf("Invalid entity name: %v", err))
			return
		}

		if effectiveDatabase != "" {
			if err := validator.ValidateDatabase(effectiveDatabase); err != nil {
				ui.Error(fmt.Sprintf("Invalid database: %v", err))
				return
			}
		}

		ui.Header(fmt.Sprintf("Generating repository for entity '%s'", entity))

		if effectiveDatabase != "" && !interfaceOnly {
			if configIntegration.HasConfigFile() && !cmd.Flags().Changed("database") {
				ui.KeyValueFromConfig("Database", effectiveDatabase)
			} else {
				ui.KeyValue("Database", effectiveDatabase)
			}
		}
		if interfaceOnly {
			ui.Feature("Interface only", false)
		}
		if implementation {
			ui.Feature("Implementation only", false)
		}
		if cache {
			ui.Feature("Including cache", false)
		}
		if transactions {
			ui.Feature("Including transactions", false)
		}
		if fields != "" {
			ui.Feature(fmt.Sprintf("Custom fields: %s", fields), false)
		}

		if configIntegration.HasConfigFile() {
			configIntegration.PrintConfigSummary()
		}

		// Initialize safety manager
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		backup, _ := cmd.Flags().GetBool("backup")
		sm := NewSafetyManager(dryRun, force, backup)

		if dryRun {
			ui.DryRun("Previewing changes without creating files")
		}

		generateRepository(entity, effectiveDatabase, interfaceOnly, implementation, cache, transactions, fields, sm)

		if dryRun {
			sm.PrintSummary()
			return
		}

		// The generated implementation needs its database driver in go.mod —
		// `goca init` only adds the driver for the project's *initial*
		// database, so a standalone `goca repository --database X` targeting
		// a different one (most visibly dynamodb/elasticsearch, which are
		// never a transitive dependency of anything else already present)
		// would otherwise leave the project failing to build until the user
		// ran `go get` manually.
		if effectiveDatabase != "" && !interfaceOnly {
			projectRoot, _ := os.Getwd()
			depMgr := NewDependencyManager(projectRoot, false)
			for _, dep := range depMgr.GetRequiredDependenciesForDatabase(effectiveDatabase) {
				if err := depMgr.AddDependency(dep); err != nil {
					ui.Warning(fmt.Sprintf("Could not add dependency %s: %v", dep.Module, err))
				}
			}
		}

		ui.Success(fmt.Sprintf("Repository for '%s' generated successfully!", entity))
	},
}

func generateRepository(entity, database string, interfaceOnly, implementation, cache, transactions bool, fields string, sm ...*SafetyManager) {
	// Create repository directory if it doesn't exist
	repoDir := "internal/repository"
	_ = os.MkdirAll(repoDir, 0o755)

	// Parse the provided fields, or recover them from the already-generated
	// entity. Using the field-aware (GORM) generators keeps the repository
	// consistent with the DI container and the feature command, which expect a
	// New<DB><Entity>Repository(*gorm.DB) constructor.
	var parsedFields []Field
	if fields != "" {
		parsedFields = parseFields(fields)
	} else if fs := readEntityFieldsString(entity); fs != "" {
		parsedFields = parseFields(fs)
	}

	// Generate interface unless the caller asked for the implementation only
	// (--implementation), which assumes the interface already exists (e.g. from
	// a prior --interface-only run). This mirrors --interface-only, which emits
	// only the contract.
	// - interfaceOnly=true:            only the interface
	// - implementation=true:           only the implementation
	// - neither (default):             both
	if !implementation {
		if len(parsedFields) > 0 {
			generateRepositoryInterfaceWithFields(repoDir, entity, parsedFields, transactions, sm...)
		} else {
			generateRepositoryInterface(repoDir, entity, transactions, sm...)
		}
	}

	// Generate implementation if not interface-only and database is specified
	if !interfaceOnly && database != "" {
		if len(parsedFields) > 0 {
			generateRepositoryImplementationWithFields(repoDir, entity, database, parsedFields, cache, transactions, sm...)
		} else {
			generateRepositoryImplementation(repoDir, entity, database, cache, transactions, sm...)
		}
	}

	// Generate cache decorator when --cache is enabled
	if cache {
		generateCacheDecorator(entity, parsedFields, sm...)
		if err := generateCachePackage(sm...); err != nil {
			ui.Warning(fmt.Sprintf("Could not generate cache package: %v", err))
		}
	}
}

func generateRepositoryInterface(dir, entity string, transactions bool, sm ...*SafetyManager) {
	filename := filepath.Join(dir, "interfaces.go")

	// A custom template only covers the base CRUD interface (no transaction
	// methods) and, like the DTO file, cannot safely merge into an existing
	// multi-entity interfaces.go. Only use it for a fresh file without
	// --transactions.
	if !transactions {
		if _, err := os.Stat(filename); err != nil {
			if custom, ok := renderCustomTemplate("repository/repo", buildRepositoryTemplateData(entity)); ok {
				if err := writeGoFileMerged(filename, custom, sm...); err != nil {
					fmt.Printf("Error writing file %s: %v\n", filename, err)
				}
				return
			}
		}
	}

	// Get the module name from go.mod
	moduleName := getModuleName()

	var content strings.Builder

	// Check if interfaces.go already exists
	if _, err := os.Stat(filename); err == nil {
		// File exists, read its content
		existingContent, err := os.ReadFile(filename)
		if err == nil {
			existingStr := string(existingContent)
			// Check if the interface already exists
			interfaceName := fmt.Sprintf("type %sRepository interface", entity)
			if strings.Contains(existingStr, interfaceName) {
				// Interface already exists, don't regenerate
				return
			}

			// Add the existing content without the final newline
			content.WriteString(strings.TrimSuffix(existingStr, "\n"))
			content.WriteString("\n\n")
		}
	} else {
		// File doesn't exist, create header
		content.WriteString("package repository\n\n")
		if transactions {
			content.WriteString(fmt.Sprintf("import (\n\t\"%s/internal/domain\"\n\t\"gorm.io/gorm\"\n)\n\n", getImportPath(moduleName)))
		} else {
			content.WriteString(fmt.Sprintf("import \"%s/internal/domain\"\n\n", getImportPath(moduleName)))
		}
	}

	content.WriteString(fmt.Sprintf("type %sRepository interface {\n", entity))
	content.WriteString(fmt.Sprintf("\tSave(%s *domain.%s) error\n", strings.ToLower(entity), entity))
	content.WriteString(fmt.Sprintf("\tFindByID(id int) (*domain.%s, error)\n", entity))
	content.WriteString(fmt.Sprintf("\tUpdate(%s *domain.%s) error\n", strings.ToLower(entity), entity))
	content.WriteString("\tDelete(id int) error\n")
	content.WriteString(fmt.Sprintf("\tFindAll() ([]domain.%s, error)\n", entity))

	if transactions {
		content.WriteString(fmt.Sprintf("\tSaveWithTx(tx *gorm.DB, %s *domain.%s) error\n", strings.ToLower(entity), entity))
		content.WriteString(fmt.Sprintf("\tUpdateWithTx(tx *gorm.DB, %s *domain.%s) error\n", strings.ToLower(entity), entity))
		content.WriteString("\tDeleteWithTx(tx *gorm.DB, id int) error\n")
	}

	content.WriteString("}\n")

	if err := writeGoFileMerged(filename, content.String(), sm...); err != nil {
		fmt.Printf("Error writing file %s: %v\n", filename, err)
	}
}

func generateRepositoryImplementation(dir, entity, database string, cache, transactions bool, sm ...*SafetyManager) {
	switch database {
	case DBPostgres:
		generatePostgresRepository(dir, entity, cache, transactions, sm...)
	case DBPostgresJSON:
		generatePostgresJSONRepository(dir, entity, cache, transactions, sm...)
	case DBMySQL:
		generateMySQLRepository(dir, entity, cache, transactions, sm...)
	case DBMongoDB:
		generateMongoRepository(dir, entity, cache, transactions, sm...)
	case DBSQLite:
		// SQLite projects run on GORM (gorm.io/driver/sqlite) just like the
		// other SQL backends: the generated main.go/DI container hold a
		// *gorm.DB, so the repository must be GORM-based (NewPostgres<Entity>
		// Repository(*gorm.DB)) to wire correctly. A bespoke database/sql
		// implementation (*sql.DB) could never be injected by that container.
		generatePostgresRepository(dir, entity, cache, transactions, sm...)
	case DBSQLServer:
		generateSQLServerRepository(dir, entity, cache, transactions, sm...)
	case DBElasticsearch:
		generateElasticsearchRepository(dir, entity, cache, transactions, sm...)
	case DBDynamoDB:
		generateDynamoDBRepository(dir, entity, cache, transactions, sm...)
	default:
		fmt.Printf("Database not supported: %s\n", database)
		os.Exit(1)
	}
}

func init() {
	repositoryCmd.Flags().StringP(DatabaseFlag, "d", "", DatabaseFlagUsage)
	repositoryCmd.Flags().BoolP(InterfaceOnlyFlag, "i", false, InterfaceOnlyFlagUsage)
	repositoryCmd.Flags().BoolP(ImplementationFlag, "", false, ImplementationFlagUsage)
	repositoryCmd.Flags().BoolP(CacheFlag, "c", false, CacheFlagUsage)
	repositoryCmd.Flags().BoolP(TransactionsFlag, "t", false, TransactionsFlagUsage)
	repositoryCmd.Flags().StringP("fields", "f", "", "Entity fields \"field:type,field2:type\"")
	repositoryCmd.Flags().Bool("dry-run", false, "Preview changes without creating files")
	repositoryCmd.Flags().Bool("force", false, "Overwrite existing files without asking")
	repositoryCmd.Flags().Bool("backup", false, "Backup existing files before overwriting")
}
