package domain

import "errors"

var (
	ErrInvalidProductData       = errors.New("datos de product inválidos")
	ErrInvalidProductName       = errors.New("el nombre es requerido")
	ErrInvalidProductNameLength = errors.New("el nombre debe tener entre 2 y 100 caracteres")
	ErrInvalidProductPrice      = errors.New("el precio es requerido")
	ErrInvalidProductPriceRange = errors.New("el precio debe ser mayor a 0 y menor a 999,999,999.99")
)
