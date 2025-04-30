package service

import (
	"errors"
	"example.com/internal/model"
	"example.com/internal/repository"
	"example.com/internal/utils"
	"github.com/google/uuid"
)

type AuthService interface {
	Signup(email, phone, password string) (string, error)
	Login(email, password string) (string, error)
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{userRepository}
}

func (service *authService) Signup(email, phone, password string) (string, error) {
	// Validate email and phone
	if _, err := service.userRepository.FindByEmail(email); err == nil {
		return "", errors.New("email already registered")
	}
	if _, err := service.userRepository.FindByPhone(phone); err == nil {
		return "", errors.New("phone already registered")
	}
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := &model.User{
		ID:           uuid.NewString(),
		Email:        email,
		Phone:        phone,
		PasswordHash: hashPassword,
	}
	if err := service.userRepository.Create(user); err != nil {
		return "", err
	}
	token, err := utils.GenerateToken(user.ID, user.Email)

	if err != nil {
		return "", err
	}
	return token, nil
}

func (service *authService) Login(email, password string) (string, error) {
	user, err := service.userRepository.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}
