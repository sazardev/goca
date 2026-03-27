package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ═══════════════════════════════════════════════════════════════════════════════
// ARCHITECTURE CHECKS
// ═══════════════════════════════════════════════════════════════════════════════

func runArchitectureChecks() []analyzeResult {
	return []analyzeResult{
		checkLayerDirsExist(),
		checkDomainHasNoExternalImports(),
		checkUsecaseDoesNotImportHandler(),
		checkRepositoryImplsExist(),
		checkDIContainerExists(),
		checkHandlerDoesNotImportRepository(),
	}
}

// checkLayerDirsExist verifies the four CA layers are present.
func checkLayerDirsExist() analyzeResult {
	layers := []string{
		"internal/domain",
		"internal/usecase",
		"internal/repository",
		"internal/handler",
	}
	var missing []string
	for _, l := range layers {
		if !dirExists(l) {
			missing = append(missing, l)
		}
	}
	if len(missing) == 0 {
		return analyzeResult{
			category: "Architecture",
			rule:     "layer-dirs",
			status:   "✓",
			message:  "All 4 Clean Architecture layers present",
		}
	}
	return analyzeResult{
		category:   "Architecture",
		rule:       "layer-dirs",
		status:     "✗",
		message:    fmt.Sprintf("Missing: %s", strings.Join(missing, ", ")),
		suggestion: "Run: goca feature <name> to scaffold the layer structure",
	}
}

// checkDomainHasNoExternalImports ensures internal/domain only imports stdlib.
func checkDomainHasNoExternalImports() analyzeResult {
	moduleName := getModuleName()
	domainDir := "internal/domain"
	if !dirExists(domainDir) {
		return analyzeResult{
			category: "Architecture",
			rule:     "domain-purity",
			status:   "⚠",
			message:  "internal/domain not found — skipped",
		}
	}

	files, _ := analyzeGoFiles(domainDir, true)
	for _, f := range files {
		violations := findForbiddenImports(f, moduleName, []string{
			"gorm.io", "database/sql", "github.com", "google.golang.org",
		})
		if len(violations) > 0 {
			return analyzeResult{
				category:   "Architecture",
				rule:       "domain-purity",
				status:     "✗",
				file:       f,
				message:    fmt.Sprintf("Domain imports external package: %s", violations[0]),
				suggestion: "Domain entities must not depend on ORMs or external frameworks",
			}
		}
	}
	return analyzeResult{
		category: "Architecture",
		rule:     "domain-purity",
		status:   "✓",
		message:  "Domain layer imports only stdlib",
	}
}

// checkUsecaseDoesNotImportHandler ensures use cases don't depend on handlers.
func checkUsecaseDoesNotImportHandler() analyzeResult {
	moduleName := getModuleName()
	ucDir := "internal/usecase"
	if !dirExists(ucDir) {
		return analyzeResult{
			category: "Architecture",
			rule:     "usecase-no-handler",
			status:   "⚠",
			message:  "internal/usecase not found — skipped",
		}
	}

	files, _ := analyzeGoFiles(ucDir, true)
	handlerPkg := moduleName + "/internal/handler"
	for _, f := range files {
		content := analyzeReadFile(f)
		if strings.Contains(content, handlerPkg) {
			return analyzeResult{
				category:   "Architecture",
				rule:       "usecase-no-handler",
				status:     "✗",
				file:       f,
				message:    "UseCase imports handler package (dependency inversion violation)",
				suggestion: "Use cases must only depend on repository interfaces, not handlers",
			}
		}
	}
	return analyzeResult{
		category: "Architecture",
		rule:     "usecase-no-handler",
		status:   "✓",
		message:  "Use cases do not import handler layer",
	}
}

