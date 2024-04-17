package cryptoutil

import (
	"medichat-be/apperror"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"uid"`
}

type JWTProvider interface {
	CreateToken(userID int64) (string, error)
	VerifyToken(token string) (JWTClaims, error)
}

type jwtProviderHS256 struct {
	issuer    string
	secretKey string
	lifespan  time.Duration
}

func NewJWTProviderHS256(issuer string, secretKey string, lifespan time.Duration) *jwtProviderHS256 {
	return &jwtProviderHS256{
		issuer:    issuer,
		secretKey: secretKey,
		lifespan:  lifespan,
	}
}

func (p *jwtProviderHS256) CreateToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    p.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.lifespan)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	})

	signed, err := token.SignedString([]byte(p.secretKey))

	if err != nil {
		return "", apperror.Wrap(err)
	}

	return signed, nil
}

func (p *jwtProviderHS256) VerifyToken(tokenstr string) (JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenstr,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(p.secretKey), nil
		},
		jwt.WithIssuer(p.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return JWTClaims{}, apperror.NewInvalidToken(err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return JWTClaims{}, apperror.NewTypeAssertionFailed(claims, token)
	}

	return *claims, nil
}

type jwtProviderAny struct {
	providers []JWTProvider
}

func NewJWTProviderAny(providers []JWTProvider) *jwtProviderAny {
	return &jwtProviderAny{
		providers: providers,
	}
}

func (p *jwtProviderAny) CreateToken(userID int64) (string, error) {
	return "", apperror.NewInternalFmt("uninplemented")
}

func (p *jwtProviderAny) VerifyToken(tokenstr string) (JWTClaims, error) {
	var err error

	for _, prov := range p.providers {
		claims, err := prov.VerifyToken(tokenstr)
		if err == nil {
			return claims, nil
		}
	}

	return JWTClaims{}, err
}
