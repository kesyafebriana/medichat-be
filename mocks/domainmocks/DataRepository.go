// Code generated by mockery v2.10.4. DO NOT EDIT.

package domainmocks

import (
	context "context"
	domain "medichat-be/domain"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// DataRepository is an autogenerated mock type for the DataRepository type
type DataRepository struct {
	mock.Mock
}

// AccountRepository provides a mock function with given fields:
func (_m *DataRepository) AccountRepository() domain.AccountRepository {
	ret := _m.Called()

	var r0 domain.AccountRepository
	if rf, ok := ret.Get(0).(func() domain.AccountRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.AccountRepository)
		}
	}

	return r0
}

// AdminRepository provides a mock function with given fields:
func (_m *DataRepository) AdminRepository() domain.AdminRepository {
	ret := _m.Called()

	var r0 domain.AdminRepository
	if rf, ok := ret.Get(0).(func() domain.AdminRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.AdminRepository)
		}
	}

	return r0
}

// Atomic provides a mock function with given fields: ctx, fn
func (_m *DataRepository) Atomic(ctx context.Context, fn domain.AtomicFuncAny) (interface{}, error) {
	ret := _m.Called(ctx, fn)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, domain.AtomicFuncAny) interface{}); ok {
		r0 = rf(ctx, fn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.AtomicFuncAny) error); ok {
		r1 = rf(ctx, fn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DoctorRepository provides a mock function with given fields:
func (_m *DataRepository) DoctorRepository() domain.DoctorRepository {
	ret := _m.Called()

	var r0 domain.DoctorRepository
	if rf, ok := ret.Get(0).(func() domain.DoctorRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.DoctorRepository)
		}
	}

	return r0
}

// PharmacyManagerRepository provides a mock function with given fields:
func (_m *DataRepository) PharmacyManagerRepository() domain.PharmacyManagerRepository {
	ret := _m.Called()

	var r0 domain.PharmacyManagerRepository
	if rf, ok := ret.Get(0).(func() domain.PharmacyManagerRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.PharmacyManagerRepository)
		}
	}

	return r0
}

// RefreshTokenRepository provides a mock function with given fields:
func (_m *DataRepository) RefreshTokenRepository() domain.RefreshTokenRepository {
	ret := _m.Called()

	var r0 domain.RefreshTokenRepository
	if rf, ok := ret.Get(0).(func() domain.RefreshTokenRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.RefreshTokenRepository)
		}
	}

	return r0
}

// ResetPasswordTokenRepository provides a mock function with given fields:
func (_m *DataRepository) ResetPasswordTokenRepository() domain.ResetPasswordTokenRepository {
	ret := _m.Called()

	var r0 domain.ResetPasswordTokenRepository
	if rf, ok := ret.Get(0).(func() domain.ResetPasswordTokenRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.ResetPasswordTokenRepository)
		}
	}

	return r0
}

// Sleep provides a mock function with given fields: ctx, duration
func (_m *DataRepository) Sleep(ctx context.Context, duration time.Duration) error {
	ret := _m.Called(ctx, duration)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Duration) error); ok {
		r0 = rf(ctx, duration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository provides a mock function with given fields:
func (_m *DataRepository) UserRepository() domain.UserRepository {
	ret := _m.Called()

	var r0 domain.UserRepository
	if rf, ok := ret.Get(0).(func() domain.UserRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.UserRepository)
		}
	}

	return r0
}

// VerifyEmailTokenRepository provides a mock function with given fields:
func (_m *DataRepository) VerifyEmailTokenRepository() domain.VerifyEmailTokenRepository {
	ret := _m.Called()

	var r0 domain.VerifyEmailTokenRepository
	if rf, ok := ret.Get(0).(func() domain.VerifyEmailTokenRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.VerifyEmailTokenRepository)
		}
	}

	return r0
}
