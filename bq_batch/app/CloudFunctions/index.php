<?php

use Psr\Http\Message\ServerRequestInterface;

use Google\CloudFunctions\CloudEvent;
use Google\Cloud\BigQuery\BigQueryClient;
use Google\Cloud\Core\ExponentialBackoff;
use Google\Cloud\Storage\StorageClient;

use Google\Cloud\SecretManager\V1\SecretManagerServiceClient;

const MAX_RUN_PROCS              = 30;
const BIGQUERY_TIMEOUT           = 60 * 60;
const BIGQUERY_WAIT_DEFAULT      = 3;       // テーブル メタデータ更新オペレーションの最大レート - テーブルあたり10秒ごとに5件のオペレーション
const BIGQUERY_PROJECT_DAY_LIMIT = 100000;  // 1 日あたりのプロジェクトあたり読み込みジョブ数（失敗を含む）
const BIGQUERY_TABLE_DAY_LIMIT   = 1500;    // 1日あたりの最大テーブル オペレーション数

const STATUS_NEW        = 'new';
const STATUS_LOCKED     = 'locked';     // logs only
const STATUS_FINISHED   = 'finished';
const STATUS_TERMINATED = 'terminated'; // processes only

const JOB_EXCEEDED_KEY = 'job_exceeded';

