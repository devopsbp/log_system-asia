<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Illuminate\Contracts\Filesystem\Cloud;
use Illuminate\Contracts\Filesystem\FileNotFoundException;
use Illuminate\Support\Facades\Storage;
use League\Flysystem\FileExistsException;
use Google\Cloud\BigQuery\BigQueryClient;
use Google\Cloud\Storage\StorageClient;
use Carbon\Carbon;

class S3ToBigquery extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 's3_to_bigquery {--debug} {--uuid=} {--s3-key=} {--s3-secret=} {--s3-region=} {--s3-bucket=} {--s3-root=} {--bq-project-id=} {--bq-dataset-id=} {--CRON}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 's3_to_bigquery pre_format and csv_upload';

    // .envでsb_bq_service_account.jsonを指定している
    const GCP_KEY_PATH = 'GOOGLE_APPLICATION_CREDENTIALS';

    const EXTRA_FIELD_NAME = '_extra_fields';
    const RAW_FIELD_NAME = '_fluentd_raw';
    const GCS_UPLOAD_MESSAGE = 'completed';

    // TODO デバッグ用に最小サイズを小さくしておく
    const MAX_RUN_PROCS = 1;                  // podのリソースにあわせて調整する
    const SIZE_LIMIT = 4500 * 1000 * 1000; // byte
    const TIME_LIMIT = 10 * 60;            // sec
    const CRON_TIME_LIMIT = 3 * 60;            // sec
    const FORCE_INTERRUPTION_RATE = 1.0;                // rate
    const MEM_LIMIT_RATE = 80;                 // %
    const DISK_LIMIT_RATE = 75;                 // %
    const MEM_MARGIN_PROCS = 10;                 // procs

    // 処理１回あたりの時間制限（秒）
    const TIME_LIMIT_PER_ONCE = 2 * 60;
    // 処理１回あたりの処理ファイル数
    const FILE_COUNT_LIMIT_PER_ONCE = 300;
    // S3ファイル取得時に取得ディレクトリからのファイル数がこの値より少なかった場合には
    // ディレクトリをさらに１つ超えてファイル取得を行う
    const GET_FILES_OVER_DIR_MIN_BORDER = 1000;
    // ファイルリストの更新を行う間隔（秒）
    const S3_FILE_LIST_REFRESH_INTERVAL = 15 * 60;

    const HEALTH_CHECK_LOG_PATH = 'storage/logs/s3-to-bq_health.log';

    const DEFAULT_BQ_TABLE = [
        'table' => null,
        'temp_fp' => null,
        'temp_path' => null,
        'schema' => null,
    ];

    private bool $_debug = false;
    private ?Cloud $_s3 = null;
    private array $_storage_config = [];
    private array $_bigquery_config = [];

    // 中断シグナル受信フラグ
    private bool $terminated = false;

    /**
     * Create a new command instance.
     *
     * @throws \Exception
     */
    public function __construct()
    {
        // メモリ上限を無制限に設定する
        ini_set('memory_limit', '-1');

        parent::__construct();

        // BigQueryの設定を反映する
        $this->_bigquery_config['keyFilePath'] = getenv(self::GCP_KEY_PATH);
        if (!$this->_bigquery_config['keyFilePath']
            || $this->_bigquery_config['keyFilePath'] !== getenv(self::GCP_KEY_PATH, true)
            || !is_file($this->_bigquery_config['keyFilePath'])
            || !is_readable($this->_bigquery_config['keyFilePath'])) {
            throw new \Exception("invalid gcp key : {$this->_bigquery_config['keyFilePath']}");
        }
    }

    /**
     * Execute the console command.
     *
     * @return int
     * @throws \Exception
     */
    public function handle(): int
    {
        if (extension_loaded('pcntl')) {
            $this->installSignalHandlers();
        }

        $this->_storage_config['driver'] = config('app.raw_storage_driver');
        $this->_storage_config['region'] = (config('app.sb_log_s3_aws_region') ?: ($this->option('s3-region') ?: $this->ask('aws_s3_region')));
        $this->_storage_config['bucket'] = (config('app.sb_log_s3_aws_bucket') ?: ($this->option('s3-bucket') ?: $this->ask('aws_s3_bucket')));
        $this->_storage_config['root'] = (config('app.sb_log_s3_aws_prefix') ?: ($this->option('s3-root')));

        $this->_bigquery_config['projectId'] = (config('app.sb_log_bq_gcp_project_id') ?: ($this->option('bq-project-id') ?: $this->ask('gcp_project_id')));
        $this->_bigquery_config['datasetId'] = (config('app.sb_log_bq_gcp_dataset_id') ?: ($this->option('bq-dataset-id') ?: $this->ask('gcp_dataset_id')));

        $this->_debug = (bool)$this->option('debug');

        // リソースの空き状況をチェックする(リソースに空きが無ければ終了する)
        if (!$this->_checkContinue()) {
            return 0;
        }

        $fputcsv_escape = PHP_VERSION_ID < 70400 ? "\0" : '';

        $this->logInfo("[escape] {$fputcsv_escape}");
        $this->logInfo('[format] ' . (new \DateTime())->format('Y-m-d H:i:s.u'));

        foreach ($this->_storage_config + $this->_bigquery_config as $key => $value) {
            if (!$value
                && $value !== null) {
                throw new \Exception("invalid env params : {$key}={$value}");
            }
        }

        $this->_storage_config['key'] = $this->option('s3-key') ?? null;
        $this->_storage_config['secret'] = $this->option('s3-secret') ?? null;

        $this->logInfo("[aws key (_storage_config)] {$this->_storage_config['key']}");
        $this->logInfo("[aws secret (_storage_config)] {$this->_storage_config['secret']}");

        $start_time = time();

        if (!($bq = new BigQueryClient($this->_bigquery_config))
            || !($bq_dataset = $bq->dataset($this->_bigquery_config['datasetId']))
            || !($this->_s3 = Storage::createS3Driver($this->_storage_config))) {
            throw new \Exception('GCP BigQuery invalid authentication');
        }

        // BiqQuery 上のテーブル一覧を取得しておく
        $bq_table_list = array();
        foreach ($bq_dataset->tables() as $table_info) {
            $bq_table_list[] = $table_info->info()['tableReference']['tableId'];
        }

        $_uuid = null;
        if (!($_uuid = $this->option('uuid'))) {
            if ($_uuid !== null) {
                $this->logError('error : invalid process uuid');
                return 0;
            }
            $_uuid = $this->option('CRON') ? (date('Ymd') . '_CRON_' . date('His')) : (get_current_user() . '_' . date('YmdHis'));
        }
        $this->logInfo("[uuid] {$_uuid}");

        /** テキストデータテーブル */
        $text_table_local = self::makeTextTable();

        /** 処理を行うS3上のログデータファイルリスト */
        $files = array();
        /** １件でも処理を行ったか */
        $had_processing = false;

        $bq_table_cache = array();
        $bq_schema_cache = array();

        $size_limit = (self::SIZE_LIMIT * self::FORCE_INTERRUPTION_RATE);
        $time_limit = ($this->_getTimeLimit() * self::FORCE_INTERRUPTION_RATE);

        $last_file_list_refresh_time = 0;

        while (true) {
            $pscount = 0;
            $pid = pcntl_fork();
            //pcntl_fork失敗
            if ($pid == -1) 
            {
                $this->logError('do not fork process');
                exit(-1);
            }
            //親なら
            else if ($pid)
            {
                $this->logInfo("child process pid = $pid");
                $pscount++;
                $mypid = getmypid();
                $this->logLine("parent process start: mypid = $mypid, child process pid = $pid" . $this->_memoryUsage());
            }
            //子なら
            else
            {
                $this->logLine("child process start:" . $this->_memoryUsage());

                try {
                    if (!$had_processing || empty($files) || (self::S3_FILE_LIST_REFRESH_INTERVAL < (time() - $last_file_list_refresh_time))) {
                        // 「１件も処理を行っていない／行わなかった」
                        // or 「ファイルリストが空だった」
                        // or 「ファイルリストの取得から一定時間経過した」
                        // -> S3からファイルリストの取得を行う
                        $files = $this->_searchTargetFiles(config('app.sb_raw_prefix'));
                        $last_file_list_refresh_time = time();
                    }

                    // $this->_debugInsertEOL();
                    // $this->_debugMemoryUsage();
                    $this->logInfo('find=' . count($files));

                    $success = 0;
                    $failure = 0;
                    $skip = 0;
                    $total_rows = 0;
                    $total_size = 0;
                    $bq_tables = [];
                    $skip_files = [];

                    // このタイミングでリセット
                    $had_processing = false;

                    // Cloud Storageへアップする際のフォルダを指定する
                    $uuid_prefix = $this->_generateUUID();

                    /** @var float $work_start_time S3 ファイル読み取りの開始時間 */
                    $work_start_time = microtime(true);

                    $this->logInfo('[S3] file count: ' . count($files));
                    $this->logLine('child process [S3] file count:' . count($files) . $this->_memoryUsage());

                    while ($files) {
//                      {
//                          // 対象とするS3のログサイズが 規定値 x FORCE_INTERRUPTION_RATE を超えている、
//                          // もしくは実行時間が 規定秒数 x FORCE_INTERRUPTION_RATE を超えている場合は、処理を中断する
//                          $total_elapsed_time = (time() - $start_time);
//
//                          if (($total_size > $size_limit) || ($total_elapsed_time > $time_limit)) {
//                              $this->logInfo('force interruption');
//                              $this->logInfo("over size or time. [size] {$total_size}/{$size_limit} [time] {$total_elapsed_time}/{$time_limit}");
//                              $this->terminated = true;
//                              break;
//                          }
//                      }

                        // 中断シグナルを受信した場合は処理を中断する
                        if ($this->terminated) {
                            $this->logLine('!! interrupted !!');
                            break;
                        }

                        $elapsed_time = (microtime(true) - $work_start_time);
                        $this->logInfo(sprintf("success count: %d, elapsed time: %.3f(s)", $success, $elapsed_time));

                        if ($elapsed_time > self::TIME_LIMIT_PER_ONCE) {
                            // S3 ファイルの読み取りを開始してからの時間が TIME_LIMIT を超過していたら切り上げる
                            $this->logInfo('work stripped. (by time)');

                            break;
                        }

                        if ($success >= self::FILE_COUNT_LIMIT_PER_ONCE) {
                            // ファイル処理数が FILE_COUNT_LIMIT を超えたら切り上げる
                            $this->logInfo('work stripped. (by success count)');

                            break;
                        }

                        if (!$file = array_shift($files)) {
                            // 処理対象ログファイルが存在しない -> 中断
                            break;
                        }

                        $date1 = new \DateTime();
                        $this->logLine('[file] ' . $file);

                        $work_file_path = sprintf(
                            "%s/%s/%s",
                            config('app.s3_work_directory'),
                            $uuid_prefix,
                            str_after($file, config('app.sb_raw_prefix'))
                        );

                        try {
                            // 作業ファイルをコピーする
                            $this->logLine("[file] copy '{$file}' to '{$work_file_path}'");

                            try {
                                if (!($this->_s3->copy($file, $work_file_path))) {
                                    throw ($this->_s3->exists($work_file_path)
                                        ? new \League\Flysystem\FileNotFoundException($work_file_path)
                                        : new \League\Flysystem\FileExistsException($file)
                                    );
                                }
                            } catch (\League\Flysystem\FileNotFoundException|\League\Flysystem\FileExistsException $exception) {
                                // 元ファイルがない or 移動先に既に同じファイルがある

                                $date2 = new \DateTime();
                                $this->logWarn(sprintf(
                                    "[S3] [warn#1] file already moved : %s (%s)",
                                    $file, $date2->diff($date1)->format('%R%S.%F')
                                ));
                                $this->logWarn('[S3] [warn#1] ' . $exception->getMessage());

                                $skip++;
                                $skip_files[] = $file;

                                continue;
                            }

                            // 作業ディレクトリにコピーしたので、元ファイルを削除する
                            $this->logLine("[file] delete original. '$file'");

                            if (!($this->_s3->delete($file))) {
                                // 「自分が削除する前に元ファイルが消えている」≒「別プロセスが処理した」
                                $date2 = new \DateTime();
                                $this->logWarn(sprintf(
                                    "[S3] [warn#2] original file not exist: %s (%s)",
                                    $file, $date2->diff($date1)->format('%R%S.%F')
                                ));

                                // 作業ファイルを削除する
                                if ($this->_s3->exists($work_file_path)) {
                                    $this->logLine("delete work file. '$work_file_path'");
                                    $this->_s3->delete($work_file_path);
                                }

                                $skip++;
                                $skip_files[] = $file;

                                continue;
                            }

                            // コピーした作業ファイルからデータを取得する
                            $this->logLine("[file] get data '$work_file_path'");
                            try {
                                $data = $this->_s3->get($work_file_path);
                            } catch (\Illuminate\Contracts\Filesystem\FileNotFoundException $file_not_found_exception) {
                                // ファイルが存在しなかった
                                $this->logWarn(sprintf("[S3] [warn#3] file not found. msg: %s", $file_not_found_exception->getMessage()));

                                $skip++;
                                $skip_files[] = $file;

                                continue;
                            }

                            // データの取得に失敗した
                            if (!$data) {
                                $this->logWarn("[S3] [warn#4] failed to get file data: {$work_file_path}");

                                // 作業ファイルを削除する
                                if ($this->_s3->exists($work_file_path)) {
                                    $this->logLine("delete work file. '$work_file_path'");
                                    $this->_s3->delete($work_file_path);
                                }

                                $skip++;
                                $skip_files[] = $file;

                                continue;
                            }
                        } catch (\Exception $e) {
                            $this->logError("error#0 Failed to prepare for processing. msg: {$e->getMessage()}, (from: {$e->getFile()}:L.{$e->getLine()})");

                            $skip++;
                            $skip_files[] = $file;

                            continue;
                        }

                        $next_size = strlen($data);
                        $total_size += $next_size;

                        $data = explode("\n", $data);
                        if (count($data) <= 1) {
                            $this->logError("error : not LF : {$file}");
                            $data = explode("\r", reset($data));
                        }
                        $data[0] = preg_replace('/^' . pack('H*', 'EFBBBF') . '/u', '', ltrim($data[0]), 1);
                        // $this->_debugMemoryUsage();
                        $this->logLine('lines=' . count($data));
                        $rows = 0;

                        // サフィックスにとしてテーブル名に追加するシャーディングID
                        $shard_id = substr(Carbon::today()->year, 0, 3);

                        $dbg_convert_container_log_time = 0.0;
                        $dbg_get_bq_table_info_time = 0.0;
                        $dbg_json_decode_time = 0.0;
                        $dbg_log_generate_uuid_time = 0.0;
                        $dbg_get_name_from_master_time = 0.0;
                        $dbg_apply_schema_time = 0.0;
                        $dbg_cleanup_work_files = 0.0;
                        $dbg_put_csv_time = 0.0;

                        $records = array();

                        $dbg_job_begin_time = microtime(true);
                        $dbg_job_time_first_half_begin = $dbg_job_begin_time;
                        $dbg_local_row_count = 0;
                        $dbg_local_logs = array();
                        foreach ($data as $line_key => $line_raw) {
                            $dbg_local_row_count++;
                            $dbg_local_convert_container_log_time = 0.0;
                            $dbg_local_get_bq_table_info_time = 0.0;
                            $dbg_local_json_decode_time = 0.0;
                            $dbg_local_log_generate_uuid_time = 0.0;
                            $dbg_local_get_name_from_master_time = 0.0;
                            $dbg_local_apply_schema_time = 0.0;

                            if (!$line_raw = trim($line_raw)) {
                                if (count($data) > $line_key + 1) {
                                    $this->logWarn("warning : invalid line : {$work_file_path}({$line_key})");
                                }
                                continue;
                            }

                            $work_time_begin = microtime(true);
                            $line = false;
                            $is_error = false;
                            // コンテナログへの対応
                            if ($pos = strpos($line_raw, "LogMetrics: ")) {
                                // 本来来るはずのないフォーマットなので無視する
                                $this->logWarn("warning : old format log : " . $line_raw);
                                continue;
                                // コンテナログをSBLOGフォーマットに変更する
                                // $line_raw = $this->_convertContainerLog2SBLog($line_raw);
                                // $line = json_decode($line_raw, true, 512, JSON_BIGINT_AS_STRING);
                                // $is_error = !$line || !isset($line['log_tag']);
                            } else {
                                // 一旦Json読み取り
                                $tmp_raw = json_decode($line_raw, true, 512, JSON_BIGINT_AS_STRING | JSON_UNESCAPED_UNICODE);
                                $is_error = !$tmp_raw || !isset($tmp_raw['log_tag']);
                                if ($is_error === false) {
                                    if (empty($tmp_raw['message'])) {
                                        //旧フォーマットのためSBTAG追加不要
                                        $line = $tmp_raw;
                                    } else {
                                        // SBLOG本体の抽出
                                        $line = json_decode($tmp_raw['message'], true, 512, JSON_BIGINT_AS_STRING | JSON_UNESCAPED_UNICODE);
                                        if ($line) {
                                            // SBTAG追加
                                            $line['log_tag'] = $tmp_raw['log_tag'];
                                        } else {
                                            $is_error = true;
                                        }
                                    }
                                }
                            }
                            $dbg_local_convert_container_log_time += (microtime(true) - $work_time_begin);
                            $dbg_convert_container_log_time += $dbg_local_convert_container_log_time;

                            $work_time_begin = microtime(true);
                            if ($is_error) {
                                // DEBUG用
                                $this->logLine("[line_raw] : {$line_raw}");
                                $this->logWarn("table log_tag missing {$work_file_path}({$line_key})");

                                // jsonデコードできない、もしくはlog_tagが含まれいない場合は読み飛ばす
                                continue;
                            }
                            $dbg_local_json_decode_time += (microtime(true) - $work_time_begin);
                            $dbg_json_decode_time += $dbg_local_json_decode_time;

                            // シャーディング対応のためテーブル名にシャーディングIDを追加する
                            $table_id = $line['log_tag'] . "_" . $shard_id;

                            $work_time_begin = microtime(true);
                            // 突合用のuuidを追加
                            $line['uuid'] = $this->_generateUUID();
                            $dbg_local_log_generate_uuid_time += (microtime(true) - $work_time_begin);
                            $dbg_log_generate_uuid_time += $dbg_local_log_generate_uuid_time;

                            if (($pos = strrpos($table_id, '.')) !== false) {
                                $table_id = trim(substr($table_id, $pos + 1));
                            }

                            // 定義されていないテーブルの場合は読み飛ばす
                            // TODO 定義済みテーブルリストから読み込むように修正したほうが良い
                            if ($table_id == "laravel") {
                                $dbg_local_logs[] = sprintf(
                                    '[%5d] local work time: [elapsed] %.8f(s) $table_id == "laravel"',
                                    $dbg_local_row_count,
                                    (microtime(true) - $dbg_job_begin_time),
                                );
                                continue;
                            }

                            if(!in_array($table_id, $bq_table_list)) {
                                $this->logWarn("[BQ] [warn] table '{$table_id}' does not exist in '{$this->_bigquery_config['datasetId']}'");
                                continue;
                            }

                            $work_time_begin = microtime(true);
                            if (!isset($bq_tables[$table_id])) {
                                try {
                                    // メモリ上にテーブル領域を作成する
                                    $bq_tables[$table_id] = $this->_createBqTable();

                                    if (!isset($bq_table_cache[$table_id])) {
                                        $bq_table_cache[$table_id] = $bq_dataset->table($table_id); // ←ここで500エラー食らったのでいつか例外処理いれる}
                                        $this->logInfo("[$table_id] get table cache from bq.");
                                    }
                                    $bq_tables[$table_id]['table'] = $bq_table_cache[$table_id];

                                    if (!isset($bq_schema_cache[$table_id])) {
                                        $bq_schema_cache[$table_id] = $bq_tables[$table_id]['table']->info()['schema'];
                                        $this->logInfo("[$table_id] get schema cache from bq table.");
                                    }
                                    $bq_tables[$table_id]['schema'] = $bq_schema_cache[$table_id];
                                } catch (\Exception $e) {
                                    // BigQueryにテーブルが存在しない場合はbq_tablesから削除する
                                    unset($bq_tables[$table_id]);

                                    $this->logWarn(sprintf("[BQ] failed to get bq table info. %s", $e->getMessage()));
                                    $dbg_local_get_bq_table_info_time += (microtime(true) - $work_time_begin);
                                    $dbg_local_logs[] = sprintf(
                                        '[%5d] local work time: [elapsed] %.8f(s) [get bq tbl inf] %.8f [EXCEP] %s',
                                        $dbg_local_row_count,
                                        (microtime(true) - $dbg_job_begin_time),
                                        $dbg_local_get_bq_table_info_time,
                                        preg_replace('/[ '.PHP_EOL.']+/', ' ', $e->getMessage())
                                    );

                                    $dbg_get_bq_table_info_time += $dbg_local_get_bq_table_info_time;

                                    continue;
                                }
                            }
                            $dbg_local_get_bq_table_info_time += (microtime(true) - $work_time_begin);
                            $dbg_get_bq_table_info_time += $dbg_local_get_bq_table_info_time;

                            // 名称文字列の取得／追加
                            if ($bq_tables[$table_id]) {
                                $work_time_begin = microtime(true);
                                foreach ($bq_tables[$table_id]['schema']['fields'] as $field) {
                                    if (is_array($field) && isset($field['name'])
                                        && (preg_match('/__name$/', $field['name']) == 1)) {
                                        // 「***__name」の取得に必要なキーを生成
                                        // （カラム名から "__name" を削除 ）
                                        $column_name = preg_replace('/__name$/', '', $field['name']);

                                        // キーがログデータに存在するか確認する
                                        if (!array_key_exists($column_name, $line)) {
                                            continue;
                                        }

                                        $primary_key_value = $line[$column_name];

                                        if (empty($primary_key_value)) {
                                            continue;
                                        }

                                        try {
                                            $ex_info = \App\SBBigQueryMasterInfo::getDataInfo($column_name);

                                            $master_name = $ex_info['master_name'];
                                            $key_name = $ex_info['text_key'];

                                            if (isset($text_table_local[$master_name])
                                                && isset($text_table_local[$master_name][$key_name])
                                                && isset($text_table_local[$master_name][$key_name][$line[$column_name]])
                                            ) {
                                                // 名称カラムに取得したデータを埋め込む
                                                $line["{$column_name}__name"] = $text_table_local[$master_name][$key_name][$line[$column_name]];
                                            }
                                        } catch (\Exception $exception) {
                                            $this->logError($exception->getMessage());
                                        }
                                    }
                                }
                                $dbg_local_get_name_from_master_time += (microtime(true) - $work_time_begin);
                                $dbg_get_name_from_master_time += $dbg_local_get_name_from_master_time;

                                $work_time_begin = microtime(true);
                                $line = $this->_applySchema($bq_tables[$table_id]['schema']['fields'], $line, $line_raw);
                                $dbg_local_apply_schema_time += (microtime(true) - $work_time_begin);
                                $dbg_apply_schema_time += $dbg_local_apply_schema_time;

                                $records[$table_id][] = $line;

                                $rows++;
                            }

                            // 計測した時間を保存
//                          $dbg_local_logs[] = sprintf(
//                              '[%5d] local work time: [elapsed] %.8f(s)'
//                              .' [cnv cntnr log] %.8f(s)'
//                              .' [get bq tbl inf] %.8f(s)'
//                              .' [json dec] %.8f(s)'
//                              .' [gen uuid] %.8f(s)'
//                              .' [itm nm frm Mstr] %.8f(s)'
//                              .' [aply schm] %.8f(s)',
//                              $dbg_local_row_count,
//                              (microtime(true) - $dbg_job_begin_time),
//                              $dbg_local_convert_container_log_time,
//                              $dbg_local_get_bq_table_info_time,
//                              $dbg_local_json_decode_time,
//                              $dbg_local_log_generate_uuid_time,
//                              $dbg_local_get_name_from_master_time,
//                              $dbg_local_apply_schema_time,
//                          );
                        }
                        unset($data);
                        $dbg_job_time_first_half = (microtime(true) - $dbg_job_time_first_half_begin);

                        $dbg_job_time_second_half_begin = microtime(true);
                        // 処理済みファイルの削除
                        $work_time_begin = microtime(true);
                        try {
                            $finished_path = sprintf("%s/%s", config('app.sb_backup_prefix'), str_after($file, config('app.sb_raw_prefix')));

                            // 必ず処理済みディレクトリに移す
                            if (!$this->_s3->move($work_file_path, $finished_path)) {
                                throw new FileExistsException($finished_path);
                            }
                        } catch (FileExistsException $file_exists_exception) {
                            // 同時に別プロセスも処理しており、作業が重複してしまった
                            $this->logWarn(sprintf("[warn#4] [S3] move failed. path to '%s'", $file_exists_exception->getMessage()));
                            $this->logWarn(sprintf("[warn#4] [S3] file is already processed "));

                            // 勿体ないが自分の作業は無かったことにする（早いもの勝ち）
                            $this->logWarn(sprintf("[warn#4] [S3] delete my worked file: '%s", $work_file_path));

                            // 作業していたファイルをＳ３から削除
                            $this->_s3->delete($work_file_path);

                            foreach ($records as $table_id => $rows) {
                                // csv ファイル削除
                                if (is_resource($bq_tables[$table_id]['temp_fp'])) {
                                    fclose($bq_tables[$table_id]['temp_fp']);
                                }
                                unlink($bq_tables[$table_id]['temp_fp']);
                            }
                            unset($records);

                            $dbg_cleanup_work_files += (microtime(true) - $work_time_begin);
                            $this->logInfo(sprintf("[warn#4] elapsed time: '%s", $dbg_cleanup_work_files));

                            continue;
                        } catch (\Exception $e) {
                            // error
                            $this->logError(sprintf(
                                "error#4: %s message:%s\n(%s:%s)",
                                $_uuid,
                                $e->getMessage(),
                                $e->getFile(),
                                $e->getLine()
                            ));
                            $this->logError($e->getTraceAsString());

                            unset($records);

                            $dbg_cleanup_work_files += (microtime(true) - $work_time_begin);
                            $this->logInfo(sprintf("error#4: elapsed time: '%s", $dbg_cleanup_work_files));

                            continue;
                        }
                        $dbg_cleanup_work_files += (microtime(true) - $work_time_begin);

                        // 読み込んだデータをCSVに書き込む
                        $work_time_begin = microtime(true);
                        foreach ($records as $table_id => $record_list) {
                            foreach ($record_list as $row) {
                                fputcsv($bq_tables[$table_id]['temp_fp'], $row, ",", '"', $fputcsv_escape);
                            }
                        }
                        unset($records);
                        $dbg_put_csv_time = (microtime(true) - $work_time_begin);

                        $dbg_job_time_second_half = (microtime(true) - $dbg_job_time_second_half_begin);
                        $total_job_time = (microtime(true) - $dbg_job_begin_time);

                        // 処理時間統計（デバッグ）
                        $this->logInfo(sprintf(
                            'total work time: %f'
                            .' [1st half] %s(s)'
                            .' [2nd half] %s(s)',
                            $total_job_time,
                            $dbg_job_time_first_half,
                            $dbg_job_time_second_half
                        ));
                        if($dbg_job_time_first_half > 1.0) {
                            foreach ($dbg_local_logs as $local_log) {
                                $this->logInfo($local_log);
                            }
                        }
                        unset($dbg_local_logs);

                        $this->logInfo(sprintf(
                            'total work time:'
                            .' [cnv cntnr log] %s(s)'
                            .' [get bq tbl inf] %s(s)'
                            .' [json dec] %s(s)'
                            .' [gen uuid] %s(s)'
                            .' [itm nm frm Mstr] %s(s)'
                            .' [aply schm] %s(s)'
                            .' [put csv] %s(s)'
                            .' [clnup wrk fle] %s(s)',
                            $dbg_convert_container_log_time,
                            $dbg_get_bq_table_info_time,
                            $dbg_json_decode_time,
                            $dbg_log_generate_uuid_time,
                            $dbg_get_name_from_master_time,
                            $dbg_apply_schema_time,
                            $dbg_put_csv_time,
                            $dbg_cleanup_work_files
                        ));

                        $this->logInfo("rows= {$rows}");

                        if ($rows) {
                            $total_rows += $rows;
                        }

                        $success++;
                        $had_processing = true;
                    }

                    $format_finished_time = time();
                    // $this->_debugInsertEOL();

                    $this->logInfo('total_rows=' . $total_rows);
                    // $this->_debugMemoryUsage();
                    $this->logLine('child process total_rows=' . $total_rows . $this->_memoryUsage());

                    // 処理したログがなければ戻る
                    if (!$total_rows) {
                        // 処理するログがなかった場合は、60秒待機する(S3へのアクセスが多発しないようにするため)
                        $this->logInfo('sleep 60sec');
                        $this->logLine("child process sleep before:" . $this->_memoryUsage());
                        sleep(60);
                        $this->logLine("child process sleep after:" . $this->_memoryUsage());
                        //continue;
                        //子終了
                        exit(0);
                    }

                    // メモリ上のBigQuery作業領域にBigQueryのテーブルとスキーマ情報を登録する
                    foreach ($bq_tables as $table_id => $bq_table) {
                        if (!isset($bq_table_cache[$table_id])) {
                            $bq_table_cache[$table_id] = $bq_dataset->table($table_id); // ←ここで500エラー食らったのでいつか例外処理いれる}
                            $this->logInfo("[$table_id] get table cache from bq.");
                        }
                        $bq_tables[$table_id]['table'] = $bq_table_cache[$table_id];

                        if (!isset($bq_schema_cache[$table_id])) {
                            $bq_schema_cache[$table_id] = $bq_tables[$table_id]['table']->info()['schema'];
                            $this->logInfo("[$table_id] get schema cache from bq table.");
                        }
                        $bq_tables[$table_id]['schema'] = $bq_schema_cache[$table_id];
                    }

                    $this->logLine('format_time  : ' . (string)($format_finished_time - $start_time));
                    $this->logInfo('[send] ' . (new \DateTime())->format('Y-m-d H:i:s.u'));

                    //転送したログファイルリスト
                    $log_file_list = [];

                    foreach ($bq_tables as $table_id => $bq_table) {
                        $this->logLine('table_put=' . $table_id);

                        if (!$bq_table['temp_fp']
                            || empty($filesize = fstat($bq_table['temp_fp'])['size'])) {
                            $this->logLine('empty bq_table : ' . $table_id);

                            if (is_resource($bq_table['temp_fp'])) {
                                fclose($bq_table['temp_fp']);
                            }
                            continue;
                        }

                        rewind($bq_table['temp_fp']);
                        file_put_contents($bq_table['temp_path'] = self::getCsvPath($table_id, $_uuid), $bq_table['temp_fp'], FILE_APPEND);

                        while ($bq_table['temp_path']) {
                            {
                                clearstatcache(true, $bq_table['temp_path']);

                                // CSVファイルをgzip圧縮する
                                $bq_table['temp_path_compress'] = $bq_table['temp_path'] . ".gz";
                                $csv_data = file_get_contents($bq_table['temp_path']);
                                $compress_csv_data = gzencode($csv_data, 1);
                                unset($csv_data);
                                file_put_contents($bq_table['temp_path_compress'], $compress_csv_data);
                                clearstatcache(true, $bq_table['temp_path_compress']);

                                if (is_resource($bq_table['temp_fp'])) {
                                    fclose($bq_table['temp_fp']);
                                }
                                unlink($bq_table['temp_path']);
                            }

                            {
                                // google cloud storage にアップロード

                                $storage = new StorageClient(
                                    [
                                        'projectId' => config('app.sb_log_cloud_storage_project_id'),
                                        'keyFile' => json_decode(file_get_contents($this->_bigquery_config['keyFilePath'], TRUE), true),
                                    ]
                                );
                                $bucket = $storage->bucket(config('app.sb_log_cloud_storage_bucket_name'));

                                $file_name = basename($bq_table['temp_path_compress']);
                                $options = [
                                    'name' => sprintf("%s/%s/%s", config('app.sb_log_cloud_storage_path'), $uuid_prefix, $file_name)
                                ];

                                // アップロード
                                $object = $bucket->upload(fopen($bq_table['temp_path_compress'], 'r'), $options);
                                unlink($bq_table['temp_path_compress']);

                                //転送したログファイルをリストに追加
                                $log_file_list[] = $bq_table['temp_path_compress'];
                            }

                            $this->logLine('empty check temp_path_compress=' . $bq_table['temp_path_compress']);
                            if (!empty($bq_table['temp_path_compress'])) {
                                clearstatcache(true, $bq_table['temp_path_compress']);
                                if (file_exists($bq_table['temp_path_compress'])) {
                                    $this->logLine('delete temp_path_compress=' . $bq_table['temp_path_compress']);
                                    unlink($bq_table['temp_path_compress']);
                                }
                            }
                            break;
                        }
                    }

                    {
                        //転送したログファイルリストをgoogle cloud storage に保存
                        $contents = json_encode($log_file_list);
                        $storage = new StorageClient(
                            [
                                'projectId' => config('app.sb_log_cloud_storage_project_id'),
                                'keyFile' => json_decode(file_get_contents($this->_bigquery_config['keyFilePath'], TRUE), true),
                            ]
                        );
                        $bucket = $storage->bucket(config('app.sb_log_cloud_storage_bucket_name'));

                        $stream = fopen('data://text/plain,' . self::GCS_UPLOAD_MESSAGE, 'r');
                        $objectName = sprintf("%s/%s/%s", config('app.sb_log_cloud_storage_path'), $uuid_prefix, config('app.sb_log_cloud_storage_log_file_list'));
                        $bucket->upload($stream, [
                            'name' => $objectName,
                        ]);
                        $this->logInfo('upload log file list: ' . $objectName);
                    }

                    // $this->_debugInsertEOL();
                    $this->logInfo('send to gcloud: ' . config('app.sb_log_cloud_storage_path') . '/' . $uuid_prefix);
                    $this->logInfo("[s3] success:{$success}, skip:{$skip}, failure:{$failure}");
                    $this->logLine(sprintf('total_received_size: %.2f MiB', $total_size / 1024 / 1024));
                    $this->logLine(sprintf('peak_memory_consume: %.2f MiB', memory_get_peak_usage() / 1024 / 1024));

                    if (!$this->_checkContinue()) {
                        $this->terminated = true;
                    }
                } catch (\Exception $e) {
                    $this->logError("error#2: {$_uuid} message:{$e->getMessage()}({$e->getFile()}:{$e->getLine()})");
                }

                if (isset($bq_tables)) {
                    foreach ($bq_tables as $bq_table) {
                        if ($bq_table
                            && !empty($bq_table['temp_path'])
                            && file_exists($bq_table['temp_path'])) {
                            clearstatcache(true, $bq_table['temp_path']);
                            if (file_exists($bq_table['temp_path'])) {
                                unlink($bq_table['temp_path']);
                            }
                        }
                    }
                }

                // スキップファイルを除外する
                $files = array_diff($files, $skip_files);

                // 中断シグナルを受信した場合は処理を中断する
                if ($this->terminated) {
                    $this->logLine('!! interrupted !!');
                    $this->logLine("child process interrupted:" . $this->_memoryUsage());
                    //break;
                    exit(-1);
                }

                $this->logLine("child process finish:" . $this->_memoryUsage());
                //子正常終了
                exit(0);
            }

            //子が終了するまで待つ
            while ($pscount > 0){
                $this->logLine("parent process waitpid before:" . $this->_memoryUsage());
                pcntl_waitpid( -1, $status, WUNTRACED);
                $this->logLine("parent process waitpid after:" . $this->_memoryUsage());
                $this->logInfo("status = $status pscount = $pscount");
                $pscount--;
            }

            $this->logInfo("[SUCCESS] finished");

            system('uptime' . PHP_EOL);
        }
        return 0;
    }

    /**
     * メモリ使用量
     * @return string
     */
    private function _memoryUsage(): string
    {
        return sprintf( " [memcheck] %d b", memory_get_usage() );
    }

    /**
     * シグナルハンドラを登録
     * @return void
     */
    private function installSignalHandlers(): void
    {
        pcntl_async_signals(true);
        $handler = fn() => $this->terminated = true;
        foreach ([SIGTERM, SIGINT, SIGQUIT] as $signal) {
            pcntl_signal($signal, $handler);
        }
    }

    /**
     * 空きリソースがあるかの判定をおこなう。
     *
     * @return bool true リソースの空きがある/false リソースの空きが無い
     */
    private function _checkContinue(): bool
    {
        $this->logInfo('usleep(mt_rand(0, 500000));');
        // TODO 連続失敗回数から指数バックオフの遅延をおこなうようにした方がよい
        usleep(mt_rand(0, 500000)); // ウェイトさせる(マイクロ秒) 0～0.5秒のウェイト

        $proc_num = self::getPhpProcessNum();
        list($mem_use_rate, $mem_free) = self::getPhpMemory();
        $disk_use_rate = $this->_getPrivateDiskUseRate();

        $check_process_num = (int)$proc_num / self::MAX_RUN_PROCS * 100;
        $check_memory_usage = (int)$mem_use_rate / self::MEM_LIMIT_RATE * 100;
        $check_disk_usage = (int)$disk_use_rate / self::DISK_LIMIT_RATE * 100;

        // ログファイルがオープンできるまで0.1秒ウェイトしつつ待つ
        while (!$fp = fopen(self::HEALTH_CHECK_LOG_PATH, 'ab')) {
            $this->logInfo('usleep(100000);');
            usleep(100000);
        }
        fwrite($fp, sprintf('%.2f    %.2f    %.2f', $check_process_num, $check_memory_usage, $check_disk_usage) . PHP_EOL);
        fclose($fp);

        // リソース制限に達しているかのチェックをおこなう（現在process_numは未使用）
        //if ($check_process_num >= 120
        if ($check_memory_usage >= 100
            // || $check_memory_available >= 100
            || $check_disk_usage >= 100) {
            $this->logWarn('warning : ' . (int)$mem_use_rate . '% memory usage, ' . (int)$disk_use_rate . '% disk useage');
            //if ($check_process_num >= 120) {
            //    $this->logWarn('fail : check_process_num ' . $check_process_num . ' >= 120');
            //}

            if ($check_memory_usage >= 100) {
                $this->logWarn('fail : check_memory_usage ' . $check_memory_usage . ' >= 100');
            }

            if ($check_disk_usage >= 100) {
                $this->logWarn('fail : check_disk_usage ' . $check_disk_usage . ' >= 100');
            }
            return false;
        }
        return true;
    }

    /**
     * 指定パス内にあるディレクトリを取得する
     *
     * @param string $path
     * @return array
     */
    private function _searchTargetFiles(string $path): array
    {
        $dirs = $this->_s3->directories($path);

        if (count($dirs)) {
            // "elasticsearch-failed"フォルダは除外する
            $index = array_search('elasticsearch-failed', $dirs);
            if ($index !== FALSE) {
                array_splice($dirs, $index);
            }

            // S3_BACKUP_PREFIX フォルダは除外する
            $index = array_search(config('app.sb_backup_prefix'), $dirs);
            if ($index !== FALSE) {
                array_splice($dirs, $index);
            }

            // work フォルダは除外する
            $index = array_search(config('app.s3_work_directory'), $dirs);
            if ($index !== FALSE) {
                array_splice($dirs, $index);
            }
        }

        // 返却するファイルリスト
        $res_files = array();

        if ($dirs) {
            foreach ($dirs as $dir_path) {
                if ($files = $this->_searchTargetFiles($dir_path)) {
                    // 抽出したファイル群から対象ログファイルのみ返却ファイルリストに追加する
                    foreach ($files as $key => $file) {
                        if (str_contains($file, config('app.s3_file_prefix'))) {
                            $res_files[] = $file;
                        }
                    }
                }

                if (self::GET_FILES_OVER_DIR_MIN_BORDER < count($res_files)) {
                    break;
                }
            }

            shuffle($res_files);
        } else {
            $res_files = $this->_s3->files($path);
        }

        if ($this->_debug) {
            print('.');
        }

        return $res_files;
    }

    private function _sortByLastmodify($files)
    {
        // ファイル名から抜粋した時刻情報を格納する配列
        $file_lastmodify_array = array();

        // ファイル名から時刻情報を抜粋する
        foreach ($files as $file) {
            preg_match('/[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-[0-9][0-9]-[0-9][0-9]-[0-9][0-9]/', $file, $last_modify);
            $file_lastmodify_array[] = $last_modify[0];
        }

        // 抜粋した時刻情報を元にファイルリストをソートする
        array_multisort($file_lastmodify_array, SORT_ASC, $files);

        return ($files);
    }

    /**
     * @return float|int
     */
    private function _getTimeLimit()
    {
        if ($this->option('CRON')) {
            return self::CRON_TIME_LIMIT;
        }
        return self::TIME_LIMIT;
    }

    private function _debugMemoryUsage(): void
    {
        if (!$this->_debug) {
            return;
        }
        $this->comment('memory_usage : ' . memory_get_usage());
    }

    private function _debugInsertEOL(): void
    {
        if (!$this->_debug) {
            return;
        }
        $this->comment('');
    }

    /**
     * メモリ上にBigQueryに投入するデータを準備する一時テーブルを作成する
     *
     * @param bool $is_temp = true
     * @return array
     */
    private function _createBqTable(bool $is_temp = true): array
    {
        $bq_table = self::DEFAULT_BQ_TABLE;
        $bq_table['temp_fp'] = ($is_temp ? fopen('php://memory', 'r+b') : false);
        return $bq_table;
    }

    /**
     * @param $value
     * @param string $type
     * @return false|float|int|string|null
     * @throws \Exception
     */
    private function _preFormatData($value, string $type)
    {
        switch (true) {
            case $value === null:
                return '';
            case is_array($value):
            case is_object($value):
                if (($json = json_encode($value)) === false
                    && ($json = json_encode(trim($value))) === false) {
                    return '';
                }
                return $json;
            case is_bool($value):
                // 'false' が空文字になってしまうので明示的に返却
                return ($value ? 'true' : 'false');
            case $type == 'TIMESTAMP':
                if ($value
                    && is_numeric($value)
                    && (int)$value) {
                    return (new \DateTime())->setTimezone(new \DateTimeZone('UTC'))->setTimestamp((int)$value)->format('Y-m-d H:i:s');
                }
                if ($datetime = \DateTime::createFromFormat('Y.m.d-H.i.s', $value)) {
                    return $datetime->setTimezone(new \DateTimeZone('UTC'))->format('Y-m-d H:i:s');
                }
                if (strtotime($value)) {
                    return (new \DateTime($value))->setTimezone(new \DateTimeZone('UTC'))->format('Y-m-d H:i:s');
                }
                if ($value === 0
                    || trim($value) === '0') {
                    return '1970-01-01 00:00:00';
                }
                break;
            case $type == 'INTEGER':
                if (!is_numeric($value)) {
                    return '';
                }
                if ((int)$value) {
                    return (int)$value;
                }
                if ($value === 0
                    || trim($value) === '0') {
                    return 0;
                }
                return '';
            case $type == 'FLOAT':
            case $type == 'NUMERIC':
                if (!is_numeric($value)) {
                    return '';
                }
                if ((float)$value) {
                    return $value;
                }
                if ($value === 0.0
                    || trim($value) === '0.0'
                    || $value === 0
                    || trim($value) === '0') {
                    return 0.0;
                }
                return '';
            case $type == 'STRING':
                return (string)$value;
            default:
                return $value;
        }
        return $value;
    }

    /**
     * S3のデータをBigQueryのスキーマを適合させる
     *
     * @param array $fields
     * @param array $line
     * @param string $line_raw
     * @return array
     */
    private function _applySchema(array $fields, array $line, string $line_raw): array
    {
        $result = [];
        $extra_key = null;
        foreach ($fields as $key => $field) {
            if ($field['name'] == self::RAW_FIELD_NAME) {
                $result[] = $line_raw;
                continue;
            }

            if ($field['name'] == self::EXTRA_FIELD_NAME) {
                $result[] = '{}';
                $extra_key = $key;
                continue;
            }

            if (isset($line[$field['name']])) {
                $result[] = $this->_preFormatData($line[$field['name']], $field['type']);
                array_pull($line, $field['name']);
                continue;
            }

            $result[] = '';
        }

        if ($extra_key) {
            array_splice($result, $extra_key, 1, json_encode($line ?: new \stdClass()));
        }

        return $result;
    }

    /**
     * PHPのプロセスリストを返す
     *
     * @return array プロセス数 TODO もう一つはなんだ？
     */
    public static function getPhpProcess(): array
    {
        return explode(',', exec('ps -ww ax o comm= o rss= | grep "^php[0-9 ]\\?" | awk \'{v+=$2;n++}END{print (n?n:0)","(n?v/n:0)}\''));
    }

    /**
     * BigQueryに読み込ませるCSVファイルのパスを返す
     *
     * @param string $table_id
     * @param string $process_uuid
     * @param bool $send_flg = false
     * @return string CSVファイルのパス
     */
    public static function getCsvPath(string $table_id, string $process_uuid, bool $send_flg = false): string
    {
        return storage_path('app/s3_to_bq/' . ($send_flg ? 'send/' : '') . $table_id . '_' . $process_uuid . '.csv');
    }

    /**
     * 実行中のプロセス数を返す。
     *
     * @return string
     */
    public static function getPhpProcessNum(): string
    {
        // exec('fuser -v ' . PHP_BINARY . ' 2>&1 | grep "php[0-9 ]\\?" | wc -l')
        return exec('ps -ww ax o comm= o pid= | grep "^php[0-9 ]\\?" | wc -l');
    }

    /**
     * メモリの使用率と空きメモリ容量を返す。
     *
     * @return array
     */
    public static function getPhpMemory(): array
    {
        return explode(',', exec('free -tk | grep Total | awk \'{print ($3/($2?$2:1)*100)","$4}\''));
    }

    /**
     * 指定したパスのストレージの未使用率を返す。
     *
     * @return float
     */
    private function _getPrivateDiskUseRate(): float
    {
        // exec("df -P | sed -e's/^.* \([0-9]*\). .*$/\\1/g' | sort -nr | head -1");
        $total = disk_total_space(storage_path());
        return ($total - disk_free_space(storage_path())) / $total * 100;
    }

    /**
     * コンテナログフォーマットをSB_LOG形式に変換する
     *
     * @param string $line_raw
     * @return string
     */
    private function _convertContainerLog2SBLog(string $line_raw): string
    {
        // 余計な"}"と改行コードを除去する
        $line_raw = str_replace('}\\\\n\"', '', $line_raw);
        // エスケープを除去する
        $line_raw = str_replace('\\', '', $line_raw);

        // 改行、エスケープを除去しているのでstrposを取得し直す
        $pos = strpos($line_raw, "LogMetrics: ");
        $line_raw = trim(substr($line_raw, $pos + strlen("LogMetrics: ")));

        // log_tagを取得する
        $pos = strpos($line_raw, " ");
        $log_tag = trim(substr($line_raw, 0, $pos));

        $line_raw = trim(substr($line_raw, strlen($log_tag) + 1));
        // 先頭が"\""であれば除去する
        if (substr($line_raw, 0, 1) === '"') {
            $line_raw = trim(substr($line_raw, 1, strrpos($line_raw, "\"}") - 1));
        } else {
            $line_raw = trim(substr($line_raw, 0, strrpos($line_raw, "\"}")));
        }

        $sb_log_tag = '"log_tag":"' . $log_tag . '"';
        $time = trim(substr($line_raw, strpos($line_raw, "\"time\"") + strlen("\"time\"") + 1));
        $time = explode("\"", $time)[1];
        $log_time = str_replace("T", " ", explode(".", $time)[0]);

        $line_raw = str_replace("\"time\":\"$time\"", "\"log_time\":\"$log_time\",\"log_tag\":\"sb.game.$log_tag\"", $line_raw);
        return ($line_raw);
    }

    /**
     * UUID生成
     * ※Laravel 5.6以降にアップデートする場合は、Str::uuid()に変更する
     */
    private function _generateUUID(): string
    {
        $pattern = "xxxxxxxx_xxxx_4xxx_yxxx_xxxxxxxxxxxx";
        $character_array = str_split($pattern);
        $uuid = "";

        foreach ($character_array as $character) {
            switch ($character) {
                case "x":
                    $uuid .= dechex(random_int(0, 15));
                    break;
                case "y":
                    $uuid .= dechex(random_int(8, 11));
                    break;
                default:
                    $uuid .= $character;
            }
        }

        return ($uuid);
    }

    /**
     * テキストマスタデータ諸々をＤＢから取得し、ローカル向けのテキストテーブルを作成する
     *
     * @return array
     */
    private static function makeTextTable(): array
    {
        // マスタデータ取得のため、 main DB へ接続
        $con = \DB::connection('main_db');

        $datas = array();
        foreach (\App\SBBigQueryMasterInfo::MASTER_LIST as $master_info) {
            if (isset($datas[$master_info['master_name']][$master_info['text_key']])) {
                // 既に同じキーで設定されていたらスキップ
                continue;
            }

            if (empty($master_info['master_name'])
                || empty($master_info['key'])
                || empty($master_info['text_key'])
                || empty($master_info['text_table_name'])) {
                // 情報が足りなければスキップ
                continue;
            }

            // テキストテーブルを取得
            $local_table = $con->table($master_info['master_name'])
                ->join('master_text_data', "{$master_info['master_name']}.{$master_info['text_key']}", '=', 'master_text_data.text_id')
                ->where('master_text_data.text_table_name', $master_info['text_table_name'])
                ->where('master_text_data.locale', config('app.sb_log_text_location'))
                ->select("{$master_info['master_name']}.{$master_info['key']} AS key_id", 'master_text_data.text AS text')
                ->get();

            // id をキーにしてデータ保管
            $json = json_encode($local_table);
            $local_table = array_column(json_decode($json, true), 'text', 'key_id');

            $datas[$master_info['master_name']][$master_info['text_key']] = $local_table;
        }

        return $datas;
    }

    public function logLine($string, $style = null, $verbosity = null)
    {
        $this->line(self::getFormattedTime() . " {$string}", $style, $verbosity);
    }

    public function logInfo($string, $verbosity = null)
    {
        $this->info(self::getFormattedTime() . " $string", $verbosity);
    }

    public function logWarn($string, $verbosity = null)
    {
        $this->warn(self::getFormattedTime() . " $string", $verbosity);
    }

    public function logError($string, $verbosity = null)
    {
        $this->error(self::getFormattedTime() . " $string", $verbosity);
    }

    private static function getFormattedTime(\DateTime $dateTime = null): string
    {
        if (!isset($dateTime)) {
            $dateTime = new \DateTime();
        }

        return ("[{$dateTime->format('Y-m-d H:i:s.u')}]");
    }
}