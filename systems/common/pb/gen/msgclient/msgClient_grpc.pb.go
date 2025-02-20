//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2023-present, Ukama Inc.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: msgclient/msgClient.proto

package msgclient

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MsgClientService_RegisterService_FullMethodName    = "/ukama.msgClient.v1.MsgClientService/RegisterService"
	MsgClientService_StartMsgBusHandler_FullMethodName = "/ukama.msgClient.v1.MsgClientService/StartMsgBusHandler"
	MsgClientService_StopMsgBusHandler_FullMethodName  = "/ukama.msgClient.v1.MsgClientService/StopMsgBusHandler"
	MsgClientService_UnregisterService_FullMethodName  = "/ukama.msgClient.v1.MsgClientService/UnregisterService"
	MsgClientService_PublishMsg_FullMethodName         = "/ukama.msgClient.v1.MsgClientService/PublishMsg"
	MsgClientService_CreateShovel_FullMethodName       = "/ukama.msgClient.v1.MsgClientService/CreateShovel"
	MsgClientService_RemoveShovel_FullMethodName       = "/ukama.msgClient.v1.MsgClientService/RemoveShovel"
)

// MsgClientServiceClient is the client API for MsgClientService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClientServiceClient interface {
	// / Use this rpc to register system to MsgClient
	RegisterService(ctx context.Context, in *RegisterServiceReq, opts ...grpc.CallOption) (*RegisterServiceResp, error)
	// / Call this rpc to  StartMsgBus after registration
	StartMsgBusHandler(ctx context.Context, in *StartMsgBusHandlerReq, opts ...grpc.CallOption) (*StartMsgBusHandlerResp, error)
	// / Call this rpc to  StopMsgBus
	StopMsgBusHandler(ctx context.Context, in *StopMsgBusHandlerReq, opts ...grpc.CallOption) (*StopMsgBusHandlerResp, error)
	// / Unregister service from MsgClient
	UnregisterService(ctx context.Context, in *UnregisterServiceReq, opts ...grpc.CallOption) (*UnregisterServiceResp, error)
	// / Call this rpc to publisg events
	PublishMsg(ctx context.Context, in *PublishMsgRequest, opts ...grpc.CallOption) (*PublishMsgResponse, error)
	// / Create a shovel
	CreateShovel(ctx context.Context, in *CreateShovelRequest, opts ...grpc.CallOption) (*CreateShovelResponse, error)
	// / Remove shovel
	RemoveShovel(ctx context.Context, in *RemoveShovelRequest, opts ...grpc.CallOption) (*RemoveShovelResponse, error)
}

type msgClientServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClientServiceClient(cc grpc.ClientConnInterface) MsgClientServiceClient {
	return &msgClientServiceClient{cc}
}

