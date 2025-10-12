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
			fmt.Printf("Warning: Could not load configuration: %v\n", err)
			fmt.Println("Using default values. Consider running 'goca init --config' to generate .goca.yaml")
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

		// Validar con el nuevo validador robusto
		validator := NewFieldValidator()

		if err := validator.ValidateEntityName(entity); err != nil {
			fmt.Printf("Error in entity name: %v\n", err)
			return
		}

		if effectiveDatabase != "" {
			if err := validator.ValidateDatabase(effectiveDatabase); err != nil {
				fmt.Printf("Error in database: %v\n", err)
				return
			}
		}

		fmt.Printf("Generating repository for entity '%s'\n", entity)

		if effectiveDatabase != "" && !interfaceOnly {
			fmt.Printf("Database: %s", effectiveDatabase)
			if configIntegration.HasConfigFile() && !cmd.Flags().Changed("database") {
				fmt.Printf(" (from config)")
			}
			fmt.Println()
		}
		if interfaceOnly {
			fmt.Println("✓ Interface only")
		}
		if implementation {
			fmt.Println("✓ Implementation only")
		}
		if cache {
			fmt.Println("✓ Including cache")
		}
		if transactions {
			fmt.Println("✓ Including transactions")
		}
		if fields != "" {
			fmt.Printf("✓ Custom fields: %s\n", fields)
		}

		if configIntegration.HasConfigFile() {
			configIntegration.PrintConfigSummary()
		}

		generateRepository(entity, effectiveDatabase, interfaceOnly, implementation, cache, transactions, fields)
		fmt.Printf("\nRepository for '%s' generated successfully!\n", entity)
	},
}

func generateRepository(entity, database string, interfaceOnly, implementation, cache, transactions bool, fields string) {
	// Crear directorio repositories si no existe
	repoDir := "internal/repository"
	_ = os.MkdirAll(repoDir, 0755)

	// Parse fields if provided
	var parsedFields []Field
	if fields != "" {
		parsedFields = parseFields(fields)
	}

	// Generate interface:
	// - Always generate interface UNLESS explicitly skipped
	// - interfaceOnly=true: only interface
	// - interfaceOnly=false, implementation=true: both
	// - interfaceOnly=false, implementation=false: both (default)
	if fields != "" {
		generateRepositoryInterfaceWithFields(repoDir, entity, parsedFields, transactions)
	} else {
		generateRepositoryInterface(repoDir, entity, transactions)
	}

	// Generate implementation if not interface-only and database is specified
	if !interfaceOnly && database != "" {
		if fields != "" {
			generateRepositoryImplementationWithFields(repoDir, entity, database, parsedFields, cache, transactions)
		} else {
			generateRepositoryImplementation(repoDir, entity, database, cache, transactions)
		}
	}
}

func generateRepositoryInterface(dir, entity string, transactions bool) {
	filename := filepath.Join(dir, "interfaces.go")

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
		content.WriteString(fmt.Sprintf("import \"%s/internal/domain\"\n\n", getImportPath(moduleName)))
	}

	content.WriteString(fmt.Sprintf("type %sRepository interface {\n", entity))
	content.WriteString(fmt.Sprintf("\tSave(%s *domain.%s) error\n", strings.ToLower(entity), entity))
	content.WriteString(fmt.Sprintf("\tFindByID(id int) (*domain.%s, error)\n", entity))
	content.WriteString(fmt.Sprintf("\tFindByEmail(email string) (*domain.%s, error)\n", entity))
	content.WriteString(fmt.Sprintf("\tUpdate(%s *domain.%s) error\n", strings.ToLower(entity), entity))
	content.WriteString("\tDelete(id int) error\n")
	content.WriteString(fmt.Sprintf("\tFindAll() ([]domain.%s, error)\n", entity))

	if transactions {
		content.WriteString(fmt.Sprintf("\tSaveWithTx(tx interface{}, %s *domain.%s) error\n", strings.ToLower(entity), entity))
		content.WriteString(fmt.Sprintf("\tUpdateWithTx(tx interface{}, %s *domain.%s) error\n", strings.ToLower(entity), entity))
		content.WriteString("\tDeleteWithTx(tx interface{}, id int) error\n")
	}

	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error writing file %s: %v\n", filename, err)
	}
}

