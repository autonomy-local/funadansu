#!/usr/bin/env bash
# 旧 access-control の runn シナリオを実行します。
#
# 接続先は、サービスごとの環境変数で切り替えます（既定はローカルの旧サービス）。
#   acm（access-control-manage）: FUNADANSU_E2E_ACCESS_MANAGE_URL（既定は http://127.0.0.1:3014）
#   ac （access-control）        : FUNADANSU_E2E_ACCESS_CONTROL_URL（既定は http://127.0.0.1:3015）
# 置き換えた一時ファイルは、実行後に消します。
#
#   test/e2e/access-control/run.sh [シナリオ...]     （省略すると全部）
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
acm_url="${FUNADANSU_E2E_ACCESS_MANAGE_URL:-http://127.0.0.1:3014}"
ac_url="${FUNADANSU_E2E_ACCESS_CONTROL_URL:-http://127.0.0.1:3015}"

if [ "$#" -eq 0 ]; then
  set -- "$here"/[0-9A-Za-z]*.yml
fi

tmp="$(mktemp -d "$here/.run-XXXXXX")"
trap 'rm -rf "$tmp"' EXIT

status=0
for scenario in "$@"; do
  name="$(basename "$scenario")"
  # URL は | を含みうるため、区切りには # を使う
  sed -e "s#__FUNADANSU_E2E_ACM_URL__#$acm_url#" -e "s#__FUNADANSU_E2E_AC_URL__#$ac_url#" "$scenario" > "$tmp/$name"
  runn run "$tmp/$name" < /dev/null || status=1
done
exit "$status"
