package testing

import (
	"os/exec"
	"strings"
)

// CLIRunner handles execution of Goca CLI commands in tests
type CLIRunner struct {
	suite *TestSuite
}

// NewCLIRunner creates a new CLI runner for the test suite
func NewCLIRunner(suite *TestSuite) *CLIRunner {
	return &CLIRunner{suite: suite}
}

// Run executes a Goca CLI command and returns the output
func (r *CLIRunner) Run(args ...string) (string, error) {
	// Build the full command with goca binary
	cmd := exec.Command("goca", args...)
	cmd.Dir = r.suite.tempDir

	output, err := cmd.CombinedOutput()
	return string(output), err
}

// RunSuccess executes a command and expects it to succeed
func (r *CLIRunner) RunSuccess(args ...string) string {
	r.suite.t.Helper()
	output, err := r.Run(args...)
	if err != nil {
		r.suite.t.Fatalf("Command 'goca %s' failed: %v\nOutput: %s",
			strings.Join(args, " "), err, output)
	}
	return output
}

// RunFailure executes a command and expects it to fail
func (r *CLIRunner) RunFailure(args ...string) string {
	r.suite.t.Helper()
	output, err := r.Run(args...)
	if err == nil {
		r.suite.t.Fatalf("Command 'goca %s' should have failed but succeeded\nOutput: %s",
			strings.Join(args, " "), output)
	}
	return output
}

// InitProject runs goca init with various options
func (r *CLIRunner) InitProject(name, module string, options ...string) string {
	args := []string{"init", name}
	if module != "" {
		args = append(args, "--module", module)
	}
	args = append(args, options...)
	return r.RunSuccess(args...)
}

// GenerateFeature runs goca feature with specified parameters
func (r *CLIRunner) GenerateFeature(name, fields string, options ...string) string {
	args := []string{"feature", name, "--fields", fields}
	args = append(args, options...)
	return r.RunSuccess(args...)
}

// GenerateEntity runs goca entity with specified parameters
func (r *CLIRunner) GenerateEntity(name, fields string, options ...string) string {
	args := []string{"entity", name, "--fields", fields}
	args = append(args, options...)
	return r.RunSuccess(args...)
}

// GenerateUseCase runs goca usecase with specified parameters
func (r *CLIRunner) GenerateUseCase(name, entity string, options ...string) string {
	args := []string{"usecase", name, "--entity", entity}
	args = append(args, options...)
	return r.RunSuccess(args...)
}

// GenerateRepository runs goca repository with specified parameters
func (r *CLIRunner) GenerateRepository(entity string, options ...string) string {
	args := []string{"repository", entity}
	args = append(args, options...)
	return r.RunSuccess(args...)
}

// GenerateHandler runs goca handler with specified parameters
func (r *CLIRunner) GenerateHandler(entity string, options ...string) string {
	args := []string{"handler", entity}
	args = append(args, options...)
	return r.RunSuccess(args...)
}

// GenerateMessages runs goca messages with specified parameters
func (r *CLIRunner) GenerateMessages(entity string, options ...string) string {
	args := []string{"messages", entity}
	args = append(args, options...)
	return r.RunSuccess(args...)
}

// GenerateDI runs goca di with specified parameters
func (r *CLIRunner) GenerateDI(features string, options ...string) string {
	args := []string{"di", "--features", features}
	args = append(args, options...)
	return r.RunSuccess(args...)
}

// GenerateInterfaces runs goca interfaces with specified parameters
func (r *CLIRunner) GenerateInterfaces(entity string, options ...string) string {
	args := []string{"interfaces", entity}
	args = append(args, options...)
	return r.RunSuccess(args...)
}

// Version runs goca version command
func (r *CLIRunner) Version() string {
	return r.RunSuccess("version")
}

// Help runs goca help command
func (r *CLIRunner) Help(command ...string) string {
	args := append([]string{"help"}, command...)
	return r.RunSuccess(args...)
}
