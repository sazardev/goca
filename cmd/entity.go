package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
		tests, _ := cmd.Flags().GetBool("tests")

		// Initialize configuration integration
		configIntegration := NewConfigIntegration()
		if err := configIntegration.LoadConfigForProject(); err != nil {
			ui.Warning(fmt.Sprintf("Could not load configuration: %v", err))
			ui.Dim("Using default values. Consider running 'goca init --config' to generate .goca.yaml")
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
			validator.errorHandler.HandleError(err, "parameter validation")
		}

		validator.errorHandler.ValidateRequiredFlag(fields, "fields")

		ui.Header(fmt.Sprintf("Generating entity '%s'", entityName))
		ui.KeyValue("Fields", fields)

		if effectiveValidation {
			ui.Feature("Including validations", configIntegration.HasConfigFile())
		}
		if effectiveBusinessRules {
			ui.Feature("Including business rules", configIntegration.HasConfigFile())
		}
		if effectiveTimestamps {
			ui.Feature("Including timestamps", configIntegration.HasConfigFile() && !cmd.Flags().Changed("timestamps"))
		}
		if effectiveSoftDelete {
			ui.Feature("Including soft delete", configIntegration.HasConfigFile() && !cmd.Flags().Changed("soft-delete"))
		}

		if configIntegration.HasConfigFile() {
			configIntegration.PrintConfigSummary()
		}

		// Get naming convention for files
		fileNamingConvention := "lowercase" // default
		if configIntegration.config != nil {
			fileNamingConvention = configIntegration.GetNamingConvention("file")
		}

		// Initialize safety manager
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		backup, _ := cmd.Flags().GetBool("backup")
		sm := NewSafetyManager(dryRun, force, backup)

		if dryRun {
			ui.DryRun("Previewing changes without creating files")
		}

		if err := generateEntity(entityName, fields, effectiveValidation, effectiveBusinessRules, effectiveTimestamps, effectiveSoftDelete, tests, fileNamingConvention, sm); err != nil {
			os.Exit(1)
		}

		if dryRun {
			sm.PrintSummary()
			return
		}

		ui.Success(fmt.Sprintf("Entity '%s' generated successfully!", entityName))

		rows := [][]string{
			{fmt.Sprintf("internal/domain/%s.go", strings.ToLower(entityName)), "Entity"},
		}
		if effectiveValidation {
			rows = append(rows, []string{"internal/domain/errors.go", "Domain errors"})
		}
		rows = append(rows, []string{fmt.Sprintf("internal/domain/%s_seeds.go", strings.ToLower(entityName)), "Seed data"})
		if tests {
			rows = append(rows, []string{fmt.Sprintf("internal/domain/%s_test.go", strings.ToLower(entityName)), "Unit tests"})
		}
		ui.Table([]string{"File", "Description"}, rows)
		ui.Blank()
		ui.Success("All set! Your entity is ready to use.")
	},
}

