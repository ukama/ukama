//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2023-present, Ukama Inc.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.3
// source: events/registryNode.proto

package events

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// added a new node
type EventRegistryNodeCreate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId    string  `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Name      string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Type      string  `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Org       string  `protobuf:"bytes,4,opt,name=org,proto3" json:"org,omitempty"`
	Latitude  float64 `protobuf:"fixed64,5,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude float64 `protobuf:"fixed64,6,opt,name=longitude,proto3" json:"longitude,omitempty"`
}

func (x *EventRegistryNodeCreate) Reset() {
	*x = EventRegistryNodeCreate{}
	mi := &file_events_registryNode_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventRegistryNodeCreate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventRegistryNodeCreate) ProtoMessage() {}

func (x *EventRegistryNodeCreate) ProtoReflect() protoreflect.Message {
	mi := &file_events_registryNode_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventRegistryNodeCreate.ProtoReflect.Descriptor instead.
func (*EventRegistryNodeCreate) Descriptor() ([]byte, []int) {
	return file_events_registryNode_proto_rawDescGZIP(), []int{0}
}

func (x *EventRegistryNodeCreate) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *EventRegistryNodeCreate) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *EventRegistryNodeCreate) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *EventRegistryNodeCreate) GetOrg() string {
	if x != nil {
		return x.Org
	}
	return ""
}

func (x *EventRegistryNodeCreate) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *EventRegistryNodeCreate) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

// updated a node
type EventRegistryNodeUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId    string  `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Name      string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Latitude  float64 `protobuf:"fixed64,3,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude float64 `protobuf:"fixed64,4,opt,name=longitude,proto3" json:"longitude,omitempty"`
}

func (x *EventRegistryNodeUpdate) Reset() {
	*x = EventRegistryNodeUpdate{}
	mi := &file_events_registryNode_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventRegistryNodeUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventRegistryNodeUpdate) ProtoMessage() {}

func (x *EventRegistryNodeUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_events_registryNode_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventRegistryNodeUpdate.ProtoReflect.Descriptor instead.
func (*EventRegistryNodeUpdate) Descriptor() ([]byte, []int) {
	return file_events_registryNode_proto_rawDescGZIP(), []int{1}
}

func (x *EventRegistryNodeUpdate) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *EventRegistryNodeUpdate) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *EventRegistryNodeUpdate) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *EventRegistryNodeUpdate) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

// updated a node state
type EventRegistryNodeStatusUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId string `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	// Types that are assignable to Status:
	//
	//	*EventRegistryNodeStatusUpdate_Connectivity
	//	*EventRegistryNodeStatusUpdate_State
	Status isEventRegistryNodeStatusUpdate_Status `protobuf_oneof:"status"`
}

func (x *EventRegistryNodeStatusUpdate) Reset() {
	*x = EventRegistryNodeStatusUpdate{}
	mi := &file_events_registryNode_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventRegistryNodeStatusUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventRegistryNodeStatusUpdate) ProtoMessage() {}

func (x *EventRegistryNodeStatusUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_events_registryNode_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventRegistryNodeStatusUpdate.ProtoReflect.Descriptor instead.
func (*EventRegistryNodeStatusUpdate) Descriptor() ([]byte, []int) {
	return file_events_registryNode_proto_rawDescGZIP(), []int{2}
}

func (x *EventRegistryNodeStatusUpdate) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (m *EventRegistryNodeStatusUpdate) GetStatus() isEventRegistryNodeStatusUpdate_Status {
	if m != nil {
		return m.Status
	}
	return nil
}

func (x *EventRegistryNodeStatusUpdate) GetConnectivity() string {
	if x, ok := x.GetStatus().(*EventRegistryNodeStatusUpdate_Connectivity); ok {
		return x.Connectivity
	}
	return ""
}

func (x *EventRegistryNodeStatusUpdate) GetState() string {
	if x, ok := x.GetStatus().(*EventRegistryNodeStatusUpdate_State); ok {
		return x.State
	}
	return ""
}

type isEventRegistryNodeStatusUpdate_Status interface {
	isEventRegistryNodeStatusUpdate_Status()
}

type EventRegistryNodeStatusUpdate_Connectivity struct {
	Connectivity string `protobuf:"bytes,2,opt,name=connectivity,proto3,oneof"`
}

type EventRegistryNodeStatusUpdate_State struct {
	State string `protobuf:"bytes,3,opt,name=state,proto3,oneof"`
}

func (*EventRegistryNodeStatusUpdate_Connectivity) isEventRegistryNodeStatusUpdate_Status() {}

func (*EventRegistryNodeStatusUpdate_State) isEventRegistryNodeStatusUpdate_Status() {}

// removed a node
type EventRegistryNodeDelete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId string `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
}

func (x *EventRegistryNodeDelete) Reset() {
	*x = EventRegistryNodeDelete{}
	mi := &file_events_registryNode_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventRegistryNodeDelete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventRegistryNodeDelete) ProtoMessage() {}

func (x *EventRegistryNodeDelete) ProtoReflect() protoreflect.Message {
	mi := &file_events_registryNode_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventRegistryNodeDelete.ProtoReflect.Descriptor instead.
func (*EventRegistryNodeDelete) Descriptor() ([]byte, []int) {
	return file_events_registryNode_proto_rawDescGZIP(), []int{3}
}

func (x *EventRegistryNodeDelete) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

// Assigned to a site
type EventRegistryNodeAssign struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId  string `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Type    string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Network string `protobuf:"bytes,4,opt,name=network,proto3" json:"network,omitempty"`
	Site    string `protobuf:"bytes,5,opt,name=site,proto3" json:"site,omitempty"`
}

