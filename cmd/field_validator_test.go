package cmd

import (
	"testing"
)

func TestFieldValidatorComplexTypes(t *testing.T) {
	validator := NewFieldValidator()

	// Test cases for complex field types
	testCases := []struct {
		name      string
		fieldType string
		shouldErr bool
	}{
		// Basic types
		{"basic string", "string", false},
		{"basic int", "int", false},
		{"basic bool", "bool", false},

		// Slice types
		{"string slice", "[]string", false},
		{"int slice", "[]int", false},
		{"pointer slice", "[]*User", false},
		{"nested slice", "[][]string", false},

		// Array types
		{"string array", "[10]string", false},
		{"int array", "[5]int", false},

		// Pointer types
		{"string pointer", "*string", false},
		{"int pointer", "*int", false},
		{"custom pointer", "*User", false},

		// Map types - valid
		{"string map", "map[string]string", false},
		{"int key map", "map[int]string", false},
		{"interface map", "map[string]interface{}", false},
		{"pointer value map", "map[string]*User", false},
		{"nested map", "map[string]map[int]string", false},

		// Map types - invalid keys (not comparable)
		{"slice key map", "map[[]string]int", true},
		{"map key map", "map[map[string]int]string", true},
		{"function key map", "map[func()]string", true},

		// Channel types
		{"basic channel", "chan string", false},
		{"receive channel", "<-chan int", false},
		{"send channel", "chan<- bool", false},

		// Function types
		{"simple function", "func()", false},
		{"function with params", "func(string) error", false},
		{"complex function", "func(int, string) (bool, error)", false},

		// Interface types
		{"empty interface", "interface{}", false},
		{"qualified interface", "io.Reader", false},

		// Qualified types
		{"time type", "time.Time", false},
		{"context type", "context.Context", false},

		// Custom types
		{"custom struct", "User", false},
		{"custom type", "ProductID", false},

		// Invalid types
		{"empty type", "", true},
		{"invalid syntax", "map[string", true},
		{"invalid array", "[abc]string", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateFieldType(tc.fieldType)

			if tc.shouldErr && err == nil {
				t.Errorf("Expected error for field type '%s', but got none", tc.fieldType)
			}

			if !tc.shouldErr && err != nil {
				t.Errorf("Expected no error for field type '%s', but got: %v", tc.fieldType, err)
			}
		})
	}
}

func TestFieldValidatorComplexFieldStrings(t *testing.T) {
	validator := NewFieldValidator()

	// Test complete field strings with complex types
	testCases := []struct {
		name      string
		fields    string
		shouldErr bool
	}{
		{
			"basic fields",
			"name:string,age:int,active:bool",
			false,
		},
		{
			"complex fields",
			"metadata:map[string]interface{},tags:[]string,owner:*User",
			false,
		},
		{
			"channel fields",
			"events:chan Event,results:<-chan Result",
			false,
		},
		{
			"function fields",
			"validator:func(interface{}) error,handler:func(string) (bool, error)",
			false,
		},
		{
			"nested complex types",
			"data:map[string][]interface{},cache:map[int]*[]User",
			false,
		},
		{
			"invalid map key",
			"invalid:map[[]string]int",
			true,
		},
		{
			"duplicate fields",
			"name:string,name:int",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateFields(tc.fields)

			if tc.shouldErr && err == nil {
				t.Errorf("Expected error for fields '%s', but got none", tc.fields)
			}

			if !tc.shouldErr && err != nil {
				t.Errorf("Expected no error for fields '%s', but got: %v", tc.fields, err)
			}
		})
	}
}
