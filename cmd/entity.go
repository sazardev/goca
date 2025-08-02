package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var entityCmd = &cobra.Command{
	Use:   "entity <name>",
	Short: "Generar entidad de dominio pura",
	Long: `Crea entidades de dominio siguiendo los principios DDD, 
sin dependencias externas y con validaciones de negocio.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entityName := args[0]

		fields, _ := cmd.Flags().GetString("fields")
		validation, _ := cmd.Flags().GetBool("validation")
		businessRules, _ := cmd.Flags().GetBool("business-rules")
		timestamps, _ := cmd.Flags().GetBool("timestamps")
		softDelete, _ := cmd.Flags().GetBool("soft-delete")

		// Validar campos con el nuevo validador robusto
		validator := NewFieldValidator()

		if err := validator.ValidateEntityName(entityName); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

		if fields == "" {
			fmt.Println("❌ Error: --fields flag es requerido")
			os.Exit(1)
		}

		if err := validator.ValidateFields(fields); err != nil {
			fmt.Printf("❌ Error en campos: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🚀 Generando entidad '%s'\n", entityName)
		fmt.Printf("📋 Campos: %s\n", fields)

		if validation {
			fmt.Println("✓ Incluyendo validaciones")
		}
		if businessRules {
			fmt.Println("✓ Incluyendo reglas de negocio")
		}
		if timestamps {
			fmt.Println("✓ Incluyendo timestamps")
		}
		if softDelete {
			fmt.Println("✓ Incluyendo soft delete")
		}

		generateEntity(entityName, fields, validation, businessRules, timestamps, softDelete)
		fmt.Printf("\n✅ Entidad '%s' generada exitosamente!\n", entityName)
	},
}

func generateEntity(entityName, fields string, validation, businessRules, timestamps, softDelete bool) {
	// Crear directorio domain si no existe
	domainDir := "internal/domain"
	_ = os.MkdirAll(domainDir, 0755)

	// Parse fields
	fieldsList := parseFields(fields)

	// Add timestamps if requested
	if timestamps {
		fieldsList = append(fieldsList, Field{Name: "CreatedAt", Type: "time.Time", Tag: "`json:\"created_at\" gorm:\"autoCreateTime\"`"})
		fieldsList = append(fieldsList, Field{Name: "UpdatedAt", Type: "time.Time", Tag: "`json:\"updated_at\" gorm:\"autoUpdateTime\"`"})
	}

	// Add soft delete if requested
	if softDelete {
		fieldsList = append(fieldsList, Field{Name: "DeletedAt", Type: "gorm.DeletedAt", Tag: "`json:\"deleted_at,omitempty\" gorm:\"index\"`"})
	}

	// Generate entity file
	generateEntityFile(domainDir, entityName, fieldsList, validation, businessRules, timestamps, softDelete)

	// Generate errors file if validation is enabled
	if validation {
		generateErrorsFile(domainDir, entityName, fieldsList)
	}
}

type Field struct {
	Name string
	Type string
	Tag  string
}

func parseFields(fields string) []Field {
	validator := NewFieldValidator()
	fieldsList, err := validator.ParseFieldsWithValidation(fields)
	if err != nil {
		fmt.Printf("❌ Error en validación de campos: %v\n", err)
		os.Exit(1)
	}
	return fieldsList
}

func getGormTag(fieldName, fieldType string) string {
	switch fieldType {
	case "string":
		if fieldName == "Email" {
			return "type:varchar(255);uniqueIndex;not null"
		}
		if fieldName == "Title" || fieldName == "Name" {
			return "type:varchar(255);not null"
		}
		if fieldName == "Description" {
			return "type:text"
		}
		return "type:varchar(255)"
	case "int":
		return "type:integer;not null;default:0"
	case "bool":
		return "type:boolean;not null;default:false"
	case "float64":
		return "type:decimal(10,2);not null;default:0"
	default:
		return "not null"
	}
}

// hasStringBusinessRules checks if any field will require the strings package for business rules
func hasStringBusinessRules(fields []Field) bool {
	for _, field := range fields {
		if field.Name == "Email" {
			return true
		}
	}
	return false
}

