<?php
declare(strict_types=1);

namespace App\Console\Commands;

use Aws\CloudWatch\CloudWatchClient;
use Illuminate\Console\Command;
use Illuminate\Contracts\Filesystem\Cloud;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Storage;

/**
 * S3のファイル数をCloudWatchに送信する
 */
class PutS3LogFileNumberMetrics extends Command
{
    const LIMIT_FILE_COUNT = 2500;
    const LIMIT_DIR_DEPTH = 5;

    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'put_s3_metrics:log_files';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Put count number of log files in S3 bucket to CloudWatch metrics';

    private string $region;
    private ?Cloud $s3;

    private int $_fileCount = 0;

    public function __construct()
    {
        parent::__construct();
        $this->region = config('app.sb_log_s3_aws_region', '');
        $this->s3 = Storage::createS3Driver([
            'region' => $this->region,
            'bucket' => config('app.sb_log_s3_aws_bucket'),
            'key' => null,
            'secret' => null,
            'version' => '2006-03-01',
        ]);
    }

    /**
     * Execute the console command.
     * @return mixed
     */
    public function handle()
    {
        try {
            $this->countS3LogFiles();
            $this->putMetrics();
        }
        catch (\Exception $exception) {
            Log::error($exception);
        }
    }

    /**
     * S3のログファイル数をカウント
     * @return int
     */
    private function countS3LogFiles(): void
    {
        $prefix = config('app.sb_raw_prefix');
        $this->_summarizeFileCountNumber($prefix);
    }

    /**
     * 指定パス内にあるディレクトリのファイル数を取得する
     * @param string $path
     * @return void
     */
    private function _summarizeFileCountNumber(string $path): void
    {
        // top levelのディレクトリ取得
        $dirs = $this->s3->directories($path);
        // "elasticsearch-failed"フォルダは除外する
        $index = array_search('elasticsearch-failed', $dirs);
        if ($index !== FALSE)
        {
            array_splice($dirs, $index);
        }
        // S3_BACKUP_PREFIX フォルダは除外する
        $index = array_search(config('app.sb_backup_prefix'), $dirs);
        if ($index !== FALSE)
        {
            array_splice($dirs, $index);
        }

        foreach ($dirs as $dir_path) {
            $this->_countInDirectory($dir_path, 1);
        }
    }

    /**
     * ディレクトリに含まれるファイルの数
     * @param string $dirPath
     * @param int $depth
     * @return void
     */
    private function _countInDirectory(string $dirPath, int $depth): void
    {
        if ($depth > self::LIMIT_DIR_DEPTH) {
            return;
        }

        $dirs = $this->s3->directories($dirPath);
        foreach ($dirs as $dir) {
            $this->_countInDirectory($dir, ($depth + 1));
            if ($this->_fileCount >= self::LIMIT_FILE_COUNT) {
                return;
            }
        }

        $files = $this->s3->files($dirPath);
        $this->_fileCount += count($files);
    }

    /**
     * CloudWatchにメトリクスをputする
     * @return void
     */
    private function putMetrics()
    {
        $client = new CloudWatchClient([
            'region' => $this->region,
            'version' => '2010-08-01'
        ]);

        $metricData = [];
        $metricData[] = [
            'MetricName' => 'log_file_count',
            'Timestamp' => time(),
            'Dimensions' => [
                [
                    'Name' => 'env',
                    'Value' => config('app.env')
                ],
            ],
            'Unit' => 'Count',
            // 上限で終了した場合の誤差を丸める
            'Value' => min($this->_fileCount, self::LIMIT_FILE_COUNT),
        ];

        $client->putMetricData([
            'Namespace' => 'bq_batch',
            'MetricData' => $metricData,
        ]);
    }
}