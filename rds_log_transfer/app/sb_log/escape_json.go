package sb_log

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	shard_db "rds_log_transfer/shard-db"
	"regexp"
	"strconv"
	"strings"
)

// EscapeJsonValue sb_log.logのJsonをパースできるように置換する
func EscapeJsonValue(log string, boolKeys *[]string) string {
	// Bool bool型なのに0,1の数値で記録されている値を置き換える
	for _, key := range *boolKeys {
		log = replaceBooleanNum(key, log)
	}
	return log
}

func escapeJsonString(key string, value string) string {
	reg, err := regexp.Compile(fmt.Sprintf("\"%s\":(\\{.+\\})([,\\}])", key))
	if err != nil {
		return value
	}
	result := reg.FindStringSubmatch(value)
	if len(result) > 2 {
		matched := result[1]
		replaceMatched := strings.ReplaceAll(matched, "\"", "\\\"")
		value = strings.Replace(value,
			fmt.Sprintf("\"%s\":%s", key, matched),
			fmt.Sprintf("\"%s\":\"%s\"", key, replaceMatched), 1)
	}

	return value
}

func replaceBooleanNum(key string, value string) string {
	reg, err := regexp.Compile(fmt.Sprintf("\"%s\":([01])", key))
	if err != nil {
		return value
	}
	result := reg.FindStringSubmatch(value)
	if len(result) > 1 {
		matched := result[1]
		var replace string
		if matched == "0" {
			replace = "false"
		} else {
			replace = "true"
		}
		value = strings.Replace(value,
			fmt.Sprintf("\"%s\":%s", key, matched),
			fmt.Sprintf("\"%s\":%s", key, replace), 1)
	}

	return value
}

func replaceTimestampString(key string, value string) string {
	reg, err := regexp.Compile(fmt.Sprintf("\"%s\":(\"\\d{4}.+? UTC\")", key))
	if err != nil {
		return value
	}
	result := reg.FindStringSubmatch(value)
	if len(result) > 1 {
		matched := result[1]
		replaceMatched := strings.ReplaceAll(matched, " UTC", "Z")
		replaceMatched = strings.Replace(replaceMatched, " ", "T", 1)
		value = strings.Replace(value,
			fmt.Sprintf("\"%s\":%s", key, matched),
			fmt.Sprintf("\"%s\":%s", key, replaceMatched), 1)
	}

	return value
}

func dropUnknownKey(key string, value string) string {
	reg, err := regexp.Compile(fmt.Sprintf("\"%s\":(\".+?\")", key))
	if err != nil {
		return value
	}
	result := reg.FindStringSubmatch(value)
	if len(result) > 1 {
		matched := result[1]
		value = strings.Replace(value,
			fmt.Sprintf(",\"%s\":%s", key, matched),
			"", 1)
	}

	return value
}

func GetRequestHeaderLocationStr(record *shard_db.SbLogRecord) string {
	rhl := JsonRequestHeaderLocation{
		SessionId: record.SessionId,
		MapId:     record.MapId,
		PosX:      record.PosX,
		PosY:      record.PosY,
		PosZ:      record.PosZ,
	}
	rhlJson, err := json.Marshal(rhl)
	if err != nil {
		return ""
	}

	return string(rhlJson)
}

// GetUUID seed(shard-db index)とレコードのIDからUUIDを生成する
func GetUUID(seed int, recordId uint64) string {
	hashBase := strconv.Itoa(seed) + strconv.FormatUint(recordId, 10)
	hash := sha256.Sum256([]byte(hashBase))
	return hex.EncodeToString(hash[:])
}

func EncodeJsonValue(v interface{}) string {
	obj, err := json.Marshal(v)
	if err != nil {
		fmt.Print("marshal failed\n")
	}
	return string(obj)
}
