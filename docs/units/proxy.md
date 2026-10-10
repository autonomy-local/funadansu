# 単位：proxy（入口、認証）

## 概要

| 項目 | 内容 |
| --- | --- |
| 単位 | proxy（入口、認証） |
| 旧実装 | pxr-block-proxy-service（旧 OpenAPI は `openapi/legacy/block-proxy.json`） |
| OpenAPI | `openapi/legacy/block-proxy.json` |
| 旧スキーマ | なし（ProxyLog のテーブルは旧 DB にある） |
| 状態 | 作業中（フェーズ0：旧 E2E の分類。E2E は 0 本） |
| 関連 Issue | #22（P0-6 / 6.2） |

## この単位が担うこと

クライアントからのリクエストを受けて、認証と認可を検査し、ブロック（データ通帳）の単位で下流のサービスへ中継する入口。Funadansu では、JWT の署名検査もここで行う（D-001）。

## 依存

| 方向 | 単位 | 使うもの |
| --- | --- | --- |
| この単位が使う | operator | セッションの確認（旧実装）。Funadansu では JWT の検査に置き換える（D-001） |
| この単位が使う | access-control | 静的認可（`/access-control/token`、`/access-control/collate`） |
| この単位が使う | catalog | 権限の定義（`/catalog`） |
| この単位を使う | 各単位 | 中継の先 |

## 外に出す関数の一覧

まだ書いていません（フェーズ1で決める）。

## エンドポイント

E2E はまだ書いていません（下の「旧 E2E の状況」を参照）。

| メソッド | パス | 旧 E2E | 状態 | 備考 |
| --- | --- | --- | --- | --- |
| GET、POST、PUT、DELETE | `/pxr-block-proxy` | `src/tests/` の 11 ファイル（350 件） | 単体テストへ（D-008） | 下流をスタブに差し替えて中継の結果を見る |
| GET、POST、PUT、DELETE | `/pxr-block-proxy/ind` | 同上 | 単体テストへ（D-008） | 個人用の入口 |
| GET、POST、PUT、DELETE | `/pxr-block-proxy/reverse` | `ReverseProxy.spec.ts` など | 単体テストへ（D-008） | 逆向きの中継。`StubReverseProxyAPI` に依る |

## 旧 E2E の分類

| 試験ファイル | 件数 | 下流のスタブ | 分類 |
| --- | --- | --- | --- |
| `AbnormalProxy.auth.spec.ts` | 10 | `StubCatalogService`、`StubOperatorService` | 単体テスト（D-008） |
| `AbnormalProxy.catalog.spec.ts` | 24 | `StubOperatorService` | 単体テスト（D-008） |
| `AbnormalProxy.reverse.spec.ts` | 8 | `StubCatalogService`、`StubAccessControlService`、`StubOperatorService`、`StubService` | 単体テスト（D-008） |
| `AbnormalProxy.token.spec.ts` | 16 | `StubCatalogService`、`StubOperatorService` | 単体テスト（D-008） |
| `AbnormalReverseProxy.token.spec.ts` | 4 | `StubAccessControlServiceWithRejection`、`StubOperatorService`、`StubCatalogService` | 単体テスト（D-008） |
| `BufferRequest.spec.ts` | 2 | 複数（`StubServer.ts`） | 単体テスト（D-008） |
| `Proxy.auth.spec.ts` | 10 | 複数（`StubServer.ts`） | 単体テスト（D-008） |
| `Proxy.spec.ts` | 68 | 複数（`StubServer.ts`） | 単体テスト（D-008） |
| `ReverseProxy.spec.ts` | 82 | 複数。`jest.spyOn` で `doPostRequest` と `doPutRequest` を差し替える箇所あり | 単体テスト（D-008。jest.spyOn の箇所は D-006 と同じ扱い） |
| `UsingMultipleApiKeysOnceRequest.ind.spec.ts` | 62 | 複数（`StubServer.ts`） | 単体テスト（D-008） |
| `UsingMultipleApiKeysOnceRequest.spec.ts` | 64 | 複数（`StubServer.ts`） | 単体テスト（D-008） |

合計 350 件。11 ファイルすべてが `StubServer.ts` のスタブを使います。10 ファイルは、`beforeAll` で `config/port.json`（と `permission.json`）を `src/tests/mocks.*.json` に差し替えています（`Proxy.auth.spec.ts` は差し替えていません）。

## 旧 E2E の状況

| 時点 | 通過（シナリオ） | 失敗（シナリオ） | 備考 |
| --- | --- | --- | --- |
| フェーズ0（着手前） | 0 | | 全部 RED（旧実装のベースラインは未取得） |
| フェーズ0（分類後） | 0 | 0 | シナリオは 0 本。350 件は D-008 で単体テストへ |

## 差分台帳への記録

- D-008（下流のスタブに依る試験。E2E から外す）
- D-001（認証とセッション。operator と共通）

## 移植と新規

まだ何も移していません。

## 作業の記録

| 日付 | 誰が | やったこと |
| --- | --- | --- |
| 2026-10-10 | Claude | 旧 block-proxy の試験 11 ファイル（350 件）を分類。全件が下流のスタブに依るため、D-008 として単体テストへ回した。E2E は 0 本 |

## 残っていること

- [ ] 下流に依らない試験（認証の拒否、入力の検査など）を、旧実装の起動した状態で確かめる。GREEN になったものだけ `test/e2e/proxy/` に runn のシナリオとして書く（desc に元のテストへのリンクを張る）
- [ ] 残りの 350 件は、Go 側の単体テストで置き換える（フェーズ1以降）
- [ ] 外に出す関数の一覧を決める（フェーズ1）
- [ ] Funadansu の proxy（Hono）に向けた E2E の URL の切り替え（`FUNADANSU_E2E_PROXY_URL` の形で、operator と同じ）
