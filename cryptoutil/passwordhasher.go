package cryptoutil

import (
	"medichat-be/apperror"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword string, password string) error
}

type passwordHasherBcrypt struct {
	hashCost int
}

func NewPasswordHasherBcrypt(hashCost int) *passwordHasherBcrypt {
	return &passwordHasherBcrypt{
		hashCost: hashCost,
	}
}

func (ph *passwordHasherBcrypt) HashPassword(password string) (string, error) {
	hpBytes, err := bcrypt.GenerateFromPassword([]byte(password), ph.hashCost)
	if err != nil {
		return "", apperror.Wrap(err)
	}

	return string(hpBytes), nil
}

func (ph *passwordHasherBcrypt) CheckPassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return apperror.NewWrongPassword(err)
	}
	if err != nil {
		return apperror.Wrap(err)
	}

	return nil
}
