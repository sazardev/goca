package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleError_TestMode(t *testing.T) {
	t.Parallel()
	eh := &ErrorHandler{TestMode: true}

	// Should not panic or os.Exit in test mode
	eh.HandleError(fmt.Errorf("test error"), "test context")
}

func TestHandleError_NilError(t *testing.T) {
	t.Parallel()
	eh := &ErrorHandler{TestMode: true}

	// nil error should be a no-op
	eh.HandleError(nil, "test context")
}

func TestHandleValidationError_TestMode(t *testing.T) {
	t.Parallel()
	eh := &ErrorHandler{TestMode: true}

	eh.HandleValidationError(fmt.Errorf("invalid"), "name")
}

func TestHandleValidationError_NilError(t *testing.T) {
	t.Parallel()
	eh := &ErrorHandler{TestMode: true}

	eh.HandleValidationError(nil, "name")
}

func TestHandleWarning(t *testing.T) {
	t.Parallel()
	eh := &ErrorHandler{TestMode: true}

	// Should not panic (uses fmt.Printf when ui is nil)
	eh.HandleWarning("something warned", "test context")
}

func TestHandleSuccess(t *testing.T) {
	t.Parallel()
	eh := &ErrorHandler{TestMode: true}

	eh.HandleSuccess("it worked")
}

func TestHandleInfo(t *testing.T) {
	t.Parallel()
	eh := &ErrorHandler{TestMode: true}

	eh.HandleInfo("some info")
}

func TestValidateRequiredFlag_NotCovered(t *testing.T) {
	t.Parallel()
	eh := &ErrorHandler{TestMode: true}

	// ValidateRequiredFlag with empty value should return error
	err := eh.ValidateRequiredFlag("", "name")
	assert.Error(t, err)

	// ValidateRequiredFlag with value should return nil
	err = eh.ValidateRequiredFlag("value", "name")
	assert.NoError(t, err)
}
