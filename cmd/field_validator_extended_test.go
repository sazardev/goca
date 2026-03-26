package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFieldValidator(t *testing.T) {
	t.Parallel()
	v := NewFieldValidator()
	assert.NotNil(t, v)
}

func TestFieldValidator_ValidateFieldName(t *testing.T) {
	t.Parallel()
	v := NewFieldValidator()

	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "Name", false},
		{"valid snake_case", "user_name", false},
		{"empty", "", true},
		{"starts with number", "1name", true},
		{"single char valid", "a", false},
		{"contains dot", "user.name", true},
		{"contains dash", "user-name", true},
		{"contains space", "user name", true},
		{"valid long name", "myFieldName", false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := v.ValidateFieldName(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFieldValidator_ValidateField(t *testing.T) {
	t.Parallel()
	v := NewFieldValidator()

	t.Run("valid field", func(t *testing.T) {
		t.Parallel()
		field, err := v.ValidateField("name:string")
		require.NoError(t, err)
		assert.Equal(t, "Name", field.Name)
		assert.Equal(t, "string", field.Type)
	})

	t.Run("empty field", func(t *testing.T) {
		t.Parallel()
		_, err := v.ValidateField("")
		assert.Error(t, err)
	})

	t.Run("missing colon", func(t *testing.T) {
		t.Parallel()
		_, err := v.ValidateField("namestring")
		assert.Error(t, err)
	})

	t.Run("too many colons", func(t *testing.T) {
		t.Parallel()
		_, err := v.ValidateField("name:string:extra")
		assert.Error(t, err)
	})

	t.Run("invalid type", func(t *testing.T) {
		t.Parallel()
		_, err := v.ValidateField("name:!!!invalid")
		assert.Error(t, err)
	})
}

func TestFieldValidator_ValidateReservedNames(t *testing.T) {
	t.Parallel()
	v := NewFieldValidator()

	t.Run("non-reserved name", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, v.ValidateReservedNames("name"))
	})
}

func TestFieldValidator_ParseFieldsWithValidation(t *testing.T) {
	t.Parallel()
	v := NewFieldValidator()

	t.Run("simple fields", func(t *testing.T) {
		t.Parallel()
		fields, err := v.ParseFieldsWithValidation("Name:string,Age:int")
		require.NoError(t, err)
		// Should include ID + 2 fields
		assert.GreaterOrEqual(t, len(fields), 2)
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		_, err := v.ParseFieldsWithValidation("")
		assert.Error(t, err)
	})
}

func TestFieldValidator_SmartSplitFields(t *testing.T) {
	t.Parallel()
	v := NewFieldValidator()

	t.Run("simple split", func(t *testing.T) {
		t.Parallel()
		result := v.smartSplitFields("Name:string,Age:int")
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "Name:string", result[0])
		assert.Equal(t, "Age:int", result[1])
	})

	t.Run("with nested brackets", func(t *testing.T) {
		t.Parallel()
		result := v.smartSplitFields("Name:string,Tags:[]string")
		assert.Equal(t, 2, len(result))
	})

	t.Run("with parentheses", func(t *testing.T) {
		t.Parallel()
		result := v.smartSplitFields("Callback:func(string),Name:string")
		assert.Equal(t, 2, len(result))
	})
}

func TestFieldValidator_ValidateFieldType_Extended(t *testing.T) {
	t.Parallel()
	v := NewFieldValidator()

	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty", "", true},
		{"time.Time", "time.Time", false},
		{"custom type", "User", false},
		{"qualified type", "domain.Product", false},
		{"byte", "byte", false},
		{"rune", "rune", false},
		{"uintptr", "uintptr", false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := v.ValidateFieldType(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCapitalizeFirst(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected string
	}{
		{"name", "Name"},
		{"", ""},
		{"Name", "Name"},
		{"a", "A"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, capitalizeFirst(tc.input))
		})
	}
}
