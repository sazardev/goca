package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var featureCmd = &cobra.Command{
	Use:   "feature <name>",
	Short: "Generar feature completo con Clean Architecture",
	Long: `Genera todas las capas necesarias para un feature completo, 
incluyendo dominio, casos de uso, repositorio y handlers en una sola operación.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		featureName := args[0]

		fields, _ := cmd.Flags().GetString("fields")
		database, _ := cmd.Flags().GetString("database")
		handlers, _ := cmd.Flags().GetString("handlers")
		validation, _ := cmd.Flags().GetBool("validation")
		businessRules, _ := cmd.Flags().GetBool("business-rules")

		if fields == "" {
			fmt.Println("Error: --fields flag es requerido")
			return
		}

		fmt.Printf("🚀 Generando feature completo '%s'\n", featureName)
		fmt.Printf("📋 Campos: %s\n", fields)
		fmt.Printf("🗄️  Base de datos: %s\n", database)
		fmt.Printf("🌐 Handlers: %s\n", handlers)

		if validation {
			fmt.Println("✅ Incluyendo validaciones")
		}
		if businessRules {
			fmt.Println("🧠 Incluyendo reglas de negocio")
		}

		generateCompleteFeature(featureName, fields, database, handlers, validation, businessRules)

		// 6. Auto-integrate with DI and main.go
		fmt.Println("6️⃣  Integrando automáticamente...")
		autoIntegrateFeature(featureName, handlers)

		fmt.Printf("\n🎉 Feature '%s' generado e integrado exitosamente!\n", featureName)
		fmt.Println("\n📂 Estructura generada:")
		printFeatureStructure(featureName, handlers)

		fmt.Println("\n✅ ¡Todo listo! El feature ya está:")
		fmt.Println("   🔗 Conectado en el contenedor DI")
		fmt.Println("   🛣️  Rutas registradas en el servidor")
		fmt.Println("   ⚡ Listo para usar inmediatamente")

		fmt.Println("\n� Próximos pasos:")
		fmt.Println("   1. Ejecutar: go mod tidy")
		fmt.Printf("   2. Iniciar servidor: go run main.go\n")
		fmt.Printf("   3. Probar endpoints: curl http://localhost:8080/api/v1/%ss\n", strings.ToLower(featureName))

		fmt.Println("\n💡 Comandos útiles adicionales:")
		fmt.Println("   goca integrate --all     # Integrar features existentes")
		fmt.Printf("   goca feature Product --fields \"name:string,price:float64\"  # Agregar otro feature\n")
	},
}

func generateCompleteFeature(featureName, fields, database, handlers string, validation, businessRules bool) {
	fmt.Println("\n🔄 Generando capas...")

	// 1. Generate Entity (Domain layer)
	fmt.Println("1️⃣  Generando entidad de dominio...")
	generateEntity(featureName, fields, true, businessRules, false, false)

	// 2. Generate Use Case
	fmt.Println("2️⃣  Generando casos de uso...")
	generateUseCase(featureName+"UseCase", featureName, "create,read,update,delete,list", validation, false)

	// 3. Generate Repository
	fmt.Println("3️⃣  Generando repositorio...")
	generateRepository(featureName, database, false, true, false, false)

	// 4. Generate Handlers
	fmt.Println("4️⃣  Generando handlers...")
	handlerTypes := strings.Split(handlers, ",")
	for _, handlerType := range handlerTypes {
		handlerType = strings.TrimSpace(handlerType)
		fmt.Printf("   📡 Generando handler %s...\n", handlerType)
		generateHandler(featureName, handlerType, true, validation, handlerType == "http")
	}

	// 5. Generate Messages
	fmt.Println("5️⃣  Generando mensajes...")
	generateMessages(featureName, true, true, true)

	fmt.Println("✅ Todas las capas generadas exitosamente!")
}

func printFeatureStructure(featureName, handlers string) {
	featureLower := strings.ToLower(featureName)

	fmt.Printf(`%s/
├── domain/
│   ├── %s.go          # Entidad pura
│   ├── errors.go      # Errores de dominio
│   └── validations.go # Validaciones de negocio
├── usecase/
│   ├── dto.go              # DTOs de entrada/salida
│   ├── %s_usecase.go       # Interfaz de casos de uso
│   ├── %s_service.go       # Implementación de casos de uso
│   └── interfaces.go       # Contratos hacia otras capas
├── repository/
│   ├── interfaces.go       # Contratos de persistencia
│   └── postgres_%s_repo.go # Implementación PostgreSQL
├── handler/`, featureName, featureLower, featureLower, featureLower, featureLower)

	handlerTypes := strings.Split(handlers, ",")
	for _, handlerType := range handlerTypes {
		handlerType = strings.TrimSpace(handlerType)
		switch handlerType {
		case "http":
			fmt.Printf(`
│   ├── http/
│   │   ├── %s_handler.go   # Controlador HTTP
│   │   └── routes.go       # Rutas HTTP`, featureLower)
		case "grpc":
			fmt.Printf(`
│   ├── grpc/
│   │   ├── %s.proto        # Definición gRPC
│   │   └── %s_server.go    # Servidor gRPC`, featureLower, featureLower)
		case "cli":
			fmt.Printf(`
│   ├── cli/
│   │   └── %s_commands.go  # Comandos CLI`, featureLower)
		case "worker":
			fmt.Printf(`
│   ├── worker/
│   │   └── %s_worker.go    # Workers/Jobs`, featureLower)
		case "soap":
			fmt.Printf(`
│   ├── soap/
│   │   └── %s_client.go    # Cliente SOAP`, featureLower)
		}
	}

	fmt.Printf(`
└── messages/
    ├── errors.go       # Mensajes de error
    └── responses.go    # Mensajes de respuesta
`)
}

