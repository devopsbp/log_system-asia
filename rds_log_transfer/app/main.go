package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"rds_log_transfer/config"
	shard_db "rds_log_transfer/shard-db"
	"rds_log_transfer/transfer_worker"
	"runtime"
	"sync"
	"syscall"
	"time"
)

const loggerFormat = "20060102"

func main() {
	logger, loggerDate := getLogger()

	exitSigCh := make(chan os.Signal, 10)
	signal.Notify(exitSigCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	shardDbNum, err := config.GetShardDBNum()
	if err != nil {
		panic(err)
	}

	// イベント受信ch
	workerEvCh := make(transfer_worker.WorkerEvChan, 20)
	// シグナル等によるworker中断用のコンテクスト生成
	cancelCtx := context.Background()
	cancelCtx, cancelWorkers := context.WithCancel(cancelCtx)
	// workerのWaitGroup
	wg := sync.WaitGroup{}

	// create workers
	for i := 0; i < shardDbNum; i++ {
		shardDbOption, err := shard_db.GetShardDbOption(i)
		if err != nil {
			panic(err)
		}
		wg.Add(1)

		worker := transfer_worker.New(i, *shardDbOption, workerEvCh, cancelCtx, &wg)
		go worker.Run()
	}

	// workerの終了を検知するチャネル
	workerEndCh := make(chan struct{})
	go func() {
		wg.Wait()
		workerEndCh <- struct{}{}
	}()

	memStats := runtime.MemStats{}
	// listen worker
listenLoop:
	for {
		// rotate output path
		if time.Now().Format(loggerFormat) != loggerDate {
			logger, loggerDate = getLogger()
		}

		runtime.ReadMemStats(&memStats)
		fmt.Print("main: memory allocated: ", memStats.Alloc, "\n")

		// listen worker-ch
		select {
		case v := <-workerEvCh:
			switch v.Event {
			case transfer_worker.EventNotice:
				fmt.Printf("worker%d: %s\n", v.WorkerNo, v.Message)
			case transfer_worker.EventInfo:
				logger.Infof("worker%d: %s", v.WorkerNo, v.Message)
			case transfer_worker.EventError:
				logger.Errorf("worker%d: %s", v.WorkerNo, v.Message)
			}
		case <-exitSigCh: // 中断シグナル
			logger.Info("main: Exit signal received. cancel workers.")
			cancelWorkers()
		case <-workerEndCh: // Worker全終了
			break listenLoop
		}
	}

	logger.Info("main: exit rds_log_transfer.")
}

// logger取得
func getLogger() (*zap.SugaredLogger, string) {
	// logs directory
	if _, err := os.Stat("logs"); err != nil {
		os.Mkdir("logs", 0666)
	}

	env, _ := config.LoadBqBatchEnv("APP_ENV")
	var logConf zap.Config
	if env == "production" {
		logConf = zap.NewProductionConfig()
	} else {
		logConf = zap.NewDevelopmentConfig()
	}

	// output paths
	date := time.Now().Format(loggerFormat)
	logConf.OutputPaths = []string{"stdout", fmt.Sprintf("logs/info_%s.log", date)}
	logConf.ErrorOutputPaths = []string{"stderr", fmt.Sprintf("logs/error_%s.log", date)}

	logger, err := logConf.Build()
	if err != nil {
		panic(err)
	}

	sugaredLogger := logger.Sugar()
	return sugaredLogger, date
}
