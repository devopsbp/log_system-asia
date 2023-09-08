<?php
require_once('DatabaseUnix.php');

use CloudEvents\V1\CloudEventInterface;
use Google\Cloud\BigQuery\BigQueryClient;
use Google\Cloud\Core\Exception\GoogleException;
use Google\Cloud\Core\ExponentialBackoff;
use Google\Cloud\Samples\CloudSQL\Postgres\DatabaseUnix;
use Google\Cloud\Storage\Bucket;
use Google\Cloud\Storage\StorageClient;
use Google\CloudFunctions\FunctionsFramework;

const STATUS_NEW      = 'new';
const STATUS_WORKING  = 'working';
const STATUS_FINISHED = 'finished';
const STATUS_ERROR    = 'error';

const LOG_DIRECTORY_LOAD_STATUS_TABLE_NAME = 'sb_logs_cloud_to_bq_work_directory';
const LOG_FILE_LOAD_STATUS_TABLE_NAME      = 'sb_logs_cloud_to_bq_work_files';

const REGEX_LOG_FILE_NAME   = '/.*(SB_LOG_.*$)/';
const REGEX_TABLE_NAME      = '/(SB_LOG_.*_\d{3})(_.*\.csv\.gz)/';
const REGEX_RETRY_FILE_NAME = '/(SB_.*)(\.csv\.gz)/';

const MAX_SEND_UUID_COUNT = 9900;

// Register the function with Functions Framework.
// This enables omitting the `FUNCTIONS_SIGNATURE_TYPE=cloudevent` environment
// variable when deploying. The `FUNCTION_TARGET` environment variable should
// match the first parameter.
FunctionsFramework::cloudEvent('RecoverInterruptedLogs', 'RecoverInterruptedLogs');