// autoIntegrateFeature automatically integrates the feature with DI and main.go
func autoIntegrateFeature(featureName, handlers string) {
	fmt.Println("   🔄 Actualizando contenedor DI...")
	updateDIContainer(featureName)

	fmt.Println("   🛣️  Registrando rutas HTTP...")
	if strings.Contains(handlers, "http") {
		updateMainRoutes(featureName)
	}

	fmt.Println("   ✅ Integración completada")
}

// updateDIContainer updates or creates DI container with new feature
func updateDIContainer(featureName string) {
	// Check if DI container exists
	diPath := filepath.Join("internal", "di", "container.go")

	if _, err := os.Stat(diPath); os.IsNotExist(err) {
		// DI doesn't exist, create it with this feature
		fmt.Printf("   📦 Creando contenedor DI para %s...\n", featureName)
		generateDI(featureName, "postgres", false)
	} else {
		// DI exists, update it to include new feature
		fmt.Printf("   🔄 Actualizando contenedor DI existente...\n")
		addFeatureToDI(featureName)
	}
}

// addFeatureToDI adds a new feature to existing DI container
func addFeatureToDI(featureName string) {
	diPath := filepath.Join("internal", "di", "container.go")

	// Read existing content
	content, err := os.ReadFile(diPath)
	if err != nil {
		fmt.Printf("   ⚠️  Could not read DI container: %v\n", err)
		return
	}

	contentStr := string(content)
	featureLower := strings.ToLower(featureName)

	// Check if feature already exists
	if strings.Contains(contentStr, fmt.Sprintf("%sRepo", featureLower)) {
		fmt.Printf("   ✅ %s ya está en el contenedor DI\n", featureName)
		return
	}

	fmt.Printf("   ➕ Agregando %s al contenedor DI...\n", featureName)

	// Add repository field
	repoField := fmt.Sprintf("\t%sRepo    repository.%sRepository\n", featureLower, featureName)
	contentStr = strings.Replace(contentStr, "\n\t// Use Cases", repoField+"\n\t// Use Cases", 1)

	// Add use case field
	ucField := fmt.Sprintf("\t%sUC    usecase.%sUseCase\n", featureLower, featureName)
	contentStr = strings.Replace(contentStr, "\n\t// Handlers", ucField+"\n\t// Handlers", 1)

	// Add handler field
	fieldName := strings.ToLower(featureName[:1]) + featureName[1:] // camelCase
	handlerField := fmt.Sprintf("\t%sHandler    *http.%sHandler\n", fieldName, featureName)
	contentStr = strings.Replace(contentStr, "}\n\nfunc NewContainer", handlerField+"}\n\nfunc NewContainer", 1)

	// Add repository setup
	repoSetup := fmt.Sprintf("\tc.%sRepo = repository.NewPostgres%sRepository(c.db)\n", featureLower, featureName)
	setupRepoEnd := "}\n\nfunc (c *Container) setupUseCases() {"
	contentStr = strings.Replace(contentStr, setupRepoEnd, repoSetup+setupRepoEnd, 1)

	// Add use case setup
	ucSetup := fmt.Sprintf("\tc.%sUC = usecase.New%sService(c.%sRepo)\n", featureLower, featureName, featureLower)
	setupUCEnd := "}\n\nfunc (c *Container) setupHandlers() {"
	contentStr = strings.Replace(contentStr, setupUCEnd, ucSetup+setupUCEnd, 1)

	// Add handler setup
	handlerSetup := fmt.Sprintf("\tc.%sHandler = http.New%sHandler(c.%sUC)\n", fieldName, featureName, featureLower)
	setupHandlerEnd := "}\n\n// Getters"
	contentStr = strings.Replace(contentStr, setupHandlerEnd, handlerSetup+setupHandlerEnd, 1)

	// Add getters
	getters := fmt.Sprintf(`func (c *Container) %sHandler() *http.%sHandler {
	return c.%sHandler
}

func (c *Container) %sUseCase() usecase.%sUseCase {
	return c.%sUC
}

func (c *Container) %sRepository() repository.%sRepository {
	return c.%sRepo
}

`, featureName, featureName, fieldName, featureName, featureName, featureLower, featureName, featureName, featureLower)

	// Add getters at the end
	contentStr = contentStr + getters

	// Write updated content
	if err := os.WriteFile(diPath, []byte(contentStr), 0644); err != nil {
		fmt.Printf("   ⚠️  Could not update DI container: %v\n", err)
		return
	}

	fmt.Printf("   ✅ %s integrado en el contenedor DI\n", featureName)
}

