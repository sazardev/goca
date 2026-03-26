package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	integrationTestDatabase  string
	integrationTestFields    string
	integrationTestFixtures  bool
	integrationTestContainer bool
)

// testIntegrationCmd represents the test-integration command
var testIntegrationCmd = &cobra.Command{
	Use:   "test-integration [entity]",
	Short: "Generate integration tests for a feature",
	Long: `Generate comprehensive integration tests for a feature that verify
the interaction between different layers of the Clean Architecture.

Integration tests will verify:
- Use case and repository interaction
- Handler and use case interaction
- Database integration
- CRUD operations end-to-end

Examples:
  goca test-integration User
  goca test-integration Product --database postgres
  goca test-integration Order --fixtures --container`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entityName := args[0]

		// Validate entity name
		validator := NewCommandValidator()
		if err := validator.fieldValidator.ValidateEntityName(entityName); err != nil {
			validator.errorHandler.HandleError(err, "test-integration")
			return
		}

		// Check if we're in a Goca project
		if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
			validator.errorHandler.HandleError(fmt.Errorf("not in a Go project directory. Run 'goca init' first"), "test-integration")
			return
		}

		// Initialize safety manager
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		backup, _ := cmd.Flags().GetBool("backup")
		sm := NewSafetyManager(dryRun, force, backup)

		if dryRun {
			ui.DryRun("Previewing changes without creating files")
		}

		// Parse fields if provided
		var fields []Field
		if integrationTestFields != "" {
			fields = parseFields(integrationTestFields)
		}

		// Generate integration tests
		if err := generateIntegrationTests(entityName, integrationTestDatabase, integrationTestFixtures, integrationTestContainer, fields, sm); err != nil {
			validator.errorHandler.HandleError(err, "test-integration")
			return
		}

		if dryRun {
			sm.PrintSummary()
			return
		}

		ui.Success(fmt.Sprintf("Integration tests generated successfully for '%s'", entityName))
		ui.Blank()
		ui.Section("Generated files")
		ui.FileCreated(fmt.Sprintf("internal/testing/integration/%s_integration_test.go", strings.ToLower(entityName)))
		if integrationTestFixtures {
			ui.FileCreated(fmt.Sprintf("internal/testing/integration/fixtures/%s_fixtures.go", strings.ToLower(entityName)))
		}
		ui.Blank()
		ui.Section("Run tests")
		ui.Dim("   go test ./internal/testing/integration -v")
	},
}

func init() {
	rootCmd.AddCommand(testIntegrationCmd)

	testIntegrationCmd.Flags().StringVar(&integrationTestDatabase, "database", "postgres", "Database type for integration tests")
	testIntegrationCmd.Flags().StringVar(&integrationTestFields, "fields", "", "Entity fields (e.g. \"Name:string,Email:string,Age:int\")")
	testIntegrationCmd.Flags().BoolVar(&integrationTestFixtures, "fixtures", true, "Generate test fixtures")
	testIntegrationCmd.Flags().BoolVar(&integrationTestContainer, "container", false, "Use test containers for database")
	testIntegrationCmd.Flags().Bool("dry-run", false, "Preview changes without creating files")
	testIntegrationCmd.Flags().Bool("force", false, "Overwrite existing files without confirmation")
	testIntegrationCmd.Flags().Bool("backup", false, "Create backup of existing files before overwriting")
}

// generateIntegrationTests generates integration test files
func generateIntegrationTests(entityName, database string, withFixtures, withContainer bool, fields []Field, sm ...*SafetyManager) error {
	// Create integration test directory
	integrationDir := filepath.Join("internal", "testing", "integration")
	if err := os.MkdirAll(integrationDir, 0755); err != nil {
		return fmt.Errorf("failed to create integration directory: %v", err)
	}

	// Generate main integration test file
	testFile := filepath.Join(integrationDir, strings.ToLower(entityName)+"_integration_test.go")
	content := generateIntegrationTestContent(entityName, database, withContainer, fields)
	if err := writeFile(testFile, content, sm...); err != nil {
		return fmt.Errorf("failed to write integration test file: %v", err)
	}

	// Generate fixtures if requested
	if withFixtures {
		fixturesDir := filepath.Join(integrationDir, "fixtures")
		if err := os.MkdirAll(fixturesDir, 0755); err != nil {
			return fmt.Errorf("failed to create fixtures directory: %v", err)
		}

		fixtureFile := filepath.Join(fixturesDir, strings.ToLower(entityName)+"_fixtures.go")
		fixtureContent := generateFixtureContent(entityName, fields)
		if err := writeFile(fixtureFile, fixtureContent, sm...); err != nil {
			return fmt.Errorf("failed to write fixture file: %v", err)
		}
	}

	// Generate database helpers if they don't exist
	helpersFile := filepath.Join(integrationDir, "helpers.go")
	if _, err := os.Stat(helpersFile); os.IsNotExist(err) {
		helpersContent := generateHelpersContent(database, withContainer, entityName)
		if err := writeFile(helpersFile, helpersContent, sm...); err != nil {
			return fmt.Errorf("failed to write helpers file: %v", err)
		}
	}

	return nil
}

