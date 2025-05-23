// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.28.3
// source: cdr.proto

package gen

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

type CDR struct {
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

func (x *CDR) Reset() {
	*x = CDR{}
	mi := &file_cdr_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CDR) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CDR) ProtoMessage() {}

func (x *CDR) ProtoReflect() protoreflect.Message {
	mi := &file_cdr_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CDR.ProtoReflect.Descriptor instead.
func (*CDR) Descriptor() ([]byte, []int) {
	return file_cdr_proto_rawDescGZIP(), []int{0}
}

func (x *CDR) GetSession() uint64 {
	if x != nil {
		return x.Session
	}
	return 0
}

func (x *CDR) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *CDR) GetImsi() string {
	if x != nil {
		return x.Imsi
	}
	return ""
}

func (x *CDR) GetPolicy() string {
	if x != nil {
		return x.Policy
	}
	return ""
}

func (x *CDR) GetApnName() string {
	if x != nil {
		return x.ApnName
	}
	return ""
}

func (x *CDR) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *CDR) GetStartTime() uint64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *CDR) GetEndTime() uint64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

func (x *CDR) GetLastUpdatedAt() uint64 {
	if x != nil {
		return x.LastUpdatedAt
	}
	return 0
}

func (x *CDR) GetTxBytes() uint64 {
	if x != nil {
		return x.TxBytes
	}
	return 0
}

func (x *CDR) GetRxBytes() uint64 {
	if x != nil {
		return x.RxBytes
	}
	return 0
}

func (x *CDR) GetTotalBytes() uint64 {
	if x != nil {
		return x.TotalBytes
	}
	return 0
}

type CDRResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CDRResp) Reset() {
	*x = CDRResp{}
	mi := &file_cdr_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CDRResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CDRResp) ProtoMessage() {}

func (x *CDRResp) ProtoReflect() protoreflect.Message {
	mi := &file_cdr_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CDRResp.ProtoReflect.Descriptor instead.
func (*CDRResp) Descriptor() ([]byte, []int) {
	return file_cdr_proto_rawDescGZIP(), []int{1}
}

type RecordReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Imsi      string `protobuf:"bytes,1,opt,name=Imsi,json=imsi,proto3" json:"Imsi,omitempty"`
	StartTime uint64 `protobuf:"varint,2,opt,name=StartTime,json=start_time,proto3" json:"StartTime,omitempty"`
	EndTime   uint64 `protobuf:"varint,3,opt,name=EndTime,json=end_time,proto3" json:"EndTime,omitempty"`
	Policy    string `protobuf:"bytes,4,opt,name=Policy,json=policy,proto3" json:"Policy,omitempty"`
	SessionId uint64 `protobuf:"varint,5,opt,name=SessionId,json=session_id,proto3" json:"SessionId,omitempty"`
}

func (x *RecordReq) Reset() {
	*x = RecordReq{}
	mi := &file_cdr_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RecordReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordReq) ProtoMessage() {}

func (x *RecordReq) ProtoReflect() protoreflect.Message {
	mi := &file_cdr_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordReq.ProtoReflect.Descriptor instead.
func (*RecordReq) Descriptor() ([]byte, []int) {
	return file_cdr_proto_rawDescGZIP(), []int{2}
}

func (x *RecordReq) GetImsi() string {
	if x != nil {
		return x.Imsi
	}
	return ""
}

func (x *RecordReq) GetStartTime() uint64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *RecordReq) GetEndTime() uint64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

func (x *RecordReq) GetPolicy() string {
	if x != nil {
		return x.Policy
	}
	return ""
}

func (x *RecordReq) GetSessionId() uint64 {
	if x != nil {
		return x.SessionId
	}
	return 0
}

type RecordResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cdr []*CDR `protobuf:"bytes,1,rep,name=cdr,json=cdrs,proto3" json:"cdr,omitempty"`
}

func (x *RecordResp) Reset() {
	*x = RecordResp{}
	mi := &file_cdr_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RecordResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordResp) ProtoMessage() {}