// updateMainRoutes updates main.go to include new feature routes
func updateMainRoutes(featureName string) {
	// Try multiple possible locations for main.go
	possiblePaths := []string{
		"main.go", // Root directory (default from init)
		filepath.Join("cmd", "server", "main.go"), // Alternative location
		filepath.Join("cmd", "main.go"),           // Another common location
	}

	var mainPath string
	var found bool

	// Find main.go in one of the possible locations
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			mainPath = path
			found = true
			break
		}
	}

	if !found {
		fmt.Println("   ⚠️  main.go not found in any expected location, skipping route registration")
		fmt.Println("   💡 You can manually add the routes to your main.go file")
		printManualIntegrationInstructions(featureName)
		return
	}

	fmt.Printf("   📍 Encontrado main.go en: %s\n", mainPath)

	// Read existing content
	content, err := os.ReadFile(mainPath)
	if err != nil {
		fmt.Printf("   ⚠️  Could not read main.go: %v\n", err)
		printManualIntegrationInstructions(featureName)
		return
	}

	contentStr := string(content)
	featureLower := strings.ToLower(featureName)

	// Check if feature routes already exist
	if strings.Contains(contentStr, fmt.Sprintf("/%ss", featureLower)) {
		fmt.Println("   ✅ Las rutas ya están registradas")
		return
	}

	// Get module name
	moduleName := getModuleName()
	if moduleName == "" {
		fmt.Println("   ⚠️  Could not determine module name from go.mod")
		printManualIntegrationInstructions(featureName)
		return
	}

	// Check if this is a basic main.go that needs complete setup
	needsCompleteSetup := !strings.Contains(contentStr, "di.NewContainer") &&
		!strings.Contains(contentStr, "internal/di")

	if needsCompleteSetup {
		fmt.Println("   🔧 Configurando main.go completo con DI...")
		if !updateMainGoWithCompleteSetup(mainPath, featureName, moduleName) {
			printManualIntegrationInstructions(featureName)
			return
		}
	} else {
		fmt.Println("   🔗 Agregando rutas al main.go existente...")
		if !updateMainGoWithRoutes(mainPath, featureName, moduleName, contentStr) {
			printManualIntegrationInstructions(featureName)
			return
		}
	}

	fmt.Println("   ✅ Rutas registradas exitosamente")
}

