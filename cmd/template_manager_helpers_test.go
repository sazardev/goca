package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCamelCase(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello", "hello"},
		{"hello_world", "helloWorld"},
		{"hello-world", "helloWorld"},
		{"hello world", "helloWorld"},
		{"HELLO_WORLD", "helloWorld"},
		{"my_long_variable_name", "myLongVariableName"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, toCamelCase(tc.input))
		})
	}
}

func TestToPascalCase(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello", "Hello"},
		{"hello_world", "HelloWorld"},
		{"hello-world", "HelloWorld"},
		{"hello world", "HelloWorld"},
		{"my_long_name", "MyLongName"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, toPascalCase(tc.input))
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello", "hello"},
		{"HelloWorld", "hello_world"},
		{"myLongName", "my_long_name"},
		{"ABC", "a_b_c"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, toSnakeCase(tc.input))
		})
	}
}

func TestToKebabCase(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"HelloWorld", "hello-world"},
		{"myLongName", "my-long-name"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, toKebabCase(tc.input))
		})
	}
}

func TestToPlural_TemplateManager(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"user", "users"},
		{"category", "categories"},
		{"bus", "buses"},
		{"box", "boxes"},
		{"match", "matches"},
		{"wish", "wishes"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, toPlural(tc.input))
		})
	}
}

func TestToSingular(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"users", "user"},
		{"categories", "category"},
		{"buses", "bus"},
		{"product", "product"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, toSingular(tc.input))
		})
	}
}

func TestExecuteTemplateString(t *testing.T) {
	t.Parallel()

	tm := NewTemplateManager(&TemplateConfig{Directory: "templates"}, "")

	t.Run("simple template", func(t *testing.T) {
		t.Parallel()
		result, err := tm.ExecuteTemplateString("Hello {{.Name}}", map[string]string{"Name": "World"})
		assert.NoError(t, err)
		assert.Equal(t, "Hello World", result)
	})

	t.Run("invalid template", func(t *testing.T) {
		t.Parallel()
		_, err := tm.ExecuteTemplateString("{{.Invalid", nil)
		assert.Error(t, err)
	})

	t.Run("template with functions", func(t *testing.T) {
		t.Parallel()
		result, err := tm.ExecuteTemplateString("{{toPlural \"user\"}}", nil)
		assert.NoError(t, err)
		assert.Equal(t, "users", result)
	})
}
