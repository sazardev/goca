package domain

type TestUserFeature struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (t *TestUserFeature) Validate() error {
	if t.Name == "" {
		return ErrInvalidTestUserFeatureName
	}
	if t.Email == "" {
		return ErrInvalidTestUserFeatureEmail
	}
	return nil
}
