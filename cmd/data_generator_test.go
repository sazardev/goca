package cmd

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSeededGenerator(seed int64) *DataGenerator {
	return &DataGenerator{rand: rand.New(rand.NewSource(seed))}
}

func TestMin(t *testing.T) {
	t.Parallel()
	cases := []struct {
		a, b, want int
	}{
		{1, 2, 1},
		{5, 3, 3},
		{0, 0, 0},
		{-1, 1, -1},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, min(tc.a, tc.b))
	}
}

func TestFormatSQLValue(t *testing.T) {
	t.Parallel()
	g := newSeededGenerator(42)

	cases := []struct {
		name  string
		value interface{}
		want  string
	}{
		{"string", "hello", "'hello'"},
		{"string with quote", "it's", "'it''s'"},
		{"int", 42, "42"},
		{"int64", int64(99), "99"},
		{"uint", uint(7), "7"},
		{"float64", 3.14, "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"default", []byte("x"), "'[120]'"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := g.formatSQLValue(tc.value)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGenerateSampleData_SkipsID(t *testing.T) {
	t.Parallel()
	g := newSeededGenerator(1)
	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
	}
	data := g.GenerateSampleData(fields, "User")
	_, hasID := data["ID"]
	assert.False(t, hasID, "ID should be skipped")
	_, hasName := data["Name"]
	assert.True(t, hasName, "Name should be present")
}

func TestGenerateValueForField_ContextAware(t *testing.T) {
	t.Parallel()
	g := newSeededGenerator(1)

	t.Run("email field", func(t *testing.T) {
		t.Parallel()
		v := g.generateValueForField(Field{Name: "Email", Type: "string"}, "user")
		s, ok := v.(string)
		require.True(t, ok)
		assert.Contains(t, s, "@")
	})

	t.Run("age field", func(t *testing.T) {
		t.Parallel()
		v := g.generateValueForField(Field{Name: "Age", Type: "int"}, "user")
		age, ok := v.(int)
		require.True(t, ok)
		assert.True(t, age >= 18 && age <= 77)
	})

	t.Run("price field", func(t *testing.T) {
		t.Parallel()
		v := g.generateValueForField(Field{Name: "Price", Type: "float64"}, "product")
		_, ok := v.(float64)
		assert.True(t, ok)
	})

	t.Run("status field", func(t *testing.T) {
		t.Parallel()
		v := g.generateValueForField(Field{Name: "Status", Type: "string"}, "order")
		s, ok := v.(string)
		require.True(t, ok)
		assert.NotEmpty(t, s)
	})

	t.Run("phone field", func(t *testing.T) {
		t.Parallel()
		v := g.generateValueForField(Field{Name: "Phone", Type: "string"}, "user")
		s, ok := v.(string)
		require.True(t, ok)
		assert.True(t, strings.HasPrefix(s, "+34"))
	})

	t.Run("url field", func(t *testing.T) {
		t.Parallel()
		v := g.generateValueForField(Field{Name: "WebsiteURL", Type: "string"}, "company")
		s, ok := v.(string)
		require.True(t, ok)
		assert.True(t, strings.HasPrefix(s, "https://"))
	})

	t.Run("code field", func(t *testing.T) {
		t.Parallel()
		v := g.generateValueForField(Field{Name: "SKU", Type: "string"}, "product")
		s, ok := v.(string)
		require.True(t, ok)
		assert.Contains(t, s, "-")
	})

	t.Run("description field", func(t *testing.T) {
		t.Parallel()
		v := g.generateValueForField(Field{Name: "Description", Type: "string"}, "product")
		s, ok := v.(string)
		require.True(t, ok)
		assert.NotEmpty(t, s)
	})

	t.Run("address field", func(t *testing.T) {
		t.Parallel()
		v := g.generateValueForField(Field{Name: "Address", Type: "string"}, "user")
		s, ok := v.(string)
		require.True(t, ok)
		assert.Contains(t, s, "Madrid")
	})
}

func TestGenerateByType(t *testing.T) {
	t.Parallel()
	g := newSeededGenerator(1)

	cases := []struct {
		fieldType string
		checkType string
	}{
		{"string", "string"},
		{"int", "int"},
		{"int64", "int"},
		{"uint", "uint"},
		{"float64", "float64"},
		{"bool", "bool"},
		{"[]byte", "[]byte"},
		{"custom", "string"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.fieldType, func(t *testing.T) {
			t.Parallel()
			v := g.generateByType(tc.fieldType)
			assert.NotNil(t, v)
		})
	}
}

func TestGenerateTestData(t *testing.T) {
	t.Parallel()
	g := newSeededGenerator(42)
	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Email", Type: "string"},
	}
	dataset := g.GenerateTestData("User", fields, 5)
	assert.Len(t, dataset, 5)
	for _, row := range dataset {
		_, hasID := row["ID"]
		assert.False(t, hasID)
		assert.NotNil(t, row["Name"])
		assert.NotNil(t, row["Email"])
	}
}

func TestGenerateInsertSQL(t *testing.T) {
	t.Parallel()
	g := newSeededGenerator(42)
	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Price", Type: "float64"},
	}
	sql := g.GenerateInsertSQL("products", fields, 2)
	assert.Contains(t, sql, "INSERT INTO products")
	assert.Contains(t, sql, "name, price")
	lines := strings.Split(strings.TrimSpace(sql), "\n")
	// 1 comment line + 2 insert lines
	assert.Equal(t, 3, len(lines))
}

func TestGenerateName_Entities(t *testing.T) {
	t.Parallel()
	g := newSeededGenerator(1)

	t.Run("user entity", func(t *testing.T) {
		t.Parallel()
		v := g.generateName("user", "name")
		assert.NotEmpty(t, v)
	})

	t.Run("product name", func(t *testing.T) {
		t.Parallel()
		v := g.generateName("product", "name")
		assert.NotEmpty(t, v)
	})

	t.Run("project entity", func(t *testing.T) {
		t.Parallel()
		v := g.generateName("project", "name")
		assert.NotEmpty(t, v)
	})

	t.Run("order entity", func(t *testing.T) {
		t.Parallel()
		v := g.generateName("order", "name")
		assert.Contains(t, v, "Order")
	})

	t.Run("generic entity", func(t *testing.T) {
		t.Parallel()
		v := g.generateName("widget", "name")
		assert.Contains(t, v, "Sample")
	})
}

func TestGenerateStatus_Entities(t *testing.T) {
	t.Parallel()
	g := newSeededGenerator(1)

	entities := []string{"user", "order", "project", "widget"}
	for _, e := range entities {
		e := e
		t.Run(e, func(t *testing.T) {
			t.Parallel()
			v := g.generateStatus(e)
			assert.NotEmpty(t, v)
		})
	}
}

func TestGenerateDescription_Entities(t *testing.T) {
	t.Parallel()
	g := newSeededGenerator(1)

	t.Run("product", func(t *testing.T) {
		t.Parallel()
		v := g.generateDescription("product")
		assert.NotEmpty(t, v)
	})

	t.Run("project", func(t *testing.T) {
		t.Parallel()
		v := g.generateDescription("project")
		assert.NotEmpty(t, v)
	})

	t.Run("generic", func(t *testing.T) {
		t.Parallel()
		v := g.generateDescription("widget")
		assert.Contains(t, v, "widget")
	})
}