// checkRepositoryImplsExist checks that for each domain entity there is a repo implementation.
func checkRepositoryImplsExist() analyzeResult {
	domainDir := "internal/domain"
	repoDir := "internal/repository"
	if !dirExists(domainDir) || !dirExists(repoDir) {
		return analyzeResult{
			category: "Architecture",
			rule:     "repo-impl-coverage",
			status:   "⚠",
			message:  "domain or repository dir missing — skipped",
		}
	}

	entities := collectEntityNames(domainDir)
	if len(entities) == 0 {
		return analyzeResult{
			category: "Architecture",
			rule:     "repo-impl-coverage",
			status:   "⚠",
			message:  "No domain entity files found",
		}
	}

	var noRepo []string
	for _, entity := range entities {
		lower := strings.ToLower(entity)
		found := false
		_ = filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if strings.Contains(strings.ToLower(path), lower) && strings.HasSuffix(path, ".go") {
				found = true
			}
			return nil
		})
		if !found {
			noRepo = append(noRepo, entity)
		}
	}

	if len(noRepo) == 0 {
		return analyzeResult{
			category: "Architecture",
			rule:     "repo-impl-coverage",
			status:   "✓",
			message:  fmt.Sprintf("Repository implementations found for all %d entities", len(entities)),
		}
	}
	return analyzeResult{
		category:   "Architecture",
		rule:       "repo-impl-coverage",
		status:     "⚠",
		message:    fmt.Sprintf("Entities with no repository: %s", strings.Join(noRepo, ", ")),
		suggestion: "Run: goca repository <Entity> --database postgres",
	}
}

// checkDIContainerExists verifies the DI container is wired.
func checkDIContainerExists() analyzeResult {
	if dirExists("internal/di") {
		files, _ := analyzeGoFiles("internal/di", true)
		if len(files) > 0 {
			return analyzeResult{
				category: "Architecture",
				rule:     "di-container",
				status:   "✓",
				message:  fmt.Sprintf("DI container present (%d files)", len(files)),
			}
		}
	}
	return analyzeResult{
		category:   "Architecture",
		rule:       "di-container",
		status:     "⚠",
		message:    "No DI container found in internal/di",
		suggestion: "Run: goca di <Entity> to generate the dependency injection container",
	}
}

