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

// detectExistingFeatures scans the project for existing features.
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

// createOrUpdateDIContainer creates or updates the DI container with all features.
func createOrUpdateDIContainer(features []string, sm ...*SafetyManager) {
	diPath := filepath.Join("internal", "di", "container.go")

	if _, err := os.Stat(diPath); os.IsNotExist(err) {
		ui.Dim("   Creating DI container...")
		generateDI(strings.Join(features, ","), DBPostgres, false, false, sm...)
	} else {
		ui.Dim("   Updating existing DI container...")
		for _, feature := range features {
			addFeatureToDI(feature, false)
		}
	}
}

// updateMainGoWithAllFeatures updates main.go to include all features.
func updateMainGoWithAllFeatures(features []string, sm ...*SafetyManager) {
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

	// Check if this needs complete rewrite
	needsCompleteRewrite := !strings.Contains(contentStr, "di.NewContainer") ||
		!strings.Contains(contentStr, "/internal/di")

	if needsCompleteRewrite {
		ui.Dim("   Rewriting complete main.go...")
		createCompleteMainGoWithFeatures(mainPath, features, moduleName, sm...)
	} else {
		ui.Dim("   Adding missing features...")
		addMissingFeaturesToMain(mainPath, features, contentStr, moduleName, sm...)
	}
}

