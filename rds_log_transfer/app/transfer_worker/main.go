package transfer_worker

import (
	"context"
	"database/sql"
	"fmt"
	"rds_log_transfer/bq"
	shard_db "rds_log_transfer/shard-db"
	"strings"
	"sync"
	"time"
)

type EvChanValues struct {
	WorkerNo int
	Message  string
	Event    int
	Err      error
}

type WorkerEvChan chan EvChanValues

// sb_logテーブルから一度に取得する件数
const selectLimit = 2000

// channel events
const (
	EventNotice int = iota
	EventInfo
	EventError
)

// Worker object struct
type Worker struct {
	no        int
	dbOption  shard_db.ShardDbOption
	evCh      WorkerEvChan
	cancelCtx context.Context
	wg        *sync.WaitGroup
	readerDb  *sql.DB
	writerDb  *sql.DB
}

// New worker生成
func New(no int, dbOption shard_db.ShardDbOption, evCh WorkerEvChan, cancelCtx context.Context, wg *sync.WaitGroup) *Worker {
	return &Worker{
		no:        no,
		dbOption:  dbOption,
		evCh:      evCh,
		cancelCtx: cancelCtx,
		wg:        wg,
	}
}

// Run Worker実行
func (w *Worker) Run() {
	// Retry Workerを同時に実行
	retryEndCh := make(chan struct{})
	retryWorker := newRetryWorker(w, retryEndCh)
	go retryWorker.run()

	// 後処理
	defer func() {
		recovered := recover()
		if recovered != nil {
			w.pushError(fmt.Errorf("an error occurred. error:%v", recovered))
		}
		w.onExit(&retryWorker)
	}()

	err := w.setupDb()
	if err != nil {
		w.pushError(err)
		return
	}

	// worker loop
	for {
		// オフセットID取得
		offsetId, err := w.getOffsetId()
		if err != nil {
			w.pushError(err)
		}

		// RDSログ取得
		sbLogRecords, err := w.fetchSbLogRows(offsetId)
		if err != nil {
			w.pushError(err)
		}

		w.pushMessage(fmt.Sprintf("fetched %d rows", len(*sbLogRecords)))

		if len(*sbLogRecords) > 0 {
			// BQに投入
			err = w.sendToBq(sbLogRecords)
			if err != nil {
				w.pushError(err)
			}

			// IDオフセットを更新
			err = w.updateOffsetId(sbLogRecords)
			if err != nil {
				w.pushError(err)
			}

			// mainに通知
			w.pushInfoLog(fmt.Sprintf("processed %d records", len(*sbLogRecords)))
		}

		// check onExit signal
		select {
		case <-w.cancelCtx.Done():
			w.pushInfoLog("signal received. exiting.")
			return
		default:
		}

		if !(len(*sbLogRecords) > 0) {
			time.Sleep(time.Second * 1)
		}
	}
}

// create db connection
func (w *Worker) createDbConnections() (*sql.DB, *sql.DB, error) {
	readerDb, err := sql.Open("mysql", w.dbOption.GetDSN(shard_db.DsnOption{IsWriter: false}))
	if err != nil {
		return nil, nil, err
	}
	writerDb, err := sql.Open("mysql", w.dbOption.GetDSN(shard_db.DsnOption{IsWriter: true}))
	if err != nil {
		readerDb.Close()
		return nil, nil, err
	}

	return writerDb, readerDb, nil
}

// set db connection
func (w *Worker) setupDb() error {
	if w.readerDb != nil {
		return nil
	}

	writerDb, readerDb, err := w.createDbConnections()
	if err != nil {
		return err
	}

	w.writerDb = writerDb
	w.readerDb = readerDb

	return nil
}

func (w *Worker) onExit(retryWorker *RetryWorker) {
	if w.readerDb != nil {
		err := w.readerDb.Close()
		if err != nil {
			w.pushError(err)
		}
	}
	if w.writerDb != nil {
		err := w.writerDb.Close()
		if err != nil {
			w.pushError(err)
		}
	}

	// retry workerの終了を待ってからWaitGroupに通知
	w.pushMessage("wait to end for retry worker")
	select {
	case <-retryWorker.retryEndCh:
		w.wg.Done()
		return
	}
}

// main にメッセージ送信
func (w *Worker) pushMessage(message string) {
	w.evCh <- EvChanValues{
		WorkerNo: w.no,
		Message:  message,
		Event:    EventNotice,
	}
}

// main にログ送信
func (w *Worker) pushInfoLog(message string) {
	w.evCh <- EvChanValues{
		WorkerNo: w.no,
		Message:  message,
		Event:    EventInfo,
	}
}

// main にエラー通知
func (w *Worker) pushError(err error) {
	w.evCh <- EvChanValues{
		WorkerNo: w.no,
		Message:  err.Error(),
		Event:    EventError,
		Err:      err,
	}
}

