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
	Short: "Generate use cases with DTOs",
	Long: `Creates application services with well-defined DTOs, 
clear interfaces and encapsulated business logic.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		usecaseName := args[0]

		// Initialize configuration integration
		configIntegration := NewConfigIntegration()
		configIntegration.LoadConfigForProject()

		// Get CLI flag values
		entity, _ := cmd.Flags().GetString("entity")
		operations, _ := cmd.Flags().GetString("operations")
		dtoValidation, _ := cmd.Flags().GetBool("dto-validation")
		async, _ := cmd.Flags().GetBool("async")

		if entity == "" {
			fmt.Println("Error: --entity flag is required")
			os.Exit(1)
		}

		// Merge only explicitly changed CLI flags with config
		flags := map[string]interface{}{}
		if cmd.Flags().Changed("operations") {
			flags["operations"] = operations
		}
		if cmd.Flags().Changed("dto-validation") {
			flags["validation"] = dtoValidation
		}
		if cmd.Flags().Changed("async") {
			flags["async"] = async
		}

		if len(flags) > 0 {
			configIntegration.MergeWithCLIFlags(flags)
		}

		// Calculate effective values (config overrides CLI defaults)
		effectiveDtoValidation := dtoValidation
		if !cmd.Flags().Changed("dto-validation") && configIntegration.config != nil {
			effectiveDtoValidation = configIntegration.config.Generation.Validation.Enabled
		}

		effectiveBusinessRules := false
		if configIntegration.config != nil {
			effectiveBusinessRules = configIntegration.config.Generation.BusinessRules.Enabled
		}

		// Print configuration summary
		fmt.Printf("Generating use case '%s' for entity '%s'\n", usecaseName, entity)
		fmt.Printf("Operations: %s\n", operations)

		if configIntegration.config != nil {
			if !cmd.Flags().Changed("dto-validation") {
				fmt.Printf("  DTO Validation: %v (from config)\n", effectiveDtoValidation)
			}
			if effectiveBusinessRules {
				fmt.Printf("  Business Rules: enabled (from config)\n")
			}
		}

		if effectiveDtoValidation {
			fmt.Println("✓ Including DTO validations")
		}
		if async {
			fmt.Println("✓ Including asynchronous operations")
		}

		generateUseCase(usecaseName, entity, operations, effectiveDtoValidation, async)
		fmt.Printf("\nUse case '%s' generated successfully!\n", usecaseName)
	},
}

func generateUseCase(usecaseName, entity, operations string, dtoValidation, async bool) {
	generateUseCaseWithFields(usecaseName, entity, operations, dtoValidation, async, "")
}

func generateUseCaseWithFields(usecaseName, entity, operations string, dtoValidation, async bool, fields string) {
	// Create usecase directory
	usecaseDir := filepath.Join(DirInternal, DirUseCase)
	_ = os.MkdirAll(usecaseDir, 0755)

	// Parse operations
	ops := parseOperations(operations)

	// Generate files
	generateDTOFileWithFields(usecaseDir, entity, ops, dtoValidation, fields)
	generateUseCaseInterface(usecaseDir, usecaseName, entity, ops)
	generateUseCaseServiceWithFields(usecaseDir, usecaseName, entity, ops, async, fields)
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

func generateDTOFileWithFields(dir, entity string, operations []string, validation bool, fields string) {
	filename := filepath.Join(dir, "dto.go")

	// Get the module name from go.mod
	moduleName := getModuleName()

	var content strings.Builder

	// Check if dto.go already exists
	if _, err := os.Stat(filename); err == nil {
		// File exists, read its content
		existingContent, err := os.ReadFile(filename)
		if err == nil {
			existingStr := string(existingContent)
			// Check if DTOs for this entity already exist
			createDTOName := fmt.Sprintf("type Create%sInput struct", entity)
			if strings.Contains(existingStr, createDTOName) {
				// DTOs already exist, don't regenerate
				return
			}

			// Add the existing content without the final newline
			content.WriteString(strings.TrimSuffix(existingStr, "\n"))
			content.WriteString("\n\n")
		}
	} else {
		// File doesn't exist, create header
		content.WriteString("package usecase\n\n")
		content.WriteString("import (\n")
		content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
		content.WriteString(")\n\n")
	}

	// Generate DTOs for each operation
	for _, op := range operations {
		switch op {
		case OpCreate:
			if fields != "" {
				generateCreateDTOWithFields(&content, entity, validation, fields)
			} else {
				generateCreateDTO(&content, entity, validation)
			}
		case OpUpdate:
			if fields != "" {
				generateUpdateDTOWithFields(&content, entity, validation, fields)
			} else {
				generateUpdateDTO(&content, entity, validation)
			}
		case OpRead, OperationGet:
			// Read operations typically don't need input DTOs, just output
		case OpList:
			generateListDTO(&content, entity)
		}
	}

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating DTO file: %v\n", err)
	}
}

func generateCreateDTO(content *strings.Builder, entity string, validation bool) {
	entityLower := strings.ToLower(entity)

	fmt.Fprintf(content, "type Create%sInput struct {\n", entity)
	// Campos estándar cuando no se especifican fields personalizados
	content.WriteString("\tNombre      string `json:\"nombre\"")
	if validation {
		content.WriteString(" validate:\"required,min=2\"")
	}
	content.WriteString("`\n")
	content.WriteString("\tDescripcion string `json:\"descripcion\"")
	if validation {
		content.WriteString(" validate:\"required,min=5\"")
	}
	content.WriteString("`\n")
	content.WriteString("}\n\n")

	fmt.Fprintf(content, "type Create%sOutput struct {\n", entity)
	fmt.Fprintf(content, "\t%s    domain.%s `json:\"%s\"`\n", entity, entity, entityLower)
	content.WriteString("\tMessage string      `json:\"message\"`\n")
	content.WriteString("}\n\n")
}

