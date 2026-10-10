# proxy（Hono の入口の見本）

Funadansu の入口です。**同じコード**（`src/app.ts`）を、Cloudflare Workers と コンテナ（Bun）の両方で動かします。実行環境の違いは `src/entry/` だけに閉じ込めます（[docs/conventions.md](../../docs/conventions.md) の 3 節）。

## 構成

| 場所 | 中身 |
| --- | --- |
| `src/app.ts` | `createApp(deps)`。認証と中継の配線 |
| `src/auth/jwt.ts` | JWT の検査。署名方式は固定（HS256）。トークンの `alg` は信じない（[ADR 0007](../../docs/decisions/0007-セキュリティの設計方針.md)） |
| `src/relay.ts` | Go のサービスへの中継。`Authorization` を落とさず渡す（サービスでも署名を検査するため） |
| `src/config.ts` | 環境変数の読み込み。エラーには値を出さず、キーの名前だけを出す |
| `src/entry/bun.ts` | コンテナ（Bun）の入口 |
| `src/entry/workers.ts` | Cloudflare Workers の入口（`wrangler.toml` の `main`） |

## 経路

| メソッド | パス | 認証 | 行き先 |
| --- | --- | --- | --- |
| 任意 | `/bootstrap/*` | JWT（`Authorization: Bearer`）が必要 | Go の同じパス |
| 任意 | それ以外 | — | 404 |

`/healthz` と `/readyz` は proxy では公開しません（[docs/conventions.md](../../docs/conventions.md) の 11 節）。

## 設定（環境変数）

| 変数 | 意味 |
| --- | --- |
| `FUNADANSU_JWT_KEY` | JWT の検証に使う鍵。**秘密情報。Git に入れない** |
| `FUNADANSU_JWT_KID` | 鍵の識別子。トークンのヘッダーの `kid` と一致しなければ 401 |
| `FUNADANSU_UPSTREAM_URL` | Go のサービスのオリジン（例：`http://127.0.0.1:8080`） |
| `FUNADANSU_PROXY_ADDR` | Bun の待ち受け先。既定は `127.0.0.1:8787`（Bun のみ） |

ローカルの Workers（`wrangler dev`）では、`.dev.vars.example` を `.dev.vars` にコピーして鍵を入れます。`.dev.vars` は `.gitignore` で除外しています。

## 動かす

```sh
npm install                       # 依存（package-lock.json で固定）
npm test                          # 単体テストと中継のテスト（bun test）
npm run typecheck

# Go のサービス（bootstrap/ の README を参照）を 127.0.0.1:8080 で起動してから：
export FUNADANSU_JWT_KEY=...      # ダミーの鍵（例：openssl rand -base64 32）
export FUNADANSU_JWT_KID=dev-dummy-1
export FUNADANSU_UPSTREAM_URL=http://127.0.0.1:8080
bun run dev:bun                   # 127.0.0.1:8787 で待ち受ける

npm run dev:workers               # 127.0.0.1:8788 で待ち受ける（.dev.vars が必要）
```

### コンテナイメージ（Bun）

リポジトリのルートで、Nix でイメージを作ります（実行時の依存は `hono` だけです）。

```sh
nix build .#proxy-image           # result に tar.gz ができる（65532 で動く。待ち受けは 0.0.0.0:8787）
docker load < result
docker run --rm -p 127.0.0.1:8787:8787 \
  -e FUNADANSU_JWT_KEY=... -e FUNADANSU_JWT_KID=... -e FUNADANSU_UPSTREAM_URL=http://... \
  funadansu-proxy:<タグ>
curl -i http://127.0.0.1:8787/bootstrap/operators/1   # 401（認証が要る）
```

proxy は `/healthz` を公開しないため、起動の確認は認証の応答で見ます。

## テストの方針

- `test/jwt.test.ts`：署名、`kid`、`alg`、期限の検査（表駆動）
- `test/app.test.ts`：認証の入口と、実際の HTTP サーバーへの中継（パス、クエリ、本文、`Authorization`、Go に届かないときの 503）
- 鍵は実行のたびに乱数で作ります。テストにも Git にも鍵は残しません。

## 既定の選択と残り

- 署名方式は HS256 のダミーです。本物の鍵と署名方式は、フェーズ1で決めます（鍵とシークレットの配布は Issue の「やらないこと」。フェーズ0の 5.2、5.3 で扱う）。
- 中継のタイムアウトは 10 秒です（`createApp` の `upstreamTimeoutMs`）。
- 旧 API の窓口（operator の移行期間のアダプター）は、この見本にはありません。
