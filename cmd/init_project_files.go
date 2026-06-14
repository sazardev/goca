package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func createGoMod(projectName, module, database string, auth bool, sm ...*SafetyManager) {
	var dependencies string

	// Base dependencies (common to all)
	baseDeps := `github.com/gorilla/mux v1.8.0`

	// Add database-specific dependencies
	switch database {
	case DBPostgres, DBPostgresJSON:
		baseDeps += `
	gorm.io/gorm v1.25.5
	gorm.io/driver/postgres v1.5.4`
	case DBMySQL:
		baseDeps += `
	gorm.io/gorm v1.25.5
	gorm.io/driver/mysql v1.5.2`
	case DBSQLite:
		baseDeps += `
	gorm.io/gorm v1.25.5
	gorm.io/driver/sqlite v1.5.4`
	case DBSQLServer:
		baseDeps += `
	gorm.io/gorm v1.25.5
	gorm.io/driver/sqlserver v1.5.2`
	case DBMongoDB:
		baseDeps += `
	go.mongodb.org/mongo-driver v1.12.1`
	case DBDynamoDB:
		baseDeps += `
	github.com/aws/aws-sdk-go-v2 v1.21.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.21.5
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.39`
	case DBElasticsearch:
		baseDeps += `
	github.com/elastic/go-elasticsearch/v8 v8.10.1`
	default: // sqlite as the safe fallback (matches createMainGo/databaseURLBody)
		baseDeps += `
	gorm.io/gorm v1.25.5
	gorm.io/driver/sqlite v1.5.4`
	}

	// Add JWT dependency if auth is enabled
	if auth {
		baseDeps += `
	github.com/golang-jwt/jwt/v4 v4.5.2`
	}

	dependencies = fmt.Sprintf(`require (
	%s
)`, baseDeps)

	content := fmt.Sprintf(`module %s

go 1.21

%s
`, module, dependencies)

	if err := writeFile(filepath.Join(projectName, "go.mod"), content, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error writing go.mod: %v", err))
		return
	}
} // downloadDependencies downloads Go module dependencies for the project

func downloadDependencies(projectName string) error {
	// First run go mod tidy to resolve dependencies and create go.sum
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}

	// Then download the dependencies
	cmd = exec.Command("go", "mod", "download")
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod download failed: %w", err)
	}

	return nil
}

func createGitignore(projectName string, sm ...*SafetyManager) {
	content := `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with go test -c
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# Environment variables
.env
.env.local

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Application specific
*.log
/tmp/
/dist/
`
	if err := writeFile(filepath.Join(projectName, ".gitignore"), content, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error writing .gitignore: %v", err))
		return
	}
}

// initializeGitRepository initializes a Git repository in the project directory.
func initializeGitRepository(projectName string) error {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return errors.New("git is not installed or not in PATH")
	}

	projectPath := filepath.Join(".", projectName)

	// Initialize git repository
	cmdInit := exec.Command("git", "init")
	cmdInit.Dir = projectPath
	if err := cmdInit.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// Add all files to staging
	cmdAdd := exec.Command("git", "add", ".")
	cmdAdd.Dir = projectPath
	if err := cmdAdd.Run(); err != nil {
		return fmt.Errorf("failed to add files to git: %w", err)
	}

	// Create initial commit. Disable GPG signing and detach stdin so the
	// command never blocks waiting for a passphrase/prompt (INIT-GIT).
	commitMessage := "Initial commit - Goca Clean Architecture project"
	cmdCommit := exec.Command("git", "-c", "commit.gpgsign=false", "commit", "-m", commitMessage)
	cmdCommit.Dir = projectPath
	cmdCommit.Stdin = nil
	if err := cmdCommit.Run(); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	return nil
}

