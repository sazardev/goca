package domain

import "errors"

var (
	ErrInvalidTestUserFeatureData  = errors.New("invalid testuserfeature data")
	ErrInvalidTestUserFeatureName  = errors.New("testuserfeature name is invalid")
	ErrInvalidTestUserFeatureEmail = errors.New("testuserfeature email is invalid")
)
