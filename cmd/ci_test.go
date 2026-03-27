package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// ci_templates.go — template output tests
// ---------------------------------------------------------------------------

func TestGenerateTestWorkflow_NoDatabase(t *testing.T) {
	t.Parallel()
	data := CITemplateData{
		ProjectName: "myapp",
		Module:      "github.com/user/myapp",
		GoVersion:   "1.25",
	}
	out := generateTestWorkflow(data)
	assert.Contains(t, out, "name: Test")
	assert.Contains(t, out, "go-version: '1.25'")
	assert.Contains(t, out, "go test -race")
	assert.Contains(t, out, "go vet ./...")
	assert.Contains(t, out, "go build ./...")
	assert.NotContains(t, out, "postgres")
	assert.NotContains(t, out, "mysql")
}

func TestGenerateTestWorkflow_Postgres(t *testing.T) {
	t.Parallel()
	data := CITemplateData{
		ProjectName: "myapp",
		GoVersion:   "1.25",
		Database:    "postgres",
	}
	out := generateTestWorkflow(data)
	assert.Contains(t, out, "postgres:16")
	assert.Contains(t, out, "POSTGRES_USER")
	assert.Contains(t, out, "DATABASE_URL")
}

func TestGenerateTestWorkflow_PostgresJSON(t *testing.T) {
	t.Parallel()
	data := CITemplateData{GoVersion: "1.25", Database: "postgres-json"}
	out := generateTestWorkflow(data)
	assert.Contains(t, out, "postgres:16")
}

func TestGenerateTestWorkflow_MySQL(t *testing.T) {
	t.Parallel()
	data := CITemplateData{GoVersion: "1.25", Database: "mysql"}
	out := generateTestWorkflow(data)
	assert.Contains(t, out, "mysql:8")
	assert.Contains(t, out, "MYSQL_ROOT_PASSWORD")
	assert.Contains(t, out, "DATABASE_URL")
}

func TestGenerateBuildWorkflow_NoDocker(t *testing.T) {
	t.Parallel()
	data := CITemplateData{
		ProjectName: "myapp",
		GoVersion:   "1.25",
	}
	out := generateBuildWorkflow(data)
	assert.Contains(t, out, "name: Build")
	assert.Contains(t, out, "go-version: '1.25'")
	assert.Contains(t, out, "CGO_ENABLED=0")
	assert.Contains(t, out, "bin/myapp")
	assert.NotContains(t, out, "docker")
}

func TestGenerateBuildWorkflow_WithDocker(t *testing.T) {
	t.Parallel()
	data := CITemplateData{
		ProjectName: "myapp",
		GoVersion:   "1.25",
		WithDocker:  true,
	}
	out := generateBuildWorkflow(data)
	assert.Contains(t, out, "Docker Buildx")
	assert.Contains(t, out, "docker build")
	assert.Contains(t, out, "myapp:latest")
}

func TestGenerateDeployWorkflow(t *testing.T) {
	t.Parallel()
	data := CITemplateData{
		ProjectName: "myapp",
		GoVersion:   "1.25",
	}
	out := generateDeployWorkflow(data)
	assert.Contains(t, out, "name: Deploy")
	assert.Contains(t, out, "tags:")
	assert.Contains(t, out, "'v*'")
	assert.Contains(t, out, "bin/myapp")
	assert.Contains(t, out, "myapp-${{ github.ref_name }}")
}

// ---------------------------------------------------------------------------
// ci_helpers.go — helper function tests
// ---------------------------------------------------------------------------

func TestDetectGoVersionFromMod_ValidFile(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)
	require.NoError(t, os.WriteFile("go.mod", []byte("module github.com/test/app\n\ngo 1.25\n"), 0644))

	v := detectGoVersionFromMod()
	assert.Equal(t, "1.25", v)
}

func TestDetectGoVersionFromMod_NoFile(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	v := detectGoVersionFromMod()
	assert.Equal(t, "1.25", v)
}

func TestDetectGoVersionFromMod_NoDirective(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)
	require.NoError(t, os.WriteFile("go.mod", []byte("module github.com/test/app\n"), 0644))

	v := detectGoVersionFromMod()
	assert.Equal(t, "1.25", v)
}

func TestBuildCITemplateData_DefaultGoVersion(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)
	require.NoError(t, os.WriteFile("go.mod", []byte("module github.com/test/app\n\ngo 1.24\n"), 0644))

	data := buildCITemplateData("")
	assert.Equal(t, "1.24", data.GoVersion)
	assert.Equal(t, "github.com/test/app", data.Module)
}

func TestBuildCITemplateData_ExplicitGoVersion(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)
	require.NoError(t, os.WriteFile("go.mod", []byte("module github.com/test/app\n\ngo 1.24\n"), 0644))

	data := buildCITemplateData("1.22")
	assert.Equal(t, "1.22", data.GoVersion, "explicit version should override go.mod")
}

func TestGenerateCIPipeline_UnsupportedProvider(t *testing.T) {
	err := generateCIPipeline("jenkins", CITemplateData{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported CI provider")
}

func TestGenerateGitHubActions_DryRun(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	cleanup := ensureTestUI(t)
	defer cleanup()

	sm := NewSafetyManager(true, false, false)
	data := CITemplateData{
		ProjectName: "myapp",
		GoVersion:   "1.25",
	}

	err := generateGitHubActions(data, sm)
	assert.NoError(t, err)

	pending := sm.GetPendingFiles()
	assert.Len(t, pending, 2)
	assert.Contains(t, pending[0].Path, "test.yml")
	assert.Contains(t, pending[1].Path, "build.yml")
}

func TestGenerateGitHubActions_DryRun_WithDeploy(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	cleanup := ensureTestUI(t)
	defer cleanup()

	sm := NewSafetyManager(true, false, false)
	data := CITemplateData{
		ProjectName: "myapp",
		GoVersion:   "1.25",
		WithDeploy:  true,
	}

	err := generateGitHubActions(data, sm)
	assert.NoError(t, err)

	pending := sm.GetPendingFiles()
	assert.Len(t, pending, 3)
	assert.Contains(t, pending[2].Path, "deploy.yml")
}

func TestGenerateGitHubActions_RealFiles(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dir := t.TempDir()
	os.Chdir(dir)

	cleanup := ensureTestUI(t)
	defer cleanup()

	data := CITemplateData{
		ProjectName: "testproject",
		GoVersion:   "1.25",
		Database:    "postgres",
		WithDeploy:  true,
		WithDocker:  true,
	}

	err := generateGitHubActions(data)
	assert.NoError(t, err)

	// Verify test.yml
	content, err := os.ReadFile(filepath.Join(".github", "workflows", "test.yml"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "postgres:16")
	assert.Contains(t, string(content), "go test -race")

	// Verify build.yml
	content, err = os.ReadFile(filepath.Join(".github", "workflows", "build.yml"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "docker build")

	// Verify deploy.yml
	content, err = os.ReadFile(filepath.Join(".github", "workflows", "deploy.yml"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "tags:")
}
