package cmd

import (
	"fmt"
	"strings"
)

// CommandValidator centralizes validation logic for all commands
type CommandValidator struct {
	fieldValidator *FieldValidator
	errorHandler   *ErrorHandler
}

// NewCommandValidator creates a new command validator
func NewCommandValidator() *CommandValidator {
	return &CommandValidator{
		fieldValidator: NewFieldValidator(),
		errorHandler:   NewErrorHandler(),
	}
}

// NewTestCommandValidator creates a command validator for testing
func NewTestCommandValidator() *CommandValidator {
	errorHandler := NewErrorHandler()
	errorHandler.TestMode = true
	return &CommandValidator{
		fieldValidator: NewFieldValidator(),
		errorHandler:   errorHandler,
	}
}

// ValidateEntityCommand validates common entity command parameters
func (v *CommandValidator) ValidateEntityCommand(entityName, fields string) error {
	// Validate entity name
	if err := v.fieldValidator.ValidateEntityName(entityName); err != nil {
		return fmt.Errorf("nombre de entidad: %w", err)
	}

	// Validate fields if provided
	if fields != "" {
		if err := v.fieldValidator.ValidateFields(fields); err != nil {
			return fmt.Errorf("campos: %w", err)
		}
	}

	return nil
}

// ValidateFeatureCommand validates feature command parameters
func (v *CommandValidator) ValidateFeatureCommand(featureName, fields, database, handlers string) error {
	// Validate feature name (same as entity)
	if err := v.fieldValidator.ValidateEntityName(featureName); err != nil {
		return fmt.Errorf("nombre de feature: %w", err)
	}

	// Validate required fields
	v.errorHandler.ValidateRequiredFlag(fields, "fields")

	// Validate fields
	if err := v.fieldValidator.ValidateFields(fields); err != nil {
		return fmt.Errorf("fields: %w", err)
	}

	// Validate database if provided
	if database != "" {
		if err := v.fieldValidator.ValidateDatabase(database); err != nil {
			return fmt.Errorf("database: %w", err)
		}
	}

	// Validate handlers if provided
	if handlers != "" {
		if err := v.fieldValidator.ValidateHandlers(handlers); err != nil {
			return fmt.Errorf("handlers: %w", err)
		}
	}

	return nil
}

// ValidateRepositoryCommand validates repository command parameters
func (v *CommandValidator) ValidateRepositoryCommand(entityName, database string) error {
	// Validate entity name
	if err := v.fieldValidator.ValidateEntityName(entityName); err != nil {
		return fmt.Errorf("entity name: %w", err)
	}

	// Validate database if provided
	if database != "" {
		if err := v.fieldValidator.ValidateDatabase(database); err != nil {
			return fmt.Errorf("database: %w", err)
		}
	}

	return nil
}

// ValidateUseCaseCommand validates use case command parameters
func (v *CommandValidator) ValidateUseCaseCommand(usecaseName, entity, operations string) error {
	// Validate use case name
	if err := v.fieldValidator.ValidateEntityName(usecaseName); err != nil {
		return fmt.Errorf("nombre de caso de uso: %w", err)
	}

	// Validate entity name
	if entity == "" {
		return fmt.Errorf("entidad es requerida")
	}
	if err := v.fieldValidator.ValidateEntityName(entity); err != nil {
		return fmt.Errorf("nombre de entidad: %w", err)
	}

	// Validate operations if provided
	if operations != "" {
		if err := v.fieldValidator.ValidateOperations(operations); err != nil {
			return fmt.Errorf("operaciones: %w", err)
		}
	}

	return nil
}

// ValidateHandlerCommand validates handler command parameters
func (v *CommandValidator) ValidateHandlerCommand(entity, handlerType string) error {
	// Validate entity name
	if err := v.fieldValidator.ValidateEntityName(entity); err != nil {
		return fmt.Errorf("nombre de entidad: %w", err)
	}

	// Validate handler type
	if handlerType != "" {
		validHandlers := []string{HandlerHTTP, HandlerGRPC, HandlerCLI, HandlerWorker}
		found := false
		for _, valid := range validHandlers {
			if handlerType == valid {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("tipo de handler no válido: %s. Opciones: %s",
				handlerType, strings.Join(validHandlers, ", "))
		}
	}

	return nil
}
