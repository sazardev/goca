package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// doctorCheck holds the result of a single project health check.
type doctorCheck struct {
	name       string
	status     string // "✓", "✗", "⚠"
	message    string
	suggestion string
}

var doctorFix bool

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check project health and Clean Architecture structure",
	Long: `doctor runs a series of health checks on the current project to verify
that it follows Clean Architecture conventions and has the expected structure.

It checks:
  - go.mod presence and module path
  - .goca.yaml configuration
  - Clean Architecture directory structure
  - go build ./... compiles without errors
  - go vet ./... has no warnings
  - Dependency injection container

Use --fix to automatically create missing directories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Header("Goca Doctor — Project Health Check")
		ui.Blank()

		checks := runAllChecks()
		printChecks(checks)

		failed := countByStatus(checks, "✗")
		warned := countByStatus(checks, "⚠")
		passed := countByStatus(checks, "✓")

		ui.Blank()
		ui.Info(fmt.Sprintf("Results: %d passed, %d warnings, %d failed", passed, warned, failed))

		if failed > 0 {
			return fmt.Errorf("%d health check(s) failed — see suggestions above", failed)
		}
		return nil
	},
}

func init() {
	doctorCmd.Flags().BoolVar(&doctorFix, "fix", false, "Automatically create missing directories")
}

func runAllChecks() []doctorCheck {
	return []doctorCheck{
		checkGoMod(),
		checkGocaYaml(),
		checkProjectStructure(),
		checkGoBuild(),
		checkGoVet(),
		checkDIContainer(),
	}
}

func checkGoMod() doctorCheck {
	check := doctorCheck{name: "go.mod"}

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		check.status = "✗"
		check.message = "go.mod not found"
		check.suggestion = "Run: go mod init <module-path>"
		return check
	}

	content, err := os.ReadFile("go.mod")
	if err != nil {
		check.status = "⚠"
		check.message = "go.mod found but unreadable"
		check.suggestion = "Check file permissions on go.mod"
		return check
	}

	if strings.Contains(string(content), "module ") {
		check.status = "✓"
		check.message = "go.mod present with module declaration"
	} else {
		check.status = "⚠"
		check.message = "go.mod found but no module declaration"
		check.suggestion = "Add 'module <path>' to go.mod"
	}
	return check
}

func checkGocaYaml() doctorCheck {
	check := doctorCheck{name: ".goca.yaml"}

	if _, err := os.Stat(".goca.yaml"); os.IsNotExist(err) {
		check.status = "⚠"
		check.message = ".goca.yaml not found (optional but recommended)"
		check.suggestion = "Run: goca init <project-name> to create it"
		return check
	}

	content, err := os.ReadFile(".goca.yaml")
	if err != nil {
		check.status = "⚠"
		check.message = ".goca.yaml found but unreadable"
		check.suggestion = "Check file permissions on .goca.yaml"
		return check
	}

	if len(strings.TrimSpace(string(content))) == 0 {
		check.status = "⚠"
		check.message = ".goca.yaml is empty"
		check.suggestion = "Run: goca init <project-name> to regenerate"
		return check
	}

	check.status = "✓"
	check.message = ".goca.yaml present and non-empty"
	return check
}

func checkProjectStructure() doctorCheck {
	check := doctorCheck{name: "Clean Architecture dirs"}

	requiredDirs := []string{
		"internal/domain",
		"internal/usecase",
		"internal/repository",
		"internal/handler",
	}

	missing := make([]string, 0, len(requiredDirs))
	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			missing = append(missing, dir)
		}
	}

	if len(missing) == 0 {
		check.status = "✓"
		check.message = "All Clean Architecture directories present"
		return check
	}

	if doctorFix {
		for _, dir := range missing {
			if err := os.MkdirAll(dir, 0755); err != nil {
				check.status = "✗"
				check.message = fmt.Sprintf("Failed to create %s: %v", dir, err)
				check.suggestion = fmt.Sprintf("Create directory manually: mkdir -p %s", dir)
				return check
			}
			ui.Debug(fmt.Sprintf("Created directory: %s", dir))
		}
		check.status = "✓"
		check.message = fmt.Sprintf("Created %d missing directories", len(missing))
		return check
	}

	check.status = "⚠"
	check.message = fmt.Sprintf("%d directories missing: %s", len(missing), strings.Join(missing, ", "))
	check.suggestion = "Run: goca doctor --fix  or  goca feature <name>"
	return check
}

func checkGoBuild() doctorCheck {
	check := doctorCheck{name: "go build ./..."}

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		check.status = "⚠"
		check.message = "Skipped (no go.mod found)"
		return check
	}

	//nolint:gosec // args are not user-controlled
	out, err := exec.Command("go", "build", "./...").CombinedOutput()
	if err != nil {
		check.status = "✗"
		check.message = "Build failed"
		check.suggestion = "Fix build errors: " + strings.TrimSpace(string(out))
		return check
	}

	check.status = "✓"
	check.message = "Project compiles without errors"
	return check
}

func checkGoVet() doctorCheck {
	check := doctorCheck{name: "go vet ./..."}

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		check.status = "⚠"
		check.message = "Skipped (no go.mod found)"
		return check
	}

	//nolint:gosec // args are not user-controlled
	out, err := exec.Command("go", "vet", "./...").CombinedOutput()
	if err != nil {
		check.status = "⚠"
		check.message = "go vet reported issues"
		check.suggestion = "Fix vet warnings: " + strings.TrimSpace(string(out))
		return check
	}

	check.status = "✓"
	check.message = "go vet passes with no warnings"
	return check
}

func checkDIContainer() doctorCheck {
	check := doctorCheck{name: "DI container"}

	candidates := []string{
		"internal/di",
		filepath.Join("internal", "di"),
	}

	for _, dir := range candidates {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			check.status = "✓"
			check.message = fmt.Sprintf("DI container found at %s", dir)
			return check
		}
	}

	check.status = "⚠"
	check.message = "No DI container directory found"
	check.suggestion = "Run: goca di  to generate the dependency injection container"
	return check
}

func printChecks(checks []doctorCheck) {
	rows := make([][]string, len(checks))
	for i, c := range checks {
		suggestion := c.suggestion
		if suggestion == "" {
			suggestion = "—"
		}
		rows[i] = []string{c.status, c.name, c.message, suggestion}
	}
	ui.Table([]string{"", "Check", "Details", "Suggestion"}, rows)
}

func countByStatus(checks []doctorCheck, status string) int {
	count := 0
	for _, c := range checks {
		if c.status == status {
			count++
		}
	}
	return count
}
