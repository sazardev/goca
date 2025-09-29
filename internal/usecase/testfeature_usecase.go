package usecase

import "github.com/sazardev/goca/internal/domain"

type TestFeatureUseCase interface {
	CreateTestFeature(input CreateTestFeatureInput) (CreateTestFeatureOutput, error)
	GetTestFeature(id int) (*domain.TestFeature, error)
	UpdateTestFeature(id int, input UpdateTestFeatureInput) error
	DeleteTestFeature(id int) error
	ListTestFeatures() (ListTestFeatureOutput, error)
}