func (x *EventRegistryNodeAssign) Reset() {
	*x = EventRegistryNodeAssign{}
	mi := &file_events_registryNode_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventRegistryNodeAssign) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventRegistryNodeAssign) ProtoMessage() {}

func (x *EventRegistryNodeAssign) ProtoReflect() protoreflect.Message {
	mi := &file_events_registryNode_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventRegistryNodeAssign.ProtoReflect.Descriptor instead.
func (*EventRegistryNodeAssign) Descriptor() ([]byte, []int) {
	return file_events_registryNode_proto_rawDescGZIP(), []int{4}
}

func (x *EventRegistryNodeAssign) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *EventRegistryNodeAssign) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *EventRegistryNodeAssign) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *EventRegistryNodeAssign) GetSite() string {
	if x != nil {
		return x.Site
	}
	return ""
}

// Release from site
type EventRegistryNodeRelease struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId  string `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Type    string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Network string `protobuf:"bytes,4,opt,name=network,proto3" json:"network,omitempty"`
	Site    string `protobuf:"bytes,5,opt,name=site,proto3" json:"site,omitempty"`
}

func (x *EventRegistryNodeRelease) Reset() {
	*x = EventRegistryNodeRelease{}
	mi := &file_events_registryNode_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventRegistryNodeRelease) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventRegistryNodeRelease) ProtoMessage() {}

func (x *EventRegistryNodeRelease) ProtoReflect() protoreflect.Message {
	mi := &file_events_registryNode_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventRegistryNodeRelease.ProtoReflect.Descriptor instead.
func (*EventRegistryNodeRelease) Descriptor() ([]byte, []int) {
	return file_events_registryNode_proto_rawDescGZIP(), []int{5}
}

func (x *EventRegistryNodeRelease) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *EventRegistryNodeRelease) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *EventRegistryNodeRelease) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *EventRegistryNodeRelease) GetSite() string {
	if x != nil {
		return x.Site
	}
	return ""
}

// Attach node
type EventRegistryNodeAttach struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId    string   `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Nodegroup []string `protobuf:"bytes,2,rep,name=nodegroup,proto3" json:"nodegroup,omitempty"`
}

func (x *EventRegistryNodeAttach) Reset() {
	*x = EventRegistryNodeAttach{}
	mi := &file_events_registryNode_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventRegistryNodeAttach) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventRegistryNodeAttach) ProtoMessage() {}

func (x *EventRegistryNodeAttach) ProtoReflect() protoreflect.Message {
	mi := &file_events_registryNode_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventRegistryNodeAttach.ProtoReflect.Descriptor instead.
func (*EventRegistryNodeAttach) Descriptor() ([]byte, []int) {
	return file_events_registryNode_proto_rawDescGZIP(), []int{6}
}

func (x *EventRegistryNodeAttach) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *EventRegistryNodeAttach) GetNodegroup() []string {
	if x != nil {
		return x.Nodegroup
	}
	return nil
}

