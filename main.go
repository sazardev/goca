package main

import (
	"context"
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
	"github.com/sazardev/goca/internal/di"
	// "github.com/sazardev/goca/pkg/config"
	"github.com/sazardev/goca/internal/domain"
)

var db *gorm.DB

func main() {
	// Load configuration
	cfg := config.Load()

	log.Printf("Starting application")
	log.Printf("Environment: %s", cfg.Environment)
	
	// Connect to database with retry
	var err error
	db, err = connectToDatabase(cfg)
	if err != nil {
		log.Printf("⚠️  Database connection failed: %v", err)
		log.Printf("📝 Server will start in degraded mode. Check your database configuration.")
		log.Printf("💡 To fix: Configure database environment variables in .env file")
		db = nil // Ensure db is nil for health checks
	} else {
		log.Printf("✅ Database connected successfully")
		
		// Run auto-migrations if database is connected
		if err := runAutoMigrations(db); err != nil {
			log.Printf("⚠️  Auto-migration failed: %v", err)
			log.Printf("💡 You may need to run migrations manually")
		} else {
			log.Printf("✅ Database schema is up to date")
		}
	}
	
	// Setup DI container (even if db is nil, for degraded mode)
	container := di.NewContainer(db)

	// Setup router
	router := mux.NewRouter()
	
	// Health check endpoints with comprehensive checks
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/health/ready", readinessHandler).Methods("GET")
	router.HandleFunc("/health/live", livenessHandler).Methods("GET")

	// Cliente routes
	if db != nil {
		clienteHandler := container.ClienteHandler()
		router.HandleFunc("/api/v1/clientes", clienteHandler.CreateCliente).Methods("POST")
		router.HandleFunc("/api/v1/clientes/{id}", clienteHandler.GetCliente).Methods("GET")
		router.HandleFunc("/api/v1/clientes/{id}", clienteHandler.UpdateCliente).Methods("PUT")
		router.HandleFunc("/api/v1/clientes/{id}", clienteHandler.DeleteCliente).Methods("DELETE")
		router.HandleFunc("/api/v1/clientes", clienteHandler.ListClientes).Methods("GET")
	} else {
		// Degraded mode routes
		router.HandleFunc("/api/v1/clientes", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Service temporarily unavailable - database not connected", http.StatusServiceUnavailable)
		})
	}
	
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
		log.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %v", err)
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
		log.Printf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server exited")
}

func connectToDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDatabaseURL()
	
	log.Printf("Connecting to database at %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	
	// Check if this is development mode without database
	if cfg.Environment == "development" && cfg.Database.Password == "" {
		log.Println("⚠️  Development mode detected: No database password set")
		log.Println("📝 To connect to PostgreSQL, set environment variables:")
		log.Println("   DB_HOST=localhost")
		log.Println("   DB_PORT=5432") 
		log.Println("   DB_USER=postgres")
		log.Println("   DB_PASSWORD=your_password")
		log.Println("   DB_NAME=your_database")
		log.Println("🚀 Server will continue without database connection...")
		return nil, fmt.Errorf("development mode: database not configured")
	}
	
	// Retry connection up to 5 times
	for i := 0; i < 5; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Attempt %d: Failed to open database connection: %v", i+1, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		
		// Get underlying sql.DB for connection pool configuration
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Attempt %d: Failed to get underlying SQL DB: %v", i+1, err)
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
		
		log.Printf("Attempt %d: Database ping failed: %v", i+1, err)
		sqlDBClose, _ := db.DB()
		if sqlDBClose != nil {
			sqlDBClose.Close()
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	
	return nil, fmt.Errorf("failed to connect to database after 5 attempts: %w", err)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := "healthy"
	httpStatus := http.StatusOK
	
	if db == nil {
		status = "degraded"
		// Still return 200 for basic health check in degraded mode
	} else if err := checkDatabase(); err != nil {
		status = "degraded"
		log.Printf("Database health check failed: %v", err)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	fmt.Fprintf(w, "{\"status\":\"cliente\",\"timestamp\":\"Cliente\"}", status, time.Now().Format(time.RFC3339))
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "{\"status\":\"not_ready\",\"reason\":\"database_not_connected\"}", http.StatusServiceUnavailable)
		return
	}
	
	if err := checkDatabase(); err != nil {
		http.Error(w, fmt.Sprintf("{\"status\":\"not_ready\",\"reason\":\"database_check_failed\",\"error\":\"%!s(MISSING)\"}", err.Error()), http.StatusServiceUnavailable)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"status\":\"ready\",\"timestamp\":\"%!s(MISSING)\"}", time.Now().Format(time.RFC3339))
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
		return fmt.Errorf("failed to get underlying sql DB: %w", err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	return sqlDB.PingContext(ctx)
}

func runAutoMigrations(database *gorm.DB) error {
	if database == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	log.Println("🔄 Running GORM auto-migrations...")
	
	// Create a slice of all domain entities to migrate
	entities := []interface{}{
		// Add domain entities here as they are created
		&domain.%!s(MISSING){},
	}
	
	// Run auto-migration for all entities
	for _, entity := range entities {
		if err := database.AutoMigrate(entity); err != nil {
			return fmt.Errorf("failed to auto-migrate entity %T: %w", entity, err)
		}
		log.Printf("✅ Auto-migrated entity: %T", entity)
	}
	
	log.Println("✅ GORM auto-migrations completed successfully")
	return nil
}
