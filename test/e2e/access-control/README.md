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
