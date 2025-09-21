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
	Short: "Generate pure domain entity",
	Long: `Creates domain entities following DDD principles, 
without external dependencies and with complete business validations.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entityName := args[0]

		fields, _ := cmd.Flags().GetString("fields")
		validation, _ := cmd.Flags().GetBool("validation")
		businessRules, _ := cmd.Flags().GetBool("business-rules")
		timestamps, _ := cmd.Flags().GetBool("timestamps")
		softDelete, _ := cmd.Flags().GetBool("soft-delete")

		// Usar validador centralizado
		validator := NewCommandValidator()

		if err := validator.ValidateEntityCommand(entityName, fields); err != nil {
			validator.errorHandler.HandleError(err, "validación de parámetros")
		}

		validator.errorHandler.ValidateRequiredFlag(fields, "fields")

		fmt.Printf("🏗️  Generando entidad '%s'\n", entityName)
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
			fmt.Println("✓ Incluyendo eliminación suave")
		}

		generateEntity(entityName, fields, validation, businessRules, timestamps, softDelete)

		// Generar datos de semilla automáticamente
		if fields != "" {
			generateSeedData("internal/domain", entityName, parseFields(fields))
			fmt.Println("🌱 Datos de semilla generados")
		}

		fmt.Printf("\n✅ Entidad '%s' generada exitosamente!\n", entityName)
		fmt.Printf("📁 Archivos creados:\n")
		fmt.Printf("   - internal/domain/%s.go\n", strings.ToLower(entityName))
		if validation {
			fmt.Printf("   - internal/domain/errors.go\n")
		}
		fmt.Printf("   - internal/domain/%s_seeds.go\n", strings.ToLower(entityName))
		fmt.Println("\n🎉 ¡Todo listo! Tu entidad está lista para usar.")
	},
}

func generateEntity(entityName, fields string, validation, businessRules, timestamps, softDelete bool) {
	// Crear directorio domain si no existe
	domainDir := "internal/domain"
	_ = os.MkdirAll(domainDir, 0755)

	// Parse fields - ahora genera campos reales basados en el input
	fieldsList := parseFields(fields)

	// Add ID field always as first field
	idField := Field{Name: "ID", Type: "uint", Tag: "`json:\"id\" gorm:\"primaryKey\"`"}
	fieldsList = append([]Field{idField}, fieldsList...)

	// Add timestamps if requested
	if timestamps {
		fieldsList = append(fieldsList, Field{Name: "CreatedAt", Type: "time.Time", Tag: "`json:\"created_at\" gorm:\"autoCreateTime\"`"})
		fieldsList = append(fieldsList, Field{Name: "UpdatedAt", Type: "time.Time", Tag: "`json:\"updated_at\" gorm:\"autoUpdateTime\"`"})
	}

	// Add soft delete if requested
	if softDelete {
		fieldsList = append(fieldsList, Field{Name: "DeletedAt", Type: "gorm.DeletedAt", Tag: "`json:\"deleted_at,omitempty\" gorm:\"index\"`"})
	}

	// Generate entity file with real field-based content
	generateEntityFile(domainDir, entityName, fieldsList, validation, businessRules, timestamps, softDelete)

	// Generate errors file if validation is enabled - now with real field validations
	if validation {
		generateErrorsFile(domainDir, entityName, fieldsList)
	}

	// Generate seed data automatically
	generateSeedData(domainDir, entityName, fieldsList)
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

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error writing entity file: %v\n", err)
		return
	}
}

