# 旧 notification の E2E（runn）

旧実装（pxr-notification-service）の SuperTest のうち、スタブもモックも使わず、E2E の候補だったものを runn のシナリオに写したものです。

## 書き換えた本数

| シナリオ | 元のテスト | 状態 |
| --- | --- | --- |
| `Notification.validator.yml` | `src/tests/Notification.validator.spec.ts`（31 件） | GREEN（旧実装に対して） |

元のテストは、旧実装をコミット `30b7ce8` で起動し、スタブなしで流して GREEN を確かめました。シナリオは、その試験の要求と応答を記録して写したもの（自動生成）です。各ステップの `test` は、状態と本文を照合します。

スタブ（`StubServer` など）やモックに依る 55 件は、E2E には入れず、単体テストに回します（差分台帳 D-011）。

## 動かす前に

1. 旧 notification を、ポート 3004 で起動します（`pxr-notification-service` を `NODE_ENV=test` で起動）。
2. テストデータを入れる DB（`pxr_pod`、スキーマ `pxr_notification`）を、`db/pxr_notification/001_tables.sql` から作ります。
3. 実行します（runn v1.11.1）。
   `test/e2e/notification/run.sh`
   接続先を変えるときは `FUNADANSU_E2E_NOTIFICATION_URL` を入れます。

## 注意

- runn は標準入力を読むことがあるため、`run.sh` は `< /dev/null` で流します。
- 元の試験は、`Application.start()` で旧実装をプロセス内で起動します。E2E は、プロセスとして起動した旧実装に HTTP で当てます。
