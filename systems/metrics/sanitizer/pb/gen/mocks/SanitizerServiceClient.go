// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	gen "github.com/ukama/ukama/systems/metrics/sanitizer/pb/gen"
	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// SanitizerServiceClient is an autogenerated mock type for the SanitizerServiceClient type
type SanitizerServiceClient struct {
	mock.Mock
}

// Sanitize provides a mock function with given fields: ctx, in, opts
func (_m *SanitizerServiceClient) Sanitize(ctx context.Context, in *gen.SanitizeRequest, opts ...grpc.CallOption) (*gen.SanitizeResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Sanitize")
	}

	var r0 *gen.SanitizeResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.SanitizeRequest, ...grpc.CallOption) (*gen.SanitizeResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.SanitizeRequest, ...grpc.CallOption) *gen.SanitizeResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.SanitizeResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.SanitizeRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSanitizerServiceClient creates a new instance of SanitizerServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSanitizerServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *SanitizerServiceClient {
	mock := &SanitizerServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
