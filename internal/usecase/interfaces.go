package usecase

import "github.com/sazardev/goca/internal/domain"

type TestPerfRepository interface {
	Save(testperf *domain.TestPerf) error
	FindByID(id int) (*domain.TestPerf, error)
	FindByEmail(email string) (*domain.TestPerf, error)
	Update(testperf *domain.TestPerf) error
	Delete(id int) error
	FindAll() ([]domain.TestPerf, error)
}
