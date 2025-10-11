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

		// Initialize configuration integration
		configIntegration := NewConfigIntegration()
		if err := configIntegration.LoadConfigForProject(); err != nil {
			fmt.Printf("⚠️  Warning: Could not load configuration: %v\n", err)
			fmt.Println("📝 Using default values. Consider running 'goca init --config' to generate .goca.yaml")
		}

		// Merge CLI flags with configuration (CLI flags take precedence)
		// Only include flags that were explicitly set by the user
		flags := map[string]interface{}{}
		if cmd.Flags().Changed("validation") {
			flags["validation"] = validation
		}
		if cmd.Flags().Changed("business-rules") {
			flags["business-rules"] = businessRules
		}
		if cmd.Flags().Changed("timestamps") {
			flags["timestamps"] = timestamps
		}
		if cmd.Flags().Changed("soft-delete") {
			flags["soft-delete"] = softDelete
		}
		if len(flags) > 0 {
			configIntegration.MergeWithCLIFlags(flags)
		}

		// Get effective values from configuration
		// Use CLI flag if explicitly set, otherwise use config
		effectiveValidation := validation
		if !cmd.Flags().Changed("validation") && configIntegration.config != nil {
			effectiveValidation = configIntegration.config.Generation.Validation.Enabled
		}

		effectiveBusinessRules := businessRules
		if !cmd.Flags().Changed("business-rules") && configIntegration.config != nil {
			effectiveBusinessRules = configIntegration.config.Generation.BusinessRules.Enabled
		}

		// Get timestamps and soft delete from config if not explicitly set via CLI
		effectiveTimestamps := timestamps
		effectiveSoftDelete := softDelete
		if !cmd.Flags().Changed("timestamps") && configIntegration.config != nil {
			effectiveTimestamps = configIntegration.config.Database.Features.Timestamps
		}
		if !cmd.Flags().Changed("soft-delete") && configIntegration.config != nil {
			effectiveSoftDelete = configIntegration.config.Database.Features.SoftDelete
		} // Usar validador centralizado
		validator := NewCommandValidator()

		if err := validator.ValidateEntityCommand(entityName, fields); err != nil {
			validator.errorHandler.HandleError(err, "validación de parámetros")
		}

		validator.errorHandler.ValidateRequiredFlag(fields, "fields")

		fmt.Printf("🏗️  Generating entity '%s'\n", entityName)
		fmt.Printf("📋 Fields: %s\n", fields)

		if effectiveValidation {
			fmt.Print("✅ Including validations")
			if configIntegration.HasConfigFile() {
				fmt.Printf(" (from config)")
			}
			fmt.Println()
		}
		if effectiveBusinessRules {
			fmt.Print("🧠 Including business rules")
			if configIntegration.HasConfigFile() {
				fmt.Printf(" (from config)")
			}
			fmt.Println()
		}
		if effectiveTimestamps {
			fmt.Print("⏰ Including timestamps")
			if configIntegration.HasConfigFile() && !cmd.Flags().Changed("timestamps") {
				fmt.Printf(" (from config)")
			}
			fmt.Println()
		}
		if effectiveSoftDelete {
			fmt.Print("🗑️  Including soft delete")
			if configIntegration.HasConfigFile() && !cmd.Flags().Changed("soft-delete") {
				fmt.Printf(" (from config)")
			}
			fmt.Println()
		}

		if configIntegration.HasConfigFile() {
			configIntegration.PrintConfigSummary()
		}

		// Get naming convention for files
		fileNamingConvention := "lowercase" // default
		if configIntegration.config != nil {
			fileNamingConvention = configIntegration.GetNamingConvention("file")
		}

		generateEntity(entityName, fields, effectiveValidation, effectiveBusinessRules, effectiveTimestamps, effectiveSoftDelete, fileNamingConvention)

		// Generar datos de semilla automáticamente
		if fields != "" {
			generateSeedData("internal/domain", entityName, parseFields(fields))
			fmt.Println("🌱 Seed data generated")
		}

		fmt.Printf("\n✅ Entity '%s' generated successfully!\n", entityName)
		fmt.Printf("📁 Files created:\n")
		fmt.Printf("   - internal/domain/%s.go\n", strings.ToLower(entityName))
		if effectiveValidation {
			fmt.Printf("   - internal/domain/errors.go\n")
		}
		fmt.Printf("   - internal/domain/%s_seeds.go\n", strings.ToLower(entityName))
		fmt.Println("\n🎉 All set! Your entity is ready to use.")
	},
}

