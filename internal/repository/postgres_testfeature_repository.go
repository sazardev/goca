package repository

import (
	"github.com/sazardev/goca/internal/domain"

	"gorm.io/gorm"
)

type postgresTestFeatureRepository struct {
	db *gorm.DB
}

func NewPostgresTestFeatureRepository(db *gorm.DB) TestFeatureRepository {
	return &postgresTestFeatureRepository{
		db: db,
	}
}

func (p *postgresTestFeatureRepository) Save(testfeature *domain.TestFeature) error {
	result := p.db.Create(testfeature)
	return result.Error
}

func (p *postgresTestFeatureRepository) FindByID(id int) (*domain.TestFeature, error) {
	testfeature := &domain.TestFeature{}
	result := p.db.First(testfeature, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return testfeature, nil
}

func (p *postgresTestFeatureRepository) Update(testfeature *domain.TestFeature) error {
	result := p.db.Save(testfeature)
	return result.Error
}

func (p *postgresTestFeatureRepository) Delete(id int) error {
	result := p.db.Delete(&domain.TestFeature{}, id)
	return result.Error
}

func (p *postgresTestFeatureRepository) FindAll() ([]domain.TestFeature, error) {
	var testfeatures []domain.TestFeature
	result := p.db.Find(&testfeatures)
	if result.Error != nil {
		return nil, result.Error
	}
	return testfeatures, nil
}

func (p *postgresTestFeatureRepository) FindByName(name string) (*domain.TestFeature, error) {
	testfeature := &domain.TestFeature{}
	result := p.db.Where("name = ?", name).First(testfeature)
	if result.Error != nil {
		return nil, result.Error
	}
	return testfeature, nil
}

func (p *postgresTestFeatureRepository) FindByEmail(email string) (*domain.TestFeature, error) {
	testfeature := &domain.TestFeature{}
	result := p.db.Where("email = ?", email).First(testfeature)
	if result.Error != nil {
		return nil, result.Error
	}
	return testfeature, nil
}

func (p *postgresTestFeatureRepository) FindByAge(age int) (*domain.TestFeature, error) {
	testfeature := &domain.TestFeature{}
	result := p.db.Where("age = ?", age).First(testfeature)
	if result.Error != nil {
		return nil, result.Error
	}
	return testfeature, nil
}
