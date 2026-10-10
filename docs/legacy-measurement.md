# 旧実装の規模（ファイル数・行数・依存パッケージ数）

フェーズ0の 8.1（[Issue #26](https://github.com/autonomy-local/funadansu/issues/26)）の記録です。移行の前後を比べるための基準で、Funadansu 側の測定はしていません。

- 対象：[openapi/README.md](../openapi/README.md) に載せた**公開の 15 サービス**（旧実装の 15 のマイクロサービス）。固定したコミットは同じ表のもの
- 対象の外：非公開レポの群と、公開されていない 4 つの呼び出し先（[MIGRATION.md](../MIGRATION.md) の「対象外」）
- 測定の日付：2026-10-10
- 版：下の各サービスのコミット（既定のブランチの最新ではなく、openapi/ の取得時に固定したもの）

## 結果

### 全体

| 項目 | 値 | 数えた範囲 |
| --- | --- | --- |
| サービス数 | 15 | 公開リポジトリ |
| ファイル数（追跡されているもの） | 2,445 | `git ls-files` の件数の合計 |
| ソース（TypeScript・JavaScript）の行数 | 742,858 | うち本体 149,020、テスト 593,838 |
| 本体のソースファイル数 | 1,371 | 同上（テストを除く） |
| SQL の行数 | 10,453 | `*.sql` |
| 直接の依存（`package.json`） | 856 件（実行時 349、開発時 507） | 合計。同じパッケージの重複を含む |
| 直接の依存の種類（重複なし） | 73 件（実行時のみ 33 件） | 15 サービスの和集合 |
| 依存パッケージ（lockfile の全エントリ） | 14,611 件 | 15 サービスの `package-lock.json` の `packages` の件数の合計。推移的な依存を含み、サービス間の重複を含む |

行数は、テストのファイル（`test`・`tests`・`__tests__` のディレクトリ、`*.spec.*`・`*.test.*`）と、`*.d.ts` の扱いが効きます。テストの行数が大きいのは、旧 E2E ではなく、サービスごとの単体テストとスタブのためです。

### サービスごと

| 旧サービス（公開リポジトリ） | 単位 | コミット | 日付 | ファイル数 | 本体の行 | テストの行 | 直接の依存（実行時＋開発時） | lockfile のエントリ |
| --- | --- | --- | --- | ---: | ---: | ---: | --- | ---: |
| pxr-operator-service | operator | `453e638` | 2025-03-26 | 188 | 12,146 | 21,027 | 25＋34 | 1,034 |
| pxr-access-control-service | access-control | `5cd2d6a` | 2025-03-04 | 97 | 3,844 | 6,953 | 23＋34 | 978 |
| pxr-access-control-manage-service | access-control | `3aabd9a` | 2025-03-04 | 150 | 5,545 | 12,239 | 23＋33 | 1,058 |
| pxr-catalog-service | catalog | `14a647d` | 2025-03-04 | 285 | 20,735 | 112,691 | 24＋35 | 988 |
| pxr-catalog-update-service | catalog | `085f312` | 2025-03-12 | 215 | 12,660 | 74,415 | 23＋34 | 833 |
| pxr-notification-service | notification | `30b7ce8` | 2025-03-04 | 103 | 3,552 | 4,830 | 22＋33 | 984 |
| pxr-book-manage-service | book-manage | `3e2082a` | 2025-03-06 | 444 | 34,507 | 184,788 | 25＋34 | 1,071 |
| pxr-book-operate-service | book-operate | `9d51020` | 2025-03-04 | 230 | 23,789 | 122,963 | 23＋35 | 1,058 |
| pxr-binary-manage-service | binary | `865ae2b` | 2025-03-04 | 112 | 4,078 | 6,107 | 23＋34 | 986 |
| pxr-block-proxy-service | block-proxy | `16aebfe` | 2025-03-04 | 76 | 3,427 | 8,412 | 24＋35 | 997 |
| pxr-certificate-manage-service | certificate | `92cd428` | 2025-03-04 | 66 | 2,177 | 2,016 | 22＋34 | 986 |
| pxr-certification-authority-service | certification-authority | `6d3d587` | 2025-03-04 | 157 | 7,501 | 9,896 | 24＋34 | 991 |
| pxr-ctoken-ledger-service | ctoken | `643c993` | 2025-03-04 | 108 | 6,471 | 12,893 | 23＋34 | 986 |
| pxr-identity-verificate-service | identity-verify | `5628165` | 2025-03-26 | 147 | 5,301 | 6,049 | 23＋33 | 830 |
| pxr-local-ctoken-service | local-ctoken | `1083d74` | 2025-03-04 | 67 | 3,287 | 8,559 | 22＋31 | 831 |

注：サービスごとのファイル数は `git ls-files` の件数で、合計の 2,445 と一致します。単位の欄は [MIGRATION.md](../MIGRATION.md) の単位に対応させた参考の情報です（access-control と catalog は、それぞれ 2 つのサービスを含みます）。

## 数えた方法

コミットは [openapi/README.md](../openapi/README.md) の表と同じです。作業用のディレクトリ（リポジトリの外）に、各サービスを 1 コミットだけ取って、次を流しました。依存の数え方（6・7）だけ `node` を使い、ほかは `git` と `grep`・`wc` です。

```sh
# 1. 固定したコミットを取る（openapi/README.md の表の SHA）
git clone --filter=blob:none --no-checkout https://github.com/Personal-Data-Linkage-Module/<サービス>.git
cd <サービス> && git fetch --depth 1 origin <SHA> && git checkout <SHA>

# 2. ファイル数（追跡されているもの）
git ls-files | wc -l

# 3. ソースの行数（本体。テストと *.d.ts を除く）
git ls-files | grep -E '\.(ts|js|mjs|cjs|tsx|jsx)$' | grep -vE '\.d\.ts$' \
  | grep -viE '(^|/)(test|tests|__tests__)/|\.(spec|test)\.' | xargs cat | wc -l

# 4. ソースの行数（テスト）：上の 3 の grep -v を grep -i に変えたもの

# 5. SQL の行数
git ls-files | grep -E '\.sql$' | xargs cat | wc -l

# 6. 直接の依存（実行時と開発時）
node -e 'const p=require("./package.json");console.log(Object.keys(p.dependencies||{}).length, Object.keys(p.devDependencies||{}).length)'

# 7. lockfile のエントリ（ルートの "" を除く）
node -e 'const l=require("./package-lock.json");console.log(Object.keys(l.packages).filter(k=>k!=="").length)'
```

重複のない直接の依存は、15 サービスの `package.json` の `dependencies` と `devDependencies` のキーを集めて `sort -u` で数えました（73 件）。

## 読み方

- 本体の行数（約 15 万行）と、テストの行数（約 59 万行）は別々に見てください。移行で書き直すのは本体の側で、テストは旧 E2E と単体テストの仕様として残す量の目安です（[docs/legacy-unit-test-inventory.md](legacy-unit-test-inventory.md)）。
- 直接の依存は、サービスごとにほぼ同じ数です（22〜25 件の実行時、31〜35 件の開発時）。重複のない種類は 73 件です。
- lockfile の数は推移的な依存を含むため、同じパッケージが多くのサービスで数えられています。「何を入れているか」の目安には、直接の依存の種類を使ってください。
- 旧実装の側は、この測定で何も変えていません（取得のみ）。
