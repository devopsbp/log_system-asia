package currency_history

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{"first_purchase", "is_used_coupon"}

type CurrencyHistoryItem struct {
	UserId                *string               `json:"user_id,omitempty"`
	CharacterId           *string               `json:"character_id,omitempty"`
	Platform              *json.Number          `json:"platform,omitempty"`
	UpdateType            *string               `json:"update_type,omitempty"`
	ChangeAmountPaid      *json.Number          `json:"change_amount_paid,omitempty"`
	ChangeAmountFree      *json.Number          `json:"change_amount_free,omitempty"`
	AfterAmountPaid       *json.Number          `json:"after_amount_paid,omitempty"`
	AfterAmountFree       *json.Number          `json:"after_amount_free,omitempty"`
	ActionType            *string               `json:"action_type,omitempty"`
	LetterId              *string               `json:"letter_id,omitempty"`
	Price                 *json.Number          `json:"price,omitempty"`
	ExpiredAt             *string               `json:"expired_at,omitempty"`
	RequestHeaderLocation *interface{}          `json:"request_header_location,omitempty"`
	Reason                *string               `json:"reason,omitempty"`
	FirstPurchase         *bool                 `json:"first_purchase,omitempty"`
	AmountPaid            *interface{}          `json:"amount_paid,omitempty"`
	AmountPaidCurrency    *string               `json:"amount_paid_currency,omitempty"`
	AddCount              *json.Number          `json:"add_count,omitempty"`
	ItemId                *json.Number          `json:"item_id,omitempty"`
	IsUsedCoupon          *bool                 `json:"is_used_coupon,omitempty"`
	DiscountAmount        *json.Number          `json:"discount_amount,omitempty"`
	LogTime               *string               `json:"log_time,omitempty"`
	Uuid                  *string               `json:"-"`
	Record                *shard_db.SbLogRecord `json:"-"`
}

func (i CurrencyHistoryItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                 i.Record.UserId,
		"character_id":            i.Record.CharacterId,
		"platform":                i.Platform,
		"update_type":             i.UpdateType,
		"change_amount_paid":      i.ChangeAmountPaid,
		"change_amount_free":      i.ChangeAmountFree,
		"after_amount_paid":       i.AfterAmountPaid,
		"after_amount_free":       i.AfterAmountFree,
		"action_type":             i.ActionType,
		"letter_id":               i.LetterId,
		"price":                   i.Price,
		"expired_at":              i.ExpiredAt,
		"request_header_location": sb_log.GetRequestHeaderLocationStr(i.Record),
		"reason":                  i.Reason,
		"first_purchase":          i.FirstPurchase,
		"amount_paid":             sb_log.EncodeJsonValue(i.AmountPaid),
		"amount_paid_currency":    i.AmountPaidCurrency,
		"add_count":               i.AddCount,
		"item_id":                 i.Record.ItemId,
		"is_used_coupon":          i.IsUsedCoupon,
		"discount_amount":         i.DiscountAmount,
		"log_time":                i.LogTime,
		"uuid":                    i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*CurrencyHistoryItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item CurrencyHistoryItem
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
