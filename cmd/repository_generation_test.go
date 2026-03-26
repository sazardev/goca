package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDryRunGenerateRepository(t *testing.T) {
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { require.NoError(t, os.Chdir(origDir)) }()

	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))
	sm := &SafetyManager{DryRun: true}

	t.Run("full repository postgres", func(t *testing.T) {
		generateRepository("Product", "postgres", false, false, false, false, "Name:string,Price:float64", sm)
	})

	t.Run("full repository mysql", func(t *testing.T) {
		generateRepository("Product", "mysql", false, false, false, false, "", sm)
	})

	t.Run("interface only", func(t *testing.T) {
		generateRepository("Product", "postgres", true, false, false, false, "", sm)
	})

	t.Run("implementation only", func(t *testing.T) {
		generateRepository("Product", "postgres", false, true, false, false, "", sm)
	})

	t.Run("with cache", func(t *testing.T) {
		generateRepository("Product", "postgres", false, false, true, false, "", sm)
	})

	t.Run("with transactions", func(t *testing.T) {
		generateRepository("Product", "postgres", false, false, false, true, "", sm)
	})

	t.Run("generate interface file", func(t *testing.T) {
		generateRepositoryInterface(tmpDir, "Product", false, sm)
	})

	t.Run("generate interface with transactions", func(t *testing.T) {
		generateRepositoryInterface(tmpDir, "Product", true, sm)
	})

	t.Run("generate implementation postgres", func(t *testing.T) {
		generateRepositoryImplementation(tmpDir, "Product", "postgres", false, false, sm)
	})

	t.Run("generate implementation mysql cache", func(t *testing.T) {
		generateRepositoryImplementation(tmpDir, "Product", "mysql", true, false, sm)
	})

	t.Run("generate implementation mongodb", func(t *testing.T) {
		generateRepositoryImplementation(tmpDir, "Product", "mongodb", false, true, sm)
	})
}
