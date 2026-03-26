package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsRequiredField(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		want bool
	}{
		{"Name", true},
		{"name", true},
		{"EMAIL", true},
		{"Title", true},
		{"ID", false},
		{"Description", false},
		{"Price", false},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, isRequiredField(tc.name))
		})
	}
}

func TestMakePlural(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input string
		want  string
	}{
		{"Category", "Categories"},
		{"Bus", "Buses"},
		{"Box", "Boxes"},
		{"Church", "Churches"},
		{"Dish", "Dishes"},
		{"Product", "Products"},
		{"User", "Users"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, makePlural(tc.input))
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		input []string
		want  []string
	}{
		{"no duplicates", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"with duplicates", []string{"time", "errors", "time"}, []string{"time", "errors"}},
		{"empty", []string{}, nil},
		{"all same", []string{"x", "x", "x"}, []string{"x"}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := removeDuplicates(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestNewTemplateGenerator(t *testing.T) {
	t.Parallel()
	gen := NewTemplateGenerator()
	require.NotNil(t, gen)
	assert.NotNil(t, gen.fieldValidator)
}

func TestGenerateImports(t *testing.T) {
	t.Parallel()
	gen := NewTemplateGenerator()

	t.Run("validation feature adds errors", func(t *testing.T) {
		t.Parallel()
		features := FeatureFlags{Validation: true}
		imports := gen.generateImports(features, nil)
		assert.Contains(t, imports, "errors")
	})

	t.Run("timestamps adds time", func(t *testing.T) {
		t.Parallel()
		features := FeatureFlags{Timestamps: true}
		imports := gen.generateImports(features, nil)
		assert.Contains(t, imports, "time")
	})

	t.Run("softDelete adds time", func(t *testing.T) {
		t.Parallel()
		features := FeatureFlags{SoftDelete: true}
		imports := gen.generateImports(features, nil)
		assert.Contains(t, imports, "time")
	})

	t.Run("time.Time field adds time", func(t *testing.T) {
		t.Parallel()
		fields := []FieldData{{Type: "time.Time"}}
		imports := gen.generateImports(FeatureFlags{}, fields)
		assert.Contains(t, imports, "time")
	})

	t.Run("[]byte field adds bytes", func(t *testing.T) {
		t.Parallel()
		fields := []FieldData{{Type: "[]byte"}}
		imports := gen.generateImports(FeatureFlags{}, fields)
		assert.Contains(t, imports, "bytes")
	})

	t.Run("deduplication", func(t *testing.T) {
		t.Parallel()
		features := FeatureFlags{Timestamps: true, SoftDelete: true}
		imports := gen.generateImports(features, nil)
		count := 0
		for _, imp := range imports {
			if imp == "time" {
				count++
			}
		}
		assert.Equal(t, 1, count)
	})
}

func TestGenerateMethods(t *testing.T) {
	t.Parallel()
	gen := NewTemplateGenerator()
	entity := EntityData{Name: "Product", NameLower: "product", NamePlural: "Products"}

	t.Run("validation feature adds Validate method", func(t *testing.T) {
		t.Parallel()
		features := FeatureFlags{Validation: true}
		methods := gen.generateMethods(entity, nil, features)
		assert.NotEmpty(t, methods)
		assert.Equal(t, "Validate", methods[0].Name)
	})

	t.Run("no features no methods", func(t *testing.T) {
		t.Parallel()
		methods := gen.generateMethods(entity, nil, FeatureFlags{})
		assert.Empty(t, methods)
	})
}

func TestGenerateValidations(t *testing.T) {
	t.Parallel()
	gen := NewTemplateGenerator()
	fields := []FieldData{
		{Name: "Name", IsRequired: true},
		{Name: "Price", IsRequired: false},
		{Name: "Email", IsRequired: true},
	}

	t.Run("with validation feature", func(t *testing.T) {
		t.Parallel()
		validations := gen.generateValidations(fields, FeatureFlags{Validation: true})
		assert.Len(t, validations, 2) // Name and Email
		assert.Equal(t, "Name", validations[0].Field)
		assert.Equal(t, "Email", validations[1].Field)
	})

	t.Run("without validation feature", func(t *testing.T) {
		t.Parallel()
		validations := gen.generateValidations(fields, FeatureFlags{})
		assert.Empty(t, validations)
	})
}

func TestGetTemplate(t *testing.T) {
	t.Parallel()
	gen := NewTemplateGenerator()

	validNames := []string{"entity", "usecase", "repository", "handler"}
	for _, name := range validNames {
		name := name
		t.Run("valid_"+name, func(t *testing.T) {
			t.Parallel()
			tmpl, err := gen.getTemplate(name)
			require.NoError(t, err)
			assert.NotNil(t, tmpl)
		})
	}

	t.Run("invalid template", func(t *testing.T) {
		t.Parallel()
		_, err := gen.getTemplate("nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
