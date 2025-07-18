// Code generated by mockery v2.53.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	gen "github.com/ukama/ukama/systems/inventory/component/pb/gen"
)

// ComponentServiceServer is an autogenerated mock type for the ComponentServiceServer type
type ComponentServiceServer struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *ComponentServiceServer) Get(_a0 context.Context, _a1 *gen.GetRequest) (*gen.GetResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *gen.GetResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetRequest) (*gen.GetResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetRequest) *gen.GetResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUser provides a mock function with given fields: _a0, _a1
func (_m *ComponentServiceServer) GetByUser(_a0 context.Context, _a1 *gen.GetByUserRequest) (*gen.GetByUserResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetByUser")
	}

	var r0 *gen.GetByUserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetByUserRequest) (*gen.GetByUserResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetByUserRequest) *gen.GetByUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetByUserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetByUserRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SyncComponents provides a mock function with given fields: _a0, _a1
func (_m *ComponentServiceServer) SyncComponents(_a0 context.Context, _a1 *gen.SyncComponentsRequest) (*gen.SyncComponentsResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SyncComponents")
	}

	var r0 *gen.SyncComponentsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.SyncComponentsRequest) (*gen.SyncComponentsResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.SyncComponentsRequest) *gen.SyncComponentsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.SyncComponentsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.SyncComponentsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mustEmbedUnimplementedComponentServiceServer provides a mock function with no fields
func (_m *ComponentServiceServer) mustEmbedUnimplementedComponentServiceServer() {
	_m.Called()
}

// NewComponentServiceServer creates a new instance of ComponentServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewComponentServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *ComponentServiceServer {
	mock := &ComponentServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
