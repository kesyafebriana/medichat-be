package cryptoutil_test

import (
	"medichat-be/apperror"
	"medichat-be/cryptoutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	hashCost        = bcrypt.MinCost
	myPassword      = "SuperSecretPassword12345"
	myWrongPassword = "SupersecretPassword12345"
	myLongPassword  = "SuperSecretPassword12345678901234567890123456789012345678901234567890123456789012345678901234567890"
	invalidHash     = "18298099307012"
)

func Test_passwordHasherBcrypt_HashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  int
	}{
		{
			name:     "should return no error",
			password: myPassword,
			wantErr:  0,
		},
		{
			name:     "should return error when password is too long",
			password: myLongPassword,
			wantErr:  apperror.CodeInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ph := cryptoutil.NewPasswordHasherBcrypt(hashCost)

			// when
			_, err := ph.HashPassword(tt.password)

			// then
			if tt.wantErr != 0 {
				apperror.AssertErrorIsCode(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func Test_passwordHasherBcrypt_CheckPassword(t *testing.T) {
	tests := []struct {
		name             string
		originalPassword string
		password         string
		wantErr          int
	}{
		{
			name:             "should successfully check password",
			originalPassword: myPassword,
			password:         myPassword,
			wantErr:          0,
		},
		{
			name:             "should return error when the password is wrong",
			originalPassword: myPassword,
			password:         myWrongPassword,
			wantErr:          apperror.CodeUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ph := cryptoutil.NewPasswordHasherBcrypt(hashCost)
			hashed, _ := ph.HashPassword(tt.originalPassword)

			// when
			err := ph.CheckPassword(hashed, tt.password)

			// then
			if tt.wantErr != 0 {
				apperror.AssertErrorIsCode(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}

	t.Run("should return internal error when given an invalid hash", func(t *testing.T) {
		// given
		ph := cryptoutil.NewPasswordHasherBcrypt(hashCost)

		// when
		err := ph.CheckPassword(invalidHash, myPassword)

		// then
		apperror.AssertErrorIsCode(t, err, apperror.CodeInternal)
	})
}
