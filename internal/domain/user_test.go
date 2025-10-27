package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUser_Validate tests the Validate method with various scenarios
func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid entity",
			user: User{
				Name:  "John Doe",
				Email: "test@example.com",
				Age:   25,
			},
			wantErr: false,
		},
		{
			name: "invalid user - empty string",
			user: User{
				Name:  "",
				Email: "test@example.com",
				Age:   25,
			},
			wantErr: true,
			errMsg:  "name",
		},
		{
			name: "invalid user - empty string",
			user: User{
				Name:  "John Doe",
				Email: "",
				Age:   25,
			},
			wantErr: true,
			errMsg:  "email",
		},
		{
			name: "invalid user - negative number",
			user: User{
				Name:  "John Doe",
				Email: "test@example.com",
				Age:   -1,
			},
			wantErr: true,
			errMsg:  "age",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUser_Initialization tests entity field initialization
func TestUser_Initialization(t *testing.T) {
	user := &User{
		Name:  "John Doe",
		Email: "test@example.com",
		Age:   25,
	}

	assert.Equal(t, "John Doe", user.Name, "Name should be set correctly")
	assert.Equal(t, "test@example.com", user.Email, "Email should be set correctly")
	assert.Equal(t, 25, user.Age, "Age should be set correctly")
}

// TestUser_Name_EdgeCases tests edge cases for Name field
func TestUser_Name_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid value", value: "Valid Name", wantErr: false},
		{name: "empty string", value: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Name:  tt.value,
				Email: "test@example.com",
				Age:   25,
			}

			err := user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUser_Email_EdgeCases tests edge cases for Email field
func TestUser_Email_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid value", value: "Valid Name", wantErr: false},
		{name: "empty string", value: "", wantErr: true},
		{name: "valid email", value: "test@example.com", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Name:  "John Doe",
				Email: tt.value,
				Age:   25,
			}

			err := user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUser_Age_NumericValidation tests numeric validation for Age
func TestUser_Age_NumericValidation(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{name: "positive value", value: 10, wantErr: false},
		{name: "zero value", value: 0, wantErr: false},
		{name: "negative value", value: -1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Name:  "John Doe",
				Email: "test@example.com",
				Age:   tt.value,
			}

			err := user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