func generateBusinessRules(content *strings.Builder, entityName string, fields []Field) {
	entityVar := strings.ToLower(string(entityName[0]))

	// Generate some common business rules based on field types
	for _, field := range fields {
		switch field.Name {
		case "Age":
			fmt.Fprintf(content, "func (%s *%s) IsAdult() bool {\n", entityVar, entityName)
			fmt.Fprintf(content, "\treturn %s.Age >= 18\n", entityVar)
			content.WriteString("}\n\n")

		case "Price":
			fmt.Fprintf(content, "func (%s *%s) IsExpensive() bool {\n", entityVar, entityName)
			fmt.Fprintf(content, "\treturn %s.Price > 1000.0\n", entityVar)
			content.WriteString("}\n\n")

		case "Email":
			fmt.Fprintf(content, "func (%s *%s) HasValidEmail() bool {\n", entityVar, entityName)
			fmt.Fprintf(content, "\treturn strings.Contains(%s.Email, \"@\")\n", entityVar)
			content.WriteString("}\n\n")

		case "Status":
			fmt.Fprintf(content, "func (%s *%s) IsActive() bool {\n", entityVar, entityName)
			fmt.Fprintf(content, "\treturn %s.Status == \"active\"\n", entityVar)
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
		fmt.Printf("⚠️  Archivo errors.go ya existe, agregando nuevos errores para %s\n", entityName)

		if existingContent, err := os.ReadFile(filename); err == nil {
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
	generalError := fmt.Sprintf("\tErrInvalid%sData = errors.New(\"datos de %s inválidos\")",
		entityName, strings.ToLower(entityName))
	if !contains(existingErrors, generalError) {
		content.WriteString(generalError + "\n")
	}

	// Add existing errors
	for _, err := range existingErrors {
		content.WriteString(err + "\n")
	}

	// Generate specific validation errors for all fields with Spanish messages
	for _, field := range fields {
		if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
			continue
		}

		// Generate field-specific errors based on type and name
		fieldLower := strings.ToLower(field.Name)

		// Basic required field error
		requiredError := fmt.Sprintf("\tErrInvalid%s%s = errors.New(\"%s es requerido\")",
			entityName, field.Name, getSpanishFieldName(fieldLower))
		if !contains(existingErrors, requiredError) {
			content.WriteString(requiredError + "\n")
		}

		// Type-specific errors
		switch field.Type {
		case "string":
			if strings.Contains(fieldLower, "email") {
				emailError := fmt.Sprintf("\tErrInvalid%s%sFormat = errors.New(\"formato de %s inválido\")",
					entityName, field.Name, getSpanishFieldName(fieldLower))
				if !contains(existingErrors, emailError) {
					content.WriteString(emailError + "\n")
				}
			}

			if strings.Contains(fieldLower, "name") {
				lengthError := fmt.Sprintf("\tErrInvalid%s%sLength = errors.New(\"%s debe tener entre 2 y 100 caracteres\")",
					entityName, field.Name, getSpanishFieldName(fieldLower))
				if !contains(existingErrors, lengthError) {
					content.WriteString(lengthError + "\n")
				}
			}

		case "int", "int64", "uint", "uint64":
			if strings.Contains(fieldLower, "age") {
				ageError := fmt.Sprintf("\tErrInvalid%s%sRange = errors.New(\"%s debe ser mayor a 0\")",
					entityName, field.Name, getSpanishFieldName(fieldLower))
				if !contains(existingErrors, ageError) {
					content.WriteString(ageError + "\n")
				}
			} else {
				rangeError := fmt.Sprintf("\tErrInvalid%s%sRange = errors.New(\"%s debe ser un número positivo\")",
					entityName, field.Name, getSpanishFieldName(fieldLower))
				if !contains(existingErrors, rangeError) {
					content.WriteString(rangeError + "\n")
				}
			}

		case "float64", "float32":
			if strings.Contains(fieldLower, "price") || strings.Contains(fieldLower, "amount") {
				priceError := fmt.Sprintf("\tErrInvalid%s%sRange = errors.New(\"%s debe ser mayor a 0 y menor a 999,999,999.99\")",
					entityName, field.Name, getSpanishFieldName(fieldLower))
				if !contains(existingErrors, priceError) {
					content.WriteString(priceError + "\n")
				}
			} else {
				numberError := fmt.Sprintf("\tErrInvalid%s%sRange = errors.New(\"%s debe ser un número positivo\")",
					entityName, field.Name, getSpanishFieldName(fieldLower))
				if !contains(existingErrors, numberError) {
					content.WriteString(numberError + "\n")
				}
			}
		}
	}

	content.WriteString(")\n")
	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error writing errors file: %v\n", err)
	}
}

// getSpanishFieldName converts common field names to Spanish for error messages
func getSpanishFieldName(fieldName string) string {
	fieldTranslations := map[string]string{
		"name":        "el nombre",
		"email":       "el email",
		"age":         "la edad",
		"price":       "el precio",
		"amount":      "el monto",
		"description": "la descripción",
		"title":       "el título",
		"status":      "el estado",
		"category":    "la categoría",
		"stock":       "el stock",
		"quantity":    "la cantidad",
		"phone":       "el teléfono",
		"address":     "la dirección",
		"password":    "la contraseña",
	}

	for key, value := range fieldTranslations {
		if strings.Contains(fieldName, key) {
			return value
		}
	}

	return "el campo " + fieldName
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

// generateSeedData creates seed data based on actual fields
func generateSeedData(dir, entityName string, fields []Field) {
	filename := filepath.Join(dir, strings.ToLower(entityName)+"_seeds.go")

	var content strings.Builder
	content.WriteString("package domain\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"time\"\n")
	content.WriteString(")\n\n")

	content.WriteString(fmt.Sprintf("// Get%sSeeds retorna datos de ejemplo para %s\n", entityName, strings.ToLower(entityName)))
	content.WriteString(fmt.Sprintf("func Get%sSeeds() []%s {\n", entityName, entityName))
	content.WriteString(fmt.Sprintf("\treturn []%s{\n", entityName))

	// Generate 3 sample records based on actual fields
	for i := 1; i <= 3; i++ {
		content.WriteString("\t\t{\n")
		for _, field := range fields {
			if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
				continue // Skip auto-managed fields
			}

			sampleValue := generateSampleValue(field, i)
			content.WriteString(fmt.Sprintf("\t\t\t%s: %s,\n", field.Name, sampleValue))
		}
		content.WriteString("\t\t},\n")
	}

	content.WriteString("\t}\n")
	content.WriteString("}\n\n")

	// Generate SQL INSERT statements as comments
	content.WriteString(fmt.Sprintf("// GetSQL%sSeeds retorna sentencias SQL INSERT para %s\n", entityName, strings.ToLower(entityName)))
	content.WriteString(fmt.Sprintf("func GetSQL%sSeeds() string {\n", entityName))
	content.WriteString(fmt.Sprintf("\treturn `-- Datos de ejemplo para tabla %s\n", strings.ToLower(entityName)))

	// Generate SQL INSERT statements
	for i := 1; i <= 3; i++ {
		content.WriteString(fmt.Sprintf("INSERT INTO %s (", strings.ToLower(entityName)+"s"))

		// Field names
		fieldNames := []string{}
		for _, field := range fields {
			if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
				continue
			}
			fieldNames = append(fieldNames, strings.ToLower(field.Name))
		}
		content.WriteString(strings.Join(fieldNames, ", "))
		content.WriteString(") VALUES (")

		// Field values
		values := []string{}
		for _, field := range fields {
			if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
				continue
			}
			sqlValue := generateSQLSampleValue(field, i)
			values = append(values, sqlValue)
		}
		content.WriteString(strings.Join(values, ", "))
		content.WriteString(");\\n")
	}

	content.WriteString("`\n")
	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error writing seed file: %v\n", err)
	}
}

