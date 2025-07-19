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

		if fields == "" {
			fmt.Println("Error: --fields flag es requerido")
			os.Exit(1)
		}

		fmt.Printf("Generando entidad '%s'\n", entityName)
		fmt.Printf("Campos: %s\n", fields)

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
	// Create domain directory
	domainDir := filepath.Join("internal", "domain")
	os.MkdirAll(domainDir, 0755)

	// Parse fields
	fieldsList := parseFields(fields)

	// Add timestamps if requested
	if timestamps {
		fieldsList = append(fieldsList, Field{Name: "CreatedAt", Type: "time.Time", Tag: "`json:\"created_at\"`"})
		fieldsList = append(fieldsList, Field{Name: "UpdatedAt", Type: "*time.Time", Tag: "`json:\"updated_at,omitempty\"`"})
	}

	// Add soft delete if requested
	if softDelete {
		fieldsList = append(fieldsList, Field{Name: "DeletedAt", Type: "*time.Time", Tag: "`json:\"deleted_at,omitempty\"`"})
	}

	// Generate entity file
	generateEntityFile(domainDir, entityName, fieldsList, validation, businessRules, timestamps, softDelete)

	// Generate errors file if validation is enabled
	if validation {
		generateErrorsFile(domainDir, entityName)
	}
}

type Field struct {
	Name string
	Type string
	Tag  string
}

func parseFields(fields string) []Field {
	var fieldsList []Field

	// Always add ID field
	fieldsList = append(fieldsList, Field{
		Name: "ID",
		Type: "int",
		Tag:  "`json:\"id\"`",
	})

	parts := strings.Split(fields, ",")
	for _, part := range parts {
		fieldParts := strings.Split(strings.TrimSpace(part), ":")
		if len(fieldParts) == 2 {
			fieldName := strings.ToUpper(string(fieldParts[0][0])) + strings.ToLower(fieldParts[0][1:])
			fieldType := strings.TrimSpace(fieldParts[1])
			tag := fmt.Sprintf("`json:\"%s\"`", strings.ToLower(fieldName))

			fieldsList = append(fieldsList, Field{
				Name: fieldName,
				Type: fieldType,
				Tag:  tag,
			})
		}
	}

	return fieldsList
}

func generateEntityFile(dir, entityName string, fields []Field, validation, businessRules, timestamps, softDelete bool) {
	entityLower := strings.ToLower(entityName)
	filename := filepath.Join(dir, entityLower+".go")

	var content strings.Builder

	// Package and imports
	content.WriteString("package domain\n\n")

	if validation || timestamps || softDelete {
		content.WriteString("import (\n")
		if validation {
			content.WriteString("\t\"errors\"\n")
		}
		if timestamps || softDelete {
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

	writeFile(filename, content.String())
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

func generateErrorsFile(dir, entityName string) {
	filename := filepath.Join(dir, "errors.go")

	var content strings.Builder
	content.WriteString("package domain\n\n")
	content.WriteString("import \"errors\"\n\n")
	content.WriteString("var (\n")
	content.WriteString(fmt.Sprintf("\tErrInvalid%sData = errors.New(\"invalid %s data\")\n",
		entityName, strings.ToLower(entityName)))

	// Check if file exists and read existing errors
	if _, err := os.Stat(filename); err == nil {
		// File exists, we should append only new errors
		fmt.Printf("⚠️  Archivo errors.go ya existe, agregando nuevos errores para %s\n", entityName)
	}

	content.WriteString(fmt.Sprintf("\tErrInvalid%sName = errors.New(\"%s name is required\")\n",
		entityName, strings.ToLower(entityName)))
	content.WriteString(fmt.Sprintf("\tErrInvalid%sEmail = errors.New(\"%s email is required\")\n",
		entityName, strings.ToLower(entityName)))
	content.WriteString(fmt.Sprintf("\tErrInvalid%sAge = errors.New(\"%s age must be positive\")\n",
		entityName, strings.ToLower(entityName)))
	content.WriteString(")\n")

	writeFile(filename, content.String())
}

func init() {
	entityCmd.Flags().StringP("fields", "f", "", "Campos de la entidad \"name:type,email:string\" (requerido)")
	entityCmd.Flags().BoolP("validation", "v", false, "Agregar validaciones de dominio")
	entityCmd.Flags().BoolP("business-rules", "b", false, "Incluir métodos de reglas de negocio")
	entityCmd.Flags().BoolP("timestamps", "t", false, "Agregar campos created_at y updated_at")
	entityCmd.Flags().BoolP("soft-delete", "s", false, "Agregar funcionalidad de soft delete")

	entityCmd.MarkFlagRequired("fields")
}