func generateEntity(entityName, fields string, validation, businessRules, timestamps, softDelete, tests bool, fileNamingConvention string, sm ...*SafetyManager) error {
	// Create domain directory if it doesn't exist
	domainDir := "internal/domain"
	_ = os.MkdirAll(domainDir, 0o755)

	// Parse fields - generates real fields based on the input
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

	// Generate entity file with real field-based content. This is the primary
	// artifact: if it cannot be written, abort without performing partial side
	// effects (errors/seeds/tests) so the caller can fail cleanly.
	if err := generateEntityFile(domainDir, entityName, fieldsList, validation, businessRules, timestamps, softDelete, fileNamingConvention, sm...); err != nil {
		return err
	}

	// Generate errors file if validation is enabled - now with real field validations
	if validation {
		generateErrorsFile(domainDir, entityName, fieldsList, sm...)
	}

	// Generate seed data automatically
	generateSeedData(domainDir, entityName, fieldsList, sm...)

	// Generate unit tests if requested
	if tests {
		generateEntityTests(domainDir, entityName, fieldsList, validation, businessRules, fileNamingConvention, sm...)
	}

	return nil
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
		ui.Error(fmt.Sprintf("Error in field validation: %v", err))
		os.Exit(1)
	}

	// Normalize field names to idiomatic Go PascalCase (handling snake_case,
	// kebab-case and common initialisms like ID/URL/API). The struct tag keeps
	// the snake_case form for json/gorm. The ID field is left untouched.
	for i := range fieldsList {
		if fieldsList[i].Name == "ID" {
			continue
		}
		snake := strings.ToLower(strings.ReplaceAll(fieldsList[i].Name, "-", "_"))
		fieldsList[i].Name = toGoFieldName(fieldsList[i].Name)
		fieldsList[i].Tag = rebuildFieldTag(fieldsList[i].Tag, snake)
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

// commonInitialisms maps lowercase word fragments to their idiomatic Go
// capitalization so field names like user_id -> UserID and api_url -> APIURL.
var commonInitialisms = map[string]string{
	"id": "ID", "url": "URL", "uri": "URI", "api": "API", "http": "HTTP",
	"https": "HTTPS", "json": "JSON", "xml": "XML", "sql": "SQL", "uuid": "UUID",
	"html": "HTML", "ip": "IP", "ssh": "SSH", "tcp": "TCP", "udp": "UDP",
	"db": "DB", "ui": "UI", "ttl": "TTL", "ascii": "ASCII", "cpu": "CPU",
}

// toGoFieldName converts an arbitrary field name (snake_case, kebab-case,
// camelCase or already PascalCase) into idiomatic Go PascalCase, honoring
// common initialisms. Examples: last_login -> LastLogin, user_id -> UserID.
func toGoFieldName(name string) string {
	// Split on separators; also split camelCase boundaries.
	var words []string
	var cur strings.Builder
	flush := func() {
		if cur.Len() > 0 {
			words = append(words, cur.String())
			cur.Reset()
		}
	}
	runes := []rune(name)
	for i, r := range runes {
		if r == '_' || r == '-' || r == ' ' {
			flush()
			continue
		}
		// Split on lower->upper boundary (camelCase) to re-segment words.
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := runes[i-1]
			if prev >= 'a' && prev <= 'z' {
				flush()
			}
		}
		cur.WriteRune(r)
	}
	flush()

	var result strings.Builder
	for _, w := range words {
		lower := strings.ToLower(w)
		if init, ok := commonInitialisms[lower]; ok {
			result.WriteString(init)
			continue
		}
		result.WriteString(strings.ToUpper(lower[:1]) + lower[1:])
	}
	return result.String()
}

// rebuildFieldTag rewrites the json and gorm tag keys of an existing struct tag
// so they use the supplied snake_case field name, preserving any additional
// gorm options (e.g. uniqueIndex;not null) that were generated from the field.
func rebuildFieldTag(tag, snake string) string {
	inner := strings.Trim(tag, "`")
	// Replace json:"..." value with the snake_case name.
	inner = replaceTagJSONName(inner, snake)
	return "`" + inner + "`"
}

// replaceTagJSONName replaces the value of the json tag key with the given name,
// keeping any options such as ",omitempty".
func replaceTagJSONName(inner, snake string) string {
	const key = `json:"`
	idx := strings.Index(inner, key)
	if idx < 0 {
		return inner
	}
	start := idx + len(key)
	end := strings.Index(inner[start:], `"`)
	if end < 0 {
		return inner
	}
	end += start
	val := inner[start:end]
	// Preserve options after the first comma (e.g. ",omitempty").
	opts := ""
	if c := strings.Index(val, ","); c >= 0 {
		opts = val[c:]
	}
	return inner[:start] + snake + opts + inner[end:]
}

func getValidateTag(fieldName, fieldType string) string {
	switch {
	case fieldType == FieldString:
		if fieldName == "Email" || strings.EqualFold(fieldName, "email") {
			return "required,email"
		}
		return "required"
	case isSignedNumericType(fieldType):
		return "required,gte=0"
	case isUnsignedIntType(fieldType):
		// Unsigned integers are always >= 0; gte=0 would be redundant.
		return "required"
	case fieldType == "bool":
		return "" // Booleans don't usually need validation
	case isSliceType(fieldType) || isPointerType(fieldType):
		// required on a slice/pointer is dubious and the runtime Validate() body
		// has no coherent check for these, so emit no validate tag (see ENTITY-9).
		return ""
	default:
		return "required"
	}
}

// isNumericType reports whether fieldType is any Go numeric subtype.
func isNumericType(fieldType string) bool {
	return isSignedNumericType(fieldType) || isUnsignedIntType(fieldType)
}

// isSignedNumericType reports whether fieldType is a signed integer (incl. rune)
// or a float, i.e. a type for which a "< 0" runtime check is meaningful.
func isSignedNumericType(fieldType string) bool {
	switch fieldType {
	case "int", "int8", "int16", "int32", "int64", "rune", "float32", "float64":
		return true
	}
	return false
}

// isUnsignedIntType reports whether fieldType is an unsigned integer subtype.
func isUnsignedIntType(fieldType string) bool {
	switch fieldType {
	case "uint", "uint8", "uint16", "uint32", "uint64", "byte":
		return true
	}
	return false
}

