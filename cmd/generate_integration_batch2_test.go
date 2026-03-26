package cmd

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// ensureTestUI initializes the global ui variable if nil, for tests that
// call functions which use ui.* directly. Returns a cleanup func.
func ensureTestUI(t *testing.T) func() {
	t.Helper()
	if ui != nil {
		return func() {}
	}
	ui = NewUIRenderer(io.Discard, true, 0)
	return func() { ui = nil }
}

// TestDryRunGenerators_Batch2 tests additional generators that require os.Chdir.
// These cannot be parallel.
func TestDryRunGenerators_Batch2(t *testing.T) {
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { require.NoError(t, os.Chdir(origDir)) }()

	cleanup := ensureTestUI(t)
	defer cleanup()

	t.Run("generateDI manual", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		generateDI("Product,Order", "postgres", false, sm)
	})

	t.Run("generateDI wire", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		generateDI("Product", "mysql", true, sm)
	})

	t.Run("generateCompleteFeature", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		generateCompleteFeature("Product", "Name:string,Price:float64", "postgres", "http", true, false, "lowercase", sm)
	})

	t.Run("generateCompleteFeature grpc", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		generateCompleteFeature("Order", "Total:float64", "mysql", "grpc", false, true, "snake", sm)
	})

	t.Run("generateEntityTests", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		fields := []Field{
			{Name: "Name", Type: "string"},
			{Name: "Price", Type: "float64"},
			{Name: "Active", Type: "bool"},
		}
		generateEntityTests(dir, "Product", fields, true, true, "lowercase", sm)
	})

	t.Run("generateEntityTests snake", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		fields := []Field{
			{Name: "Email", Type: "string"},
			{Name: "Age", Type: "int"},
		}
		generateEntityTests(dir, "User", fields, false, false, "snake", sm)
	})

	t.Run("createProjectStructure", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		ci := &ConfigIntegration{config: &GocaConfig{}}
		createProjectStructure("myproject", "github.com/user/myproject", "postgres", false, "rest", ci, false, "", sm)
	})

	t.Run("addEntityToAutoMigration", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		// No main.go file, should return an error or handle gracefully
		_ = addEntityToAutoMigration("Product", sm)
	})
}
