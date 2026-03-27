package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── helpers ─────────────────────────────────────────────────────────────────

// setupAnalyzeProject creates a minimal Clean Architecture project in a temp dir
// and returns a teardown function that restores cwd.
func setupAnalyzeProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	dirs := []string{
		"internal/domain",
		"internal/usecase",
		"internal/repository",
		"internal/handler",
		"internal/di",
		"internal/mocks",
	}
	for _, d := range dirs {
		require.NoError(t, os.MkdirAll(filepath.Join(dir, d), 0755))
	}

	// go.mod
	writeTestFile(t, dir, "go.mod", "module github.com/test/myapp\n\ngo 1.22\n")
	// go.sum (presence check)
	writeTestFile(t, dir, "go.sum", "")
	// .goca.yaml
	writeTestFile(t, dir, ".goca.yaml", "project: myapp\nmodule: github.com/test/myapp\n")
	// main.go
	writeTestFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")

	// domain entity
	writeTestFile(t, dir, "internal/domain/user.go", `package domain

import "context"

// User is the core domain entity.
type User struct {
	ID    uint
	Name  string
	Email string
}

// Validate validates the user.
func (u *User) Validate(ctx context.Context) error {
	if u.Name == "" {
		return ErrUserNameRequired
	}
	return nil
}

var ErrUserNameRequired = errorf("user name is required")

func errorf(s string) error { return &domainError{s} }
type domainError struct{ msg string }
func (e *domainError) Error() string { return e.msg }
`)

	// usecase
	writeTestFile(t, dir, "internal/usecase/user_service.go", `package usecase

import "context"

// UserService handles user business logic.
type UserService struct{}

// CreateUser creates a new user.
func (s *UserService) CreateUser(ctx context.Context) error { return nil }
`)

	// repository
	writeTestFile(t, dir, "internal/repository/user_repository.go", `package repository

import "context"

// UserRepository persists users.
type UserRepository struct{}

// Save stores a user.
func (r *UserRepository) Save(ctx context.Context) error { return nil }
`)

	// di container
	writeTestFile(t, dir, "internal/di/container.go", `package di

// Container holds all dependencies.
type Container struct{}
`)

	// mock
	writeTestFile(t, dir, "internal/mocks/user_mock.go", `package mocks

// MockUserRepository is a testify mock.
type MockUserRepository struct{}
`)

	// test file using table-driven
	writeTestFile(t, dir, "internal/domain/user_test.go", `package domain

import "testing"

func TestUser_Validate(t *testing.T) {
	cases := []struct {
		name    string
		input   User
		wantErr bool
	}{
		{"valid", User{Name: "A"}, false},
		{"empty name", User{}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantErr {
				t.Log("expected error")
			}
		})
	}
}
`)

	old, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	return dir
}

func writeTestFile(t *testing.T, base, rel, content string) {
	t.Helper()
	path := filepath.Join(base, rel)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
}

// ─── resolveAnalyzeCategories ─────────────────────────────────────────────────

func TestResolveAnalyzeCategories_AllWhenNoneSet(t *testing.T) {
	t.Parallel()
	opts := analyzeOptions{}
	result := resolveAnalyzeCategories(opts)
	assert.True(t, result.architecture)
	assert.True(t, result.quality)
	assert.True(t, result.security)
	assert.True(t, result.standards)
	assert.True(t, result.tests)
	assert.True(t, result.deps)
}

func TestResolveAnalyzeCategories_KeepsExplicit(t *testing.T) {
	t.Parallel()
	opts := analyzeOptions{security: true}
	result := resolveAnalyzeCategories(opts)
	assert.True(t, result.security)
	assert.False(t, result.architecture)
	assert.False(t, result.quality)
}

// ─── analyzeCountByStatus ─────────────────────────────────────────────────────

func TestAnalyzeCountByStatus(t *testing.T) {
	t.Parallel()
	results := []analyzeResult{
		{status: "✓"},
		{status: "✓"},
		{status: "✗"},
		{status: "⚠"},
	}
	assert.Equal(t, 2, analyzeCountByStatus(results, "✓"))
	assert.Equal(t, 1, analyzeCountByStatus(results, "✗"))
	assert.Equal(t, 1, analyzeCountByStatus(results, "⚠"))
	assert.Equal(t, 0, analyzeCountByStatus(results, "?"))
}

// ─── analyzeResultsByCategory ─────────────────────────────────────────────────

func TestAnalyzeResultsByCategory(t *testing.T) {
	t.Parallel()
	results := []analyzeResult{
		{category: "Architecture", rule: "a"},
		{category: "Security", rule: "b"},
		{category: "Architecture", rule: "c"},
	}
	m := analyzeResultsByCategory(results)
	assert.Len(t, m["Architecture"], 2)
	assert.Len(t, m["Security"], 1)
}

