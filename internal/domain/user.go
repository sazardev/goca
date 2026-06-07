package domain

import (
	"time"
)

type User struct {
	ID        uint       `gorm:"primaryKey;autoIncrement"               json:"id"`
	Name      string     `gorm:"type:varchar(255);not null"             json:"name"                 validate:"required"`
	Email     string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"                validate:"required,email"`
	Age       int        `gorm:"type:integer;not null;default:0"        json:"age"                  validate:"required,gte=0"`
	CreatedAt time.Time  `gorm:"autoCreateTime"                         json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"                         json:"updated_at"`
	DeletedAt *time.Time `gorm:"index"                                  json:"deleted_at,omitempty"`
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
	now := time.Now()
	u.DeletedAt = &now
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}