func (x *RecordResp) ProtoReflect() protoreflect.Message {
	mi := &file_cdr_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordResp.ProtoReflect.Descriptor instead.
func (*RecordResp) Descriptor() ([]byte, []int) {
	return file_cdr_proto_rawDescGZIP(), []int{3}
}

func (x *RecordResp) GetCdr() []*CDR {
	if x != nil {
		return x.Cdr
	}
	return nil
}

type UsageReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Imsi      string `protobuf:"bytes,1,opt,name=Imsi,json=imsi,proto3" json:"Imsi,omitempty"`
	StartTime uint64 `protobuf:"varint,2,opt,name=StartTime,json=start_time,proto3" json:"StartTime,omitempty"`
	EndTime   uint64 `protobuf:"varint,3,opt,name=EndTime,json=end_time,proto3" json:"EndTime,omitempty"`
	Policy    string `protobuf:"bytes,4,opt,name=Policy,json=policy,proto3" json:"Policy,omitempty"`
	SessionId uint64 `protobuf:"varint,5,opt,name=SessionId,json=session_id,proto3" json:"SessionId,omitempty"`
}

func (x *UsageReq) Reset() {
	*x = UsageReq{}
	mi := &file_cdr_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UsageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsageReq) ProtoMessage() {}

func (x *UsageReq) ProtoReflect() protoreflect.Message {
	mi := &file_cdr_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsageReq.ProtoReflect.Descriptor instead.
func (*UsageReq) Descriptor() ([]byte, []int) {
	return file_cdr_proto_rawDescGZIP(), []int{4}
}

func (x *UsageReq) GetImsi() string {
	if x != nil {
		return x.Imsi
	}
	return ""
}

func (x *UsageReq) GetStartTime() uint64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *UsageReq) GetEndTime() uint64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

func (x *UsageReq) GetPolicy() string {
	if x != nil {
		return x.Policy
	}
	return ""
}

func (x *UsageReq) GetSessionId() uint64 {
	if x != nil {
		return x.SessionId
	}
	return 0
}

type UsageResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Imsi   string `protobuf:"bytes,1,opt,name=Imsi,json=imsi,proto3" json:"Imsi,omitempty"`
	Usage  uint64 `protobuf:"varint,2,opt,name=usage,proto3" json:"usage,omitempty"`
	Policy string `protobuf:"bytes,3,opt,name=policy,proto3" json:"policy,omitempty"`
}

func (x *UsageResp) Reset() {
	*x = UsageResp{}
	mi := &file_cdr_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UsageResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsageResp) ProtoMessage() {}

func (x *UsageResp) ProtoReflect() protoreflect.Message {
	mi := &file_cdr_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsageResp.ProtoReflect.Descriptor instead.
func (*UsageResp) Descriptor() ([]byte, []int) {
	return file_cdr_proto_rawDescGZIP(), []int{5}
}

func (x *UsageResp) GetImsi() string {
	if x != nil {
		return x.Imsi
	}
	return ""
}

func (x *UsageResp) GetUsage() uint64 {
	if x != nil {
		return x.Usage
	}
	return 0
}

func (x *UsageResp) GetPolicy() string {
	if x != nil {
		return x.Policy
	}
	return ""
}

type CycleUsageReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Imsi string `protobuf:"bytes,1,opt,name=Imsi,json=imsi,proto3" json:"Imsi,omitempty"`
}

func (x *CycleUsageReq) Reset() {
	*x = CycleUsageReq{}
	mi := &file_cdr_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CycleUsageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CycleUsageReq) ProtoMessage() {}

func (x *CycleUsageReq) ProtoReflect() protoreflect.Message {
	mi := &file_cdr_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CycleUsageReq.ProtoReflect.Descriptor instead.
func (*CycleUsageReq) Descriptor() ([]byte, []int) {
	return file_cdr_proto_rawDescGZIP(), []int{6}
}

func (x *CycleUsageReq) GetImsi() string {
	if x != nil {
		return x.Imsi
	}
	return ""
}

type CycleUsageResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Imsi             string `protobuf:"bytes,1,opt,name=Imsi,json=imsi,proto3" json:"Imsi,omitempty"`
	Historical       uint64 `protobuf:"varint,2,opt,name=historical,proto3" json:"historical,omitempty"`
	Usage            uint64 `protobuf:"varint,3,opt,name=usage,proto3" json:"usage,omitempty"`
	LastSessionUsage uint64 `protobuf:"varint,4,opt,name=LastSessionUsage,proto3" json:"LastSessionUsage,omitempty"`
	LastSessionId    uint64 `protobuf:"varint,5,opt,name=LastSessionId,proto3" json:"LastSessionId,omitempty"`
	LastNodeId       string `protobuf:"bytes,6,opt,name=lastNodeId,proto3" json:"lastNodeId,omitempty"`
	LastCDRUpdatedAt uint64 `protobuf:"varint,7,opt,name=LastCDRUpdatedAt,proto3" json:"LastCDRUpdatedAt,omitempty"`
	Policy           string `protobuf:"bytes,8,opt,name=Policy,proto3" json:"Policy,omitempty"`
}