func generateUpdateDTO(content *strings.Builder, entity string, validation bool) {
	fmt.Fprintf(content, "type Update%sInput struct {\n", entity)

	// Generar campos estándar en español para cuando no se especifican fields personalizados
	content.WriteString("\tNombre      *string `json:\"nombre,omitempty\"")
	if validation {
		content.WriteString(" validate:\"omitempty,min=2\"")
	}
	content.WriteString("`\n")
	content.WriteString("\tDescripcion *string `json:\"descripcion,omitempty\"")
	if validation {
		content.WriteString(" validate:\"omitempty,min=5\"")
	}
	content.WriteString("`\n")
	content.WriteString("}\n\n")
}

func generateListDTO(content *strings.Builder, entity string) {
	entityLower := strings.ToLower(entity)
	fmt.Fprintf(content, "type List%sOutput struct {\n", entity)
	fmt.Fprintf(content, "\t%ss   []domain.%s `json:\"%ss\"`\n", entity, entity, entityLower)
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
	content.WriteString(fmt.Sprintf("import \"%s/internal/domain\"\n\n", getImportPath(moduleName)))

	// DEBUG: Print what interface name is being used
	fmt.Printf("DEBUG: Generating interface with name: %s\n", usecaseName)
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

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating use case file: %v\n", err)
	}
}

