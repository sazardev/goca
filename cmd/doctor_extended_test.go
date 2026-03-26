package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Doctor check functions with filesystem ---

func TestCheckGoMod_NoFile(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	check := checkGoMod()
	assert.Equal(t, "✗", check.status)
	assert.Contains(t, check.message, "not found")
}

func TestCheckGoMod_ValidFile(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)
	require.NoError(t, os.WriteFile("go.mod", []byte("module github.com/test/project\n\ngo 1.21\n"), 0644))

	check := checkGoMod()
	assert.Equal(t, "✓", check.status)
	assert.Contains(t, check.message, "module declaration")
}

func TestCheckGoMod_NoModuleDeclaration(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)
	require.NoError(t, os.WriteFile("go.mod", []byte("go 1.21\n"), 0644))

	check := checkGoMod()
	assert.Equal(t, "⚠", check.status)
	assert.Contains(t, check.message, "no module declaration")
}

func TestCheckGocaYaml_NoFile(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	check := checkGocaYaml()
	assert.Equal(t, "⚠", check.status)
	assert.Contains(t, check.message, "not found")
}

func TestCheckGocaYaml_ValidFile(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)
	require.NoError(t, os.WriteFile(".goca.yaml", []byte("project: test\ndatabase: postgres\n"), 0644))

	check := checkGocaYaml()
	assert.Equal(t, "✓", check.status)
	assert.Contains(t, check.message, "non-empty")
}

func TestCheckGocaYaml_EmptyFile(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)
	require.NoError(t, os.WriteFile(".goca.yaml", []byte("  \n"), 0644))

	check := checkGocaYaml()
	assert.Equal(t, "⚠", check.status)
	assert.Contains(t, check.message, "empty")
}

func TestCheckProjectStructure_AllPresent(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)
	for _, d := range []string{"internal/domain", "internal/usecase", "internal/repository", "internal/handler"} {
		require.NoError(t, os.MkdirAll(d, 0755))
	}

	check := checkProjectStructure()
	assert.Equal(t, "✓", check.status)
}

func TestCheckProjectStructure_Missing(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	check := checkProjectStructure()
	assert.Equal(t, "⚠", check.status)
	assert.Contains(t, check.message, "missing")
}

func TestCheckProjectStructure_FixMode(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)

	// Set up UI for ui.Debug calls
	cleanup := setupDiscardUI(t)
	defer cleanup()

	oldFix := doctorFix
	doctorFix = true
	defer func() { doctorFix = oldFix }()

	check := checkProjectStructure()
	assert.Equal(t, "✓", check.status)
	assert.Contains(t, check.message, "Created")

	// Verify directories were actually created
	for _, d := range []string{"internal/domain", "internal/usecase", "internal/repository", "internal/handler"} {
		_, err := os.Stat(d)
		assert.NoError(t, err, "expected %s to exist", d)
	}
}

func TestCheckGoBuild_NoGoMod(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	check := checkGoBuild()
	assert.Equal(t, "⚠", check.status)
	assert.Contains(t, check.message, "no go.mod")
}

func TestCheckGoVet_NoGoMod(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	check := checkGoVet()
	assert.Equal(t, "⚠", check.status)
	assert.Contains(t, check.message, "no go.mod")
}

func TestCheckDIContainer_Found(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)
	require.NoError(t, os.MkdirAll("internal/di", 0755))

	check := checkDIContainer()
	assert.Equal(t, "✓", check.status)
	assert.Contains(t, check.message, "DI container found")
}

func TestCheckDIContainer_NotFound(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	check := checkDIContainer()
	assert.Equal(t, "⚠", check.status)
	assert.Contains(t, check.message, "No DI container")
}

func TestRunAllChecks(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	checks := runAllChecks()
	assert.Len(t, checks, 6)
	// In empty dir: gomod=fail, gocayaml=warn, structure=warn, build=warn, vet=warn, di=warn
	for _, c := range checks {
		assert.NotEmpty(t, c.name)
		assert.NotEmpty(t, c.status)
		assert.NotEmpty(t, c.message)
	}
}

func TestPrintChecks(t *testing.T) {
	// Not parallel: uses global ui
	buf := &bytes.Buffer{}
	testUI := NewUIRenderer(buf, true, 2)
	oldUI := ui
	ui = testUI
	defer func() { ui = oldUI }()

	checks := []doctorCheck{
		{name: "test1", status: "✓", message: "All good", suggestion: ""},
		{name: "test2", status: "✗", message: "Failed", suggestion: "Fix it"},
		{name: "test3", status: "⚠", message: "Warning", suggestion: "Check it"},
	}
	printChecks(checks) // Should not panic
	out := buf.String()
	assert.Contains(t, out, "test1")
	assert.Contains(t, out, "test2")
}

