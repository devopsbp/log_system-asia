package bq

import (
	"cloud.google.com/go/bigquery"
	"fmt"
	"rds_log_transfer/config"
	"strconv"
	"time"
)

func InsertRows[T bigquery.ValueSaver](clientWithContext *ClientWithContext, tableName string, items *[]*T) error {
	datasetId, ok := config.LoadBqBatchEnv("SB_LOG_BQ_GCP_RDS_DATASET_ID")
	if !ok {
		return fmt.Errorf("don't load LoadBqBatchEnv")
	}

	inserter := clientWithContext.Client.Dataset(datasetId).Table(tableName + tableNameSuffix()).Inserter()

	if err := inserter.Put(clientWithContext.Ctx, *items); err != nil {
		return err
	}
	return nil
}

func tableNameSuffix() string {
	year := time.Now().Year()
	return "_" + strconv.Itoa(year)[0:3]
}
