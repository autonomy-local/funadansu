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
| 4.5 | bootstrap：疎通（ローカル、リモート）。リモートは 5.3 の構成が要る（[Issue #20](https://github.com/autonomy-local/funadansu/issues/20)） | 未着手 |
| 5.1 | GCP のプロジェクト、Cloudflare、Pulumi の状態の置き場所（外部の準備。[Issue #17](https://github.com/autonomy-local/funadansu/issues/17)） | 未着手（外部の準備が要る） |
| 5.2 | sops と age の設定、JWT の鍵（dev のダミー）の暗号化（[docs/operations/secrets.md](docs/operations/secrets.md)） | ✅ 完了（[Issue #18](https://github.com/autonomy-local/funadansu/issues/18)） |
| 5.3 | Pulumi のスタック（dev）：Artifact Registry、Cloud Run、Workers、シークレットの配布（外部の準備が要る。[Issue #19](https://github.com/autonomy-local/funadansu/issues/19)） | 未着手（外部の準備が要る） |
| 5.4 | 最小のデプロイの手順（ローカルの CI → `pulumi up`・`nixos-rebuild`）と運用の手引きの初版（`docs/operations/`） | 未着手 |
| 6.1 | 旧 OpenAPI を `openapi/legacy/` に集める（取得元のコミットを [openapi/README.md](openapi/README.md) に記録） | ✅ 完了 |
| 6.2 | 旧基盤の E2E の書き換え（operator と access-control、catalog、notification、book-manage の一部）。残りの単位は試験の分類のみ | 作業中（[docs/units/](docs/units/README.md)） |
| 6.3 | 旧基盤の E2E（URL の切り替え、Funadansu に向けて全部 RED） | ✅ 完了 |
| 6.4 | （推奨）旧実装のベースライン（[docs/e2e-baseline.md](docs/e2e-baseline.md)。厳密 11 / 32、差分台帳 D-018） | ✅ 完了 |
| 7.1 | CI（GitHub Actions、Nix の開発シェル。`check` は Go・proxy・Bootstrap のシナリオ、`legacy-e2e` は旧 E2E の通過数を出し、落ちても CI は落とさない。[.github/workflows/ci.yml](.github/workflows/ci.yml)） | ✅ 完了（[Issue #25](https://github.com/autonomy-local/funadansu/issues/25)。Nix の開発シェルでの緑は、PR の CI で確かめる。この環境からは flake の入力を取れなかった） |
| 7.2 | 脆弱性と秘密情報の検査（govulncheck、CodeQL・gosec、秘密情報の検出、Actions の SHA 固定と最小権限、Renovate） | 未着手 |
| 7.3 | （フェーズ1までに）SBOM、署名、OpenSSF Scorecard | 未着手 |
| 8.1 | 旧実装のファイル数・行数・依存パッケージ数（[docs/legacy-measurement.md](docs/legacy-measurement.md)） | ✅ 完了 |
| 8.2 | 旧構成（15コンテナ、ALB、Nginx、Aurora）の月額の見積もり（[docs/legacy-cost-estimate.md](docs/legacy-cost-estimate.md)） | ✅ 完了 |
| 8.3 | 測定の基準（所要時間。旧 catalog の一括登録、[docs/legacy-timing.md](docs/legacy-timing.md)。検証用の DB で測った。[Issue #28](https://github.com/autonomy-local/funadansu/issues/28)） | ✅ 完了 |

## 単位の状況

| 順 | 単位 | フェーズ1 | フェーズ2 | フェーズ3 | 旧 E2E（通過／全体） | 引き継ぎファイル |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | proxy | 未着手 | — | — | 0 / 0（分類のみ。D-008） | [proxy.md](docs/units/proxy.md) |
| 2 | operator | 未着手 | 未着手 | 未着手 | 6 / 11（シナリオ） | [operator.md](docs/units/operator.md) |
| 3 | access-control | 未着手 | 未着手 | 未着手 | 137 / 137（候補 137 件。access-control-manage の 75 件と、access-control の 62 件。単体テストへ回すもの 184 件は D-009） | [access-control.md](docs/units/access-control.md) |
| 4 | catalog | 未着手 | 未着手 | 未着手 | 1187 / 1187（pxr-catalog-service の候補 7 ファイル 1144 件と、catalog-update の候補のうち 43 件。catalog-update の残り 19 件は単体テストへ。D-010） | [catalog.md](docs/units/catalog.md) |
| 5 | notification | 未着手 | 未着手 | 未着手 | 31 / 31（候補の確認。単体テストへ回すもの 55 件は D-011） | [notification.md](docs/units/notification.md) |
| 6 | identity-verify | 未着手 | 未着手 | 未着手 | 0 / 0（分類のみ。D-012） | [identity-verify.md](docs/units/identity-verify.md) |
| 7 | certificate | 未着手 | 未着手 | 未着手 | 0 / 0（分類のみ。D-013） | [certificate.md](docs/units/certificate.md) |
| 8 | binary | 未着手 | 未着手 | 未着手 | 0 / 0（分類のみ。D-014） | [binary.md](docs/units/binary.md) |
| 9 | ctoken | 未着手 | 未着手 | 未着手 | 0 / 0（分類のみ。D-015） | [ctoken.md](docs/units/ctoken.md) |
| 10 | book-manage | 未着手 | 未着手 | 未着手 | 4 / 4（候補の確認。単体テストへ回すもの 1920 件は D-016） | [book-manage.md](docs/units/book-manage.md) |
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
| 2026-10-10 | 旧実装（ベースライン） | 11 | 21 | フェーズ0の 6.4。旧実装の 7 サービス（固定したコミット）に、同じ 32 本を流した。DDL のままの厳密な数字。NOT NULL を外した診断版は 24 / 8。内訳と原因は [docs/e2e-baseline.md](docs/e2e-baseline.md) |
| 2026-10-10 | Funadansu | 0 | 32 | フェーズ0の 6.3。runn のシナリオ 32 本（operator 11、access-control 7、catalog 7、catalog-update 4、notification 1、book-manage 2）をすべて RED。接続先は proxy（127.0.0.1:8787、その背後に Go を 8080）、DB は PostgreSQL 16 の `funadansu`。内訳は各単位の「旧 E2E の状況」 |

## リリース

| タグ | 日付 | フェーズ | 内容 |
| --- | --- | --- | --- |
| | | | |

## この文書の更新

- 単位の作業が終わったら、その単位の行と E2E の推移を更新します
- フェーズのタスクが終わったら、状態を「✅ 完了」にします
- 外部からの Issue を受けるようになったら（公開と MIT への移行の時点）、状況の管理を Issue に移し、この文書は索引にします