function RecoverInterruptedLogs(CloudEventInterface $event): void
{
    $is_terminated = false;

    // -- 環境設定の読み込み ---
    $env = setEnv();

    // -- CloudSQL 向けに環境変数を設定 ---
    putenv("DB_USER={$env['SB_LOG_DB_USER']}");
    putenv("DB_PASS=".trim(getenv('SB_DB_PW')));
    putenv("DB_NAME={$env['SB_LOG_DB_NAME']}");
    putenv("INSTANCE_UNIX_SOCKET=/cloudsql/{$env['SB_LOG_DB_INSTANCE_UNIX_SOCKET']}");

    // -- 実行開始時間を保管 ---
    $execute_start_time = time();
    $execute_limit_time = $execute_start_time + $env['SB_LOG_LOAD_FUNCTION_EXECUTE_LIMIT_TIME'];

    // ---
    // Google Cloud Storage からログデータを取得
    $storage_client = new StorageClient(
        [
            'projectId' => $env['SB_LOG_CLOUD_STORAGE_PROJECT_ID'],
        ]
    );

    $storage_bucket = $storage_client->bucket($env['SB_LOG_CLOUD_STORAGE_BUCKET_NAME']);

    $bq_config = [
        'projectId' => $env['SB_LOG_BQ_GCP_PROJECT_ID'],
    ];

    $bq_client = new BigQueryClient($bq_config);
    // ---

    $pdo = DatabaseUnix::initUnixDatabaseConnection();

    // リトライ対象になる境界時刻
    // 最終更新時刻がこの時刻以前で処理が終了していないものはリトライ対象になる
    $retry_check_border_time = (new DateTime())->modify("-{$env['SB_LOG_RELOAD_INTERVAL_SECONDS']} second");

    /* ***
     * ファイル読み込み処理（LoadToBQ）でステータス更新が行われる際に CloudSQL のクエリエラーで停止し、
     * ステータス登録が行われない場合があるため、エラーステータスの処理は ファイル → ディレクトリ の順で処理する
     * */

    // まずファイルのエラーステータスから処理を行う
    $error_files = getRetryFiles($pdo, $retry_check_border_time->format('Y-m-d H:i:s'));

    $bq_table_cache = array();

    foreach ($error_files as $file_info) {
        if ($execute_limit_time < time()) {
            logWarning("Execution time exceeded stipulated limit.");
            $is_terminated = true;
            break;
        }

        $log_file = $storage_bucket->object($file_info['log_file']);

        if (!isset($log_file)) {
            // failed to get log file.
            continue;
        }

        $log_file_name = preg_replace(REGEX_LOG_FILE_NAME, '$1', $log_file->name());
        $table_name    = preg_replace(REGEX_TABLE_NAME, '$1', $log_file_name);
        logDefault(sprintf('[GCS] log file name: %s. table name: %s', $log_file_name, $table_name));

        $target_dataset_name = $file_info['target_dataset'];

        if (!isset($bq_table_cache[$table_name])) {
            // 投入先ログテーブルからスキーマを取得する
            $bq_dataset = $bq_client->dataset($target_dataset_name);
            $table = $bq_dataset->table($table_name);
            $table_fields = $table->info()['schema']['fields'];
            $table_schema = array_column($table_fields, 'type', 'name');

            $bq_table_cache[$table_name] = $table_schema;
        }
        $table_schema = $bq_table_cache[$table_name];
        logDefault(sprintf('[BQ] [%s] schema: {\'%s\'}', $table_name, implode("','", array_keys($table_schema))));

        //  CloudStorage からエラー中断したログデータをダウンロード
        $dl_file_path = "./" . $log_file_name;

        try {
            $log_file->downloadToFile($dl_file_path);

            // ログファイルと投入済みレコードから差分データを抽出する
            $log_data_from_csv  = makeDiffLogData($dl_file_path, $table_schema, $target_dataset_name, $table_name, $bq_client);
            $log_data_row_count = ((!is_null($log_data_from_csv)) ? count($log_data_from_csv) : 0);

            if ($log_data_row_count != 0) {
                // 差分データがある（１件以上のデータがある）なら、新しい投入用データを作成する
                logDefault("[BQ] new file row count: $log_data_row_count");

                $temp_log_file = preg_replace(REGEX_RETRY_FILE_NAME, "$1_retry$2", $log_file_name);

                // CSV からデータを作成し、 CloudStorage にアップロードする
                if (!is_null($temp_log_file) && makeRetryFile($temp_log_file, $log_data_from_csv)) {
                    $upload_directory = sprintf("%s/retry_%s", $env['SB_LOG_CLOUD_STORAGE_RETRY_FILE_PATH'], makeUUID());
                    $upload_target    = "$upload_directory/$log_file_name";

                    // ファイルアップロード
                    uploadRetryLogFile($storage_bucket, $upload_target, $upload_directory, $temp_log_file, $pdo, $target_dataset_name);
                } else {
                    logError("[$log_file_name] upload file is not exist.",);
                    continue;
                }
            } else {
                // 再投入レコードが無い場合はスキップ
                logDefault("[BQ] already loaded '$log_file_name' into '$target_dataset_name' all rows.");
            }

            // ダウンロードファイルの削除
            unlink($dl_file_path);
        } catch (\Google\Cloud\Core\Exception\NotFoundException $e) {
            // 取得するファイルがクラウド上になかった場合は警告メッセージを出してスキップ
            logError(sprintf("[%s] file is not found. [msg] %s", $log_file_name, $e->getMessage()));
        } catch (\Exception $e) {
            logError(sprintf("[%s] exception! [msg] %s", $log_file_name, $e->getMessage()));
            continue;
        }

        // ファイル読み込みステータスを「完了済み」に更新する
        updateFileStatusFinished($pdo, $file_info['log_file'], $target_dataset_name);
    }

    if (!$is_terminated) {
        // 次にディレクトリステータスの処理を行う
        $retry_directories = getRetryDataList($pdo, $retry_check_border_time->format('Y-m-d H:i:s'));
        foreach ($retry_directories as $dir_dataset_key => $retry_data_list) {
            if ($execute_limit_time < time()) {
                logWarning("Execution time exceeded stipulated limit.");
                $is_terminated = true;
                break;
            }

            logDefault("[$dir_dataset_key] start check status.");

            /*
             * 1. ファイルステータスと照らし合わせ、全て「処理完了」となっているか
             *
             *  ここまででエラーステータスのファイルを処理しているので、
             * ディレクトリ，データセットに対応するファイルが全て処理されていれば「完了」とする
             * */
            if (checkAllRetried($pdo, $dir_dataset_key, $retry_data_list)) {
                // 全て「完了」になったので終了
                continue;
            }

            /*
             * 2. 処理ディレクトリ，データセットに対してファイルの登録があるか
             *
             *  ファイルの登録がない場合はＢＱ投入（ LoadToBQ ）がまだ行われていないため、
             * ステータスをリセットして新規処理対象とする
             * */
            if (checkRestartable($pdo, $dir_dataset_key, $retry_data_list)) {
                // リトライすることとなったら終了
                continue;
            }
        }
    }

    // 明示的にコネクションを破棄
    $pdo = null;

    if ($is_terminated) {
        logWarning("function terminated.");
    } else {
        logDefault("function complete.");
    }
}

