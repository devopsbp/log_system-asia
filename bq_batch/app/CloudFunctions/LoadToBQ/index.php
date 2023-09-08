<?php
require_once('DatabaseUnix.php');

use CloudEvents\V1\CloudEventInterface;
use Google\Cloud\BigQuery\BigQueryClient;
use Google\Cloud\BigQuery\Dataset;
use Google\Cloud\Core\ExponentialBackoff;
use Google\Cloud\Samples\CloudSQL\Postgres\DatabaseUnix;
use Google\Cloud\Storage\StorageClient;
use Google\CloudFunctions\FunctionsFramework;

const STATUS_NEW      = 'new';
const STATUS_WORKING  = 'working';
const STATUS_FINISHED = 'finished';
const STATUS_ERROR    = 'error';

const LOG_DIRECTORY_LOAD_STATUS_TABLE_NAME = 'sb_logs_cloud_to_bq_work_directory';
const LOG_FILE_LOAD_STATUS_TABLE_NAME      = 'sb_logs_cloud_to_bq_work_files';

const MAX_RETRY_COUNT = 10;

// Register the function with Functions Framework.
// This enables omitting the `FUNCTIONS_SIGNATURE_TYPE=cloudevent` environment
// variable when deploying. The `FUNCTION_TARGET` environment variable should
// match the first parameter.
FunctionsFramework::cloudEvent('LoadToBQ', 'LoadToBQ');

