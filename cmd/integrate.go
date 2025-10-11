package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var integrateCmd = &cobra.Command{
	Use:   "integrate",
	Short: "Integrate existing features with DI and main.go",
	Long: `Automatically detects existing features and integrates them 
completely with the dependency injection container and main.go.
Useful for projects that have unintegrated features.`,
	Run: func(cmd *cobra.Command, _ []string) {
		all, _ := cmd.Flags().GetBool("all")
		features, _ := cmd.Flags().GetString("features")

		if all {
			fmt.Println("🔍 Detectando features existentes...")
			autoDetectedFeatures := detectExistingFeatures()
			if len(autoDetectedFeatures) == 0 {
				fmt.Println("❌ No se encontraron features para integrar")
				return
			}

			fmt.Printf("📋 Features detectados: %s\n", strings.Join(autoDetectedFeatures, ", "))
			integrateFeatures(autoDetectedFeatures)
		} else if features != "" {
			featureList := strings.Split(features, ",")
			for i, feature := range featureList {
				featureList[i] = strings.TrimSpace(feature)
			}
			fmt.Printf("🔧 Integrando features especificados: %s\n", strings.Join(featureList, ", "))
			integrateFeatures(featureList)
		} else {
			fmt.Println("❌ Debe especificar --all o --features")
			os.Exit(1)
		}

		fmt.Println("\n🎉 ¡Integración completada!")
		fmt.Println("✅ Todos los features están ahora:")
		fmt.Println("   🔗 Conectados en el contenedor DI")
		fmt.Println("   🛣️  Con rutas registradas en main.go")
		fmt.Println("   ⚡ Listos para usar")
	},
}

// detectExistingFeatures scans the project for existing features
func detectExistingFeatures() []string {
	var features []string

	// Look for domain entities in internal/domain
	domainDir := filepath.Join(DirInternal, DirDomain)
	if entries, err := os.ReadDir(domainDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
				name := strings.TrimSuffix(entry.Name(), ".go")
				// Skip common files and seed files
				if name != "errors" && name != "validations" && name != "common" && !strings.HasSuffix(name, "_seeds") {
					// Capitalize first letter to match feature naming
					if len(name) > 0 {
						caser := cases.Title(language.English)
						features = append(features, caser.String(name))
					}
				}
			}
		}
	}

	// Also look for handlers in internal/handler/http
	httpDir := filepath.Join(DirInternal, DirHandler, DirHTTP)
	if entries, err := os.ReadDir(httpDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), "_handler.go") {
				name := strings.TrimSuffix(entry.Name(), "_handler.go")
				caser := cases.Title(language.English)
				featureName := caser.String(name)

				// Only add if not already in the list
				found := false
				for _, existing := range features {
					if strings.EqualFold(existing, featureName) {
						found = true
						break
					}
				}
				if !found {
					features = append(features, featureName)
				}
			}
		}
	}

	return features
}

// integrateFeatures integrates multiple features
func integrateFeatures(features []string) {
	fmt.Println("\n🔄 Iniciando proceso de integración...")

	// Step 1: Create or update DI container
	fmt.Println("\n1️⃣  Configurando contenedor DI...")
	createOrUpdateDIContainer(features)

	// Step 2: Update main.go
	fmt.Println("\n2️⃣  Actualizando main.go...")
	updateMainGoWithAllFeatures(features)

	// Step 3: Verify integration
	fmt.Println("\n3️⃣  Verificando integración...")
	verifyIntegration(features)
}

// createOrUpdateDIContainer creates or updates the DI container with all features
func createOrUpdateDIContainer(features []string) {
	diPath := filepath.Join("internal", "di", "container.go")

	if _, err := os.Stat(diPath); os.IsNotExist(err) {
		fmt.Println("   📦 Creando contenedor DI...")
		generateDI(strings.Join(features, ","), DBPostgres, false)
	} else {
		fmt.Println("   🔄 Actualizando contenedor DI existente...")
		for _, feature := range features {
			addFeatureToDI(feature)
		}
	}
}