/**
 * 取得した CSV ファイルと BigQuery の投入済みレコードを比較し
 * 「 BigQuery に未投入のレコード」データを抽出する
 *
 * @param string $dl_file_path
 * @param array $table_schema
 * @param string $target_dataset_name
 * @param string $table_name
 * @param BigQueryClient $bq_client
 * @return array|null
 * @throws GoogleException
 */
function makeDiffLogData(
    string         $dl_file_path,
    array          $table_schema,
    string         $target_dataset_name,
    string         $table_name,
    BigQueryClient $bq_client
): array|null
{
    $log_data_from_csv = array();
    try {
        if (!$handler = gzopen($dl_file_path, 'r')) {
            throw new Exception("file open error.");
        }

        // CSV データを行ごとに配列データとして読み取る
        $row_tmp = '';
        while ($line = gzgets($handler, (4*1024))) {
            $row_tmp .= $line;
            if(!preg_match('/\\n$/', $row_tmp) ){
                // 改行まで取得する
                continue;
            }

            $line = $row_tmp;
            $columns = str_getcsv(rtrim($row_tmp));
            $row_tmp = '';

            if (count($table_schema) != count($columns)) {
                // カラム数が不一致のためデータフォーマットが不正
                // → スキップ
                logWarning(sprintf(
                    '[%s] does not match schema format. [schema] \'%s\' [row] \'%s\'',
                    $table_name, implode("','", array_keys($table_schema)), implode("','", $columns)
                ));

                continue;
            }

            // スキーマと照らし合わせる
            $row_data = array_combine(array_keys($table_schema), $columns);
            $is_correct_format = false;

            foreach ($row_data as $column_name => $column_value) {
                $column_type = $table_schema[$column_name];

                if (!empty($column_value)) {
                    switch ($column_type) {
                        case 'TIMESTAMP':
                            $is_correct_format = date_parse_from_format($column_value, '%Y-%m-%d %H:%M:%S');
                            break;
                        case 'BOOLEAN':
                            $is_correct_format = (($column_value == 'true') || ($column_value == 'false'));
                            break;
                        case 'FLOAT':
                            $is_correct_format = (
                                is_numeric($column_value)
                                && ((float)$column_value || ($column_value == '0.0') || ($column_value == '0'))
                            );
                            break;
                        case 'INTEGER':
                            $is_correct_format = (
                                is_numeric($column_value)
                                && ((int)$column_value || ($column_value == '0'))
                            );
                            break;
                        case 'STRING':
                            // 『"』が正しく閉じられているか確認（出現数が２で割り切れるなら「閉じられている」と判断）
                            $double_quote_count_mod = (substr_count($column_value, '"') % 2);
                            if ($double_quote_count_mod == 0) {
                                $is_correct_format = true;
                            } else {
                                logWarning(sprintf('[%s] data is not closed by quote character. \'%s\'', $table_name, $column_value));
                            }
                            break;
                        default:
                            break;
                    }
                } else {
                    // 値が空である
                    switch ($column_type) {
                        case 'STRING':
                            // STRING 型の場合は空でも許容
                            $is_correct_format = true;
                            break;
                        default:
                            break;
                    }
                }

                if (!$is_correct_format) {
                    // カラムのデータが型に合わない
                    // → スキップ
                    logWarning(sprintf(
                        '[%s] does not match schema type. column \'%s\' (value: \'%s\') is not \'%s\' type.',
                        $table_name, $column_name, $column_value, $column_type
                    ));

                    break;
                }
            }

            if ($is_correct_format) {
                // "uuid" をキーとして行データを格納する
                $log_data_from_csv[$row_data['uuid']] = $line;
            }
        }
    } catch (Exception $e) {
        logError($e->getMessage());
    } finally {
        if ($handler) {
            gzclose($handler);
        }
    }

    logDefault(sprintf("log data row count: %s", count($log_data_from_csv)));

    if (count($log_data_from_csv) > 0) {
        try {
            // 条件に指定する uuid のパラメータ数が多くなりクエリ上限にかかる可能性があるので
            // 一定数ごとにクエリを分割する
            // https://cloud.google.com/bigquery/quotas?hl=ja#query_jobs
            foreach (array_chunk(array_keys($log_data_from_csv), MAX_SEND_UUID_COUNT) as $chunk_records) {
                // BigQuery のテーブルから投入済みレコードを取得する
                $query = sprintf(
                    "SELECT `uuid` FROM %s.%s WHERE `uuid` IN ( '%s' )",
                    $target_dataset_name,
                    $table_name,
                    implode("', '", $chunk_records)
                );

                logDefault("[BQ] SEND_UUID_COUNT= " . count($chunk_records));
                logDefault("[BQ] SQL= $query");

                $bq_query = $bq_client->query($query);
                $res = $bq_client->runQuery($bq_query);

                // ※ 時間はかからない想定なので再試行は３回まで
                $backoff = new ExponentialBackoff(3);
                $backoff->execute(
                    function () use ($res) {
                        $res->reload();
                        if (!$res->isComplete()) {
                            throw new Exception('Job has not yet completed', 500);
                        }
                    }
                );

                // CSV から読み取ったデータから、 BigQuery に投入済みデータを除外する
                foreach ($res->rows() as $row) {
                    if (isset($log_data_from_csv[$row['uuid']])) {
                        unset($log_data_from_csv[$row['uuid']]);
                    }
                }
            }
        } catch (Exception $e) {
            // エラーが発生した場合はスキップする
            logError("[BQ] failed to get exist logs. [ErrorMessage] {$e->getMessage()}");

            // 中途半端な差分が発生してしまうとレコード重複の原因となるので
            // null を返して終了する
            return null;
        }
    }

    return $log_data_from_csv;
}

