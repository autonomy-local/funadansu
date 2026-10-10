# bootstrap

Funadansu の規約（[docs/conventions.md](../docs/conventions.md)）を、動くコードで示した見本です。新しく単位を移すときは、ここの形に合わせます。

## 構成

| 場所 | 中身 |
| --- | --- |
| `cmd/funadansu/main.go` | 配線をする唯一の場所（設定、ログ、DB、単位、HTTP の順に組み立てる） |
| `internal/platform/` | 共通の基盤。設定、ログ、エラーの形、HTTP のミドルウェアとタイムアウト、ヘルスチェック、DB の接続 |
| `internal/operator/` | 見本の単位。`handler.go`（入口）、`service.go`（業務）、`store/`（sqlc の生成物。手で直さない） |
| `sqlc.yaml` | sqlc の設定。SQL の正本は `internal/operator/store/query.sql`、型の元は `../db/pxr_operator` |

## 動かす

`nix develop` の中で実行します（Go は `GOTOOLCHAIN=local` で使います）。

```sh
db-start                 # ローカルの PostgreSQL を起動（localhost:5432）
db-init                  # 旧スキーマを作る（初回だけ。何度流しても安全）

cd bootstrap
go test ./...            # 単体テストとエンドポイントのテスト（db-start が必要）
sqlc generate            # query.sql を変えたら生成し直す

FUNADANSU_DATABASE_URL='postgres://postgres@localhost:5432/funadansu?sslmode=disable' \
  go run ./cmd/funadansu # 127.0.0.1:8080 で待ち受ける
```

確かめる：

```sh
curl http://127.0.0.1:8080/healthz    # 生存の確認（DB に触れない）
curl http://127.0.0.1:8080/readyz     # DB の ping
curl http://127.0.0.1:8080/bootstrap/operators/1
```

`/bootstrap/operators/{id}` は見本のための経路で、旧 API の形ではありません。`pxr_operator.operator` の `id` と `pxr_id` だけを返し、パスワードや個人情報の列は読みません。見つからなければ 404、整数でなければ 400 を返します。

## 設定（環境変数）

| 変数 | 既定 | 意味 |
| --- | --- | --- |
| `FUNADANSU_DATABASE_URL` | （必須） | PostgreSQL の接続先 |
| `FUNADANSU_ADDR` | `127.0.0.1:8080` | 待ち受け先。Cloud Run では `0.0.0.0:$PORT` を設定 |
| `FUNADANSU_DB_MAX_CONNS` | `5` | 接続プールの上限 |
| `FUNADANSU_DB_STATEMENT_TIMEOUT` | `5s` | 接続ごとの `statement_timeout` |
| `FUNADANSU_READ_HEADER_TIMEOUT` ほか | `5s` / `15s` / `30s` / `60s` | `ReadHeader`、`Read`、`Write`、`Idle` のタイムアウト |
| `FUNADANSU_SHUTDOWN_TIMEOUT` | `8s` | 終了の時限 |
| `FUNADANSU_READY_TIMEOUT` | `2s` | `/readyz` の DB の ping の時限 |
| `FUNADANSU_MAX_BODY_BYTES` | `1048576` | 本文の上限 |
| `FUNADANSU_MEMORY_LIMIT_RATIO` | `0.9` | cgroup の上限に掛ける GOMEMLIMIT の比率 |
| `FUNADANSU_MEMORY_LIMIT_BYTES` | （無し） | GOMEMLIMIT の明示の値 |

設定の形が不正なときは、値を出さずにキーの名前だけを標準エラーに出して終了します。

## ログ

標準出力に、1行1件の JSON で出ます。`kind` は `access`、`application`、`system` のどれかです（[docs/conventions.md](../docs/conventions.md) の 9 節）。