// isSliceType reports whether fieldType is a slice type.
func isSliceType(fieldType string) bool {
	return strings.HasPrefix(fieldType, "[]")
}

// isPointerType reports whether fieldType is a pointer type.
func isPointerType(fieldType string) bool {
	return strings.HasPrefix(fieldType, "*")
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

// hasStringBusinessRules checks if any field will require the strings package for business rules.
func hasStringBusinessRules(fields []Field) bool {
	for _, field := range fields {
		if field.Name == "Email" {
			return true
		}
	}
	return false
}

func generateEntityFile(dir, entityName string, fields []Field, validation, businessRules, timestamps, softDelete bool, fileNamingConvention string, sm ...*SafetyManager) error {
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
	// Emit stub definitions for unknown custom/named types referenced by fields
	// (e.g. status:UserStatus) so the generated package compiles (ENTITY-1).
	writeCustomTypeStubs(&content, entityName, fields)

	if validation {
		writeValidationMethod(&content, entityName, fields)
	}

	if businessRules {
		generateBusinessRules(&content, entityName, fields)
	}

	if softDelete {
		writeSoftDeleteMethods(&content, entityName)
	}

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing entity file: %v", err))
		return err
	}
	return nil
}

// writeEntityHeader writes package declaration and imports.
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
	needsStrings := (businessRules && hasStringBusinessRules(fields)) || hasEmailField(fields)
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

// writeEntityStruct writes the entity struct definition.
func writeEntityStruct(content *strings.Builder, entityName string, fields []Field) {
	fmt.Fprintf(content, "type %s struct {\n", entityName)
	for _, field := range fields {
		fmt.Fprintf(content, "\t%s %s %s\n", field.Name, field.Type, field.Tag)
	}
	content.WriteString("}\n\n")
}

// writeCustomTypeStubs emits a stub "type X string" definition for each unknown
// custom/named type referenced by a field, so the generated domain package
// compiles even when the user passes a type that does not yet exist.
func writeCustomTypeStubs(content *strings.Builder, entityName string, fields []Field) {
	seen := map[string]bool{}
	var stubs []string
	for _, field := range fields {
		base := customTypeBase(field.Type)
		if base == "" || base == entityName || seen[base] {
			continue
		}
		seen[base] = true
		stubs = append(stubs, base)
	}
	for _, t := range stubs {
		fmt.Fprintf(content, "// %s is a generated stub type. Replace with your own definition.\n", t)
		fmt.Fprintf(content, "type %s string\n\n", t)
	}
}

// customTypeBase returns the underlying unqualified custom type name referenced
// by fieldType (unwrapping leading []/*/[N]), or "" when the base is a builtin,
// qualified (package.Type), composite or otherwise not a stubbable custom type.
func customTypeBase(fieldType string) string {
	t := fieldType
	for {
		switch {
		case strings.HasPrefix(t, "[]"):
			t = strings.TrimPrefix(t, "[]")
		case strings.HasPrefix(t, "*"):
			t = strings.TrimPrefix(t, "*")
		default:
			goto unwrapped
		}
	}
unwrapped:
	// Skip arrays, maps, channels, funcs, interfaces and qualified types.
	if strings.ContainsAny(t, ".[]{}() <>") || strings.HasPrefix(t, "map") ||
		strings.HasPrefix(t, "chan") || strings.HasPrefix(t, "func") ||
		strings.HasPrefix(t, "interface") {

		return ""
	}
	// Known builtin types never need a stub.
	for _, vt := range ValidFieldTypes {
		if t == vt {
			return ""
		}
	}
	// A stubbable custom type is an exported identifier.
	if t == "" || t[0] < 'A' || t[0] > 'Z' {
		return ""
	}
	return t
}

// writeValidationMethod writes the Validate method for the entity.
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

// writeFieldValidation writes validation logic for a specific field.
func writeFieldValidation(content *strings.Builder, entityVar, entityName string, field Field) {
	switch field.Type {
	case FieldString:
		fmt.Fprintf(content, "\tif %s.%s == \"\" {\n", entityVar, field.Name)
		fmt.Fprintf(content, "\t\treturn ErrInvalid%s%s\n", entityName, field.Name)
		content.WriteString("\t}\n")
		if isEmailFieldName(field.Name) {
			// Minimal email-format validation (no external dependency).
			fmt.Fprintf(content, "\tif !strings.Contains(%s.%s, \"@\") || !strings.Contains(%s.%s, \".\") {\n",
				entityVar, field.Name, entityVar, field.Name)
			fmt.Fprintf(content, "\t\treturn ErrInvalid%s%s\n", entityName, field.Name)
			content.WriteString("\t}\n")
		}
	default:
		// "< 0" is only meaningful for signed integers and floats; unsigned
		// integers are always non-negative so no runtime check is emitted.
		if isSignedNumericType(field.Type) {
			fmt.Fprintf(content, "\tif %s.%s < 0 {\n", entityVar, field.Name)
			fmt.Fprintf(content, "\t\treturn ErrInvalid%s%s\n", entityName, field.Name)
			content.WriteString("\t}\n")
		}
	}
}

