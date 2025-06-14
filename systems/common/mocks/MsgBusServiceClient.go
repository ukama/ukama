// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

// MsgBusServiceClient is an autogenerated mock type for the MsgBusServiceClient type
type MsgBusServiceClient struct {
	mock.Mock
}

// PublishRequest provides a mock function with given fields: route, msg
func (_m *MsgBusServiceClient) PublishRequest(route string, msg protoreflect.ProtoMessage) error {
	ret := _m.Called(route, msg)

	if len(ret) == 0 {
		panic("no return value specified for PublishRequest")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, protoreflect.ProtoMessage) error); ok {
		r0 = rf(route, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Register provides a mock function with no fields
func (_m *MsgBusServiceClient) Register() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with no fields
func (_m *MsgBusServiceClient) Start() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with no fields
func (_m *MsgBusServiceClient) Stop() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Stop")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMsgBusServiceClient creates a new instance of MsgBusServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMsgBusServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MsgBusServiceClient {
	mock := &MsgBusServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
