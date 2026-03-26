package cmd

import (
	"fmt"
	"path/filepath"
)

func createMainGo(projectName, module, database string, sm ...*SafetyManager) {
	// For MongoDB and other NoSQL databases, generate a different main.go
	if database == DBMongoDB {
		createMongoDBMainGo(projectName, module, sm...)
		return
	}

	// For DynamoDB and Elasticsearch, generate specific implementations
	if database == DBDynamoDB {
		createDynamoDBMainGo(projectName, module, sm...)
		return
	}

	if database == DBElasticsearch {
		createElasticsearchMainGo(projectName, module, sm...)
		return
	}

	// Determine database driver import based on database type (GORM databases)
	var dbDriverImport string
	var dbDriverPackage string
	var requiresGorm bool

	switch database {
	case DBPostgres, DBPostgresJSON:
		dbDriverImport = `"gorm.io/driver/postgres"`
		dbDriverPackage = "postgres"
		requiresGorm = true
	case DBMySQL:
		dbDriverImport = `"gorm.io/driver/mysql"`
		dbDriverPackage = "mysql"
		requiresGorm = true
	case DBSQLite:
		dbDriverImport = `"gorm.io/driver/sqlite"`
		dbDriverPackage = "sqlite"
		requiresGorm = true
	case DBSQLServer:
		dbDriverImport = `"gorm.io/driver/sqlserver"`
		dbDriverPackage = "sqlserver"
		requiresGorm = true
	default:
		dbDriverImport = `"gorm.io/driver/sqlite"`
		dbDriverPackage = "sqlite"
		requiresGorm = true
	}

	// Build imports
	importLines := `"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"`

	if requiresGorm {
		importLines += `
	"gorm.io/gorm"`
	}

	importLines += fmt.Sprintf(`
	%s
	"%s/pkg/config"
	"%s/pkg/logger"`, dbDriverImport, module, module)

	// Generate main.go content with database-specific connection
	content := fmt.Sprintf(`package main

import (
	%s
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
		log.Println("To connect to database, set environment variables:")
		log.Println("   DB_HOST=localhost")
		log.Println("   DB_PORT=<port>") 
		log.Println("   DB_USER=<user>")
		log.Println("   DB_PASSWORD=your_password")
		log.Println("   DB_NAME=your_database")
		log.Println("Server will continue without database connection...")
		return nil, fmt.Errorf("development mode: database not configured")
	}
	
	// Retry connection up to 5 times
	for i := 0; i < 5; i++ {
		db, err := gorm.Open(%s.Open(dsn), &gorm.Config{})
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

`, importLines, dbDriverPackage)

	if err := writeGoFile(filepath.Join(projectName, "cmd", "server", "main.go"), content, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error writing main.go: %v", err))
		return
	}
}

func createMongoDBMainGo(projectName, module string, sm ...*SafetyManager) {
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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	mongoClient *mongo.Client
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Initialize logger
	logger.Init()
	
	log.Printf("Starting application v%%s (built: %%s)", Version, BuildTime)
	log.Printf("Environment: %%s", cfg.Environment)
	
	// Connect to MongoDB with retry
	var err error
	mongoClient, err = connectToMongoDB(cfg)
	if err != nil {
		log.Printf("Warning: MongoDB connection failed: %%v", err)
		log.Printf("Server will start in degraded mode. Check your database configuration.")
		log.Printf("Tip: Configure MongoDB environment variables in .env file")
		mongoClient = nil
	} else {
		log.Printf("MongoDB connected successfully")
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
	
	// Disconnect MongoDB
	if mongoClient != nil {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %%v", err)
		}
	}
	
	log.Println("Server exited")
}

