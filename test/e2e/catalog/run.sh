#!/usr/bin/env bash
# 旧 catalog の runn シナリオを実行します。
#
# 接続先（ct）は FUNADANSU_E2E_CATALOG_URL で、DB の接続先は FUNADANSU_E2E_CATALOG_DB で渡します（DSN に秘密を書かないため）。
# 既定の接続先は http://127.0.0.1:3001 です。置き換えた一時ファイルは、実行後に消します。
#
#   export FUNADANSU_E2E_CATALOG_DB='postgres://pxr_catalog_user:<パスワード>@localhost:5432/pxr_pod?sslmode=disable'
#   test/e2e/catalog/run.sh [シナリオ...]     （省略すると全部）
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
: "${FUNADANSU_E2E_CATALOG_DB:?FUNADANSU_E2E_CATALOG_DB を設定してください}"
catalog_url="${FUNADANSU_E2E_CATALOG_URL:-http://127.0.0.1:3001}"

if [ "$#" -eq 0 ]; then
  set -- "$here"/[0-9]*.yml
fi

tmp="$(mktemp -d "$here/.run-XXXXXX")"
trap 'rm -rf "$tmp"' EXIT

status=0
for scenario in "$@"; do
  name="$(basename "$scenario")"
  # DSN と URL は | を含みうるため、区切りには # を使う
  sed -e "s#__FUNADANSU_E2E_CT_URL__#$catalog_url#" -e "s#__FUNADANSU_E2E_CT_DB__#$FUNADANSU_E2E_CATALOG_DB#" "$scenario" > "$tmp/$name"
  runn run "$tmp/$name" < /dev/null || status=1
done
exit "$status"
