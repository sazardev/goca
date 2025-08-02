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
	Short: "Generar handlers para diferentes protocolos",
	Long: `Crea adaptadores de entrega que manejan diferentes protocolos 
(HTTP, gRPC, CLI) manteniendo la separación de capas.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entity := args[0]

		handlerType, _ := cmd.Flags().GetString("type")
		middleware, _ := cmd.Flags().GetBool("middleware")
		validation, _ := cmd.Flags().GetBool("validation")
		swagger, _ := cmd.Flags().GetBool("swagger")

		fmt.Printf("Generando handler '%s' para entidad '%s'\n", handlerType, entity)

		if middleware {
			fmt.Println("✓ Incluyendo middleware")
		}
		if validation {
			fmt.Println("✓ Incluyendo validación")
		}
		if swagger && handlerType == HandlerHTTP {
			fmt.Println("✓ Incluyendo documentación Swagger")
		}

		generateHandler(entity, handlerType, middleware, validation, swagger)
		fmt.Printf("\n✅ Handler '%s' para '%s' generado exitosamente!\n", handlerType, entity)
	},
}

func generateHandler(entity, handlerType string, middleware, validation, swagger bool) {
	switch handlerType {
	case HandlerHTTP:
		generateHTTPHandler(entity, middleware, validation, swagger)
	case HandlerGRPC:
		generateGRPCHandler(entity)
	case HandlerCLI:
		generateCLIHandler(entity)
	case "worker":
		generateWorkerHandler(entity)
	case "soap":
		generateSOAPHandler(entity)
	default:
		fmt.Printf("Tipo de handler no soportado: %s\n", handlerType)
		os.Exit(1)
	}
}

func generateHTTPHandler(entity string, middleware, validation, swagger bool) {
	// Crear directorio handlers si no existe
	handlerDir := filepath.Join(DirInternal, DirHandler, DirHTTP)
	_ = os.MkdirAll(handlerDir, 0755) // Generate handler file
	generateHTTPHandlerFile(handlerDir, entity, validation)

	// Generate routes file
	generateHTTPRoutesFile(handlerDir, entity, middleware)

	// Generate DTOs for HTTP if validation is enabled
	if validation {
		generateHTTPDTOFile(handlerDir, entity)
	}

	// Generate Swagger docs if requested
	if swagger {
		generateSwaggerFile(handlerDir, entity)
	}
}

func generateHTTPHandlerFile(dir, entity string, validation bool) {
	filename := filepath.Join(dir, strings.ToLower(entity)+"_handler.go")

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
		content.WriteString(fmt.Sprintf("\t\"%s/internal/messages\"\n", importPath))
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
	generateCreateHandlerMethod(&content, entity, handlerName)
	generateGetHandlerMethod(&content, entity, handlerName)
	generateUpdateHandlerMethod(&content, entity, handlerName)
	generateDeleteHandlerMethod(&content, entity, handlerName)
	generateListHandlerMethod(&content, entity, handlerName)

	writeGoFile(filename, content.String())
}

func generateCreateHandlerMethod(content *strings.Builder, entity, handlerName string) {
	handlerVar := strings.ToLower(string(handlerName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) Create%s(w http.ResponseWriter, r *http.Request) {\n",
		handlerVar, handlerName, entity))
	content.WriteString(fmt.Sprintf("\tvar input usecase.Create%sInput\n\n", entity))

	content.WriteString("\tif err := json.NewDecoder(r.Body).Decode(&input); err != nil {\n")
	content.WriteString("\t\thttp.Error(w, \"Invalid request body\", http.StatusBadRequest)\n")
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\toutput, err := %s.usecase.Create%s(input)\n", handlerVar, entity))
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

	content.WriteString(fmt.Sprintf("func (%s *%s) Get%s(w http.ResponseWriter, r *http.Request) {\n",
		handlerVar, handlerName, entity))
	content.WriteString("\tvars := mux.Vars(r)\n")
	content.WriteString("\tid, err := strconv.Atoi(vars[\"id\"])\n")
	content.WriteString("\tif err != nil {\n")
	content.WriteString(fmt.Sprintf("\t\thttp.Error(w, \"Invalid %s ID\", http.StatusBadRequest)\n", strings.ToLower(entity)))
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\t%s, err := %s.usecase.Get%s(id)\n", strings.ToLower(entity), handlerVar, entity))
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\thttp.Error(w, err.Error(), http.StatusNotFound)\n")
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString("\tw.Header().Set(\"Content-Type\", \"application/json\")\n")
	content.WriteString(fmt.Sprintf("\tjson.NewEncoder(w).Encode(%s)\n", strings.ToLower(entity)))
	content.WriteString("}\n\n")
}

func generateUpdateHandlerMethod(content *strings.Builder, entity, handlerName string) {
	handlerVar := strings.ToLower(string(handlerName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) Update%s(w http.ResponseWriter, r *http.Request) {\n",
		handlerVar, handlerName, entity))
	content.WriteString("\tvars := mux.Vars(r)\n")
	content.WriteString("\tid, err := strconv.Atoi(vars[\"id\"])\n")
	content.WriteString("\tif err != nil {\n")
	content.WriteString(fmt.Sprintf("\t\thttp.Error(w, \"Invalid %s ID\", http.StatusBadRequest)\n", strings.ToLower(entity)))
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\tvar input usecase.Update%sInput\n", entity))
	content.WriteString("\tif err := json.NewDecoder(r.Body).Decode(&input); err != nil {\n")
	content.WriteString("\t\thttp.Error(w, \"Invalid request body\", http.StatusBadRequest)\n")
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\tif err := %s.usecase.Update%s(id, input); err != nil {\n", handlerVar, entity))
	content.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString("\tw.WriteHeader(http.StatusNoContent)\n")
	content.WriteString("}\n\n")
}

func generateDeleteHandlerMethod(content *strings.Builder, entity, handlerName string) {
	handlerVar := strings.ToLower(string(handlerName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) Delete%s(w http.ResponseWriter, r *http.Request) {\n",
		handlerVar, handlerName, entity))
	content.WriteString("\tvars := mux.Vars(r)\n")
	content.WriteString("\tid, err := strconv.Atoi(vars[\"id\"])\n")
	content.WriteString("\tif err != nil {\n")
	content.WriteString(fmt.Sprintf("\t\thttp.Error(w, \"Invalid %s ID\", http.StatusBadRequest)\n", strings.ToLower(entity)))
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\tif err := %s.usecase.Delete%s(id); err != nil {\n", handlerVar, entity))
	content.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString("\tw.WriteHeader(http.StatusNoContent)\n")
	content.WriteString("}\n\n")
}

func generateListHandlerMethod(content *strings.Builder, entity, handlerName string) {
	handlerVar := strings.ToLower(string(handlerName[0]))

	content.WriteString(fmt.Sprintf("func (%s *%s) List%ss(w http.ResponseWriter, r *http.Request) {\n",
		handlerVar, handlerName, entity))
	content.WriteString(fmt.Sprintf("\toutput, err := %s.usecase.List%ss()\n", handlerVar, entity))
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	content.WriteString("\t\treturn\n")
	content.WriteString("\t}\n\n")

	content.WriteString("\tw.Header().Set(\"Content-Type\", \"application/json\")\n")
	content.WriteString("\tjson.NewEncoder(w).Encode(output)\n")
	content.WriteString("}\n\n")
}

func generateHTTPRoutesFile(dir, entity string, middleware bool) {
	filename := filepath.Join(dir, "routes.go")

	// Get the module name from go.mod
	moduleName := getModuleName()
	importPath := getImportPath(moduleName)

	var content strings.Builder
	content.WriteString("package http\n\n")
	content.WriteString("import (\n")
	if middleware {
		content.WriteString("\t\"log\"\n")
		content.WriteString("\t\"net/http\"\n\n")
	}
	content.WriteString("\t\"github.com/gorilla/mux\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/usecase\"\n", importPath))
	content.WriteString(")\n\n")

	entityLower := strings.ToLower(entity)
	pluralEntity := entityLower + "s"

	content.WriteString(fmt.Sprintf("func Setup%sRoutes(router *mux.Router, uc usecase.%sUseCase) {\n",
		entity, entity))
	content.WriteString(fmt.Sprintf("\thandler := New%sHandler(uc)\n\n", entity))

	if middleware {
		content.WriteString("\t// Apply middleware\n")
		content.WriteString(fmt.Sprintf("\t%sRouter := router.PathPrefix(\"/%s\").Subrouter()\n", entityLower, pluralEntity))
		content.WriteString(fmt.Sprintf("\t%sRouter.Use(corsMiddleware)\n", entityLower))
		content.WriteString(fmt.Sprintf("\t%sRouter.Use(loggingMiddleware)\n\n", entityLower))

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

	if middleware {
		content.WriteString("\n// Middleware functions\n")
		generateMiddlewareFunctions(&content)
	}

	writeGoFile(filename, content.String())
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

func generateHTTPDTOFile(dir, entity string) {
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

	writeGoFile(filename, content.String())
}

func generateSwaggerFile(dir, entity string) {
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
          description: %s creado exitosamente
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

	writeFile(filename, content)
}

func generateGRPCHandler(entity string) {
	// Create gRPC directory
	grpcDir := filepath.Join(DirInternal, DirHandler, DirGRPC)
	_ = os.MkdirAll(grpcDir, 0755)

	generateProtoFile(grpcDir, entity)
	generateGRPCServerFile(grpcDir, entity)
}

func generateProtoFile(dir, entity string) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, entityLower+".proto")

	var content strings.Builder
	content.WriteString("syntax = \"proto3\";\n\n")
	content.WriteString(fmt.Sprintf("package %s;\n\n", entityLower))
	content.WriteString(fmt.Sprintf("option go_package = \"./%s\";\n\n", entityLower))

	content.WriteString(fmt.Sprintf("service %sService {\n", entity))
	content.WriteString(fmt.Sprintf("  rpc Create%s(Create%sRequest) returns (Create%sResponse);\n", entity, entity, entity))
	content.WriteString(fmt.Sprintf("  rpc Get%s(Get%sRequest) returns (%sResponse);\n", entity, entity, entity))
	content.WriteString(fmt.Sprintf("  rpc Update%s(Update%sRequest) returns (Update%sResponse);\n", entity, entity, entity))
	content.WriteString(fmt.Sprintf("  rpc Delete%s(Delete%sRequest) returns (Delete%sResponse);\n", entity, entity, entity))
	content.WriteString(fmt.Sprintf("  rpc List%ss(List%ssRequest) returns (List%ssResponse);\n", entity, entity, entity))
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message %s {\n", entity))
	content.WriteString("  int32 id = 1;\n")
	content.WriteString("  string name = 2;\n")
	content.WriteString("  string email = 3;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Create%sRequest {\n", entity))
	content.WriteString("  string name = 1;\n")
	content.WriteString("  string email = 2;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Create%sResponse {\n", entity))
	content.WriteString(fmt.Sprintf("  %s %s = 1;\n", entity, entityLower))
	content.WriteString("  string message = 2;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Get%sRequest {\n", entity))
	content.WriteString("  int32 id = 1;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message %sResponse {\n", entity))
	content.WriteString(fmt.Sprintf("  %s %s = 1;\n", entity, entityLower))
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Update%sRequest {\n", entity))
	content.WriteString("  int32 id = 1;\n")
	content.WriteString("  string name = 2;\n")
	content.WriteString("  string email = 3;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Update%sResponse {\n", entity))
	content.WriteString("  string message = 1;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Delete%sRequest {\n", entity))
	content.WriteString("  int32 id = 1;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Delete%sResponse {\n", entity))
	content.WriteString("  string message = 1;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message List%ssRequest {\n", entity))
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message List%ssResponse {\n", entity))
	content.WriteString(fmt.Sprintf("  repeated %s %ss = 1;\n", entity, entityLower))
	content.WriteString("  int32 total = 2;\n")
	content.WriteString("}\n")

	writeGoFile(filename, content.String())
}

func generateGRPCServerFile(dir, entity string) {
	// Get the module name from go.mod
	moduleName := getModuleName()
	importPath := getImportPath(moduleName)

	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, entityLower+"_server.go")

	var content strings.Builder
	content.WriteString("package grpc\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"context\"\n\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/usecase\"\n", importPath))
	content.WriteString(fmt.Sprintf("\tpb \"%s/internal/handler/grpc/%s\"\n", importPath, entityLower))
	content.WriteString(")\n\n")

	content.WriteString(fmt.Sprintf("type %sServer struct {\n", entity))
	content.WriteString(fmt.Sprintf("\tpb.Unimplemented%sServiceServer\n", entity))
	content.WriteString(fmt.Sprintf("\tusecase usecase.%sUseCase\n", entity))
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("func New%sServer(uc usecase.%sUseCase) *%sServer {\n", entity, entity, entity))
	content.WriteString(fmt.Sprintf("\treturn &%sServer{usecase: uc}\n", entity))
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("func (s *%sServer) Create%s(ctx context.Context, req *pb.Create%sRequest) (*pb.Create%sResponse, error) {\n", entity, entity, entity, entity))
	content.WriteString(fmt.Sprintf("\tinput := usecase.Create%sInput{\n", entity))
	content.WriteString("\t\tName:  req.Name,\n")
	content.WriteString("\t\tEmail: req.Email,\n")
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\toutput, err := s.usecase.Create%s(input)\n", entity))
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\treturn &pb.Create%sResponse{\n", entity))
	content.WriteString(fmt.Sprintf("\t\t%s: &pb.%s{\n", entity, entity))
	content.WriteString(fmt.Sprintf("\t\t\tId:    int32(output.%s.ID),\n", entity))
	content.WriteString(fmt.Sprintf("\t\t\tName:  output.%s.Name,\n", entity))
	content.WriteString(fmt.Sprintf("\t\t\tEmail: output.%s.Email,\n", entity))
	content.WriteString("\t\t},\n")
	content.WriteString("\t\tMessage: output.Message,\n")
	content.WriteString("\t}, nil\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("func (s *%sServer) Get%s(ctx context.Context, req *pb.Get%sRequest) (*pb.%sResponse, error) {\n", entity, entity, entity, entity))
	content.WriteString(fmt.Sprintf("\t%s, err := s.usecase.Get%s(int(req.Id))\n", entityLower, entity))
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\treturn &pb.%sResponse{\n", entity))
	content.WriteString(fmt.Sprintf("\t\t%s: &pb.%s{\n", entity, entity))
	content.WriteString(fmt.Sprintf("\t\t\tId:    int32(%s.ID),\n", entityLower))
	content.WriteString(fmt.Sprintf("\t\t\tName:  %s.Name,\n", entityLower))
	content.WriteString(fmt.Sprintf("\t\t\tEmail: %s.Email,\n", entityLower))
	content.WriteString("\t\t},\n")
	content.WriteString("\t}, nil\n")
	content.WriteString("}\n")

	writeGoFile(filename, content.String())
}

func generateCLIHandler(entity string) {
	// Create CLI directory
	cliDir := filepath.Join(DirInternal, DirHandler, DirCLI)
	_ = os.MkdirAll(cliDir, 0755)

	// Get the module name from go.mod
	moduleName := getModuleName()
	importPath := getImportPath(moduleName)

	entityLower := strings.ToLower(entity)
	filename := filepath.Join(cliDir, entityLower+"_commands.go")

	var content strings.Builder
	content.WriteString("package cli\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"fmt\"\n")
	content.WriteString("\t\"strconv\"\n\n")
	content.WriteString("\t\"github.com/spf13/cobra\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/usecase\"\n", importPath))
	content.WriteString(")\n\n")

	// CLI struct
	content.WriteString(fmt.Sprintf("type %sCLI struct {\n", entity))
	content.WriteString(fmt.Sprintf("\tusecase usecase.%sUseCase\n", entity))
	content.WriteString("}\n\n")

	// Constructor
	content.WriteString(fmt.Sprintf("func New%sCLI(uc usecase.%sUseCase) *%sCLI {\n", entity, entity, entity))
	content.WriteString(fmt.Sprintf("\treturn &%sCLI{usecase: uc}\n", entity))
	content.WriteString("}\n\n")

	// Create command
	content.WriteString(fmt.Sprintf("func (c *%sCLI) Create%sCommand() *cobra.Command {\n", entity, entity))
	content.WriteString("\tcmd := &cobra.Command{\n")
	content.WriteString("\t\tUse:   \"crear\",\n")
	content.WriteString(fmt.Sprintf("\t\tShort: \"Crear un nuevo %s\",\n", entityLower))
	content.WriteString("\t\tRun: func(cmd *cobra.Command, args []string) {\n")
	content.WriteString("\t\t\tname, _ := cmd.Flags().GetString(\"name\")\n")
	content.WriteString("\t\t\temail, _ := cmd.Flags().GetString(\"email\")\n\n")

	content.WriteString(fmt.Sprintf("\t\t\tinput := usecase.Create%sInput{\n", entity))
	content.WriteString("\t\t\t\tName:  name,\n")
	content.WriteString("\t\t\t\tEmail: email,\n")
	content.WriteString("\t\t\t}\n\n")

	content.WriteString(fmt.Sprintf("\t\t\toutput, err := c.usecase.Create%s(input)\n", entity))
	content.WriteString("\t\t\tif err != nil {\n")
	content.WriteString("\t\t\t\tfmt.Printf(\"Error: %v\\n\", err)\n")
	content.WriteString("\t\t\t\treturn\n")
	content.WriteString("\t\t\t}\n\n")

	content.WriteString(fmt.Sprintf("\t\t\tfmt.Printf(\"%s created: %%+v\\n\", output.%s)\n", entity, entity))
	content.WriteString("\t\t},\n")
	content.WriteString("\t}\n\n")

	content.WriteString("\tcmd.Flags().StringP(\"name\", \"n\", \"\", \"Name of the ")
	content.WriteString(entityLower)
	content.WriteString("\")\n")
	content.WriteString("\tcmd.Flags().StringP(\"email\", \"e\", \"\", \"Email of the ")
	content.WriteString(entityLower)
	content.WriteString("\")\n")
	content.WriteString("\tcmd.MarkFlagRequired(\"name\")\n")
	content.WriteString("\tcmd.MarkFlagRequired(\"email\")\n\n")
	content.WriteString("\treturn cmd\n")
	content.WriteString("}\n\n")

	// Get command
	content.WriteString(fmt.Sprintf("func (c *%sCLI) Get%sCommand() *cobra.Command {\n", entity, entity))
	content.WriteString("\treturn &cobra.Command{\n")
	content.WriteString("\t\tUse:   \"obtener [id]\",\n")
	content.WriteString(fmt.Sprintf("\t\tShort: \"Obtener %s por ID\",\n", entityLower))
	content.WriteString("\t\tArgs:  cobra.ExactArgs(1),\n")
	content.WriteString("\t\tRun: func(cmd *cobra.Command, args []string) {\n")
	content.WriteString("\t\t\tid, err := strconv.Atoi(args[0])\n")
	content.WriteString("\t\t\tif err != nil {\n")
	content.WriteString("\t\t\t\tfmt.Printf(\"Invalid ID: %v\\n\", err)\n")
	content.WriteString("\t\t\t\treturn\n")
	content.WriteString("\t\t\t}\n\n")

	content.WriteString(fmt.Sprintf("\t\t\t%s, err := c.usecase.Get%s(id)\n", entityLower, entity))
	content.WriteString("\t\t\tif err != nil {\n")
	content.WriteString("\t\t\t\tfmt.Printf(\"Error: %v\\n\", err)\n")
	content.WriteString("\t\t\t\treturn\n")
	content.WriteString("\t\t\t}\n\n")

	content.WriteString(fmt.Sprintf("\t\t\tfmt.Printf(\"%s: %%+v\\n\", %s)\n", entity, entityLower))
	content.WriteString("\t\t},\n")
	content.WriteString("\t}\n")
	content.WriteString("}\n")

	writeGoFile(filename, content.String())
}

func generateWorkerHandler(entity string) {
	// Create worker directory
	workerDir := filepath.Join(DirInternal, DirHandler, DirWorker)
	_ = os.MkdirAll(workerDir, 0755)

	// Get the module name from go.mod
	moduleName := getModuleName()

	entityLower := strings.ToLower(entity)
	filename := filepath.Join(workerDir, entityLower+"_worker.go")

	content := fmt.Sprintf(`package worker

