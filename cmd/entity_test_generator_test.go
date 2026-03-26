package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValidFieldValue(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		field    Field
		expected string
	}{
		{"string email", Field{Name: "Email", Type: "string"}, "\"test@example.com\""},
		{"string name", Field{Name: "Name", Type: "string"}, "\"John Doe\""},
		{"string description", Field{Name: "Description", Type: "string"}, "\"A valid description\""},
		{"string status", Field{Name: "Status", Type: "string"}, "\"active\""},
		{"string other", Field{Name: "Title", Type: "string"}, "\"valid value\""},
		{"int age", Field{Name: "Age", Type: "int"}, "25"},
		{"int64 quantity", Field{Name: "Quantity", Type: "int64"}, "10"},
		{"int other", Field{Name: "Count", Type: "int"}, "1"},
		{"float64 price", Field{Name: "Price", Type: "float64"}, "99.99"},
		{"float64 total", Field{Name: "Total", Type: "float64"}, "100.50"},
		{"float64 other", Field{Name: "Rate", Type: "float64"}, "10.5"},
		{"bool", Field{Name: "Active", Type: "bool"}, "true"},
		{"time.Time", Field{Name: "Birthday", Type: "time.Time"}, "time.Now()"},
		{"unknown type", Field{Name: "Data", Type: "[]byte"}, "\"\""},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, getValidFieldValue(tc.field))
		})
	}
}

func TestGetInvalidFieldValue(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		field    Field
		expected string
	}{
		{"string", Field{Name: "Name", Type: "string"}, "\"\""},
		{"int", Field{Name: "Count", Type: "int"}, "-1"},
		{"int64", Field{Name: "ID", Type: "int64"}, "-1"},
		{"float64", Field{Name: "Price", Type: "float64"}, "-1.0"},
		{"bool", Field{Name: "Active", Type: "bool"}, "false"},
		{"unknown", Field{Name: "Data", Type: "[]byte"}, "\"\""},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, getInvalidFieldValue(tc.field))
		})
	}
}

func TestGetInvalidDescription(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		field    Field
		expected string
	}{
		{"string", Field{Name: "Name", Type: "string"}, "empty string"},
		{"int", Field{Name: "Count", Type: "int"}, "negative number"},
		{"int64", Field{Name: "ID", Type: "int64"}, "negative number"},
		{"float64", Field{Name: "Price", Type: "float64"}, "negative number"},
		{"bool", Field{Name: "Active", Type: "bool"}, "false value"},
		{"unknown", Field{Name: "Data", Type: "[]byte"}, "invalid value"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, getInvalidDescription(tc.field))
		})
	}
}

func TestGenerateValidationTests(t *testing.T) {
	t.Parallel()
	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Price", Type: "float64"},
	}

	var sb strings.Builder
	generateValidationTests(&sb, "Product", fields)
	result := sb.String()

	assert.Contains(t, result, "TestProduct_Validate")
	assert.Contains(t, result, "valid entity")
	assert.Contains(t, result, "wantErr")
	assert.Contains(t, result, "John Doe")
	assert.Contains(t, result, "99.99")
}

func TestGenerateConstructorTests(t *testing.T) {
	t.Parallel()
	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Age", Type: "int"},
	}

	var sb strings.Builder
	generateConstructorTests(&sb, "User", fields)
	result := sb.String()

	assert.Contains(t, result, "TestUser_Initialization")
	assert.Contains(t, result, "user := &User{")
	assert.Contains(t, result, "assert.Equal")
	assert.Contains(t, result, "John Doe")
}

func TestGenerateFieldTests(t *testing.T) {
	t.Parallel()
	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Price", Type: "float64"},
		{Name: "Active", Type: "bool"},
	}

	var sb strings.Builder
	generateFieldTests(&sb, "Product", fields)
	result := sb.String()

	assert.Contains(t, result, "TestProduct_Name_EdgeCases")
	assert.Contains(t, result, "TestProduct_Price_NumericValidation")
	assert.Contains(t, result, "empty string")
	assert.Contains(t, result, "positive value")
	assert.Contains(t, result, "negative value")
}

func TestGenerateInvalidTestCase(t *testing.T) {
	t.Parallel()
	allFields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Price", Type: "float64"},
	}
	invalidField := Field{Name: "Name", Type: "string"}

	var sb strings.Builder
	generateInvalidTestCase(&sb, "Product", "product", invalidField, allFields)
	result := sb.String()

	assert.Contains(t, result, "invalid product")
	assert.Contains(t, result, "Name: \"\"")
	assert.Contains(t, result, "wantErr: true")
}