func generateRepositoryImplementation(dir, entity, database string, cache, transactions bool) {
	switch database {
	case DBPostgres:
		generatePostgresRepository(dir, entity, cache, transactions)
	case DBMySQL:
		generateMySQLRepository(dir, entity, cache, transactions)
	case DBMongoDB:
		generateMongoRepository(dir, entity, cache, transactions)
	default:
		fmt.Printf("Database not supported: %s\n", database)
		os.Exit(1)
	}
}

func generatePostgresRepository(dir, entity string, cache, transactions bool) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, "postgres_"+entityLower+"_repository.go")

	// Get the module name from go.mod
	moduleName := getModuleName()

	var content strings.Builder
	content.WriteString("package repository\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"gorm.io/gorm\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
	if cache {
		content.WriteString("\t\"time\"\n")
		content.WriteString("\t\"encoding/json\"\n")
		content.WriteString("\t\"github.com/go-redis/redis/v8\"\n")
		content.WriteString("\t\"context\"\n")
	}
	content.WriteString(")\n\n")

	// Repository struct
	repoName := fmt.Sprintf("postgres%sRepository", entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n", repoName))
	content.WriteString("\tdb *gorm.DB\n")
	if cache {
		content.WriteString("\tcache *redis.Client\n")
		content.WriteString("\tcacheTTL time.Duration\n")
	}
	content.WriteString("}\n\n")

	// Constructor
	content.WriteString(fmt.Sprintf("func NewPostgres%sRepository(db *gorm.DB", entity))
	if cache {
		content.WriteString(", cache *redis.Client")
	}
	content.WriteString(fmt.Sprintf(") %sRepository {\n", entity))
	content.WriteString(fmt.Sprintf("\treturn &%s{\n", repoName))
	content.WriteString("\t\tdb: db,\n")
	if cache {
		content.WriteString("\t\tcache: cache,\n")
		content.WriteString("\t\tcacheTTL: 5 * time.Minute,\n")
	}
	content.WriteString("\t}\n")
	content.WriteString("}\n\n")

	// Generate methods
	generatePostgresSaveMethod(&content, entity, repoName, cache)
	generatePostgresFindByIDMethod(&content, entity, repoName, cache)
	generatePostgresFindByEmailMethod(&content, entity, repoName)
	generatePostgresUpdateMethod(&content, entity, repoName, cache)
	generatePostgresDeleteMethod(&content, entity, repoName, cache)
	generatePostgresFindAllMethod(&content, entity, repoName)

	if transactions {
		generatePostgresTransactionMethods(&content, entity, repoName)
	}

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating PostgreSQL repository file: %v\n", err)
	}
}

func generatePostgresSaveMethod(content *strings.Builder, entity, repoName string, cache bool) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	fmt.Fprintf(content, "func (%s *%s) Save(%s *domain.%s) error {\n",
		repoVar, repoName, entityLower, entity)
	fmt.Fprintf(content, "\tresult := %s.db.Create(%s)\n", repoVar, entityLower)

	if cache {
		content.WriteString("\tif result.Error == nil {\n")
		fmt.Fprintf(content, "\t\t%s.invalidateCache(%s.ID)\n", repoVar, entityLower)
		content.WriteString("\t}\n")
	}

	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")
}

func generatePostgresFindByIDMethod(content *strings.Builder, entity, repoName string, cache bool) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	fmt.Fprintf(content, "func (%s *%s) FindByID(id int) (*domain.%s, error) {\n",
		repoVar, repoName, entity)

	if cache {
		content.WriteString("\t// Try cache first\n")
		fmt.Fprintf(content, "\tif %s := %s.getFromCache(id); %s != nil {\n",
			entityLower, repoVar, entityLower)
		fmt.Fprintf(content, "\t\treturn %s, nil\n", entityLower)
		content.WriteString("\t}\n\n")
	}

	fmt.Fprintf(content, "\t%s := &domain.%s{}\n", entityLower, entity)
	fmt.Fprintf(content, "\tresult := %s.db.First(%s, id)\n", repoVar, entityLower)
	content.WriteString("\tif result.Error != nil {\n")
	content.WriteString("\t\treturn nil, result.Error\n")
	content.WriteString("\t}\n\n")

	if cache {
		fmt.Fprintf(content, "\t%s.setCache(%s)\n", repoVar, entityLower)
	}

	fmt.Fprintf(content, "\treturn %s, nil\n", entityLower)
	content.WriteString("}\n\n")
}