func generateUseCaseServiceWithFields(dir, usecaseName, entity string, operations []string, async bool, fields string) {
	// Get the module name from go.mod
	moduleName := getModuleName()

	entityLower := strings.ToLower(entity)
	filename := filepath.Join(dir, entityLower+"_service.go")

	var content strings.Builder
	content.WriteString("package usecase\n\n")
	content.WriteString("import (\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/domain\"\n", getImportPath(moduleName)))
	content.WriteString(fmt.Sprintf("\t\"%s/internal/messages\"\n", getImportPath(moduleName)))
	content.WriteString(fmt.Sprintf("\t\"%s/internal/repository\"\n", getImportPath(moduleName)))
	content.WriteString(")\n\n")

	// Service struct
	serviceName := fmt.Sprintf("%sService", entityLower)
	content.WriteString(fmt.Sprintf("type %s struct {\n", serviceName))
	content.WriteString(fmt.Sprintf("\trepo repository.%sRepository\n", entity))
	if async {
		content.WriteString("\t// Canal para procesamiento asíncrono\n")
		content.WriteString("\tasyncChannel chan AsyncTask\n")
		content.WriteString("\t// Logger para operaciones asíncronas\n")
		content.WriteString("\tlogger       Logger\n")
	}
	content.WriteString("}\n\n")

	// Constructor
	interfaceName := strings.Replace(usecaseName, "Service", "UseCase", 1)
	content.WriteString(fmt.Sprintf("func New%s(repo repository.%sRepository) %s {\n",
		strings.ToUpper(string(serviceName[0]))+serviceName[1:], entity, interfaceName))
	content.WriteString(fmt.Sprintf("\treturn &%s{repo: repo}\n", serviceName))
	content.WriteString("}\n\n")

	// Generate methods for each operation
	for _, op := range operations {
		switch op {
		case "create":
			if fields != "" {
				generateCreateMethodWithFields(&content, serviceName, entity, fields)
			} else {
				generateCreateMethod(&content, serviceName, entity)
			}
		case "read", "get":
			generateGetMethod(&content, serviceName, entity)
		case "update":
			if fields != "" {
				generateUpdateMethodWithFields(&content, serviceName, entity, fields)
			} else {
				generateUpdateMethod(&content, serviceName, entity)
			}
		case "delete":
			generateDeleteMethod(&content, serviceName, entity)
		case "list":
			generateListMethod(&content, serviceName, entity)
		}
	}

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating use case service with fields: %v\n", err)
	}
}

func generateCreateMethod(content *strings.Builder, serviceName, entity string) {
	entityLower := strings.ToLower(entity)
	serviceVar := string(serviceName[0])

	fmt.Fprintf(content, "func (%s *%s) Create%s(input Create%sInput) (Create%sOutput, error) {\n",
		serviceVar, serviceName, entity, entity, entity)
	fmt.Fprintf(content, "\t%s := domain.%s{\n", entityLower, entity)
	content.WriteString("\t\t// Mapeo automático de campos - ajustar según tu entidad\n")
	content.WriteString("\t\t// Nombre: input.Nombre,\n")
	content.WriteString("\t\t// Email: input.Email,\n")
	content.WriteString("\t\t// Edad: input.Edad,\n")
	content.WriteString("\t}\n\n")

	fmt.Fprintf(content, "\tif err := %s.Validate(); err != nil {\n", entityLower)
	fmt.Fprintf(content, "\t\treturn Create%sOutput{}, err\n", entity)
	content.WriteString("\t}\n\n")

	fmt.Fprintf(content, "\tif err := %s.repo.Save(&%s); err != nil {\n", serviceVar, entityLower)
	fmt.Fprintf(content, "\t\treturn Create%sOutput{}, err\n", entity)
	content.WriteString("\t}\n\n")

	fmt.Fprintf(content, "\treturn Create%sOutput{\n", entity)
	fmt.Fprintf(content, "\t\t%s:    %s,\n", entity, entityLower)
	fmt.Fprintf(content, "\t\tMessage: messages.%sCreatedSuccessfully,\n", entity)
	content.WriteString("\t}, nil\n")
	content.WriteString("}\n\n")
}

func generateCreateMethodWithFields(content *strings.Builder, serviceName, entity, fields string) {
	entityLower := strings.ToLower(entity)
	serviceVar := string(serviceName[0])
	fieldsList := parseFields(fields)

	fmt.Fprintf(content, "func (%s *%s) Create%s(input Create%sInput) (Create%sOutput, error) {\n",
		serviceVar, serviceName, entity, entity, entity)
	fmt.Fprintf(content, "\t%s := domain.%s{\n", entityLower, entity)

	// Map fields from input to entity
	for _, field := range fieldsList {
		if field.Name == "ID" {
			continue // Skip ID, it's auto-generated
		}
		fmt.Fprintf(content, "\t\t%s: input.%s,\n", field.Name, field.Name)
	}

	content.WriteString("\t}\n\n")

	fmt.Fprintf(content, "\tif err := %s.Validate(); err != nil {\n", entityLower)
	fmt.Fprintf(content, "\t\treturn Create%sOutput{}, err\n", entity)
	content.WriteString("\t}\n\n")

	fmt.Fprintf(content, "\tif err := %s.repo.Save(&%s); err != nil {\n", serviceVar, entityLower)
	fmt.Fprintf(content, "\t\treturn Create%sOutput{}, err\n", entity)
	content.WriteString("\t}\n\n")

	fmt.Fprintf(content, "\treturn Create%sOutput{\n", entity)
	content.WriteString("\t\tID:      " + entityLower + ".ID,\n")

	// Map fields from entity to output
	for _, field := range fieldsList {
		if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
			continue
		}
		fmt.Fprintf(content, "\t\t%s: %s.%s,\n", field.Name, entityLower, field.Name)
	}

	fmt.Fprintf(content, "\t\tMessage: messages.%sCreatedSuccessfully,\n", entity)
	content.WriteString("\t}, nil\n")
	content.WriteString("}\n\n")
}

