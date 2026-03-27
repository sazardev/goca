package cmd

import (
	"path/filepath"
	"strings"
)

// generateCachePackage creates internal/cache/redis.go with a Redis client
// factory that reads connection details from environment variables.
func generateCachePackage(sm ...*SafetyManager) error {
	cacheDir := filepath.Join(DirInternal, "cache")
	filename := filepath.Join(cacheDir, "redis.go")

	var b strings.Builder

	b.WriteString("package cache\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"context\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strconv\"\n\n")
	b.WriteString("\t\"github.com/redis/go-redis/v9\"\n")
	b.WriteString(")\n\n")

	b.WriteString("// NewRedisClient creates a Redis client using environment variables:\n")
	b.WriteString("//   REDIS_URL      — address (default \"localhost:6379\")\n")
	b.WriteString("//   REDIS_PASSWORD  — password (default \"\")\n")
	b.WriteString("//   REDIS_DB       — database index (default 0)\n")
	b.WriteString("func NewRedisClient() (*redis.Client, error) {\n")
	b.WriteString("\taddr := os.Getenv(\"REDIS_URL\")\n")
	b.WriteString("\tif addr == \"\" {\n")
	b.WriteString("\t\taddr = \"localhost:6379\"\n")
	b.WriteString("\t}\n\n")
	b.WriteString("\tpassword := os.Getenv(\"REDIS_PASSWORD\")\n\n")
	b.WriteString("\tdb := 0\n")
	b.WriteString("\tif dbStr := os.Getenv(\"REDIS_DB\"); dbStr != \"\" {\n")
	b.WriteString("\t\tparsed, err := strconv.Atoi(dbStr)\n")
	b.WriteString("\t\tif err != nil {\n")
	b.WriteString("\t\t\treturn nil, fmt.Errorf(\"invalid REDIS_DB: %w\", err)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t\tdb = parsed\n")
	b.WriteString("\t}\n\n")
	b.WriteString("\tclient := redis.NewClient(&redis.Options{\n")
	b.WriteString("\t\tAddr:     addr,\n")
	b.WriteString("\t\tPassword: password,\n")
	b.WriteString("\t\tDB:       db,\n")
	b.WriteString("\t})\n\n")
	b.WriteString("\tif err := client.Ping(context.Background()).Err(); err != nil {\n")
	b.WriteString("\t\treturn nil, fmt.Errorf(\"redis connection failed: %w\", err)\n")
	b.WriteString("\t}\n\n")
	b.WriteString("\treturn client, nil\n")
	b.WriteString("}\n")

	return writeGoFile(filename, b.String(), sm...)
}
