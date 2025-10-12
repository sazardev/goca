package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var initCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Initialize Clean Architecture project",
	Long: `Creates the base structure of a Go project following Clean Architecture principles, 
including directories, configuration files and layer structure.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		module, _ := cmd.Flags().GetString("module")
		database, _ := cmd.Flags().GetString("database")
		auth, _ := cmd.Flags().GetBool("auth")
		api, _ := cmd.Flags().GetString("api")
		config, _ := cmd.Flags().GetBool("config")

		if module == "" {
			fmt.Println("Error: --module flag is required")
			os.Exit(1)
		}

		fmt.Printf("Initializing project '%s' with module '%s'\n", projectName, module)
		fmt.Printf("Database: %s\n", database)
		fmt.Printf("API: %s\n", api)
		if auth {
			fmt.Println("Including authentication")
		}
		if config {
			fmt.Println("Generating YAML configuration")
		}

		// Create configuration integration
		configIntegration := NewConfigIntegration()

		// Merge CLI flags with configuration
		flags := map[string]interface{}{
			"database": database,
			"auth":     auth,
			"api":      api,
		}
		configIntegration.MergeWithCLIFlags(flags)

		createProjectStructure(projectName, module, database, auth, api, configIntegration, config)

		fmt.Printf("\nProject '%s' created successfully!\n", projectName)
		fmt.Printf("Directory: ./%s\n", projectName)

		if config {
			fmt.Printf("Configuration file: ./%s/.goca.yaml\n", projectName)
		}

		fmt.Println("\nNext steps:")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Println("  go mod tidy")

		if config {
			fmt.Println("  # Edit .goca.yaml to customize your project")
			fmt.Println("  goca feature User --fields \"name:string,email:string\"")
		} else {
			fmt.Println("  goca feature User --fields \"name:string,email:string\"")
			fmt.Printf("  # Or generate config: goca init %s --module %s --config\n", projectName, module)
		}
	},
}

func createProjectStructure(projectName, module, database string, auth bool, api string, configIntegration *ConfigIntegration, generateConfig bool) {
	// Create main directories
	dirs := []string{
		filepath.Join(projectName, "cmd", "server"),
		filepath.Join(projectName, "internal", "domain"),
		filepath.Join(projectName, "internal", "usecase"),
		filepath.Join(projectName, "internal", "repository"),
		filepath.Join(projectName, "internal", "handler"),
		filepath.Join(projectName, "pkg", "config"),
		filepath.Join(projectName, "pkg", "logger"),
	}

	if auth {
		dirs = append(dirs, filepath.Join(projectName, "pkg", "auth"))
	}

	for _, dir := range dirs {
		_ = os.MkdirAll(dir, 0755)
	}

	// Create go.mod
	createGoMod(projectName, module, database, auth)

	// Create main.go
	createMainGo(projectName, module, api)

	// Create .gitignore
	createGitignore(projectName)

	// Create README.md
	createReadme(projectName, module)

	// Create config
	createConfig(projectName, module, database)

	// Create environment files
	createEnvFiles(projectName, database)

	// Create migrations
	createMigrations(projectName)

	// Create Makefile and Docker files
	createMakefile(projectName)
	createDockerfiles(projectName, database)

	// Create logger
	createLogger(projectName, module)

	if auth {
		createAuth(projectName, module)
	}

	// Generate .goca.yaml configuration file if requested
	if generateConfig && configIntegration != nil {
		configPath := filepath.Join(projectName, ".goca.yaml")
		if err := configIntegration.GenerateConfigFile(projectName, projectName, module, database); err != nil {
			fmt.Printf("Warning: Failed to generate config file: %v\n", err)
		} else {
			fmt.Printf("Generated configuration file: %s\n", configPath)
		}
	}

	// Download dependencies after creating go.mod
	if err := downloadDependencies(projectName); err != nil {
		fmt.Printf("Warning: Failed to download dependencies: %v\n", err)
		fmt.Printf("Tip: Run 'go mod download' manually in the project directory\n")
	}

	// Initialize Git repository
	fmt.Println("Initializing Git repository...")
	if err := initializeGitRepository(projectName); err != nil {
		fmt.Printf("Warning: Failed to initialize Git repository: %v\n", err)
		fmt.Printf("Tip: You can initialize Git manually with 'git init' in the project directory\n")
	} else {
		fmt.Println("Git repository initialized with initial commit")
	}
}

func createGoMod(projectName, module, database string, auth bool) {
	var dependencies string

	// Base dependencies
	baseDeps := `github.com/gorilla/mux v1.8.0
	gorm.io/gorm v1.25.5
	gorm.io/driver/postgres v1.5.4`

	switch database {
	case DBMySQL:
		baseDeps += `
	gorm.io/driver/mysql v1.5.2`
	case DBMongoDB:
		baseDeps += `
	go.mongodb.org/mongo-driver v1.12.1`
	default: // postgres - already included above
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

	// Add replace directive for test modules to make them resolvable locally
	if strings.Contains(module, "github.com/goca/testproject") {
		content += `