func generateEntity(entityName, fields string, validation, businessRules, timestamps, softDelete bool, fileNamingConvention string) {
	// Crear directorio domain si no existe
	domainDir := "internal/domain"
	_ = os.MkdirAll(domainDir, 0755)

	// Parse fields - ahora genera campos reales basados en el input
	// Note: ParseFieldsWithValidation already adds the ID field
	fieldsList := parseFieldsWithValidation(fields, validation)

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
	generateEntityFile(domainDir, entityName, fieldsList, validation, businessRules, timestamps, softDelete, fileNamingConvention)

	// Generate errors file if validation is enabled - now with real field validations
	if validation {
		generateErrorsFile(domainDir, entityName, fieldsList)
	}

	// Generate seed data automatically
	generateSeedData(domainDir, entityName, fieldsList)
}

// Field represents a single field definition for an entity structure.
type Field struct {
	Name string
	Type string
	Tag  string
}

func parseFields(fields string) []Field {
	return parseFieldsWithValidation(fields, false)
}

func parseFieldsWithValidation(fields string, withValidation bool) []Field {
	validator := NewFieldValidator()
	fieldsList, err := validator.ParseFieldsWithValidation(fields)
	if err != nil {
		fmt.Printf("❌ Error in field validation: %v\n", err)
		os.Exit(1)
	}

	// If validation is enabled, add validate tags to the field tags
	if withValidation {
		for i := range fieldsList {
			if fieldsList[i].Name != "ID" {
				// Parse existing tag and add validation tag
				existingTag := fieldsList[i].Tag
				// Remove backticks
				existingTag = strings.Trim(existingTag, "`")

				// Add validation tag based on field type
				validateTag := getValidateTag(fieldsList[i].Name, fieldsList[i].Type)
				if validateTag != "" {
					existingTag += fmt.Sprintf(" validate:\"%s\"", validateTag)
				}

				fieldsList[i].Tag = "`" + existingTag + "`"
			}
		}
	}

	return fieldsList
}

func getValidateTag(fieldName, fieldType string) string {
	switch fieldType {
	case FieldString:
		if fieldName == "Email" || strings.ToLower(fieldName) == "email" {
			return "required,email"
		}
		return "required"
	case "int", "int64", "uint", "uint64":
		return "required,gte=0"
	case "float64":
		return "required,gte=0"
	case "bool":
		return "" // Booleans don't usually need validation
	default:
		return "required"
	}
}

