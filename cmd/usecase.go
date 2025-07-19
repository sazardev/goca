package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var usecaseCmd = &cobra.Command{
	Use:   "usecase <name>",
	Short: "Generar casos de uso con DTOs",
	Long: `Crea servicios de aplicación con DTOs bien definidos, 
interfaces claras y lógica de negocio encapsulada.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		usecaseName := args[0]

		entity, _ := cmd.Flags().GetString("entity")
		operations, _ := cmd.Flags().GetString("operations")
		dtoValidation, _ := cmd.Flags().GetBool("dto-validation")
		async, _ := cmd.Flags().GetBool("async")

		if entity == "" {
			fmt.Println("Error: --entity flag es requerido")
			os.Exit(1)
		}

		fmt.Printf("Generando caso de uso '%s' para entidad '%s'\n", usecaseName, entity)
		fmt.Printf("Operaciones: %s\n", operations)

		if dtoValidation {
			fmt.Println("✓ Incluyendo validaciones en DTOs")
		}
		if async {
			fmt.Println("✓ Incluyendo operaciones asíncronas")
		}

		generateUseCase(usecaseName, entity, operations, dtoValidation, async)
		fmt.Printf("\n✅ Caso de uso '%s' generado exitosamente!\n", usecaseName)
	},
}

func generateUseCase(usecaseName, entity, operations string, dtoValidation, async bool) {
	// Create usecase directory
	usecaseDir := filepath.Join("internal", "usecase")
	os.MkdirAll(usecaseDir, 0755)

	// Parse operations
	ops := parseOperations(operations)

	// Generate files
	generateDTOFile(usecaseDir, entity, ops, dtoValidation)
	generateUseCaseInterface(usecaseDir, usecaseName, entity, ops)
	generateUseCaseService(usecaseDir, usecaseName, entity, ops, async)
	generateUseCaseInterfaces(usecaseDir, entity)
}

func parseOperations(operations string) []string {
	if operations == "" {
		return []string{"create", "read"}
	}

	ops := strings.Split(operations, ",")
	var result []string
	for _, op := range ops {
		result = append(result, strings.TrimSpace(op))
	}
	return result
}

func generateDTOFile(dir, entity string, operations []string, validation bool) {
	filename := filepath.Join(dir, "dto.go")

	// Get the module name from go.mod
	moduleName := getModuleName()

	var content strings.Builder
	content.WriteString("package usecase\n\n")

	// Always import domain package since DTOs reference it
	content.WriteString("import (\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", moduleName))
	content.WriteString(")\n\n")

	// Generate DTOs for each operation
	for _, op := range operations {
		switch op {
		case "create":
			generateCreateDTO(&content, entity, validation)
		case "update":
			generateUpdateDTO(&content, entity, validation)
		case "read", "get":
			// Read operations typically don't need input DTOs, just output
		case "list":
			generateListDTO(&content, entity)
		}
	}

	writeFile(filename, content.String())
}

func generateCreateDTO(content *strings.Builder, entity string, validation bool) {
	entityLower := strings.ToLower(entity)

	content.WriteString(fmt.Sprintf("type Create%sInput struct {\n", entity))
	content.WriteString("\tName  string `json:\"name\"")
	if validation {
		content.WriteString(" validate:\"required,min=2\"")
	}
	content.WriteString("`\n")
	content.WriteString("\tEmail string `json:\"email\"")
	if validation {
		content.WriteString(" validate:\"required,email\"")
	}
	content.WriteString("`\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("type Create%sOutput struct {\n", entity))
	content.WriteString(fmt.Sprintf("\t%s    domain.%s `json:\"%s\"`\n", entity, entity, entityLower))
	content.WriteString("\tMessage string      `json:\"message\"`\n")
	content.WriteString("}\n\n")
}

func generateUpdateDTO(content *strings.Builder, entity string, validation bool) {
	content.WriteString(fmt.Sprintf("type Update%sInput struct {\n", entity))
	content.WriteString("\tName  string `json:\"name,omitempty\"")
	if validation {
		content.WriteString(" validate:\"omitempty,min=2\"")
	}
	content.WriteString("`\n")
	content.WriteString("\tEmail string `json:\"email,omitempty\"")
	if validation {
		content.WriteString(" validate:\"omitempty,email\"")
	}
	content.WriteString("`\n")
	content.WriteString("}\n\n")
}

func generateListDTO(content *strings.Builder, entity string) {
	entityLower := strings.ToLower(entity)
	content.WriteString(fmt.Sprintf("type List%sOutput struct {\n", entity))
	content.WriteString(fmt.Sprintf("\t%ss   []domain.%s `json:\"%ss\"`\n", entity, entity, entityLower))
	content.WriteString("\tTotal   int           `json:\"total\"`\n")
	content.WriteString("\tMessage string        `json:\"message\"`\n")
	content.WriteString("}\n\n")
}

func generateUseCaseInterface(dir, usecaseName, entity string, operations []string) {
	// Get the module name from go.mod
	moduleName := getModuleName()

	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, entityLower+"_usecase.go")

	var content strings.Builder
	content.WriteString("package usecase\n\n")
	content.WriteString(fmt.Sprintf("import \"%s/internal/domain\"\n\n", moduleName))

	content.WriteString(fmt.Sprintf("type %s interface {\n", usecaseName))

	for _, op := range operations {
		switch op {
		case "create":
			content.WriteString(fmt.Sprintf("\tCreate%s(input Create%sInput) (Create%sOutput, error)\n",
				entity, entity, entity))
		case "read", "get":
			content.WriteString(fmt.Sprintf("\tGet%s(id int) (*domain.%s, error)\n", entity, entity))
		case "update":
			content.WriteString(fmt.Sprintf("\tUpdate%s(id int, input Update%sInput) error\n", entity, entity))
		case "delete":
			content.WriteString(fmt.Sprintf("\tDelete%s(id int) error\n", entity))
		case "list":
			content.WriteString(fmt.Sprintf("\tList%ss() (List%sOutput, error)\n", entity, entity))
		}
	}

	content.WriteString("}\n")

	writeFile(filename, content.String())
}

