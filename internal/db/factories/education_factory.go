package factories

import (
	"example.com/internal/model"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"time"
)

func EducationFactory(profileID string, overrides ...func(*model.Education)) *model.Education {
	end := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
	educationModel := &model.Education{
		ID:          uuid.NewString(),
		ProfileID:   profileID,
		Institution: faker.Word() + " University",
		Description: faker.Sentence(),
		StartDate:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     &end,
		CreatedAt:   time.Now(),
	}
	for _, fn := range overrides {
		fn(educationModel)
	}
	return educationModel
}
