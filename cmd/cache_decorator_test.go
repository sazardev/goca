package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCacheDecorator_Basic(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
	require.NoError(t, os.MkdirAll(filepath.Join(DirInternal, DirRepository), 0755))

	sm := NewSafetyManager(false, true, false)
	fields := []Field{
		{Name: "Name", Type: "string"},
		{Name: "Price", Type: "float64"},
	}

	generateCacheDecorator("Product", fields, sm)

	content, err := os.ReadFile(filepath.Join(DirInternal, DirRepository, "cached_product_repository.go"))
	require.NoError(t, err)
	src := string(content)

	assert.Contains(t, src, "type CachedProductRepository struct")
	assert.Contains(t, src, "inner    ProductRepository")
	assert.Contains(t, src, "cache    *redis.Client")
	assert.Contains(t, src, "func NewCachedProductRepository(inner ProductRepository, cache *redis.Client, ttl time.Duration)")
	assert.Contains(t, src, "func (r *CachedProductRepository) Save(product *domain.Product) error")
	assert.Contains(t, src, "func (r *CachedProductRepository) FindByID(id int) (*domain.Product, error)")
	assert.Contains(t, src, "func (r *CachedProductRepository) FindAll() ([]domain.Product, error)")
	assert.Contains(t, src, "func (r *CachedProductRepository) Update(product *domain.Product) error")
	assert.Contains(t, src, "func (r *CachedProductRepository) Delete(id int) error")
}

func TestGenerateCacheDecorator_CacheKeyHelpers(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
	require.NoError(t, os.MkdirAll(filepath.Join(DirInternal, DirRepository), 0755))

	sm := NewSafetyManager(false, true, false)
	generateCacheDecorator("Order", nil, sm)

	content, err := os.ReadFile(filepath.Join(DirInternal, DirRepository, "cached_order_repository.go"))
	require.NoError(t, err)
	src := string(content)

	assert.Contains(t, src, `return fmt.Sprintf("order:%d", id)`)
	assert.Contains(t, src, `return "order:list"`)
}

func TestGenerateCacheDecorator_SearchMethods(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
	require.NoError(t, os.MkdirAll(filepath.Join(DirInternal, DirRepository), 0755))

	sm := NewSafetyManager(false, true, false)
	fields := []Field{
		{Name: "Email", Type: "string"},
		{Name: "Age", Type: "int"},
	}

	generateCacheDecorator("User", fields, sm)

	content, err := os.ReadFile(filepath.Join(DirInternal, DirRepository, "cached_user_repository.go"))
	require.NoError(t, err)
	src := string(content)

	// Search methods should delegate to inner (no caching)
	assert.Contains(t, src, "func (r *CachedUserRepository) FindByEmail")
	assert.Contains(t, src, "return r.inner.FindByEmail")
}

func TestGenerateCacheDecorator_NoFields(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
	require.NoError(t, os.MkdirAll(filepath.Join(DirInternal, DirRepository), 0755))

	sm := NewSafetyManager(false, true, false)
	generateCacheDecorator("Item", nil, sm)

	content, err := os.ReadFile(filepath.Join(DirInternal, DirRepository, "cached_item_repository.go"))
	require.NoError(t, err)
	src := string(content)

	// Default FindByEmail delegate when no fields
	assert.Contains(t, src, "func (r *CachedItemRepository) FindByEmail")
}

func TestGenerateCacheDecorator_DryRun(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
	require.NoError(t, os.MkdirAll(filepath.Join(DirInternal, DirRepository), 0755))

	sm := NewSafetyManager(true, false, false)
	generateCacheDecorator("Product", nil, sm)

	// File should NOT be created in dry-run
	_, err := os.ReadFile(filepath.Join(DirInternal, DirRepository, "cached_product_repository.go"))
	assert.True(t, os.IsNotExist(err))

	// But should be in pending files
	pending := sm.GetPendingFiles()
	found := false
	for _, p := range pending {
		if strings.Contains(p.Path, "cached_product_repository") {
			found = true
			break
		}
	}
	assert.True(t, found, "cache decorator should be in pending files")
}

func TestGenerateCacheDecorator_Imports(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	require.NoError(t, os.WriteFile("go.mod", []byte("module myapp\n\ngo 1.21\n"), 0644))
	require.NoError(t, os.MkdirAll(filepath.Join(DirInternal, DirRepository), 0755))

	sm := NewSafetyManager(false, true, false)
	generateCacheDecorator("Product", nil, sm)

	content, err := os.ReadFile(filepath.Join(DirInternal, DirRepository, "cached_product_repository.go"))
	require.NoError(t, err)
	src := string(content)

	assert.Contains(t, src, `"context"`)
	assert.Contains(t, src, `"encoding/json"`)
	assert.Contains(t, src, `"time"`)
	assert.Contains(t, src, `"myapp/internal/domain"`)
	assert.Contains(t, src, `"github.com/redis/go-redis/v9"`)
}

func TestGenerateCacheDecorator_Invalidation(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))
	require.NoError(t, os.MkdirAll(filepath.Join(DirInternal, DirRepository), 0755))

	sm := NewSafetyManager(false, true, false)
	generateCacheDecorator("Product", nil, sm)

	content, err := os.ReadFile(filepath.Join(DirInternal, DirRepository, "cached_product_repository.go"))
	require.NoError(t, err)
	src := string(content)

	// Save should invalidate list cache
	assert.Contains(t, src, "r.cache.Del(r.ctx, r.listCacheKey())")
	// Update should invalidate both item and list cache
	assert.Contains(t, src, "r.cache.Del(r.ctx, r.cacheKey(int(product.ID)), r.listCacheKey())")
	// Delete should invalidate both item and list cache
	assert.Contains(t, src, "r.cache.Del(r.ctx, r.cacheKey(id), r.listCacheKey())")
}

func TestGenerateCacheSearchMethodDelegate(t *testing.T) {
	t.Parallel()

	var b strings.Builder
	m := SearchMethod{
		MethodName: "FindByEmail",
		FieldName:  "Email",
		FieldType:  "string",
		ReturnType: "(*domain.User, error)",
		IsUnique:   true,
	}

	generateCacheSearchMethodDelegate(&b, "User", m)
	output := b.String()

	assert.Contains(t, output, "func (r *CachedUserRepository) FindByEmail(email string) (*domain.User, error)")
	assert.Contains(t, output, "return r.inner.FindByEmail(email)")
}
