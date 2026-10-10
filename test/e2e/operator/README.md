# 旧 operator の E2E（runn）

旧実装（pxr-operator-service）の SuperTest を、runn のシナリオに書き換えたものです。各シナリオの `desc` に、元のテストへのリンク（コミット固定）を書いています。

## 書き換えた本数

| シナリオ | 元のテスト | 状態 |
| --- | --- | --- |
| `03-01.OperatorGetById.yml` | `03-01.OperatorGetById.spec.ts`（12 件） | GREEN |
| `05-03.OperatorDelete.RequestFromAnotherUser.yml` | `05-03.OperatorDelete.RequestFromAnotherUser.spec.ts`（1 件） | GREEN |
| `07-01.OperatorOneTimeLogin.yml` | `07-01.OperatorOneTimeLogin.spec.ts`（12 件） | 2 件が赤（期限の日付。差分台帳 D-004） |
| `08-01.OperatorLogout.yml` | `08-01.OperatorLogout.spec.ts`（10 件） | 期限の日付で赤になる（D-004） |
| `09-01.OperatorSession.yml` | `09-01.OperatorSession.spec.ts`（10 件） | 期限の日付で赤になる（D-004）。期限だけ 2030 にした診断版は全件 GREEN |

残りの operator のテストは、まだ書き換えていません（10-01、14-01、15-01、16-01、16-02）。

## 動かす前に

1. 旧 operator のサービスを、ポート 3000 で起動します（`pxr-operator-service` の `dist` を、`NODE_ENV=test` で起動）。
2. 旧 operator が使う catalog のスタブを、ポート 3001 で起動します。`pxr-operator-service/src/tests/catalog/pxr-setting.json` を渡します。
   `node test/e2e/stubs/catalog-stub.js <pxr-setting.json のパス>`
3. テストデータを入れる DB を用意します。旧 operator の `ormconfig.json` と同じ `pxr_pod`、ユーザー `pxr_operator_user` です。スキーマは `db/` の DDL から作ります（NOT NULL の差は D-005）。
4. 接続先の DSN を環境変数で渡します。
   `export FUNADANSU_E2E_OPERATOR_DB='postgres://pxr_operator_user:<パスワード>@localhost:5432/pxr_pod?sslmode=disable'`
5. 実行します（runn v1.11.1）。
   `runn run test/e2e/operator/03-01.OperatorGetById.yml`

## 注意

- 元のテストは `127.0.0.1` のホストから呼ばれる前提で、旧 operator の CSRF の検査を外しています。シナリオも `127.0.0.1` に向けます。
- 各リクエストに `Connection: close` を付けています。古い keep-alive の接続を使い回すと、EOF で止まるためです。
- 元のテストは、期限の日付を書き換えずに流します。期限切れで落ちる件は、そのまま赤で残します。
