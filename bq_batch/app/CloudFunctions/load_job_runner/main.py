import datetime
import gc
import itertools
import json
import os
import textwrap
import time
import re

import sqlalchemy
from google.cloud import pubsub_v1
from google.cloud import storage

import connect_unix

BQ_WORK_TABLE = 'sb_logs_cloud_to_bq_work_directory'

STORAGE_TIME_FORMAT  = '%Y-%m-%dT%H:%M:%S.%fZ'
DATABASE_TIME_FORMAT = '%Y-%m-%d %H:%M:%S'

STATUS_NEW      = 'new'
STATUS_WORKING  = 'working'

# Publishes a message to a Cloud Pub/Sub topic.
def load_job_runner(event, context):
    # シークレットから環境変数を設定
    set_env()

    # CloudSQL 用の環境変数を設定
    os.environ["SB_LOG_INSTANCE_UNIX_SOCKET"] = f"/cloudsql/{os.environ['SB_LOG_DB_INSTANCE_UNIX_SOCKET']}"

    load_start_interval = int(os.environ['SB_LOG_LOAD_START_INTERVAL_SECONDS'])

    print("init storage client.")
    storage_client = storage.Client(
        project=os.environ['SB_LOG_CLOUD_STORAGE_PROJECT_ID'],
    )

    print("get storage objects.")
    blobs = storage_client.list_blobs(
        bucket_or_name=os.environ['SB_LOG_CLOUD_STORAGE_BUCKET_NAME'],
        prefix=os.environ['SB_LOG_CLOUD_STORAGE_PATH']
    )

    reg_search = rf"{os.environ['SB_LOG_CLOUD_STORAGE_PATH'].rstrip('/')}/[^/]*$"

    # 「処理済み」「リトライ」ディレクトリは無視する
    reg_ignore_list = [
        rf"^{os.environ['SB_LOG_CLOUD_STORAGE_LOADED_FILE_PATH']}",
        rf"^{os.environ['SB_LOG_CLOUD_STORAGE_RETRY_FILE_PATH']}"
    ]

    # "logfilelist.json"が存在するディレクトリリスト作成
    logfilelist_dir_list = []
    blob_list_obj = list(blobs)
    for blob in blob_list_obj:
        # "logfilelist.json"なら
        if blob.name.endswith(os.environ['SB_LOG_CLOUD_STORAGE_LOG_FILE_LIST']):
            # ディレクトリ名を追加
            logfilelist_dir_list.append(os.path.dirname(blob.name))

    target_objects: dict = dict()

    blobs_count: int = 0
    for blob in blob_list_obj:
        blobs_count += 1

        if blob.name.endswith('/'):
            continue

        dir_str: str = os.path.dirname(blob.name)

        # 指定ディレクトリ直下のディレクトリか
        if not re.match(reg_search, dir_str):
            continue

        # 無視リストに含まれるか
        ignore: bool = False
        for reg_ignore in reg_ignore_list:
            ignore = re.match(reg_ignore, dir_str)
            if ignore:
                break

        # 無視対象だった場合にはスキップ
        if ignore:
            continue

        # "logfilelist.json"が存在しないディレクトリはまだ転送中
        if dir_str not in logfilelist_dir_list:
            continue

        # "logfilelist.json"は除外
        if blob.name.endswith(os.environ['SB_LOG_CLOUD_STORAGE_LOG_FILE_LIST']):
            continue

        if dir_str not in target_objects:
            target_objects[dir_str] = list()

        target_objects[dir_str].append(
            # tuple type
            (
                blob,
                datetime.datetime.strptime(blob._properties['timeCreated'], STORAGE_TIME_FORMAT)
            )
        )

    print(f"total blobs count: {blobs_count}")

    con = connect_unix.connect_unix_socket()

    current_time = datetime.datetime.utcnow()
    current_str = format(current_time, DATABASE_TIME_FORMAT)

    try:
        query = textwrap.dedent(f"SELECT log_directory FROM {BQ_WORK_TABLE}")

        # print(query)
        stmt = sqlalchemy.text(query)
        works_already = con.execute(stmt)

    except Exception as e:
        err = 'Error: {}'.format(str(e))
        print(err)
        return err

    for item in works_already:
        dir_name = item['log_directory']

        # DBに登録済みのディレクトリは除外する
        if dir_name in target_objects:
            del target_objects[dir_name]

    # CloudStorageにアップロードされてから読み込み開始可能と判断するまでの期間（秒）
    ignore_time = datetime.timedelta(seconds=load_start_interval)

    for key, val in list(target_objects.items()):
        upload_time = max(val, key=(lambda item: item[1]))[1]
        if current_time < (upload_time + ignore_time):
            # 読み込み開始可能でない場合は対象から除外する
            del target_objects[key]

    # 処理する最大数で区切る
    list_split_count: int = int(os.environ['SB_LOG_CLOUD_STORAGE_LIST_SPLIT_COUNT'])
    target_objects = itertools.islice(target_objects, list_split_count)
    target_objects, objs_tmp = itertools.tee(target_objects)

    object_count: int = 0
    for item in objs_tmp:
        object_count += 1
        print(f"[Storage] new work directory '{item}'")

    del objs_tmp

    target_objects = list(target_objects)

    if object_count > 0:
        try:
            send_split_count = 100
            for i in range(0, object_count, send_split_count):
                insert_values = dict()
                for target_dataset in os.environ['ARR_SB_LOG_BQ_GCP_DATASET_ID'].split(','):
                    start: int = i
                    end: int = i+send_split_count
                    for item in target_objects[start: end]:
                        # 「＜ターゲットデータセット＞:＜ログディレクトリパス＞」をユニークキーフォーマットとする
                        uniq_key = f"{target_dataset}:{item}"
                        insert_values[uniq_key] = (f"""
                            ('{uniq_key}', '{item}', '{target_dataset}', '{STATUS_NEW}', '{current_str}', '{current_str}')\
                        """)

                # 重複を除外するため、すでに DB に登録されているユニークキーを取得する
                query = f"""\
                    SELECT
                        uniq_key
                    FROM
                        {BQ_WORK_TABLE}
                    WHERE
                        uniq_key IN ('{''', '''.join(insert_values.keys())}')
                """

                stmt = sqlalchemy.text(textwrap.dedent(query))
                # print(stmt)
                result = con.execute(stmt)

                # 登録済みユニークキーを持つレコードは追加リストから除外する
                for item in result:
                    print(f"delete key {item['uniq_key']}")
                    del(insert_values[item['uniq_key']])

                if insert_values:
                    # 新規処理ディレクトリレコードを追加する
                    query = f"""\
                    INSERT INTO
                        {BQ_WORK_TABLE} ( uniq_key, log_directory, target_dataset, status, log_time, created_at )
                    VALUES
                        {', '.join(insert_values.values())}
                    """

                    stmt = sqlalchemy.text(textwrap.dedent(query))
                    # print(stmt)
                    con.execute(stmt)

                else:
                    print('value not exist.')

                del insert_values

        except Exception as e:
            err = 'Error: {}'.format(str(e))
            # print(err)
            return err

    else:
        print('[CloudStorage] new directory is none.')

    del target_objects

    work_dir_list: dict = dict()

    try:
        query = textwrap.dedent(f"""\
        SELECT 
            target_dataset,
            log_directory 
        FROM
            {BQ_WORK_TABLE}
        WHERE
            log_directory IN (
                SELECT 
                    log_directory
                FROM 
                    {BQ_WORK_TABLE}
                WHERE
                    status = '{STATUS_NEW}'
                GROUP BY
                    log_directory,
                    log_time
                ORDER BY 
                    log_time
                LIMIT
                    {list_split_count}
            )
        ORDER BY
            target_dataset,
            log_directory
        ;
        """)

        stmt = sqlalchemy.text(textwrap.dedent(query))
        # print(stmt)
        result = con.execute(stmt)

        for item in result:
            dataset_name = item['target_dataset']
            dir_name = item['log_directory']
            if dataset_name not in work_dir_list:
                work_dir_list[dataset_name] = list()

            work_dir_list[dataset_name].append(dir_name)

    except Exception as e:
        err = 'Error: {}'.format(str(e))
        print(err)
        return err

    if result.rowcount == 0:
        mrg = '[CloudStorage] work directory is none.'
        print(mrg)
        return mrg

    # 処理するディレクトリを「処理中」に変更
    cond: str = ''
    for work_dir_per_dataset in work_dir_list:
        if len(work_dir_list[work_dir_per_dataset]) == 0:
            continue

        if cond:
            cond += ' OR '

        cond += f"( target_dataset = '{work_dir_per_dataset}' AND log_directory IN ("

        for index, elm in enumerate(work_dir_list[work_dir_per_dataset]):
            if index != 0:
                cond += ', '

            cond += f"'{elm}'"

        cond += ' ) )'

    query = textwrap.dedent(f"""\
    UPDATE
        {BQ_WORK_TABLE}
    SET
        status = '{STATUS_WORKING}',
        log_time = '{current_str}'
    WHERE ( 
        {cond}
    )
    """)

    stmt = sqlalchemy.text(textwrap.dedent(query))
    # print(stmt)

    try:
        con.execute(stmt)
    except Exception as e:
        err = 'Error: {}'.format(str(e))
        print(err)
        return err

    # Instantiates a Pub/Sub client
    publisher = pubsub_v1.PublisherClient()

    # References an existing topic
    topic_path = publisher.topic_path(os.environ['SB_LOG_GCP_PROJECT_ID'], os.environ['SB_LOG_BQ_BATCH_PUB_TOPIC'])

    try:
        publish_count: int = 0

        for work_dir_per_dataset in work_dir_list:
            dir_list = work_dir_list[work_dir_per_dataset]
            dataset_list = [work_dir_per_dataset]

            message_json = json.dumps({
                'data': {
                    'work_dir': dir_list,
                    'target_dataset': dataset_list,
                },
            })
            message_bytes = message_json.encode('utf-8')
            print(f'Load job start: {dir_list} -> {dataset_list}')

            # Publishes a message
            try:
                publish_future = publisher.publish(topic_path, data=message_bytes)
                publish_future.result()
            except Exception as e:
                print(e)

                cond: str = f"target_dataset = '{work_dir_per_dataset}' AND log_directory IN ("
                for work_dir in dir_list:
                    if cond:
                        cond += ', '

                    cond += f"'{work_dir}'"

                query = textwrap.dedent(f"""\
                UPDATE
                    {BQ_WORK_TABLE}
                SET
                    error = '{e}',
                    log_time = '{current_str}'
                WHERE ( {cond} )
                """)

                stmt = sqlalchemy.text(textwrap.dedent(query))
                # print(stmt)

                try:
                    con.execute(stmt)
                except Exception as e:
                    err = 'Error: {}'.format(str(e))
                    print(err)
                    continue

            publish_count += 1
            time.sleep(0.1)

        print(f"start job count: {publish_count}")

    except Exception as pubsub_exception:
        err = 'Error: {}'.format(str(pubsub_exception.with_traceback()))
        print(err)

    del blobs
    del blob_list_obj
    del work_dir_list

    del con

    gc.collect()

def set_env() -> None:
    with open(os.environ['SB_SECRET_ID']) as f:
        for line in f:
            sp = line.split('=', 1)
            if len(sp) > 1:
                os.environ[sp[0]] = sp[1].strip()
