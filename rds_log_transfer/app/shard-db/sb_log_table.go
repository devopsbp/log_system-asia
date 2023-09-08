package shard_db

import (
	"database/sql"
	"strings"
)

// sb_logテーブルの定義

const SbLogTableName = "sb_log"

type SbLogRecord struct {
	Id          uint64
	LogTime     string
	LogTag      string
	UserId      string
	CharacterId string
	ItemId      int
	UniqueId    string
	Log         string
	SessionId   string
	MapId       string
	PosX        string
	PosY        string
	PosZ        string
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}

// GetBQTableName BQテーブル名取得
func (rec *SbLogRecord) GetBQTableName() string {
	// remove sb_log prefix
	tableName := strings.Replace(rec.LogTag, "SB_LOG_", "", 1)
	return tableName
}