// checkHandlerDoesNotImportRepository ensures handlers use use cases, not repos directly.
func checkHandlerDoesNotImportRepository() analyzeResult {
	moduleName := getModuleName()
	handlerDir := "internal/handler"
	if !dirExists(handlerDir) {
		return analyzeResult{
			category: "Architecture",
			rule:     "handler-no-repo",
			status:   "⚠",
			message:  "internal/handler not found — skipped",
		}
	}

	files, _ := analyzeGoFiles(handlerDir, true)
	repoPkg := moduleName + "/internal/repository"
	for _, f := range files {
		content := analyzeReadFile(f)
		if strings.Contains(content, repoPkg) {
			return analyzeResult{
				category:   "Architecture",
				rule:       "handler-no-repo",
				status:     "✗",
				file:       f,
				message:    "Handler imports repository directly (skips use case layer)",
				suggestion: "Handlers must depend only on use case interfaces, not repositories",
			}
		}
	}
	return analyzeResult{
		category: "Architecture",
		rule:     "handler-no-repo",
		status:   "✓",
		message:  "Handlers do not import repository layer directly",
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// QUALITY CHECKS
// ═══════════════════════════════════════════════════════════════════════════════

func runQualityChecks() []analyzeResult {
	return []analyzeResult{
		checkEmptyGoFiles(),
		checkPackageNamingConvention(),
		checkNoTODOsInGeneratedCode(),
		checkExportedFunctionsHaveDocs(),
		checkMainGoExists(),
	}
}

func checkEmptyGoFiles() analyzeResult {
	files, _ := analyzeGoFiles("internal", true)
	var empty []string
	for _, f := range files {
		content := strings.TrimSpace(analyzeReadFile(f))
		lines := strings.Split(content, "\n")
		// Only package declaration = effectively empty
		nonBlank := 0
		for _, l := range lines {
			l = strings.TrimSpace(l)
			if l != "" && !strings.HasPrefix(l, "package ") {
				nonBlank++
			}
		}
		if nonBlank == 0 {
			empty = append(empty, f)
		}
	}
	if len(empty) == 0 {
		return analyzeResult{
			category: "Quality",
			rule:     "no-empty-files",
			status:   "✓",
			message:  "No effectively empty .go files found",
		}
	}
	return analyzeResult{
		category:   "Quality",
		rule:       "no-empty-files",
		status:     "⚠",
		file:       empty[0],
		message:    fmt.Sprintf("%d empty file(s) — first: %s", len(empty), empty[0]),
		suggestion: "Remove or complete stub files that contain only a package declaration",
	}
}

func checkPackageNamingConvention() analyzeResult {
	files, _ := analyzeGoFiles("internal", true)
	for _, f := range files {
		content := analyzeReadFile(f)
		for _, line := range strings.Split(content, "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "package ") {
				continue
			}
			pkg := strings.TrimPrefix(line, "package ")
			pkg = strings.TrimSpace(pkg)
			// Package names must be lowercase, no underscores
			if pkg != strings.ToLower(pkg) || strings.Contains(pkg, "_") {
				return analyzeResult{
					category:   "Quality",
					rule:       "package-naming",
					status:     "✗",
					file:       f,
					message:    fmt.Sprintf("Package %q violates Go naming (must be lowercase, no underscores)", pkg),
					suggestion: "Rename package to all-lowercase single word per Go conventions",
				}
			}
			break
		}
	}
	return analyzeResult{
		category: "Quality",
		rule:     "package-naming",
		status:   "✓",
		message:  "All package names follow Go conventions",
	}
}

func checkNoTODOsInGeneratedCode() analyzeResult {
	files, _ := analyzeGoFiles("internal", false)
	var found []string
	for _, f := range files {
		content := analyzeReadFile(f)
		if strings.Contains(content, "// TODO") || strings.Contains(content, "// FIXME") {
			found = append(found, f)
		}
	}
	if len(found) == 0 {
		return analyzeResult{
			category: "Quality",
			rule:     "no-todos",
			status:   "✓",
			message:  "No TODO/FIXME comments in generated code",
		}
	}
	return analyzeResult{
		category:   "Quality",
		rule:       "no-todos",
		status:     "⚠",
		file:       found[0],
		message:    fmt.Sprintf("%d file(s) with TODO/FIXME — first: %s", len(found), found[0]),
		suggestion: "Resolve or track TODOs as issues before shipping to production",
	}
}

func checkExportedFunctionsHaveDocs() analyzeResult {
	if !dirExists("internal") {
		return analyzeResult{
			category: "Quality",
			rule:     "exported-docs",
			status:   "⚠",
			message:  "internal/ not found — skipped",
		}
	}

	fset := token.NewFileSet()
	files, _ := analyzeGoFiles("internal", true)
	var undocumented []string
	for _, f := range files {
		node, err := parser.ParseFile(fset, f, nil, parser.ParseComments)
		if err != nil {
			continue
		}
		for _, decl := range node.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if fn.Name == nil || !fn.Name.IsExported() {
				continue
			}
			if fn.Doc == nil || strings.TrimSpace(fn.Doc.Text()) == "" {
				undocumented = append(undocumented, fmt.Sprintf("%s: %s()", f, fn.Name.Name))
			}
		}
	}

	if len(undocumented) == 0 {
		return analyzeResult{
			category: "Quality",
			rule:     "exported-docs",
			status:   "✓",
			message:  "All exported functions have doc comments",
		}
	}
	return analyzeResult{
		category:   "Quality",
		rule:       "exported-docs",
		status:     "⚠",
		message:    fmt.Sprintf("%d exported function(s) without doc comment", len(undocumented)),
		suggestion: fmt.Sprintf("Add godoc comment to: %s (and others)", undocumented[0]),
	}
}

