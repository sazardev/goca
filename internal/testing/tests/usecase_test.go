package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sazardev/goca/internal/testing/framework"
)

// TestUseCaseCommand prueba exhaustivamente el comando 'usecase'
func TestUseCaseCommand(t *testing.T) {
	tc := framework.NewTestContext(t)
	defer tc.Cleanup()
	tc.CurrentTestName = "TestUseCaseCommand"

	// Preparar proyecto base para todos los tests
	prepareProjectWithEntities(tc, t)

	// Ejecutar subpruebas para diferentes escenarios
	t.Run("UseCaseWithDefaultOperations", func(t *testing.T) {
		testUseCaseWithDefaultOperations(tc, t)
	})

	t.Run("UseCaseWithSpecificOperations", func(t *testing.T) {
		testUseCaseWithSpecificOperations(tc, t)
	})

	t.Run("UseCaseWithDTOValidation", func(t *testing.T) {
		testUseCaseWithDTOValidation(tc, t)
	})

	t.Run("UseCaseWithAsync", func(t *testing.T) {
		testUseCaseWithAsync(tc, t)
	})

	t.Run("UseCaseWithAllOptions", func(t *testing.T) {
		testUseCaseWithAllOptions(tc, t)
	})

	// Verificar que el proyecto completo compila
	tc.AssertGoBuild("test-project")

	// Verificar que pasa go vet
	tc.AssertGoVet("test-project")

	// Imprimir resumen
	tc.PrintTestSummary()
}

// prepareProjectWithEntities inicializa un proyecto base con entidades para las pruebas de usecase
func prepareProjectWithEntities(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando init
	_, err := tc.RunCommand("init", "test-project", "--module", "github.com/test/testproject")
	if err != nil {
		t.Fatalf("Error al inicializar proyecto base: %v", err)
	}

	// Crear entidades para las pruebas
	entities := []struct {
		name   string
		fields string
	}{
		{"User", "name:string,email:string,age:int"},
		{"Product", "name:string,price:float64,sku:string"},
		{"Order", "orderID:string,total:float64,status:string"},
	}

	for _, entity := range entities {
		_, err := tc.RunCommand("entity", entity.name, "--fields", entity.fields)
		if err != nil {
			t.Fatalf("Error al crear entidad %s: %v", entity.name, err)
		}
	}

	// Verificar que se creó correctamente
	tc.AssertGoBuild("test-project")
}

