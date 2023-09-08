package season_pass_rank_up

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type SeasonPassRankUpItem struct {
	UserId                *string               `json:"user_id,omitempty"`
	CharacterId           *string               `json:"character_id,omitempty"`
	CharacterLevel        *json.Number          `json:"character_level,omitempty"`
	SeasonId              *json.Number          `json:"season_id,omitempty"`
	SeasonIdName          *string               `json:"season_id__name,omitempty"`
	BuyDate               *string               `json:"buy_date,omitempty"`
	BuyRank               *json.Number          `json:"buy_rank,omitempty"`
	ResultRank            *json.Number          `json:"result_rank,omitempty"`
	ChangePoints          *json.Number          `json:"change_points,omitempty"`
	ChangedPoints         *json.Number          `json:"changed_points,omitempty"`
	ChangeFree            *json.Number          `json:"change_free,omitempty"`
	ChangePaid            *json.Number          `json:"change_paid,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i SeasonPassRankUpItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                 i.Record.UserId,
		"character_id":            i.Record.CharacterId,
		"character_level":         i.CharacterLevel,
		"season_id":               i.SeasonId,
		"season_id__name":         i.SeasonIdName,
		"buy_date":                i.BuyDate,
		"buy_rank":                i.BuyRank,
		"result_rank":             i.ResultRank,
		"change_points":           i.ChangePoints,
		"changed_points":          i.ChangedPoints,
		"change_free":             i.ChangeFree,
		"change_paid":             i.ChangePaid,
		"request_header_location": sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                i.LogTime,
		"uuid":                    i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*SeasonPassRankUpItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item SeasonPassRankUpItem
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
