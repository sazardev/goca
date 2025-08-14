package tests

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/sazardev/goca/internal/testing/framework"
)

// TestHandlerCommand prueba exhaustivamente el comando 'handler'
func TestHandlerCommand(t *testing.T) {
	// Crear contexto de test
	tc := framework.NewTestContext(t)
	defer tc.Cleanup()
	tc.CurrentTestName = "TestHandlerCommand"

	// Preparar proyecto base para todos los tests
	prepareFullProjectBase(tc, t)

	// Ejecutar subpruebas para diferentes escenarios
	t.Run("HTTPHandler", func(t *testing.T) {
		testHTTPHandler(tc, t)
	})

	t.Run("GRPCHandler", func(t *testing.T) {
		testGRPCHandler(tc, t)
	})

	t.Run("CLIHandler", func(t *testing.T) {
		testCLIHandler(tc, t)
	})

	t.Run("WorkerHandler", func(t *testing.T) {
		testWorkerHandler(tc, t)
	})

	t.Run("HandlerWithMiddleware", func(t *testing.T) {
		testHandlerWithMiddleware(tc, t)
	})

	t.Run("HandlerWithValidation", func(t *testing.T) {
		testHandlerWithValidation(tc, t)
	})

	t.Run("HandlerWithSwagger", func(t *testing.T) {
		testHandlerWithSwagger(tc, t)
	})

	t.Run("HandlerWithAllOptions", func(t *testing.T) {
		testHandlerWithAllOptions(tc, t)
	})

	// Verificar que el proyecto completo compila
	tc.AssertGoBuild("test-project")

	// Verificar que pasa go vet
	tc.AssertGoVet("test-project")

	// Imprimir resumen
	tc.PrintTestSummary()
}

// prepareFullProjectBase inicializa un proyecto completo para las pruebas de handler
func prepareFullProjectBase(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando init
	_, err := tc.RunCommand("init", "test-project", "--module", "github.com/test/testproject", "--api", "rest", "--database", "postgres")
	if err != nil {
		t.Fatalf("Error al inicializar proyecto base: %v", err)
	}

	// Crear entidades y casos de uso para las pruebas
	entities := []struct {
		name   string
		fields string
	}{
		{"User", "name:string,email:string,password:string,age:int"},
		{"Product", "name:string,price:float64,sku:string"},
		{"Order", "orderID:string,total:float64,status:string"},
		{"Customer", "name:string,email:string,phone:string"},
	}

	for _, entity := range entities {
		_, err := tc.RunCommand("entity", entity.name, "--fields", entity.fields)
		if err != nil {
			t.Fatalf("Error al crear entidad %s: %v", entity.name, err)
		}

		_, err = tc.RunCommand("usecase", entity.name+"Service", "--entity", entity.name)
		if err != nil {
			t.Fatalf("Error al crear caso de uso para %s: %v", entity.name, err)
		}
	}

	// Verificar que se creó correctamente
	tc.AssertGoBuild("test-project")
}

