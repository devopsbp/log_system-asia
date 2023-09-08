package transfer_worker

import (
	"database/sql"
	"fmt"
	"rds_log_transfer/bq"
	shard_db "rds_log_transfer/shard-db"
	"strings"
	"time"
)

// 通常Workerでエラーになったレコードをリトライする

type RetryWorker struct {
	mainWorker *Worker
	retryEndCh chan struct{}
}

func newRetryWorker(worker *Worker, retryEndCh chan struct{}) RetryWorker {
	return RetryWorker{
		mainWorker: worker,
		retryEndCh: retryEndCh,
	}
}

// リトライ間隔
const retryInterval = time.Hour * 2
const retrySelectLimit = 500

func (w *RetryWorker) run() {
	defer func() {
		// RetryWorkerが終了したことを通知
		w.retryEndCh <- struct{}{}
	}()

	// 別スレッドでWaitする
	retryExecCh := make(chan struct{})
	go waitRetry(retryExecCh)
	for {
		select {
		// リトライ実行
		case <-retryExecCh:
			w.mainWorker.pushMessage("retry exec")
			err := w.retry()
			if err != nil {
				w.mainWorker.pushError(err)
			}
		// 終了
		case <-w.mainWorker.cancelCtx.Done():
			return
		}
	}
}

// リトライ実行を通知する
func waitRetry(execCh chan struct{}) {
	for {
		execCh <- struct{}{}
		time.Sleep(retryInterval)
	}
}

func (w *RetryWorker) retry() error {
	// リトライはインターバルがあるので実行の度コネクションを生成・破棄
	writerDb, readerDb, err := w.mainWorker.createDbConnections()
	if err != nil {
		return err
	}
	defer func() {
		err := writerDb.Close()
		if err != nil {
			w.mainWorker.pushError(err)
		}
		err = readerDb.Close()
		if err != nil {
			w.mainWorker.pushError(err)
		}
	}()

	var offsetId uint64 = 0

	for {
		// エラーのレコード取得
		errorRecords, err := w.fetchErrorRows(readerDb, offsetId)
		if err != nil {
			return err
		}

		// エラーが無ければ終了
		if errorRecords == nil || len(*errorRecords) == 0 {
			w.mainWorker.pushMessage("no error records")
			return nil
		}

		// BQ投入
		successIds, err := w.sendToBq(errorRecords)
		if err != nil {
			return err
		}

		// エラーログテーブルから削除
		err = w.removeErrorRows(successIds, writerDb)
		if err != nil {
			return err
		}

		offsetId = (*errorRecords)[len(*errorRecords)-1].Id
	}
}

// sb_log_errorと結合しリトライ対象のレコードを取得する
func (w *RetryWorker) fetchErrorRows(readerDb *sql.DB, offsetId uint64) (*[]*shard_db.SbLogRecord, error) {
	// クエリ実行
	rows, err := readerDb.Query(fmt.Sprintf("SELECT l.`id`, `log_time`, `log_tag`, `user_id`, `character_id`,"+
		" `item_id`, `unique_id`, `log`, `session_id`, `map_id`, `pos_x`, `pos_y`, `pos_z` FROM %s l"+
		" INNER JOIN `sb_log_error` e ON l.`id` = e.`error_id` "+
		" WHERE l.`id` > %d AND e.`reason` = %d ORDER BY l.`id` ASC LIMIT %d;",
		shard_db.SbLogTableName, offsetId, shard_db.ErrorReasonBigQueryError, retrySelectLimit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]*shard_db.SbLogRecord, 0)

	// 結果取得
	for rows.Next() {
		record := shard_db.SbLogRecord{}
		err := rows.Scan(
			&record.Id, &record.LogTime, &record.LogTag,
			&record.UserId, &record.CharacterId, &record.ItemId,
			&record.UniqueId, &record.Log, &record.SessionId, &record.MapId,
			&record.PosX, &record.PosY, &record.PosZ,
		)
		if err != nil {
			w.mainWorker.pushError(err)
		}

		records = append(records, &record)
	}
	rows.Close()

	return &records, nil
}

func (w *RetryWorker) sendToBq(errorRecords *[]*shard_db.SbLogRecord) (*[]uint64, error) {
	clientWithContext, err := bq.GetClientWithContext()
	if err != nil {
		return nil, err
	}
	defer clientWithContext.Client.Close()

	// log_tagでグループ
	var groupedRecords = make(map[string][]*shard_db.SbLogRecord)
	for _, record := range *errorRecords {
		groupedRecords[record.LogTag] = append(groupedRecords[record.LogTag], record)
	}

	// 投入成功したID
	successIds := make([]uint64, 0)

	// 送信
	for tag, records := range groupedRecords {
		parseFailedRecords, err := sendBq(tag, records, w.mainWorker.no, clientWithContext)
		// エラーレコード記録
		// errorが返れば投入失敗 全レコード失敗と判断
		if err != nil {
			w.mainWorker.pushError(err)
		} else {
			// 失敗したIDを弾いて成功IDリストに追加
			var failedIds map[uint64]bool
			if parseFailedRecords != nil && len(parseFailedRecords) > 0 {
				for _, failedRecord := range parseFailedRecords {
					failedIds[failedRecord.Id] = true
				}
			}

			for _, record := range records {
				_, isFailed := failedIds[record.Id]
				if !isFailed {
					successIds = append(successIds, record.Id)
				}
			}
		}
	} // end for

	return &successIds, nil
}

func (w *RetryWorker) removeErrorRows(ids *[]uint64, writerDb *sql.DB) error {
	recordCount := len(*ids)
	if recordCount == 0 {
		return nil
	}

	values := make([]interface{}, 0, recordCount)
	sqlDelete := "DELETE FROM `sb_log_error` WHERE `error_id` IN ("

	for _, id := range *ids {
		sqlDelete = sqlDelete + "?,"
		values = append(values, id)
	}

	sqlDelete = strings.TrimSuffix(sqlDelete, ",") + ")"

	prepare, err := writerDb.Prepare(sqlDelete)
	if err != nil {
		return fmt.Errorf("retry error_log prepare failed. %v", err.Error())
	}

	result, err := prepare.Exec(values...)
	if err != nil {
		return fmt.Errorf("retry error_log delete failed. %v", err.Error())
	}
	deletedRows, err := result.RowsAffected()
	if err != nil || deletedRows == 0 {
		return fmt.Errorf("retry error_log delete failed. %v", err.Error())
	}

	w.mainWorker.pushInfoLog(fmt.Sprintf("deleted error records for %d rows", deletedRows))

	return nil
}