// createCompleteMainGo creates a new main.go with all features.
func createCompleteMainGo(features []string, sm ...*SafetyManager) {
	mainPath := "main.go"
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
	Version   string            `+"`"+`json:"version"`+"`"+"\n}\n\nvar (\n\t// Build information (set by build flags)\n\tVersion   = \"dev\"\n\tBuildTime = \"unknown\"\n\tdb        *gorm.DB\n)\n\nfunc main() {\n\t// Load configuration\n\tcfg := config.Load()\n\t\n\t// Initialize logger\n\tlogger.Init()\n\t\n\tlog.Printf(\"Starting application v%%s (built: %%s)\", Version, BuildTime)\n\tlog.Printf(\"Environment: %%s\", cfg.Environment)\n\t\n\t// Connect to database with retry\n\tvar err error\n\tdb, err = connectToDatabase(cfg)\n\tif err != nil {\n\t\tlog.Printf(\"Warning: Database connection failed: %%v\", err)\n\t\tlog.Printf(\"Server will start in degraded mode. Check your database configuration.\")\n\t\tlog.Printf(\"Tip: Configure database environment variables in .env file\")\n\t\tdb = nil // Ensure db is nil for health checks\n\t} else {\n\t\tlog.Printf(\"Database connected successfully\")\n\t\t\n\t\t// Run auto-migrations if database is connected\n\t\tif err := runAutoMigrations(db); err != nil {\n\t\t\tlog.Printf(\"Warning: Auto-migration failed: %%v\", err)\n\t\t\tlog.Printf(\"Tip: You may need to run migrations manually\")\n\t\t} else {\n\t\t\tlog.Printf(\"Database schema is up to date\")\n\t\t}\n\t\n\t// Setup DI container\n\t:= di.NewContainer(db)\n\t\n\t// Setup router\n\t:= mux.NewRouter()\n\t\n\t// Health check endpoints\n\trouter.HandleFunc(\"/health\", healthCheckHandler).Methods(\"GET\")\n\trouter.HandleFunc(\"/health/ready\", readinessHandler).Methods(\"GET\")\n\trouter.HandleFunc(\"/health/live\", livenessHandler).Methods(\"GET\")\n%s\n\t// Setup HTTP server with timeouts\n\tserver := &http.Server{\n\t\tAddr:         \":\" + cfg.Port,\n\t\tHandler:      router,\n\t\tReadTimeout:  cfg.Server.ReadTimeout,\n\t\tWriteTimeout: cfg.Server.WriteTimeout,\n\t\tIdleTimeout:  cfg.Server.IdleTimeout,\n\t}\n\t\n\t// Start server in goroutine\n\tgo func() {\n\t\tlog.Printf(\"Server starting on port %%s\", cfg.Port)\n\t\tif err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {\n\t\t\tlog.Fatalf(\"Server startup failed: %%v\", err)\n\t\t}\n\t}()\n\t\n\t// Wait for interrupt signal to gracefully shutdown\n\tquit := make(chan os.Signal, 1)\n\tsignal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)\n\t<-quit\n\t\n\tlog.Println(\"Shutting down server...\")\n\t\n\t// Graceful shutdown with timeout\n\tctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)\n\tdefer cancel()\n\t\n\tif err := server.Shutdown(ctx); err != nil {\n\t\tlog.Printf(\"Server forced to shutdown: %%v\", err)\n\t}\n\t\n\tlog.Println(\"Server exited\")\n}\n\nfunc connectToDatabase(cfg *config.Config) (*gorm.DB, error) {\n\tdsn := cfg.GetDatabaseURL()\n\t\n\tlog.Printf(\"Connecting to database at %%s:%%s/%%s\", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)\n\t\n\t// Check if this is development mode without database\n\tif cfg.Environment == \"development\" && cfg.Database.Password == \"\" {\n\t\tlog.Println(\"Warning: Development mode detected: No database password set\")\n\t\tlog.Println(\"To connect to PostgreSQL, set environment variables:\")\n\t\tlog.Println(\"   DB_HOST=localhost\")\n\t\tlog.Println(\"   DB_PORT=5432\") \n\t\tlog.Println(\"   DB_USER=postgres\")\n\t\tlog.Println(\"   DB_PASSWORD=your_password\")\n\t\tlog.Println(\"   DB_NAME=your_database\")\n\t\tlog.Println(\"Server will continue without database connection...\")\n\t\treturn nil, fmt.Errorf(\"development mode: database not configured\")\n\t}\n\t\n\t// Retry connection up to 5 times\n\tfor i := 0; i < 5; i++ {\n\t\tdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})\n\t\tif err != nil {\n\t\t\tlog.Printf(\"Attempt %%d: Failed to open database connection: %%v\", i+1, err)\n\t\t\ttime.Sleep(time.Duration(i+1) * time.Second)\n\t\t\tcontinue\n\t\t}\n\t\t\n\t\t// Get underlying sql.DB for connection pool configuration\n\t\tsqlDB, err := db.DB()\n\t\tif err != nil {\n\t\t\tlog.Printf(\"Attempt %%d: Failed to get underlying SQL DB: %%v\", i+1, err)\n\t\t\ttime.Sleep(time.Duration(i+1) * time.Second)\n\t\t\tcontinue\n\t\t}\n\t\t\n\t\t// Configure connection pool\n\t\tsqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)\n\t\tsqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)\n\t\tsqlDB.SetConnMaxLifetime(cfg.Database.MaxLifetime)\n\t\t\n\t\t// Test the connection\n\t\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n\t\terr = sqlDB.PingContext(ctx)\n\t\tcancel()\n\t\t\n\t\tif err == nil {\n\t\t\treturn db, nil\n\t\t}\n\t\t\n\t\tlog.Printf(\"Attempt %%d: Database ping failed: %%v\", i+1, err)\n\t\tsqlDBClose, _ := db.DB()\n\t\tif sqlDBClose != nil {\n\t\t\tsqlDBClose.Close()\n\t\t}\n\t\ttime.Sleep(time.Duration(i+1) * time.Second)\n\t}\n\t\n\treturn nil, fmt.Errorf(\"failed to connect to database after 5 attempts\")\n}\n\nfunc runAutoMigrations(db *gorm.DB) error {\n\tlog.Println(\"Running database auto-migrations...\")\n\t\n\t// Import and register all domain models here\n\t// Example: db.AutoMigrate(&domain.User{}, &domain.Product{})\n\t\n\treturn nil\n}\n\nfunc healthCheckHandler(w http.ResponseWriter, r *http.Request) {\n\tstatus := HealthStatus{\n\t\tStatus:    \"healthy\",\n\t\tTimestamp: time.Now(),\n\t\tServices:  make(map[string]string),\n\t\tVersion:   Version,\n\t}\n\t\n\t// Check database connection\n\tif db != nil {\n\t\tsqlDB, err := db.DB()\n\t\tif err != nil || sqlDB.Ping() != nil {\n\t\t\tstatus.Services[\"database\"] = \"unhealthy\"\n\t\t\tstatus.Status = \"degraded\"\n\t\t} else {\n\t\t\tstatus.Services[\"database\"] = \"healthy\"\n\t\t} else {\n\t\tstatus.Services[\"database\"] = \"not configured\"\n\t\tstatus.Status = \"degraded\"\n\t}\n\t\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\tif status.Status == \"healthy\" {\n\t\tw.WriteHeader(http.StatusOK)\n\t} else {\n\t\tw.WriteHeader(http.StatusServiceUnavailable)\n\t}\n\tjson.NewEncoder(w).Encode(status)\n}\n\nfunc readinessHandler(w http.ResponseWriter, r *http.Request) {\n\t// Check if application is ready to serve traffic\n\tif db != nil {\n\t\tsqlDB, _ := db.DB()\n\t\tif sqlDB != nil && sqlDB.Ping() == nil {\n\t\t\tw.WriteHeader(http.StatusOK)\n\t\t\tw.Write([]byte(\"ready\"))\n\t\t\treturn\n\t\t}\n\tw.WriteHeader(http.StatusServiceUnavailable)\n\tw.Write([]byte(\"not ready\"))\n}\n\nfunc livenessHandler(w http.ResponseWriter, r *http.Request) {\n\t// Application is alive if it can respond\n\tw.WriteHeader(http.StatusOK)\n\tw.Write([]byte(\"alive\"))", moduleName, moduleName, moduleName, routesSB.String())

	if err := writeFile(mainPath, newMainContent, sm...); err != nil {
		ui.Warning(fmt.Sprintf("Could not create main.go: %v", err))
	} else {
		ui.Info(fmt.Sprintf("main.go created with %d features", len(features)))
	}
}

// addMissingFeaturesToMain adds missing feature routes to existing main.go.
func addMissingFeaturesToMain(mainPath string, features []string, contentStr, moduleName string, sm ...*SafetyManager) {
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
		diSetup := "\n\t// Setup DI container\n\t:="
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
		if err := writeFile(mainPath, newContent, sm...); err != nil {
			ui.Warning(fmt.Sprintf("Could not update main.go: %v", err))
		} else {
			ui.Info(fmt.Sprintf("%d features added to main.go", addedFeatures))
		}
	} else {
		ui.Dim("   All features are already integrated")
	}
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
