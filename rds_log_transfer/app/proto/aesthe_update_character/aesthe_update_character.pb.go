// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: aesthe_update_character.proto

package aesthe_update_character

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

type AestheUpdateCharacterMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId                *string                `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	CharacterId           *string                `protobuf:"bytes,2,opt,name=character_id,json=characterId" json:"character_id,omitempty"`
	CourseId              *int32                 `protobuf:"zigzag32,3,opt,name=course_id,json=courseId" json:"course_id,omitempty"`
	CourseId_Name         *string                `protobuf:"bytes,4,opt,name=course_id__name,json=courseIdName" json:"course_id__name,omitempty"`
	TokenId               *int32                 `protobuf:"zigzag32,5,opt,name=token_id,json=tokenId" json:"token_id,omitempty"`
	TokenId_Name          *string                `protobuf:"bytes,6,opt,name=token_id__name,json=tokenIdName" json:"token_id__name,omitempty"`
	TokenAmount           *int32                 `protobuf:"zigzag32,7,opt,name=token_amount,json=tokenAmount" json:"token_amount,omitempty"`
	CharacterCreate       *string                `protobuf:"bytes,8,opt,name=character_create,json=characterCreate" json:"character_create,omitempty"`
	RequestHeaderLocation *string                `protobuf:"bytes,9,opt,name=request_header_location,json=requestHeaderLocation" json:"request_header_location,omitempty"`
	LogTime               *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=log_time,json=logTime" json:"log_time,omitempty"`
	Uuid                  *string                `protobuf:"bytes,11,opt,name=uuid" json:"uuid,omitempty"`
}

func (x *AestheUpdateCharacterMsg) Reset() {
	*x = AestheUpdateCharacterMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aesthe_update_character_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AestheUpdateCharacterMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AestheUpdateCharacterMsg) ProtoMessage() {}

func (x *AestheUpdateCharacterMsg) ProtoReflect() protoreflect.Message {
	mi := &file_aesthe_update_character_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AestheUpdateCharacterMsg.ProtoReflect.Descriptor instead.
func (*AestheUpdateCharacterMsg) Descriptor() ([]byte, []int) {
	return file_aesthe_update_character_proto_rawDescGZIP(), []int{0}
}

func (x *AestheUpdateCharacterMsg) GetUserId() string {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return ""
}

func (x *AestheUpdateCharacterMsg) GetCharacterId() string {
	if x != nil && x.CharacterId != nil {
		return *x.CharacterId
	}
	return ""
}

func (x *AestheUpdateCharacterMsg) GetCourseId() int32 {
	if x != nil && x.CourseId != nil {
		return *x.CourseId
	}
	return 0
}

func (x *AestheUpdateCharacterMsg) GetCourseId_Name() string {
	if x != nil && x.CourseId_Name != nil {
		return *x.CourseId_Name
	}
	return ""
}

func (x *AestheUpdateCharacterMsg) GetTokenId() int32 {
	if x != nil && x.TokenId != nil {
		return *x.TokenId
	}
	return 0
}

func (x *AestheUpdateCharacterMsg) GetTokenId_Name() string {
	if x != nil && x.TokenId_Name != nil {
		return *x.TokenId_Name
	}
	return ""
}

func (x *AestheUpdateCharacterMsg) GetTokenAmount() int32 {
	if x != nil && x.TokenAmount != nil {
		return *x.TokenAmount
	}
	return 0
}

func (x *AestheUpdateCharacterMsg) GetCharacterCreate() string {
	if x != nil && x.CharacterCreate != nil {
		return *x.CharacterCreate
	}
	return ""
}

func (x *AestheUpdateCharacterMsg) GetRequestHeaderLocation() string {
	if x != nil && x.RequestHeaderLocation != nil {
		return *x.RequestHeaderLocation
	}
	return ""
}

func (x *AestheUpdateCharacterMsg) GetLogTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LogTime
	}
	return nil
}

func (x *AestheUpdateCharacterMsg) GetUuid() string {
	if x != nil && x.Uuid != nil {
		return *x.Uuid
	}
	return ""
}

var File_aesthe_update_character_proto protoreflect.FileDescriptor

var file_aesthe_update_character_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x61, 0x65, 0x73, 0x74, 0x68, 0x65, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f,
	0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x17, 0x61, 0x65, 0x73, 0x74, 0x68, 0x65, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x63,
	0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xab, 0x03, 0x0a, 0x18, 0x41, 0x65,
	0x73, 0x74, 0x68, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63,
	0x74, 0x65, 0x72, 0x4d, 0x73, 0x67, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6f, 0x75, 0x72, 0x73, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x11, 0x52, 0x08, 0x63, 0x6f, 0x75, 0x72, 0x73, 0x65, 0x49, 0x64, 0x12,
	0x25, 0x0a, 0x0f, 0x63, 0x6f, 0x75, 0x72, 0x73, 0x65, 0x5f, 0x69, 0x64, 0x5f, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6f, 0x75, 0x72, 0x73, 0x65,
	0x49, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x11, 0x52, 0x07, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x49,
	0x64, 0x12, 0x23, 0x0a, 0x0e, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x5f, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x49, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f,
	0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x11, 0x52, 0x0b, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x29, 0x0a, 0x10, 0x63, 0x68, 0x61,
	0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x12, 0x36, 0x0a, 0x17, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x35, 0x0a, 0x08,
	0x6c, 0x6f, 0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x6c, 0x6f, 0x67, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x42, 0x1f, 0x5a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x61, 0x65, 0x73, 0x74, 0x68, 0x65, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x63,
	0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72,
}

var (
	file_aesthe_update_character_proto_rawDescOnce sync.Once
	file_aesthe_update_character_proto_rawDescData = file_aesthe_update_character_proto_rawDesc
)

func file_aesthe_update_character_proto_rawDescGZIP() []byte {
	file_aesthe_update_character_proto_rawDescOnce.Do(func() {
		file_aesthe_update_character_proto_rawDescData = protoimpl.X.CompressGZIP(file_aesthe_update_character_proto_rawDescData)
	})
	return file_aesthe_update_character_proto_rawDescData
}

var file_aesthe_update_character_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_aesthe_update_character_proto_goTypes = []interface{}{
	(*AestheUpdateCharacterMsg)(nil), // 0: aesthe_update_character.AestheUpdateCharacterMsg
	(*timestamppb.Timestamp)(nil),    // 1: google.protobuf.Timestamp
}
var file_aesthe_update_character_proto_depIdxs = []int32{
	1, // 0: aesthe_update_character.AestheUpdateCharacterMsg.log_time:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_aesthe_update_character_proto_init() }
func file_aesthe_update_character_proto_init() {
	if File_aesthe_update_character_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aesthe_update_character_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AestheUpdateCharacterMsg); i {
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
			RawDescriptor: file_aesthe_update_character_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aesthe_update_character_proto_goTypes,
		DependencyIndexes: file_aesthe_update_character_proto_depIdxs,
		MessageInfos:      file_aesthe_update_character_proto_msgTypes,
	}.Build()
	File_aesthe_update_character_proto = out.File
	file_aesthe_update_character_proto_rawDesc = nil
	file_aesthe_update_character_proto_goTypes = nil
	file_aesthe_update_character_proto_depIdxs = nil
}