func generatePostgresFindByEmailMethod(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	fmt.Fprintf(content, "func (%s *%s) FindByEmail(email string) (*domain.%s, error) {\n",
		repoVar, repoName, entity)
	fmt.Fprintf(content, "\t%s := &domain.%s{}\n", entityLower, entity)
	fmt.Fprintf(content, "\tresult := %s.db.Where(\"email = ?\", email).First(%s)\n", repoVar, entityLower)
	content.WriteString("\tif result.Error != nil {\n")
	content.WriteString("\t\treturn nil, result.Error\n")
	content.WriteString("\t}\n")
	fmt.Fprintf(content, "\treturn %s, nil\n", entityLower)
	content.WriteString("}\n\n")
}

func generatePostgresUpdateMethod(content *strings.Builder, entity, repoName string, cache bool) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	fmt.Fprintf(content, "func (%s *%s) Update(%s *domain.%s) error {\n",
		repoVar, repoName, entityLower, entity)
	fmt.Fprintf(content, "\tresult := %s.db.Save(%s)\n", repoVar, entityLower)

	if cache {
		content.WriteString("\tif result.Error == nil {\n")
		fmt.Fprintf(content, "\t\t%s.invalidateCache(%s.ID)\n", repoVar, entityLower)
		content.WriteString("\t}\n")
	}

	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")
}

func generatePostgresDeleteMethod(content *strings.Builder, entity, repoName string, cache bool) {
	repoVar := strings.ToLower(string(repoName[0]))

	fmt.Fprintf(content, "func (%s *%s) Delete(id int) error {\n",
		repoVar, repoName)
	fmt.Fprintf(content, "\tresult := %s.db.Delete(&domain.%s{}, id)\n", repoVar, entity)

	if cache {
		content.WriteString("\tif result.Error == nil {\n")
		fmt.Fprintf(content, "\t\t%s.invalidateCache(id)\n", repoVar)
		content.WriteString("\t}\n")
	}

	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")
}

func generatePostgresFindAllMethod(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	fmt.Fprintf(content, "func (%s *%s) FindAll() ([]domain.%s, error) {\n",
		repoVar, repoName, entity)
	fmt.Fprintf(content, "\tvar %ss []domain.%s\n", entityLower, entity)
	fmt.Fprintf(content, "\tresult := %s.db.Find(&%ss)\n", repoVar, entityLower)
	content.WriteString("\tif result.Error != nil {\n")
	content.WriteString("\t\treturn nil, result.Error\n")
	content.WriteString("\t}\n\n")

	fmt.Fprintf(content, "\treturn %ss, nil\n", entityLower)
	content.WriteString("}\n\n")
}

func generatePostgresTransactionMethods(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	// SaveWithTx
	fmt.Fprintf(content, "func (%s *%s) SaveWithTx(tx interface{}, %s *domain.%s) error {\n",
		repoVar, repoName, entityLower, entity)
	content.WriteString("\tgormTx := tx.(*gorm.DB)\n")
	fmt.Fprintf(content, "\tresult := gormTx.Create(%s)\n", entityLower)
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// UpdateWithTx
	fmt.Fprintf(content, "func (%s *%s) UpdateWithTx(tx interface{}, %s *domain.%s) error {\n",
		repoVar, repoName, entityLower, entity)
	content.WriteString("\tgormTx := tx.(*gorm.DB)\n")
	fmt.Fprintf(content, "\tresult := gormTx.Save(%s)\n", entityLower)
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// DeleteWithTx
	fmt.Fprintf(content, "func (%s *%s) DeleteWithTx(tx interface{}, id int) error {\n",
		repoVar, repoName)
	content.WriteString("\tgormTx := tx.(*gorm.DB)\n")
	fmt.Fprintf(content, "\tresult := gormTx.Delete(&domain.%s{}, id)\n", entity)
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")
}

