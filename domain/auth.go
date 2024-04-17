package domain

import "time"

type AuthTokens struct {
	AccessToken     string
	RefreshToken    string
	AccessExpiresAt time.Time
	RefreshExpireAt time.Time
}
