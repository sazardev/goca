package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var interfacesCmd = &cobra.Command{
	Use:   "interfaces <entity>",
	Short: "Generar solo interfaces para TDD",
	Long: `Genera únicamente las interfaces de contratos entre capas, 
útil para desarrollo dirigido por pruebas (TDD).`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entity := args[0]

		usecase, _ := cmd.Flags().GetBool("usecase")
		repository, _ := cmd.Flags().GetBool("repository")
		handler, _ := cmd.Flags().GetBool("handler")
		all, _ := cmd.Flags().GetBool("all")

		// If all is true, enable all options
		if all {
			usecase = true
			repository = true
			handler = true
		}

		// If no specific flags, generate all by default
		if !usecase && !repository && !handler {
			usecase = true
			repository = true
			handler = true
		}

		fmt.Printf("Generando interfaces para entidad '%s'\n", entity)

		if usecase {
			fmt.Println("✓ Generando interfaces de casos de uso")
		}
		if repository {
			fmt.Println("✓ Generando interfaces de repositorio")
		}
		if handler {
			fmt.Println("✓ Generando interfaces de handlers")
		}

		generateInterfaces(entity, usecase, repository, handler)
		fmt.Printf("\n✅ Interfaces para '%s' generadas exitosamente!\n", entity)
	},
}

func generateInterfaces(entity string, usecase, repository, handler bool) {
	// Create interfaces directory
	interfacesDir := filepath.Join("internal", "interfaces")
	_ = os.MkdirAll(interfacesDir, 0755)

	if usecase {
		generateUseCaseInterfaceFile(interfacesDir, entity)
	}

	if repository {
		generateRepositoryInterfaceFile(interfacesDir, entity)
	}

	if handler {
		generateHandlerInterfaceFile(interfacesDir, entity)
	}
}

