// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: weapon_fusion.proto

package weapon_fusion

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

type WeaponFusionMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId                *string                `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	CharacterId           *string                `protobuf:"bytes,2,opt,name=character_id,json=characterId" json:"character_id,omitempty"`
	CharacterClass        *string                `protobuf:"bytes,3,opt,name=character_class,json=characterClass" json:"character_class,omitempty"`
	UniqueId              *string                `protobuf:"bytes,4,opt,name=unique_id,json=uniqueId" json:"unique_id,omitempty"`
	ItemId                *int32                 `protobuf:"zigzag32,5,opt,name=item_id,json=itemId" json:"item_id,omitempty"`
	ItemId_Name           *string                `protobuf:"bytes,6,opt,name=item_id__name,json=itemIdName" json:"item_id__name,omitempty"`
	PerkId1               *int32                 `protobuf:"zigzag32,7,opt,name=perk_id1,json=perkId1" json:"perk_id1,omitempty"`
	EffectValues1         *string                `protobuf:"bytes,8,opt,name=effect_values1,json=effectValues1" json:"effect_values1,omitempty"`
	PerkId2               *int32                 `protobuf:"zigzag32,9,opt,name=perk_id2,json=perkId2" json:"perk_id2,omitempty"`
	EffectValues2         *string                `protobuf:"bytes,10,opt,name=effect_values2,json=effectValues2" json:"effect_values2,omitempty"`
	PerkId3               *int32                 `protobuf:"zigzag32,11,opt,name=perk_id3,json=perkId3" json:"perk_id3,omitempty"`
	EffectValues3         *string                `protobuf:"bytes,12,opt,name=effect_values3,json=effectValues3" json:"effect_values3,omitempty"`
	PerkId4               *int32                 `protobuf:"zigzag32,13,opt,name=perk_id4,json=perkId4" json:"perk_id4,omitempty"`
	EffectValues4         *string                `protobuf:"bytes,14,opt,name=effect_values4,json=effectValues4" json:"effect_values4,omitempty"`
	PresentUnlockedSlot   *int32                 `protobuf:"zigzag32,15,opt,name=present_unlocked_slot,json=presentUnlockedSlot" json:"present_unlocked_slot,omitempty"`
	UsedItemUniqueId      *int32                 `protobuf:"zigzag32,16,opt,name=used_item_unique_id,json=usedItemUniqueId" json:"used_item_unique_id,omitempty"`
	UsedItemId            *int32                 `protobuf:"zigzag32,17,opt,name=used_item_id,json=usedItemId" json:"used_item_id,omitempty"`
	UsedSupportItemId1    *int32                 `protobuf:"zigzag32,18,opt,name=used_support_item_id1,json=usedSupportItemId1" json:"used_support_item_id1,omitempty"`
	UsedSupportItemId2    *int32                 `protobuf:"zigzag32,19,opt,name=used_support_item_id2,json=usedSupportItemId2" json:"used_support_item_id2,omitempty"`
	FusionSlotNo          *int32                 `protobuf:"zigzag32,20,opt,name=fusion_slot_no,json=fusionSlotNo" json:"fusion_slot_no,omitempty"`
	ResultPerkId          *int32                 `protobuf:"zigzag32,21,opt,name=result_perk_id,json=resultPerkId" json:"result_perk_id,omitempty"`
	ResultEffectValues    *string                `protobuf:"bytes,22,opt,name=result_effect_values,json=resultEffectValues" json:"result_effect_values,omitempty"`
	RequestHeaderLocation *string                `protobuf:"bytes,23,opt,name=request_header_location,json=requestHeaderLocation" json:"request_header_location,omitempty"`
	LogTime               *timestamppb.Timestamp `protobuf:"bytes,24,opt,name=log_time,json=logTime" json:"log_time,omitempty"`
	Uuid                  *string                `protobuf:"bytes,25,opt,name=uuid" json:"uuid,omitempty"`
}

func (x *WeaponFusionMsg) Reset() {
	*x = WeaponFusionMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_weapon_fusion_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WeaponFusionMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WeaponFusionMsg) ProtoMessage() {}

func (x *WeaponFusionMsg) ProtoReflect() protoreflect.Message {
	mi := &file_weapon_fusion_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WeaponFusionMsg.ProtoReflect.Descriptor instead.
func (*WeaponFusionMsg) Descriptor() ([]byte, []int) {
	return file_weapon_fusion_proto_rawDescGZIP(), []int{0}
}

func (x *WeaponFusionMsg) GetUserId() string {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return ""
}

func (x *WeaponFusionMsg) GetCharacterId() string {
	if x != nil && x.CharacterId != nil {
		return *x.CharacterId
	}
	return ""
}

func (x *WeaponFusionMsg) GetCharacterClass() string {
	if x != nil && x.CharacterClass != nil {
		return *x.CharacterClass
	}
	return ""
}

func (x *WeaponFusionMsg) GetUniqueId() string {
	if x != nil && x.UniqueId != nil {
		return *x.UniqueId
	}
	return ""
}

func (x *WeaponFusionMsg) GetItemId() int32 {
	if x != nil && x.ItemId != nil {
		return *x.ItemId
	}
	return 0
}

func (x *WeaponFusionMsg) GetItemId_Name() string {
	if x != nil && x.ItemId_Name != nil {
		return *x.ItemId_Name
	}
	return ""
}

func (x *WeaponFusionMsg) GetPerkId1() int32 {
	if x != nil && x.PerkId1 != nil {
		return *x.PerkId1
	}
	return 0
}

func (x *WeaponFusionMsg) GetEffectValues1() string {
	if x != nil && x.EffectValues1 != nil {
		return *x.EffectValues1
	}
	return ""
}

func (x *WeaponFusionMsg) GetPerkId2() int32 {
	if x != nil && x.PerkId2 != nil {
		return *x.PerkId2
	}
	return 0
}

func (x *WeaponFusionMsg) GetEffectValues2() string {
	if x != nil && x.EffectValues2 != nil {
		return *x.EffectValues2
	}
	return ""
}

func (x *WeaponFusionMsg) GetPerkId3() int32 {
	if x != nil && x.PerkId3 != nil {
		return *x.PerkId3
	}
	return 0
}

func (x *WeaponFusionMsg) GetEffectValues3() string {
	if x != nil && x.EffectValues3 != nil {
		return *x.EffectValues3
	}
	return ""
}

func (x *WeaponFusionMsg) GetPerkId4() int32 {
	if x != nil && x.PerkId4 != nil {
		return *x.PerkId4
	}
	return 0
}

func (x *WeaponFusionMsg) GetEffectValues4() string {
	if x != nil && x.EffectValues4 != nil {
		return *x.EffectValues4
	}
	return ""
}

func (x *WeaponFusionMsg) GetPresentUnlockedSlot() int32 {
	if x != nil && x.PresentUnlockedSlot != nil {
		return *x.PresentUnlockedSlot
	}
	return 0
}

func (x *WeaponFusionMsg) GetUsedItemUniqueId() int32 {
	if x != nil && x.UsedItemUniqueId != nil {
		return *x.UsedItemUniqueId
	}
	return 0
}

func (x *WeaponFusionMsg) GetUsedItemId() int32 {
	if x != nil && x.UsedItemId != nil {
		return *x.UsedItemId
	}
	return 0
}

func (x *WeaponFusionMsg) GetUsedSupportItemId1() int32 {
	if x != nil && x.UsedSupportItemId1 != nil {
		return *x.UsedSupportItemId1
	}
	return 0
}

func (x *WeaponFusionMsg) GetUsedSupportItemId2() int32 {
	if x != nil && x.UsedSupportItemId2 != nil {
		return *x.UsedSupportItemId2
	}
	return 0
}

func (x *WeaponFusionMsg) GetFusionSlotNo() int32 {
	if x != nil && x.FusionSlotNo != nil {
		return *x.FusionSlotNo
	}
	return 0
}

func (x *WeaponFusionMsg) GetResultPerkId() int32 {
	if x != nil && x.ResultPerkId != nil {
		return *x.ResultPerkId
	}
	return 0
}

func (x *WeaponFusionMsg) GetResultEffectValues() string {
	if x != nil && x.ResultEffectValues != nil {
		return *x.ResultEffectValues
	}
	return ""
}

func (x *WeaponFusionMsg) GetRequestHeaderLocation() string {
	if x != nil && x.RequestHeaderLocation != nil {
		return *x.RequestHeaderLocation
	}
	return ""
}

func (x *WeaponFusionMsg) GetLogTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LogTime
	}
	return nil
}

func (x *WeaponFusionMsg) GetUuid() string {
	if x != nil && x.Uuid != nil {
		return *x.Uuid
	}
	return ""
}

var File_weapon_fusion_proto protoreflect.FileDescriptor

var file_weapon_fusion_proto_rawDesc = []byte{
	0x0a, 0x13, 0x77, 0x65, 0x61, 0x70, 0x6f, 0x6e, 0x5f, 0x66, 0x75, 0x73, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x77, 0x65, 0x61, 0x70, 0x6f, 0x6e, 0x5f, 0x66, 0x75,
	0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc3, 0x07, 0x0a, 0x0f, 0x57, 0x65, 0x61, 0x70, 0x6f, 0x6e,
	0x46, 0x75, 0x73, 0x69, 0x6f, 0x6e, 0x4d, 0x73, 0x67, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63,
	0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74,
	0x65, 0x72, 0x5f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e,
	0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x12, 0x1b,
	0x0a, 0x09, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x69,
	0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x11, 0x52, 0x06, 0x69, 0x74,
	0x65, 0x6d, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0d, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x5f,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x74, 0x65,
	0x6d, 0x49, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x65, 0x72, 0x6b, 0x5f,
	0x69, 0x64, 0x31, 0x18, 0x07, 0x20, 0x01, 0x28, 0x11, 0x52, 0x07, 0x70, 0x65, 0x72, 0x6b, 0x49,
	0x64, 0x31, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x73, 0x31, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65, 0x66, 0x66, 0x65,
	0x63, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x31, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x65, 0x72,
	0x6b, 0x5f, 0x69, 0x64, 0x32, 0x18, 0x09, 0x20, 0x01, 0x28, 0x11, 0x52, 0x07, 0x70, 0x65, 0x72,
	0x6b, 0x49, 0x64, 0x32, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x5f, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x73, 0x32, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65, 0x66,
	0x66, 0x65, 0x63, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x32, 0x12, 0x19, 0x0a, 0x08, 0x70,
	0x65, 0x72, 0x6b, 0x5f, 0x69, 0x64, 0x33, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x11, 0x52, 0x07, 0x70,
	0x65, 0x72, 0x6b, 0x49, 0x64, 0x33, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74,
	0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x33, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x33, 0x12, 0x19, 0x0a,
	0x08, 0x70, 0x65, 0x72, 0x6b, 0x5f, 0x69, 0x64, 0x34, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x11, 0x52,
	0x07, 0x70, 0x65, 0x72, 0x6b, 0x49, 0x64, 0x34, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x66, 0x66, 0x65,
	0x63, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x34, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x34, 0x12,
	0x32, 0x0a, 0x15, 0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x5f, 0x75, 0x6e, 0x6c, 0x6f, 0x63,
	0x6b, 0x65, 0x64, 0x5f, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x11, 0x52, 0x13,
	0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x53,
	0x6c, 0x6f, 0x74, 0x12, 0x2d, 0x0a, 0x13, 0x75, 0x73, 0x65, 0x64, 0x5f, 0x69, 0x74, 0x65, 0x6d,
	0x5f, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x10, 0x20, 0x01, 0x28, 0x11,
	0x52, 0x10, 0x75, 0x73, 0x65, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x55, 0x6e, 0x69, 0x71, 0x75, 0x65,
	0x49, 0x64, 0x12, 0x20, 0x0a, 0x0c, 0x75, 0x73, 0x65, 0x64, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f,
	0x69, 0x64, 0x18, 0x11, 0x20, 0x01, 0x28, 0x11, 0x52, 0x0a, 0x75, 0x73, 0x65, 0x64, 0x49, 0x74,
	0x65, 0x6d, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x15, 0x75, 0x73, 0x65, 0x64, 0x5f, 0x73, 0x75, 0x70,
	0x70, 0x6f, 0x72, 0x74, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x31, 0x18, 0x12, 0x20,
	0x01, 0x28, 0x11, 0x52, 0x12, 0x75, 0x73, 0x65, 0x64, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74,
	0x49, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x31, 0x12, 0x31, 0x0a, 0x15, 0x75, 0x73, 0x65, 0x64, 0x5f,
	0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x32,
	0x18, 0x13, 0x20, 0x01, 0x28, 0x11, 0x52, 0x12, 0x75, 0x73, 0x65, 0x64, 0x53, 0x75, 0x70, 0x70,
	0x6f, 0x72, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x32, 0x12, 0x24, 0x0a, 0x0e, 0x66, 0x75,
	0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x6c, 0x6f, 0x74, 0x5f, 0x6e, 0x6f, 0x18, 0x14, 0x20, 0x01,
	0x28, 0x11, 0x52, 0x0c, 0x66, 0x75, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x6c, 0x6f, 0x74, 0x4e, 0x6f,
	0x12, 0x24, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x70, 0x65, 0x72, 0x6b, 0x5f,
	0x69, 0x64, 0x18, 0x15, 0x20, 0x01, 0x28, 0x11, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x50, 0x65, 0x72, 0x6b, 0x49, 0x64, 0x12, 0x30, 0x0a, 0x14, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x5f, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x16,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x45, 0x66, 0x66, 0x65,
	0x63, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x36, 0x0a, 0x17, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x17, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x35, 0x0a, 0x08, 0x6c, 0x6f, 0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x18, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07,
	0x6c, 0x6f, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x19, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x42, 0x15, 0x5a, 0x13, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x65, 0x61, 0x70, 0x6f, 0x6e, 0x5f, 0x66, 0x75, 0x73, 0x69,
	0x6f, 0x6e,
}

var (
	file_weapon_fusion_proto_rawDescOnce sync.Once
	file_weapon_fusion_proto_rawDescData = file_weapon_fusion_proto_rawDesc
)

func file_weapon_fusion_proto_rawDescGZIP() []byte {
	file_weapon_fusion_proto_rawDescOnce.Do(func() {
		file_weapon_fusion_proto_rawDescData = protoimpl.X.CompressGZIP(file_weapon_fusion_proto_rawDescData)
	})
	return file_weapon_fusion_proto_rawDescData
}

var file_weapon_fusion_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_weapon_fusion_proto_goTypes = []interface{}{
	(*WeaponFusionMsg)(nil),       // 0: weapon_fusion.WeaponFusionMsg
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_weapon_fusion_proto_depIdxs = []int32{
	1, // 0: weapon_fusion.WeaponFusionMsg.log_time:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_weapon_fusion_proto_init() }
func file_weapon_fusion_proto_init() {
	if File_weapon_fusion_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_weapon_fusion_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WeaponFusionMsg); i {
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
			RawDescriptor: file_weapon_fusion_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_weapon_fusion_proto_goTypes,
		DependencyIndexes: file_weapon_fusion_proto_depIdxs,
		MessageInfos:      file_weapon_fusion_proto_msgTypes,
	}.Build()
	File_weapon_fusion_proto = out.File
	file_weapon_fusion_proto_rawDesc = nil
	file_weapon_fusion_proto_goTypes = nil
	file_weapon_fusion_proto_depIdxs = nil
}
