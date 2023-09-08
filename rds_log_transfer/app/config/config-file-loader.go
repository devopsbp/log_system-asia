package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"strings"
)

const configPath = ".env/bq-batch.config.ini"

var cfg *ini.File = nil

func getBqBatchConfig(section string, key string) (string, error) {
	if cfg == nil {
		configFile, err := ini.Load(configPath)
		if err != nil {
			return "", err
		}
		cfg = configFile
		fmt.Print("config loaded\n")
	}

	value := cfg.Section(section).Key(key).Value()
	return value, nil
}

func GetDatasetIds() ([]string, error) {
	datasetIdsString, err := getBqBatchConfig("gcp", "dataset_ids")
	if err != nil {
		return []string{}, err
	}
	return strings.Split(datasetIdsString, ","), nil
}
