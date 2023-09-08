package bq

import (
	"cloud.google.com/go/bigquery/storage/managedwriter"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/api/option"
	proto2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"rds_log_transfer/config"
)

func Send[T proto.GeneratedMessage](tableName string, msgs *[]T) error {
	projectId, ok := config.LoadBqBatchEnv("SB_LOG_BQ_GCP_PROJECT_ID")
	if !ok {
		return fmt.Errorf("don't load LoadBqBatchEnv")
	}
	datasetId, ok := config.LoadBqBatchEnv("SB_LOG_BQ_GCP_DATASET_ID")
	if !ok {
		return fmt.Errorf("don't load LoadBqBatchEnv")
	}
	ctx := context.Background()
	c, err := managedwriter.NewClient(ctx, projectId, option.WithCredentialsFile(config.GoogleAccountJsonPath))
	if err != nil {
		return err
	}

	defer c.Close()

	dp := protodesc.ToDescriptorProto(proto.MessageV2((*msgs)[0]).ProtoReflect().Descriptor())
	fullTableName := fmt.Sprintf("projects/%s/datasets/%s/tables/%s",
		projectId,
		datasetId,
		tableName+tableNameSuffix())

	managedStream, err := c.NewManagedStream(ctx,
		managedwriter.WithDestinationTable(fullTableName),
		managedwriter.WithSchemaDescriptor(dp))
	if err != nil {
		fmt.Print("error managedStream\n")
		return err
	}

	encoded := make([][]byte, len(*msgs))

	for k, v := range *msgs {
		msg := proto.MessageV2(v)
		b, err := proto2.Marshal(msg)
		if err != nil {
			fmt.Printf("error proto marshal %v\n", err.Error())
		}
		encoded[k] = b
	}

	result, err := managedStream.AppendRows(ctx, encoded)
	if err != nil {
		return err
	}

	_, err = result.GetResult(ctx)
	if err != nil {
		return err
	}

	return nil
}
