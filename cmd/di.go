package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	dbPostgres = "postgres"
	dbMySQL    = "mysql"
	dbMongoDB  = "mongodb"
)

var diCmd = &cobra.Command{
	Use:   "di",
	Short: "Generate dependency injection container",
	Long: `Creates a dependency injection container that automatically connects 
all layers of the system using Google Wire.`,
	Run: func(cmd *cobra.Command, _ []string) {
		features, _ := cmd.Flags().GetString("features")
		database, _ := cmd.Flags().GetString("database")
		wire, _ := cmd.Flags().GetBool("wire")

		if features == "" {
			ui.Error("--features flag is required")
			os.Exit(1)
		}

		ui.Header(fmt.Sprintf("Generating DI container for features: %s", features))
		ui.KeyValue("Database", database)

		if wire {
			ui.Feature("Using Google Wire", false)
		}

		// Initialize safety manager
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		backup, _ := cmd.Flags().GetBool("backup")
		sm := NewSafetyManager(dryRun, force, backup)

		if dryRun {
			ui.DryRun("Previewing changes without creating files")
		}

		generateDI(features, database, wire, sm)

		if dryRun {
			sm.PrintSummary()
			return
		}

		ui.Success("DI container generated successfully!")
	},
}

func generateDI(features, database string, wire bool, sm ...*SafetyManager) {
	diDir := "internal/di"
	// Create di directory if it doesn't exist
	_ = os.MkdirAll(diDir, 0755)

	// Parse features
	featureList := strings.Split(features, ",")
	for i, feature := range featureList {
		featureList[i] = strings.TrimSpace(feature)
	}

	if wire {
		generateWireDI(diDir, featureList, database, sm...)
	} else {
		generateManualDI(diDir, featureList, database, sm...)
	}
}

func generateManualDI(dir string, features []string, database string, sm ...*SafetyManager) {
	// Get the module name from go.mod
	moduleName := getModuleName()

	// Use relative imports for local development and testing
	importPath := getImportPath(moduleName)

	filename := filepath.Join(dir, "container.go")

	var content strings.Builder
	content.WriteString("package di\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"gorm.io/gorm\"\n\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/repository\"\n", importPath))
	content.WriteString(fmt.Sprintf("\t\"%s/internal/usecase\"\n", importPath))
	content.WriteString(fmt.Sprintf("\t\"%s/internal/handler/http\"\n", importPath))
	content.WriteString(")\n\n") // Container struct
	content.WriteString("type Container struct {\n")
	content.WriteString("\tdb *gorm.DB\n\n")

	// Repositories
	content.WriteString("\t// Repositories\n")
	for _, feature := range features {
		featureLower := strings.ToLower(feature)
		content.WriteString(fmt.Sprintf("\t%sRepo    repository.%sRepository\n", featureLower, feature))
	}

	// Use Cases
	content.WriteString("\n\t// Use Cases\n")
	for _, feature := range features {
		featureLower := strings.ToLower(feature)
		content.WriteString(fmt.Sprintf("\t%sUC    usecase.%sUseCase\n", featureLower, feature))
	}

	// Handlers
	content.WriteString("\n\t// Handlers\n")
	for _, feature := range features {
		fieldName := strings.ToLower(feature[:1]) + feature[1:] // camelCase
		content.WriteString(fmt.Sprintf("\t%sHandler    *http.%sHandler\n", fieldName, feature))
	}

	content.WriteString("}\n\n")

	// Constructor
	content.WriteString("func NewContainer(db *gorm.DB) *Container {\n")
	content.WriteString("\tc := &Container{db: db}\n")
	content.WriteString("\tc.setupRepositories()\n")
	content.WriteString("\tc.setupUseCases()\n")
	content.WriteString("\tc.setupHandlers()\n")
	content.WriteString("\treturn c\n")
	content.WriteString("}\n\n")

	// Setup methods
	generateSetupRepositories(&content, features, database)
	generateSetupUseCases(&content, features)
	generateSetupHandlers(&content, features)

	// Getters
	generateGetters(&content, features)

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing DI file: %v", err))
		return
	}
}

func generateSetupRepositories(content *strings.Builder, features []string, database string) {
	content.WriteString("func (c *Container) setupRepositories() {\n")

	for _, feature := range features {
		featureLower := strings.ToLower(feature)
		switch database {
		case dbPostgres:
			fmt.Fprintf(content, "\tc.%sRepo = repository.NewPostgres%sRepository(c.db)\n",
				featureLower, feature)
		case dbMySQL:
			fmt.Fprintf(content, "\tc.%sRepo = repository.NewMySQL%sRepository(c.db)\n",
				featureLower, feature)
		case dbMongoDB:
			fmt.Fprintf(content, "\tc.%sRepo = repository.NewMongo%sRepository(c.db)\n",
				featureLower, feature)
		default:
			fmt.Fprintf(content, "\tc.%sRepo = repository.NewPostgres%sRepository(c.db)\n",
				featureLower, feature)
		}
	}

	content.WriteString("}\n\n")
}

