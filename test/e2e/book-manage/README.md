# 旧 book-manage の E2E（runn）

旧実装（pxr-book-manage-service、コミット 3e2082a）の SuperTest のうち、スタブもモックも使わず、E2E の候補だったものを runn のシナリオに写したものです。

## 書き換えた本数

| シナリオ | 元のテスト | 状態 |
| --- | --- | --- |
| `14-01.Identification.yml` | `src/tests/14-01.Identification.spec.ts`（2 件） | GREEN（旧実装に対して） |
| `14-02.Identification.yml` | `src/tests/14-02.Identification.spec.ts`（2 件） | GREEN（旧実装に対して） |

シナリオには、元の試験の DB 準備（`beforeAll` の `initialData.sql`、「正常」の `INSERT`）を `db` のステップとして入れています。

スタブやモックに依る 1920 件は、E2E には入れず、単体テストに回します（差分台帳 D-016）。

## 動かす前に

1. 旧 book-manage を、ポート 3005 で起動します（`pxr-book-manage-service` を `NODE_ENV=test` で起動）。
2. DB（`pxr_pod`、スキーマ `pxr_book_manage`）を、`db/pxr_book_manage/001_tables.sql` から作ります。DDL の差（`region_use` の欠け、NOT NULL）は [単位の引き継ぎ](../../../docs/units/book-manage.md) の「DDL の差」を見てください。
3. 接続先を環境変数で渡します。
   `export FUNADANSU_E2E_BOOK_MANAGE_DB='postgres://pxr_book_manage_user:<パスワード>@localhost:5432/pxr_pod?sslmode=disable'`
   接続先（bm）は `FUNADANSU_E2E_BOOK_MANAGE_URL` で切り替えます（既定は `http://127.0.0.1:3005`）。
4. 実行します。`test/e2e/book-manage/run.sh`

## 注意

- runn は標準入力を読むことがあるため、`run.sh` は `< /dev/null` で流します。
- 各要求の `session` ヘッダは、元の試験のものを写しています（認証はヘッダの値を使うため、オペレーターへの問い合わせは起きません）。
