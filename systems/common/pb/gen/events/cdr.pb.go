// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.3
// source: events/cdr.proto

package events

import (
	_ "github.com/mwitkow/go-proto-validators"
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

type NodeChanged struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Imsi              string `protobuf:"bytes,1,opt,name=Imsi,json=imsi,proto3" json:"Imsi,omitempty"`
	Policy            string `protobuf:"bytes,2,opt,name=Policy,json=policy,proto3" json:"Policy,omitempty"`
	TotalUsage        uint64 `protobuf:"varint,3,opt,name=TotalUsage,proto3" json:"TotalUsage,omitempty"`
	UsageTillLastNode uint64 `protobuf:"varint,4,opt,name=UsageTillLastNode,proto3" json:"UsageTillLastNode,omitempty"`
	NodeId            string `protobuf:"bytes,5,opt,name=NodeId,proto3" json:"NodeId,omitempty"`
	OldNodeId         string `protobuf:"bytes,6,opt,name=OldNodeId,proto3" json:"OldNodeId,omitempty"`
}

func (x *NodeChanged) Reset() {
	*x = NodeChanged{}
	mi := &file_events_cdr_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NodeChanged) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeChanged) ProtoMessage() {}

func (x *NodeChanged) ProtoReflect() protoreflect.Message {
	mi := &file_events_cdr_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeChanged.ProtoReflect.Descriptor instead.
func (*NodeChanged) Descriptor() ([]byte, []int) {
	return file_events_cdr_proto_rawDescGZIP(), []int{0}
}

func (x *NodeChanged) GetImsi() string {
	if x != nil {
		return x.Imsi
	}
	return ""
}

func (x *NodeChanged) GetPolicy() string {
	if x != nil {
		return x.Policy
	}
	return ""
}

func (x *NodeChanged) GetTotalUsage() uint64 {
	if x != nil {
		return x.TotalUsage
	}
	return 0
}

func (x *NodeChanged) GetUsageTillLastNode() uint64 {
	if x != nil {
		return x.UsageTillLastNode
	}
	return 0
}

func (x *NodeChanged) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *NodeChanged) GetOldNodeId() string {
	if x != nil {
		return x.OldNodeId
	}
	return ""
}