replace github.com/goca/testproject => ./
`
	}

	if err := writeFile(filepath.Join(projectName, "go.mod"), content); err != nil {
		fmt.Printf("Error escribiendo go.mod: %v\n", err)
		return
	}
} // downloadDependencies downloads Go module dependencies for the project
func downloadDependencies(projectName string) error {
	// First run go mod tidy to resolve dependencies and create go.sum
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %v", err)
	}

	// Then download the dependencies
	cmd = exec.Command("go", "mod", "download")
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod download failed: %v", err)
	}

	return nil
}

func createMainGo(projectName, module, _ string) {
	content := fmt.Sprintf(`package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"%s/pkg/config"
	"%s/pkg/logger"
)

type HealthStatus struct {
	Status    string            `+"`"+`json:"status"`+"`"+`
	Timestamp time.Time         `+"`"+`json:"timestamp"`+"`"+`
	Services  map[string]string `+"`"+`json:"services"`+"`"+`
	Version   string            `+"`"+`json:"version"`+"`"+`
}

var (
	// Build information (set by build flags)
	Version   = "dev"
	BuildTime = "unknown"
	db        *gorm.DB
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Initialize logger
	logger.Init()
	
	log.Printf("Starting application v%%s (built: %%s)", Version, BuildTime)
	log.Printf("Environment: %%s", cfg.Environment)
	
	// Connect to database with retry
	var err error
	db, err = connectToDatabase(cfg)
	if err != nil {
		log.Printf("Warning: Database connection failed: %%v", err)
		log.Printf("Server will start in degraded mode. Check your database configuration.")
		log.Printf("Tip: Configure database environment variables in .env file")
		db = nil // Ensure db is nil for health checks
	} else {
		log.Printf("Database connected successfully")
		
		// Run auto-migrations if database is connected
		if err := runAutoMigrations(db); err != nil {
			log.Printf("Warning: Auto-migration failed: %%v", err)
			log.Printf("Tip: You may need to run migrations manually")
		} else {
			log.Printf("Database schema is up to date")
		}
	}
	
	// Setup router
	router := mux.NewRouter()
	
	// Health check endpoint with comprehensive checks
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/health/ready", readinessHandler).Methods("GET")
	router.HandleFunc("/health/live", livenessHandler).Methods("GET")
	
	// Setup HTTP server with timeouts
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	
	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %%v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	
	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %%v", err)
	}
	
	log.Println("Server exited")
}

func connectToDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDatabaseURL()
	
	log.Printf("Connecting to database at %%s:%%s/%%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	
	// Check if this is development mode without database
	if cfg.Environment == "development" && cfg.Database.Password == "" {
		log.Println("Warning: Development mode detected: No database password set")
		log.Println("To connect to PostgreSQL, set environment variables:")
		log.Println("   DB_HOST=localhost")
		log.Println("   DB_PORT=5432") 
		log.Println("   DB_USER=postgres")
		log.Println("   DB_PASSWORD=your_password")
		log.Println("   DB_NAME=your_database")
		log.Println("Server will continue without database connection...")
		return nil, fmt.Errorf("development mode: database not configured")
	}
	
	// Retry connection up to 5 times
	for i := 0; i < 5; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Attempt %%d: Failed to open database connection: %%v", i+1, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		
		// Get underlying sql.DB for connection pool configuration
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Attempt %%d: Failed to get underlying SQL DB: %%v", i+1, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		
		// Configure connection pool
		sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(cfg.Database.MaxLifetime)
		
		// Test the connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = sqlDB.PingContext(ctx)
		cancel()
		
		if err == nil {
			return db, nil
		}
		
		log.Printf("Attempt %%d: Database ping failed: %%v", i+1, err)
		sqlDBClose, _ := db.DB()
		if sqlDBClose != nil {
			sqlDBClose.Close()
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	
	return nil, fmt.Errorf("failed to connect to database after 5 attempts")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
		Version:   Version,
	}
	
	// Check database
	if err := checkDatabase(); err != nil {
		status.Status = "degraded"
		status.Services["database"] = fmt.Sprintf("error: %%v", err)
		// Don't fail the whole health check for database issues in development
		log.Printf("Database health check failed: %%v", err)
	} else {
		status.Services["database"] = "healthy"
	}
	
	// Always return 200 for basic health check - let readiness handle critical dependencies
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	// Check if all dependencies are ready
	if err := checkDatabase(); err != nil {
		http.Error(w, fmt.Sprintf("Database not ready: %%v", err), http.StatusServiceUnavailable)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	// Basic liveness check
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alive"))
}

func checkDatabase() error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql DB: %%w", err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	return sqlDB.PingContext(ctx)
}

