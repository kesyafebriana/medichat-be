package cryptoutil_test

import (
	"medichat-be/apperror"
	"medichat-be/cryptoutil"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var (
	jwtSecretA = "194h2BV2198321.4?"
	jwtSecretB = "195285jJe29.195/1"

	randomGarbage = "10ihi409jcb2h09j4kng049909m"

	lifespanNormal = 15 * time.Minute
	lifespanNeg    = -1 * time.Minute

	issuerA = "issuer-a"
	issuerB = "issuer-b"

	myUserID = int64(0)

	myClaims = cryptoutil.JWTClaims{
		UserID: myUserID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: issuerA,
		},
	}
)

func Test_jwtProviderHS256_CreateToken(t *testing.T) {
	tests := []struct {
		name string

		userID    int64
		issuer    string
		secretKey string
		lifespan  time.Duration

		wantErr int
	}{
		{
			name:   "should successfully return a signed token",
			userID: myUserID,

			issuer:    issuerA,
			secretKey: jwtSecretA,
			lifespan:  lifespanNormal,

			wantErr: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := cryptoutil.NewJWTProviderHS256(tt.issuer, tt.secretKey, tt.lifespan)

			_, err := p.CreateToken(tt.userID)

			if tt.wantErr != 0 {
				apperror.AssertErrorIsCode(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func Test_jwtProviderHS256_VerifyToken(t *testing.T) {
	tests := []struct {
		name string

		userID    int64
		issuer    string
		secretKey string
		lifespan  time.Duration

		verifierIssuer    string
		verifierSecretKey string
		verifierLifespan  time.Duration

		want    cryptoutil.JWTClaims
		wantErr int
	}{
		{
			name:   "should successfully verify a token",
			userID: myUserID,

			issuer:    issuerA,
			secretKey: jwtSecretA,
			lifespan:  lifespanNormal,

			verifierIssuer:    issuerA,
			verifierSecretKey: jwtSecretA,
			verifierLifespan:  lifespanNormal,

			want: myClaims,
		},
		{
			name:   "should return invalid token when verify with different secret",
			userID: myUserID,

			issuer:    issuerA,
			secretKey: jwtSecretA,
			lifespan:  lifespanNormal,

			verifierIssuer:    issuerA,
			verifierSecretKey: jwtSecretB,
			verifierLifespan:  lifespanNormal,

			wantErr: apperror.CodeUnauthorized,
		},
		{
			name:   "should return invalid token when verify with different issuer",
			userID: myUserID,

			issuer:    issuerA,
			secretKey: jwtSecretA,
			lifespan:  lifespanNormal,

			verifierIssuer:    issuerB,
			verifierSecretKey: jwtSecretA,
			verifierLifespan:  lifespanNormal,

			wantErr: apperror.CodeUnauthorized,
		},
		{
			name:   "should return invalid token when verify expired token",
			userID: myUserID,

			issuer:    issuerA,
			secretKey: jwtSecretA,
			lifespan:  lifespanNeg,

			verifierIssuer:    issuerA,
			verifierSecretKey: jwtSecretA,
			verifierLifespan:  lifespanNormal,

			wantErr: apperror.CodeUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := cryptoutil.NewJWTProviderHS256(tt.issuer, tt.secretKey, tt.lifespan)
			v := cryptoutil.NewJWTProviderHS256(tt.verifierIssuer, tt.verifierSecretKey, tt.verifierLifespan)

			token, _ := p.CreateToken(tt.userID)

			claims, err := v.VerifyToken(token)
			claims.ExpiresAt = nil
			claims.IssuedAt = nil

			assert.Equal(t, tt.want, claims)
			if tt.wantErr != 0 {
				apperror.AssertErrorIsCode(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}

	t.Run("should return invalid", func(t *testing.T) {
		v := cryptoutil.NewJWTProviderHS256(issuerA, jwtSecretA, lifespanNormal)

		_, err := v.VerifyToken(randomGarbage)

		apperror.AssertErrorIsCode(t, err, apperror.CodeUnauthorized)
	})
}
