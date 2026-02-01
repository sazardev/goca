package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sazardev/goca/internal/testing/framework"
)

// TestInitSQLiteDriverFix verifica específicamente el fix del issue #31
// Este test asegura que SQLite NO genera dependencias de PostgreSQL
func TestInitSQLiteDriverFix(t *testing.T) {
	tc := framework.NewTestContext(t)
	defer tc.Cleanup()

	projectName := "test-sqlite-fix"
	projectPath := filepath.Join(tc.TempDir, projectName)

	// Ejecutar init con SQLite
	_, err := tc.RunCommand("init", projectName, "--module", "github.com/test/sqlite-fix", "--database", "sqlite")
	if err != nil {
		t.Fatalf("Error ejecutando goca init: %v", err)
	}

	// Leer go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Error leyendo go.mod: %v", err)
	}

	goModStr := string(goModContent)

	// VERIFICACIÓN CRÍTICA: NO debe contener driver de PostgreSQL
	if strings.Contains(goModStr, "gorm.io/driver/postgres") {
		t.Errorf("❌ BUG #31 NO RESUELTO: go.mod contiene driver de PostgreSQL cuando se especificó SQLite!\n\nContenido de go.mod:\n%s", goModStr)
	}

	// DEBE contener driver de SQLite
	if !strings.Contains(goModStr, "gorm.io/driver/sqlite") {
		t.Errorf("❌ go.mod NO contiene el driver de SQLite esperado!\n\nContenido de go.mod:\n%s", goModStr)
	}

	// Leer main.go
	mainGoPath := filepath.Join(projectPath, "cmd", "server", "main.go")
	mainGoContent, err := os.ReadFile(mainGoPath)
	if err != nil {
		t.Fatalf("Error leyendo main.go: %v", err)
	}

	mainGoStr := string(mainGoContent)

	// VERIFICACIÓN CRÍTICA: NO debe importar driver de PostgreSQL
	if strings.Contains(mainGoStr, `"gorm.io/driver/postgres"`) {
		t.Errorf("❌ BUG #31 NO RESUELTO: main.go importa driver de PostgreSQL cuando se especificó SQLite!")
	}

	// DEBE importar driver de SQLite
	if !strings.Contains(mainGoStr, `"gorm.io/driver/sqlite"`) {
		t.Errorf("❌ main.go NO importa el driver de SQLite esperado!")
	}

	// DEBE usar sqlite.Open() en vez de postgres.Open()
	if strings.Contains(mainGoStr, "postgres.Open(") {
		t.Errorf("❌ BUG #31 NO RESUELTO: main.go usa postgres.Open() cuando se especificó SQLite!")
	}

	if !strings.Contains(mainGoStr, "sqlite.Open(") {
		t.Errorf("❌ main.go NO usa sqlite.Open() como se espera!")
	}

	t.Logf("✅ Issue #31 RESUELTO: SQLite se configura correctamente sin dependencias de PostgreSQL")
}

// TestInitMySQLDriverFix verifica que MySQL también funciona correctamente
func TestInitMySQLDriverFix(t *testing.T) {
	tc := framework.NewTestContext(t)
	defer tc.Cleanup()

	projectName := "test-mysql-fix"
	projectPath := filepath.Join(tc.TempDir, projectName)

	// Ejecutar init con MySQL
	_, err := tc.RunCommand("init", projectName, "--module", "github.com/test/mysql-fix", "--database", "mysql")
	if err != nil {
		t.Fatalf("Error ejecutando goca init: %v", err)
	}

	// Leer go.mod
	goModContent, err := os.ReadFile(filepath.Join(projectPath, "go.mod"))
	if err != nil {
		t.Fatalf("Error leyendo go.mod: %v", err)
	}

	goModStr := string(goModContent)

	// Verificaciones
	if !strings.Contains(goModStr, "gorm.io/driver/mysql") {
		t.Errorf("go.mod NO contiene driver de MySQL")
	}

	if strings.Contains(goModStr, "gorm.io/driver/postgres") {
		t.Errorf("go.mod contiene driver de PostgreSQL (no debería)")
	}

	// Leer main.go
	mainGoContent, err := os.ReadFile(filepath.Join(projectPath, "cmd", "server", "main.go"))
	if err != nil {
		t.Fatalf("Error leyendo main.go: %v", err)
	}

	mainGoStr := string(mainGoContent)

	if !strings.Contains(mainGoStr, `"gorm.io/driver/mysql"`) {
		t.Errorf("main.go NO importa driver de MySQL")
	}

	if !strings.Contains(mainGoStr, "mysql.Open(") {
		t.Errorf("main.go NO usa mysql.Open()")
	}

	t.Logf("✅ MySQL se configura correctamente")
}

