package domain

import "errors"

var (
	// User errors
	ErrInvalidUserData       = errors.New("invalid user data")
	ErrInvalidUserName       = errors.New("name is required")
	ErrInvalidUserNameLength = errors.New("name must be between 2 and 100 characters")
	ErrInvalidUserEmail      = errors.New("email is required")
	ErrInvalidUserAge        = errors.New("age is required")
	ErrInvalidUserAgeRange   = errors.New("age must be greater than 0")

	// Product errors
	ErrInvalidProductData        = errors.New("invalid product data")
	ErrInvalidProductName        = errors.New("name is required")
	ErrInvalidProductNameLength  = errors.New("name must be between 2 and 100 characters")
	ErrInvalidProductPrice       = errors.New("price is required")
	ErrInvalidProductPriceRange  = errors.New("price must be greater than 0 and less than 999,999,999.99")
	ErrInvalidProductDescription = errors.New("description is required")

	// Order errors
	ErrInvalidOrderData             = errors.New("invalid order data")
	ErrInvalidOrderCustomer_id      = errors.New("customer_id is required")
	ErrInvalidOrderCustomer_idRange = errors.New("customer_id must be a positive number")
	ErrInvalidOrderTotal            = errors.New("total is required")
	ErrInvalidOrderTotalRange       = errors.New("total must be a positive number")
	ErrInvalidOrderStatus           = errors.New("status is required")
)
