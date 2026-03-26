package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper to set up UI and return cleanup
func setupTestUI(t *testing.T) (*bytes.Buffer, func()) {
	t.Helper()
	buf := &bytes.Buffer{}
	testUI := NewUIRenderer(buf, true, 2)
	oldUI := ui
	ui = testUI
	return buf, func() { ui = oldUI }
}

// helper to set up UI with discard writer and return cleanup
func setupDiscardUI(t *testing.T) func() {
	t.Helper()
	testUI := NewUIRenderer(io.Discard, true, 0)
	oldUI := ui
	ui = testUI
	return func() { ui = oldUI }
}

// helper to set up a project directory with go.mod
func setupProjectDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	gomod := `module github.com/test/myproject

go 1.21

require (
	github.com/gorilla/mux v1.8.0
	gorm.io/gorm v1.25.5
)
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0644))
	return dir
}

// --- Init project file generators (DryRun) ---

func TestCreateGoMod_DryRun(t *testing.T) {
	// Not parallel: os.Chdir
	cleanup := setupDiscardUI(t)
	defer cleanup()

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)
	projDir := filepath.Join(dir, "testproj")
	require.NoError(t, os.MkdirAll(projDir, 0755))

	sm := NewSafetyManager(true, false, false)
	tests := []struct {
		name     string
		database string
		auth     bool
	}{
		{"postgres", DBPostgres, false},
		{"mysql", DBMySQL, false},
		{"sqlite", DBSQLite, true},
		{"mongodb", DBMongoDB, false},
		{"dynamodb", DBDynamoDB, false},
		{"elasticsearch", DBElasticsearch, false},
		{"postgres_json", DBPostgresJSON, false},
		{"sqlserver", DBSQLServer, false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			sm2 := NewSafetyManager(true, false, false)
			createGoMod("testproj", "github.com/test/proj", tc.database, tc.auth, sm2)
			assert.NotEmpty(t, sm2.GetPendingFiles())
		})
	}
	_ = sm
}

func TestCreateMainGo_DryRun(t *testing.T) {
	// Not parallel: os.Chdir
	cleanup := setupDiscardUI(t)
	defer cleanup()

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)
	projDir := filepath.Join(dir, "testproj")
	require.NoError(t, os.MkdirAll(filepath.Join(projDir, "cmd", "server"), 0755))

	tests := []struct {
		name     string
		database string
	}{
		{"postgres", DBPostgres},
		{"mysql", DBMySQL},
		{"sqlite", DBSQLite},
		{"mongodb", DBMongoDB},
		{"dynamodb", DBDynamoDB},
		{"elasticsearch", DBElasticsearch},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			sm := NewSafetyManager(true, false, false)
			createMainGo("testproj", "github.com/test/proj", tc.database, sm)
			assert.NotEmpty(t, sm.GetPendingFiles(), "expected pending files for %s", tc.database)
		})
	}
}

func TestCreateGitignore_DryRun(t *testing.T) {
	cleanup := setupDiscardUI(t)
	defer cleanup()

	dir := t.TempDir()
	sm := NewSafetyManager(true, false, false)
	createGitignore(dir, sm)
	assert.Len(t, sm.GetPendingFiles(), 1)
	assert.Contains(t, sm.GetPendingFiles()[0].Path, ".gitignore")
}

func TestCreateReadme_DryRun(t *testing.T) {
	cleanup := setupDiscardUI(t)
	defer cleanup()

	dir := t.TempDir()
	sm := NewSafetyManager(true, false, false)
	createReadme(dir, "github.com/test/proj", sm)
	assert.Len(t, sm.GetPendingFiles(), 1)
	assert.Contains(t, sm.GetPendingFiles()[0].Path, "README.md")
}

func TestCreateConfig_DryRun(t *testing.T) {
	cleanup := setupDiscardUI(t)
	defer cleanup()

	dir := t.TempDir()
	sm := NewSafetyManager(true, false, false)
	createConfig(dir, "", "postgres", sm)
	assert.NotEmpty(t, sm.GetPendingFiles())
}

func TestCreateLogger_DryRun(t *testing.T) {
	cleanup := setupDiscardUI(t)
	defer cleanup()

	dir := t.TempDir()
	sm := NewSafetyManager(true, false, false)
	createLogger(dir, "", sm)
	assert.Len(t, sm.GetPendingFiles(), 1)
}

func TestCreateAuth_DryRun(t *testing.T) {
	cleanup := setupDiscardUI(t)
	defer cleanup()

	dir := t.TempDir()
	sm := NewSafetyManager(true, false, false)
	createAuth(dir, "github.com/test/proj", sm)
	assert.Len(t, sm.GetPendingFiles(), 1)
}

func TestCreateEnvFiles_DryRun(t *testing.T) {
	cleanup := setupDiscardUI(t)
	defer cleanup()

	dir := t.TempDir()
	tests := []struct {
		name     string
		database string
	}{
		{"postgres", "postgres"},
		{"mysql", "mysql"},
		{"sqlite", "sqlite"},
		{"mongodb", "mongodb"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			sm := NewSafetyManager(true, false, false)
			createEnvFiles(dir, tc.database, sm)
			assert.NotEmpty(t, sm.GetPendingFiles())
		})
	}
}

// --- Entity generator functions (DryRun) ---

func TestGenerateErrorsFile_DryRun(t *testing.T) {
	// Not parallel: os.Chdir
	cleanup := setupDiscardUI(t)
	defer cleanup()

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := setupProjectDir(t)
	os.Chdir(dir)

	fields := []Field{
		{Name: "Name", Type: "string"},
		{Name: "Email", Type: "string"},
	}
	sm := NewSafetyManager(true, false, false)
	generateErrorsFile(dir, "Product", fields, sm)
	assert.NotEmpty(t, sm.GetPendingFiles())
}

func TestGenerateSeedData_DryRun(t *testing.T) {
	// Not parallel: os.Chdir
	cleanup := setupDiscardUI(t)
	defer cleanup()

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := setupProjectDir(t)
	os.Chdir(dir)

	fields := []Field{
		{Name: "Name", Type: "string"},
		{Name: "Price", Type: "float64"},
		{Name: "Quantity", Type: "int"},
	}
	sm := NewSafetyManager(true, false, false)
	generateSeedData(dir, "Product", fields, sm)
	assert.NotEmpty(t, sm.GetPendingFiles())
}

// --- UseCase generator functions (DryRun) ---

func TestGenerateDTOFileWithFields_DryRun(t *testing.T) {
	// Not parallel: os.Chdir
	cleanup := setupDiscardUI(t)
	defer cleanup()

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := setupProjectDir(t)
	os.Chdir(dir)

	sm := NewSafetyManager(true, false, false)
	generateDTOFileWithFields(dir, "Product", []string{"create", "get", "list", "update", "delete"}, true, "Name:string,Price:float64", sm)
	assert.NotEmpty(t, sm.GetPendingFiles())
}

func TestGenerateUseCaseServiceWithFields_DryRun(t *testing.T) {
	// Not parallel: os.Chdir
	cleanup := setupDiscardUI(t)
	defer cleanup()

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := setupProjectDir(t)
	os.Chdir(dir)

	sm := NewSafetyManager(true, false, false)
	generateUseCaseServiceWithFields(dir, "Product", "Product", []string{"create", "get", "list", "update", "delete"}, false, "Name:string,Price:float64", sm)
	assert.NotEmpty(t, sm.GetPendingFiles())
}

func TestGenerateUseCaseServiceWithFields_Async_DryRun(t *testing.T) {
	// Not parallel: os.Chdir
	cleanup := setupDiscardUI(t)
	defer cleanup()

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := setupProjectDir(t)
	os.Chdir(dir)

	sm := NewSafetyManager(true, false, false)
	generateUseCaseServiceWithFields(dir, "Order", "Order", []string{"create", "get"}, true, "Total:float64,Status:string", sm)
	assert.NotEmpty(t, sm.GetPendingFiles())
}

// --- DI generator functions (DryRun) ---

func TestGenerateWireFile_DryRun(t *testing.T) {
	// Not parallel: os.Chdir
	cleanup := setupDiscardUI(t)
	defer cleanup()

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := setupProjectDir(t)
	os.Chdir(dir)

	sm := NewSafetyManager(true, false, false)
	generateWireFile(dir, []string{"Product", "User"}, "postgres", sm)
	assert.NotEmpty(t, sm.GetPendingFiles())
}

func TestGenerateWireGenTemplate_DryRun(t *testing.T) {
	// Not parallel: os.Chdir
	cleanup := setupDiscardUI(t)
	defer cleanup()

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := setupProjectDir(t)
	os.Chdir(dir)

	sm := NewSafetyManager(true, false, false)
	generateWireGenTemplate(dir, []string{"Product"}, sm)
	assert.NotEmpty(t, sm.GetPendingFiles())
}

// --- Repository generator functions (DryRun) ---

func TestGenerateRepositoryInterface_DryRun(t *testing.T) {
	// Not parallel: os.Chdir
	cleanup := setupDiscardUI(t)
	defer cleanup()

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := setupProjectDir(t)
	os.Chdir(dir)

	t.Run("without transactions", func(t *testing.T) {
		sm := NewSafetyManager(true, false, false)
		generateRepositoryInterface(dir, "Product", false, sm)
		assert.NotEmpty(t, sm.GetPendingFiles())
	})

	t.Run("with transactions", func(t *testing.T) {
		sm := NewSafetyManager(true, false, false)
		generateRepositoryInterface(dir, "Order", true, sm)
		assert.NotEmpty(t, sm.GetPendingFiles())
	})
}

func TestGenerateRepositoryInterfaceWithFields_DryRun(t *testing.T) {
	// Not parallel: os.Chdir
	cleanup := setupDiscardUI(t)
	defer cleanup()

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := setupProjectDir(t)
	os.Chdir(dir)

	fields := []Field{
		{Name: "Name", Type: "string"},
		{Name: "Price", Type: "float64"},
	}
	sm := NewSafetyManager(true, false, false)
	generateRepositoryInterfaceWithFields(dir, "Product", fields, true, sm)
	assert.NotEmpty(t, sm.GetPendingFiles())
}

// --- Handler generator functions (DryRun) ---

func TestGenerateProtoFile_DryRun(t *testing.T) {
	// Not parallel: os.Chdir
	cleanup := setupDiscardUI(t)
	defer cleanup()

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := setupProjectDir(t)
	os.Chdir(dir)

	sm := NewSafetyManager(true, false, false)
	generateProtoFile(dir, "Product", "snake_case", sm)
	assert.NotEmpty(t, sm.GetPendingFiles())
}

// --- Init project files pure helpers ---

func TestGetDatabasePort_Extended(t *testing.T) {
	t.Parallel()
	tests := []struct {
		db   string
		port string
	}{
		{"postgres", "5432"},
		{"mysql", "3306"},
		{"mongodb", "27017"},
		{"sqlite", "5432"}, // falls through to default
		{"unknown", "5432"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.db, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.port, getDatabasePort(tc.db))
		})
	}
}

func TestGetDatabaseUser_Extended(t *testing.T) {
	t.Parallel()
	tests := []struct {
		db   string
		user string
	}{
		{"postgres", "postgres"},
		{"mysql", "root"},
		{"mongodb", "admin"},
		{"sqlite", "postgres"}, // falls through to default
		{"unknown", "postgres"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.db, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.user, getDatabaseUser(tc.db))
		})
	}
}
