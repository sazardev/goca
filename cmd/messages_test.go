package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDryRunGenerateMessages(t *testing.T) {
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { require.NoError(t, os.Chdir(origDir)) }()

	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))
	sm := &SafetyManager{DryRun: true}

	t.Run("generate all messages", func(t *testing.T) {
		generateMessages("Product", true, true, true, sm)
	})

	t.Run("generate usecase messages", func(t *testing.T) {
		generateUseCaseMessages("Product", sm)
	})

	t.Run("generate response messages", func(t *testing.T) {
		generateResponseMessages(tmpDir, "Product", sm)
	})

	t.Run("generate constants", func(t *testing.T) {
		generateConstants(tmpDir, "Product", sm)
	})
}
