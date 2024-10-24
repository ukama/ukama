// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	gen "github.com/ukama/ukama/systems/registry/site/pb/gen"
	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// SiteServiceClient is an autogenerated mock type for the SiteServiceClient type
type SiteServiceClient struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, in, opts
func (_m *SiteServiceClient) Add(ctx context.Context, in *gen.AddRequest, opts ...grpc.CallOption) (*gen.AddResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *gen.AddResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.AddRequest, ...grpc.CallOption) (*gen.AddResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.AddRequest, ...grpc.CallOption) *gen.AddResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.AddResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.AddRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, in, opts
func (_m *SiteServiceClient) Get(ctx context.Context, in *gen.GetRequest, opts ...grpc.CallOption) (*gen.GetResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *gen.GetResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetRequest, ...grpc.CallOption) (*gen.GetResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetRequest, ...grpc.CallOption) *gen.GetResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSites provides a mock function with given fields: ctx, in, opts
func (_m *SiteServiceClient) GetSites(ctx context.Context, in *gen.GetSitesRequest, opts ...grpc.CallOption) (*gen.GetSitesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *gen.GetSitesResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetSitesRequest, ...grpc.CallOption) (*gen.GetSitesResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetSitesRequest, ...grpc.CallOption) *gen.GetSitesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetSitesResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetSitesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, in, opts
func (_m *SiteServiceClient) Update(ctx context.Context, in *gen.UpdateRequest, opts ...grpc.CallOption) (*gen.UpdateResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *gen.UpdateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.UpdateRequest, ...grpc.CallOption) (*gen.UpdateResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.UpdateRequest, ...grpc.CallOption) *gen.UpdateResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.UpdateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.UpdateRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSiteServiceClient creates a new instance of SiteServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSiteServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *SiteServiceClient {
	mock := &SiteServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