// --- Config debug extended tests ---

func TestValidateConfigSilent_Coverage(t *testing.T) {
	// Not parallel: uses global ui
	buf := &bytes.Buffer{}
	testUI := NewUIRenderer(buf, true, 2)
	oldUI := ui
	ui = testUI
	defer func() { ui = oldUI }()

	t.Run("valid config", func(t *testing.T) {
		config := map[string]interface{}{
			"project":  "myapp",
			"module":   "github.com/test/myapp",
			"database": "postgres",
		}
		validateConfigSilent(config) // Should not panic
	})

	t.Run("empty config", func(t *testing.T) {
		validateConfigSilent(map[string]interface{}{}) // Should not panic
	})
}

func TestGetCurrentDir(t *testing.T) {
	result := getCurrentDir()
	assert.NotEmpty(t, result)
}

func TestGetCurrentProjectName(t *testing.T) {
	result := getCurrentProjectName()
	assert.NotEmpty(t, result)
}

func TestGetCurrentModuleName(t *testing.T) {
	result := getCurrentModuleName()
	assert.NotEmpty(t, result)
}

func TestShowCurrentConfig(t *testing.T) {
	// Not parallel: uses global ui
	buf := &bytes.Buffer{}
	testUI := NewUIRenderer(buf, true, 2)
	oldUI := ui
	ui = testUI
	defer func() { ui = oldUI }()

	// showCurrentConfig reads .goca.yaml from CWD - in test dir it won't exist
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	showCurrentConfig() // Should not panic even without .goca.yaml
}

func TestValidateConfiguration(t *testing.T) {
	// Not parallel: uses global ui + filesystem
	buf := &bytes.Buffer{}
	testUI := NewUIRenderer(buf, true, 2)
	oldUI := ui
	ui = testUI
	defer func() { ui = oldUI }()

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	validateConfiguration() // Should not panic in empty directory
}

func TestShowTemplateOptions(t *testing.T) {
	// Not parallel: uses global ui
	buf := &bytes.Buffer{}
	testUI := NewUIRenderer(buf, true, 2)
	oldUI := ui
	ui = testUI
	defer func() { ui = oldUI }()

	showTemplateOptions() // Should not panic
	out := buf.String()
	assert.Contains(t, out, "default")
}

// --- Integrate functions with filesystem setup ---

func TestPrintFeatureStructure(t *testing.T) {
	// Not parallel: uses global ui
	buf := &bytes.Buffer{}
	testUI := NewUIRenderer(buf, true, 2)
	oldUI := ui
	ui = testUI
	defer func() { ui = oldUI }()

	printFeatureStructure("Product", "http")
	out := buf.String()
	assert.Contains(t, out, "product")
}

func TestPrintManualIntegrationInstructions(t *testing.T) {
	// Not parallel: uses global ui
	buf := &bytes.Buffer{}
	testUI := NewUIRenderer(buf, true, 2)
	oldUI := ui
	ui = testUI
	defer func() { ui = oldUI }()

	printManualIntegrationInstructions("Product")
	out := buf.String()
	assert.Contains(t, out, "Product")
}

// --- Template manager functions ---

func TestWriteFile_Real_NoSM(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "sub", "test.go")
	err := writeFile(f, "package sub\n")
	assert.NoError(t, err)
	content, err := os.ReadFile(f)
	require.NoError(t, err)
	assert.Equal(t, "package sub\n", string(content))
}

func TestWriteFile_DryRun_SM(t *testing.T) {
	sm := NewSafetyManager(true, false, false)
	err := writeFile("/tmp/fake/file.go", "package main\n", sm)
	assert.NoError(t, err)
	assert.Len(t, sm.GetPendingFiles(), 1)
}

func TestWriteGoFile_Real_NoSM(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.go")
	err := writeGoFile(f, "package main\n")
	assert.NoError(t, err)
	content, err := os.ReadFile(f)
	require.NoError(t, err)
	assert.Equal(t, "package main\n", string(content))
}

func TestWriteGoFile_DryRun_SM(t *testing.T) {
	sm := NewSafetyManager(true, false, false)
	err := writeGoFile("/tmp/fake/file.go", "package main\n", sm)
	assert.NoError(t, err)
	assert.Len(t, sm.GetPendingFiles(), 1)
}
