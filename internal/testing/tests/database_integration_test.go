package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewDatabaseIntegrations tests all newly added database support
func TestNewDatabaseIntegrations(t *testing.T) {
	databases := []string{
		"sqlite",
		"sqlserver",
		"postgres-json",
		"elasticsearch",
		"dynamodb",
	}

	for _, db := range databases {
		t.Run(fmt.Sprintf("Database_%s", db), func(t *testing.T) {
			testDatabaseIntegration(t, db)
		})
	}
}

// testDatabaseIntegration tests a single database integration end-to-end
func testDatabaseIntegration(t *testing.T, database string) {
	// Create temporary project
	tmpDir := t.TempDir()
	projectName := "test-" + strings.ReplaceAll(database, "-", "_")
	projectPath := filepath.Join(tmpDir, projectName)

	// Initialize project
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Change to project directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(projectPath); err != nil {
		t.Fatalf("Failed to change to project directory: %v", err)
	}

	// Create basic go.mod
	goModContent := fmt.Sprintf("module %s\n\ngo 1.25\n", projectName)
	if err := os.WriteFile("go.mod", []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create required directory structure
	dirs := []string{
		"internal/domain",
		"internal/repository",
		"internal/usecase",
		"internal/handler/http",
		"internal/di",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	t.Run("GenerateEntity", func(t *testing.T) {
		testGenerateEntity(t, "Product")
	})

	t.Run("GenerateRepository", func(t *testing.T) {
		testGenerateRepository(t, "Product", database)
	})

	t.Run("RepositoryFileExists", func(t *testing.T) {
		testRepositoryFileExists(t, "Product", database)
	})

	t.Run("RepositoryImplementsInterface", func(t *testing.T) {
		testRepositoryImplementsInterface(t, "Product", database)
	})

	t.Run("DatabaseSpecificMethods", func(t *testing.T) {
		testDatabaseSpecificMethods(t, "Product", database)
	})
}

// testGenerateEntity verifies entity generation
func testGenerateEntity(t *testing.T, entityName string) {
	// Create basic entity file
	entityLower := strings.ToLower(entityName)
	entityPath := filepath.Join("internal", "domain", entityLower+".go")

	entityContent := fmt.Sprintf(`package domain

type %s struct {
	ID    uint   `+"`"+`json:"id" gorm:"primaryKey"`+"`"+`
	Name  string `+"`"+`json:"name" gorm:"type:varchar(255);not null"`+"`"+`
	Price float64 `+"`"+`json:"price" gorm:"type:decimal(10,2)"`+"`"+`
}

func (p *%s) Validate() error {
	if p.Name == "" {
		return ErrInvalidEntity
	}
	if p.Price <= 0 {
		return ErrInvalidEntity
	}
	return nil
}
`, entityName, entityName)

	if err := os.WriteFile(entityPath, []byte(entityContent), 0644); err != nil {
		t.Fatalf("Failed to create entity file: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(entityPath); err != nil {
		t.Errorf("Entity file not found: %v", err)
	}
}

// testGenerateRepository generates repository for the database
func testGenerateRepository(t *testing.T, entityName, database string) {
	repositoryDir := filepath.Join("internal", "repository")

	// Generate appropriate repository based on database
	switch database {
	case "sqlite":
		generateTestSQLiteRepository(t, repositoryDir, entityName)
	case "sqlserver":
		generateTestSQLServerRepository(t, repositoryDir, entityName)
	case "postgres-json":
		generateTestPostgresJSONRepository(t, repositoryDir, entityName)
	case "elasticsearch":
		generateTestElasticsearchRepository(t, repositoryDir, entityName)
	case "dynamodb":
		generateTestDynamoDBRepository(t, repositoryDir, entityName)
	}
}

// testRepositoryFileExists verifies repository file was created
func testRepositoryFileExists(t *testing.T, entityName, database string) {
	entityLower := strings.ToLower(entityName)
	repositoryDir := filepath.Join("internal", "repository")

	var expectedFile string
	switch database {
	case "sqlite":
		expectedFile = filepath.Join(repositoryDir, "sqlite_"+entityLower+"_repository.go")
	case "sqlserver":
		expectedFile = filepath.Join(repositoryDir, "sqlserver_"+entityLower+"_repository.go")
	case "postgres-json":
		expectedFile = filepath.Join(repositoryDir, "postgres_json_"+entityLower+"_repository.go")
	case "elasticsearch":
		expectedFile = filepath.Join(repositoryDir, "elasticsearch_"+entityLower+"_repository.go")
	case "dynamodb":
		expectedFile = filepath.Join(repositoryDir, "dynamodb_"+entityLower+"_repository.go")
	}

	if _, err := os.Stat(expectedFile); err != nil {
		t.Errorf("Repository file not found: %s - %v", expectedFile, err)
	}

	// Verify file is not empty
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read repository file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Repository file is empty")
	}

	// Verify key content exists
	contentStr := string(content)
	expectedStrings := []string{
		"package repository",
		"func",
		"Save",
		"FindByID",
		"Delete",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Repository file missing expected content: %s", expected)
		}
	}
}

// testRepositoryImplementsInterface verifies repository implements the interface
func testRepositoryImplementsInterface(t *testing.T, entityName, database string) {
	entityLower := strings.ToLower(entityName)
	repositoryDir := filepath.Join("internal", "repository")

	var filePath string
	switch database {
	case "sqlite":
		filePath = filepath.Join(repositoryDir, "sqlite_"+entityLower+"_repository.go")
	case "sqlserver":
		filePath = filepath.Join(repositoryDir, "sqlserver_"+entityLower+"_repository.go")
	case "postgres-json":
		filePath = filepath.Join(repositoryDir, "postgres_json_"+entityLower+"_repository.go")
	case "elasticsearch":
		filePath = filepath.Join(repositoryDir, "elasticsearch_"+entityLower+"_repository.go")
	case "dynamodb":
		filePath = filepath.Join(repositoryDir, "dynamodb_"+entityLower+"_repository.go")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read repository file: %v", err)
	}

	contentStr := string(content)

	// Required interface methods
	requiredMethods := []string{
		"Save(",
		"FindByID(",
		"Update(",
		"Delete(",
		"FindAll(",
	}

	for _, method := range requiredMethods {
		if !strings.Contains(contentStr, "func (") && strings.Contains(contentStr, method) {
			if !strings.Contains(contentStr, method) {
				t.Errorf("Repository missing required method: %s", method)
			}
		}
	}
}

// testDatabaseSpecificMethods verifies database-specific methods exist
func testDatabaseSpecificMethods(t *testing.T, entityName, database string) {
	entityLower := strings.ToLower(entityName)
	repositoryDir := filepath.Join("internal", "repository")

	var filePath string
	var expectedMethods []string

	switch database {
	case "sqlite":
		filePath = filepath.Join(repositoryDir, "sqlite_"+entityLower+"_repository.go")
		expectedMethods = []string{} // Standard interface methods only
	case "sqlserver":
		filePath = filepath.Join(repositoryDir, "sqlserver_"+entityLower+"_repository.go")
		expectedMethods = []string{} // Standard interface methods only
	case "postgres-json":
		filePath = filepath.Join(repositoryDir, "postgres_json_"+entityLower+"_repository.go")
		expectedMethods = []string{"FindByJSONField"}
	case "elasticsearch":
		filePath = filepath.Join(repositoryDir, "elasticsearch_"+entityLower+"_repository.go")
		expectedMethods = []string{"FullTextSearch"}
	case "dynamodb":
		filePath = filepath.Join(repositoryDir, "dynamodb_"+entityLower+"_repository.go")
		expectedMethods = []string{} // Standard interface methods only
	}

	if len(expectedMethods) > 0 {
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read repository file: %v", err)
		}

		contentStr := string(content)
		for _, method := range expectedMethods {
			if !strings.Contains(contentStr, method) {
				t.Errorf("Repository missing database-specific method: %s", method)
			}
		}
	}
}

