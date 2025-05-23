// Code generated by mockery v2.53.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	gen "github.com/ukama/ukama/systems/node/software/pb/gen"
)

// softwareManager is an autogenerated mock type for the softwareManager type
type softwareManager struct {
	mock.Mock
}

// UpdateSoftware provides a mock function with given fields: space, name, tag, nodeId
func (_m *softwareManager) UpdateSoftware(space string, name string, tag string, nodeId string) (*gen.UpdateSoftwareResponse, error) {
	ret := _m.Called(space, name, tag, nodeId)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSoftware")
	}

	var r0 *gen.UpdateSoftwareResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, string, string) (*gen.UpdateSoftwareResponse, error)); ok {
		return rf(space, name, tag, nodeId)
	}
	if rf, ok := ret.Get(0).(func(string, string, string, string) *gen.UpdateSoftwareResponse); ok {
		r0 = rf(space, name, tag, nodeId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.UpdateSoftwareResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, string, string) error); ok {
		r1 = rf(space, name, tag, nodeId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// newSoftwareManager creates a new instance of softwareManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newSoftwareManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *softwareManager {
	mock := &softwareManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
