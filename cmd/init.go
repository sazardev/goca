package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var initCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Inicializar proyecto Clean Architecture",
	Long: `Crea la estructura base de un proyecto Go siguiendo los principios de Clean Architecture, 
incluyendo directorios, archivos de configuración y estructura de capas.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		module, _ := cmd.Flags().GetString("module")
		database, _ := cmd.Flags().GetString("database")
		auth, _ := cmd.Flags().GetBool("auth")
		api, _ := cmd.Flags().GetString("api")

		if module == "" {
			fmt.Println("Error: --module flag es requerido")
			os.Exit(1)
		}

		fmt.Printf("Inicializando proyecto '%s' con módulo '%s'\n", projectName, module)
		fmt.Printf("Base de datos: %s\n", database)
		fmt.Printf("API: %s\n", api)
		if auth {
			fmt.Println("Incluyendo autenticación")
		}

		createProjectStructure(projectName, module, database, auth, api)
		fmt.Printf("\n✅ Proyecto '%s' creado exitosamente!\n", projectName)
		fmt.Printf("📁 Directorio: ./%s\n", projectName)
		fmt.Println("\nPróximos pasos:")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Println("  go mod tidy")
		fmt.Println("  goca feature User --fields \"name:string,email:string\"")
	},
}

func createProjectStructure(projectName, module, database string, auth bool, api string) {
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

	// Create logger
	createLogger(projectName, module)

	if auth {
		createAuth(projectName, module)
	}

	// Download dependencies after creating go.mod
	if err := downloadDependencies(projectName); err != nil {
		fmt.Printf("⚠️  Warning: Failed to download dependencies: %v\n", err)
		fmt.Printf("💡 Run 'go mod download' manually in the project directory\n")
	}
}

func createGoMod(projectName, module, database string, auth bool) {
	var dependencies string

	// Base dependencies
	baseDeps := `github.com/gorilla/mux v1.8.0`

	switch database {
	case "mysql":
		baseDeps += `
	github.com/go-sql-driver/mysql v1.7.1`
	case "mongodb":
		baseDeps += `
	go.mongodb.org/mongo-driver v1.12.1`
	default: // postgres
		baseDeps += `
	github.com/lib/pq v1.10.9`
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

	writeFile(filepath.Join(projectName, "go.mod"), content)
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
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"%s/pkg/config"
	"%s/pkg/logger"
	
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Initialize logger
	logger.Init()
	
	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	// Setup router
	router := mux.NewRouter()
	
	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	log.Printf("Server starting on port %%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
`, module, module)

	writeGoFile(filepath.Join(projectName, "cmd", "server", "main.go"), content)
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
	writeFile(filepath.Join(projectName, ".gitignore"), content)
}

func createReadme(projectName, _ string) {
	content := fmt.Sprintf(`# %s

Proyecto generado con Goca - Go Clean Architecture Code Generator

## 🏗️ Arquitectura

Este proyecto sigue los principios de Clean Architecture:

- **Domain**: Entidades y reglas de negocio
- **Use Cases**: Lógica de aplicación
- **Repository**: Abstracción de datos
- **Handler**: Adaptadores de entrega

## 🚀 Inicio Rápido

1. Instalar dependencias:
`+"```bash\n"+`   go mod tidy
`+"```\n"+`

2. Configurar variables de entorno:
`+"```bash\n"+`   cp .env.example .env
`+"```\n"+`

3. Ejecutar la aplicación:
`+"```bash\n"+`   go run cmd/server/main.go
`+"```\n"+`

## 📁 Estructura del Proyecto

`+"```\n"+`%s/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   ├── usecase/
│   ├── repository/
│   └── handler/
├── pkg/
│   ├── config/
│   └── logger/
├── go.mod
└── README.md
`+"```\n"+`

## 🔧 Comandos Útiles

Generar un nuevo feature:
`+"```bash\n"+`goca feature User --fields "name:string,email:string"
`+"```\n"+`

Generar solo una entidad:
`+"```bash\n"+`goca entity Product --fields "name:string,price:float64"
`+"```\n"+`
`, cases.Title(language.English).String(projectName), projectName)

	writeFile(filepath.Join(projectName, "README.md"), content)
}

func createConfig(projectName, _, database string) {
	content := fmt.Sprintf(`package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	LogLevel    string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost/%s?sslmode=disable"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
`, projectName)

	writeGoFile(filepath.Join(projectName, "pkg", "config", "config.go"), content)
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
	writeGoFile(filepath.Join(projectName, "pkg", "logger", "logger.go"), content)
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
	writeGoFile(filepath.Join(projectName, "pkg", "auth", "jwt.go"), content)
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("module", "m", "", "Nombre del módulo Go (ej: github.com/user/project)")
	initCmd.Flags().StringP("database", "d", "postgres", "Tipo de base de datos (postgres, mysql, mongodb)")
	initCmd.Flags().StringP("api", "a", "rest", "Tipo de API (rest, graphql, grpc)")
	initCmd.Flags().Bool("auth", false, "Incluir sistema de autenticación")
	_ = initCmd.MarkFlagRequired("module")
}
