package repository

import (
	"github.com/sazardev/goca/internal/domain"

	"gorm.io/gorm"
)

type postgresTestPerfRepository struct {
	db *gorm.DB
}

func NewPostgresTestPerfRepository(db *gorm.DB) TestPerfRepository {
	return &postgresTestPerfRepository{
		db: db,
	}
}

func (p *postgresTestPerfRepository) Save(testperf *domain.TestPerf) error {
	result := p.db.Create(testperf)
	return result.Error
}

func (p *postgresTestPerfRepository) FindByID(id int) (*domain.TestPerf, error) {
	testperf := &domain.TestPerf{}
	result := p.db.First(testperf, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return testperf, nil
}

func (p *postgresTestPerfRepository) Update(testperf *domain.TestPerf) error {
	result := p.db.Save(testperf)
	return result.Error
}

func (p *postgresTestPerfRepository) Delete(id int) error {
	result := p.db.Delete(&domain.TestPerf{}, id)
	return result.Error
}

func (p *postgresTestPerfRepository) FindAll() ([]domain.TestPerf, error) {
	var testperfs []domain.TestPerf
	result := p.db.Find(&testperfs)
	if result.Error != nil {
		return nil, result.Error
	}
	return testperfs, nil
}

func (p *postgresTestPerfRepository) FindByName(name string) (*domain.TestPerf, error) {
	testperf := &domain.TestPerf{}
	result := p.db.Where("name = ?", name).First(testperf)
	if result.Error != nil {
		return nil, result.Error
	}
	return testperf, nil
}

func (p *postgresTestPerfRepository) FindByEmail(email string) (*domain.TestPerf, error) {
	testperf := &domain.TestPerf{}
	result := p.db.Where("email = ?", email).First(testperf)
	if result.Error != nil {
		return nil, result.Error
	}
	return testperf, nil
}

func (p *postgresTestPerfRepository) FindByAge(age int) (*domain.TestPerf, error) {
	testperf := &domain.TestPerf{}
	result := p.db.Where("age = ?", age).First(testperf)
	if result.Error != nil {
		return nil, result.Error
	}
	return testperf, nil
}
