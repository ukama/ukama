// Code generated by mockery v2.53.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UnsafeSimServiceServer is an autogenerated mock type for the UnsafeSimServiceServer type
type UnsafeSimServiceServer struct {
	mock.Mock
}

// mustEmbedUnimplementedSimServiceServer provides a mock function with no fields
func (_m *UnsafeSimServiceServer) mustEmbedUnimplementedSimServiceServer() {
	_m.Called()
}

// NewUnsafeSimServiceServer creates a new instance of UnsafeSimServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUnsafeSimServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *UnsafeSimServiceServer {
	mock := &UnsafeSimServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