func runAutoMigrations(database *gorm.DB) error {
	if database == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	// Auto-migrate domain entities using GORM
	log.Println("Running GORM auto-migrations...")
	
	// Create a slice of all domain entities to migrate
	entities := []interface{}{
		// Add domain entities here as they are created
		// Example: &domain.User{}, &domain.Product{}
	}
	
	// Run auto-migration for all entities
	for _, entity := range entities {
		if err := database.AutoMigrate(entity); err != nil {
			return fmt.Errorf("failed to auto-migrate entity %%T: %%w", entity, err)
		}
	}
	
	// For now, just ensure the connection works
	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %%w", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %%w", err)
	}
	
	log.Println("GORM auto-migrations completed successfully")
	return nil
}

`, module, module)

	if err := writeGoFile(filepath.Join(projectName, "cmd", "server", "main.go"), content); err != nil {
		fmt.Printf("Error writing main.go: %v\n", err)
		return
	}
}

func createGitignore(projectName string) {
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
	if err := writeFile(filepath.Join(projectName, ".gitignore"), content); err != nil {
		fmt.Printf("Error escribiendo .gitignore: %v\n", err)
		return
	}
}

// initializeGitRepository initializes a Git repository in the project directory
func initializeGitRepository(projectName string) error {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed or not in PATH")
	}

	projectPath := filepath.Join(".", projectName)

	// Initialize git repository
	cmdInit := exec.Command("git", "init")
	cmdInit.Dir = projectPath
	if err := cmdInit.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %v", err)
	}

	// Add all files to staging
	cmdAdd := exec.Command("git", "add", ".")
	cmdAdd.Dir = projectPath
	if err := cmdAdd.Run(); err != nil {
		return fmt.Errorf("failed to add files to git: %v", err)
	}

	// Create initial commit
	commitMessage := "Initial commit - Goca Clean Architecture project"
	cmdCommit := exec.Command("git", "commit", "-m", commitMessage)
	cmdCommit.Dir = projectPath
	if err := cmdCommit.Run(); err != nil {
		return fmt.Errorf("failed to create initial commit: %v", err)
	}

	return nil
}

func createReadme(projectName, module string) {
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

### 2. Configure database (PostgreSQL):

#### Option A: Using Docker (Recommended)
`+"```bash\n"+`# Run PostgreSQL
docker run --name postgres-dev \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=%s \
  -p 5432:5432 \
  -d postgres:15

# Or using docker-compose
docker-compose up -d postgres
`+"```\n"+`

#### Option B: Local PostgreSQL
`+"```bash\n"+`# Create database
createdb %s
`+"```\n"+`

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

### Error: "dial tcp [::1]:5432: connection refused"
PostgreSQL database is not running. 

**Solution:**
`+"```bash\n"+`# With Docker
docker run --name postgres-dev \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=%s \
  -p 5432:5432 \
  -d postgres:15

# Verify it's running
docker ps
`+"```\n"+`

### Error: "database not configured"
Database environment variables are not configured.

**Solution:**
`+"```bash\n"+`# Configure in .env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=%s
`+"```\n"+`

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
1. Verify PostgreSQL is running
2. Verify environment variables in .env
3. Test connection manually: `+"`"+`psql -h localhost -U postgres -d %s`+"`"+`

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
`, cases.Title(language.English).String(projectName), projectName, projectName, projectName, projectName, projectName, projectName, projectName)

	if err := writeFile(filepath.Join(projectName, "README.md"), content); err != nil {
		fmt.Printf("Error escribiendo README.md: %v\n", err)
		return
	}
}

func createConfig(projectName, _, database string) {
	content := fmt.Sprintf(`package config

import (
	"fmt"
	"log"
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
			Port:         getEnv("DB_PORT", "5432"),
			User:         getEnv("DB_USER", "postgres"),
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
	return fmt.Sprintf("host=%%s port=%%s user=%%s password=%%s dbname=%%s sslmode=%%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
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
`, projectName)

	if err := writeGoFile(filepath.Join(projectName, "pkg", "config", "config.go"), content); err != nil {
		fmt.Printf("Error escribiendo config.go: %v\n", err)
		return
	}
}

