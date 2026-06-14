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
			// If this entity's interface already exists, leave the file as-is
			// (rewriting it would drop the package/imports header).
			if strings.Contains(existingStr, fmt.Sprintf("type %sRepository interface", entity)) {
				return
			}
			// Otherwise append to the existing content (without the trailing newline).
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

	// Generate dynamic search methods based on actual fields
	searchMethods := generateSearchMethods(fields, entity)
	for _, method := range searchMethods {
		content.WriteString(method.generateSearchMethodSignature() + "\n")
	}

	content.WriteString(fmt.Sprintf("\tUpdate(%s *domain.%s) error\n", strings.ToLower(entity), entity))
	content.WriteString("\tDelete(id int) error\n")
	content.WriteString(fmt.Sprintf("\tFindAll() ([]domain.%s, error)\n", entity))

	if transactions {
		content.WriteString(fmt.Sprintf("\tSaveWithTx(tx *gorm.DB, %s *domain.%s) error\n", strings.ToLower(entity), entity))
		content.WriteString(fmt.Sprintf("\tUpdateWithTx(tx *gorm.DB, %s *domain.%s) error\n", strings.ToLower(entity), entity))
		content.WriteString("\tDeleteWithTx(tx *gorm.DB, id int) error\n")
	}

	content.WriteString("}\n\n")

	if err := writeGoFileMerged(filename, content.String(), sm...); err != nil {
		fmt.Printf("Error writing file %s: %v\n", filename, err)
	}
}

// generateRepositoryImplementationWithFields generates repository implementations with dynamic methods.
func generateRepositoryImplementationWithFields(dir, entity, database string, fields []Field, cache, transactions bool, sm ...*SafetyManager) {
	switch database {
	case DBPostgres:
		generatePostgresRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
	case DBMySQL:
		generateMySQLRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
	case DBMongoDB:
		generateMongoRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
	case DBPostgresJSON:
		generatePostgresJSONRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
	case DBSQLServer:
		generateSQLServerRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
	case DBSQLite:
		generateSQLiteRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
	case DBElasticsearch:
		generateElasticsearchRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
	case DBDynamoDB:
		generateDynamoDBRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
	default:
		generatePostgresRepositoryWithFields(dir, entity, fields, cache, transactions, sm...)
	}
}

// The dedicated DB generators (in repository_other_db.go) already emit a full
// CRUD set. For the field-aware path the interface additionally declares the
// per-field finders, so we generate the base repository and then append the
// finder implementations rendered in the backend's native style.

func generatePostgresJSONRepositoryWithFields(dir, entity string, fields []Field, cache, transactions bool, sm ...*SafetyManager) {
	generatePostgresJSONRepository(dir, entity, cache, transactions, sm...)
	appendGormFinders(dir, "postgres_json_"+strings.ToLower(entity)+"_repository.go", fmt.Sprintf("postgresJSON%sRepository", entity), entity, fields, sm...)
}

func generateSQLServerRepositoryWithFields(dir, entity string, fields []Field, cache, transactions bool, sm ...*SafetyManager) {
	generateSQLServerRepository(dir, entity, cache, transactions, sm...)
	appendGormFinders(dir, "sqlserver_"+strings.ToLower(entity)+"_repository.go", fmt.Sprintf("sqlserver%sRepository", entity), entity, fields, sm...)
}

func generateSQLiteRepositoryWithFields(dir, entity string, fields []Field, cache, transactions bool, sm ...*SafetyManager) {
	generateSQLiteRepository(dir, entity, cache, transactions, sm...)
	appendSQLiteFinders(dir, entity, fields, sm...)
}

func generateElasticsearchRepositoryWithFields(dir, entity string, fields []Field, cache, transactions bool, sm ...*SafetyManager) {
	generateElasticsearchRepository(dir, entity, cache, transactions, sm...)
	appendDelegatingFinders(dir, "elasticsearch_"+strings.ToLower(entity)+"_repository.go", fmt.Sprintf("elasticsearch%sRepository", entity), "e", entity, fields, sm...)
}

func generateDynamoDBRepositoryWithFields(dir, entity string, fields []Field, cache, transactions bool, sm ...*SafetyManager) {
	generateDynamoDBRepository(dir, entity, cache, transactions, sm...)
	appendDelegatingFinders(dir, "dynamodb_"+strings.ToLower(entity)+"_repository.go", fmt.Sprintf("dynamodb%sRepository", entity), "d", entity, fields, sm...)
}

// appendGormFinders appends GORM-based per-field finder implementations to an
// already-generated repository file whose receiver exposes a `db *gorm.DB`.
func appendGormFinders(dir, file, repoName, entity string, fields []Field, sm ...*SafetyManager) {
	methods := generateSearchMethods(fields, entity)
	if len(methods) == 0 {
		return
	}
	var b strings.Builder
	for _, m := range methods {
		b.WriteString(m.generateSearchMethodImplementation(strings.ToLower(string(repoName[0])), repoName, entity))
	}
	appendToRepoFile(filepath.Join(dir, file), b.String(), nil, sm...)
}

