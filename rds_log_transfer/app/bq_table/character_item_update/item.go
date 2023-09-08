package character_item_update

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{"before_used_flag", "after_used_flag", "before_locked", "after_locked"}

type CharacterItemUpdateItem struct {
	UserId                *string               `json:"user_id,omitempty"`
	CharacterId           *string               `json:"character_id,omitempty"`
	CharacterClass        *interface{}          `json:"character_class,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	ItemType              *json.Number          `json:"item_type,omitempty"`
	ItemId                *json.Number          `json:"item_id,omitempty"`
	ItemIdName            *string               `json:"item_id__name,omitempty"`
	UniqueId              *string               `json:"unique_id,omitempty"`
	BeforeUsedFlag        *bool                 `json:"before_used_flag,omitempty"`
	AfterUsedFlag         *bool                 `json:"after_used_flag,omitempty"`
	BeforeLocked          *bool                 `json:"before_locked,omitempty"`
	AfterLocked           *bool                 `json:"after_locked,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i CharacterItemUpdateItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                 i.Record.UserId,
		"character_id":            i.Record.CharacterId,
		"character_class":         sb_log.EncodeJsonValue(i.CharacterClass),
		"request_header_location": sb_log.GetRequestHeaderLocationStr(i.Record),
		"item_type":               i.ItemType,
		"item_id":                 i.Record.ItemId,
		"item_id__name":           i.ItemIdName,
		"unique_id":               i.Record.UniqueId,
		"before_used_flag":        i.BeforeUsedFlag,
		"after_used_flag":         i.AfterUsedFlag,
		"before_locked":           i.BeforeLocked,
		"after_locked":            i.AfterLocked,
		"log_time":                i.LogTime,
		"uuid":                    i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*CharacterItemUpdateItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item CharacterItemUpdateItem
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
