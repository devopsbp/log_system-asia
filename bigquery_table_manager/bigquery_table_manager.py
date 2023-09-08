import json

from google.oauth2 import service_account
from google.cloud import bigquery
from datetime import datetime
from dateutil.relativedelta import relativedelta

class BigQueryTableManager:

    def __init__(self, dataset_ids, migration_file_paths, credentials_path):
        self.credentials = self.create_credentials(credentials_path)
        self._bq_client = bigquery.Client(
            credentials=self.credentials,
            project=self.credentials.project_id
        )
        self.dataset_ids = dataset_ids

        ref_list = [ self._bq_client.dataset(dataset_id) for dataset_id in self.dataset_ids ]
        self.dataset_ref_map = dict(zip(self.dataset_ids,ref_list))

        self.migration_file_paths = migration_file_paths
        self.migrations_data = None

    def load_migration_files(self):
        json_data_list = []
        for json_file_path in self.migration_file_paths:
            with open(json_file_path) as json_file:
                json_data_list.append(json.load(json_file))

        return json_data_list

    def get_tables(self, dataset_id):
        return list(self._bq_client.list_tables(dataset_id))

    def create_table(self, table_ref, table_schema):
        schema = self.create_schema(table_schema)
        table = bigquery.Table(table_ref, schema=schema)

        # パーティション方法を指定する
        table.time_partitioning = bigquery.TimePartitioning(
            type_ = bigquery.TimePartitioningType.DAY,
            # 負荷試験ではダミーログのlog_timeが固定であるため、fluent_timeを使用する
            field = "log_time",
            #field = "fluent_time",
        )

        job = self._bq_client.create_table(table, exists_ok=True)

        print("Created table {}.{}.{}".format(job.project, job.dataset_id, job.table_id))

    def update_table(self, table_ref, new_schema):
        table = self._bq_client.get_table(table_ref)
        schema = self.update_schema(table.schema, new_schema)
        if table.schema != schema:
            table.schema = schema
            table = self._bq_client.update_table(table, ["schema"])

            # パーティション方法を指定する
            table.time_partitioning = bigquery.TimePartitioning(
                type_ = bigquery.TimePartitioningType.DAY,
                # 負荷試験ではダミーログのlog_timeが固定であるため、fluent_timeを使用する
                #field = "log_time",
                field = "fluent_time",
            )

            print("Updated table {}".format(table_ref))
        else:
            print("No changes found on table {}".format(table_ref))

    def delete_table(self, table_ref):
        table = self._bq_client.get_table(table_ref)
        self._bq_client.delete_table(table)

    @staticmethod
    def create_credentials(key_path):
        return service_account.Credentials.from_service_account_file(
            key_path,
            scopes=["https://www.googleapis.com/auth/cloud-platform"],
        )

    def create_schema(self, table_schema):
        return [self.create_schema_field(schema_field["name"], schema_field["type"]) for schema_field in table_schema]

    @staticmethod
    def create_schema_field(field_name, field_type):
        return bigquery.SchemaField(field_name, field_type, mode="NULLABLE")

    def update_schema(self, current_schema, new_schema):
        current_schema_fields = [field.name for field in current_schema]
        result_schema = current_schema[:]
        for field in new_schema:
            if field["name"] not in current_schema_fields:
                result_schema.append(self.create_schema_field(field["name"], field["type"]))
        return result_schema

    def migrate(self):
        migration_data_list = self.load_migration_files()

        for dataset_id in self.dataset_ids:
            for migration_data in migration_data_list:
                current_tables = self.get_tables(dataset_id)
                current_table_ids = [table.table_id for table in current_tables]

                dataset_ref = self.dataset_ref_map[dataset_id]

                for migration_table in migration_data["table_data"]:
                    table_name = migration_table["table_name"]
                    table_schema = migration_table["table_schema"]

                    if table_name.startswith("SB_LOG"):
                        # シャーディング対応でテーブルを10年(西暦の先頭3桁)毎に作成するため、現時点と次のテーブルを作成する
                        current_shard_id = str(datetime.utcnow().date().year)[:3]
                        next_shard_id = str(datetime.utcnow().date().year + 10)[:3]
                        table_name_shard_current = table_name+"_"+current_shard_id
                        table_name_shard_next = table_name+"_"+next_shard_id

                        table_ref = dataset_ref.table(table_name_shard_current)
                        if table_name_shard_current not in current_table_ids:
                            self.create_table(table_ref, table_schema)
                        else:
                            self.update_table(table_ref, table_schema)

                        table_ref = dataset_ref.table(table_name_shard_next)
                        if table_name_shard_next not in current_table_ids:
                            self.create_table(table_ref, table_schema)
                        else:
                            self.update_table(table_ref, table_schema)

                    else:
                        table_ref = dataset_ref.table(table_name)
                        if table_name not in current_table_ids:
                            self.create_table(table_ref, table_schema)
                        else:
                            self.update_table(table_ref, table_schema)
