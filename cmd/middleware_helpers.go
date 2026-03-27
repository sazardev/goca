package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

const dirMiddleware = "middleware"

// middlewarePackageExists checks whether the standalone middleware package
// (internal/middleware/middleware.go) has been generated.
func middlewarePackageExists() bool {
	_, err := os.Stat(filepath.Join(DirInternal, dirMiddleware, "middleware.go"))
	return err == nil
}

// middlewareTypeToFile maps a middleware type key to its output filename and
// the generator function that produces the file content.
var middlewareTypeToFile = map[string]struct {
	filename string
	generate func() string
}{
	"cors":       {filename: "cors.go", generate: generateCORSMiddleware},
	"logging":    {filename: "logging.go", generate: generateLoggingMiddleware},
	"auth":       {filename: "auth.go", generate: generateAuthMiddleware},
	"rate-limit": {filename: "rate_limit.go", generate: generateRateLimitMiddleware},
	"recovery":   {filename: "recovery.go", generate: generateRecoveryMiddleware},
	"request-id": {filename: "request_id.go", generate: generateRequestIDMiddleware},
	"timeout":    {filename: "timeout.go", generate: generateTimeoutMiddleware},
}

// generateMiddlewarePackage creates the internal/middleware/ package with the
// requested middleware types and the Chain() helper.
func generateMiddlewarePackage(name string, types []string, sm *SafetyManager) error {
	module := getModuleName()
	middlewareDir := filepath.Join(DirInternal, dirMiddleware)

	// Always write the chain helper first.
	chainPath := filepath.Join(middlewareDir, "middleware.go")
	if err := writeGoFile(chainPath, generateChainMiddleware(module), sm); err != nil {
		return fmt.Errorf("writing middleware.go: %w", err)
	}
	ui.FileCreated(chainPath)

	// Write each requested middleware type.
	for _, t := range types {
		entry, ok := middlewareTypeToFile[t]
		if !ok {
			return fmt.Errorf("unknown middleware type: %s", t)
		}
		filePath := filepath.Join(middlewareDir, entry.filename)
		if err := writeGoFile(filePath, entry.generate(), sm); err != nil {
			return fmt.Errorf("writing %s: %w", entry.filename, err)
		}
		ui.FileCreated(filePath)
	}

	return nil
}