// generateIntegrationTestContent generates the main integration test file content
func generateIntegrationTestContent(entityName, database string, withContainer bool, fields []Field) string {
	lowerEntity := strings.ToLower(entityName)

	content := fmt.Sprintf(`package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/repository"
	"github.com/sazardev/goca/internal/usecase"
)

// Test%[1]sIntegration tests the complete %[1]s feature integration
func Test%[1]sIntegration(t *testing.T) {
	// Setup test database
	db := setupTestDatabase(t, "%[2]s")
	defer cleanupTestDatabase(t, db)

	// Initialize dependencies
	repo := repository.NewPostgres%[1]sRepository(db)
	service := usecase.New%[1]sService(repo)

	ctx := context.Background()

	t.Run("CreateAndRetrieve%[1]s", func(t *testing.T) {
		// Create test data
		input := usecase.Create%[1]sInput{
			// TODO: Add fields based on entity structure
		}

		// Create %[3]s
		output, err := service.Create%[1]s(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, output)
		assert.NotZero(t, output.ID)

		// Retrieve %[3]s
		retrieved, err := service.Get%[1]sByID(ctx, output.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, output.ID, retrieved.ID)
	})

	t.Run("Update%[1]s", func(t *testing.T) {
		// Create initial %[3]s
		input := usecase.Create%[1]sInput{
			// TODO: Add fields
		}
		created, err := service.Create%[1]s(ctx, input)
		require.NoError(t, err)

		// Update %[3]s
		updateInput := usecase.Update%[1]sInput{
			ID: created.ID,
			// TODO: Add updated fields
		}
		updated, err := service.Update%[1]s(ctx, updateInput)
		require.NoError(t, err)
		assert.Equal(t, created.ID, updated.ID)
		// TODO: Add field assertions
	})

	t.Run("Delete%[1]s", func(t *testing.T) {
		// Create %[3]s to delete
		input := usecase.Create%[1]sInput{
			// TODO: Add fields
		}
		created, err := service.Create%[1]s(ctx, input)
		require.NoError(t, err)

		// Delete %[3]s
		err = service.Delete%[1]s(ctx, created.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = service.Get%[1]sByID(ctx, created.ID)
		assert.Error(t, err)
	})

	t.Run("List%[1]s", func(t *testing.T) {
		// Create multiple %[3]ss
		for i := 0; i < 3; i++ {
			input := usecase.Create%[1]sInput{
				// TODO: Add fields with variation
			}
			_, err := service.Create%[1]s(ctx, input)
			require.NoError(t, err)
		}

		// List all %[3]ss
		list, err := service.List%[1]s(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 3)
	})

	t.Run("Transaction%[1]s", func(t *testing.T) {
		// Test transaction rollback
		input := usecase.Create%[1]sInput{
			// TODO: Add invalid data to trigger error
		}

		// This should fail and rollback
		_, err := service.Create%[1]s(ctx, input)
		assert.Error(t, err)

		// Verify no %[3]s was created
		list, err := service.List%[1]s(ctx)
		require.NoError(t, err)
		// TODO: Verify count hasn't increased
	})
}

// Test%[1]sRepositoryIntegration tests repository layer directly
func Test%[1]sRepositoryIntegration(t *testing.T) {
	db := setupTestDatabase(t, "%[2]s")
	defer cleanupTestDatabase(t, db)

	repo := repository.NewPostgres%[1]sRepository(db)
	ctx := context.Background()

	t.Run("SaveAndFindByID", func(t *testing.T) {
		%[3]s := &domain.%[1]s{
			// TODO: Add fields
		}

		// Save
		err := repo.Save(ctx, %[3]s)
		require.NoError(t, err)
		assert.NotZero(t, %[3]s.ID)

		// Find by ID
		found, err := repo.FindByID(ctx, %[3]s.ID)
		require.NoError(t, err)
		assert.Equal(t, %[3]s.ID, found.ID)
	})

	t.Run("FindAll", func(t *testing.T) {
		// Create test entities
		for i := 0; i < 5; i++ {
			%[3]s := &domain.%[1]s{
				// TODO: Add fields
			}
			err := repo.Save(ctx, %[3]s)
			require.NoError(t, err)
		}

		// Find all
		all, err := repo.FindAll(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(all), 5)
	})

	t.Run("DeleteByID", func(t *testing.T) {
		%[3]s := &domain.%[1]s{
			// TODO: Add fields
		}
		err := repo.Save(ctx, %[3]s)
		require.NoError(t, err)

		// Delete
		err = repo.Delete(ctx, %[3]s.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.FindByID(ctx, %[3]s.ID)
		assert.Error(t, err)
	})
}
`, entityName, database, lowerEntity)
	return replaceIntegrationTestTODOs(content, fields, entityName)
}

