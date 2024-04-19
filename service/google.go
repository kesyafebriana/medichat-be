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
	oauth2Service  domain.OAuth2Service
	accountService domain.AccountService
}

type GoogleServiceOpts struct {
	DataRepository domain.DataRepository
	OAuth2Service  domain.OAuth2Service
	AccountService domain.AccountService
}

func NewGoogleService(opts GoogleServiceOpts) *googleService {
	return &googleService{
		dataRepository: opts.DataRepository,
		oauth2Service:  opts.OAuth2Service,
		accountService: opts.AccountService,
	}
}

func (s *googleService) OAuth2Callback(
	ctx context.Context,
	state string,
	opts domain.OAuth2CallbackOpts,
) (domain.AuthTokens, error) {
	rtRepo := s.dataRepository.RefreshTokenRepository()

	gTokens, err := s.oauth2Service.Callback(ctx, state, opts)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}

	account, err := s.EnsureRegisteredByToken(ctx, gTokens.AccessToken)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}

	tokens, err := s.accountService.CreateTokensForAccount(account.ID, account.Role)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}

	rToken := domain.RefreshToken{
		Account:   account,
		Token:     tokens.RefreshToken,
		ClientIP:  opts.ClientIP,
		ExpiredAt: tokens.RefreshExpireAt,
	}

	_, err = rtRepo.Add(ctx, rToken)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}

	return tokens, nil
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
		accountRepo := dr.AccountRepository()

		newAccount := domain.Account{
			Email:         profile.Email,
			EmailVerified: profile.VerifiedEmail,
			Role:          domain.AccountRoleUser,
			AccountType:   domain.AccountTypeGoogle,
		}

		account, err := accountRepo.GetByEmail(ctx, newAccount.Email)
		if err != nil && !apperror.IsErrorCode(err, apperror.CodeNotFound) {
			return domain.Account{}, apperror.Wrap(err)
		}
		if err == nil {
			if account.AccountType != domain.AccountTypeGoogle {
				return domain.Account{}, apperror.NewAppError(
					apperror.CodeUnauthorized,
					"not a google account",
					nil,
				)
			}

			return account, nil
		}

		account, err = accountRepo.Add(ctx, domain.AccountWithCredentials{
			Account:        newAccount,
			HashedPassword: nil,
		})
		if err != nil {
			return domain.Account{}, apperror.Wrap(err)
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
