<?php
require_once('DatabaseUnix.php');

use CloudEvents\V1\CloudEventInterface;
use Google\Cloud\Samples\CloudSQL\Postgres\DatabaseUnix;
use Google\Cloud\Storage\StorageClient;
use Google\CloudFunctions\FunctionsFramework;

const STATUS_FINISHED = 'finished';

const LOG_DIRECTORY_LOAD_STATUS_TABLE_NAME = 'sb_logs_cloud_to_bq_work_directory';

const REMOVE_ACTION_LIST = [
    'NONE'   => 'NONE',   // 特に処理を行わない（ログ出力のみ行う）
    'MOVE'   => 'MOVE',   // 'SB_LOG_CLOUD_STORAGE_LOADED_FILE_PATH' で設定されたディレクトリへ移動
    'DELETE' => 'DELETE', // 削除
];

// Register the function with Functions Framework.
// This enables omitting the `FUNCTIONS_SIGNATURE_TYPE=cloudevent` environment
// variable when deploying. The `FUNCTION_TARGET` environment variable should
// match the first parameter.
FunctionsFramework::cloudEvent('DeleteExpiredLogs', 'DeleteExpiredLogs');

function DeleteExpiredLogs(CloudEventInterface $event): void
{
    // -- 環境設定の読み込み ---
    $env = setEnv();

    // -- CloudSQL 向けに環境変数を設定 ---
    putenv("DB_USER={$env['SB_LOG_DB_USER']}");
    putenv("DB_PASS=".trim(getenv('SB_DB_PW')));
    putenv("DB_NAME={$env['SB_LOG_DB_NAME']}");
    putenv("INSTANCE_UNIX_SOCKET=/cloudsql/{$env['SB_LOG_DB_INSTANCE_UNIX_SOCKET']}");

    // ---

    // Google Cloud Storage からログデータを取得
    $storage = new StorageClient(
        [
            'projectId' => $env['SB_LOG_CLOUD_STORAGE_PROJECT_ID'],
        ]
    );

    $bucket = $storage->bucket($env['SB_LOG_CLOUD_STORAGE_BUCKET_NAME']);

    logDefault("[GCS] bucket= {$env['SB_LOG_CLOUD_STORAGE_BUCKET_NAME']}");
    logDefault("[GCS] prefix= {$env['SB_LOG_CLOUD_STORAGE_PATH']}");

    $pdo = DatabaseUnix::initUnixDatabaseConnection();
    $finished_dir_list = createFinishedList(pdo_: $pdo);

    // 「削除」のアクションを環境変数から取得
    $finished_action = strtoupper($env['SB_LOG_FINISHED_REMOVE_ACTION']);
    if (!in_array($finished_action, array_keys(REMOVE_ACTION_LIST))) {
        $finished_action = REMOVE_ACTION_LIST['NONE'];
    }

    $worked_list = array();
    foreach ($finished_dir_list as $index => $dir_path) {
        $move_list = $bucket->objects(['prefix' => $dir_path]);

        $object_list = iterator_to_array($move_list);
        if (count($object_list) == 0) {
            // リストが空の場合はスキップ
            logDefault("[CloudStorage] [$dir_path] no delete targets.");
            $worked_list[] = $dir_path;
            continue;
        }

        logDefault(sprintf(
            "[CloudStorage] directory '%s' is deletable. ( files: %s )",
            $dir_path,
            implode(', ', array_map(fn($row) => preg_replace('/^.*\//', '', $row->name()), $object_list))
        ));

        if (REMOVE_ACTION_LIST['NONE'] == $finished_action) {
            // 何もしない設定ならここまで
            continue;
        }

        $preg_str = str_replace('/', '\/', $env['SB_LOG_CLOUD_STORAGE_RETRY_FILE_PATH']);
        $replace_str = (
            preg_match("/^$preg_str/", $dir_path)
                ? $env['SB_LOG_CLOUD_STORAGE_RETRY_FILE_PATH']
                : $env['SB_LOG_CLOUD_STORAGE_PATH']
        );

        $preg_str = str_replace('/', '\/', $replace_str);

        foreach ($move_list as $item) {
            try {
                if (REMOVE_ACTION_LIST['MOVE'] == $finished_action) {
                    $new_path = preg_replace("/^$preg_str/", $env['SB_LOG_CLOUD_STORAGE_LOADED_FILE_PATH'], $item->name());

                    logDefault("[CloudStorage] [$dir_path] move file '{$item->name()}' to '$new_path'");

                    $item->copy(
                        $env['SB_LOG_CLOUD_STORAGE_BUCKET_NAME'],
                        ['name' => $new_path]
                    );
                }

                if ((REMOVE_ACTION_LIST['MOVE'] == $finished_action
                    || REMOVE_ACTION_LIST['DELETE'] == $finished_action)) {
                    logDefault("[CloudStorage] [$dir_path] delete '{$item->name()}'");
                    $item->delete();
                }
            } catch (\Exception $excp) {
                // ジョブが並列で実行された時、削除しようとしたディレクトリが他のジョブで既に削除されている可能性がある
                logWarning('exception!! [msg] ' . $excp->getMessage());
            }
        }

        $worked_list[] = $dir_path;
    }

    // 移動／削除 処理が終わったディレクトリは DB からも削除する
    removeWorkedDirectoryOnDB(pdo_: $pdo, remove_dir_list_: $worked_list);

    // コネクションを明示的に切断する
    $pdo = null;
    logDefault("function complete.");
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

function createFinishedList(PDO $pdo_): array
{
    $log_status_query = preg_replace(
        '/[ ' . PHP_EOL . ']+/', ' ',
        sprintf(
            <<<EOC
            SELECT
                tbl.log_directory
            FROM (
                SELECT 
                    log_directory, 
                    COUNT((status = '%s') or null) AS finished_count, 
                    COUNT(status)                  AS total_count 
                FROM 
                    %s
                GROUP BY 
                    log_directory
            ) AS tbl
            WHERE
                ( tbl.finished_count = tbl.total_count )
            EOC,
            STATUS_FINISHED,
            LOG_DIRECTORY_LOAD_STATUS_TABLE_NAME
        )
    );

    $err_msg = "Failed to get finished directories.";
    if (!tryRunQuery($pdo_, $log_status_query, $err_msg, $res)) {
        throw new Exception($err_msg);
    }

    $res_array = array_column($res->fetchAll(PDO::FETCH_ASSOC), 'log_directory');

    // 並列で実行された時に処理対象の取り合いにならないようにシャッフルする
    shuffle($res_array);

    return $res_array;
}

function removeWorkedDirectoryOnDB(PDO $pdo_, array $remove_dir_list_): void
{
    if(empty($remove_dir_list_)) {
        return;
    }

    $log_status_query = preg_replace(
        '/[ ' . PHP_EOL . ']+/', ' ',
        sprintf(
            <<<EOC
            DELETE
            FROM
                %s
            WHERE
                log_directory IN ( '%s' )
            EOC,
            LOG_DIRECTORY_LOAD_STATUS_TABLE_NAME,
            implode("', '", $remove_dir_list_ )
        )
    );

    logDefault('[' . LOG_DIRECTORY_LOAD_STATUS_TABLE_NAME . '] SQL= ' . $log_status_query);

    $err_msg = "Failed to remove finished directories.";
    if (!tryRunQuery($pdo_, $log_status_query, $err_msg)) {
        throw new Exception($err_msg);
    }
}

/**
 * CloudSQL へのクエリの実行，エラーチェックを行う
 *
 * @param $pdo_
 * @param $query_
 * @param $error_message_
 *  エラーが発生した際のメッセージ
 * @param $result_
 * @return bool
 *  実行に成功した場合 "true" を返す
 */
function tryRunQuery($pdo_, $query_, $error_message_, &$result_ = null): bool
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
