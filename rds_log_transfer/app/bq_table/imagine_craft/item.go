package imagine_craft

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type ImagineCraftItem struct {
	UserId                *string               `json:"user_id,omitempty"`
	CharacterId           *string               `json:"character_id,omitempty"`
	CharacterClass        *interface{}          `json:"character_class,omitempty"`
	UniqueId              *string               `json:"unique_id,omitempty"`
	ImagineId             *json.Number          `json:"imagine_id,omitempty"`
	ImagineIdName         *string               `json:"imagine_id__name,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	PerkId                *json.Number          `json:"perk_id,omitempty"`
	PerkIdName            *string               `json:"perk_id__name,omitempty"`
	PerkValue             *json.Number          `json:"perk_value,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i ImagineCraftItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                 i.Record.UserId,
		"character_id":            i.Record.CharacterId,
		"character_class":         sb_log.EncodeJsonValue(i.CharacterClass),
		"unique_id":               i.Record.UniqueId,
		"imagine_id":              i.ImagineId,
		"imagine_id__name":        i.ImagineIdName,
		"request_header_location": sb_log.GetRequestHeaderLocationStr(i.Record),
		"perk_id":                 i.PerkId,
		"perk_id__name":           i.PerkIdName,
		"perk_value":              i.PerkValue,
		"log_time":                i.LogTime,
		"uuid":                    i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*ImagineCraftItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item ImagineCraftItem
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
