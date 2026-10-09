# 規約

Funadansu の構成、命名、書き方の決定です。人と AI エージェントが同じ書き方をできるよう、ここに決めて書きます。
見本は `bootstrap/`（P0-4）に置き、各決定の末尾で示します。

この文書は、アーキテクチャやコードベースの改善に合わせて随時更新します。決めたことが変わったら、ここを直し、変えた理由を書き添えます。

この文書は、次の順で書き足します。

| 作業 | 内容 | 状態 |
| --- | --- | --- |
| 3.1 | ディレクトリとパッケージの構成、命名、DI を使わない配線（本文の 1〜3 節） | 決定（この版） |
| 3.2 | ID の型、エラーの形、ログの形 | 未着手 |
| 3.3 | テストの方針（本文の 5 節）、「外に出す関数の一覧」の書式（本文の 6 節） | 決定（この版） |

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

## 5. テストの方針

### 決定

テストは次の3種類を使い分けます。

- **Go の単体テスト**（`*_test.go`、標準の `testing`）：`service.go` の業務の処理と、`platform/` の部品を確かめます。表駆動（`t.Run` で各行を試す）で書きます。
- **エンドポイントのテスト**（標準の `net/http/httptest`）：`handler.go` の入口を、旧 API のパスと形のまま確かめます。ステータス、JSON の形、エラーの形を見ます。
- **API のシナリオテスト**（runn、YAML）：複数のエンドポイントをつなぐ流れを確かめます。テストデータもここで投入します。置き場は `test/api/<単位>/` です。
- **旧基盤の E2E**（`test/e2e/`）は、互換を確かめる**外の物差し**です。Funadansu の実装に合わせて書き換えません。フェーズ0では Funadansu に向けて全部 RED から始めます（[MIGRATION.md](../MIGRATION.md) の 6.x）。
- **モックは極力使いません**。DB を使う部分は、モックの代わりに、`db-start` で起動するローカルの PostgreSQL で試します。
- **テストデータ**は、単体テストでは各テストの中に書き、API のシナリオでは runn の中に書きます。共通の fixture ファイルは、必要になってから決めます。
- **旧 E2E が仕様と食い違う場合**は、テストを直さず、[差分台帳](spec-deviations.md) に書き、Issue で相談します（AGENTS.md の「互換性」）。
- **変更には、それを確かめるテストを付けます**。テストの実行コマンド（`make test`、`make e2e BASE_URL=...`）は、フェーズ0で確定させます。

### 既定の選択

- runn は nix の開発シェルで入れ、バージョンを固定します（Go の依存には入れません）。
- 旧 E2E は、フェーズ0の 6.x で `test/e2e/` に置きます。エンドポイントのテストと区別するため、単体テストの中では旧 E2E を呼びません。

### 見本

```go
// internal/operator/service_test.go（実際の DB を使う。モックは使わない）
func TestFindOperator(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		want    Operator
		wantErr error
	}{
		{name: "存在する", id: "op-1", want: Operator{ID: "op-1"}},
		{name: "存在しない", id: "op-x", wantErr: ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(testStore(t)) // 各テストで DB を用意する（testStore は同じパッケージの helper）
			got, err := svc.FindOperator(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got = %+v, want %+v", got, tt.want)
			}
		})
	}
}
```

```go
// internal/operator/handler_test.go（httptest で入口を確かめる）
func TestLoginHandler(t *testing.T) {
	h := NewHandler(NewService(testStore(t)))
	req := httptest.NewRequest(http.MethodPost, "/operator/login", strings.NewReader(`{"id":"op-1"}`))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
```

```yaml
# test/api/operator/login.yml（runn。テストデータは steps の中に書く）
desc: オペレーターのログイン
runners:
  req: http://localhost:8080
steps:
  login:
    req:
      /operator/login:
        post:
          body:
            application/json:
              id: op-1
    test: current.res.status == 200
```

見本の置き場：`bootstrap/internal/operator/`、`test/api/`（P0-4 で作ります）。上の例は、形を示すための見本で、実際の API の内容は P0-4 で確かめます。

## 6. 「外に出す関数の一覧」の書式

### 決定

- **場所**：`docs/units/<単位>.md` の「外に出す関数の一覧」の節に書きます。書式は [_template.md](units/_template.md) の表を使います（関数、入力、出力、説明）。
- **関数の欄**：Go の関数名と、引数・戻り値を、実際のシグネチャどおりに書きます。`ctx` は必ず最初の引数です（本文 3 節）。
- **入力・出力の欄**：引数と戻り値の型を書きます。型が別の単位のものなら、その単位の一覧に載っている型だけを使います。
- **対象**：他の単位が呼んでよい関数だけを書きます。単位の中だけで使う関数は書きません。1行に1関数です。
- **配線**：呼ぶ側の単位は、一覧に載っている関数だけを、`cmd/funadansu/main.go` でコンストラクタの引数として渡します（本文 3 節）。一覧に無い関数は渡しません。
- **変更の手順**：一覧を変えたら、同じ作業の中で、呼ぶ側の単位の引き継ぎファイルの「依存」の表も直します。関数を消すとき、シグネチャを変えるときは、呼ぶ側を先に直します。
- **確かめ方**：レビューで、一覧と `main.go` の配線を照らし合わせます。

### 既定の選択

- 一覧は**手で書きます**。生成はしません（一覧の正しさは、人がレビューで確かめます）。
- 「依存」の表（引き継ぎファイルの 2 節）と、この一覧はセットで直します。

### 見本

次は架空の単位（`operator`）の引き継ぎファイルに書く例です。

```markdown
## 外に出す関数の一覧

| 関数 | 入力 | 出力 | 説明 |
| --- | --- | --- | --- |
| `Find(ctx context.Context, id OperatorID) (Operator, error)` | `id`：オペレーターの ID | `Operator`、`error`（見つからないときは `ErrNotFound`） | オペレーターを ID で探す |
```

呼ぶ側（例：`book-manage`）は、`main.go` で次のように渡します。

```go
// cmd/funadansu/main.go（配線はここだけ）
operators := operator.NewService(store.New(pool))
bookManage := bookmanage.NewService(operators.Find)
```

見本の置き場：`bootstrap/cmd/funadansu/main.go`、`docs/units/_template.md`（P0-4 で作ります）。
