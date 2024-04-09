package cryptoutil

import (
	"context"
	"medichat-be/apperror"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuth2Provider interface {
	GetAuthURL(state string) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
}

type oAuth2ProviderImpl struct {
	config oauth2.Config
}

func (p *oAuth2ProviderImpl) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state)
}

func (p *oAuth2ProviderImpl) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	tok, err := p.config.Exchange(ctx, code)
	if reterr, ok := err.(*oauth2.RetrieveError); ok && reterr.ErrorCode == "invalid_grant" {
		return nil, apperror.NewAppError(
			apperror.CodeBadRequest,
			"invalid oauth2 grant",
			err,
		)
	}
	if err != nil {
		return tok, apperror.Wrap(err)
	}
	return tok, nil
}

type GoogleAuthProviderOpts struct {
	RedirectURL  string
	ClientID     string
	ClientSecret string
}

func NewGoogleAuthProvider(opts GoogleAuthProviderOpts) *oAuth2ProviderImpl {
	return &oAuth2ProviderImpl{
		config: oauth2.Config{
			RedirectURL:  opts.RedirectURL,
			ClientID:     opts.ClientID,
			ClientSecret: opts.ClientSecret,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}
