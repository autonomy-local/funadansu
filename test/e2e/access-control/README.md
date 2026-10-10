# 旧 access-control の E2E（runn）

旧実装（pxr-access-control-manage-service、pxr-access-control-service）の SuperTest のうち、スタブもモックも使わず、E2E の候補だったものを runn のシナリオに写したものです。

## 書き換えた本数

| シナリオ | 元のテスト | 状態 |
| --- | --- | --- |
| `CreateAPIKey.validator.yml` | access-control-manage の `src/tests/CreateAPIKey.validator.spec.ts`（75 件） | GREEN（旧実装に対して） |

残りの候補（access-control の 00-00、02-01、03-01、03-02、03-08、03-09、計 62 件）は、まだ確かめていません。単位の引き継ぎ（[access-control.md](../../../docs/units/access-control.md)）を見てください。

スタブやモックに依る 184 件は、E2E には入れず、単体テストに回します（差分台帳 D-009）。

## 動かす前に

1. access-control-manage を、ポート 3014 で、access-control を、ポート 3015 で起動します（`NODE_ENV=test`）。
2. DB（`pxr_pod`、スキーマ `pxr_access_manage` と `pxr_access_control`）を、`db/` の DDL から作ります。
3. 実行します。`test/e2e/access-control/run.sh`
   接続先は `FUNADANSU_E2E_ACCESS_MANAGE_URL`（既定は `http://127.0.0.1:3014`）と `FUNADANSU_E2E_ACCESS_CONTROL_URL`（既定は `http://127.0.0.1:3015`）で切り替えます。

## 注意

- runn は標準入力を読むことがあるため、`run.sh` は `< /dev/null` で流します。
- access-control-manage の `config/default-raw.crt` は公開のリポジトリにありません。テスト環境では証明書の検査を飛ばすため、試験の環境では中身を問わず置くだけで動きます（リポジトリには入れていません）。

## 旧 access-control（pxr-access-control-service）の試験

| シナリオ | 元のテスト | 本数 | 状態 |
| --- | --- | --- | --- |
| `00-00.Validation.yml` | `00-00.Validation.spec.ts` | 8 | GREEN（旧実装に対して） |
| `02-01.AccessControl.yml` | `02-01.AccessControl.spec.ts` | 30 | GREEN（旧実装に対して） |
| `03-01.Collate.yml` | `03-01.Collate.spec.ts` | 12 | GREEN（旧実装に対して） |
| `03-02.CollateOk.yml` | `03-02.CollateOk.spec.ts` | 10 | GREEN（旧実装に対して） |
| `03-08.Collate.NoMacthError.yml` | `03-08.Collate.NoMacthError.spec.ts` | 1 | GREEN（旧実装に対して） |
| `03-09.Collate.Macth.yml` | `03-09.Collate.Macth.spec.ts` | 1 | GREEN（旧実装に対して） |

- 各シナリオは、先頭で元の試験の DB 準備（`initialData.sql`）を流し、試験の順番どおりに続けて流します（前のデータを前提にする試験があるため）
- 応答の `apiToken` は実行ごとに変わるため、後のステップでは `steps[N].res.body[0].apiToken` で受け渡します（02-01）
- 検査は、元の試験が見ている項目だけにしています（応答の本文を見ていない試験は、ステータスだけを見る）
- 旧 access-control は、DB の接続先を `FUNADANSU_E2E_ACCESS_CONTROL_DB`（`pxr_access_control_user` のロールの DSN）で渡します。接続先の URL は `FUNADANSU_E2E_ACCESS_CONTROL_URL`（既定は `http://127.0.0.1:3015`）
- スタブに依る 184 件は E2E には入れず、単体テストに回します（D-009）