func generateUseCaseInterfaceFile(dir, entity string) {
	// Get the module name from go.mod
	moduleName := getModuleName()

	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, entityLower+"_usecase.go")

	var content strings.Builder
	content.WriteString("package interfaces\n\n")
	content.WriteString(fmt.Sprintf("import \"%s/internal/domain\"\n\n", moduleName))

	// DTO interfaces
	content.WriteString(fmt.Sprintf("// %s UseCase DTOs\n", entity))
	content.WriteString(fmt.Sprintf("type Create%sInput interface {\n", entity))
	content.WriteString("\tGetName() string\n")
	content.WriteString("\tGetEmail() string\n")
	content.WriteString("\tValidate() error\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("type Create%sOutput interface {\n", entity))
	content.WriteString(fmt.Sprintf("\tGet%s() domain.%s\n", entity, entity))
	content.WriteString("\tGetMessage() string\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("type Update%sInput interface {\n", entity))
	content.WriteString("\tGetName() string\n")
	content.WriteString("\tGetEmail() string\n")
	content.WriteString("\tValidate() error\n")
	content.WriteString("}\n\n")

	// UseCase interface
	content.WriteString(fmt.Sprintf("// %s UseCase interface\n", entity))
	content.WriteString(fmt.Sprintf("type %sUseCase interface {\n", entity))
	content.WriteString(fmt.Sprintf("\tCreate%s(input Create%sInput) (Create%sOutput, error)\n",
		entity, entity, entity))
	content.WriteString(fmt.Sprintf("\tGet%s(id int) (*domain.%s, error)\n", entity, entity))
	content.WriteString(fmt.Sprintf("\tUpdate%s(id int, input Update%sInput) error\n", entity, entity))
	content.WriteString(fmt.Sprintf("\tDelete%s(id int) error\n", entity))
	content.WriteString(fmt.Sprintf("\tList%ss() ([]domain.%s, error)\n", entity, entity))
	content.WriteString("}\n")

	writeGoFile(filename, content.String())
}

func generateRepositoryInterfaceFile(dir, entity string) {
	generateRepositoryInterfaceFileWithFields(dir, entity, "")
}

func generateRepositoryInterfaceFileWithFields(dir, entity, fields string) {
	// Get the module name from go.mod
	moduleName := getModuleName()

	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, entityLower+"_repository.go")

	var content strings.Builder
	content.WriteString("package interfaces\n\n")
	content.WriteString(fmt.Sprintf("import \"%s/internal/domain\"\n\n", moduleName))

	content.WriteString(fmt.Sprintf("// %s Repository interface\n", entity))
	content.WriteString(fmt.Sprintf("type %sRepository interface {\n", entity))

	// Basic CRUD operations
	content.WriteString(fmt.Sprintf("\tSave(%s *domain.%s) error\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tFindByID(id int) (*domain.%s, error)\n", entity))

	// Generate dynamic query methods based on fields
	if fields != "" {
		validator := NewFieldValidator()
		fieldsList, err := validator.ParseFieldsWithValidation(fields)
		if err == nil {
			queryMethods := validator.GenerateQueryMethodsForFields(entity, fieldsList)
			for _, method := range queryMethods {
				if method.MethodName != "FindByID" { // Skip ID as it's already added
					content.WriteString(fmt.Sprintf("\t%s(%s %s) (*domain.%s, error)\n",
						method.MethodName, method.Field, method.Type, entity))
				}
			}
		}
	} else {
		// Fallback to basic query method if no fields specified
		content.WriteString(fmt.Sprintf("\tFindByEmail(email string) (*domain.%s, error)\n", entity))
	}

	content.WriteString(fmt.Sprintf("\tUpdate(%s *domain.%s) error\n", entityLower, entity))
	content.WriteString("\tDelete(id int) error\n")
	content.WriteString(fmt.Sprintf("\tFindAll() ([]domain.%s, error)\n", entity))

	// Query operations
	content.WriteString("\tCount() (int, error)\n")
	content.WriteString("\tExists(id int) (bool, error)\n")

	// Batch operations
	content.WriteString(fmt.Sprintf("\tSaveBatch(%ss []domain.%s) error\n", entityLower, entity))
	content.WriteString("\tDeleteBatch(ids []int) error\n")

	// Transaction operations
	content.WriteString(fmt.Sprintf("\tSaveWithTx(tx interface{}, %s *domain.%s) error\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tUpdateWithTx(tx interface{}, %s *domain.%s) error\n", entityLower, entity))
	content.WriteString("\tDeleteWithTx(tx interface{}, id int) error\n")

	content.WriteString("}\n")

	writeGoFile(filename, content.String())
}

func generateHandlerInterfaceFile(dir, entity string) {
	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, entityLower+"_handler.go")

	var content strings.Builder
	content.WriteString("package interfaces\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"net/http\"\n")
	content.WriteString("\t\"context\"\n")
	content.WriteString(")\n\n")

	// HTTP Handler interface
	content.WriteString(fmt.Sprintf("// %s HTTP Handler interface\n", entity))
	content.WriteString(fmt.Sprintf("type %sHTTPHandler interface {\n", entity))
	content.WriteString(fmt.Sprintf("\tCreate%s(w http.ResponseWriter, r *http.Request)\n", entity))
	content.WriteString(fmt.Sprintf("\tGet%s(w http.ResponseWriter, r *http.Request)\n", entity))
	content.WriteString(fmt.Sprintf("\tUpdate%s(w http.ResponseWriter, r *http.Request)\n", entity))
	content.WriteString(fmt.Sprintf("\tDelete%s(w http.ResponseWriter, r *http.Request)\n", entity))
	content.WriteString(fmt.Sprintf("\tList%ss(w http.ResponseWriter, r *http.Request)\n", entity))
	content.WriteString("}\n\n")

	// gRPC Handler interface
	content.WriteString(fmt.Sprintf("// %s gRPC Handler interface\n", entity))
	content.WriteString(fmt.Sprintf("type %sGRPCHandler interface {\n", entity))
	content.WriteString(fmt.Sprintf("\tCreate%s(ctx context.Context, req *Create%sRequest) (*Create%sResponse, error)\n",
		entity, entity, entity))
	content.WriteString(fmt.Sprintf("\tGet%s(ctx context.Context, req *Get%sRequest) (*%sResponse, error)\n",
		entity, entity, entity))
	content.WriteString(fmt.Sprintf("\tUpdate%s(ctx context.Context, req *Update%sRequest) (*Update%sResponse, error)\n",
		entity, entity, entity))
	content.WriteString(fmt.Sprintf("\tDelete%s(ctx context.Context, req *Delete%sRequest) (*Delete%sResponse, error)\n",
		entity, entity, entity))
	content.WriteString(fmt.Sprintf("\tList%ss(ctx context.Context, req *List%ssRequest) (*List%ssResponse, error)\n",
		entity, entity, entity))
	content.WriteString("}\n\n")

	// CLI Handler interface
	content.WriteString(fmt.Sprintf("// %s CLI Handler interface\n", entity))
	content.WriteString(fmt.Sprintf("type %sCLIHandler interface {\n", entity))
	content.WriteString(fmt.Sprintf("\tCreate%sCommand() interface{}\n", entity))
	content.WriteString(fmt.Sprintf("\tGet%sCommand() interface{}\n", entity))
	content.WriteString(fmt.Sprintf("\tUpdate%sCommand() interface{}\n", entity))
	content.WriteString(fmt.Sprintf("\tDelete%sCommand() interface{}\n", entity))
	content.WriteString(fmt.Sprintf("\tList%ssCommand() interface{}\n", entity))
	content.WriteString("}\n\n")

	// Request/Response interfaces for gRPC
	content.WriteString("// gRPC Request/Response interfaces\n")
	generateGRPCRequestResponseInterfaces(&content, entity)

	writeGoFile(filename, content.String())
}

func generateGRPCRequestResponseInterfaces(content *strings.Builder, entity string) {
	// Create Request interface
	content.WriteString(fmt.Sprintf("type Create%sRequest interface {\n", entity))
	content.WriteString("\tGetName() string\n")
	content.WriteString("\tGetEmail() string\n")
	content.WriteString("}\n\n")

	// Create Response interface
	content.WriteString(fmt.Sprintf("type Create%sResponse interface {\n", entity))
	content.WriteString(fmt.Sprintf("\tGet%s() *%s\n", entity, entity))
	content.WriteString("\tGetMessage() string\n")
	content.WriteString("}\n\n")

	// Get Request interface
	content.WriteString(fmt.Sprintf("type Get%sRequest interface {\n", entity))
	content.WriteString("\tGetId() int32\n")
	content.WriteString("}\n\n")

	// Get Response interface
	content.WriteString(fmt.Sprintf("type %sResponse interface {\n", entity))
	content.WriteString(fmt.Sprintf("\tGet%s() *%s\n", entity, entity))
	content.WriteString("}\n\n")

	// Update Request interface
	content.WriteString(fmt.Sprintf("type Update%sRequest interface {\n", entity))
	content.WriteString("\tGetId() int32\n")
	content.WriteString("\tGetName() string\n")
	content.WriteString("\tGetEmail() string\n")
	content.WriteString("}\n\n")

	// Update Response interface
	content.WriteString(fmt.Sprintf("type Update%sResponse interface {\n", entity))
	content.WriteString("\tGetMessage() string\n")
	content.WriteString("}\n\n")

	// Delete Request interface
	content.WriteString(fmt.Sprintf("type Delete%sRequest interface {\n", entity))
	content.WriteString("\tGetId() int32\n")
	content.WriteString("}\n\n")

	// Delete Response interface
	content.WriteString(fmt.Sprintf("type Delete%sResponse interface {\n", entity))
	content.WriteString("\tGetMessage() string\n")
	content.WriteString("}\n\n")

	// List Request interface
	content.WriteString(fmt.Sprintf("type List%ssRequest interface {\n", entity))
	content.WriteString("\t// No fields for basic list\n")
	content.WriteString("}\n\n")

	// List Response interface
	content.WriteString(fmt.Sprintf("type List%ssResponse interface {\n", entity))
	content.WriteString(fmt.Sprintf("\tGet%ss() []*%s\n", entity, entity))
	content.WriteString("\tGetTotal() int32\n")
	content.WriteString("}\n")
}

func init() {
	interfacesCmd.Flags().BoolP("usecase", "u", false, "Generar interfaces de casos de uso")
	interfacesCmd.Flags().BoolP("repository", "r", false, "Generar interfaces de repositorio")
	interfacesCmd.Flags().BoolP("handler", "", false, "Generar interfaces de handlers")
	interfacesCmd.Flags().BoolP("all", "a", false, "Generar todas las interfaces")
}