func createReadme(projectName, module, database string, sm ...*SafetyManager) {
	dbDisplay := getDatabaseDisplayName(database)
	dbSection := getReadmeDatabaseSection(database, projectName)
	dbTroubleshooting := getReadmeDatabaseTroubleshooting(database, projectName)
	content := fmt.Sprintf(`# %s

Generated with Goca - Go Clean Architecture Code Generator

## Architecture

This project follows Clean Architecture principles:

- **Domain**: Entities and business rules
- **Use Cases**: Application logic
- **Repository**: Data abstraction
- **Handler**: Delivery adapters

## Quick Start

### 1. Install dependencies:
`+"```bash\n"+`go mod tidy
`+"```\n"+`

### 2. Configure database (%s):

%s

### 3. Configure environment variables:
`+"```bash\n"+`# Copy example file
cp .env.example .env

# Edit with your credentials
# DB_PASSWORD=password
# DB_NAME=%s
`+"```\n"+`

### 4. Run the application:
`+"```bash\n"+`go run cmd/server/main.go
`+"```\n"+`

### 5. Test endpoints:
`+"```bash\n"+`# Health check
curl http://localhost:8080/health

# Create user (if you have the User feature)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
`+"```\n"+`

## Project Structure

`+"```\n"+`%s/
├── cmd/
│   └── server/           # Application entry point
│       └── main.go
├── internal/
│   ├── domain/           # Entities and business rules
│   ├── usecase/          # Application logic
│   ├── repository/       # Persistence implementations
│   ├── handler/          # HTTP/gRPC adapters
│   │   └── http/
│   └── messages/         # Error and response messages
├── pkg/
│   ├── config/           # Application configuration
│   └── logger/           # Logging system
├── migrations/           # Database migrations
├── .env                  # Environment variables
├── .env.example          # Configuration example
├── docker-compose.yml    # Docker services
├── Makefile              # Useful commands
├── go.mod
└── README.md
`+"```\n"+`

## Useful Commands

### Generate new features:
`+"```bash\n"+`# Complete feature with all layers
goca feature User --fields "name:string,email:string"

# Feature with validations
goca feature Product --fields "name:string,price:float64" --validation

# Integrate existing features
goca integrate --all
`+"```\n"+`

### Development commands:
`+"```bash\n"+`# Run application
make run

# Run tests
make test

# Build for production
make build

# Linting and formatting
make lint
make fmt
`+"```\n"+`

## Troubleshooting

%s

### Error: "command not found: goca"
Goca CLI is not installed or not in PATH.

**Solution:**
`+"```bash\n"+`# Reinstall Goca
go install github.com/sazardev/goca@latest

# Verify installation
goca version
`+"```\n"+`

### Health Check shows "degraded"
Application runs but cannot connect to database.

**Solution:**
1. Verify %s is running
2. Verify environment variables in .env
3. Verify your connection settings in .env

## Additional Resources

- [Goca Documentation](https://github.com/sazardev/goca)
- [Clean Architecture Principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Complete Tutorial](https://github.com/sazardev/goca/wiki/Complete-Tutorial)

## Contributing

This project was generated with Goca. To contribute:

1. Add new features with `+"`"+`goca feature`+"`"+`
2. Maintain layer separation
3. Write tests for new functionality
4. Follow Clean Architecture conventions

---

Generated with [Goca](https://github.com/sazardev/goca)
`, cases.Title(language.English).String(projectName), dbDisplay, dbSection, projectName, projectName, dbTroubleshooting, dbDisplay)

	if err := writeFile(filepath.Join(projectName, "README.md"), content, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error writing README.md: %v", err))
		return
	}
}

// getDatabaseDisplayName returns a human-friendly name for the database.
func getDatabaseDisplayName(database string) string {
	switch database {
	case DBPostgres, DBPostgresJSON:
		return "PostgreSQL"
	case DBMySQL:
		return "MySQL"
	case DBMongoDB:
		return "MongoDB"
	case DBSQLite:
		return "SQLite"
	case DBSQLServer:
		return "SQL Server"
	case DBDynamoDB:
		return "DynamoDB"
	case DBElasticsearch:
		return "Elasticsearch"
	default:
		return database
	}
}