function LoadToBQ(CloudEvent $event): void
{
	// -- 環境設定の読み込み ---
	$data = file(getenv('SB_SECRET_ID'));

	$env = array();
	foreach ($data as $row)
	{
		$split = explode('=', $row);

		$key = $split[0];
		$val = rtrim($split[1]);

		if (str_starts_with($key, 'ARR_'))
		{
			$val = explode(',', $val);
		}

		$env[$split[0]] = $val;
	}

	// ---

	$key_file_path = getenv('GCP_KEY_PATH');

	// Google Cloud Storage からログデータを取得
	$storage = new StorageClient(
		[
			'projectId' => $env['SB_LOG_CLOUD_STORAGE_PROJECT_ID'],
			'keyFile' => json_decode(file_get_contents($key_file_path, TRUE), true),
		]
	);

	$bucket = $storage->bucket($env['SB_LOG_CLOUD_STORAGE_BUCKET_NAME']);
	$saved_log_data = $bucket->objects(
		[
			'prefix' => "{$env['SB_LOG_CLOUD_STORAGE_PATH']}/"
		]
	);

	$uuid_ = uniqid((new \DateTime())->format('YmdHis_'), true);

	if (!is_file($key_file_path)
		|| !is_readable($key_file_path))
	{
		throw new \Exception('invalid gcp key : ' . $key_file_path);
	}

	foreach ( $env['ARR_SB_LOG_BQ_GCP_DATASET_ID'] as $dataset_id )
	{
		$bq_config = [];

		// BigQueryの設定を反映する
		$bq_config['keyFilePath'] = $key_file_path;

		$bq_config['projectId'] = $env['SB_LOG_BQ_GCP_PROJECT_ID'];
		$bq_config['datasetId'] = $dataset_id;

		if (!isset($bq_config['projectId']) || !isset($bq_config['datasetId']))
		{
			throw new \Exception('invalid gcp settings.');
		}

		$bq_client = new BigQueryClient($bq_config);
		$bq_dataset = $bq_client->dataset($bq_config['datasetId']);

		foreach ($saved_log_data as $log_data)
		{
			try
			{
				// CloudStorage に保存されているログデータのパスを生成
				$table_resource_url = 'https://storage.cloud.google.com/'
					. $env['SB_LOG_CLOUD_STORAGE_BUCKET_NAME'] . '/'
					. $log_data->name();

				preg_match('/SB_LOG_.*/', $table_resource_url, $match_pattern);

				if (empty($match_pattern))
				{
					logWarning("$table_resource_url is invalid name format.");
					continue;
				}

				// ファイル名が「SB_LOG_<ログ名>_<シャードサフィックス（年頭３桁）>_<ログ作成時の情報>.<拡張子>」となっているので、
				// 先頭の「SB_LOG_<ログ名>_<シャードサフィックス（年頭３桁）>」（＝データセット名）を取り出す
				$table_id = preg_replace('/(_\d{3})(_.*)/', '$1', $match_pattern[0]);

				$status_manage_key = substr($log_data->name(), 0, 512);

				$target_table_ = "{$bq_config['projectId']}.{$bq_config['datasetId']}.cloud_to_bq_sb_logs";

				$bq_status_query = "SELECT `status`, `log_time` FROM `$target_table_` WHERE `log_file` = '$status_manage_key'";

				$bq_status_job = $bq_client->query($bq_status_query);
				$bq_status_job_res = $bq_client->runQuery($bq_status_job);

				$bq_status_job_res_row = $bq_status_job_res->rows()->current();

				if (empty($bq_status_job_res_row) == false)
				{
					if ($bq_status_job_res_row['status'] == STATUS_FINISHED)
					{
						// 処理済み
						continue;
					}

					$pend_limit_time = new DateTime($bq_status_job_res_row['log_time'] . ' +3600 seconds');
					$current_time = new DateTime();

					if ($pend_limit_time < $current_time)
					{
						// 読み込み中の可能性あり
						continue;
					}

					// 「処理済みでない」且つ「前回更新時刻から 3600s 経過している」
					$bq_status_query = sprintf(
						"UPDATE `%s` SET `status` = '%s', `log_time` = '%s' WHERE `log_file` = '%s'"
						, $target_table_
						, STATUS_NEW
						, $current_time->format('Y-m-d H:i:s')
						, $status_manage_key
					);

					$bq_status_job = $bq_client->query($bq_status_query);
					$bq_status_job_res = $bq_client->runQuery($bq_status_job);
				} else
				{
					$format_current_time = (new \DateTime())->format('Y-m-d H:i:s');

					$bq_status_query = sprintf(
						"INSERT INTO `%s` ( `log_file`, `status`, `log_time`, `created_at` ) VALUES ( '%s', '%s', '%s', '%s' )"
						, $target_table_
						, $status_manage_key
						, STATUS_NEW
						, $format_current_time
						, $format_current_time
					);

					$bq_status_job = $bq_client->query($bq_status_query);
					$bq_status_job_res = $bq_client->runQuery($bq_status_job);
				}

				$table = $bq_dataset->table($table_id);
				$schema = $table->info()['schema'];

				$load_config = $table->loadFromStorage($table_resource_url)
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

				logDefault("table_job={$table_id}");

				$job = $table->runJob($load_config);

				$backoff = new ExponentialBackoff(100);
				$backoff->execute(
					function () use ($job)
					{
						$job->reload();
						if (!$job->isComplete())
						{
							throw new \Exception('Job has not yet completed', 500);
						}
					}
				);

				if (isset($job->info()['status']['errorResult']))
				{
					$job_error_result = $job->info()['status']['errorResult'];
					$job_error_reason = $job_error_result['reason'];

					logWarning("warning : {$job_error_reason}");

					$all_exceeded_num = 1;
					$table_exceeded_num = 1;

					$retry_count = 0;
					$search_time = time();

					$project_exceeded_wait = BIGQUERY_WAIT_DEFAULT;
					$table_id_prefix = substr($table_id, 0, 512);

					{
						$target_table_ = sprintf(
							"%s.%s.%s",
							$bq_config['projectId'],
							$bq_config['datasetId'],
							'cloud_to_bq_exceeded'
						);

						$bq_status_query = sprintf(
							"SELECT COUNT(1) `all`, COUNT(`target` = '%s' OR NULL) `table` FROM `" . $target_table_ . "` WHERE `key` = '%s' AND `timestamp` > %s ",
							$table_id_prefix, JOB_EXCEEDED_KEY, time()
						);

						$bq_status_job = $bq_client->query($bq_status_query);
						$bq_status_job_res = $bq_client->runQuery($bq_status_job);

						$res_current = $bq_status_job_res->rows()->current();

						if (!empty($res_current))
						{
							$all_exceeded_num = (int)$res_current['all'] ?: 1;
							$table_exceeded_num = (int)$res_current['table'] ?: 1;
						}
					}

					if ($job_error_reason == 'quotaExceeded')
					{
						logWarning("warning : limit exceeded ... wait {$project_exceeded_wait} sec");
						logDefault("waiting_files={$all_exceeded_num}({$table_exceeded_num})");

						$project_exceeded_wait = (int)max((24 * 60 * 60) / (BIGQUERY_PROJECT_DAY_LIMIT / $all_exceeded_num)
							(24 * 60 * 60) / (BIGQUERY_TABLE_DAY_LIMIT / $table_exceeded_num), BIGQUERY_WAIT_DEFAULT);
						$project_exceeded_wait = min($project_exceeded_wait, (BIGQUERY_WAIT_DEFAULT + $retry_count) * MAX_RUN_PROCS);

						$uuid_suffix = str_split($uuid_, 32);
						$uuid_prefix = array_shift($uuid_suffix);
						$uuid_suffix = implode('', $uuid_suffix) . '_';

						while (++$retry_count)
						{
							$timestamp = (time() + $project_exceeded_wait);

							$uuid_ = $uuid_suffix . $retry_count;

							if (strlen($uuid_prefix . $uuid_) > 64)
							{
								$uuid_ = substr(md5($uuid_), 0, 32);
							}

							$uuid_ = substr($uuid_prefix . $uuid_, 0, 64);

							$target_table_ = sprintf(
								"%s.%s.%s",
								$bq_config['projectId'],
								$bq_config['datasetId'],
								'cloud_to_bq_exceeded');

							$bq_status_query = sprintf(
								'INSERT INTO `' . $target_table_ . '` (`uuid`, `key`, `target`, `timestamp`) '
								. "SELECT * FROM (SELECT `uuid`, '%s' `key`, '%s' `target`, '%s' `timestamp` FROM `" . $target_table_ . "` WHERE `key` = '%s' AND `timestamp` < %s LIMIT 1) `t` "
								. "UNION ALL (SELECT '%s' `uuid`, '%s' `key`, '%s' `target`, '%s' `timestamp`) LIMIT 1 "
								. "ON DUPLICATE KEY UPDATE `uuid` = '%s', `target` = '%s', `timestamp` = %s",
								JOB_EXCEEDED_KEY,
								$table_id_prefix,
								$timestamp,
								JOB_EXCEEDED_KEY,
								$search_time,
								$uuid_,
								JOB_EXCEEDED_KEY,
								$table_id_prefix,
								$timestamp,
								$uuid_,
								$table_id_prefix,
								$timestamp
							);

							$bq_status_job = $bq_client->query($bq_status_query);
							$bq_status_job_res = $bq_client->runQuery($bq_status_job);

							logWarning('warning : retry insert ... wait ' . BIGQUERY_WAIT_DEFAULT . ' sec');
							sleep(BIGQUERY_WAIT_DEFAULT);
						}

						$project_exceeded_wait += $all_exceeded_num;

						logWarning('warning : limit exceeded ... wait ' . $project_exceeded_wait . ' sec');
						logDefault('waiting_files=' . $all_exceeded_num . '(' . $table_exceeded_num . ')');

						sleep($project_exceeded_wait);
						continue;

					} elseif ($job_error_reason == 'rateLimitExceeded')
					{
						$project_exceeded_wait = mt_rand(1000000, $all_exceeded_num * 1000000 * $retry_count);
						$project_exceeded_wait = min($project_exceeded_wait, (BIGQUERY_WAIT_DEFAULT + $retry_count) * MAX_RUN_PROCS);

						logWarning('warning : limit exceeded ... wait ' . sprintf('%.2f', ($project_exceeded_wait / 1000000)) . ' sec');
						logDefault("waiting_files={$all_exceeded_num}({$table_exceeded_num})");

						// レートリミットに引っかかったので、遅延させる
						usleep($project_exceeded_wait);

						continue;

					} elseif ($job_error_reason == 'backendError'
						|| $job_error_reason == 'internalError'
						|| $job_error_reason == 'badGateway')
					{
						$project_exceeded_wait = min($all_exceeded_num, BIGQUERY_WAIT_DEFAULT) + $retry_count;

						logWarning("warning : temporary failure response ... wait {$project_exceeded_wait} sec");
						logDefault("waiting_files={$all_exceeded_num}({$table_exceeded_num})");

						// レートリミットに引っかかったので、遅延させる
						usleep($project_exceeded_wait);

						continue;
					} elseif($job->info()['status']['errorResult']['reason'] == 'invalid')
					{
						# CSVファイルのデータに問題がある場合は、処理済みにして次のログに進む
						$this->error('warning: invalid');
						$this->error('error : ' . var_export($job->info()));

						// 処理したCSVファイルの終了処理をおこなう
						logfileFinishProcess($env, $bq_config, $status_manage_key, $bq_client, $log_data->name());
					} else {
						logError('error : ' . var_export($job->info(), true));

						// sleep(BIGQUERY_TIMEOUT);
						continue;
					}
				} else
				{
					logDefault('send_rows=' . $job->info()['statistics']['load']['outputRows']);

					// 処理したCSVファイルの終了処理をおこなう
					logfileFinishProcess($env, $bq_config, $status_manage_key, $bq_client, $log_data->name());
				}
			} catch (\Exception $e)
			{
				if ($e->getMessage() == 'Service Unavailable'
					|| $e->getMessage() == 'Bad Gateway'
					|| $e->getMessage() == 'Internal Server Error'
					|| $e->getMessage() == 'Gateway Timeout'
					|| $e->getCode() == 504
					|| $e->getCode() == 503
					|| $e->getCode() == 502
					|| $e->getCode() == 500)
				{
					$project_exceeded_wait = BIGQUERY_WAIT_DEFAULT + $retry_count;

					logWarning("warning : http request failure ... wait {$project_exceeded_wait} sec");

					sleep($project_exceeded_wait);
					continue;

				}
				elseif ($e->getMessage() == 'Forbidden'
					|| $e->getCode() == 403)
				{
					$project_exceeded_wait = (BIGQUERY_WAIT_DEFAULT + $retry_count) * MAX_RUN_PROCS;

					logWarning("warning : bigquery resource limit ... wait {$project_exceeded_wait} sec");

					sleep($project_exceeded_wait);
					continue;
				}

				logError("error exception : " . var_export($e, true));
				throw $e;
			}
		}
	}

	logDefault('function complete.');
}

