package character_class_add_exp

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type CharacterClassAddExpItem struct {
	UserId                *string               `json:"user_id,omitempty"`
	CharacterId           *string               `json:"character_id,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	ClassType             *json.Number          `json:"class_type,omitempty"`
	RequestExp            *json.Number          `json:"request_exp,omitempty"`
	AddExp                *json.Number          `json:"add_exp,omitempty"`
	BeforeExp             *json.Number          `json:"before_exp,omitempty"`
	AfterExp              *json.Number          `json:"after_exp,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i CharacterClassAddExpItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                 i.Record.UserId,
		"character_id":            i.Record.CharacterId,
		"request_header_location": sb_log.GetRequestHeaderLocationStr(i.Record),
		"class_type":              i.ClassType,
		"request_exp":             i.RequestExp,
		"add_exp":                 i.AddExp,
		"before_exp":              i.BeforeExp,
		"after_exp":               i.AfterExp,
		"log_time":                i.LogTime,
		"uuid":                    i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*CharacterClassAddExpItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item CharacterClassAddExpItem
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
