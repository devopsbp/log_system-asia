// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: character_buy_item.proto

package character_buy_item

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

type CharacterBuyItemMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId                *string                `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	CharacterId           *string                `protobuf:"bytes,2,opt,name=character_id,json=characterId" json:"character_id,omitempty"`
	CharacterClass        *string                `protobuf:"bytes,3,opt,name=character_class,json=characterClass" json:"character_class,omitempty"`
	RequestHeaderLocation *string                `protobuf:"bytes,4,opt,name=request_header_location,json=requestHeaderLocation" json:"request_header_location,omitempty"`
	CounterId             *string                `protobuf:"bytes,5,opt,name=counter_id,json=counterId" json:"counter_id,omitempty"`
	CounterId_Name        *string                `protobuf:"bytes,6,opt,name=counter_id__name,json=counterIdName" json:"counter_id__name,omitempty"`
	ItemList              *string                `protobuf:"bytes,7,opt,name=item_list,json=itemList" json:"item_list,omitempty"`
	UseAmount             *int32                 `protobuf:"zigzag32,8,opt,name=use_amount,json=useAmount" json:"use_amount,omitempty"`
	PresentAmount         *int32                 `protobuf:"zigzag32,9,opt,name=present_amount,json=presentAmount" json:"present_amount,omitempty"`
	TokenId               *int32                 `protobuf:"zigzag32,10,opt,name=token_id,json=tokenId" json:"token_id,omitempty"`
	LogTime               *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=log_time,json=logTime" json:"log_time,omitempty"`
	Uuid                  *string                `protobuf:"bytes,12,opt,name=uuid" json:"uuid,omitempty"`
}

func (x *CharacterBuyItemMsg) Reset() {
	*x = CharacterBuyItemMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_character_buy_item_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CharacterBuyItemMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CharacterBuyItemMsg) ProtoMessage() {}

func (x *CharacterBuyItemMsg) ProtoReflect() protoreflect.Message {
	mi := &file_character_buy_item_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CharacterBuyItemMsg.ProtoReflect.Descriptor instead.
func (*CharacterBuyItemMsg) Descriptor() ([]byte, []int) {
	return file_character_buy_item_proto_rawDescGZIP(), []int{0}
}

func (x *CharacterBuyItemMsg) GetUserId() string {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return ""
}

func (x *CharacterBuyItemMsg) GetCharacterId() string {
	if x != nil && x.CharacterId != nil {
		return *x.CharacterId
	}
	return ""
}

func (x *CharacterBuyItemMsg) GetCharacterClass() string {
	if x != nil && x.CharacterClass != nil {
		return *x.CharacterClass
	}
	return ""
}

func (x *CharacterBuyItemMsg) GetRequestHeaderLocation() string {
	if x != nil && x.RequestHeaderLocation != nil {
		return *x.RequestHeaderLocation
	}
	return ""
}

func (x *CharacterBuyItemMsg) GetCounterId() string {
	if x != nil && x.CounterId != nil {
		return *x.CounterId
	}
	return ""
}

func (x *CharacterBuyItemMsg) GetCounterId_Name() string {
	if x != nil && x.CounterId_Name != nil {
		return *x.CounterId_Name
	}
	return ""
}

func (x *CharacterBuyItemMsg) GetItemList() string {
	if x != nil && x.ItemList != nil {
		return *x.ItemList
	}
	return ""
}

func (x *CharacterBuyItemMsg) GetUseAmount() int32 {
	if x != nil && x.UseAmount != nil {
		return *x.UseAmount
	}
	return 0
}

func (x *CharacterBuyItemMsg) GetPresentAmount() int32 {
	if x != nil && x.PresentAmount != nil {
		return *x.PresentAmount
	}
	return 0
}

func (x *CharacterBuyItemMsg) GetTokenId() int32 {
	if x != nil && x.TokenId != nil {
		return *x.TokenId
	}
	return 0
}

func (x *CharacterBuyItemMsg) GetLogTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LogTime
	}
	return nil
}

func (x *CharacterBuyItemMsg) GetUuid() string {
	if x != nil && x.Uuid != nil {
		return *x.Uuid
	}
	return ""
}

var File_character_buy_item_proto protoreflect.FileDescriptor

var file_character_buy_item_proto_rawDesc = []byte{
	0x0a, 0x18, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x62, 0x75, 0x79, 0x5f,
	0x69, 0x74, 0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x63, 0x68, 0x61, 0x72,
	0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x62, 0x75, 0x79, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xc3, 0x03, 0x0a, 0x13, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x42, 0x75, 0x79,
	0x49, 0x74, 0x65, 0x6d, 0x4d, 0x73, 0x67, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72,
	0x5f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x68,
	0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x12, 0x36, 0x0a, 0x17,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x6c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x10, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x5f, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x49, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x69, 0x74, 0x65, 0x6d, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x69, 0x74, 0x65, 0x6d, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x73, 0x65,
	0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x11, 0x52, 0x09, 0x75,
	0x73, 0x65, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x72, 0x65, 0x73,
	0x65, 0x6e, 0x74, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x11,
	0x52, 0x0d, 0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x19, 0x0a, 0x08, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x11, 0x52, 0x07, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x64, 0x12, 0x35, 0x0a, 0x08, 0x6c, 0x6f,
	0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x6c, 0x6f, 0x67, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x75, 0x75, 0x69, 0x64, 0x42, 0x1a, 0x5a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63,
	0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x62, 0x75, 0x79, 0x5f, 0x69, 0x74, 0x65,
	0x6d,
}

var (
	file_character_buy_item_proto_rawDescOnce sync.Once
	file_character_buy_item_proto_rawDescData = file_character_buy_item_proto_rawDesc
)

func file_character_buy_item_proto_rawDescGZIP() []byte {
	file_character_buy_item_proto_rawDescOnce.Do(func() {
		file_character_buy_item_proto_rawDescData = protoimpl.X.CompressGZIP(file_character_buy_item_proto_rawDescData)
	})
	return file_character_buy_item_proto_rawDescData
}

var file_character_buy_item_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_character_buy_item_proto_goTypes = []interface{}{
	(*CharacterBuyItemMsg)(nil),   // 0: character_buy_item.CharacterBuyItemMsg
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_character_buy_item_proto_depIdxs = []int32{
	1, // 0: character_buy_item.CharacterBuyItemMsg.log_time:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_character_buy_item_proto_init() }
func file_character_buy_item_proto_init() {
	if File_character_buy_item_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_character_buy_item_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CharacterBuyItemMsg); i {
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
			RawDescriptor: file_character_buy_item_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_character_buy_item_proto_goTypes,
		DependencyIndexes: file_character_buy_item_proto_depIdxs,
		MessageInfos:      file_character_buy_item_proto_msgTypes,
	}.Build()
	File_character_buy_item_proto = out.File
	file_character_buy_item_proto_rawDesc = nil
	file_character_buy_item_proto_goTypes = nil
	file_character_buy_item_proto_depIdxs = nil
}