func (w *Worker) getOffsetId() (uint64, error) {
	// writer のコネクション
	db, err := sql.Open("mysql", w.dbOption.GetDSN(shard_db.DsnOption{IsWriter: true}))
	if err != nil {
		w.pushError(err)
		return 0, err
	}

	// 後処理
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			w.pushError(err)
		}
	}(db)

	var offsetId uint64
	err = db.QueryRow("SELECT offset_id FROM sb_log_offset").Scan(&offsetId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return offsetId, nil
}

// DBからSBログ取得
func (w *Worker) fetchSbLogRows(offsetId uint64) (*[]*shard_db.SbLogRecord, error) {
	// クエリ実行
	rows, err := w.readerDb.Query(fmt.Sprintf("SELECT `id`, `log_time`, `log_tag`, `user_id`, `character_id`,"+
		" `item_id`, `unique_id`, `log`, `session_id`, `map_id`, `pos_x`, `pos_y`, `pos_z` FROM %s"+
		" WHERE `id` > %d ORDER BY `id` ASC LIMIT %d;",
		shard_db.SbLogTableName, offsetId, selectLimit))
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
			w.pushError(err)
		}

		records = append(records, &record)
	}
	rows.Close()

	return &records, nil
}

// 処理済みのIDをオフセットテーブルに記録
func (w *Worker) updateOffsetId(sbLogRecords *[]*shard_db.SbLogRecord) error {
	lastId := (*sbLogRecords)[len(*sbLogRecords)-1].Id
	result, err := w.writerDb.Exec("UPDATE sb_log_offset SET `offset_id` = ?", lastId)
	if err != nil {
		return fmt.Errorf("updateOffsetId failed. %v", err.Error())
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("updateOffsetId failed. %v", err.Error())
	}

	// レコード無し時Insert
	if affectedRows == 0 {
		result, err := w.writerDb.Exec("INSERT INTO sb_log_offset (`offset_id`) VALUES (?)", lastId)
		if err != nil {
			return fmt.Errorf("sendBq offset_id failed. %v", err.Error())
		}
		insertedRows, err := result.RowsAffected()
		if err != nil || insertedRows == 0 {
			return fmt.Errorf("sendBq offset_id failed. %v", err.Error())
		}
	}

	return nil
}

func (w *Worker) sendToBq(sbLogRecords *[]*shard_db.SbLogRecord) error {
	clientWithContext, err := bq.GetClientWithContext()
	if err != nil {
		return err
	}
	defer clientWithContext.Client.Close()

	// log_tagでグループ
	var groupedRecords = make(map[string][]*shard_db.SbLogRecord)
	for _, record := range *sbLogRecords {
		groupedRecords[record.LogTag] = append(groupedRecords[record.LogTag], record)
	}

	// 送信
	for tag, records := range groupedRecords {
		parseFailedRecords, err := sendBq(tag, records, w.no, clientWithContext)
		// エラーレコード記録
		// errorが返れば投入失敗 全レコード失敗と判断
		if err != nil {
			w.pushError(err)
			err = w.insertErrorLogIds(&records, shard_db.ErrorReasonBigQueryError)
		} else if len(parseFailedRecords) > 0 {
			// JSONのパースに失敗したレコード
			err = w.insertErrorLogIds(&parseFailedRecords, shard_db.ErrorReasonParseFailed)
		}
		if err != nil { // これにも失敗した場合は止める
			panic(err)
		}
	} // end for

	return nil
}

// BQ投入に失敗したレコードのIDを記録する
func (w *Worker) insertErrorLogIds(records *[]*shard_db.SbLogRecord, reason shard_db.ErrorReason) error {
	recordCount := len(*records)
	if recordCount == 0 {
		return nil
	}

	values := make([]interface{}, 0, recordCount*2)
	sqlInsert := "INSERT INTO sb_log_error (`error_id`, `reason`) VALUES "

	for _, record := range *records {
		sqlInsert = sqlInsert + "(?, ?),"
		values = append(values, record.Id)
		values = append(values, reason)
	}

	prepare, err := w.writerDb.Prepare(strings.TrimSuffix(sqlInsert, ","))
	if err != nil {
		return fmt.Errorf("sendBq error prepare failed. %v", err.Error())
	}

	result, err := prepare.Exec(values...)
	if err != nil {
		return fmt.Errorf("sendBq error id failed. %v", err.Error())
	}
	insertedRows, err := result.RowsAffected()
	if err != nil || insertedRows == 0 {
		return fmt.Errorf("sendBq error id failed. %v", err.Error())
	}

	w.pushInfoLog(fmt.Sprintf("inserted error records for %d rows", insertedRows))

	return nil
}
