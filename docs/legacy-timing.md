# 旧実装の所要時間（測定の基準）

フェーズ0の 8.3（[Issue #28](https://github.com/autonomy-local/funadansu/issues/28)）の記録です。移行の前後を比べるための基準で、Funadansu 側の測定はしていません。

## 測った操作

- 操作：旧 catalog の一括登録（`test/e2e/catalog/04-09.Catalog.bulk.yml` のシナリオ全体）
- シナリオは、DB の準備（表の空化、初期データ、有効化の SQL）と、一括登録の要求の列からなる。要求の数は、シナリオ中の `post:` が 612 件
- 測ったのは、`run.sh` を 1 回流す時間（runn の起動、DB の準備、要求と照合、後片付けを含む壁時計）。要求ごとの時間ではない

## 測定の条件

| 項目 | 内容 |
| --- | --- |
| 日付 | 2026-10-11 |
| 機械 | 4 vCPU、メモリ 16 GB（この作業環境の仮想マシン） |
| 旧実装 | pxr-catalog-service `14a647d`、pxr-operator-service `453e638`（[openapi/README.md](../openapi/README.md) と同じ固定のコミット） |
| ランタイム | Node.js 22.22.0、`npm ci` のあと `tsc` でビルド、`NODE_ENV=development` で `node dist/index.js` |
| DB | PostgreSQL 16.15（ローカル、localhost の TCP）。DB 名 `pxr_pod`、DDL は `db/pxr_operator`、`db/pxr_catalog` |
| 認証 | 旧 operator（3000）と旧 catalog（3001）を同時に起動。catalog は operator のセッションを使う |
| テストツール | runn v1.11.1（`go install`） |
| データ量 | 上の 612 件の要求（シナリオに書かれたもの）。DB は毎回、シナリオの先頭で空にしてから入れる |

**注意：DB の NOT NULL の差（D-005）**
旧 DDL（`db/`）のままでは、シナリオの DB 準備（`activate.sql`）が NOT NULL に当たって止まります（`pq: null value in column "attributes"`）。そのため、この測定は、**主キー以外の NOT NULL を外した、ローカルの検証用 DB**で流しました（[docs/e2e-baseline.md](e2e-baseline.md) の診断版と同じ扱い。DDL は変えていません）。厳密な DDL での所要時間は、この測定では出していません。

## 結果

| 回 | 所要時間（壁時計） | 結果 |
| --- | ---: | --- |
| 1 | 8.16 秒 | 通過（1 / 1） |
| 2 | 7.02 秒 | 通過（1 / 1） |
| 3 | 6.48 秒 | 通過（1 / 1） |

- 3 回の中央値は 7.02 秒。最初の回は、ビルド直後の JIT とキャッシュが温まっていないため、少し長い
- 同じ条件で、厳密な DDL のまま流すと、0 秒台で DB の準備の段階で止まる（測定にならない）

## 測っていないこと

- 一括登録の 1 件あたりの時間、件数を変えたときの伸び
- 旧 operator の一括の操作（一括の登録は catalog の操作のみを測った）
- 旧構成（15 コンテナ、ALB、Nginx、Aurora）での時間。コストは [docs/legacy-cost-estimate.md](legacy-cost-estimate.md) を見る
- Funadansu 側の時間（移行の後に、同じ条件で測る）

## 再現の手順

1. 上の 2 リポジトリを、固定したコミットで取得し、`npm ci` して `tsc` でビルドする
2. `db/` の DDL を、DB `pxr_pod` に流し、ロール（`pxr_operator_user`、`pxr_catalog_user`）に権限を与える
3. 検証用に、主キー以外の NOT NULL を外す（上の注意）
4. `NODE_ENV=development` で、operator（3000）と catalog（3001）を起動する
5. `FUNADANSU_E2E_CATALOG_DB='postgres://pxr_catalog_user:<パスワード>@localhost:5432/pxr_pod?sslmode=disable'` と `FUNADANSU_E2E_CATALOG_URL=http://127.0.0.1:3001` を入れ、`test/e2e/catalog/run.sh test/e2e/catalog/04-09.Catalog.bulk.yml` を 3 回流して、壁時計で測る
