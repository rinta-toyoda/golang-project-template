package factories

import (
	"math/rand"
	"time"

	"example.com/internal/model"
	"github.com/google/uuid"
)

func UserSkillFactory(userID, skillID, skillableType, skillableID string, overrides ...func(*model.UserSkill)) *model.UserSkill {
	userSkill := &model.UserSkill{
		ID:            uuid.NewString(),
		UserID:        userID,
		SkillID:       skillID,
		SkillableType: skillableType, // "Education" or "Experience"
		SkillableID:   skillableID,
		Years:         float64(rand.Intn(50)+10) / 10.0, // 1桁小数のランダムな年数
		CreatedAt:     time.Now(),
	}

	for _, fn := range overrides {
		fn(userSkill)
	}

	return userSkill
}
