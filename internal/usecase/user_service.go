package usecase

import (
"github.com/sazardev/goca/internal/domain"
"github.com/sazardev/goca/internal/repository"
)

type userService struct {
repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserUseCase {
return &userService{repo: repo}
}

func (u *userService) Create(input CreateUserInput) (*CreateUserOutput, error) {
user := domain.User{
Name:  input.Name,
Email: input.Email,
Age:   input.Age,
}

if err := user.Validate(); err != nil {
return nil, err
}

if err := u.repo.Save(&user); err != nil {
return nil, err
}

return &CreateUserOutput{
User:    user,
Message: "User created successfully",
}, nil
}

func (u *userService) GetByID(id uint) (*domain.User, error) {
return u.repo.FindByID(int(id))
}

func (u *userService) Update(id uint, input UpdateUserInput) (*domain.User, error) {
user, err := u.repo.FindByID(int(id))
if err != nil {
return nil, err
}

if input.Name != nil {
user.Name = *input.Name
}
if input.Email != nil {
user.Email = *input.Email
}
if input.Age != nil {
user.Age = *input.Age
}

err = u.repo.Update(user)
if err != nil {
return nil, err
}

return user, nil
}

func (u *userService) Delete(id uint) error {
return u.repo.Delete(int(id))
}

func (u *userService) List() (*ListUserOutput, error) {
users, err := u.repo.FindAll()
if err != nil {
return nil, err
}

return &ListUserOutput{
Users:   users,
Total:   len(users),
Message: "Users listed successfully",
}, nil
}
