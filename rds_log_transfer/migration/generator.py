import json
import os
import string
import jinja2
import yaml


def conv_pascal(text: str) -> str:
    return (string.capwords(text, sep="_")).replace("_", "")


# jsonの文字列型
jsonColumns: dict[str, list[str]]
with open("json_columns.yaml", mode="r") as jsonColumnsFile:
    jsonColumns = yaml.safe_load(jsonColumnsFile)
    print(jsonColumns)


def main():
    # bqのマイグレーションファイル bigquery_table_manager/migration_file.json
    try:
        jsonFile = open("migration_file.json", mode="r")
    except FileNotFoundError:
        print("migration_file.json not exists.")
        exit(1)

    migrations = json.loads(jsonFile.read())

    # RDS記録対象のログリスト
    targetTables: list[str] = []
    with open("target_logs.txt", mode="r") as targetTableFile:
        for line in targetTableFile:
            targetTables.append(line.strip())

    if len(targetTables) == 0:
        print("target_logs is empty.")
        exit(1)

    for migration in migrations["table_data"]:
        if migration["table_name"] not in targetTables:
            continue

        table_migration = TableMigration(migration)
        table_migration.generate_code()
        print(
            f"case \"{table_migration.table_name}\":\nparseFailedRecords, err = {table_migration.get_package_name()}.InsertBqItems(clientWithContext, &records, w.no)")


class TableMigration:
    def __init__(self, migration):
        self._migration = migration

    @property
    def table_name(self) -> str:
        return self._migration["table_name"]

    def get_package_name(self) -> str:
        tableName = self.table_name
        return tableName.replace("SB_LOG_", "")

    def get_struct_columns(self):
        columns = []
        for schema in self._migration["table_schema"]:
            if schema["name"] == "uuid":
                continue

            field = conv_pascal(schema["name"])
            type = ""
            if schema["type"] == "INTEGER":
                type = "json.Number"
            elif schema["type"] == "STRING":
                if self.table_name in jsonColumns and schema["name"] in jsonColumns[self.table_name]:
                    type = "interface{}"
                elif schema["name"] in ["achievement_id", "season_id", "counter_id"]:
                    type = "json.Number"
                else:
                    type = "string"
            elif schema["type"] == "BOOLEAN":
                type = "bool"
            elif schema["type"] == "FLOAT":
                type = "json.Number"
            elif schema["type"] == "TIMESTAMP":
                type = "string"

            column = {
                "field": field,
                "type": type,
                "key": schema["name"],
            }
            columns.append(column)

        return columns

    def get_bq_saver_columns(self):
        columns = []
        for schema in self._migration["table_schema"]:
            name = schema["name"]
            field = conv_pascal(schema["name"])
            if name in ["user_id", "character_id", "unique_id", "item_id"]:
                value = f"i.Record.{field}"
            elif name == "request_header_location":
                value = "sb_log.GetRequestHeaderLocationStr(i.Record)"
            elif self.table_name in jsonColumns and name in jsonColumns[self.table_name]:
                value = f"sb_log.EncodeJsonValue(i.{field})"
            else:
                value = f"i.{field}"

            column = {
                "name": name,
                "value": value,
            }
            columns.append(column)

        return columns

    def get_json_columns(self) -> list[str]:
        table_name = self.table_name
        if table_name not in jsonColumns:
            return []

        columns = []
        for key in jsonColumns[table_name]:
            columns.append(key)

        return columns

    def get_bool_columns(self) -> list[str]:
        columns = []
        for schema in self._migration["table_schema"]:
            if schema["type"] == "BOOLEAN":
                columns.append(schema["name"])
        return columns

    def generate_code(self):
        package_name = self.get_package_name()

        j2 = jinja2.Environment(
            loader=jinja2.FileSystemLoader("./", encoding='utf-8'),
        )
        template = j2.get_template("template.go.j2")

        output_path = f"../app/bq_table/{self.get_package_name()}"
        os.makedirs(output_path, exist_ok=True)
        with open(f"{output_path}/item.go", mode="w", newline="\n") as o:
            o.write(template.render({
                "package_name": package_name,
                "pascal_name": conv_pascal(package_name),
                "struct_columns": self.get_struct_columns(),
                "bq_saver_columns": self.get_bq_saver_columns(),
                "json_columns": self.get_json_columns(),
                "bool_columns": self.get_bool_columns(),
            }))


if __name__ == "__main__":
    main()
