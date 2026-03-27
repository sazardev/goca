package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCachePackage_Basic(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	require.NoError(t, os.MkdirAll(filepath.Join(DirInternal, "cache"), 0755))

	sm := NewSafetyManager(false, true, false)
	err := generateCachePackage(sm)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(DirInternal, "cache", "redis.go"))
	require.NoError(t, err)
	src := string(content)

	assert.Contains(t, src, "package cache")
	assert.Contains(t, src, "func NewRedisClient() (*redis.Client, error)")
	assert.Contains(t, src, `os.Getenv("REDIS_URL")`)
	assert.Contains(t, src, `os.Getenv("REDIS_PASSWORD")`)
	assert.Contains(t, src, `os.Getenv("REDIS_DB")`)
	assert.Contains(t, src, `addr = "localhost:6379"`)
	assert.Contains(t, src, "client.Ping(context.Background())")
}

func TestGenerateCachePackage_DryRun(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	require.NoError(t, os.MkdirAll(filepath.Join(DirInternal, "cache"), 0755))

	sm := NewSafetyManager(true, false, false)
	err := generateCachePackage(sm)
	require.NoError(t, err)

	// File should NOT be created in dry-run
	_, err = os.ReadFile(filepath.Join(DirInternal, "cache", "redis.go"))
	assert.True(t, os.IsNotExist(err))

	// Should be in pending files
	pending := sm.GetPendingFiles()
	found := false
	for _, p := range pending {
		if strings.Contains(p.Path, "redis.go") {
			found = true
			break
		}
	}
	assert.True(t, found, "redis.go should be in pending files")
}

func TestGenerateCachePackage_Imports(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	require.NoError(t, os.MkdirAll(filepath.Join(DirInternal, "cache"), 0755))

	sm := NewSafetyManager(false, true, false)
	err := generateCachePackage(sm)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(DirInternal, "cache", "redis.go"))
	require.NoError(t, err)
	src := string(content)

	assert.Contains(t, src, `"context"`)
	assert.Contains(t, src, `"fmt"`)
	assert.Contains(t, src, `"os"`)
	assert.Contains(t, src, `"strconv"`)
	assert.Contains(t, src, `"github.com/redis/go-redis/v9"`)
}

func TestGenerateSetupRepositories_WithCache(t *testing.T) {
	t.Parallel()

	var b strings.Builder
	features := []string{"Product", "User"}
	generateSetupRepositories(&b, features, "postgres", true)
	output := b.String()

	assert.Contains(t, output, "baseProductRepo := repository.NewPostgresProductRepository(c.db)")
	assert.Contains(t, output, "c.productRepo = repository.NewCachedProductRepository(baseProductRepo, c.redisClient, 5*time.Minute)")
	assert.Contains(t, output, "baseUserRepo := repository.NewPostgresUserRepository(c.db)")
	assert.Contains(t, output, "c.userRepo = repository.NewCachedUserRepository(baseUserRepo, c.redisClient, 5*time.Minute)")
}

func TestGenerateSetupRepositories_WithCacheMySQL(t *testing.T) {
	t.Parallel()

	var b strings.Builder
	generateSetupRepositories(&b, []string{"Order"}, "mysql", true)
	output := b.String()

	assert.Contains(t, output, "baseOrderRepo := repository.NewMySQLOrderRepository(c.db)")
	assert.Contains(t, output, "c.orderRepo = repository.NewCachedOrderRepository(baseOrderRepo, c.redisClient, 5*time.Minute)")
}

func TestGenerateManualDI_WithCache(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))

	sm := NewSafetyManager(false, true, false)
	generateManualDI(filepath.Join(dir, "di"), []string{"Product"}, "postgres", true, sm)

	content, err := os.ReadFile(filepath.Join(dir, "di", "container.go"))
	require.NoError(t, err)
	src := string(content)

	// Should have redis imports
	assert.Contains(t, src, `"time"`)
	assert.Contains(t, src, `"github.com/redis/go-redis/v9"`)

	// Should have redis client field
	assert.Contains(t, src, "redisClient *redis.Client")

	// Constructor should accept redis client
	assert.Contains(t, src, "func NewContainer(db *gorm.DB, redisClient *redis.Client) *Container")
	assert.Contains(t, src, "redisClient: redisClient")

	// Repository setup should use cache decorator
	assert.Contains(t, src, "baseProductRepo := repository.NewPostgresProductRepository(c.db)")
	assert.Contains(t, src, "c.productRepo = repository.NewCachedProductRepository(baseProductRepo, c.redisClient, 5*time.Minute)")
}

func TestGenerateManualDI_WithoutCache(t *testing.T) {
	cleanup := ensureTestUI(t)
	defer cleanup()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	require.NoError(t, os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644))

	sm := NewSafetyManager(false, true, false)
	generateManualDI(filepath.Join(dir, "di"), []string{"Product"}, "postgres", false, sm)

	content, err := os.ReadFile(filepath.Join(dir, "di", "container.go"))
	require.NoError(t, err)
	src := string(content)

	// Should NOT have redis imports
	assert.NotContains(t, src, `"github.com/redis/go-redis/v9"`)

	// Constructor should only accept db
	assert.Contains(t, src, "func NewContainer(db *gorm.DB) *Container")
	assert.NotContains(t, src, "redisClient")

	// Direct repo wiring
	assert.Contains(t, src, "c.productRepo = repository.NewPostgresProductRepository(c.db)")
}

func TestAddSetupMethodsToDI_WithCache(t *testing.T) {
	t.Parallel()

	content := `func (c *Container) setupRepositories() {
}

func (c *Container) setupUseCases() {
}

func (c *Container) setupHandlers() {
}

// Getters`

	result := addSetupMethodsToDI(content, "Product", "product", true)
	assert.Contains(t, result, "baseProductRepo := repository.NewPostgresProductRepository(c.db)")
	assert.Contains(t, result, "c.productRepo = repository.NewCachedProductRepository(baseProductRepo, c.redisClient, 5*time.Minute)")
}

func TestAddSetupMethodsToDI_WithoutCache(t *testing.T) {
	t.Parallel()

	content := `func (c *Container) setupRepositories() {
}

func (c *Container) setupUseCases() {
}

func (c *Container) setupHandlers() {
}

// Getters`

	result := addSetupMethodsToDI(content, "Order", "order", false)
	assert.Contains(t, result, "c.orderRepo = repository.NewPostgresOrderRepository(c.db)")
	assert.NotContains(t, result, "baseOrderRepo")
	assert.NotContains(t, result, "NewCachedOrderRepository")
}
