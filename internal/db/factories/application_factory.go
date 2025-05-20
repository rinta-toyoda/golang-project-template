package factories

import (
	"example.com/internal/model"
	"github.com/google/uuid"
	"time"
)

func ApplicationFactory(userID, jobID string, overrides ...func(*model.Application)) *model.Application {
	applicationModel := &model.Application{
		ID:        uuid.NewString(),
		UserID:    userID,
		JobID:     jobID,
		AppliedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	for _, fn := range overrides {
		fn(applicationModel)
	}
	return applicationModel
}
