# 旧スキーマの DDL

旧 PostgreSQL のスキーマ（15スキーマ、108 表）を、空の表として作るための DDL です。`db-init`（`nix develop` の中）が `db/*/*.sql` を順に流します。

旧実装（PxR）は TypeORM の `synchronize: false` で、テーブルの DDL をリポジトリに持っていません。そのため、各サービスの TypeORM エンティティから起こしました（MIT License、Copyright 2022 NEC Corporation。表示は [NOTICE](../NOTICE) を参照）。

| スキーマ | 旧サービス | ファイル | 表の数 |
| --- | --- | --- | --- |
| pxr_access_control | pxr-access-control-service | [pxr_access_control/001_tables.sql](pxr_access_control/001_tables.sql) | 3 |
| pxr_access_manage | pxr-access-control-manage-service | [pxr_access_manage/001_tables.sql](pxr_access_manage/001_tables.sql) | 3 |
| pxr_binary_manage | pxr-binary-manage-service | [pxr_binary_manage/001_tables.sql](pxr_binary_manage/001_tables.sql) | 2 |
| pxr_block_proxy | pxr-block-proxy-service | [pxr_block_proxy/001_log_called_api.sql](pxr_block_proxy/001_log_called_api.sql) | 1 |
| pxr_book_manage | pxr-book-manage-service | [pxr_book_manage/001_tables.sql](pxr_book_manage/001_tables.sql) | 24 |
| pxr_book_operate | pxr-book-operate-service | [pxr_book_operate/001_tables.sql](pxr_book_operate/001_tables.sql) | 20 |
| pxr_catalog | pxr-catalog-service | [pxr_catalog/001_tables.sql](pxr_catalog/001_tables.sql) | 15 |
| pxr_catalog_update | pxr-catalog-update-service | [pxr_catalog_update/001_tables.sql](pxr_catalog_update/001_tables.sql) | 15 |
| pxr_certificate_manage | pxr-certificate-manage-service | [pxr_certificate_manage/001_tables.sql](pxr_certificate_manage/001_tables.sql) | 1 |
| pxr_certification_authority | pxr-certification-authority-service | [pxr_certification_authority/001_tables.sql](pxr_certification_authority/001_tables.sql) | 1 |
| pxr_ctoken_ledger | pxr-ctoken-ledger-service | [pxr_ctoken_ledger/001_tables.sql](pxr_ctoken_ledger/001_tables.sql) | 5 |
| pxr_identify_verify | pxr-identity-verificate-service | [pxr_identify_verify/001_tables.sql](pxr_identify_verify/001_tables.sql) | 2 |
| pxr_local_ctoken | pxr-local-ctoken-service | [pxr_local_ctoken/001_tables.sql](pxr_local_ctoken/001_tables.sql) | 2 |
| pxr_notification | pxr-notification-service | [pxr_notification/001_tables.sql](pxr_notification/001_tables.sql) | 4 |
| pxr_operator | pxr-operator-service | [pxr_operator/001_tables.sql](pxr_operator/001_tables.sql) | 10 |

## 決めていること

- 各ファイルの先頭に、出どころ（旧サービスとエンティティ）と、判断を入れた箇所を注記しています。
- ライセンスのないリポジトリ（pxr-linkage）の内容は、コピーも参照もしていません。
- 旧 DB の実物との差は、ここでは直しません。動いている旧 DB と比べて、フェーズ3の移行の作業で扱います。
- 外部キーは作っていません（各 PR の判断に従う）。
- 旧実装の ormconfig の database は `pxr_pod` ですが、Funadansu では `PGDATABASE`（既定は `funadansu`）の中に作ります。
