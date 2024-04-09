package domain

import "context"

type GoogleUserProfile struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Picture       string
	Locale        string
}

type GoogleService interface {
	EnsureRegistered(ctx context.Context, profile GoogleUserProfile) (Account, error)
	EnsureRegisteredByToken(ctx context.Context, accessToken string) (Account, error)
	GetProfileByAccessToken(ctx context.Context, accessToken string) (GoogleUserProfile, error)
}
