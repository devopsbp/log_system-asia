package shard_db

// sb_log_errorテーブル

// ErrorReason sb_log_errorテーブル エラー事由
type ErrorReason int

const (
	ErrorReasonParseFailed   ErrorReason = 1
	ErrorReasonBigQueryError             = 2
)
