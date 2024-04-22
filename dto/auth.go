package dto

import (
	"medichat-be/domain"
	"time"
)

type AuthTokensResponse struct {
	AccessToken      string    `json:"access_token,omitempty"`
	RefreshToken     string    `json:"refresh_token,omitempty"`
	AccessExpiresAt  time.Time `json:"access_expires_at,omitempty"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at,omitempty"`
}

func NewAuthTokensResponse(t domain.AuthTokens) AuthTokensResponse {
	return AuthTokensResponse{
		AccessToken:      t.AccessToken,
		RefreshToken:     t.RefreshToken,
		AccessExpiresAt:  t.AccessExpiresAt,
		RefreshExpiresAt: t.RefreshExpireAt,
	}
}
