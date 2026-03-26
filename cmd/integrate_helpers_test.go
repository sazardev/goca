package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractImportSection(t *testing.T) {
	t.Parallel()

	t.Run("with imports", func(t *testing.T) {
		t.Parallel()
		content := `package main

import (
	"fmt"
	"os"
)

func main() {}`
		result := extractImportSection(content)
		assert.Contains(t, result, "import (")
		assert.Contains(t, result, "\"fmt\"")
		assert.Contains(t, result, "\"os\"")
		assert.Contains(t, result, ")")
	})

	t.Run("no imports", func(t *testing.T) {
		t.Parallel()
		content := "package main\n\nfunc main() {}"
		result := extractImportSection(content)
		assert.Empty(t, result)
	})

	t.Run("single line import not matched", func(t *testing.T) {
		t.Parallel()
		content := `package main

import "fmt"

func main() {}`
		result := extractImportSection(content)
		// extractImportSection only matches `import (`
		assert.Empty(t, result)
	})
}
