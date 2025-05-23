// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ConfigStoreProvider is an autogenerated mock type for the ConfigStoreProvider type
type ConfigStoreProvider struct {
	mock.Mock
}

// HandleConfigCommitReq provides a mock function with given fields: ctx, rVer
func (_m *ConfigStoreProvider) HandleConfigCommitReq(ctx context.Context, rVer string) error {
	ret := _m.Called(ctx, rVer)

	if len(ret) == 0 {
		panic("no return value specified for HandleConfigCommitReq")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, rVer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HandleConfigCommitReqForNode provides a mock function with given fields: ctx, rVer, nodeid
func (_m *ConfigStoreProvider) HandleConfigCommitReqForNode(ctx context.Context, rVer string, nodeid string) error {
	ret := _m.Called(ctx, rVer, nodeid)

	if len(ret) == 0 {
		panic("no return value specified for HandleConfigCommitReqForNode")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, rVer, nodeid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HandleConfigStoreEvent provides a mock function with given fields: ctx
func (_m *ConfigStoreProvider) HandleConfigStoreEvent(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for HandleConfigStoreEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewConfigStoreProvider creates a new instance of ConfigStoreProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConfigStoreProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConfigStoreProvider {
	mock := &ConfigStoreProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
