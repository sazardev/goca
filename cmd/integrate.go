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

		// Initialize safety manager
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		backup, _ := cmd.Flags().GetBool("backup")
		sm := NewSafetyManager(dryRun, force, backup)

		if dryRun {
			ui.DryRun("Previewing changes without creating files")
		}

		if all {
			ui.Info("Detecting existing features...")
			autoDetectedFeatures := detectExistingFeatures()
			if len(autoDetectedFeatures) == 0 {
				ui.Warning("No features found to integrate")
				return
			}

			ui.Info(fmt.Sprintf("Features detected: %s", strings.Join(autoDetectedFeatures, ", ")))
			integrateFeatures(autoDetectedFeatures, sm)
		} else if features != "" {
			featureList := strings.Split(features, ",")
			for i, feature := range featureList {
				featureList[i] = strings.TrimSpace(feature)
			}
			ui.Info(fmt.Sprintf("Integrating specified features: %s", strings.Join(featureList, ", ")))
			integrateFeatures(featureList, sm)
		} else {
			ui.Error("Must specify --all or --features")
			os.Exit(1)
		}

		if dryRun {
			sm.PrintSummary()
			return
		}

		ui.Success("Integration completed!")
		ui.Println("All features are now:")
		ui.Dim("   - Connected in the DI container")
		ui.Dim("   - With routes registered in main.go")
		ui.Dim("   - Ready to use")
	},
}

// detectExistingFeatures scans the project for COMPLETE features. A feature is
// only integrated when all four layers exist for it: a domain entity, a usecase
// service, a repository implementation, and an HTTP handler. Orphan domain
// entities (a domain/*.go with no usecase/repository/handler) are skipped, so
// the generated DI container never references constructors that do not exist.
func detectExistingFeatures() []string {
	var features []string

	domainDir := filepath.Join(DirInternal, DirDomain)
	entries, err := os.ReadDir(domainDir)
	if err != nil {
		return features
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".go")
		// Skip shared/common files, seed files and test files.
		if name == "errors" || name == "validations" || name == "common" ||
			strings.HasSuffix(name, "_seeds") || strings.HasSuffix(name, "_test") || name == "" {

			continue
		}

		if hasAllFeatureLayers(name) {
			// Reconstruct the PascalCase feature name from the snake_case file
			// name (e.g. user_profile.go -> UserProfile).
			features = append(features, snakeToPascal(name))
		} else {
			ui.Dim(fmt.Sprintf("   Skipping %s: incomplete feature (missing usecase/repository/handler)", snakeToPascal(name)))
		}
	}

	return features
}

// hasAllFeatureLayers reports whether the entity (named by its domain file's
// base name, snake_case or concatenated) has a usecase service, a repository
// implementation, and an HTTP handler.
//
// File naming differs across goca's generators (snake_case vs concatenated
// lowercase depending on the naming convention), so layer files are matched by
// their normalized key: lowercased with underscores/hyphens removed.
func hasAllFeatureLayers(entityFileBase string) bool {
	key := normalizeNameKey(entityFileBase)

	usecaseDir := filepath.Join(DirInternal, DirUseCase)
	repoDir := filepath.Join(DirInternal, DirRepository)
	httpDir := filepath.Join(DirInternal, DirHandler, DirHTTP)

	hasUseCase := dirHasLayerFile(usecaseDir, key, "_service.go", nil)
	hasHandler := dirHasLayerFile(httpDir, key, "_handler.go", nil)
	// Repository implementations are driver-prefixed
	// (postgres_/mysql_/mongo_<entity>_repository.go). Ignore the cache
	// decorator (cached_*): a base implementation is still required.
	hasRepo := dirHasLayerFile(repoDir, key, "_repository.go", func(fn string) bool {
		return strings.HasPrefix(fn, "cached_")
	})

	return hasUseCase && hasRepo && hasHandler
}

// repoDriverPrefixes are the driver tokens prepended to repository
// implementation filenames (e.g. postgres_user_repository.go).
var repoDriverPrefixes = []string{"postgres_", "postgresjson_", "mysql_", "mongo_", "sqlite_", "sqlserver_"}

