package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var handlerCmd = &cobra.Command{
	Use:   "handler <entity>",
	Short: "Generate handlers for different protocols",
	Long: `Creates delivery adapters that handle different protocols 
(HTTP, gRPC, CLI) maintaining layer separation.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entity := args[0]

		// Initialize configuration integration
		configIntegration := NewConfigIntegration()
		configIntegration.LoadConfigForProject()

		// Get CLI flag values
		handlerType, _ := cmd.Flags().GetString("type")
		middleware, _ := cmd.Flags().GetBool("middleware")
		validation, _ := cmd.Flags().GetBool("validation")
		swagger, _ := cmd.Flags().GetBool("swagger")

		// Merge only explicitly changed CLI flags with config
		flags := map[string]interface{}{}
		if cmd.Flags().Changed("type") {
			flags["handlers"] = handlerType
		}
		if cmd.Flags().Changed("middleware") {
			flags["middleware"] = middleware
		}
		if cmd.Flags().Changed("validation") {
			flags["validation"] = validation
		}
		if cmd.Flags().Changed("swagger") {
			flags["swagger"] = swagger
		}

		if len(flags) > 0 {
			configIntegration.MergeWithCLIFlags(flags)
		}

		// Calculate effective values (config overrides CLI defaults)
		effectiveHandlerType := handlerType
		if !cmd.Flags().Changed("type") && configIntegration.config != nil {
			handlers := configIntegration.GetHandlerTypes(handlerType)
			if len(handlers) > 0 {
				effectiveHandlerType = handlers[0]
			}
		}

		effectiveMiddleware := middleware
		effectiveValidation := validation
		if !cmd.Flags().Changed("validation") && configIntegration.config != nil {
			effectiveValidation = configIntegration.config.Generation.Validation.Enabled
		}

		effectiveSwagger := swagger
		if !cmd.Flags().Changed("swagger") && configIntegration.config != nil {
			effectiveSwagger = configIntegration.config.Generation.Documentation.Swagger.Enabled
		}

		// Get naming convention from config
		fileNamingConvention := "lowercase" // default
		if configIntegration.config != nil {
			fileNamingConvention = configIntegration.GetNamingConvention("file")
		}

		// Print configuration summary
		ui.Header(fmt.Sprintf("Generating handler '%s' for entity '%s'", effectiveHandlerType, entity))
		if configIntegration.config != nil {
			if !cmd.Flags().Changed("type") {
				handlers := configIntegration.GetHandlerTypes(handlerType)
				if len(handlers) > 0 {
					ui.KeyValueFromConfig("Handler type", effectiveHandlerType)
				}
			}
			if !cmd.Flags().Changed("validation") {
				ui.KeyValueFromConfig("Validation", fmt.Sprintf("%v", effectiveValidation))
			}
			if !cmd.Flags().Changed("swagger") {
				ui.KeyValueFromConfig("Swagger", fmt.Sprintf("%v", effectiveSwagger))
			}
		}

		if effectiveMiddleware {
			ui.Feature("Including middleware", false)
		}
		if effectiveValidation {
			ui.Feature("Including validation", false)
		}
		if effectiveSwagger && effectiveHandlerType == HandlerHTTP {
			ui.Feature("Including Swagger documentation", false)
		}

		// Initialize safety manager
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		backup, _ := cmd.Flags().GetBool("backup")
		sm := NewSafetyManager(dryRun, force, backup)

		if dryRun {
			ui.DryRun("Previewing changes without creating files")
		}

		generateHandler(entity, effectiveHandlerType, effectiveMiddleware, effectiveValidation, effectiveSwagger, fileNamingConvention, sm)

		if dryRun {
			sm.PrintSummary()
			return
		}

		// Add required dependencies
		projectRoot, _ := os.Getwd()
		depMgr := NewDependencyManager(projectRoot, false)
		features := map[string]bool{"validation": effectiveValidation}
		requiredDeps := depMgr.GetRequiredDependenciesForFeature(effectiveHandlerType, features)
		for _, dep := range requiredDeps {
			if err := depMgr.AddDependency(dep); err != nil {
				ui.Warning(fmt.Sprintf("Could not add dependency %s: %v", dep.Module, err))
			}
		}
		if len(requiredDeps) > 0 {
			if err := depMgr.UpdateGoMod(); err != nil {
				ui.Warning(fmt.Sprintf("Could not update go.mod: %v", err))
			}
		}

		ui.Success(fmt.Sprintf("Handler '%s' for '%s' generated successfully!", effectiveHandlerType, entity))
	},
}

func generateHandler(entity, handlerType string, middleware, validation, swagger bool, fileNamingConvention string, sm ...*SafetyManager) {
	switch handlerType {
	case HandlerHTTP:
		generateHTTPHandler(entity, middleware, validation, swagger, fileNamingConvention, sm...)
	case HandlerGRPC:
		generateGRPCHandler(entity, fileNamingConvention, sm...)
	case HandlerCLI:
		generateCLIHandler(entity, fileNamingConvention, sm...)
	case "worker":
		generateWorkerHandler(entity, fileNamingConvention, sm...)
	case "soap":
		generateSOAPHandler(entity, fileNamingConvention, sm...)
	default:
		ui.Error(fmt.Sprintf("Unsupported handler type: %s", handlerType))
		os.Exit(1)
	}
}

func generateHTTPHandler(entity string, middleware, validation, swagger bool, fileNamingConvention string, sm ...*SafetyManager) {
	// Create handlers directory if it doesn't exist
	handlerDir := filepath.Join(DirInternal, DirHandler, DirHTTP)
	_ = os.MkdirAll(handlerDir, 0755)

	// Generate handler file
	generateHTTPHandlerFile(handlerDir, entity, validation, fileNamingConvention, sm...)

	// Generate routes file
	generateHTTPRoutesFile(handlerDir, entity, middleware, sm...)

	// Generate DTOs for HTTP if validation is enabled
	if validation {
		generateHTTPDTOFile(handlerDir, entity, sm...)
	}

	// Generate Swagger docs if requested
	if swagger {
		generateSwaggerFile(handlerDir, entity, sm...)
	}
}

func generateHTTPHandlerFile(dir, entity string, validation bool, fileNamingConvention string, sm ...*SafetyManager) {
	// Apply naming convention to filename
	var filename string
	if fileNamingConvention == "snake_case" {
		filename = filepath.Join(dir, toSnakeCase(entity)+"_handler.go")
	} else if fileNamingConvention == "kebab-case" {
		filename = filepath.Join(dir, toKebabCase(entity)+"-handler.go")
	} else {
		filename = filepath.Join(dir, strings.ToLower(entity)+"_handler.go")
	}

	// Get the module name from go.mod
	moduleName := getModuleName()
	importPath := getImportPath(moduleName)

	var content strings.Builder
	content.WriteString("package " + DirHTTP + "\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"encoding/json\"\n")
	content.WriteString("\t\"net/http\"\n")
	content.WriteString("\t\"strconv\"\n\n")
	content.WriteString("\t\"github.com/gorilla/mux\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/usecase\"\n", importPath))
	if validation {
		content.WriteString("\t\"github.com/go-playground/validator/v10\"\n")
	}
	content.WriteString(")\n\n")

	// Handler struct
	handlerName := fmt.Sprintf("%sHandler", entity)
	content.WriteString(fmt.Sprintf("type %s struct {\n", handlerName))
	content.WriteString(fmt.Sprintf("\tusecase usecase.%sUseCase\n", entity))
	content.WriteString("}\n\n")

	// Constructor
	content.WriteString(fmt.Sprintf("func New%s(uc usecase.%sUseCase) *%s {\n",
		handlerName, entity, handlerName))
	content.WriteString(fmt.Sprintf("\treturn &%s{usecase: uc}\n", handlerName))
	content.WriteString("}\n\n")

	// Generate HTTP methods
	generateCreateHandlerMethod(&content, entity, handlerName, validation)
	generateGetHandlerMethod(&content, entity, handlerName)
	generateUpdateHandlerMethod(&content, entity, handlerName, validation)
	generateDeleteHandlerMethod(&content, entity, handlerName)
	generateListHandlerMethod(&content, entity, handlerName)

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing handler file: %v", err))
		return
	}
}

func generateCreateHandlerMethod(content *strings.Builder, entity, handlerName string, validation bool) {
	handlerVar := strings.ToLower(string(handlerName[0]))

	fmt.Fprintf(content, "func (%s *%s) Create%s(w http.ResponseWriter, r *http.Request) {\n",
		handlerVar, handlerName, entity)
	fmt.Fprintf(content, "\tvar input usecase.Create%sInput\n\n", entity)

	content.WriteString("\tif err := json.NewDecoder(r.Body).Decode(&input); err != nil {\n")
	content.WriteString("\t\thttp.Error(w, \"Invalid request body\", http.StatusBadRequest)\n")
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	if validation {
		content.WriteString("\tif err := validator.New().Struct(input); err != nil {\n")
		content.WriteString("\t\thttp.Error(w, err.Error(), http.StatusUnprocessableEntity)\n")
		content.WriteString("\t\treturn\n")
		content.WriteString("\t}\n\n")
	}

	fmt.Fprintf(content, "\toutput, err := %s.usecase.Create%s(input)\n", handlerVar, entity)
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString("\tw.Header().Set(\"Content-Type\", \"application/json\")\n")
	content.WriteString("\tw.WriteHeader(http.StatusCreated)\n")
	content.WriteString("\tjson.NewEncoder(w).Encode(output)\n")
	content.WriteString("}\n\n")
}

func generateGetHandlerMethod(content *strings.Builder, entity, handlerName string) {
	handlerVar := strings.ToLower(string(handlerName[0]))

	fmt.Fprintf(content, "func (%s *%s) Get%s(w http.ResponseWriter, r *http.Request) {\n",
		handlerVar, handlerName, entity)
	content.WriteString("\tvars := mux.Vars(r)\n")
	content.WriteString("\tid, err := strconv.Atoi(vars[\"id\"])\n")
	content.WriteString("\tif err != nil {\n")
	fmt.Fprintf(content, "\t\thttp.Error(w, \"Invalid %s ID\", http.StatusBadRequest)\n", strings.ToLower(entity))
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	fmt.Fprintf(content, "\t%s, err := %s.usecase.Get%s(id)\n", strings.ToLower(entity), handlerVar, entity)
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\thttp.Error(w, err.Error(), http.StatusNotFound)\n")
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString("\tw.Header().Set(\"Content-Type\", \"application/json\")\n")
	fmt.Fprintf(content, "\tjson.NewEncoder(w).Encode(%s)\n", strings.ToLower(entity))
	content.WriteString("}\n\n")
}

func generateUpdateHandlerMethod(content *strings.Builder, entity, handlerName string, validation bool) {
	handlerVar := strings.ToLower(string(handlerName[0]))

	fmt.Fprintf(content, "func (%s *%s) Update%s(w http.ResponseWriter, r *http.Request) {\n",
		handlerVar, handlerName, entity)
	content.WriteString("\tvars := mux.Vars(r)\n")
	content.WriteString("\tid, err := strconv.Atoi(vars[\"id\"])\n")
	content.WriteString("\tif err != nil {\n")
	fmt.Fprintf(content, "\t\thttp.Error(w, \"Invalid %s ID\", http.StatusBadRequest)\n", strings.ToLower(entity))
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	fmt.Fprintf(content, "\tvar input usecase.Update%sInput\n", entity)
	content.WriteString("\tif err := json.NewDecoder(r.Body).Decode(&input); err != nil {\n")
	content.WriteString("\t\thttp.Error(w, \"Invalid request body\", http.StatusBadRequest)\n")
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	if validation {
		content.WriteString("\tif err := validator.New().Struct(input); err != nil {\n")
		content.WriteString("\t\thttp.Error(w, err.Error(), http.StatusUnprocessableEntity)\n")
		content.WriteString("\t\treturn\n")
		content.WriteString("\t}\n\n")
	}

	fmt.Fprintf(content, "\tif err := %s.usecase.Update%s(id, input); err != nil {\n", handlerVar, entity)
	content.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString("\tw.WriteHeader(http.StatusNoContent)\n")
	content.WriteString("}\n\n")
}

func generateDeleteHandlerMethod(content *strings.Builder, entity, handlerName string) {
	handlerVar := strings.ToLower(string(handlerName[0]))

	fmt.Fprintf(content, "func (%s *%s) Delete%s(w http.ResponseWriter, r *http.Request) {\n",
		handlerVar, handlerName, entity)
	content.WriteString("\tvars := mux.Vars(r)\n")
	content.WriteString("\tid, err := strconv.Atoi(vars[\"id\"])\n")
	content.WriteString("\tif err != nil {\n")
	fmt.Fprintf(content, "\t\thttp.Error(w, \"Invalid %s ID\", http.StatusBadRequest)\n", strings.ToLower(entity))
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	fmt.Fprintf(content, "\tif err := %s.usecase.Delete%s(id); err != nil {\n", handlerVar, entity)
	content.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString("\tw.WriteHeader(http.StatusNoContent)\n")
	content.WriteString("}\n\n")
}

func generateListHandlerMethod(content *strings.Builder, entity, handlerName string) {
	handlerVar := strings.ToLower(string(handlerName[0]))

	fmt.Fprintf(content, "func (%s *%s) List%ss(w http.ResponseWriter, r *http.Request) {\n",
		handlerVar, handlerName, entity)
	fmt.Fprintf(content, "\toutput, err := %s.usecase.List%ss()\n", handlerVar, entity)
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString("\tw.Header().Set(\"Content-Type\", \"application/json\")\n")
	content.WriteString("\tjson.NewEncoder(w).Encode(output)\n")
	content.WriteString("}\n\n")
}

func generateHTTPRoutesFile(dir, entity string, middleware bool, sm ...*SafetyManager) {
	filename := filepath.Join(dir, "routes.go")

	// Get the module name from go.mod
	moduleName := getModuleName()
	importPath := getImportPath(moduleName)

	// Detect whether the standalone middleware package exists.
	middlewarePkgExists := middlewarePackageExists()

	var content strings.Builder
	content.WriteString("package http\n\n")
	content.WriteString("import (\n")
	if middleware && !middlewarePkgExists {
		content.WriteString("\t\"log\"\n")
		content.WriteString("\t\"net/http\"\n\n")
	}
	content.WriteString("\t\"github.com/gorilla/mux\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/usecase\"\n", importPath))
	if middleware && middlewarePkgExists {
		content.WriteString(fmt.Sprintf("\t\"%s/internal/middleware\"\n", importPath))
	}
	content.WriteString(")\n\n")

	entityLower := strings.ToLower(entity)
	pluralEntity := entityLower + "s"

	content.WriteString(fmt.Sprintf("func Setup%sRoutes(router *mux.Router, uc usecase.%sUseCase) {\n",
		entity, entity))
	content.WriteString(fmt.Sprintf("\thandler := New%sHandler(uc)\n\n", entity))

	if middleware {
		content.WriteString("\t// Apply middleware\n")
		content.WriteString(fmt.Sprintf("\t%sRouter := router.PathPrefix(\"/%s\").Subrouter()\n", entityLower, pluralEntity))
		if middlewarePkgExists {
			content.WriteString(fmt.Sprintf("\t%sRouter.Use(middleware.CORS(middleware.DefaultCORSConfig()))\n", entityLower))
			content.WriteString(fmt.Sprintf("\t%sRouter.Use(middleware.Logging())\n\n", entityLower))
		} else {
			content.WriteString(fmt.Sprintf("\t%sRouter.Use(corsMiddleware)\n", entityLower))
			content.WriteString(fmt.Sprintf("\t%sRouter.Use(loggingMiddleware)\n\n", entityLower))
		}

		content.WriteString(fmt.Sprintf("\t%sRouter.HandleFunc(\"\", handler.Create%s).Methods(\"POST\")\n",
			entityLower, entity))
		content.WriteString(fmt.Sprintf("\t%sRouter.HandleFunc(\"/{id}\", handler.Get%s).Methods(\"GET\")\n",
			entityLower, entity))
		content.WriteString(fmt.Sprintf("\t%sRouter.HandleFunc(\"/{id}\", handler.Update%s).Methods(\"PUT\")\n",
			entityLower, entity))
		content.WriteString(fmt.Sprintf("\t%sRouter.HandleFunc(\"/{id}\", handler.Delete%s).Methods(\"DELETE\")\n",
			entityLower, entity))
		content.WriteString(fmt.Sprintf("\t%sRouter.HandleFunc(\"\", handler.List%ss).Methods(\"GET\")\n",
			entityLower, entity))
	} else {
		content.WriteString(fmt.Sprintf("\trouter.HandleFunc(\"/%s\", handler.Create%s).Methods(\"POST\")\n",
			pluralEntity, entity))
		content.WriteString(fmt.Sprintf("\trouter.HandleFunc(\"/%s/{id}\", handler.Get%s).Methods(\"GET\")\n",
			pluralEntity, entity))
		content.WriteString(fmt.Sprintf("\trouter.HandleFunc(\"/%s/{id}\", handler.Update%s).Methods(\"PUT\")\n",
			pluralEntity, entity))
		content.WriteString(fmt.Sprintf("\trouter.HandleFunc(\"/%s/{id}\", handler.Delete%s).Methods(\"DELETE\")\n",
			pluralEntity, entity))
		content.WriteString(fmt.Sprintf("\trouter.HandleFunc(\"/%s\", handler.List%ss).Methods(\"GET\")\n",
			pluralEntity, entity))
	}

	content.WriteString("}\n")

	if middleware && !middlewarePkgExists {
		content.WriteString("\n// Middleware functions\n")
		generateMiddlewareFunctions(&content)
	}

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing routes file: %v", err))
		return
	}
}

func generateMiddlewareFunctions(content *strings.Builder) {
	content.WriteString(`
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
`)
}

func generateHTTPDTOFile(dir, entity string, sm ...*SafetyManager) {
	filename := filepath.Join(dir, "dto.go")

	var content strings.Builder
	content.WriteString("package http\n\n")

	content.WriteString(fmt.Sprintf("// HTTP-specific DTOs for %s\n", entity))
	content.WriteString(fmt.Sprintf("type HTTP%sRequest struct {\n", entity))
	content.WriteString("\tName  string `json:\"name\" validate:\"required\"`\n")
	content.WriteString("\tEmail string `json:\"email\" validate:\"required,email\"`\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("type HTTP%sResponse struct {\n", entity))
	content.WriteString("\tID      int    `json:\"id\"`\n")
	content.WriteString("\tName    string `json:\"name\"`\n")
	content.WriteString("\tEmail   string `json:\"email\"`\n")
	content.WriteString("\tMessage string `json:\"message,omitempty\"`\n")
	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing types file: %v", err))
		return
	}
}

