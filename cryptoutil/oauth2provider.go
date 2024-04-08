package cryptoutil

import (
	"context"

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
	return p.config.Exchange(ctx, code)
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
