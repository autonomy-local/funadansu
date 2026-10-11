#!/usr/bin/env bash
# 旧基盤の E2E（test/e2e/ の runn シナリオ）を、Funadansu に向けて全部流し、通過数を出します。
#
# 接続先は proxy（既定 127.0.0.1:8787）と、その背後の Go（127.0.0.1:8080）です。DB は FUNADANSU_DATABASE_URL で指定します。
# 既定では Go と proxy を、このスクリプトが一時的に起動し、終わったら止めます（JWT の鍵は乱数で作り、残しません）。
# 旧基盤の E2E は、フェーズ0の判定で「全部 RED」が正です。落ちたシナリオがあっても、このスクリプトは 0 で終わります。
# 終了コードが 0 でないのは、サービスが起動しなかったときだけです。
#
# 前提：nix develop の中で実行し、先に db-start と db-init を済ませておきます。
set -euo pipefail

root="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$root"

export FUNADANSU_DATABASE_URL="${FUNADANSU_DATABASE_URL:-postgres://postgres@localhost:5432/funadansu?sslmode=disable}"
proxy_url="${FUNADANSU_TEST_PROXY_URL:-http://127.0.0.1:8787}"
go_url="${FUNADANSU_TEST_GO_URL:-http://127.0.0.1:8080}"

# 各単位の E2E の接続先は、proxy に向ける（Funadansu の入口）。DB は、テストデータを入れるため同じ DB を指す。
export FUNADANSU_E2E_OPERATOR_URL="$proxy_url" FUNADANSU_E2E_OPERATOR_DB="$FUNADANSU_DATABASE_URL"
export FUNADANSU_E2E_CATALOG_URL="$proxy_url" FUNADANSU_E2E_CATALOG_DB="$FUNADANSU_DATABASE_URL"
export FUNADANSU_E2E_CATALOG_UPDATE_URL="$proxy_url" FUNADANSU_E2E_CATALOG_UPDATE_DB="$FUNADANSU_DATABASE_URL"
export FUNADANSU_E2E_NOTIFICATION_URL="$proxy_url"
export FUNADANSU_E2E_BOOK_MANAGE_URL="$proxy_url" FUNADANSU_E2E_BOOK_MANAGE_DB="$FUNADANSU_DATABASE_URL"
export FUNADANSU_E2E_ACCESS_MANAGE_URL="$proxy_url" FUNADANSU_E2E_ACCESS_CONTROL_URL="$proxy_url" FUNADANSU_E2E_ACCESS_CONTROL_DB="$FUNADANSU_DATABASE_URL"

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

# Go と proxy を、既定の接続先で起動する（既に外で起動しているときは、その URL を使う）。
if [ -z "${FUNADANSU_TEST_GO_URL:-}" ]; then
  (cd bootstrap && go build -o "$workdir/funadansu" ./cmd/funadansu)
  FUNADANSU_ADDR=127.0.0.1:8080 "$workdir/funadansu" >"$workdir/go.log" 2>&1 &
  pids+=("$!")
fi
if [ -z "${FUNADANSU_TEST_PROXY_URL:-}" ]; then
  export FUNADANSU_JWT_KID="e2e-$(node -e 'process.stdout.write(require("node:crypto").randomBytes(4).toString("hex"))')"
  export FUNADANSU_JWT_KEY="$(node -e 'process.stdout.write(require("node:crypto").randomBytes(32).toString("base64url"))')"
  (cd bootstrap/proxy && if [ ! -d node_modules ]; then npm ci --no-audit --no-fund; fi)
  (
    cd bootstrap/proxy
    FUNADANSU_UPSTREAM_URL="$go_url" \
      FUNADANSU_PROXY_ADDR=127.0.0.1:8787 \
      exec bun run src/entry/bun.ts
  ) >"$workdir/proxy.log" 2>&1 &
  pids+=("$!")
fi

# 起動を待つ。応答しなければ、ログを出して止める（E2E の結果は出さない）。
wait_for() {
  local url="$1" name="$2" i
  for i in $(seq 1 60); do
    if [ "$(curl -s -o /dev/null -w '%{http_code}' "$url")" != "000" ]; then
      return 0
    fi
    sleep 0.5
  done
  echo "$name が起動しませんでした。ログ：" >&2
  cat "$workdir/go.log" "$workdir/proxy.log" 2>/dev/null >&2
  return 1
}
if [ -z "${FUNADANSU_TEST_GO_URL:-}" ]; then wait_for "$go_url/healthz" go; fi
if [ -z "${FUNADANSU_TEST_PROXY_URL:-}" ]; then wait_for "$proxy_url/healthz" proxy; fi

# 単位ごとに、run.sh へシナリオを1本ずつ渡す（通過数を数えるため）。
# 対象の並びは、各 run.sh の既定（引数なしのとき）と同じ glob を使う。
declare -A glob=(
  [access-control]='[0-9A-Za-z]*.yml'
  [book-manage]='[0-9]*.yml'
  [catalog]='[0-9]*.yml'
  [catalog-update]='[0-9]*.yml'
  [notification]='[0-9A-Za-z]*.yml'
  [operator]='[0-9][0-9]-*.yml'
)

pass=0
fail=0
summary="${GITHUB_STEP_SUMMARY:-/dev/null}"
{
  echo "## 旧基盤の E2E（Funadansu に向けて）"
  echo
  echo "| 単位 | シナリオ | 結果 |"
  echo "| --- | --- | --- |"
} >>"$summary"

for unit in access-control book-manage catalog catalog-update notification operator; do
  shopt -s nullglob
  for scenario in "test/e2e/$unit"/${glob[$unit]}; do
    name="$(basename "$scenario")"
    if timeout 300 "test/e2e/$unit/run.sh" "$scenario" >"$workdir/$unit-$name.log" 2>&1; then
      pass=$((pass + 1))
      result="通過"
    else
      fail=$((fail + 1))
      result="失敗"
    fi
    echo "$unit/$name: $result"
    echo "| $unit | $name | $result |" >>"$summary"
  done
  shopt -u nullglob
done

total=$((pass + fail))
echo "旧 E2E：通過 $pass / 全体 $total（失敗 $fail）"
{
  echo
  echo "**通過 $pass / 全体 $total**（フェーズ0のベースラインは 0 / 32。落ちても CI は落としません）"
} >>"$summary"
