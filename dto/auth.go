package dto

import (
	"medichat-be/domain"
	"time"
)

type AuthTokensResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func NewAuthTokensResponse(t domain.AuthTokens) AuthTokensResponse {
	return AuthTokensResponse{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiresAt:    t.ExpiresAt,
	}
}
