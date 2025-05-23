// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	db "github.com/ukama/ukama/systems/init/lookup/internal/db"
)

// SystemRepo is an autogenerated mock type for the SystemRepo type
type SystemRepo struct {
	mock.Mock
}

// Add provides a mock function with given fields: sys
func (_m *SystemRepo) Add(sys *db.System) error {
	ret := _m.Called(sys)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*db.System) error); ok {
		r0 = rf(sys)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: sys, org
func (_m *SystemRepo) Delete(sys string, org uint) error {
	ret := _m.Called(sys, org)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, uint) error); ok {
		r0 = rf(sys, org)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByName provides a mock function with given fields: sys, org
func (_m *SystemRepo) GetByName(sys string, org uint) (*db.System, error) {
	ret := _m.Called(sys, org)

	if len(ret) == 0 {
		panic("no return value specified for GetByName")
	}

	var r0 *db.System
	var r1 error
	if rf, ok := ret.Get(0).(func(string, uint) (*db.System, error)); ok {
		return rf(sys, org)
	}
	if rf, ok := ret.Get(0).(func(string, uint) *db.System); ok {
		r0 = rf(sys, org)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*db.System)
		}
	}

	if rf, ok := ret.Get(1).(func(string, uint) error); ok {
		r1 = rf(sys, org)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: sys, org
func (_m *SystemRepo) Update(sys *db.System, org uint) error {
	ret := _m.Called(sys, org)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*db.System, uint) error); ok {
		r0 = rf(sys, org)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewSystemRepo creates a new instance of SystemRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSystemRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *SystemRepo {
	mock := &SystemRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
