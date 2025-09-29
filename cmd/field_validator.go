package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// FieldValidator validates field definitions and entity names
type FieldValidator struct{}

// NewFieldValidator creates a new field validator instance
func NewFieldValidator() *FieldValidator {
	return &FieldValidator{}
}

// ValidateFields validates the complete fields string
func (v *FieldValidator) ValidateFields(fields string) error {
	if fields == "" {
		return fmt.Errorf("%s", ErrEmptyFields)
	}

	fieldParts := v.smartSplitFields(fields)
	if len(fieldParts) == 0 {
		return fmt.Errorf("%s", ErrEmptyFields)
	}

	fieldNames := make(map[string]bool)

	for _, fieldPart := range fieldParts {
		field, err := v.ValidateField(strings.TrimSpace(fieldPart))
		if err != nil {
			return err
		}

		// Check for duplicate field names
		if fieldNames[field.Name] {
			return fmt.Errorf("campo duplicado: %s", field.Name)
		}
		fieldNames[field.Name] = true

		// Validate reserved field names
		if err := v.ValidateReservedNames(field.Name); err != nil {
			return err
		}
	}

	return nil
}

// smartSplitFields splits fields while respecting parentheses and brackets
func (v *FieldValidator) smartSplitFields(fields string) []string {
	var result []string
	var current strings.Builder
	parenLevel := 0
	bracketLevel := 0

	for _, char := range fields {
		switch char {
		case '(':
			parenLevel++
			current.WriteRune(char)
		case ')':
			parenLevel--
			current.WriteRune(char)
		case '[':
			bracketLevel++
			current.WriteRune(char)
		case ']':
			bracketLevel--
			current.WriteRune(char)
		case ',':
			if parenLevel == 0 && bracketLevel == 0 {
				// We're at the top level, this comma separates fields
				if current.Len() > 0 {
					result = append(result, current.String())
					current.Reset()
				}
			} else {
				// We're inside parentheses or brackets, keep the comma
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	// Add the last field
	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// ValidateField validates a single field definition
func (v *FieldValidator) ValidateField(fieldDef string) (*Field, error) {
	if fieldDef == "" {
		return nil, fmt.Errorf("definición de campo vacía")
	}

	parts := strings.Split(fieldDef, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%s. Recibido: '%s'", ErrInvalidFieldSyntax, fieldDef)
	}

	fieldName := strings.TrimSpace(parts[0])
	fieldType := strings.TrimSpace(parts[1])

	// Validate field name
	if err := v.ValidateFieldName(fieldName); err != nil {
		return nil, err
	}

	// Validate field type
	if err := v.ValidateFieldType(fieldType); err != nil {
		return nil, err
	}

	return &Field{
		Name: capitalizeFirst(fieldName),
		Type: fieldType,
	}, nil
}

// ValidateFieldName validates a field name
func (v *FieldValidator) ValidateFieldName(name string) error {
	if name == "" {
		return fmt.Errorf("nombre de campo no puede estar vacío")
	}

	if len(name) < MinFieldNameLength || len(name) > MaxFieldNameLength {
		return fmt.Errorf("nombre de campo debe tener entre %d y %d caracteres", MinFieldNameLength, MaxFieldNameLength)
	}

	// Check if it starts with a letter
	if !unicode.IsLetter(rune(name[0])) {
		return fmt.Errorf("nombre de campo debe empezar con una letra: %s", name)
	}

	// Check if it contains only valid characters (letters, numbers, underscore)
	validName := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("nombre de campo contiene caracteres inválidos: %s. Solo se permiten letras, números y guiones bajos", name)
	}

	return nil
}

// ValidateFieldType validates a field type with comprehensive Go type support
func (v *FieldValidator) ValidateFieldType(fieldType string) error {
	if fieldType == "" {
		return fmt.Errorf("tipo de campo no puede estar vacío")
	}

	// Check if it's a valid basic type
	for _, validType := range ValidFieldTypes {
		if fieldType == validType {
			return nil
		}
	}

	// Enhanced validation for complex Go types
	return v.validateComplexType(fieldType)
}

// validateComplexType handles validation of complex Go types
func (v *FieldValidator) validateComplexType(fieldType string) error {
	// Trim leading/trailing whitespace but preserve internal spaces for channel/func parsing
	fieldType = strings.TrimSpace(fieldType)

	// Check for slice types ([]string, []int, []*User, etc.)
	if strings.HasPrefix(fieldType, "[]") {
		baseType := strings.TrimPrefix(fieldType, "[]")
		return v.ValidateFieldType(baseType)
	}

	// Check for array types ([10]string, [5]int, etc.)
	arrayPattern := regexp.MustCompile(`^\[\d+\](.+)$`)
	if arrayPattern.MatchString(fieldType) {
		matches := arrayPattern.FindStringSubmatch(fieldType)
		if len(matches) > 1 {
			return v.ValidateFieldType(matches[1])
		}
	}

	// Check for pointer types (*string, *int, *User, etc.)
	if strings.HasPrefix(fieldType, "*") {
		baseType := strings.TrimPrefix(fieldType, "*")
		return v.ValidateFieldType(baseType)
	}

	// Enhanced map type validation (map[string]interface{}, map[int]*User, etc.)
	if strings.HasPrefix(fieldType, "map[") {
		return v.validateMapType(fieldType)
	}

	// Check for channel types (chan string, <-chan int, chan<- bool, etc.)
	if v.isChannelType(fieldType) {
		return v.validateChannelType(fieldType)
	}

	// Check for function types (func(string) error, func(int, string) (bool, error), etc.)
	if strings.HasPrefix(fieldType, "func") {
		return v.validateFunctionType(fieldType)
	}

	// Check for interface types (interface{}, io.Reader, etc.)
	if v.isInterfaceType(fieldType) {
		return nil // Interface types are generally valid
	}

	// Check for qualified types (package.Type, time.Time, etc.)
	if v.isQualifiedType(fieldType) {
		return nil // Assume qualified types are valid
	}

	// Check for custom struct types (User, Product, etc.)
	if v.isCustomType(fieldType) {
		return nil // Custom types are valid
	}

	return fmt.Errorf("%s: %s. Tipos válidos incluyen: %s, slices, pointers, maps, channels, functions, interfaces y tipos personalizados",
		ErrInvalidFieldType, fieldType, strings.Join(ValidFieldTypes, ", "))
}

// validateMapType validates map types with comprehensive key/value type support
func (v *FieldValidator) validateMapType(fieldType string) error {
	// Pattern for map[keyType]valueType
	mapPattern := regexp.MustCompile(`^map\[([^\]]+)\](.+)$`)
	matches := mapPattern.FindStringSubmatch(fieldType)

	if len(matches) != 3 {
		return fmt.Errorf("formato de map inválido: %s. Formato esperado: map[keyType]valueType", fieldType)
	}

	keyType := matches[1]
	valueType := matches[2]

	// Validate key type (must be comparable in Go)
	if err := v.validateMapKeyType(keyType); err != nil {
		return fmt.Errorf("tipo de clave de map inválido: %w", err)
	}

	// Validate value type (can be any type)
	if err := v.ValidateFieldType(valueType); err != nil {
		return fmt.Errorf("tipo de valor de map inválido: %w", err)
	}

	return nil
}

// validateMapKeyType validates that a type can be used as a map key (must be comparable)
func (v *FieldValidator) validateMapKeyType(keyType string) error {
	// Go comparable types that can be map keys
	comparableTypes := []string{
		"string", "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		"byte", "rune", "float32", "float64", "bool",
	}

	// Check basic comparable types
	for _, comparable := range comparableTypes {
		if keyType == comparable {
			return nil
		}
	}

	// Pointer types are comparable
	if strings.HasPrefix(keyType, "*") {
		return nil
	}

	// Array types are comparable if their element type is comparable
	arrayPattern := regexp.MustCompile(`^\[\d+\](.+)$`)
	if arrayPattern.MatchString(keyType) {
		matches := arrayPattern.FindStringSubmatch(keyType)
		if len(matches) > 1 {
			return v.validateMapKeyType(matches[1])
		}
	}

	// Interface types are comparable
	if v.isInterfaceType(keyType) {
		return nil
	}

	// Custom types and qualified types are assumed comparable
	if v.isCustomType(keyType) || v.isQualifiedType(keyType) {
		return nil
	}

	// Slices, maps, and functions are not comparable
	if strings.HasPrefix(keyType, "[]") || strings.HasPrefix(keyType, "map[") || strings.HasPrefix(keyType, "func") {
		return fmt.Errorf("tipo %s no es comparable y no puede usarse como clave de map", keyType)
	}

	return nil
}

// validateChannelType validates channel types
func (v *FieldValidator) validateChannelType(fieldType string) error {
	// Handle specific channel patterns

	// Pattern for send-only: chan<- Type
	sendOnlyPattern := regexp.MustCompile(`^chan<-\s*(.+)$`)
	if sendOnlyPattern.MatchString(fieldType) {
		matches := sendOnlyPattern.FindStringSubmatch(fieldType)
		if len(matches) > 1 {
			elementType := strings.TrimSpace(matches[1])
			return v.ValidateFieldType(elementType)
		}
	}

	// Pattern for receive-only: <-chan Type
	receiveOnlyPattern := regexp.MustCompile(`^<-chan\s*(.+)$`)
	if receiveOnlyPattern.MatchString(fieldType) {
		matches := receiveOnlyPattern.FindStringSubmatch(fieldType)
		if len(matches) > 1 {
			elementType := strings.TrimSpace(matches[1])
			return v.ValidateFieldType(elementType)
		}
	}

	// Pattern for bidirectional: chan Type
	bidirectionalPattern := regexp.MustCompile(`^chan\s+(.+)$`)
	if bidirectionalPattern.MatchString(fieldType) {
		matches := bidirectionalPattern.FindStringSubmatch(fieldType)
		if len(matches) > 1 {
			elementType := strings.TrimSpace(matches[1])
			return v.ValidateFieldType(elementType)
		}
	}

	return fmt.Errorf("formato de channel inválido: %s", fieldType)
}

// validateFunctionType validates function types
func (v *FieldValidator) validateFunctionType(fieldType string) error {
	// More flexible function signature validation
	// Allow: func(), func(params), func(params) return, func(params) (returns)
	funcPattern := regexp.MustCompile(`^func\s*(\([^)]*\))?\s*(.+)?$`)
	if funcPattern.MatchString(fieldType) {
		return nil // Basic function signature is valid
	}

	return fmt.Errorf("formato de función inválido: %s", fieldType)
}

// isChannelType checks if a type is a channel type
func (v *FieldValidator) isChannelType(fieldType string) bool {
	return strings.HasPrefix(fieldType, "chan ") ||
		strings.HasPrefix(fieldType, "chan<-") ||
		strings.HasPrefix(fieldType, "<-chan")
}

// isInterfaceType checks if a type is an interface type
func (v *FieldValidator) isInterfaceType(fieldType string) bool {
	return fieldType == "interface{}" || strings.HasSuffix(fieldType, "interface{}")
}

// isQualifiedType checks if a type is a qualified type (package.Type)
func (v *FieldValidator) isQualifiedType(fieldType string) bool {
	return strings.Contains(fieldType, ".") && !strings.HasPrefix(fieldType, ".")
}

// isCustomType checks if a type is a custom type (starts with uppercase letter)
func (v *FieldValidator) isCustomType(fieldType string) bool {
	if len(fieldType) == 0 {
		return false
	}
	return unicode.IsUpper(rune(fieldType[0])) && regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`).MatchString(fieldType)
}

// ValidateEntityName validates an entity name
func (v *FieldValidator) ValidateEntityName(name string) error {
	if name == "" {
		return fmt.Errorf("nombre de entidad no puede estar vacío")
	}

	if len(name) < MinEntityNameLength || len(name) > MaxEntityNameLength {
		return fmt.Errorf("nombre de entidad debe tener entre %d y %d caracteres", MinEntityNameLength, MaxEntityNameLength)
	}

	// Check if it starts with a capital letter
	if !unicode.IsUpper(rune(name[0])) {
		return fmt.Errorf("nombre de entidad debe empezar con mayúscula: %s", name)
	}

	// Check if it contains only valid characters (letters and numbers)
	validName := regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("nombre de entidad contiene caracteres inválidos: %s. Solo se permiten letras y números", name)
	}

	return nil
}

// ValidateDatabase validates a database type
func (v *FieldValidator) ValidateDatabase(database string) error {
	for _, validDB := range ValidDatabases {
		if database == validDB {
			return nil
		}
	}
	return fmt.Errorf("%s. Recibido: %s", ErrInvalidDatabase, database)
}

// ValidateHandlers validates handler types
func (v *FieldValidator) ValidateHandlers(handlers string) error {
	if handlers == "" {
		return nil // Empty is valid, will use default
	}

	handlerList := strings.Split(handlers, ",")
	for _, handler := range handlerList {
		handler = strings.TrimSpace(handler)
		found := false
		for _, validHandler := range ValidHandlers {
			if handler == validHandler {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%s. Recibido: %s", ErrInvalidHandler, handler)
		}
	}
	return nil
}

// ValidateOperations validates operation types
func (v *FieldValidator) ValidateOperations(operations string) error {
	if operations == "" {
		return nil // Empty is valid, will use default
	}

	opList := strings.Split(operations, ",")
	for _, op := range opList {
		op = strings.TrimSpace(op)
		found := false
		for _, validOp := range ValidOperations {
			if op == validOp {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%s. Recibido: %s", ErrInvalidOperation, op)
		}
	}
	return nil
}

// ValidateReservedNames checks for Go reserved words and common conflicts
func (v *FieldValidator) ValidateReservedNames(name string) error {
	lowerName := strings.ToLower(name)

	// Go reserved words
	goReserved := []string{
		"break", "case", "chan", "const", "continue", "default", "defer", "else",
		"fallthrough", "for", "func", "go", "goto", "if", "import", "interface",
		"map", "package", "range", "return", "select", "struct", "switch", "type", "var",
	}

	// Common field names that might cause conflicts
	conflictNames := []string{
		"id", "ID", "string", "int", "bool", "true", "false", "nil", "len", "cap",
		"make", "new", "delete", "copy", "append", "panic", "recover", "print", "println",
		"error", "Error",
	}

	for _, reserved := range goReserved {
		if lowerName == reserved {
			return fmt.Errorf("'%s' es una palabra reservada de Go", name)
		}
	}

	for _, conflict := range conflictNames {
		if lowerName == strings.ToLower(conflict) {
			return fmt.Errorf("'%s' puede causar conflictos. Usa un nombre diferente", name)
		}
	}

	return nil
}

// ParseFieldsWithValidation parses and validates fields string
func (v *FieldValidator) ParseFieldsWithValidation(fields string) ([]Field, error) {
	if err := v.ValidateFields(fields); err != nil {
		return nil, err
	}

	var fieldsList []Field

	// Always add ID field with GORM tags
	fieldsList = append(fieldsList, Field{
		Name: "ID",
		Type: "uint",
		Tag:  "`json:\"id\" gorm:\"primaryKey;autoIncrement\"`",
	})

	parts := strings.Split(fields, ",")
	for _, part := range parts {
		field, err := v.ValidateField(strings.TrimSpace(part))
		if err != nil {
			return nil, err
		}

		// Generate GORM tag based on field type
		gormTag := getGormTag(field.Name, field.Type)
		tag := fmt.Sprintf("`json:\"%s\" gorm:\"%s\"`", strings.ToLower(field.Name), gormTag)

		fieldsList = append(fieldsList, Field{
			Name: field.Name,
			Type: field.Type,
			Tag:  tag,
		})
	}

	return fieldsList, nil
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

// GenerateQueryMethodsForFields generates appropriate query methods based on field types
func (v *FieldValidator) GenerateQueryMethodsForFields(entity string, fields []Field) []QueryMethod {
	var methods []QueryMethod
	entityLower := strings.ToLower(entity)

	// Always include basic methods
	methods = append(methods, QueryMethod{
		Name:       fmt.Sprintf("FindBy%sID", entity),
		MethodName: "FindByID",
		Field:      "id",
		Type:       "uint",
	})

	for _, field := range fields {
		if field.Name == "ID" {
			continue
		}

		fieldLower := strings.ToLower(field.Name)

		// Generate FindBy methods for string fields that might be unique
		if field.Type == FieldString && v.isLikelyUniqueField(fieldLower) {
			methods = append(methods, QueryMethod{
				Name:       fmt.Sprintf("FindBy%s%s", entity, field.Name),
				MethodName: fmt.Sprintf("FindBy%s", field.Name),
				Field:      fieldLower,
				Type:       field.Type,
			})
		}

		// Generate FindBy methods for commonly queried fields
		if v.isCommonQueryField(entityLower, fieldLower) {
			methods = append(methods, QueryMethod{
				Name:       fmt.Sprintf("FindBy%s%s", entity, field.Name),
				MethodName: fmt.Sprintf("FindBy%s", field.Name),
				Field:      fieldLower,
				Type:       field.Type,
			})
		}
	}

	return methods
}

// isLikelyUniqueField checks if a field is likely to be unique
func (v *FieldValidator) isLikelyUniqueField(fieldName string) bool {
	uniqueFields := []string{"email", "username", "slug", "code", "sku", "token", "uuid"}
	for _, unique := range uniqueFields {
		if strings.Contains(fieldName, unique) {
			return true
		}
	}
	return false
}

// isCommonQueryField checks if a field is commonly used for queries
func (v *FieldValidator) isCommonQueryField(entityName, fieldName string) bool {
	// Check general common query fields
	for _, common := range CommonQueryFields {
		if fieldName == common {
			return true
		}
	}

	// General common query fields
	generalCommon := []string{"name", "title", "status", "type", "category", "author_id", "user_id"}
	for _, common := range generalCommon {
		if fieldName == common || strings.HasSuffix(fieldName, "_id") {
			return true
		}
	}

	return false
}

// QueryMethod represents a dynamic query method
type QueryMethod struct {
	Name       string // Full method name (FindByUserEmail)
	MethodName string // Method name only (FindByEmail)
	Field      string // Field name (email)
	Type       string // Field type (string)
}