import (
	"encoding/json"
	"log"
	
	"%s/internal/usecase"
)

type %sWorker struct {
	usecase usecase.%sUseCase
}

func New%sWorker(uc usecase.%sUseCase) *%sWorker {
	return &%sWorker{usecase: uc}
}

func (w *%sWorker) Process%sJob(jobData []byte) error {
	var input usecase.Create%sInput
	
	if err := json.Unmarshal(jobData, &input); err != nil {
		log.Printf("Failed to unmarshal job data: %%v", err)
		return err
	}
	
	output, err := w.usecase.Create%s(input)
	if err != nil {
		log.Printf("Failed to process %s job: %%v", err)
		return err
	}
	
	log.Printf("%s job completed: %%v", output.Message)
	return nil
}

func (w *%sWorker) ProcessBatch%sJob(jobData []byte) error {
	var inputs []usecase.Create%sInput
	
	if err := json.Unmarshal(jobData, &inputs); err != nil {
		log.Printf("Failed to unmarshal batch job data: %%v", err)
		return err
	}
	
	for _, input := range inputs {
		if _, err := w.usecase.Create%s(input); err != nil {
			log.Printf("Failed to process %s in batch: %%v", err)
			// Continue processing other items
		}
	}
	
	log.Printf("Batch %s job completed")
	return nil
}
`, moduleName, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entityLower, entity, entity, entity, entity, entity, entityLower, entity)

	writeFile(filename, content)
}

func generateSOAPHandler(entity string) {
	// Create SOAP directory
	soapDir := filepath.Join(DirInternal, DirHandler, DirSOAP)
	_ = os.MkdirAll(soapDir, 0755)

	// Get the module name from go.mod
	moduleName := getModuleName()

	entityLower := strings.ToLower(entity)
	filename := filepath.Join(soapDir, entityLower+"_client.go")

	content := fmt.Sprintf(`package soap

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	
	"%s/internal/usecase"
)

