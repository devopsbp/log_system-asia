// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: soul_rank_up.proto

package soul_rank_up

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

type SoulRankUpMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId                *string                `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	CharacterId           *string                `protobuf:"bytes,2,opt,name=character_id,json=characterId" json:"character_id,omitempty"`
	AdventurerRank        *int32                 `protobuf:"zigzag32,3,opt,name=adventurer_rank,json=adventurerRank" json:"adventurer_rank,omitempty"`
	ClassLevel            *int32                 `protobuf:"zigzag32,4,opt,name=class_level,json=classLevel" json:"class_level,omitempty"`
	CharacterClass        *string                `protobuf:"bytes,5,opt,name=character_class,json=characterClass" json:"character_class,omitempty"`
	RequestHeaderLocation *string                `protobuf:"bytes,6,opt,name=request_header_location,json=requestHeaderLocation" json:"request_header_location,omitempty"`
	LogTime               *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=log_time,json=logTime" json:"log_time,omitempty"`
	Uuid                  *string                `protobuf:"bytes,8,opt,name=uuid" json:"uuid,omitempty"`
}

func (x *SoulRankUpMsg) Reset() {
	*x = SoulRankUpMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_soul_rank_up_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SoulRankUpMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SoulRankUpMsg) ProtoMessage() {}

func (x *SoulRankUpMsg) ProtoReflect() protoreflect.Message {
	mi := &file_soul_rank_up_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SoulRankUpMsg.ProtoReflect.Descriptor instead.
func (*SoulRankUpMsg) Descriptor() ([]byte, []int) {
	return file_soul_rank_up_proto_rawDescGZIP(), []int{0}
}

func (x *SoulRankUpMsg) GetUserId() string {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return ""
}

func (x *SoulRankUpMsg) GetCharacterId() string {
	if x != nil && x.CharacterId != nil {
		return *x.CharacterId
	}
	return ""
}

func (x *SoulRankUpMsg) GetAdventurerRank() int32 {
	if x != nil && x.AdventurerRank != nil {
		return *x.AdventurerRank
	}
	return 0
}

func (x *SoulRankUpMsg) GetClassLevel() int32 {
	if x != nil && x.ClassLevel != nil {
		return *x.ClassLevel
	}
	return 0
}

func (x *SoulRankUpMsg) GetCharacterClass() string {
	if x != nil && x.CharacterClass != nil {
		return *x.CharacterClass
	}
	return ""
}

func (x *SoulRankUpMsg) GetRequestHeaderLocation() string {
	if x != nil && x.RequestHeaderLocation != nil {
		return *x.RequestHeaderLocation
	}
	return ""
}

func (x *SoulRankUpMsg) GetLogTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LogTime
	}
	return nil
}

func (x *SoulRankUpMsg) GetUuid() string {
	if x != nil && x.Uuid != nil {
		return *x.Uuid
	}
	return ""
}

var File_soul_rank_up_proto protoreflect.FileDescriptor

var file_soul_rank_up_proto_rawDesc = []byte{
	0x0a, 0x12, 0x73, 0x6f, 0x75, 0x6c, 0x5f, 0x72, 0x61, 0x6e, 0x6b, 0x5f, 0x75, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x73, 0x6f, 0x75, 0x6c, 0x5f, 0x72, 0x61, 0x6e, 0x6b, 0x5f,
	0x75, 0x70, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xc1, 0x02, 0x0a, 0x0d, 0x53, 0x6f, 0x75, 0x6c, 0x52, 0x61, 0x6e, 0x6b,
	0x55, 0x70, 0x4d, 0x73, 0x67, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x21,
	0x0a, 0x0c, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x27, 0x0a, 0x0f, 0x61, 0x64, 0x76, 0x65, 0x6e, 0x74, 0x75, 0x72, 0x65, 0x72, 0x5f,
	0x72, 0x61, 0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x11, 0x52, 0x0e, 0x61, 0x64, 0x76, 0x65,
	0x6e, 0x74, 0x75, 0x72, 0x65, 0x72, 0x52, 0x61, 0x6e, 0x6b, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6c,
	0x61, 0x73, 0x73, 0x5f, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x11, 0x52,
	0x0a, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x27, 0x0a, 0x0f, 0x63,
	0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x43,
	0x6c, 0x61, 0x73, 0x73, 0x12, 0x36, 0x0a, 0x17, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x35, 0x0a, 0x08,
	0x6c, 0x6f, 0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x6c, 0x6f, 0x67, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x42, 0x14, 0x5a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x73, 0x6f, 0x75, 0x6c, 0x5f, 0x72, 0x61, 0x6e, 0x6b, 0x5f, 0x75, 0x70,
}

var (
	file_soul_rank_up_proto_rawDescOnce sync.Once
	file_soul_rank_up_proto_rawDescData = file_soul_rank_up_proto_rawDesc
)

func file_soul_rank_up_proto_rawDescGZIP() []byte {
	file_soul_rank_up_proto_rawDescOnce.Do(func() {
		file_soul_rank_up_proto_rawDescData = protoimpl.X.CompressGZIP(file_soul_rank_up_proto_rawDescData)
	})
	return file_soul_rank_up_proto_rawDescData
}

var file_soul_rank_up_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_soul_rank_up_proto_goTypes = []interface{}{
	(*SoulRankUpMsg)(nil),         // 0: soul_rank_up.SoulRankUpMsg
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_soul_rank_up_proto_depIdxs = []int32{
	1, // 0: soul_rank_up.SoulRankUpMsg.log_time:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_soul_rank_up_proto_init() }
func file_soul_rank_up_proto_init() {
	if File_soul_rank_up_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_soul_rank_up_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SoulRankUpMsg); i {
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
			RawDescriptor: file_soul_rank_up_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_soul_rank_up_proto_goTypes,
		DependencyIndexes: file_soul_rank_up_proto_depIdxs,
		MessageInfos:      file_soul_rank_up_proto_msgTypes,
	}.Build()
	File_soul_rank_up_proto = out.File
	file_soul_rank_up_proto_rawDesc = nil
	file_soul_rank_up_proto_goTypes = nil
	file_soul_rank_up_proto_depIdxs = nil
}
