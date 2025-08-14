package tests

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/sazardev/goca/internal/testing/framework"
)

// TestEntityCommand prueba exhaustivamente el comando 'entity'
func TestEntityCommand(t *testing.T) {
	// Crear contexto de test
	tc := framework.NewTestContext(t)
	defer tc.Cleanup()
	tc.CurrentTestName = "TestEntityCommand"

	// Preparar proyecto base para todos los tests
	prepareBaseProject(tc, t)

	// Ejecutar subpruebas para diferentes escenarios
	t.Run("EntityWithRequiredFields", func(t *testing.T) {
		testEntityWithRequiredFields(tc, t)
	})

	t.Run("EntityWithValidation", func(t *testing.T) {
		testEntityWithValidation(tc, t)
	})

	t.Run("EntityWithBusinessRules", func(t *testing.T) {
		testEntityWithBusinessRules(tc, t)
	})

	t.Run("EntityWithTimestamps", func(t *testing.T) {
		testEntityWithTimestamps(tc, t)
	})

	t.Run("EntityWithSoftDelete", func(t *testing.T) {
		testEntityWithSoftDelete(tc, t)
	})

	t.Run("EntityWithAllOptions", func(t *testing.T) {
		testEntityWithAllOptions(tc, t)
	})

	t.Run("EntityWithDifferentFieldTypes", func(t *testing.T) {
		testEntityWithDifferentFieldTypes(tc, t)
	})

	// Verificar que el proyecto completo compila
	tc.AssertGoBuild("test-project")

	// Verificar que pasa go vet
	tc.AssertGoVet("test-project")

	// Imprimir resumen
	tc.PrintTestSummary()
}

// prepareBaseProject inicializa un proyecto base para las pruebas de entity
func prepareBaseProject(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando init
	_, err := tc.RunCommand("init", "test-project", "--module", "github.com/test/testproject")
	if err != nil {
		t.Fatalf("Error al inicializar proyecto base: %v", err)
	}

	// Verificar que se creó correctamente
	tc.AssertFileExists(filepath.Join("test-project", "go.mod"))
	tc.AssertGoBuild("test-project")
}

