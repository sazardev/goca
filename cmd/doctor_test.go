package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountByStatus(t *testing.T) {
	t.Parallel()

	checks := []doctorCheck{
		{name: "a", status: "✓"},
		{name: "b", status: "✓"},
		{name: "c", status: "✗"},
		{name: "d", status: "⚠"},
		{name: "e", status: "✓"},
	}

	cases := []struct {
		status   string
		expected int
	}{
		{"✓", 3},
		{"✗", 1},
		{"⚠", 1},
		{"?", 0},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.status, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, countByStatus(checks, tc.status))
		})
	}
}

func TestCountByStatus_Empty(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 0, countByStatus(nil, "✓"))
	assert.Equal(t, 0, countByStatus([]doctorCheck{}, "✓"))
}
