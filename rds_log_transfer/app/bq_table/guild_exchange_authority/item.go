package guild_exchange_authority

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type GuildExchangeAuthorityItem struct {
	UserId                *string               `json:"user_id,omitempty"`
	CharacterId           *string               `json:"character_id,omitempty"`
	CharacterName         *string               `json:"character_name,omitempty"`
	ToUserId              *string               `json:"to_user_id,omitempty"`
	ToCharacterId         *string               `json:"to_character_id,omitempty"`
	ToCharacterName       *string               `json:"to_character_name,omitempty"`
	GuildId               *string               `json:"guild_id,omitempty"`
	GuildName             *string               `json:"guild_name,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i GuildExchangeAuthorityItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                 i.Record.UserId,
		"character_id":            i.Record.CharacterId,
		"character_name":          i.CharacterName,
		"to_user_id":              i.ToUserId,
		"to_character_id":         i.ToCharacterId,
		"to_character_name":       i.ToCharacterName,
		"guild_id":                i.GuildId,
		"guild_name":              i.GuildName,
		"request_header_location": sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                i.LogTime,
		"uuid":                    i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*GuildExchangeAuthorityItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item GuildExchangeAuthorityItem
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