func generateUseCaseService(dir, usecaseName, entity string, operations []string, async bool) {
	// Get the module name from go.mod
	moduleName := getModuleName()

	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, entityLower+"_service.go")

	var content strings.Builder
	content.WriteString("package usecase\n\n")
	content.WriteString("import (\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", moduleName))
	content.WriteString(fmt.Sprintf("\t\"%s/internal/messages\"\n", moduleName))
	content.WriteString(")\n\n")

	// Service struct
	serviceName := fmt.Sprintf("%sService", entityLower)
	content.WriteString(fmt.Sprintf("type %s struct {\n", serviceName))
	content.WriteString(fmt.Sprintf("\trepo %sRepository\n", entity))
	content.WriteString("}\n\n")

	// Constructor
	interfaceName := strings.Replace(usecaseName, "Service", "UseCase", 1)
	content.WriteString(fmt.Sprintf("func New%s(repo %sRepository) %s {\n",
		strings.ToUpper(string(serviceName[0]))+serviceName[1:], entity, interfaceName))
	content.WriteString(fmt.Sprintf("\treturn &%s{repo: repo}\n", serviceName))
	content.WriteString("}\n\n")

	// Generate methods for each operation
	for _, op := range operations {
		switch op {
		case "create":
			generateCreateMethod(&content, serviceName, entity)
		case "read", "get":
			generateGetMethod(&content, serviceName, entity)
		case "update":
			generateUpdateMethod(&content, serviceName, entity)
		case "delete":
			generateDeleteMethod(&content, serviceName, entity)
		case "list":
			generateListMethod(&content, serviceName, entity)
		}
	}

	writeFile(filename, content.String())
}