func generateMySQLRepository(dir, entity string, cache, transactions bool) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, "mysql_"+entityLower+"_repository.go")

	// Get the module name from go.mod
	moduleName := getModuleName()

	var content strings.Builder
	content.WriteString("package repository\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"gorm.io/gorm\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
	if cache {
		content.WriteString("\t\"time\"\n")
		content.WriteString("\t\"encoding/json\"\n")
		content.WriteString("\t\"github.com/go-redis/redis/v8\"\n")
		content.WriteString("\t\"context\"\n")
	}
	content.WriteString(")\n\n")

	// MySQL repository structure (using GORM)
	repoName := fmt.Sprintf("mysql%sRepository", entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n", repoName))
	content.WriteString("\tdb *gorm.DB\n")
	if cache {
		content.WriteString("\tredis *redis.Client\n")
	}
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("func NewMySQL%sRepository(db *gorm.DB", entity))
	if cache {
		content.WriteString(", redis *redis.Client")
	}
	content.WriteString(fmt.Sprintf(") %sRepository {\n", entity))
	content.WriteString(fmt.Sprintf("\treturn &%s{\n", repoName))
	content.WriteString("\t\tdb: db,\n")
	if cache {
		content.WriteString("\t\tredis: redis,\n")
	}
	content.WriteString("\t}\n")
	content.WriteString("}\n\n")

	// Save method
	content.WriteString(fmt.Sprintf("func (r *%s) Save(%s *domain.%s) error {\n",
		repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := r.db.Create(%s)\n", entityLower))
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// FindByID method
	content.WriteString(fmt.Sprintf("func (r *%s) FindByID(id int) (*domain.%s, error) {\n",
		repoName, entity))
	content.WriteString(fmt.Sprintf("\t%s := &domain.%s{}\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := r.db.First(%s, id)\n", entityLower))
	content.WriteString("\tif result.Error != nil {\n")
	content.WriteString("\t\treturn nil, result.Error\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\treturn %s, nil\n", entityLower))
	content.WriteString("}\n\n")

	// FindByEmail method
	content.WriteString(fmt.Sprintf("func (r *%s) FindByEmail(email string) (*domain.%s, error) {\n",
		repoName, entity))
	content.WriteString(fmt.Sprintf("\t%s := &domain.%s{}\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := r.db.Where(\"email = ?\", email).First(%s)\n", entityLower))
	content.WriteString("\tif result.Error != nil {\n")
	content.WriteString("\t\treturn nil, result.Error\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\treturn %s, nil\n", entityLower))
	content.WriteString("}\n\n")

	// Update method
	content.WriteString(fmt.Sprintf("func (r *%s) Update(%s *domain.%s) error {\n",
		repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := r.db.Save(%s)\n", entityLower))
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// Delete method
	content.WriteString(fmt.Sprintf("func (r *%s) Delete(id int) error {\n", repoName))
	content.WriteString(fmt.Sprintf("\tresult := r.db.Delete(&domain.%s{}, id)\n", entity))
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// FindAll method
	content.WriteString(fmt.Sprintf("func (r *%s) FindAll() ([]domain.%s, error) {\n",
		repoName, entity))
	content.WriteString(fmt.Sprintf("\tvar %ss []domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := r.db.Find(&%ss)\n", entityLower))
	content.WriteString("\tif result.Error != nil {\n")
	content.WriteString("\t\treturn nil, result.Error\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\treturn %ss, nil\n", entityLower))
	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating MySQL repository file: %v\n", err)
	}
}