type SessionCreated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Imsi      string `protobuf:"bytes,1,opt,name=Imsi,json=imsi,proto3" json:"Imsi,omitempty"`
	Policy    string `protobuf:"bytes,2,opt,name=Policy,json=policy,proto3" json:"Policy,omitempty"`
	Usage     uint64 `protobuf:"varint,3,opt,name=Usage,proto3" json:"Usage,omitempty"`
	StartTime uint64 `protobuf:"varint,4,opt,name=StartTime,proto3" json:"StartTime,omitempty"`
	NodeId    string `protobuf:"bytes,5,opt,name=NodeId,proto3" json:"NodeId,omitempty"`
	SessionId uint64 `protobuf:"varint,6,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
}

func (x *SessionCreated) Reset() {
	*x = SessionCreated{}
	mi := &file_events_cdr_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SessionCreated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionCreated) ProtoMessage() {}

func (x *SessionCreated) ProtoReflect() protoreflect.Message {
	mi := &file_events_cdr_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionCreated.ProtoReflect.Descriptor instead.
func (*SessionCreated) Descriptor() ([]byte, []int) {
	return file_events_cdr_proto_rawDescGZIP(), []int{1}
}

func (x *SessionCreated) GetImsi() string {
	if x != nil {
		return x.Imsi
	}
	return ""
}

func (x *SessionCreated) GetPolicy() string {
	if x != nil {
		return x.Policy
	}
	return ""
}

func (x *SessionCreated) GetUsage() uint64 {
	if x != nil {
		return x.Usage
	}
	return 0
}

func (x *SessionCreated) GetStartTime() uint64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *SessionCreated) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *SessionCreated) GetSessionId() uint64 {
	if x != nil {
		return x.SessionId
	}
	return 0
}

type CDRReported struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Session       uint64 `protobuf:"varint,1,opt,name=Session,json=session_id,proto3" json:"Session,omitempty"`
	NodeId        string `protobuf:"bytes,2,opt,name=NodeId,json=node_id,proto3" json:"NodeId,omitempty"`
	Imsi          string `protobuf:"bytes,3,opt,name=Imsi,json=imsi,proto3" json:"Imsi,omitempty"`
	Policy        string `protobuf:"bytes,4,opt,name=Policy,json=policy,proto3" json:"Policy,omitempty"`
	ApnName       string `protobuf:"bytes,5,opt,name=ApnName,json=apn_name,proto3" json:"ApnName,omitempty"`
	Ip            string `protobuf:"bytes,6,opt,name=Ip,json=ue_ip,proto3" json:"Ip,omitempty"`
	StartTime     uint64 `protobuf:"varint,7,opt,name=StartTime,json=start_time,proto3" json:"StartTime,omitempty"`
	EndTime       uint64 `protobuf:"varint,8,opt,name=EndTime,json=end_time,proto3" json:"EndTime,omitempty"`
	LastUpdatedAt uint64 `protobuf:"varint,9,opt,name=LastUpdatedAt,json=last_updated_at,proto3" json:"LastUpdatedAt,omitempty"`
	TxBytes       uint64 `protobuf:"varint,10,opt,name=TxBytes,json=tx_bytes,proto3" json:"TxBytes,omitempty"`
	RxBytes       uint64 `protobuf:"varint,11,opt,name=RxBytes,json=rx_bytes,proto3" json:"RxBytes,omitempty"`
	TotalBytes    uint64 `protobuf:"varint,12,opt,name=TotalBytes,json=total_bytes,proto3" json:"TotalBytes,omitempty"`
}

func (x *CDRReported) Reset() {
	*x = CDRReported{}
	mi := &file_events_cdr_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CDRReported) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CDRReported) ProtoMessage() {}

func (x *CDRReported) ProtoReflect() protoreflect.Message {
	mi := &file_events_cdr_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CDRReported.ProtoReflect.Descriptor instead.
func (*CDRReported) Descriptor() ([]byte, []int) {
	return file_events_cdr_proto_rawDescGZIP(), []int{2}
}

func (x *CDRReported) GetSession() uint64 {
	if x != nil {
		return x.Session
	}
	return 0
}

func (x *CDRReported) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *CDRReported) GetImsi() string {
	if x != nil {
		return x.Imsi
	}
	return ""
}

func (x *CDRReported) GetPolicy() string {
	if x != nil {
		return x.Policy
	}
	return ""
}

func (x *CDRReported) GetApnName() string {
	if x != nil {
		return x.ApnName
	}
	return ""
}

func (x *CDRReported) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *CDRReported) GetStartTime() uint64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *CDRReported) GetEndTime() uint64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

func (x *CDRReported) GetLastUpdatedAt() uint64 {
	if x != nil {
		return x.LastUpdatedAt
	}
	return 0
}

func (x *CDRReported) GetTxBytes() uint64 {
	if x != nil {
		return x.TxBytes
	}
	return 0
}

func (x *CDRReported) GetRxBytes() uint64 {
	if x != nil {
		return x.RxBytes
	}
	return 0
}

func (x *CDRReported) GetTotalBytes() uint64 {
	if x != nil {
		return x.TotalBytes
	}
	return 0
}

type SessionDestroyed struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Imsi         string `protobuf:"bytes,1,opt,name=Imsi,json=imsi,proto3" json:"Imsi,omitempty"`
	Policy       string `protobuf:"bytes,2,opt,name=Policy,json=policy,proto3" json:"Policy,omitempty"`
	Usage        uint32 `protobuf:"varint,3,opt,name=Usage,proto3" json:"Usage,omitempty"`
	NodeId       string `protobuf:"bytes,6,opt,name=NodeId,proto3" json:"NodeId,omitempty"`
	SessionId    uint64 `protobuf:"varint,7,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	SessionUsage uint64 `protobuf:"varint,8,opt,name=SessionUsage,proto3" json:"SessionUsage,omitempty"`
	TotalUsage   uint64 `protobuf:"varint,9,opt,name=TotalUsage,proto3" json:"TotalUsage,omitempty"`
}