func generateCreateMethod(content *strings.Builder, serviceName, entity string) {
	entityLower := strings.ToLower(entity)
	serviceVar := string(serviceName[0])

	content.WriteString(fmt.Sprintf("func (%s *%s) Create%s(input Create%sInput) (Create%sOutput, error) {\n",
		serviceVar, serviceName, entity, entity, entity))
	content.WriteString(fmt.Sprintf("\t%s := domain.%s{\n", entityLower, entity))
	content.WriteString("\t\tName:  input.Name,\n")
	content.WriteString("\t\tEmail: input.Email,\n")
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\tif err := %s.Validate(); err != nil {\n", entityLower))
	content.WriteString(fmt.Sprintf("\t\treturn Create%sOutput{}, err\n", entity))
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\tif err := %s.repo.Save(&%s); err != nil {\n", serviceVar, entityLower))
	content.WriteString(fmt.Sprintf("\t\treturn Create%sOutput{}, err\n", entity))
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\treturn Create%sOutput{\n", entity))
	content.WriteString(fmt.Sprintf("\t\t%s:    %s,\n", entity, entityLower))
	content.WriteString(fmt.Sprintf("\t\tMessage: messages.%sCreatedSuccessfully,\n", entity))
	content.WriteString("\t}, nil\n")
	content.WriteString("}\n\n")
}

func generateGetMethod(content *strings.Builder, serviceName, entity string) {
	serviceVar := string(serviceName[0])

	content.WriteString(fmt.Sprintf("func (%s *%s) Get%s(id int) (*domain.%s, error) {\n",
		serviceVar, serviceName, entity, entity))
	content.WriteString(fmt.Sprintf("\treturn %s.repo.FindByID(id)\n", serviceVar))
	content.WriteString("}\n\n")
}

func generateUpdateMethod(content *strings.Builder, serviceName, entity string) {
	serviceVar := string(serviceName[0])

	content.WriteString(fmt.Sprintf("func (%s *%s) Update%s(id int, input Update%sInput) error {\n",
		serviceVar, serviceName, entity, entity))
	content.WriteString(fmt.Sprintf("\t%s, err := %s.repo.FindByID(id)\n", strings.ToLower(entity), serviceVar))
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn err\n")
	content.WriteString("\t}\n\n")

	content.WriteString("\tif input.Name != \"\" {\n")
	content.WriteString(fmt.Sprintf("\t\t%s.Name = input.Name\n", strings.ToLower(entity)))
	content.WriteString("\t}\n")
	content.WriteString("\tif input.Email != \"\" {\n")
	content.WriteString(fmt.Sprintf("\t\t%s.Email = input.Email\n", strings.ToLower(entity)))
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\treturn %s.repo.Update(%s)\n", serviceVar, strings.ToLower(entity)))
	content.WriteString("}\n\n")
}

func generateDeleteMethod(content *strings.Builder, serviceName, entity string) {
	serviceVar := string(serviceName[0])

	content.WriteString(fmt.Sprintf("func (%s *%s) Delete%s(id int) error {\n",
		serviceVar, serviceName, entity))
	content.WriteString(fmt.Sprintf("\treturn %s.repo.Delete(id)\n", serviceVar))
	content.WriteString("}\n\n")
}

func generateListMethod(content *strings.Builder, serviceName, entity string) {
	serviceVar := string(serviceName[0])
	entityLower := strings.ToLower(entity)

	content.WriteString(fmt.Sprintf("func (%s *%s) List%ss() (List%sOutput, error) {\n",
		serviceVar, serviceName, entity, entity))
	content.WriteString(fmt.Sprintf("\t%ss, err := %s.repo.FindAll()\n", entityLower, serviceVar))
	content.WriteString("\tif err != nil {\n")
	content.WriteString(fmt.Sprintf("\t\treturn List%sOutput{}, err\n", entity))
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\treturn List%sOutput{\n", entity))
	content.WriteString(fmt.Sprintf("\t\t%ss:   %ss,\n", entity, entityLower))
	content.WriteString(fmt.Sprintf("\t\tTotal:   len(%ss),\n", entityLower))
	content.WriteString(fmt.Sprintf("\t\tMessage: messages.%ssListedSuccessfully,\n", entity))
	content.WriteString("\t}, nil\n")
	content.WriteString("}\n\n")
}

func generateUseCaseInterfaces(dir, entity string) {
	// Get the module name from go.mod
	moduleName := getModuleName()

	filename := filepath.Join(dir, "interfaces.go")

	var content strings.Builder
	content.WriteString("package usecase\n\n")
	content.WriteString(fmt.Sprintf("import \"%s/internal/domain\"\n\n", moduleName))

	content.WriteString(fmt.Sprintf("type %sRepository interface {\n", entity))
	content.WriteString(fmt.Sprintf("\tSave(user *domain.%s) error\n", entity))
	content.WriteString(fmt.Sprintf("\tFindByID(id int) (*domain.%s, error)\n", entity))
	content.WriteString(fmt.Sprintf("\tFindByEmail(email string) (*domain.%s, error)\n", entity))
	content.WriteString(fmt.Sprintf("\tUpdate(user *domain.%s) error\n", entity))
	content.WriteString("\tDelete(id int) error\n")
	content.WriteString(fmt.Sprintf("\tFindAll() ([]domain.%s, error)\n", entity))
	content.WriteString("}\n")

	writeFile(filename, content.String())
}

func init() {
	usecaseCmd.Flags().StringP("entity", "e", "", "Entidad asociada al caso de uso (requerido)")
	usecaseCmd.Flags().StringP("operations", "o", "create,read", "Operaciones CRUD \"create,read,update,delete,list\"")
	usecaseCmd.Flags().BoolP("dto-validation", "d", false, "DTOs con validaciones específicas")
	usecaseCmd.Flags().BoolP("async", "a", false, "Incluir operaciones asíncronas")

	usecaseCmd.MarkFlagRequired("entity")
}
