# 差分台帳

公開されている OpenAPI 仕様（`openapi/`）、元の実装（PxR）、Funadansu の間で、振る舞いが異なる箇所を記録します。

## 書き方

- 1つの差分につき1つの項目。番号は `D-001` から順に振ります
- **評価の言葉を使わず、何を、なぜ変えたかだけ**を書きます
- 脆弱性やその疑いに関わる差分は、ここには書きません。[SECURITY.md](../SECURITY.md) の窓口で扱い、元の実装の側で修正が公開された後に、必要な範囲で追記します

### 種別

| 種別 | 意味 |
| --- | --- |
| 仕様に合わせた | 元の実装が仕様と異なっていたため、Funadansu は仕様に合わせた |
| 意図した変更 | API の形は変えずに、内部の仕組みや性能を変えた |
| 非互換 | API の形や結果が元の実装と異なる（原則として作らない。作る場合は理由を必ず書く） |
| 要確認 | 仕様と実装のどちらに合わせるかを決めていない |

## 一覧

| 番号 | 単位 | 対象 | 種別 | 状態 |
| --- | --- | --- | --- | --- |
| D-001 | proxy、operator | 認証とセッション | 意図した変更 | 方針決定済み |
| D-002 | book-operate | 大量件数の一括登録 | 要確認 | フェーズ0で計測 |
| D-003 | operator | 利用者の情報の検索（`pxrId`） | 要確認 | 仕様の確認待ち |
| D-004 | operator | 旧 E2E の日付（テストデータの期限） | 要確認 | 方針決定済み（赤のまま残す） |
| D-005 | operator | DB の列の NOT NULL（旧 DB の実物との差） | 要確認 | フェーズ3で比べる |
| D-006 | operator | jest.mock に依る旧 E2E（E2E から外す） | 要確認 | 方針決定済み（単体テストへ） |

## 項目

### D-001 認証とセッション

| 項目 | 内容 |
| --- | --- |
| 単位 | proxy、operator |
| 対象 | ログイン、セッション確認のエンドポイント |
| 仕様 | 旧 operator の OpenAPI（`pxr-operator-service/config/openapi.json`）の `POST /login`、`POST /session`、`POST /logout`、`POST /ind/login`、`POST /ind/session`、`POST /ind/logout`（`openapi/legacy/operator.json`） |
| 元の実装 | 独自のセッションを、リクエストごとにデータベースで確認する |
| Funadansu | ログイン時にパスワードを確認して JWT を発行し、リクエストごとには JWT の署名を検査する。エンドポイントの形は変えない |
| 種別 | 意図した変更 |
| 理由 | リクエストごとのデータベースへのアクセスをなくし、入口（proxy）で認証を完結させるため |
| 影響 | 切り替えた時点で、元の実装で発行されたセッションは無効になり、利用者は一度ログインし直す。元の実装のセッションのテーブルはフェーズ1・2では使わず、変更もしない |
| 関連 | なし（Issue は未起票） |

### D-002 大量件数の一括登録

| 項目 | 内容 |
| --- | --- |
| 単位 | book-operate |
| 対象 | 利用者の一括登録 |
| 仕様 | 旧 book-operate の OpenAPI（`pxr-book-operate-service/config/openapi.json`）の `POST /user/batch`（`openapi/legacy/book-operate.json`） |
| 元の実装 | 約 200 件の一括登録で、応答時間の上限を超える事例がある（フェーズ0で計測して確かめる） |
| Funadansu | 仕様どおりの件数で完了することを目標とする |
| 種別 | 要確認 |
| 理由 | — |
| 関連 | なし（Issue は未起票） |

### D-003 利用者の情報の検索（`pxrId`）

| 項目 | 内容 |
| --- | --- |
| 単位 | operator |
| 対象 | `/user/info`、`/ind/user/info` の `pxrId` のパラメータ |
| 仕様 | 旧 operator の OpenAPI（`openapi/legacy/operator.json`）の `pxrId` は `number`（`GET`、`DELETE` の query、`/ind/user/info` の `GET`、`PUT` の query） |
| 元の実装 | 未確認（旧実装の検索の型はまだ見ていない）。他の旧仕様では `pxrId` は `string`、または文字列の配列（`openapi/legacy/identity-verificate.json` など） |
| Funadansu | 未決。仕様と PXR-ID（文字列）のどちらに合わせるかを決める |
| 種別 | 要確認 |
| 理由 | 仕様は `number` だが、PXR-ID は業務上の識別子で、旧スキーマでは `varchar(255)`（[規約](conventions.md)）。数値の検索にすると、文字列の PXR-ID と合わない |
| 影響 | 合わせ方によって、`pxrId` を送るクライアントの振る舞いが変わる |
| 関連 | 単位 operator の移行の Issue で決める |

