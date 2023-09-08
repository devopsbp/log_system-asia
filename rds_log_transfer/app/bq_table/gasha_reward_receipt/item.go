package gasha_reward_receipt

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type GashaRewardReceiptItem struct {
	UserId                *string               `json:"user_id,omitempty"`
	CharacterId           *string               `json:"character_id,omitempty"`
	Gender                *json.Number          `json:"gender,omitempty"`
	AdventurerRank        *json.Number          `json:"adventurer_rank,omitempty"`
	Platform              *json.Number          `json:"platform,omitempty"`
	GashaId               *string               `json:"gasha_id,omitempty"`
	GashaIdName           *string               `json:"gasha_id__name,omitempty"`
	GashaName             *string               `json:"gasha_name,omitempty"`
	RewardItems           *interface{}          `json:"reward_items,omitempty"`
	BonusItems            *interface{}          `json:"bonus_items,omitempty"`
	ChangeAmountPaid      *json.Number          `json:"change_amount_paid,omitempty"`
	ChangeAmountFree      *json.Number          `json:"change_amount_free,omitempty"`
	ChangeAmountCoin      *json.Number          `json:"change_amount_coin,omitempty"`
	ChangeTokenId         *json.Number          `json:"change_token_id,omitempty"`
	ChangeTokenIdName     *string               `json:"change_token_id__name,omitempty"`
	ChangeAmountToken     *json.Number          `json:"change_amount_token,omitempty"`
	GashaCount            *json.Number          `json:"gasha_count,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i GashaRewardReceiptItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                 i.Record.UserId,
		"character_id":            i.Record.CharacterId,
		"gender":                  i.Gender,
		"adventurer_rank":         i.AdventurerRank,
		"platform":                i.Platform,
		"gasha_id":                i.GashaId,
		"gasha_id__name":          i.GashaIdName,
		"gasha_name":              i.GashaName,
		"reward_items":            sb_log.EncodeJsonValue(i.RewardItems),
		"bonus_items":             sb_log.EncodeJsonValue(i.BonusItems),
		"change_amount_paid":      i.ChangeAmountPaid,
		"change_amount_free":      i.ChangeAmountFree,
		"change_amount_coin":      i.ChangeAmountCoin,
		"change_token_id":         i.ChangeTokenId,
		"change_token_id__name":   i.ChangeTokenIdName,
		"change_amount_token":     i.ChangeAmountToken,
		"gasha_count":             i.GashaCount,
		"request_header_location": sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                i.LogTime,
		"uuid":                    i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*GashaRewardReceiptItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item GashaRewardReceiptItem
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
