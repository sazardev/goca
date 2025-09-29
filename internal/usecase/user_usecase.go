package usecase

import "github.com/sazardev/goca/internal/domain"

type UserUseCase interface {
	CreateUser(input CreateUserInput) (CreateUserOutput, error)
	GetUser(id int) (*domain.User, error)
	UpdateUser(id int, input UpdateUserInput) error
	DeleteUser(id int) error
	ListUsers() (ListUserOutput, error)
}