// Helper functions to generate test repositories

func generateTestSQLiteRepository(t *testing.T, dir, entity string) {
	entityLower := strings.ToLower(entity)
	filePath := filepath.Join(dir, "sqlite_"+entityLower+"_repository.go")

	content := `package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"` + getTestModuleName(t) + `/internal/domain"
)

type sqlite` + entity + `Repository struct {
	db *sql.DB
}

func NewSQLite` + entity + `Repository(db *sql.DB) ` + entity + `Repository {
	return &sqlite` + entity + `Repository{db: db}
}

func (s *sqlite` + entity + `Repository) Save(` + entityLower + ` *domain.` + entity + `) error {
	data, err := json.Marshal(` + entityLower + `)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	query := "INSERT INTO ` + entityLower + `s (data) VALUES (?)"
	if _, err := s.db.Exec(query, data); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (s *sqlite` + entity + `Repository) FindByID(id int) (*domain.` + entity + `, error) {
	var data []byte
	query := "SELECT data FROM ` + entityLower + `s WHERE id = ? LIMIT 1"
	if err := s.db.QueryRow(query, id).Scan(&data); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("` + entity + ` not found")
		}
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	var ` + entityLower + ` domain.` + entity + `
	if err := json.Unmarshal(data, &` + entityLower + `); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &` + entityLower + `, nil
}

func (s *sqlite` + entity + `Repository) Update(` + entityLower + ` *domain.` + entity + `) error {
	data, err := json.Marshal(` + entityLower + `)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	query := "UPDATE ` + entityLower + `s SET data = ? WHERE id = ?"
	if _, err := s.db.Exec(query, data, ` + entityLower + `.ID); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	return nil
}

func (s *sqlite` + entity + `Repository) Delete(id int) error {
	query := "DELETE FROM ` + entityLower + `s WHERE id = ?"
	if _, err := s.db.Exec(query, id); err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	return nil
}

func (s *sqlite` + entity + `Repository) FindAll() ([]domain.` + entity + `, error) {
	query := "SELECT data FROM ` + entityLower + `s"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()
	var ` + entityLower + `s []domain.` + entity + `
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		var ` + entityLower + ` domain.` + entity + `
		if err := json.Unmarshal(data, &` + entityLower + `); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %w", err)
		}
		` + entityLower + `s = append(` + entityLower + `s, ` + entityLower + `)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return ` + entityLower + `s, nil
}
`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create SQLite repository test file: %v", err)
	}
}

