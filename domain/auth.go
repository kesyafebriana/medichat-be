package domain

import "time"

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
