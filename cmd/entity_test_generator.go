package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// generateEntityTests generates unit tests for entity validation and business logic
func generateEntityTests(domainDir, entityName string, fields []Field, validation, businessRules bool, fileNamingConvention string) {
	// Determine filename based on naming convention
	var filename string
	switch fileNamingConvention {
	case "snake":
		filename = toSnakeCase(entityName) + "_test.go"
	default: // "lowercase" or any other
		filename = strings.ToLower(entityName) + "_test.go"
	}

	testFile := filepath.Join(domainDir, filename)

	var content strings.Builder

	// Package declaration
	content.WriteString("package domain\n\n")

	// Imports - no time import since we don't test timestamp fields
	content.WriteString("import (\n")
	content.WriteString("\t\"testing\"\n")
	content.WriteString("\n\t\"github.com/stretchr/testify/assert\"\n")
	content.WriteString(")\n\n") // Generate validation tests if validation is enabled
	if validation {
		generateValidationTests(&content, entityName, fields)
	}

	// Generate constructor/initialization tests
	generateConstructorTests(&content, entityName, fields)

	// Generate field-specific tests
	generateFieldTests(&content, entityName, fields)

	// Write file
	if err := os.WriteFile(testFile, []byte(content.String()), 0644); err != nil {
		fmt.Printf("Error writing test file: %v\n", err)
		return
	}

	fmt.Printf("✓ Generated test file: %s\n", testFile)
}

// generateValidationTests creates table-driven tests for Validate() method
func generateValidationTests(content *strings.Builder, entityName string, fields []Field) {
	entityLower := strings.ToLower(string(entityName[0])) + entityName[1:]

	content.WriteString(fmt.Sprintf("// Test%s_Validate tests the Validate method with various scenarios\n", entityName))
	content.WriteString(fmt.Sprintf("func Test%s_Validate(t *testing.T) {\n", entityName))
	content.WriteString("\ttests := []struct {\n")
	content.WriteString("\t\tname    string\n")
	content.WriteString(fmt.Sprintf("\t\t%s    %s\n", entityLower, entityName))
	content.WriteString("\t\twantErr bool\n")
	content.WriteString("\t\terrMsg   string\n")
	content.WriteString("\t}{\n")

	// Valid case
	content.WriteString("\t\t{\n")
	content.WriteString("\t\t\tname: \"valid entity\",\n")
	content.WriteString(fmt.Sprintf("\t\t\t%s: %s{\n", entityLower, entityName))

	// Generate valid field values
	for _, field := range fields {
		if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
			continue
		}

		validValue := getValidFieldValue(field)
		content.WriteString(fmt.Sprintf("\t\t\t\t%s: %s,\n", field.Name, validValue))
	}

	content.WriteString("\t\t\t},\n")
	content.WriteString("\t\t\twantErr: false,\n")
	content.WriteString("\t\t},\n")

	// Invalid cases for each field
	for _, field := range fields {
		if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
			continue
		}

		generateInvalidTestCase(content, entityName, entityLower, field, fields)
	}

	content.WriteString("\t}\n\n")

	// Test execution
	content.WriteString("\tfor _, tt := range tests {\n")
	content.WriteString("\t\tt.Run(tt.name, func(t *testing.T) {\n")
	content.WriteString(fmt.Sprintf("\t\t\terr := tt.%s.Validate()\n", entityLower))
	content.WriteString("\t\t\tif tt.wantErr {\n")
	content.WriteString("\t\t\t\tassert.Error(t, err)\n")
	content.WriteString("\t\t\t\tif tt.errMsg != \"\" {\n")
	content.WriteString("\t\t\t\t\tassert.Contains(t, err.Error(), tt.errMsg)\n")
	content.WriteString("\t\t\t\t}\n")
	content.WriteString("\t\t\t} else {\n")
	content.WriteString("\t\t\t\tassert.NoError(t, err)\n")
	content.WriteString("\t\t\t}\n")
	content.WriteString("\t\t})\n")
	content.WriteString("\t}\n")
	content.WriteString("}\n\n")
}

