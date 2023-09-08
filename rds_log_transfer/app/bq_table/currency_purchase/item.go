package currency_purchase

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{"first_purchase_for_payment"}

type CurrencyPurchaseItem struct {
	UserId                  *string               `json:"user_id,omitempty"`
	CharacterId             *string               `json:"character_id,omitempty"`
	AgencyId                *json.Number          `json:"agency_id,omitempty"`
	PaymentType             *string               `json:"payment_type,omitempty"`
	OrderId                 *string               `json:"order_id,omitempty"`
	ItemId                  *json.Number          `json:"item_id,omitempty"`
	ItemName                *string               `json:"item_name,omitempty"`
	Count                   *json.Number          `json:"count,omitempty"`
	Price                   *json.Number          `json:"price,omitempty"`
	Step                    *json.Number          `json:"step,omitempty"`
	Url                     *string               `json:"url,omitempty"`
	RequestDate             *string               `json:"request_date,omitempty"`
	RequestParams           *string               `json:"request_params,omitempty"`
	ResultCode              *json.Number          `json:"result_code,omitempty"`
	AgencyResultCode        *string               `json:"agency_result_code,omitempty"`
	AgencyResultMessage     *string               `json:"agency_result_message,omitempty"`
	ResponseDate            *string               `json:"response_date,omitempty"`
	StatusCode              *json.Number          `json:"status_code,omitempty"`
	Result                  *string               `json:"result,omitempty"`
	UnitPrice               *json.Number          `json:"unit_price,omitempty"`
	Currency                *string               `json:"currency,omitempty"`
	FirstPurchaseForPayment *bool                 `json:"first_purchase_for_payment,omitempty"`
	RequestHeaderLocation   *interface{}          `json:"request_header_location,omitempty"`
	LogTime                 *string               `json:"log_time,omitempty"`
	Uuid                    *string               `json:"-"`
	Record                  *shard_db.SbLogRecord `json:"-"`
}

func (i CurrencyPurchaseItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                    i.Record.UserId,
		"character_id":               i.Record.CharacterId,
		"agency_id":                  i.AgencyId,
		"payment_type":               i.PaymentType,
		"order_id":                   i.OrderId,
		"item_id":                    i.Record.ItemId,
		"item_name":                  i.ItemName,
		"count":                      i.Count,
		"price":                      i.Price,
		"step":                       i.Step,
		"url":                        i.Url,
		"request_date":               i.RequestDate,
		"request_params":             i.RequestParams,
		"result_code":                i.ResultCode,
		"agency_result_code":         i.AgencyResultCode,
		"agency_result_message":      i.AgencyResultMessage,
		"response_date":              i.ResponseDate,
		"status_code":                i.StatusCode,
		"result":                     i.Result,
		"unit_price":                 i.UnitPrice,
		"currency":                   i.Currency,
		"first_purchase_for_payment": i.FirstPurchaseForPayment,
		"request_header_location":    sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                   i.LogTime,
		"uuid":                       i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*CurrencyPurchaseItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item CurrencyPurchaseItem
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
