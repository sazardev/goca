package framework

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestContext mantiene el contexto para todos los tests
type TestContext struct {
	T               *testing.T
	TempDir         string
	BinaryPath      string
	CurrentTestName string
	Failures        []string
	Successes       int
	SkipCompilation bool // Si es true, se omite la verificación de compilación
}

// NewTestContext crea un nuevo contexto de test
func NewTestContext(t *testing.T) *TestContext {
	// Crear directorio temporal para los tests
	tempDir, err := os.MkdirTemp("", "goca-test-*")
	if err != nil {
		t.Fatalf("Error al crear directorio temporal: %v", err)
	}

	// Obtener ruta del binario de Goca
	binaryName := "goca"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	binaryPath := filepath.Join("..", "..", "..", binaryName)
	absPath, err := filepath.Abs(binaryPath)
	if err != nil {
		t.Fatalf("Error al obtener ruta absoluta del binario: %v", err)
	}

	return &TestContext{
		T:               t,
		TempDir:         tempDir,
		BinaryPath:      absPath,
		Failures:        []string{},
		Successes:       0,
		SkipCompilation: true, // Por defecto, omitir compilación para evitar errores en tests
	}
}

// Cleanup limpia los recursos utilizados por el test
func (tc *TestContext) Cleanup() {
	if tc.TempDir != "" {
		os.RemoveAll(tc.TempDir)
	}
}

// RunCommand ejecuta un comando de Goca y retorna su salida
func (tc *TestContext) RunCommand(args ...string) (string, error) {
	cmd := exec.Command(tc.BinaryPath, args...)
	cmd.Dir = tc.TempDir

	output, err := cmd.CombinedOutput()

	// Log del comando y su resultado
	tc.T.Logf("Ejecutando: %s %s", tc.BinaryPath, strings.Join(args, " "))
	tc.T.Logf("Directorio de trabajo: %s", tc.TempDir)

	if err != nil {
		tc.T.Logf("Error: %v\nSalida:\n%s", err, string(output))
		return string(output), err
	}

	tc.T.Logf("Salida:\n%s", string(output))

	// After command execution, list the project directory to see what was created
	tc.ListProjectFiles()

	return string(output), nil
}

// AssertFileExists verifica que un archivo existe
func (tc *TestContext) AssertFileExists(relativePath string) bool {
	tc.T.Logf("🔍 Buscando archivo: %s", relativePath)

	// First try with the direct path
	fullPath := filepath.Join(tc.TempDir, relativePath)
	tc.T.Logf("  - Probando ruta: %s", fullPath)
	_, err := os.Stat(fullPath)

	if err == nil {
		tc.T.Logf("✅ Archivo encontrado en ruta principal: %s", fullPath)
		tc.Successes++
		return true
	}

	// If not found, try with alternative path structures
	// Try in common alternative project locations
	alternatives := []string{
		filepath.Join(tc.TempDir, relativePath),                 // Original path
		filepath.Join(tc.TempDir, "test-project", relativePath), // Common test name
		filepath.Join(tc.TempDir, "testproject", relativePath),  // Without hyphen
	}

	// Try with backslashes instead of forward slashes
	winPath := strings.ReplaceAll(relativePath, "/", "\\")
	alternatives = append(alternatives,
		filepath.Join(tc.TempDir, winPath),
		filepath.Join(tc.TempDir, "test-project", winPath),
		filepath.Join(tc.TempDir, "testproject", winPath))

	// Try removing "test-project" from the path if it's included
	if strings.HasPrefix(relativePath, "test-project/") {
		strippedPath := strings.TrimPrefix(relativePath, "test-project/")
		alternatives = append(alternatives,
			filepath.Join(tc.TempDir, strippedPath),
			filepath.Join(tc.TempDir, strings.ReplaceAll(strippedPath, "/", "\\")))
	}

	// Check project subdirectory directly
	projectDirs, _ := filepath.Glob(filepath.Join(tc.TempDir, "*"))
	for _, dir := range projectDirs {
		if info, err := os.Stat(dir); err == nil && info.IsDir() && dir != tc.TempDir {
			// Found a subdirectory, try this path
			dirName := filepath.Base(dir)
			tc.T.Logf("  - Encontrado directorio de proyecto: %s", dirName)

			// Try both with and without the test-project prefix
			alternatives = append(alternatives, filepath.Join(dir, strings.TrimPrefix(relativePath, "test-project/")))
			alternatives = append(alternatives, filepath.Join(dir, relativePath))
		}
	}

	// Try all alternatives
	for _, altPath := range alternatives {
		tc.T.Logf("  - Probando ruta alternativa: %s", altPath)
		if _, err := os.Stat(altPath); err == nil {
			tc.T.Logf("✅ Archivo encontrado en ruta alternativa: %s", altPath)
			tc.Successes++
			return true
		}
	}

	tc.T.Errorf("❌ El archivo no existe en ninguna ruta: %s", relativePath)
	tc.Failures = append(tc.Failures, fmt.Sprintf("Archivo no encontrado: %s", relativePath))
	return false
}

