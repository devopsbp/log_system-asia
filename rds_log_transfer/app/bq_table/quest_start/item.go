package quest_start

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type QuestStartItem struct {
	UserId                 *string               `json:"user_id,omitempty"`
	CharacterId            *string               `json:"character_id,omitempty"`
	ClassType              *json.Number          `json:"class_type,omitempty"`
	Level                  *json.Number          `json:"level,omitempty"`
	RequestHeaderLocation  *interface{}          `json:"request_header_location,omitempty"`
	QuestId                *json.Number          `json:"quest_id,omitempty"`
	QuestName              *string               `json:"quest_name,omitempty"`
	QuestUniqueId          *string               `json:"quest_unique_id,omitempty"`
	AutoOrderQuestUniqueId *string               `json:"auto_order_quest_unique_id,omitempty"`
	LogTime                *string               `json:"log_time,omitempty"`
	Uuid                   *string               `json:"-"`
	Record                 *shard_db.SbLogRecord `json:"-"`
}

func (i QuestStartItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                    i.Record.UserId,
		"character_id":               i.Record.CharacterId,
		"class_type":                 i.ClassType,
		"level":                      i.Level,
		"request_header_location":    sb_log.GetRequestHeaderLocationStr(i.Record),
		"quest_id":                   i.QuestId,
		"quest_name":                 i.QuestName,
		"quest_unique_id":            i.QuestUniqueId,
		"auto_order_quest_unique_id": i.AutoOrderQuestUniqueId,
		"log_time":                   i.LogTime,
		"uuid":                       i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*QuestStartItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item QuestStartItem
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
