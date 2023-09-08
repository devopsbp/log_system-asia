// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: defense_battle.proto

package defense_battle

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

type DefenseBattleMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId                      *string                `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	CharacterId                 *string                `protobuf:"bytes,2,opt,name=character_id,json=characterId" json:"character_id,omitempty"`
	CharacterClass              *string                `protobuf:"bytes,3,opt,name=character_class,json=characterClass" json:"character_class,omitempty"`
	CharacterClassLevel         *int32                 `protobuf:"zigzag32,4,opt,name=character_class_level,json=characterClassLevel" json:"character_class_level,omitempty"`
	UnalteredContributionScores *string                `protobuf:"bytes,5,opt,name=unaltered_contribution_scores,json=unalteredContributionScores" json:"unaltered_contribution_scores,omitempty"`
	CharacterBattleScore        *float32               `protobuf:"fixed32,6,opt,name=character_battle_score,json=characterBattleScore" json:"character_battle_score,omitempty"`
	EquippedWeapons             *string                `protobuf:"bytes,7,opt,name=equipped_weapons,json=equippedWeapons" json:"equipped_weapons,omitempty"`
	EquippedBattleImagines      *string                `protobuf:"bytes,8,opt,name=equipped_battle_imagines,json=equippedBattleImagines" json:"equipped_battle_imagines,omitempty"`
	EquippedInnerImagines       *string                `protobuf:"bytes,9,opt,name=equipped_inner_imagines,json=equippedInnerImagines" json:"equipped_inner_imagines,omitempty"`
	EquippedTacticalSkills      *string                `protobuf:"bytes,10,opt,name=equipped_tactical_skills,json=equippedTacticalSkills" json:"equipped_tactical_skills,omitempty"`
	EquippedTacticalAbilities   *string                `protobuf:"bytes,11,opt,name=equipped_tactical_abilities,json=equippedTacticalAbilities" json:"equipped_tactical_abilities,omitempty"`
	EquippedClassSkills         *string                `protobuf:"bytes,12,opt,name=equipped_class_skills,json=equippedClassSkills" json:"equipped_class_skills,omitempty"`
	CharacterDamage             *float32               `protobuf:"fixed32,13,opt,name=character_damage,json=characterDamage" json:"character_damage,omitempty"`
	NumberOfDeaths              *int32                 `protobuf:"zigzag32,14,opt,name=number_of_deaths,json=numberOfDeaths" json:"number_of_deaths,omitempty"`
	NumberOfEnemiesDefeated     *int32                 `protobuf:"zigzag32,15,opt,name=number_of_enemies_defeated,json=numberOfEnemiesDefeated" json:"number_of_enemies_defeated,omitempty"`
	RequestHeaderLocation       *string                `protobuf:"bytes,16,opt,name=request_header_location,json=requestHeaderLocation" json:"request_header_location,omitempty"`
	LogTime                     *timestamppb.Timestamp `protobuf:"bytes,17,opt,name=log_time,json=logTime" json:"log_time,omitempty"`
	Uuid                        *string                `protobuf:"bytes,18,opt,name=uuid" json:"uuid,omitempty"`
}

func (x *DefenseBattleMsg) Reset() {
	*x = DefenseBattleMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_defense_battle_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DefenseBattleMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DefenseBattleMsg) ProtoMessage() {}

func (x *DefenseBattleMsg) ProtoReflect() protoreflect.Message {
	mi := &file_defense_battle_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DefenseBattleMsg.ProtoReflect.Descriptor instead.
func (*DefenseBattleMsg) Descriptor() ([]byte, []int) {
	return file_defense_battle_proto_rawDescGZIP(), []int{0}
}

func (x *DefenseBattleMsg) GetUserId() string {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return ""
}

func (x *DefenseBattleMsg) GetCharacterId() string {
	if x != nil && x.CharacterId != nil {
		return *x.CharacterId
	}
	return ""
}

func (x *DefenseBattleMsg) GetCharacterClass() string {
	if x != nil && x.CharacterClass != nil {
		return *x.CharacterClass
	}
	return ""
}

func (x *DefenseBattleMsg) GetCharacterClassLevel() int32 {
	if x != nil && x.CharacterClassLevel != nil {
		return *x.CharacterClassLevel
	}
	return 0
}

func (x *DefenseBattleMsg) GetUnalteredContributionScores() string {
	if x != nil && x.UnalteredContributionScores != nil {
		return *x.UnalteredContributionScores
	}
	return ""
}

func (x *DefenseBattleMsg) GetCharacterBattleScore() float32 {
	if x != nil && x.CharacterBattleScore != nil {
		return *x.CharacterBattleScore
	}
	return 0
}

func (x *DefenseBattleMsg) GetEquippedWeapons() string {
	if x != nil && x.EquippedWeapons != nil {
		return *x.EquippedWeapons
	}
	return ""
}

func (x *DefenseBattleMsg) GetEquippedBattleImagines() string {
	if x != nil && x.EquippedBattleImagines != nil {
		return *x.EquippedBattleImagines
	}
	return ""
}

func (x *DefenseBattleMsg) GetEquippedInnerImagines() string {
	if x != nil && x.EquippedInnerImagines != nil {
		return *x.EquippedInnerImagines
	}
	return ""
}

func (x *DefenseBattleMsg) GetEquippedTacticalSkills() string {
	if x != nil && x.EquippedTacticalSkills != nil {
		return *x.EquippedTacticalSkills
	}
	return ""
}

func (x *DefenseBattleMsg) GetEquippedTacticalAbilities() string {
	if x != nil && x.EquippedTacticalAbilities != nil {
		return *x.EquippedTacticalAbilities
	}
	return ""
}

func (x *DefenseBattleMsg) GetEquippedClassSkills() string {
	if x != nil && x.EquippedClassSkills != nil {
		return *x.EquippedClassSkills
	}
	return ""
}

func (x *DefenseBattleMsg) GetCharacterDamage() float32 {
	if x != nil && x.CharacterDamage != nil {
		return *x.CharacterDamage
	}
	return 0
}

func (x *DefenseBattleMsg) GetNumberOfDeaths() int32 {
	if x != nil && x.NumberOfDeaths != nil {
		return *x.NumberOfDeaths
	}
	return 0
}

func (x *DefenseBattleMsg) GetNumberOfEnemiesDefeated() int32 {
	if x != nil && x.NumberOfEnemiesDefeated != nil {
		return *x.NumberOfEnemiesDefeated
	}
	return 0
}

func (x *DefenseBattleMsg) GetRequestHeaderLocation() string {
	if x != nil && x.RequestHeaderLocation != nil {
		return *x.RequestHeaderLocation
	}
	return ""
}

func (x *DefenseBattleMsg) GetLogTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LogTime
	}
	return nil
}

func (x *DefenseBattleMsg) GetUuid() string {
	if x != nil && x.Uuid != nil {
		return *x.Uuid
	}
	return ""
}

var File_defense_battle_proto protoreflect.FileDescriptor

var file_defense_battle_proto_rawDesc = []byte{
	0x0a, 0x14, 0x64, 0x65, 0x66, 0x65, 0x6e, 0x73, 0x65, 0x5f, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x64, 0x65, 0x66, 0x65, 0x6e, 0x73, 0x65, 0x5f,
	0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x85, 0x07, 0x0a, 0x10, 0x44, 0x65, 0x66, 0x65,
	0x6e, 0x73, 0x65, 0x42, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x4d, 0x73, 0x67, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x68, 0x61,
	0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x68, 0x61, 0x72,
	0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x43, 0x6c, 0x61, 0x73,
	0x73, 0x12, 0x32, 0x0a, 0x15, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x63,
	0x6c, 0x61, 0x73, 0x73, 0x5f, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x11,
	0x52, 0x13, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x43, 0x6c, 0x61, 0x73, 0x73,
	0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x42, 0x0a, 0x1d, 0x75, 0x6e, 0x61, 0x6c, 0x74, 0x65, 0x72,
	0x65, 0x64, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x73, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1b, 0x75, 0x6e,
	0x61, 0x6c, 0x74, 0x65, 0x72, 0x65, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x12, 0x34, 0x0a, 0x16, 0x63, 0x68, 0x61,
	0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x5f, 0x73, 0x63,
	0x6f, 0x72, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x02, 0x52, 0x14, 0x63, 0x68, 0x61, 0x72, 0x61,
	0x63, 0x74, 0x65, 0x72, 0x42, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x12,
	0x29, 0x0a, 0x10, 0x65, 0x71, 0x75, 0x69, 0x70, 0x70, 0x65, 0x64, 0x5f, 0x77, 0x65, 0x61, 0x70,
	0x6f, 0x6e, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x65, 0x71, 0x75, 0x69, 0x70,
	0x70, 0x65, 0x64, 0x57, 0x65, 0x61, 0x70, 0x6f, 0x6e, 0x73, 0x12, 0x38, 0x0a, 0x18, 0x65, 0x71,
	0x75, 0x69, 0x70, 0x70, 0x65, 0x64, 0x5f, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x5f, 0x69, 0x6d,
	0x61, 0x67, 0x69, 0x6e, 0x65, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x16, 0x65, 0x71,
	0x75, 0x69, 0x70, 0x70, 0x65, 0x64, 0x42, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x49, 0x6d, 0x61, 0x67,
	0x69, 0x6e, 0x65, 0x73, 0x12, 0x36, 0x0a, 0x17, 0x65, 0x71, 0x75, 0x69, 0x70, 0x70, 0x65, 0x64,
	0x5f, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x69, 0x6e, 0x65, 0x73, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x65, 0x71, 0x75, 0x69, 0x70, 0x70, 0x65, 0x64, 0x49,
	0x6e, 0x6e, 0x65, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x69, 0x6e, 0x65, 0x73, 0x12, 0x38, 0x0a, 0x18,
	0x65, 0x71, 0x75, 0x69, 0x70, 0x70, 0x65, 0x64, 0x5f, 0x74, 0x61, 0x63, 0x74, 0x69, 0x63, 0x61,
	0x6c, 0x5f, 0x73, 0x6b, 0x69, 0x6c, 0x6c, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x16,
	0x65, 0x71, 0x75, 0x69, 0x70, 0x70, 0x65, 0x64, 0x54, 0x61, 0x63, 0x74, 0x69, 0x63, 0x61, 0x6c,
	0x53, 0x6b, 0x69, 0x6c, 0x6c, 0x73, 0x12, 0x3e, 0x0a, 0x1b, 0x65, 0x71, 0x75, 0x69, 0x70, 0x70,
	0x65, 0x64, 0x5f, 0x74, 0x61, 0x63, 0x74, 0x69, 0x63, 0x61, 0x6c, 0x5f, 0x61, 0x62, 0x69, 0x6c,
	0x69, 0x74, 0x69, 0x65, 0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x19, 0x65, 0x71, 0x75,
	0x69, 0x70, 0x70, 0x65, 0x64, 0x54, 0x61, 0x63, 0x74, 0x69, 0x63, 0x61, 0x6c, 0x41, 0x62, 0x69,
	0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x15, 0x65, 0x71, 0x75, 0x69, 0x70, 0x70,
	0x65, 0x64, 0x5f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x5f, 0x73, 0x6b, 0x69, 0x6c, 0x6c, 0x73, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x65, 0x71, 0x75, 0x69, 0x70, 0x70, 0x65, 0x64, 0x43,
	0x6c, 0x61, 0x73, 0x73, 0x53, 0x6b, 0x69, 0x6c, 0x6c, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x63, 0x68,
	0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x0f, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x44,
	0x61, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x5f,
	0x6f, 0x66, 0x5f, 0x64, 0x65, 0x61, 0x74, 0x68, 0x73, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x11, 0x52,
	0x0e, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x4f, 0x66, 0x44, 0x65, 0x61, 0x74, 0x68, 0x73, 0x12,
	0x3b, 0x0a, 0x1a, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x6f, 0x66, 0x5f, 0x65, 0x6e, 0x65,
	0x6d, 0x69, 0x65, 0x73, 0x5f, 0x64, 0x65, 0x66, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x0f, 0x20,
	0x01, 0x28, 0x11, 0x52, 0x17, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x4f, 0x66, 0x45, 0x6e, 0x65,
	0x6d, 0x69, 0x65, 0x73, 0x44, 0x65, 0x66, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x36, 0x0a, 0x17,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x6c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x35, 0x0a, 0x08, 0x6c, 0x6f, 0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x11, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x07, 0x6c, 0x6f, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75,
	0x75, 0x69, 0x64, 0x18, 0x12, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x42,
	0x16, 0x5a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x65, 0x66, 0x65, 0x6e, 0x73, 0x65,
	0x5f, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65,
}

var (
	file_defense_battle_proto_rawDescOnce sync.Once
	file_defense_battle_proto_rawDescData = file_defense_battle_proto_rawDesc
)

func file_defense_battle_proto_rawDescGZIP() []byte {
	file_defense_battle_proto_rawDescOnce.Do(func() {
		file_defense_battle_proto_rawDescData = protoimpl.X.CompressGZIP(file_defense_battle_proto_rawDescData)
	})
	return file_defense_battle_proto_rawDescData
}

var file_defense_battle_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_defense_battle_proto_goTypes = []interface{}{
	(*DefenseBattleMsg)(nil),      // 0: defense_battle.DefenseBattleMsg
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_defense_battle_proto_depIdxs = []int32{
	1, // 0: defense_battle.DefenseBattleMsg.log_time:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_defense_battle_proto_init() }
func file_defense_battle_proto_init() {
	if File_defense_battle_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_defense_battle_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DefenseBattleMsg); i {
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
			RawDescriptor: file_defense_battle_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_defense_battle_proto_goTypes,
		DependencyIndexes: file_defense_battle_proto_depIdxs,
		MessageInfos:      file_defense_battle_proto_msgTypes,
	}.Build()
	File_defense_battle_proto = out.File
	file_defense_battle_proto_rawDesc = nil
	file_defense_battle_proto_goTypes = nil
	file_defense_battle_proto_depIdxs = nil
}
