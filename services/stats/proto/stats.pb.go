// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: services/stats/proto/stats.proto

package stats

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type GetServiceTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppID     string `protobuf:"bytes,1,opt,name=appID,proto3" json:"appID,omitempty"`
	AppSECRET string `protobuf:"bytes,2,opt,name=appSECRET,proto3" json:"appSECRET,omitempty"`
}

func (x *GetServiceTokenRequest) Reset() {
	*x = GetServiceTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_stats_proto_stats_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetServiceTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetServiceTokenRequest) ProtoMessage() {}

func (x *GetServiceTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_stats_proto_stats_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetServiceTokenRequest.ProtoReflect.Descriptor instead.
func (*GetServiceTokenRequest) Descriptor() ([]byte, []int) {
	return file_services_stats_proto_stats_proto_rawDescGZIP(), []int{0}
}

func (x *GetServiceTokenRequest) GetAppID() string {
	if x != nil {
		return x.AppID
	}
	return ""
}

func (x *GetServiceTokenRequest) GetAppSECRET() string {
	if x != nil {
		return x.AppSECRET
	}
	return ""
}

type GetServiceTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *GetServiceTokenResponse) Reset() {
	*x = GetServiceTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_stats_proto_stats_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetServiceTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetServiceTokenResponse) ProtoMessage() {}

