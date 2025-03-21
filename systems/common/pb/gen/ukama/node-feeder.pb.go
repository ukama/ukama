// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.3
// source: ukama/node-feeder.proto

package ukama

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

type NodeFeederMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Target     string `protobuf:"bytes,1,opt,name=Target,proto3" json:"Target,omitempty"`
	HTTPMethod string `protobuf:"bytes,2,opt,name=HTTPMethod,proto3" json:"HTTPMethod,omitempty"`
	Path       string `protobuf:"bytes,3,opt,name=Path,proto3" json:"Path,omitempty"`
	Msg        []byte `protobuf:"bytes,4,opt,name=msg,proto3" json:"msg,omitempty"` // Use bytes for binary data
}

func (x *NodeFeederMessage) Reset() {
	*x = NodeFeederMessage{}
	mi := &file_ukama_node_feeder_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NodeFeederMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeFeederMessage) ProtoMessage() {}

func (x *NodeFeederMessage) ProtoReflect() protoreflect.Message {
	mi := &file_ukama_node_feeder_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeFeederMessage.ProtoReflect.Descriptor instead.
func (*NodeFeederMessage) Descriptor() ([]byte, []int) {
	return file_ukama_node_feeder_proto_rawDescGZIP(), []int{0}
}

func (x *NodeFeederMessage) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *NodeFeederMessage) GetHTTPMethod() string {
	if x != nil {
		return x.HTTPMethod
	}
	return ""
}

func (x *NodeFeederMessage) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *NodeFeederMessage) GetMsg() []byte {
	if x != nil {
		return x.Msg
	}
	return nil
}

var File_ukama_node_feeder_proto protoreflect.FileDescriptor

var file_ukama_node_feeder_proto_rawDesc = []byte{
	0x0a, 0x17, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2d, 0x66, 0x65, 0x65,
	0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x75, 0x6b, 0x61, 0x6d, 0x61,
	0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x66, 0x65, 0x65, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x22, 0x71,
	0x0a, 0x11, 0x4e, 0x6f, 0x64, 0x65, 0x46, 0x65, 0x65, 0x64, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x48,
	0x54, 0x54, 0x50, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x48, 0x54, 0x54, 0x50, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x50,
	0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x50, 0x61, 0x74, 0x68, 0x12,
	0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6d, 0x73,
	0x67, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2f, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2f, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x62, 0x2f, 0x67, 0x65,
	0x6e, 0x2f, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ukama_node_feeder_proto_rawDescOnce sync.Once
	file_ukama_node_feeder_proto_rawDescData = file_ukama_node_feeder_proto_rawDesc
)

func file_ukama_node_feeder_proto_rawDescGZIP() []byte {
	file_ukama_node_feeder_proto_rawDescOnce.Do(func() {
		file_ukama_node_feeder_proto_rawDescData = protoimpl.X.CompressGZIP(file_ukama_node_feeder_proto_rawDescData)
	})
	return file_ukama_node_feeder_proto_rawDescData
}

var file_ukama_node_feeder_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_ukama_node_feeder_proto_goTypes = []any{
	(*NodeFeederMessage)(nil), // 0: ukama.nodefeeder.v1.NodeFeederMessage
}
var file_ukama_node_feeder_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ukama_node_feeder_proto_init() }
func file_ukama_node_feeder_proto_init() {
	if File_ukama_node_feeder_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ukama_node_feeder_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ukama_node_feeder_proto_goTypes,
		DependencyIndexes: file_ukama_node_feeder_proto_depIdxs,
		MessageInfos:      file_ukama_node_feeder_proto_msgTypes,
	}.Build()
	File_ukama_node_feeder_proto = out.File
	file_ukama_node_feeder_proto_rawDesc = nil
	file_ukama_node_feeder_proto_goTypes = nil
	file_ukama_node_feeder_proto_depIdxs = nil
}
