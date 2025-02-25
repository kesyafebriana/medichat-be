// Code generated by mockery v2.10.4. DO NOT EDIT.

package domainmocks

import (
	context "context"
	domain "medichat-be/domain"

	mock "github.com/stretchr/testify/mock"
)

// AccountService is an autogenerated mock type for the AccountService type
type AccountService struct {
	mock.Mock
}

// CheckResetPasswordToken provides a mock function with given fields: ctx, email, tokenStr
func (_m *AccountService) CheckResetPasswordToken(ctx context.Context, email string, tokenStr string) error {
	ret := _m.Called(ctx, email, tokenStr)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, email, tokenStr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckVerifyEmailToken provides a mock function with given fields: ctx, email, tokenStr
func (_m *AccountService) CheckVerifyEmailToken(ctx context.Context, email string, tokenStr string) error {
	ret := _m.Called(ctx, email, tokenStr)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, email, tokenStr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateTokensForAccount provides a mock function with given fields: accountID, role
func (_m *AccountService) CreateTokensForAccount(accountID int64, role string) (domain.AuthTokens, error) {
	ret := _m.Called(accountID, role)

	var r0 domain.AuthTokens
	if rf, ok := ret.Get(0).(func(int64, string) domain.AuthTokens); ok {
		r0 = rf(accountID, role)
	} else {
		r0 = ret.Get(0).(domain.AuthTokens)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, string) error); ok {
		r1 = rf(accountID, role)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProfile provides a mock function with given fields: ctx
func (_m *AccountService) GetProfile(ctx context.Context) (interface{}, error) {
	ret := _m.Called(ctx)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context) interface{}); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetResetPasswordToken provides a mock function with given fields: ctx, email
func (_m *AccountService) GetResetPasswordToken(ctx context.Context, email string) (string, error) {
	ret := _m.Called(ctx, email)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVerifyEmailToken provides a mock function with given fields: ctx, email
func (_m *AccountService) GetVerifyEmailToken(ctx context.Context, email string) (string, error) {
	ret := _m.Called(ctx, email)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: ctx, creds
func (_m *AccountService) Login(ctx context.Context, creds domain.AccountLoginCredentials) (domain.AuthTokens, error) {
	ret := _m.Called(ctx, creds)

	var r0 domain.AuthTokens
	if rf, ok := ret.Get(0).(func(context.Context, domain.AccountLoginCredentials) domain.AuthTokens); ok {
		r0 = rf(ctx, creds)
	} else {
		r0 = ret.Get(0).(domain.AuthTokens)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.AccountLoginCredentials) error); ok {
		r1 = rf(ctx, creds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RefreshTokens provides a mock function with given fields: ctx, creds
func (_m *AccountService) RefreshTokens(ctx context.Context, creds domain.AccountRefreshTokensCredentials) (domain.AuthTokens, error) {
	ret := _m.Called(ctx, creds)

	var r0 domain.AuthTokens
	if rf, ok := ret.Get(0).(func(context.Context, domain.AccountRefreshTokensCredentials) domain.AuthTokens); ok {
		r0 = rf(ctx, creds)
	} else {
		r0 = ret.Get(0).(domain.AuthTokens)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.AccountRefreshTokensCredentials) error); ok {
		r1 = rf(ctx, creds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: ctx, creds
func (_m *AccountService) Register(ctx context.Context, creds domain.AccountRegisterCredentials) (domain.Account, error) {
	ret := _m.Called(ctx, creds)

	var r0 domain.Account
	if rf, ok := ret.Get(0).(func(context.Context, domain.AccountRegisterCredentials) domain.Account); ok {
		r0 = rf(ctx, creds)
	} else {
		r0 = ret.Get(0).(domain.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.AccountRegisterCredentials) error); ok {
		r1 = rf(ctx, creds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResetPassword provides a mock function with given fields: ctx, creds
func (_m *AccountService) ResetPassword(ctx context.Context, creds domain.AccountResetPasswordCredentials) error {
	ret := _m.Called(ctx, creds)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.AccountResetPasswordCredentials) error); ok {
		r0 = rf(ctx, creds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyEmail provides a mock function with given fields: ctx, creds
func (_m *AccountService) VerifyEmail(ctx context.Context, creds domain.AccountVerifyEmailCredentials) error {
	ret := _m.Called(ctx, creds)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.AccountVerifyEmailCredentials) error); ok {
		r0 = rf(ctx, creds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