func checkMainGoExists() analyzeResult {
	candidates := []string{"main.go", "cmd/main.go"}
	for _, c := range candidates {
		if fileExists(c) {
			return analyzeResult{
				category: "Quality",
				rule:     "main-go",
				status:   "✓",
				message:  fmt.Sprintf("Entry point found at %s", c),
			}
		}
	}
	return analyzeResult{
		category:   "Quality",
		rule:       "main-go",
		status:     "⚠",
		message:    "No main.go found at root or cmd/",
		suggestion: "Run: goca init <project-name> to generate project entry point",
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// SECURITY CHECKS
// ═══════════════════════════════════════════════════════════════════════════════

func runSecurityChecks() []analyzeResult {
	return []analyzeResult{
		checkNoHardcodedSecrets(),
		checkNoRawSQLStringFormat(),
		checkNoUnsafePackage(),
		checkNoInsecureHTTPClient(),
		checkEnvVarsForSensitiveConfig(),
	}
}

func checkNoHardcodedSecrets() analyzeResult {
	patterns := []string{
		"password = \"", "secret = \"", "token = \"",
		"apikey = \"", "api_key = \"", "passwd = \"",
		`password:"`, `secret:"`, `token:"`,
	}
	files, _ := analyzeGoFiles("internal", false)
	for _, f := range files {
		lower := strings.ToLower(analyzeReadFile(f))
		for _, p := range patterns {
			if strings.Contains(lower, p) {
				return analyzeResult{
					category:   "Security",
					rule:       "no-hardcoded-secrets",
					status:     "✗",
					file:       f,
					message:    fmt.Sprintf("Possible hardcoded secret pattern %q", p),
					suggestion: "Use os.Getenv() or a secrets manager — never hardcode credentials",
				}
			}
		}
	}
	return analyzeResult{
		category: "Security",
		rule:     "no-hardcoded-secrets",
		status:   "✓",
		message:  "No hardcoded secret patterns detected",
	}
}

func checkNoRawSQLStringFormat() analyzeResult {
	files, _ := analyzeGoFiles("internal", true)
	dangerous := []string{
		`fmt.Sprintf("SELECT`, `fmt.Sprintf("INSERT`,
		`fmt.Sprintf("UPDATE`, `fmt.Sprintf("DELETE`,
		`fmt.Sprintf("DROP`,
	}
	for _, f := range files {
		content := analyzeReadFile(f)
		for _, pattern := range dangerous {
			if strings.Contains(content, pattern) {
				return analyzeResult{
					category:   "Security",
					rule:       "no-sql-injection",
					status:     "✗",
					file:       f,
					message:    "SQL query built with fmt.Sprintf — potential injection risk (OWASP A03)",
					suggestion: "Use parameterised queries: db.QueryContext(ctx, sql, args...)",
				}
			}
		}
	}
	return analyzeResult{
		category: "Security",
		rule:     "no-sql-injection",
		status:   "✓",
		message:  "No fmt.Sprintf SQL construction patterns found",
	}
}

func checkNoUnsafePackage() analyzeResult {
	files, _ := analyzeGoFiles("internal", true)
	for _, f := range files {
		content := analyzeReadFile(f)
		if strings.Contains(content, `"unsafe"`) {
			return analyzeResult{
				category:   "Security",
				rule:       "no-unsafe",
				status:     "✗",
				file:       f,
				message:    "Uses the 'unsafe' package",
				suggestion: "Avoid unsafe in generated application code — use typed interfaces instead",
			}
		}
	}
	return analyzeResult{
		category: "Security",
		rule:     "no-unsafe",
		status:   "✓",
		message:  "No 'unsafe' package usage in generated code",
	}
}

func checkNoInsecureHTTPClient() analyzeResult {
	insecurePattern := "InsecureSkipVerify: true"
	files, _ := analyzeGoFiles("internal", true)
	for _, f := range files {
		if strings.Contains(analyzeReadFile(f), insecurePattern) {
			return analyzeResult{
				category:   "Security",
				rule:       "no-tls-skip",
				status:     "✗",
				file:       f,
				message:    "TLS verification disabled (InsecureSkipVerify: true)",
				suggestion: "Never disable TLS verification in production code",
			}
		}
	}
	return analyzeResult{
		category: "Security",
		rule:     "no-tls-skip",
		status:   "✓",
		message:  "No TLS verification bypass detected",
	}
}

func checkEnvVarsForSensitiveConfig() analyzeResult {
	// Verify sensitive config is read from environment, not config files or constants
	files, _ := analyzeGoFiles("internal", true)
	usesEnv := false
	for _, f := range files {
		if strings.Contains(analyzeReadFile(f), "os.Getenv") {
			usesEnv = true
			break
		}
	}
	// Check cache package separately
	if !usesEnv && dirExists("internal/cache") {
		cacheFiles, _ := analyzeGoFiles("internal/cache", true)
		for _, f := range cacheFiles {
			if strings.Contains(analyzeReadFile(f), "os.Getenv") {
				usesEnv = true
				break
			}
		}
	}
	if usesEnv {
		return analyzeResult{
			category: "Security",
			rule:     "env-config",
			status:   "✓",
			message:  "Sensitive config read from environment variables",
		}
	}
	return analyzeResult{
		category:   "Security",
		rule:       "env-config",
		status:     "⚠",
		message:    "No os.Getenv() usage found — config values may be hardcoded",
		suggestion: "Store DSN, secrets, and ports in environment variables",
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// STANDARDS CHECKS
// ═══════════════════════════════════════════════════════════════════════════════

func runStandardsChecks() []analyzeResult {
	return []analyzeResult{
		checkFileNamingSnakeCase(),
		checkNoInitFunctionsInDomain(),
		checkGoModHasCorrectModule(),
		checkGocaYamlPresent(),
		checkContextPropagation(),
	}
}

func checkFileNamingSnakeCase() analyzeResult {
	files, _ := analyzeGoFiles("internal", false)
	var violations []string
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, ".go")
		name = strings.TrimSuffix(name, "_test")
		if strings.Contains(name, "-") {
			violations = append(violations, f)
		}
	}
	if len(violations) == 0 {
		return analyzeResult{
			category: "Standards",
			rule:     "file-naming",
			status:   "✓",
			message:  "All Go files use snake_case naming",
		}
	}
	return analyzeResult{
		category:   "Standards",
		rule:       "file-naming",
		status:     "⚠",
		file:       violations[0],
		message:    fmt.Sprintf("%d file(s) use kebab-case — first: %s", len(violations), violations[0]),
		suggestion: "Go convention: use snake_case for file names (user_repository.go)",
	}
}

func checkNoInitFunctionsInDomain() analyzeResult {
	domainDir := "internal/domain"
	if !dirExists(domainDir) {
		return analyzeResult{
			category: "Standards",
			rule:     "no-init-in-domain",
			status:   "⚠",
			message:  "internal/domain not found — skipped",
		}
	}
	files, _ := analyzeGoFiles(domainDir, true)
	for _, f := range files {
		content := analyzeReadFile(f)
		if strings.Contains(content, "\nfunc init()") || strings.Contains(content, "\nfunc init() {") {
			return analyzeResult{
				category:   "Standards",
				rule:       "no-init-in-domain",
				status:     "⚠",
				file:       f,
				message:    "Domain entity uses init() function",
				suggestion: "Avoid init() in domain — use explicit constructors instead",
			}
		}
	}
	return analyzeResult{
		category: "Standards",
		rule:     "no-init-in-domain",
		status:   "✓",
		message:  "No init() functions in domain layer",
	}
}

func checkGoModHasCorrectModule() analyzeResult {
	if !fileExists("go.mod") {
		return analyzeResult{
			category:   "Standards",
			rule:       "go-mod-valid",
			status:     "✗",
			message:    "go.mod not found",
			suggestion: "Run: go mod init <module-path>",
		}
	}
	content := analyzeReadFile("go.mod")
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			mod := strings.TrimPrefix(line, "module ")
			mod = strings.TrimSpace(mod)
			if mod == "" {
				return analyzeResult{
					category:   "Standards",
					rule:       "go-mod-valid",
					status:     "✗",
					message:    "go.mod has empty module declaration",
					suggestion: "Set module path: module github.com/yourname/yourproject",
				}
			}
			return analyzeResult{
				category: "Standards",
				rule:     "go-mod-valid",
				status:   "✓",
				message:  fmt.Sprintf("go.mod module: %s", mod),
			}
		}
	}
	return analyzeResult{
		category:   "Standards",
		rule:       "go-mod-valid",
		status:     "✗",
		message:    "go.mod missing module declaration",
		suggestion: "Add: module <your-module-path>",
	}
}

