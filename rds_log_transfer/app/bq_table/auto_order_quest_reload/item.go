package auto_order_quest_reload

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type AutoOrderQuestReloadItem struct {
	UserId                 *string               `json:"user_id,omitempty"`
	CharacterId            *string               `json:"character_id,omitempty"`
	CharacterLevel         *json.Number          `json:"character_level,omitempty"`
	QuestAutoOrderType     *json.Number          `json:"quest_auto_order_type,omitempty"`
	SeasonId               *json.Number          `json:"season_id,omitempty"`
	SeasonIdName           *string               `json:"season_id__name,omitempty"`
	QuestReloadId          *string               `json:"quest_reload_id,omitempty"`
	TicketId               *json.Number          `json:"ticket_id,omitempty"`
	TicketCount            *json.Number          `json:"ticket_count,omitempty"`
	RequestHeaderLocation  *interface{}          `json:"request_header_location,omitempty"`
	QuestUniqueId          *string               `json:"quest_unique_id,omitempty"`
	AutoOrderQuestUniqueId *string               `json:"auto_order_quest_unique_id,omitempty"`
	LogTime                *string               `json:"log_time,omitempty"`
	Uuid                   *string               `json:"-"`
	Record                 *shard_db.SbLogRecord `json:"-"`
}

func (i AutoOrderQuestReloadItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                    i.Record.UserId,
		"character_id":               i.Record.CharacterId,
		"character_level":            i.CharacterLevel,
		"quest_auto_order_type":      i.QuestAutoOrderType,
		"season_id":                  i.SeasonId,
		"season_id__name":            i.SeasonIdName,
		"quest_reload_id":            i.QuestReloadId,
		"ticket_id":                  i.TicketId,
		"ticket_count":               i.TicketCount,
		"request_header_location":    sb_log.GetRequestHeaderLocationStr(i.Record),
		"quest_unique_id":            i.QuestUniqueId,
		"auto_order_quest_unique_id": i.AutoOrderQuestUniqueId,
		"log_time":                   i.LogTime,
		"uuid":                       i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*AutoOrderQuestReloadItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item AutoOrderQuestReloadItem
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
