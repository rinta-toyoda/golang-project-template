package factories

import (
	"example.com/internal/model"
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func OrganizationProfileFactory(organizationID string, overrides ...func(*model.OrganizationProfile)) *model.OrganizationProfile {
	profileModel := &model.OrganizationProfile{
		ID:             uuid.NewString(),
		OrganizationID: organizationID,
		Name:           faker.Word(),
		Address:        fmt.Sprintf("%d %s Ave", rand.Intn(900)+100, faker.LastName()),
		Description:    faker.Sentence(),
		CreatedAt:      time.Now(),
	}
	for _, fn := range overrides {
		fn(profileModel)
	}
	return profileModel
}