// isEmailFieldName reports whether a field name denotes an email field.
func isEmailFieldName(name string) bool {
	return strings.Contains(strings.ToLower(name), "email")
}

// hasEmailField reports whether any non-system field is an email field.
func hasEmailField(fields []Field) bool {
	for _, f := range fields {
		if f.Type == FieldString && !isSystemField(f.Name) && isEmailFieldName(f.Name) {
			return true
		}
	}
	return false
}

// writeSoftDeleteMethods writes soft delete helper methods.
func writeSoftDeleteMethods(content *strings.Builder, entityName string) {
	entityVar := strings.ToLower(string(entityName[0]))

	fmt.Fprintf(content, "func (%s *%s) SoftDelete() {\n", entityVar, entityName)
	fmt.Fprintf(content, "\t%s.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}\n", entityVar)
	content.WriteString("}\n\n")

	fmt.Fprintf(content, "func (%s *%s) IsDeleted() bool {\n", entityVar, entityName)
	fmt.Fprintf(content, "\treturn %s.DeletedAt.Valid\n", entityVar)
	content.WriteString("}\n\n")
}

// isSystemField checks if a field is a system-managed field.
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

func generateErrorsFile(dir, entityName string, fields []Field, sm ...*SafetyManager) {
	filename := filepath.Join(dir, "errors.go")
	existingErrors := readExistingErrors(filename, entityName)

	var content strings.Builder
	writeErrorsHeader(&content)
	writeEntityErrors(&content, entityName, fields, existingErrors)

	if err := writeGoFileMerged(filename, content.String(), sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing errors file: %v", err))
	}
}

