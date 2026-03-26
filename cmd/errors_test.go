package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewErrorHandler(t *testing.T) {
	t.Parallel()
	eh := NewErrorHandler()
	require.NotNil(t, eh)
	assert.False(t, eh.TestMode)
}

func TestValidateRequiredFlag(t *testing.T) {
	t.Parallel()
	eh := NewErrorHandler()
	eh.TestMode = true

	cases := []struct {
		name     string
		value    string
		flagName string
		wantErr  bool
	}{
		{name: "populated value", value: "hello", flagName: "name", wantErr: false},
		{name: "empty value", value: "", flagName: "fields", wantErr: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := eh.ValidateRequiredFlag(tc.value, tc.flagName)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.flagName)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandleErrorWithReturn(t *testing.T) {
	t.Parallel()
	eh := NewErrorHandler()
	eh.TestMode = true

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()
		err := eh.HandleErrorWithReturn(nil, "test")
		assert.NoError(t, err)
	})

	t.Run("non-nil error", func(t *testing.T) {
		t.Parallel()
		err := eh.HandleErrorWithReturn(assert.AnError, "context")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context")
	})
}
