package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUIRenderer(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	require.NotNil(t, renderer)
	assert.True(t, renderer.noColor)
	assert.True(t, renderer.interactive)
	assert.Equal(t, 1, renderer.verbosity)
}

func TestNewUIRenderer_NilWriter(t *testing.T) {
	t.Parallel()
	renderer := NewUIRenderer(nil, true, 0)
	require.NotNil(t, renderer)
	// Should default to os.Stdout
	assert.NotNil(t, renderer.writer)
}

func TestUIRenderer_SetInteractive(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	assert.True(t, renderer.IsInteractive())

	renderer.SetInteractive(false)
	assert.False(t, renderer.IsInteractive())
}

func TestUIRenderer_Success(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.Success("it worked")
	assert.Contains(t, buf.String(), "it worked")
}

func TestUIRenderer_Error(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.Error("something failed")
	assert.Contains(t, buf.String(), "something failed")
}

func TestUIRenderer_Warning(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.Warning("be careful")
	assert.Contains(t, buf.String(), "be careful")
}

func TestUIRenderer_Info(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.Info("note this")
	assert.Contains(t, buf.String(), "note this")
}

func TestUIRenderer_VerbositySuppression(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 0)

	renderer.Warning("hidden")
	renderer.Info("hidden")
	renderer.Header("hidden")
	renderer.Step(1, "hidden")
	renderer.Blank()
	renderer.Dim("hidden")
	renderer.Section("hidden")
	renderer.KeyValue("k", "v")
	renderer.Feature("f", false)

	assert.Empty(t, buf.String())

	// Success and Error always show regardless of verbosity
	renderer.Success("visible")
	assert.Contains(t, buf.String(), "visible")
	buf.Reset()
	renderer.Error("visible")
	assert.Contains(t, buf.String(), "visible")
}

func TestUIRenderer_Header(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.Header("My Header")
	assert.Contains(t, buf.String(), "My Header")
}

func TestUIRenderer_Step(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.Step(1, "first step")
	output := buf.String()
	assert.Contains(t, output, "1.")
	assert.Contains(t, output, "first step")
}

func TestUIRenderer_KeyValue(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.KeyValue("Name", "Goca")
	output := buf.String()
	assert.Contains(t, output, "Name")
	assert.Contains(t, output, "Goca")
}

func TestUIRenderer_FileCreated(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.FileCreated("/path/to/file.go")
	assert.Contains(t, buf.String(), "/path/to/file.go")
}

func TestUIRenderer_DryRun(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.DryRun("would create file")
	output := buf.String()
	assert.Contains(t, output, "DRY-RUN")
	assert.Contains(t, output, "would create file")
}

func TestUIRenderer_Println(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.Println("hello world")
	assert.Equal(t, "hello world\n", buf.String())
}

func TestUIRenderer_Printf(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.Printf("hello %s", "world")
	assert.Equal(t, "hello world", buf.String())
}

func TestUIRenderer_Table(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.Table(
		[]string{"File", "Status"},
		[][]string{
			{"test.go", "created"},
			{"main.go", "exists"},
		},
	)
	output := buf.String()
	assert.Contains(t, output, "File")
	assert.Contains(t, output, "test.go")
	assert.Contains(t, output, "main.go")
}

func TestUIRenderer_Table_Empty(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.Table([]string{}, nil)
	assert.Empty(t, buf.String())
}

func TestUIRenderer_NextSteps(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.NextSteps([]string{"Run build", "Run tests"})
	output := buf.String()
	assert.Contains(t, output, "Run build")
	assert.Contains(t, output, "Run tests")
}

func TestUIRenderer_Feature(t *testing.T) {
	t.Parallel()

	t.Run("without config", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		renderer := NewUIRenderer(&buf, true, 1)
		renderer.Feature("Including validation", false)
		assert.Contains(t, buf.String(), "Including validation")
		assert.NotContains(t, buf.String(), "from config")
	})

	t.Run("with config", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		renderer := NewUIRenderer(&buf, true, 1)
		renderer.Feature("Including validation", true)
		output := buf.String()
		assert.Contains(t, output, "Including validation")
		assert.Contains(t, output, "from config")
	})
}

func TestUIRenderer_Debug_Verbose(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 2)
	renderer.Debug("debug msg")
	assert.Contains(t, buf.String(), "debug msg")
}

func TestUIRenderer_Debug_Normal(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.Debug("debug msg")
	assert.Empty(t, buf.String())
}

func TestUIRenderer_Trace_Verbose(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 2)
	renderer.Trace("trace msg")
	assert.Contains(t, buf.String(), "trace msg")
}

func TestUIRenderer_Spinner_NoColor(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	stop := renderer.Spinner("loading")
	assert.Contains(t, buf.String(), "loading...")
	stop()
	assert.Contains(t, buf.String(), "done")
}

func TestUIRenderer_FileBackedUp(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.FileBackedUp("/old/path", "/new/path")
	output := buf.String()
	assert.Contains(t, output, "/old/path")
	assert.Contains(t, output, "/new/path")
}

func TestUIRenderer_KeyValueFromConfig(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	renderer := NewUIRenderer(&buf, true, 1)
	renderer.KeyValueFromConfig("DB", "postgres")
	output := buf.String()
	assert.Contains(t, output, "DB")
	assert.Contains(t, output, "postgres")
	assert.Contains(t, output, "from config")
}
