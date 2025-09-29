package usecase

import (
	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/messages"
	"github.com/sazardev/goca/internal/repository"
)

type testfeatureService struct {
	repo repository.TestFeatureRepository
}

func NewTestfeatureService(repo repository.TestFeatureRepository) TestFeatureUseCase {
	return &testfeatureService{repo: repo}
}

func (t *testfeatureService) CreateTestFeature(input CreateTestFeatureInput) (CreateTestFeatureOutput, error) {
	testfeature := domain.TestFeature{
		Name:  input.Name,
		Email: input.Email,
		Age:   input.Age,
	}

	if err := testfeature.Validate(); err != nil {
		return CreateTestFeatureOutput{}, err
	}

	if err := t.repo.Save(&testfeature); err != nil {
		return CreateTestFeatureOutput{}, err
	}

	return CreateTestFeatureOutput{
		TestFeature: testfeature,
		Message:     messages.TestFeatureCreatedSuccessfully,
	}, nil
}

func (t *testfeatureService) GetTestFeature(id int) (*domain.TestFeature, error) {
	return t.repo.FindByID(id)
}

func (t *testfeatureService) UpdateTestFeature(id int, input UpdateTestFeatureInput) error {
	testfeature, err := t.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Actualizar campos según tu entidad
	if input.Nombre != "" {
		testfeature.Nombre = input.Nombre
	}
	if input.Email != "" {
		testfeature.Email = input.Email
	}
	// Agregar más campos según necesites

	return t.repo.Update(testfeature)
}

func (t *testfeatureService) DeleteTestFeature(id int) error {
	return t.repo.Delete(id)
}

func (t *testfeatureService) ListTestFeatures() (ListTestFeatureOutput, error) {
	testfeatures, err := t.repo.FindAll()
	if err != nil {
		return ListTestFeatureOutput{}, err
	}

	return ListTestFeatureOutput{
		TestFeatures: testfeatures,
		Total:        len(testfeatures),
		Message:      messages.TestFeaturesListedSuccessfully,
	}, nil
}
