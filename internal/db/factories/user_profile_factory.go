package factories

import (
	"example.com/internal/model"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"time"
)

func UserProfileFactory(userID string, overrides ...func(*model.UserProfile)) *model.UserProfile {
	userProfileModel := &model.UserProfile{
		ID:          uuid.NewString(),
		UserID:      userID,
		FirstName:   faker.FirstName(),
		MiddleName:  "",
		LastName:    faker.LastName(),
		Address:     faker.Sentence(),
		Description: faker.Sentence(),
		CreatedAt:   time.Now(),
	}
	for _, fn := range overrides {
		fn(userProfileModel)
	}
	return userProfileModel
}