func createLogger(projectName, _ string) {
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
	if err := writeGoFile(filepath.Join(projectName, "pkg", "logger", "logger.go"), content); err != nil {
		fmt.Printf("Error escribiendo logger.go: %v\n", err)
		return
	}
}

func createAuth(projectName, module string) {
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
	if err := writeGoFile(filepath.Join(projectName, "pkg", "auth", "jwt.go"), content); err != nil {
		fmt.Printf("Error creating JWT file: %v\n", err)
	}
}

func createEnvFiles(projectName, database string) {
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

	if err := writeFile(filepath.Join(projectName, ".env.example"), envExampleContent); err != nil {
		fmt.Printf("Error creating .env.example file: %v\n", err)
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

	if err := writeFile(filepath.Join(projectName, ".env"), envContent); err != nil {
		fmt.Printf("Error creating .env file: %v\n", err)
	}
}

func getDatabasePort(database string) string {
	switch database {
	case "mysql":
		return "3306"
	case "mongodb":
		return "27017"
	default: // postgres
		return "5432"
	}
}

func getDatabaseUser(database string) string {
	switch database {
	case "mysql":
		return "root"
	case "mongodb":
		return "admin"
	default: // postgres
		return "postgres"
	}
}

func createMigrations(projectName string) {
	// Create migrations directory
	migrationDir := filepath.Join(projectName, "migrations")
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		fmt.Printf("Error creating migrations directory: %v\n", err)
		return
	}

	// Create initial migration
	migrationContent := `-- Initial migration
-- This file contains the initial database schema

-- Enable UUID extension (PostgreSQL)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Example: Users table
-- Uncomment and modify based on your needs
/*
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
*/

-- Add your tables here
-- Remember to create corresponding down migration files for rollbacks
`

	if err := writeFile(filepath.Join(migrationDir, "001_initial.up.sql"), migrationContent); err != nil {
		fmt.Printf("Error creating up migration: %v\n", err)
	}

	// Create down migration
	downMigrationContent := `-- Down migration for initial schema
-- This file should reverse changes made in 001_initial.up.sql

-- Example: Drop users table
-- DROP TABLE IF EXISTS users;

-- Add your down migration here
`

	if err := writeFile(filepath.Join(migrationDir, "001_initial.down.sql"), downMigrationContent); err != nil {
		fmt.Printf("Error creating down migration: %v\n", err)
	}

	// Create README for migrations
	migrationReadme := `# Database Migrations

This directory contains database migration files.

## Structure
- *.up.sql - Migration files (apply changes)
- *.down.sql - Rollback files (reverse changes)

## Usage

### Using golang-migrate tool
bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path ./migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" up

# Rollback last migration
migrate -path ./migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" down 1


### Manual execution
bash
# Apply migration
psql -h localhost -U postgres -d your_db -f migrations/001_initial.up.sql

# Rollback migration  
psql -h localhost -U postgres -d your_db -f migrations/001_initial.down.sql


## Creating new migrations
1. Create new files: 002_description.up.sql and 002_description.down.sql
2. Add your changes in the .up.sql file
3. Add the reverse changes in the .down.sql file
`

	if err := writeFile(filepath.Join(migrationDir, "README.md"), migrationReadme); err != nil {
		fmt.Printf("Error creating migration README: %v\n", err)
	}
}

func createMakefile(projectName string) {
	makefileContent := fmt.Sprintf(`# Makefile for %s
.PHONY: help build run test clean docker-build docker-run deps lint migrate-up migrate-down

# Variables
APP_NAME := %s
DOCKER_IMAGE := %s:latest
MIGRATE_PATH := ./migrations
DATABASE_URL := postgres://postgres:@localhost/%s?sslmode=disable

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%%-20s\033[0m %%s\n", $$1, $$2}'

deps: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o bin/$(APP_NAME) cmd/server/main.go

run: ## Run the application
	go run cmd/server/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	golangci-lint run

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

# Database migrations
migrate-install: ## Install migrate tool
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up: ## Apply database migrations
	migrate -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" up

migrate-down: ## Rollback last migration
	migrate -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" down 1

migrate-force: ## Force migration version (use with caution)
	migrate -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" force $(VERSION)

# Docker commands
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run application in Docker
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

docker-compose-up: ## Start services with Docker Compose
	docker-compose up -d

docker-compose-down: ## Stop services with Docker Compose
	docker-compose down

# Development helpers
dev-db: ## Start development database
	docker run --name %s-postgres -e POSTGRES_DB=%s -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15

dev-db-stop: ## Stop development database
	docker stop %s-postgres && docker rm %s-postgres

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

mod-upgrade: ## Upgrade dependencies
	go get -u ./...
	go mod tidy

# Production helpers
build-prod: ## Build for production
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/$(APP_NAME) cmd/server/main.go

# Security
sec-scan: ## Run security scan
	gosec ./...

# API documentation
api-docs: ## Generate API documentation
	swag init -g cmd/server/main.go
`, projectName, projectName, projectName, projectName, projectName, projectName, projectName, projectName)

	if err := writeFile(filepath.Join(projectName, "Makefile"), makefileContent); err != nil {
		fmt.Printf("Error creating Makefile: %v\n", err)
	}
}