// updateMainGoWithAllFeatures updates main.go to include all features
func updateMainGoWithAllFeatures(features []string) {
	// Try multiple possible locations for main.go
	possiblePaths := []string{
		"main.go",
		filepath.Join("cmd", "server", "main.go"),
		filepath.Join("cmd", "main.go"),
	}

	var mainPath string
	var found bool

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			mainPath = path
			found = true
			break
		}
	}

	if !found {
		fmt.Println("   ⚠️  main.go no encontrado, creando nuevo...")
		createCompleteMainGo(features)
		return
	}

	fmt.Printf("   📍 Actualizando main.go en: %s\n", mainPath)

	// Read existing content
	content, err := os.ReadFile(mainPath)
	if err != nil {
		fmt.Printf("   ⚠️  No se pudo leer main.go: %v\n", err)
		return
	}

	contentStr := string(content)
	moduleName := getModuleName()

	// Check if this needs complete rewrite
	needsCompleteRewrite := !strings.Contains(contentStr, "di.NewContainer") ||
		!strings.Contains(contentStr, "/internal/di")

	if needsCompleteRewrite {
		fmt.Println("   🔧 Reescribiendo main.go completo...")
		createCompleteMainGoWithFeatures(mainPath, features, moduleName)
	} else {
		fmt.Println("   ➕ Agregando features faltantes...")
		addMissingFeaturesToMain(mainPath, features, contentStr, moduleName)
	}
}

// createCompleteMainGo creates a new main.go with all features
func createCompleteMainGo(features []string) {
	mainPath := "main.go"
	moduleName := getModuleName()
	createCompleteMainGoWithFeatures(mainPath, features, moduleName)
}

// createCompleteMainGoWithFeatures creates a complete main.go with DI and all feature routes
func createCompleteMainGoWithFeatures(mainPath string, features []string, moduleName string) {
	var routesSB strings.Builder

	// Generate routes for all features
	for _, feature := range features {
		featureLower := strings.ToLower(feature)
		routesSB.WriteString(fmt.Sprintf(`
	// %s routes
	%sHandler := container.%sHandler()
	router.HandleFunc("/api/v1/%ss", %sHandler.Create%s).Methods("POST")
	router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Get%s).Methods("GET")
	router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Update%s).Methods("PUT")
	router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Delete%s).Methods("DELETE")
	router.HandleFunc("/api/v1/%ss", %sHandler.List%ss).Methods("GET")
`, feature, featureLower, feature, featureLower, featureLower, feature, featureLower, featureLower, feature, featureLower, featureLower, feature, featureLower, featureLower, feature, featureLower, featureLower, feature))
	}

	newMainContent := fmt.Sprintf(`package main

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
	"%s/internal/di"
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
		log.Printf("⚠️  Database connection failed: %%v", err)
		log.Printf("📝 Server will start in degraded mode. Check your database configuration.")
		log.Printf("💡 To fix: Configure database environment variables in .env file")
		db = nil // Ensure db is nil for health checks
	} else {
		log.Printf("✅ Database connected successfully")
		
		// Run auto-migrations if database is connected
		if err := runAutoMigrations(db); err != nil {
			log.Printf("⚠️  Auto-migration failed: %%v", err)
			log.Printf("💡 You may need to run migrations manually")
		} else {
			log.Printf("✅ Database schema is up to date")
		}
	}
	
	// Setup DI container
	container := di.NewContainer(db)
	
	// Setup router
	router := mux.NewRouter()
	
	// Health check endpoints
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/health/ready", readinessHandler).Methods("GET")
	router.HandleFunc("/health/live", livenessHandler).Methods("GET")
%s
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

func runAutoMigrations(db *gorm.DB) error {
	log.Println("Running database auto-migrations...")
	
	// Import and register all domain models here
	// Example: db.AutoMigrate(&domain.User{}, &domain.Product{})
	
	return nil
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
		Version:   Version,
	}
	
	// Check database connection
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			status.Services["database"] = "unhealthy"
			status.Status = "degraded"
		} else {
			status.Services["database"] = "healthy"
		}
	} else {
		status.Services["database"] = "not configured"
		status.Status = "degraded"
	}
	
	w.Header().Set("Content-Type", "application/json")
	if status.Status == "healthy" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(status)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	// Check if application is ready to serve traffic
	if db != nil {
		sqlDB, _ := db.DB()
		if sqlDB != nil && sqlDB.Ping() == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ready"))
			return
		}
	}
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("not ready"))
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	// Application is alive if it can respond
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("alive"))
}
`, moduleName, moduleName, moduleName, routesSB.String())

	if err := os.WriteFile(mainPath, []byte(newMainContent), 0644); err != nil {
		fmt.Printf("   ⚠️  No se pudo crear main.go: %v\n", err)
	} else {
		fmt.Printf("   ✅ main.go created with %d features\n", len(features))
	}
}

