package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteSoftDeleteMethods(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	writeSoftDeleteMethods(&b, "Product")
	output := b.String()
	assert.Contains(t, output, "func (p *Product) SoftDelete()")
	assert.Contains(t, output, "func (p *Product) IsDeleted() bool")
	assert.Contains(t, output, "DeletedAt")
}

func TestGenerateBusinessRules(t *testing.T) {
	t.Parallel()

	t.Run("Age field", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateBusinessRules(&b, "User", []Field{{Name: "Age", Type: "int"}})
		assert.Contains(t, b.String(), "IsAdult")
	})

	t.Run("Price field", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateBusinessRules(&b, "Product", []Field{{Name: "Price", Type: "float64"}})
		assert.Contains(t, b.String(), "IsExpensive")
	})

	t.Run("Email field", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateBusinessRules(&b, "User", []Field{{Name: "Email", Type: "string"}})
		assert.Contains(t, b.String(), "HasValidEmail")
	})

	t.Run("Status field", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateBusinessRules(&b, "User", []Field{{Name: "Status", Type: "string"}})
		assert.Contains(t, b.String(), "IsActive")
	})

	t.Run("No rules for unknown field", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateBusinessRules(&b, "User", []Field{{Name: "Custom", Type: "string"}})
		assert.Empty(t, b.String())
	})
}

func TestContains(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		slice []string
		item  string
		want  bool
	}{
		{"found", []string{"a", "b", "c"}, "b", true},
		{"not found", []string{"a", "b"}, "c", false},
		{"empty slice", []string{}, "a", false},
		{"with whitespace", []string{"  a  "}, "a", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, contains(tc.slice, tc.item))
		})
	}
}

func TestGetFieldDisplayName(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input string
		want  string
	}{
		{"email", "email"},
		{"name", "name"},
		{"price", "price"},
		{"unknown_field", "unknown_field"},
		{"user_email_address", "email|address"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			result := getFieldDisplayName(tc.input)
			if strings.Contains(tc.want, "|") {
				options := strings.Split(tc.want, "|")
				assert.Contains(t, options, result)
			} else {
				assert.Equal(t, tc.want, result)
			}
		})
	}
}

func TestWriteErrorsHeader(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	writeErrorsHeader(&b)
	output := b.String()
	assert.Contains(t, output, "package domain")
	assert.Contains(t, output, "import \"errors\"")
	assert.Contains(t, output, "var (")
}

func TestWriteGeneralError(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	writeGeneralError(&b, "User", nil)
	assert.Contains(t, b.String(), "ErrInvalidUserData")
}

func TestWriteGeneralError_AlreadyExists(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	existing := []string{"\tErrInvalidUserData = errors.New(\"invalid user data\")"}
	writeGeneralError(&b, "User", existing)
	assert.Empty(t, b.String())
}

func TestWriteExistingErrors(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	existing := []string{"\tErrInvalidFoo = errors.New(\"foo\")", "\tErrInvalidBar = errors.New(\"bar\")"}
	writeExistingErrors(&b, existing)
	assert.Contains(t, b.String(), "ErrInvalidFoo")
	assert.Contains(t, b.String(), "ErrInvalidBar")
}

func TestWriteEntityErrors(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	fields := []Field{
		{Name: "Name", Type: "string"},
		{Name: "Age", Type: "int"},
	}
	writeEntityErrors(&b, "User", fields, nil)
	output := b.String()
	assert.Contains(t, output, "ErrInvalidUserData")
	assert.Contains(t, output, "ErrInvalidUserName")
	assert.Contains(t, output, "ErrInvalidUserAge")
	assert.Contains(t, output, ")")
}

func TestWriteFieldErrors(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Email", Type: "string"},
		{Name: "Price", Type: "float64"},
	}
	writeFieldErrors(&b, "Product", fields, nil)
	output := b.String()
	assert.NotContains(t, output, "ID") // System field skipped
	assert.Contains(t, output, "Email")
	assert.Contains(t, output, "Price")
}

func TestWriteRequiredFieldError(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	writeRequiredFieldError(&b, "User", Field{Name: "Name", Type: "string"}, nil)
	assert.Contains(t, b.String(), "ErrInvalidUserName")
	assert.Contains(t, b.String(), "name is required")
}

func TestWriteTypeSpecificErrors(t *testing.T) {
	t.Parallel()

	t.Run("string field", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		writeTypeSpecificErrors(&b, "User", Field{Name: "Name", Type: "string"}, nil)
		assert.Contains(t, b.String(), "Length")
	})

	t.Run("int field", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		writeTypeSpecificErrors(&b, "User", Field{Name: "Age", Type: "int"}, nil)
		assert.Contains(t, b.String(), "Range")
	})

	t.Run("float field", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		writeTypeSpecificErrors(&b, "Product", Field{Name: "Price", Type: "float64"}, nil)
		assert.Contains(t, b.String(), "Range")
	})
}