func generateSwaggerFile(dir, entity string, sm ...*SafetyManager) {
	filename := filepath.Join(dir, "swagger.yaml")
	entityLower := strings.ToLower(entity)

	content := fmt.Sprintf(`openapi: 3.0.0
info:
  title: %s API
  version: 1.0.0
  description: API for managing %s entities

paths:
  /%ss:
    get:
      summary: List all %ss
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/%s'
    post:
      summary: Create a new %s
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Create%sRequest'
      responses:
        '201':
          description: %s created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/%s'

  /%ss/{id}:
    get:
      summary: Get %s by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/%s'

components:
  schemas:
    %s:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
        email:
          type: string
    
    Create%sRequest:
      type: object
      required:
        - name
        - email
      properties:
        name:
          type: string
        email:
          type: string
`, entity, entityLower, entityLower, entityLower, entity, entityLower, entity, entity, entity, entityLower, entityLower, entity, entity, entity)

	if err := writeFile(filename, content, sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing swagger file: %v", err))
		return
	}
}

func init() {
	handlerCmd.Flags().StringP("type", "t", "http", "Handler type (http, grpc, cli, worker, soap)")
	handlerCmd.Flags().BoolP("middleware", "m", false, "Include middleware setup")
	handlerCmd.Flags().Bool("validation", false, "Input validation in handler")
	handlerCmd.Flags().BoolP("swagger", "s", false, "Generate Swagger documentation (HTTP only)")
	handlerCmd.Flags().Bool("dry-run", false, "Preview changes without creating files")
	handlerCmd.Flags().Bool("force", false, "Overwrite existing files without asking")
	handlerCmd.Flags().Bool("backup", false, "Backup existing files before overwriting")
}
