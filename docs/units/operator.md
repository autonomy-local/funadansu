# 単位：operator（認証・権限）

## 概要

| 項目 | 内容 |
| --- | --- |
| 単位 | operator（認証・権限） |
| 旧実装 | pxr-operator-service（旧 OpenAPI は `openapi/legacy/operator.json`） |
| OpenAPI | `openapi/legacy/operator.json` |
| 旧スキーマ | `pxr_operator`（DDL は `db/pxr_operator/`） |
| 状態 | 作業中（フェーズ0：旧 E2E の書き換え） |
| 関連 Issue | #22（P0-6 / 6.2） |

## この単位が担うこと

オペレーター（運営、個人、ワークフロー、アプリケーション）の管理、ログイン、セッション、パスワードの変更、SMS 検証コード、利用者情報の登録と取得。

## 依存

| 方向 | 単位 | 使うもの |
| --- | --- | --- |
| この単位が使う | catalog | 権限と設定の catalog（未決。フェーズ1で決める） |
| この単位を使う | proxy | ログインとセッションの検査（D-001） |

## 外に出す関数の一覧

まだ書いていません（フェーズ1で決める）。

## エンドポイント

E2E は `test/e2e/operator/`（[README](../../test/e2e/operator/README.md)）。

| シナリオ | 旧 E2E の元 | 状態 |
| --- | --- | --- |
| 03-01 OperatorGetById | 03-01.OperatorGetById.spec.ts | GREEN |
| 05-03 OperatorDelete.RequestFromAnotherUser | 05-03.OperatorDelete.RequestFromAnotherUser.spec.ts | GREEN |
| 05-04 OperatorCancelDelete | 05-04.OperatorCancelDelete.spec.ts | GREEN |
| 06-04 OperatorLogin | 06-04.OperatorLogin.spec.ts | GREEN |
| 14-01 RequestArrayValidator | 14-01.RequestArrayValidator.spec.ts | GREEN |
| 15-01 IdentifyCode | 15-01.IdentifyCode.spec.ts | GREEN |
| 07-01 OperatorOneTimeLogin | 07-01.OperatorOneTimeLogin.spec.ts | 期限の日付で赤（D-004） |
| 08-01 OperatorLogout | 08-01.OperatorLogout.spec.ts | 期限の日付で赤（D-004） |
| 09-01 OperatorSession | 09-01.OperatorSession.spec.ts | 期限の日付で赤（D-004） |
| 16-01 IndSmsVerificate | 16-01.IndSmsVerificate.spec.ts | 期限の日付で赤（D-004） |
| 16-02 IndSmsVerificateVerifiy | 16-02.IndSmsVerificateVerifiy.spec.ts | 期限の日付で赤（D-004） |

## 旧 E2E の状況

| 時点 | 通過（シナリオ） | 失敗（シナリオ） | 備考 |
| --- | --- | --- | --- |
| フェーズ0（着手前） | 0 | | 全部 RED |
| フェーズ0（書き換え後） | 6 | 5 | 赤は、期限の日付（D-004）だけ。期限を 2030 にした診断版では全件通る |

## 単体テストへ回したもの

| 項目 | 内容 |
| --- | --- |
| jest.mock で失敗を起こす試験 | D-006（12-01 など） |
| catalog のスタブに依る試験 | D-007。スタブの仕様は [docs/legacy-catalog-stub.md](../legacy-catalog-stub.md) |

## 差分台帳への記録

- D-001（認証とセッション）
- D-003（利用者の情報の検索 `pxrId`）
- D-004（旧 E2E の日付）
- D-005（DB の列の NOT NULL）
- D-006（jest.mock に依る試験）
- D-007（catalog のスタブに依る試験）

## 移植と新規

まだありません。

## 作業の記録

| 日付 | 誰が | やったこと |
| --- | --- | --- |
| 2026-10-10 | Claude | 旧 E2E の書き換え（6 本 GREEN、5 本は期限で赤）。jest.mock とスタブに依る試験を単体テストへ回す方針を決めた |

## 残っていること

- [ ] 旧 operator のスタブ依存の試験を単体テストに移す（D-007。フェーズ1以降）
- [ ] jest.mock で失敗を起こす試験を単体テストに移す（D-006。フェーズ1以降）
- [ ] 期限の日付で赤になる 5 本（D-004）の扱いを、フェーズ1の Issue で決める
- [ ] 外に出す関数の一覧を決める
