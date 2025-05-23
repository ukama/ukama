//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2023-present, Ukama Inc.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.3
// source: accounting.proto

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

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounting_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_accounting_proto_msgTypes[0]
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
	return file_accounting_proto_rawDescGZIP(), []int{0}
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

	Accounting *Accounting `protobuf:"bytes,1,opt,name=accounting,proto3" json:"accounting,omitempty"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounting_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_accounting_proto_msgTypes[1]
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
	return file_accounting_proto_rawDescGZIP(), []int{1}
}

func (x *GetResponse) GetAccounting() *Accounting {
	if x != nil {
		return x.Accounting
	}
	return nil
}

type GetByUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=userId,json=user_id,proto3" json:"userId,omitempty"`
}

func (x *GetByUserRequest) Reset() {
	*x = GetByUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounting_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByUserRequest) ProtoMessage() {}

func (x *GetByUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_accounting_proto_msgTypes[2]
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
	return file_accounting_proto_rawDescGZIP(), []int{2}
}

func (x *GetByUserRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GetByUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Accounting []*Accounting `protobuf:"bytes,1,rep,name=accounting,proto3" json:"accounting,omitempty"`
}

func (x *GetByUserResponse) Reset() {
	*x = GetByUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounting_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByUserResponse) ProtoMessage() {}

func (x *GetByUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_accounting_proto_msgTypes[3]
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
	return file_accounting_proto_rawDescGZIP(), []int{3}
}

func (x *GetByUserResponse) GetAccounting() []*Accounting {
	if x != nil {
		return x.Accounting
	}
	return nil
}

type SyncAcountingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SyncAcountingRequest) Reset() {
	*x = SyncAcountingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounting_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncAcountingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncAcountingRequest) ProtoMessage() {}

func (x *SyncAcountingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_accounting_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncAcountingRequest.ProtoReflect.Descriptor instead.
func (*SyncAcountingRequest) Descriptor() ([]byte, []int) {
	return file_accounting_proto_rawDescGZIP(), []int{4}
}

type SyncAcountingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SyncAcountingResponse) Reset() {
	*x = SyncAcountingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounting_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncAcountingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncAcountingResponse) ProtoMessage() {}

func (x *SyncAcountingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_accounting_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncAcountingResponse.ProtoReflect.Descriptor instead.
func (*SyncAcountingResponse) Descriptor() ([]byte, []int) {
	return file_accounting_proto_rawDescGZIP(), []int{5}
}

type Accounting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        string `protobuf:"bytes,2,opt,name=userId,json=user_id,proto3" json:"userId,omitempty"`
	Item          string `protobuf:"bytes,3,opt,name=item,proto3" json:"item,omitempty"`
	Description   string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Inventory     string `protobuf:"bytes,5,opt,name=inventory,proto3" json:"inventory,omitempty"`
	OpexFee       string `protobuf:"bytes,6,opt,name=opexFee,json=opex_fee,proto3" json:"opexFee,omitempty"`
	Vat           string `protobuf:"bytes,7,opt,name=vat,proto3" json:"vat,omitempty"`
	EffectiveDate string `protobuf:"bytes,8,opt,name=effectiveDate,json=effective_date,proto3" json:"effectiveDate,omitempty"`
}

func (x *Accounting) Reset() {
	*x = Accounting{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounting_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Accounting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Accounting) ProtoMessage() {}

func (x *Accounting) ProtoReflect() protoreflect.Message {
	mi := &file_accounting_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Accounting.ProtoReflect.Descriptor instead.
func (*Accounting) Descriptor() ([]byte, []int) {
	return file_accounting_proto_rawDescGZIP(), []int{6}
}

func (x *Accounting) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Accounting) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Accounting) GetItem() string {
	if x != nil {
		return x.Item
	}
	return ""
}

func (x *Accounting) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Accounting) GetInventory() string {
	if x != nil {
		return x.Inventory
	}
	return ""
}

func (x *Accounting) GetOpexFee() string {
	if x != nil {
		return x.OpexFee
	}
	return ""
}

func (x *Accounting) GetVat() string {
	if x != nil {
		return x.Vat
	}
	return ""
}

func (x *Accounting) GetEffectiveDate() string {
	if x != nil {
		return x.EffectiveDate
	}
	return ""
}

var File_accounting_proto protoreflect.FileDescriptor