// readExistingErrors reads existing error definitions from the file.
func readExistingErrors(filename, entityName string) []string {
	var existingErrors []string

	if _, err := os.Stat(filename); err == nil {
		ui.Warning(fmt.Sprintf("errors.go already exists, adding new errors for %s", entityName))

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

// writeErrorsHeader writes the package declaration and imports.
func writeErrorsHeader(content *strings.Builder) {
	content.WriteString("package domain\n\n")
	content.WriteString("import \"errors\"\n\n")
	content.WriteString("var (\n")
}

// writeEntityErrors writes all error definitions for the entity.
func writeEntityErrors(content *strings.Builder, entityName string, fields []Field, existingErrors []string) {
	var tmp strings.Builder
	writeGeneralError(&tmp, entityName, existingErrors)
	writeExistingErrors(&tmp, existingErrors)
	writeFieldErrors(&tmp, entityName, fields, existingErrors)

	// Deduplicate by error identifier. For example a field named "data" yields
	// ErrInvalid<Entity>Data, which collides with the general
	// ErrInvalid<Entity>Data constant and would otherwise be redeclared.
	seen := make(map[string]bool)
	for _, line := range strings.Split(tmp.String(), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if id := errorIdentifier(line); id != "" {
			if seen[id] {
				continue
			}
			seen[id] = true
		}
		content.WriteString(line + "\n")
	}
	content.WriteString(")\n")
}

// errorIdentifier extracts the Err… constant name from a declaration line.
func errorIdentifier(line string) string {
	t := strings.TrimSpace(line)
	if !strings.HasPrefix(t, "Err") {
		return ""
	}
	if i := strings.IndexAny(t, " ="); i > 0 {
		return t[:i]
	}
	return ""
}

// writeGeneralError writes the general entity error.
func writeGeneralError(content *strings.Builder, entityName string, existingErrors []string) {
	generalError := fmt.Sprintf("\tErrInvalid%sData = errors.New(\"invalid %s data\")",
		entityName, strings.ToLower(entityName))
	if !contains(existingErrors, generalError) {
		content.WriteString(generalError + "\n")
	}
}

// writeExistingErrors writes previously defined errors.
func writeExistingErrors(content *strings.Builder, existingErrors []string) {
	for _, err := range existingErrors {
		content.WriteString(err + "\n")
	}
}

// writeFieldErrors writes validation errors for all fields.
func writeFieldErrors(content *strings.Builder, entityName string, fields []Field, existingErrors []string) {
	for _, field := range fields {
		if isSystemField(field.Name) {
			continue
		}

		// Only declare ErrInvalid<Entity><Field> when Validate() actually emits a
		// check for it (string emptiness or signed-numeric "< 0"); otherwise the
		// constant would be declared but never used (ENTITY-9).
		if fieldHasBaseValidation(field.Type) {
			writeRequiredFieldError(content, entityName, field, existingErrors)
		}
		writeTypeSpecificErrors(content, entityName, field, existingErrors)
	}
}

// fieldHasBaseValidation reports whether writeFieldValidation emits a check that
// references the ErrInvalid<Entity><Field> "required" constant for this type.
func fieldHasBaseValidation(fieldType string) bool {
	return fieldType == FieldString || isSignedNumericType(fieldType)
}

// writeRequiredFieldError writes the required field error.
func writeRequiredFieldError(content *strings.Builder, entityName string, field Field, existingErrors []string) {
	fieldLower := strings.ToLower(field.Name)
	requiredError := fmt.Sprintf("\tErrInvalid%s%s = errors.New(\"%s is required\")",
		entityName, field.Name, fieldLower)
	if !contains(existingErrors, requiredError) {
		content.WriteString(requiredError + "\n")
	}
}

// writeTypeSpecificErrors writes type-specific validation errors.
func writeTypeSpecificErrors(content *strings.Builder, entityName string, field Field, existingErrors []string) {
	fieldLower := strings.ToLower(field.Name)

	switch {
	case field.Type == FieldString:
		writeStringFieldErrors(content, entityName, field, fieldLower, existingErrors)
	case field.Type == "float32" || field.Type == "float64":
		// Range error only for signed types where the runtime "< 0" check exists.
		writeFloatFieldErrors(content, entityName, field, fieldLower, existingErrors)
	case isSignedNumericType(field.Type):
		writeIntegerFieldErrors(content, entityName, field, fieldLower, existingErrors)
	}
}

// writeStringFieldErrors writes string-specific validation errors.
func writeStringFieldErrors(content *strings.Builder, entityName string, field Field, fieldLower string, existingErrors []string) {
	if strings.Contains(fieldLower, FieldEmailType) {
		emailError := fmt.Sprintf("\tErrInvalid%s%sFormat = errors.New(\"invalid %s format\")",
			entityName, field.Name, getFieldDisplayName(fieldLower))
		if !contains(existingErrors, emailError) {
			content.WriteString(emailError + "\n")
		}
	}

	if strings.Contains(fieldLower, "name") {
		lengthError := fmt.Sprintf("\tErrInvalid%s%sLength = errors.New(\"%s must be between 2 and 100 characters\")",
			entityName, field.Name, fieldLower)
		if !contains(existingErrors, lengthError) {
			content.WriteString(lengthError + "\n")
		}
	}
}

// writeIntegerFieldErrors writes integer-specific validation errors.
func writeIntegerFieldErrors(content *strings.Builder, entityName string, field Field, fieldLower string, existingErrors []string) {
	var rangeError string
	if strings.Contains(fieldLower, "age") {
		rangeError = fmt.Sprintf("\tErrInvalid%s%sRange = errors.New(\"%s must be greater than 0\")",
			entityName, field.Name, fieldLower)
	} else {
		rangeError = fmt.Sprintf("\tErrInvalid%s%sRange = errors.New(\"%s must be a positive number\")",
			entityName, field.Name, fieldLower)
	}
	if !contains(existingErrors, rangeError) {
		content.WriteString(rangeError + "\n")
	}
}

// writeFloatFieldErrors writes float-specific validation errors.
func writeFloatFieldErrors(content *strings.Builder, entityName string, field Field, fieldLower string, existingErrors []string) {
	var rangeError string
	if strings.Contains(fieldLower, "price") || strings.Contains(fieldLower, "amount") {
		rangeError = fmt.Sprintf("\tErrInvalid%s%sRange = errors.New(\"%s must be greater than 0 and less than 999,999,999.99\")",
			entityName, field.Name, fieldLower)
	} else {
		rangeError = fmt.Sprintf("\tErrInvalid%s%sRange = errors.New(\"%s must be a positive number\")",
			entityName, field.Name, fieldLower)
	}
	if !contains(existingErrors, rangeError) {
		content.WriteString(rangeError + "\n")
	}
}

// getFieldDisplayName converts field names to human-readable display names for error messages.
func getFieldDisplayName(fieldName string) string {
	fieldTranslations := map[string]string{
		"name":        "name",
		"email":       "email",
		"age":         "age",
		"price":       "price",
		"amount":      "amount",
		"description": "description",
		"title":       "title",
		"status":      "status",
		"category":    "category",
		"stock":       "stock",
		"quantity":    "quantity",
		"phone":       "phone",
		"address":     "address",
		"password":    "password",
	}

	for key, value := range fieldTranslations {
		if strings.Contains(fieldName, key) {
			return value
		}
	}

	return fieldName
}

// Helper function to check if a slice contains a string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == strings.TrimSpace(item) {
			return true
		}
	}
	return false
}

