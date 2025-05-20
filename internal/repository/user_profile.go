package repository

import (
	"example.com/internal/model"

	"gorm.io/gorm"
)

type UserProfileRepository interface {
	Create(userProfile *model.UserProfile) error
}

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (repository *userProfileRepository) Create(userProfile *model.UserProfile) error {
	return repository.db.Create(userProfile).Error
}
