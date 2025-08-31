package factories

import (
	"example.com/internal/model"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"time"
)

func UserFactory(overrides ...func(*model.User)) *model.User {
	userModel := &model.User{
		ID:           uuid.NewString(),
		UserName:     faker.Username(),
		Email:        faker.Email(),
		PasswordHash: faker.Password(),
		CreatedAt:    time.Now(),
	}
	for _, fn := range overrides {
		fn(userModel)
	}
	return userModel
}