func checkGocaYamlPresent() analyzeResult {
	if !fileExists(".goca.yaml") {
		return analyzeResult{
			category:   "Standards",
			rule:       "goca-yaml",
			status:     "⚠",
			message:    ".goca.yaml not found",
			suggestion: "Run: goca init <project-name> to create project configuration",
		}
	}
	content := strings.TrimSpace(analyzeReadFile(".goca.yaml"))
	if len(content) == 0 {
		return analyzeResult{
			category:   "Standards",
			rule:       "goca-yaml",
			status:     "⚠",
			message:    ".goca.yaml is empty",
			suggestion: "Run: goca upgrade to regenerate the config file",
		}
	}
	return analyzeResult{
		category: "Standards",
		rule:     "goca-yaml",
		status:   "✓",
		message:  ".goca.yaml present and configured",
	}
}

func checkContextPropagation() analyzeResult {
	repoDir := "internal/repository"
	ucDir := "internal/usecase"
	if !dirExists(repoDir) && !dirExists(ucDir) {
		return analyzeResult{
			category: "Standards",
			rule:     "context-propagation",
			status:   "⚠",
			message:  "repository and usecase dirs missing — skipped",
		}
	}

	checkDirs := []string{repoDir, ucDir}
	var noContext []string
	for _, dir := range checkDirs {
		if !dirExists(dir) {
			continue
		}
		files, _ := analyzeGoFiles(dir, true)
		for _, f := range files {
			content := analyzeReadFile(f)
			// If file has func signatures with DB/repo calls but no context import
			if strings.Contains(content, "func (") && !strings.Contains(content, "context.Context") {
				noContext = append(noContext, f)
			}
		}
	}

	if len(noContext) == 0 {
		return analyzeResult{
			category: "Standards",
			rule:     "context-propagation",
			status:   "✓",
			message:  "context.Context used in repository and use case layers",
		}
	}
	return analyzeResult{
		category:   "Standards",
		rule:       "context-propagation",
		status:     "⚠",
		file:       noContext[0],
		message:    fmt.Sprintf("%d file(s) with methods that may lack context.Context", len(noContext)),
		suggestion: "Pass context.Context as first parameter to all I/O methods",
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// TEST CHECKS
// ═══════════════════════════════════════════════════════════════════════════════

func runTestChecks() []analyzeResult {
	return []analyzeResult{
		checkTestFilesExist(),
		checkTableDrivenTests(),
		checkMocksExist(),
		checkTestHelpersTempDir(),
	}
}

func checkTestFilesExist() analyzeResult {
	layers := []string{"internal/domain", "internal/usecase", "internal/repository"}
	var noTests []string
	for _, layer := range layers {
		if !dirExists(layer) {
			continue
		}
		src, _ := analyzeGoFiles(layer, true)
		tst, _ := analyzeGoFiles(layer, false)
		testCount := len(tst) - len(src)
		if testCount == 0 && len(src) > 0 {
			noTests = append(noTests, layer)
		}
	}
	if len(noTests) == 0 {
		return analyzeResult{
			category: "Tests",
			rule:     "test-files-exist",
			status:   "✓",
			message:  "Test files present in all non-empty layers",
		}
	}
	return analyzeResult{
		category:   "Tests",
		rule:       "test-files-exist",
		status:     "⚠",
		message:    fmt.Sprintf("No *_test.go files in: %s", strings.Join(noTests, ", ")),
		suggestion: "Run: goca test-integration <Entity> to generate integration tests",
	}
}

func checkTableDrivenTests() analyzeResult {
	files, _ := analyzeGoFiles(".", false)
	// Filter to only _test.go files
	var testFiles []string
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			testFiles = append(testFiles, f)
		}
	}
	if len(testFiles) == 0 {
		return analyzeResult{
			category: "Tests",
			rule:     "table-driven-tests",
			status:   "⚠",
			message:  "No test files found",
		}
	}

	tableCount := 0
	for _, f := range testFiles {
		content := analyzeReadFile(f)
		if strings.Contains(content, "[]struct{") || strings.Contains(content, "[]struct {") {
			tableCount++
		}
	}
	if tableCount > 0 {
		return analyzeResult{
			category: "Tests",
			rule:     "table-driven-tests",
			status:   "✓",
			message:  fmt.Sprintf("%d test file(s) use table-driven pattern", tableCount),
		}
	}
	return analyzeResult{
		category:   "Tests",
		rule:       "table-driven-tests",
		status:     "⚠",
		message:    "No table-driven tests found ([]struct{...} pattern)",
		suggestion: "Prefer table-driven tests for functions with multiple input variants",
	}
}

