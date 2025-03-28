// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	db "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
)

// GutiRepo is an autogenerated mock type for the GutiRepo type
type GutiRepo struct {
	mock.Mock
}

// GetImsi provides a mock function with given fields: guti
func (_m *GutiRepo) GetImsi(guti string) (string, error) {
	ret := _m.Called(guti)

	if len(ret) == 0 {
		panic("no return value specified for GetImsi")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(guti)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(guti)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(guti)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: guti
func (_m *GutiRepo) Update(guti *db.Guti) error {
	ret := _m.Called(guti)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*db.Guti) error); ok {
		r0 = rf(guti)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewGutiRepo creates a new instance of GutiRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGutiRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *GutiRepo {
	mock := &GutiRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
