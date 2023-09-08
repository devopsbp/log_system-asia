// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: guild_create.proto

package guild_create

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

type GuildCreateMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId                *string                `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	CharacterId           *string                `protobuf:"bytes,2,opt,name=character_id,json=characterId" json:"character_id,omitempty"`
	CharacterName         *string                `protobuf:"bytes,3,opt,name=character_name,json=characterName" json:"character_name,omitempty"`
	GuildId               *string                `protobuf:"bytes,4,opt,name=guild_id,json=guildId" json:"guild_id,omitempty"`
	GuildName             *string                `protobuf:"bytes,5,opt,name=guild_name,json=guildName" json:"guild_name,omitempty"`
	GuildTag              *string                `protobuf:"bytes,6,opt,name=guild_tag,json=guildTag" json:"guild_tag,omitempty"`
	Comment               *string                `protobuf:"bytes,7,opt,name=comment" json:"comment,omitempty"`
	Times                 *string                `protobuf:"bytes,8,opt,name=times" json:"times,omitempty"`
	Styles                *string                `protobuf:"bytes,9,opt,name=styles" json:"styles,omitempty"`
	Accepttype            *bool                  `protobuf:"varint,10,opt,name=accepttype" json:"accepttype,omitempty"`
	RequestHeaderLocation *string                `protobuf:"bytes,11,opt,name=request_header_location,json=requestHeaderLocation" json:"request_header_location,omitempty"`
	LogTime               *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=log_time,json=logTime" json:"log_time,omitempty"`
	Uuid                  *string                `protobuf:"bytes,13,opt,name=uuid" json:"uuid,omitempty"`
}

func (x *GuildCreateMsg) Reset() {
	*x = GuildCreateMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_guild_create_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuildCreateMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuildCreateMsg) ProtoMessage() {}

func (x *GuildCreateMsg) ProtoReflect() protoreflect.Message {
	mi := &file_guild_create_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuildCreateMsg.ProtoReflect.Descriptor instead.
func (*GuildCreateMsg) Descriptor() ([]byte, []int) {
	return file_guild_create_proto_rawDescGZIP(), []int{0}
}

func (x *GuildCreateMsg) GetUserId() string {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return ""
}

func (x *GuildCreateMsg) GetCharacterId() string {
	if x != nil && x.CharacterId != nil {
		return *x.CharacterId
	}
	return ""
}

func (x *GuildCreateMsg) GetCharacterName() string {
	if x != nil && x.CharacterName != nil {
		return *x.CharacterName
	}
	return ""
}

func (x *GuildCreateMsg) GetGuildId() string {
	if x != nil && x.GuildId != nil {
		return *x.GuildId
	}
	return ""
}

func (x *GuildCreateMsg) GetGuildName() string {
	if x != nil && x.GuildName != nil {
		return *x.GuildName
	}
	return ""
}

func (x *GuildCreateMsg) GetGuildTag() string {
	if x != nil && x.GuildTag != nil {
		return *x.GuildTag
	}
	return ""
}

func (x *GuildCreateMsg) GetComment() string {
	if x != nil && x.Comment != nil {
		return *x.Comment
	}
	return ""
}

func (x *GuildCreateMsg) GetTimes() string {
	if x != nil && x.Times != nil {
		return *x.Times
	}
	return ""
}

func (x *GuildCreateMsg) GetStyles() string {
	if x != nil && x.Styles != nil {
		return *x.Styles
	}
	return ""
}

func (x *GuildCreateMsg) GetAccepttype() bool {
	if x != nil && x.Accepttype != nil {
		return *x.Accepttype
	}
	return false
}

func (x *GuildCreateMsg) GetRequestHeaderLocation() string {
	if x != nil && x.RequestHeaderLocation != nil {
		return *x.RequestHeaderLocation
	}
	return ""
}

func (x *GuildCreateMsg) GetLogTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LogTime
	}
	return nil
}

func (x *GuildCreateMsg) GetUuid() string {
	if x != nil && x.Uuid != nil {
		return *x.Uuid
	}
	return ""
}

var File_guild_create_proto protoreflect.FileDescriptor

var file_guild_create_proto_rawDesc = []byte{
	0x0a, 0x12, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xb5, 0x03, 0x0a, 0x0e, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x4d, 0x73, 0x67, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x68, 0x61, 0x72,
	0x61, 0x63, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x75, 0x69,
	0x6c, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x75, 0x69,
	0x6c, 0x64, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x74, 0x61, 0x67,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x61, 0x67,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x65,
	0x70, 0x74, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x61, 0x63,
	0x63, 0x65, 0x70, 0x74, 0x74, 0x79, 0x70, 0x65, 0x12, 0x36, 0x0a, 0x17, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x35, 0x0a, 0x08, 0x6c, 0x6f, 0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07,
	0x6c, 0x6f, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x42, 0x14, 0x5a, 0x12, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65,
}

var (
	file_guild_create_proto_rawDescOnce sync.Once
	file_guild_create_proto_rawDescData = file_guild_create_proto_rawDesc
)

func file_guild_create_proto_rawDescGZIP() []byte {
	file_guild_create_proto_rawDescOnce.Do(func() {
		file_guild_create_proto_rawDescData = protoimpl.X.CompressGZIP(file_guild_create_proto_rawDescData)
	})
	return file_guild_create_proto_rawDescData
}

var file_guild_create_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_guild_create_proto_goTypes = []interface{}{
	(*GuildCreateMsg)(nil),        // 0: guild_create.GuildCreateMsg
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_guild_create_proto_depIdxs = []int32{
	1, // 0: guild_create.GuildCreateMsg.log_time:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_guild_create_proto_init() }
func file_guild_create_proto_init() {
	if File_guild_create_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_guild_create_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuildCreateMsg); i {
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
			RawDescriptor: file_guild_create_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_guild_create_proto_goTypes,
		DependencyIndexes: file_guild_create_proto_depIdxs,
		MessageInfos:      file_guild_create_proto_msgTypes,
	}.Build()
	File_guild_create_proto = out.File
	file_guild_create_proto_rawDesc = nil
	file_guild_create_proto_goTypes = nil
	file_guild_create_proto_depIdxs = nil
}
