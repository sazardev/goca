package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSafetyManager(t *testing.T) {
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

func TestSafetyManager_DryRun_WriteFile(t *testing.T) {
	t.Parallel()
	sm := NewSafetyManager(true, false, false)

	err := sm.WriteFile("/fake/path/test.go", "package main")
	require.NoError(t, err)

	pending := sm.GetPendingFiles()
	require.Len(t, pending, 1)
	assert.Equal(t, "/fake/path/test.go", pending[0].Path)
	assert.Equal(t, "create", pending[0].Action)
	assert.Equal(t, len("package main"), pending[0].Size)

	created := sm.GetCreatedFiles()
	assert.Contains(t, created, "/fake/path/test.go")
}

func TestSafetyManager_WriteFile_CreatesFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	sm := NewSafetyManager(false, false, false)

	filePath := filepath.Join(dir, "sub", "test.go")
	err := sm.WriteFile(filePath, "package test")
	require.NoError(t, err)

	content, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, "package test", string(content))
	assert.Contains(t, sm.GetCreatedFiles(), filePath)
}

func TestSafetyManager_CheckFileConflict_NoFile(t *testing.T) {
	t.Parallel()
	sm := NewSafetyManager(false, false, false)
	err := sm.CheckFileConflict("/nonexistent/path/file.go")
	assert.NoError(t, err)
}

func TestSafetyManager_CheckFileConflict_ExistingFile_NoForce(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	filePath := filepath.Join(dir, "existing.go")
	require.NoError(t, os.WriteFile(filePath, []byte("old"), 0644))

	sm := NewSafetyManager(false, false, false)
	err := sm.CheckFileConflict(filePath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestSafetyManager_CheckFileConflict_DryRun(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	filePath := filepath.Join(dir, "existing.go")
	require.NoError(t, os.WriteFile(filePath, []byte("old"), 0644))

	sm := NewSafetyManager(true, false, false)
	err := sm.CheckFileConflict(filePath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "dry-run")
	assert.Contains(t, sm.GetConflicts(), filePath)
}

func TestSafetyManager_CheckFileConflict_Force_WithBackup(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	filePath := filepath.Join(dir, "existing.go")
	require.NoError(t, os.WriteFile(filePath, []byte("old content"), 0644))

	sm := NewSafetyManager(false, true, true)
	sm.BackupDir = filepath.Join(dir, ".goca-backup")

	err := sm.CheckFileConflict(filePath)
	require.NoError(t, err)

	// Verify backup was created
	backupPath := filepath.Join(sm.BackupDir, "existing.go.backup")
	content, err := os.ReadFile(backupPath)
	require.NoError(t, err)
	assert.Equal(t, "old content", string(content))
}

func TestNewNameConflictDetector(t *testing.T) {
	t.Parallel()
	ncd := NewNameConflictDetector("/some/path")
	require.NotNil(t, ncd)
	assert.Equal(t, "/some/path", ncd.projectRoot)
}

func TestNameConflictDetector_ScanExistingEntities_NoDomainDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	ncd := NewNameConflictDetector(dir)
	err := ncd.ScanExistingEntities()
	assert.NoError(t, err)
}

func TestNameConflictDetector_ScanAndCheck(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	domainDir := filepath.Join(dir, "internal", "domain")
	require.NoError(t, os.MkdirAll(domainDir, 0755))

	// Create entity files
	require.NoError(t, os.WriteFile(filepath.Join(domainDir, "user.go"), []byte("package domain"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(domainDir, "product.go"), []byte("package domain"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(domainDir, "errors.go"), []byte("package domain"), 0644))

	ncd := NewNameConflictDetector(dir)
	require.NoError(t, ncd.ScanExistingEntities())

	// Existing entities should conflict
	assert.Error(t, ncd.CheckNameConflict("user"))
	assert.Error(t, ncd.CheckNameConflict("Product"))

	// Non-existing should not
	assert.NoError(t, ncd.CheckNameConflict("order"))

	// errors.go should be excluded
	assert.NoError(t, ncd.CheckNameConflict("errors"))

	entities := ncd.GetExistingEntities()
	assert.Contains(t, entities, "user")
	assert.Contains(t, entities, "product")
}
