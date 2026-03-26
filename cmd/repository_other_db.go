package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
)

func generatePostgresJSONRepository(dir, entity string, cache, transactions bool, sm ...*SafetyManager) {
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

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		fmt.Printf("Error creating PostgreSQL JSON repository file: %v\n", err)
	}
}

// generateSQLServerRepository generates a repository for SQL Server with GORM + mssql
func generateSQLServerRepository(dir, entity string, cache, transactions bool, sm ...*SafetyManager) {
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

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		fmt.Printf("Error creating SQL Server repository file: %v\n", err)
	}
}

// generateElasticsearchRepository generates a repository for Elasticsearch with full-text search
func generateElasticsearchRepository(dir, entity string, cache, transactions bool, sm ...*SafetyManager) {
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

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		fmt.Printf("Error creating Elasticsearch repository file: %v\n", err)
	}
}

// generateDynamoDBRepository generates a repository for DynamoDB with AWS SDK v2
func generateDynamoDBRepository(dir, entity string, cache, transactions bool, sm ...*SafetyManager) {
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

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		fmt.Printf("Error creating DynamoDB repository file: %v\n", err)
	}
}

// generateSQLiteRepository generates a repository for SQLite with database/sql
func generateSQLiteRepository(dir, entity string, cache, transactions bool, sm ...*SafetyManager) {
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
	content.WriteString("\tvar data []byte\n")
	content.WriteString(fmt.Sprintf("\tdata, err := json.Marshal(%s)\n", entityLower))
	content.WriteString(fmt.Sprintf("\tif err != nil {\n\t\treturn fmt.Errorf(\"failed to marshal: %%w\", err)\n\t}\n"))
	content.WriteString(fmt.Sprintf("\tquery := \"INSERT INTO %ss (data) VALUES (?)\"\n", entityLower))
	content.WriteString("\tif _, err := s.db.Exec(query, data); err != nil {\n")
	content.WriteString(fmt.Sprintf("\t\treturn fmt.Errorf(\"failed to insert: %%w\", err)\n\t}\n"))
	content.WriteString("\treturn nil\n")
	content.WriteString("}\n\n")

	// FindByID method
	content.WriteString(fmt.Sprintf("func (s *%s) FindByID(id int) (*domain.%s, error) {\n", repoName, entity))
	content.WriteString("\tvar data []byte\n")
	content.WriteString(fmt.Sprintf("\tquery := \"SELECT data FROM %ss WHERE id = ? LIMIT 1\"\n", entityLower))
	content.WriteString("\tif err := s.db.QueryRow(query, id).Scan(&data); err != nil {\n")
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
	content.WriteString("\treturn nil\n")
	content.WriteString("}\n\n")

	// Delete method
	content.WriteString(fmt.Sprintf("func (s *%s) Delete(id int) error {\n", repoName))
	content.WriteString(fmt.Sprintf("\tquery := \"DELETE FROM %ss WHERE id = ?\"\n", entityLower))
	content.WriteString("\tif _, err := s.db.Exec(query, id); err != nil {\n")
	content.WriteString(fmt.Sprintf("\t\treturn fmt.Errorf(\"failed to delete: %%w\", err)\n\t}\n"))
	content.WriteString("\treturn nil\n")
	content.WriteString("}\n\n")

	// FindAll method
	content.WriteString(fmt.Sprintf("func (s *%s) FindAll() ([]domain.%s, error) {\n", repoName, entity))
	content.WriteString(fmt.Sprintf("\tquery := \"SELECT data FROM %ss\"\n", entityLower))
	content.WriteString("\trows, err := s.db.Query(query)\n")
	content.WriteString(fmt.Sprintf("\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"failed to query: %%w\", err)\n\t}\n"))
	content.WriteString("\tdefer rows.Close()\n")
	content.WriteString(fmt.Sprintf("\tvar %ss []domain.%s\n", entityLower, entity))
	content.WriteString("\tfor rows.Next() {\n")
	content.WriteString("\t\tvar data []byte\n")
	content.WriteString(fmt.Sprintf("\t\tif err := rows.Scan(&data); err != nil {\n\t\t\treturn nil, fmt.Errorf(\"failed to scan: %%w\", err)\n\t\t}\n"))
	content.WriteString(fmt.Sprintf("\t\tvar %s domain.%s\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\t\tif err := json.Unmarshal(data, &%s); err != nil {\n\t\t\treturn nil, fmt.Errorf(\"failed to unmarshal: %%w\", err)\n\t\t}\n", entityLower))
	content.WriteString(fmt.Sprintf("\t\t%ss = append(%ss, %s)\n", entityLower, entityLower, entityLower))
	content.WriteString("\t}\n")
	content.WriteString(fmt.Sprintf("\tif err := rows.Err(); err != nil {\n\t\treturn nil, fmt.Errorf(\"rows error: %%w\", err)\n\t}\n"))
	content.WriteString(fmt.Sprintf("\treturn %ss, nil\n", entityLower))
	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		fmt.Printf("Error creating SQLite repository file: %v\n", err)
	}
}