func checkMocksExist() analyzeResult {
	mockDirs := []string{"internal/mocks", "mocks"}
	for _, d := range mockDirs {
		if dirExists(d) {
			files, _ := analyzeGoFiles(d, false)
			if len(files) > 0 {
				return analyzeResult{
					category: "Tests",
					rule:     "mocks-exist",
					status:   "✓",
					message:  fmt.Sprintf("Mocks found at %s (%d files)", d, len(files)),
				}
			}
		}
	}
	return analyzeResult{
		category:   "Tests",
		rule:       "mocks-exist",
		status:     "⚠",
		message:    "No mock directory found",
		suggestion: "Run: goca mocks to generate testify mocks for all interfaces",
	}
}

func checkTestHelpersTempDir() analyzeResult {
	// Validate tests use t.TempDir() for filesystem ops rather than hard paths
	var testFiles []string
	files, _ := analyzeGoFiles(".", false)
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			testFiles = append(testFiles, f)
		}
	}
	for _, f := range testFiles {
		content := analyzeReadFile(f)
		// Warn if test creates files under /tmp directly
		if strings.Contains(content, `"/tmp/`) || strings.Contains(content, `"/var/tmp/`) {
			return analyzeResult{
				category:   "Tests",
				rule:       "test-tempdir",
				status:     "⚠",
				file:       f,
				message:    "Test uses hardcoded /tmp path instead of t.TempDir()",
				suggestion: "Use t.TempDir() — it's auto-cleaned and test-isolated",
			}
		}
	}
	return analyzeResult{
		category: "Tests",
		rule:     "test-tempdir",
		status:   "✓",
		message:  "No hardcoded /tmp paths in test files",
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// DEPENDENCY CHECKS
// ═══════════════════════════════════════════════════════════════════════════════

func runDependencyChecks() []analyzeResult {
	return []analyzeResult{
		checkGoModTidy(),
		checkNoReplaceDirectives(),
		checkGoVersionDeclared(),
		checkKnownInsecureModules(),
	}
}

func checkGoModTidy() analyzeResult {
	if !fileExists("go.mod") {
		return analyzeResult{
			category:   "Dependencies",
			rule:       "go-mod-tidy",
			status:     "✗",
			message:    "go.mod not found",
			suggestion: "Run: go mod init <module-path>",
		}
	}
	// Check if go.sum exists when go.mod has external deps
	content := analyzeReadFile("go.mod")
	if strings.Contains(content, "require") && !fileExists("go.sum") {
		return analyzeResult{
			category:   "Dependencies",
			rule:       "go-mod-tidy",
			status:     "⚠",
			message:    "go.mod has require block but go.sum is missing",
			suggestion: "Run: go mod tidy",
		}
	}
	return analyzeResult{
		category: "Dependencies",
		rule:     "go-mod-tidy",
		status:   "✓",
		message:  "go.mod and go.sum present",
	}
}

func checkNoReplaceDirectives() analyzeResult {
	if !fileExists("go.mod") {
		return analyzeResult{
			category: "Dependencies",
			rule:     "no-replace",
			status:   "⚠",
			message:  "go.mod not found — skipped",
		}
	}
	content := analyzeReadFile("go.mod")
	if strings.Contains(content, "\nreplace ") {
		return analyzeResult{
			category:   "Dependencies",
			rule:       "no-replace",
			status:     "⚠",
			message:    "go.mod contains replace directive(s)",
			suggestion: "Replace directives are for local development; remove before production release",
		}
	}
	return analyzeResult{
		category: "Dependencies",
		rule:     "no-replace",
		status:   "✓",
		message:  "No replace directives in go.mod",
	}
}

func checkGoVersionDeclared() analyzeResult {
	if !fileExists("go.mod") {
		return analyzeResult{
			category:   "Dependencies",
			rule:       "go-version",
			status:     "✗",
			message:    "go.mod not found",
			suggestion: "Run: go mod init <module-path>",
		}
	}
	content := analyzeReadFile("go.mod")
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "go ") {
			version := strings.TrimPrefix(line, "go ")
			return analyzeResult{
				category: "Dependencies",
				rule:     "go-version",
				status:   "✓",
				message:  fmt.Sprintf("Go version declared: %s", strings.TrimSpace(version)),
			}
		}
	}
	return analyzeResult{
		category:   "Dependencies",
		rule:       "go-version",
		status:     "⚠",
		message:    "No Go version declared in go.mod",
		suggestion: "Add: go 1.22 (or your target version) to go.mod",
	}
}

