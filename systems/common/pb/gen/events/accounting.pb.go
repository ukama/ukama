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
// source: events/accounting.proto

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

type UserAccountingEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId     string            `protobuf:"bytes,1,opt,name=userId,json=user_id,proto3" json:"userId,omitempty"`
	Accounting []*UserAccounting `protobuf:"bytes,2,rep,name=accounting,json=user_accounting,proto3" json:"accounting,omitempty"`
}

func (x *UserAccountingEvent) Reset() {
	*x = UserAccountingEvent{}
	mi := &file_events_accounting_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserAccountingEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserAccountingEvent) ProtoMessage() {}

func (x *UserAccountingEvent) ProtoReflect() protoreflect.Message {
	mi := &file_events_accounting_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserAccountingEvent.ProtoReflect.Descriptor instead.
func (*UserAccountingEvent) Descriptor() ([]byte, []int) {
	return file_events_accounting_proto_rawDescGZIP(), []int{0}
}

func (x *UserAccountingEvent) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserAccountingEvent) GetAccounting() []*UserAccounting {
	if x != nil {
		return x.Accounting
	}
	return nil
}

type UserAccounting struct {
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

func (x *UserAccounting) Reset() {
	*x = UserAccounting{}
	mi := &file_events_accounting_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserAccounting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserAccounting) ProtoMessage() {}

func (x *UserAccounting) ProtoReflect() protoreflect.Message {
	mi := &file_events_accounting_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserAccounting.ProtoReflect.Descriptor instead.
func (*UserAccounting) Descriptor() ([]byte, []int) {
	return file_events_accounting_proto_rawDescGZIP(), []int{1}
}

func (x *UserAccounting) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UserAccounting) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserAccounting) GetItem() string {
	if x != nil {
		return x.Item
	}
	return ""
}

func (x *UserAccounting) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *UserAccounting) GetInventory() string {
	if x != nil {
		return x.Inventory
	}
	return ""
}

func (x *UserAccounting) GetOpexFee() string {
	if x != nil {
		return x.OpexFee
	}
	return ""
}

func (x *UserAccounting) GetVat() string {
	if x != nil {
		return x.Vat
	}
	return ""
}

func (x *UserAccounting) GetEffectiveDate() string {
	if x != nil {
		return x.EffectiveDate
	}
	return ""
}

var File_events_accounting_proto protoreflect.FileDescriptor

var file_events_accounting_proto_rawDesc = []byte{
	0x0a, 0x17, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x75, 0x6b, 0x61, 0x6d, 0x61,
	0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x22, 0x74, 0x0a, 0x13, 0x55, 0x73,
	0x65, 0x72, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x12, 0x17, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x12, 0x44, 0x0a, 0x0a, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f,
	0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67, 0x52,
	0x0f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x69, 0x6e, 0x67,
	0x22, 0xe1, 0x01, 0x0a, 0x0e, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x69, 0x6e, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x69, 0x74, 0x65, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x74, 0x65, 0x6d,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79,
	0x12, 0x19, 0x0a, 0x07, 0x6f, 0x70, 0x65, 0x78, 0x46, 0x65, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x78, 0x5f, 0x66, 0x65, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x76,
	0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x76, 0x61, 0x74, 0x12, 0x25, 0x0a,
	0x0d, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x44, 0x61, 0x74, 0x65, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f,
	0x64, 0x61, 0x74, 0x65, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2f, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2f, 0x73,
	0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x62,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_events_accounting_proto_rawDescOnce sync.Once
	file_events_accounting_proto_rawDescData = file_events_accounting_proto_rawDesc
)

func file_events_accounting_proto_rawDescGZIP() []byte {
	file_events_accounting_proto_rawDescOnce.Do(func() {
		file_events_accounting_proto_rawDescData = protoimpl.X.CompressGZIP(file_events_accounting_proto_rawDescData)
	})
	return file_events_accounting_proto_rawDescData
}

var file_events_accounting_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_events_accounting_proto_goTypes = []any{
	(*UserAccountingEvent)(nil), // 0: ukama.events.v1.UserAccountingEvent
	(*UserAccounting)(nil),      // 1: ukama.events.v1.UserAccounting
}
var file_events_accounting_proto_depIdxs = []int32{
	1, // 0: ukama.events.v1.UserAccountingEvent.accounting:type_name -> ukama.events.v1.UserAccounting
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_events_accounting_proto_init() }
func file_events_accounting_proto_init() {
	if File_events_accounting_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_events_accounting_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_events_accounting_proto_goTypes,
		DependencyIndexes: file_events_accounting_proto_depIdxs,
		MessageInfos:      file_events_accounting_proto_msgTypes,
	}.Build()
	File_events_accounting_proto = out.File
	file_events_accounting_proto_rawDesc = nil
	file_events_accounting_proto_goTypes = nil
	file_events_accounting_proto_depIdxs = nil
}
