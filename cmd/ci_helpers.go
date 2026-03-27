package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// buildCITemplateData reads project metadata from go.mod and .goca.yaml and
// returns a populated CITemplateData. If goVersion is empty the version is
// extracted from go.mod; if go.mod is absent it defaults to "1.25".
func buildCITemplateData(goVersion string) CITemplateData {
	data := CITemplateData{
		ProjectName: filepath.Base(getCurrentDir()),
		Module:      getModuleName(),
		GoVersion:   goVersion,
	}

	if data.GoVersion == "" {
		data.GoVersion = detectGoVersionFromMod()
	}

	// Try to read database from .goca.yaml via ConfigIntegration
	ci := NewConfigIntegration()
	ci.LoadConfigForProject()
	if ci.config != nil {
		data.Database = ci.GetDatabaseType("")
	}

	return data
}

// detectGoVersionFromMod reads the "go X.Y" directive from go.mod and returns
// the version string. Falls back to "1.25" when the file is missing or the
// directive is not found.
func detectGoVersionFromMod() string {
	f, err := os.Open("go.mod")
	if err != nil {
		return "1.25"
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "go ") {
			return strings.TrimPrefix(line, "go ")
		}
	}
	return "1.25"
}

// generateCIPipeline dispatches to the appropriate provider generator.
func generateCIPipeline(provider string, data CITemplateData, sm ...*SafetyManager) error {
	switch provider {
	case "github-actions":
		return generateGitHubActions(data, sm...)
	default:
		return fmt.Errorf("unsupported CI provider: %s (supported: github-actions)", provider)
	}
}

// generateGitHubActions writes GitHub Actions workflow files under
// .github/workflows/.
func generateGitHubActions(data CITemplateData, sm ...*SafetyManager) error {
	workflowDir := filepath.Join(".github", "workflows")

	// test.yml — always generated
	testYAML := generateTestWorkflow(data)
	ui.Step(1, "Generating test workflow")
	if err := writeFile(filepath.Join(workflowDir, "test.yml"), testYAML, sm...); err != nil {
		return fmt.Errorf("writing test.yml: %w", err)
	}
	ui.FileCreated(filepath.Join(workflowDir, "test.yml"))

	// build.yml — always generated
	buildYAML := generateBuildWorkflow(data)
	ui.Step(2, "Generating build workflow")
	if err := writeFile(filepath.Join(workflowDir, "build.yml"), buildYAML, sm...); err != nil {
		return fmt.Errorf("writing build.yml: %w", err)
	}
	ui.FileCreated(filepath.Join(workflowDir, "build.yml"))

	// deploy.yml — only when requested
	if data.WithDeploy {
		deployYAML := generateDeployWorkflow(data)
		ui.Step(3, "Generating deploy workflow")
		if err := writeFile(filepath.Join(workflowDir, "deploy.yml"), deployYAML, sm...); err != nil {
			return fmt.Errorf("writing deploy.yml: %w", err)
		}
		ui.FileCreated(filepath.Join(workflowDir, "deploy.yml"))
	}

	return nil
}
