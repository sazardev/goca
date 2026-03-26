package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSliceClosingBrace(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		content  string
		startPos int
		expected int
	}{
		{
			"simple brace",
			"{ }",
			1,
			2,
		},
		{
			"nested braces",
			"{ { } }",
			1,
			6,
		},
		{
			"no closing brace",
			"{ { ",
			1,
			-1,
		},
		{
			"entities slice pattern",
			`entities := []interface{}{
		&domain.Product{},
		&domain.User{},
	}`,
			26,
			67,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := findSliceClosingBrace(tc.content, tc.startPos)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsEntityInMigrationList(t *testing.T) {
	t.Parallel()

	content := `func main() {
	entities := []interface{}{
		&domain.Product{},
		// &domain.Old{},
		&domain.User{},
	}
}`

	t.Run("entity present", func(t *testing.T) {
		t.Parallel()
		assert.True(t, isEntityInMigrationList(content, "&domain.Product{}"))
	})

	t.Run("entity absent", func(t *testing.T) {
		t.Parallel()
		assert.False(t, isEntityInMigrationList(content, "&domain.Order{}"))
	})

	t.Run("entity in comment", func(t *testing.T) {
		t.Parallel()
		assert.False(t, isEntityInMigrationList(content, "&domain.Old{}"))
	})

	t.Run("no entities slice", func(t *testing.T) {
		t.Parallel()
		assert.False(t, isEntityInMigrationList("package main", "&domain.Product{}"))
	})
}

func TestAddEntityToEntitiesSlice(t *testing.T) {
	t.Parallel()

	t.Run("adds entity", func(t *testing.T) {
		t.Parallel()
		content := `entities := []interface{}{
		&domain.Product{},
	}`
		result, err := addEntityToEntitiesSlice(content, "&domain.User{}")
		assert.NoError(t, err)
		assert.Contains(t, result, "&domain.User{}")
		assert.Contains(t, result, "&domain.Product{}")
	})

	t.Run("no entities slice", func(t *testing.T) {
		t.Parallel()
		_, err := addEntityToEntitiesSlice("package main", "&domain.User{}")
		assert.Error(t, err)
	})
}
