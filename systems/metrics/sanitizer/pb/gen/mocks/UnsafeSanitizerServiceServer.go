// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UnsafeSanitizerServiceServer is an autogenerated mock type for the UnsafeSanitizerServiceServer type
type UnsafeSanitizerServiceServer struct {
	mock.Mock
}

// mustEmbedUnimplementedSanitizerServiceServer provides a mock function with no fields
func (_m *UnsafeSanitizerServiceServer) mustEmbedUnimplementedSanitizerServiceServer() {
	_m.Called()
}

// NewUnsafeSanitizerServiceServer creates a new instance of UnsafeSanitizerServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUnsafeSanitizerServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *UnsafeSanitizerServiceServer {
	mock := &UnsafeSanitizerServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
