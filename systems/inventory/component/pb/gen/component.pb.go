//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2023-present, Ukama Inc.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: component.proto

package gen

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

type ComponentCategory int32

const (
	ComponentCategory_ALL      ComponentCategory = 0
	ComponentCategory_ACCESS   ComponentCategory = 1
	ComponentCategory_BACKHAUL ComponentCategory = 2
	ComponentCategory_POWER    ComponentCategory = 3
	ComponentCategory_SWITCH   ComponentCategory = 4
	ComponentCategory_SPECTRUM ComponentCategory = 5
)

// Enum value maps for ComponentCategory.
var (
	ComponentCategory_name = map[int32]string{
		0: "ALL",
		1: "ACCESS",
		2: "BACKHAUL",
		3: "POWER",
		4: "SWITCH",
		5: "SPECTRUM",
	}
	ComponentCategory_value = map[string]int32{
		"ALL":      0,
		"ACCESS":   1,
		"BACKHAUL": 2,
		"POWER":    3,
		"SWITCH":   4,
		"SPECTRUM": 5,
	}
)

func (x ComponentCategory) Enum() *ComponentCategory {
	p := new(ComponentCategory)
	*p = x
	return p
}

func (x ComponentCategory) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ComponentCategory) Descriptor() protoreflect.EnumDescriptor {
	return file_component_proto_enumTypes[0].Descriptor()
}

func (ComponentCategory) Type() protoreflect.EnumType {
	return &file_component_proto_enumTypes[0]
}

func (x ComponentCategory) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ComponentCategory.Descriptor instead.
func (ComponentCategory) EnumDescriptor() ([]byte, []int) {
	return file_component_proto_rawDescGZIP(), []int{0}
}

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_component_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_component_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_component_proto_rawDescGZIP(), []int{0}
}