// TestInitPostgreSQLStillWorks verifica que el fix no rompió PostgreSQL
func TestInitPostgreSQLStillWorks(t *testing.T) {
	tc := framework.NewTestContext(t)
	defer tc.Cleanup()

	projectName := "test-postgres-regression"
	projectPath := filepath.Join(tc.TempDir, projectName)

	// Ejecutar init con PostgreSQL (default behavior)
	_, err := tc.RunCommand("init", projectName, "--module", "github.com/test/postgres-regression", "--database", "postgres")
	if err != nil {
		t.Fatalf("Error ejecutando goca init: %v", err)
	}

	// Leer go.mod
	goModContent, err := os.ReadFile(filepath.Join(projectPath, "go.mod"))
	if err != nil {
		t.Fatalf("Error leyendo go.mod: %v", err)
	}

	goModStr := string(goModContent)

	// Verificaciones - PostgreSQL debe seguir funcionando
	if !strings.Contains(goModStr, "gorm.io/driver/postgres") {
		t.Errorf("❌ REGRESIÓN: go.mod NO contiene driver de PostgreSQL")
	}

	// Leer main.go
	mainGoContent, err := os.ReadFile(filepath.Join(projectPath, "cmd", "server", "main.go"))
	if err != nil {
		t.Fatalf("Error leyendo main.go: %v", err)
	}

	mainGoStr := string(mainGoContent)

	if !strings.Contains(mainGoStr, `"gorm.io/driver/postgres"`) {
		t.Errorf("❌ REGRESIÓN: main.go NO importa driver de PostgreSQL")
	}

	if !strings.Contains(mainGoStr, "postgres.Open(") {
		t.Errorf("❌ REGRESIÓN: main.go NO usa postgres.Open()")
	}

	t.Logf("✅ PostgreSQL sigue funcionando correctamente (sin regresión)")
}

// TestInitDefaultDatabase tests that the default database is SQLite
func TestInitDefaultDatabase(t *testing.T) {
	tc := framework.NewTestContext(t)
	defer tc.Cleanup()

	projectName := "test-default-db"
	projectPath := filepath.Join(tc.TempDir, projectName)

	// Execute init command without database flag (should use default)
	_, err := tc.RunCommand("init", projectName, "--module", "testmodule")
	if err != nil {
		t.Fatalf("Failed to execute init command: %v", err)
	}

	// Verify project was created
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Fatalf("Project directory was not created: %s", projectPath)
	}

	// Read .goca.yaml to verify default database
	configPath := filepath.Join(projectPath, ".goca.yaml")
	configContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	// Check that database is sqlite
	if !strings.Contains(string(configContent), "type: \"sqlite\"") {
		t.Errorf("Default database is not sqlite. Config content:\n%s", string(configContent))
	}

	// Verify go.mod has sqlite driver
	goModPath := filepath.Join(projectPath, "go.mod")
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	if !strings.Contains(string(goModContent), "gorm.io/driver/sqlite") {
		t.Errorf("go.mod does not contain sqlite driver")
	}

	t.Logf("✓ Default database is SQLite as expected")
}

// TestInitMongoDBNoGorm tests that MongoDB projects don't use GORM
func TestInitMongoDBNoGorm(t *testing.T) {
	tc := framework.NewTestContext(t)
	defer tc.Cleanup()

	projectName := "test-mongodb"
	projectPath := filepath.Join(tc.TempDir, projectName)

	// Execute init command with MongoDB
	_, err := tc.RunCommand("init", projectName, "--module", "testmodule", "--database", "mongodb")
	if err != nil {
		t.Fatalf("Failed to execute init command: %v", err)
	}

	// Verify project was created
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Fatalf("Project directory was not created: %s", projectPath)
	}

	// Read main.go
	mainGoPath := filepath.Join(projectPath, "cmd", "server", "main.go")
	mainGoContent, err := os.ReadFile(mainGoPath)
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}

	mainGoStr := string(mainGoContent)

	// Verify MongoDB driver is imported
	if !strings.Contains(mainGoStr, "go.mongodb.org/mongo-driver/mongo") {
		t.Errorf("main.go does not import MongoDB driver")
	}

	// Verify GORM is NOT imported
	if strings.Contains(mainGoStr, "gorm.io/gorm") {
		t.Errorf("main.go incorrectly imports GORM for MongoDB")
	}

	// Verify mongo client variable is used (not gorm.DB)
	if strings.Contains(mainGoStr, "*gorm.DB") {
		t.Errorf("main.go incorrectly uses gorm.DB type for MongoDB")
	}

	if !strings.Contains(mainGoStr, "*mongo.Client") {
		t.Errorf("main.go does not use mongo.Client type")
	}

	// Verify go.mod has mongo driver and NOT gorm
	goModPath := filepath.Join(projectPath, "go.mod")
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	goModStr := string(goModContent)

	if !strings.Contains(goModStr, "go.mongodb.org/mongo-driver") {
		t.Errorf("go.mod does not contain mongo-driver")
	}

	if strings.Contains(goModStr, "gorm.io/gorm") {
		t.Errorf("go.mod incorrectly contains gorm for MongoDB project")
	}

	t.Logf("✓ MongoDB project correctly uses mongo-driver without GORM")
}
