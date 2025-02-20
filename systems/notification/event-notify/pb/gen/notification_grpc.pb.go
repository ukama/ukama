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
// source: notification.proto

package gen

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
	EventToNotifyService_Get_FullMethodName          = "/ukama.notification.eventToNotify.v1.EventToNotifyService/Get"
	EventToNotifyService_GetAll_FullMethodName       = "/ukama.notification.eventToNotify.v1.EventToNotifyService/GetAll"
	EventToNotifyService_UpdateStatus_FullMethodName = "/ukama.notification.eventToNotify.v1.EventToNotifyService/UpdateStatus"
)

// EventToNotifyServiceClient is the client API for EventToNotifyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventToNotifyServiceClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*GetAllResponse, error)
	UpdateStatus(ctx context.Context, in *UpdateStatusRequest, opts ...grpc.CallOption) (*UpdateStatusResponse, error)
}

type eventToNotifyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventToNotifyServiceClient(cc grpc.ClientConnInterface) EventToNotifyServiceClient {
	return &eventToNotifyServiceClient{cc}
}

func (c *eventToNotifyServiceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, EventToNotifyService_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventToNotifyServiceClient) GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*GetAllResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllResponse)
	err := c.cc.Invoke(ctx, EventToNotifyService_GetAll_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventToNotifyServiceClient) UpdateStatus(ctx context.Context, in *UpdateStatusRequest, opts ...grpc.CallOption) (*UpdateStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateStatusResponse)
	err := c.cc.Invoke(ctx, EventToNotifyService_UpdateStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventToNotifyServiceServer is the server API for EventToNotifyService service.
// All implementations must embed UnimplementedEventToNotifyServiceServer
// for forward compatibility.
type EventToNotifyServiceServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	GetAll(context.Context, *GetAllRequest) (*GetAllResponse, error)
	UpdateStatus(context.Context, *UpdateStatusRequest) (*UpdateStatusResponse, error)
	mustEmbedUnimplementedEventToNotifyServiceServer()
}

// UnimplementedEventToNotifyServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEventToNotifyServiceServer struct{}

func (UnimplementedEventToNotifyServiceServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedEventToNotifyServiceServer) GetAll(context.Context, *GetAllRequest) (*GetAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedEventToNotifyServiceServer) UpdateStatus(context.Context, *UpdateStatusRequest) (*UpdateStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStatus not implemented")
}
func (UnimplementedEventToNotifyServiceServer) mustEmbedUnimplementedEventToNotifyServiceServer() {}
func (UnimplementedEventToNotifyServiceServer) testEmbeddedByValue()                              {}

// UnsafeEventToNotifyServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventToNotifyServiceServer will
// result in compilation errors.
type UnsafeEventToNotifyServiceServer interface {
	mustEmbedUnimplementedEventToNotifyServiceServer()
}

func RegisterEventToNotifyServiceServer(s grpc.ServiceRegistrar, srv EventToNotifyServiceServer) {
	// If the following call pancis, it indicates UnimplementedEventToNotifyServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EventToNotifyService_ServiceDesc, srv)
}

func _EventToNotifyService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventToNotifyServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventToNotifyService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventToNotifyServiceServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventToNotifyService_GetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventToNotifyServiceServer).GetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventToNotifyService_GetAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventToNotifyServiceServer).GetAll(ctx, req.(*GetAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventToNotifyService_UpdateStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventToNotifyServiceServer).UpdateStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventToNotifyService_UpdateStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventToNotifyServiceServer).UpdateStatus(ctx, req.(*UpdateStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EventToNotifyService_ServiceDesc is the grpc.ServiceDesc for EventToNotifyService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventToNotifyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ukama.notification.eventToNotify.v1.EventToNotifyService",
	HandlerType: (*EventToNotifyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _EventToNotifyService_Get_Handler,
		},
		{
			MethodName: "GetAll",
			Handler:    _EventToNotifyService_GetAll_Handler,
		},
		{
			MethodName: "UpdateStatus",
			Handler:    _EventToNotifyService_UpdateStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notification.proto",
}
