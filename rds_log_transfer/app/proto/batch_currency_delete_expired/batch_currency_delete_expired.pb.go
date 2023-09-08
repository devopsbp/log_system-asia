// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: batch_currency_delete_expired.proto

package batch_currency_delete_expired

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BatchCurrencyDeleteExpiredMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageLevel          *string                `protobuf:"bytes,1,opt,name=message_level,json=messageLevel" json:"message_level,omitempty"`
	Message               *string                `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	RequestHeaderLocation *string                `protobuf:"bytes,3,opt,name=request_header_location,json=requestHeaderLocation" json:"request_header_location,omitempty"`
	LogTime               *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=log_time,json=logTime" json:"log_time,omitempty"`
	Uuid                  *string                `protobuf:"bytes,5,opt,name=uuid" json:"uuid,omitempty"`
}

func (x *BatchCurrencyDeleteExpiredMsg) Reset() {
	*x = BatchCurrencyDeleteExpiredMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_batch_currency_delete_expired_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchCurrencyDeleteExpiredMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchCurrencyDeleteExpiredMsg) ProtoMessage() {}

func (x *BatchCurrencyDeleteExpiredMsg) ProtoReflect() protoreflect.Message {
	mi := &file_batch_currency_delete_expired_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchCurrencyDeleteExpiredMsg.ProtoReflect.Descriptor instead.
func (*BatchCurrencyDeleteExpiredMsg) Descriptor() ([]byte, []int) {
	return file_batch_currency_delete_expired_proto_rawDescGZIP(), []int{0}
}

func (x *BatchCurrencyDeleteExpiredMsg) GetMessageLevel() string {
	if x != nil && x.MessageLevel != nil {
		return *x.MessageLevel
	}
	return ""
}

func (x *BatchCurrencyDeleteExpiredMsg) GetMessage() string {
	if x != nil && x.Message != nil {
		return *x.Message
	}
	return ""
}

func (x *BatchCurrencyDeleteExpiredMsg) GetRequestHeaderLocation() string {
	if x != nil && x.RequestHeaderLocation != nil {
		return *x.RequestHeaderLocation
	}
	return ""
}

func (x *BatchCurrencyDeleteExpiredMsg) GetLogTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LogTime
	}
	return nil
}

func (x *BatchCurrencyDeleteExpiredMsg) GetUuid() string {
	if x != nil && x.Uuid != nil {
		return *x.Uuid
	}
	return ""
}

var File_batch_currency_delete_expired_proto protoreflect.FileDescriptor

var file_batch_currency_delete_expired_proto_rawDesc = []byte{
	0x0a, 0x23, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79,
	0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1d, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x65, 0x78, 0x70,
	0x69, 0x72, 0x65, 0x64, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe1, 0x01, 0x0a, 0x1d, 0x42, 0x61, 0x74, 0x63, 0x68, 0x43,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x78, 0x70,
	0x69, 0x72, 0x65, 0x64, 0x4d, 0x73, 0x67, 0x12, 0x23, 0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x5f, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x36, 0x0a, 0x17, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x35,
	0x0a, 0x08, 0x6c, 0x6f, 0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x6c, 0x6f,
	0x67, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x42, 0x25, 0x5a, 0x23, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63,
	0x79, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64,
}

var (
	file_batch_currency_delete_expired_proto_rawDescOnce sync.Once
	file_batch_currency_delete_expired_proto_rawDescData = file_batch_currency_delete_expired_proto_rawDesc
)

func file_batch_currency_delete_expired_proto_rawDescGZIP() []byte {
	file_batch_currency_delete_expired_proto_rawDescOnce.Do(func() {
		file_batch_currency_delete_expired_proto_rawDescData = protoimpl.X.CompressGZIP(file_batch_currency_delete_expired_proto_rawDescData)
	})
	return file_batch_currency_delete_expired_proto_rawDescData
}

var file_batch_currency_delete_expired_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_batch_currency_delete_expired_proto_goTypes = []interface{}{
	(*BatchCurrencyDeleteExpiredMsg)(nil), // 0: batch_currency_delete_expired.BatchCurrencyDeleteExpiredMsg
	(*timestamppb.Timestamp)(nil),         // 1: google.protobuf.Timestamp
}
var file_batch_currency_delete_expired_proto_depIdxs = []int32{
	1, // 0: batch_currency_delete_expired.BatchCurrencyDeleteExpiredMsg.log_time:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_batch_currency_delete_expired_proto_init() }
func file_batch_currency_delete_expired_proto_init() {
	if File_batch_currency_delete_expired_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_batch_currency_delete_expired_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchCurrencyDeleteExpiredMsg); i {
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
			RawDescriptor: file_batch_currency_delete_expired_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_batch_currency_delete_expired_proto_goTypes,
		DependencyIndexes: file_batch_currency_delete_expired_proto_depIdxs,
		MessageInfos:      file_batch_currency_delete_expired_proto_msgTypes,
	}.Build()
	File_batch_currency_delete_expired_proto = out.File
	file_batch_currency_delete_expired_proto_rawDesc = nil
	file_batch_currency_delete_expired_proto_goTypes = nil
	file_batch_currency_delete_expired_proto_depIdxs = nil
}
