#!/usr/bin/env bash
# 旧 book-manage の runn シナリオを実行します。
#
# 接続先（book-manage）は、環境変数 FUNADANSU_E2E_BOOK_MANAGE_URL で切り替えます（既定は http://127.0.0.1:3005）。
# 置き換えた一時ファイルは、実行後に消します。
#
#   test/e2e/book-manage/run.sh [シナリオ...]     （省略すると全部）
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
book_manage_url="${FUNADANSU_E2E_BOOK_MANAGE_URL:-http://127.0.0.1:3005}"

if [ "$#" -eq 0 ]; then
  set -- "$here"/[0-9]*.yml
fi

tmp="$(mktemp -d "$here/.run-XXXXXX")"
trap 'rm -rf "$tmp"' EXIT

status=0
for scenario in "$@"; do
  name="$(basename "$scenario")"
  # URL は | を含みうるため、区切りには # を使う
  sed -e "s#__FUNADANSU_E2E_BM_URL__#$book_manage_url#" -e "s#__FUNADANSU_E2E_BM_DB__#$FUNADANSU_E2E_BOOK_MANAGE_DB#" "$scenario" > "$tmp/$name"
  runn run "$tmp/$name" < /dev/null || status=1
done
exit "$status"
