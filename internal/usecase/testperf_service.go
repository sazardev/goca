package usecase

import (
	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/messages"
	"github.com/sazardev/goca/internal/repository"
)

type testperfService struct {
	repo repository.TestPerfRepository
}

func NewTestperfService(repo repository.TestPerfRepository) TestPerfUseCase {
	return &testperfService{repo: repo}
}

func (t *testperfService) CreateTestPerf(input CreateTestPerfInput) (CreateTestPerfOutput, error) {
	testperf := domain.TestPerf{
		Name:   input.Name,
		Email:  input.Email,
		Age:    input.Age,
		Score:  input.Score,
		Active: input.Active,
	}

	if err := testperf.Validate(); err != nil {
		return CreateTestPerfOutput{}, err
	}

	if err := t.repo.Save(&testperf); err != nil {
		return CreateTestPerfOutput{}, err
	}

	return CreateTestPerfOutput{
		TestPerf: testperf,
		Message:  messages.TestPerfCreatedSuccessfully,
	}, nil
}

func (t *testperfService) GetTestPerf(id int) (*domain.TestPerf, error) {
	return t.repo.FindByID(id)
}

func (t *testperfService) UpdateTestPerf(id int, input UpdateTestPerfInput) error {
	testperf, err := t.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Actualizar campos según tu entidad
	if input.Nombre != "" {
		testperf.Nombre = input.Nombre
	}
	if input.Email != "" {
		testperf.Email = input.Email
	}
	// Agregar más campos según necesites

	return t.repo.Update(testperf)
}

func (t *testperfService) DeleteTestPerf(id int) error {
	return t.repo.Delete(id)
}

func (t *testperfService) ListTestPerfs() (ListTestPerfOutput, error) {
	testperfs, err := t.repo.FindAll()
	if err != nil {
		return ListTestPerfOutput{}, err
	}

	return ListTestPerfOutput{
		TestPerfs: testperfs,
		Total:     len(testperfs),
		Message:   messages.TestPerfsListedSuccessfully,
	}, nil
}
