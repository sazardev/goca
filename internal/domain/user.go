package domain

type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name" gorm:"type:varchar(255);not null"`
	Email string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Age   int    `json:"age" gorm:"type:integer;not null;default:0"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrInvalidUserName
	}
	if u.Email == "" {
		return ErrInvalidUserEmail
	}
	if u.Age < 0 {
		return ErrInvalidUserAge
	}
	return nil
}
