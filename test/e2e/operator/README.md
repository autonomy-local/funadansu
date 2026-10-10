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
| `10-01.PasswordReset.yml` | `10-01.PasswordReset.spec.ts` | 期限の日付で赤になる（D-004）。期限だけ 2030 にした診断版は GREEN。カタログのスタブ（3001）が必要 |
| `14-01.RequestArrayValidator.yml` | `14-01.RequestArrayValidator.spec.ts`（10 件） | GREEN（DB に触る前に 400 を返すため、DB の初期化は省く） |
| `15-01.IdentifyCode.yml` | `15-01.IdentifyCode.spec.ts`（5 件） | GREEN |
| `16-01.IndSmsVerificate.yml` | `16-01.IndSmsVerificate.spec.ts`（42 件） | 赤（期限の日付。D-004）。期限だけ 2030 にした診断版は全件 GREEN。元のテストが期待する `message` のうち 2 つは message.json に無いキーで、シナリオでは `message` が無いことを確かめる（下の注意） |
| `16-02.IndSmsVerificateVerifiy.yml` | `16-02.IndSmsVerificateVerifiy.spec.ts`（7 件） | 赤（期限の日付。D-004）。期限だけ 2030 にした診断版は全件 GREEN |
| `05-04.OperatorCancelDelete.yml` | `05-04.OperatorCancelDelete.spec.ts`（3 件） | GREEN（公式の期限のまま） |
| `06-04.OperatorLogin.yml` | `06-04.OperatorLogin.spec.ts`（1 件） | GREEN（公式の期限のまま） |

残りの operator のテストは、まだ書き換えていません。jest.mock でリポジトリの失敗を起こす試験（12-01 など）は、E2E には入れず、単体テストに回します（差分台帳 D-006）。13-01 は設定値の上書きなので、別に確かめます。catalog のスタブに依る試験は、スタブを揃えるまで入れていません。`src/tests` の supertest を使う 28 本のうち、10 本を書き換えました。残りは 01-01、01-02、02-01、04-01 から 04-04、05-01、05-02、05-04、06-01 から 06-04、11-01、12-01、13-01、17-01 です。

## 動かす前に

1. 旧 operator のサービスを、ポート 3000 で起動します（`pxr-operator-service` の `dist` を、`NODE_ENV=test` で起動）。
2. 旧 operator が使う catalog のスタブを、ポート 3001 で起動します。`pxr-operator-service/src/tests/catalog/pxr-setting.json` を渡します。
   `node test/e2e/stubs/catalog-stub.js <pxr-setting.json のパス>`
3. テストデータを入れる DB を用意します。旧 operator の `ormconfig.json` と同じ `pxr_pod`、ユーザー `pxr_operator_user` です。スキーマは `db/` の DDL から作ります（NOT NULL の差は D-005）。
4. 接続先の DSN を環境変数で渡します。
   `export FUNADANSU_E2E_OPERATOR_DB='postgres://pxr_operator_user:<パスワード>@localhost:5432/pxr_pod?sslmode=disable'`
5. 実行します（runn v1.11.1）。DSN はシナリオに書かず、`run.sh` が `__FUNADANSU_E2E_OPERATOR_DB__` を環境変数の値に置き換えた一時ファイルで流します（runn は runners の中で環境変数を展開しないため）。
   `test/e2e/operator/run.sh`（引数でシナリオを絞れます）

## 注意

- 元のテストは `127.0.0.1` のホストから呼ばれる前提で、旧 operator の CSRF の検査を外しています。シナリオも `127.0.0.1` に向けます。
- 各リクエストに `Connection: close` を付けています。古い keep-alive の接続を使い回すと、EOF で止まるためです。
- 元のテストは、期限の日付を書き換えずに流します。期限切れで落ちる件は、そのまま赤で残します。
- runn は、最初に失敗したステップで止まります。そのため、赤になったシナリオは最初の 1 件しか表示されません。全件の結果は、期限だけ 2030 にした診断版で確かめています（診断版は、リポジトリには入れていません）。
- 16-01 は、元のテストが `message.PHONE_NUMBER_FIELD_IS_NOT_STRING` と `message.REQUIRED_PHONE_NUMBER` を期待しています。旧リポジトリの `config/message.json` にこの 2 つのキーはなく、元のテストは `message` が未定義であることを期待しています。シナリオも、`message` が無いことを確かめます。
