#!/usr/bin/env bash
# API のシナリオ（runn）を実行します。Go のサービスと proxy を一時的に起動し、終わったら止めます。
#
# 前提：nix develop の中で実行し、先に db-start で PostgreSQL を起動し、db-init で旧スキーマを作っておきます。
# 接続先は FUNADANSU_DATABASE_URL で変えられます（既定は db-start の localhost）。
# 使う経路：Go は 127.0.0.1:8080、proxy は 127.0.0.1:8787（シナリオの runners と合わせます）。
set -euo pipefail

root="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$root"

export FUNADANSU_DATABASE_URL="${FUNADANSU_DATABASE_URL:-postgres://postgres@localhost:5432/funadansu?sslmode=disable}"
export FUNADANSU_JWT_KID="runn-$(node -e 'process.stdout.write(require("node:crypto").randomBytes(4).toString("hex"))')"
export FUNADANSU_JWT_KEY="$(node -e 'process.stdout.write(require("node:crypto").randomBytes(32).toString("base64url"))')"
export FUNADANSU_TEST_TOKEN="$(node test/api/token.mjs)"

workdir="$(mktemp -d)"
pids=()
cleanup() {
  for pid in "${pids[@]}"; do
    kill "$pid" 2>/dev/null || true
  done
  wait 2>/dev/null || true
  rm -rf "$workdir"
}
trap cleanup EXIT

# Go のサービスを、ビルドしてから起動する（bootstrap/ の README と同じ設定）。
(cd bootstrap && go build -o "$workdir/funadansu" ./cmd/funadansu)
FUNADANSU_ADDR=127.0.0.1:8080 "$workdir/funadansu" >"$workdir/go.log" 2>&1 &
pids+=("$!")

# proxy（Bun）を起動する。依存は package-lock.json で固定する。
(cd bootstrap/proxy && if [ ! -d node_modules ]; then npm ci --no-audit --no-fund; fi)
(
  cd bootstrap/proxy
  FUNADANSU_UPSTREAM_URL=http://127.0.0.1:8080 \
    FUNADANSU_PROXY_ADDR=127.0.0.1:8787 \
    FUNADANSU_JWT_KEY="$FUNADANSU_JWT_KEY" \
    FUNADANSU_JWT_KID="$FUNADANSU_JWT_KID" \
    exec bun run src/entry/bun.ts
) >"$workdir/proxy.log" 2>&1 &
pids+=("$!")

# 両方が HTTP で応答するまで待つ（proxy は /healthz を公開しないため、応答の番号は見ない）。
# 起動に失敗したら、ログを出して止める。
wait_for() {
  local url="$1" name="$2" i
  for i in $(seq 1 60); do
    if [ "$(curl -s -o /dev/null -w '%{http_code}' "$url")" != "000" ]; then
      return 0
    fi
    sleep 0.5
  done
  echo "$name が起動しませんでした。ログ：" >&2
  cat "$workdir/go.log" "$workdir/proxy.log" >&2
  return 1
}
wait_for http://127.0.0.1:8080/healthz go
wait_for http://127.0.0.1:8787/healthz proxy

# シナリオを順に実行する（テストデータを共有するため並列にしない）。
# DB の接続先は、シナリオの既定より FUNADANSU_DATABASE_URL を優先する。
runn run "test/api/bootstrap/*.yml" --concurrent off --runner "db:$FUNADANSU_DATABASE_URL"