type %sSOAPClient struct {
	usecase usecase.%sUseCase
	baseURL string
}

func New%sSOAPClient(uc usecase.%sUseCase, baseURL string) *%sSOAPClient {
	return &%sSOAPClient{
		usecase: uc,
		baseURL: baseURL,
	}
}

type SOAPEnvelope struct {
	XMLName xml.Name    `+"`xml:\"soap:Envelope\"`"+`
	Header  interface{} `+"`xml:\"soap:Header\"`"+`
	Body    SOAPBody    `+"`xml:\"soap:Body\"`"+`
}

type SOAPBody struct {
	XMLName xml.Name    `+"`xml:\"soap:Body\"`"+`
	Content interface{} `+"`xml:\",innerxml\"`"+`
}

type Create%sRequest struct {
	XMLName xml.Name `+"`xml:\"Create%sRequest\"`"+`
	Name    string   `+"`xml:\"name\"`"+`
	Email   string   `+"`xml:\"email\"`"+`
}

type Create%sResponse struct {
	XMLName xml.Name `+"`xml:\"Create%sResponse\"`"+`
	ID      int      `+"`xml:\"id\"`"+`
	Message string   `+"`xml:\"message\"`"+`
}

func (c *%sSOAPClient) Create%s(name, email string) (*Create%sResponse, error) {
	request := Create%sRequest{
		Name:  name,
		Email: email,
	}
	
	envelope := SOAPEnvelope{
		Body: SOAPBody{Content: request},
	}
	
	xmlData, err := xml.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SOAP request: %%v", err)
	}
	
	resp, err := http.Post(
		c.baseURL+"/Create%s",
		"text/xml; charset=utf-8",
		strings.NewReader(string(xmlData)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send SOAP request: %%v", err)
	}
	defer resp.Body.Close()
	
	var responseEnvelope SOAPEnvelope
	if err := xml.NewDecoder(resp.Body).Decode(&responseEnvelope); err != nil {
		return nil, fmt.Errorf("failed to decode SOAP response: %%v", err)
	}
	
	var createResponse Create%sResponse
	if err := xml.Unmarshal([]byte(responseEnvelope.Body.Content.(string)), &createResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %%v", err)
	}
	
	return &createResponse, nil
}
`, moduleName, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity)

	writeFile(filename, content)
}

func init() {
	handlerCmd.Flags().StringP("type", "t", "http", "Tipo de handler (http, grpc, cli, worker, soap)")
	handlerCmd.Flags().BoolP("middleware", "m", false, "Incluir setup de middleware")
	handlerCmd.Flags().BoolP("validation", "v", false, "Validación de entrada en handler")
	handlerCmd.Flags().BoolP("swagger", "s", false, "Generar documentación Swagger (solo HTTP)")
}
