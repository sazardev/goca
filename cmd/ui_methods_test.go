package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for UI functions not yet fully covered

func TestUIRenderer_FileBackedUp_Full(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	u := NewUIRenderer(buf, true, 2)
	u.FileBackedUp("old.go", "old.go.bak")
	out := buf.String()
	assert.Contains(t, out, "old.go")
	assert.Contains(t, out, "old.go.bak")
}

func TestUIRenderer_KeyValueFromConfig_Full(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	u := NewUIRenderer(buf, true, 1)
	u.KeyValueFromConfig("handler", "http")
	out := buf.String()
	assert.Contains(t, out, "handler")
	assert.Contains(t, out, "http")
	assert.Contains(t, out, "from config")
}

func TestUIRenderer_Dim_Verbosity(t *testing.T) {
	t.Parallel()
	t.Run("verbosity 0 skips", func(t *testing.T) {
		t.Parallel()
		buf := &bytes.Buffer{}
		u := NewUIRenderer(buf, true, 0)
		u.Dim("subtle text")
		assert.Empty(t, buf.String())
	})
	t.Run("verbosity 1 shows", func(t *testing.T) {
		t.Parallel()
		buf := &bytes.Buffer{}
		u := NewUIRenderer(buf, true, 1)
		u.Dim("subtle text")
		assert.Contains(t, buf.String(), "subtle text")
	})
}

func TestUIRenderer_Section_Verbosity(t *testing.T) {
	t.Parallel()
	t.Run("verbosity 0 skips", func(t *testing.T) {
		t.Parallel()
		buf := &bytes.Buffer{}
		u := NewUIRenderer(buf, true, 0)
		u.Section("Title")
		assert.Empty(t, buf.String())
	})
	t.Run("verbosity 1 shows", func(t *testing.T) {
		t.Parallel()
		buf := &bytes.Buffer{}
		u := NewUIRenderer(buf, true, 1)
		u.Section("Title")
		assert.Contains(t, buf.String(), "Title")
	})
}

func TestUIRenderer_Blank_Full(t *testing.T) {
	t.Parallel()
	t.Run("verbosity 0 skips", func(t *testing.T) {
		t.Parallel()
		buf := &bytes.Buffer{}
		u := NewUIRenderer(buf, true, 0)
		u.Blank()
		assert.Empty(t, buf.String())
	})
	t.Run("verbosity 1 prints", func(t *testing.T) {
		t.Parallel()
		buf := &bytes.Buffer{}
		u := NewUIRenderer(buf, true, 1)
		u.Blank()
		assert.Equal(t, "\n", buf.String())
	})
}
