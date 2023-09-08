package skill_acquisition

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type SkillAcquisitionItem struct {
	UserId                *string               `json:"user_id,omitempty"`
	CharacterId           *string               `json:"character_id,omitempty"`
	CharacterClass        *interface{}          `json:"character_class,omitempty"`
	SkillId               *json.Number          `json:"skill_id,omitempty"`
	SkillIdName           *string               `json:"skill_id__name,omitempty"`
	SkillLevel            *json.Number          `json:"skill_level,omitempty"`
	ClassType             *json.Number          `json:"class_type,omitempty"`
	SkillType             *json.Number          `json:"skill_type,omitempty"`
	SkillPointRemain      *json.Number          `json:"skill_point_remain,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i SkillAcquisitionItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                 i.Record.UserId,
		"character_id":            i.Record.CharacterId,
		"character_class":         sb_log.EncodeJsonValue(i.CharacterClass),
		"skill_id":                i.SkillId,
		"skill_id__name":          i.SkillIdName,
		"skill_level":             i.SkillLevel,
		"class_type":              i.ClassType,
		"skill_type":              i.SkillType,
		"skill_point_remain":      i.SkillPointRemain,
		"request_header_location": sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                i.LogTime,
		"uuid":                    i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*SkillAcquisitionItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item SkillAcquisitionItem
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
