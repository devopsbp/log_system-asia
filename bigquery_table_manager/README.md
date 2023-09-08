# BigQuery Table Manager

## Requirements
|ツール|バージョン|備考|
|---|---|---|
|python|3.5||

## Steps before running the script

1. Install Python dependencies.
```
pip install -r requirements.txt
```

2. Copy the service account JSON file.  
The file can be found at `build\web_provisioning\ansible-playbooks\keys\local\sb_bigquery\sb_bq_service_account.json`

3. Set the settings at `config.ini` file.
```
[gcp]
dataset_id = # <-- The BigQuery Dataset name [SB_LOG_dev|SB_LOG_cbt]

[path]
migration_file = # <-- The migration file (in general let it as migration_file.json)
service_account_key = # <-- The path of the service account key obtained in step 2
```

## Directory structure
```
<root>
    ┣━ bigquery_table_manager.py
    ┣━ config.ini
    ┣━ create_or_update_bigquery_tables.py
    ┣━ migration_file.json
    ┣━ README.md ←　本ファイル
    ┗━ requirements.txt
```

## Run

By running the script `create_or_update_bigquery_tables.py` the changes on migration_file will be applied to `dataset_id`.
The `dataset_id` must be under the `project_id` found in the service account key `sb_bq_service_account.json`.

```
python create_or_update_bigquery_tables.py
```

The script will:
1. Create any table present in `migration_file.json` and not present in the `dataset_id`.
2. Add new rows present in `migration_file.json` and not in existent in the table in the `dataset_id`.

The script will not:
1. Delete any table not present in `migration_file.json` and present in the `dataset_id`.
2. Delete any rows not present in `migration_file.json` and existent in the table in the `dataset_id`.
3. Change the `date_type` of any existent table in the `dataset_id` even if it is different in the `migration_file.json`.