func (x *CycleUsageResp) Reset() {
	*x = CycleUsageResp{}
	mi := &file_cdr_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CycleUsageResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CycleUsageResp) ProtoMessage() {}

func (x *CycleUsageResp) ProtoReflect() protoreflect.Message {
	mi := &file_cdr_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CycleUsageResp.ProtoReflect.Descriptor instead.
func (*CycleUsageResp) Descriptor() ([]byte, []int) {
	return file_cdr_proto_rawDescGZIP(), []int{7}
}

func (x *CycleUsageResp) GetImsi() string {
	if x != nil {
		return x.Imsi
	}
	return ""
}

func (x *CycleUsageResp) GetHistorical() uint64 {
	if x != nil {
		return x.Historical
	}
	return 0
}

func (x *CycleUsageResp) GetUsage() uint64 {
	if x != nil {
		return x.Usage
	}
	return 0
}

func (x *CycleUsageResp) GetLastSessionUsage() uint64 {
	if x != nil {
		return x.LastSessionUsage
	}
	return 0
}

func (x *CycleUsageResp) GetLastSessionId() uint64 {
	if x != nil {
		return x.LastSessionId
	}
	return 0
}

func (x *CycleUsageResp) GetLastNodeId() string {
	if x != nil {
		return x.LastNodeId
	}
	return ""
}

func (x *CycleUsageResp) GetLastCDRUpdatedAt() uint64 {
	if x != nil {
		return x.LastCDRUpdatedAt
	}
	return 0
}

func (x *CycleUsageResp) GetPolicy() string {
	if x != nil {
		return x.Policy
	}
	return ""
}

type UsageForPeriodReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Imsi      string `protobuf:"bytes,1,opt,name=Imsi,json=imsi,proto3" json:"Imsi,omitempty"`
	StartTime uint64 `protobuf:"varint,4,opt,name=StartTime,proto3" json:"StartTime,omitempty"`
	EndTime   uint64 `protobuf:"varint,5,opt,name=EndTime,proto3" json:"EndTime,omitempty"`
}

func (x *UsageForPeriodReq) Reset() {
	*x = UsageForPeriodReq{}
	mi := &file_cdr_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UsageForPeriodReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsageForPeriodReq) ProtoMessage() {}

func (x *UsageForPeriodReq) ProtoReflect() protoreflect.Message {
	mi := &file_cdr_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsageForPeriodReq.ProtoReflect.Descriptor instead.
func (*UsageForPeriodReq) Descriptor() ([]byte, []int) {
	return file_cdr_proto_rawDescGZIP(), []int{8}
}

func (x *UsageForPeriodReq) GetImsi() string {
	if x != nil {
		return x.Imsi
	}
	return ""
}

func (x *UsageForPeriodReq) GetStartTime() uint64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *UsageForPeriodReq) GetEndTime() uint64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

type UsageForPeriodResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Usage uint64 `protobuf:"varint,1,opt,name=Usage,proto3" json:"Usage,omitempty"`
}

func (x *UsageForPeriodResp) Reset() {
	*x = UsageForPeriodResp{}
	mi := &file_cdr_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UsageForPeriodResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsageForPeriodResp) ProtoMessage() {}

func (x *UsageForPeriodResp) ProtoReflect() protoreflect.Message {
	mi := &file_cdr_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsageForPeriodResp.ProtoReflect.Descriptor instead.
func (*UsageForPeriodResp) Descriptor() ([]byte, []int) {
	return file_cdr_proto_rawDescGZIP(), []int{9}
}

func (x *UsageForPeriodResp) GetUsage() uint64 {
	if x != nil {
		return x.Usage
	}
	return 0
}

type QueryUsageReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Imsi     string   `protobuf:"bytes,1,opt,name=Imsi,json=imsi,proto3" json:"Imsi,omitempty"`
	NodeId   string   `protobuf:"bytes,2,opt,name=NodeId,proto3" json:"NodeId,omitempty"`
	Session  uint64   `protobuf:"varint,3,opt,name=Session,proto3" json:"Session,omitempty"`
	From     uint64   `protobuf:"varint,4,opt,name=From,proto3" json:"From,omitempty"`
	To       uint64   `protobuf:"varint,5,opt,name=To,proto3" json:"To,omitempty"`
	Policies []string `protobuf:"bytes,6,rep,name=Policies,proto3" json:"Policies,omitempty"`
	Count    uint32   `protobuf:"varint,7,opt,name=Count,proto3" json:"Count,omitempty"`
	Sort     bool     `protobuf:"varint,8,opt,name=Sort,proto3" json:"Sort,omitempty"`
}