func getGormTag(fieldName, fieldType string) string {
	switch fieldType {
	case FieldString:
		if fieldName == FieldEmailType {
			return "type:varchar(255);uniqueIndex;not null"
		}
		if fieldName == "Title" || fieldName == "Name" {
			return "type:varchar(255);not null"
		}
		if fieldName == "Description" {
			return "type:text"
		}
		return "type:varchar(255)"
	case FieldInt:
		return "type:integer;not null;default:0"
	case FieldBool:
		return "type:boolean;not null;default:false"
	case FieldFloat64:
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

func generateEntityFile(dir, entityName string, fields []Field, validation, businessRules, timestamps, softDelete bool, fileNamingConvention string) {
	// Apply naming convention to filename
	var filename string
	if fileNamingConvention == "snake_case" {
		filename = filepath.Join(dir, toSnakeCase(entityName)+".go")
	} else if fileNamingConvention == "kebab-case" {
		filename = filepath.Join(dir, toKebabCase(entityName)+".go")
	} else {
		// Default to lowercase
		filename = filepath.Join(dir, strings.ToLower(entityName)+".go")
	}

	var content strings.Builder

	writeEntityHeader(&content, fields, businessRules, timestamps, softDelete)
	writeEntityStruct(&content, entityName, fields)

	if validation {
		writeValidationMethod(&content, entityName, fields)
	}

	if businessRules {
		generateBusinessRules(&content, entityName, fields)
	}

	if softDelete {
		writeSoftDeleteMethods(&content, entityName)
	}

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error writing entity file: %v\n", err)
		return
	}
}

// writeEntityHeader writes package declaration and imports
func writeEntityHeader(content *strings.Builder, fields []Field, businessRules, timestamps, softDelete bool) {
	content.WriteString("package domain\n\n")

	// Check if any field is time.Time
	hasTimeField := false
	for _, field := range fields {
		if field.Type == "time.Time" {
			hasTimeField = true
			break
		}
	}

	needsTime := timestamps || softDelete || hasTimeField
	needsStrings := businessRules && hasStringBusinessRules(fields)
	needsGorm := softDelete // Need gorm.io/gorm for gorm.DeletedAt

	if needsTime || needsStrings || needsGorm {
		content.WriteString("import (\n")
		if needsStrings {
			content.WriteString("\t\"strings\"\n")
		}
		if needsTime {
			content.WriteString("\t\"time\"\n")
		}
		if needsGorm {
			content.WriteString("\n\t\"gorm.io/gorm\"\n")
		}
		content.WriteString(")\n\n")
	}
}

// writeEntityStruct writes the entity struct definition
func writeEntityStruct(content *strings.Builder, entityName string, fields []Field) {
	fmt.Fprintf(content, "type %s struct {\n", entityName)
	for _, field := range fields {
		fmt.Fprintf(content, "\t%s %s %s\n", field.Name, field.Type, field.Tag)
	}
	content.WriteString("}\n\n")
}

// writeValidationMethod writes the Validate method for the entity
func writeValidationMethod(content *strings.Builder, entityName string, fields []Field) {
	entityVar := strings.ToLower(string(entityName[0]))
	fmt.Fprintf(content, "func (%s *%s) Validate() error {\n", entityVar, entityName)

	for _, field := range fields {
		if isSystemField(field.Name) {
			continue
		}

		writeFieldValidation(content, entityVar, entityName, field)
	}

	content.WriteString("\treturn nil\n")
	content.WriteString("}\n\n")
}

// writeFieldValidation writes validation logic for a specific field
func writeFieldValidation(content *strings.Builder, entityVar, entityName string, field Field) {
	switch field.Type {
	case FieldString:
		fmt.Fprintf(content, "\tif %s.%s == \"\" {\n", entityVar, field.Name)
		fmt.Fprintf(content, "\t\treturn ErrInvalid%s%s\n", entityName, field.Name)
		content.WriteString("\t}\n")
	case "int", "int64", "float64":
		fmt.Fprintf(content, "\tif %s.%s < 0 {\n", entityVar, field.Name)
		fmt.Fprintf(content, "\t\treturn ErrInvalid%s%s\n", entityName, field.Name)
		content.WriteString("\t}\n")
	}
}

// writeSoftDeleteMethods writes soft delete helper methods
func writeSoftDeleteMethods(content *strings.Builder, entityName string) {
	entityVar := strings.ToLower(string(entityName[0]))

	fmt.Fprintf(content, "func (%s *%s) SoftDelete() {\n", entityVar, entityName)
	content.WriteString("\tnow := time.Now()\n")
	fmt.Fprintf(content, "\t%s.DeletedAt = &now\n", entityVar)
	content.WriteString("}\n\n")

	fmt.Fprintf(content, "func (%s *%s) IsDeleted() bool {\n", entityVar, entityName)
	fmt.Fprintf(content, "\treturn %s.DeletedAt != nil\n", entityVar)
	content.WriteString("}\n\n")
}

// isSystemField checks if a field is a system-managed field
func isSystemField(fieldName string) bool {
	systemFields := []string{"ID", StringCreatedAt, "UpdatedAt", "DeletedAt"}
	for _, sf := range systemFields {
		if fieldName == sf {
			return true
		}
	}
	return false
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
	existingErrors := readExistingErrors(filename, entityName)

	var content strings.Builder
	writeErrorsHeader(&content)
	writeEntityErrors(&content, entityName, fields, existingErrors)

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error writing errors file: %v\n", err)
	}
}

