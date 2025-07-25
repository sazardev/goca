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
	Short: "Generar contenedor de inyección de dependencias",
	Long: `Crea un contenedor de inyección de dependencias que conecta 
automáticamente todas las capas del sistema.`,
	Run: func(cmd *cobra.Command, _ []string) {
		features, _ := cmd.Flags().GetString("features")
		database, _ := cmd.Flags().GetString("database")
		wire, _ := cmd.Flags().GetBool("wire")

		if features == "" {
			fmt.Println("Error: --features flag es requerido")
			os.Exit(1)
		}

		fmt.Printf("Generando contenedor DI para features: %s\n", features)
		fmt.Printf("Base de datos: %s\n", database)

		if wire {
			fmt.Println("✓ Usando Google Wire")
		}

		generateDI(features, database, wire)
		fmt.Printf("\n✅ Contenedor DI generado exitosamente!\n")
	},
}

func generateDI(features, database string, wire bool) {
	diDir := "internal/di"
	// Crear directorio di si no existe
	_ = os.MkdirAll(diDir, 0755)

	// Parse features
	featureList := strings.Split(features, ",")
	for i, feature := range featureList {
		featureList[i] = strings.TrimSpace(feature)
	}

	if wire {
		generateWireDI(diDir, featureList, database)
	} else {
		generateManualDI(diDir, featureList, database)
	}
}

func generateManualDI(dir string, features []string, database string) {
	// Get the module name from go.mod
	moduleName := getModuleName()

	// Use relative imports for local development and testing
	importPath := getImportPath(moduleName)

	filename := filepath.Join(dir, "container.go")

	var content strings.Builder
	content.WriteString("package di\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"database/sql\"\n")
	content.WriteString("\t\"log\"\n\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/repository\"\n", importPath))
	content.WriteString(fmt.Sprintf("\t\"%s/internal/usecase\"\n", importPath))
	content.WriteString(fmt.Sprintf("\t\"%s/internal/handler/http\"\n", importPath))
	content.WriteString(")\n\n") // Container struct
	content.WriteString("type Container struct {\n")
	content.WriteString("\tdb *sql.DB\n\n")

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
		featureLower := strings.ToLower(feature)
		content.WriteString(fmt.Sprintf("\t%sHandler    *http.%sHandler\n", featureLower, feature))
	}

	content.WriteString("}\n\n")

	// Constructor
	content.WriteString("func NewContainer(db *sql.DB) *Container {\n")
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

	writeGoFile(filename, content.String())
}

func generateSetupRepositories(content *strings.Builder, features []string, database string) {
	content.WriteString("func (c *Container) setupRepositories() {\n")

	for _, feature := range features {
		featureLower := strings.ToLower(feature)
		switch database {
		case dbPostgres:
			content.WriteString(fmt.Sprintf("\tc.%sRepo = repository.NewPostgres%sRepository(c.db)\n",
				featureLower, feature))
		case dbMySQL:
			content.WriteString(fmt.Sprintf("\tc.%sRepo = repository.NewMySQL%sRepository(c.db)\n",
				featureLower, feature))
		case dbMongoDB:
			content.WriteString(fmt.Sprintf("\tc.%sRepo = repository.NewMongo%sRepository(c.db)\n",
				featureLower, feature))
		default:
			content.WriteString(fmt.Sprintf("\tc.%sRepo = repository.NewPostgres%sRepository(c.db)\n",
				featureLower, feature))
		}
	}

	content.WriteString("}\n\n")
}

func generateSetupUseCases(content *strings.Builder, features []string) {
	content.WriteString("func (c *Container) setupUseCases() {\n")

	for _, feature := range features {
		featureLower := strings.ToLower(feature)
		content.WriteString(fmt.Sprintf("\tc.%sUC = usecase.New%sService(c.%sRepo)\n",
			featureLower, feature, featureLower))
	}

	content.WriteString("}\n\n")
}

func generateSetupHandlers(content *strings.Builder, features []string) {
	content.WriteString("func (c *Container) setupHandlers() {\n")

	for _, feature := range features {
		featureLower := strings.ToLower(feature)
		content.WriteString(fmt.Sprintf("\tc.%sHandler = http.New%sHandler(c.%sUC)\n",
			featureLower, feature, featureLower))
	}

	content.WriteString("}\n\n")
}

