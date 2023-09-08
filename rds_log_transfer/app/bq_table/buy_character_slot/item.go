package buy_character_slot

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type BuyCharacterSlotItem struct {
	UserId                  *string               `json:"user_id,omitempty"`
	AdditionalCharacterSlot *json.Number          `json:"additional_character_slot,omitempty"`
	CharacterSlot           *json.Number          `json:"character_slot,omitempty"`
	PurchaseAmount          *json.Number          `json:"purchase_amount,omitempty"`
	RequestHeaderLocation   *interface{}          `json:"request_header_location,omitempty"`
	LogTime                 *string               `json:"log_time,omitempty"`
	Uuid                    *string               `json:"-"`
	Record                  *shard_db.SbLogRecord `json:"-"`
}

func (i BuyCharacterSlotItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                   i.Record.UserId,
		"additional_character_slot": i.AdditionalCharacterSlot,
		"character_slot":            i.CharacterSlot,
		"purchase_amount":           i.PurchaseAmount,
		"request_header_location":   sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                  i.LogTime,
		"uuid":                      i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*BuyCharacterSlotItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item BuyCharacterSlotItem
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
