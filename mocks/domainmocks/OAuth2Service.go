// Code generated by mockery v2.10.4. DO NOT EDIT.

package domainmocks

import (
	context "context"
	domain "medichat-be/domain"

	mock "github.com/stretchr/testify/mock"
)

// OAuth2Service is an autogenerated mock type for the OAuth2Service type
type OAuth2Service struct {
	mock.Mock
}

// Callback provides a mock function with given fields: ctx, state, opts
func (_m *OAuth2Service) Callback(ctx context.Context, state string, opts domain.OAuth2CallbackOpts) (domain.AuthTokens, error) {
	ret := _m.Called(ctx, state, opts)

	var r0 domain.AuthTokens
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.OAuth2CallbackOpts) domain.AuthTokens); ok {
		r0 = rf(ctx, state, opts)
	} else {
		r0 = ret.Get(0).(domain.AuthTokens)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, domain.OAuth2CallbackOpts) error); ok {
		r1 = rf(ctx, state, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAuthURL provides a mock function with given fields: ctx, state
func (_m *OAuth2Service) GetAuthURL(ctx context.Context, state string) (string, error) {
	ret := _m.Called(ctx, state)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, state)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, state)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
