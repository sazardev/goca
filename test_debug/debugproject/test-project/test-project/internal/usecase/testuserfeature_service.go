package usecase

import (
	"github.com/test/debugproject/internal/domain"
	"github.com/test/debugproject/internal/messages"
	"github.com/test/debugproject/internal/repository"
)

type testuserfeatureService struct {
	repo repository.TestUserFeatureRepository
}

func NewTestuserfeatureService(repo repository.TestUserFeatureRepository) TestUserFeatureUseCase {
	return &testuserfeatureService{repo: repo}
}

func (t *testuserfeatureService) CreateTestUserFeature(input CreateTestUserFeatureInput) (CreateTestUserFeatureOutput, error) {
	testuserfeature := domain.TestUserFeature{
		// TODO: Map fields from input to entity
		// Example: Name: input.Name,
	}

	if err := testuserfeature.Validate(); err != nil {
		return CreateTestUserFeatureOutput{}, err
	}

	if err := t.repo.Save(&testuserfeature); err != nil {
		return CreateTestUserFeatureOutput{}, err
	}

	return CreateTestUserFeatureOutput{
		TestUserFeature: testuserfeature,
		Message:         messages.TestUserFeatureCreatedSuccessfully,
	}, nil
}

func (t *testuserfeatureService) GetTestUserFeature(id int) (*domain.TestUserFeature, error) {
	return t.repo.FindByID(id)
}

func (t *testuserfeatureService) UpdateTestUserFeature(id int, input UpdateTestUserFeatureInput) error {
	testuserfeature, err := t.repo.FindByID(id)
	if err != nil {
		return err
	}

	// TODO: Update fields based on your entity
	// Example: if input.Name != "" {
	//     entity.Name = input.Name
	// }

	return t.repo.Update(testuserfeature)
}

func (t *testuserfeatureService) DeleteTestUserFeature(id int) error {
	return t.repo.Delete(id)
}

func (t *testuserfeatureService) ListTestUserFeatures() (ListTestUserFeatureOutput, error) {
	testuserfeatures, err := t.repo.FindAll()
	if err != nil {
		return ListTestUserFeatureOutput{}, err
	}

	return ListTestUserFeatureOutput{
		TestUserFeatures: testuserfeatures,
		Total:            len(testuserfeatures),
		Message:          messages.TestUserFeaturesListedSuccessfully,
	}, nil
}