// dirHasLayerFile reports whether dir contains a file that, after stripping the
// given suffix (and any leading repository driver prefix), normalizes to key.
// Files for which skip returns true are ignored.
func dirHasLayerFile(dir, key, suffix string, skip func(string) bool) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fn := strings.ToLower(e.Name())
		if skip != nil && skip(fn) {
			continue
		}
		if !strings.HasSuffix(fn, suffix) {
			continue
		}
		base := strings.TrimSuffix(fn, suffix)
		// Strip a leading driver prefix so postgres_user_repository.go matches
		// the "user" entity.
		for _, p := range repoDriverPrefixes {
			if strings.HasPrefix(base, p) {
				base = strings.TrimPrefix(base, p)
				break
			}
		}
		if normalizeNameKey(base) == key {
			return true
		}
	}
	return false
}

// normalizeNameKey lowercases a name and removes underscores and hyphens so
// "user_profile", "user-profile" and "userprofile" all compare equal.
func normalizeNameKey(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "-", "")
	return s
}

// integrateFeatures integrates multiple features.
func integrateFeatures(features []string, sm ...*SafetyManager) {
	ui.Blank()
	ui.Info("Starting integration process...")

	// Step 1: Create or update DI container
	ui.Step(1, "Configuring DI container...")
	createOrUpdateDIContainer(features, sm...)

	// Step 2: Update main.go
	ui.Step(2, "Updating main.go...")
	updateMainGoWithAllFeatures(features, sm...)

	// Step 3: Verify integration
	ui.Step(3, "Verifying integration...")
	verifyIntegration(features)
}

// snakeToPascal converts a snake_case identifier to PascalCase
// (e.g. "user_profile" -> "UserProfile", "user" -> "User").
func snakeToPascal(s string) string {
	caser := cases.Title(language.English)
	parts := strings.Split(s, "_")
	for i, p := range parts {
		parts[i] = caser.String(p)
	}
	return strings.Join(parts, "")
}

// createOrUpdateDIContainer creates or updates the DI container with all features.
func createOrUpdateDIContainer(features []string, sm ...*SafetyManager) {
	diPath := filepath.Join("internal", "di", "container.go")

	configIntegration := NewConfigIntegration()
	_ = configIntegration.LoadConfigForProject()
	database := configIntegration.GetDatabaseType("")

	if _, err := os.Stat(diPath); os.IsNotExist(err) {
		ui.Dim("   Creating DI container...")
		generateDI(strings.Join(features, ","), database, false, false, sm...)
	} else {
		ui.Dim("   Updating existing DI container...")
		for _, feature := range features {
			addFeatureToDI(feature, database, false)
		}
	}
}

// updateMainGoWithAllFeatures updates main.go to include all features.
func updateMainGoWithAllFeatures(features []string, sm ...*SafetyManager) {
	// Try multiple possible locations for main.go. The init-generated
	// entrypoint lives at cmd/server/main.go, so prefer it to avoid creating a
	// duplicate `package main` at the project root.
	possiblePaths := []string{
		filepath.Join("cmd", "server", "main.go"),
		filepath.Join("cmd", "main.go"),
		"main.go",
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
		ui.Warning("main.go not found, creating new...")
		createCompleteMainGo(features, sm...)
		return
	}

	ui.Dim(fmt.Sprintf("   Updating main.go at: %s", mainPath))

	// Read existing content
	content, err := os.ReadFile(mainPath)
	if err != nil {
		ui.Warning(fmt.Sprintf("Could not read main.go: %v", err))
		return
	}

	contentStr := string(content)
	moduleName := getModuleName()

	// A complete rewrite is only needed when main.go is unrecognizable (no mux
	// router to anchor injection). The standard init-generated main.go is
	// extended in place, which preserves its correct DB driver wiring.
	needsCompleteRewrite := !strings.Contains(contentStr, "mux.NewRouter()")

	if needsCompleteRewrite {
		ui.Dim("   Rewriting complete main.go...")
		createCompleteMainGoWithFeatures(mainPath, features, moduleName, sm...)
	} else {
		ui.Dim("   Wiring features into main.go...")
		addMissingFeaturesToMain(mainPath, features, contentStr, moduleName, sm...)
	}
}