// appendSQLiteFinders appends raw-SQL per-field finders for the SQLite repo.
func appendSQLiteFinders(dir, entity string, fields []Field, sm ...*SafetyManager) {
	methods := generateSearchMethods(fields, entity)
	if len(methods) == 0 {
		return
	}
	entityLower := strings.ToLower(entity)
	repoName := fmt.Sprintf("sqlite%sRepository", entity)
	var b strings.Builder
	for _, m := range methods {
		paramName := strings.ToLower(m.FieldName)
		fmt.Fprintf(&b, "func (s *%s) %s(%s %s) %s {\n", repoName, m.MethodName, paramName, m.FieldType, m.ReturnType)
		b.WriteString("\tvar data []byte\n")
		fmt.Fprintf(&b, "\tquery := \"SELECT data FROM %ss WHERE json_extract(data, '$.%s') = ? LIMIT 1\"\n", entityLower, m.FieldName)
		fmt.Fprintf(&b, "\tif err := s.db.QueryRow(query, %s).Scan(&data); err != nil {\n", paramName)
		fmt.Fprintf(&b, "\t\tif err == sql.ErrNoRows {\n\t\t\treturn nil, fmt.Errorf(\"%s not found\")\n\t\t}\n", entity)
		b.WriteString("\t\treturn nil, fmt.Errorf(\"failed to query: %w\", err)\n\t}\n")
		fmt.Fprintf(&b, "\tvar %s domain.%s\n", entityLower, entity)
		fmt.Fprintf(&b, "\tif err := json.Unmarshal(data, &%s); err != nil {\n", entityLower)
		b.WriteString("\t\treturn nil, fmt.Errorf(\"failed to unmarshal: %w\", err)\n\t}\n")
		fmt.Fprintf(&b, "\treturn &%s, nil\n", entityLower)
		b.WriteString("}\n\n")
	}
	appendToRepoFile(filepath.Join(dir, "sqlite_"+entityLower+"_repository.go"), b.String(), nil, sm...)
}

// appendDelegatingFinders appends per-field finders that reuse FindAll and filter
// in memory — used for backends (Elasticsearch, DynamoDB) where a dedicated query
// per field is out of scope but the interface still requires the method.
func appendDelegatingFinders(dir, file, repoName, recv, entity string, fields []Field, sm ...*SafetyManager) {
	methods := generateSearchMethods(fields, entity)
	if len(methods) == 0 {
		return
	}
	entityLower := strings.ToLower(entity)
	var b strings.Builder
	for _, m := range methods {
		paramName := strings.ToLower(m.FieldName)
		fmt.Fprintf(&b, "func (%s *%s) %s(%s %s) %s {\n", recv, repoName, m.MethodName, paramName, m.FieldType, m.ReturnType)
		fmt.Fprintf(&b, "\titems, err := %s.FindAll()\n", recv)
		b.WriteString("\tif err != nil {\n\t\treturn nil, err\n\t}\n")
		b.WriteString("\tfor i := range items {\n")
		fmt.Fprintf(&b, "\t\tif items[i].%s == %s {\n", m.FieldName, paramName)
		b.WriteString("\t\t\treturn &items[i], nil\n\t\t}\n")
		b.WriteString("\t}\n")
		fmt.Fprintf(&b, "\treturn nil, fmt.Errorf(\"%s not found\")\n", entityLower)
		b.WriteString("}\n\n")
	}
	appendToRepoFile(filepath.Join(dir, file), b.String(), []string{"fmt"}, sm...)
}

// appendToRepoFile appends generated method source to an existing repository
// file, respecting the SafetyManager (dry-run) if provided. ensureImports lists
// stdlib import paths the appended code requires; any not already present are
// injected into the file's import block.
func appendToRepoFile(path, methods string, ensureImports []string, sm ...*SafetyManager) {
	if strings.TrimSpace(methods) == "" {
		return
	}
	if len(sm) > 0 && sm[0] != nil && sm[0].DryRun {
		return
	}
	existing, err := os.ReadFile(path)
	if err != nil {
		return
	}
	src := string(existing)
	for _, imp := range ensureImports {
		quoted := "\"" + imp + "\""
		if strings.Contains(src, quoted) {
			continue
		}
		// Inject right after the opening of the import block.
		if idx := strings.Index(src, "import (\n"); idx != -1 {
			pos := idx + len("import (\n")
			src = src[:pos] + "\t" + quoted + "\n" + src[pos:]
		}
	}
	combined := strings.TrimRight(src, "\n") + "\n\n" + methods
	if err := writeGoFile(path, combined); err != nil {
		fmt.Printf("Error appending finders to %s: %v\n", path, err)
	}
}