func (x *SessionDestroyed) Reset() {
	*x = SessionDestroyed{}
	mi := &file_events_cdr_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SessionDestroyed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionDestroyed) ProtoMessage() {}

func (x *SessionDestroyed) ProtoReflect() protoreflect.Message {
	mi := &file_events_cdr_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionDestroyed.ProtoReflect.Descriptor instead.
func (*SessionDestroyed) Descriptor() ([]byte, []int) {
	return file_events_cdr_proto_rawDescGZIP(), []int{3}
}

func (x *SessionDestroyed) GetImsi() string {
	if x != nil {
		return x.Imsi
	}
	return ""
}

func (x *SessionDestroyed) GetPolicy() string {
	if x != nil {
		return x.Policy
	}
	return ""
}

func (x *SessionDestroyed) GetUsage() uint32 {
	if x != nil {
		return x.Usage
	}
	return 0
}

func (x *SessionDestroyed) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *SessionDestroyed) GetSessionId() uint64 {
	if x != nil {
		return x.SessionId
	}
	return 0
}

func (x *SessionDestroyed) GetSessionUsage() uint64 {
	if x != nil {
		return x.SessionUsage
	}
	return 0
}

func (x *SessionDestroyed) GetTotalUsage() uint64 {
	if x != nil {
		return x.TotalUsage
	}
	return 0
}

var File_events_cdr_proto protoreflect.FileDescriptor

