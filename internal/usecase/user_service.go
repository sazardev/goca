package usecase

import (
	"github.com/sazardev/goca/internal/domain"
	"github.com/sazardev/goca/internal/messages"
	"github.com/sazardev/goca/internal/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserUseCase {
	return &userService{repo: repo}
}

func (u *userService) CreateUser(input CreateUserInput) (CreateUserOutput, error) {
	user := domain.User{
		Name:  input.Name,
		Email: input.Email,
		Age:   input.Age,
	}

	if err := user.Validate(); err != nil {
		return CreateUserOutput{}, err
	}

	if err := u.repo.Save(&user); err != nil {
		return CreateUserOutput{}, err
	}

	return CreateUserOutput{
		User:    user,
		Message: messages.UserCreatedSuccessfully,
	}, nil
}

func (u *userService) GetUser(id int) (*domain.User, error) {
	return u.repo.FindByID(id)
}

func (u *userService) UpdateUser(id int, input UpdateUserInput) error {
	user, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Actualizar campos según tu entidad
	if input.Nombre != "" {
		user.Nombre = input.Nombre
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	// Agregar más campos según necesites

	return u.repo.Update(user)
}

func (u *userService) DeleteUser(id int) error {
	return u.repo.Delete(id)
}

func (u *userService) ListUsers() (ListUserOutput, error) {
	users, err := u.repo.FindAll()
	if err != nil {
		return ListUserOutput{}, err
	}

	return ListUserOutput{
		Users:   users,
		Total:   len(users),
		Message: messages.UsersListedSuccessfully,
	}, nil
}