var file_accounting_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x1d, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74,
	0x6f, 0x72, 0x79, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x76,
	0x31, 0x22, 0x1c, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x58, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x49,
	0x0a, 0x0a, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x29, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e,
	0x74, 0x6f, 0x72, 0x79, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x0a, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x22, 0x2b, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x22, 0x5e, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x42, 0x79, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x49, 0x0a, 0x0a, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x29, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72,
	0x79, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x0a, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x22, 0x16, 0x0a, 0x14, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x17,
	0x0a, 0x15, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xdd, 0x01, 0x0a, 0x0a, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69,
	0x74, 0x65, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f,
	0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74,
	0x6f, 0x72, 0x79, 0x12, 0x19, 0x0a, 0x07, 0x6f, 0x70, 0x65, 0x78, 0x46, 0x65, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x78, 0x5f, 0x66, 0x65, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x76, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x76, 0x61, 0x74,
	0x12, 0x25, 0x0a, 0x0d, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x44, 0x61, 0x74,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x69,
	0x76, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x32, 0xde, 0x02, 0x0a, 0x11, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5c, 0x0a,
	0x03, 0x47, 0x65, 0x74, 0x12, 0x29, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76,
	0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e,
	0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2a, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72,
	0x79, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6e, 0x0a, 0x09, 0x47,
	0x65, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x12, 0x2f, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61,
	0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e, 0x75, 0x6b, 0x61, 0x6d,
	0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x7b, 0x0a, 0x0e, 0x53,
	0x79, 0x6e, 0x63, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x33, 0x2e,
	0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x79,
	0x6e, 0x63, 0x41, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x34, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e,
	0x74, 0x6f, 0x72, 0x79, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2f, 0x75, 0x6b, 0x61,
	0x6d, 0x61, 0x2f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x69, 0x6e, 0x76, 0x65, 0x6e,
	0x74, 0x6f, 0x72, 0x79, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x2f,
	0x70, 0x62, 0x2f, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_accounting_proto_rawDescOnce sync.Once
	file_accounting_proto_rawDescData = file_accounting_proto_rawDesc
)

func file_accounting_proto_rawDescGZIP() []byte {
	file_accounting_proto_rawDescOnce.Do(func() {
		file_accounting_proto_rawDescData = protoimpl.X.CompressGZIP(file_accounting_proto_rawDescData)
	})
	return file_accounting_proto_rawDescData
}

var file_accounting_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_accounting_proto_goTypes = []interface{}{
	(*GetRequest)(nil),            // 0: ukama.inventory.accounting.v1.GetRequest
	(*GetResponse)(nil),           // 1: ukama.inventory.accounting.v1.GetResponse
	(*GetByUserRequest)(nil),      // 2: ukama.inventory.accounting.v1.GetByUserRequest
	(*GetByUserResponse)(nil),     // 3: ukama.inventory.accounting.v1.GetByUserResponse
	(*SyncAcountingRequest)(nil),  // 4: ukama.inventory.accounting.v1.SyncAcountingRequest
	(*SyncAcountingResponse)(nil), // 5: ukama.inventory.accounting.v1.SyncAcountingResponse
	(*Accounting)(nil),            // 6: ukama.inventory.accounting.v1.Accounting
}
var file_accounting_proto_depIdxs = []int32{
	6, // 0: ukama.inventory.accounting.v1.GetResponse.accounting:type_name -> ukama.inventory.accounting.v1.Accounting
	6, // 1: ukama.inventory.accounting.v1.GetByUserResponse.accounting:type_name -> ukama.inventory.accounting.v1.Accounting
	0, // 2: ukama.inventory.accounting.v1.AccountingService.Get:input_type -> ukama.inventory.accounting.v1.GetRequest
	2, // 3: ukama.inventory.accounting.v1.AccountingService.GetByUser:input_type -> ukama.inventory.accounting.v1.GetByUserRequest
	4, // 4: ukama.inventory.accounting.v1.AccountingService.SyncAccounting:input_type -> ukama.inventory.accounting.v1.SyncAcountingRequest
	1, // 5: ukama.inventory.accounting.v1.AccountingService.Get:output_type -> ukama.inventory.accounting.v1.GetResponse
	3, // 6: ukama.inventory.accounting.v1.AccountingService.GetByUser:output_type -> ukama.inventory.accounting.v1.GetByUserResponse
	5, // 7: ukama.inventory.accounting.v1.AccountingService.SyncAccounting:output_type -> ukama.inventory.accounting.v1.SyncAcountingResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_accounting_proto_init() }
func file_accounting_proto_init() {
	if File_accounting_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_accounting_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_accounting_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_accounting_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_accounting_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_accounting_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncAcountingRequest); i {
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
		file_accounting_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncAcountingResponse); i {
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
		file_accounting_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Accounting); i {
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
			RawDescriptor: file_accounting_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_accounting_proto_goTypes,
		DependencyIndexes: file_accounting_proto_depIdxs,
		MessageInfos:      file_accounting_proto_msgTypes,
	}.Build()
	File_accounting_proto = out.File
	file_accounting_proto_rawDesc = nil
	file_accounting_proto_goTypes = nil
	file_accounting_proto_depIdxs = nil
}