func generateSetupUseCases(content *strings.Builder, features []string) {
	content.WriteString("func (c *Container) setupUseCases() {\n")

	for _, feature := range features {
		featureLower := strings.ToLower(feature)
		fmt.Fprintf(content, "\tc.%sUC = usecase.New%sService(c.%sRepo)\n",
			featureLower, feature, featureLower)
	}

	content.WriteString("}\n\n")
}

func generateSetupHandlers(content *strings.Builder, features []string) {
	content.WriteString("func (c *Container) setupHandlers() {\n")

	for _, feature := range features {
		fieldName := strings.ToLower(feature[:1]) + feature[1:] // camelCase
		featureLower := strings.ToLower(feature)
		fmt.Fprintf(content, "\tc.%sHandler = http.New%sHandler(c.%sUC)\n",
			fieldName, feature, featureLower)
	}

	content.WriteString("}\n\n")
}

func generateGetters(content *strings.Builder, features []string) {
	content.WriteString("// Getters\n")

	for _, feature := range features {
		fieldName := strings.ToLower(feature[:1]) + feature[1:] // camelCase
		featureLower := strings.ToLower(feature)

		// Handler getter
		fmt.Fprintf(content, "func (c *Container) %sHandler() *http.%sHandler {\n",
			feature, feature)
		fmt.Fprintf(content, "\treturn c.%sHandler\n", fieldName)
		content.WriteString("}\n\n")

		// UseCase getter
		fmt.Fprintf(content, "func (c *Container) %sUseCase() usecase.%sUseCase {\n",
			feature, feature)
		fmt.Fprintf(content, "\treturn c.%sUC\n", featureLower)
		content.WriteString("}\n\n")

		// Repository getter
		fmt.Fprintf(content, "func (c *Container) %sRepository() repository.%sRepository {\n",
			feature, feature)
		fmt.Fprintf(content, "\treturn c.%sRepo\n", featureLower)
		content.WriteString("}\n\n")
	}
}

func generateWireDI(dir string, features []string, database string, sm ...*SafetyManager) {
	// Generate wire.go
	generateWireFile(dir, features, database, sm...)

	// Generate wire_gen.go template
	generateWireGenTemplate(dir, features, sm...)
}

func generateWireFile(dir string, features []string, database string, sm ...*SafetyManager) {
	moduleName := getModuleName()
	importPath := getImportPath(moduleName)
	filename := filepath.Join(dir, "wire.go")

	var content strings.Builder
	writeWireHeader(&content)
	writeWireImports(&content, importPath)
	writeWireSets(&content, features, database)
	writeWireFunctions(&content, features)

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing Wire file: %v", err))
		return
	}
}

// writeWireHeader writes the Wire file header with build tags
func writeWireHeader(content *strings.Builder) {
	content.WriteString("//go:build wireinject\n")
	content.WriteString("// +build wireinject\n\n")
	content.WriteString("package di\n\n")
}

// writeWireImports writes the import section for Wire file
func writeWireImports(content *strings.Builder, importPath string) {
	content.WriteString("import (\n")
	content.WriteString("\t\"database/sql\"\n\n")
	content.WriteString("\t\"github.com/google/wire\"\n")
	fmt.Fprintf(content, "\t\"%s/internal/repository\"\n", importPath)
	fmt.Fprintf(content, "\t\"%s/internal/usecase\"\n", importPath)
	fmt.Fprintf(content, "\t\"%s/internal/handler/http\"\n", importPath)
	content.WriteString(")\n\n")
}

// writeWireSets writes all Wire sets (Repository, UseCase, Handler, All)
func writeWireSets(content *strings.Builder, features []string, database string) {
	content.WriteString("// Wire sets\n")
	content.WriteString("var (\n")

	writeRepositorySet(content, features, database)
	writeUseCaseSet(content, features)
	writeHandlerSet(content, features)
	writeAllSet(content)

	content.WriteString(")\n\n")
}

// writeRepositorySet writes the Repository Wire set
func writeRepositorySet(content *strings.Builder, features []string, database string) {
	content.WriteString("\tRepositorySet = wire.NewSet(\n")
	for _, feature := range features {
		switch database {
		case dbPostgres:
			fmt.Fprintf(content, "\t\trepository.NewPostgres%sRepository,\n", feature)
		case dbMySQL:
			fmt.Fprintf(content, "\t\trepository.NewMySQL%sRepository,\n", feature)
		case dbMongoDB:
			fmt.Fprintf(content, "\t\trepository.NewMongo%sRepository,\n", feature)
		default:
			fmt.Fprintf(content, "\t\trepository.NewPostgres%sRepository,\n", feature)
		}
	}
	content.WriteString("\t)\n\n")
}