function LoadToBQ(CloudEventInterface $event): void
{
    // -- 環境設定の読み込み ---
    $env = setEnv();

    // -- 実行開始時間を保管 ---
    $execute_start_time = time();
    $execute_limit_time = $execute_start_time + $env['SB_LOG_LOAD_FUNCTION_EXECUTE_LIMIT_TIME'];

    // -- CloudSQL 向けに環境変数を設定 ---
    putenv("DB_USER={$env['SB_LOG_DB_USER']}");
    putenv("DB_PASS=".trim(getenv('SB_DB_PW')));
    putenv("DB_NAME={$env['SB_LOG_DB_NAME']}");
    putenv("INSTANCE_UNIX_SOCKET=/cloudsql/{$env['SB_LOG_DB_INSTANCE_UNIX_SOCKET']}");

    // -- 受け取ったメッセージから読み込み情報を解析 -

    $cloudEventData = $event->getData();
    $pubSubData = base64_decode($cloudEventData['message']['data']);
    $pubSubData = json_decode($pubSubData, true);

    if (!isset($pubSubData) || !isset($pubSubData['data'])) {
        logError("event data is empty");
        return;
    }

    logDefault(preg_replace('/['.PHP_EOL.' ]+/', ' ', var_export($pubSubData, true)));

    $target_dataset_list = $pubSubData['data']['target_dataset'];
    foreach ($target_dataset_list as $dataset) {
        logDefault("load target dataset = $dataset");
    }

    logDefault("[GCS] bucket: {$env['SB_LOG_CLOUD_STORAGE_BUCKET_NAME']}");

    // ---
    // Google Cloud Storage からログデータを取得
    $storage = new StorageClient(
        [
            'projectId' => $env['SB_LOG_CLOUD_STORAGE_PROJECT_ID'],
        ]
    );
    $bucket = $storage->bucket($env['SB_LOG_CLOUD_STORAGE_BUCKET_NAME']);

    $arg_dir_list = $pubSubData['data']['work_dir'];
    $log_files = array();

    logDefault(sprintf("[GCS] get objects from storage prefix: '%s'", implode("','", $arg_dir_list)));
    foreach ($arg_dir_list as $work_dir) {
        $saved_log_data = $bucket->objects(
            [
                'prefix' => $work_dir
            ]
        );

        $log_files = createLogFilesMap(saved_log_data_: $saved_log_data, work_dir_: $work_dir, append_log_files_base_: $log_files);
    }

    $pdo = DatabaseUnix::initUnixDatabaseConnection();

    foreach ($target_dataset_list as $dataset_id) {
        $work_dir_list = $arg_dir_list;

        logDefault("[BQ] dataset_id: $dataset_id");

        $bq_config = [
            'projectId' => $env['SB_LOG_BQ_GCP_PROJECT_ID'],
            'datasetId' => $dataset_id,
        ];

        if (!isset($bq_config['projectId']) || !isset($bq_config['datasetId'])) {
            logError('invalid gcp settings.');
            throw new Exception('invalid gcp settings.');
        }

        $bq_client = new BigQueryClient($bq_config);
        $bq_dataset = $bq_client->dataset(id: $bq_config['datasetId']);

        $work_time = now();

        // 既に読み込み開始済みのリストを取得し、今回作業するリストから除外する
        $skip_dirs = getAlreadyStartedList(pdo_: $pdo, work_dir_list_: $work_dir_list, dataset_id_: $dataset_id);
        $work_dir_list = array_diff($work_dir_list, $skip_dirs);

        array_walk(
            $log_files,
            function (&$files_per_directories, $log_table) use ($skip_dirs) {
                foreach ($skip_dirs as $skip_dir) {
                    if (isset($files_per_directories[$log_table][$skip_dir])) {
                        unset($files_per_directories[$log_table][$skip_dir]);
                    }
                }
            }
        );

        foreach ($log_files as $table => $directories) {
            foreach ($directories as $directory => $files) {
                if (!tryAddLoadStartStatus(pdo_: $pdo, dataset_id_: $dataset_id,
                    work_dir_: $directory, log_files_: $files, work_time_: $work_time)) {
                    // 読み込み開始ステータスの登録に失敗したらスキップしておく
                    $skip_dirs[] = $directory;
                }
            }
        }

        foreach ($skip_dirs as $skip_dir) {
            logWarning("[$dataset_id] work skipped '$skip_dir'");
        }
        unset($skip_dirs);

        $load_jobs = array();
        foreach ($log_files as $table => $files_per_directories) {
            $load_jobs[$table] = new Fiber('LoadToBigQueryFromStorage');
            $load_jobs[$table]->start(
                $bq_dataset, $table, $env['SB_LOG_CLOUD_STORAGE_BUCKET_NAME'],
                $pdo, $dataset_id, $files_per_directories, $work_time, $execute_limit_time
            );
        }

        $is_all_done = false;
        while (!$is_all_done) {
            $is_all_done = true;
            foreach ($load_jobs as $job){
                if($job->isSuspended()) {
                    $job->resume();
                }

                $is_all_done &= $job->isTerminated();
            }
        }

        $has_error = false;
        foreach ($load_jobs as $job) {
            $has_error |= ($job->getReturn() == false);

            if($has_error){
                break;
            }
        }

        $worked_dir_list = array();
        foreach ($log_files as $table => $dirs) {
            foreach ($dirs as $dir => $file) {
                $worked_dir_list[] = $dir;
            }
        }
        $worked_dir_list = array_unique($worked_dir_list);

        foreach ($worked_dir_list as $worked_dir) {
            setDirectoryLoadStatus(pdo_: $pdo, work_dir_: $worked_dir, dataset_id_: $dataset_id, has_error_: $has_error);
        }

        // 処理経過時間の確認
        $elapsed_time = (time() - $execute_start_time);
        logDefault("elapsed time: $elapsed_time(s) / limit time: {$env['SB_LOG_LOAD_FUNCTION_EXECUTE_LIMIT_TIME']}(s)");

        // Functions のタイムアウトとは別で設定されている「実行制限時間」を超過していないかチェックする
        if ($env['SB_LOG_LOAD_FUNCTION_EXECUTE_LIMIT_TIME'] < $elapsed_time) {
            // 超過していた場合は中断する
            // ※ 警告表示のみで正常終了とする
            logWarning("Execution time exceeded stipulated limit.");
            break;
        }
    }

    $pdo = null;

    $elapsed_time = (time() - $execute_start_time);
    logDefault(sprintf("function complete. (time: %s s)", $elapsed_time));
}

