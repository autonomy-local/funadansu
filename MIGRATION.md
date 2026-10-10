# 移行の状況

旧実装（PxR、TypeScript の15のマイクロサービス）から Funadansu への移行の状況です。
作業を始める人（人・AI エージェント）は、まずここで**いまのフェーズ**と**担当する単位の状態**を確かめてください。

- フェーズの定義：[docs/migration-phases.md](docs/migration-phases.md)
- 単位の一覧と引き継ぎファイル：[docs/units/](docs/units/README.md)
- 差分台帳：[docs/spec-deviations.md](docs/spec-deviations.md)
- セキュリティ：[ADR 0007](docs/decisions/0007-セキュリティの設計方針.md)、[脅威モデル](docs/security/threat-model.md)、[運用の手引き](docs/operations/README.md)

## いまのフェーズ

**フェーズ0**：各単位の移行作業ができる状態にする

## 最終目的

旧基盤の API が提供する機能を満たした、Cloudflare（Hono）と Go のアプリケーションがデプロイできること。
判定：旧基盤の E2E が、フェーズ0で取ったベースライン以上に通ること。

## フェーズの状況

| フェーズ | 目的 | 状態 | リリース |
| --- | --- | --- | --- |
| 0 | 各単位の移行作業ができる状態 | **作業中** | — |
| 1 | 既存の DB（旧スキーマ）のまま、アプリだけ差し替えられる | 未着手 | — |
| 2 | エンドポイントはそのまま、CQRS で分離 | 未着手 | — |
| 3 | DB の統合（統合スキーマ、CockroachDB への対応、読み取り側の複製） | 未着手 | — |
| 4 | 新しい API の設計と移行（旧 API の窓口は残す） | 未着手 | — |

## フェーズ0のタスク

ゴール：**bootstrap が動き**（Hono の proxy〈Cloudflare とコンテナ〉と、コンテナの Go がつながり、bootstrap のエンドポイントのテストが通る）、**旧基盤の E2E が Funadansu に向けて動いて全部 RED** になっていること。

| ID | タスク | 状態 |
| --- | --- | --- |
| 1.1 | リポジトリの作成、`main` の保護 | ✅ 完了 |
| 1.2 | README、LICENSE、ADDITIONAL-PERMISSIONS、NOTICE、ADOPTERS、SECURITY.md | ✅ 完了 |
| 1.3 | CONTRIBUTING、AGENTS.md・CLAUDE.md、MIGRATION.md、`docs/units/`、差分台帳、`docs/decisions/` | ✅ 完了 |
| 1.4 | 旧構想（autonomy-tellus）からのリンクとアーカイブ | ✅ 完了 |
| 2.1 | Nix の開発シェル（`nix develop`） | ✅ 完了 |
| 2.2 | ローカルの PostgreSQL と旧スキーマの投入 | ✅ 完了（`db-start`、`db-init`。[db/README.md](db/README.md)） |
| 3.1 | 規約の決定：ディレクトリとパッケージの構成、命名、DI を使わない配線（`docs/conventions.md` の本文 1〜3 節） | ✅ 完了 |
| 3.2 | 規約の決定：ID・エラー・ログの形（`docs/conventions.md` の本文 7〜9 節） | ✅ 完了 |
| 3.3 | 規約の決定：テストの方針、「外に出す関数の一覧」の書式（`docs/conventions.md` の本文 5・6 節） | ✅ 完了 |
| 3.4 | セキュリティの設計方針（ADR 0007）と脅威モデル（`docs/security/threat-model.md`） | ✅ 完了 |
| 4.1 | bootstrap：Go のサービスの骨組み（`bootstrap/`。platform と見本のエンドポイント1本） | ✅ 完了 |
| 4.2 | Hono の proxy の骨組み（Workers と Bun の両方で動く、Go への中継、JWT の検査の見本） | ✅ 完了 |
| 4.3 | bootstrap：Go と proxy（Bun）のコンテナイメージを Nix でビルド | ✅ 完了 |
| 4.4 | bootstrap：エンドポイントのテスト（接続先を URL で切り替え。`test/api/run.sh`、[bootstrap/README.md](bootstrap/README.md)） | ✅ 完了（ローカルで GREEN） |
| 4.5 | bootstrap：疎通（ローカル、リモート） | 未着手 |
| 5.1〜5.3 | IaC（Pulumi、sops、dev のスタック） | 未着手 |
| 5.4 | 最小のデプロイの手順（ローカルの CI → `pulumi up`・`nixos-rebuild`）と運用の手引きの初版（`docs/operations/`） | 未着手 |
| 6.1 | 旧 OpenAPI を `openapi/legacy/` に集める（取得元のコミットを [openapi/README.md](openapi/README.md) に記録） | ✅ 完了 |
| 6.2〜6.4 | 旧基盤の E2E（URL の切り替え、Funadansu に向けて全部 RED、旧実装のベースライン） | 作業中（operator と proxy。残りの単位は試験の分類のみ。[docs/units/](docs/units/README.md)） |
| 7.1 | CI | 未着手 |
| 7.2 | 脆弱性と秘密情報の検査（govulncheck、CodeQL・gosec、秘密情報の検出、Actions の SHA 固定と最小権限、Renovate） | 未着手 |
| 7.3 | （フェーズ1までに）SBOM、署名、OpenSSF Scorecard | 未着手 |
| 8.1〜8.3 | 測定の基準（規模、費用、所要時間） | 未着手 |