/**
 * ログデータから再投入ファイルを作成する
 *
 * @param string $upload_file_path_
 * @param array $log_data_
 * @return bool
 */
function makeRetryFile(string $upload_file_path_, array $log_data_): bool
{
    $has_error = false;

    try {
        // 再投入ファイルを作成
        if (!$handler = gzopen("./$upload_file_path_", 'wb')) {
            throw new Exception("file open error.");
        }

        foreach ($log_data_ as $row) {
            gzputs($handler, $row);
        }
    } catch (Exception $e) {
        logError($e->getMessage());
        $has_error = true;
    } finally {
        if ($handler) {
            gzclose($handler);
        }
    }

    return (!$has_error && file_exists("./$upload_file_path_"));
}

/**
 *
 *
 * @param Bucket $storage_bucket_
 * @param string $upload_target_
 * @param string $upload_directory_
 * @param string $upload_file_path_
 * @param PDO $pdo_
 * @param mixed $target_dataset_
 */
function uploadRetryLogFile(Bucket $storage_bucket_, string $upload_target_, string $upload_directory_, string $upload_file_path_, PDO $pdo_, mixed $target_dataset_): void
{
    logDefault("[CloudStorage] new log file upload to '$upload_target_'");

    // CloudStorage にアップロード
    $storage_bucket_->upload(
        fopen("./$upload_file_path_", 'r'),
        ['name' => $upload_target_]
    );

    unlink("./$upload_file_path_");

    // 改めて作成し直した CSV ログデータをアップロードしたディレクトリを
    // 'new' ステータスで DB に登録する
    $now_ = now();
    $query = sprintf(
        "INSERT INTO %s ( uniq_key, log_directory, target_dataset, status, log_time, created_at ) VALUES ( '%s:%s', '%s', '%s', '%s', '%s', '%s' )",
        LOG_DIRECTORY_LOAD_STATUS_TABLE_NAME,
        $target_dataset_, $upload_directory_,
        $upload_directory_,
        $target_dataset_,
        STATUS_NEW,
        $now_,
        $now_
    );

    $errMessage = sprintf("[%s] new record insert error.", LOG_DIRECTORY_LOAD_STATUS_TABLE_NAME);
    if (!tryRunQuery($pdo_, $query, $errMessage)) {
        logError($errMessage);
    }
}

