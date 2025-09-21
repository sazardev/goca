package usecase

import "github.com/sazardev/goca/internal/domain"

type TestPerfUseCase interface {
	CreateTestPerf(input CreateTestPerfInput) (CreateTestPerfOutput, error)
	GetTestPerf(id int) (*domain.TestPerf, error)
	UpdateTestPerf(id int, input UpdateTestPerfInput) error
	DeleteTestPerf(id int) error
	ListTestPerfs() (ListTestPerfOutput, error)
}
