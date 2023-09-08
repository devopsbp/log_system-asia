package bq

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"rds_log_transfer/config"
)

type ClientWithContext struct {
	Client *bigquery.Client
	Ctx    context.Context
}

func GetClientWithContext() (*ClientWithContext, error) {
	projectId, ok := config.LoadBqBatchEnv("SB_LOG_BQ_GCP_PROJECT_ID")
	if !ok {
		return nil, fmt.Errorf("don't load LoadBqBatchEnv")
	}
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectId, option.WithCredentialsFile(config.GoogleAccountJsonPath))
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %w", err)
	}

	clientWithContext := ClientWithContext{client, ctx}
	return &clientWithContext, nil
}
