package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
	Email     string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,email"`
	Age       int            `json:"age" gorm:"type:integer;not null;default:0" validate:"required,gte=0"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
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

func (u *User) SoftDelete() {
	u.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt.Valid
}