// createCompleteMainGo creates a new main.go with all features.
//
// It writes to cmd/server/main.go (the init-generated entrypoint location) so
// the result lives alongside the canonical layout instead of creating a
// duplicate `package main` at the project root.
func createCompleteMainGo(features []string, sm ...*SafetyManager) {
	mainPath := filepath.Join("cmd", "server", "main.go")
	if err := os.MkdirAll(filepath.Dir(mainPath), 0o755); err != nil {
		ui.Warning(fmt.Sprintf("Could not create %s directory: %v", filepath.Dir(mainPath), err))
	}
	moduleName := getModuleName()
	createCompleteMainGoWithFeatures(mainPath, features, moduleName, sm...)
}

// createCompleteMainGoWithFeatures creates a complete main.go with DI and all feature routes.
func createCompleteMainGoWithFeatures(mainPath string, features []string, moduleName string, sm ...*SafetyManager) {
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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"%s/internal/di"
	"%s/pkg/config"
	"%s/pkg/logger"
)

type HealthStatus struct {
	Status    string            `+"`json:\"status\"`"+`
	Timestamp time.Time         `+"`json:\"timestamp\"`"+`
	Services  map[string]string `+"`json:\"services\"`"+`
	Version   string            `+"`json:\"version\"`"+`
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
		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Attempt %%d: Failed to open database connection: %%v", i+1, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		// Get underlying sql.DB for connection pool configuration
		sqlDB, err := conn.DB()
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
			return conn, nil
		}

		log.Printf("Attempt %%d: Database ping failed: %%v", i+1, err)
		sqlDB.Close()
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

	if err := writeFile(mainPath, newMainContent, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Could not create main.go: %v", err))
	} else {
		ui.Info(fmt.Sprintf("main.go created with %d features", len(features)))
	}
}

// addMissingFeaturesToMain adds missing feature routes to existing main.go.
func addMissingFeaturesToMain(mainPath string, features []string, contentStr, moduleName string, sm ...*SafetyManager) {
	newContent := contentStr
	changed := false

	// The generated DI container is GORM-based and takes a *gorm.DB. NoSQL
	// main.go files (e.g. MongoDB, which expose a mongoClient instead of db)
	// are not auto-wired here; doing so would reference an undefined `db` and
	// break compilation. Leave such projects untouched.
	if !strings.Contains(newContent, "*gorm.DB") {
		ui.Warning("main.go is not GORM-based; skipping automatic DI/route wiring")
		ui.Dim("   Wire the container and routes manually for this database.")
		return
	}

	// Add the DI import to the existing import block (before pkg/config) so the
	// project's real driver/imports are preserved.
	diImportPath := fmt.Sprintf("%s/internal/di", getImportPath(moduleName))
	if !strings.Contains(newContent, diImportPath) {
		cfgImport := fmt.Sprintf("\t\"%s/pkg/config\"", getImportPath(moduleName))
		if strings.Contains(newContent, cfgImport) {
			newContent = strings.Replace(newContent, cfgImport,
				fmt.Sprintf("\t\"%s\"\n%s", diImportPath, cfgImport), 1)
			changed = true
		}
	}

	// Create the DI container right after the router is created.
	if !strings.Contains(newContent, "di.NewContainer") {
		routerPattern := "router := mux.NewRouter()"
		if strings.Contains(newContent, routerPattern) {
			newContent = strings.Replace(newContent, routerPattern,
				routerPattern+"\n\n\t// Setup dependency injection container\n\tcontainer := di.NewContainer(db)", 1)
			changed = true
		}
	}

	// Register the domain entities with GORM auto-migration so the tables exist.
	// (Only for the GORM-based main.go, which contains the entities placeholder.)
	migratePlaceholder := "// Add domain entities here as they are created\n\t\t// Example: &domain.User{}, &domain.Product{}"
	if strings.Contains(newContent, migratePlaceholder) {
		var entitiesSB strings.Builder
		for i, feature := range features {
			if i > 0 {
				entitiesSB.WriteString("\n\t\t")
			}
			entitiesSB.WriteString(fmt.Sprintf("&domain.%s{},", feature))
		}
		newContent = strings.Replace(newContent, migratePlaceholder, entitiesSB.String(), 1)

		// Ensure the domain package is imported (Replace is a no-op if the
		// anchor import is absent).
		domainImportPath := fmt.Sprintf("%s/internal/domain", getImportPath(moduleName))
		if !strings.Contains(newContent, domainImportPath) {
			cfgImport := fmt.Sprintf("\t\"%s/pkg/config\"", getImportPath(moduleName))
			newContent = strings.Replace(newContent, cfgImport,
				fmt.Sprintf("\t\"%s\"\n%s", domainImportPath, cfgImport), 1)
		}
		changed = true
	}

	// Insert feature route blocks before the HTTP server setup.
	addedFeatures := 0
	for _, feature := range features {
		featureLower := strings.ToLower(feature)

		if strings.Contains(newContent, fmt.Sprintf("/api/v1/%ss", featureLower)) {
			continue
		}

		routeBlock := fmt.Sprintf(`	// %s routes
	%sHandler := container.%sHandler()
	router.HandleFunc("/api/v1/%ss", %sHandler.Create%s).Methods("POST")
	router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Get%s).Methods("GET")
	router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Update%s).Methods("PUT")
	router.HandleFunc("/api/v1/%ss/{id}", %sHandler.Delete%s).Methods("DELETE")
	router.HandleFunc("/api/v1/%ss", %sHandler.List%ss).Methods("GET")

