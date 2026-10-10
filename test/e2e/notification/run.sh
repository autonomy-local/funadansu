#!/usr/bin/env bash
# 旧 notification の runn シナリオを実行します。
#
# 接続先（notification）は、環境変数 FUNADANSU_E2E_NOTIFICATION_URL で切り替えます（既定は http://127.0.0.1:3004）。
# 置き換えた一時ファイルは、実行後に消します。
#
#   test/e2e/notification/run.sh [シナリオ...]     （省略すると全部）
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
notification_url="${FUNADANSU_E2E_NOTIFICATION_URL:-http://127.0.0.1:3004}"

if [ "$#" -eq 0 ]; then
  set -- "$here"/[0-9A-Za-z]*.yml
fi

tmp="$(mktemp -d "$here/.run-XXXXXX")"
trap 'rm -rf "$tmp"' EXIT

status=0
for scenario in "$@"; do
  name="$(basename "$scenario")"
  # URL は | を含みうるため、区切りには # を使う
  sed -e "s#__FUNADANSU_E2E_NOTIFICATION_URL__#$notification_url#" "$scenario" > "$tmp/$name"
  runn run "$tmp/$name" < /dev/null || status=1
done
exit "$status"
