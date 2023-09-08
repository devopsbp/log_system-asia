package character_ban

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type CharacterBanItem struct {
	AdminUserId           *string               `json:"admin_user_id,omitempty"`
	TargetUserId          *string               `json:"target_user_id,omitempty"`
	TargetCharacterId     *string               `json:"target_character_id,omitempty"`
	BanType               *json.Number          `json:"ban_type,omitempty"`
	BanType2              *json.Number          `json:"ban_type2,omitempty"`
	RankingExclusion      *json.Number          `json:"ranking_exclusion,omitempty"`
	BanFinishedAt         *string               `json:"ban_finished_at,omitempty"`
	BanDurationType       *json.Number          `json:"ban_duration_type,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i CharacterBanItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"admin_user_id":           i.AdminUserId,
		"target_user_id":          i.TargetUserId,
		"target_character_id":     i.TargetCharacterId,
		"ban_type":                i.BanType,
		"ban_type2":               i.BanType2,
		"ranking_exclusion":       i.RankingExclusion,
		"ban_finished_at":         i.BanFinishedAt,
		"ban_duration_type":       i.BanDurationType,
		"request_header_location": sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                i.LogTime,
		"uuid":                    i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*CharacterBanItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item CharacterBanItem
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