// generateInvalidTestCase creates test cases for invalid field values
func generateInvalidTestCase(content *strings.Builder, entityName, entityLower string, invalidField Field, allFields []Field) {
	fieldNameLower := strings.ToLower(invalidField.Name)

	content.WriteString("\t\t{\n")
	content.WriteString(fmt.Sprintf("\t\t\tname: \"invalid %s - %s\",\n", entityLower, getInvalidDescription(invalidField)))
	content.WriteString(fmt.Sprintf("\t\t\t%s: %s{\n", entityLower, entityName))

	// Generate field values
	for _, field := range allFields {
		if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
			continue
		}

		// Use invalid value for the field being tested
		if field.Name == invalidField.Name {
			invalidValue := getInvalidFieldValue(field)
			content.WriteString(fmt.Sprintf("\t\t\t\t%s: %s,\n", field.Name, invalidValue))
		} else {
			// Use valid values for other fields
			validValue := getValidFieldValue(field)
			content.WriteString(fmt.Sprintf("\t\t\t\t%s: %s,\n", field.Name, validValue))
		}
	}

	content.WriteString("\t\t\t},\n")
	content.WriteString("\t\t\twantErr: true,\n")
	content.WriteString(fmt.Sprintf("\t\t\terrMsg: \"%s\",\n", fieldNameLower))
	content.WriteString("\t\t},\n")
}

// generateConstructorTests creates tests for entity initialization
func generateConstructorTests(content *strings.Builder, entityName string, fields []Field) {
	entityLower := strings.ToLower(string(entityName[0])) + entityName[1:]

	content.WriteString(fmt.Sprintf("// Test%s_Initialization tests entity field initialization\n", entityName))
	content.WriteString(fmt.Sprintf("func Test%s_Initialization(t *testing.T) {\n", entityName))
	content.WriteString(fmt.Sprintf("\t%s := &%s{\n", entityLower, entityName))

	// Generate sample initialization values
	for _, field := range fields {
		if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
			continue
		}

		validValue := getValidFieldValue(field)
		content.WriteString(fmt.Sprintf("\t\t%s: %s,\n", field.Name, validValue))
	}

	content.WriteString("\t}\n\n")

	// Assert all fields are set correctly
	for _, field := range fields {
		if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
			continue
		}

		expectedValue := getValidFieldValue(field)
		content.WriteString(fmt.Sprintf("\tassert.Equal(t, %s, %s.%s, \"%s should be set correctly\")\n",
			expectedValue, entityLower, field.Name, field.Name))
	}

	content.WriteString("}\n\n")
}

