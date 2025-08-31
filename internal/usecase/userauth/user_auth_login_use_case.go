package userauth

import (
	"example.com/internal/service"
	"example.com/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserAuthLoginUseCase interface {
	Call(context *gin.Context, identifier, password string) error
}
type userAuthLoginUseCase struct {
	userAuthService service.UserAuthService
}

func NewUserAuthLoginUseCase(userAuthService service.UserAuthService) UserAuthLoginUseCase {
	return &userAuthLoginUseCase{userAuthService: userAuthService}
}

func (u *userAuthLoginUseCase) Call(context *gin.Context, identifier, password string) error {
	userId, err := u.userAuthService.Login(identifier, password)
	if err != nil {
		return err
	}
	err = utils.SaveUserSession(context, userId)

	return err
}
