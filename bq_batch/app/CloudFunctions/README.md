### CloudFunctions ディレクトリについて

このディレクトリでは [GoogleCloudFunctions](https://console.cloud.google.com/functions/list?env=gen2&authuser=1&hl=ja&project=sbtest-258905) にて使用している実行ファイルが管理されています。  
実際に実行されるファイルは CloudFunctions ジョブ上に配置されますが、バージョン管理のためこちらでも管理を行います。

### 各ファイルについて

* **index.php**

    メインの実行関数が記載されているファイルです。  
    CloudFunctions より、 `LoadToBQ` メソッドが実行されます。


* **composer.json**

    `index.php` の実行に必要なパッケージが記載されているファイルです。  
    パッケージの追加が必要な場合は、このファイルに記載することで利用可能になります。

### GCP SecretManager で管理しているシークレット（環境変数）について 

扱いは AWS SecretsManager と同じです。  
SecretManager に管理されているシークレットを CloudFunctions ジョブが参照することで環境変数を利用できるようにしています。

以下変数を管理しています

* **SB_LOG_BQ_GCP_PROJECT_ID**

    BigQuery のプロジェクトIDを指定します。


* **ARR_SB_LOG_BQ_GCP_DATASET_ID**

    ログデータを投入する BigQuery のデータセットを指定します。 （複数のデータセットを指定可）  
    複数データセットを対象にログデータを出力する際には、「カンマ区切り（スペース不要）」で記載してください。
    
        例） SB_LOG_dev6_1_1 と SB_LOG_dev6_1_2 のデータセットに出力
          ARR_SB_LOG_BQ_GCP_DATASET_ID=SB_LOG_dev6_1_1,SB_LOG_dev6_1_2


* **SB_LOG_CLOUD_STORAGE_PROJECT_ID**

    CloudStorage のプロジェクトIDを指定します。


* **SB_LOG_CLOUD_STORAGE_BUCKET_NAME**

    CloudStorage のバケット名を指定します。 


* **SB_LOG_CLOUD_STORAGE_PATH**

    CloudStorage 上のログデータが配置されているパスを指定します。  
    フルパスではなく、 `SB_LOG_CLOUD_STORAGE_BUCKET_NAME` より後のパスを指定します。

 
* **SB_LOG_CLOUD_STORAGE_LOADED_FILE_PATH**

    すべてのデータセットに読み込まれたログデータの移動先パスを指定します。  
    `SB_LOG_CLOUD_STORAGE_PATH` と同様に `SB_LOG_CLOUD_STORAGE_BUCKET_NAME` より後のパスを指定します。

    ※ 主に開発環境用なので、不使用になる可能性があります