// Dettach node
type EventRegistryNodeDettach struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId    string   `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Nodegroup []string `protobuf:"bytes,2,rep,name=nodegroup,proto3" json:"nodegroup,omitempty"`
}

func (x *EventRegistryNodeDettach) Reset() {
	*x = EventRegistryNodeDettach{}
	mi := &file_events_registryNode_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventRegistryNodeDettach) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventRegistryNodeDettach) ProtoMessage() {}

func (x *EventRegistryNodeDettach) ProtoReflect() protoreflect.Message {
	mi := &file_events_registryNode_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventRegistryNodeDettach.ProtoReflect.Descriptor instead.
func (*EventRegistryNodeDettach) Descriptor() ([]byte, []int) {
	return file_events_registryNode_proto_rawDescGZIP(), []int{7}
}

func (x *EventRegistryNodeDettach) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *EventRegistryNodeDettach) GetNodegroup() []string {
	if x != nil {
		return x.Nodegroup
	}
	return nil
}

var File_events_registryNode_proto protoreflect.FileDescriptor

var file_events_registryNode_proto_rawDesc = []byte{
	0x0a, 0x19, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x4e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x75, 0x6b, 0x61,
	0x6d, 0x61, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x22, 0xa5, 0x01, 0x0a,
	0x17, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x4e, 0x6f,
	0x64, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6f, 0x72, 0x67, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6f, 0x72, 0x67, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61,
	0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6c, 0x61,
	0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74,
	0x75, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69,
	0x74, 0x75, 0x64, 0x65, 0x22, 0x7f, 0x0a, 0x17, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x4e, 0x6f, 0x64, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c,
	0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6c,
	0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69,
	0x74, 0x75, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67,
	0x69, 0x74, 0x75, 0x64, 0x65, 0x22, 0x7f, 0x0a, 0x1d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x4e, 0x6f, 0x64, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x24,
	0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x76, 0x69, 0x74, 0x79, 0x12, 0x16, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x42, 0x08, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x31, 0x0a, 0x17, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x4e, 0x6f, 0x64, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0x73, 0x0a, 0x17, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x4e, 0x6f, 0x64, 0x65, 0x41, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69,
	0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x69, 0x74, 0x65, 0x22, 0x74,
	0x0a, 0x18, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x4e,
	0x6f, 0x64, 0x65, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f,
	0x64, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65,
	0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x73, 0x69, 0x74, 0x65, 0x22, 0x4f, 0x0a, 0x17, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x4e, 0x6f, 0x64, 0x65, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x12,
	0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x6f, 0x64, 0x65,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x50, 0x0a, 0x18, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x4e, 0x6f, 0x64, 0x65, 0x44, 0x65, 0x74, 0x74, 0x61, 0x63,
	0x68, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x6f, 0x64,
	0x65, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x6f,
	0x64, 0x65, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2f, 0x75, 0x6b, 0x61, 0x6d,
	0x61, 0x2f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2f, 0x70, 0x62, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_events_registryNode_proto_rawDescOnce sync.Once
	file_events_registryNode_proto_rawDescData = file_events_registryNode_proto_rawDesc
)

func file_events_registryNode_proto_rawDescGZIP() []byte {
	file_events_registryNode_proto_rawDescOnce.Do(func() {
		file_events_registryNode_proto_rawDescData = protoimpl.X.CompressGZIP(file_events_registryNode_proto_rawDescData)
	})
	return file_events_registryNode_proto_rawDescData
}

var file_events_registryNode_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_events_registryNode_proto_goTypes = []any{
	(*EventRegistryNodeCreate)(nil),       // 0: ukama.events.v1.EventRegistryNodeCreate
	(*EventRegistryNodeUpdate)(nil),       // 1: ukama.events.v1.EventRegistryNodeUpdate
	(*EventRegistryNodeStatusUpdate)(nil), // 2: ukama.events.v1.EventRegistryNodeStatusUpdate
	(*EventRegistryNodeDelete)(nil),       // 3: ukama.events.v1.EventRegistryNodeDelete
	(*EventRegistryNodeAssign)(nil),       // 4: ukama.events.v1.EventRegistryNodeAssign
	(*EventRegistryNodeRelease)(nil),      // 5: ukama.events.v1.EventRegistryNodeRelease
	(*EventRegistryNodeAttach)(nil),       // 6: ukama.events.v1.EventRegistryNodeAttach
	(*EventRegistryNodeDettach)(nil),      // 7: ukama.events.v1.EventRegistryNodeDettach
}
var file_events_registryNode_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_events_registryNode_proto_init() }
func file_events_registryNode_proto_init() {
	if File_events_registryNode_proto != nil {
		return
	}
	file_events_registryNode_proto_msgTypes[2].OneofWrappers = []any{
		(*EventRegistryNodeStatusUpdate_Connectivity)(nil),
		(*EventRegistryNodeStatusUpdate_State)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_events_registryNode_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_events_registryNode_proto_goTypes,
		DependencyIndexes: file_events_registryNode_proto_depIdxs,
		MessageInfos:      file_events_registryNode_proto_msgTypes,
	}.Build()
	File_events_registryNode_proto = out.File
	file_events_registryNode_proto_rawDesc = nil
	file_events_registryNode_proto_goTypes = nil
	file_events_registryNode_proto_depIdxs = nil
}
