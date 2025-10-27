package domain

import "errors"

var (
	// User errors
	ErrInvalidUserData       = errors.New("datos de user inválidos")
	ErrInvalidUserName       = errors.New("el nombre es requerido")
	ErrInvalidUserNameLength = errors.New("el nombre debe tener entre 2 y 100 caracteres")
	ErrInvalidUserEmail      = errors.New("el email es requerido")
	ErrInvalidUserAge        = errors.New("la edad es requerido")
	ErrInvalidUserAgeRange   = errors.New("la edad debe ser mayor a 0")

	// Product errors
	ErrInvalidProductData        = errors.New("datos de product inválidos")
	ErrInvalidProductName        = errors.New("el nombre es requerido")
	ErrInvalidProductNameLength  = errors.New("el nombre debe tener entre 2 y 100 caracteres")
	ErrInvalidProductPrice       = errors.New("el precio es requerido")
	ErrInvalidProductPriceRange  = errors.New("el precio debe ser mayor a 0 y menor a 999,999,999.99")
	ErrInvalidProductDescription = errors.New("la descripción es requerido")

	// Order errors
	ErrInvalidOrderData             = errors.New("datos de order inválidos")
	ErrInvalidOrderCustomer_id      = errors.New("el campo customer_id es requerido")
	ErrInvalidOrderCustomer_idRange = errors.New("el campo customer_id debe ser un número positivo")
	ErrInvalidOrderTotal            = errors.New("el campo total es requerido")
	ErrInvalidOrderTotalRange       = errors.New("el campo total debe ser un número positivo")
	ErrInvalidOrderStatus           = errors.New("el estado es requerido")
)
