package tests

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/sazardev/goca/internal/testing/framework"
)

// TestInitCommand prueba exhaustivamente el comando 'init'
func TestInitCommand(t *testing.T) {
	// Crear contexto de test
	tc := framework.NewTestContext(t)
	defer tc.Cleanup()
	tc.CurrentTestName = "TestInitCommand"

	// Ejecutar subpruebas para diferentes escenarios
	t.Run("InitWithDefaultOptions", func(t *testing.T) {
		testInitWithDefaultOptions(tc, t)
	})

	t.Run("InitWithModuleFlag", func(t *testing.T) {
		testInitWithModuleFlag(tc, t)
	})

	t.Run("InitWithDatabaseFlag", func(t *testing.T) {
		testInitWithDatabaseFlag(tc, t)
	})

	t.Run("InitWithAPIFlag", func(t *testing.T) {
		testInitWithAPIFlag(tc, t)
	})

	t.Run("InitWithAuthFlag", func(t *testing.T) {
		testInitWithAuthFlag(tc, t)
	})

	t.Run("InitWithAllOptions", func(t *testing.T) {
		testInitWithAllOptions(tc, t)
	})

	// Imprimir resumen
	tc.PrintTestSummary()
}

// testInitWithDefaultOptions prueba el comando init con opciones por defecto
func testInitWithDefaultOptions(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("init", "test-default", "--module", "github.com/test/default")
	if err != nil {
		t.Fatalf("Error al ejecutar comando init: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Initializing project 'test-default'") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar estructura de archivos y directorios
	essentialFiles := []string{
		"go.mod",
		"cmd/server/main.go",
		"internal/domain",
		"internal/usecase",
		"internal/repository",
		"internal/handler",
	}

	for _, file := range essentialFiles {
		tc.AssertFileExists(filepath.Join("test-default", file))
	}

	// Verificar contenido del go.mod
	tc.AssertFileContains(filepath.Join("test-default", "go.mod"),
		"module github.com/test/default")

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-default")

	// Verificar que pasa go vet
	tc.AssertGoVet("test-default")
}

// testInitWithModuleFlag prueba el comando init con el flag --module
func testInitWithModuleFlag(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	moduleName := "github.com/custom/module"
	output, err := tc.RunCommand("init", "test-module", "--module", moduleName)
	if err != nil {
		t.Fatalf("Error al ejecutar comando init: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, moduleName) {
		t.Errorf("Nombre de módulo no encontrado en la salida: %s", output)
	}

	// Verificar go.mod
	tc.AssertFileContains(filepath.Join("test-module", "go.mod"),
		"module "+moduleName)

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-module")
}

// testInitWithDatabaseFlag prueba el comando init con el flag --database
func testInitWithDatabaseFlag(tc *framework.TestContext, t *testing.T) {
	databases := []string{"postgres", "mysql", "mongodb"}

	for _, db := range databases {
		projectName := "test-db-" + db

		// Ejecutar comando
		output, err := tc.RunCommand("init", projectName, "--module", "github.com/test/"+projectName, "--database", db)
		if err != nil {
			t.Fatalf("Error al ejecutar comando init con database=%s: %v", db, err)
		}

		// Verificar salida
		if !strings.Contains(output, "Database: "+db) {
			t.Errorf("Base de datos no encontrada en la salida para %s: %s", db, output)
		}

		// Verificar archivos específicos para la base de datos
		configPath := filepath.Join(projectName, "pkg", "config", "config.go")
		tc.AssertFileExists(configPath)
		// La configuración de la base de datos está ahora en pkg/config/config.go
		// No verificamos el contenido exacto para evitar dependencias de implementación

		// Verificar que el proyecto compila
		tc.AssertGoBuild(projectName)
	}
}

// testInitWithAPIFlag prueba el comando init con el flag --api
func testInitWithAPIFlag(tc *framework.TestContext, t *testing.T) {
	apiTypes := []string{"rest", "grpc", "graphql"}

	for _, api := range apiTypes {
		projectName := "test-api-" + api

		// Ejecutar comando
		output, err := tc.RunCommand("init", projectName, "--module", "github.com/test/"+projectName, "--api", api)
		if err != nil {
			t.Fatalf("Error al ejecutar comando init con api=%s: %v", api, err)
		}

		// Verificar salida
		if !strings.Contains(output, "API: "+api) {
			t.Errorf("Tipo de API no encontrado en la salida para %s: %s", api, output)
		}

		// Verificar archivos específicos para el tipo de API
		mainPath := filepath.Join(projectName, "cmd", "server", "main.go")
		tc.AssertFileExists(mainPath)

		switch api {
		case "rest":
			tc.AssertFileContains(mainPath, "gorilla")
		case "grpc":
			tc.AssertFileContains(mainPath, "grpc")
		case "graphql":
			tc.AssertFileContains(mainPath, "graphql")
		}

		// Verificar que el proyecto compila
		tc.AssertGoBuild(projectName)
	}
}

// testInitWithAuthFlag prueba el comando init con el flag --auth
func testInitWithAuthFlag(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	projectName := "test-auth"
	output, err := tc.RunCommand("init", projectName, "--module", "github.com/test/"+projectName, "--auth")
	if err != nil {
		t.Fatalf("Error al ejecutar comando init con auth: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Including authentication") {
		t.Errorf("Mensaje de autenticación no encontrado en la salida: %s", output)
	}

	// Verificar archivos relacionados con autenticación
	authPaths := []string{
		filepath.Join(projectName, "pkg", "auth", "jwt.go"),
		filepath.Join(projectName, "pkg", "config", "config.go"),
	}

	for _, path := range authPaths {
		tc.AssertFileExists(path)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild(projectName)
}

// testInitWithAllOptions prueba el comando init con todas las opciones
func testInitWithAllOptions(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	projectName := "test-full"
	output, err := tc.RunCommand("init", projectName,
		"--module", "github.com/test/"+projectName,
		"--database", "postgres",
		"--api", "rest",
		"--auth")

	if err != nil {
		t.Fatalf("Error al ejecutar comando init con todas las opciones: %v", err)
	}

	// Verificar salida para cada opción
	expectedOutputs := []string{
		"Initializing project 'test-full'",
		"Database: postgres",
		"API: rest",
		"Including authentication",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Salida esperada no encontrada: %s", expected)
		}
	}

	// Verificar existencia de archivos clave
	keyPaths := []string{
		filepath.Join(projectName, "go.mod"),
		filepath.Join(projectName, "cmd", "server", "main.go"),
		filepath.Join(projectName, "pkg", "config", "config.go"),
		filepath.Join(projectName, "pkg", "auth", "jwt.go"),
		filepath.Join(projectName, "internal", "handler"),
	}

	for _, path := range keyPaths {
		tc.AssertFileExists(path)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild(projectName)

	// Verificar que pasa go vet
	tc.AssertGoVet(projectName)
}
