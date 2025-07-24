package cmd

import (
	"bufio"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
)

// getModuleName reads the module name from go.mod file
func getModuleName() string {
	goMod, err := os.Open("go.mod")
	if err != nil {
		return "myproject" // fallback
	}
	defer goMod.Close()

	scanner := bufio.NewScanner(goMod)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}
	return "myproject" // fallback
}

// writeFile creates a file with the given content, creating directories if needed
func writeFile(path, content string) {
	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0755)

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", path, err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", path, err)
	}
}

// writeGoFile creates a Go file with auto-formatting
func writeGoFile(path, content string) {
	// Format Go code if it's a .go file
	if strings.HasSuffix(path, ".go") {
		formatted, err := format.Source([]byte(content))
		if err != nil {
			fmt.Printf("Warning: Failed to format Go code for %s: %v\n", path, err)
			// Continue with unformatted code
		} else {
			content = string(formatted)
		}
	}

	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0755)

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", path, err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", path, err)
	}
}
