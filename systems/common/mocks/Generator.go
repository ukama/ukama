// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	uuid "github.com/ukama/ukama/systems/common/uuid"
)

// Generator is an autogenerated mock type for the Generator type
type Generator struct {
	mock.Mock
}

// NewV1 provides a mock function with no fields
func (_m *Generator) NewV1() uuid.UUID {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NewV1")
	}

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func() uuid.UUID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	return r0
}

// NewV2 provides a mock function with given fields: domain
func (_m *Generator) NewV2(domain byte) uuid.UUID {
	ret := _m.Called(domain)

	if len(ret) == 0 {
		panic("no return value specified for NewV2")
	}

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(byte) uuid.UUID); ok {
		r0 = rf(domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	return r0
}

// NewV3 provides a mock function with given fields: ns, name
func (_m *Generator) NewV3(ns uuid.UUID, name string) uuid.UUID {
	ret := _m.Called(ns, name)

	if len(ret) == 0 {
		panic("no return value specified for NewV3")
	}

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(uuid.UUID, string) uuid.UUID); ok {
		r0 = rf(ns, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	return r0
}

// NewV4 provides a mock function with no fields
func (_m *Generator) NewV4() uuid.UUID {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NewV4")
	}

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func() uuid.UUID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	return r0
}

// NewV5 provides a mock function with given fields: ns, name
func (_m *Generator) NewV5(ns uuid.UUID, name string) uuid.UUID {
	ret := _m.Called(ns, name)

	if len(ret) == 0 {
		panic("no return value specified for NewV5")
	}

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(uuid.UUID, string) uuid.UUID); ok {
		r0 = rf(ns, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	return r0
}

// NewGenerator creates a new instance of Generator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *Generator {
	mock := &Generator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
