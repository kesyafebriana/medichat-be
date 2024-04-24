package service

import (
	"context"
	"medichat-be/apperror"
	"medichat-be/cryptoutil"
	"medichat-be/domain"
)

type oAuth2Service struct {
	oAuth2Provider cryptoutil.OAuth2Provider
}

type OAuth2ServiceOpts struct {
	OAuth2Provider cryptoutil.OAuth2Provider
}

func NewOAuth2Service(opts OAuth2ServiceOpts) *oAuth2Service {
	return &oAuth2Service{
		oAuth2Provider: opts.OAuth2Provider,
	}
}

func (s *oAuth2Service) GetAuthURL(ctx context.Context, state string) (string, error) {
	url := s.oAuth2Provider.GetAuthURL(state)
	return url, nil
}

func (s *oAuth2Service) Callback(ctx context.Context, state string, opts domain.OAuth2CallbackOpts) (domain.AuthTokens, error) {
	if state != opts.State {
		return domain.AuthTokens{}, apperror.NewAppError(
			apperror.CodeUnauthorized,
			"invalid state",
			nil,
		)
	}

	token, err := s.oAuth2Provider.Exchange(ctx, opts.Code)
	if err != nil {
		return domain.AuthTokens{}, apperror.Wrap(err)
	}

	return domain.AuthTokens{
		AccessToken:     token.AccessToken,
		RefreshToken:    token.RefreshToken,
		AccessExpiresAt: token.Expiry,
	}, nil
}
