//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2023-present, Ukama Inc.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: dsubscriber.proto

package gen

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UpdateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Dsubscriber   *Dsubscriber           `protobuf:"bytes,1,opt,name=dsubscriber,proto3" json:"dsubscriber,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateRequest) Reset() {
	*x = UpdateRequest{}
	mi := &file_dsubscriber_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRequest) ProtoMessage() {}

func (x *UpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dsubscriber_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRequest.ProtoReflect.Descriptor instead.
func (*UpdateRequest) Descriptor() ([]byte, []int) {
	return file_dsubscriber_proto_rawDescGZIP(), []int{0}
}

func (x *UpdateRequest) GetDsubscriber() *Dsubscriber {
	if x != nil {
		return x.Dsubscriber
	}
	return nil
}

type UpdateResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Dsubscriber   *Dsubscriber           `protobuf:"bytes,1,opt,name=dsubscriber,proto3" json:"dsubscriber,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateResponse) Reset() {
	*x = UpdateResponse{}
	mi := &file_dsubscriber_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateResponse) ProtoMessage() {}

func (x *UpdateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dsubscriber_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateResponse.ProtoReflect.Descriptor instead.
func (*UpdateResponse) Descriptor() ([]byte, []int) {
	return file_dsubscriber_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateResponse) GetDsubscriber() *Dsubscriber {
	if x != nil {
		return x.Dsubscriber
	}
	return nil
}

type Dsubscriber struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Iccid         string                 `protobuf:"bytes,1,opt,name=iccid,proto3" json:"iccid,omitempty"`
	Profile       string                 `protobuf:"bytes,2,opt,name=profile,proto3" json:"profile,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Dsubscriber) Reset() {
	*x = Dsubscriber{}
	mi := &file_dsubscriber_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Dsubscriber) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Dsubscriber) ProtoMessage() {}

func (x *Dsubscriber) ProtoReflect() protoreflect.Message {
	mi := &file_dsubscriber_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Dsubscriber.ProtoReflect.Descriptor instead.
func (*Dsubscriber) Descriptor() ([]byte, []int) {
	return file_dsubscriber_proto_rawDescGZIP(), []int{2}
}

func (x *Dsubscriber) GetIccid() string {
	if x != nil {
		return x.Iccid
	}
	return ""
}

func (x *Dsubscriber) GetProfile() string {
	if x != nil {
		return x.Profile
	}
	return ""
}

var File_dsubscriber_proto protoreflect.FileDescriptor

var file_dsubscriber_proto_rawDesc = string([]byte{
	0x0a, 0x11, 0x64, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x64, 0x75, 0x6d, 0x6d, 0x79,
	0x2e, 0x64, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x22,
	0x5a, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x49, 0x0a, 0x0b, 0x64, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x64, 0x75,
	0x6d, 0x6d, 0x79, 0x2e, 0x64, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x44, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x52, 0x0b,
	0x64, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x22, 0x5b, 0x0a, 0x0e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x49, 0x0a,
	0x0b, 0x64, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x27, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x64, 0x75, 0x6d, 0x6d, 0x79,
	0x2e, 0x64, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x44, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x52, 0x0b, 0x64, 0x73, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x22, 0x3d, 0x0a, 0x0b, 0x44, 0x73, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x63, 0x63, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x63, 0x63, 0x69, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x32, 0x75, 0x0a, 0x12, 0x44, 0x73, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5f, 0x0a,
	0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x29, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e,
	0x64, 0x75, 0x6d, 0x6d, 0x79, 0x2e, 0x64, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x64, 0x75, 0x6d, 0x6d, 0x79,
	0x2e, 0x64, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x42,
	0x5a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x6b, 0x61,
	0x6d, 0x61, 0x2f, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x75, 0x6d, 0x6d, 0x79, 0x2f,
	0x64, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x2f, 0x70, 0x62, 0x2f, 0x67,
	0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_dsubscriber_proto_rawDescOnce sync.Once
	file_dsubscriber_proto_rawDescData []byte
)

func file_dsubscriber_proto_rawDescGZIP() []byte {
	file_dsubscriber_proto_rawDescOnce.Do(func() {
		file_dsubscriber_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_dsubscriber_proto_rawDesc), len(file_dsubscriber_proto_rawDesc)))
	})
	return file_dsubscriber_proto_rawDescData
}

var file_dsubscriber_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_dsubscriber_proto_goTypes = []any{
	(*UpdateRequest)(nil),  // 0: ukama.dummy.dsubscriber.v1.UpdateRequest
	(*UpdateResponse)(nil), // 1: ukama.dummy.dsubscriber.v1.UpdateResponse
	(*Dsubscriber)(nil),    // 2: ukama.dummy.dsubscriber.v1.Dsubscriber
}
var file_dsubscriber_proto_depIdxs = []int32{
	2, // 0: ukama.dummy.dsubscriber.v1.UpdateRequest.dsubscriber:type_name -> ukama.dummy.dsubscriber.v1.Dsubscriber
	2, // 1: ukama.dummy.dsubscriber.v1.UpdateResponse.dsubscriber:type_name -> ukama.dummy.dsubscriber.v1.Dsubscriber
	0, // 2: ukama.dummy.dsubscriber.v1.DsubscriberService.Update:input_type -> ukama.dummy.dsubscriber.v1.UpdateRequest
	1, // 3: ukama.dummy.dsubscriber.v1.DsubscriberService.Update:output_type -> ukama.dummy.dsubscriber.v1.UpdateResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_dsubscriber_proto_init() }
func file_dsubscriber_proto_init() {
	if File_dsubscriber_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_dsubscriber_proto_rawDesc), len(file_dsubscriber_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_dsubscriber_proto_goTypes,
		DependencyIndexes: file_dsubscriber_proto_depIdxs,
		MessageInfos:      file_dsubscriber_proto_msgTypes,
	}.Build()
	File_dsubscriber_proto = out.File
	file_dsubscriber_proto_goTypes = nil
	file_dsubscriber_proto_depIdxs = nil
}