// readExistingErrors reads existing error definitions from the file
func readExistingErrors(filename, entityName string) []string {
	var existingErrors []string

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
	return existingErrors
}

// writeErrorsHeader writes the package declaration and imports
func writeErrorsHeader(content *strings.Builder) {
	content.WriteString("package domain\n\n")
	content.WriteString("import \"errors\"\n\n")
	content.WriteString("var (\n")
}

// writeEntityErrors writes all error definitions for the entity
func writeEntityErrors(content *strings.Builder, entityName string, fields []Field, existingErrors []string) {
	writeGeneralError(content, entityName, existingErrors)
	writeExistingErrors(content, existingErrors)
	writeFieldErrors(content, entityName, fields, existingErrors)
	content.WriteString(")\n")
}

// writeGeneralError writes the general entity error
func writeGeneralError(content *strings.Builder, entityName string, existingErrors []string) {
	generalError := fmt.Sprintf("\tErrInvalid%sData = errors.New(\"datos de %s inválidos\")",
		entityName, strings.ToLower(entityName))
	if !contains(existingErrors, generalError) {
		content.WriteString(generalError + "\n")
	}
}

// writeExistingErrors writes previously defined errors
func writeExistingErrors(content *strings.Builder, existingErrors []string) {
	for _, err := range existingErrors {
		content.WriteString(err + "\n")
	}
}

// writeFieldErrors writes validation errors for all fields
func writeFieldErrors(content *strings.Builder, entityName string, fields []Field, existingErrors []string) {
	for _, field := range fields {
		if isSystemField(field.Name) {
			continue
		}

		writeRequiredFieldError(content, entityName, field, existingErrors)
		writeTypeSpecificErrors(content, entityName, field, existingErrors)
	}
}

// writeRequiredFieldError writes the required field error
func writeRequiredFieldError(content *strings.Builder, entityName string, field Field, existingErrors []string) {
	fieldLower := strings.ToLower(field.Name)
	requiredError := fmt.Sprintf("\tErrInvalid%s%s = errors.New(\"%s es requerido\")",
		entityName, field.Name, getSpanishFieldName(fieldLower))
	if !contains(existingErrors, requiredError) {
		content.WriteString(requiredError + "\n")
	}
}

// writeTypeSpecificErrors writes type-specific validation errors
func writeTypeSpecificErrors(content *strings.Builder, entityName string, field Field, existingErrors []string) {
	fieldLower := strings.ToLower(field.Name)

	switch field.Type {
	case FieldString:
		writeStringFieldErrors(content, entityName, field, fieldLower, existingErrors)
	case "int", "int64", "uint", "uint64":
		writeIntegerFieldErrors(content, entityName, field, fieldLower, existingErrors)
	case "float64", "float32":
		writeFloatFieldErrors(content, entityName, field, fieldLower, existingErrors)
	}
}