func checkKnownInsecureModules() analyzeResult {
	// Modules with known issues or that have been superseded
	insecureModules := map[string]string{
		"github.com/dgrijalva/jwt-go": "Use github.com/golang-jwt/jwt/v5 instead (dgrijalva/jwt-go is unmaintained)",
		"github.com/go-redis/redis":   "Use github.com/redis/go-redis/v9 (go-redis v8 is deprecated)",
		"gopkg.in/mgo.v2":             "Use go.mongodb.org/mongo-driver instead (mgo.v2 is unmaintained)",
		"github.com/astaxie/beego":    "Use github.com/beego/beego/v2 (astaxie path is deprecated)",
	}

	if !fileExists("go.mod") {
		return analyzeResult{
			category: "Dependencies",
			rule:     "no-insecure-modules",
			status:   "⚠",
			message:  "go.mod not found — skipped",
		}
	}
	content := analyzeReadFile("go.mod")
	for mod, advice := range insecureModules {
		if strings.Contains(content, mod) {
			return analyzeResult{
				category:   "Dependencies",
				rule:       "no-insecure-modules",
				status:     "✗",
				message:    fmt.Sprintf("Known-insecure/deprecated module: %s", mod),
				suggestion: advice,
			}
		}
	}
	return analyzeResult{
		category: "Dependencies",
		rule:     "no-insecure-modules",
		status:   "✓",
		message:  "No known-insecure or deprecated modules detected",
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// SHARED HELPERS
// ═══════════════════════════════════════════════════════════════════════════════

// findForbiddenImports parses a Go file and returns any imports from forbiddenPrefixes,
// excluding the project's own module (moduleName).
func findForbiddenImports(filePath, moduleName string, forbiddenPrefixes []string) []string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil
	}
	var found []string
	for _, imp := range node.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if strings.HasPrefix(path, moduleName) {
			continue // own module is allowed
		}
		for _, prefix := range forbiddenPrefixes {
			if strings.HasPrefix(path, prefix) {
				found = append(found, path)
				break
			}
		}
	}
	return found
}

// collectEntityNames scans domainDir for .go files and extracts PascalCase struct names
// that look like domain entities (exported structs).
func collectEntityNames(domainDir string) []string {
	files, _ := analyzeGoFiles(domainDir, true)
	var names []string
	seen := make(map[string]bool)
	for _, f := range files {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, f, nil, 0)
		if err != nil {
			continue
		}
		for _, decl := range node.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if _, ok := ts.Type.(*ast.StructType); !ok {
					continue
				}
				name := ts.Name.Name
				if ts.Name.IsExported() && !seen[name] {
					names = append(names, name)
					seen[name] = true
				}
			}
		}
	}
	return names
}
