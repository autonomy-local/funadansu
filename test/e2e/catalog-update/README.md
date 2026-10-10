# 旧 catalog-update の E2E（runn）

旧実装（pxr-catalog-update-service、コミット 085f312）の SuperTest のうち、スタブもモックも使わず、E2E の候補だったものを runn のシナリオに写したものです。

## 書き換えた本数

| シナリオ | 元のテスト | 本数 | 状態 |
| --- | --- | --- | --- |
| `05-05.Join.ParamNg.yml` | `src/tests/05-05.Join.ParamNg.spec.ts` | 18 | GREEN（旧実装に対して） |
| `06-08.JoinRemove.ParamNg.yml` | `src/tests/06-08.JoinRemove.ParamNg.spec.ts` | 18 | GREEN（旧実装に対して） |
| `07-07.JoinApproval.ParamNg.yml` | `src/tests/07-07.JoinApproval.ParamNg.spec.ts` | 5 | GREEN（旧実装に対して） |
| `09-01.GetJoin.yml` | `src/tests/09-01.GetJoin.spec.ts` のうち、スタブを使わない 2 件（試験 19・20） | 2 | GREEN（旧実装に対して） |

シナリオには、元の試験の DB 準備（`beforeAll` の `initialData.sql`、「正常」の `INSERT`）を `db` のステップとして入れています。各シナリオの最初に、表を空にします（試験は DB の状態を持ち越すため）。

スタブに依る 18 件と、catalog の停止に依る 1 件（09-01 の「カタログサービスへの接続に失敗」）は、E2E には入れず、単体テストに回します（差分台帳 D-010）。

## 動かす前に

1. 旧 catalog-update を、ポート 3002 で `TZ=Asia/Tokyo` と `NODE_ENV=test` を付けて起動します（試験の時刻の期待が JST のため）。
2. 旧 catalog を、ポート 3001 で起動します（catalog-update の一部の API が、catalog を呼ぶため）。
3. DB（`pxr_pod`、スキーマ `pxr_catalog_update`）を、`db/pxr_catalog_update/001_tables.sql` から作ります。
4. 接続先を環境変数で渡します。
   `export FUNADANSU_E2E_CATALOG_UPDATE_DB='postgres://pxr_catalog_update_user:<パスワード>@localhost:5432/pxr_pod?sslmode=disable'`
   接続先（cu）は `FUNADANSU_E2E_CATALOG_UPDATE_URL` で切り替えます（既定は `http://127.0.0.1:3002`）。
5. 実行します。`test/e2e/catalog-update/run.sh`（引数でシナリオを絞れます）

## 注意

- runn は標準入力を読むことがあるため、`run.sh` は `< /dev/null` で流します。
- 検査は、元の試験が見ている項目だけにしています。たとえば「承認コードがない」（07-07）は、元の試験どおりステータス 404 だけを見ます。記録した本文（既定の `{}`）は、supertest が本文を持たないときの既定値のため、比べていません。
- 試験の時刻（`applicantDate` など、実行時刻で変わる値）は、元の試験が見ていないため比べていません。
- `run.sh` は、接続先の値を置き換えた一時ファイルを作り、実行後に消します。