// generateFieldTests creates specific tests for field constraints
func generateFieldTests(content *strings.Builder, entityName string, fields []Field) {
	entityLower := strings.ToLower(string(entityName[0])) + entityName[1:]

	for _, field := range fields {
		if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
			continue
		}

		// Test for string fields (length, empty)
		if field.Type == "string" {
			content.WriteString(fmt.Sprintf("// Test%s_%s_EdgeCases tests edge cases for %s field\n",
				entityName, field.Name, field.Name))
			content.WriteString(fmt.Sprintf("func Test%s_%s_EdgeCases(t *testing.T) {\n", entityName, field.Name))
			content.WriteString("\ttests := []struct {\n")
			content.WriteString("\t\tname    string\n")
			content.WriteString("\t\tvalue   string\n")
			content.WriteString("\t\twantErr bool\n")
			content.WriteString("\t}{\n")

			// Valid cases
			content.WriteString("\t\t{name: \"valid value\", value: \"Valid Name\", wantErr: false},\n")
			content.WriteString("\t\t{name: \"empty string\", value: \"\", wantErr: true},\n")

			// Email-specific tests
			if strings.Contains(strings.ToLower(field.Name), "email") {
				content.WriteString("\t\t{name: \"valid email\", value: \"test@example.com\", wantErr: false},\n")
				content.WriteString("\t\t{name: \"invalid email format\", value: \"notanemail\", wantErr: true},\n")
			}

			content.WriteString("\t}\n\n")
			content.WriteString("\tfor _, tt := range tests {\n")
			content.WriteString("\t\tt.Run(tt.name, func(t *testing.T) {\n")
			content.WriteString(fmt.Sprintf("\t\t\t%s := &%s{\n", entityLower, entityName))

			// Set all required fields with valid values
			for _, f := range fields {
				if f.Name == "ID" || f.Name == "CreatedAt" || f.Name == "UpdatedAt" || f.Name == "DeletedAt" {
					continue
				}

				if f.Name == field.Name {
					content.WriteString(fmt.Sprintf("\t\t\t\t%s: tt.value,\n", f.Name))
				} else {
					validValue := getValidFieldValue(f)
					content.WriteString(fmt.Sprintf("\t\t\t\t%s: %s,\n", f.Name, validValue))
				}
			}

			content.WriteString("\t\t\t}\n\n")
			content.WriteString(fmt.Sprintf("\t\t\terr := %s.Validate()\n", entityLower))
			content.WriteString("\t\t\tif tt.wantErr {\n")
			content.WriteString("\t\t\t\tassert.Error(t, err)\n")
			content.WriteString("\t\t\t} else {\n")
			content.WriteString("\t\t\t\tassert.NoError(t, err)\n")
			content.WriteString("\t\t\t}\n")
			content.WriteString("\t\t})\n")
			content.WriteString("\t}\n")
			content.WriteString("}\n\n")
		}

		// Test for numeric fields (negative, zero, positive)
		if field.Type == "int" || field.Type == "int64" || field.Type == "float64" {
			content.WriteString(fmt.Sprintf("// Test%s_%s_NumericValidation tests numeric validation for %s\n",
				entityName, field.Name, field.Name))
			content.WriteString(fmt.Sprintf("func Test%s_%s_NumericValidation(t *testing.T) {\n", entityName, field.Name))
			content.WriteString("\ttests := []struct {\n")
			content.WriteString("\t\tname    string\n")
			content.WriteString(fmt.Sprintf("\t\tvalue   %s\n", field.Type))
			content.WriteString("\t\twantErr bool\n")
			content.WriteString("\t}{\n")

			if field.Type == "int" || field.Type == "int64" {
				content.WriteString("\t\t{name: \"positive value\", value: 10, wantErr: false},\n")
				content.WriteString("\t\t{name: \"zero value\", value: 0, wantErr: false},\n")
				content.WriteString("\t\t{name: \"negative value\", value: -1, wantErr: true},\n")
			} else if field.Type == "float64" {
				content.WriteString("\t\t{name: \"positive value\", value: 10.5, wantErr: false},\n")
				content.WriteString("\t\t{name: \"zero value\", value: 0.0, wantErr: false},\n")
				content.WriteString("\t\t{name: \"negative value\", value: -1.5, wantErr: true},\n")
			}

			content.WriteString("\t}\n\n")
			content.WriteString("\tfor _, tt := range tests {\n")
			content.WriteString("\t\tt.Run(tt.name, func(t *testing.T) {\n")
			content.WriteString(fmt.Sprintf("\t\t\t%s := &%s{\n", entityLower, entityName))

			// Set all required fields with valid values
			for _, f := range fields {
				if f.Name == "ID" || f.Name == "CreatedAt" || f.Name == "UpdatedAt" || f.Name == "DeletedAt" {
					continue
				}

				if f.Name == field.Name {
					content.WriteString(fmt.Sprintf("\t\t\t\t%s: tt.value,\n", f.Name))
				} else {
					validValue := getValidFieldValue(f)
					content.WriteString(fmt.Sprintf("\t\t\t\t%s: %s,\n", f.Name, validValue))
				}
			}

			content.WriteString("\t\t\t}\n\n")
			content.WriteString(fmt.Sprintf("\t\t\terr := %s.Validate()\n", entityLower))
			content.WriteString("\t\t\tif tt.wantErr {\n")
			content.WriteString("\t\t\t\tassert.Error(t, err)\n")
			content.WriteString("\t\t\t} else {\n")
			content.WriteString("\t\t\t\tassert.NoError(t, err)\n")
			content.WriteString("\t\t\t}\n")
			content.WriteString("\t\t})\n")
			content.WriteString("\t}\n")
			content.WriteString("}\n\n")
		}
	}
}

// Helper functions to generate test values

func getValidFieldValue(field Field) string {
	switch field.Type {
	case "string":
		if strings.Contains(strings.ToLower(field.Name), "email") {
			return "\"test@example.com\""
		}
		if strings.Contains(strings.ToLower(field.Name), "name") {
			return "\"John Doe\""
		}
		if strings.Contains(strings.ToLower(field.Name), "description") {
			return "\"A valid description\""
		}
		if strings.Contains(strings.ToLower(field.Name), "status") {
			return "\"active\""
		}
		return "\"valid value\""
	case "int", "int64":
		if strings.Contains(strings.ToLower(field.Name), "age") {
			return "25"
		}
		if strings.Contains(strings.ToLower(field.Name), "quantity") {
			return "10"
		}
		return "1"
	case "float64":
		if strings.Contains(strings.ToLower(field.Name), "price") {
			return "99.99"
		}
		if strings.Contains(strings.ToLower(field.Name), "total") {
			return "100.50"
		}
		return "10.5"
	case "bool":
		return "true"
	case "time.Time":
		return "time.Now()"
	default:
		return "\"\""
	}
}

func getInvalidFieldValue(field Field) string {
	switch field.Type {
	case "string":
		return "\"\""
	case "int", "int64":
		return "-1"
	case "float64":
		return "-1.0"
	case "bool":
		return "false"
	default:
		return "\"\""
	}
}

func getInvalidDescription(field Field) string {
	switch field.Type {
	case "string":
		return "empty string"
	case "int", "int64":
		return "negative number"
	case "float64":
		return "negative number"
	case "bool":
		return "false value"
	default:
		return "invalid value"
	}
}
