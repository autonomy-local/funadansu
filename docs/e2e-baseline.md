# 旧 E2E のベースライン（旧実装で流した結果）

フェーズ0の 6.4（[Issue #24](https://github.com/autonomy-local/funadansu/issues/24)）の記録です。`test/e2e/` の 32 本のシナリオを、旧実装の各サービス（PxR の公開リポジトリ）に向けて、そのまま流しました。

判定「旧基盤の E2E が、フェーズ0で取ったベースライン以上に通ること」の基準は、下の**厳密**の数字です（Funadansu と同じ DDL の DB で流した結果）。

## 結果

| 版 | 通過 | 失敗 | 内容 |
| --- | --- | --- | --- |
| **厳密（ベースライン）** | **11** | **21** | DDL（`db/`）のまま。Funadansu と同じ DB で流した |
| 診断版 | 24 | 8 | 同じ DDL から、主キー以外の NOT NULL を外した（検証用 DB だけ。DDL は変えていない） |

単位ごとの内訳（通過／全体）

| 単位 | 旧サービス | 厳密 | 診断版 | 厳密で落ちたシナリオ | 診断版で落ちたシナリオ |
| --- | --- | --- | --- | --- | --- |
| operator | pxr-operator-service | 2 / 11 | 6 / 11 | 03-01、05-03、05-04、06-04、07-01、08-01、09-01、16-01、16-02 | 07-01、08-01、09-01、16-01、16-02 |
| access-control | pxr-access-control-service、pxr-access-control-manage-service | 4 / 7 | 6 / 7 | 02-01、03-01、03-02 | 02-01 |
| catalog | pxr-catalog-service | 0 / 7 | 7 / 7 | 04-01、04-02、04-03、04-09、05-01、08-01、10-01 | なし |
| catalog-update | pxr-catalog-update-service | 4 / 4 | 4 / 4 | なし | なし |
| notification | pxr-notification-service | 1 / 1 | 1 / 1 | なし | なし |
| book-manage | pxr-book-manage-service | 0 / 2 | 0 / 2 | 14-01、14-02 | 14-01、14-02 |
| **合計** | | **11 / 32** | **24 / 32** | | |

## 落ちた理由

### 厳密（21 本）

- **18 本：DB の NOT NULL**（`pq: null value in column ... violates not-null constraint`）。旧 E2E の準備（INSERT）が、DDL の NOT NULL に当たる。DDL は ORM のエンティティから起こしたため、旧 DB の実物と違う（[D-005](spec-deviations.md#d-005-db-の列の-not-null旧-db-の実物との差)）。
- **2 本：表が無い**（`relation "pxr_book_manage.region_use" does not exist`）。book-manage の 14-01、14-02。DDL に無い表を、シナリオが参照している（[D-019](spec-deviations.md#d-019-book-manage-の-region_use-表が-ddl-に無い)）。
- **1 本：401**。access-control の 02-01。原因は未特定（後述の「残り」）。

### 診断版（8 本）

- **5 本**（operator の 07-01、08-01、09-01、16-01、16-02）。テストデータの期限が 2025-01-27 に固定されているため、期限切れで落ちるものがある（[D-004](spec-deviations.md#d-004-旧-e2-e-の日付テストデータの期限)）。07-01 と 09-01 は、前提のカタログのスタブ（ポート 3001）を立てていないため、期限と切り分けていない。
- **1 本：401**。access-control の 02-01（厳密と同じ）。
- **2 本：表が無い**。book-manage の 14-01、14-02（厳密と同じ。D-019）。

## 残り（確かめていないこと）

- access-control の 02-01 の 401 の原因。旧 operator のセッションか、続けて流す前提のデータ（access-control-manage が作る API キー）のどちらかと見られるが、切り分けていません。
- operator の 07-01、09-01 が前提とするカタログのスタブは、この作業では立てていません（実物の pxr-catalog-service を 3001 で動かしている）。この 2 本の失敗が期限によるものかは、スタブを立てて確かめる必要があります。
- 元のテストの期限は、書き換えていません（D-004 の方針どおり、赤のまま残す）。

## 旧実装を動かした環境

| 項目 | 内容 |
| --- | --- |
| 旧実装 | 固定したコミットは [openapi/README.md](../openapi/README.md) の表と同じ（`pxr-operator-service` `453e638`、`pxr-catalog-service` `14a647d`、`pxr-catalog-update-service` `085f312`、`pxr-notification-service` `30b7ce8`、`pxr-book-manage-service` `3e2082a`、`pxr-access-control-service` `5cd2d6a`、`pxr-access-control-manage-service` `3aabd9a`） |
| ランタイム | Node.js 22.22.0、npm ci、`tsc` でビルド、`node dist/index.js` |
| DB | PostgreSQL 16（ローカル）。DB 名は旧 ormconfig と同じ `pxr_pod`。DDL は `db/` を流す（`db-init` と同じ内容） |
| テスト | runn v1.11.1。`test/e2e/*/run.sh` を既定の接続先（旧実装のポート 3000、3001、3002、3004、3005、3014、3015）のまま流す |

旧実装の側で、次の 3 点だけ手を入れました（コミットはしていません。ローカルの検証用）。

1. 旧実装は `NODE_ENV=test` のとき listen しません（supertest 用）。E2E で叩くため、`config/test.yml` の内容を `config/development.yml` にコピーして、`NODE_ENV=development` で起動しました。
2. `pxr-access-control-manage-service` は `config/default-raw.crt`（クライアント証明書の既定値）を読みますが、公開リポジトリに無いため、ダミーの文字列を置きました。HTTP の経路でしか使わないため、この試験の結果には影響しません。
3. 旧 ormconfig のロール（`pxr_*_user`）を、ローカルの PostgreSQL に作りました。旧 operator と book-manage の `ddl/testData.sql` は、DDL との差で一部の INSERT が落ちました（NOT NULL と同じ原因）。

## 再現の手順

1. 上の 7 リポジトリを、固定したコミットで取得する。
2. `db/` の DDL を `pxr_pod` に流す。ロールを作り、各 `pxr_*` スキーマに権限を与える。
3. 1. の手入れを行い、`NODE_ENV=development` で 7 サービスを起動する（operator 3000、catalog 3001、catalog-update 3002、notification 3004、book-manage 3005、access-control-manage 3014、access-control 3015）。
4. `FUNADANSU_E2E_*_DB` に、各ロールの DSN を入れて、`test/e2e/{operator,access-control,catalog,catalog-update,notification,book-manage}/run.sh` を流す。

## 測定していないこと

- 8.3（一括登録などの所要時間）は、この Issue の終了条件に入っていないため、測っていません。
