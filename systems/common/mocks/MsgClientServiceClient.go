// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	msgclient "github.com/ukama/ukama/systems/common/pb/gen/msgclient"
)

// MsgClientServiceClient is an autogenerated mock type for the MsgClientServiceClient type
type MsgClientServiceClient struct {
	mock.Mock
}

// CreateShovel provides a mock function with given fields: ctx, in, opts
func (_m *MsgClientServiceClient) CreateShovel(ctx context.Context, in *msgclient.CreateShovelRequest, opts ...grpc.CallOption) (*msgclient.CreateShovelResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CreateShovel")
	}

	var r0 *msgclient.CreateShovelResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.CreateShovelRequest, ...grpc.CallOption) (*msgclient.CreateShovelResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.CreateShovelRequest, ...grpc.CallOption) *msgclient.CreateShovelResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*msgclient.CreateShovelResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *msgclient.CreateShovelRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PublishMsg provides a mock function with given fields: ctx, in, opts
func (_m *MsgClientServiceClient) PublishMsg(ctx context.Context, in *msgclient.PublishMsgRequest, opts ...grpc.CallOption) (*msgclient.PublishMsgResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for PublishMsg")
	}

	var r0 *msgclient.PublishMsgResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.PublishMsgRequest, ...grpc.CallOption) (*msgclient.PublishMsgResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.PublishMsgRequest, ...grpc.CallOption) *msgclient.PublishMsgResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*msgclient.PublishMsgResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *msgclient.PublishMsgRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterService provides a mock function with given fields: ctx, in, opts
func (_m *MsgClientServiceClient) RegisterService(ctx context.Context, in *msgclient.RegisterServiceReq, opts ...grpc.CallOption) (*msgclient.RegisterServiceResp, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for RegisterService")
	}

	var r0 *msgclient.RegisterServiceResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.RegisterServiceReq, ...grpc.CallOption) (*msgclient.RegisterServiceResp, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.RegisterServiceReq, ...grpc.CallOption) *msgclient.RegisterServiceResp); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*msgclient.RegisterServiceResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *msgclient.RegisterServiceReq, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveShovel provides a mock function with given fields: ctx, in, opts
func (_m *MsgClientServiceClient) RemoveShovel(ctx context.Context, in *msgclient.RemoveShovelRequest, opts ...grpc.CallOption) (*msgclient.RemoveShovelResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for RemoveShovel")
	}

	var r0 *msgclient.RemoveShovelResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.RemoveShovelRequest, ...grpc.CallOption) (*msgclient.RemoveShovelResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.RemoveShovelRequest, ...grpc.CallOption) *msgclient.RemoveShovelResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*msgclient.RemoveShovelResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *msgclient.RemoveShovelRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StartMsgBusHandler provides a mock function with given fields: ctx, in, opts
func (_m *MsgClientServiceClient) StartMsgBusHandler(ctx context.Context, in *msgclient.StartMsgBusHandlerReq, opts ...grpc.CallOption) (*msgclient.StartMsgBusHandlerResp, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for StartMsgBusHandler")
	}

	var r0 *msgclient.StartMsgBusHandlerResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.StartMsgBusHandlerReq, ...grpc.CallOption) (*msgclient.StartMsgBusHandlerResp, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.StartMsgBusHandlerReq, ...grpc.CallOption) *msgclient.StartMsgBusHandlerResp); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*msgclient.StartMsgBusHandlerResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *msgclient.StartMsgBusHandlerReq, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StopMsgBusHandler provides a mock function with given fields: ctx, in, opts
func (_m *MsgClientServiceClient) StopMsgBusHandler(ctx context.Context, in *msgclient.StopMsgBusHandlerReq, opts ...grpc.CallOption) (*msgclient.StopMsgBusHandlerResp, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for StopMsgBusHandler")
	}

	var r0 *msgclient.StopMsgBusHandlerResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.StopMsgBusHandlerReq, ...grpc.CallOption) (*msgclient.StopMsgBusHandlerResp, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.StopMsgBusHandlerReq, ...grpc.CallOption) *msgclient.StopMsgBusHandlerResp); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*msgclient.StopMsgBusHandlerResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *msgclient.StopMsgBusHandlerReq, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnregisterService provides a mock function with given fields: ctx, in, opts
func (_m *MsgClientServiceClient) UnregisterService(ctx context.Context, in *msgclient.UnregisterServiceReq, opts ...grpc.CallOption) (*msgclient.UnregisterServiceResp, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for UnregisterService")
	}

	var r0 *msgclient.UnregisterServiceResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.UnregisterServiceReq, ...grpc.CallOption) (*msgclient.UnregisterServiceResp, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *msgclient.UnregisterServiceReq, ...grpc.CallOption) *msgclient.UnregisterServiceResp); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*msgclient.UnregisterServiceResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *msgclient.UnregisterServiceReq, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMsgClientServiceClient creates a new instance of MsgClientServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMsgClientServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MsgClientServiceClient {
	mock := &MsgClientServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
