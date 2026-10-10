# 旧 operator の catalog スタブの仕様

旧 operator の試験は、catalog の API（ポート 3001、`catalog_url`）を、試験の中で立てたスタブに差し替えて動かしていました。ここでは、そのスタブの **仕様（エンドポイントと応答の形）だけ** を抜き出します。スタブのコードは持ち込みません。

この仕様は、単体テストで catalog の応答を作るときの根拠に使います。E2E では使いません（差分台帳 D-007）。

## 対象

旧 operator（`pxr-operator-service`）の `src/tests/StubServer.ts` と、試験ファイルの中にある `_StubCatalogServer`。

## エンドポイント

| パス | 応答 |
| --- | --- |
| `GET /catalog?ns=<名前空間>` | 名前空間ごとに、下の表の JSON を返す |
| `GET /catalog/:code` | カタログ項目（`{ catalogItem: {...} }`）を返す。`code` は利用者情報の項目のコード（30019、30020、30021、30036 など） |
| `GET /catalog/:code/:ver` | 上と同じ |

## `GET /catalog?ns=` の名前空間

| ns | 応答 |
| --- | --- |
| `catalog/ext/test-org/setting/global` | `pxr-setting.json`。`CatalogServer(1)` のときは `pxr-setting_IdServiceOn.json` |
| `catalog/model/auth/member` | `member.json` |
| `catalog/model/auth/catalog` | `catalog.json` |
| `catalog/ext/test-org/setting/actor/data-trader/actor_1000020` | `settings.json` |
| `catalog/model/setting/actor/data-trader` | `data-trader.json` |
| `catalog/model/actor/wf/部署/store` | カタログ項目の配列（`catalogItem` を持つ） |
| `catalog/model/actor/wf/部署/workflow` | 同上 |
| `catalog/ext/test-org/actor/wf/actor_1000004/store` | 同上 |
| `catalog/ext/test-org/actor/wf/actor_1000004/workflow` | 同上 |

上にない名前空間は、試験によって応答が違います（試験ファイル内の `_StubCatalogServer` が決めています）。例：`catalog/model/actor/wf/部署/role`、`catalog/model/format`、`catalog/model/actor/app/application`。

## 試験ごとのスタブ

- `CatalogServer`（`StubServer.ts`）：大半の試験が使う。`number` で IdService の設定を切り替える
- `CatalogServer2`（`StubServer.ts`）：利用者情報の試験（13-01）が 1 件使う
- `_StubCatalogServer`：試験ファイルの中で定義されている。01-01 や 06-01 など。応答は試験ごとに違う

## 扱い

スタブを立てて動かす E2E は持ちません。このスタブに依る試験は、単体テストへ移し、上の応答の形だけを単体テストに写します（差分台帳 D-007）。
