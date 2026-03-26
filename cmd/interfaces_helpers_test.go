package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateGRPCRequestResponseInterfaces(t *testing.T) {
	t.Parallel()

	t.Run("Product entity", func(t *testing.T) {
		t.Parallel()
		var sb strings.Builder
		generateGRPCRequestResponseInterfaces(&sb, "Product")
		result := sb.String()

		assert.Contains(t, result, "type CreateProductRequest interface")
		assert.Contains(t, result, "type CreateProductResponse interface")
		assert.Contains(t, result, "type GetProductRequest interface")
		assert.Contains(t, result, "type ProductResponse interface")
		assert.Contains(t, result, "type UpdateProductRequest interface")
		assert.Contains(t, result, "type UpdateProductResponse interface")
		assert.Contains(t, result, "type DeleteProductRequest interface")
		assert.Contains(t, result, "type DeleteProductResponse interface")
		assert.Contains(t, result, "type ListProductsRequest interface")
		assert.Contains(t, result, "type ListProductsResponse interface")
		assert.Contains(t, result, "GetName() string")
		assert.Contains(t, result, "GetEmail() string")
		assert.Contains(t, result, "GetId() int32")
		assert.Contains(t, result, "GetProduct() *Product")
		assert.Contains(t, result, "GetMessage() string")
		assert.Contains(t, result, "GetProducts() []*Product")
		assert.Contains(t, result, "GetTotal() int32")
	})

	t.Run("User entity", func(t *testing.T) {
		t.Parallel()
		var sb strings.Builder
		generateGRPCRequestResponseInterfaces(&sb, "User")
		result := sb.String()

		assert.Contains(t, result, "type CreateUserRequest interface")
		assert.Contains(t, result, "GetUser() *User")
		assert.Contains(t, result, "GetUsers() []*User")
	})
}

func TestDryRunGenerateInterfaces(t *testing.T) {
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { require.NoError(t, os.Chdir(origDir)) }()

	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))
	sm := &SafetyManager{DryRun: true}

	t.Run("generate usecase interface", func(t *testing.T) {
		generateUseCaseInterfaceFile(tmpDir, "Product", sm)
		// DryRun — no files created, just exercises code paths
	})

	t.Run("generate handler interface", func(t *testing.T) {
		generateHandlerInterfaceFile(tmpDir, "Product", sm)
	})

	t.Run("generate repository interface with fields", func(t *testing.T) {
		generateRepositoryInterfaceFileWithFields(tmpDir, "Product", "Name:string,Price:float64", sm)
	})

	t.Run("generate all interfaces", func(t *testing.T) {
		generateInterfaces("Product", true, true, true, sm)
	})
}
