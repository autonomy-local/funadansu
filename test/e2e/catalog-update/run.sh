#!/usr/bin/env bash
# 旧 catalog-update の runn シナリオを実行します。
#
# 接続先（cu）は FUNADANSU_E2E_CATALOG_UPDATE_URL で、DB の接続先は FUNADANSU_E2E_CATALOG_UPDATE_DB で渡します（DSN に秘密を書かないため）。
# 既定の接続先は http://127.0.0.1:3002 です。置き換えた一時ファイルは、実行後に消します。
# 旧実装は TZ=Asia/Tokyo で起動してください（試験の時刻の期待が JST のため）。
#
#   export FUNADANSU_E2E_CATALOG_UPDATE_DB='postgres://pxr_catalog_update_user:<パスワード>@localhost:5432/pxr_pod?sslmode=disable'
#   test/e2e/catalog-update/run.sh [シナリオ...]     （省略すると全部）
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
: "${FUNADANSU_E2E_CATALOG_UPDATE_DB:?FUNADANSU_E2E_CATALOG_UPDATE_DB を設定してください}"
catalog_update_url="${FUNADANSU_E2E_CATALOG_UPDATE_URL:-http://127.0.0.1:3002}"

if [ "$#" -eq 0 ]; then
  set -- "$here"/[0-9]*.yml
fi

tmp="$(mktemp -d "$here/.run-XXXXXX")"
trap 'rm -rf "$tmp"' EXIT

status=0
for scenario in "$@"; do
  name="$(basename "$scenario")"
  # DSN と URL は | を含みうるため、区切りには # を使う
  sed -e "s#__FUNADANSU_E2E_CU_URL__#$catalog_update_url#" -e "s#__FUNADANSU_E2E_CU_DB__#$FUNADANSU_E2E_CATALOG_UPDATE_DB#" "$scenario" > "$tmp/$name"
  runn run "$tmp/$name" < /dev/null || status=1
done
exit "$status"
