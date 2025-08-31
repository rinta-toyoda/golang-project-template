package security

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}

type bcryptHasher struct {
	cost int
}

func NewBcryptHasher() PasswordHasher {
	return &bcryptHasher{
		cost: bcrypt.DefaultCost,
	}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	return string(bytes), err
}

func (h *bcryptHasher) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