var file_events_cdr_proto_rawDesc = []byte{
	0x0a, 0x10, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x63, 0x64, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0f, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x2e, 0x76, 0x31, 0x1a, 0x0f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd4, 0x01, 0x0a, 0x0b, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x64, 0x12, 0x1e, 0x0a, 0x04, 0x49, 0x6d, 0x73, 0x69, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x0a, 0xe2, 0xdf, 0x1f, 0x06, 0x58, 0x01, 0x70, 0x05, 0x78, 0x10, 0x52, 0x04,
	0x69, 0x6d, 0x73, 0x69, 0x12, 0x21, 0x0a, 0x06, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xe2, 0xdf, 0x1f, 0x05, 0x58, 0x01, 0x90, 0x01, 0x04, 0x52,
	0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x6f, 0x74, 0x61, 0x6c,
	0x55, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x54, 0x6f, 0x74,
	0x61, 0x6c, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2c, 0x0a, 0x11, 0x55, 0x73, 0x61, 0x67, 0x65,
	0x54, 0x69, 0x6c, 0x6c, 0x4c, 0x61, 0x73, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x11, 0x55, 0x73, 0x61, 0x67, 0x65, 0x54, 0x69, 0x6c, 0x6c, 0x4c, 0x61, 0x73,
	0x74, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x4f, 0x6c, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x4f, 0x6c, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0xbd, 0x01, 0x0a, 0x0e,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1e,
	0x0a, 0x04, 0x49, 0x6d, 0x73, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe2, 0xdf,
	0x1f, 0x06, 0x58, 0x01, 0x70, 0x05, 0x78, 0x10, 0x52, 0x04, 0x69, 0x6d, 0x73, 0x69, 0x12, 0x21,
	0x0a, 0x06, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09,
	0xe2, 0xdf, 0x1f, 0x05, 0x58, 0x01, 0x90, 0x01, 0x04, 0x52, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x55, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x05, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x72, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0xe2, 0x02, 0x0a, 0x0b,
	0x43, 0x44, 0x52, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x12, 0x1b, 0x0a, 0x07, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69,
	0x64, 0x12, 0x1e, 0x0a, 0x04, 0x49, 0x6d, 0x73, 0x69, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x0a, 0xe2, 0xdf, 0x1f, 0x06, 0x58, 0x01, 0x70, 0x05, 0x78, 0x10, 0x52, 0x04, 0x69, 0x6d, 0x73,
	0x69, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x19, 0x0a, 0x07, 0x41, 0x70, 0x6e,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x61, 0x70, 0x6e, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x11, 0x0a, 0x02, 0x49, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x75, 0x65, 0x5f, 0x69, 0x70, 0x12, 0x1d, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x72, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x07, 0x45, 0x6e, 0x64, 0x54, 0x69, 0x6d,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x12, 0x26, 0x0a, 0x0d, 0x4c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x12, 0x19, 0x0a, 0x07, 0x54, 0x78, 0x42,
	0x79, 0x74, 0x65, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x74, 0x78, 0x5f, 0x62,
	0x79, 0x74, 0x65, 0x73, 0x12, 0x19, 0x0a, 0x07, 0x52, 0x78, 0x42, 0x79, 0x74, 0x65, 0x73, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x72, 0x78, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x12,
	0x1f, 0x0a, 0x0a, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x42, 0x79, 0x74, 0x65, 0x73, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73,
	0x22, 0xe5, 0x01, 0x0a, 0x10, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x73, 0x74,
	0x72, 0x6f, 0x79, 0x65, 0x64, 0x12, 0x1e, 0x0a, 0x04, 0x49, 0x6d, 0x73, 0x69, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x0a, 0xe2, 0xdf, 0x1f, 0x06, 0x58, 0x01, 0x70, 0x05, 0x78, 0x10, 0x52,
	0x04, 0x69, 0x6d, 0x73, 0x69, 0x12, 0x21, 0x0a, 0x06, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xe2, 0xdf, 0x1f, 0x05, 0x58, 0x01, 0x90, 0x01, 0x04,
	0x52, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x55, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x49, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x53, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x55,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x6f, 0x74, 0x61,
	0x6c, 0x55, 0x73, 0x61, 0x67, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x54, 0x6f,
	0x74, 0x61, 0x6c, 0x55, 0x73, 0x61, 0x67, 0x65, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2f, 0x75, 0x6b, 0x61,
	0x6d, 0x61, 0x2f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x70, 0x62, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_events_cdr_proto_rawDescOnce sync.Once
	file_events_cdr_proto_rawDescData = file_events_cdr_proto_rawDesc
)

func file_events_cdr_proto_rawDescGZIP() []byte {
	file_events_cdr_proto_rawDescOnce.Do(func() {
		file_events_cdr_proto_rawDescData = protoimpl.X.CompressGZIP(file_events_cdr_proto_rawDescData)
	})
	return file_events_cdr_proto_rawDescData
}

var file_events_cdr_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_events_cdr_proto_goTypes = []any{
	(*NodeChanged)(nil),      // 0: ukama.events.v1.NodeChanged
	(*SessionCreated)(nil),   // 1: ukama.events.v1.SessionCreated
	(*CDRReported)(nil),      // 2: ukama.events.v1.CDRReported
	(*SessionDestroyed)(nil), // 3: ukama.events.v1.SessionDestroyed
}
var file_events_cdr_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_events_cdr_proto_init() }
func file_events_cdr_proto_init() {
	if File_events_cdr_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_events_cdr_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_events_cdr_proto_goTypes,
		DependencyIndexes: file_events_cdr_proto_depIdxs,
		MessageInfos:      file_events_cdr_proto_msgTypes,
	}.Build()
	File_events_cdr_proto = out.File
	file_events_cdr_proto_rawDesc = nil
	file_events_cdr_proto_goTypes = nil
	file_events_cdr_proto_depIdxs = nil
}