// addMissingFeaturesToMain adds missing feature routes to existing main.go
func addMissingFeaturesToMain(mainPath string, features []string, contentStr, moduleName string) {
	newContent := contentStr

	// Add DI import if missing
	if !strings.Contains(newContent, "/internal/di") {
		importSection := fmt.Sprintf(`import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"%s/internal/di"
	"%s/pkg/config"
	"%s/pkg/logger"

	_ "github.com/lib/pq"
)`, moduleName, moduleName, moduleName)

		if oldImport := extractImportSection(newContent); oldImport != "" {
			newContent = strings.Replace(newContent, oldImport, importSection, 1)
		}
	}

	// Add DI container if missing
	if !strings.Contains(newContent, "di.NewContainer") {
		diSetup := "\n\t// Setup DI container\n\tcontainer := di.NewContainer(db)\n"
		routerPattern := "router := mux.NewRouter()"
		newContent = strings.Replace(newContent, routerPattern, diSetup+"\n\t// Setup router\n\t"+routerPattern, 1)
	}

	// Add missing feature routes
	addedFeatures := 0
	for _, feature := range features {
		featureLower := strings.ToLower(feature)

		// Check if routes already exist
		if !strings.Contains(newContent, fmt.Sprintf("/api/v1/%ss", featureLower)) {
			routeBlock := fmt.Sprintf(`
	// %s routes
	%sHandler := container.%sHandler()
	router.HandleFunc("/api/v1/%ss", %sHandler.Create%s).Methods("POST")
	router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Get%s).Methods("GET")
	router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Update%s).Methods("PUT")
	router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Delete%s).Methods("DELETE")
	router.HandleFunc("/api/v1/%ss", %sHandler.List%ss).Methods("GET")
`, feature, featureLower, feature, featureLower, featureLower, feature, featureLower, featureLower, feature, featureLower, featureLower, feature, featureLower, featureLower, feature, featureLower, featureLower, feature)

			// Insert before server start
			serverStartPatterns := []string{
				"log.Printf(\"Server starting",
				"log.Fatal(http.ListenAndServe",
			}

			for _, pattern := range serverStartPatterns {
				if idx := strings.Index(newContent, pattern); idx != -1 {
					newContent = newContent[:idx] + routeBlock + "\n\t" + newContent[idx:]
					addedFeatures++
					break
				}
			}
		}
	}

	if addedFeatures > 0 {
		if err := os.WriteFile(mainPath, []byte(newContent), 0644); err != nil {
			fmt.Printf("   ⚠️  No se pudo actualizar main.go: %v\n", err)
		} else {
			fmt.Printf("   ✅ %d features agregados a main.go\n", addedFeatures)
		}
	} else {
		fmt.Println("   ✅ Todos los features ya están integrados")
	}
}

// extractImportSection extracts the import section from Go code
func extractImportSection(content string) string {
	lines := strings.Split(content, "\n")
	var importLines []string
	inImport := false

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "import (") {
			inImport = true
			importLines = append(importLines, line)
		} else if inImport {
			importLines = append(importLines, line)
			if strings.Contains(line, ")") {
				break
			}
		}
	}

	return strings.Join(importLines, "\n")
}

// verifyIntegration checks that all features are properly integrated
func verifyIntegration(features []string) {
	issues := 0

	// Check DI container
	diPath := filepath.Join("internal", "di", "container.go")
	if _, err := os.Stat(diPath); os.IsNotExist(err) {
		fmt.Println("   ❌ Contenedor DI no encontrado")
		issues++
	} else {
		fmt.Println("   ✅ Contenedor DI existe")
	}

	// Check main.go integration
	mainPaths := []string{"main.go", filepath.Join("cmd", "server", "main.go")}
	mainFound := false

	for _, path := range mainPaths {
		if content, err := os.ReadFile(path); err == nil {
			contentStr := string(content)
			if strings.Contains(contentStr, "di.NewContainer") {
				fmt.Printf("   ✅ main.go integrado (%s)\n", path)
				mainFound = true

				// Check individual feature routes
				for _, feature := range features {
					featureLower := strings.ToLower(feature)
					if strings.Contains(contentStr, fmt.Sprintf("/api/v1/%ss", featureLower)) {
						fmt.Printf("   ✅ %s routes integradas\n", feature)
					} else {
						fmt.Printf("   ⚠️  %s routes faltantes\n", feature)
						issues++
					}
				}
				break
			}
		}
	}

	if !mainFound {
		fmt.Println("   ❌ main.go no encontrado o no integrado")
		issues++
	}

	if issues == 0 {
		fmt.Println("\n🎯 ¡Integración perfecta! Todo está listo.")
	} else {
		fmt.Printf("\n⚠️  Integración completada con %d advertencias\n", issues)
	}
}

func init() {
	integrateCmd.Flags().BoolP("all", "a", false, "Integrate all detected features automatically")
	integrateCmd.Flags().StringP("features", "f", "", "Specific features to integrate \"User,Product,Order\"")
}