// generateFixtureContent generates test fixtures
func generateFixtureContent(entityName string, fields []Field) string {
	lowerEntity := strings.ToLower(entityName)

	content := fmt.Sprintf(`package fixtures

import (
	"github.com/sazardev/goca/internal/domain"
)

// New%[1]sFixture creates a new %[1]s fixture with default values
func New%[1]sFixture() *domain.%[1]s {
	return &domain.%[1]s{
		// TODO: Add default field values for testing
		// Example:
		// Name: "Test %[1]s",
		// Email: "test@example.com",
	}
}

// New%[1]sFixtureWithCustomFields creates a %[1]s fixture with custom values
func New%[1]sFixtureWithCustomFields(fields map[string]interface{}) *domain.%[1]s {
	%[2]s := New%[1]sFixture()

	// Override fields based on provided map
	// TODO: Implement field overrides
	// Example:
	// if name, ok := fields["name"].(string); ok {
	//     %[2]s.Name = name
	// }

	return %[2]s
}

// New%[1]sFixtureList creates multiple %[1]s fixtures
func New%[1]sFixtureList(count int) []*domain.%[1]s {
	fixtures := make([]*domain.%[1]s, count)
	for i := 0; i < count; i++ {
		fixtures[i] = New%[1]sFixture()
		// TODO: Vary fixture data for each instance
		// Example: 
		// fixtures[i].Name = fmt.Sprintf("Test %[1]s %%d", i+1)
	}
	return fixtures
}
`, entityName, lowerEntity)
	return replaceFixtureTODOs(content, fields, entityName)
}

