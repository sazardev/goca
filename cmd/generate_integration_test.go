package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newDryRunSafety creates a SafetyManager in dry-run mode for testing
func newDryRunSafety() *SafetyManager {
	return &SafetyManager{DryRun: true}
}

func TestGenerateEntityFile_DryRun(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	fields := []Field{
		{Name: "ID", Type: "uint", Tag: "`json:\"id\" gorm:\"primaryKey\"`"},
		{Name: "Name", Type: "string", Tag: "`json:\"name\"`"},
	}

	sm := newDryRunSafety()
	generateEntityFile(dir, "Product", fields, true, false, false, false, "", sm)
}

func TestGenerateEntityFile_WithAllOptions(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	fields := []Field{
		{Name: "ID", Type: "uint", Tag: "`json:\"id\" gorm:\"primaryKey\"`"},
		{Name: "Name", Type: "string", Tag: "`json:\"name\"`"},
		{Name: "Price", Type: "float64", Tag: "`json:\"price\"`"},
		{Name: "CreatedAt", Type: "time.Time", Tag: "`json:\"created_at\"`"},
	}

	sm := newDryRunSafety()
	generateEntityFile(dir, "Product", fields, true, true, true, true, "snake_case", sm)
	generateEntityFile(dir, "Product", fields, false, false, false, false, "kebab-case", sm)
}

func TestGenerateSwaggerFile_DryRun(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	sm := newDryRunSafety()
	generateSwaggerFile(filepath.Join(dir, "handler"), "Product", sm)
}

// TestDryRunGenerators_Sequential tests generators that need os.Chdir (not parallelizable)
func TestDryRunGenerators_Sequential(t *testing.T) {
	// Save original directory
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { require.NoError(t, os.Chdir(origDir)) }()

	sm := newDryRunSafety()

	// Subtest: generateHTTPHandler
	t.Run("HTTPHandler", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
		generateHTTPHandler("Product", false, false, false, "", sm)
	})

	// Subtest: generateHTTPHandlerFile
	t.Run("HTTPHandlerFile", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
		generateHTTPHandlerFile(filepath.Join(dir, "handler"), "Product", false, "", sm)
	})

	// Subtest: generateHTTPRoutesFile
	t.Run("HTTPRoutesFile", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
		generateHTTPRoutesFile(filepath.Join(dir, "handler"), "Product", true, sm)
	})

	// Subtest: generateHTTPDTOFile
	t.Run("HTTPDTOFile", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
		generateHTTPDTOFile(filepath.Join(dir, "handler"), "Product", sm)
	})

	// Subtest: generateManualDI
	t.Run("ManualDI", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
		generateManualDI(filepath.Join(dir, "di"), []string{"Product", "User"}, "postgres", sm)
	})

	// Subtest: generateWireDI
	t.Run("WireDI", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
		generateWireDI(filepath.Join(dir, "di"), []string{"Product"}, "postgres", sm)
	})

	// Subtest: generateUseCaseWithFields
	t.Run("UseCaseWithFields", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
		require.NoError(t, os.MkdirAll("internal/usecase", 0755))
		generateUseCaseWithFields("ProductService", "Product", "create,read,update,delete,list", true, false, "Name:string,Price:float64", sm)
	})

	// Subtest: generateUseCase
	t.Run("UseCase", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
		require.NoError(t, os.MkdirAll("internal/usecase", 0755))
		generateUseCase("ProductService", "Product", "create,read", false, false, sm)
	})

	// Subtest: generateEntity
	t.Run("Entity", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
		require.NoError(t, os.MkdirAll("internal/domain", 0755))
		generateEntity("Product", "Name:string,Price:float64", true, true, true, false, false, "", sm)
	})

	// Subtest: generateManualDI with mysql
	t.Run("ManualDI_MySQL", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
		generateManualDI(filepath.Join(dir, "di"), []string{"Product"}, "mysql", sm)
	})

	// Subtest: generateManualDI with mongodb
	t.Run("ManualDI_MongoDB", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
		generateManualDI(filepath.Join(dir, "di"), []string{"Product"}, "mongodb", sm)
	})

	// Subtest: HTTPHandler with all options
	t.Run("HTTPHandler_AllOptions", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
		generateHTTPHandler("Product", true, true, true, "snake_case", sm)
	})
}

func TestWriteGoFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	t.Run("writes valid go file", func(t *testing.T) {
		t.Parallel()
		path := filepath.Join(dir, "test_write.go")
		content := "package test\n\nfunc Hello() {}\n"
		err := writeGoFile(path, content)
		assert.NoError(t, err)

		data, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Contains(t, string(data), "package test")
	})

	t.Run("writes with safety manager dry run", func(t *testing.T) {
		t.Parallel()
		path := filepath.Join(dir, "dry_run.go")
		sm := newDryRunSafety()
		err := writeGoFile(path, "package test\n", sm)
		assert.NoError(t, err)
		_, err = os.Stat(path)
		assert.True(t, os.IsNotExist(err))
	})
}

func TestWriteFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	t.Run("writes file", func(t *testing.T) {
		t.Parallel()
		path := filepath.Join(dir, "test.txt")
		err := writeFile(path, "hello world")
		assert.NoError(t, err)

		data, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "hello world", string(data))
	})

	t.Run("creates directories", func(t *testing.T) {
		t.Parallel()
		path := filepath.Join(dir, "sub", "dir", "test.txt")
		err := writeFile(path, "nested")
		assert.NoError(t, err)

		data, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "nested", string(data))
	})
}

func TestGetModuleName(t *testing.T) {
	// go test runs from cmd/ dir, no go.mod there — fallback returned
	name := getModuleName()
	assert.NotEmpty(t, name)
}

func TestGetImportPath_Extended(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected string
	}{
		{"github.com/user/project", "github.com/user/project"},
		{"", ""},
		{"myproject", "myproject"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, getImportPath(tc.input))
		})
	}
}