function logfileFinishProcess($env, $bq_config, $status_manage_key, $bq_client, $csv_log_filename): void
{
	$current_time = new DateTime();
	$target_table_ = "{$bq_config['projectId']}.{$bq_config['datasetId']}.cloud_to_bq_sb_logs";

	$bq_status_query = "UPDATE `$target_table_` " .
		"SET `status` = '" . STATUS_FINISHED . "', `log_time` = '" . $current_time->format('Y-m-d H:i:s') . "' " .
		"WHERE `log_file` = '$status_manage_key'";

	$bq_status_job     = $bq_client->query($bq_status_query);
	$bq_status_job_res = $bq_client->runQuery($bq_status_job);

	{
		// すべてのデータセットに読み込まれたかチェック
		// 指定ログファイルの読み込み完了件数を取得

		$bq_status_query = "";
		foreach ($env['ARR_SB_LOG_BQ_GCP_DATASET_ID'] as $bq_tb)
		{
			$target_table_ = "{$bq_config['projectId']}.{$bq_tb}.cloud_to_bq_sb_logs";

			if (!empty($bq_status_query))
			{
				$bq_status_query .= " UNION ALL";
			}

			$bq_status_query .= sprintf(
				" ( SELECT '%s' AS DATASET_NAME, `status` FROM `%s` WHERE `log_file` = '%s' AND `status` = '%s' LIMIT 1 )",
				$bq_tb,
				$target_table_,
				$status_manage_key,
				STATUS_FINISHED
			);
		}

		$bq_status_job     = $bq_client->query($bq_status_query);
		$bq_status_job_res = $bq_client->runQuery($bq_status_job);

		$finished_count = iterator_count( $bq_status_job_res->rows() );

		// 設定された読み込み先数と完了件数を比較して、すべてのデータセットへの読み込みを終えたかを判断する
		if( count( $env['ARR_SB_LOG_BQ_GCP_DATASET_ID'] ) <= $finished_count )
		{
			// 処理済みのCSVファイルを削除(退避)する
			deleteLogFile($bucket, $env['SB_LOG_CLOUD_STORAGE_PATH'], $env['SB_LOG_CLOUD_STORAGE_LOADED_FILE_PATH'], $csv_log_filename);
		}
	}
}