// ─── Architecture checks ──────────────────────────────────────────────────────

func TestCheckLayerDirsExist_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkLayerDirsExist()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckLayerDirsExist_Fail(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkLayerDirsExist()
	assert.Equal(t, "✗", r.status)
}

func TestCheckDomainHasNoExternalImports_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkDomainHasNoExternalImports()
	assert.NotEqual(t, "✗", r.status, r.message)
}

func TestCheckDomainHasNoExternalImports_Fail(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "internal/domain"), 0755))
	writeTestFile(t, dir, "internal/domain/user.go", `package domain

import "gorm.io/gorm"

type User struct{ gorm.Model }
`)
	writeTestFile(t, dir, "go.mod", "module github.com/test/app\n\ngo 1.22\n")
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkDomainHasNoExternalImports()
	assert.Equal(t, "✗", r.status)
}

func TestCheckUsecaseDoesNotImportHandler_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkUsecaseDoesNotImportHandler()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckHandlerDoesNotImportRepository_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkHandlerDoesNotImportRepository()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckDIContainerExists_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkDIContainerExists()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckRepositoryImplsExist_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkRepositoryImplsExist()
	// At minimum should not error out (may warn if entity names don't match)
	assert.NotEqual(t, "✗", r.status, r.message)
}

// ─── Quality checks ───────────────────────────────────────────────────────────

func TestCheckEmptyGoFiles_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkEmptyGoFiles()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckEmptyGoFiles_Warn(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "internal/domain"), 0755))
	writeTestFile(t, dir, "internal/domain/stub.go", "package domain\n")
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkEmptyGoFiles()
	assert.Equal(t, "⚠", r.status)
}

func TestCheckPackageNamingConvention_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkPackageNamingConvention()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckPackageNamingConvention_Fail(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "internal/domain"), 0755))
	writeTestFile(t, dir, "internal/domain/bad.go", "package my_Domain\n\ntype User struct{}\n")
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkPackageNamingConvention()
	assert.Equal(t, "✗", r.status)
}

func TestCheckMainGoExists_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkMainGoExists()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckMainGoExists_Warn(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkMainGoExists()
	assert.Equal(t, "⚠", r.status)
}

// ─── Security checks ──────────────────────────────────────────────────────────

func TestCheckNoHardcodedSecrets_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkNoHardcodedSecrets()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckNoHardcodedSecrets_Fail(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "internal/domain"), 0755))
	writeTestFile(t, dir, "internal/domain/cfg.go", `package domain

const dbPassword = "supersecret123"
var password = "hardcoded"
`)
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkNoHardcodedSecrets()
	assert.Equal(t, "✗", r.status)
}

func TestCheckNoRawSQLStringFormat_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkNoRawSQLStringFormat()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckNoRawSQLStringFormat_Fail(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "internal/repository"), 0755))
	writeTestFile(t, dir, "internal/repository/repo.go", `package repository

import "fmt"

func query(id string) string {
	return fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", id)
}
`)
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkNoRawSQLStringFormat()
	assert.Equal(t, "✗", r.status)
}

func TestCheckNoUnsafePackage_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkNoUnsafePackage()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckNoUnsafePackage_Fail(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "internal/domain"), 0755))
	writeTestFile(t, dir, "internal/domain/hack.go", `package domain

import "unsafe"

var _ = unsafe.Sizeof(0)
`)
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkNoUnsafePackage()
	assert.Equal(t, "✗", r.status)
}

func TestCheckNoInsecureHTTPClient_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkNoInsecureHTTPClient()
	assert.Equal(t, "✓", r.status, r.message)
}

// ─── Standards checks ─────────────────────────────────────────────────────────

func TestCheckGoModHasCorrectModule_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkGoModHasCorrectModule()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckGoModHasCorrectModule_Fail(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkGoModHasCorrectModule()
	assert.Equal(t, "✗", r.status)
}

func TestCheckGocaYamlPresent_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkGocaYamlPresent()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckGocaYamlPresent_Warn(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkGocaYamlPresent()
	assert.Equal(t, "⚠", r.status)
}

func TestCheckFileNamingSnakeCase_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkFileNamingSnakeCase()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckFileNamingSnakeCase_Warn(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "internal/domain"), 0755))
	writeTestFile(t, dir, "internal/domain/user-repository.go", "package domain\n\ntype X struct{}\n")
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkFileNamingSnakeCase()
	assert.Equal(t, "⚠", r.status)
}

// ─── Test checks ──────────────────────────────────────────────────────────────