/**
 * ファイルの読み込みステータスを「完了」に更新する
 *
 * @param PDO $pdo_
 * @param string $log_file_
 * @param string $target_dataset_
 * @return void
 */
function updateFileStatusFinished(PDO $pdo_, string $log_file_, string $target_dataset_): void
{
    $query = sprintf(
        "UPDATE %s SET status = '%s', log_time = '%s' WHERE ( log_file = '%s' AND target_dataset = '%s' )",
        LOG_FILE_LOAD_STATUS_TABLE_NAME,
        STATUS_FINISHED,
        now(),
        $log_file_,
        $target_dataset_
    );

    $errMessage = sprintf("[%s] status update failed.", LOG_FILE_LOAD_STATUS_TABLE_NAME);
    if (!tryRunQuery($pdo_, $query, $errMessage)) {
        logError($errMessage);
    }
}

/**
 * ファイルステータスと照らし合わせ全て「処理完了」となっているかを確認し
 * 全て「完了」となっていた場合ディレクトリステータスを「完了」に更新する
 *
 * @param PDO $pdo_
 * @param string $dir_dataset_key_
 * @param array $retry_data_list_
 * @return bool
 *      全て「完了」となっているなら true を返す
 * @throws Exception
 */
function checkAllRetried(PDO $pdo_, string $dir_dataset_key_, array $retry_data_list_): bool
{
    $keys = str_getcsv($dir_dataset_key_);
    $dir_name     = $keys[0];
    $dataset_name = $keys[1];

    $is_all_retried = (count($retry_data_list_) > 0);
    foreach ($retry_data_list_ as $data_info) {
        $is_all_retried &= ($data_info['file_status'] == STATUS_FINISHED);

        if (!$is_all_retried) {
            break;
        }
    }

    logDefault("[$dir_dataset_key_] \$is_all_retried: $is_all_retried");

    if (!$is_all_retried) {
        // 全て「完了」でないなら false を返す
        return false;
    }

    // ファイルステータスがすべて「完了」なので、ディレクトリステータスを「完了」に更新
    $query = sprintf(
        "UPDATE %s SET status = '%s', log_time = '%s' WHERE ( log_directory = '%s' AND target_dataset = '%s' )",
        LOG_DIRECTORY_LOAD_STATUS_TABLE_NAME,
        STATUS_FINISHED,
        now(),
        $dir_name,
        $dataset_name
    );

    $errMessage = sprintf("[%s] [%s] status update failed.", LOG_FILE_LOAD_STATUS_TABLE_NAME, $dir_dataset_key_);
    if (!tryRunQuery($pdo_, $query, $errMessage)) {
        logError($errMessage);
        throw new Exception(message: $errMessage);
    }

    return true;
}

