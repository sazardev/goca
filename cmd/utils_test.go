package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetImportPath(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		module   string
		expected string
	}{
		{name: "github module", module: "github.com/user/myproject", expected: "github.com/user/myproject"},
		{name: "test project module", module: "github.com/goca/testproject", expected: "github.com/goca/testproject"},
		{name: "simple module", module: "myproject", expected: "myproject"},
		{name: "gitlab module", module: "gitlab.com/org/project", expected: "gitlab.com/org/project"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, getImportPath(tc.module))
		})
	}
}

func TestIsSearchableField(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		fieldName string
		fieldType string
		expected  bool
	}{
		{name: "email string", fieldName: "Email", fieldType: "string", expected: true},
		{name: "username string", fieldName: "Username", fieldType: "string", expected: true},
		{name: "name string", fieldName: "Name", fieldType: "string", expected: true},
		{name: "code string", fieldName: "Code", fieldType: "string", expected: true},
		{name: "sku string", fieldName: "SKU", fieldType: "string", expected: true},
		{name: "slug string", fieldName: "Slug", fieldType: "string", expected: true},
		{name: "phone string", fieldName: "Phone", fieldType: "string", expected: true},
		{name: "title string", fieldName: "Title", fieldType: "string", expected: true},
		{name: "age int", fieldName: "Age", fieldType: "int", expected: true},
		{name: "count uint", fieldName: "Count", fieldType: "uint", expected: true},
		{name: "data bytes", fieldName: "Data", fieldType: "[]byte", expected: false},
		{name: "any interface", fieldName: "Meta", fieldType: "interface{}", expected: false},
		{name: "price string", fieldName: "Price", fieldType: "string", expected: true},
		{name: "active bool", fieldName: "Active", fieldType: "bool", expected: false},
		{name: "random float", fieldName: "Random", fieldType: "float64", expected: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, isSearchableField(tc.fieldName, tc.fieldType))
		})
	}
}

func TestIsUniqueField(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		fieldName string
		expected  bool
	}{
		{name: "email", fieldName: "Email", expected: true},
		{name: "username", fieldName: "Username", expected: true},
		{name: "code", fieldName: "Code", expected: true},
		{name: "sku", fieldName: "ProductSKU", expected: true},
		{name: "slug", fieldName: "Slug", expected: true},
		{name: "document", fieldName: "Document", expected: true},
		{name: "license", fieldName: "DriverLicense", expected: true},
		{name: "name not unique", fieldName: "Name", expected: false},
		{name: "age not unique", fieldName: "Age", expected: false},
		{name: "price not unique", fieldName: "Price", expected: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, isUniqueField(tc.fieldName))
		})
	}
}

func TestGenerateSearchMethods(t *testing.T) {
	t.Parallel()

	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Email", Type: "string"},
		{Name: "Name", Type: "string"},
		{Name: "Age", Type: "int"},
		{Name: "Data", Type: "[]byte"},
	}

	methods := generateSearchMethods(fields, "User")

	// ID should be skipped
	for _, m := range methods {
		assert.NotEqual(t, "FindByID", m.MethodName)
	}

	// Email should be searchable and unique
	var emailMethod *SearchMethod
	for i, m := range methods {
		if m.MethodName == "FindByEmail" {
			emailMethod = &methods[i]
			break
		}
	}
	require.NotNil(t, emailMethod, "FindByEmail should be generated")
	assert.True(t, emailMethod.IsUnique)
	assert.Equal(t, "string", emailMethod.FieldType)

	// Name should be searchable but not unique
	var nameMethod *SearchMethod
	for i, m := range methods {
		if m.MethodName == "FindByName" {
			nameMethod = &methods[i]
			break
		}
	}
	require.NotNil(t, nameMethod, "FindByName should be generated")
	assert.False(t, nameMethod.IsUnique)

	// Data ([]byte) should NOT have a search method
	for _, m := range methods {
		assert.NotEqual(t, "FindByData", m.MethodName)
	}
}