// getReadmeDatabaseSection returns the DB-specific "Configure database" section.
func getReadmeDatabaseSection(database, projectName string) string {
	fence := func(s string) string { return "```bash\n" + s + "\n```" }
	switch database {
	case DBSQLite:
		return "SQLite is file-based; no server is required. The database file is created automatically on first run."
	case DBMySQL:
		return "#### Option A: Using Docker (Recommended)\n" +
			fence(fmt.Sprintf("# Run MySQL\ndocker run --name mysql-dev \\\n  -e MYSQL_ROOT_PASSWORD=password \\\n  -e MYSQL_DATABASE=%s \\\n  -p 3306:3306 \\\n  -d mysql:8.0\n\n# Or using docker-compose\ndocker-compose up -d", projectName))
	case DBMongoDB:
		return "#### Option A: Using Docker (Recommended)\n" +
			fence(fmt.Sprintf("# Run MongoDB\ndocker run --name mongo-dev \\\n  -e MONGO_INITDB_ROOT_USERNAME=admin \\\n  -e MONGO_INITDB_ROOT_PASSWORD=password \\\n  -e MONGO_INITDB_DATABASE=%s \\\n  -p 27017:27017 \\\n  -d mongo:7.0\n\n# Or using docker-compose\ndocker-compose up -d", projectName))
	case DBSQLServer:
		return "#### Option A: Using Docker (Recommended)\n" +
			fence("# Run SQL Server\ndocker run --name sqlserver-dev \\\n  -e ACCEPT_EULA=Y \\\n  -e MSSQL_SA_PASSWORD=Your_password123 \\\n  -p 1433:1433 \\\n  -d mcr.microsoft.com/mssql/server:2022-latest\n\n# Or using docker-compose\ndocker-compose up -d")
	case DBDynamoDB:
		return "#### Option A: Using Docker (Recommended)\n" +
			fence("# Run DynamoDB Local\ndocker run --name dynamodb-dev \\\n  -p 8000:8000 \\\n  -d amazon/dynamodb-local:latest\n\n# Or using docker-compose\ndocker-compose up -d")
	case DBElasticsearch:
		return "#### Option A: Using Docker (Recommended)\n" +
			fence("# Run Elasticsearch\ndocker run --name elasticsearch-dev \\\n  -e discovery.type=single-node \\\n  -e xpack.security.enabled=false \\\n  -p 9200:9200 \\\n  -d docker.elastic.co/elasticsearch/elasticsearch:8.10.1\n\n# Or using docker-compose\ndocker-compose up -d")
	default: // postgres, postgres-json
		return "#### Option A: Using Docker (Recommended)\n" +
			fence(fmt.Sprintf("# Run PostgreSQL\ndocker run --name postgres-dev \\\n  -e POSTGRES_PASSWORD=password \\\n  -e POSTGRES_DB=%s \\\n  -p 5432:5432 \\\n  -d postgres:15\n\n# Or using docker-compose\ndocker-compose up -d", projectName)) +
			"\n\n#### Option B: Local PostgreSQL\n" +
			fence(fmt.Sprintf("# Create database\ncreatedb %s", projectName))
	}
}

// getReadmeDatabaseTroubleshooting returns the DB-specific troubleshooting block.
func getReadmeDatabaseTroubleshooting(database, projectName string) string {
	fence := func(s string) string { return "```bash\n" + s + "\n```" }
	display := getDatabaseDisplayName(database)
	port := getDatabasePort(database)
	user := getDatabaseUser(database)

	if database == DBSQLite {
		return "### Error: \"unable to open database file\"\n" +
			"SQLite cannot create or access the database file.\n\n" +
			"**Solution:** Ensure the application has write permission in its working directory."
	}

	envBlock := fmt.Sprintf("# Configure in .env\nDB_HOST=localhost\nDB_PORT=%s\nDB_USER=%s\nDB_PASSWORD=password\nDB_NAME=%s", port, user, projectName)

	return fmt.Sprintf("### Error: \"connection refused\"\n%s database is not running.\n\n**Solution:** Start the database service (see step 2) and verify it is reachable on port %s.\n\n### Error: \"database not configured\"\nDatabase environment variables are not configured.\n\n**Solution:**\n%s", display, port, fence(envBlock))
}

func createConfig(projectName, _, database string, sm ...*SafetyManager) {
	dbURLBody := databaseURLBody(database)
	// "fmt" is only needed when the DSN is built with fmt.Sprintf (every driver
	// except SQLite, whose URL is a plain file path).
	fmtImport := ""
	if strings.Contains(dbURLBody, "fmt.") {
		fmtImport = "\t\"fmt\"\n"
	}
	content := fmt.Sprintf(`package config

import (
%s	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port        string
	Environment string
	LogLevel    string
	Database    DatabaseConfig
	Server      ServerConfig
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

type ServerConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "%s"),
			User:         getEnv("DB_USER", "%s"),
			Password:     getEnv("DB_PASSWORD", ""),
			Name:         getEnv("DB_NAME", "%s"),
			SSLMode:      getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
			MaxLifetime:  getEnvAsDuration("DB_MAX_LIFETIME", "5m"),
		},
		Server: ServerConfig{
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", "10s"),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", "10s"),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", "60s"),
		},
	}
}

func (c *Config) GetDatabaseURL() string {
%s
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Invalid integer value for %%s: %%s, using default: %%d", key, value, defaultValue)
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("Warning: Invalid duration value for %%s: %%s, using default: %%s", key, value, defaultValue)
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
`, fmtImport, getConfigDefaultPort(database), getDatabaseUser(database), projectName, dbURLBody)

	if err := writeGoFile(filepath.Join(projectName, "pkg", "config", "config.go"), content, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error writing config.go: %v", err))
		return
	}
}

