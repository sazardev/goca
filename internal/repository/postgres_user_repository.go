package repository

import (
	"github.com/sazardev/goca/internal/domain"

	"gorm.io/gorm"
)

type postgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) UserRepository {
	return &postgresUserRepository{
		db: db,
	}
}

func (p *postgresUserRepository) Save(user *domain.User) error {
	result := p.db.Create(user)
	return result.Error
}

func (p *postgresUserRepository) FindByID(id int) (*domain.User, error) {
	user := &domain.User{}
	result := p.db.First(user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (p *postgresUserRepository) Update(user *domain.User) error {
	result := p.db.Save(user)
	return result.Error
}

func (p *postgresUserRepository) Delete(id int) error {
	result := p.db.Delete(&domain.User{}, id)
	return result.Error
}

func (p *postgresUserRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	result := p.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (p *postgresUserRepository) FindByName(name string) (*domain.User, error) {
	user := &domain.User{}
	result := p.db.Where("name = ?", name).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (p *postgresUserRepository) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	result := p.db.Where("email = ?", email).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (p *postgresUserRepository) FindByAge(age int) (*domain.User, error) {
	user := &domain.User{}
	result := p.db.Where("age = ?", age).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