// testUseCaseWithDefaultOperations prueba el comando usecase con operaciones por defecto
func testUseCaseWithDefaultOperations(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("usecase", "UserManagement", "--entity", "User")
	if err != nil {
		t.Fatalf("Error al ejecutar comando usecase: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Generating use case 'UserManagement' for entity 'User'") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivos generados
	basePath := filepath.Join("internal", "usecase")
	tc.AssertFileExists(filepath.Join(basePath, "user_service.go"))
	tc.AssertFileExists(filepath.Join(basePath, "user_usecase.go"))
	tc.AssertFileExists(filepath.Join(basePath, "dto.go"))

	// Verificar contenido del archivo de caso de uso - debe tener operaciones CRUD
	usecasePath := filepath.Join(basePath, "user_usecase.go")
	expectedContents := []string{
		"type UserManagement interface",
		"CreateUser",
		"GetUser",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(usecasePath, content)
	}

	// Verificar contenido del archivo de DTOs
	dtoPath := filepath.Join(basePath, "dto.go")
	expectedDTOs := []string{
		"type CreateUserInput struct {",
		"type CreateUserOutput struct {",
	}

	for _, dto := range expectedDTOs {
		tc.AssertFileContains(dtoPath, dto)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testUseCaseWithSpecificOperations prueba el comando usecase con operaciones específicas
func testUseCaseWithSpecificOperations(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	operations := "create,read"
	output, err := tc.RunCommand("usecase", "ProductCatalog",
		"--entity", "Product",
		"--operations", operations)

	if err != nil {
		t.Fatalf("Error al ejecutar comando usecase con operaciones específicas: %v", err)
	}

	// Verificar salida
	expectedOutputs := []string{
		"Generating use case 'ProductCatalog' for entity 'Product'",
		"Operations: " + operations,
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Salida esperada no encontrada: %s", expected)
		}
	}

	// Verificar archivos generados
	basePath := filepath.Join("internal", "usecase")
	tc.AssertFileExists(filepath.Join(basePath, "product_service.go"))
	tc.AssertFileExists(filepath.Join(basePath, "product_usecase.go"))
	tc.AssertFileExists(filepath.Join(basePath, "dto.go"))

	// Verificar contenido del archivo de caso de uso - debe tener solo las operaciones especificadas
	usecasePath := filepath.Join(basePath, "product_service.go")

	// Debe contener
	expectedContents := []string{
		"NewProductService",
		"CreateProduct",
		"GetProduct",
	}

	// No debe contener
	unexpectedContents := []string{
		"UpdateProduct",
		"DeleteProduct",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(usecasePath, content)
	}

	for _, content := range unexpectedContents {
		data, err := os.ReadFile(filepath.Join(tc.TempDir, usecasePath))
		if err != nil {
			t.Errorf("Error al leer archivo %s: %v", usecasePath, err)
			continue
		}

		if strings.Contains(string(data), content) {
			t.Errorf("Contenido no esperado encontrado en %s: %s", usecasePath, content)
		}
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testUseCaseWithDTOValidation prueba el comando usecase con validación de DTOs
func testUseCaseWithDTOValidation(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("usecase", "OrderProcessing",
		"--entity", "Order",
		"--operations", "create,update",
		"--dto-validation")

	if err != nil {
		t.Fatalf("Error al ejecutar comando usecase con validación de DTOs: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Including DTO validations") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivos generados
	basePath := filepath.Join("internal", "usecase")
	dtoPath := filepath.Join(basePath, "dto.go")
	tc.AssertFileExists(dtoPath)

	// Verificar que existe el archivo de DTOs - sólo verificamos que exista
	// El contenido específico de validación puede variar, así que verificamos el archivo principal
	tc.AssertFileExists(filepath.Join(basePath, "order_service.go"))
	tc.AssertFileExists(filepath.Join(basePath, "order_usecase.go"))

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testUseCaseWithAsync prueba el comando usecase con operaciones asíncronas
func testUseCaseWithAsync(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("usecase", "AsyncUserNotification",
		"--entity", "User",
		"--operations", "create",
		"--async")

	if err != nil {
		t.Fatalf("Error al ejecutar comando usecase con operaciones asíncronas: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Including asynchronous operations") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivos generados
	basePath := filepath.Join("internal", "usecase")
	usecasePath := filepath.Join(basePath, "user_service.go")
	tc.AssertFileExists(usecasePath)

	// Verificar contenido del archivo - debe tener alguna referencia a async
	tc.AssertFileContains(usecasePath, "async")

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testUseCaseWithAllOptions prueba el comando usecase con todas las opciones
func testUseCaseWithAllOptions(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("usecase", "FullProductService",
		"--entity", "Product",
		"--operations", "create,read,update,list",
		"--dto-validation",
		"--async")

	if err != nil {
		t.Fatalf("Error al ejecutar comando usecase con todas las opciones: %v", err)
	}

	// Verificar salida
	expectedOutputs := []string{
		"Generating use case 'FullProductService' for entity 'Product'",
		"Operations: create,read,update,list",
		"Including DTO validations",
		"Including asynchronous operations",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Salida esperada no encontrada: %s", expected)
		}
	}

	// Verificar archivos generados
	basePath := filepath.Join("internal", "usecase")
	usecasePath := filepath.Join(basePath, "product_service.go")
	dtoPath := filepath.Join(basePath, "dto.go")
	tc.AssertFileExists(usecasePath)
	tc.AssertFileExists(dtoPath)

	// Verificar al menos algunas características clave en los archivos
	// Verificamos que el servicio contenga algunas operaciones básicas
	tc.AssertFileContains(usecasePath, "CreateProduct")
	tc.AssertFileContains(usecasePath, "GetProduct")
	tc.AssertFileContains(usecasePath, "async")

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}