// generateHelpersContent generates database helpers for integration tests
func generateHelpersContent(database string, withContainer bool, entityName string) string {
	containerSetup := ""
	if withContainer {
		containerSetup = `
	// Using test containers
	// TODO: Implement test container setup
	// Example with testcontainers-go:
	// ctx := context.Background()
	// req := testcontainers.ContainerRequest{
	//     Image:        "postgres:15",
	//     ExposedPorts: []string{"5432/tcp"},
	//     Env: map[string]string{
	//         "POSTGRES_USER":     "test",
	//         "POSTGRES_PASSWORD": "test",
	//         "POSTGRES_DB":       "testdb",
	//     },
	//     WaitingFor: wait.ForLog("database system is ready to accept connections"),
	// }
	// container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
	//     ContainerRequest: req,
	//     Started:          true,
	// })
	// if err != nil {
	//     t.Fatalf("failed to start container: %v", err)
	// }`
	}

	content := fmt.Sprintf(`package integration

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq" // PostgreSQL driver
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDatabase initializes a test database connection
func setupTestDatabase(t *testing.T, dbType string) *gorm.DB {
	t.Helper()

	var dsn string
	var dialector gorm.Dialector

	switch dbType {
	case "postgres":
		%[1]s
		// Using in-memory or test database
		dsn = "host=localhost user=test password=test dbname=goca_test port=5432 sslmode=disable"
		dialector = postgres.Open(dsn)
	case "mysql":
		// TODO: Add MySQL setup
		t.Skip("MySQL integration tests not implemented yet")
	case "sqlite":
		// SQLite is perfect for fast integration tests
		dialector = postgres.Open(":memory:")
	default:
		t.Fatalf("unsupported database type: %%s", dbType)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %%v", err)
	}

	// Run migrations
	// TODO: Add auto-migration for test entities
	// Example: db.AutoMigrate(&domain.User{}, &domain.Product{})

	return db
}

// cleanupTestDatabase cleans up the test database
func cleanupTestDatabase(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlDB, err := db.DB()
	if err != nil {
		t.Logf("failed to get SQL DB: %%v", err)
		return
	}

	// Drop all tables
	// TODO: Add table cleanup based on entities
	// Example: db.Migrator().DropTable(&domain.User{}, &domain.Product{})

	if err := sqlDB.Close(); err != nil {
		t.Logf("failed to close database: %%v", err)
	}
}

// truncateTables truncates all tables in the test database
func truncateTables(t *testing.T, db *gorm.DB, tables ...string) {
	t.Helper()

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %%s CASCADE", table)).Error; err != nil {
			t.Logf("failed to truncate table %%s: %%v", table, err)
		}
	}
}

// beginTransaction starts a new transaction for testing
func beginTransaction(t *testing.T, db *gorm.DB) *gorm.DB {
	t.Helper()

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %%v", tx.Error)
	}

	return tx
}

// rollbackTransaction rolls back a transaction
func rollbackTransaction(t *testing.T, tx *gorm.DB) {
	t.Helper()

	if err := tx.Rollback().Error; err != nil {
		t.Logf("failed to rollback transaction: %%v", err)
	}
}

// commitTransaction commits a transaction
func commitTransaction(t *testing.T, tx *gorm.DB) {
	t.Helper()

	if err := tx.Commit().Error; err != nil {
		t.Fatalf("failed to commit transaction: %%v", err)
	}
}

// seedTestData seeds the database with test data
func seedTestData(t *testing.T, db *gorm.DB) {
	t.Helper()

	// TODO: Add seed data for tests
	// Example:
	// users := []*domain.User{
	//     {Name: "Test User 1", Email: "user1@test.com"},
	//     {Name: "Test User 2", Email: "user2@test.com"},
	// }
	// for _, user := range users {
	//     if err := db.Create(user).Error; err != nil {
	//         t.Fatalf("failed to seed user: %%v", err)
	//     }
	// }
}
`, containerSetup)
	return replaceHelperTODOs(content, entityName)
}

// skipTestField returns true for auto-managed fields that should not appear in test input.
func skipTestField(name string) bool {
	switch name {
	case "ID", "CreatedAt", "UpdatedAt", "DeletedAt":
		return true
	}
	return false
}

// testLiteral returns a Go literal string for use in generated test code.
func testLiteral(name, typ, entity string) string {
	lower := strings.ToLower(name)
	switch {
	case lower == "email":
		return fmt.Sprintf(`"test@%s.com"`, strings.ToLower(entity))
	case lower == "name" || lower == "title":
		return fmt.Sprintf(`"Test %s"`, entity)
	case lower == "description":
		return fmt.Sprintf(`"A test %s description"`, strings.ToLower(entity))
	case lower == "status":
		return `"active"`
	case lower == "phone":
		return `"+1234567890"`
	case lower == "address":
		return `"123 Test Street"`
	case lower == "price" || lower == "amount" || lower == "total" || lower == "cost":
		return "9.99"
	case lower == "age" || lower == "quantity" || lower == "count":
		return "1"
	case typ == "string":
		return fmt.Sprintf(`"test_%s"`, strings.ToLower(name))
	case typ == "int" || typ == "int64" || typ == "uint" || typ == "uint64" || typ == "int32" || typ == "uint32":
		return "1"
	case typ == "float64" || typ == "float32":
		return "9.99"
	case typ == "bool":
		return "true"
	case typ == "time.Time":
		return "time.Now()"
	default:
		return fmt.Sprintf(`"test_%s"`, strings.ToLower(name))
	}
}

// updatedTestLiteral returns an updated Go literal for testing update operations.
func updatedTestLiteral(name, typ, entity string) string {
	lower := strings.ToLower(name)
	switch {
	case lower == "email":
		return fmt.Sprintf(`"updated@%s.com"`, strings.ToLower(entity))
	case lower == "name" || lower == "title":
		return fmt.Sprintf(`"Updated %s"`, entity)
	case lower == "description":
		return fmt.Sprintf(`"Updated %s description"`, strings.ToLower(entity))
	case lower == "status":
		return `"inactive"`
	case lower == "phone":
		return `"+0987654321"`
	case lower == "address":
		return `"456 Updated Avenue"`
	case lower == "price" || lower == "amount" || lower == "total" || lower == "cost":
		return "19.99"
	case lower == "age" || lower == "quantity" || lower == "count":
		return "2"
	case typ == "string":
		return fmt.Sprintf(`"updated_%s"`, strings.ToLower(name))
	case typ == "int" || typ == "int64" || typ == "uint" || typ == "uint64" || typ == "int32" || typ == "uint32":
		return "2"
	case typ == "float64" || typ == "float32":
		return "19.99"
	case typ == "bool":
		return "false"
	case typ == "time.Time":
		return "time.Now()"
	default:
		return fmt.Sprintf(`"updated_%s"`, strings.ToLower(name))
	}
}

