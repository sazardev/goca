package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
)

// generateCacheDecorator produces internal/repository/cached_<entity>_repository.go
// implementing the <Entity>Repository interface with a Redis caching layer.
func generateCacheDecorator(entity string, fields []Field, sm ...*SafetyManager) {
	entityLower := strings.ToLower(entity)
	repoDir := filepath.Join(DirInternal, DirRepository)
	filename := filepath.Join(repoDir, "cached_"+entityLower+"_repository.go")

	moduleName := getModuleName()
	importPath := getImportPath(moduleName)

	var b strings.Builder

	b.WriteString("package repository\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"context\"\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"time\"\n\n")
	b.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n\n", importPath))
	b.WriteString("\t\"github.com/redis/go-redis/v9\"\n")
	b.WriteString(")\n\n")

	// Struct
	b.WriteString(fmt.Sprintf("// Cached%sRepository is a caching decorator around %sRepository.\n", entity, entity))
	b.WriteString(fmt.Sprintf("// Read operations check Redis first; write operations delegate then invalidate.\n"))
	b.WriteString(fmt.Sprintf("type Cached%sRepository struct {\n", entity))
	b.WriteString(fmt.Sprintf("\tinner    %sRepository\n", entity))
	b.WriteString("\tcache    *redis.Client\n")
	b.WriteString("\tcacheTTL time.Duration\n")
	b.WriteString("\tctx      context.Context\n")
	b.WriteString("}\n\n")

	// Constructor
	b.WriteString(fmt.Sprintf("// NewCached%sRepository creates a caching decorator that wraps inner.\n", entity))
	b.WriteString(fmt.Sprintf("func NewCached%sRepository(inner %sRepository, cache *redis.Client, ttl time.Duration) *Cached%sRepository {\n", entity, entity, entity))
	b.WriteString(fmt.Sprintf("\treturn &Cached%sRepository{\n", entity))
	b.WriteString("\t\tinner:    inner,\n")
	b.WriteString("\t\tcache:    cache,\n")
	b.WriteString("\t\tcacheTTL: ttl,\n")
	b.WriteString("\t\tctx:      context.Background(),\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n\n")

	// Cache key helpers
	cachePrefix := entityLower
	b.WriteString(fmt.Sprintf("func (r *Cached%sRepository) cacheKey(id int) string {\n", entity))
	b.WriteString(fmt.Sprintf("\treturn fmt.Sprintf(\"%s:%%d\", id)\n", cachePrefix))
	b.WriteString("}\n\n")
	b.WriteString(fmt.Sprintf("func (r *Cached%sRepository) listCacheKey() string {\n", entity))
	b.WriteString(fmt.Sprintf("\treturn \"%s:list\"\n", cachePrefix))
	b.WriteString("}\n\n")

	// Save — delegate + invalidate
	b.WriteString(fmt.Sprintf("func (r *Cached%sRepository) Save(%s *domain.%s) error {\n", entity, entityLower, entity))
	b.WriteString(fmt.Sprintf("\tif err := r.inner.Save(%s); err != nil {\n", entityLower))
	b.WriteString("\t\treturn err\n")
	b.WriteString("\t}\n")
	b.WriteString("\tr.cache.Del(r.ctx, r.listCacheKey())\n")
	b.WriteString("\treturn nil\n")
	b.WriteString("}\n\n")

	// FindByID — check cache → miss → delegate → set
	b.WriteString(fmt.Sprintf("func (r *Cached%sRepository) FindByID(id int) (*domain.%s, error) {\n", entity, entity))
	b.WriteString("\tkey := r.cacheKey(id)\n")
	b.WriteString("\tcached, err := r.cache.Get(r.ctx, key).Bytes()\n")
	b.WriteString("\tif err == nil {\n")
	b.WriteString(fmt.Sprintf("\t\tvar result domain.%s\n", entity))
	b.WriteString("\t\tif json.Unmarshal(cached, &result) == nil {\n")
	b.WriteString("\t\t\treturn &result, nil\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n\n")
	b.WriteString("\tentity, err := r.inner.FindByID(id)\n")
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\treturn nil, err\n")
	b.WriteString("\t}\n\n")
	b.WriteString("\tif data, mErr := json.Marshal(entity); mErr == nil {\n")
	b.WriteString("\t\tr.cache.Set(r.ctx, key, data, r.cacheTTL)\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn entity, nil\n")
	b.WriteString("}\n\n")

	// Dynamic search methods from fields — delegate only (no caching for search)
	if len(fields) > 0 {
		searchMethods := generateSearchMethods(fields, entity)
		for _, m := range searchMethods {
			generateCacheSearchMethodDelegate(&b, entity, m)
		}
	} else {
		// Default FindByEmail when no fields are specified
		b.WriteString(fmt.Sprintf("func (r *Cached%sRepository) FindByEmail(email string) (*domain.%s, error) {\n", entity, entity))
		b.WriteString("\treturn r.inner.FindByEmail(email)\n")
		b.WriteString("}\n\n")
	}

	// Update — delegate + invalidate
	b.WriteString(fmt.Sprintf("func (r *Cached%sRepository) Update(%s *domain.%s) error {\n", entity, entityLower, entity))
	b.WriteString(fmt.Sprintf("\tif err := r.inner.Update(%s); err != nil {\n", entityLower))
	b.WriteString("\t\treturn err\n")
	b.WriteString("\t}\n")
	b.WriteString(fmt.Sprintf("\tr.cache.Del(r.ctx, r.cacheKey(int(%s.ID)), r.listCacheKey())\n", entityLower))
	b.WriteString("\treturn nil\n")
	b.WriteString("}\n\n")

	// Delete — delegate + invalidate
	b.WriteString(fmt.Sprintf("func (r *Cached%sRepository) Delete(id int) error {\n", entity))
	b.WriteString("\tif err := r.inner.Delete(id); err != nil {\n")
	b.WriteString("\t\treturn err\n")
	b.WriteString("\t}\n")
	b.WriteString("\tr.cache.Del(r.ctx, r.cacheKey(id), r.listCacheKey())\n")
	b.WriteString("\treturn nil\n")
	b.WriteString("}\n\n")

	// FindAll — check cache → miss → delegate → set
	b.WriteString(fmt.Sprintf("func (r *Cached%sRepository) FindAll() ([]domain.%s, error) {\n", entity, entity))
	b.WriteString("\tkey := r.listCacheKey()\n")
	b.WriteString("\tcached, err := r.cache.Get(r.ctx, key).Bytes()\n")
	b.WriteString("\tif err == nil {\n")
	b.WriteString(fmt.Sprintf("\t\tvar result []domain.%s\n", entity))
	b.WriteString("\t\tif json.Unmarshal(cached, &result) == nil {\n")
	b.WriteString("\t\t\treturn result, nil\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n\n")
	b.WriteString("\tentities, err := r.inner.FindAll()\n")
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\treturn nil, err\n")
	b.WriteString("\t}\n\n")
	b.WriteString("\tif data, mErr := json.Marshal(entities); mErr == nil {\n")
	b.WriteString("\t\tr.cache.Set(r.ctx, key, data, r.cacheTTL)\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn entities, nil\n")
	b.WriteString("}\n")

	if err := writeGoFile(filename, b.String(), sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing cache decorator: %v", err))
	}
}

// generateCacheSearchMethodDelegate generates a delegate-only method for a search method.
func generateCacheSearchMethodDelegate(b *strings.Builder, entity string, m SearchMethod) {
	paramName := strings.ToLower(m.FieldName)

	b.WriteString(fmt.Sprintf("func (r *Cached%sRepository) %s(%s %s) %s {\n",
		entity, m.MethodName, paramName, m.FieldType, m.ReturnType))
	b.WriteString(fmt.Sprintf("\treturn r.inner.%s(%s)\n", m.MethodName, paramName))
	b.WriteString("}\n\n")
}
