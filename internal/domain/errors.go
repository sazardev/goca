// Package domain contains domain-specific errors and business rule violations.
// This package defines all error constants used throughout the domain layer.
package domain

import "errors"

var (
	// ErrInvalidTestPerfData represents invalid test performance data error.
	ErrInvalidTestPerfData = errors.New("datos de testperf inválidos")
	// ErrInvalidTestPerfName represents missing name error.
	ErrInvalidTestPerfName = errors.New("el nombre es requerido")
	// ErrInvalidTestPerfNameLength represents invalid name length error.
	ErrInvalidTestPerfNameLength = errors.New("el nombre debe tener entre 2 y 100 caracteres")
	// ErrInvalidTestPerfEmail represents missing email error.
	ErrInvalidTestPerfEmail = errors.New("el email es requerido")
	// ErrInvalidTestPerfEmailFormat represents invalid email format error.
	ErrInvalidTestPerfEmailFormat = errors.New("formato de el email inválido")
	// ErrInvalidTestPerfAge represents missing age error.
	ErrInvalidTestPerfAge = errors.New("la edad es requerido")
	// ErrInvalidTestPerfAgeRange represents invalid age range error.
	ErrInvalidTestPerfAgeRange = errors.New("la edad debe ser mayor a 0")
	// ErrInvalidTestPerfScore represents missing score error.
	ErrInvalidTestPerfScore = errors.New("el campo score es requerido")
	// ErrInvalidTestPerfScoreRange represents invalid score range error.
	ErrInvalidTestPerfScoreRange = errors.New("el campo score debe ser un número positivo")
	// ErrInvalidTestPerfActive represents missing active field error.
	ErrInvalidTestPerfActive = errors.New("el campo active es requerido")
)
