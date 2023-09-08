package weapon_fusion

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type WeaponFusionItem struct {
	UserId                *string               `json:"user_id,omitempty"`
	CharacterId           *string               `json:"character_id,omitempty"`
	CharacterClass        *interface{}          `json:"character_class,omitempty"`
	UniqueId              *string               `json:"unique_id,omitempty"`
	ItemId                *json.Number          `json:"item_id,omitempty"`
	ItemIdName            *string               `json:"item_id__name,omitempty"`
	PerkId1               *json.Number          `json:"perk_id1,omitempty"`
	EffectValues1         *string               `json:"effect_values1,omitempty"`
	PerkId2               *json.Number          `json:"perk_id2,omitempty"`
	EffectValues2         *string               `json:"effect_values2,omitempty"`
	PerkId3               *json.Number          `json:"perk_id3,omitempty"`
	EffectValues3         *string               `json:"effect_values3,omitempty"`
	PerkId4               *json.Number          `json:"perk_id4,omitempty"`
	EffectValues4         *string               `json:"effect_values4,omitempty"`
	PresentUnlockedSlot   *json.Number          `json:"present_unlocked_slot,omitempty"`
	UsedItemUniqueId      *string               `json:"used_item_unique_id,omitempty"`
	UsedItemId            *json.Number          `json:"used_item_id,omitempty"`
	UsedSupportItemId1    *json.Number          `json:"used_support_item_id1,omitempty"`
	UsedSupportItemId2    *json.Number          `json:"used_support_item_id2,omitempty"`
	FusionSlotNo          *json.Number          `json:"fusion_slot_no,omitempty"`
	ResultPerkId          *json.Number          `json:"result_perk_id,omitempty"`
	ResultEffectValues    *string               `json:"result_effect_values,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i WeaponFusionItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                 i.Record.UserId,
		"character_id":            i.Record.CharacterId,
		"character_class":         sb_log.EncodeJsonValue(i.CharacterClass),
		"unique_id":               i.Record.UniqueId,
		"item_id":                 i.Record.ItemId,
		"item_id__name":           i.ItemIdName,
		"perk_id1":                i.PerkId1,
		"effect_values1":          i.EffectValues1,
		"perk_id2":                i.PerkId2,
		"effect_values2":          i.EffectValues2,
		"perk_id3":                i.PerkId3,
		"effect_values3":          i.EffectValues3,
		"perk_id4":                i.PerkId4,
		"effect_values4":          i.EffectValues4,
		"present_unlocked_slot":   i.PresentUnlockedSlot,
		"used_item_unique_id":     i.UsedItemUniqueId,
		"used_item_id":            i.UsedItemId,
		"used_support_item_id1":   i.UsedSupportItemId1,
		"used_support_item_id2":   i.UsedSupportItemId2,
		"fusion_slot_no":          i.FusionSlotNo,
		"result_perk_id":          i.ResultPerkId,
		"result_effect_values":    i.ResultEffectValues,
		"request_header_location": sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                i.LogTime,
		"uuid":                    i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*WeaponFusionItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item WeaponFusionItem
		item.Record = record
		uuid := sb_log.GetUUID(uuidSeed, record.Id)
		item.Uuid = &uuid
		// escape JSON value
		jsonEscapedLog := sb_log.EscapeJsonValue(record.Log, &boolKeys)
		err := json.Unmarshal([]byte(jsonEscapedLog), &item)
		if err != nil {
			failedRecords = append(failedRecords, record)
		}
		items = append(items, &item)
	}

	err := bq.InsertRows(clientWithContext, tableName, &items)
	if err != nil {
		return failedRecords, err
	}

	return failedRecords, nil
}
