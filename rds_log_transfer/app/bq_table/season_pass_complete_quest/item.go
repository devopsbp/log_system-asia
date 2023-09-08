package season_pass_complete_quest

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type SeasonPassCompleteQuestItem struct {
	UserId                *string               `json:"user_id,omitempty"`
	CharacterId           *string               `json:"character_id,omitempty"`
	CharacterLevel        *json.Number          `json:"character_level,omitempty"`
	SeasonId              *json.Number          `json:"season_id,omitempty"`
	SeasonIdName          *string               `json:"season_id__name,omitempty"`
	SeasonPassQuestId     *string               `json:"season_pass_quest_id,omitempty"`
	SeasonPassQuestIdName *string               `json:"season_pass_quest_id__name,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i SeasonPassCompleteQuestItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                    i.Record.UserId,
		"character_id":               i.Record.CharacterId,
		"character_level":            i.CharacterLevel,
		"season_id":                  i.SeasonId,
		"season_id__name":            i.SeasonIdName,
		"season_pass_quest_id":       i.SeasonPassQuestId,
		"season_pass_quest_id__name": i.SeasonPassQuestIdName,
		"request_header_location":    sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                   i.LogTime,
		"uuid":                       i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*SeasonPassCompleteQuestItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item SeasonPassCompleteQuestItem
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
