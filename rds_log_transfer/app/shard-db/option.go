package shard_db

import (
	"fmt"
	"rds_log_transfer/config"
	"strconv"
)

// ShardDbOption シャードDB接続情報
type ShardDbOption struct {
	HostWriter string
	HostReader string
	Database   string
	Username   string
	Password   string
}

type DsnOption struct {
	IsWriter  bool
	ParseTime bool
}

func (o *ShardDbOption) GetDSN(option DsnOption) string {
	var host string
	if option.IsWriter {
		host = o.HostWriter
	} else {
		host = o.HostReader
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", o.Username, o.Password, host, o.Database)
	if option.ParseTime {
		dsn = dsn + "?parseTime=true"
	}

	return dsn
}

var optionLoaded = false
var shardDbOptions []ShardDbOption

func GetShardDbOptions() (*[]ShardDbOption, error) {
	if optionLoaded {
		return &shardDbOptions, nil
	}

	// シャードDBの数取得
	dbShardNum, err := config.GetShardDBNum()
	if err != nil {
		return &shardDbOptions, err
	}

	// シャードDBの接続情報リスト
	shardDbOptions = make([]ShardDbOption, dbShardNum)
	for i := 0; i < dbShardNum; i++ {
		suffix := strconv.Itoa(i)
		wHost, ok := config.LoadDatastoreEnv("DB_HOST" + suffix)
		if !ok || len(wHost) == 0 {
			return nil, fmt.Errorf("%s is not defined in env", "DB_HOST"+suffix)
		}
		rHost, ok := config.LoadDatastoreEnv("DB_SLAVE_HOST" + suffix)
		if !ok || len(rHost) == 0 {
			return nil, fmt.Errorf("%s is not defined in env", "DB_SLAVE_HOST"+suffix)
		}
		db, ok := config.LoadDatastoreEnv("DB_DATABASE" + suffix)
		if !ok || len(db) == 0 {
			return nil, fmt.Errorf("%s is not defined in env", "DB_DATABASE"+suffix)
		}
		user, ok := config.LoadDatastoreEnv("DB_USERNAME")
		if !ok || len(user) == 0 {
			return nil, fmt.Errorf("%s is not defined in env", "DB_USERNAME"+suffix)
		}
		pass, ok := config.LoadDatastoreEnv("DB_PASSWORD")
		if !ok || len(pass) == 0 {
			return nil, fmt.Errorf("%s is not defined in env", "DB_PASSWORD"+suffix)
		}

		shardDbOptions[i] = ShardDbOption{
			HostWriter: wHost,
			HostReader: rHost,
			Database:   db,
			Username:   user,
			Password:   pass,
		}
	}

	optionLoaded = true
	return &shardDbOptions, nil
}

func GetShardDbOption(index int) (*ShardDbOption, error) {
	options, err := GetShardDbOptions()
	if err != nil {
		return nil, err
	}
	if index >= len(*options) {
		return nil, fmt.Errorf("index[%d] out of bounds options", index)
	}
	return &(*options)[index], nil
}
