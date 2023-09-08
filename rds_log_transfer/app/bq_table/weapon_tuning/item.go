package weapon_tuning

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type WeaponTuningItem struct {
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
	TuningTicketUniqueId  *string               `json:"tuning_ticket_unique_id,omitempty"`
	TuningTicketId        *json.Number          `json:"tuning_ticket_id,omitempty"`
	TuningTicketCount     *json.Number          `json:"tuning_ticket_count,omitempty"`
	ProtectTicketId1      *json.Number          `json:"protect_ticket_id1,omitempty"`
	ProtectTicketCount1   *json.Number          `json:"protect_ticket_count1,omitempty"`
	ProtectSlotNo1        *json.Number          `json:"protect_slot_no1,omitempty"`
	ProtectPerkId1        *json.Number          `json:"protect_perk_id1,omitempty"`
	ProtectEffectValues1  *string               `json:"protect_effect_values1,omitempty"`
	ProtectTicketId2      *json.Number          `json:"protect_ticket_id2,omitempty"`
	ProtectTicketCount2   *json.Number          `json:"protect_ticket_count2,omitempty"`
	ProtectSlotNo2        *json.Number          `json:"protect_slot_no2,omitempty"`
	ProtectPerkId2        *json.Number          `json:"protect_perk_id2,omitempty"`
	ProtectEffectValues2  *string               `json:"protect_effect_values2,omitempty"`
	ProtectTicketId3      *json.Number          `json:"protect_ticket_id3,omitempty"`
	ProtectTicketCount3   *json.Number          `json:"protect_ticket_count3,omitempty"`
	ProtectSlotNo3        *json.Number          `json:"protect_slot_no3,omitempty"`
	ProtectPerkId3        *json.Number          `json:"protect_perk_id3,omitempty"`
	ProtectEffectValues3  *string               `json:"protect_effect_values3,omitempty"`
	AfterPerkId1          *json.Number          `json:"after_perk_id1,omitempty"`
	AfterEffectValues1    *string               `json:"after_effect_values1,omitempty"`
	AfterPerkId2          *json.Number          `json:"after_perk_id2,omitempty"`
	AfterEffectValues2    *string               `json:"after_effect_values2,omitempty"`
	AfterPerkId3          *json.Number          `json:"after_perk_id3,omitempty"`
	AfterEffectValues3    *string               `json:"after_effect_values3,omitempty"`
	AfterPerkId4          *json.Number          `json:"after_perk_id4,omitempty"`
	AfterEffectValues4    *string               `json:"after_effect_values4,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i WeaponTuningItem) Save() (map[string]bigquery.Value, string, error) {
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
		"tuning_ticket_unique_id": i.TuningTicketUniqueId,
		"tuning_ticket_id":        i.TuningTicketId,
		"tuning_ticket_count":     i.TuningTicketCount,
		"protect_ticket_id1":      i.ProtectTicketId1,
		"protect_ticket_count1":   i.ProtectTicketCount1,
		"protect_slot_no1":        i.ProtectSlotNo1,
		"protect_perk_id1":        i.ProtectPerkId1,
		"protect_effect_values1":  i.ProtectEffectValues1,
		"protect_ticket_id2":      i.ProtectTicketId2,
		"protect_ticket_count2":   i.ProtectTicketCount2,
		"protect_slot_no2":        i.ProtectSlotNo2,
		"protect_perk_id2":        i.ProtectPerkId2,
		"protect_effect_values2":  i.ProtectEffectValues2,
		"protect_ticket_id3":      i.ProtectTicketId3,
		"protect_ticket_count3":   i.ProtectTicketCount3,
		"protect_slot_no3":        i.ProtectSlotNo3,
		"protect_perk_id3":        i.ProtectPerkId3,
		"protect_effect_values3":  i.ProtectEffectValues3,
		"after_perk_id1":          i.AfterPerkId1,
		"after_effect_values1":    i.AfterEffectValues1,
		"after_perk_id2":          i.AfterPerkId2,
		"after_effect_values2":    i.AfterEffectValues2,
		"after_perk_id3":          i.AfterPerkId3,
		"after_effect_values3":    i.AfterEffectValues3,
		"after_perk_id4":          i.AfterPerkId4,
		"after_effect_values4":    i.AfterEffectValues4,
		"request_header_location": sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                i.LogTime,
		"uuid":                    i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*WeaponTuningItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item WeaponTuningItem
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