// generateSampleValue creates realistic sample data based on field type and name
func generateSampleValue(field Field, index int) string {
	fieldLower := strings.ToLower(field.Name)

	switch field.Type {
	case "string":
		switch {
		case strings.Contains(fieldLower, "name"):
			names := []string{"Juan Pérez", "María García", "Carlos López"}
			return fmt.Sprintf("\"%s\"", names[(index-1)%len(names)])
		case strings.Contains(fieldLower, "email"):
			emails := []string{"juan@ejemplo.com", "maria@ejemplo.com", "carlos@ejemplo.com"}
			return fmt.Sprintf("\"%s\"", emails[(index-1)%len(emails)])
		case strings.Contains(fieldLower, "description"):
			descriptions := []string{"Descripción detallada del primer elemento", "Información completa del segundo item", "Detalles específicos del tercer registro"}
			return fmt.Sprintf("\"%s\"", descriptions[(index-1)%len(descriptions)])
		case strings.Contains(fieldLower, "title"):
			titles := []string{"Título Principal", "Elemento Secundario", "Item Terciario"}
			return fmt.Sprintf("\"%s\"", titles[(index-1)%len(titles)])
		case strings.Contains(fieldLower, "status"):
			statuses := []string{"activo", "pendiente", "completado"}
			return fmt.Sprintf("\"%s\"", statuses[(index-1)%len(statuses)])
		case strings.Contains(fieldLower, "category"):
			categories := []string{"tecnología", "educación", "salud"}
			return fmt.Sprintf("\"%s\"", categories[(index-1)%len(categories)])
		default:
			return fmt.Sprintf("\"Ejemplo %s %d\"", field.Name, index)
		}

	case "int", "int64", "uint", "uint64":
		switch {
		case strings.Contains(fieldLower, "age"):
			ages := []int{25, 30, 35}
			return fmt.Sprintf("%d", ages[(index-1)%len(ages)])
		case strings.Contains(fieldLower, "stock"):
			stocks := []int{100, 50, 75}
			return fmt.Sprintf("%d", stocks[(index-1)%len(stocks)])
		case strings.Contains(fieldLower, "quantity"):
			quantities := []int{10, 5, 15}
			return fmt.Sprintf("%d", quantities[(index-1)%len(quantities)])
		default:
			return fmt.Sprintf("%d", index*10)
		}

	case "float64", "float32":
		switch {
		case strings.Contains(fieldLower, "price"):
			prices := []float64{99.99, 149.50, 199.99}
			return fmt.Sprintf("%.2f", prices[(index-1)%len(prices)])
		case strings.Contains(fieldLower, "amount"):
			amounts := []float64{1000.00, 2500.50, 3750.75}
			return fmt.Sprintf("%.2f", amounts[(index-1)%len(amounts)])
		default:
			return fmt.Sprintf("%.2f", float64(index)*10.50)
		}

	case "bool":
		return fmt.Sprintf("%t", index%2 == 1)

	case "time.Time":
		return "time.Now()"

	default:
		// Generar valores por defecto inteligentes según el tipo
		switch {
		case strings.Contains(field.Type, "int"):
			return fmt.Sprintf("%d", index*10)
		case strings.Contains(field.Type, "string"):
			return fmt.Sprintf("\"Valor%d\"", index)
		case strings.Contains(field.Type, "float"):
			return fmt.Sprintf("%.2f", float64(index)*10.5)
		case strings.Contains(field.Type, "bool"):
			return fmt.Sprintf("%t", index%2 == 1)
		default:
			return "nil // Tipo personalizado"
		}
	}
}