func TestWriteStringFieldErrors(t *testing.T) {
	t.Parallel()

	t.Run("name field only", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		writeStringFieldErrors(&b, "User", Field{Name: "UserName", Type: "string"}, "username", nil)
		assert.Contains(t, b.String(), "Length")
	})

	t.Run("name field", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		writeStringFieldErrors(&b, "User", Field{Name: "Name", Type: "string"}, "name", nil)
		assert.Contains(t, b.String(), "Length")
	})
}

func TestWriteIntegerFieldErrors(t *testing.T) {
	t.Parallel()

	t.Run("age field", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		writeIntegerFieldErrors(&b, "User", Field{Name: "Age", Type: "int"}, "age", nil)
		assert.Contains(t, b.String(), "greater than 0")
	})

	t.Run("generic int field", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		writeIntegerFieldErrors(&b, "User", Field{Name: "Count", Type: "int"}, "count", nil)
		assert.Contains(t, b.String(), "positive number")
	})
}

func TestWriteFloatFieldErrors(t *testing.T) {
	t.Parallel()

	t.Run("price field", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		writeFloatFieldErrors(&b, "Product", Field{Name: "Price", Type: "float64"}, "price", nil)
		assert.Contains(t, b.String(), "999,999,999.99")
	})

	t.Run("generic float field", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		writeFloatFieldErrors(&b, "Product", Field{Name: "Rate", Type: "float64"}, "rate", nil)
		assert.Contains(t, b.String(), "positive number")
	})
}

func TestGetNonSystemFieldNames(t *testing.T) {
	t.Parallel()
	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "CreatedAt", Type: "time.Time"},
		{Name: "Email", Type: "string"},
	}
	result := getNonSystemFieldNames(fields)
	assert.Equal(t, []string{"name", "email"}, result)
}

func TestGenerateSampleValue(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		field     Field
		index     int
		checkFunc func(t *testing.T, result string)
	}{
		{"string field", Field{Name: "Name", Type: "string"}, 1, func(t *testing.T, r string) {
			assert.Contains(t, r, "John Smith")
		}},
		{"int field", Field{Name: "Age", Type: "int"}, 1, func(t *testing.T, r string) {
			assert.Equal(t, "25", r)
		}},
		{"float field", Field{Name: "Price", Type: "float64"}, 1, func(t *testing.T, r string) {
			assert.Equal(t, "99.99", r)
		}},
		{"bool field", Field{Name: "Active", Type: "bool"}, 1, func(t *testing.T, r string) {
			assert.Equal(t, "true", r)
		}},
		{"time field", Field{Name: "Start", Type: "time.Time"}, 1, func(t *testing.T, r string) {
			assert.Equal(t, "time.Now()", r)
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, _ := generateSampleValue(tc.field, tc.index)
			tc.checkFunc(t, result)
		})
	}
}

func TestGenerateStringSampleValue(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		index    int
		contains string
	}{
		{"Name", 1, "John Smith"},
		{"Email", 1, "Sample"},
		{"Description", 1, "Detailed"},
		{"Title", 1, "Main Title"},
		{"Status", 1, "active"},
		{"Category", 1, "technology"},
		{"Custom", 1, "Sample"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := generateStringSampleValue(tc.name, tc.index)
			assert.Contains(t, result, tc.contains)
		})
	}
}

func TestGenerateIntSampleValue(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "25", generateIntSampleValue("Age", 1))
	assert.Equal(t, "100", generateIntSampleValue("Stock", 1))
	assert.Equal(t, "10", generateIntSampleValue("Quantity", 1))
	assert.Equal(t, "10", generateIntSampleValue("Other", 1))
}

func TestGenerateFloatSampleValue(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "99.99", generateFloatSampleValue("Price", 1))
	assert.Equal(t, "1000.00", generateFloatSampleValue("Amount", 1))
	assert.Equal(t, "10.50", generateFloatSampleValue("Other", 1))
}