func init() {
	rootCmd.AddCommand(featureCmd)
	featureCmd.Flags().StringP("fields", "f", "", "Campos de la entidad \"field:type,field2:type\" (requerido)")
	featureCmd.Flags().StringP("database", "d", "postgres", "Tipo de base de datos (postgres, mysql, mongodb)")
	featureCmd.Flags().StringP("handlers", "", "http", "Tipos de handlers \"http,grpc,cli\"")
	featureCmd.Flags().BoolP("validation", "v", false, "Incluir validaciones en todas las capas")
	featureCmd.Flags().BoolP("business-rules", "b", false, "Incluir métodos de reglas de negocio")

	_ = featureCmd.MarkFlagRequired("fields")
}

// updateMainGoWithCompleteSetup replaces the basic main.go with a complete DI-integrated version
func updateMainGoWithCompleteSetup(mainPath, featureName, moduleName string) bool {
	featureLower := strings.ToLower(featureName)

	newMainContent := fmt.Sprintf(`package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"%s/internal/di"
	"%s/pkg/config"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Load configuration
	cfg := config.Load()

	log.Printf("Starting application")
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
		defer db.Close()
		log.Printf("✅ Database connected successfully")
		
		// Run auto-migrations if database is connected
		if err := runAutoMigrations(db); err != nil {
			log.Printf("⚠️  Auto-migration failed: %%v", err)
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

	// %s routes
	if db != nil {
		%sHandler := container.%sHandler()
		router.HandleFunc("/api/v1/%ss", %sHandler.Create%s).Methods("POST")
		router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Get%s).Methods("GET")
		router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Update%s).Methods("PUT")
		router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Delete%s).Methods("DELETE")
		router.HandleFunc("/api/v1/%ss", %sHandler.List%ss).Methods("GET")
	} else {
		// Degraded mode routes
		router.HandleFunc("/api/v1/%ss", func(w http.ResponseWriter, r *http.Request) {
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

func connectToDatabase(cfg *config.Config) (*sql.DB, error) {
	dbURL := cfg.GetDatabaseURL()
	
	log.Printf("Connecting to database at %%s:%%s/%%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	
	var db *sql.DB
	var err error
	
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
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			log.Printf("Attempt %%d: Failed to open database connection: %%v", i+1, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		
		// Configure connection pool
		db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
		db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
		db.SetConnMaxLifetime(cfg.Database.MaxLifetime)
		
		// Test the connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = db.PingContext(ctx)
		cancel()
		
		if err == nil {
			return db, nil
		}
		
		log.Printf("Attempt %%d: Database ping failed: %%v", i+1, err)
		db.Close()
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	
	return nil, fmt.Errorf("failed to connect to database after 5 attempts: %%w", err)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := "healthy"
	httpStatus := http.StatusOK
	
	if db == nil {
		status = "degraded"
		// Still return 200 for basic health check in degraded mode
	} else if err := checkDatabase(); err != nil {
		status = "degraded"
		log.Printf("Database health check failed: %%v", err)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	fmt.Fprintf(w, "{\"status\":\"%s\",\"timestamp\":\"%s\"}", status, time.Now().Format(time.RFC3339))
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "{\"status\":\"not_ready\",\"reason\":\"database_not_connected\"}", http.StatusServiceUnavailable)
		return
	}
	
	if err := checkDatabase(); err != nil {
		http.Error(w, fmt.Sprintf("{\"status\":\"not_ready\",\"reason\":\"database_check_failed\",\"error\":\"%s\"}", err.Error()), http.StatusServiceUnavailable)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"status\":\"ready\",\"timestamp\":\"%s\"}", time.Now().Format(time.RFC3339))
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
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	return db.PingContext(ctx)
}

func runAutoMigrations(database *sql.DB) error {
	if database == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	// Check if migrations table exists
	createMigrationsTable := "CREATE TABLE IF NOT EXISTS schema_migrations (version VARCHAR(255) PRIMARY KEY, applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)"
	
	_, err := database.Exec(createMigrationsTable)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %%w", err)
	}
	
	// Run basic auto-migrations for generated features
	migrations := []struct {
		version string
		sql     string
	}{
		{
			version: "001_create_users_table",
			sql: "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, email VARCHAR(255) UNIQUE NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP); CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);",
		},
	}
	
	for _, migration := range migrations {
		// Check if migration already applied
		var count int
		err := database.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", migration.version).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %%w", err)
		}
		
		if count == 0 {
			// Apply migration
			_, err := database.Exec(migration.sql)
			if err != nil {
				return fmt.Errorf("failed to apply migration %%s: %%w", migration.version, err)
			}
			
			// Record migration as applied
			_, err = database.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.version)
			if err != nil {
				return fmt.Errorf("failed to record migration %%s: %%w", migration.version, err)
			}
			
			log.Printf("✅ Applied migration: %%s", migration.version)
		}
	}
	
	return nil
}
`, moduleName, moduleName, featureName, featureLower, featureName, featureLower, featureLower, featureName, featureLower, featureLower, featureName, featureLower, featureLower, featureName, featureLower, featureLower, featureName, featureLower, featureLower, featureName, featureLower)

	if err := os.WriteFile(mainPath, []byte(newMainContent), 0644); err != nil {
		fmt.Printf("   ⚠️  Could not update main.go: %v\n", err)
		return false
	}

	return true
}

