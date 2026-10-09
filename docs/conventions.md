# 規約

Funadansu の構成、命名、書き方の決定です。人と AI エージェントが同じ書き方をできるよう、ここに決めて書きます。
見本は `bootstrap/`（P0-4）に置き、各決定の末尾で示します。

この文書は、アーキテクチャやコードベースの改善に合わせて随時更新します。決めたことが変わったら、ここを直し、変えた理由を書き添えます。

この文書は、次の順で書き足します。

| 作業 | 内容 | 状態 |
| --- | --- | --- |
| 3.1 | ディレクトリとパッケージの構成、命名、DI を使わない配線（本文の 1〜3 節） | 決定（この版） |
| 3.2 | ID の型、エラーの形、ログの形 | 未着手 |
| 3.3 | テストの方針、「外に出す関数の一覧」の書式 | 未着手 |

## 1. ディレクトリとパッケージの構成

### 決定

- Go は、**1つのモジュールで1つのバイナリ**を作ります。入口は `cmd/funadansu/` だけです。
- 共通の基盤は `internal/platform/` に置きます（HTTP、エラー、ログ、ヘルスチェック、設定の読み込み）。
- 単位（ドメイン）ごとに `internal/<パッケージ名>/` を1つ作ります。単位の一覧は [docs/units/README.md](units/README.md) です。
- 単位の中は次のとおりです。
  - `handler.go`：HTTP の入口。旧 API のパスと形をそのまま受けます。
  - `service.go`：業務の処理。
  - `store/`：sqlc の入力（`query.sql`）と生成物（`*.go`）。SQL が正本で、生成物は手で直しません。`store/` は同じ単位の中からだけ import します。
- proxy（Hono）は TypeScript で `proxy/` に置きます。
- OpenAPI の仕様は `openapi/` に置きます（正本）。旧 DDL は `db/` に置きます。
- 単位の移植か新規かの区別は、コミットメッセージと引き継ぎファイルの「移植と新規」に書きます。

### 既定の選択（確認が要るもの）

- Issue には細かい配置が無いため、上の形を既定とします。見本は `bootstrap/` の中に、同じ形で置きます（`bootstrap/cmd/funadansu/`、`bootstrap/internal/platform/`、`bootstrap/internal/<単位>/`）。
- 単位の移植が始まったあとに、この形を `bootstrap/` の外（リポジトリ直下）へ移すかは、フェーズ1で改めて決めます。

### 見本

- `bootstrap/cmd/funadansu/main.go`：配線をする唯一の場所（「DI を使わない配線」を参照）
- `bootstrap/internal/platform/`：共通の基盤
- `bootstrap/internal/<単位>/`：handler、service、store の組み方

## 2. 命名

### 決定

- **Go のパッケージ名**は、小文字だけで、ハイフンもアンダースコアも使いません。ディレクトリ名とパッケージ名は一致させます。
  - 例：`identity-verify` → `internal/identityverify/`、`book-manage` → `internal/bookmanage/`、`ctoken` → `internal/ctoken/`
- **Go の識別子**は標準の書き方に従います（`ID`、`URL`、`JWT` のように頭字語は全部大文字）。
- **ファイル名**は小文字で、単語の区切りが必要なときだけ `_` を使います（例：`service_test.go`）。
- **単位の名前**は、文書とファイル名（`docs/units/identity-verify.md`、`openapi/identity-verify.yaml`）ではハイフン区切り、Go のパッケージ名では区切りなしにします。
- **TypeScript（proxy）**のファイル名は kebab-case、識別子は camelCase にします。
- **API のパス**は、旧実装と同じ形のまま使います。名前を変えません（AGENTS.md の互換性の決め）。
- **SQL の名前**は、旧スキーマの名前（`pxr_operator` など）をそのまま使います。新しく作るスキーマは `funadansu_` の接頭辞を付けます（例：`funadansu_auth`）。
- **sqlc のクエリ名**は、sqlc の既定（PascalCase、`-- name: FindOperator :one`）に従います。

### 見本

- `bootstrap/internal/<単位>/` のディレクトリ名、`bootstrap/internal/<単位>/store/query.sql` の `-- name:` の書き方

## 3. DI を使わない配線の書き方

DI コンテナは使いません。依存は**コンストラクタの引数**で受け取り、`main` で手で配線します。

### 決定

- **依存の受け取り**：`New<型名>` の関数で、必要なものを引数に取ります。`init()` と、パッケージ変数に置く DB や設定の値は使いません。
- **インターフェース**は、**使う側のパッケージ**に、使う分だけ定義します。作る側のパッケージでインターフェースを先に作りません。
- **配線**は `cmd/funadansu/main.go` だけで行います。ほかのパッケージは、受け取った依存を使うだけで、どこから来たかを知りません。
- **設定と接続**は `main` で一度だけ読み、開いた接続は `main` で閉じます（`defer`）。
- **単位どうしの呼び出し**は、呼ぶ側のコンストラクタに、相手の関数を値として渡します。相手が「外に出す関数の一覧」（引き継ぎファイル）に無い関数は、渡しません。
- **context** は、関数の最初の引数で受け渡します。
- **proxy（Hono）**も同じ形です。実行環境（Workers と Bun）の違いは、`proxy/src/entry/` の入口だけに閉じ込め、`createApp(deps)` に依存を渡します。

### 既定の選択

- インターフェースの名前は、役割で付けます（例：`Store`、`AccessChecker`）。実装の名前（`Postgres...`）は付けません。
- 1つの単位に、依存を渡すための構造体は作りません。引数が多くなったら、関数を分けることを先に考えます。

### 見本

```go
// internal/operator/service.go（使う側でインターフェースを定義する）
type Store interface {
	FindOperator(ctx context.Context, id string) (Operator, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}
```

```go
// cmd/funadansu/main.go（配線はここだけ）
pool := platform.OpenDB(cfg.DatabaseURL)
defer pool.Close()

operators := operator.NewService(store.New(pool))
```

見本の置き場：`bootstrap/internal/platform/`、`bootstrap/internal/<単位>/service.go`、`bootstrap/cmd/funadansu/main.go`（P0-4 で作ります）。

## 4. 移植した関数の出典

旧実装（`legacy/` の該当サービス）から移した関数には、出典を残します。NOTICE の表示と、MIT License の条件を満たすためです。

### 決定

- **関数レベルの出典**：移植元の GitHub の該当コードへのリンクと、関数名を、関数の直前のコメントに書きます。リンクは、旧実装のコミットのハッシュを含む固定のリンクにします（`main` のリンクは、後で中身が変わるため使いません）。
- **集約や変更がある場合**：複数の旧関数をまとめたとき、または処理を変えたときは、Go のコードコメントで、どの関数から何を変えたかを書きます。
- **書く場所**：コメントは関数の直前に置きます。PR 本文の「移植か、新規か」と引き継ぎファイルの「移植と新規」にも、同じ出典を書きます。
- **新規に書いた関数**には、出典のコメントを付けません。

### 見本

```go
// Migrated from: https://github.com/autonomy-local/<旧リポジトリ>/blob/<コミット>/src/<ファイル>#L10-L40
// function: findOperatorByID
// Changed: 旧実装の二つのクエリを、store の一つのクエリにまとめた。
func (s *Service) FindOperator(ctx context.Context, id string) (Operator, error) {
	// ...
}
```

見本の置き場：`bootstrap/internal/<単位>/`（P0-4 で作ります）。
