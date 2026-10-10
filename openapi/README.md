# OpenAPI 仕様

Funadansu の互換の基準です。[AGENTS.md](../AGENTS.md) のとおり、`openapi/` の仕様が正です。

## 旧実装の仕様（`legacy/`）

`openapi/legacy/` には、旧実装（PxR、TypeScript の各サービス）が配っている OpenAPI 仕様を、そのまま写しています。修正や統合はしていません（MIGRATION.md の 6.1）。

- 元のファイル：各リポジトリの `config/openapi.json`（OpenAPI 3.0.2、JSON）
- 取得の日付：2026-10-10。リポジトリの既定のブランチの最新（浅いクローン。`--depth 1`）から取りました
- 取得元は [Personal-Data-Linkage-Module](https://github.com/Personal-Data-Linkage-Module) の公開リポジトリです。PxR の仕様は MIT License で、表示は [NOTICE](../NOTICE) にあります

| ファイル | 旧サービス（リポジトリ） | 取得したコミット | コミットの日付 | 仕様の題名 | パス数 |
| --- | --- | --- | --- | --- | --- |
| [access-control-manage.json](legacy/access-control-manage.json) | pxr-access-control-manage-service | `3aabd9ad922bdad35133d3955d53b0a545bdc4fa` | 2025-03-04 | Access-Control-Manage Service | 4 |
| [access-control.json](legacy/access-control.json) | pxr-access-control-service | `5cd2d6a70feb78ba898af2a64cd5eb2ade7e1148` | 2025-03-04 | access-control | 3 |
| [binary-manage.json](legacy/binary-manage.json) | pxr-binary-manage-service | `865ae2b732fc14abd00793bec761dd3be4f44a98` | 2025-03-04 | binary-manage | 9 |
| [block-proxy.json](legacy/block-proxy.json) | pxr-block-proxy-service | `16aebfe86fbd8205dc7386a89608e1898434f1c0` | 2025-03-04 | pxr-block-proxy-service | 3 |
| [book-manage.json](legacy/book-manage.json) | pxr-book-manage-service | `3e2082aee9f3f58beb3dd57a1df4e00873d3e8c5` | 2025-03-06 | book-manage | 59 |
| [book-operate.json](legacy/book-operate.json) | pxr-book-operate-service | `9d51020413bb37a76a69181a56e285672902924d` | 2025-03-04 | book-operate | 17 |
| [catalog-update.json](legacy/catalog-update.json) | pxr-catalog-update-service | `085f312209a74f5ddd95a57e15ec95b14857ab71` | 2025-03-12 | catalog-update | 20 |
| [catalog.json](legacy/catalog.json) | pxr-catalog-service | `14a647d0e01247844a86442b136ebde6a0f1cef3` | 2025-03-04 | catalog | 33 |
| [certificate-manage.json](legacy/certificate-manage.json) | pxr-certificate-manage-service | `92cd428c8ca2e72f5cc60782c66db34a5b38f75d` | 2025-03-04 | certificate-manage | 3 |
| [certification-authority.json](legacy/certification-authority.json) | pxr-certification-authority-service | `6d3d5874edad0d00d78889293a3f2bac69806cee` | 2025-03-04 | certification-authority | 10 |
| [ctoken-ledger.json](legacy/ctoken-ledger.json) | pxr-ctoken-ledger-service | `643c9937303900bd429d7fa27956ce882353080c` | 2025-03-04 | CToken台帳サービス | 4 |
| [identity-verificate.json](legacy/identity-verificate.json) | pxr-identity-verificate-service | `56281650e38091951ed1cae7b1e3e195dc2ca641` | 2025-03-26 | identity-verificate-service | 10 |
| [local-ctoken.json](legacy/local-ctoken.json) | pxr-local-ctoken-service | `1083d74320b66e826400654ffee6f6854906f31e` | 2025-03-04 | local-ctoken | 3 |
| [notification.json](legacy/notification.json) | pxr-notification-service | `30b7ce86b1f3279a5146d2ce91936825e20a002b` | 2025-03-04 | notification-service | 2 |
| [operator.json](legacy/operator.json) | pxr-operator-service | `453e6385a7545be45d6dccc284004457cb8d28f6` | 2025-03-26 | operator | 21 |

### 対象の外

次の4つは、組織の公開リポジトリにないため、仕様はありません（MIGRATION.md の「対象外」）。

- info-account-manage、audit、outsideStoreService、share-trigger

### 元の仕様と、ずれが見つかったところ

写しの時点で見つかった要確認の項目は、[差分台帳](../docs/spec-deviations.md) に書きます。