func generateGetMethod(content *strings.Builder, serviceName, entity string) {
	serviceVar := string(serviceName[0])

	fmt.Fprintf(content, "func (%s *%s) Get%s(id int) (*domain.%s, error) {\n",
		serviceVar, serviceName, entity, entity)
	fmt.Fprintf(content, "\treturn %s.repo.FindByID(id)\n", serviceVar)
	content.WriteString("}\n\n")
}

func generateUpdateMethodWithFields(content *strings.Builder, serviceName, entity, fields string) {
	serviceVar := string(serviceName[0])
	entityVar := strings.ToLower(entity)
	fieldsList := parseFields(fields)

	fmt.Fprintf(content, "func (%s *%s) Update%s(id int, input Update%sInput) error {\n",
		serviceVar, serviceName, entity, entity)
	fmt.Fprintf(content, "\t%s, err := %s.repo.FindByID(id)\n", entityVar, serviceVar)
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn err\n")
	content.WriteString("\t}\n\n")

	// Update fields based on actual entity fields
	// In UpdateInput DTOs, fields are always pointers (optional updates)
	for _, field := range fieldsList {
		if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
			continue
		}

		// Fields in UpdateInput are always pointers, check if not nil
		fmt.Fprintf(content, "\tif input.%s != nil {\n", field.Name)
		fmt.Fprintf(content, "\t\t%s.%s = *input.%s\n", entityVar, field.Name, field.Name)
		content.WriteString("\t}\n")
	}

	content.WriteString("\n")
	fmt.Fprintf(content, "\treturn %s.repo.Update(%s)\n", serviceVar, entityVar)
	content.WriteString("}\n\n")
}

func generateUpdateMethod(content *strings.Builder, serviceName, entity string) {
	serviceVar := string(serviceName[0])
	entityVar := strings.ToLower(entity)

	fmt.Fprintf(content, "func (%s *%s) Update%s(id int, input Update%sInput) error {\n",
		serviceVar, serviceName, entity, entity)
	fmt.Fprintf(content, "\t%s, err := %s.repo.FindByID(id)\n", entityVar, serviceVar)
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn err\n")
	content.WriteString("\t}\n\n")

	content.WriteString("\t// Actualizar campos según tu entidad\n")
	content.WriteString("\tif input.Nombre != \"\" {\n")
	fmt.Fprintf(content, "\t\t%s.Nombre = input.Nombre\n", entityVar)
	content.WriteString("\t}\n")
	content.WriteString("\tif input.Email != \"\" {\n")
	fmt.Fprintf(content, "\t\t%s.Email = input.Email\n", entityVar)
	content.WriteString("\t}\n")
	content.WriteString("\t// Agregar más campos según necesites\n\n")

	fmt.Fprintf(content, "\treturn %s.repo.Update(%s)\n", serviceVar, entityVar)
	content.WriteString("}\n\n")
}

func generateDeleteMethod(content *strings.Builder, serviceName, entity string) {
	serviceVar := string(serviceName[0])

	fmt.Fprintf(content, "func (%s *%s) Delete%s(id int) error {\n",
		serviceVar, serviceName, entity)
	fmt.Fprintf(content, "\treturn %s.repo.Delete(id)\n", serviceVar)
	content.WriteString("}\n\n")
}

