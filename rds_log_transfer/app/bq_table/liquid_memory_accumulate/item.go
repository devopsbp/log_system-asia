package liquid_memory_accumulate

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type LiquidMemoryAccumulateItem struct {
	UserId                 *string               `json:"user_id,omitempty"`
	CharacterId            *string               `json:"character_id,omitempty"`
	CharacterClass         *interface{}          `json:"character_class,omitempty"`
	LiquidMemoryId         *json.Number          `json:"liquid_memory_id,omitempty"`
	LiquidMemoryType       *json.Number          `json:"liquid_memory_type,omitempty"`
	LiquidMemoryLevel      *json.Number          `json:"liquid_memory_level,omitempty"`
	BeforeAccumulateAmount *json.Number          `json:"before_accumulate_amount,omitempty"`
	AfterAccumulatedAmount *json.Number          `json:"after_accumulated_amount,omitempty"`
	AddAccumulatedAmount   *json.Number          `json:"add_accumulated_amount,omitempty"`
	RequestHeaderLocation  *interface{}          `json:"request_header_location,omitempty"`
	LogTime                *string               `json:"log_time,omitempty"`
	Uuid                   *string               `json:"-"`
	Record                 *shard_db.SbLogRecord `json:"-"`
}

func (i LiquidMemoryAccumulateItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                  i.Record.UserId,
		"character_id":             i.Record.CharacterId,
		"character_class":          sb_log.EncodeJsonValue(i.CharacterClass),
		"liquid_memory_id":         i.LiquidMemoryId,
		"liquid_memory_type":       i.LiquidMemoryType,
		"liquid_memory_level":      i.LiquidMemoryLevel,
		"before_accumulate_amount": i.BeforeAccumulateAmount,
		"after_accumulated_amount": i.AfterAccumulatedAmount,
		"add_accumulated_amount":   i.AddAccumulatedAmount,
		"request_header_location":  sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                 i.LogTime,
		"uuid":                     i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*LiquidMemoryAccumulateItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item LiquidMemoryAccumulateItem
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
