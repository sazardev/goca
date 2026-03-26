package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func generateRepositoryInterfaceWithFields(dir, entity string, fields []Field, transactions bool, sm ...*SafetyManager) {
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

	// Generate dynamic search methods based on actual fields
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

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		fmt.Printf("Error writing file %s: %v\n", filename, err)
	}
}

// generateRepositoryImplementationWithFields generates repository implementations with dynamic methods
func generateRepositoryImplementationWithFields(dir, entity, database string, fields []Field, cache, transactions bool, sm ...*SafetyManager) {
	switch database {
	case DBPostgres:
		generatePostgresRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
	case DBMySQL:
		generateMySQLRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
	case DBMongoDB:
		generateMongoRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
	default:
		generatePostgresRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
	}
}

// generatePostgresRepositoryWithFields generates PostgreSQL repository with dynamic methods
func generatePostgresRepositoryWithFields(dir, entity string, fields []Field, cache, transactions bool, sm ...*SafetyManager) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, "postgres_"+entityLower+"_repository.go")

	// Get the module name from go.mod
	moduleName := getModuleName()

	var content strings.Builder
	content.WriteString("package repository\n\n")
	content.WriteString("import (\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
	if cache {
		content.WriteString("\t// Cache imports (Redis, etc.)\n")
		content.WriteString("\t// \"github.com/go-redis/redis/v8\"\n")
	}
	if transactions {
		content.WriteString("\t// SQL transaction support\n")
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

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		fmt.Printf("Error creating PostgreSQL repository with fields: %v\n", err)
	}
}

// generateBasicCRUDMethods generates basic CRUD methods
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

// generateTransactionMethods generates methods that support transactions
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

// generateMySQLRepositoryWithFields generates MySQL repository with dynamic methods
func generateMySQLRepositoryWithFields(dir, entity string, fields []Field, cache, transactions bool, sm ...*SafetyManager) {
	// For MySQL we use the same logic as PostgreSQL since both use GORM
	generatePostgresRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
}

// generateMongoRepositoryWithFields generates MongoDB repository with dynamic methods
func generateMongoRepositoryWithFields(dir, entity string, fields []Field, cache, transactions bool, sm ...*SafetyManager) {
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

	// Generate basic MongoDB methods (simplified for brevity)
	generateBasicMongoCRUDMethods(&content, entity, repoName)

	// Generate dynamic search methods for MongoDB
	searchMethods := generateSearchMethods(fields, entity)
	for _, method := range searchMethods {
		content.WriteString(generateMongoSearchMethodImplementation(method, repoName, entity))
	}

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		fmt.Printf("Error creating MongoDB repository with fields: %v\n", err)
	}
}

// generateBasicMongoCRUDMethods generates basic CRUD methods for MongoDB
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
	content.WriteString("// Other basic CRUD methods for MongoDB...\n\n")
}

// generateMongoSearchMethodImplementation generates search method implementation for MongoDB
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

// generatePostgresJSONRepository generates a repository for PostgreSQL with JSONB support
