package testing

import (
	"fmt"
	"strings"
)

// TestError represents a test error with detailed context
type TestError struct {
	Type        string
	File        string
	Expected    string
	Actual      string
	Description string
	Severity    ErrorSeverity
}

// ErrorSeverity represents the severity level of a test error
type ErrorSeverity int

const (
	ErrorSeverityLow ErrorSeverity = iota
	ErrorSeverityMedium
	ErrorSeverityHigh
	ErrorSeverityCritical
)

func (s ErrorSeverity) String() string {
	switch s {
	case ErrorSeverityLow:
		return "LOW"
	case ErrorSeverityMedium:
		return "MEDIUM"
	case ErrorSeverityHigh:
		return "HIGH"
	case ErrorSeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// Error implements the error interface
func (e *TestError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Severity, e.Type, e.Description)
}

// DetailedError returns a detailed error message with context
func (e *TestError) DetailedError() string {
	var parts []string

	parts = append(parts, fmt.Sprintf("Error Type: %s", e.Type))
	parts = append(parts, fmt.Sprintf("Severity: %s", e.Severity))

	if e.File != "" {
		parts = append(parts, fmt.Sprintf("File: %s", e.File))
	}

	if e.Description != "" {
		parts = append(parts, fmt.Sprintf("Description: %s", e.Description))
	}

	if e.Expected != "" {
		parts = append(parts, fmt.Sprintf("Expected: %s", e.Expected))
	}

	if e.Actual != "" {
		parts = append(parts, fmt.Sprintf("Actual: %s", e.Actual))
	}

	return strings.Join(parts, "\n")
}

// Error constructors for different types of validation errors

// NewFileError creates a file-related error
func NewFileError(file, operation, description string) *TestError {
	return &TestError{
		Type:        "FILE_ERROR",
		File:        file,
		Description: fmt.Sprintf("%s: %s", operation, description),
		Severity:    ErrorSeverityHigh,
	}
}

// NewStructureError creates a project structure error
func NewStructureError(path, description string) *TestError {
	return &TestError{
		Type:        "STRUCTURE_ERROR",
		File:        path,
		Description: description,
		Severity:    ErrorSeverityHigh,
	}
}

// NewSyntaxError creates a syntax error
func NewSyntaxError(file, description string) *TestError {
	return &TestError{
		Type:        "SYNTAX_ERROR",
		File:        file,
		Description: description,
		Severity:    ErrorSeverityCritical,
	}
}

// NewComplianceError creates a compliance error
func NewComplianceError(file, component, description string) *TestError {
	return &TestError{
		Type:        "COMPLIANCE_ERROR",
		File:        file,
		Description: fmt.Sprintf("%s compliance: %s", component, description),
		Severity:    ErrorSeverityMedium,
	}
}

// NewDependencyError creates a dependency violation error
func NewDependencyError(file, fromLayer, toLayer, description string) *TestError {
	return &TestError{
		Type:        "DEPENDENCY_ERROR",
		File:        file,
		Description: fmt.Sprintf("Invalid dependency from %s to %s: %s", fromLayer, toLayer, description),
		Severity:    ErrorSeverityHigh,
	}
}

// NewLocationError creates a file location error
func NewLocationError(file, expectedLocation, description string) *TestError {
	return &TestError{
		Type:        "LOCATION_ERROR",
		File:        file,
		Expected:    expectedLocation,
		Description: description,
		Severity:    ErrorSeverityMedium,
	}
}

// NewNamingError creates a naming convention error
func NewNamingError(file, element, expected, actual string) *TestError {
	return &TestError{
		Type:        "NAMING_ERROR",
		File:        file,
		Expected:    expected,
		Actual:      actual,
		Description: fmt.Sprintf("Naming convention violation for %s", element),
		Severity:    ErrorSeverityLow,
	}
}

// NewImportError creates an import-related error
func NewImportError(file, importType, description string) *TestError {
	return &TestError{
		Type:        "IMPORT_ERROR",
		File:        file,
		Description: fmt.Sprintf("%s import issue: %s", importType, description),
		Severity:    ErrorSeverityMedium,
	}
}

// NewGenerationError creates a code generation error
func NewGenerationError(component, description string) *TestError {
	return &TestError{
		Type:        "GENERATION_ERROR",
		Description: fmt.Sprintf("Failed to generate %s: %s", component, description),
		Severity:    ErrorSeverityCritical,
	}
}

// NewFlagError creates a CLI flag error
func NewFlagError(flag, expected, actual, description string) *TestError {
	return &TestError{
		Type:        "FLAG_ERROR",
		Expected:    expected,
		Actual:      actual,
		Description: fmt.Sprintf("Flag %s: %s", flag, description),
		Severity:    ErrorSeverityHigh,
	}
}

// ErrorSummary provides a summary of test errors
type ErrorSummary struct {
	Total    int
	Critical int
	High     int
	Medium   int
	Low      int
	ByType   map[string]int
	ByFile   map[string]int
}

// NewErrorSummary creates an error summary from a list of errors
func NewErrorSummary(errors []*TestError) *ErrorSummary {
	summary := &ErrorSummary{
		Total:  len(errors),
		ByType: make(map[string]int),
		ByFile: make(map[string]int),
	}

	for _, err := range errors {
		switch err.Severity {
		case ErrorSeverityCritical:
			summary.Critical++
		case ErrorSeverityHigh:
			summary.High++
		case ErrorSeverityMedium:
			summary.Medium++
		case ErrorSeverityLow:
			summary.Low++
		}

		summary.ByType[err.Type]++
		if err.File != "" {
			summary.ByFile[err.File]++
		}
	}

	return summary
}

// String returns a formatted summary
func (s *ErrorSummary) String() string {
	if s.Total == 0 {
		return "No errors found - all tests passed! ✅"
	}

	var parts []string
	parts = append(parts, fmt.Sprintf("Total Errors: %d", s.Total))

	if s.Critical > 0 {
		parts = append(parts, fmt.Sprintf("Critical: %d", s.Critical))
	}
	if s.High > 0 {
		parts = append(parts, fmt.Sprintf("High: %d", s.High))
	}
	if s.Medium > 0 {
		parts = append(parts, fmt.Sprintf("Medium: %d", s.Medium))
	}
	if s.Low > 0 {
		parts = append(parts, fmt.Sprintf("Low: %d", s.Low))
	}

	return strings.Join(parts, " | ")
}

// HasCriticalErrors returns true if there are critical errors
func (s *ErrorSummary) HasCriticalErrors() bool {
	return s.Critical > 0
}

// HasHighPriorityErrors returns true if there are critical or high severity errors
func (s *ErrorSummary) HasHighPriorityErrors() bool {
	return s.Critical > 0 || s.High > 0
}
