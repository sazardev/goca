package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkipTestField(t *testing.T) {
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
		{name: "Age", field: "Age", expected: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, skipTestField(tc.field))
		})
	}
}

func TestTestLiteral(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		field    string
		typ      string
		entity   string
		expected string
	}{
		{name: "email", field: "Email", typ: "string", entity: "User", expected: `"test@user.com"`},
		{name: "name", field: "Name", typ: "string", entity: "Product", expected: `"Test Product"`},
		{name: "title", field: "Title", typ: "string", entity: "Post", expected: `"Test Post"`},
		{name: "description", field: "Description", typ: "string", entity: "Product", expected: `"A test product description"`},
		{name: "status", field: "Status", typ: "string", entity: "Order", expected: `"active"`},
		{name: "phone", field: "Phone", typ: "string", entity: "User", expected: `"+1234567890"`},
		{name: "address", field: "Address", typ: "string", entity: "User", expected: `"123 Test Street"`},
		{name: "price float", field: "Price", typ: "float64", entity: "Product", expected: "9.99"},
		{name: "age int", field: "Age", typ: "int", entity: "User", expected: "1"},
		{name: "quantity int", field: "Quantity", typ: "int", entity: "Item", expected: "1"},
		{name: "count uint", field: "Count", typ: "uint", entity: "Item", expected: "1"},
		{name: "active bool", field: "Active", typ: "bool", entity: "User", expected: "true"},
		{name: "timestamp", field: "Deadline", typ: "time.Time", entity: "Task", expected: "time.Now()"},
		{name: "generic string", field: "Bio", typ: "string", entity: "User", expected: `"test_bio"`},
		{name: "generic int64", field: "Score", typ: "int64", entity: "User", expected: "1"},
		{name: "generic float32", field: "Weight", typ: "float32", entity: "Item", expected: "9.99"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, testLiteral(tc.field, tc.typ, tc.entity))
		})
	}
}

func TestUpdatedTestLiteral(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		field    string
		typ      string
		entity   string
		expected string
	}{
		{name: "email", field: "Email", typ: "string", entity: "User", expected: `"updated@user.com"`},
		{name: "name", field: "Name", typ: "string", entity: "Product", expected: `"Updated Product"`},
		{name: "status", field: "Status", typ: "string", entity: "Order", expected: `"inactive"`},
		{name: "price", field: "Price", typ: "float64", entity: "Product", expected: "19.99"},
		{name: "age", field: "Age", typ: "int", entity: "User", expected: "2"},
		{name: "active", field: "Active", typ: "bool", entity: "User", expected: "false"},
		{name: "generic string", field: "Bio", typ: "string", entity: "User", expected: `"updated_bio"`},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, updatedTestLiteral(tc.field, tc.typ, tc.entity))
		})
	}
}

func TestNeedsFmtImport(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		fields   []Field
		expected bool
	}{
		{name: "has string field", fields: []Field{{Name: "Name", Type: "string"}}, expected: true},
		{name: "only int", fields: []Field{{Name: "Age", Type: "int"}}, expected: false},
		{name: "skip ID string", fields: []Field{{Name: "ID", Type: "string"}}, expected: false},
		{name: "mixed with string", fields: []Field{{Name: "ID", Type: "uint"}, {Name: "Email", Type: "string"}}, expected: true},
		{name: "empty", fields: []Field{}, expected: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, needsFmtImport(tc.fields))
		})
	}
}

func TestBuildTestFieldInit(t *testing.T) {
	t.Parallel()

	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Age", Type: "int"},
	}

	result := buildTestFieldInit(fields, "User", "\t\t\t")
	assert.Contains(t, result, `Name: "Test User",`)
	assert.Contains(t, result, "Age: 1,")
	assert.NotContains(t, result, "ID:")
}

func TestBuildTestFieldInitUpdated(t *testing.T) {
	t.Parallel()

	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Email", Type: "string"},
		{Name: "Price", Type: "float64"},
	}

	result := buildTestFieldInitUpdated(fields, "Product", "\t\t\t")
	assert.Contains(t, result, `Email: "updated@product.com",`)
	assert.Contains(t, result, "Price: 19.99,")
	assert.NotContains(t, result, "ID:")
}

func TestBuildTestFieldAssertions(t *testing.T) {
	t.Parallel()

	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Email", Type: "string"},
	}

	result := buildTestFieldAssertions(fields, "\t\t")
	assert.Contains(t, result, "assert.Equal(t, updateInput.Name, updated.Name)")
	assert.Contains(t, result, "assert.Equal(t, updateInput.Email, updated.Email)")
	assert.NotContains(t, result, "updateInput.ID")
}

func TestBuildTestFieldInitVaried(t *testing.T) {
	t.Parallel()

	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Age", Type: "int"},
	}

	result := buildTestFieldInitVaried(fields, "User", "\t\t\t\t", true)
	assert.Contains(t, result, "fmt.Sprintf")
	assert.Contains(t, result, "i + 1")
	assert.NotContains(t, result, "ID:")
}

func TestBuildFixtureOverrides(t *testing.T) {
	t.Parallel()

	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Age", Type: "int"},
	}

	result := buildFixtureOverrides(fields, "user", "\t")
	assert.Contains(t, result, `fields["name"].(string)`)
	assert.Contains(t, result, `fields["age"].(int)`)
	assert.Contains(t, result, "user.Name = v")
	assert.NotContains(t, result, `fields["id"]`)
}

func TestBuildFixtureVariation(t *testing.T) {
	t.Parallel()

	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Age", Type: "int"},
	}

	result := buildFixtureVariation(fields, "User", "\t\t", true)
	assert.Contains(t, result, "fmt.Sprintf")
	assert.Contains(t, result, "fixtures[i].Name")
	assert.Contains(t, result, "fixtures[i].Age = i + 1")
}

func TestGenerateIntegrationTestContent_WithFields(t *testing.T) {
	t.Parallel()

	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Email", Type: "string"},
		{Name: "Age", Type: "int"},
	}

	content := generateIntegrationTestContent("User", "postgres", false, fields)

	// Should have field initializers instead of TODOs
	assert.Contains(t, content, `Name: "Test User"`)
	assert.Contains(t, content, `Email: "test@user.com"`)
	assert.Contains(t, content, "Age: 1")
	// Should have fmt import for varied data
	assert.Contains(t, content, `"fmt"`)
	// Should have assertions
	assert.Contains(t, content, "assert.Equal(t, updateInput.Name, updated.Name)")
}

func TestGenerateIntegrationTestContent_WithoutFields(t *testing.T) {
	t.Parallel()

	content := generateIntegrationTestContent("User", "postgres", false, nil)

	// Should retain TODOs as guidance
	assert.Contains(t, content, "// TODO: Add fields")
}

func TestGenerateFixtureContent_WithFields(t *testing.T) {
	t.Parallel()

	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Price", Type: "float64"},
	}

	content := generateFixtureContent("Product", fields)

	assert.Contains(t, content, `Name: "Test Product"`)
	assert.Contains(t, content, "Price: 9.99")
	assert.Contains(t, content, `fields["name"].(string)`)
	assert.Contains(t, content, `fields["price"].(float64)`)
}

func TestGenerateFixtureContent_WithoutFields(t *testing.T) {
	t.Parallel()

	content := generateFixtureContent("Product", nil)
	assert.Contains(t, content, "// TODO: Add default field values")
}
