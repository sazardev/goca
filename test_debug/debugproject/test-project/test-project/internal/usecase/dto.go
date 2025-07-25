package usecase

import (
	"github.com/test/debugproject/internal/domain"
)

type CreateTestUserFeatureInput struct {
	// TODO: Add specific fields for your TestUserFeature entity
	// Example: Name string `json:"name" validate:"required,min=2"`
}

type CreateTestUserFeatureOutput struct {
	TestUserFeature domain.TestUserFeature `json:"testuserfeature"`
	Message         string                 `json:"message"`
}

type UpdateTestUserFeatureInput struct {
	// TODO: Add specific fields for your TestUserFeature entity
	// Example: Name string `json:"name,omitempty" validate:"omitempty,min=2"`
}

type ListTestUserFeatureOutput struct {
	TestUserFeatures []domain.TestUserFeature `json:"testuserfeatures"`
	Total            int                      `json:"total"`
	Message          string                   `json:"message"`
}
