package usecase

import (
	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/repository"
)

type testFeatureService struct {
	repo repository.TestFeatureRepository
}

func NewTestFeatureService(repo repository.TestFeatureRepository) TestFeatureUseCase {
	return &testFeatureService{repo: repo}
}

func (t *testFeatureService) Create(input CreateTestFeatureInput) (*CreateTestFeatureOutput, error) {
	testfeature := domain.TestFeature{
		Name:  input.Name,
		Email: input.Email,
		Age:   input.Age,
	}

	if err := testfeature.Validate(); err != nil {
		return nil, err
	}

	if err := t.repo.Save(&testfeature); err != nil {
		return nil, err
	}

	return &CreateTestFeatureOutput{
		TestFeature: testfeature,
		Message:     "TestFeature created successfully",
	}, nil
}

func (t *testFeatureService) GetByID(id uint) (*domain.TestFeature, error) {
	return t.repo.FindByID(int(id))
}

func (t *testFeatureService) Update(id uint, input UpdateTestFeatureInput) (*domain.TestFeature, error) {
	testfeature, err := t.repo.FindByID(int(id))
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		testfeature.Name = *input.Name
	}
	if input.Email != nil {
		testfeature.Email = *input.Email
	}
	if input.Age != nil {
		testfeature.Age = *input.Age
	}

	err = t.repo.Update(testfeature)
	if err != nil {
		return nil, err
	}

	return testfeature, nil
}

func (t *testFeatureService) Delete(id uint) error {
	return t.repo.Delete(int(id))
}

func (t *testFeatureService) List() (*ListTestFeatureOutput, error) {
	testfeatures, err := t.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return &ListTestFeatureOutput{
		TestFeatures: testfeatures,
		Total:        len(testfeatures),
		Message:      "TestFeatures listed successfully",
	}, nil
}
