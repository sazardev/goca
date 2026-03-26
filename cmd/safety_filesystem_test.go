package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- SafetyManager filesystem tests ---

func TestNewSafetyManager_Full(t *testing.T) {
	t.Parallel()
	sm := NewSafetyManager(true, false, true)
	assert.True(t, sm.DryRun)
	assert.False(t, sm.Force)
	assert.True(t, sm.Backup)
	assert.Equal(t, ".goca-backup", sm.BackupDir)
	assert.Empty(t, sm.GetConflicts())
	assert.Empty(t, sm.GetCreatedFiles())
	assert.Empty(t, sm.GetPendingFiles())
}

func TestSafetyManager_CheckFileConflict(t *testing.T) {
	t.Run("no conflict non-existent file", func(t *testing.T) {
		sm := NewSafetyManager(false, false, false)
		err := sm.CheckFileConflict(filepath.Join(t.TempDir(), "nonexistent.go"))
		assert.NoError(t, err)
	})

	t.Run("conflict dryrun", func(t *testing.T) {
		dir := t.TempDir()
		f := filepath.Join(dir, "existing.go")
		require.NoError(t, os.WriteFile(f, []byte("package x"), 0644))

		sm := NewSafetyManager(true, false, false)
		err := sm.CheckFileConflict(f)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "dry-run")
		assert.Len(t, sm.GetConflicts(), 1)
	})

	t.Run("conflict without force", func(t *testing.T) {
		dir := t.TempDir()
		f := filepath.Join(dir, "existing.go")
		require.NoError(t, os.WriteFile(f, []byte("package x"), 0644))

		sm := NewSafetyManager(false, false, false)
		err := sm.CheckFileConflict(f)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "--force")
	})

	t.Run("conflict with force", func(t *testing.T) {
		dir := t.TempDir()
		f := filepath.Join(dir, "existing.go")
		require.NoError(t, os.WriteFile(f, []byte("package x"), 0644))

		sm := NewSafetyManager(false, true, false)
		err := sm.CheckFileConflict(f)
		assert.NoError(t, err)
	})

	t.Run("conflict with force and backup", func(t *testing.T) {
		dir := t.TempDir()
		f := filepath.Join(dir, "existing.go")
		require.NoError(t, os.WriteFile(f, []byte("package x"), 0644))

		sm := NewSafetyManager(false, true, true)
		sm.BackupDir = filepath.Join(dir, ".backup")
		err := sm.CheckFileConflict(f)
		assert.NoError(t, err)
		// Verify backup was created
		_, err = os.Stat(filepath.Join(dir, ".backup", "existing.go.backup"))
		assert.NoError(t, err)
	})
}

func TestSafetyManager_BackupFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "myfile.go")
	require.NoError(t, os.WriteFile(f, []byte("original content"), 0644))

	sm := NewSafetyManager(false, false, true)
	sm.BackupDir = filepath.Join(dir, ".backup")
	err := sm.BackupFile(f)
	assert.NoError(t, err)

	backup := filepath.Join(dir, ".backup", "myfile.go.backup")
	content, err := os.ReadFile(backup)
	require.NoError(t, err)
	assert.Equal(t, "original content", string(content))
}

func TestSafetyManager_WriteFile_DryRun(t *testing.T) {
	t.Parallel()
	sm := NewSafetyManager(true, false, false)
	err := sm.WriteFile("/tmp/fake/file.go", "package main")
	assert.NoError(t, err)
	assert.Len(t, sm.GetPendingFiles(), 1)
	assert.Equal(t, "create", sm.GetPendingFiles()[0].Action)
	assert.Len(t, sm.GetCreatedFiles(), 1)
}

func TestSafetyManager_WriteFile_Real(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "sub", "file.go")

	sm := NewSafetyManager(false, false, false)
	err := sm.WriteFile(f, "package sub\n")
	assert.NoError(t, err)

	content, err := os.ReadFile(f)
	require.NoError(t, err)
	assert.Equal(t, "package sub\n", string(content))
	assert.Contains(t, sm.GetCreatedFiles(), f)
}

func TestSafetyManager_PrintSummary(t *testing.T) {
	// Not parallel: modifies global ui

	t.Run("styled dryrun", func(t *testing.T) {
		buf := &bytes.Buffer{}
		testUI := NewUIRenderer(buf, true, 2)
		oldUI := ui
		ui = testUI
		defer func() { ui = oldUI }()

		sm := NewSafetyManager(true, false, false)
		sm.pendingFiles = []DryRunEntry{{Path: "test.go", Action: "create", Size: 100}}
		sm.createdFiles = []string{"test.go"}
		sm.PrintSummary()
		out := buf.String()
		assert.Contains(t, out, "DRY-RUN")
		assert.Contains(t, out, "test.go")
	})

	t.Run("styled non-dryrun", func(t *testing.T) {
		buf := &bytes.Buffer{}
		testUI := NewUIRenderer(buf, true, 2)
		oldUI := ui
		ui = testUI
		defer func() { ui = oldUI }()

		sm := NewSafetyManager(false, false, false)
		sm.createdFiles = []string{"created.go"}
		sm.PrintSummary()
		out := buf.String()
		assert.Contains(t, out, "created.go")
	})

	t.Run("plain when ui nil", func(t *testing.T) {
		oldUI := ui
		ui = nil
		defer func() { ui = oldUI }()

		sm := NewSafetyManager(true, false, false)
		sm.createdFiles = []string{"test.go"}
		sm.conflicts = []string{"conflict.go"}
		sm.PrintSummary() // Should use printSummaryPlain without panic
	})
}

// --- NameConflictDetector tests ---

func TestNameConflictDetector_ScanExistingEntities_WithFiles(t *testing.T) {
	dir := t.TempDir()
	domainDir := filepath.Join(dir, "internal", "domain")
	require.NoError(t, os.MkdirAll(domainDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(domainDir, "product.go"), []byte("package domain"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(domainDir, "user.go"), []byte("package domain"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(domainDir, "errors.go"), []byte("package domain"), 0644))

	ncd := NewNameConflictDetector(dir)
	err := ncd.ScanExistingEntities()
	require.NoError(t, err)

	entities := ncd.GetExistingEntities()
	assert.Contains(t, entities, "product")
	assert.Contains(t, entities, "user")
	assert.NotContains(t, entities, "errors")
}

func TestNameConflictDetector_ScanNoDomainDir(t *testing.T) {
	dir := t.TempDir()
	ncd := NewNameConflictDetector(dir)
	err := ncd.ScanExistingEntities()
	assert.NoError(t, err) // No error when domain dir doesn't exist
	assert.Empty(t, ncd.GetExistingEntities())
}

func TestNameConflictDetector_CheckNameConflict(t *testing.T) {
	dir := t.TempDir()
	domainDir := filepath.Join(dir, "internal", "domain")
	require.NoError(t, os.MkdirAll(domainDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(domainDir, "product.go"), []byte("package domain"), 0644))

	ncd := NewNameConflictDetector(dir)
	require.NoError(t, ncd.ScanExistingEntities())

	assert.Error(t, ncd.CheckNameConflict("Product"))
	assert.Error(t, ncd.CheckNameConflict("product"))
	assert.NoError(t, ncd.CheckNameConflict("Order"))
}
