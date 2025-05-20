package factories

import (
	"example.com/internal/model"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"time"
)

func JobFactory(organizationProfileID string, overrides ...func(*model.Job)) *model.Job {
	jobModel := &model.Job{
		ID:                    uuid.NewString(),
		OrganizationProfileID: organizationProfileID,
		Name:                  faker.Word() + " Engineer",
		Description:           faker.Sentence(),
		StartDate:             time.Now().AddDate(0, -6, 0),
		EndDate:               nil,
		IsActive:              true,
		CreatedAt:             time.Now(),
	}
	for _, fn := range overrides {
		fn(jobModel)
	}
	return jobModel
}