// needsFmtImport checks whether any string fields exist that require fmt in generated tests.
func needsFmtImport(fields []Field) bool {
	for _, f := range fields {
		if skipTestField(f.Name) {
			continue
		}
		if f.Type == "string" {
			return true
		}
	}
	return false
}

// buildTestFieldInit generates Go struct literal field initializers for test code.
func buildTestFieldInit(fields []Field, entity, indent string) string {
	var lines []string
	for _, f := range fields {
		if skipTestField(f.Name) {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s%s: %s,", indent, f.Name, testLiteral(f.Name, f.Type, entity)))
	}
	return strings.Join(lines, "\n")
}

// buildTestFieldInitUpdated generates updated field initializers for test update operations.
func buildTestFieldInitUpdated(fields []Field, entity, indent string) string {
	var lines []string
	for _, f := range fields {
		if skipTestField(f.Name) {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s%s: %s,", indent, f.Name, updatedTestLiteral(f.Name, f.Type, entity)))
	}
	return strings.Join(lines, "\n")
}

// buildTestFieldAssertions generates assert.Equal calls for each field.
func buildTestFieldAssertions(fields []Field, indent string) string {
	var lines []string
	for _, f := range fields {
		if skipTestField(f.Name) {
			continue
		}
		lines = append(lines, fmt.Sprintf("%sassert.Equal(t, updateInput.%s, updated.%s)", indent, f.Name, f.Name))
	}
	return strings.Join(lines, "\n")
}

// buildTestFieldInitVaried generates field initializers with variation for loop-based test creation.
func buildTestFieldInitVaried(fields []Field, entity, indent string, useFmt bool) string {
	var lines []string
	for _, f := range fields {
		if skipTestField(f.Name) {
			continue
		}
		lower := strings.ToLower(f.Name)
		switch {
		case f.Type == "string" && useFmt:
			if lower == "email" {
				lines = append(lines, fmt.Sprintf(`%s%s: fmt.Sprintf("test%%d@%s.com", i+1),`, indent, f.Name, strings.ToLower(entity)))
			} else {
				lines = append(lines, fmt.Sprintf(`%s%s: fmt.Sprintf("Test %s %%d", i+1),`, indent, f.Name, f.Name))
			}
		case f.Type == "int" || f.Type == "int64" || f.Type == "uint" || f.Type == "uint64":
			lines = append(lines, fmt.Sprintf(`%s%s: i + 1,`, indent, f.Name))
		case f.Type == "float64" || f.Type == "float32":
			lines = append(lines, fmt.Sprintf(`%s%s: float64(i+1) * 9.99,`, indent, f.Name))
		default:
			lines = append(lines, fmt.Sprintf(`%s%s: %s,`, indent, f.Name, testLiteral(f.Name, f.Type, entity)))
		}
	}
	return strings.Join(lines, "\n")
}

// buildFixtureOverrides generates field override code for custom fixture creation.
func buildFixtureOverrides(fields []Field, lowerEntity, indent string) string {
	var lines []string
	for _, f := range fields {
		if skipTestField(f.Name) {
			continue
		}
		lowerName := strings.ToLower(f.Name)
		lines = append(lines, fmt.Sprintf(`%sif v, ok := fields["%s"].(%s); ok {`, indent, lowerName, f.Type))
		lines = append(lines, fmt.Sprintf("%s\t%s.%s = v", indent, lowerEntity, f.Name))
		lines = append(lines, fmt.Sprintf("%s}", indent))
	}
	return strings.Join(lines, "\n")
}

// buildFixtureVariation generates variation code for fixture lists.
func buildFixtureVariation(fields []Field, entity, indent string, useFmt bool) string {
	var lines []string
	for _, f := range fields {
		if skipTestField(f.Name) {
			continue
		}
		switch {
		case f.Type == "string" && useFmt:
			lower := strings.ToLower(f.Name)
			if lower == "email" {
				lines = append(lines, fmt.Sprintf(`%sfixtures[i].%s = fmt.Sprintf("test%%d@%s.com", i+1)`, indent, f.Name, strings.ToLower(entity)))
			} else {
				lines = append(lines, fmt.Sprintf(`%sfixtures[i].%s = fmt.Sprintf("Test %s %%d", i+1)`, indent, f.Name, f.Name))
			}
		case f.Type == "int" || f.Type == "int64" || f.Type == "uint" || f.Type == "uint64":
			lines = append(lines, fmt.Sprintf(`%sfixtures[i].%s = i + 1`, indent, f.Name))
		case f.Type == "float64" || f.Type == "float32":
			lines = append(lines, fmt.Sprintf(`%sfixtures[i].%s = float64(i+1) * 9.99`, indent, f.Name))
		}
	}
	return strings.Join(lines, "\n")
}

// replaceIntegrationTestTODOs replaces TODO placeholders with field-aware code in integration test content.
func replaceIntegrationTestTODOs(content string, fields []Field, entityName string) string {
	if len(fields) == 0 {
		return content
	}

	useFmt := needsFmtImport(fields)
	if useFmt {
		content = strings.Replace(content, "\t\"context\"\n", "\t\"context\"\n\t\"fmt\"\n", 1)
	}

	createInit := buildTestFieldInit(fields, entityName, "\t\t\t")
	updateInit := buildTestFieldInitUpdated(fields, entityName, "\t\t\t")
	assertions := buildTestFieldAssertions(fields, "\t\t")
	variedInit := buildTestFieldInitVaried(fields, entityName, "\t\t\t\t", useFmt)

	// Replace in order: most specific patterns first to avoid partial matches
	content = strings.ReplaceAll(content, "\t\t\t// TODO: Add fields based on entity structure", createInit)
	content = strings.ReplaceAll(content, "\t\t\t\t// TODO: Add fields with variation", variedInit)
	content = strings.ReplaceAll(content, "\t\t\t// TODO: Add updated fields", updateInit)
	content = strings.ReplaceAll(content, "\t\t// TODO: Add field assertions", assertions)
	content = strings.ReplaceAll(content, "\t\t\t// TODO: Add invalid data to trigger error", "\t\t\t// Empty input to trigger validation error")
	content = strings.ReplaceAll(content, "\t\t// TODO: Verify count hasn't increased", "")
	content = strings.ReplaceAll(content, "\t\t\t// TODO: Add fields", createInit)

	return content
}

// replaceFixtureTODOs replaces TODO placeholders with field-aware code in fixture content.
func replaceFixtureTODOs(content string, fields []Field, entityName string) string {
	if len(fields) == 0 {
		return content
	}

	lowerEntity := strings.ToLower(entityName)
	useFmt := needsFmtImport(fields)

	if useFmt {
		content = strings.Replace(content,
			"\t\"github.com/sazardev/goca/internal/domain\"",
			"\t\"fmt\"\n\n\t\"github.com/sazardev/goca/internal/domain\"", 1)
	}

	// Default field values
	defaults := buildTestFieldInit(fields, entityName, "\t\t")
	content = strings.Replace(content,
		"\t\t// TODO: Add default field values for testing\n\t\t// Example:\n\t\t// Name: \"Test "+entityName+"\",\n\t\t// Email: \"test@example.com\",",
		defaults, 1)

	// Field overrides
	overrides := buildFixtureOverrides(fields, lowerEntity, "\t")
	content = strings.Replace(content,
		"\t// TODO: Implement field overrides\n\t// Example:\n\t// if name, ok := fields[\"name\"].(string); ok {\n\t//     "+lowerEntity+".Name = name\n\t// }",
		overrides, 1)

	// Varied fixture data
	variation := buildFixtureVariation(fields, entityName, "\t\t", useFmt)
	content = strings.Replace(content,
		"\t\t// TODO: Vary fixture data for each instance\n\t\t// Example: \n\t\t// fixtures[i].Name = fmt.Sprintf(\"Test "+entityName+" %d\", i+1)",
		variation, 1)

	return content
}

// replaceHelperTODOs replaces TODO placeholders with entity-specific guidance in helper content.
func replaceHelperTODOs(content, entityName string) string {
	content = strings.Replace(content,
		"// TODO: Add auto-migration for test entities",
		"// Auto-migrate entity (requires domain import): db.AutoMigrate(&domain."+entityName+"{})", 1)
	content = strings.Replace(content,
		"// TODO: Add table cleanup based on entities",
		"// Drop entity table (requires domain import): db.Migrator().DropTable(&domain."+entityName+"{})", 1)
	return content
}