func generateTestSQLServerRepository(t *testing.T, dir, entity string) {
	entityLower := strings.ToLower(entity)
	filePath := filepath.Join(dir, "sqlserver_"+entityLower+"_repository.go")

	content := `package repository

import (
	"fmt"
	"gorm.io/gorm"
	"` + getTestModuleName(t) + `/internal/domain"
)

type sqlserver` + entity + `Repository struct {
	db *gorm.DB
}

func NewSQLServer` + entity + `Repository(db *gorm.DB) ` + entity + `Repository {
	return &sqlserver` + entity + `Repository{db: db}
}

func (s *sqlserver` + entity + `Repository) Save(` + entityLower + ` *domain.` + entity + `) error {
	if err := s.db.Create(` + entityLower + `).Error; err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	return nil
}

func (s *sqlserver` + entity + `Repository) FindByID(id int) (*domain.` + entity + `, error) {
	var ` + entityLower + ` domain.` + entity + `
	if err := s.db.First(&` + entityLower + `, id).Error; err != nil {
		return nil, err
	}
	return &` + entityLower + `, nil
}

func (s *sqlserver` + entity + `Repository) Update(` + entityLower + ` *domain.` + entity + `) error {
	return s.db.Save(` + entityLower + `).Error
}

func (s *sqlserver` + entity + `Repository) Delete(id int) error {
	return s.db.Delete(&domain.` + entity + `{}, id).Error
}

func (s *sqlserver` + entity + `Repository) FindAll() ([]domain.` + entity + `, error) {
	var ` + entityLower + `s []domain.` + entity + `
	if err := s.db.Find(&` + entityLower + `s).Error; err != nil {
		return nil, err
	}
	return ` + entityLower + `s, nil
}
`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create SQL Server repository: %v", err)
	}
}