func generateMongoRepository(dir, entity string, cache, transactions bool) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, "mongo_"+entityLower+"_repository.go")

	// Get the module name from go.mod
	moduleName := getModuleName()

	var content strings.Builder
	content.WriteString("package repository\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"context\"\n")
	content.WriteString("\t\"time\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
	if cache {
		content.WriteString("\t// Imports para cache de MongoDB\n")
		content.WriteString("\t// \"github.com/go-redis/redis/v8\"\n")
	}
	if transactions {
		content.WriteString("\t// Soporte para transacciones de MongoDB\n")
		content.WriteString("\t// \"go.mongodb.org/mongo-driver/mongo/options\"\n")
		content.WriteString("\t// \"go.mongodb.org/mongo-driver/mongo/writeconcern\"\n")
	}
	content.WriteString("\n\t\"go.mongodb.org/mongo-driver/mongo\"\n")
	content.WriteString("\t\"go.mongodb.org/mongo-driver/bson\"\n")
	content.WriteString("\t\"go.mongodb.org/mongo-driver/bson/primitive\"\n")
	content.WriteString(")\n\n")

	// MongoDB repository structure
	repoName := fmt.Sprintf("mongo%sRepository", entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n", repoName))
	content.WriteString("\tcollection *mongo.Collection\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("func NewMongo%sRepository(db *mongo.Database) %sRepository {\n", entity, entity))
	content.WriteString(fmt.Sprintf("\treturn &%s{\n", repoName))
	content.WriteString(fmt.Sprintf("\t\tcollection: db.Collection(\"%ss\"),\n", entityLower))
	content.WriteString("\t}\n")
	content.WriteString("}\n\n")

	// Basic Save method for MongoDB
	content.WriteString(fmt.Sprintf("func (r *%s) Save(%s *domain.%s) error {\n",
		repoName, entityLower, entity))
	content.WriteString("\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n")
	content.WriteString("\tdefer cancel()\n\n")
	content.WriteString(fmt.Sprintf("\tresult, err := r.collection.InsertOne(ctx, %s)\n", entityLower))
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn err\n")
	content.WriteString("\t}\n\n")
	content.WriteString("\tif oid, ok := result.InsertedID.(primitive.ObjectID); ok {\n")
	content.WriteString(fmt.Sprintf("\t\t%s.ID = int(oid.Timestamp().Unix())\n", entityLower))
	content.WriteString("\t}\n")
	content.WriteString("\treturn nil\n")
	content.WriteString("}\n\n")

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating MongoDB repository file: %v\n", err)
	}
}

// generateRepositoryInterfaceWithFields genera interfaces de repository con métodos dinámicos basados en campos
func generateRepositoryInterfaceWithFields(dir, entity string, fields []Field, transactions bool) {
	filename := filepath.Join(dir, "interfaces.go")

	// Get the module name from go.mod
	moduleName := getModuleName()

	var content strings.Builder

	// Check if file exists to preserve existing interfaces
	if _, err := os.Stat(filename); err == nil {
		existingContent, err := os.ReadFile(filename)
		if err == nil {
			existingStr := string(existingContent)
			// Remove the final package declaration to avoid duplication
			if !strings.Contains(existingStr, fmt.Sprintf("type %sRepository interface", entity)) {
				// Add the existing content without the final newline
				content.WriteString(strings.TrimSuffix(existingStr, "\n"))
				content.WriteString("\n\n")
			}
		}
	} else {
		// File doesn't exist, create header
		content.WriteString("package repository\n\n")
		content.WriteString(fmt.Sprintf("import \"%s/internal/domain\"\n\n", getImportPath(moduleName)))
	}

	content.WriteString(fmt.Sprintf("type %sRepository interface {\n", entity))
	content.WriteString(fmt.Sprintf("\tSave(%s *domain.%s) error\n", strings.ToLower(entity), entity))
	content.WriteString(fmt.Sprintf("\tFindByID(id int) (*domain.%s, error)\n", entity))

	// Generar métodos de búsqueda dinámicos basados en los campos reales
	searchMethods := generateSearchMethods(fields, entity)
	for _, method := range searchMethods {
		content.WriteString(method.generateSearchMethodSignature() + "\n")
	}

	content.WriteString(fmt.Sprintf("\tUpdate(%s *domain.%s) error\n", strings.ToLower(entity), entity))
	content.WriteString("\tDelete(id int) error\n")
	content.WriteString(fmt.Sprintf("\tFindAll() ([]domain.%s, error)\n", entity))

	if transactions {
		content.WriteString(fmt.Sprintf("\tSaveWithTx(tx interface{}, %s *domain.%s) error\n", strings.ToLower(entity), entity))
		content.WriteString(fmt.Sprintf("\tUpdateWithTx(tx interface{}, %s *domain.%s) error\n", strings.ToLower(entity), entity))
		content.WriteString("\tDeleteWithTx(tx interface{}, id int) error\n")
	}

	content.WriteString("}\n\n")

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error writing file %s: %v\n", filename, err)
	}
}

// generateRepositoryImplementationWithFields generates repository implementations with dynamic methods
func generateRepositoryImplementationWithFields(dir, entity, database string, fields []Field, cache, transactions bool) {
	switch database {
	case DBPostgres:
		generatePostgresRepositoryWithFields(dir, entity, fields, cache, transactions)
	case DBMySQL:
		generateMySQLRepositoryWithFields(dir, entity, fields, cache, transactions)
	case DBMongoDB:
		generateMongoRepositoryWithFields(dir, entity, fields, cache, transactions)
	default:
		generatePostgresRepositoryWithFields(dir, entity, fields, cache, transactions)
	}
}

// generatePostgresRepositoryWithFields genera repository PostgreSQL con métodos dinámicos
func generatePostgresRepositoryWithFields(dir, entity string, fields []Field, cache, transactions bool) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, "postgres_"+entityLower+"_repository.go")

	// Get the module name from go.mod
	moduleName := getModuleName()

	var content strings.Builder
	content.WriteString("package repository\n\n")
	content.WriteString("import (\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
	if cache {
		content.WriteString("\t// Imports para cache (Redis, etc.)\n")
		content.WriteString("\t// \"github.com/go-redis/redis/v8\"\n")
	}
	if transactions {
		content.WriteString("\t// Soporte para transacciones SQL\n")
		content.WriteString("\t// \"database/sql/driver\"\n")
	}
	content.WriteString("\n\t\"gorm.io/gorm\"\n")
	content.WriteString(")\n\n")

	// Repository structure
	repoName := fmt.Sprintf("postgres%sRepository", entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n", repoName))
	content.WriteString("\tdb *gorm.DB\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("func NewPostgres%sRepository(db *gorm.DB) %sRepository {\n", entity, entity))
	content.WriteString(fmt.Sprintf("\treturn &%s{\n", repoName))
	content.WriteString("\t\tdb: db,\n")
	content.WriteString("\t}\n")
	content.WriteString("}\n\n")

	// Generate basic CRUD methods
	generateBasicCRUDMethods(&content, entity, repoName)

	// Generate dynamic search methods based on fields
	searchMethods := generateSearchMethods(fields, entity)
	for _, method := range searchMethods {
		content.WriteString(method.generateSearchMethodImplementation("p", repoName, entity))
	}

	if transactions {
		generateTransactionMethods(&content, entity, repoName)
	}

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating PostgreSQL repository with fields: %v\n", err)
	}
}

// generateBasicCRUDMethods genera los métodos CRUD básicos
func generateBasicCRUDMethods(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)

	// Save method
	fmt.Fprintf(content, "func (p *%s) Save(%s *domain.%s) error {\n", repoName, entityLower, entity)
	fmt.Fprintf(content, "\tresult := p.db.Create(%s)\n", entityLower)
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// FindByID method
	fmt.Fprintf(content, "func (p *%s) FindByID(id int) (*domain.%s, error) {\n", repoName, entity)
	fmt.Fprintf(content, "\t%s := &domain.%s{}\n", entityLower, entity)
	fmt.Fprintf(content, "\tresult := p.db.First(%s, id)\n", entityLower)
	content.WriteString("\tif result.Error != nil {\n")
	content.WriteString("\t\treturn nil, result.Error\n")
	content.WriteString("\t}\n")
	fmt.Fprintf(content, "\treturn %s, nil\n", entityLower)
	content.WriteString("}\n\n")

	// Update method
	fmt.Fprintf(content, "func (p *%s) Update(%s *domain.%s) error {\n", repoName, entityLower, entity)
	fmt.Fprintf(content, "\tresult := p.db.Save(%s)\n", entityLower)
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// Delete method
	fmt.Fprintf(content, "func (p *%s) Delete(id int) error {\n", repoName)
	fmt.Fprintf(content, "\tresult := p.db.Delete(&domain.%s{}, id)\n", entity)
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// FindAll method
	fmt.Fprintf(content, "func (p *%s) FindAll() ([]domain.%s, error) {\n", repoName, entity)
	fmt.Fprintf(content, "\tvar %ss []domain.%s\n", entityLower, entity)
	fmt.Fprintf(content, "\tresult := p.db.Find(&%ss)\n", entityLower)
	content.WriteString("\tif result.Error != nil {\n")
	content.WriteString("\t\treturn nil, result.Error\n")
	content.WriteString("\t}\n")
	fmt.Fprintf(content, "\treturn %ss, nil\n", entityLower)
	content.WriteString("}\n\n")
}

// generateTransactionMethods genera métodos que soportan transacciones
func generateTransactionMethods(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)

	// SaveWithTx
	fmt.Fprintf(content, "func (p *%s) SaveWithTx(tx interface{}, %s *domain.%s) error {\n",
		repoName, entityLower, entity)
	content.WriteString("\tgormTx := tx.(*gorm.DB)\n")
	fmt.Fprintf(content, "\tresult := gormTx.Create(%s)\n", entityLower)
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// UpdateWithTx
	fmt.Fprintf(content, "func (p *%s) UpdateWithTx(tx interface{}, %s *domain.%s) error {\n",
		repoName, entityLower, entity)
	content.WriteString("\tgormTx := tx.(*gorm.DB)\n")
	fmt.Fprintf(content, "\tresult := gormTx.Save(%s)\n", entityLower)
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// DeleteWithTx
	fmt.Fprintf(content, "func (p *%s) DeleteWithTx(tx interface{}, id int) error {\n", repoName)
	content.WriteString("\tgormTx := tx.(*gorm.DB)\n")
	fmt.Fprintf(content, "\tresult := gormTx.Delete(&domain.%s{}, id)\n", entity)
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")
}