func (x *GetRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Component *Component `protobuf:"bytes,1,opt,name=component,proto3" json:"component,omitempty"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_component_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_component_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_component_proto_rawDescGZIP(), []int{1}
}

func (x *GetResponse) GetComponent() *Component {
	if x != nil {
		return x.Component
	}
	return nil
}

type GetByUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   string            `protobuf:"bytes,1,opt,name=userId,json=user_id,proto3" json:"userId,omitempty"`
	Category ComponentCategory `protobuf:"varint,2,opt,name=category,proto3,enum=ukama.inventory.component.v1.ComponentCategory" json:"category,omitempty"`
}

func (x *GetByUserRequest) Reset() {
	*x = GetByUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_component_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByUserRequest) ProtoMessage() {}

func (x *GetByUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_component_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByUserRequest.ProtoReflect.Descriptor instead.
func (*GetByUserRequest) Descriptor() ([]byte, []int) {
	return file_component_proto_rawDescGZIP(), []int{2}
}

func (x *GetByUserRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetByUserRequest) GetCategory() ComponentCategory {
	if x != nil {
		return x.Category
	}
	return ComponentCategory_ALL
}

type GetByUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Components []*Component `protobuf:"bytes,1,rep,name=components,proto3" json:"components,omitempty"`
}

func (x *GetByUserResponse) Reset() {
	*x = GetByUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_component_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByUserResponse) ProtoMessage() {}

func (x *GetByUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_component_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByUserResponse.ProtoReflect.Descriptor instead.
func (*GetByUserResponse) Descriptor() ([]byte, []int) {
	return file_component_proto_rawDescGZIP(), []int{3}
}

func (x *GetByUserResponse) GetComponents() []*Component {
	if x != nil {
		return x.Components
	}
	return nil
}

type SyncComponentsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SyncComponentsRequest) Reset() {
	*x = SyncComponentsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_component_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncComponentsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncComponentsRequest) ProtoMessage() {}

func (x *SyncComponentsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_component_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncComponentsRequest.ProtoReflect.Descriptor instead.
func (*SyncComponentsRequest) Descriptor() ([]byte, []int) {
	return file_component_proto_rawDescGZIP(), []int{4}
}

type SyncComponentsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SyncComponentsResponse) Reset() {
	*x = SyncComponentsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_component_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncComponentsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncComponentsResponse) ProtoMessage() {}

func (x *SyncComponentsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_component_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncComponentsResponse.ProtoReflect.Descriptor instead.
func (*SyncComponentsResponse) Descriptor() ([]byte, []int) {
	return file_component_proto_rawDescGZIP(), []int{5}
}

type Component struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Inventory     string            `protobuf:"bytes,2,opt,name=inventory,json=inventory_id,proto3" json:"inventory,omitempty"`
	Category      ComponentCategory `protobuf:"varint,3,opt,name=category,proto3,enum=ukama.inventory.component.v1.ComponentCategory" json:"category,omitempty"`
	Type          string            `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	UserId        string            `protobuf:"bytes,5,opt,name=userId,json=user_id,proto3" json:"userId,omitempty"`
	Description   string            `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	DatasheetURL  string            `protobuf:"bytes,7,opt,name=datasheetURL,json=datasheet_url,proto3" json:"datasheetURL,omitempty"`
	ImagesURL     string            `protobuf:"bytes,8,opt,name=imagesURL,json=images_url,proto3" json:"imagesURL,omitempty"`
	PartNumber    string            `protobuf:"bytes,9,opt,name=partNumber,json=part_number,proto3" json:"partNumber,omitempty"`
	Manufacturer  string            `protobuf:"bytes,10,opt,name=manufacturer,proto3" json:"manufacturer,omitempty"`
	Managed       string            `protobuf:"bytes,11,opt,name=managed,proto3" json:"managed,omitempty"`
	Warranty      uint32            `protobuf:"varint,12,opt,name=warranty,proto3" json:"warranty,omitempty"`
	Specification string            `protobuf:"bytes,13,opt,name=specification,proto3" json:"specification,omitempty"`
}

func (x *Component) Reset() {
	*x = Component{}
	if protoimpl.UnsafeEnabled {
		mi := &file_component_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Component) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Component) ProtoMessage() {}

func (x *Component) ProtoReflect() protoreflect.Message {
	mi := &file_component_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Component.ProtoReflect.Descriptor instead.
func (*Component) Descriptor() ([]byte, []int) {
	return file_component_proto_rawDescGZIP(), []int{6}
}

func (x *Component) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Component) GetInventory() string {
	if x != nil {
		return x.Inventory
	}
	return ""
}

func (x *Component) GetCategory() ComponentCategory {
	if x != nil {
		return x.Category
	}
	return ComponentCategory_ALL
}

func (x *Component) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Component) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Component) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Component) GetDatasheetURL() string {
	if x != nil {
		return x.DatasheetURL
	}
	return ""
}

func (x *Component) GetImagesURL() string {
	if x != nil {
		return x.ImagesURL
	}
	return ""
}

func (x *Component) GetPartNumber() string {
	if x != nil {
		return x.PartNumber
	}
	return ""
}

func (x *Component) GetManufacturer() string {
	if x != nil {
		return x.Manufacturer
	}
	return ""
}

func (x *Component) GetManaged() string {
	if x != nil {
		return x.Managed
	}
	return ""
}

func (x *Component) GetWarranty() uint32 {
	if x != nil {
		return x.Warranty
	}
	return 0
}

func (x *Component) GetSpecification() string {
	if x != nil {
		return x.Specification
	}
	return ""
}

var File_component_proto protoreflect.FileDescriptor

var file_component_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x1c, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f,
	0x72, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x22,
	0x1c, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x54, 0x0a,
	0x0b, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a, 0x09,
	0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x27, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72,
	0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e,
	0x65, 0x6e, 0x74, 0x22, 0x78, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x12, 0x4b, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x2f, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e,
	0x74, 0x6f, 0x72, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x22, 0x5c, 0x0a,
	0x11, 0x47, 0x65, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x47, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69,
	0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65,
	0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x52,
	0x0a, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x17, 0x0a, 0x15, 0x53,
	0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x18, 0x0a, 0x16, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6d, 0x70,
	0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xbd,
	0x03, 0x0a, 0x09, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x09,
	0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x12, 0x4b, 0x0a,
	0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x2f, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72,
	0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x17,
	0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0c, 0x64, 0x61, 0x74,
	0x61, 0x73, 0x68, 0x65, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x64, 0x61, 0x74, 0x61, 0x73, 0x68, 0x65, 0x65, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x12, 0x1d,
	0x0a, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x55, 0x52, 0x4c, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x5f, 0x75, 0x72, 0x6c, 0x12, 0x1f, 0x0a,
	0x0a, 0x70, 0x61, 0x72, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x70, 0x61, 0x72, 0x74, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x22,
	0x0a, 0x0c, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x65, 0x72, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72,
	0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x64, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x77, 0x61, 0x72, 0x72, 0x61, 0x6e, 0x74, 0x79, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08,
	0x77, 0x61, 0x72, 0x72, 0x61, 0x6e, 0x74, 0x79, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x70, 0x65, 0x63,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2a, 0x5b,
	0x0a, 0x11, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x4c, 0x4c, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06,
	0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x42, 0x41, 0x43, 0x4b,
	0x48, 0x41, 0x55, 0x4c, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x50, 0x4f, 0x57, 0x45, 0x52, 0x10,
	0x03, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x57, 0x49, 0x54, 0x43, 0x48, 0x10, 0x04, 0x12, 0x0c, 0x0a,
	0x08, 0x53, 0x50, 0x45, 0x43, 0x54, 0x52, 0x55, 0x4d, 0x10, 0x05, 0x32, 0xd9, 0x02, 0x0a, 0x10,
	0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x5a, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x28, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e,
	0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e,
	0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x29, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74,
	0x6f, 0x72, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6c, 0x0a, 0x09,
	0x47, 0x65, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x12, 0x2e, 0x2e, 0x75, 0x6b, 0x61, 0x6d,
	0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x70,
	0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x75, 0x6b, 0x61, 0x6d,
	0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x70,
	0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x7b, 0x0a, 0x0e, 0x53, 0x79,
	0x6e, 0x63, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x33, 0x2e, 0x75,
	0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x63,
	0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x79, 0x6e, 0x63,
	0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x34, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74,
	0x6f, 0x72, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2f, 0x75, 0x6b, 0x61, 0x6d,
	0x61, 0x2f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74,
	0x6f, 0x72, 0x79, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x62,
	0x2f, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_component_proto_rawDescOnce sync.Once
	file_component_proto_rawDescData = file_component_proto_rawDesc
)

func file_component_proto_rawDescGZIP() []byte {
	file_component_proto_rawDescOnce.Do(func() {
		file_component_proto_rawDescData = protoimpl.X.CompressGZIP(file_component_proto_rawDescData)
	})
	return file_component_proto_rawDescData
}

var file_component_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_component_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_component_proto_goTypes = []interface{}{
	(ComponentCategory)(0),         // 0: ukama.inventory.component.v1.ComponentCategory
	(*GetRequest)(nil),             // 1: ukama.inventory.component.v1.GetRequest
	(*GetResponse)(nil),            // 2: ukama.inventory.component.v1.GetResponse
	(*GetByUserRequest)(nil),       // 3: ukama.inventory.component.v1.GetByUserRequest
	(*GetByUserResponse)(nil),      // 4: ukama.inventory.component.v1.GetByUserResponse
	(*SyncComponentsRequest)(nil),  // 5: ukama.inventory.component.v1.SyncComponentsRequest
	(*SyncComponentsResponse)(nil), // 6: ukama.inventory.component.v1.SyncComponentsResponse
	(*Component)(nil),              // 7: ukama.inventory.component.v1.Component
}
var file_component_proto_depIdxs = []int32{
	7, // 0: ukama.inventory.component.v1.GetResponse.component:type_name -> ukama.inventory.component.v1.Component
	0, // 1: ukama.inventory.component.v1.GetByUserRequest.category:type_name -> ukama.inventory.component.v1.ComponentCategory
	7, // 2: ukama.inventory.component.v1.GetByUserResponse.components:type_name -> ukama.inventory.component.v1.Component
	0, // 3: ukama.inventory.component.v1.Component.category:type_name -> ukama.inventory.component.v1.ComponentCategory
	1, // 4: ukama.inventory.component.v1.ComponentService.Get:input_type -> ukama.inventory.component.v1.GetRequest
	3, // 5: ukama.inventory.component.v1.ComponentService.GetByUser:input_type -> ukama.inventory.component.v1.GetByUserRequest
	5, // 6: ukama.inventory.component.v1.ComponentService.SyncComponents:input_type -> ukama.inventory.component.v1.SyncComponentsRequest
	2, // 7: ukama.inventory.component.v1.ComponentService.Get:output_type -> ukama.inventory.component.v1.GetResponse
	4, // 8: ukama.inventory.component.v1.ComponentService.GetByUser:output_type -> ukama.inventory.component.v1.GetByUserResponse
	6, // 9: ukama.inventory.component.v1.ComponentService.SyncComponents:output_type -> ukama.inventory.component.v1.SyncComponentsResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_component_proto_init() }
func file_component_proto_init() {
	if File_component_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_component_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_component_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_component_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByUserRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_component_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByUserResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_component_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncComponentsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_component_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncComponentsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_component_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Component); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_component_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_component_proto_goTypes,
		DependencyIndexes: file_component_proto_depIdxs,
		EnumInfos:         file_component_proto_enumTypes,
		MessageInfos:      file_component_proto_msgTypes,
	}.Build()
	File_component_proto = out.File
	file_component_proto_rawDesc = nil
	file_component_proto_goTypes = nil
	file_component_proto_depIdxs = nil
}