// generateSeedData creates seed data based on actual fields.
func generateSeedData(dir, entityName string, fields []Field, sm ...*SafetyManager) {
	filename := filepath.Join(dir, strings.ToLower(entityName)+"_seeds.go")

	// Build the body first so the import block reflects what is actually emitted.
	var body strings.Builder
	writeGoSeeds(&body, entityName, fields)
	writeSQLSeeds(&body, entityName, fields)

	var content strings.Builder
	writeSeedFileHeader(&content, body.String())
	content.WriteString(body.String())

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing seed file: %v", err))
	}
}

// writeSeedFileHeader writes the package declaration and imports for the seed
// file. The "time" import is only added when the generated body actually
// references the time package, to avoid an unused-import compile error.
func writeSeedFileHeader(content *strings.Builder, body string) {
	content.WriteString("package domain\n\n")

	if strings.Contains(body, "time.") {
		content.WriteString("import \"time\"\n\n")
	}
}

// writeGoSeeds writes the Go struct seed data function.
func writeGoSeeds(content *strings.Builder, entityName string, fields []Field) {
	fmt.Fprintf(content, "// Get%sSeeds returns sample data for %s\n", entityName, strings.ToLower(entityName))
	fmt.Fprintf(content, "func Get%sSeeds() []%s {\n", entityName, entityName)
	fmt.Fprintf(content, "\treturn []%s{\n", entityName)

	// Generate 3 sample records based on actual fields
	for i := 1; i <= 3; i++ {
		writeGoSeedRecord(content, fields, i)
	}

	content.WriteString("\t}\n")
	content.WriteString("}\n\n")
}

// writeGoSeedRecord writes a single Go seed record.
func writeGoSeedRecord(content *strings.Builder, fields []Field, recordNum int) {
	content.WriteString("\t\t{\n")
	for _, field := range fields {
		if isSystemField(field.Name) {
			continue // Skip auto-managed fields
		}

		sampleValue, ok := generateSampleValue(field, recordNum)
		if !ok {
			// No reliable sample for this type; rely on the Go zero value.
			continue
		}
		fmt.Fprintf(content, "\t\t\t%s: %s,\n", field.Name, sampleValue)
	}
	content.WriteString("\t\t},\n")
}

// writeSQLSeeds writes the SQL INSERT seed data function.
func writeSQLSeeds(content *strings.Builder, entityName string, fields []Field) {
	fmt.Fprintf(content, "// GetSQL%sSeeds returns SQL INSERT statements for %s\n", entityName, strings.ToLower(entityName))
	fmt.Fprintf(content, "func GetSQL%sSeeds() string {\n", entityName)
	fmt.Fprintf(content, "\treturn `-- Sample data for table %s\n", strings.ToLower(entityName))

	// Generate SQL INSERT statements
	for i := 1; i <= 3; i++ {
		writeSQLInsertStatement(content, entityName, fields, i)
	}

	content.WriteString("`\n")
	content.WriteString("}\n")
}

// writeSQLInsertStatement writes a single SQL INSERT statement.
func writeSQLInsertStatement(content *strings.Builder, entityName string, fields []Field, recordNum int) {
	fmt.Fprintf(content, "INSERT INTO %s (", strings.ToLower(entityName)+"s")

	// Field names
	fieldNames := getNonSystemFieldNames(fields)
	content.WriteString(strings.Join(fieldNames, ", "))
	content.WriteString(") VALUES (")

	// Field values
	values := getSQLFieldValues(fields, recordNum)
	content.WriteString(strings.Join(values, ", "))
	content.WriteString(");\n")
}

// getNonSystemFieldNames returns field names excluding system fields.
func getNonSystemFieldNames(fields []Field) []string {
	var fieldNames []string
	for _, field := range fields {
		if !isSystemField(field.Name) {
			fieldNames = append(fieldNames, strings.ToLower(field.Name))
		}
	}
	return fieldNames
}

// getSQLFieldValues returns SQL-formatted field values.
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