// writeUseCaseSet writes the UseCase Wire set
func writeUseCaseSet(content *strings.Builder, features []string) {
	content.WriteString("\tUseCaseSet = wire.NewSet(\n")
	for _, feature := range features {
		fmt.Fprintf(content, "\t\tusecase.New%sService,\n", feature)
	}
	content.WriteString("\t)\n\n")
}

// writeHandlerSet writes the Handler Wire set
func writeHandlerSet(content *strings.Builder, features []string) {
	content.WriteString("\tHandlerSet = wire.NewSet(\n")
	for _, feature := range features {
		fmt.Fprintf(content, "\t\thttp.New%sHandler,\n", feature)
	}
	content.WriteString("\t)\n\n")
}

// writeAllSet writes the combined All Wire set
func writeAllSet(content *strings.Builder) {
	content.WriteString("\tAllSet = wire.NewSet(\n")
	content.WriteString("\t\tRepositorySet,\n")
	content.WriteString("\t\tUseCaseSet,\n")
	content.WriteString("\t\tHandlerSet,\n")
	content.WriteString("\t)\n")
}

// writeWireFunctions writes Wire initialization functions
func writeWireFunctions(content *strings.Builder, features []string) {
	for _, feature := range features {
		fmt.Fprintf(content, "func Initialize%sHandler(db *sql.DB) *http.%sHandler {\n",
			feature, feature)
		content.WriteString("\twire.Build(AllSet)\n")
		fmt.Fprintf(content, "\treturn &http.%sHandler{}\n", feature)
		content.WriteString("}\n\n")
	}

	content.WriteString("func InitializeContainer(db *sql.DB) *Container {\n")
	content.WriteString("\twire.Build(\n")
	content.WriteString("\t\tAllSet,\n")
	content.WriteString("\t\tNewWireContainer,\n")
	content.WriteString("\t)\n")
	content.WriteString("\treturn &Container{}\n")
	content.WriteString("}\n")
}

func generateWireGenTemplate(dir string, features []string, sm ...*SafetyManager) {
	// Get the module name from go.mod
	moduleName := getModuleName()

	// Use relative imports for local development and testing
	importPath := getImportPath(moduleName)

	filename := filepath.Join(dir, "wire_container.go")

	var content strings.Builder
	content.WriteString("package di\n\n")
	content.WriteString("import (\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/handler/http\"\n", importPath))
	content.WriteString(")\n\n") // Wire container struct
	content.WriteString("type Container struct {\n")
	for _, feature := range features {
		featureLower := strings.ToLower(feature)
		content.WriteString(fmt.Sprintf("\t%sHandler *http.%sHandler\n", featureLower, feature))
	}
	content.WriteString("}\n\n")

	// Constructor
	content.WriteString("func NewWireContainer(\n")
	for _, feature := range features {
		featureLower := strings.ToLower(feature)
		content.WriteString(fmt.Sprintf("\t%sHandler *http.%sHandler,\n", featureLower, feature))
	}
	content.WriteString(") *Container {\n")
	content.WriteString("\treturn &Container{\n")
	for _, feature := range features {
		featureLower := strings.ToLower(feature)
		content.WriteString(fmt.Sprintf("\t\t%sHandler: %sHandler,\n", featureLower, featureLower))
	}
	content.WriteString("\t}\n")
	content.WriteString("}\n\n")

	// Getters
	for _, feature := range features {
		featureLower := strings.ToLower(feature)
		content.WriteString(fmt.Sprintf("func (c *Container) %sHandler() *http.%sHandler {\n",
			feature, feature))
		content.WriteString(fmt.Sprintf("\treturn c.%sHandler\n", featureLower))
		content.WriteString("}\n\n")
	}

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing Wire generator file: %v", err))
		return
	}
}

func init() {
	diCmd.Flags().StringP("features", "f", "", "Project features (crud,auth,validation,etc)")
	diCmd.Flags().StringP("database", "d", "postgres", "Database type (postgres, mysql, mongodb)")
	diCmd.Flags().BoolP("wire", "w", false, "Use Google Wire for dependency injection")
	diCmd.Flags().Bool("dry-run", false, "Preview changes without creating files")
	diCmd.Flags().Bool("force", false, "Overwrite existing files without asking")
	diCmd.Flags().Bool("backup", false, "Backup existing files before overwriting")
	_ = diCmd.MarkFlagRequired("features")
}
