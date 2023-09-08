package save_appearance_weapon_sticker

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var boolKeys = []string{}

type SaveAppearanceWeaponStickerItem struct {
	UserId                              *string               `json:"user_id,omitempty"`
	CharacterId                         *string               `json:"character_id,omitempty"`
	WeaponId                            *json.Number          `json:"weapon_id,omitempty"`
	WeaponIdName                        *string               `json:"weapon_id__name,omitempty"`
	UniqueId                            *string               `json:"unique_id,omitempty"`
	WeaponAppearanceWeaponId            *json.Number          `json:"weapon_appearance_weapon_id,omitempty"`
	WeaponAppearanceWeaponIdName        *string               `json:"weapon_appearance_weapon_id__name,omitempty"`
	WeaponAppearanceWeaponStickerId     *json.Number          `json:"weapon_appearance_weapon_sticker_id,omitempty"`
	WeaponAppearanceWeaponStickerIdName *string               `json:"weapon_appearance_weapon_sticker_id__name,omitempty"`
	StickerType                         *json.Number          `json:"sticker_type,omitempty"`
	StickerReleaseTime                  *string               `json:"sticker_release_time,omitempty"`
	RequestHeaderLocation               *interface{}          `json:"request_header_location,omitempty"`
	LogTime                             *string               `json:"log_time,omitempty"`
	Uuid                                *string               `json:"-"`
	Record                              *shard_db.SbLogRecord `json:"-"`
}

func (i SaveAppearanceWeaponStickerItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id":                                   i.Record.UserId,
		"character_id":                              i.Record.CharacterId,
		"weapon_id":                                 i.WeaponId,
		"weapon_id__name":                           i.WeaponIdName,
		"unique_id":                                 i.Record.UniqueId,
		"weapon_appearance_weapon_id":               i.WeaponAppearanceWeaponId,
		"weapon_appearance_weapon_id__name":         i.WeaponAppearanceWeaponIdName,
		"weapon_appearance_weapon_sticker_id":       i.WeaponAppearanceWeaponStickerId,
		"weapon_appearance_weapon_sticker_id__name": i.WeaponAppearanceWeaponStickerIdName,
		"sticker_type":                              i.StickerType,
		"sticker_release_time":                      i.StickerReleaseTime,
		"request_header_location":                   sb_log.GetRequestHeaderLocationStr(i.Record),
		"log_time":                                  i.LogTime,
		"uuid":                                      i.Uuid,
	}, bigquery.NoDedupeID, nil
}

func InsertBqItems(clientWithContext *bq.ClientWithContext, records *[]*shard_db.SbLogRecord, uuidSeed int) ([]*shard_db.SbLogRecord, error) {
	var tableName string
	var items = make([]*SaveAppearanceWeaponStickerItem, 0)
	var failedRecords = make([]*shard_db.SbLogRecord, 0)
	for _, record := range *records {
		if tableName == "" {
			tableName = record.LogTag
		}
		var item SaveAppearanceWeaponStickerItem
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