// generateSQLSampleValue creates SQL-compatible sample values
func generateSQLSampleValue(field Field, index int) string {
	fieldLower := strings.ToLower(field.Name)

	switch field.Type {
	case "string":
		switch {
		case strings.Contains(fieldLower, "name"):
			names := []string{"Juan Pérez", "María García", "Carlos López"}
			return fmt.Sprintf("'%s'", names[(index-1)%len(names)])
		case strings.Contains(fieldLower, "email"):
			emails := []string{"juan@ejemplo.com", "maria@ejemplo.com", "carlos@ejemplo.com"}
			return fmt.Sprintf("'%s'", emails[(index-1)%len(emails)])
		case strings.Contains(fieldLower, "description"):
			descriptions := []string{"Descripción detallada del primer elemento", "Información completa del segundo item", "Detalles específicos del tercer registro"}
			return fmt.Sprintf("'%s'", descriptions[(index-1)%len(descriptions)])
		case strings.Contains(fieldLower, "status"):
			statuses := []string{"activo", "pendiente", "completado"}
			return fmt.Sprintf("'%s'", statuses[(index-1)%len(statuses)])
		default:
			return fmt.Sprintf("'Ejemplo %s %d'", field.Name, index)
		}

	case "int", "int64", "uint", "uint64":
		switch {
		case strings.Contains(fieldLower, "age"):
			ages := []int{25, 30, 35}
			return fmt.Sprintf("%d", ages[(index-1)%len(ages)])
		case strings.Contains(fieldLower, "stock"):
			stocks := []int{100, 50, 75}
			return fmt.Sprintf("%d", stocks[(index-1)%len(stocks)])
		default:
			return fmt.Sprintf("%d", index*10)
		}

	case "float64", "float32":
		switch {
		case strings.Contains(fieldLower, "price"):
			prices := []float64{99.99, 149.50, 199.99}
			return fmt.Sprintf("%.2f", prices[(index-1)%len(prices)])
		default:
			return fmt.Sprintf("%.2f", float64(index)*10.50)
		}

	case "bool":
		return fmt.Sprintf("%t", index%2 == 1)

	case "time.Time":
		return "NOW()"

	default:
		return "NULL"
	}
}

func init() {
	rootCmd.AddCommand(entityCmd)
	entityCmd.Flags().StringP("fields", "f", "", "Campos de la entidad \"campo:tipo,campo2:tipo\" (requerido)")
	entityCmd.Flags().BoolP("validation", "v", false, "Incluir validaciones de negocio")
	entityCmd.Flags().BoolP("business-rules", "b", false, "Incluir reglas de negocio avanzadas")
	entityCmd.Flags().BoolP("timestamps", "t", false, "Incluir campos CreatedAt y UpdatedAt")
	entityCmd.Flags().BoolP("soft-delete", "s", false, "Incluir eliminación suave (DeletedAt)")
	_ = entityCmd.MarkFlagRequired("fields")
}
