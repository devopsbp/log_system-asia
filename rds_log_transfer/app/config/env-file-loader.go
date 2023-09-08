package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"strconv"
)

var datastoreEnv, _ = godotenv.Read(".env/datastore.env")
var bqBatchEnv, _ = godotenv.Read(".env/bq-batch.env")

func LoadDatastoreEnv(key string) (string, bool) {
	val, ok := datastoreEnv[key]
	return val, ok
}

func LoadBqBatchEnv(key string) (string, bool) {
	val, ok := bqBatchEnv[key]
	return val, ok
}

// GetShardDBNum シャードDBの数取得
func GetShardDBNum() (int, error) {
	dbShardNumValue, envOk := LoadDatastoreEnv("DB_SHARD_NUM")
	if !envOk {
		return 0, fmt.Errorf("DB_SHARD_NUM is not defined")
	}
	dbShardNum, err := strconv.Atoi(dbShardNumValue)
	if err != nil { // invalid value
		return 0, err
	}
	return dbShardNum, nil
}
