package weapon_craft

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{"is_great_success"}

type WeaponCraftItem struct {
	UserId                *string               `json:"user_id,omitempty"`
	CharacterId           *string               `json:"character_id,omitempty"`
	CharacterClass        *interface{}          `json:"character_class,omitempty"`
	Level                 *json.Number          `json:"level,omitempty"`
	ItemId                *json.Number          `json:"item_id,omitempty"`
	ItemIdName            *string               `json:"item_id__name,omitempty"`
	AbilityUnlock         *json.Number          `json:"ability_unlock,omitempty"`
	UniqueId              *string               `json:"unique_id,omitempty"`
	VitalType             *json.Number          `json:"vital_type,omitempty"`
	VitalValue            *json.Number          `json:"vital_value,omitempty"`
	IsGreatSuccess        *bool                 `json:"is_great_success,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i WeaponCraftItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                 i.Record.UserId,
		"character_id":            i.Record.CharacterId,
		"character_class":         sb_log.EncodeJsonValue(i.CharacterClass),
		"level":                   i.Level,
		"item_id":                 i.Record.ItemId,
		"item_id__name":           i.ItemIdName,
		"ability_unlock":          i.AbilityUnlock,
		"unique_id":               i.Record.UniqueId,
		"vital_type":              i.VitalType,
		"vital_value":             i.VitalValue,
		"is_great_success":        i.IsGreatSuccess,
		"request_header_location": sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                i.LogTime,
		"uuid":                    i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*WeaponCraftItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item WeaponCraftItem
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
