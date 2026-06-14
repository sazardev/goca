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

	// Password-based drivers skip connecting in development when no password
	// is set (avoids connection-retry spam); file-based SQLite always connects.
	degradedBlock := ""
	if database != DBSQLite {
		degradedBlock = "// Check if this is development mode without database\n\tif cfg.Environment == \"development\" && cfg.Database.Password == \"\" {\n\t\tlog.Println(\"Warning: Development mode detected: No database password set\")\n\t\tlog.Println(\"To connect to database, set environment variables:\")\n\t\tlog.Println(\"   DB_HOST=localhost\")\n\t\tlog.Println(\"   DB_PORT=<port>\") \n\t\tlog.Println(\"   DB_USER=<user>\")\n\t\tlog.Println(\"   DB_PASSWORD=your_password\")\n\t\tlog.Println(\"   DB_NAME=your_database\")\n\t\tlog.Println(\"Server will continue without database connection...\")\n\t\treturn nil, fmt.Errorf(\"development mode: database not configured\")\n\t}"
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
	Version   string            `+"`"+`json:"version"`+"`"+"\n}\n\nvar (\n\t// Build information (set by build flags)\n\tVersion   = \"dev\"\n\tBuildTime = \"unknown\"\n\tdb        *gorm.DB\n)\n\nfunc main() {\n\t// Load configuration\n\tcfg := config.Load()\n\t\n\t// Initialize logger\n\tlogger.Init()\n\t\n\tlog.Printf(\"Starting application v%%s (built: %%s)\", Version, BuildTime)\n\tlog.Printf(\"Environment: %%s\", cfg.Environment)\n\t\n\t// Connect to database with retry\n\tvar err error\n\tdb, err = connectToDatabase(cfg)\n\tif err != nil {\n\t\tlog.Printf(\"Warning: Database connection failed: %%v\", err)\n\t\tlog.Printf(\"Server will start in degraded mode. Check your database configuration.\")\n\t\tlog.Printf(\"Tip: Configure database environment variables in .env file\")\n\t\tdb = nil // Ensure db is nil for health checks\n\t} else {\n\t\tlog.Printf(\"Database connected successfully\")\n\t\t\n\t\t// Run auto-migrations if database is connected\n\t\tif err := runAutoMigrations(db); err != nil {\n\t\t\tlog.Printf(\"Warning: Auto-migration failed: %%v\", err)\n\t\t\tlog.Printf(\"Tip: You may need to run migrations manually\")\n\t\t} else {\n\t\t\tlog.Printf(\"Database schema is up to date\")\n\t\t}\n\t}\n\t\n\t// Setup router\n\trouter := mux.NewRouter()\n\t\n\t// Health check endpoint with comprehensive checks\n\trouter.HandleFunc(\"/health\", healthCheckHandler).Methods(\"GET\")\n\trouter.HandleFunc(\"/health/ready\", readinessHandler).Methods(\"GET\")\n\trouter.HandleFunc(\"/health/live\", livenessHandler).Methods(\"GET\")\n\t\n\t// Setup HTTP server with timeouts\n\tserver := &http.Server{\n\t\tAddr:         \":\" + cfg.Port,\n\t\tHandler:      router,\n\t\tReadTimeout:  cfg.Server.ReadTimeout,\n\t\tWriteTimeout: cfg.Server.WriteTimeout,\n\t\tIdleTimeout:  cfg.Server.IdleTimeout,\n\t}\n\t\n\t// Start server in goroutine\n\tgo func() {\n\t\tlog.Printf(\"Server starting on port %%s\", cfg.Port)\n\t\tif err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {\n\t\t\tlog.Fatalf(\"Server startup failed: %%v\", err)\n\t\t}\n\t}()\n\t\n\t// Wait for interrupt signal to gracefully shutdown\n\tquit := make(chan os.Signal, 1)\n\tsignal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)\n\t<-quit\n\t\n\tlog.Println(\"Shutting down server...\")\n\t\n\t// Graceful shutdown with timeout\n\tctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)\n\tdefer cancel()\n\t\n\tif err := server.Shutdown(ctx); err != nil {\n\t\tlog.Printf(\"Server forced to shutdown: %%v\", err)\n\t}\n\t\n\tlog.Println(\"Server exited\")\n}\n\nfunc connectToDatabase(cfg *config.Config) (*gorm.DB, error) {\n\tdsn := cfg.GetDatabaseURL()\n\t\n\tlog.Printf(\"Connecting to database at %%s:%%s/%%s\", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)\n\t\n\t%s\n\t\n\t// Retry connection up to 5 times\n\tfor i := 0; i < 5; i++ {\n\t\tdb, err := gorm.Open(%s.Open(dsn), &gorm.Config{})\n\t\tif err != nil {\n\t\t\tlog.Printf(\"Attempt %%d: Failed to open database connection: %%v\", i+1, err)\n\t\t\ttime.Sleep(time.Duration(i+1) * time.Second)\n\t\t\tcontinue\n\t\t}\n\t\t\n\t\t// Get underlying sql.DB for connection pool configuration\n\t\tsqlDB, err := db.DB()\n\t\tif err != nil {\n\t\t\tlog.Printf(\"Attempt %%d: Failed to get underlying SQL DB: %%v\", i+1, err)\n\t\t\ttime.Sleep(time.Duration(i+1) * time.Second)\n\t\t\tcontinue\n\t\t}\n\t\t\n\t\t// Configure connection pool\n\t\tsqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)\n\t\tsqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)\n\t\tsqlDB.SetConnMaxLifetime(cfg.Database.MaxLifetime)\n\t\t\n\t\t// Test the connection\n\t\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n\t\terr = sqlDB.PingContext(ctx)\n\t\tcancel()\n\t\t\n\t\tif err == nil {\n\t\t\treturn db, nil\n\t\t}\n\t\t\n\t\tlog.Printf(\"Attempt %%d: Database ping failed: %%v\", i+1, err)\n\t\tsqlDBClose, _ := db.DB()\n\t\tif sqlDBClose != nil {\n\t\t\tsqlDBClose.Close()\n\t\t}\n\t\ttime.Sleep(time.Duration(i+1) * time.Second)\n\t}\n\t\n\treturn nil, fmt.Errorf(\"failed to connect to database after 5 attempts\")\n}\n\nfunc healthCheckHandler(w http.ResponseWriter, r *http.Request) {\n\tstatus := HealthStatus{\n\t\tStatus:    \"healthy\",\n\t\tTimestamp: time.Now(),\n\t\tServices:  make(map[string]string),\n\t\tVersion:   Version,\n\t}\n\t\n\t// Check database\n\tif err := checkDatabase(); err != nil {\n\t\tstatus.Status = \"degraded\"\n\t\tstatus.Services[\"database\"] = fmt.Sprintf(\"error: %%v\", err)\n\t\t// Don't fail the whole health check for database issues in development\n\t\tlog.Printf(\"Database health check failed: %%v\", err)\n\t} else {\n\t\tstatus.Services[\"database\"] = \"healthy\"\n\t}\n\t\n\t// Always return 200 for basic health check - let readiness handle critical dependencies\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\tjson.NewEncoder(w).Encode(status)\n}\n\nfunc readinessHandler(w http.ResponseWriter, r *http.Request) {\n\t// Check if all dependencies are ready\n\tif err := checkDatabase(); err != nil {\n\t\thttp.Error(w, fmt.Sprintf(\"Database not ready: %%v\", err), http.StatusServiceUnavailable)\n\t\treturn\n\t}\n\t\n\tw.WriteHeader(http.StatusOK)\n\tw.Write([]byte(\"Ready\"))\n}\n\nfunc livenessHandler(w http.ResponseWriter, r *http.Request) {\n\t// Basic liveness check\n\tw.WriteHeader(http.StatusOK)\n\tw.Write([]byte(\"Alive\"))\n}\n\nfunc checkDatabase() error {\n\tif db == nil {\n\t\treturn fmt.Errorf(\"database connection is nil\")\n\t}\n\t\n\tsqlDB, err := db.DB()\n\tif err != nil {\n\t\treturn fmt.Errorf(\"failed to get underlying sql DB: %%w\", err)\n\t}\n\t\n\tctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)\n\tdefer cancel()\n\t\n\treturn sqlDB.PingContext(ctx)\n}\n\nfunc runAutoMigrations(database *gorm.DB) error {\n\tif database == nil {\n\t\treturn fmt.Errorf(\"database connection is nil\")\n\t}\n\t\n\t// Auto-migrate domain entities using GORM\n\tlog.Println(\"Running GORM auto-migrations...\")\n\t\n\t// Create a slice of all domain entities to migrate\n\tentities := []interface{}{\n\t\t// Add domain entities here as they are created\n\t\t// Example: &domain.User{}, &domain.Product{}\n\t}\n\t\n\t// Run auto-migration for all entities\n\tfor _, entity := range entities {\n\t\tif err := database.AutoMigrate(entity); err != nil {\n\t\t\treturn fmt.Errorf(\"failed to auto-migrate entity %%T: %%w\", entity, err)\n\t\t}\n\t}\n\t\n\t// For now, just ensure the connection works\n\tsqlDB, err := database.DB()\n\tif err != nil {\n\t\treturn fmt.Errorf(\"failed to get underlying SQL DB: %%w\", err)\n\t}\n\t\n\tif err := sqlDB.Ping(); err != nil {\n\t\treturn fmt.Errorf(\"database ping failed: %%w\", err)\n\t}\n\t\n\tlog.Println(\"GORM auto-migrations completed successfully\")\n\treturn nil\n}\n\n", importLines, degradedBlock, dbDriverPackage)

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
	Version   string            `+"`"+`json:"version"`+"`"+"\n}\n\nvar (\n\t// Build information (set by build flags)\n\tVersion   = \"dev\"\n\tBuildTime = \"unknown\"\n\tmongoClient *mongo.Client\n)\n\nfunc main() {\n\t// Load configuration\n\tcfg := config.Load()\n\t\n\t// Initialize logger\n\tlogger.Init()\n\t\n\tlog.Printf(\"Starting application v%%s (built: %%s)\", Version, BuildTime)\n\tlog.Printf(\"Environment: %%s\", cfg.Environment)\n\t\n\t// Connect to MongoDB with retry\n\tvar err error\n\tmongoClient, err = connectToMongoDB(cfg)\n\tif err != nil {\n\t\tlog.Printf(\"Warning: MongoDB connection failed: %%v\", err)\n\t\tlog.Printf(\"Server will start in degraded mode. Check your database configuration.\")\n\t\tlog.Printf(\"Tip: Configure MongoDB environment variables in .env file\")\n\t\tmongoClient = nil\n\t} else {\n\t\tlog.Printf(\"MongoDB connected successfully\")\n\t}\n\t\n\t// Setup router\n\trouter := mux.NewRouter()\n\t\n\t// Health check endpoint with comprehensive checks\n\trouter.HandleFunc(\"/health\", healthCheckHandler).Methods(\"GET\")\n\trouter.HandleFunc(\"/health/ready\", readinessHandler).Methods(\"GET\")\n\trouter.HandleFunc(\"/health/live\", livenessHandler).Methods(\"GET\")\n\t\n\t// Setup HTTP server with timeouts\n\tserver := &http.Server{\n\t\tAddr:         \":\" + cfg.Port,\n\t\tHandler:      router,\n\t\tReadTimeout:  cfg.Server.ReadTimeout,\n\t\tWriteTimeout: cfg.Server.WriteTimeout,\n\t\tIdleTimeout:  cfg.Server.IdleTimeout,\n\t}\n\t\n\t// Start server in goroutine\n\tgo func() {\n\t\tlog.Printf(\"Server starting on port %%s\", cfg.Port)\n\t\tif err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {\n\t\t\tlog.Fatalf(\"Server startup failed: %%v\", err)\n\t\t}\n\t}()\n\t\n\t// Wait for interrupt signal to gracefully shutdown\n\tquit := make(chan os.Signal, 1)\n\tsignal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)\n\t<-quit\n\t\n\tlog.Println(\"Shutting down server...\")\n\t\n\t// Graceful shutdown with timeout\n\tctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)\n\tdefer cancel()\n\t\n\tif err := server.Shutdown(ctx); err != nil {\n\t\tlog.Printf(\"Server forced to shutdown: %%v\", err)\n\t}\n\t\n\t// Disconnect MongoDB\n\tif mongoClient != nil {\n\t\tif err := mongoClient.Disconnect(ctx); err != nil {\n\t\t\tlog.Printf(\"Error disconnecting from MongoDB: %%v\", err)\n\t\t}\n\t}\n\t\n\tlog.Println(\"Server exited\")\n}\n\nfunc connectToMongoDB(cfg *config.Config) (*mongo.Client, error) {\n\tdsn := cfg.GetDatabaseURL()\n\t\n\tlog.Printf(\"Connecting to MongoDB at %%s\", cfg.Database.Host)\n\t\n\t// Check if this is development mode without database\n\tif cfg.Environment == \"development\" && cfg.Database.Password == \"\" {\n\t\tlog.Println(\"Warning: Development mode detected: No database password set\")\n\t\tlog.Println(\"To connect to MongoDB, set environment variables:\")\n\t\tlog.Println(\"   DB_HOST=localhost\")\n\t\tlog.Println(\"   DB_PORT=27017\")\n\t\tlog.Println(\"   DB_USER=<user>\")\n\t\tlog.Println(\"   DB_PASSWORD=your_password\")\n\t\tlog.Println(\"   DB_NAME=your_database\")\n\t\tlog.Println(\"Server will continue without database connection...\")\n\t\treturn nil, fmt.Errorf(\"development mode: database not configured\")\n\t}\n\t\n\t// Create MongoDB client options\n\tclientOptions := options.Client().ApplyURI(dsn)\n\t\n\t// Retry connection up to 5 times\n\tfor i := 0; i < 5; i++ {\n\t\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)\n\t\tclient, err := mongo.Connect(ctx, clientOptions)\n\t\t\n\t\tif err != nil {\n\t\t\tcancel()\n\t\t\tlog.Printf(\"Attempt %%d: Failed to connect to MongoDB: %%v\", i+1, err)\n\t\t\ttime.Sleep(time.Duration(i+1) * time.Second)\n\t\t\tcontinue\n\t\t}\n\t\t\n\t\t// Ping the database\n\t\terr = client.Ping(ctx, readpref.Primary())\n\t\tcancel()\n\t\t\n\t\tif err == nil {\n\t\t\treturn client, nil\n\t\t}\n\t\t\n\t\tlog.Printf(\"Attempt %%d: MongoDB ping failed: %%v\", i+1, err)\n\t\tclient.Disconnect(context.Background())\n\t\ttime.Sleep(time.Duration(i+1) * time.Second)\n\t}\n\t\n\treturn nil, fmt.Errorf(\"failed to connect to MongoDB after 5 attempts\")\n}\n\nfunc healthCheckHandler(w http.ResponseWriter, r *http.Request) {\n\tstatus := HealthStatus{\n\t\tStatus:    \"healthy\",\n\t\tTimestamp: time.Now(),\n\t\tServices:  make(map[string]string),\n\t\tVersion:   Version,\n\t}\n\t\n\t// Check database\n\tif err := checkMongoDB(); err != nil {\n\t\tstatus.Status = \"degraded\"\n\t\tstatus.Services[\"database\"] = fmt.Sprintf(\"error: %%v\", err)\n\t\tlog.Printf(\"MongoDB health check failed: %%v\", err)\n\t} else {\n\t\tstatus.Services[\"database\"] = \"healthy\"\n\t}\n\t\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\tjson.NewEncoder(w).Encode(status)\n}\n\nfunc readinessHandler(w http.ResponseWriter, r *http.Request) {\n\t// Check if all dependencies are ready\n\tif err := checkMongoDB(); err != nil {\n\t\thttp.Error(w, fmt.Sprintf(\"MongoDB not ready: %%v\", err), http.StatusServiceUnavailable)\n\t\treturn\n\t}\n\t\n\tw.WriteHeader(http.StatusOK)\n\tw.Write([]byte(\"Ready\"))\n}\n\nfunc livenessHandler(w http.ResponseWriter, r *http.Request) {\n\t// Basic liveness check\n\tw.WriteHeader(http.StatusOK)\n\tw.Write([]byte(\"Alive\"))\n}\n\nfunc checkMongoDB() error {\n\tif mongoClient == nil {\n\t\treturn fmt.Errorf(\"MongoDB client is nil\")\n\t}\n\t\n\tctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)\n\tdefer cancel()\n\t\n\treturn mongoClient.Ping(ctx, readpref.Primary())\n}\n", module, module)

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