func connectToMongoDB(cfg *config.Config) (*mongo.Client, error) {
	dsn := cfg.GetDatabaseURL()
	
	log.Printf("Connecting to MongoDB at %%s", cfg.Database.Host)
	
	// Check if this is development mode without database
	if cfg.Environment == "development" && cfg.Database.Password == "" {
		log.Println("Warning: Development mode detected: No database password set")
		log.Println("To connect to MongoDB, set environment variables:")
		log.Println("   DB_HOST=localhost")
		log.Println("   DB_PORT=27017")
		log.Println("   DB_USER=<user>")
		log.Println("   DB_PASSWORD=your_password")
		log.Println("   DB_NAME=your_database")
		log.Println("Server will continue without database connection...")
		return nil, fmt.Errorf("development mode: database not configured")
	}
	
	// Create MongoDB client options
	clientOptions := options.Client().ApplyURI(dsn)
	
	// Retry connection up to 5 times
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(ctx, clientOptions)
		
		if err != nil {
			cancel()
			log.Printf("Attempt %%d: Failed to connect to MongoDB: %%v", i+1, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		
		// Ping the database
		err = client.Ping(ctx, readpref.Primary())
		cancel()
		
		if err == nil {
			return client, nil
		}
		
		log.Printf("Attempt %%d: MongoDB ping failed: %%v", i+1, err)
		client.Disconnect(context.Background())
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	
	return nil, fmt.Errorf("failed to connect to MongoDB after 5 attempts")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
		Version:   Version,
	}
	
	// Check database
	if err := checkMongoDB(); err != nil {
		status.Status = "degraded"
		status.Services["database"] = fmt.Sprintf("error: %%v", err)
		log.Printf("MongoDB health check failed: %%v", err)
	} else {
		status.Services["database"] = "healthy"
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	// Check if all dependencies are ready
	if err := checkMongoDB(); err != nil {
		http.Error(w, fmt.Sprintf("MongoDB not ready: %%v", err), http.StatusServiceUnavailable)
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

func checkMongoDB() error {
	if mongoClient == nil {
		return fmt.Errorf("MongoDB client is nil")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	return mongoClient.Ping(ctx, readpref.Primary())
}
`, module, module)

	if err := writeGoFile(filepath.Join(projectName, "cmd", "server", "main.go"), content, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error writing main.go: %v", err))
		return
	}
}

func createDynamoDBMainGo(projectName, module string, sm ...*SafetyManager) {
	// Placeholder for DynamoDB implementation
	// For now, create a simple main.go that warns about missing implementation
	content := fmt.Sprintf(`package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"%s/pkg/config"
	"%s/pkg/logger"
)

type HealthStatus struct {
	Status    string            `+"`"+`json:"status"`+"`"+`
	Timestamp time.Time         `+"`"+`json:"timestamp"`+"`"+`
	Version   string            `+"`"+`json:"version"`+"`"+`
}

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	cfg := config.Load()
	logger.Init()
	
	log.Printf("Starting application v%%s (built: %%s)", Version, BuildTime)
	log.Println("Warning: DynamoDB integration is not yet fully implemented")
	log.Println("Please implement DynamoDB client initialization in main.go")
	
	router := mux.NewRouter()
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	
	go func() {
		log.Printf("Server starting on port %%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %%v", err)
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server exited")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   Version,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
`, module, module)

	if err := writeGoFile(filepath.Join(projectName, "cmd", "server", "main.go"), content, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error writing main.go: %v", err))
		return
	}
}

func createElasticsearchMainGo(projectName, module string, sm ...*SafetyManager) {
	// Placeholder for Elasticsearch implementation
	content := fmt.Sprintf(`package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"%s/pkg/config"
	"%s/pkg/logger"
)

type HealthStatus struct {
	Status    string            `+"`"+`json:"status"`+"`"+`
	Timestamp time.Time         `+"`"+`json:"timestamp"`+"`"+`
	Version   string            `+"`"+`json:"version"`+"`"+`
}

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	cfg := config.Load()
	logger.Init()
	
	log.Printf("Starting application v%%s (built: %%s)", Version, BuildTime)
	log.Println("Warning: Elasticsearch integration is not yet fully implemented")
	log.Println("Please implement Elasticsearch client initialization in main.go")
	
	router := mux.NewRouter()
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	
	go func() {
		log.Printf("Server starting on port %%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %%v", err)
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server exited")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   Version,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
`, module, module)

	if err := writeGoFile(filepath.Join(projectName, "cmd", "server", "main.go"), content, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Error writing main.go: %v", err))
		return
	}
}

