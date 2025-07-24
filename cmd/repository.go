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

		database, _ := cmd.Flags().GetString("database")
		interfaceOnly, _ := cmd.Flags().GetBool("interface-only")
		implementation, _ := cmd.Flags().GetBool("implementation")
		cache, _ := cmd.Flags().GetBool("cache")
		transactions, _ := cmd.Flags().GetBool("transactions")

		fmt.Printf("Generando repositorio para entidad '%s'\n", entity)

		if database != "" && !interfaceOnly {
			fmt.Printf("Base de datos: %s\n", database)
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

		generateRepository(entity, database, interfaceOnly, implementation, cache, transactions)
		fmt.Printf("\n✅ Repositorio para '%s' generado exitosamente!\n", entity)
	},
}

func generateRepository(entity, database string, interfaceOnly, implementation, cache, transactions bool) {
	// Crear directorio repositories si no existe
	repoDir := "internal/infrastructure/repository"
	_ = os.MkdirAll(repoDir, 0755) // Generate interface if not interface-only or if implementation is requested
	if !interfaceOnly || implementation {
		generateRepositoryInterface(repoDir, entity, transactions)
	}

	// Generate implementation if not interface-only and database is specified
	if !interfaceOnly && database != "" {
		generateRepositoryImplementation(repoDir, entity, database, cache, transactions)
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
		content.WriteString(fmt.Sprintf("import \"%s/internal/domain\"\n\n", moduleName))
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

	writeFile(filename, content.String())
}

func generateRepositoryImplementation(dir, entity, database string, cache, transactions bool) {
	switch database {
	case "postgres":
		generatePostgresRepository(dir, entity, cache, transactions)
	case "mysql":
		generateMySQLRepository(dir, entity, cache, transactions)
	case "mongodb":
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
	content.WriteString("\t\"database/sql\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", moduleName))
	if cache {
		content.WriteString("\t\"time\"\n")
		content.WriteString("\t\"encoding/json\"\n")
		content.WriteString("\t\"github.com/go-redis/redis/v8\"\n")
		content.WriteString("\t\"context\"\n")
	}
	content.WriteString("\n\t_ \"github.com/lib/pq\"\n")
	content.WriteString(")\n\n")

	// Repository struct
	repoName := fmt.Sprintf("postgres%sRepository", entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n", repoName))
	content.WriteString("\tdb *sql.DB\n")
	if cache {
		content.WriteString("\tcache *redis.Client\n")
		content.WriteString("\tcacheTTL time.Duration\n")
	}
	content.WriteString("}\n\n")

	// Constructor
	content.WriteString(fmt.Sprintf("func NewPostgres%sRepository(db *sql.DB", entity))
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

	writeFile(filename, content.String())
}

func generatePostgresSaveMethod(content *strings.Builder, entity, repoName string, cache bool) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) Save(%s *domain.%s) error {\n",
		repoVar, repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\t// TODO: Customize this query based on your %s entity fields\n", entity))
	content.WriteString(fmt.Sprintf("\tquery := `INSERT INTO %ss DEFAULT VALUES RETURNING id`\n", entityLower))
	content.WriteString(fmt.Sprintf("\terr := %s.db.QueryRow(query).Scan(&%s.ID)\n",
		repoVar, entityLower))

	if cache {
		content.WriteString("\tif err == nil {\n")
		content.WriteString(fmt.Sprintf("\t\t%s.invalidateCache(%s.ID)\n", repoVar, entityLower))
		content.WriteString("\t}\n")
	}

	content.WriteString("\treturn err\n")
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
	content.WriteString(fmt.Sprintf("\t// TODO: Customize this query based on your %s entity fields\n", entity))
	content.WriteString(fmt.Sprintf("\tquery := `SELECT id FROM %ss WHERE id = $1`\n", entityLower))
	content.WriteString(fmt.Sprintf("\terr := %s.db.QueryRow(query, id).Scan(&%s.ID)\n",
		repoVar, entityLower))
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn nil, err\n")
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
	content.WriteString(fmt.Sprintf("\t// TODO: Customize this query based on your %s entity fields\n", entity))
	content.WriteString(fmt.Sprintf("\tquery := `SELECT id FROM %ss WHERE id = $1 LIMIT 1`\n", entityLower))
	content.WriteString(fmt.Sprintf("\terr := %s.db.QueryRow(query, email).Scan(&%s.ID)\n",
		repoVar, entityLower))
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\treturn %s, nil\n", entityLower))
	content.WriteString("}\n\n")
}

func generatePostgresUpdateMethod(content *strings.Builder, entity, repoName string, cache bool) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) Update(%s *domain.%s) error {\n",
		repoVar, repoName, entityLower, entity))
	content.WriteString(fmt.Sprintf("\t// TODO: Customize this query based on your %s entity fields\n", entity))
	content.WriteString(fmt.Sprintf("\tquery := `UPDATE %ss SET id = $1 WHERE id = $2`\n", entityLower))
	content.WriteString(fmt.Sprintf("\t_, err := %s.db.Exec(query, %s.ID, %s.ID)\n",
		repoVar, entityLower, entityLower))

	if cache {
		content.WriteString("\tif err == nil {\n")
		content.WriteString(fmt.Sprintf("\t\t%s.invalidateCache(%s.ID)\n", repoVar, entityLower))
		content.WriteString("\t}\n")
	}

	content.WriteString("\treturn err\n")
	content.WriteString("}\n\n")
}

