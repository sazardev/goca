package domain

type TestFeature struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name" gorm:"type:varchar(255);not null"`
	Email string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Age   int    `json:"age" gorm:"type:integer;not null;default:0"`
}

func (t *TestFeature) Validate() error {
	if t.Name == "" {
		return ErrInvalidTestFeatureName
	}
	if t.Email == "" {
		return ErrInvalidTestFeatureEmail
	}
	if t.Age < 0 {
		return ErrInvalidTestFeatureAge
	}
	return nil
}