func generateGetters(content *strings.Builder, features []string) {
	content.WriteString("// Getters\n")

	for _, feature := range features {
		featureLower := strings.ToLower(feature)

		// Handler getter
		content.WriteString(fmt.Sprintf("func (c *Container) %sHandler() *http.%sHandler {\n",
			feature, feature))
		content.WriteString(fmt.Sprintf("\treturn c.%sHandler\n", featureLower))
		content.WriteString("}\n\n")

		// UseCase getter
		content.WriteString(fmt.Sprintf("func (c *Container) %sUseCase() usecase.%sUseCase {\n",
			feature, feature))
		content.WriteString(fmt.Sprintf("\treturn c.%sUC\n", featureLower))
		content.WriteString("}\n\n")

		// Repository getter
		content.WriteString(fmt.Sprintf("func (c *Container) %sRepository() repository.%sRepository {\n",
			feature, feature))
		content.WriteString(fmt.Sprintf("\treturn c.%sRepo\n", featureLower))
		content.WriteString("}\n\n")
	}
}

func generateWireDI(dir string, features []string, database string) {
	// Generate wire.go
	generateWireFile(dir, features, database)

	// Generate wire_gen.go template
	generateWireGenTemplate(dir, features)
}

func generateWireFile(dir string, features []string, database string) {
	// Get the module name from go.mod
	moduleName := getModuleName()

	// Use relative imports for local development and testing
	importPath := getImportPath(moduleName)

	filename := filepath.Join(dir, "wire.go")

	var content strings.Builder
	content.WriteString("//go:build wireinject\n")
	content.WriteString("// +build wireinject\n\n")
	content.WriteString("package di\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"database/sql\"\n\n")
	content.WriteString("\t\"github.com/google/wire\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/repository\"\n", importPath))
	content.WriteString(fmt.Sprintf("\t\"%s/internal/usecase\"\n", importPath))
	content.WriteString(fmt.Sprintf("\t\"%s/internal/handler/http\"\n", importPath))
	content.WriteString(")\n\n") // Wire sets
	content.WriteString("// Wire sets\n")
	content.WriteString("var (\n")

	// Repository set
	content.WriteString("\tRepositorySet = wire.NewSet(\n")
	for _, feature := range features {
		switch database {
		case dbPostgres:
			content.WriteString(fmt.Sprintf("\t\trepository.NewPostgres%sRepository,\n", feature))
		case dbMySQL:
			content.WriteString(fmt.Sprintf("\t\trepository.NewMySQL%sRepository,\n", feature))
		case dbMongoDB:
			content.WriteString(fmt.Sprintf("\t\trepository.NewMongo%sRepository,\n", feature))
		default:
			content.WriteString(fmt.Sprintf("\t\trepository.NewPostgres%sRepository,\n", feature))
		}
	}
	content.WriteString("\t)\n\n")

	// UseCase set
	content.WriteString("\tUseCaseSet = wire.NewSet(\n")
	for _, feature := range features {
		content.WriteString(fmt.Sprintf("\t\tusecase.New%sService,\n", feature))
	}
	content.WriteString("\t)\n\n")

	// Handler set
	content.WriteString("\tHandlerSet = wire.NewSet(\n")
	for _, feature := range features {
		content.WriteString(fmt.Sprintf("\t\thttp.New%sHandler,\n", feature))
	}
	content.WriteString("\t)\n\n")

	// All set
	content.WriteString("\tAllSet = wire.NewSet(\n")
	content.WriteString("\t\tRepositorySet,\n")
	content.WriteString("\t\tUseCaseSet,\n")
	content.WriteString("\t\tHandlerSet,\n")
	content.WriteString("\t)\n")
	content.WriteString(")\n\n")

	// Wire functions
	for _, feature := range features {
		content.WriteString(fmt.Sprintf("func Initialize%sHandler(db *sql.DB) *http.%sHandler {\n",
			feature, feature))
		content.WriteString("\twire.Build(AllSet)\n")
		content.WriteString(fmt.Sprintf("\treturn &http.%sHandler{}\n", feature))
		content.WriteString("}\n\n")
	}

	// Main container function
	content.WriteString("func InitializeContainer(db *sql.DB) *Container {\n")
	content.WriteString("\twire.Build(\n")
	content.WriteString("\t\tAllSet,\n")
	content.WriteString("\t\tNewWireContainer,\n")
	content.WriteString("\t)\n")
	content.WriteString("\treturn &Container{}\n")
	content.WriteString("}\n")

	writeGoFile(filename, content.String())
}

func generateWireGenTemplate(dir string, features []string) {
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

	writeGoFile(filename, content.String())
}

func init() {
	rootCmd.AddCommand(diCmd)
	diCmd.Flags().StringP("features", "f", "", "Características del proyecto (crud,auth,validation,etc)")
	diCmd.Flags().StringP("database", "d", "postgres", "Tipo de base de datos (postgres, mysql, mongodb)")
	diCmd.Flags().BoolP("wire", "w", false, "Usar Google Wire para inyección de dependencias")
	_ = diCmd.MarkFlagRequired("features")
}
