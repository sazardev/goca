package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
)

func generatePostgresRepository(dir, entity string, cache, transactions bool, sm ...*SafetyManager) {
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

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
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

func generateMySQLRepository(dir, entity string, cache, transactions bool, sm ...*SafetyManager) {
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

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		fmt.Printf("Error creating MySQL repository file: %v\n", err)
	}
}

func generateMongoRepository(dir, entity string, cache, transactions bool, sm ...*SafetyManager) {
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
		content.WriteString("\t// MongoDB cache imports\n")
		content.WriteString("\t// \"github.com/go-redis/redis/v8\"\n")
	}
	if transactions {
		content.WriteString("\t// MongoDB transaction support\n")
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

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		fmt.Printf("Error creating MongoDB repository file: %v\n", err)
	}
}

// generateRepositoryInterfaceWithFields generates repository interfaces with dynamic methods based on fields