// generateSampleValue creates realistic sample data based on field type and name.
// The second return value reports whether a type-correct sample could be produced;
// when false, callers should omit the field and rely on the Go zero value.
func generateSampleValue(field Field, index int) (string, bool) {
	switch field.Type {
	case FieldString:
		return generateStringSampleValue(field.Name, index), true
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "byte", "rune":
		return generateIntSampleValue(field.Name, index), true
	case "float64", "float32":
		return generateFloatSampleValue(field.Name, index), true
	case FieldBool:
		return strconv.FormatBool(index%2 == 1), true
	case "time.Time":
		return "time.Now()", true
	case "[]byte":
		return fmt.Sprintf("[]byte(\"sample%d\")", index), true
	}
	return generateDefaultSampleValue(field.Type, index)
}

// generateStringSampleValue generates string sample values based on field name.
func generateStringSampleValue(fieldName string, index int) string {
	fieldLower := strings.ToLower(fieldName)

	switch {
	case strings.Contains(fieldLower, "name"):
		names := []string{"John Smith", "Jane Doe", "Bob Johnson"}
		return fmt.Sprintf("\"%s\"", names[(index-1)%len(names)])
	case strings.Contains(fieldLower, FieldEmailType):
		emails := []string{"john@example.com", "jane@example.com", "bob@example.com"}
		return fmt.Sprintf("\"%s\"", emails[(index-1)%len(emails)])
	case strings.Contains(fieldLower, "description"):
		descriptions := []string{"Detailed description of the first item", "Complete information about the second item", "Specific details of the third record"}
		return fmt.Sprintf("\"%s\"", descriptions[(index-1)%len(descriptions)])
	case strings.Contains(fieldLower, "title"):
		titles := []string{"Main Title", "Secondary Item", "Third Entry"}
		return fmt.Sprintf("\"%s\"", titles[(index-1)%len(titles)])
	case strings.Contains(fieldLower, "status"):
		statuses := []string{"active", "pending", "completed"}
		return fmt.Sprintf("\"%s\"", statuses[(index-1)%len(statuses)])
	case strings.Contains(fieldLower, "category"):
		categories := []string{"technology", "education", "health"}
		return fmt.Sprintf("\"%s\"", categories[(index-1)%len(categories)])
	default:
		return fmt.Sprintf("\"Sample %s %d\"", fieldName, index)
	}
}

// generateIntSampleValue generates integer sample values based on field name.
func generateIntSampleValue(fieldName string, index int) string {
	fieldLower := strings.ToLower(fieldName)

	switch {
	case strings.Contains(fieldLower, "age"):
		ages := []int{25, 30, 35}
		return strconv.Itoa(ages[(index-1)%len(ages)])
	case strings.Contains(fieldLower, "stock"):
		stocks := []int{100, 50, 75}
		return strconv.Itoa(stocks[(index-1)%len(stocks)])
	case strings.Contains(fieldLower, "quantity"):
		quantities := []int{10, 5, 15}
		return strconv.Itoa(quantities[(index-1)%len(quantities)])
	default:
		return strconv.Itoa(index * 10)
	}
}

// generateFloatSampleValue generates float sample values based on field name.
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

// generateDefaultSampleValue generates default sample values for composite or
// unknown types. The boolean reports whether a type-correct, compilable literal
// could be produced; when false the field should be omitted from the seed record.
func generateDefaultSampleValue(fieldType string, index int) (string, bool) {
	switch {
	case strings.HasPrefix(fieldType, "[]"):
		// Slice of a scalar element type -> single-element literal.
		elem := strings.TrimPrefix(fieldType, "[]")
		if sample, ok := scalarSampleLiteral(elem, index); ok {
			return fmt.Sprintf("%s{%s}", fieldType, sample), true
		}
		// Element type has no known literal; emit an empty slice (always valid).
		return fmt.Sprintf("%s{}", fieldType), true
	case strings.HasPrefix(fieldType, "map["):
		// Empty map literal is always valid.
		return fmt.Sprintf("%s{}", fieldType), true
	case strings.HasPrefix(fieldType, "*"):
		// Pointer types: nil is always valid.
		return "nil", true
	case strings.Contains(fieldType, FieldInt):
		return strconv.Itoa(index * 10), true
	case strings.Contains(fieldType, FieldString):
		return fmt.Sprintf("\"Sample%d\"", index), true
	case strings.Contains(fieldType, "float"):
		return fmt.Sprintf("%.2f", float64(index)*10.5), true
	case strings.Contains(fieldType, FieldBool):
		return strconv.FormatBool(index%2 == 1), true
	default:
		// Unknown named/custom type: no reliable inline literal, omit the field.
		return "", false
	}
}

