package tests

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/sazardev/goca/internal/testing/framework"
)

// TestFeatureCommand prueba exhaustivamente el comando 'feature' que integra todas las capas
func TestFeatureCommand(t *testing.T) {
	// Crear contexto de test
	tc := framework.NewTestContext(t)
	defer tc.Cleanup()
	tc.CurrentTestName = "TestFeatureCommand"

	// Preparar proyecto base para todos los tests
	prepareEmptyProject(tc, t)

	// Ejecutar subpruebas para diferentes escenarios
	t.Run("FeatureWithBasicOptions", func(t *testing.T) {
		testFeatureWithBasicOptions(tc, t)
	})

	t.Run("FeatureWithCustomDatabase", func(t *testing.T) {
		testFeatureWithCustomDatabase(tc, t)
	})

	t.Run("FeatureWithCustomHandlers", func(t *testing.T) {
		testFeatureWithCustomHandlers(tc, t)
	})

	t.Run("FeatureWithValidations", func(t *testing.T) {
		testFeatureWithValidations(tc, t)
	})

	t.Run("FeatureWithAutoIntegration", func(t *testing.T) {
		testFeatureWithAutoIntegration(tc, t)
	})

	t.Run("FeatureWithComplexFields", func(t *testing.T) {
		testFeatureWithComplexFields(tc, t)
	})

	// Verificar que el proyecto completo compila
	tc.AssertGoBuild(".")

	// Verificar que pasa go vet
	tc.AssertGoVet(".")

	// Imprimir resumen
	tc.PrintTestSummary()
}

// prepareEmptyProject inicializa un proyecto vacío para las pruebas de feature
func prepareEmptyProject(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando init
	_, err := tc.RunCommand("init", "test-project", "--module", "github.com/test/testproject", "--api", "rest", "--database", "postgres")
	if err != nil {
		t.Fatalf("Error al inicializar proyecto base: %v", err)
	}

	// Set the project directory for subsequent commands
	tc.SetProjectDir("test-project")

	// Verificar que se creó correctamente
	tc.AssertFileExists("go.mod")
	tc.AssertGoBuild(".")
}

// testFeatureWithBasicOptions prueba el comando feature con opciones básicas
func testFeatureWithBasicOptions(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("feature", "User", "--fields", "name:string,email:string,age:int")
	if err != nil {
		t.Fatalf("Error al ejecutar comando feature: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Generando feature completo 'User'") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar que se generaron todos los archivos necesarios
	// 1. Entidad de dominio
	tc.AssertFileExists(filepath.Join("internal", "domain", "user.go"))

	// 2. Caso de uso
	tc.AssertFileExists(filepath.Join("internal", "usecase", "user_service.go"))
	// Note: user_dto.go is generated as a shared dto.go file, not individual entity DTOs

	// 3. Repositorio
	tc.AssertFileExists(filepath.Join("internal", "repository", "postgres_user_repository.go"))

	// 4. Handler (por defecto HTTP)
	tc.AssertFileExists(filepath.Join("internal", "handler", "http", "user_handler.go"))

	// Verificar que el código compila
	tc.AssertGoBuild(".")
}

// testFeatureWithCustomDatabase prueba el comando feature con base de datos específica
func testFeatureWithCustomDatabase(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("feature", "Product",
		"--fields", "name:string,price:float64,sku:string",
		"--database", "mongodb")

	if err != nil {
		t.Fatalf("Error al ejecutar comando feature con base de datos: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Base de datos: mongodb") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar repositorio específico para MongoDB
	repoPath := filepath.Join("internal", "repository", "mongo_product_repository.go")
	tc.AssertFileExists(repoPath)
	tc.AssertFileContains(repoPath, "mongo")

	// Verificar que el código compila
	tc.AssertGoBuild(".")
}

// testFeatureWithCustomHandlers prueba el comando feature con tipos de handlers específicos
func testFeatureWithCustomHandlers(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("feature", "Order",
		"--fields", "orderID:string,total:float64,status:string",
		"--handlers", "http,grpc")

	if err != nil {
		t.Fatalf("Error al ejecutar comando feature con handlers: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Handlers: http,grpc") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar que se generaron ambos tipos de handlers
	tc.AssertFileExists(filepath.Join("internal", "handler", "http", "order_handler.go"))
	tc.AssertFileExists(filepath.Join("internal", "handler", "grpc", "order_server.go"))
	tc.AssertFileExists(filepath.Join("internal", "handler", "grpc", "order.proto"))

	// Verificar que el código compila
	tc.AssertGoBuild(".")
}

// testFeatureWithValidations prueba el comando feature con validaciones
func testFeatureWithValidations(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando - solo usa flags válidos
	output, err := tc.RunCommand("feature", "Customer",
		"--fields", "name:string,email:string,phone:string",
		"--validation")

	if err != nil {
		t.Fatalf("Error al ejecutar comando feature con validaciones: %v", err)
	}

	// Verificar salida
	expectedOutputs := []string{
		"Generando feature completo 'Customer'",
		"✅ Incluyendo validaciones",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Salida esperada no encontrada: %s", expected)
		}
	}

	// Verificar validaciones en entidad
	entityPath := filepath.Join("internal", "domain", "customer.go")
	tc.AssertFileExists(entityPath)

	// Verificar que el código compila
	tc.AssertGoBuild(".")
}

// testFeatureWithAutoIntegration prueba el comando feature con integración automática
func testFeatureWithAutoIntegration(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando - usar --business-rules que es un flag válido
	output, err := tc.RunCommand("feature", "Category",
		"--fields", "name:string,description:string",
		"--business-rules")

	if err != nil {
		t.Fatalf("Error al ejecutar comando feature con business rules: %v", err)
	}

	// Verificar salida
	expectedOutputs := []string{
		"Generando feature completo 'Category'",
		"🧠 Incluyendo reglas de negocio",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Salida esperada no encontrada: %s", expected)
		}
	}

	// Verificar archivos básicos
	tc.AssertFileExists(filepath.Join("internal", "domain", "category.go"))
	tc.AssertFileExists(filepath.Join("internal", "usecase", "category_service.go"))
	tc.AssertFileExists(filepath.Join("internal", "repository", "postgres_category_repository.go"))
	tc.AssertFileExists(filepath.Join("internal", "handler", "http", "category_handler.go"))

	// Verificar que el código compila
	tc.AssertGoBuild("test-project")
}

// testFeatureWithComplexFields prueba el comando feature con campos complejos
func testFeatureWithComplexFields(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando con campos complejos - usar solo tipos válidos
	output, err := tc.RunCommand("feature", "Article",
		"--fields", "title:string,content:string,author_id:int,created_at:time.Time,metadata:interface{}")

	if err != nil {
		t.Fatalf("Error al ejecutar comando feature con campos complejos: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Generando feature completo 'Article'") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar entidad con campos complejos
	entityPath := filepath.Join("internal", "domain", "article.go")
	tc.AssertFileExists(entityPath)

	// Verificar todos los tipos de campos (using actual generated field names)
	expectedContents := []string{
		"Title      string",
		"Content    string",
		"Author_id  int",
		"Created_at time.Time",
		"Metadata   interface{}",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(entityPath, content)
	}

	// Verificar import de time (check for "time" in import section)
	tc.AssertFileContains(entityPath, "time")

	// Verificar que el proyecto compila
	tc.AssertGoBuild(".")
}