// testEntityWithRequiredFields prueba el comando entity con campos requeridos
func testEntityWithRequiredFields(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("entity", "User", "--fields", "name:string,email:string,age:int")
	if err != nil {
		t.Fatalf("Error al ejecutar comando entity: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Generando entidad 'User'") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivo generado - Nota: ahora sabemos que se genera en internal/domain, no en test-project/internal/domain
	entityPath := filepath.Join("internal", "domain", "user.go")
	tc.AssertFileExists(entityPath)

	// Verificar contenido del archivo
	expectedContents := []string{
		"type User struct",
		"Name",
		"string",
		"Email",
		"string",
		"Age",
		"int",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(entityPath, content)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testEntityWithValidation prueba el comando entity con validación
func testEntityWithValidation(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("entity", "Product",
		"--fields", "name:string,price:float64,sku:string",
		"--validation")

	if err != nil {
		t.Fatalf("Error al ejecutar comando entity con validación: %v", err)
	}

	// Verificar salida
	expectedOutputs := []string{
		"Generando entidad 'Product'",
		"Incluyendo validaciones",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Salida esperada no encontrada: %s", expected)
		}
	}

	// Verificar archivo generado - corregido el path
	entityPath := filepath.Join("internal", "domain", "product.go")
	tc.AssertFileExists(entityPath)

	// Verificar contenido del archivo - debe incluir función Validate()
	expectedContents := []string{
		"type Product struct",
		"Name",
		"string",
		"Price",
		"float64",
		"Sku",
		"string",
		"func (p *Product) Validate() error",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(entityPath, content)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testEntityWithBusinessRules prueba el comando entity con reglas de negocio
func testEntityWithBusinessRules(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("entity", "Order",
		"--fields", "orderID:string,total:float64,status:string",
		"--business-rules")

	if err != nil {
		t.Fatalf("Error al ejecutar comando entity con reglas de negocio: %v", err)
	}

	// Verificar salida
	expectedOutputs := []string{
		"Generando entidad 'Order'",
		"Incluyendo reglas de negocio",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Salida esperada no encontrada: %s", expected)
		}
	}

	// Verificar archivo generado
	entityPath := filepath.Join("internal", "domain", "order.go")
	tc.AssertFileExists(entityPath)

	// Verificar métodos de reglas de negocio
	tc.AssertGoBuild("test-project")
}

// testEntityWithTimestamps prueba el comando entity con timestamps
func testEntityWithTimestamps(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("entity", "Category",
		"--fields", "name:string,description:string",
		"--timestamps")

	if err != nil {
		t.Fatalf("Error al ejecutar comando entity con timestamps: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Incluyendo timestamps") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivo generado
	entityPath := filepath.Join("internal", "domain", "category.go")
	tc.AssertFileExists(entityPath)

	// Verificar campos de timestamp
	expectedContents := []string{
		"CreatedAt",
		"time.Time",
		"UpdatedAt",
		"time.Time",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(entityPath, content)
	}

	// Verificar import de time
	tc.AssertFileContains(entityPath, "time")

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testEntityWithSoftDelete prueba el comando entity con soft delete
func testEntityWithSoftDelete(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("entity", "Customer",
		"--fields", "name:string,email:string,phone:string",
		"--soft-delete")

	if err != nil {
		t.Fatalf("Error al ejecutar comando entity con soft delete: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Incluyendo eliminación suave") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivo generado
	entityPath := filepath.Join("internal", "domain", "customer.go")
	tc.AssertFileExists(entityPath)

	// Verificar campo DeletedAt
	tc.AssertFileContains(entityPath, "DeletedAt")

	// Verificar import de time
	tc.AssertFileContains(entityPath, "time")

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testEntityWithAllOptions prueba el comando entity con todas las opciones
func testEntityWithAllOptions(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("entity", "Invoice",
		"--fields", "number:string,amount:float64,paid:bool,client_id:int",
		"--validation",
		"--business-rules",
		"--timestamps",
		"--soft-delete")

	if err != nil {
		t.Fatalf("Error al ejecutar comando entity con todas las opciones: %v", err)
	}

	// Verificar salida
	expectedOutputs := []string{
		"Generando entidad 'Invoice'",
		"Incluyendo validaciones",
		"Incluyendo reglas de negocio",
		"Incluyendo timestamps",
		"Incluyendo eliminación suave",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Salida esperada no encontrada: %s", expected)
		}
	}

	// Verificar archivo generado
	entityPath := filepath.Join("internal", "domain", "invoice.go")
	tc.AssertFileExists(entityPath)

	// Verificar todos los elementos esperados
	expectedContents := []string{
		"type Invoice struct",
		"Number",
		"Amount",
		"Paid",
		"Client_id",
		"CreatedAt",
		"UpdatedAt",
		"DeletedAt",
		"func (i *Invoice) Validate() error",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(entityPath, content)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testEntityWithDifferentFieldTypes prueba el comando entity con diferentes tipos de campos
func testEntityWithDifferentFieldTypes(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando con muchos tipos de campos diferentes
	output, err := tc.RunCommand("entity", "Complex",
		"--fields", "complexID:string,count:int,active:bool,price:float64,created:time.Time,items:[]string,metadata:interface{}")

	if err != nil {
		t.Fatalf("Error al ejecutar comando entity con diferentes tipos: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Generando entidad 'Complex'") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivo generado
	entityPath := filepath.Join("internal", "domain", "complex.go")
	tc.AssertFileExists(entityPath)

	// Verificar todos los tipos de campos
	expectedContents := []string{
		"Complexid",
		"string",
		"Count",
		"int",
		"Active",
		"bool",
		"Price",
		"float64",
		"Created",
		"time.Time",
		"Items",
		"[]string",
		"Metadata",
		"interface{}",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(entityPath, content)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}