### D-004 旧 E2E の日付（テストデータの期限）

| 項目 | 内容 |
| --- | --- |
| 単位 | operator |
| 対象 | 旧 E2E のテストデータ（ワンタイムログインのコード、セッションの期限） |
| 仕様 | 旧 operator のテスト（`src/tests/07-01.OperatorOneTimeLogin.spec.ts`、`08-01.OperatorLogout.spec.ts`） |
| 元の実装 | テストデータの期限が 2025-01-27 に固定されている。今日の日付では期限切れで、元のテストは一部が落ちる |
| Funadansu | 元のテストの期限を変えない。赤のまま残す |
| 種別 | 要確認 |
| 理由 | テストの中身を書き換えて通すことを避けるため（AGENTS.md の互換性）。期限だけを延ばした診断版では全件通ることを確かめた |
| 影響 | 旧 E2E の本数は、この 2 件（07-01）と、期限切れの各件で赤になる |
| 関連 | なし（Issue は未起票） |

### D-005 DB の列の NOT NULL（旧 DB の実物との差）

| 項目 | 内容 |
| --- | --- |
| 単位 | operator |
| 対象 | `pxr_operator.operator` の `pxr_id`、`user_information`、`last_login_at` などの NOT NULL |
| 仕様 | 旧 operator のテスト（INSERT で、これらの列に null を入れている） |
| 元の実装 | 実物の DB では、これらの列を null にできる（テストが入れている） |
| Funadansu | `db/pxr_operator/001_tables.sql` は ORM のエンティティから起こしたため、これらの列が NOT NULL になっている。そのままでは旧 E2E が入れられない |
| 種別 | 要確認 |
| 理由 | 実物の DB との比較は、フェーズ3の移行の作業で扱う（`db/README.md` の方針）。この作業では DDL を変えず、ローカルの検証用 DB だけ NOT NULL を外した |
| 影響 | 検証用 DB と DDL が一致しない。旧 E2E の結果は、ローカルの検証用 DB で確かめたもの |
| 関連 | なし（Issue は未起票） |

### D-006 jest.mock に依る旧 E2E（E2E から外す）

| 項目 | 内容 |
| --- | --- |
| 単位 | operator |
| 対象 | 旧 operator の試験のうち、jest.mock（リポジトリや関数の差し替え）で内部の失敗を起こすもの |
| 仕様 | 旧 operator のテスト：`01-02.OperatorAdd.libgetOperatorError.spec.ts`（1 件）、`04-02.OperatorUpdate.libgetOperatorError.spec.ts`（1 件）、`04-03.OperatorUpdate.libisSessionIdExistsError.spec.ts`（1 件）、`04-04.OperatorUpdate.libisAllAuthMemberExistsOtherThisIdError.spec.ts`（1 件）、`05-02.OperatorDelete.libgetOperatorError.spec.ts`（1 件）、`06-02.OperatorLogin.GeneratesSessionIdIfAlreadyUse.spec.ts`（1 件）、`06-03.OperatorLogin.GeneratesLoginIdIfAlreadyUse.spec.ts`（1 件）、`12-01.RoleAndAuth.spec.ts`（3 件） |
| 元の実装 | 試験の中で、リポジトリや関数を差し替えて、ライブラリエラーや ID の重複を起こす |
| Funadansu | E2E（runn）には入れない。単体テストに回し、同じ失敗の経路を Go 側の単体テストで確かめる |
| 種別 | 要確認 |
| 理由 | runn は外から HTTP で叩くため、差し替えた内部の失敗は起こせない。差し替えを runn で代用すると、試験の中身を変えることになるため（AGENTS.md の互換性） |
| 影響 | E2E の本数から、この件数が外れる。単体テストで置き換えるまでは、これらの失敗の経路は検証されない。`13-01.UserInfo.spec.ts` は、`config/config.json` の `Config` を、同じ値のインラインの設定で差し替えているだけなので（18 項目すべて一致）、この項目には入れない |
| 関連 | なし（Issue は未起票） |

<!--
### D-000 <題>

| 項目 | 内容 |
| --- | --- |
| 単位 | |
| 対象 | <エンドポイントや機能> |
| 仕様 | <OpenAPI の該当箇所> |
| 元の実装 | |
| Funadansu | |
| 種別 | 仕様に合わせた / 意図した変更 / 非互換 / 要確認 |
| 理由 | |
| 影響 | |
| 関連 | <Issue> |
-->
