package repository

import (
	"example.com/internal/model"

	"gorm.io/gorm"
)

type OrganizationRepository interface {
	Create(organization *model.Organization) error
	FindByEmail(email string) (*model.Organization, error)
	FindByPhone(phone string) (*model.Organization, error)
}

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (repository *organizationRepository) Create(organization *model.Organization) error {
	return repository.db.Create(organization).Error
}

func (repository *organizationRepository) FindByEmail(email string) (*model.Organization, error) {
	var organization model.Organization
	if err := repository.db.Where("email = ?", email).First(&organization).Error; err != nil {
		return nil, err
	}
	return &organization, nil
}

func (repository *organizationRepository) FindByPhone(phone string) (*model.Organization, error) {
	var organization model.Organization
	if err := repository.db.Where("phone = ?", phone).First(&organization).Error; err != nil {
		return nil, err
	}
	return &organization, nil
}