// testHTTPHandler prueba el comando handler con tipo HTTP
func testHTTPHandler(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("handler", "User", "--type", "http")
	if err != nil {
		t.Fatalf("Error al ejecutar comando handler HTTP: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Generando handler 'http' para entidad 'User'") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivos generados
	basePath := filepath.Join("test-project", "internal", "handler", "http")
	handlerPath := filepath.Join(basePath, "user_handler.go")
	routerPath := filepath.Join(basePath, "router.go")
	tc.AssertFileExists(handlerPath)
	tc.AssertFileExists(routerPath)

	// Verificar contenido del archivo del handler
	expectedContents := []string{
		"type UserHandler struct {",
		"usecase usecase.UserServiceUseCase",
		"func NewUserHandler(",
		"r.POST(\"/users\"",
		"r.GET(\"/users/",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(handlerPath, content)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testGRPCHandler prueba el comando handler con tipo gRPC
func testGRPCHandler(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("handler", "Product", "--type", "grpc")
	if err != nil {
		t.Fatalf("Error al ejecutar comando handler gRPC: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Generando handler 'grpc' para entidad 'Product'") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivos generados
	basePath := filepath.Join("test-project", "internal", "handler", "grpc")
	handlerPath := filepath.Join(basePath, "product_handler.go")
	protoPath := filepath.Join("test-project", "proto", "product.proto")
	tc.AssertFileExists(handlerPath)
	tc.AssertFileExists(protoPath)

	// Verificar contenido del archivo del handler
	expectedContents := []string{
		"type ProductHandler struct {",
		"usecase usecase.ProductServiceUseCase",
		"func NewProductHandler(",
		"RegisterProductServiceServer",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(handlerPath, content)
	}

	// Verificar contenido del archivo .proto
	expectedProtoContents := []string{
		"syntax = \"proto3\";",
		"package proto;",
		"service ProductService {",
		"message Product {",
		"string name",
		"double price",
	}

	for _, content := range expectedProtoContents {
		tc.AssertFileContains(protoPath, content)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testCLIHandler prueba el comando handler con tipo CLI
func testCLIHandler(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("handler", "Order", "--type", "cli")
	if err != nil {
		t.Fatalf("Error al ejecutar comando handler CLI: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Generando handler 'cli' para entidad 'Order'") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivos generados
	basePath := filepath.Join("test-project", "internal", "handler", "cli")
	handlerPath := filepath.Join(basePath, "order_cmd.go")
	tc.AssertFileExists(handlerPath)

	// Verificar contenido del archivo del handler
	expectedContents := []string{
		"var orderCmd = &cobra.Command{",
		"UseCase usecase.OrderServiceUseCase",
		"func NewOrderCommands(",
		"createOrderCmd",
		"getOrderCmd",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(handlerPath, content)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testWorkerHandler prueba el comando handler con tipo worker
func testWorkerHandler(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("handler", "Customer", "--type", "worker")
	if err != nil {
		t.Fatalf("Error al ejecutar comando handler worker: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Generando handler 'worker' para entidad 'Customer'") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivos generados
	basePath := filepath.Join("test-project", "internal", "handler", "worker")
	handlerPath := filepath.Join(basePath, "customer_worker.go")
	tc.AssertFileExists(handlerPath)

	// Verificar contenido del archivo del worker
	expectedContents := []string{
		"type CustomerWorker struct {",
		"useCase usecase.CustomerServiceUseCase",
		"func NewCustomerWorker(",
		"func (w *CustomerWorker) Start() error {",
		"go w.process",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(handlerPath, content)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testHandlerWithMiddleware prueba el comando handler con middleware
func testHandlerWithMiddleware(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("handler", "User", "--type", "http", "--middleware")
	if err != nil {
		t.Fatalf("Error al ejecutar comando handler con middleware: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Incluyendo middleware") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivos generados
	middlewarePath := filepath.Join("test-project", "internal", "handler", "http", "middleware.go")
	tc.AssertFileExists(middlewarePath)

	// Verificar contenido del archivo de middleware
	expectedContents := []string{
		"package http",
		"func AuthMiddleware() gin.HandlerFunc {",
		"func LoggerMiddleware() gin.HandlerFunc {",
		"func CORSMiddleware() gin.HandlerFunc {",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(middlewarePath, content)
	}

	// Verificar que el handler usa middleware
	handlerPath := filepath.Join("test-project", "internal", "handler", "http", "user_handler.go")
	tc.AssertFileContains(handlerPath, "AuthMiddleware()")

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testHandlerWithValidation prueba el comando handler con validación
func testHandlerWithValidation(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("handler", "Product", "--type", "http", "--validation")
	if err != nil {
		t.Fatalf("Error al ejecutar comando handler con validación: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Incluyendo validación") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivos generados
	handlerPath := filepath.Join("test-project", "internal", "handler", "http", "product_handler.go")
	tc.AssertFileExists(handlerPath)

	// Verificar contenido del archivo de validación
	expectedContents := []string{
		"if err := c.ShouldBindJSON(&request); err != nil {",
		"c.JSON(http.StatusBadRequest, gin.H{\"error\": err.Error()})",
		"return",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(handlerPath, content)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testHandlerWithSwagger prueba el comando handler con documentación Swagger
func testHandlerWithSwagger(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("handler", "Order", "--type", "http", "--swagger")
	if err != nil {
		t.Fatalf("Error al ejecutar comando handler con Swagger: %v", err)
	}

	// Verificar salida
	if !strings.Contains(output, "Incluyendo documentación Swagger") {
		t.Errorf("Salida esperada no encontrada: %s", output)
	}

	// Verificar archivos generados
	handlerPath := filepath.Join("test-project", "internal", "handler", "http", "order_handler.go")
	tc.AssertFileExists(handlerPath)

	// Verificar contenido de la documentación Swagger
	expectedContents := []string{
		"// @Summary",
		"// @Description",
		"// @Accept",
		"// @Produce",
		"// @Param",
		"// @Success",
		"// @Failure",
		"// @Router",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(handlerPath, content)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}

// testHandlerWithAllOptions prueba el comando handler con todas las opciones
func testHandlerWithAllOptions(tc *framework.TestContext, t *testing.T) {
	// Ejecutar comando
	output, err := tc.RunCommand("handler", "Customer", "--type", "http", "--middleware", "--validation", "--swagger")
	if err != nil {
		t.Fatalf("Error al ejecutar comando handler con todas las opciones: %v", err)
	}

	// Verificar salida
	expectedOutputs := []string{
		"Generando handler 'http' para entidad 'Customer'",
		"Incluyendo middleware",
		"Incluyendo validación",
		"Incluyendo documentación Swagger",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Salida esperada no encontrada: %s", expected)
		}
	}

	// Verificar archivos generados
	handlerPath := filepath.Join("test-project", "internal", "handler", "http", "customer_handler.go")
	tc.AssertFileExists(handlerPath)

	// Verificar todos los elementos esperados
	expectedContents := []string{
		"// @Summary",
		"AuthMiddleware()",
		"if err := c.ShouldBindJSON(&request); err != nil {",
		"c.JSON(http.StatusBadRequest, gin.H{\"error\": err.Error()})",
	}

	for _, content := range expectedContents {
		tc.AssertFileContains(handlerPath, content)
	}

	// Verificar que el proyecto compila
	tc.AssertGoBuild("test-project")
}