func TestGenerateDefaultSampleValue(t *testing.T) {
	t.Parallel()

	intVal, ok := generateDefaultSampleValue("int", 1)
	assert.True(t, ok)
	assert.Contains(t, intVal, "10")

	strVal, ok := generateDefaultSampleValue("string", 1)
	assert.True(t, ok)
	assert.Contains(t, strVal, "Sample")

	floatVal, ok := generateDefaultSampleValue("float", 1)
	assert.True(t, ok)
	assert.Contains(t, floatVal, "10.5")

	boolVal, ok := generateDefaultSampleValue("bool", 1)
	assert.True(t, ok)
	assert.Contains(t, boolVal, "true")

	// Slices and maps produce valid composite literals.
	sliceVal, ok := generateDefaultSampleValue("[]string", 1)
	assert.True(t, ok)
	assert.Contains(t, sliceVal, "[]string{")

	mapVal, ok := generateDefaultSampleValue("map[string]int", 1)
	assert.True(t, ok)
	assert.Contains(t, mapVal, "map[string]int{")

	// Unknown named types are reported as unusable so callers can omit them.
	_, ok = generateDefaultSampleValue("custom", 1)
	assert.False(t, ok)
}

func TestGenerateSQLSampleValue(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		field    Field
		index    int
		contains string
	}{
		{"string name", Field{Name: "Name", Type: "string"}, 1, "'John Smith'"},
		{"string email", Field{Name: "Email", Type: "string"}, 1, "'john@example.com'"},
		{"string desc", Field{Name: "Description", Type: "string"}, 1, "'Detailed"},
		{"string status", Field{Name: "Status", Type: "string"}, 1, "'active'"},
		{"string other", Field{Name: "Code", Type: "string"}, 1, "'Sample"},
		{"int age", Field{Name: "Age", Type: "int"}, 1, "25"},
		{"int stock", Field{Name: "Stock", Type: "int"}, 1, "100"},
		{"int other", Field{Name: "Count", Type: "int"}, 1, "10"},
		{"float price", Field{Name: "Price", Type: "float64"}, 1, "99.99"},
		{"float other", Field{Name: "Rate", Type: "float64"}, 1, "10.50"},
		{"bool", Field{Name: "Active", Type: "bool"}, 1, "true"},
		{"time", Field{Name: "Start", Type: "time.Time"}, 1, "NOW()"},
		// ENTITY-7: composite/custom types emit a non-NULL value consistent with
		// the Go seeds and valid for NOT NULL columns.
		{"slice string", Field{Name: "Tags", Type: "[]string"}, 1, "'sample1,sample2'"},
		{"slice int", Field{Name: "Codes", Type: "[]int"}, 1, "'10,11'"},
		{"pointer", Field{Name: "LastLogin", Type: "*time.Time"}, 1, "NOW()"},
		{"custom type", Field{Name: "Status", Type: "UserStatus"}, 1, "'sample1'"},
		{"unknown", Field{Name: "X", Type: "complex"}, 1, "'sample1'"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := generateSQLSampleValue(tc.field, tc.index)
			assert.Contains(t, result, tc.contains)
		})
	}
}

func TestGetSQLFieldValues(t *testing.T) {
	t.Parallel()
	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Price", Type: "float64"},
	}
	values := getSQLFieldValues(fields, 1)
	assert.Len(t, values, 2) // ID excluded
}

func TestWriteSeedFileHeader(t *testing.T) {
	t.Parallel()

	t.Run("body without time usage", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		writeSeedFileHeader(&b, "func GetXSeeds() []X { return []X{{Name: \"a\"}} }")
		assert.Contains(t, b.String(), "package domain")
		assert.NotContains(t, b.String(), "import \"time\"")
	})

	t.Run("body using time", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		writeSeedFileHeader(&b, "func GetXSeeds() []X { return []X{{Start: time.Now()}} }")
		assert.Contains(t, b.String(), "import \"time\"")
	})
}

func TestWriteGoSeeds(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
	}
	writeGoSeeds(&b, "User", fields)
	output := b.String()
	assert.Contains(t, output, "GetUserSeeds")
	assert.Contains(t, output, "[]User")
	// Should have 3 records
	assert.Equal(t, 3, strings.Count(output, "Name:"))
}

func TestWriteGoSeedRecord(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Email", Type: "string"},
	}
	writeGoSeedRecord(&b, fields, 1)
	output := b.String()
	assert.NotContains(t, output, "ID:")
	assert.Contains(t, output, "Name:")
	assert.Contains(t, output, "Email:")
}

func TestWriteSQLSeeds(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
	}
	writeSQLSeeds(&b, "User", fields)
	output := b.String()
	assert.Contains(t, output, "GetSQLUserSeeds")
	assert.Contains(t, output, "INSERT INTO users")
}

func TestWriteSQLInsertStatement(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	fields := []Field{
		{Name: "ID", Type: "uint"},
		{Name: "Name", Type: "string"},
		{Name: "Age", Type: "int"},
	}
	writeSQLInsertStatement(&b, "User", fields, 1)
	output := b.String()
	assert.Contains(t, output, "INSERT INTO users")
	assert.Contains(t, output, "name, age")
	assert.Contains(t, output, "VALUES")
}
