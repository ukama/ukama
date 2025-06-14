// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	gen "github.com/ukama/ukama/systems/metrics/sanitizer/pb/gen"
)

// SanitizerServiceServer is an autogenerated mock type for the SanitizerServiceServer type
type SanitizerServiceServer struct {
	mock.Mock
}

// Sanitize provides a mock function with given fields: _a0, _a1
func (_m *SanitizerServiceServer) Sanitize(_a0 context.Context, _a1 *gen.SanitizeRequest) (*gen.SanitizeResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Sanitize")
	}

	var r0 *gen.SanitizeResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.SanitizeRequest) (*gen.SanitizeResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.SanitizeRequest) *gen.SanitizeResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.SanitizeResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.SanitizeRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mustEmbedUnimplementedSanitizerServiceServer provides a mock function with no fields
func (_m *SanitizerServiceServer) mustEmbedUnimplementedSanitizerServiceServer() {
	_m.Called()
}

// NewSanitizerServiceServer creates a new instance of SanitizerServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSanitizerServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *SanitizerServiceServer {
	mock := &SanitizerServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
