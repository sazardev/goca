package repository

import "github.com/test/debugproject/internal/domain"

type TestUserFeatureRepository interface {
	Save(testuserfeature *domain.TestUserFeature) error
	FindByID(id int) (*domain.TestUserFeature, error)
	FindByEmail(email string) (*domain.TestUserFeature, error)
	Update(testuserfeature *domain.TestUserFeature) error
	Delete(id int) error
	FindAll() ([]domain.TestUserFeature, error)
}