`, feature, featureLower, feature, featureLower, featureLower, feature, featureLower, featureLower, feature, featureLower, featureLower, feature, featureLower, featureLower, feature, featureLower, featureLower, feature)

		// Anchor: the HTTP server setup comment present in the generated main.go.
		markers := []string{"// Setup HTTP server", "server := &http.Server{"}
		inserted := false
		for _, marker := range markers {
			if idx := strings.Index(newContent, marker); idx != -1 {
				newContent = newContent[:idx] + routeBlock + "\t" + newContent[idx:]
				addedFeatures++
				changed = true
				inserted = true
				break
			}
		}
		if !inserted {
			ui.Warning(fmt.Sprintf("Could not find an insertion point for %s routes", feature))
		}
	}

	if !changed {
		ui.Dim("   All features are already integrated")
		return
	}

	if err := writeGoFileMerged(mainPath, newContent, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Could not update main.go: %v", err))
		return
	}
	ui.Info(fmt.Sprintf("%d feature(s) wired into main.go", addedFeatures))
}

// extractImportSection extracts the import section from Go code.
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

// verifyIntegration checks that all features are properly integrated.
func verifyIntegration(features []string) {
	issues := 0

	// Check DI container
	diPath := filepath.Join("internal", "di", "container.go")
	if _, err := os.Stat(diPath); os.IsNotExist(err) {
		ui.Warning("DI container not found")
		issues++
	} else {
		ui.Dim("   DI container exists")
	}

	// Check main.go integration
	mainPaths := []string{"main.go", filepath.Join("cmd", "server", "main.go")}
	mainFound := false

	for _, path := range mainPaths {
		if content, err := os.ReadFile(path); err == nil {
			contentStr := string(content)
			if strings.Contains(contentStr, "di.NewContainer") {
				ui.Dim(fmt.Sprintf("   main.go integrated (%s)", path))
				mainFound = true

				// Check individual feature routes
				for _, feature := range features {
					featureLower := strings.ToLower(feature)
					if strings.Contains(contentStr, fmt.Sprintf("/api/v1/%ss", featureLower)) {
						ui.Dim(fmt.Sprintf("   %s routes integrated", feature))
					} else {
						ui.Warning(fmt.Sprintf("%s routes missing", feature))
						issues++
					}
				}
				break
			}
		}
	}

	if !mainFound {
		ui.Warning("main.go not found or not integrated")
		issues++
	}

	if issues == 0 {
		ui.Success("Perfect integration! Everything is ready.")
	} else {
		ui.Warning(fmt.Sprintf("Integration completed with %d warnings", issues))
	}
}

func init() {
	integrateCmd.Flags().BoolP("all", "a", false, "Integrate all detected features automatically")
	integrateCmd.Flags().StringP("features", "f", "", "Specific features to integrate \"User,Product,Order\"")
	integrateCmd.Flags().Bool("dry-run", false, "Preview changes without creating files")
	integrateCmd.Flags().Bool("force", false, "Overwrite existing files without confirmation")
	integrateCmd.Flags().Bool("backup", false, "Create backup of existing files before overwriting")
}
