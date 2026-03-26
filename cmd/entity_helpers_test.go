package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetValidateTag(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		fieldName string
		fieldType string
		expected  string
	}{
		{name: "email field", fieldName: "Email", fieldType: "string", expected: "required,email"},
		{name: "email lower", fieldName: "email", fieldType: "string", expected: "required,email"},
		{name: "string field", fieldName: "Name", fieldType: "string", expected: "required"},
		{name: "int field", fieldName: "Age", fieldType: "int", expected: "required,gte=0"},
		{name: "int64 field", fieldName: "Count", fieldType: "int64", expected: "required,gte=0"},
		{name: "uint field", fieldName: "ID", fieldType: "uint", expected: "required,gte=0"},
		{name: "float64 field", fieldName: "Price", fieldType: "float64", expected: "required,gte=0"},
		{name: "bool field", fieldName: "Active", fieldType: "bool", expected: ""},
		{name: "time field", fieldName: "CreatedAt", fieldType: "time.Time", expected: "required"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, getValidateTag(tc.fieldName, tc.fieldType))
		})
	}
}

func TestGetGormTag(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		fieldName string
		fieldType string
		contains  string
	}{
		{name: "email unique", fieldName: "Email", fieldType: "string", contains: "uniqueIndex"},
		{name: "string varchar", fieldName: "Name", fieldType: "string", contains: "varchar"},
		{name: "int integer", fieldName: "Age", fieldType: "int", contains: "integer"},
		{name: "float decimal", fieldName: "Price", fieldType: "float64", contains: "decimal"},
		{name: "bool boolean", fieldName: "Active", fieldType: "bool", contains: "boolean"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tag := getGormTag(tc.fieldName, tc.fieldType)
			assert.Contains(t, tag, tc.contains)
		})
	}
}

func TestIsSystemField(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		field    string
		expected bool
	}{
		{name: "ID", field: "ID", expected: true},
		{name: "CreatedAt", field: "CreatedAt", expected: true},
		{name: "UpdatedAt", field: "UpdatedAt", expected: true},
		{name: "DeletedAt", field: "DeletedAt", expected: true},
		{name: "Name", field: "Name", expected: false},
		{name: "Email", field: "Email", expected: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, isSystemField(tc.field))
		})
	}
}

func TestHasStringBusinessRules(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		fields   []Field
		expected bool
	}{
		{name: "has email", fields: []Field{{Name: "Email", Type: "string"}}, expected: true},
		{name: "no email", fields: []Field{{Name: "Name", Type: "string"}, {Name: "Age", Type: "int"}}, expected: false},
		{name: "empty fields", fields: []Field{}, expected: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, hasStringBusinessRules(tc.fields))
		})
	}
}

func TestWriteEntityStruct(t *testing.T) {
	t.Parallel()

	fields := []Field{
		{Name: "ID", Type: "uint", Tag: "`json:\"id\" gorm:\"primaryKey\"`"},
		{Name: "Name", Type: "string", Tag: "`json:\"name\"`"},
	}

	var content strings.Builder
	writeEntityStruct(&content, "User", fields)
	result := content.String()

	assert.Contains(t, result, "type User struct")
	assert.Contains(t, result, "ID")
	assert.Contains(t, result, "Name")
	assert.Contains(t, result, "string")
}

func TestWriteEntityHeader(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		fields         []Field
		businessRules  bool
		timestamps     bool
		softDelete     bool
		expectContains []string
	}{
		{
			name:           "basic",
			fields:         []Field{{Name: "Name", Type: "string"}},
			businessRules:  false,
			timestamps:     false,
			softDelete:     false,
			expectContains: []string{"package domain"},
		},
		{
			name:           "with timestamps",
			fields:         []Field{{Name: "Name", Type: "string"}},
			timestamps:     true,
			expectContains: []string{"time"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var content strings.Builder
			writeEntityHeader(&content, tc.fields, tc.businessRules, tc.timestamps, tc.softDelete)
			result := content.String()
			for _, s := range tc.expectContains {
				assert.Contains(t, result, s)
			}
		})
	}
}

func TestWriteFieldValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		field    Field
		contains string
	}{
		{name: "string field", field: Field{Name: "Name", Type: "string"}, contains: "Name"},
		{name: "int field", field: Field{Name: "Age", Type: "int"}, contains: "Age"},
		{name: "email field", field: Field{Name: "Email", Type: "string"}, contains: "Email"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var content strings.Builder
			writeFieldValidation(&content, "u", "User", tc.field)
			result := content.String()
			assert.Contains(t, result, tc.contains)
		})
	}
}

func TestParseFields(t *testing.T) {
	t.Parallel()

	fields := parseFields("Name:string,Email:string,Age:int")
	require.GreaterOrEqual(t, len(fields), 3, "should have at least Name, Email, Age fields")

	// Find each expected field (ID may be auto-added)
	fieldNames := make(map[string]string)
	for _, f := range fields {
		fieldNames[f.Name] = f.Type
	}

	assert.Equal(t, "string", fieldNames["Name"])
	assert.Equal(t, "string", fieldNames["Email"])
	assert.Equal(t, "int", fieldNames["Age"])
}

func TestWriteValidationMethod(t *testing.T) {
	t.Parallel()

	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Email", Type: "string"},
	}

	var content strings.Builder
	writeValidationMethod(&content, "User", fields)
	result := content.String()

	assert.Contains(t, result, "func (u *User) Validate() error")
	assert.Contains(t, result, "Name")
	assert.Contains(t, result, "Email")
	// ID should not be validated
	assert.NotContains(t, result, fmt.Sprintf("u.ID"))
}