func generateEntityFile(dir, entityName string, fields []Field, validation, businessRules, timestamps, softDelete bool) {
	entityLower := strings.ToLower(entityName)
	filename := filepath.Join(dir, entityLower+".go")

	var content strings.Builder

	// Package and imports
	content.WriteString("package domain\n\n")

	// Determine which imports are needed
	needsTime := timestamps || softDelete
	needsStrings := businessRules && hasStringBusinessRules(fields)

	if needsTime || needsStrings {
		content.WriteString("import (\n")
		if needsStrings {
			content.WriteString("\t\"strings\"\n")
		}
		if needsTime {
			content.WriteString("\t\"time\"\n")
		}
		content.WriteString(")\n\n")
	}

	// Entity struct
	content.WriteString(fmt.Sprintf("type %s struct {\n", entityName))
	for _, field := range fields {
		content.WriteString(fmt.Sprintf("\t%s %s %s\n", field.Name, field.Type, field.Tag))
	}
	content.WriteString("}\n\n")

	// Validation method
	if validation {
		content.WriteString(fmt.Sprintf("func (%s *%s) Validate() error {\n", strings.ToLower(string(entityName[0])), entityName))

		for _, field := range fields {
			if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
				continue
			}

			switch field.Type {
			case "string":
				content.WriteString(fmt.Sprintf("\tif %s.%s == \"\" {\n", strings.ToLower(string(entityName[0])), field.Name))
				content.WriteString(fmt.Sprintf("\t\treturn ErrInvalid%s%s\n", entityName, field.Name))
				content.WriteString("\t}\n")
			case "int", "int64", "float64":
				content.WriteString(fmt.Sprintf("\tif %s.%s < 0 {\n", strings.ToLower(string(entityName[0])), field.Name))
				content.WriteString(fmt.Sprintf("\t\treturn ErrInvalid%s%s\n", entityName, field.Name))
				content.WriteString("\t}\n")
			}
		}

		content.WriteString("\treturn nil\n")
		content.WriteString("}\n\n")
	}

	// Business rules methods
	if businessRules {
		generateBusinessRules(&content, entityName, fields)
	}

	// Soft delete methods
	if softDelete {
		content.WriteString(fmt.Sprintf("func (%s *%s) SoftDelete() {\n", strings.ToLower(string(entityName[0])), entityName))
		content.WriteString("\tnow := time.Now()\n")
		content.WriteString(fmt.Sprintf("\t%s.DeletedAt = &now\n", strings.ToLower(string(entityName[0]))))
		content.WriteString("}\n\n")

		content.WriteString(fmt.Sprintf("func (%s *%s) IsDeleted() bool {\n", strings.ToLower(string(entityName[0])), entityName))
		content.WriteString(fmt.Sprintf("\treturn %s.DeletedAt != nil\n", strings.ToLower(string(entityName[0]))))
		content.WriteString("}\n\n")
	}

	writeGoFile(filename, content.String())
}

func generateBusinessRules(content *strings.Builder, entityName string, fields []Field) {
	entityVar := strings.ToLower(string(entityName[0]))

	// Generate some common business rules based on field types
	for _, field := range fields {
		switch {
		case field.Name == "Age":
			content.WriteString(fmt.Sprintf("func (%s *%s) IsAdult() bool {\n", entityVar, entityName))
			content.WriteString(fmt.Sprintf("\treturn %s.Age >= 18\n", entityVar))
			content.WriteString("}\n\n")

		case field.Name == "Price":
			content.WriteString(fmt.Sprintf("func (%s *%s) IsExpensive() bool {\n", entityVar, entityName))
			content.WriteString(fmt.Sprintf("\treturn %s.Price > 1000.0\n", entityVar))
			content.WriteString("}\n\n")

		case field.Name == "Email":
			content.WriteString(fmt.Sprintf("func (%s *%s) HasValidEmail() bool {\n", entityVar, entityName))
			content.WriteString(fmt.Sprintf("\treturn strings.Contains(%s.Email, \"@\")\n", entityVar))
			content.WriteString("}\n\n")

		case field.Name == "Status":
			content.WriteString(fmt.Sprintf("func (%s *%s) IsActive() bool {\n", entityVar, entityName))
			content.WriteString(fmt.Sprintf("\treturn %s.Status == \"active\"\n", entityVar))
			content.WriteString("}\n\n")
		}
	}
}

func generateErrorsFile(dir, entityName string, fields []Field) {
	filename := filepath.Join(dir, "errors.go")

	var content strings.Builder
	var existingErrors []string

	// Check if file exists and read existing errors
	if _, err := os.Stat(filename); err == nil {
		// File exists, read existing content
		fmt.Printf("⚠️  Archivo errors.go ya existe, agregando nuevos errores para %s\n", entityName)

		if existingContent, err := os.ReadFile(filename); err == nil {
			// Extract existing error declarations
			lines := strings.Split(string(existingContent), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "ErrInvalid") && strings.Contains(line, "errors.New") {
					existingErrors = append(existingErrors, "\t"+line)
				}
			}
		}
	}

	content.WriteString("package domain\n\n")
	content.WriteString("import \"errors\"\n\n")
	content.WriteString("var (\n")

	// Add general error if not exists
	generalError := fmt.Sprintf("\tErrInvalid%sData = errors.New(\"invalid %s data\")",
		entityName, strings.ToLower(entityName))
	if !contains(existingErrors, generalError) {
		content.WriteString(generalError + "\n")
	}

	// Add existing errors
	for _, err := range existingErrors {
		content.WriteString(err + "\n")
	}

	// Generate specific validation errors for all fields
	for _, field := range fields {
		if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
			continue
		}
		newError := fmt.Sprintf("\tErrInvalid%s%s = errors.New(\"%s %s is invalid\")",
			entityName, field.Name, strings.ToLower(entityName), strings.ToLower(field.Name))

		if !contains(existingErrors, newError) {
			content.WriteString(newError + "\n")
		}
	}
	content.WriteString(")\n")

	writeGoFile(filename, content.String())
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == strings.TrimSpace(item) {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(entityCmd)
	entityCmd.Flags().StringP("fields", "f", "", "Campos de la entidad \"field:type,field2:type\" (requerido)")
	entityCmd.Flags().BoolP("validation", "v", false, "Incluir validaciones de negocio")
	entityCmd.Flags().BoolP("business-rules", "b", false, "Incluir reglas de negocio avanzadas")
	entityCmd.Flags().BoolP("timestamps", "t", false, "Incluir campos CreatedAt y UpdatedAt")
	entityCmd.Flags().BoolP("soft-delete", "s", false, "Incluir soft delete (DeletedAt)")
	_ = entityCmd.MarkFlagRequired("fields")
}