func generateListMethod(content *strings.Builder, serviceName, entity string) {
	serviceVar := string(serviceName[0])
	entityLower := strings.ToLower(entity)

	fmt.Fprintf(content, "func (%s *%s) List%ss() (List%sOutput, error) {\n",
		serviceVar, serviceName, entity, entity)
	fmt.Fprintf(content, "\t%ss, err := %s.repo.FindAll()\n", entityLower, serviceVar)
	content.WriteString("\tif err != nil {\n")
	fmt.Fprintf(content, "\t\treturn List%sOutput{}, err\n", entity)
	content.WriteString("\t}\n\n")

	fmt.Fprintf(content, "\treturn List%sOutput{\n", entity)
	fmt.Fprintf(content, "\t\t%ss:   %ss,\n", entity, entityLower)
	fmt.Fprintf(content, "\t\tTotal:   len(%ss),\n", entityLower)
	fmt.Fprintf(content, "\t\tMessage: messages.%ssListedSuccessfully,\n", entity)
	content.WriteString("\t}, nil\n")
	content.WriteString("}\n\n")
}

func generateUseCaseInterfaces(dir, entity string) {
	// Get the module name from go.mod
	moduleName := getModuleName()

	filename := filepath.Join(dir, "interfaces.go")

	var content strings.Builder
	content.WriteString("package usecase\n\n")
	content.WriteString(fmt.Sprintf("import \"%s/internal/domain\"\n\n", getImportPath(moduleName)))

	entityLower := strings.ToLower(entity)
	content.WriteString(fmt.Sprintf("type %sRepository interface {\n", entity))
	content.WriteString(fmt.Sprintf("\tSave(%s *domain.%s) error\n", entityLower, entity))
	content.WriteString(fmt.Sprintf("\tFindByID(id int) (*domain.%s, error)\n", entity))
	content.WriteString(fmt.Sprintf("\tFindByEmail(email string) (*domain.%s, error)\n", entity))
	content.WriteString(fmt.Sprintf("\tUpdate(%s *domain.%s) error\n", entityLower, entity))
	content.WriteString("\tDelete(id int) error\n")
	content.WriteString(fmt.Sprintf("\tFindAll() ([]domain.%s, error)\n", entity))
	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating interfaces file: %v\n", err)
	}
}

