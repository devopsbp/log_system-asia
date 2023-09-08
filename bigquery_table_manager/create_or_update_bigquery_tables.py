import configparser
from bigquery_table_manager import BigQueryTableManager

config = configparser.ConfigParser()
config.read('config.ini')

DATASET_IDS = [
    line.strip() for line in str(config.get('gcp', 'dataset_ids')).split(',')
]
MIGRATION_FILE_PATHS = [
    line.strip() for line in config.get('path', 'migration_files').split(',')
]
CREDENTIALS_PATH = config.get('path', 'service_account_key')

def main():
    try:
        print('Used Dataset: {} '.format(DATASET_IDS))
        print('Migration file: {}'.format(MIGRATION_FILE_PATHS))
        print('Credentials file: {}'.format(CREDENTIALS_PATH))
        bigquery_table_manager = BigQueryTableManager(DATASET_IDS, MIGRATION_FILE_PATHS, CREDENTIALS_PATH)
        bigquery_table_manager.migrate()
    except Exception as e:
        print("Could not make migration")
        print(e)


if __name__ == "__main__":
    main()
