#!/usr/bin/env bash
# 旧 operator の runn シナリオを実行します。
#
# シナリオの db の接続先は、環境変数 FUNADANSU_E2E_OPERATOR_DB に入れた DSN で置き換えます（DSN に秘密を書かないため）。
# 置き換えた一時ファイルは、実行後に消します。
#
#   export FUNADANSU_E2E_OPERATOR_DB='postgres://...'
#   test/e2e/operator/run.sh [シナリオ...]     （省略すると全部）
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
: "${FUNADANSU_E2E_OPERATOR_DB:?FUNADANSU_E2E_OPERATOR_DB を設定してください}"

if [ "$#" -eq 0 ]; then
  set -- "$here"/[0-9][0-9]-*.yml
fi

tmp="$(mktemp -d "$here/.run-XXXXXX")"
trap 'rm -rf "$tmp"' EXIT

status=0
for scenario in "$@"; do
  name="$(basename "$scenario")"
  # DSN は | を含みうるため、区切りには # を使う
  sed "s#__FUNADANSU_E2E_OPERATOR_DB__#$FUNADANSU_E2E_OPERATOR_DB#" "$scenario" > "$tmp/$name"
  runn run "$tmp/$name" || status=1
done
exit "$status"
