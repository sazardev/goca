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