// generatePostgresRepositoryWithFields generates a PostgreSQL repository.
func generatePostgresRepositoryWithFields(dir, entity string, fields []Field, cache, transactions bool, sm ...*SafetyManager) {
	generateGormRepositoryWithFields(dir, entity, "Postgres", fields, cache, transactions, sm...)
}

// generateGormRepositoryWithFields generates a GORM-based repository for the
// given driver (e.g. "Postgres", "MySQL"). The constructor and file name are
// driver-prefixed so they match what the DI container references
// (New<Driver><Entity>Repository).
func generateGormRepositoryWithFields(dir, entity, driver string, fields []Field, cache, transactions bool, sm ...*SafetyManager) {
	entityLower := strings.ToLower(entity)
	driverLower := strings.ToLower(driver)
	filename := filepath.Join(dir, driverLower+"_"+entityLower+"_repository.go")

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
	repoName := fmt.Sprintf("%s%sRepository", driverLower, entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n", repoName))
	content.WriteString("\tdb *gorm.DB\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("func New%s%sRepository(db *gorm.DB) %sRepository {\n", driver, entity, entity))
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
		fmt.Printf("Error creating %s repository with fields: %v\n", driver, err)
	}
}

// generateBasicCRUDMethods generates basic CRUD methods.
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

// generateTransactionMethods generates methods that support transactions.
func generateTransactionMethods(content *strings.Builder, entity, repoName string) {
	entityLower := strings.ToLower(entity)

	// SaveWithTx
	fmt.Fprintf(content, "func (p *%s) SaveWithTx(tx *gorm.DB, %s *domain.%s) error {\n",
		repoName, entityLower, entity)
	fmt.Fprintf(content, "\tresult := tx.Create(%s)\n", entityLower)
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// UpdateWithTx
	fmt.Fprintf(content, "func (p *%s) UpdateWithTx(tx *gorm.DB, %s *domain.%s) error {\n",
		repoName, entityLower, entity)
	fmt.Fprintf(content, "\tresult := tx.Save(%s)\n", entityLower)
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")

	// DeleteWithTx
	fmt.Fprintf(content, "func (p *%s) DeleteWithTx(tx *gorm.DB, id int) error {\n", repoName)
	fmt.Fprintf(content, "\tresult := tx.Delete(&domain.%s{}, id)\n", entity)
	content.WriteString("\treturn result.Error\n")
	content.WriteString("}\n\n")
}

// generateMySQLRepositoryWithFields generates MySQL repository with dynamic methods.
func generateMySQLRepositoryWithFields(dir, entity string, fields []Field, cache, transactions bool, sm ...*SafetyManager) {
	// All SQL databases share the same GORM repository (the concrete driver is
	// chosen by the dialector in main.go), so they use a single
	// NewPostgres<Entity>Repository(*gorm.DB) constructor that the DI container
	// references uniformly.
	generateGormRepositoryWithFields(dir, entity, "Postgres", fields, cache, transactions, sm...)
}

// generateMongoRepositoryWithFields generates MongoDB repository with dynamic methods.
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

// generateBasicMongoCRUDMethods generates basic CRUD methods for MongoDB.
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

	// Update method
	fmt.Fprintf(content, "func (m *%s) Update(%s *domain.%s) error {\n", repoName, entityLower, entity)
	content.WriteString("\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n")
	content.WriteString("\tdefer cancel()\n")
	fmt.Fprintf(content, "\t_, err := m.collection.ReplaceOne(ctx, bson.M{\"id\": %s.ID}, %s)\n", entityLower, entityLower)
	content.WriteString("\treturn err\n")
	content.WriteString("}\n\n")

	// Delete method
	fmt.Fprintf(content, "func (m *%s) Delete(id int) error {\n", repoName)
	content.WriteString("\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n")
	content.WriteString("\tdefer cancel()\n")
	content.WriteString("\t_, err := m.collection.DeleteOne(ctx, bson.M{\"id\": id})\n")
	content.WriteString("\treturn err\n")
	content.WriteString("}\n\n")

	// FindAll method
	fmt.Fprintf(content, "func (m *%s) FindAll() ([]domain.%s, error) {\n", repoName, entity)
	content.WriteString("\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n")
	content.WriteString("\tdefer cancel()\n")
	content.WriteString("\tcursor, err := m.collection.Find(ctx, bson.M{})\n")
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n")
	content.WriteString("\tdefer cursor.Close(ctx)\n")
	fmt.Fprintf(content, "\tvar %ss []domain.%s\n", entityLower, entity)
	fmt.Fprintf(content, "\tif err := cursor.All(ctx, &%ss); err != nil {\n", entityLower)
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n")
	fmt.Fprintf(content, "\treturn %ss, nil\n", entityLower)
	content.WriteString("}\n\n")
}

// generateMongoSearchMethodImplementation generates search method implementation for MongoDB.
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