// databaseURLBody returns the body of Config.GetDatabaseURL for the given
// database driver, so the generated DSN matches the configured database
// instead of always emitting a PostgreSQL connection string.
func databaseURLBody(database string) string {
	switch database {
	case DBSQLite:
		// GORM's sqlite driver expects a file path (or :memory:).
		return "\tname := c.Database.Name\n\tif name == \"\" {\n\t\tname = \"app\"\n\t}\n\treturn name + \".db\""
	case DBMySQL:
		return "\treturn fmt.Sprintf(\"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local\",\n\t\tc.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)"
	case DBSQLServer:
		return "\treturn fmt.Sprintf(\"sqlserver://%s:%s@%s:%s?database=%s\",\n\t\tc.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)"
	case DBMongoDB:
		return "\tif c.Database.User != \"\" && c.Database.Password != \"\" {\n\t\treturn fmt.Sprintf(\"mongodb://%s:%s@%s:%s\", c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port)\n\t}\n\treturn fmt.Sprintf(\"mongodb://%s:%s\", c.Database.Host, c.Database.Port)"
	case DBPostgres, DBPostgresJSON:
		return "\treturn fmt.Sprintf(\"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s\",\n\t\tc.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name, c.Database.SSLMode)"
	default: // safe fallback: sqlite file path (matches createGoMod/createMainGo)
		return "\tname := c.Database.Name\n\tif name == \"\" {\n\t\tname = \"app\"\n\t}\n\treturn name + \".db\""
	}
}

func createLogger(projectName, _ string, sm ...*SafetyManager) {
	content := `package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func Init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(v ...interface{}) {
	InfoLogger.Println(v...)
}

func Error(v ...interface{}) {
	ErrorLogger.Println(v...)
}
`
	if err := writeGoFile(filepath.Join(projectName, "pkg", "logger", "logger.go"), content, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error writing logger.go: %v", err))
		return
	}
}

func createAuth(projectName, module string, sm ...*SafetyManager) {
	content := `package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key") // Change this!

type Claims struct {
	UserID int    ` + "`json:\"user_id\"`" + `
	Email  string ` + "`json:\"email\"`" + `
	jwt.RegisteredClaims
}

func GenerateToken(userID int, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
`
	if err := writeGoFile(filepath.Join(projectName, "pkg", "auth", "jwt.go"), content, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error creating JWT file: %v", err))
	}
}

func createEnvFiles(projectName, database string, sm ...*SafetyManager) {
	// Create .env.example
	envExampleContent := fmt.Sprintf(`# Server Configuration
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info

# Database Configuration
DB_HOST=localhost
DB_PORT=%s
DB_USER=%s
DB_PASSWORD=
DB_NAME=%s
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_MAX_LIFETIME=5m

# Server Timeouts
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
SERVER_IDLE_TIMEOUT=60s

# JWT Configuration (if using auth)
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_ISSUER=%s
JWT_EXPIRY=24h

# External Services (uncomment as needed)
# REDIS_URL=redis://localhost:6379
# ELASTICSEARCH_URL=http://localhost:9200
# SMTP_HOST=localhost
# SMTP_PORT=587
# SMTP_USER=
# SMTP_PASSWORD=
`, getDatabasePort(database), getDatabaseUser(database), projectName, projectName)

	if err := writeFile(filepath.Join(projectName, ".env.example"), envExampleContent, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error creating .env.example file: %v", err))
	}

	// Create .env file with defaults
	envContent := fmt.Sprintf(`# Server Configuration
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info

# Database Configuration
DB_HOST=localhost
DB_PORT=%s
DB_USER=%s
DB_PASSWORD=
DB_NAME=%s
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_MAX_LIFETIME=5m

# Server Timeouts
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
SERVER_IDLE_TIMEOUT=60s

# JWT Configuration
JWT_SECRET=dev-secret-change-in-production
JWT_ISSUER=%s
JWT_EXPIRY=24h
`, getDatabasePort(database), getDatabaseUser(database), projectName, projectName)

	if err := writeFile(filepath.Join(projectName, ".env"), envContent, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error creating .env file: %v", err))
	}
}

// getConfigDefaultPort returns the DB_PORT default for the generated config.go.
// SQLite has no network port, so it returns an empty string.
func getConfigDefaultPort(database string) string {
	return getDatabasePort(database)
}

func getDatabasePort(database string) string {
	switch database {
	case DBMySQL:
		return "3306"
	case DBMongoDB:
		return "27017"
	case DBSQLServer:
		return "1433"
	case DBDynamoDB:
		return "8000"
	case DBElasticsearch:
		return "9200"
	case DBSQLite:
		return ""
	default: // postgres, postgres-json
		return "5432"
	}
}

func getDatabaseUser(database string) string {
	switch database {
	case DBMySQL:
		return "root"
	case DBMongoDB:
		return "admin"
	case DBSQLServer:
		return "sa"
	case DBElasticsearch:
		return "elastic"
	case DBDynamoDB:
		return "local"
	case DBSQLite:
		return ""
	default: // postgres, postgres-json
		return "postgres"
	}
}