// generateMySQLRepositoryWithFields genera repository MySQL con métodos dinámicos
func generateMySQLRepositoryWithFields(dir, entity string, fields []Field, cache, transactions bool) {
	// Para MySQL usamos la misma lógica que PostgreSQL ya que ambos usan GORM
	generatePostgresRepositoryWithFields(dir, entity, fields, cache, transactions)
}

// generateMongoRepositoryWithFields genera repository MongoDB con métodos dinámicos
func generateMongoRepositoryWithFields(dir, entity string, fields []Field, cache, transactions bool) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, "mongo_"+entityLower+"_repository.go")

	// Get the module name from go.mod
	moduleName := getModuleName()

	var content strings.Builder
	content.WriteString("package repository\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"context\"\n")
	content.WriteString("\t\"time\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
	if cache {
		content.WriteString("\t// Imports para cache de MongoDB\n")
		content.WriteString("\t// \"github.com/go-redis/redis/v8\"\n")
	}
	if transactions {
		content.WriteString("\t// Soporte para transacciones de MongoDB\n")
		content.WriteString("\t// \"go.mongodb.org/mongo-driver/mongo/options\"\n")
		content.WriteString("\t// \"go.mongodb.org/mongo-driver/mongo/writeconcern\"\n")
	}
	content.WriteString("\n\t\"go.mongodb.org/mongo-driver/mongo\"\n")
	content.WriteString("\t\"go.mongodb.org/mongo-driver/bson\"\n")
	content.WriteString("\t\"go.mongodb.org/mongo-driver/bson/primitive\"\n")
	content.WriteString(")\n\n")

	// MongoDB repository structure
	repoName := fmt.Sprintf("mongo%sRepository", entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n", repoName))
	content.WriteString("\tcollection *mongo.Collection\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("func NewMongo%sRepository(db *mongo.Database) %sRepository {\n", entity, entity))
	content.WriteString(fmt.Sprintf("\treturn &%s{\n", repoName))
	content.WriteString(fmt.Sprintf("\t\tcollection: db.Collection(\"%ss\"),\n", entityLower))
	content.WriteString("\t}\n")
	content.WriteString("}\n\n")

	// Generate basic MongoDB methods (simplified for brevity)
	generateBasicMongoCRUDMethods(&content, entity, repoName)

	// Generate dynamic search methods for MongoDB
	searchMethods := generateSearchMethods(fields, entity)
	for _, method := range searchMethods {
		content.WriteString(generateMongoSearchMethodImplementation(method, repoName, entity))
	}

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating MongoDB repository with fields: %v\n", err)
	}
}

// generateBasicMongoCRUDMethods genera métodos CRUD básicos para MongoDB
func generateBasicMongoCRUDMethods(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)

	// Save method
	fmt.Fprintf(content, "func (m *%s) Save(%s *domain.%s) error {\n", repoName, entityLower, entity)
	content.WriteString("\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n")
	content.WriteString("\tdefer cancel()\n")
	fmt.Fprintf(content, "\t_, err := m.collection.InsertOne(ctx, %s)\n", entityLower)
	content.WriteString("\treturn err\n")
	content.WriteString("}\n\n")

	// FindByID method
	fmt.Fprintf(content, "func (m *%s) FindByID(id int) (*domain.%s, error) {\n", repoName, entity)
	content.WriteString("\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n")
	content.WriteString("\tdefer cancel()\n")
	fmt.Fprintf(content, "\t%s := &domain.%s{}\n", entityLower, entity)
	content.WriteString("\terr := m.collection.FindOne(ctx, bson.M{\"id\": id}).Decode(" + entityLower + ")\n")
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n")
	fmt.Fprintf(content, "\treturn %s, nil\n", entityLower)
	content.WriteString("}\n\n")

	// Add other basic methods (Update, Delete, FindAll) - simplified
	content.WriteString("// Otros métodos CRUD básicos para MongoDB...\n\n")
}

// generateMongoSearchMethodImplementation genera implementación de método de búsqueda para MongoDB
func generateMongoSearchMethodImplementation(method SearchMethod, repoName, entity string) string {
	paramName := strings.ToLower(method.FieldName)
	entityVar := strings.ToLower(entity)

	var implementation strings.Builder
	implementation.WriteString(fmt.Sprintf("func (m *%s) %s(%s %s) %s {\n",
		repoName, method.MethodName, paramName, method.FieldType, method.ReturnType))

	implementation.WriteString("\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n")
	implementation.WriteString("\tdefer cancel()\n")
	implementation.WriteString(fmt.Sprintf("\t%s := &domain.%s{}\n", entityVar, entity))
	implementation.WriteString(fmt.Sprintf("\terr := m.collection.FindOne(ctx, bson.M{\"%s\": %s}).Decode(%s)\n",
		strings.ToLower(method.FieldName), paramName, entityVar))
	implementation.WriteString("\tif err != nil {\n")
	implementation.WriteString("\t\treturn nil, err\n")
	implementation.WriteString("\t}\n")
	implementation.WriteString(fmt.Sprintf("\treturn %s, nil\n", entityVar))
	implementation.WriteString("}\n\n")

	return implementation.String()
}

func init() {
	repositoryCmd.Flags().StringP(DatabaseFlag, "d", "", DatabaseFlagUsage)
	repositoryCmd.Flags().BoolP(InterfaceOnlyFlag, "i", false, InterfaceOnlyFlagUsage)
	repositoryCmd.Flags().BoolP(ImplementationFlag, "", false, ImplementationFlagUsage)
	repositoryCmd.Flags().BoolP(CacheFlag, "c", false, CacheFlagUsage)
	repositoryCmd.Flags().BoolP(TransactionsFlag, "t", false, TransactionsFlagUsage)
	repositoryCmd.Flags().StringP("fields", "f", "", "Entity fields \"field:type,field2:type\"")
}