function deleteLogFile($bucket, $storage_path, $storage_loaded_file_path, $remove_log_file_name): void
{
	$new_name = str_replace(
		$storage_path,
		$storage_loaded_file_path,
		$remove_log_file_name
	);

	try
	{
		// 完了済みバケットにコピーし、削除する
		// ※ 保存が不要であればコピーは行わず削除のみ行う
		$log_data->copy($bucket, ['name' => $new_name]);
		$log_data->delete();

		logDefault("$remove_log_file_name を削除しました");
	}
	catch (\Exception $e)
	{
		if ($e->getCode() == 404)
		{
			// 同時にジョブが走った場合など、削除対象ファイルの取り合いになる可能性もあるので、
			// 削除対象ファイルが無いだけであればエラー出力のみでスルーさせる
			logDefault("$remove_log_file_name が存在しません");
		}
		else
		{
			throw $e;
		}
	}
}

function logDefault($log): void
{
	$msg = '[' . now() . '] [LOG] ' . $log . PHP_EOL;
	fwrite(fopen('php://stderr', 'wb'), $msg);
}

function logWarning($log): void
{
	$msg = '[' . now() . '] [WARN] ' . $log . PHP_EOL;
	fwrite(fopen('php://stderr', 'wb'), $msg);
}

function logError($log): void
{
	$msg = '[' . now() . '] [ERR] ' . $log . PHP_EOL;
	fwrite(fopen('php://stderr', 'wb'), $msg);
}

function now(): string
{
	return (new \DateTime())->format('Y-m-d H:i:s');
}