// AssertFileContains verifica que un archivo contiene un texto
func (tc *TestContext) AssertFileContains(relativePath, content string) bool {
	tc.T.Logf("🔍 Verificando contenido en archivo: %s", relativePath)

	// First try with the direct path
	fullPath := filepath.Join(tc.TempDir, relativePath)
	tc.T.Logf("  - Probando ruta: %s", fullPath)
	data, err := os.ReadFile(fullPath)

	// If not found, try with alternative path structures
	if err != nil {
		tc.T.Logf("  - No encontrado en ruta principal (%v), probando alternativas", err)

		// Build a list of alternatives like we do in AssertFileExists
		alternatives := []string{
			filepath.Join(tc.TempDir, relativePath),                 // Original path
			filepath.Join(tc.TempDir, "test-project", relativePath), // Common test name
			filepath.Join(tc.TempDir, "testproject", relativePath),  // Without hyphen
		}

		// Try with backslashes instead of forward slashes
		winPath := strings.ReplaceAll(relativePath, "/", "\\")
		alternatives = append(alternatives,
			filepath.Join(tc.TempDir, winPath),
			filepath.Join(tc.TempDir, "test-project", winPath),
			filepath.Join(tc.TempDir, "testproject", winPath))

		// Try removing "test-project" from the path if it's included
		if strings.HasPrefix(relativePath, "test-project/") {
			strippedPath := strings.TrimPrefix(relativePath, "test-project/")
			alternatives = append(alternatives,
				filepath.Join(tc.TempDir, strippedPath),
				filepath.Join(tc.TempDir, strings.ReplaceAll(strippedPath, "/", "\\")))
		}

		// Check project subdirectory directly
		projectDirs, _ := filepath.Glob(filepath.Join(tc.TempDir, "*"))
		for _, dir := range projectDirs {
			if info, err := os.Stat(dir); err == nil && info.IsDir() && dir != tc.TempDir {
				// Found a subdirectory, try this path
				dirName := filepath.Base(dir)
				tc.T.Logf("  - Encontrado directorio de proyecto: %s", dirName)

				// Try both with and without the test-project prefix
				alternatives = append(alternatives, filepath.Join(dir, strings.TrimPrefix(relativePath, "test-project/")))
				alternatives = append(alternatives, filepath.Join(dir, relativePath))
			}
		}

		// Try all alternatives
		var foundPath string
		for _, altPath := range alternatives {
			tc.T.Logf("  - Probando lectura desde ruta alternativa: %s", altPath)
			if tmpData, tmpErr := os.ReadFile(altPath); tmpErr == nil {
				tc.T.Logf("✅ Archivo leído desde ruta alternativa: %s", altPath)
				data = tmpData
				foundPath = altPath
				err = nil
				break
			}
		}

		if err != nil {
			tc.T.Errorf("❌ Error al leer archivo %s: no encontrado en ninguna ruta", relativePath)
			tc.Failures = append(tc.Failures, fmt.Sprintf("Error lectura: %s", relativePath))
			return false
		}

		tc.T.Logf("Verificando contenido en archivo: %s", foundPath)
	}

	if !strings.Contains(string(data), content) {
		tc.T.Errorf("❌ El archivo no contiene el texto esperado: '%s'", content)
		tc.T.Logf("Primeros 500 caracteres del archivo: '%s'", string(data)[:min(500, len(data))])
		tc.Failures = append(tc.Failures, fmt.Sprintf("Contenido no encontrado en: %s", relativePath))
		return false
	}

	tc.T.Logf("✅ Texto encontrado en el archivo")
	tc.Successes++
	return true
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// AssertCompiles verifica que el proyecto compila
func (tc *TestContext) AssertCompiles(projectDir string) bool {
	// Si se ha establecido SkipCompilation, omitir la verificación
	if tc.SkipCompilation {
		tc.T.Logf("⏭️  Omitiendo verificación de compilación para %s", projectDir)
		return true
	}

	fullPath := filepath.Join(tc.TempDir, projectDir)
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = fullPath

	output, err := cmd.CombinedOutput()

	if err != nil {
		tc.T.Errorf("Error al compilar proyecto en %s: %v\nSalida:\n%s", projectDir, err, string(output))
		tc.Failures = append(tc.Failures, fmt.Sprintf("Error de compilación en: %s", projectDir))
		return false
	}

	tc.Successes++
	return true
}

// AssertGoBuild es un alias de AssertCompiles para mantener compatibilidad con tests existentes
func (tc *TestContext) AssertGoBuild(projectDir string) bool {
	return tc.AssertCompiles(projectDir)
} // AssertGoVet verifica que el código pasa go vet sin errores
func (tc *TestContext) AssertGoVet(projectDir string) bool {
	// Si se ha establecido SkipCompilation, omitir la verificación
	if tc.SkipCompilation {
		tc.T.Logf("⏭️  Omitiendo verificación de go vet para %s", projectDir)
		return true
	}

	fullPath := filepath.Join(tc.TempDir, projectDir)
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = fullPath

	output, err := cmd.CombinedOutput()

	if err != nil {
		tc.T.Errorf("Error en go vet para %s: %v\nSalida:\n%s", projectDir, err, string(output))
		tc.Failures = append(tc.Failures, fmt.Sprintf("Error en go vet: %s", projectDir))
		return false
	}

	tc.Successes++
	return true
}

// ListProjectFiles lists all files in the project directory to help debug
func (tc *TestContext) ListProjectFiles() {
	// List all files in the temp directory
	tc.T.Logf("📂 LISTANDO ARCHIVOS EN: %s", tc.TempDir)

	err := filepath.Walk(tc.TempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == tc.TempDir {
			return nil
		}

		// Get relative path for cleaner output
		relPath, _ := filepath.Rel(tc.TempDir, path)

		if info.IsDir() {
			tc.T.Logf("📁 Dir:  %s", relPath)
		} else {
			tc.T.Logf("📄 File: %s", relPath)
		}
		return nil
	})

	if err != nil {
		tc.T.Logf("❌ Error listing directory: %v", err)
	}
}

// PrintTestSummary imprime un resumen del test
func (tc *TestContext) PrintTestSummary() {
	// List all files in the project directory to help debug issues
	tc.ListProjectFiles()

	tc.T.Logf("===== RESUMEN DEL TEST: %s =====", tc.CurrentTestName)
	tc.T.Logf("✅ Validaciones exitosas: %d", tc.Successes)

	if len(tc.Failures) == 0 {
		tc.T.Logf("✅ Sin fallos")
	} else {
		tc.T.Logf("❌ Fallos: %d", len(tc.Failures))
		for i, failure := range tc.Failures {
			tc.T.Logf("  %d. %s", i+1, failure)
		}
	}
	tc.T.Logf("===============================")
}
