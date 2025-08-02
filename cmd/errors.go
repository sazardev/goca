package cmd

import (
	"fmt"
	"os"
)

// ErrorHandler centralizes error handling for the CLI
type ErrorHandler struct{}

// NewErrorHandler creates a new error handler
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandleError handles errors with consistent formatting and exit behavior
func (e *ErrorHandler) HandleError(err error, context string) {
	if err != nil {
		fmt.Printf("❌ Error en %s: %v\n", context, err)
		os.Exit(1)
	}
}

// HandleValidationError handles validation errors with specific formatting
func (e *ErrorHandler) HandleValidationError(err error, field string) {
	if err != nil {
		fmt.Printf("❌ Error de validación en %s: %v\n", field, err)
		os.Exit(1)
	}
}

// HandleWarning handles warnings without exiting
func (e *ErrorHandler) HandleWarning(message string, context string) {
	fmt.Printf("⚠️  Advertencia en %s: %s\n", context, message)
}

// HandleSuccess handles success messages with consistent formatting
func (e *ErrorHandler) HandleSuccess(message string) {
	fmt.Printf("✅ %s\n", message)
}

// HandleInfo handles informational messages
func (e *ErrorHandler) HandleInfo(message string) {
	fmt.Printf("ℹ️  %s\n", message)
}

// ValidateRequiredFlag checks if a required flag is provided
func (e *ErrorHandler) ValidateRequiredFlag(value string, flagName string) {
	if value == "" {
		fmt.Printf("❌ Error: --%s flag es requerido\n", flagName)
		os.Exit(1)
	}
}