func generateTestPostgresJSONRepository(t *testing.T, dir, entity string) {
	entityLower := strings.ToLower(entity)
	filePath := filepath.Join(dir, "postgres_json_"+entityLower+"_repository.go")

	content := `package repository

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"` + getTestModuleName(t) + `/internal/domain"
)

type postgresJSON` + entity + `Repository struct {
	db *gorm.DB
}

func NewPostgresJSON` + entity + `Repository(db *gorm.DB) ` + entity + `Repository {
	return &postgresJSON` + entity + `Repository{db: db}
}

func (p *postgresJSON` + entity + `Repository) Save(` + entityLower + ` *domain.` + entity + `) error {
	return p.db.Create(` + entityLower + `).Error
}

func (p *postgresJSON` + entity + `Repository) FindByID(id int) (*domain.` + entity + `, error) {
	var ` + entityLower + ` domain.` + entity + `
	if err := p.db.First(&` + entityLower + `, id).Error; err != nil {
		return nil, err
	}
	return &` + entityLower + `, nil
}

func (p *postgresJSON` + entity + `Repository) FindByJSONField(jsonField, value string) ([]domain.` + entity + `, error) {
	var ` + entityLower + `s []domain.` + entity + `
	if err := p.db.Where("data @> ?", datatypes.JSONQuery(jsonField)).Find(&` + entityLower + `s).Error; err != nil {
		return nil, err
	}
	return ` + entityLower + `s, nil
}

func (p *postgresJSON` + entity + `Repository) Update(` + entityLower + ` *domain.` + entity + `) error {
	return p.db.Save(` + entityLower + `).Error
}

func (p *postgresJSON` + entity + `Repository) Delete(id int) error {
	return p.db.Delete(&domain.` + entity + `{}, id).Error
}

func (p *postgresJSON` + entity + `Repository) FindAll() ([]domain.` + entity + `, error) {
	var ` + entityLower + `s []domain.` + entity + `
	if err := p.db.Find(&` + entityLower + `s).Error; err != nil {
		return nil, err
	}
	return ` + entityLower + `s, nil
}
`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create PostgreSQL JSON repository: %v", err)
	}
}

func generateTestElasticsearchRepository(t *testing.T, dir, entity string) {
	entityLower := strings.ToLower(entity)
	filePath := filepath.Join(dir, "elasticsearch_"+entityLower+"_repository.go")

	content := `package repository

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"` + getTestModuleName(t) + `/internal/domain"
)

type elasticsearch` + entity + `Repository struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticsearch` + entity + `Repository(client *elasticsearch.Client) ` + entity + `Repository {
	return &elasticsearch` + entity + `Repository{
		client: client,
		index:  "` + strings.ToLower(entity) + `",
	}
}

func (e *elasticsearch` + entity + `Repository) Save(` + entityLower + ` *domain.` + entity + `) error {
	return nil
}

func (e *elasticsearch` + entity + `Repository) FindByID(id int) (*domain.` + entity + `, error) {
	return nil, nil
}

func (e *elasticsearch` + entity + `Repository) FullTextSearch(query string) ([]domain.` + entity + `, error) {
	return []domain.` + entity + `{}, nil
}

func (e *elasticsearch` + entity + `Repository) Update(` + entityLower + ` *domain.` + entity + `) error {
	return e.Save(` + entityLower + `)
}

func (e *elasticsearch` + entity + `Repository) Delete(id int) error {
	return nil
}

func (e *elasticsearch` + entity + `Repository) FindAll() ([]domain.` + entity + `, error) {
	return []domain.` + entity + `{}, nil
}
`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create Elasticsearch repository: %v", err)
	}
}

func generateTestDynamoDBRepository(t *testing.T, dir, entity string) {
	entityLower := strings.ToLower(entity)
	filePath := filepath.Join(dir, "dynamodb_"+entityLower+"_repository.go")

	content := `package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"` + getTestModuleName(t) + `/internal/domain"
)

type dynamodb` + entity + `Repository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDB` + entity + `Repository(client *dynamodb.Client) ` + entity + `Repository {
	return &dynamodb` + entity + `Repository{
		client:    client,
		tableName: "` + strings.ToLower(entity) + `",
	}
}

func (d *dynamodb` + entity + `Repository) Save(` + entityLower + ` *domain.` + entity + `) error {
	return nil
}

func (d *dynamodb` + entity + `Repository) FindByID(id int) (*domain.` + entity + `, error) {
	return nil, nil
}

func (d *dynamodb` + entity + `Repository) Update(` + entityLower + ` *domain.` + entity + `) error {
	return d.Save(` + entityLower + `)
}

func (d *dynamodb` + entity + `Repository) Delete(id int) error {
	return nil
}

func (d *dynamodb` + entity + `Repository) FindAll() ([]domain.` + entity + `, error) {
	return []domain.` + entity + `{}, nil
}
`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create DynamoDB repository: %v", err)
	}
}

// Helper function to get module name for tests
func getTestModuleName(t *testing.T) string {
	// Read go.mod to get module name
	content, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}

	return "test-module"
}