func (x *QueryUsageReq) Reset() {
	*x = QueryUsageReq{}
	mi := &file_cdr_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryUsageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryUsageReq) ProtoMessage() {}

func (x *QueryUsageReq) ProtoReflect() protoreflect.Message {
	mi := &file_cdr_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryUsageReq.ProtoReflect.Descriptor instead.
func (*QueryUsageReq) Descriptor() ([]byte, []int) {
	return file_cdr_proto_rawDescGZIP(), []int{10}
}

func (x *QueryUsageReq) GetImsi() string {
	if x != nil {
		return x.Imsi
	}
	return ""
}

func (x *QueryUsageReq) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *QueryUsageReq) GetSession() uint64 {
	if x != nil {
		return x.Session
	}
	return 0
}

func (x *QueryUsageReq) GetFrom() uint64 {
	if x != nil {
		return x.From
	}
	return 0
}

func (x *QueryUsageReq) GetTo() uint64 {
	if x != nil {
		return x.To
	}
	return 0
}

func (x *QueryUsageReq) GetPolicies() []string {
	if x != nil {
		return x.Policies
	}
	return nil
}

func (x *QueryUsageReq) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *QueryUsageReq) GetSort() bool {
	if x != nil {
		return x.Sort
	}
	return false
}

type QueryUsageResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Usage uint64 `protobuf:"varint,1,opt,name=Usage,proto3" json:"Usage,omitempty"`
}

func (x *QueryUsageResp) Reset() {
	*x = QueryUsageResp{}
	mi := &file_cdr_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryUsageResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryUsageResp) ProtoMessage() {}

