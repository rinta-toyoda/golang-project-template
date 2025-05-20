package factories

import (
	"example.com/internal/model"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"time"
)

func OrganizationFactory(overrides ...func(*model.Organization)) *model.Organization {
	organizationModel := &model.Organization{
		ID:           uuid.NewString(),
		Email:        faker.Email(),
		Phone:        faker.Phonenumber(),
		PasswordHash: faker.Password(),
		CreatedAt:    time.Now(),
	}
	for _, fn := range overrides {
		fn(organizationModel)
	}
	return organizationModel
}
