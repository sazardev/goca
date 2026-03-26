package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDryRunGenerators_Batch3 tests additional generators that require os.Chdir.
func TestDryRunGenerators_Batch3(t *testing.T) {
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { require.NoError(t, os.Chdir(origDir)) }()

	cleanup := ensureTestUI(t)
	defer cleanup()

	// repository_other_db.go
	t.Run("generatePostgresJSONRepository", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		generatePostgresJSONRepository(dir, "Product", false, false, sm)
	})

	t.Run("generateSQLServerRepository", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		generateSQLServerRepository(dir, "Product", false, false, sm)
	})

	t.Run("generateElasticsearchRepository", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		generateElasticsearchRepository(dir, "Product", false, false, sm)
	})

	t.Run("generateDynamoDBRepository", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		generateDynamoDBRepository(dir, "Product", false, false, sm)
	})

	t.Run("generateSQLiteRepository", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		generateSQLiteRepository(dir, "Product", true, true, sm)
	})

	// repository_fields.go
	t.Run("generateRepositoryInterfaceWithFields", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		fields := []Field{{Name: "Name", Type: "string"}, {Name: "Price", Type: "float64"}}
		generateRepositoryInterfaceWithFields(dir, "Product", fields, false, sm)
	})

	t.Run("generateRepositoryImplementationWithFields postgres", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		fields := []Field{{Name: "Name", Type: "string"}, {Name: "Price", Type: "float64"}}
		generateRepositoryImplementationWithFields(dir, "Product", "postgres", fields, false, false, sm)
	})

	t.Run("generateRepositoryImplementationWithFields mysql", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		fields := []Field{{Name: "Name", Type: "string"}}
		generateRepositoryImplementationWithFields(dir, "Product", "mysql", fields, true, true, sm)
	})

	t.Run("generateRepositoryImplementationWithFields mongodb", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		fields := []Field{{Name: "Name", Type: "string"}}
		generateRepositoryImplementationWithFields(dir, "Product", "mongodb", fields, false, false, sm)
	})

	// handler_other.go
	t.Run("generateGRPCHandler", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		generateGRPCHandler("Product", "lowercase", sm)
	})

	t.Run("generateCLIHandler", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		generateCLIHandler("Product", "lowercase", sm)
	})

	t.Run("generateWorkerHandler", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		generateWorkerHandler("Product", "lowercase", sm)
	})

	t.Run("generateSOAPHandler", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		generateSOAPHandler("Product", "lowercase", sm)
	})

	// test_integration.go
	t.Run("generateIntegrationTests", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		fields := []Field{{Name: "Name", Type: "string"}, {Name: "Price", Type: "float64"}}
		_ = generateIntegrationTests("Product", "postgres", true, true, fields, sm)
	})

	t.Run("generateIntegrationTests mysql no fixtures", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		fields := []Field{{Name: "Name", Type: "string"}}
		_ = generateIntegrationTests("Order", "mysql", false, false, fields, sm)
	})

	// integrate.go
	t.Run("integrateFeatures", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		sm := &SafetyManager{DryRun: true}
		integrateFeatures([]string{"Product", "Order"}, sm)
	})

	t.Run("detectExistingFeatures empty dir", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		features := detectExistingFeatures()
		assert.Empty(t, features)
	})

	t.Run("verifyIntegration", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		verifyIntegration([]string{"Product"})
	})
}

// Pure function tests

func TestGenerateHelpersContent(t *testing.T) {
	t.Parallel()

	t.Run("postgres with container", func(t *testing.T) {
		t.Parallel()
		result := generateHelpersContent("postgres", true, "Product")
		assert.Contains(t, result, "package")
		assert.Contains(t, result, "Product")
	})

	t.Run("mysql no container", func(t *testing.T) {
		t.Parallel()
		result := generateHelpersContent("mysql", false, "Order")
		assert.Contains(t, result, "package")
	})

	t.Run("mongodb", func(t *testing.T) {
		t.Parallel()
		result := generateHelpersContent("mongodb", false, "User")
		assert.Contains(t, result, "package")
	})
}

func TestReplaceHelperTODOs(t *testing.T) {
	t.Parallel()

	t.Run("replaces TODO", func(t *testing.T) {
		t.Parallel()
		content := "TODO: Add Product helper"
		result := replaceHelperTODOs(content, "Order")
		assert.NotEmpty(t, result)
	})

	t.Run("no TODOs", func(t *testing.T) {
		t.Parallel()
		content := "some plain content"
		result := replaceHelperTODOs(content, "Product")
		assert.Equal(t, content, result)
	})
}

func TestGenerateTransactionMethods(t *testing.T) {
	t.Parallel()
	var sb strings.Builder
	generateTransactionMethods(&sb, "Product", "ProductRepository")
	result := sb.String()
	assert.Contains(t, result, "SaveWithTx")
	assert.Contains(t, result, "UpdateWithTx")
	assert.Contains(t, result, "DeleteWithTx")
	assert.Contains(t, result, "ProductRepository")
}

func TestGenerateBasicMongoCRUDMethods(t *testing.T) {
	t.Parallel()
	var sb strings.Builder
	generateBasicMongoCRUDMethods(&sb, "Product", "ProductRepository")
	result := sb.String()
	assert.Contains(t, result, "InsertOne")
	assert.Contains(t, result, "FindOne")
	assert.Contains(t, result, "ProductRepository")
}

func TestGenerateMongoSearchMethodImplementation(t *testing.T) {
	t.Parallel()
	method := SearchMethod{
		MethodName: "FindByName",
		FieldName:  "Name",
		FieldType:  "string",
		IsUnique:   true,
	}
	result := generateMongoSearchMethodImplementation(method, "ProductRepository", "Product")
	assert.Contains(t, result, "FindByName")
	assert.Contains(t, result, "Product")
}

func TestGenerateUpdateMethod(t *testing.T) {
	t.Parallel()
	var sb strings.Builder
	generateUpdateMethod(&sb, "ProductService", "Product")
	result := sb.String()
	assert.Contains(t, result, "func (P *ProductService) UpdateProduct")
	assert.Contains(t, result, "UpdateProductInput")
	assert.Contains(t, result, "repo.FindByID")
}