/**
 * @param Dataset $bq_dataset_
 * @param string $log_table_
 * @param string $storage_bucket_name_
 * @param PDO $pdo_
 * @param string $dataset_id_
 * @param array $files_per_directories_
 * @param string $work_time_
 * @param int $limit_time_
 * @return bool
 * @throws Throwable
 */
function LoadToBigQueryFromStorage(
    Dataset $bq_dataset_,
    string  $log_table_,
    string  $storage_bucket_name_,
    PDO     $pdo_,
    string  $dataset_id_,
    array   $files_per_directories_,
    string  $work_time_,
    int     $limit_time_
): bool
{
    $is_succeeded = true;
    $retry_count = 0;

    try {
        $table = $bq_dataset_->table($log_table_);

        $options = array();

        $directory_list = array_keys($files_per_directories_);
        foreach ($directory_list as $dir) {
            // CloudStorage に保存されているログデータのパスを生成
            $table_resource_url = 'https://storage.cloud.google.com/'
                . "{$storage_bucket_name_}/{$dir}/{$log_table_}*";

            $options['configuration']['load']['sourceUris'][] = $table_resource_url;
        }

        $load_config = $table->load(null, $options)
            ->ignoreUnknownValues(false)
            ->allowJaggedRows(false)
            ->createDisposition('CREATE_NEVER')
            ->writeDisposition('WRITE_APPEND')
            ->quote('"')
            ->schemaUpdateOptions([])
            ->fieldDelimiter(',')
            ->skipLeadingRows(0)
            ->allowQuotedNewlines(true)
            ->autodetect(false)
            ->encoding('UTF-8')
            ->sourceFormat('CSV');

        logDefault("config: " . preg_replace('/[ ' . PHP_EOL . ']+/', ' ', var_export($load_config, true)));

        $job = $table->startJob($load_config);

        while ($retry_count < MAX_RETRY_COUNT) {
            if (0 < $retry_count) {
                // 待機時間
                $delay_time = (ExponentialBackoff::calculateDelay($retry_count) / 1000000);
                $delay_end_time = time() + $delay_time;

                if ($limit_time_ < $delay_end_time) {
                    // 待機時間が最大実行時間を超過する場合は待機せず「最大回数待った」ものとする
                    logWarning(sprintf(
                        "wait time is over. (current: %s, limit: %s, delay_time: %s, delay_end_time: %s, current_retry: %s)",
                        time(),
                        $limit_time_,
                        $delay_time,
                        $delay_end_time,
                        $retry_count
                    ));
                    $retry_count = MAX_RETRY_COUNT;
                } else {
                    // 規定時間待機する
                    while (time() < $delay_end_time) {
                        // 待機時間内に処理が再開されたら戻す
                        Fiber::suspend();
                    }
                }
            }

            $job->reload();

            // 完了チェック
            if (!$job->isComplete()) {
                // 未完了の場合はリトライカウンタを上げて再待機
                ++$retry_count;
                continue;
            }

            // エラーチェック
            if (isset($job->info()['status']['errorResult'])) {
                $job_error_result = $job->info()['status']['errorResult'];
                $job_error_reason = $job_error_result['reason'];
                $job_error_message = $job_error_result['message'];

                // エラーがあった場合は Exception で回収
                throw new Exception("[$job_error_reason] $job_error_message");
            }

            // エラーなしの場合は「完了」ステータスへ更新
            updateLoadFinishedStatus(
                pdo_: $pdo_, dataset_id_: $dataset_id_,
                log_file_directories_: $files_per_directories_, work_time_: $work_time_
            );

            // エラーなしの場合はリトライループを抜ける
            break;
        }

        if (MAX_RETRY_COUNT <= $retry_count) {
            // 規定回数を超えてリトライが発生してしまった場合は
            // エラー扱いとして Exception で回収
            throw new Exception('Job has not yet completed', 500);
        }
    } catch (Exception $e) {
        // 何らかのエラーが発生した場合には Exception で回収されるので
        // Exception をキャッチしたら「エラー」ステータスへ更新する

        $errMessage = preg_replace('/[ ' . PHP_EOL . ']+/', ' ', $e->getMessage());
        logError($errMessage);

        setLoadErrorStatus(
            pdo_: $pdo_, dataset_id_: $dataset_id_,
            log_file_directories_: $files_per_directories_, err_message_: $errMessage, work_time_: $work_time_
        );

        // 「失敗」扱い
        $is_succeeded = false;
    }

    return $is_succeeded;
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
 * ログファイルのパスからキーをディレクトリパスとした連想配列を生成する
 *
 * @param $saved_log_data_
 * @param string $work_dir_
 * @param array $append_log_files_base_
 * @return array
 */
function createLogFilesMap($saved_log_data_, string $work_dir_, array $append_log_files_base_ = array()): array
{
    $log_files = $append_log_files_base_;

    $work_dir_path_for_regex = str_replace('/', '\/', $work_dir_);

    foreach ($saved_log_data_ as $log_data) {
        $pattern = "/$work_dir_path_for_regex\/(SB_LOG_.*\.gz).*/";

        if (!preg_match($pattern, $log_data->name())) {
            continue;
        }

        $log_file_name = preg_replace($pattern, '$1', $log_data->name());

        $pattern = '/(SB_LOG_.*_\d{3})_.*/';
        $log_table_name = preg_replace($pattern, '$1', $log_file_name);

        if (!isset($log_table_name)) {
            continue;
        }

        if (!isset($log_files[$log_table_name])) {
            $log_files[$log_table_name] = array();
        }

        if (!isset($log_files[$log_table_name][$work_dir_])) {
            $log_files[$log_table_name][$work_dir_] = array();
        }

        $log_files[$log_table_name][$work_dir_][] = $log_file_name;
    }

    return $log_files;
}

/**
 * 指定データセットへの指定ディレクトリの読み込み処理が
 * 既に始まっている（DBにステータスの登録がある）かを確認する
 *
 * @param PDO $pdo_
 * @param array $work_dir_list_
 * @param string $dataset_id_
 * @return array
 *  DB 上に登録済み（既に読み込み処理開始済み）のリスト
 * @throws Exception
 */
function getAlreadyStartedList(PDO $pdo_, array $work_dir_list_, string $dataset_id_): array
{
    $already_started_list = array();

    if (!empty($work_dir_list_)) {
        $union_queries = array();

        $log_status_query = preg_replace('/[ '.PHP_EOL.']+/', ' ',sprintf(
            <<<EO_SQL
            SELECT
                log_directory AS log_dir
            FROM
                %s
            WHERE
                target_dataset = '%s'
                AND
                log_directory IN ('%s')
            GROUP BY
                log_directory
            EO_SQL,
            LOG_FILE_LOAD_STATUS_TABLE_NAME,
            $dataset_id_,
            implode("', '", $work_dir_list_)
        ));

        $err_msg = "Log status insert error.";
        if (!tryRunQuery($pdo_, $log_status_query, $err_msg, $res)) {
            throw new Exception($err_msg);
        }

        $already_started_list = array_column($res->fetchAll(PDO::FETCH_ASSOC), 'log_dir');
    }

    return $already_started_list;
}

/**
 * DBに読み込み開始ステータスの登録を行う
 *
 * @param PDO $pdo_
 * @param string $dataset_id_
 * @param string $work_dir_
 * @param array $log_files_
 * @param string $work_time_
 * @return bool
 *  登録に成功で "true" を返す
 */
function tryAddLoadStartStatus(PDO $pdo_, string $dataset_id_, string $work_dir_, array $log_files_, string $work_time_): bool
{
    $insert_values = array();

    foreach ($log_files_ as $log_file) {
        $status_manage_key = substr("$work_dir_/$log_file", 0, 512);
        $uniq_key = "{$dataset_id_}:{$status_manage_key}";

        $insert_values[$uniq_key] = sprintf(
            "('%s', '%s', '%s', '%s', '%s', '%s', '%s')",
            $uniq_key,
            $work_dir_,
            $status_manage_key,
            $dataset_id_,
            STATUS_WORKING,
            $work_time_,
            $work_time_,
        );
    }

    if (empty($insert_values)) {
        return false;
    }

    $log_status_query = sprintf(
        "INSERT INTO %s ( uniq_key, log_directory, log_file, target_dataset, status, log_time, created_at ) VALUES %s",
        LOG_FILE_LOAD_STATUS_TABLE_NAME,
        implode(', ', $insert_values)
    );

    return tryRunQuery($pdo_, $log_status_query, "Log status insert error.");
}

/**
 * DBに読み込み完了ステータスの登録（更新）を行う
 *
 * @param PDO $pdo_
 * @param string $dataset_id_
 * @param array $log_file_directories_
 * @param string $work_time_
 * @return void
 */
function updateLoadFinishedStatus(PDO $pdo_, string $dataset_id_, array $log_file_directories_, string $work_time_): void
{
    $update_targets = array();
    foreach ($log_file_directories_ as $dir => $files) {
        $update_targets = array_merge(
            $update_targets,
            array_map(fn($row) => substr("$dir/$row", 0, 512), $files)
        );
    }

    $log_status_query = sprintf(
        "UPDATE %s SET status = '%s', log_time = '%s' WHERE target_dataset = '%s' AND log_time = '%s' AND log_file IN ( '%s' )",
        LOG_FILE_LOAD_STATUS_TABLE_NAME,
        STATUS_FINISHED,
        now(),
        $dataset_id_,
        $work_time_,
        implode("', '", $update_targets)
    );

    tryRunQuery($pdo_, $log_status_query, "Log status update error.");
}

/**
 * DBに読み込みエラーステータスとエラーメッセージの登録（更新）を行う
 *
 * @param PDO $pdo_
 * @param string $dataset_id_
 * @param array $log_file_directories_
 * @param string $err_message_
 * @param string $work_time_
 * @return void
 */
function setLoadErrorStatus(PDO $pdo_, string $dataset_id_, array $log_file_directories_, string $err_message_, string $work_time_): void
{
    $error_message = str_replace('\'', '\"', $err_message_);
    $conditions = array();
    foreach ($log_file_directories_ as $dir => $files) {
        $conditions = array_merge($conditions, array_map(fn($row) => substr("$dir/$row", 0, 512), $files));
    }

    $log_status_query = sprintf(
        "UPDATE %s SET status = '%s', error = '%s', log_time = '%s' WHERE target_dataset = '%s' AND log_time = '%s' AND log_file IN ( '%s' ) ",
        LOG_FILE_LOAD_STATUS_TABLE_NAME,
        STATUS_ERROR,
        $error_message,
        now(),
        $dataset_id_,
        $work_time_,
        implode("', '", $conditions)
    );

    tryRunQuery($pdo_, $log_status_query, "Log error status insert error.");
}

/**
 * DBに読み込み完了ステータスの登録（更新）を行う
 * ※ ディレクトリ単位
 *
 * @param $pdo_
 * @param $work_dir_
 * @param $dataset_id_
 * @param $has_error_
 * @return void
 */
function setDirectoryLoadStatus($pdo_, $work_dir_, $dataset_id_, $has_error_)
{
    $log_status_query = sprintf(
        "UPDATE %s SET status = '%s', log_time = '%s' WHERE target_dataset = '%s' AND log_directory = '%s'",
        LOG_DIRECTORY_LOAD_STATUS_TABLE_NAME,
        ($has_error_ ? STATUS_ERROR : STATUS_FINISHED),
        now(),
        $dataset_id_,
        $work_dir_,
    );

    tryRunQuery($pdo_, $log_status_query, "Log status insert error.");
}

/**
 * CloudSQL へのクエリの実行，エラーチェックを行う
 *
 * @param PDO $pdo_
 * @param string $query_
 * @param string $error_message_
 *  エラーが発生した際のメッセージ
 * @param PDOStatement|null $result_
 * @return bool
 *  実行に成功した場合 "true" を返す
 */
function tryRunQuery(PDO $pdo_, string $query_, string $error_message_, PDOStatement &$result_ = null): bool
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
