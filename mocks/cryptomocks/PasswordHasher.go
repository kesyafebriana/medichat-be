// Code generated by mockery v2.10.4. DO NOT EDIT.

package cryptomocks

import mock "github.com/stretchr/testify/mock"

// PasswordHasher is an autogenerated mock type for the PasswordHasher type
type PasswordHasher struct {
	mock.Mock
}

// CheckPassword provides a mock function with given fields: hashedPassword, password
func (_m *PasswordHasher) CheckPassword(hashedPassword string, password string) error {
	ret := _m.Called(hashedPassword, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(hashedPassword, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HashPassword provides a mock function with given fields: password
func (_m *PasswordHasher) HashPassword(password string) (string, error) {
	ret := _m.Called(password)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