/**
 * ディレクトリに紐づいたファイルの登録があるかを確認する
 * ※ ファイルの登録がなければ、投入処理をやり直すだけで読み込みが再開，続行される
 *
 * @param PDO $pdo_
 * @param string $dir_dataset_key_
 * @param array $retry_data_list_
 * @return bool
 * @throws Exception
 */
function checkRestartable(PDO $pdo_, string $dir_dataset_key_, array $retry_data_list_): bool
{
    $keys = str_getcsv($dir_dataset_key_);
    $dir_name     = $keys[0];
    $dataset_name = $keys[1];

    // ファイルの登録件数をカウント
    $registered_file_count = count(array_filter(array_column($retry_data_list_, 'file')));

    logDefault("[$dir_dataset_key_] \$registered_file_count: $registered_file_count");

    if ($registered_file_count != 0) {
        // ファイルの登録があった場合、ファイルのエラーを解決する必要があるので
        // ここでは何もできない
        return false;
    }

    $query = sprintf(
        "UPDATE %s SET status = '%s', log_time = '%s' WHERE log_directory = '%s' AND target_dataset = '%s'",
        LOG_DIRECTORY_LOAD_STATUS_TABLE_NAME,
        STATUS_NEW,
        now(),
        $dir_name,
        $dataset_name
    );

    $errMessage = sprintf("[%s] [%s] status reset failed.", LOG_FILE_LOAD_STATUS_TABLE_NAME, $dir_dataset_key_);
    if (!tryRunQuery($pdo_, $query, $errMessage)) {
        logError($errMessage);
        throw new Exception(message: $errMessage);
    }

    return true;
}

/**
 * リトライ対象のファイルを取得する
 *
 * ※ リトライ対象：
 *  「ステータスが 'error' 」もしくは「ステータスが 'working' 且つ処理時間から規定秒数経過している」
 *
 * @param PDO $pdo_
 * @param string $border_time_str
 * @return array
 * @throws Exception
 */
function getRetryFiles(PDO $pdo_, string $border_time_str): array
{
    $query = sprintf(
        "SELECT log_file, target_dataset, status FROM %s WHERE ( status = '%s' OR ( status = '%s' AND log_time < '%s' ) )",
        LOG_FILE_LOAD_STATUS_TABLE_NAME,
        STATUS_ERROR,
        STATUS_WORKING,
        $border_time_str
    );

    $errMessage = "failed to get error directories.";
    if (!tryRunQuery(pdo_: $pdo_, query_: $query, error_message_: $errMessage, result_: $res)) {
        throw new Exception($errMessage);
    }

    return $res->fetchAll(PDO::FETCH_ASSOC);
}

/**
 * ディレクトリをベースとしたリトライ対象のリストを取得
 *
 * @param PDO $pdo_
 * @param string $border_time_str
 * @return array
 * @throws Exception
 */
