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

		// Generate integration tests
		if err := generateIntegrationTests(entityName, integrationTestDatabase, integrationTestFixtures, integrationTestContainer); err != nil {
			validator.errorHandler.HandleError(err, "test-integration")
			return
		}

		fmt.Printf("\n✅ Integration tests generated successfully for '%s'\n", entityName)
		fmt.Println("\nGenerated files:")
		fmt.Printf("   - internal/testing/integration/%s_integration_test.go\n", strings.ToLower(entityName))
		if integrationTestFixtures {
			fmt.Printf("   - internal/testing/integration/fixtures/%s_fixtures.go\n", strings.ToLower(entityName))
		}
		fmt.Println("\nRun tests:")
		fmt.Printf("   go test ./internal/testing/integration -v\n")
	},
}

func init() {
	rootCmd.AddCommand(testIntegrationCmd)

	testIntegrationCmd.Flags().StringVar(&integrationTestDatabase, "database", "postgres", "Database type for integration tests")
	testIntegrationCmd.Flags().BoolVar(&integrationTestFixtures, "fixtures", true, "Generate test fixtures")
	testIntegrationCmd.Flags().BoolVar(&integrationTestContainer, "container", false, "Use test containers for database")
}

// generateIntegrationTests generates integration test files
func generateIntegrationTests(entityName, database string, withFixtures, withContainer bool) error {
	// Create integration test directory
	integrationDir := filepath.Join("internal", "testing", "integration")
	if err := os.MkdirAll(integrationDir, 0755); err != nil {
		return fmt.Errorf("failed to create integration directory: %v", err)
	}

	// Generate main integration test file
	testFile := filepath.Join(integrationDir, strings.ToLower(entityName)+"_integration_test.go")
	content := generateIntegrationTestContent(entityName, database, withContainer)
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write integration test file: %v", err)
	}

	// Generate fixtures if requested
	if withFixtures {
		fixturesDir := filepath.Join(integrationDir, "fixtures")
		if err := os.MkdirAll(fixturesDir, 0755); err != nil {
			return fmt.Errorf("failed to create fixtures directory: %v", err)
		}

		fixtureFile := filepath.Join(fixturesDir, strings.ToLower(entityName)+"_fixtures.go")
		fixtureContent := generateFixtureContent(entityName)
		if err := os.WriteFile(fixtureFile, []byte(fixtureContent), 0644); err != nil {
			return fmt.Errorf("failed to write fixture file: %v", err)
		}
	}

	// Generate database helpers if they don't exist
	helpersFile := filepath.Join(integrationDir, "helpers.go")
	if _, err := os.Stat(helpersFile); os.IsNotExist(err) {
		helpersContent := generateHelpersContent(database, withContainer)
		if err := os.WriteFile(helpersFile, []byte(helpersContent), 0644); err != nil {
			return fmt.Errorf("failed to write helpers file: %v", err)
		}
	}

	return nil
}

// generateIntegrationTestContent generates the main integration test file content
func generateIntegrationTestContent(entityName, database string, withContainer bool) string {
	lowerEntity := strings.ToLower(entityName)

	return fmt.Sprintf(`package integration

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
}

// generateFixtureContent generates test fixtures
func generateFixtureContent(entityName string) string {
	lowerEntity := strings.ToLower(entityName)

	return fmt.Sprintf(`package fixtures

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
}

// generateHelpersContent generates database helpers for integration tests
func generateHelpersContent(database string, withContainer bool) string {
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

	return fmt.Sprintf(`package integration

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
}
