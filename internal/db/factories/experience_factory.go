package factories

import (
	"example.com/internal/model"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"time"
)

func ExperienceFactory(profileID string, overrides ...func(*model.Experience)) *model.Experience {
	experienceModel := &model.Experience{
		ID:          uuid.NewString(),
		ProfileID:   profileID,
		Company:     faker.Word() + " Inc",
		Description: faker.Sentence(),
		StartDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     nil,
		CreatedAt:   time.Now(),
	}
	for _, fn := range overrides {
		fn(experienceModel)
	}
	return experienceModel
}