func (x *GetServiceTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_stats_proto_stats_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetServiceTokenResponse.ProtoReflect.Descriptor instead.
func (*GetServiceTokenResponse) Descriptor() ([]byte, []int) {
	return file_services_stats_proto_stats_proto_rawDescGZIP(), []int{1}
}

func (x *GetServiceTokenResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type SingleStat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int32                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserUID   string               `protobuf:"bytes,2,opt,name=userUID,proto3" json:"userUID,omitempty"`
	Action    string               `protobuf:"bytes,3,opt,name=action,proto3" json:"action,omitempty"`
	Timestamp *timestamp.Timestamp `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Input     []byte               `protobuf:"bytes,5,opt,name=input,proto3" json:"input,omitempty"`
	Output    []byte               `protobuf:"bytes,6,opt,name=output,proto3" json:"output,omitempty"`
}

func (x *SingleStat) Reset() {
	*x = SingleStat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_stats_proto_stats_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SingleStat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SingleStat) ProtoMessage() {}

func (x *SingleStat) ProtoReflect() protoreflect.Message {
	mi := &file_services_stats_proto_stats_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SingleStat.ProtoReflect.Descriptor instead.
func (*SingleStat) Descriptor() ([]byte, []int) {
	return file_services_stats_proto_stats_proto_rawDescGZIP(), []int{2}
}

func (x *SingleStat) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SingleStat) GetUserUID() string {
	if x != nil {
		return x.UserUID
	}
	return ""
}

func (x *SingleStat) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *SingleStat) GetTimestamp() *timestamp.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *SingleStat) GetInput() []byte {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *SingleStat) GetOutput() []byte {
	if x != nil {
		return x.Output
	}
	return nil
}

type ListStatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageSize   int32  `protobuf:"varint,1,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
	PageNumber int32  `protobuf:"varint,2,opt,name=pageNumber,proto3" json:"pageNumber,omitempty"`
	Token      string `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *ListStatsRequest) Reset() {
	*x = ListStatsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_stats_proto_stats_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListStatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListStatsRequest) ProtoMessage() {}

func (x *ListStatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_stats_proto_stats_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListStatsRequest.ProtoReflect.Descriptor instead.
func (*ListStatsRequest) Descriptor() ([]byte, []int) {
	return file_services_stats_proto_stats_proto_rawDescGZIP(), []int{3}
}

func (x *ListStatsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListStatsRequest) GetPageNumber() int32 {
	if x != nil {
		return x.PageNumber
	}
	return 0
}

func (x *ListStatsRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type ListStatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stats      []*SingleStat `protobuf:"bytes,1,rep,name=stats,proto3" json:"stats,omitempty"`
	PageSize   int32         `protobuf:"varint,2,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
	PageNumber int32         `protobuf:"varint,3,opt,name=pageNumber,proto3" json:"pageNumber,omitempty"`
	PageCount  int32         `protobuf:"varint,4,opt,name=pageCount,proto3" json:"pageCount,omitempty"`
	Token      string        `protobuf:"bytes,5,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *ListStatsResponse) Reset() {
	*x = ListStatsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_stats_proto_stats_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListStatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListStatsResponse) ProtoMessage() {}

func (x *ListStatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_stats_proto_stats_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListStatsResponse.ProtoReflect.Descriptor instead.
func (*ListStatsResponse) Descriptor() ([]byte, []int) {
	return file_services_stats_proto_stats_proto_rawDescGZIP(), []int{4}
}

func (x *ListStatsResponse) GetStats() []*SingleStat {
	if x != nil {
		return x.Stats
	}
	return nil
}

func (x *ListStatsResponse) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListStatsResponse) GetPageNumber() int32 {
	if x != nil {
		return x.PageNumber
	}
	return 0
}

func (x *ListStatsResponse) GetPageCount() int32 {
	if x != nil {
		return x.PageCount
	}
	return 0
}

func (x *ListStatsResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type AddStatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserUID string `protobuf:"bytes,1,opt,name=userUID,proto3" json:"userUID,omitempty"`
	Action  string `protobuf:"bytes,2,opt,name=action,proto3" json:"action,omitempty"`
	Input   []byte `protobuf:"bytes,3,opt,name=input,proto3" json:"input,omitempty"`
	Output  []byte `protobuf:"bytes,4,opt,name=output,proto3" json:"output,omitempty"`
	Token   string `protobuf:"bytes,5,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *AddStatsRequest) Reset() {
	*x = AddStatsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_stats_proto_stats_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddStatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddStatsRequest) ProtoMessage() {}

func (x *AddStatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_stats_proto_stats_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddStatsRequest.ProtoReflect.Descriptor instead.
func (*AddStatsRequest) Descriptor() ([]byte, []int) {
	return file_services_stats_proto_stats_proto_rawDescGZIP(), []int{5}
}

func (x *AddStatsRequest) GetUserUID() string {
	if x != nil {
		return x.UserUID
	}
	return ""
}

func (x *AddStatsRequest) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *AddStatsRequest) GetInput() []byte {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *AddStatsRequest) GetOutput() []byte {
	if x != nil {
		return x.Output
	}
	return nil
}

func (x *AddStatsRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type StatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stats *SingleStat `protobuf:"bytes,1,opt,name=Stats,proto3" json:"Stats,omitempty"`
}

func (x *StatsResponse) Reset() {
	*x = StatsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_stats_proto_stats_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatsResponse) ProtoMessage() {}

func (x *StatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_stats_proto_stats_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatsResponse.ProtoReflect.Descriptor instead.
func (*StatsResponse) Descriptor() ([]byte, []int) {
	return file_services_stats_proto_stats_proto_rawDescGZIP(), []int{6}
}

func (x *StatsResponse) GetStats() *SingleStat {
	if x != nil {
		return x.Stats
	}
	return nil
}

var File_services_stats_proto_stats_proto protoreflect.FileDescriptor

var file_services_stats_proto_stats_proto_rawDesc = []byte{
	0x0a, 0x20, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x73,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4c, 0x0a, 0x16, 0x47, 0x65,
	0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x70, 0x70, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x70,
	0x70, 0x53, 0x45, 0x43, 0x52, 0x45, 0x54, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61,
	0x70, 0x70, 0x53, 0x45, 0x43, 0x52, 0x45, 0x54, 0x22, 0x2f, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0xb6, 0x01, 0x0a, 0x0a, 0x53, 0x69,
	0x6e, 0x67, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x55, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x55,
	0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75,
	0x74, 0x70, 0x75, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70,
	0x75, 0x74, 0x22, 0x64, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69,
	0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69,
	0x7a, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0xac, 0x01, 0x0a, 0x11, 0x4c, 0x69, 0x73,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27,
	0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53,
	0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53,
	0x69, 0x7a, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x87, 0x01, 0x0a, 0x0f, 0x41, 0x64, 0x64, 0x53,
	0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x55, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x55, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a,
	0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x69, 0x6e,
	0x70, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x22, 0x38, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x27, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x52, 0x05, 0x53, 0x74, 0x61, 0x74, 0x73, 0x32, 0xef, 0x03, 0x0a, 0x05,
	0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x50, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1d, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73,
	0x2e, 0x47, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e,
	0x47, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x17, 0x2e, 0x73,
	0x74, 0x61, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x40, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x53, 0x74,
	0x61, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x41, 0x64, 0x64, 0x53,
	0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x73, 0x74,
	0x61, 0x74, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x42, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x65, 0x77, 0x73, 0x53, 0x74, 0x61,
	0x74, 0x73, 0x12, 0x17, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x73, 0x74,
	0x61, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x0c, 0x41, 0x64, 0x64, 0x4e, 0x65, 0x77, 0x73,
	0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x41, 0x64,
	0x64, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e,
	0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x17, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x18, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x10, 0x41,
	0x64, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12,
	0x16, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x41, 0x64, 0x64, 0x53, 0x74, 0x61, 0x74, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_stats_proto_stats_proto_rawDescOnce sync.Once
	file_services_stats_proto_stats_proto_rawDescData = file_services_stats_proto_stats_proto_rawDesc
)

func file_services_stats_proto_stats_proto_rawDescGZIP() []byte {
	file_services_stats_proto_stats_proto_rawDescOnce.Do(func() {
		file_services_stats_proto_stats_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_stats_proto_stats_proto_rawDescData)
	})
	return file_services_stats_proto_stats_proto_rawDescData
}

var file_services_stats_proto_stats_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_services_stats_proto_stats_proto_goTypes = []interface{}{
	(*GetServiceTokenRequest)(nil),  // 0: stats.GetServiceTokenRequest
	(*GetServiceTokenResponse)(nil), // 1: stats.GetServiceTokenResponse
	(*SingleStat)(nil),              // 2: stats.SingleStat
	(*ListStatsRequest)(nil),        // 3: stats.ListStatsRequest
	(*ListStatsResponse)(nil),       // 4: stats.ListStatsResponse
	(*AddStatsRequest)(nil),         // 5: stats.AddStatsRequest
	(*StatsResponse)(nil),           // 6: stats.StatsResponse
	(*timestamp.Timestamp)(nil),     // 7: google.protobuf.Timestamp
}
var file_services_stats_proto_stats_proto_depIdxs = []int32{
	7,  // 0: stats.SingleStat.timestamp:type_name -> google.protobuf.Timestamp
	2,  // 1: stats.ListStatsResponse.stats:type_name -> stats.SingleStat
	2,  // 2: stats.StatsResponse.Stats:type_name -> stats.SingleStat
	0,  // 3: stats.Stats.GetServiceToken:input_type -> stats.GetServiceTokenRequest
	3,  // 4: stats.Stats.ListAccountsStats:input_type -> stats.ListStatsRequest
	5,  // 5: stats.Stats.AddAccountsStats:input_type -> stats.AddStatsRequest
	3,  // 6: stats.Stats.ListNewsStats:input_type -> stats.ListStatsRequest
	5,  // 7: stats.Stats.AddNewsStats:input_type -> stats.AddStatsRequest
	3,  // 8: stats.Stats.ListCommentsStats:input_type -> stats.ListStatsRequest
	5,  // 9: stats.Stats.AddCommentsStats:input_type -> stats.AddStatsRequest
	1,  // 10: stats.Stats.GetServiceToken:output_type -> stats.GetServiceTokenResponse
	4,  // 11: stats.Stats.ListAccountsStats:output_type -> stats.ListStatsResponse
	6,  // 12: stats.Stats.AddAccountsStats:output_type -> stats.StatsResponse
	4,  // 13: stats.Stats.ListNewsStats:output_type -> stats.ListStatsResponse
	6,  // 14: stats.Stats.AddNewsStats:output_type -> stats.StatsResponse
	4,  // 15: stats.Stats.ListCommentsStats:output_type -> stats.ListStatsResponse
	6,  // 16: stats.Stats.AddCommentsStats:output_type -> stats.StatsResponse
	10, // [10:17] is the sub-list for method output_type
	3,  // [3:10] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_services_stats_proto_stats_proto_init() }
func file_services_stats_proto_stats_proto_init() {
	if File_services_stats_proto_stats_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_services_stats_proto_stats_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetServiceTokenRequest); i {
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
		file_services_stats_proto_stats_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetServiceTokenResponse); i {
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
		file_services_stats_proto_stats_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SingleStat); i {
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
		file_services_stats_proto_stats_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListStatsRequest); i {
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
		file_services_stats_proto_stats_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListStatsResponse); i {
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
		file_services_stats_proto_stats_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddStatsRequest); i {
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
		file_services_stats_proto_stats_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatsResponse); i {
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
			RawDescriptor: file_services_stats_proto_stats_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_stats_proto_stats_proto_goTypes,
		DependencyIndexes: file_services_stats_proto_stats_proto_depIdxs,
		MessageInfos:      file_services_stats_proto_stats_proto_msgTypes,
	}.Build()
	File_services_stats_proto_stats_proto = out.File
	file_services_stats_proto_stats_proto_rawDesc = nil
	file_services_stats_proto_stats_proto_goTypes = nil
	file_services_stats_proto_stats_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// StatsClient is the client API for Stats service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StatsClient interface {
	GetServiceToken(ctx context.Context, in *GetServiceTokenRequest, opts ...grpc.CallOption) (*GetServiceTokenResponse, error)
	ListAccountsStats(ctx context.Context, in *ListStatsRequest, opts ...grpc.CallOption) (*ListStatsResponse, error)
	AddAccountsStats(ctx context.Context, in *AddStatsRequest, opts ...grpc.CallOption) (*StatsResponse, error)
	ListNewsStats(ctx context.Context, in *ListStatsRequest, opts ...grpc.CallOption) (*ListStatsResponse, error)
	AddNewsStats(ctx context.Context, in *AddStatsRequest, opts ...grpc.CallOption) (*StatsResponse, error)
	ListCommentsStats(ctx context.Context, in *ListStatsRequest, opts ...grpc.CallOption) (*ListStatsResponse, error)
	AddCommentsStats(ctx context.Context, in *AddStatsRequest, opts ...grpc.CallOption) (*StatsResponse, error)
}

type statsClient struct {
	cc grpc.ClientConnInterface
}

func NewStatsClient(cc grpc.ClientConnInterface) StatsClient {
	return &statsClient{cc}
}

func (c *statsClient) GetServiceToken(ctx context.Context, in *GetServiceTokenRequest, opts ...grpc.CallOption) (*GetServiceTokenResponse, error) {
	out := new(GetServiceTokenResponse)
	err := c.cc.Invoke(ctx, "/stats.Stats/GetServiceToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statsClient) ListAccountsStats(ctx context.Context, in *ListStatsRequest, opts ...grpc.CallOption) (*ListStatsResponse, error) {
	out := new(ListStatsResponse)
	err := c.cc.Invoke(ctx, "/stats.Stats/ListAccountsStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statsClient) AddAccountsStats(ctx context.Context, in *AddStatsRequest, opts ...grpc.CallOption) (*StatsResponse, error) {
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, "/stats.Stats/AddAccountsStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statsClient) ListNewsStats(ctx context.Context, in *ListStatsRequest, opts ...grpc.CallOption) (*ListStatsResponse, error) {
	out := new(ListStatsResponse)
	err := c.cc.Invoke(ctx, "/stats.Stats/ListNewsStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statsClient) AddNewsStats(ctx context.Context, in *AddStatsRequest, opts ...grpc.CallOption) (*StatsResponse, error) {
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, "/stats.Stats/AddNewsStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statsClient) ListCommentsStats(ctx context.Context, in *ListStatsRequest, opts ...grpc.CallOption) (*ListStatsResponse, error) {
	out := new(ListStatsResponse)
	err := c.cc.Invoke(ctx, "/stats.Stats/ListCommentsStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statsClient) AddCommentsStats(ctx context.Context, in *AddStatsRequest, opts ...grpc.CallOption) (*StatsResponse, error) {
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, "/stats.Stats/AddCommentsStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatsServer is the server API for Stats service.
type StatsServer interface {
	GetServiceToken(context.Context, *GetServiceTokenRequest) (*GetServiceTokenResponse, error)
	ListAccountsStats(context.Context, *ListStatsRequest) (*ListStatsResponse, error)
	AddAccountsStats(context.Context, *AddStatsRequest) (*StatsResponse, error)
	ListNewsStats(context.Context, *ListStatsRequest) (*ListStatsResponse, error)
	AddNewsStats(context.Context, *AddStatsRequest) (*StatsResponse, error)
	ListCommentsStats(context.Context, *ListStatsRequest) (*ListStatsResponse, error)
	AddCommentsStats(context.Context, *AddStatsRequest) (*StatsResponse, error)
}

// UnimplementedStatsServer can be embedded to have forward compatible implementations.
type UnimplementedStatsServer struct {
}

func (*UnimplementedStatsServer) GetServiceToken(context.Context, *GetServiceTokenRequest) (*GetServiceTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetServiceToken not implemented")
}
func (*UnimplementedStatsServer) ListAccountsStats(context.Context, *ListStatsRequest) (*ListStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAccountsStats not implemented")
}
func (*UnimplementedStatsServer) AddAccountsStats(context.Context, *AddStatsRequest) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAccountsStats not implemented")
}
func (*UnimplementedStatsServer) ListNewsStats(context.Context, *ListStatsRequest) (*ListStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNewsStats not implemented")
}
func (*UnimplementedStatsServer) AddNewsStats(context.Context, *AddStatsRequest) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNewsStats not implemented")
}
func (*UnimplementedStatsServer) ListCommentsStats(context.Context, *ListStatsRequest) (*ListStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCommentsStats not implemented")
}
func (*UnimplementedStatsServer) AddCommentsStats(context.Context, *AddStatsRequest) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCommentsStats not implemented")
}

func RegisterStatsServer(s *grpc.Server, srv StatsServer) {
	s.RegisterService(&_Stats_serviceDesc, srv)
}

func _Stats_GetServiceToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetServiceTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsServer).GetServiceToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stats.Stats/GetServiceToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatsServer).GetServiceToken(ctx, req.(*GetServiceTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stats_ListAccountsStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsServer).ListAccountsStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stats.Stats/ListAccountsStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatsServer).ListAccountsStats(ctx, req.(*ListStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stats_AddAccountsStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsServer).AddAccountsStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stats.Stats/AddAccountsStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatsServer).AddAccountsStats(ctx, req.(*AddStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stats_ListNewsStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsServer).ListNewsStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stats.Stats/ListNewsStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatsServer).ListNewsStats(ctx, req.(*ListStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stats_AddNewsStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsServer).AddNewsStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stats.Stats/AddNewsStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatsServer).AddNewsStats(ctx, req.(*AddStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stats_ListCommentsStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsServer).ListCommentsStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stats.Stats/ListCommentsStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatsServer).ListCommentsStats(ctx, req.(*ListStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stats_AddCommentsStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsServer).AddCommentsStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stats.Stats/AddCommentsStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatsServer).AddCommentsStats(ctx, req.(*AddStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Stats_serviceDesc = grpc.ServiceDesc{
	ServiceName: "stats.Stats",
	HandlerType: (*StatsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetServiceToken",
			Handler:    _Stats_GetServiceToken_Handler,
		},
		{
			MethodName: "ListAccountsStats",
			Handler:    _Stats_ListAccountsStats_Handler,
		},
		{
			MethodName: "AddAccountsStats",
			Handler:    _Stats_AddAccountsStats_Handler,
		},
		{
			MethodName: "ListNewsStats",
			Handler:    _Stats_ListNewsStats_Handler,
		},
		{
			MethodName: "AddNewsStats",
			Handler:    _Stats_AddNewsStats_Handler,
		},
		{
			MethodName: "ListCommentsStats",
			Handler:    _Stats_ListCommentsStats_Handler,
		},
		{
			MethodName: "AddCommentsStats",
			Handler:    _Stats_AddCommentsStats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/stats/proto/stats.proto",
}
