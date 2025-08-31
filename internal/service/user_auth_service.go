package service

import (
	"errors"
	"example.com/internal/model"
	"example.com/internal/repository"
	"example.com/internal/utils"
	"github.com/google/uuid"
)

type UserAuthService interface {
	Signup(userName, email, password string) (string, error)
	Login(identifier, password string) (string, error)
}

type userAuthService struct {
	userRepository repository.UserRepository
}

func NewUserAuthService(userRepository repository.UserRepository) UserAuthService {
	return &userAuthService{userRepository}
}

func (service *userAuthService) Signup(userName, email, password string) (string, error) {
	// Validate userName and email
	if _, err := service.userRepository.FindByUserName(userName); err == nil {
		return "", errors.New("phone already registered")
	}
	if _, err := service.userRepository.FindByEmail(email); err == nil {
		return "", errors.New("email already registered")
	}
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := &model.User{
		ID:           uuid.NewString(),
		UserName:     userName,
		Email:        email,
		PasswordHash: hashPassword,
	}
	if err := service.userRepository.Create(user); err != nil {
		return "", err
	}
	return user.ID, nil
}

func (service *userAuthService) Login(identifier, password string) (string, error) {
	user, err := service.userRepository.FindByUserNameOrEmail(identifier)
	if err != nil {
		return "", err
	}
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}
	return user.ID, nil
}