// updateMainGoWithRoutes adds routes to an existing main.go with DI setup
func updateMainGoWithRoutes(mainPath, featureName, moduleName, contentStr string) bool {
	featureLower := strings.ToLower(featureName)

	// Add import for DI if not present
	importPath := getImportPath(moduleName)
	if !strings.Contains(contentStr, fmt.Sprintf("\"%s/internal/di\"", importPath)) {
		// Try to add the import
		importPattern := "import ("
		diImport := fmt.Sprintf("import (\n\t\"database/sql\"\n\t\"log\"\n\t\"net/http\"\n\n\t\"github.com/gorilla/mux\"\n\t\"%s/internal/di\"\n\t\"%s/pkg/config\"\n\t\"%s/pkg/logger\"\n\n\t_ \"github.com/lib/pq\"\n)", importPath, importPath, importPath)

		if strings.Contains(contentStr, importPattern) {
			contentStr = strings.Replace(contentStr, importPattern, diImport[:len(importPattern)], 1)
			// Replace the rest after the opening
			afterImport := strings.SplitN(contentStr, importPattern, 2)[1]
			if closeIndex := strings.Index(afterImport, ")"); closeIndex != -1 {
				contentStr = strings.Replace(contentStr, importPattern+afterImport[:closeIndex+1], diImport, 1)
			}
		}
	}

	// Add DI container setup if not present
	if !strings.Contains(contentStr, "container := di.NewContainer(db)") {
		setupContainer := "\n\t// Setup DI container\n\tcontainer := di.NewContainer(db)\n"
		if strings.Contains(contentStr, "// Setup router") {
			contentStr = strings.Replace(contentStr, "// Setup router", setupContainer+"\t// Setup router", 1)
		} else if strings.Contains(contentStr, "router := mux.NewRouter()") {
			contentStr = strings.Replace(contentStr, "router := mux.NewRouter()", setupContainer+"\n\t// Setup router\n\trouter := mux.NewRouter()", 1)
		}
	}

	// Add feature routes
	routeRegistration := fmt.Sprintf(`
	// %s routes
	%sHandler := container.%sHandler()
	router.HandleFunc("/api/v1/%ss", %sHandler.Create%s).Methods("POST")
	router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Get%s).Methods("GET")
	router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Update%s).Methods("PUT")
	router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Delete%s).Methods("DELETE")
	router.HandleFunc("/api/v1/%ss", %sHandler.List%ss).Methods("GET")
`, featureName, featureLower, featureName, featureLower, featureLower, featureName, featureLower, featureLower, featureName, featureLower, featureLower, featureName, featureLower, featureLower, featureName, featureLower, featureLower, featureName)

	// Insert routes before server start
	serverStartPatterns := []string{
		"log.Printf(\"Server starting on port %s\", cfg.Port)",
		"log.Printf(\"Server starting on port %%s\", cfg.Port)",
		"log.Fatal(http.ListenAndServe(",
	}

	routeInserted := false
	for _, pattern := range serverStartPatterns {
		if strings.Contains(contentStr, pattern) {
			contentStr = strings.Replace(contentStr, pattern, routeRegistration+"\n\t"+pattern, 1)
			routeInserted = true
			break
		}
	}

	if !routeInserted {
		fmt.Println("   ⚠️  Could not find a place to insert routes")
		return false
	}

	// Write updated content
	if err := os.WriteFile(mainPath, []byte(contentStr), 0644); err != nil {
		fmt.Printf("   ⚠️  Could not update main.go: %v\n", err)
		return false
	}

	return true
}