// writeStringFieldErrors writes string-specific validation errors
func writeStringFieldErrors(content *strings.Builder, entityName string, field Field, fieldLower string, existingErrors []string) {
	if strings.Contains(fieldLower, FieldEmailType) {
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
}

// writeIntegerFieldErrors writes integer-specific validation errors
func writeIntegerFieldErrors(content *strings.Builder, entityName string, field Field, fieldLower string, existingErrors []string) {
	var rangeError string
	if strings.Contains(fieldLower, "age") {
		rangeError = fmt.Sprintf("\tErrInvalid%s%sRange = errors.New(\"%s debe ser mayor a 0\")",
			entityName, field.Name, getSpanishFieldName(fieldLower))
	} else {
		rangeError = fmt.Sprintf("\tErrInvalid%s%sRange = errors.New(\"%s debe ser un número positivo\")",
			entityName, field.Name, getSpanishFieldName(fieldLower))
	}
	if !contains(existingErrors, rangeError) {
		content.WriteString(rangeError + "\n")
	}
}

// writeFloatFieldErrors writes float-specific validation errors
func writeFloatFieldErrors(content *strings.Builder, entityName string, field Field, fieldLower string, existingErrors []string) {
	var rangeError string
	if strings.Contains(fieldLower, "price") || strings.Contains(fieldLower, "amount") {
		rangeError = fmt.Sprintf("\tErrInvalid%s%sRange = errors.New(\"%s debe ser mayor a 0 y menor a 999,999,999.99\")",
			entityName, field.Name, getSpanishFieldName(fieldLower))
	} else {
		rangeError = fmt.Sprintf("\tErrInvalid%s%sRange = errors.New(\"%s debe ser un número positivo\")",
			entityName, field.Name, getSpanishFieldName(fieldLower))
	}
	if !contains(existingErrors, rangeError) {
		content.WriteString(rangeError + "\n")
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
	writeSeedFileHeader(&content, fields)
	writeGoSeeds(&content, entityName, fields)
	writeSQLSeeds(&content, entityName, fields)

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error writing seed file: %v\n", err)
	}
}

// writeSeedFileHeader writes the package declaration and imports for seed file
func writeSeedFileHeader(content *strings.Builder, fields []Field) {
	content.WriteString("package domain\n\n")
	
	// Check if any field is time.Time
	hasTimeField := false
	for _, field := range fields {
		if field.Type == "time.Time" {
			hasTimeField = true
			break
		}
	}
	
	if hasTimeField {
		content.WriteString("import \"time\"\n\n")
	}
}

// writeGoSeeds writes the Go struct seed data function
func writeGoSeeds(content *strings.Builder, entityName string, fields []Field) {
	fmt.Fprintf(content, "// Get%sSeeds retorna datos de ejemplo para %s\n", entityName, strings.ToLower(entityName))
	fmt.Fprintf(content, "func Get%sSeeds() []%s {\n", entityName, entityName)
	fmt.Fprintf(content, "\treturn []%s{\n", entityName)

	// Generate 3 sample records based on actual fields
	for i := 1; i <= 3; i++ {
		writeGoSeedRecord(content, fields, i)
	}

	content.WriteString("\t}\n")
	content.WriteString("}\n\n")
}

// writeGoSeedRecord writes a single Go seed record
func writeGoSeedRecord(content *strings.Builder, fields []Field, recordNum int) {
	content.WriteString("\t\t{\n")
	for _, field := range fields {
		if isSystemField(field.Name) {
			continue // Skip auto-managed fields
		}

		sampleValue := generateSampleValue(field, recordNum)
		fmt.Fprintf(content, "\t\t\t%s: %s,\n", field.Name, sampleValue)
	}
	content.WriteString("\t\t},\n")
}

// writeSQLSeeds writes the SQL INSERT seed data function
func writeSQLSeeds(content *strings.Builder, entityName string, fields []Field) {
	fmt.Fprintf(content, "// GetSQL%sSeeds retorna sentencias SQL INSERT para %s\n", entityName, strings.ToLower(entityName))
	fmt.Fprintf(content, "func GetSQL%sSeeds() string {\n", entityName)
	fmt.Fprintf(content, "\treturn `-- Datos de ejemplo para tabla %s\n", strings.ToLower(entityName))

	// Generate SQL INSERT statements
	for i := 1; i <= 3; i++ {
		writeSQLInsertStatement(content, entityName, fields, i)
	}

	content.WriteString("`\n")
	content.WriteString("}\n")
}

// writeSQLInsertStatement writes a single SQL INSERT statement
func writeSQLInsertStatement(content *strings.Builder, entityName string, fields []Field, recordNum int) {
	fmt.Fprintf(content, "INSERT INTO %s (", strings.ToLower(entityName)+"s")

	// Field names
	fieldNames := getNonSystemFieldNames(fields)
	content.WriteString(strings.Join(fieldNames, ", "))
	content.WriteString(") VALUES (")

	// Field values
	values := getSQLFieldValues(fields, recordNum)
	content.WriteString(strings.Join(values, ", "))
	content.WriteString(");\\n")
}

// getNonSystemFieldNames returns field names excluding system fields
func getNonSystemFieldNames(fields []Field) []string {
	var fieldNames []string
	for _, field := range fields {
		if !isSystemField(field.Name) {
			fieldNames = append(fieldNames, strings.ToLower(field.Name))
		}
	}
	return fieldNames
}

// getSQLFieldValues returns SQL-formatted field values
func getSQLFieldValues(fields []Field, recordNum int) []string {
	var values []string
	for _, field := range fields {
		if !isSystemField(field.Name) {
			sqlValue := generateSQLSampleValue(field, recordNum)
			values = append(values, sqlValue)
		}
	}
	return values
}

// generateSampleValue creates realistic sample data based on field type and name
func generateSampleValue(field Field, index int) string {
	switch field.Type {
	case FieldString:
		return generateStringSampleValue(field.Name, index)
	case "int", "int64", "uint", "uint64":
		return generateIntSampleValue(field.Name, index)
	case "float64", "float32":
		return generateFloatSampleValue(field.Name, index)
	case FieldBool:
		return fmt.Sprintf("%t", index%2 == 1)
	case "time.Time":
		return "time.Now()"
	default:
		return generateDefaultSampleValue(field.Type, index)
	}
}

// generateStringSampleValue generates string sample values based on field name
func generateStringSampleValue(fieldName string, index int) string {
	fieldLower := strings.ToLower(fieldName)

	switch {
	case strings.Contains(fieldLower, "name"):
		names := []string{"Juan Pérez", "María García", "Carlos López"}
		return fmt.Sprintf("\"%s\"", names[(index-1)%len(names)])
	case strings.Contains(fieldLower, FieldEmailType):
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
		return fmt.Sprintf("\"Ejemplo %s %d\"", fieldName, index)
	}
}

// generateIntSampleValue generates integer sample values based on field name
func generateIntSampleValue(fieldName string, index int) string {
	fieldLower := strings.ToLower(fieldName)

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
}

// generateFloatSampleValue generates float sample values based on field name
func generateFloatSampleValue(fieldName string, index int) string {
	fieldLower := strings.ToLower(fieldName)

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
}

// generateDefaultSampleValue generates default sample values for unknown types
func generateDefaultSampleValue(fieldType string, index int) string {
	switch {
	case strings.Contains(fieldType, FieldInt):
		return fmt.Sprintf("%d", index*10)
	case strings.Contains(fieldType, FieldString):
		return fmt.Sprintf("\"Valor%d\"", index)
	case strings.Contains(fieldType, "float"):
		return fmt.Sprintf("%.2f", float64(index)*10.5)
	case strings.Contains(fieldType, FieldBool):
		return fmt.Sprintf("%t", index%2 == 1)
	default:
		return "nil // Tipo personalizado"
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
	entityCmd.Flags().StringP("fields", "f", "", "Entity fields \"field:type,field2:type\" (required)")
	entityCmd.Flags().BoolP("validation", "v", false, "Include business validations")
	entityCmd.Flags().BoolP("business-rules", "b", false, "Include advanced business rules")
	entityCmd.Flags().BoolP("timestamps", "t", false, "Include CreatedAt and UpdatedAt fields")
	entityCmd.Flags().BoolP("soft-delete", "s", false, "Include soft delete (DeletedAt)")
	_ = entityCmd.MarkFlagRequired("fields")
}