func (c *msgClientServiceClient) RegisterService(ctx context.Context, in *RegisterServiceReq, opts ...grpc.CallOption) (*RegisterServiceResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterServiceResp)
	err := c.cc.Invoke(ctx, MsgClientService_RegisterService_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClientServiceClient) StartMsgBusHandler(ctx context.Context, in *StartMsgBusHandlerReq, opts ...grpc.CallOption) (*StartMsgBusHandlerResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StartMsgBusHandlerResp)
	err := c.cc.Invoke(ctx, MsgClientService_StartMsgBusHandler_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClientServiceClient) StopMsgBusHandler(ctx context.Context, in *StopMsgBusHandlerReq, opts ...grpc.CallOption) (*StopMsgBusHandlerResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StopMsgBusHandlerResp)
	err := c.cc.Invoke(ctx, MsgClientService_StopMsgBusHandler_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClientServiceClient) UnregisterService(ctx context.Context, in *UnregisterServiceReq, opts ...grpc.CallOption) (*UnregisterServiceResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UnregisterServiceResp)
	err := c.cc.Invoke(ctx, MsgClientService_UnregisterService_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClientServiceClient) PublishMsg(ctx context.Context, in *PublishMsgRequest, opts ...grpc.CallOption) (*PublishMsgResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PublishMsgResponse)
	err := c.cc.Invoke(ctx, MsgClientService_PublishMsg_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClientServiceClient) CreateShovel(ctx context.Context, in *CreateShovelRequest, opts ...grpc.CallOption) (*CreateShovelResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateShovelResponse)
	err := c.cc.Invoke(ctx, MsgClientService_CreateShovel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClientServiceClient) RemoveShovel(ctx context.Context, in *RemoveShovelRequest, opts ...grpc.CallOption) (*RemoveShovelResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RemoveShovelResponse)
	err := c.cc.Invoke(ctx, MsgClientService_RemoveShovel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgClientServiceServer is the server API for MsgClientService service.
// All implementations must embed UnimplementedMsgClientServiceServer
// for forward compatibility.
type MsgClientServiceServer interface {
	// / Use this rpc to register system to MsgClient
	RegisterService(context.Context, *RegisterServiceReq) (*RegisterServiceResp, error)
	// / Call this rpc to  StartMsgBus after registration
	StartMsgBusHandler(context.Context, *StartMsgBusHandlerReq) (*StartMsgBusHandlerResp, error)
	// / Call this rpc to  StopMsgBus
	StopMsgBusHandler(context.Context, *StopMsgBusHandlerReq) (*StopMsgBusHandlerResp, error)
	// / Unregister service from MsgClient
	UnregisterService(context.Context, *UnregisterServiceReq) (*UnregisterServiceResp, error)
	// / Call this rpc to publisg events
	PublishMsg(context.Context, *PublishMsgRequest) (*PublishMsgResponse, error)
	// / Create a shovel
	CreateShovel(context.Context, *CreateShovelRequest) (*CreateShovelResponse, error)
	// / Remove shovel
	RemoveShovel(context.Context, *RemoveShovelRequest) (*RemoveShovelResponse, error)
	mustEmbedUnimplementedMsgClientServiceServer()
}

// UnimplementedMsgClientServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMsgClientServiceServer struct{}

func (UnimplementedMsgClientServiceServer) RegisterService(context.Context, *RegisterServiceReq) (*RegisterServiceResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterService not implemented")
}
func (UnimplementedMsgClientServiceServer) StartMsgBusHandler(context.Context, *StartMsgBusHandlerReq) (*StartMsgBusHandlerResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartMsgBusHandler not implemented")
}
func (UnimplementedMsgClientServiceServer) StopMsgBusHandler(context.Context, *StopMsgBusHandlerReq) (*StopMsgBusHandlerResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopMsgBusHandler not implemented")
}
func (UnimplementedMsgClientServiceServer) UnregisterService(context.Context, *UnregisterServiceReq) (*UnregisterServiceResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnregisterService not implemented")
}
func (UnimplementedMsgClientServiceServer) PublishMsg(context.Context, *PublishMsgRequest) (*PublishMsgResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishMsg not implemented")
}
func (UnimplementedMsgClientServiceServer) CreateShovel(context.Context, *CreateShovelRequest) (*CreateShovelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShovel not implemented")
}
func (UnimplementedMsgClientServiceServer) RemoveShovel(context.Context, *RemoveShovelRequest) (*RemoveShovelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveShovel not implemented")
}
func (UnimplementedMsgClientServiceServer) mustEmbedUnimplementedMsgClientServiceServer() {}
func (UnimplementedMsgClientServiceServer) testEmbeddedByValue()                          {}

// UnsafeMsgClientServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgClientServiceServer will
// result in compilation errors.
type UnsafeMsgClientServiceServer interface {
	mustEmbedUnimplementedMsgClientServiceServer()
}

func RegisterMsgClientServiceServer(s grpc.ServiceRegistrar, srv MsgClientServiceServer) {
	// If the following call pancis, it indicates UnimplementedMsgClientServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MsgClientService_ServiceDesc, srv)
}

func _MsgClientService_RegisterService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterServiceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgClientServiceServer).RegisterService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgClientService_RegisterService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgClientServiceServer).RegisterService(ctx, req.(*RegisterServiceReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgClientService_StartMsgBusHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartMsgBusHandlerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgClientServiceServer).StartMsgBusHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgClientService_StartMsgBusHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgClientServiceServer).StartMsgBusHandler(ctx, req.(*StartMsgBusHandlerReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgClientService_StopMsgBusHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopMsgBusHandlerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgClientServiceServer).StopMsgBusHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgClientService_StopMsgBusHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgClientServiceServer).StopMsgBusHandler(ctx, req.(*StopMsgBusHandlerReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgClientService_UnregisterService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnregisterServiceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgClientServiceServer).UnregisterService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgClientService_UnregisterService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgClientServiceServer).UnregisterService(ctx, req.(*UnregisterServiceReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgClientService_PublishMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishMsgRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgClientServiceServer).PublishMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgClientService_PublishMsg_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgClientServiceServer).PublishMsg(ctx, req.(*PublishMsgRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgClientService_CreateShovel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateShovelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgClientServiceServer).CreateShovel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgClientService_CreateShovel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgClientServiceServer).CreateShovel(ctx, req.(*CreateShovelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgClientService_RemoveShovel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveShovelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgClientServiceServer).RemoveShovel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MsgClientService_RemoveShovel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgClientServiceServer).RemoveShovel(ctx, req.(*RemoveShovelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MsgClientService_ServiceDesc is the grpc.ServiceDesc for MsgClientService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MsgClientService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ukama.msgClient.v1.MsgClientService",
	HandlerType: (*MsgClientServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterService",
			Handler:    _MsgClientService_RegisterService_Handler,
		},
		{
			MethodName: "StartMsgBusHandler",
			Handler:    _MsgClientService_StartMsgBusHandler_Handler,
		},
		{
			MethodName: "StopMsgBusHandler",
			Handler:    _MsgClientService_StopMsgBusHandler_Handler,
		},
		{
			MethodName: "UnregisterService",
			Handler:    _MsgClientService_UnregisterService_Handler,
		},
		{
			MethodName: "PublishMsg",
			Handler:    _MsgClientService_PublishMsg_Handler,
		},
		{
			MethodName: "CreateShovel",
			Handler:    _MsgClientService_CreateShovel_Handler,
		},
		{
			MethodName: "RemoveShovel",
			Handler:    _MsgClientService_RemoveShovel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "msgclient/msgClient.proto",
}
