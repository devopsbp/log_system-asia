package aesthe_update_character

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type AestheUpdateCharacterItem struct {
	UserId                *string               `json:"user_id,omitempty"`
	CharacterId           *string               `json:"character_id,omitempty"`
	CourseId              *json.Number          `json:"course_id,omitempty"`
	CourseIdName          *string               `json:"course_id__name,omitempty"`
	TokenId               *json.Number          `json:"token_id,omitempty"`
	TokenIdName           *string               `json:"token_id__name,omitempty"`
	TokenAmount           *json.Number          `json:"token_amount,omitempty"`
	CharacterCreate       *interface{}          `json:"character_create,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i AestheUpdateCharacterItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                 i.Record.UserId,
		"character_id":            i.Record.CharacterId,
		"course_id":               i.CourseId,
		"course_id__name":         i.CourseIdName,
		"token_id":                i.TokenId,
		"token_id__name":          i.TokenIdName,
		"token_amount":            i.TokenAmount,
		"character_create":        sb_log.EncodeJsonValue(i.CharacterCreate),
		"request_header_location": sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                i.LogTime,
		"uuid":                    i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*AestheUpdateCharacterItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item AestheUpdateCharacterItem
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