func generatePostgresDeleteMethod(content *strings.Builder, entity, repoName string, cache bool) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) Delete(id int) error {\n",
		repoVar, repoName))
	content.WriteString(fmt.Sprintf("\tquery := `DELETE FROM %ss WHERE id = $1`\n", entityLower))
	content.WriteString(fmt.Sprintf("\t_, err := %s.db.Exec(query, id)\n", repoVar))

	if cache {
		content.WriteString("\tif err == nil {\n")
		content.WriteString(fmt.Sprintf("\t\t%s.invalidateCache(id)\n", repoVar))
		content.WriteString("\t}\n")
	}

	content.WriteString("\treturn err\n")
	content.WriteString("}\n\n")
}

func generatePostgresFindAllMethod(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)
	repoVar := strings.ToLower(string(repoName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) FindAll() ([]domain.%s, error) {\n",
		repoVar, repoName, entity))
	content.WriteString(fmt.Sprintf("\t// TODO: Customize this query based on your %s entity fields\n", entity))
	content.WriteString(fmt.Sprintf("\tquery := `SELECT id FROM %ss`\n", entityLower))
	content.WriteString(fmt.Sprintf("\trows, err := %s.db.Query(query)\n", repoVar))
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n")
	content.WriteString("\tdefer rows.Close()\n\n")

	content.WriteString(fmt.Sprintf("\tvar %ss []domain.%s\n", entityLower, entity))
	content.WriteString("\tfor rows.Next() {\n")
	content.WriteString(fmt.Sprintf("\t\tvar %s domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\t\t// TODO: Scan all fields of your %s entity\n", entity))
	content.WriteString(fmt.Sprintf("\t\tif err := rows.Scan(&%s.ID); err != nil {\n", entityLower))
	content.WriteString("\t\t\treturn nil, err\n")
	content.WriteString("\t\t}\n")
	content.WriteString(fmt.Sprintf("\t\t%ss = append(%ss, %s)\n", entityLower, entityLower, entityLower))
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
	content.WriteString("\tsqlTx := tx.(*sql.Tx)\n")
	content.WriteString(fmt.Sprintf("\t// TODO: Customize this query based on your %s entity fields\n", entity))
	content.WriteString(fmt.Sprintf("\tquery := `INSERT INTO %ss DEFAULT VALUES RETURNING id`\n", entityLower))
	content.WriteString(fmt.Sprintf("\treturn sqlTx.QueryRow(query).Scan(&%s.ID)\n", entityLower))
	content.WriteString("}\n\n")

	// UpdateWithTx
	content.WriteString(fmt.Sprintf("func (%s *%s) UpdateWithTx(tx interface{}, %s *domain.%s) error {\n",
		repoVar, repoName, entityLower, entity))
	content.WriteString("\tsqlTx := tx.(*sql.Tx)\n")
	content.WriteString(fmt.Sprintf("\t// TODO: Customize this query based on your %s entity fields\n", entity))
	content.WriteString(fmt.Sprintf("\tquery := `UPDATE %ss SET updated_at = NOW() WHERE id = $1`\n", entityLower))
	content.WriteString(fmt.Sprintf("\t_, err := sqlTx.Exec(query, %s.ID)\n", entityLower))
	content.WriteString("\treturn err\n")
	content.WriteString("}\n\n")

	// DeleteWithTx
	content.WriteString(fmt.Sprintf("func (%s *%s) DeleteWithTx(tx interface{}, id int) error {\n",
		repoVar, repoName))
	content.WriteString("\tsqlTx := tx.(*sql.Tx)\n")
	content.WriteString(fmt.Sprintf("\tquery := `DELETE FROM %ss WHERE id = $1`\n", entityLower))
	content.WriteString("\t_, err := sqlTx.Exec(query, id)\n")
	content.WriteString("\treturn err\n")
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
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", moduleName))
	if cache {
		content.WriteString("\t// TODO: Add cache imports when cache is implemented\n")
	}
	if transactions {
		content.WriteString("\t// TODO: Add transaction support\n")
	}
	content.WriteString("package repository\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"database/sql\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", moduleName))
	content.WriteString("\n\t_ \"github.com/go-sql-driver/mysql\"\n")
	content.WriteString(")\n\n")

	// Similar structure to Postgres but with MySQL syntax
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
	content.WriteString(fmt.Sprintf("\t// TODO: Customize this query based on your %s entity fields\n", entity))
	content.WriteString(fmt.Sprintf("\tquery := `INSERT INTO %ss () VALUES ()`\n", entityLower))
	content.WriteString("\tresult, err := r.db.Exec(query)\n")
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn err\n")
	content.WriteString("\t}\n")
	content.WriteString("\tid, err := result.LastInsertId()\n")
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn err\n")
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\t%s.ID = int(id)\n", entityLower))
	content.WriteString("\treturn nil\n")
	content.WriteString("}\n\n")

	writeFile(filename, content.String())
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
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", moduleName))
	if cache {
		content.WriteString("\t// TODO: Add cache imports when cache is implemented\n")
	}
	if transactions {
		content.WriteString("\t// TODO: Add transaction support for MongoDB\n")
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

	writeFile(filename, content.String())
}

func init() {
	repositoryCmd.Flags().StringP("database", "d", "", "Tipo de base de datos (postgres, mysql, mongodb)")
	repositoryCmd.Flags().BoolP("interface-only", "i", false, "Solo generar interfaces")
	repositoryCmd.Flags().BoolP("implementation", "", false, "Solo generar implementación")
	repositoryCmd.Flags().BoolP("cache", "c", false, "Incluir capa de caché")
	repositoryCmd.Flags().BoolP("transactions", "t", false, "Incluir soporte para transacciones")
}