func createDockerfiles(projectName, database string) {
	// Dockerfile
	dockerfileContent := fmt.Sprintf(`# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/%s cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/bin/%s .

# Copy migrations if they exist
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./%s"]
`, projectName, projectName, projectName)

	if err := writeFile(filepath.Join(projectName, "Dockerfile"), dockerfileContent); err != nil {
		fmt.Printf("Error creating Dockerfile: %v\n", err)
	}

	// Docker Compose
	dockerComposeContent := fmt.Sprintf(`version: '3.8'

services:
  %s:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=database
      - DB_USER=%s
      - DB_PASSWORD=password
      - DB_NAME=%s
    depends_on:
      database:
        condition: service_healthy
    restart: unless-stopped

  database:
    image: %s
    environment:%s
    ports:
      - "%s:%s"
    volumes:
      - db_data:/var/lib/postgresql/data
    healthcheck:
      test: %s
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  db_data:
`, projectName, getDatabaseUser(database), projectName, getDatabaseImage(database), getDatabaseEnvVars(database, projectName), getDatabasePort(database), getDatabasePort(database), getDatabaseHealthCheck(database))

	if err := writeFile(filepath.Join(projectName, "docker-compose.yml"), dockerComposeContent); err != nil {
		fmt.Printf("Error creating docker-compose.yml: %v\n", err)
	}

	// .dockerignore
	dockerignoreContent := `# Git
.git
.gitignore

# Documentation
README.md
*.md

# Environment files
.env
.env.local
.env.example

# Build artifacts
bin/
*.exe

# Go files
*.mod
*.sum

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Logs
*.log

# Testing
coverage.out
coverage.html
`

	if err := writeFile(filepath.Join(projectName, ".dockerignore"), dockerignoreContent); err != nil {
		fmt.Printf("Error creating .dockerignore: %v\n", err)
	}
}

func getDatabaseImage(database string) string {
	switch database {
	case "mysql":
		return "mysql:8.0"
	case "mongodb":
		return "mongo:7.0"
	default:
		return "postgres:15"
	}
}

func getDatabaseEnvVars(database, projectName string) string {
	switch database {
	case "mysql":
		return fmt.Sprintf("\n      - MYSQL_ROOT_PASSWORD=password\n      - MYSQL_DATABASE=%s", projectName)
	case "mongodb":
		return fmt.Sprintf("\n      - MONGO_INITDB_ROOT_USERNAME=admin\n      - MONGO_INITDB_ROOT_PASSWORD=password\n      - MONGO_INITDB_DATABASE=%s", projectName)
	default:
		return fmt.Sprintf("\n      - POSTGRES_USER=postgres\n      - POSTGRES_PASSWORD=password\n      - POSTGRES_DB=%s", projectName)
	}
}

func getDatabaseHealthCheck(database string) string {
	switch database {
	case "mysql":
		return `["CMD", "mysqladmin", "ping", "-h", "localhost"]`
	case "mongodb":
		return `["CMD", "mongo", "--eval", "db.adminCommand('ping')"]`
	default:
		return `["CMD-SHELL", "pg_isready -U postgres"]`
	}
}

func init() {
	initCmd.Flags().StringP("module", "m", "", "Go module name (e.g: github.com/user/project)")
	initCmd.Flags().StringP("database", "d", "postgres", "Database type (postgres, mysql, mongodb)")
	initCmd.Flags().StringP("api", "a", "rest", "API type (rest, graphql, grpc)")
	initCmd.Flags().Bool("auth", false, "Include authentication system")
	initCmd.Flags().Bool("config", true, "Generate .goca.yaml configuration file")
	_ = initCmd.MarkFlagRequired("module")
}
