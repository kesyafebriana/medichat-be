package cryptoutil

import (
	"crypto/rand"
	"encoding/base64"
	"medichat-be/apperror"
)

type RandomTokenProvider interface {
	GenerateToken() (string, error)
}

type randomTokenProviderImpl struct {
	length int
}

func NewRandomTokenProvider(length int) *randomTokenProviderImpl {
	return &randomTokenProviderImpl{
		length: length,
	}
}

func (p *randomTokenProviderImpl) GenerateToken() (string, error) {
	b := make([]byte, p.length)
	_, err := rand.Read(b)
	if err != nil {
		return "", apperror.Wrap(err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
