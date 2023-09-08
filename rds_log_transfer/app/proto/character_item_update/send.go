package character_item_update

import (
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"rds_log_transfer/bq"
	"rds_log_transfer/sb_log"
	shard_db "rds_log_transfer/shard-db"
)

var jsonKeys = []string{"character_class", "request_header_location"}
var boolKeys = []string{"before_used_flag", "after_used_flag", "before_locked", "after_locked"}
var timestampKeys = []string{}

func SendBqMessage(records *[]*shard_db.SbLogRecord) error {
	var messages = make([]*CharacterItemUpdateMsg, len(*records))
	for _, record := range *records {
		var message CharacterItemUpdateMsg
		// 値のJSONをエスケープ
		jsonEscapedLog := sb_log.EscapeJsonValue(record.Log, &jsonKeys, &boolKeys, &timestampKeys)
		err := protojson.Unmarshal([]byte(jsonEscapedLog), &message)
		if err != nil {
			fmt.Printf("Unmarshal failed. %v\n", jsonEscapedLog)
			return err
		}
		messages = append(messages, &message)
	}

	err := bq.Send("SB_LOG_character_item_update", &messages)
	if err != nil {
		return err
	}

	return nil
}
