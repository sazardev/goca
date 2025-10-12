package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// DependencyManager handles go.mod updates and version checking
type DependencyManager struct {
	projectRoot string
	goModPath   string
	dryRun      bool
}

// Dependency represents a Go module dependency
type Dependency struct {
	Module  string
	Version string
	Type    string // "required", "optional", "suggested"
	Reason  string
}

// NewDependencyManager creates a new dependency manager
func NewDependencyManager(projectRoot string, dryRun bool) *DependencyManager {
	return &DependencyManager{
		projectRoot: projectRoot,
		goModPath:   filepath.Join(projectRoot, "go.mod"),
		dryRun:      dryRun,
	}
}

// CommonDependencies returns list of commonly used dependencies
func (dm *DependencyManager) CommonDependencies() map[string]Dependency {
	return map[string]Dependency{
		"validator": {
			Module:  "github.com/go-playground/validator/v10",
			Version: "v10.16.0",
			Type:    "optional",
			Reason:  "struct validation for DTOs",
		},
		"jwt": {
			Module:  "github.com/golang-jwt/jwt/v5",
			Version: "v5.2.0",
			Type:    "optional",
			Reason:  "JWT authentication",
		},
		"cors": {
			Module:  "github.com/rs/cors",
			Version: "v1.10.1",
			Type:    "optional",
			Reason:  "CORS middleware",
		},
		"uuid": {
			Module:  "github.com/google/uuid",
			Version: "v1.5.0",
			Type:    "optional",
			Reason:  "UUID generation",
		},
		"bcrypt": {
			Module:  "golang.org/x/crypto",
			Version: "v0.17.0",
			Type:    "optional",
			Reason:  "password hashing",
		},
		"testify": {
			Module:  "github.com/stretchr/testify",
			Version: "v1.8.4",
			Type:    "optional",
			Reason:  "testing assertions and mocks",
		},
		"mock": {
			Module:  "github.com/golang/mock",
			Version: "v1.6.0",
			Type:    "optional",
			Reason:  "mock generation for testing",
		},
		"grpc": {
			Module:  "google.golang.org/grpc",
			Version: "v1.60.0",
			Type:    "required",
			Reason:  "gRPC protocol support",
		},
		"protobuf": {
			Module:  "google.golang.org/protobuf",
			Version: "v1.31.0",
			Type:    "required",
			Reason:  "Protocol Buffers",
		},
	}
}

// SuggestDependencies suggests optional dependencies based on features
func (dm *DependencyManager) SuggestDependencies(features []string) []Dependency {
	suggestions := make([]Dependency, 0)
	commonDeps := dm.CommonDependencies()

	for _, feature := range features {
		switch strings.ToLower(feature) {
		case "validation":
			suggestions = append(suggestions, commonDeps["validator"])
		case "auth", "authentication":
			suggestions = append(suggestions, commonDeps["jwt"], commonDeps["bcrypt"])
		case "grpc":
			suggestions = append(suggestions, commonDeps["grpc"], commonDeps["protobuf"])
		case "testing", "test":
			suggestions = append(suggestions, commonDeps["testify"], commonDeps["mock"])
		case "uuid":
			suggestions = append(suggestions, commonDeps["uuid"])
		}
	}

	return suggestions
}

// AddDependency adds a dependency to go.mod
func (dm *DependencyManager) AddDependency(dep Dependency) error {
	if dm.dryRun {
		fmt.Printf("[DRY-RUN] Would add dependency: %s %s\n", dep.Module, dep.Version)
		return nil
	}

	// Check if dependency already exists
	exists, err := dm.DependencyExists(dep.Module)
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("Dependency %s already exists\n", dep.Module)
		return nil
	}

	// Use go get to add dependency
	cmd := exec.Command("go", "get", dep.Module+"@"+dep.Version)
	cmd.Dir = dm.projectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add dependency %s: %v\n%s", dep.Module, err, string(output))
	}

	fmt.Printf("Added dependency: %s %s\n", dep.Module, dep.Version)
	return nil
}

// DependencyExists checks if a dependency is already in go.mod
func (dm *DependencyManager) DependencyExists(module string) (bool, error) {
	content, err := os.ReadFile(dm.goModPath)
	if err != nil {
		return false, err
	}

	// Simple check: look for module name in go.mod
	return strings.Contains(string(content), module), nil
}

// CheckGoVersion verifies the Go version is compatible
func (dm *DependencyManager) CheckGoVersion(requiredVersion string) error {
	cmd := exec.Command("go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get Go version: %v", err)
	}

	// Extract version from output: "go version go1.21.0 linux/amd64"
	re := regexp.MustCompile(`go(\d+\.\d+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return fmt.Errorf("could not parse Go version from: %s", string(output))
	}

	currentVersion := matches[1]

	// Simple version comparison (assumes format X.Y)
	if !dm.isVersionCompatible(currentVersion, requiredVersion) {
		return fmt.Errorf("Go version %s required, but found %s", requiredVersion, currentVersion)
	}

	return nil
}

// isVersionCompatible checks if current version meets required version
func (dm *DependencyManager) isVersionCompatible(current, required string) bool {
	// Parse versions
	var currMajor, currMinor, reqMajor, reqMinor int
	fmt.Sscanf(current, "%d.%d", &currMajor, &currMinor)
	fmt.Sscanf(required, "%d.%d", &reqMajor, &reqMinor)

	if currMajor > reqMajor {
		return true
	}
	if currMajor == reqMajor && currMinor >= reqMinor {
		return true
	}
	return false
}

// UpdateGoMod runs go mod tidy to update go.mod and go.sum
func (dm *DependencyManager) UpdateGoMod() error {
	if dm.dryRun {
		fmt.Println("[DRY-RUN] Would run: go mod tidy")
		return nil
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dm.projectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %v\n%s", err, string(output))
	}

	fmt.Println("Updated go.mod and go.sum")
	return nil
}

// PrintDependencySuggestions prints suggested optional dependencies
func (dm *DependencyManager) PrintDependencySuggestions(suggestions []Dependency) {
	if len(suggestions) == 0 {
		return
	}

	fmt.Println("\nOPTIONAL DEPENDENCIES:")
	fmt.Println("   The following dependencies might be useful for your feature:")

	for _, dep := range suggestions {
		fmt.Printf("   %s %s\n", dep.Module, dep.Version)
		fmt.Printf("      Reason: %s\n", dep.Reason)
		fmt.Printf("      Install: go get %s@%s\n\n", dep.Module, dep.Version)
	}
}

// VerifyDependencyVersions checks if all dependencies have compatible versions
func (dm *DependencyManager) VerifyDependencyVersions() error {
	cmd := exec.Command("go", "mod", "verify")
	cmd.Dir = dm.projectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("dependency verification failed: %v\n%s", err, string(output))
	}

	fmt.Println("All dependencies verified")
	return nil
}

// GetRequiredDependenciesForFeature returns required dependencies for a feature type
func (dm *DependencyManager) GetRequiredDependenciesForFeature(featureType string, options map[string]bool) []Dependency {
	required := make([]Dependency, 0)
	commonDeps := dm.CommonDependencies()

	// Add dependencies based on feature type
	switch strings.ToLower(featureType) {
	case "grpc":
		required = append(required, commonDeps["grpc"], commonDeps["protobuf"])
	case "auth":
		required = append(required, commonDeps["jwt"], commonDeps["bcrypt"])
	}

	// Add based on options
	if options["validation"] {
		required = append(required, commonDeps["validator"])
	}

	return required
}
