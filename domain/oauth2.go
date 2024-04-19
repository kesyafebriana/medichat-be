package domain

import "context"

type OAuth2CallbackOpts struct {
	Code     string
	State    string
	ClientIP string
}

type OAuth2Service interface {
	GetAuthURL(ctx context.Context, state string) (string, error)
	Callback(ctx context.Context, state string, opts OAuth2CallbackOpts) (AuthTokens, error)
}
