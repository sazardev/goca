package cmd

import (
	"fmt"
	"os"
)

// ErrorHandler centralizes error handling for the CLI
type ErrorHandler struct {
	TestMode bool // Set to true during testing to avoid os.Exit()
}

// NewErrorHandler creates a new error handler
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		TestMode: false,
	}
}

// HandleError handles errors with consistent formatting and exit behavior
func (e *ErrorHandler) HandleError(err error, context string) {
	if err != nil {
		fmt.Printf("Error in %s: %v\n", context, err)
		if !e.TestMode {
			os.Exit(1)
		}
	}
}

// HandleValidationError handles validation errors with specific formatting
func (e *ErrorHandler) HandleValidationError(err error, field string) {
	if err != nil {
		fmt.Printf("Validation error in %s: %v\n", field, err)
		if !e.TestMode {
			os.Exit(1)
		}
	}
}

// HandleWarning handles warnings without exiting
func (e *ErrorHandler) HandleWarning(message string, context string) {
	fmt.Printf("Warning in %s: %s\n", context, message)
}

// HandleSuccess handles success messages with consistent formatting
func (e *ErrorHandler) HandleSuccess(message string) {
	fmt.Printf("%s\n", message)
}

// HandleInfo handles informational messages
func (e *ErrorHandler) HandleInfo(message string) {
	fmt.Printf("Info: %s\n", message)
}

// ValidateRequiredFlag checks if a required flag is provided
func (e *ErrorHandler) ValidateRequiredFlag(value string, flagName string) error {
	if value == "" {
		err := fmt.Errorf("--%s flag is required", flagName)
		fmt.Printf("Error: %v\n", err)
		if !e.TestMode {
			os.Exit(1)
		}
		return err
	}
	return nil
}

// HandleErrorWithReturn handles errors and returns them for testing
func (e *ErrorHandler) HandleErrorWithReturn(err error, context string) error {
	if err != nil {
		fmt.Printf("Error in %s: %v\n", context, err)
		if !e.TestMode {
			os.Exit(1)
		}
		return fmt.Errorf("%s: %w", context, err)
	}
	return nil
}
