package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommandValidator(t *testing.T) {
	t.Parallel()
	v := NewCommandValidator()
	require.NotNil(t, v)
	require.NotNil(t, v.fieldValidator)
	require.NotNil(t, v.errorHandler)
	assert.False(t, v.errorHandler.TestMode)
}

func TestNewTestCommandValidator(t *testing.T) {
	t.Parallel()
	v := NewTestCommandValidator()
	require.NotNil(t, v)
	assert.True(t, v.errorHandler.TestMode)
}

func TestValidateEntityCommand(t *testing.T) {
	t.Parallel()
	v := NewTestCommandValidator()

	cases := []struct {
		name    string
		entity  string
		fields  string
		wantErr bool
	}{
		{name: "valid entity no fields", entity: "User", fields: "", wantErr: false},
		{name: "valid entity with fields", entity: "Product", fields: "Name:string,Price:float64", wantErr: false},
		{name: "empty entity name", entity: "", fields: "", wantErr: true},
		{name: "invalid entity name with slash", entity: "foo/bar", fields: "", wantErr: true},
		{name: "invalid field type", entity: "User", fields: "Name:invalid_type_xyz", wantErr: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := v.ValidateEntityCommand(tc.entity, tc.fields)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateRepositoryCommand(t *testing.T) {
	t.Parallel()
	v := NewTestCommandValidator()

	cases := []struct {
		name     string
		entity   string
		database string
		wantErr  bool
	}{
		{name: "valid postgres", entity: "User", database: "postgres", wantErr: false},
		{name: "valid mysql", entity: "Product", database: "mysql", wantErr: false},
		{name: "valid sqlite", entity: "Order", database: "sqlite", wantErr: false},
		{name: "empty database ok", entity: "User", database: "", wantErr: false},
		{name: "invalid database", entity: "User", database: "oracle", wantErr: true},
		{name: "empty entity", entity: "", database: "postgres", wantErr: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := v.ValidateRepositoryCommand(tc.entity, tc.database)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUseCaseCommand(t *testing.T) {
	t.Parallel()
	v := NewTestCommandValidator()

	cases := []struct {
		name       string
		usecase    string
		entity     string
		operations string
		wantErr    bool
	}{
		{name: "valid full", usecase: "UserService", entity: "User", operations: "create,read,update,delete", wantErr: false},
		{name: "valid empty ops", usecase: "UserService", entity: "User", operations: "", wantErr: false},
		{name: "empty usecase", usecase: "", entity: "User", operations: "", wantErr: true},
		{name: "empty entity", usecase: "UserService", entity: "", operations: "", wantErr: true},
		{name: "invalid operations", usecase: "UserService", entity: "User", operations: "fly,swim", wantErr: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := v.ValidateUseCaseCommand(tc.usecase, tc.entity, tc.operations)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateHandlerCommand(t *testing.T) {
	t.Parallel()
	v := NewTestCommandValidator()

	cases := []struct {
		name        string
		entity      string
		handlerType string
		wantErr     bool
	}{
		{name: "valid http", entity: "User", handlerType: "http", wantErr: false},
		{name: "valid grpc", entity: "User", handlerType: "grpc", wantErr: false},
		{name: "valid cli", entity: "User", handlerType: "cli", wantErr: false},
		{name: "valid worker", entity: "User", handlerType: "worker", wantErr: false},
		{name: "empty handler ok", entity: "User", handlerType: "", wantErr: false},
		{name: "invalid handler", entity: "User", handlerType: "soap", wantErr: true},
		{name: "empty entity", entity: "", handlerType: "http", wantErr: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := v.ValidateHandlerCommand(tc.entity, tc.handlerType)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFeatureCommand(t *testing.T) {
	t.Parallel()
	v := NewTestCommandValidator()

	cases := []struct {
		name     string
		feature  string
		fields   string
		database string
		handlers string
		wantErr  bool
	}{
		{name: "valid full", feature: "User", fields: "Name:string,Email:string", database: "postgres", handlers: "http", wantErr: false},
		{name: "valid no db no handlers", feature: "Product", fields: "Name:string", database: "", handlers: "", wantErr: false},
		{name: "empty feature name", feature: "", fields: "Name:string", database: "", handlers: "", wantErr: true},
		{name: "invalid database", feature: "User", fields: "Name:string", database: "oracle", handlers: "", wantErr: true},
		{name: "invalid handlers", feature: "User", fields: "Name:string", database: "", handlers: "soap", wantErr: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := v.ValidateFeatureCommand(tc.feature, tc.fields, tc.database, tc.handlers)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