func TestCheckTestFilesExist_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkTestFilesExist()
	// domain has a test file in our setup
	assert.NotEqual(t, "✗", r.status, r.message)
}

func TestCheckTableDrivenTests_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkTableDrivenTests()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckMocksExist_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkMocksExist()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckMocksExist_Warn(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkMocksExist()
	assert.Equal(t, "⚠", r.status)
}

// ─── Dependency checks ────────────────────────────────────────────────────────

func TestCheckGoModTidy_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkGoModTidy()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckGoModTidy_MissingGoSum(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, dir, "go.mod", "module github.com/test/app\n\ngo 1.22\n\nrequire github.com/stretchr/testify v1.9.0\n")
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkGoModTidy()
	assert.Equal(t, "⚠", r.status)
}

func TestCheckNoReplaceDirectives_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkNoReplaceDirectives()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckNoReplaceDirectives_Warn(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, dir, "go.mod", "module github.com/test/app\n\ngo 1.22\n\nreplace github.com/foo/bar => ../bar\n")
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkNoReplaceDirectives()
	assert.Equal(t, "⚠", r.status)
}

func TestCheckGoVersionDeclared_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkGoVersionDeclared()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckKnownInsecureModules_Pass(t *testing.T) {
	setupAnalyzeProject(t)
	r := checkKnownInsecureModules()
	assert.Equal(t, "✓", r.status, r.message)
}

func TestCheckKnownInsecureModules_Fail(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, dir, "go.mod", "module github.com/test/app\n\ngo 1.22\n\nrequire github.com/dgrijalva/jwt-go v3.2.0+incompatible\n")
	old, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(old) })

	r := checkKnownInsecureModules()
	assert.Equal(t, "✗", r.status)
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func TestDirExists(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	assert.True(t, dirExists(dir))
	assert.False(t, dirExists(filepath.Join(dir, "notexist")))
}

func TestFileExists(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	f := filepath.Join(dir, "x.txt")
	require.NoError(t, os.WriteFile(f, []byte("hi"), 0644))
	assert.True(t, fileExists(f))
	assert.False(t, fileExists(filepath.Join(dir, "nofile.txt")))
	assert.False(t, fileExists(dir)) // dir is not a file
}

func TestAnalyzeGoFiles_SkipsTests(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.go"), []byte("package p"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a_test.go"), []byte("package p"), 0644))

	all, _ := analyzeGoFiles(dir, false)
	noTests, _ := analyzeGoFiles(dir, true)
	assert.Len(t, all, 2)
	assert.Len(t, noTests, 1)
}

func TestCollectEntityNames(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "user.go"), []byte(`package domain

type User struct{ ID uint }
type privateStruct struct{}
`), 0644))

	names := collectEntityNames(dir)
	assert.Contains(t, names, "User")
	assert.NotContains(t, names, "privateStruct")
}

// ─── Full run integration ─────────────────────────────────────────────────────

func TestRunArchitectureChecks_AllPass(t *testing.T) {
	setupAnalyzeProject(t)
	results := runArchitectureChecks()
	// All should be pass or warn, none should be fail on a well-formed project
	for _, r := range results {
		assert.NotEqual(t, "✗", r.status, "rule %s failed: %s", r.rule, r.message)
	}
}

func TestRunSecurityChecks_AllPass(t *testing.T) {
	setupAnalyzeProject(t)
	results := runSecurityChecks()
	for _, r := range results {
		assert.NotEqual(t, "✗", r.status, "rule %s failed: %s", r.rule, r.message)
	}
}

func TestRunDependencyChecks_AllPass(t *testing.T) {
	setupAnalyzeProject(t)
	results := runDependencyChecks()
	for _, r := range results {
		assert.NotEqual(t, "✗", r.status, "rule %s failed: %s", r.rule, r.message)
	}
}

func TestRunQualityChecks_AllPass(t *testing.T) {
	setupAnalyzeProject(t)
	results := runQualityChecks()
	for _, r := range results {
		assert.NotEqual(t, "✗", r.status, "rule %s failed: %s", r.rule, r.message)
	}
}

func TestRunStandardsChecks_AllPass(t *testing.T) {
	setupAnalyzeProject(t)
	results := runStandardsChecks()
	for _, r := range results {
		assert.NotEqual(t, "✗", r.status, "rule %s failed: %s", r.rule, r.message)
	}
}

func TestRunTestChecks_AllPass(t *testing.T) {
	setupAnalyzeProject(t)
	results := runTestChecks()
	for _, r := range results {
		assert.NotEqual(t, "✗", r.status, "rule %s failed: %s", r.rule, r.message)
	}
}