function getRetryDataList(PDO $pdo_, string $border_time_str): array
{
    // ディレクトリの処理ステータスが「エラー」もしくは「処理中 且つ 規定時間を経過している」データを取得する
    $query = preg_replace('/[ ' . PHP_EOL . ']+/', ' ',
        sprintf(
            <<<EOC
            SELECT
                dir_tbl.log_directory  AS directory,
                dir_tbl.target_dataset AS dataset,
                dir_tbl.status         AS dir_status,
                dir_tbl.log_time       AS dir_log_time,
                file_tbl.log_file      AS file,
                file_tbl.status        AS file_status,
                file_tbl.log_time      AS file_log_time
            FROM
                %s AS dir_tbl LEFT OUTER JOIN %s AS file_tbl 
            ON
                file_tbl.target_dataset = dir_tbl.target_dataset AND file_tbl.log_directory = dir_tbl.log_directory
            WHERE
                dir_tbl.status = '%s' 
                OR ( dir_tbl.status = '%s' AND dir_tbl.log_time < '%s' )
            EOC,
            LOG_DIRECTORY_LOAD_STATUS_TABLE_NAME,
            LOG_FILE_LOAD_STATUS_TABLE_NAME,
            STATUS_ERROR,
            STATUS_WORKING,
            $border_time_str
        )
    );

    $errMessage = "failed to get error directories.";
    if (!tryRunQuery(pdo_: $pdo_, query_: $query, error_message_: $errMessage, result_: $res)) {
        throw new Exception($errMessage);
    }

    // 「'directory','dataset'」の形をキーとした配列を生成，返却する
    // → ディレクトリ，対象データセット単位で纏める
    $res_array = array();
    foreach ($res->fetchAll(PDO::FETCH_ASSOC) as $row) {
        $arr_key = "{$row['directory']},{$row['dataset']}";

        if (!isset($res_array[$arr_key])) {
            $res_array[$arr_key] = array();
        }

        $res_array[$arr_key][] = $row;
    }

    return $res_array;
}

/**
 * CloudSQL へのクエリの実行，エラーチェックを行う
 *
 * @param PDO $pdo_
 * @param string $query_
 * @param string $error_message_
 *  エラーが発生した際のメッセージ
 * @param false|PDOStatement|null $result_
 * @return bool
 *  実行に成功した場合 "true" を返す
 */
function tryRunQuery(PDO $pdo_, string $query_, string $error_message_, false|PDOStatement &$result_ = null): bool
{
    try {
        logDefault("[CloudSQL] SQL= $query_");
        $result_ = $pdo_->query($query_);

        $err_code = $result_->errorCode();
        if (!is_null($err_code) && $err_code != 0) {
            throw new Exception(message: "{$result_->errorInfo()}", code: $err_code);
        }
    } catch (Exception $e) {
        $msg = str_replace(PHP_EOL, '', var_export($e->getMessage(), true));
        $msg = preg_replace('/ +/', ' ', $msg);
        logError("[code: {$e->getCode()}] [message] $error_message_ [SQL] $query_" . PHP_EOL . $msg);

        return false;
    }

    return true;
}

function setEnv(): array
{
    // -- 環境設定の読み込み ---
    $data = file(getenv('SB_SECRET_ID'));

    $env = array();
    foreach ($data as $row) {
        $split = explode('=', $row);

        $key = $split[0];
        $val = rtrim($split[1]);

        if (str_starts_with($key, 'ARR_')) {
            $val = explode(',', $val);
        }

        $env[$split[0]] = $val;
    }

    return $env;
}

/**
 * UUID を生成する
 *
 * @return string
 * @throws Exception
 * @link https://qiita.com/gongon143/items/8916ae5c3b394c5a7ac5
 */
function makeUUID(): string
{
    return preg_replace_callback(
        '/x|y/',
        function ($m) {
            return dechex($m[0] === 'x' ? random_int(0, 15) : random_int(8, 11));
        },
        'xxxxxxxx_xxxx_4xxx_yxxx_xxxxxxxxxxxx'
    );
}

function logDefault($log): void
{
    $msg = '[LOG] ' . $log . PHP_EOL;
    fwrite(fopen('php://stderr', 'wb'), $msg);
}

function logWarning($log): void
{
    $msg = '[WARN] ' . $log . PHP_EOL;
    fwrite(fopen('php://stderr', 'wb'), $msg);
}

function logError($log): void
{
    $msg = '[ERR] ' . $log . PHP_EOL;
    fwrite(fopen('php://stderr', 'wb'), $msg);
}

function now(): string
{
    return (new DateTime())->format('Y-m-d H:i:s');
}
