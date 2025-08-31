package userauth

import (
	"example.com/internal/service"
	"example.com/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserAuthSignupUseCase interface {
	Call(context *gin.Context, userName, email, password string) error
}
type userAuthSignupUseCase struct {
	userAuthService service.UserAuthService
}

func NewUserAuthSignupUseCase(userAuthService service.UserAuthService) UserAuthSignupUseCase {
	return &userAuthSignupUseCase{userAuthService: userAuthService}
}

func (u *userAuthSignupUseCase) Call(context *gin.Context, userName, email, password string) error {
	userId, err := u.userAuthService.Signup(userName, email, password)
	if err != nil {
		return err
	}
	err = utils.SaveUserSession(context, userId)

	return err
}
