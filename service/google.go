package service

import (
	"context"
	"encoding/json"
	"medichat-be/apperror"
	"medichat-be/domain"
	"medichat-be/dto"
	"net/http"
)

type googleService struct {
	dataRepository domain.DataRepository
}

type GoogleServiceOpts struct {
}

func NewGoogleService(opts GoogleServiceOpts) *googleService {
	return &googleService{}
}

func (s *googleService) GetProfileByAccessToken(
	ctx context.Context,
	accessToken string,
) (domain.GoogleUserProfile, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return domain.GoogleUserProfile{}, apperror.Wrap(err)
	}

	if resp.StatusCode >= 400 {
		return domain.GoogleUserProfile{}, apperror.NewInternal(nil)
	}

	var userProfileResp dto.GoogleUserProfileResponse

	err = json.NewDecoder(resp.Body).Decode(&userProfileResp)
	if err != nil {
		return domain.GoogleUserProfile{}, apperror.Wrap(err)
	}

	userProfile := userProfileResp.ToProfile()

	return userProfile, nil
}

func (s *googleService) EnsureRegisteredClosure(
	ctx context.Context,
	profile domain.GoogleUserProfile,
) domain.AtomicFunc[domain.Account] {
	return func(dr domain.DataRepository) (domain.Account, error) {
		account := domain.Account{
			Email:         profile.Email,
			EmailVerified: profile.VerifiedEmail,
			Role:          domain.AccountRoleUser,
			AccountType:   domain.AccountTypeGoogle,
		}

		return account, nil
	}
}

func (s *googleService) EnsureRegistered(
	ctx context.Context,
	profile domain.GoogleUserProfile,
) (domain.Account, error) {
	return domain.RunAtomic(
		s.dataRepository,
		ctx,
		s.EnsureRegisteredClosure(ctx, profile),
	)
}

func (s *googleService) EnsureRegisteredByToken(
	ctx context.Context,
	accessToken string,
) (domain.Account, error) {
	profile, err := s.GetProfileByAccessToken(ctx, accessToken)
	if err != nil {
		return domain.Account{}, apperror.Wrap(err)
	}

	return s.EnsureRegistered(ctx, profile)
}