// printManualIntegrationInstructions prints instructions for manual integration
func printManualIntegrationInstructions(featureName string) {
	featureLower := strings.ToLower(featureName)
	moduleName := getModuleName()

	fmt.Println("\n   📋 Instrucciones de integración manual:")
	fmt.Println("   1. Agregar import en main.go:")
	fmt.Printf("      \"%s/internal/di\"\n", moduleName)
	fmt.Println("\n   2. Agregar en main(), después de conectar la DB:")
	fmt.Println("      container := di.NewContainer(db)")
	fmt.Println("\n   3. Agregar las rutas del feature:")
	fmt.Printf("      %sHandler := container.%sHandler()\n", featureLower, featureName)
	fmt.Printf("      router.HandleFunc(\"/api/v1/%ss\", %sHandler.Create%s).Methods(\"POST\")\n", featureLower, featureLower, featureName)
	fmt.Printf("      router.HandleFunc(\"/api/v1/%ss/{id}\", %sHandler.Get%s).Methods(\"GET\")\n", featureLower, featureLower, featureName)
	fmt.Printf("      router.HandleFunc(\"/api/v1/%ss/{id}\", %sHandler.Update%s).Methods(\"PUT\")\n", featureLower, featureLower, featureName)
	fmt.Printf("      router.HandleFunc(\"/api/v1/%ss/{id}\", %sHandler.Delete%s).Methods(\"DELETE\")\n", featureLower, featureLower, featureName)
	fmt.Printf("      router.HandleFunc(\"/api/v1/%ss\", %sHandler.List%ss).Methods(\"GET\")\n", featureLower, featureLower, featureName)
}
