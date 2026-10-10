# 単位：access-control（アクセス認可）

## 概要

| 項目 | 内容 |
| --- | --- |
| 単位 | access-control（アクセス認可） |
| 旧実装 | pxr-access-control-manage-service、pxr-access-control-service |
| OpenAPI | `openapi/legacy/` の access-control-manage.json、access-control.json |
| 旧スキーマ | 未確認（フェーズ1で DDL を見る） |
| 状態 | 作業中（フェーズ0：旧 E2E の候補のうち、access-control-manage の 75 件を確認し E2E に入れた。残りは下の「残っていること」） |
| 関連 Issue | #22（P0-6 / 6.2） |

## この単位が担うこと

未整理です。フェーズ1で、旧 OpenAPI と旧実装を読んでから書きます。

## 依存

未整理です（フェーズ1で決める）。

## 外に出す関数の一覧

まだ書いていません（フェーズ1で決める）。

## 旧試験の分類

試験の一覧は [docs/legacy-unit-test-inventory.md](../legacy-unit-test-inventory.md) の該当の節。

| 項目 | 件数 |
| --- | --- |
| 試験ファイル | 29 |
| 試験（件数） | 321 |
| 単体テストへ（スタブ・モック・jest.spyOn に依る） | 184（D-009） |
| E2E の候補（supertest、依存なし） | 137（確認済み 75 件、未確認 62 件） |
| 対象外（supertest なし） | ファイル 0 件 |

## 旧 E2E の状況

| 時点 | 通過（シナリオ） | 失敗（シナリオ） | 備考 |
| --- | --- | --- | --- |
| フェーズ0（分類後） | 0 | 0 | E2E の候補は未確認。単体テストへ回すものは D-009 |
| フェーズ0（候補の確認後） | 1 | 0 | access-control-manage の CreateAPIKey.validator.spec.ts（75 件）が GREEN（`test/e2e/access-control/CreateAPIKey.validator.yml`） |

## 差分台帳への記録

- D-009（下流のスタブ・モックに依る旧試験。E2E から外す）

## 移植と新規

まだ何も移していません。

## 確認の途中で分かったこと

- access-control（旧）の 00-00.Validation.spec.ts の 200 の 2 件は、DB に既に入っている行に依ります。試験の中には行を入れる処理がなく、前の試験の残りに依るため、空の DB からは GREEN になりません。E2E に入れるには、どの行を入れるかを決めてからにします（要確認）
- access-control（旧）の 03-01、03-02 は、前の応答の `apiToken` を次の要求に使います。runn では値を受け渡すための書き方が要ります（未着手）
- 試験の DB の NOT NULL の差は、notification と book-manage と同じく、ローカルの試験用 DB だけで外しました（D-005 と同じ種類の差）

## 作業の記録

| 日付 | 誰が | やったこと |
| --- | --- | --- |
| 2026-10-10 | Claude | 旧試験を分類し、単体テストへ回すもの（D-009）を差分台帳に記録。引き継ぎファイルを作る |
| 2026-10-10 | Claude | access-control-manage の CreateAPIKey.validator.spec.ts（75 件）を、旧実装を起動して流し、75 件 GREEN を確認。runn にした |

## 残っていること

- [x] access-control-manage の候補 75 件を確かめ、GREEN のものを `test/e2e/access-control/` に入れた
- [ ] access-control（旧）の候補 62 件（00-00、02-01、03-01、03-02、03-08、03-09）を確かめる。00-00 は DB の行に依るため要確認。03-01 と 03-02 は応答の値を受け渡すため、runn の書き方を決める
- [ ] 単体テストへ回した 184 件を、Go 側の単体テストで置き換える（フェーズ1以降）
- [ ] 旧 OpenAPI と旧実装を読み、この単位が担うことと依存を書く（フェーズ1）
- [ ] 外に出す関数の一覧を決める（フェーズ1）
