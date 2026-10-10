# 旧 catalog の E2E（runn）

旧実装（pxr-catalog-service、コミット 14a647d）の SuperTest のうち、スタブもモックも使わず、E2E の候補だったものを runn のシナリオに写したものです。

## 書き換えた本数

| シナリオ | 元のテスト | 試験の数 | 状態 |
| --- | --- | --- | --- |
| `04-01.Catalog.model.yml` | `04-01.Catalog.model.spec.ts` | 303 | GREEN（旧実装に対して） |
| `04-02.Catalog.built_in.yml` | `04-02.Catalog.built_in.spec.ts` | 229 | GREEN |
| `04-03.Catalog.ext.yml` | `04-03.Catalog.ext.spec.ts` | 229 | GREEN |
| `04-09.Catalog.bulk.yml` | `04-09.Catalog.bulk.spec.ts` | 337（繰り返しで生成） | GREEN |
| `05-01.CatalogInner.yml` | `05-01.CatalogInner.spec.ts` | 20 | GREEN |
| `08-01.CatalogPublic.yml` | `08-01.CatalogPublic.spec.ts` | 4 | GREEN |
| `10-01.CatalogHistoryCode.yml` | `10-01.CatalogHistoryCode.spec.ts` | 22 | GREEN |

各シナリオは、最初に表を空にしてから（試験は DB の状態を持ち越すため）、元の試験の DB 準備（`initialData.sql` など）と要求を、順番どおりに流します。各要求の応答は、状態と本文を照合します（本文は全体を比べます）。

スタブやモックに依る試験は、E2E には入れず、単体テストに回します（差分台帳 D-010）。

## 動かす前に

1. 旧 catalog を、ポート 3001 で起動します（`pxr-catalog-service` を `NODE_ENV=test` で起動）。
2. DB（`pxr_pod`、スキーマ `pxr_catalog`）を、`db/pxr_catalog/001_tables.sql` から作ります。
3. DSN を環境変数で渡して実行します。
   `export FUNADANSU_E2E_CATALOG_DB='postgres://pxr_catalog_user:<パスワード>@localhost:5432/pxr_pod?sslmode=disable'`
   `test/e2e/catalog/run.sh`（接続先は `FUNADANSU_E2E_CATALOG_URL` で、既定は `http://127.0.0.1:3001`）

## 注意

- 表を空にする文は `TRUNCATE ... CASCADE` です。シーケンスの所有者ではない接続では `RESTART IDENTITY` が使えないため、シーケンスは試験の `SETVAL` に任せています。
- 試験の DB の NOT NULL の差は、ローカルの試験用 DB だけで外しています（D-005 と同じ種類の差。単位の引き継ぎに記録）。
- runn は標準入力を読むことがあるため、`run.sh` は `< /dev/null` で流します。
- 元の試験の中には、DB を読み戻して確かめる箇所があります。その箇所は、シナリオでは DB の読み取りとして写していません（応答の照合だけ）。