// scalarSampleLiteral returns a sample literal for a scalar element type.
func scalarSampleLiteral(elem string, index int) (string, bool) {
	switch elem {
	case FieldString:
		return fmt.Sprintf("\"Sample%d\"", index), true
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "byte", "rune":
		return strconv.Itoa(index * 10), true
	case "float64", "float32":
		return fmt.Sprintf("%.2f", float64(index)*10.5), true
	case FieldBool:
		return strconv.FormatBool(index%2 == 1), true
	}
	return "", false
}

// generateSQLSampleValue creates SQL-compatible sample values.
func generateSQLSampleValue(field Field, index int) string {
	fieldLower := strings.ToLower(field.Name)

	switch field.Type {
	case "string":
		switch {
		case strings.Contains(fieldLower, "name"):
			names := []string{"John Smith", "Jane Doe", "Bob Johnson"}
			return fmt.Sprintf("'%s'", names[(index-1)%len(names)])
		case strings.Contains(fieldLower, "email"):
			emails := []string{"john@example.com", "jane@example.com", "bob@example.com"}
			return fmt.Sprintf("'%s'", emails[(index-1)%len(emails)])
		case strings.Contains(fieldLower, "description"):
			descriptions := []string{"Detailed description of the first item", "Complete information about the second item", "Specific details of the third record"}
			return fmt.Sprintf("'%s'", descriptions[(index-1)%len(descriptions)])
		case strings.Contains(fieldLower, "status"):
			statuses := []string{"active", "pending", "completed"}
			return fmt.Sprintf("'%s'", statuses[(index-1)%len(statuses)])
		default:
			return fmt.Sprintf("'Sample %s %d'", field.Name, index)
		}

	case "int", "int64", "uint", "uint64":
		switch {
		case strings.Contains(fieldLower, "age"):
			ages := []int{25, 30, 35}
			return strconv.Itoa(ages[(index-1)%len(ages)])
		case strings.Contains(fieldLower, "stock"):
			stocks := []int{100, 50, 75}
			return strconv.Itoa(stocks[(index-1)%len(stocks)])
		default:
			return strconv.Itoa(index * 10)
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
		return strconv.FormatBool(index%2 == 1)

	case "time.Time":
		return "NOW()"

	default:
		return generateSQLCompositeSampleValue(field, index)
	}
}

// generateSQLCompositeSampleValue produces a non-NULL SQL literal for composite
// or custom field types (slices, pointers, custom named types). It keeps SQL
// seeds consistent with the Go seeds and valid for NOT NULL columns (ENTITY-7).
func generateSQLCompositeSampleValue(field Field, index int) string {
	ft := field.Type
	switch {
	case strings.HasPrefix(ft, "[]"):
		// Slices are persisted as text (e.g. comma-separated / JSON-like).
		elem := strings.TrimPrefix(ft, "[]")
		switch {
		case elem == FieldString:
			return fmt.Sprintf("'sample%d,sample%d'", index, index+1)
		case isNumericType(elem):
			return fmt.Sprintf("'%d,%d'", index*10, index*10+1)
		case elem == FieldBool:
			return "'true,false'"
		default:
			return fmt.Sprintf("'sample%d'", index)
		}
	case strings.HasPrefix(ft, "*"):
		// Pointer: emit a value of the underlying type.
		base := strings.TrimPrefix(ft, "*")
		return generateSQLSampleValue(Field{Name: field.Name, Type: base}, index)
	case ft == "time.Time":
		return "NOW()"
	default:
		// Custom/named types (e.g. UserStatus): treat as text.
		return fmt.Sprintf("'sample%d'", index)
	}
}

func init() {
	entityCmd.Flags().StringP("fields", "f", "", "Entity fields \"field:type,field2:type\" (required)")
	entityCmd.Flags().Bool("validation", false, "Include business validations")
	entityCmd.Flags().BoolP("business-rules", "b", false, "Include advanced business rules")
	entityCmd.Flags().BoolP("timestamps", "t", false, "Include CreatedAt and UpdatedAt fields")
	entityCmd.Flags().BoolP("soft-delete", "s", false, "Include soft delete (DeletedAt)")
	entityCmd.Flags().Bool("tests", true, "Generate unit tests for the entity")
	entityCmd.Flags().Bool("dry-run", false, "Preview changes without creating files")
	entityCmd.Flags().Bool("force", false, "Overwrite existing files without asking")
	entityCmd.Flags().Bool("backup", false, "Backup existing files before overwriting")
	_ = entityCmd.MarkFlagRequired("fields")
}