func generateCreateDTOWithFields(content *strings.Builder, entity string, validation bool, fields string) {
	fieldsList := parseFields(fields)

	// Generate Create Input DTO
	fmt.Fprintf(content, "// Create%sInput DTO para crear un nuevo %s\n", entity, strings.ToLower(entity))
	fmt.Fprintf(content, "type Create%sInput struct {\n", entity)

	for _, field := range fieldsList {
		// Skip auto-managed fields in create request
		if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
			continue
		}

		jsonTag := fmt.Sprintf("json:\"%s\"", strings.ToLower(field.Name))

		if validation {
			validateTag := getValidationTag(field.Type)
			fmt.Fprintf(content, "\t%s %s `%s validate:\"%s\"`\n",
				field.Name, field.Type, jsonTag, validateTag)
		} else {
			fmt.Fprintf(content, "\t%s %s `%s`\n",
				field.Name, field.Type, jsonTag)
		}
	}

	content.WriteString("}\n\n")

	// Generate validation method for the DTO
	if validation {
		fmt.Fprintf(content, "// Validate valida los datos del DTO Create%sInput\n", entity)
		fmt.Fprintf(content, "func (r *Create%sInput) Validate() error {\n", entity)

		for _, field := range fieldsList {
			if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
				continue
			}

			switch field.Type {
			case "string":
				if strings.Contains(strings.ToLower(field.Name), "email") {
					fmt.Fprintf(content, "\tif r.%s == \"\" {\n", field.Name)
					fmt.Fprintf(content, "\t\treturn errors.New(\"%s es requerido\")\n", getSpanishFieldName(strings.ToLower(field.Name)))
					content.WriteString("\t}\n")
					fmt.Fprintf(content, "\tif !strings.Contains(r.%s, \"@\") {\n", field.Name)
					fmt.Fprintf(content, "\t\treturn errors.New(\"formato de %s inválido\")\n", getSpanishFieldName(strings.ToLower(field.Name)))
					content.WriteString("\t}\n")
				} else {
					fmt.Fprintf(content, "\tif strings.TrimSpace(r.%s) == \"\" {\n", field.Name)
					fmt.Fprintf(content, "\t\treturn errors.New(\"%s es requerido\")\n", getSpanishFieldName(strings.ToLower(field.Name)))
					content.WriteString("\t}\n")
				}
			case "int", "int64", "uint", "uint64":
				fmt.Fprintf(content, "\tif r.%s < 0 {\n", field.Name)
				fmt.Fprintf(content, "\t\treturn errors.New(\"%s debe ser un número positivo\")\n", getSpanishFieldName(strings.ToLower(field.Name)))
				content.WriteString("\t}\n")
			case "float64", "float32":
				fmt.Fprintf(content, "\tif r.%s < 0 {\n", field.Name)
				fmt.Fprintf(content, "\t\treturn errors.New(\"%s debe ser un número positivo\")\n", getSpanishFieldName(strings.ToLower(field.Name)))
				content.WriteString("\t}\n")
			}
		}

		content.WriteString("\treturn nil\n")
		content.WriteString("}\n\n")
	}

	// Generate Create Output DTO
	fmt.Fprintf(content, "// Create%sOutput DTO para la respuesta de creación\n", entity)
	fmt.Fprintf(content, "type Create%sOutput struct {\n", entity)
	content.WriteString("\tID      uint   `json:\"id\"`\n")

	// Add actual fields to response
	for _, field := range fieldsList {
		if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
			continue
		}
		jsonTag := fmt.Sprintf("json:\"%s\"", strings.ToLower(field.Name))
		fmt.Fprintf(content, "\t%s %s `%s`\n", field.Name, field.Type, jsonTag)
	}

	content.WriteString("\tMessage string `json:\"message\"`\n")
	content.WriteString("}\n\n")
}

func generateUpdateDTOWithFields(content *strings.Builder, entity string, validation bool, fields string) {
	fieldsList := parseFields(fields)

	// Generate Update Input DTO (fields are optional)
	fmt.Fprintf(content, "type Update%sInput struct {\n", entity)

	for _, field := range fieldsList {
		// Skip ID field in update input (it's in the URL)
		if field.Name == "ID" {
			continue
		}

		// Make fields optional for update (pointers)
		var fieldType string
		switch field.Type {
		case "string":
			fieldType = "*string"
		case "int":
			fieldType = "*int"
		case "bool":
			fieldType = "*bool"
		case "float64":
			fieldType = "*float64"
		default:
			fieldType = "*" + field.Type
		}

		jsonTag := fmt.Sprintf("json:\"%s,omitempty\"", strings.ToLower(field.Name))

		if validation {
			validateTag := "omitempty," + getValidationTag(field.Type)
			fmt.Fprintf(content, "\t%s %s `%s validate:\"%s\"`\n",
				field.Name, fieldType, jsonTag, validateTag)
		} else {
			fmt.Fprintf(content, "\t%s %s `%s`\n",
				field.Name, fieldType, jsonTag)
		}
	}

	content.WriteString("}\n\n")
}

func getValidationTag(fieldType string) string {
	switch fieldType {
	case "string":
		return "required,min=1"
	case "int", "int64", "uint", "uint64":
		return "required,min=1"
	case "float64", "float32":
		return "required,min=0"
	case "bool":
		return ""
	case "time.Time":
		return "required"
	default:
		return "required"
	}
}

func init() {
	usecaseCmd.Flags().StringP("entity", "e", "", "Associated entity for the use case (required)")
	usecaseCmd.Flags().StringP("operations", "o", "create,read", "CRUD operations \"create,read,update,delete,list\"")
	usecaseCmd.Flags().BoolP("dto-validation", "d", false, "DTOs with specific validations")
	usecaseCmd.Flags().BoolP("async", "a", false, "Include asynchronous operations")

	_ = usecaseCmd.MarkFlagRequired("entity")
}
