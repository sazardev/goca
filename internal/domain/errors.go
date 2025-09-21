package domain

import "errors"

var (
	ErrInvalidTestPerfData        = errors.New("datos de testperf inválidos")
	ErrInvalidTestPerfName        = errors.New("el nombre es requerido")
	ErrInvalidTestPerfNameLength  = errors.New("el nombre debe tener entre 2 y 100 caracteres")
	ErrInvalidTestPerfEmail       = errors.New("el email es requerido")
	ErrInvalidTestPerfEmailFormat = errors.New("formato de el email inválido")
	ErrInvalidTestPerfAge         = errors.New("la edad es requerido")
	ErrInvalidTestPerfAgeRange    = errors.New("la edad debe ser mayor a 0")
	ErrInvalidTestPerfScore       = errors.New("el campo score es requerido")
	ErrInvalidTestPerfScoreRange  = errors.New("el campo score debe ser un número positivo")
	ErrInvalidTestPerfActive      = errors.New("el campo active es requerido")
)