func (x *QueryUsageResp) ProtoReflect() protoreflect.Message {
	mi := &file_cdr_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryUsageResp.ProtoReflect.Descriptor instead.
func (*QueryUsageResp) Descriptor() ([]byte, []int) {
	return file_cdr_proto_rawDescGZIP(), []int{11}
}

func (x *QueryUsageResp) GetUsage() uint64 {
	if x != nil {
		return x.Usage
	}
	return 0
}

var File_cdr_proto protoreflect.FileDescriptor

var file_cdr_proto_rawDesc = []byte{
	0x0a, 0x09, 0x63, 0x64, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x75, 0x6b, 0x61,
	0x6d, 0x61, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x64,
	0x72, 0x2e, 0x76, 0x31, 0x1a, 0x0f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xda, 0x02, 0x0a, 0x03, 0x43, 0x44, 0x52, 0x12, 0x1b, 0x0a,
	0x07, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a,
	0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x06, 0x4e, 0x6f,
	0x64, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6e, 0x6f, 0x64, 0x65,
	0x5f, 0x69, 0x64, 0x12, 0x1e, 0x0a, 0x04, 0x49, 0x6d, 0x73, 0x69, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x0a, 0xe2, 0xdf, 0x1f, 0x06, 0x58, 0x01, 0x70, 0x05, 0x78, 0x10, 0x52, 0x04, 0x69,
	0x6d, 0x73, 0x69, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x19, 0x0a, 0x07, 0x41,
	0x70, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x61, 0x70,
	0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x11, 0x0a, 0x02, 0x49, 0x70, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x75, 0x65, 0x5f, 0x69, 0x70, 0x12, 0x1d, 0x0a, 0x09, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x07, 0x45, 0x6e, 0x64, 0x54,
	0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0d, 0x4c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x6c, 0x61, 0x73, 0x74,
	0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x12, 0x19, 0x0a, 0x07, 0x54,
	0x78, 0x42, 0x79, 0x74, 0x65, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x74, 0x78,
	0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x12, 0x19, 0x0a, 0x07, 0x52, 0x78, 0x42, 0x79, 0x74, 0x65,
	0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x72, 0x78, 0x5f, 0x62, 0x79, 0x74, 0x65,
	0x73, 0x12, 0x1f, 0x0a, 0x0a, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x42, 0x79, 0x74, 0x65, 0x73, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x62, 0x79, 0x74,
	0x65, 0x73, 0x22, 0x09, 0x0a, 0x07, 0x43, 0x44, 0x52, 0x52, 0x65, 0x73, 0x70, 0x22, 0x9c, 0x01,
	0x0a, 0x09, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x12, 0x1e, 0x0a, 0x04, 0x49,
	0x6d, 0x73, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe2, 0xdf, 0x1f, 0x06, 0x58,
	0x01, 0x70, 0x05, 0x78, 0x10, 0x52, 0x04, 0x69, 0x6d, 0x73, 0x69, 0x12, 0x1d, 0x0a, 0x09, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x07, 0x45, 0x6e,
	0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x65, 0x6e, 0x64,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x1d, 0x0a,
	0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x3d, 0x0a, 0x0a,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x12, 0x2f, 0x0a, 0x03, 0x63, 0x64,
	0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e,
	0x75, 0x6b, 0x61, 0x6d, 0x61, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x64, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x44, 0x52, 0x52, 0x04, 0x63, 0x64, 0x72, 0x73, 0x22, 0x9b, 0x01, 0x0a, 0x08,
	0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x12, 0x1e, 0x0a, 0x04, 0x49, 0x6d, 0x73, 0x69,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe2, 0xdf, 0x1f, 0x06, 0x58, 0x01, 0x70, 0x05,
	0x78, 0x10, 0x52, 0x04, 0x69, 0x6d, 0x73, 0x69, 0x12, 0x1d, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x07, 0x45, 0x6e, 0x64, 0x54, 0x69,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x1d, 0x0a, 0x09, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x59, 0x0a, 0x09, 0x55, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x1e, 0x0a, 0x04, 0x49, 0x6d, 0x73, 0x69, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe2, 0xdf, 0x1f, 0x06, 0x58, 0x01, 0x70, 0x05, 0x78, 0x10,
	0x52, 0x04, 0x69, 0x6d, 0x73, 0x69, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x22, 0x2f, 0x0a, 0x0d, 0x43, 0x79, 0x63, 0x6c, 0x65, 0x55, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x12, 0x1e, 0x0a, 0x04, 0x49, 0x6d, 0x73, 0x69, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x0a, 0xe2, 0xdf, 0x1f, 0x06, 0x58, 0x01, 0x70, 0x05, 0x78, 0x10, 0x52,
	0x04, 0x69, 0x6d, 0x73, 0x69, 0x22, 0x9c, 0x02, 0x0a, 0x0e, 0x43, 0x79, 0x63, 0x6c, 0x65, 0x55,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x1e, 0x0a, 0x04, 0x49, 0x6d, 0x73, 0x69,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe2, 0xdf, 0x1f, 0x06, 0x58, 0x01, 0x70, 0x05,
	0x78, 0x10, 0x52, 0x04, 0x69, 0x6d, 0x73, 0x69, 0x12, 0x1e, 0x0a, 0x0a, 0x68, 0x69, 0x73, 0x74,
	0x6f, 0x72, 0x69, 0x63, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x68, 0x69,
	0x73, 0x74, 0x6f, 0x72, 0x69, 0x63, 0x61, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2a,
	0x0a, 0x10, 0x4c, 0x61, 0x73, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x55, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x4c, 0x61, 0x73, 0x74, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x4c, 0x61,
	0x73, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0d, 0x4c, 0x61, 0x73, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64,
	0x12, 0x1e, 0x0a, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64,
	0x12, 0x2a, 0x0a, 0x10, 0x4c, 0x61, 0x73, 0x74, 0x43, 0x44, 0x52, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x4c, 0x61, 0x73, 0x74,
	0x43, 0x44, 0x52, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x22, 0x6b, 0x0a, 0x11, 0x55, 0x73, 0x61, 0x67, 0x65, 0x46, 0x6f, 0x72,
	0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x52, 0x65, 0x71, 0x12, 0x1e, 0x0a, 0x04, 0x49, 0x6d, 0x73,
	0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe2, 0xdf, 0x1f, 0x06, 0x58, 0x01, 0x70,
	0x05, 0x78, 0x10, 0x52, 0x04, 0x69, 0x6d, 0x73, 0x69, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x45, 0x6e, 0x64, 0x54, 0x69,
	0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x45, 0x6e, 0x64, 0x54, 0x69, 0x6d,
	0x65, 0x22, 0x2a, 0x0a, 0x12, 0x55, 0x73, 0x61, 0x67, 0x65, 0x46, 0x6f, 0x72, 0x50, 0x65, 0x72,
	0x69, 0x6f, 0x64, 0x52, 0x65, 0x73, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x55, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x55, 0x73, 0x61, 0x67, 0x65, 0x22, 0xcb, 0x01,
	0x0a, 0x0d, 0x51, 0x75, 0x65, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x12,
	0x1e, 0x0a, 0x04, 0x49, 0x6d, 0x73, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe2,
	0xdf, 0x1f, 0x06, 0x58, 0x01, 0x70, 0x05, 0x78, 0x10, 0x52, 0x04, 0x69, 0x6d, 0x73, 0x69, 0x12,
	0x16, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x04, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x54, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x02, 0x54, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65,
	0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x05, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x6f, 0x72, 0x74, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x53, 0x6f, 0x72, 0x74, 0x22, 0x26, 0x0a, 0x0e, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x14, 0x0a,
	0x05, 0x55, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x55, 0x73,
	0x61, 0x67, 0x65, 0x32, 0xae, 0x04, 0x0a, 0x0a, 0x43, 0x44, 0x52, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x49, 0x0a, 0x07, 0x50, 0x6f, 0x73, 0x74, 0x43, 0x44, 0x52, 0x12, 0x1c, 0x2e,
	0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x2e, 0x63, 0x64, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x44, 0x52, 0x1a, 0x20, 0x2e, 0x75, 0x6b,
	0x61, 0x6d, 0x61, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63,
	0x64, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x44, 0x52, 0x52, 0x65, 0x73, 0x70, 0x12, 0x51, 0x0a,
	0x06, 0x47, 0x65, 0x74, 0x43, 0x44, 0x52, 0x12, 0x22, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e,
	0x75, 0x6b, 0x61, 0x6d, 0x61, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x64, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x1a, 0x23, 0x2e, 0x75, 0x6b,
	0x61, 0x6d, 0x61, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63,
	0x64, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x51, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x21, 0x2e, 0x75,
	0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e,
	0x63, 0x64, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x1a,
	0x22, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x2e, 0x63, 0x64, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x12, 0x6c, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x55, 0x73, 0x61, 0x67, 0x65, 0x46,
	0x6f, 0x72, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x12, 0x2a, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61,
	0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x64, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x55, 0x73, 0x61, 0x67, 0x65, 0x46, 0x6f, 0x72, 0x50, 0x65, 0x72, 0x69, 0x6f,
	0x64, 0x52, 0x65, 0x71, 0x1a, 0x2b, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x75, 0x6b, 0x61,
	0x6d, 0x61, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x64, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55,
	0x73, 0x61, 0x67, 0x65, 0x46, 0x6f, 0x72, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x52, 0x65, 0x73,
	0x70, 0x12, 0x62, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x55, 0x73, 0x61, 0x67, 0x65, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x73, 0x12, 0x26, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x75, 0x6b, 0x61,
	0x6d, 0x61, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x64, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x79, 0x63, 0x6c, 0x65, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x27, 0x2e, 0x75,
	0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e,
	0x63, 0x64, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x79, 0x63, 0x6c, 0x65, 0x55, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x5d, 0x0a, 0x0a, 0x51, 0x75, 0x65, 0x72, 0x79, 0x55, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x26, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2e, 0x75, 0x6b, 0x61, 0x6d,
	0x61, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x64, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75,
	0x65, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x27, 0x2e, 0x75, 0x6b,
	0x61, 0x6d, 0x61, 0x2e, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63,
	0x64, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x42, 0x37, 0x5a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2f, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2f, 0x73,
	0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f, 0x75, 0x6b, 0x61, 0x6d, 0x61, 0x2d, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x2f, 0x63, 0x64, 0x72, 0x2f, 0x70, 0x62, 0x2f, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cdr_proto_rawDescOnce sync.Once
	file_cdr_proto_rawDescData = file_cdr_proto_rawDesc
)

func file_cdr_proto_rawDescGZIP() []byte {
	file_cdr_proto_rawDescOnce.Do(func() {
		file_cdr_proto_rawDescData = protoimpl.X.CompressGZIP(file_cdr_proto_rawDescData)
	})
	return file_cdr_proto_rawDescData
}

var file_cdr_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_cdr_proto_goTypes = []any{
	(*CDR)(nil),                // 0: ukama.ukamaagent.cdr.v1.CDR
	(*CDRResp)(nil),            // 1: ukama.ukamaagent.cdr.v1.CDRResp
	(*RecordReq)(nil),          // 2: ukama.ukamaagent.cdr.v1.RecordReq
	(*RecordResp)(nil),         // 3: ukama.ukamaagent.cdr.v1.RecordResp
	(*UsageReq)(nil),           // 4: ukama.ukamaagent.cdr.v1.UsageReq
	(*UsageResp)(nil),          // 5: ukama.ukamaagent.cdr.v1.UsageResp
	(*CycleUsageReq)(nil),      // 6: ukama.ukamaagent.cdr.v1.CycleUsageReq
	(*CycleUsageResp)(nil),     // 7: ukama.ukamaagent.cdr.v1.CycleUsageResp
	(*UsageForPeriodReq)(nil),  // 8: ukama.ukamaagent.cdr.v1.UsageForPeriodReq
	(*UsageForPeriodResp)(nil), // 9: ukama.ukamaagent.cdr.v1.UsageForPeriodResp
	(*QueryUsageReq)(nil),      // 10: ukama.ukamaagent.cdr.v1.QueryUsageReq
	(*QueryUsageResp)(nil),     // 11: ukama.ukamaagent.cdr.v1.QueryUsageResp
}
var file_cdr_proto_depIdxs = []int32{
	0,  // 0: ukama.ukamaagent.cdr.v1.RecordResp.cdr:type_name -> ukama.ukamaagent.cdr.v1.CDR
	0,  // 1: ukama.ukamaagent.cdr.v1.CDRService.PostCDR:input_type -> ukama.ukamaagent.cdr.v1.CDR
	2,  // 2: ukama.ukamaagent.cdr.v1.CDRService.GetCDR:input_type -> ukama.ukamaagent.cdr.v1.RecordReq
	4,  // 3: ukama.ukamaagent.cdr.v1.CDRService.GetUsage:input_type -> ukama.ukamaagent.cdr.v1.UsageReq
	8,  // 4: ukama.ukamaagent.cdr.v1.CDRService.GetUsageForPeriod:input_type -> ukama.ukamaagent.cdr.v1.UsageForPeriodReq
	6,  // 5: ukama.ukamaagent.cdr.v1.CDRService.GetUsageDetails:input_type -> ukama.ukamaagent.cdr.v1.CycleUsageReq
	10, // 6: ukama.ukamaagent.cdr.v1.CDRService.QueryUsage:input_type -> ukama.ukamaagent.cdr.v1.QueryUsageReq
	1,  // 7: ukama.ukamaagent.cdr.v1.CDRService.PostCDR:output_type -> ukama.ukamaagent.cdr.v1.CDRResp
	3,  // 8: ukama.ukamaagent.cdr.v1.CDRService.GetCDR:output_type -> ukama.ukamaagent.cdr.v1.RecordResp
	5,  // 9: ukama.ukamaagent.cdr.v1.CDRService.GetUsage:output_type -> ukama.ukamaagent.cdr.v1.UsageResp
	9,  // 10: ukama.ukamaagent.cdr.v1.CDRService.GetUsageForPeriod:output_type -> ukama.ukamaagent.cdr.v1.UsageForPeriodResp
	7,  // 11: ukama.ukamaagent.cdr.v1.CDRService.GetUsageDetails:output_type -> ukama.ukamaagent.cdr.v1.CycleUsageResp
	11, // 12: ukama.ukamaagent.cdr.v1.CDRService.QueryUsage:output_type -> ukama.ukamaagent.cdr.v1.QueryUsageResp
	7,  // [7:13] is the sub-list for method output_type
	1,  // [1:7] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_cdr_proto_init() }
func file_cdr_proto_init() {
	if File_cdr_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cdr_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cdr_proto_goTypes,
		DependencyIndexes: file_cdr_proto_depIdxs,
		MessageInfos:      file_cdr_proto_msgTypes,
	}.Build()
	File_cdr_proto = out.File
	file_cdr_proto_rawDesc = nil
	file_cdr_proto_goTypes = nil
	file_cdr_proto_depIdxs = nil
}