## 単位の状況

| 順 | 単位 | フェーズ1 | フェーズ2 | フェーズ3 | 旧 E2E（通過／全体） | 引き継ぎファイル |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | proxy | 未着手 | — | — | 0 / 0（分類のみ。D-008） | [proxy.md](docs/units/proxy.md) |
| 2 | operator | 未着手 | 未着手 | 未着手 | 6 / 11（シナリオ） | [operator.md](docs/units/operator.md) |
| 3 | access-control | 未着手 | 未着手 | 未着手 | 0 / 0（分類のみ。D-009） | [access-control.md](docs/units/access-control.md) |
| 4 | catalog | 未着手 | 未着手 | 未着手 | 0 / 0（分類のみ。D-010） | [catalog.md](docs/units/catalog.md) |
| 5 | notification | 未着手 | 未着手 | 未着手 | 31 / 31（候補の確認。単体テストへ回すもの 55 件は D-011） | [notification.md](docs/units/notification.md) |
| 6 | identity-verify | 未着手 | 未着手 | 未着手 | 0 / 0（分類のみ。D-012） | [identity-verify.md](docs/units/identity-verify.md) |
| 7 | certificate | 未着手 | 未着手 | 未着手 | 0 / 0（分類のみ。D-013） | [certificate.md](docs/units/certificate.md) |
| 8 | binary | 未着手 | 未着手 | 未着手 | 0 / 0（分類のみ。D-014） | [binary.md](docs/units/binary.md) |
| 9 | ctoken | 未着手 | 未着手 | 未着手 | 0 / 0（分類のみ。D-015） | [ctoken.md](docs/units/ctoken.md) |
| 10 | book-manage | 未着手 | 未着手 | 未着手 | 0 / 0（分類のみ。D-016） | [book-manage.md](docs/units/book-manage.md) |
| 11 | book-operate | 未着手 | 未着手 | 未着手 | 0 / 0（分類のみ。D-017） | [book-operate.md](docs/units/book-operate.md) |

状態は「未着手 / 作業中 / 完了」で書きます。

## 対象外

次の 2 群は、移植の対象から明示で外します。決めたのは 2026-10-09、ヤマシタです。

- 非公開レポ：ddl、catalog、manifest、fluentd、nginx、ext-idp-linkage-plugin、test（テスト用ツールとテスト結果）を含む群。組織内で読み取れず、公開された OSS の範囲ではないため対象外とします。
- 公開レポにない呼び出し先：info-account-manage、audit、outsideStoreService、share-trigger。組織にも pxr-linkage にもソースがなく、出し口として扱うため対象外とします。

対象は、組織の公開レポにあるサービスです。

## 旧基盤の E2E の推移

| 日付 | 対象 | 通過 | 失敗 | 備考 |
| --- | --- | --- | --- | --- |
| | 旧実装（ベースライン） | | | フェーズ0の 6.4 |
| | Funadansu | 0 | | フェーズ0の 6.3（全部 RED） |

## リリース

| タグ | 日付 | フェーズ | 内容 |
| --- | --- | --- | --- |
| | | | |

## この文書の更新

- 単位の作業が終わったら、その単位の行と E2E の推移を更新します
- フェーズのタスクが終わったら、状態を「✅ 完了」にします
- 外部からの Issue を受けるようになったら（公開と MIT への移行の時点）、状況の管理を Issue に移し、この文書は索引にします
