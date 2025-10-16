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
	case DBPostgresJSON:
		generatePostgresJSONRepository(dir, entity, cache, transactions)
	case DBMySQL:
		generateMySQLRepository(dir, entity, cache, transactions)
	case DBMongoDB:
		generateMongoRepository(dir, entity, cache, transactions)
	case DBSQLite:
		generateSQLiteRepository(dir, entity, cache, transactions)
	case DBSQLServer:
		generateSQLServerRepository(dir, entity, cache, transactions)
	case DBElasticsearch:
		generateElasticsearchRepository(dir, entity, cache, transactions)
	case DBDynamoDB:
		generateDynamoDBRepository(dir, entity, cache, transactions)
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

// generatePostgresJSONRepository genera un repository para PostgreSQL con soporte JSONB
func generatePostgresJSONRepository(dir, entity string, cache, transactions bool) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, "postgres_json_"+entityLower+"_repository.go")
	moduleName := getModuleName()

	var content strings.Builder
	content.WriteString("package repository\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"encoding/json\"\n")
	content.WriteString("\t\"gorm.io/datatypes\"\n")
	content.WriteString("\t\"gorm.io/gorm\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
	content.WriteString(")\n\n")

	repoName := fmt.Sprintf("postgresJSON%sRepository", entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n\tdb *gorm.DB\n}\n\n", repoName))
	content.WriteString(fmt.Sprintf("func NewPostgresJSON%sRepository(db *gorm.DB) %sRepository {\n", entity, entity))
	content.WriteString(fmt.Sprintf("\treturn &%s{db: db}\n", repoName))
	content.WriteString("}\n\n")

	// Save method with JSONB support
	content.WriteString(fmt.Sprintf("func (p *%s) Save(%s *domain.%s) error {\n", repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\treturn p.db.Create(%s).Error\n", entityLower))
	content.WriteString("}\n\n")

	// FindByID method
	content.WriteString(fmt.Sprintf("func (p *%s) FindByID(id int) (*domain.%s, error) {\n", repoName, entity))
	content.WriteString(fmt.Sprintf("\tvar %s domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tif err := p.db.First(&%s, id).Error; err != nil {\n", entityLower))
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\treturn &%s, nil\n", entityLower))
	content.WriteString("}\n\n")

	// FindByJSONField - Query nested JSON fields
	content.WriteString(fmt.Sprintf("func (p *%s) FindByJSONField(jsonField, value string) ([]domain.%s, error) {\n", repoName, entity))
	content.WriteString(fmt.Sprintf("\tvar %ss []domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tif err := p.db.Where(\"data @> ?\", datatypes.JSONQuery(jsonField)).Find(&%ss).Error; err != nil {\n", entityLower))
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\treturn %ss, nil\n", entityLower))
	content.WriteString("}\n\n")

	// Update method
	content.WriteString(fmt.Sprintf("func (p *%s) Update(%s *domain.%s) error {\n", repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\treturn p.db.Save(%s).Error\n", entityLower))
	content.WriteString("}\n\n")

	// Delete method
	content.WriteString(fmt.Sprintf("func (p *%s) Delete(id int) error {\n", repoName))
	content.WriteString(fmt.Sprintf("\treturn p.db.Delete(&domain.%s{}, id).Error\n", entity))
	content.WriteString("}\n\n")

	// FindAll method
	content.WriteString(fmt.Sprintf("func (p *%s) FindAll() ([]domain.%s, error) {\n", repoName, entity))
	content.WriteString(fmt.Sprintf("\tvar %ss []domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tif err := p.db.Find(&%ss).Error; err != nil {\n", entityLower))
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\treturn %ss, nil\n", entityLower))
	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating PostgreSQL JSON repository file: %v\n", err)
	}
}

// generateSQLServerRepository genera un repository para SQL Server con GORM + mssql
func generateSQLServerRepository(dir, entity string, cache, transactions bool) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, "sqlserver_"+entityLower+"_repository.go")
	moduleName := getModuleName()

	var content strings.Builder
	content.WriteString("package repository\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"fmt\"\n")
	content.WriteString("\t\"gorm.io/gorm\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
	content.WriteString(")\n\n")

	repoName := fmt.Sprintf("sqlserver%sRepository", entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n\tdb *gorm.DB\n}\n\n", repoName))
	content.WriteString(fmt.Sprintf("func NewSQLServer%sRepository(db *gorm.DB) %sRepository {\n", entity, entity))
	content.WriteString(fmt.Sprintf("\treturn &%s{db: db}\n", repoName))
	content.WriteString("}\n\n")

	// Save method
	content.WriteString(fmt.Sprintf("func (s *%s) Save(%s *domain.%s) error {\n", repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\tif err := s.db.Create(%s).Error; err != nil {\n", entityLower))
	content.WriteString("\t\treturn fmt.Errorf(\"failed to save %s: %%w\", err)\n")
	content.WriteString("\t}\n")
	content.WriteString("\treturn nil\n")
	content.WriteString("}\n\n")

	// FindByID method
	content.WriteString(fmt.Sprintf("func (s *%s) FindByID(id int) (*domain.%s, error) {\n", repoName, entity))
	content.WriteString(fmt.Sprintf("\tvar %s domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tif err := s.db.WithContext(s.db.Statement.Context).First(&%s, id).Error; err != nil {\n", entityLower))
	content.WriteString("\t\tif err == gorm.ErrRecordNotFound {\n")
	content.WriteString("\t\t\treturn nil, fmt.Errorf(\"%s not found\")\n")
	content.WriteString("\t\t}\n")
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\treturn &%s, nil\n", entityLower))
	content.WriteString("}\n\n")

	// Update method
	content.WriteString(fmt.Sprintf("func (s *%s) Update(%s *domain.%s) error {\n", repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\tif err := s.db.Save(%s).Error; err != nil {\n", entityLower))
	content.WriteString("\t\treturn fmt.Errorf(\"failed to update %s: %%w\", err)\n")
	content.WriteString("\t}\n")
	content.WriteString("\treturn nil\n")
	content.WriteString("}\n\n")

	// Delete method
	content.WriteString(fmt.Sprintf("func (s *%s) Delete(id int) error {\n", repoName))
	content.WriteString(fmt.Sprintf("\tif err := s.db.Delete(&domain.%s{}, id).Error; err != nil {\n", entity))
	content.WriteString("\t\treturn fmt.Errorf(\"failed to delete %s: %%w\", err)\n")
	content.WriteString("\t}\n")
	content.WriteString("\treturn nil\n")
	content.WriteString("}\n\n")

	// FindAll method
	content.WriteString(fmt.Sprintf("func (s *%s) FindAll() ([]domain.%s, error) {\n", repoName, entity))
	content.WriteString(fmt.Sprintf("\tvar %ss []domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tif err := s.db.Find(&%ss).Error; err != nil {\n", entityLower))
	content.WriteString("\t\treturn nil, fmt.Errorf(\"failed to fetch %ss: %%w\", err)\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\treturn %ss, nil\n", entityLower))
	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating SQL Server repository file: %v\n", err)
	}
}

// generateElasticsearchRepository genera un repository para Elasticsearch con búsqueda full-text
func generateElasticsearchRepository(dir, entity string, cache, transactions bool) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, "elasticsearch_"+entityLower+"_repository.go")
	moduleName := getModuleName()

	var content strings.Builder
	content.WriteString("package repository\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"bytes\"\n")
	content.WriteString("\t\"context\"\n")
	content.WriteString("\t\"encoding/json\"\n")
	content.WriteString("\t\"fmt\"\n")
	content.WriteString("\t\"github.com/elastic/go-elasticsearch/v8\"\n")
	content.WriteString("\t\"strconv\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
	content.WriteString(")\n\n")

	repoName := fmt.Sprintf("elasticsearch%sRepository", entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n\tclient *elasticsearch.Client\n\tindex  string\n}\n\n", repoName))
	content.WriteString(fmt.Sprintf("func NewElasticsearch%sRepository(client *elasticsearch.Client) %sRepository {\n", entity, entity))
	content.WriteString(fmt.Sprintf("\treturn &%s{\n\t\tclient: client,\n\t\tindex:  \"%s\",\n\t}\n", repoName, strings.ToLower(entity)))
	content.WriteString("}\n\n")

	// Save method
	content.WriteString(fmt.Sprintf("func (e *%s) Save(%s *domain.%s) error {\n", repoName, entityLower, entity))
	content.WriteString("\tdata, err := json.Marshal(" + entityLower + ")\n")
	content.WriteString("\tif err != nil {\n\t\treturn err\n\t}\n")
	content.WriteString("\treq := esapi.IndexRequest{\n")
	content.WriteString("\t\tIndex: e.index,\n")
	content.WriteString("\t\tBody:  bytes.NewReader(data),\n")
	content.WriteString("\t}\n")
	content.WriteString("\tres, err := req.Do(context.Background(), e.client)\n")
	content.WriteString("\tif err != nil {\n\t\treturn err\n\t}\n")
	content.WriteString("\tdefer res.Body.Close()\n")
	content.WriteString("\treturn nil\n")
	content.WriteString("}\n\n")

	// FindByID method
	content.WriteString(fmt.Sprintf("func (e *%s) FindByID(id int) (*domain.%s, error) {\n", repoName, entity))
	content.WriteString("\treq := esapi.GetRequest{\n")
	content.WriteString("\t\tIndex:      e.index,\n")
	content.WriteString("\t\tDocumentID: strconv.Itoa(id),\n")
	content.WriteString("\t}\n")
	content.WriteString("\tres, err := req.Do(context.Background(), e.client)\n")
	content.WriteString("\tif err != nil {\n\t\treturn nil, err\n\t}\n")
	content.WriteString("\tdefer res.Body.Close()\n")
	content.WriteString(fmt.Sprintf("\tvar doc domain.%s\n", entity))
	content.WriteString("\tif err := json.NewDecoder(res.Body).Decode(&doc); err != nil {\n")
	content.WriteString("\t\treturn nil, err\n\t}\n")
	content.WriteString("\treturn &doc, nil\n")
	content.WriteString("}\n\n")

	// FullTextSearch method
	content.WriteString(fmt.Sprintf("func (e *%s) FullTextSearch(query string) ([]domain.%s, error) {\n", repoName, entity))
	content.WriteString("\tsearchBody := map[string]interface{}{\n")
	content.WriteString("\t\t\"query\": map[string]interface{}{\n")
	content.WriteString("\t\t\t\"multi_match\": map[string]interface{}{\n")
	content.WriteString("\t\t\t\t\"query\":  query,\n")
	content.WriteString("\t\t\t\t\"fields\": []string{\"*\"},\n")
	content.WriteString("\t\t\t},\n")
	content.WriteString("\t\t},\n")
	content.WriteString("\t}\n")
	content.WriteString("\tvar buf bytes.Buffer\n")
	content.WriteString("\tif err := json.NewEncoder(&buf).Encode(searchBody); err != nil {\n")
	content.WriteString("\t\treturn nil, err\n\t}\n")
	content.WriteString("\treq := esapi.SearchRequest{\n")
	content.WriteString("\t\tIndex: []string{e.index},\n")
	content.WriteString("\t\tBody:  &buf,\n")
	content.WriteString("\t}\n")
	content.WriteString("\tres, err := req.Do(context.Background(), e.client)\n")
	content.WriteString("\tif err != nil {\n\t\treturn nil, err\n\t}\n")
	content.WriteString("\tdefer res.Body.Close()\n")
	content.WriteString(fmt.Sprintf("\tvar results []domain.%s\n", entity))
	content.WriteString("\tvar sr map[string]interface{}\n")
	content.WriteString("\tif err := json.NewDecoder(res.Body).Decode(&sr); err != nil {\n")
	content.WriteString("\t\treturn nil, err\n\t}\n")
	content.WriteString("\treturn results, nil\n")
	content.WriteString("}\n\n")

	// FindAll method
	content.WriteString(fmt.Sprintf("func (e *%s) FindAll() ([]domain.%s, error) {\n", repoName, entity))
	content.WriteString("\tsearchBody := map[string]interface{}{\n")
	content.WriteString("\t\t\"query\": map[string]interface{}{\n")
	content.WriteString("\t\t\t\"match_all\": map[string]interface{}{},\n")
	content.WriteString("\t\t},\n")
	content.WriteString("\t}\n")
	content.WriteString("\tvar buf bytes.Buffer\n")
	content.WriteString("\tif err := json.NewEncoder(&buf).Encode(searchBody); err != nil {\n")
	content.WriteString("\t\treturn nil, err\n\t}\n")
	content.WriteString("\treq := esapi.SearchRequest{\n")
	content.WriteString("\t\tIndex: []string{e.index},\n")
	content.WriteString("\t\tBody:  &buf,\n")
	content.WriteString("\t}\n")
	content.WriteString("\tres, err := req.Do(context.Background(), e.client)\n")
	content.WriteString("\tif err != nil {\n\t\treturn nil, err\n\t}\n")
	content.WriteString("\tdefer res.Body.Close()\n")
	content.WriteString(fmt.Sprintf("\tvar results []domain.%s\n", entity))
	content.WriteString("\treturn results, nil\n")
	content.WriteString("}\n\n")

	// Delete method
	content.WriteString(fmt.Sprintf("func (e *%s) Delete(id int) error {\n", repoName))
	content.WriteString("\treq := esapi.DeleteRequest{\n")
	content.WriteString("\t\tIndex:      e.index,\n")
	content.WriteString("\t\tDocumentID: strconv.Itoa(id),\n")
	content.WriteString("\t}\n")
	content.WriteString("\tres, err := req.Do(context.Background(), e.client)\n")
	content.WriteString("\tif err != nil {\n\t\treturn err\n\t}\n")
	content.WriteString("\tdefer res.Body.Close()\n")
	content.WriteString("\treturn nil\n")
	content.WriteString("}\n\n")

	// Update method (stub)
	content.WriteString(fmt.Sprintf("func (e *%s) Update(%s *domain.%s) error {\n", repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\treturn e.Save(%s)\n", entityLower))
	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating Elasticsearch repository file: %v\n", err)
	}
}

// generateDynamoDBRepository genera un repository para DynamoDB con AWS SDK v2
func generateDynamoDBRepository(dir, entity string, cache, transactions bool) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, "dynamodb_"+entityLower+"_repository.go")
	moduleName := getModuleName()

	var content strings.Builder
	content.WriteString("package repository\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"context\"\n")
	content.WriteString("\t\"encoding/json\"\n")
	content.WriteString("\t\"fmt\"\n")
	content.WriteString("\t\"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue\"\n")
	content.WriteString("\t\"github.com/aws/aws-sdk-go-v2/service/dynamodb\"\n")
	content.WriteString("\t\"github.com/aws/aws-sdk-go-v2/service/dynamodb/types\"\n")
	content.WriteString("\t\"strconv\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
	content.WriteString(")\n\n")

	repoName := fmt.Sprintf("dynamodb%sRepository", entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n\tclient    *dynamodb.Client\n\ttableName string\n}\n\n", repoName))
	content.WriteString(fmt.Sprintf("func NewDynamoDB%sRepository(client *dynamodb.Client) %sRepository {\n", entity, entity))
	content.WriteString(fmt.Sprintf("\treturn &%s{\n\t\tclient:    client,\n\t\ttableName: \"%s\",\n\t}\n", repoName, strings.ToLower(entity)))
	content.WriteString("}\n\n")

	// Save method
	content.WriteString(fmt.Sprintf("func (d *%s) Save(%s *domain.%s) error {\n", repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\tav, err := attributevalue.MarshalMap(%s)\n", entityLower))
	content.WriteString("\tif err != nil {\n\t\treturn fmt.Errorf(\"failed to marshal: %%w\", err)\n\t}\n")
	content.WriteString("\t_, err = d.client.PutItem(context.Background(), &dynamodb.PutItemInput{\n")
	content.WriteString("\t\tTableName: &d.tableName,\n")
	content.WriteString("\t\tItem:      av,\n")
	content.WriteString("\t})\n")
	content.WriteString("\treturn err\n")
	content.WriteString("}\n\n")

	// FindByID method
	content.WriteString(fmt.Sprintf("func (d *%s) FindByID(id int) (*domain.%s, error) {\n", repoName, entity))
	content.WriteString("\tresult, err := d.client.GetItem(context.Background(), &dynamodb.GetItemInput{\n")
	content.WriteString("\t\tTableName: &d.tableName,\n")
	content.WriteString("\t\tKey: map[string]types.AttributeValue{\n")
	content.WriteString("\t\t\t\"id\": &types.AttributeValueMemberN{Value: strconv.Itoa(id)},\n")
	content.WriteString("\t\t},\n")
	content.WriteString("\t})\n")
	content.WriteString("\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"failed to get item: %%w\", err)\n\t}\n")
	content.WriteString(fmt.Sprintf("\tvar %s domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\terr = attributevalue.UnmarshalMap(result.Item, &%s)\n", entityLower))
	content.WriteString("\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"failed to unmarshal: %%w\", err)\n\t}\n")
	content.WriteString(fmt.Sprintf("\treturn &%s, nil\n", entityLower))
	content.WriteString("}\n\n")

	// Update method
	content.WriteString(fmt.Sprintf("func (d *%s) Update(%s *domain.%s) error {\n", repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\treturn d.Save(%s)\n", entityLower))
	content.WriteString("}\n\n")

	// Delete method
	content.WriteString(fmt.Sprintf("func (d *%s) Delete(id int) error {\n", repoName))
	content.WriteString("\t_, err := d.client.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{\n")
	content.WriteString("\t\tTableName: &d.tableName,\n")
	content.WriteString("\t\tKey: map[string]types.AttributeValue{\n")
	content.WriteString("\t\t\t\"id\": &types.AttributeValueMemberN{Value: strconv.Itoa(id)},\n")
	content.WriteString("\t\t},\n")
	content.WriteString("\t})\n")
	content.WriteString("\treturn err\n")
	content.WriteString("}\n\n")

	// FindAll method
	content.WriteString(fmt.Sprintf("func (d *%s) FindAll() ([]domain.%s, error) {\n", repoName, entity))
	content.WriteString("\tresult, err := d.client.Scan(context.Background(), &dynamodb.ScanInput{\n")
	content.WriteString("\t\tTableName: &d.tableName,\n")
	content.WriteString("\t})\n")
	content.WriteString("\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"failed to scan: %%w\", err)\n\t}\n")
	content.WriteString(fmt.Sprintf("\tvar %ss []domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\terr = attributevalue.UnmarshalListOfMaps(result.Items, &%ss)\n", entityLower))
	content.WriteString("\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"failed to unmarshal: %%w\", err)\n\t}\n")
	content.WriteString(fmt.Sprintf("\treturn %ss, nil\n", entityLower))
	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating DynamoDB repository file: %v\n", err)
	}
}

// generateSQLiteRepository genera un repository para SQLite con database/sql
func generateSQLiteRepository(dir, entity string, cache, transactions bool) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, "sqlite_"+entityLower+"_repository.go")
	moduleName := getModuleName()

	var content strings.Builder
	content.WriteString("package repository\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"database/sql\"\n")
	content.WriteString("\t\"encoding/json\"\n")
	content.WriteString("\t\"fmt\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
	content.WriteString(")\n\n")

	repoName := fmt.Sprintf("sqlite%sRepository", entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n\tdb *sql.DB\n}\n\n", repoName))
	content.WriteString(fmt.Sprintf("func NewSQLite%sRepository(db *sql.DB) %sRepository {\n", entity, entity))
	content.WriteString(fmt.Sprintf("\treturn &%s{db: db}\n", repoName))
	content.WriteString("}\n\n")

	// Save method
	content.WriteString(fmt.Sprintf("func (s *%s) Save(%s *domain.%s) error {\n", repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\tvar data []byte\n"))
	content.WriteString(fmt.Sprintf("\tdata, err := json.Marshal(%s)\n", entityLower))
	content.WriteString(fmt.Sprintf("\tif err != nil {\n\t\treturn fmt.Errorf(\"failed to marshal: %%w\", err)\n\t}\n"))
	content.WriteString(fmt.Sprintf("\tquery := \"INSERT INTO %ss (data) VALUES (?)\"\n", entityLower))
	content.WriteString(fmt.Sprintf("\tif _, err := s.db.Exec(query, data); err != nil {\n"))
	content.WriteString(fmt.Sprintf("\t\treturn fmt.Errorf(\"failed to insert: %%w\", err)\n\t}\n"))
	content.WriteString(fmt.Sprintf("\treturn nil\n"))
	content.WriteString("}\n\n")

	// FindByID method
	content.WriteString(fmt.Sprintf("func (s *%s) FindByID(id int) (*domain.%s, error) {\n", repoName, entity))
	content.WriteString(fmt.Sprintf("\tvar data []byte\n"))
	content.WriteString(fmt.Sprintf("\tquery := \"SELECT data FROM %ss WHERE id = ? LIMIT 1\"\n", entityLower))
	content.WriteString(fmt.Sprintf("\tif err := s.db.QueryRow(query, id).Scan(&data); err != nil {\n"))
	content.WriteString(fmt.Sprintf("\t\tif err == sql.ErrNoRows {\n\t\t\treturn nil, fmt.Errorf(\"%s not found\")\n\t\t}\n", entity))
	content.WriteString(fmt.Sprintf("\t\treturn nil, fmt.Errorf(\"failed to query: %%w\", err)\n\t}\n"))
	content.WriteString(fmt.Sprintf("\tvar %s domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tif err := json.Unmarshal(data, &%s); err != nil {\n", entityLower))
	content.WriteString(fmt.Sprintf("\t\treturn nil, fmt.Errorf(\"failed to unmarshal: %%w\", err)\n\t}\n"))
	content.WriteString(fmt.Sprintf("\treturn &%s, nil\n", entityLower))
	content.WriteString("}\n\n")

	// Update method
	content.WriteString(fmt.Sprintf("func (s *%s) Update(%s *domain.%s) error {\n", repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\tdata, err := json.Marshal(%s)\n", entityLower))
	content.WriteString(fmt.Sprintf("\tif err != nil {\n\t\treturn fmt.Errorf(\"failed to marshal: %%w\", err)\n\t}\n"))
	content.WriteString(fmt.Sprintf("\tquery := \"UPDATE %ss SET data = ? WHERE id = ?\"\n", entityLower))
	content.WriteString(fmt.Sprintf("\tif _, err := s.db.Exec(query, data, %s.ID); err != nil {\n", entityLower))
	content.WriteString(fmt.Sprintf("\t\treturn fmt.Errorf(\"failed to update: %%w\", err)\n\t}\n"))
	content.WriteString(fmt.Sprintf("\treturn nil\n"))
	content.WriteString("}\n\n")

	// Delete method
	content.WriteString(fmt.Sprintf("func (s *%s) Delete(id int) error {\n", repoName))
	content.WriteString(fmt.Sprintf("\tquery := \"DELETE FROM %ss WHERE id = ?\"\n", entityLower))
	content.WriteString(fmt.Sprintf("\tif _, err := s.db.Exec(query, id); err != nil {\n"))
	content.WriteString(fmt.Sprintf("\t\treturn fmt.Errorf(\"failed to delete: %%w\", err)\n\t}\n"))
	content.WriteString(fmt.Sprintf("\treturn nil\n"))
	content.WriteString("}\n\n")

	// FindAll method
	content.WriteString(fmt.Sprintf("func (s *%s) FindAll() ([]domain.%s, error) {\n", repoName, entity))
	content.WriteString(fmt.Sprintf("\tquery := \"SELECT data FROM %ss\"\n", entityLower))
	content.WriteString(fmt.Sprintf("\trows, err := s.db.Query(query)\n"))
	content.WriteString(fmt.Sprintf("\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"failed to query: %%w\", err)\n\t}\n"))
	content.WriteString(fmt.Sprintf("\tdefer rows.Close()\n"))
	content.WriteString(fmt.Sprintf("\tvar %ss []domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tfor rows.Next() {\n"))
	content.WriteString(fmt.Sprintf("\t\tvar data []byte\n"))
	content.WriteString(fmt.Sprintf("\t\tif err := rows.Scan(&data); err != nil {\n\t\t\treturn nil, fmt.Errorf(\"failed to scan: %%w\", err)\n\t\t}\n"))
	content.WriteString(fmt.Sprintf("\t\tvar %s domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\t\tif err := json.Unmarshal(data, &%s); err != nil {\n\t\t\treturn nil, fmt.Errorf(\"failed to unmarshal: %%w\", err)\n\t\t}\n", entityLower))
	content.WriteString(fmt.Sprintf("\t\t%ss = append(%ss, %s)\n", entityLower, entityLower, entityLower))
	content.WriteString(fmt.Sprintf("\t}\n"))
	content.WriteString(fmt.Sprintf("\tif err := rows.Err(); err != nil {\n\t\treturn nil, fmt.Errorf(\"rows error: %%w\", err)\n\t}\n"))
	content.WriteString(fmt.Sprintf("\treturn %ss, nil\n", entityLower))
	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating SQLite repository file: %v\n", err)
	}
}

func init() {
	repositoryCmd.Flags().StringP(DatabaseFlag, "d", "", DatabaseFlagUsage)
	repositoryCmd.Flags().BoolP(InterfaceOnlyFlag, "i", false, InterfaceOnlyFlagUsage)
	repositoryCmd.Flags().BoolP(ImplementationFlag, "", false, ImplementationFlagUsage)
	repositoryCmd.Flags().BoolP(CacheFlag, "c", false, CacheFlagUsage)
	repositoryCmd.Flags().BoolP(TransactionsFlag, "t", false, TransactionsFlagUsage)
	repositoryCmd.Flags().StringP("fields", "f", "", "Entity fields \"field:type,field2:type\"")
}
