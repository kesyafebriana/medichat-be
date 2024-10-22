// Code generated by mockery v2.10.4. DO NOT EDIT.

package domainmocks

import (
	context "context"
	domain "medichat-be/domain"

	mock "github.com/stretchr/testify/mock"
)

// AccountRepository is an autogenerated mock type for the AccountRepository type
type AccountRepository struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, creds
func (_m *AccountRepository) Add(ctx context.Context, creds domain.AccountWithCredentials) (domain.Account, error) {
	ret := _m.Called(ctx, creds)

	var r0 domain.Account
	if rf, ok := ret.Get(0).(func(context.Context, domain.AccountWithCredentials) domain.Account); ok {
		r0 = rf(ctx, creds)
	} else {
		r0 = ret.Get(0).(domain.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.AccountWithCredentials) error); ok {
		r1 = rf(ctx, creds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *AccountRepository) GetByEmail(ctx context.Context, email string) (domain.Account, error) {
	ret := _m.Called(ctx, email)

	var r0 domain.Account
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Account); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(domain.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByEmailAndLock provides a mock function with given fields: ctx, email
func (_m *AccountRepository) GetByEmailAndLock(ctx context.Context, email string) (domain.Account, error) {
	ret := _m.Called(ctx, email)

	var r0 domain.Account
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Account); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(domain.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *AccountRepository) GetByID(ctx context.Context, id int64) (domain.Account, error) {
	ret := _m.Called(ctx, id)

	var r0 domain.Account
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.Account); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByIDAndLock provides a mock function with given fields: ctx, id
func (_m *AccountRepository) GetByIDAndLock(ctx context.Context, id int64) (domain.Account, error) {
	ret := _m.Called(ctx, id)

	var r0 domain.Account
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.Account); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWithCredentialsByEmail provides a mock function with given fields: ctx, email
func (_m *AccountRepository) GetWithCredentialsByEmail(ctx context.Context, email string) (domain.AccountWithCredentials, error) {
	ret := _m.Called(ctx, email)

	var r0 domain.AccountWithCredentials
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.AccountWithCredentials); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(domain.AccountWithCredentials)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWithCredentialsByID provides a mock function with given fields: ctx, id
func (_m *AccountRepository) GetWithCredentialsByID(ctx context.Context, id int64) (domain.AccountWithCredentials, error) {
	ret := _m.Called(ctx, id)

	var r0 domain.AccountWithCredentials
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.AccountWithCredentials); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.AccountWithCredentials)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsExistByEmail provides a mock function with given fields: ctx, email
func (_m *AccountRepository) IsExistByEmail(ctx context.Context, email string) (bool, error) {
	ret := _m.Called(ctx, email)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsExistByID provides a mock function with given fields: ctx, id
func (_m *AccountRepository) IsExistByID(ctx context.Context, id int64) (bool, error) {
	ret := _m.Called(ctx, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int64) bool); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileSetByID provides a mock function with given fields: ctx, id
func (_m *AccountRepository) ProfileSetByID(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, a
func (_m *AccountRepository) Update(ctx context.Context, a domain.Account) (domain.Account, error) {
	ret := _m.Called(ctx, a)

	var r0 domain.Account
	if rf, ok := ret.Get(0).(func(context.Context, domain.Account) domain.Account); ok {
		r0 = rf(ctx, a)
	} else {
		r0 = ret.Get(0).(domain.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.Account) error); ok {
		r1 = rf(ctx, a)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePasswordByID provides a mock function with given fields: ctx, id, newHashedPassword
func (_m *AccountRepository) UpdatePasswordByID(ctx context.Context, id int64, newHashedPassword string) error {
	ret := _m.Called(ctx, id, newHashedPassword)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) error); ok {
		r0 = rf(ctx, id, newHashedPassword)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyEmailByID provides a mock function with given fields: ctx, id
func (_m *AccountRepository) VerifyEmailByID(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
