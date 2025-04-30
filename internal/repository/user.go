package repository

import (
	"example.com/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByPhone(phone string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (repository *userRepository) Create(user *model.User) error {
	return repository.db.Create(user).Error
}

func (repository *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := repository.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repository *userRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := repository.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
