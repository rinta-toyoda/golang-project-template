package utils

import (
	"errors"
	"example.com/internal/model"
)

type FakeUserRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindByPhone(phone string) (*model.User, error)
	Create(user *model.User) error
}

type fakeUserRepository struct {
	Email map[string]*model.User
	Phone map[string]*model.User
}

func NewFakeUserRepository() FakeUserRepository {
	return &fakeUserRepository{
		Email: make(map[string]*model.User),
		Phone: make(map[string]*model.User),
	}
}

func (f *fakeUserRepository) FindByEmail(email string) (*model.User, error) {
	user, ok := f.Email[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return user, nil
}

func (f *fakeUserRepository) FindByPhone(phone string) (*model.User, error) {
	user, ok := f.Phone[phone]
	if !ok {
		return nil, errors.New("not found")
	}
	return user, nil
}

func (f *fakeUserRepository) Create(user *model.User) error {
	f.Email[user.Email] = user
	f.Phone[user.Phone] = user
	return nil
}
