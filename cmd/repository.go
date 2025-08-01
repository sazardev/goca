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
	Short: "Generar repositorios con interfaces",
	Long: `Crea repositorios que implementan el patrón Repository con interfaces 
bien definidas e implementaciones específicas por base de datos.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entity := args[0]

		database, _ := cmd.Flags().GetString(DatabaseFlag)
		interfaceOnly, _ := cmd.Flags().GetBool(InterfaceOnlyFlag)
		implementation, _ := cmd.Flags().GetBool(ImplementationFlag)
		cache, _ := cmd.Flags().GetBool(CacheFlag)
		transactions, _ := cmd.Flags().GetBool(TransactionsFlag)
		fields, _ := cmd.Flags().GetString("fields")

		// Validar con el nuevo validador robusto
		validator := NewFieldValidator()

		if err := validator.ValidateEntityName(entity); err != nil {
			fmt.Printf("❌ Error en nombre de entidad: %v\n", err)
			return
		}

		if database != "" {
			if err := validator.ValidateDatabase(database); err != nil {
				fmt.Printf("❌ Error en base de datos: %v\n", err)
				return
			}
		}

		fmt.Printf("🚀 Generando repositorio para entidad '%s'\n", entity)

		if database != "" && !interfaceOnly {
			fmt.Printf("🗄️  Base de datos: %s\n", database)
		}
		if interfaceOnly {
			fmt.Println("✓ Solo interfaces")
		}
		if implementation {
			fmt.Println("✓ Solo implementación")
		}
		if cache {
			fmt.Println("✓ Incluyendo caché")
		}
		if transactions {
			fmt.Println("✓ Incluyendo transacciones")
		}
		if fields != "" {
			fmt.Printf("✓ Campos personalizados: %s\n", fields)
		}

		generateRepository(entity, database, interfaceOnly, implementation, cache, transactions, fields)
		fmt.Printf("\n✅ Repositorio para '%s' generado exitosamente!\n", entity)
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

	// Generate interface if not implementation-only
	if !implementation || interfaceOnly {
		if fields != "" {
			generateRepositoryInterfaceWithFields(repoDir, entity, parsedFields, transactions)
		} else {
			generateRepositoryInterface(repoDir, entity, transactions)
		}
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
		fmt.Printf("❌ Error escribiendo archivo %s: %v\n", filename, err)
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
		fmt.Printf("Base de datos no soportada: %s\n", database)
		os.Exit(1)
	}
}

func generatePostgresRepository(dir, entity string, cache, transactions bool) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, "postgres_"+entityLower+"_repo.go")

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

	writeGoFile(filename, content.String())
}

func generatePostgresSaveMethod(content *strings.Builder, entity, repoName string, cache bool) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) Save(%s *domain.%s) error {\n",
		repoVar, repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := %s.db.Create(%s)\n", repoVar, entityLower))

	if cache {
		content.WriteString("\tif result.Error == nil {\n")
		content.WriteString(fmt.Sprintf("\t\t%s.invalidateCache(%s.ID)\n", repoVar, entityLower))
		content.WriteString("\t}\n")
	}

	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")
}

func generatePostgresFindByIDMethod(content *strings.Builder, entity, repoName string, cache bool) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) FindByID(id int) (*domain.%s, error) {\n",
		repoVar, repoName, entity))

	if cache {
		content.WriteString("\t// Try cache first\n")
		content.WriteString(fmt.Sprintf("\tif %s := %s.getFromCache(id); %s != nil {\n",
			entityLower, repoVar, entityLower))
		content.WriteString(fmt.Sprintf("\t\treturn %s, nil\n", entityLower))
		content.WriteString("\t}\n\n")
	}

	content.WriteString(fmt.Sprintf("\t%s := &domain.%s{}\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := %s.db.First(%s, id)\n", repoVar, entityLower))
	content.WriteString("\tif result.Error != nil {\n")
	content.WriteString("\t\treturn nil, result.Error\n")
	content.WriteString("\t}\n\n")

	if cache {
		content.WriteString(fmt.Sprintf("\t%s.setCache(%s)\n", repoVar, entityLower))
	}

	content.WriteString(fmt.Sprintf("\treturn %s, nil\n", entityLower))
	content.WriteString("}\n\n")
}

func generatePostgresFindByEmailMethod(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) FindByEmail(email string) (*domain.%s, error) {\n",
		repoVar, repoName, entity))
	content.WriteString(fmt.Sprintf("\t%s := &domain.%s{}\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := %s.db.Where(\"email = ?\", email).First(%s)\n", repoVar, entityLower))
	content.WriteString("\tif result.Error != nil {\n")
	content.WriteString("\t\treturn nil, result.Error\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\treturn %s, nil\n", entityLower))
	content.WriteString("}\n\n")
}

func generatePostgresUpdateMethod(content *strings.Builder, entity, repoName string, cache bool) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) Update(%s *domain.%s) error {\n",
		repoVar, repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := %s.db.Save(%s)\n", repoVar, entityLower))

	if cache {
		content.WriteString("\tif result.Error == nil {\n")
		content.WriteString(fmt.Sprintf("\t\t%s.invalidateCache(%s.ID)\n", repoVar, entityLower))
		content.WriteString("\t}\n")
	}

	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")
}

func generatePostgresDeleteMethod(content *strings.Builder, entity, repoName string, cache bool) {
	repoVar := strings.ToLower(string(repoName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) Delete(id int) error {\n",
		repoVar, repoName))
	content.WriteString(fmt.Sprintf("\tresult := %s.db.Delete(&domain.%s{}, id)\n", repoVar, entity))

	if cache {
		content.WriteString("\tif result.Error == nil {\n")
		content.WriteString(fmt.Sprintf("\t\t%s.invalidateCache(id)\n", repoVar))
		content.WriteString("\t}\n")
	}

	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")
}

func generatePostgresFindAllMethod(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) FindAll() ([]domain.%s, error) {\n",
		repoVar, repoName, entity))
	content.WriteString(fmt.Sprintf("\tvar %ss []domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := %s.db.Find(&%ss)\n", repoVar, entityLower))
	content.WriteString("\tif result.Error != nil {\n")
	content.WriteString("\t\treturn nil, result.Error\n")
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\treturn %ss, nil\n", entityLower))
	content.WriteString("}\n\n")
}

func generatePostgresTransactionMethods(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	// SaveWithTx
	content.WriteString(fmt.Sprintf("func (%s *%s) SaveWithTx(tx interface{}, %s *domain.%s) error {\n",
		repoVar, repoName, entityLower, entity))
	content.WriteString("\tgormTx := tx.(*gorm.DB)\n")
	content.WriteString(fmt.Sprintf("\tresult := gormTx.Create(%s)\n", entityLower))
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// UpdateWithTx
	content.WriteString(fmt.Sprintf("func (%s *%s) UpdateWithTx(tx interface{}, %s *domain.%s) error {\n",
		repoVar, repoName, entityLower, entity))
	content.WriteString("\tgormTx := tx.(*gorm.DB)\n")
	content.WriteString(fmt.Sprintf("\tresult := gormTx.Save(%s)\n", entityLower))
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// DeleteWithTx
	content.WriteString(fmt.Sprintf("func (%s *%s) DeleteWithTx(tx interface{}, id int) error {\n",
		repoVar, repoName))
	content.WriteString("\tgormTx := tx.(*gorm.DB)\n")
	content.WriteString(fmt.Sprintf("\tresult := gormTx.Delete(&domain.%s{}, id)\n", entity))
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
	content.WriteString("\t\"database/sql\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
	if cache {
		content.WriteString("\t// Imports para cache (Redis, etc.)\n")
		content.WriteString("\t// \"github.com/go-redis/redis/v8\"\n")
	}
	if transactions {
		content.WriteString("\t// Soporte para transacciones SQL\n")
		content.WriteString("\t// \"database/sql/driver\"\n")
	}
	content.WriteString("\n\t_ \"github.com/go-sql-driver/mysql\"\n")
	content.WriteString(")\n\n")

	// MySQL repository structure
	repoName := fmt.Sprintf("mysql%sRepository", entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n", repoName))
	content.WriteString("\tdb *sql.DB\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("func NewMySQL%sRepository(db *sql.DB) %sRepository {\n", entity, entity))
	content.WriteString(fmt.Sprintf("\treturn &%s{db: db}\n", repoName))
	content.WriteString("}\n\n")

	// Basic Save method for MySQL
	content.WriteString(fmt.Sprintf("func (r *%s) Save(%s *domain.%s) error {\n",
		repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\t// Query personalizada para tu entidad %s\n", entity))
	content.WriteString(fmt.Sprintf("\tquery := `INSERT INTO %ss (nombre, email, edad, created_at) VALUES (?, ?, ?, NOW())`\n", entityLower))
	content.WriteString(fmt.Sprintf("\tresult, err := r.db.Exec(query, %s.Nombre, %s.Email, %s.Edad)\n", entityLower, entityLower, entityLower))
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn err\n")
	content.WriteString("\t}\n")
	content.WriteString("\tid, err := result.LastInsertId()\n")
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn err\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\t%s.ID = uint(id)\n", entityLower))
	content.WriteString("\treturn nil\n")
	content.WriteString("}\n\n")

	writeGoFile(filename, content.String())
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

	writeGoFile(filename, content.String())
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
		fmt.Printf("❌ Error escribiendo archivo %s: %v\n", filename, err)
	}
}

// generateRepositoryImplementationWithFields genera implementaciones de repository con métodos dinámicos
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

	writeGoFile(filename, content.String())
}

// generateBasicCRUDMethods genera los métodos CRUD básicos
func generateBasicCRUDMethods(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)

	// Save method
	content.WriteString(fmt.Sprintf("func (p *%s) Save(%s *domain.%s) error {\n", repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := p.db.Create(%s)\n", entityLower))
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// FindByID method
	content.WriteString(fmt.Sprintf("func (p *%s) FindByID(id int) (*domain.%s, error) {\n", repoName, entity))
	content.WriteString(fmt.Sprintf("\t%s := &domain.%s{}\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := p.db.First(%s, id)\n", entityLower))
	content.WriteString("\tif result.Error != nil {\n")
	content.WriteString("\t\treturn nil, result.Error\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\treturn %s, nil\n", entityLower))
	content.WriteString("}\n\n")

	// Update method
	content.WriteString(fmt.Sprintf("func (p *%s) Update(%s *domain.%s) error {\n", repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := p.db.Save(%s)\n", entityLower))
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// Delete method
	content.WriteString(fmt.Sprintf("func (p *%s) Delete(id int) error {\n", repoName))
	content.WriteString(fmt.Sprintf("\tresult := p.db.Delete(&domain.%s{}, id)\n", entity))
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// FindAll method
	content.WriteString(fmt.Sprintf("func (p *%s) FindAll() ([]domain.%s, error) {\n", repoName, entity))
	content.WriteString(fmt.Sprintf("\tvar %ss []domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tresult := p.db.Find(&%ss)\n", entityLower))
	content.WriteString("\tif result.Error != nil {\n")
	content.WriteString(fmt.Sprintf("\t\treturn nil, result.Error\n"))
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\treturn %ss, nil\n", entityLower))
	content.WriteString("}\n\n")
}

// generateTransactionMethods genera métodos que soportan transacciones
func generateTransactionMethods(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)

	// SaveWithTx
	content.WriteString(fmt.Sprintf("func (p *%s) SaveWithTx(tx interface{}, %s *domain.%s) error {\n",
		repoName, entityLower, entity))
	content.WriteString("\tgormTx := tx.(*gorm.DB)\n")
	content.WriteString(fmt.Sprintf("\tresult := gormTx.Create(%s)\n", entityLower))
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// UpdateWithTx
	content.WriteString(fmt.Sprintf("func (p *%s) UpdateWithTx(tx interface{}, %s *domain.%s) error {\n",
		repoName, entityLower, entity))
	content.WriteString("\tgormTx := tx.(*gorm.DB)\n")
	content.WriteString(fmt.Sprintf("\tresult := gormTx.Save(%s)\n", entityLower))
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// DeleteWithTx
	content.WriteString(fmt.Sprintf("func (p *%s) DeleteWithTx(tx interface{}, id int) error {\n", repoName))
	content.WriteString("\tgormTx := tx.(*gorm.DB)\n")
	content.WriteString(fmt.Sprintf("\tresult := gormTx.Delete(&domain.%s{}, id)\n", entity))
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

	writeGoFile(filename, content.String())
}

// generateBasicMongoCRUDMethods genera métodos CRUD básicos para MongoDB
func generateBasicMongoCRUDMethods(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)

	// Save method
	content.WriteString(fmt.Sprintf("func (m *%s) Save(%s *domain.%s) error {\n", repoName, entityLower, entity))
	content.WriteString("\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n")
	content.WriteString("\tdefer cancel()\n")
	content.WriteString(fmt.Sprintf("\t_, err := m.collection.InsertOne(ctx, %s)\n", entityLower))
	content.WriteString("\treturn err\n")
	content.WriteString("}\n\n")

	// FindByID method
	content.WriteString(fmt.Sprintf("func (m *%s) FindByID(id int) (*domain.%s, error) {\n", repoName, entity))
	content.WriteString("\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n")
	content.WriteString("\tdefer cancel()\n")
	content.WriteString(fmt.Sprintf("\t%s := &domain.%s{}\n", entityLower, entity))
	content.WriteString("\terr := m.collection.FindOne(ctx, bson.M{\"id\": id}).Decode(" + entityLower + ")\n")
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\treturn %s, nil\n", entityLower))
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
	repositoryCmd.Flags().StringP("fields", "f", "", "Campos de la entidad \"field:type,field2:type\"")
}
