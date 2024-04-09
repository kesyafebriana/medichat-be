package dto

import "medichat-be/domain"

type GoogleUserProfileResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (r *GoogleUserProfileResponse) ToProfile() domain.GoogleUserProfile {
	return domain.GoogleUserProfile{
		ID:            r.ID,
		Email:         r.Email,
		VerifiedEmail: r.VerifiedEmail,
		Name:          r.Name,
		GivenName:     r.GivenName,
		FamilyName:    r.FamilyName,
		Picture:       r.Picture,
		Locale:        r.Locale,
	}
}
