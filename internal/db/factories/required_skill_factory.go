package factories

import (
	"example.com/internal/model"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func RequiredSkillFactory(jobID string, overrides ...func(*model.RequiredSkill)) *model.RequiredSkill {
	skillModel := &model.RequiredSkill{
		ID:        uuid.NewString(),
		JobID:     jobID,
		Name:      faker.Word(),
		Years:     float64(rand.Intn(5) + 1),
		CreatedAt: time.Now(),
	}
	for _, fn := range overrides {
		fn(skillModel)
	}
	return skillModel
}
