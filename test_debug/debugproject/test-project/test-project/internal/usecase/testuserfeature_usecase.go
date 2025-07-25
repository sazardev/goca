package usecase

import "github.com/test/debugproject/internal/domain"

type TestUserFeatureUseCase interface {
	CreateTestUserFeature(input CreateTestUserFeatureInput) (CreateTestUserFeatureOutput, error)
	GetTestUserFeature(id int) (*domain.TestUserFeature, error)
	UpdateTestUserFeature(id int, input UpdateTestUserFeatureInput) error
	DeleteTestUserFeature(id int) error
	ListTestUserFeatures() (ListTestUserFeatureOutput, error)
}
