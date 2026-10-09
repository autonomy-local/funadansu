# 規約

Funadansu の構成、命名、書き方の決定です。人と AI エージェントが同じ書き方をできるよう、ここに決めて書きます。
見本は `bootstrap/`（P0-4）に置き、各決定の末尾で示します。

この文書は、アーキテクチャやコードベースの改善に合わせて随時更新します。決めたことが変わったら、ここを直し、変えた理由を書き添えます。

この文書は、次の順で書き足します。

| 作業 | 内容 | 状態 |
| --- | --- | --- |
| 3.1 | ディレクトリとパッケージの構成、命名、DI を使わない配線（本文の 1〜3 節） | 決定（この版） |
| 3.2 | ID の型、エラーの形、ログの形（本文の 7〜9 節） | 決定（この版） |
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

## 7. ID の型

### 決定

- **内部の主キー**は、旧スキーマと同じ `bigint GENERATED BY DEFAULT AS IDENTITY` で採番します（`db/pxr_operator/001_tables.sql` など）。Go では `int64` で受けます。アプリ側で ID を作って主キーにしません。
- **Go の型は、単位ごとに名前を付けます**（`type OperatorID int64`）。別の単位の ID を取り違えないためです。関数の引数も、この名前付きの型で書きます（本文 6 節の一覧の例と合わせます）。
- **外向けの ID** は、旧 API の JSON の形のまま使います。仕様の型に従い、オペレーターの ID は JSON の整数（`integer`）、PXR-ID は JSON の文字列（`string`）です。内部の主キーを外に出すかどうかは、単位ごとに旧 API の仕様で決めます。
- **PXR-ID** は業務上の識別子として扱い、主キーにしません。旧スキーマの `varchar(255)` のまま、その列で検索します。
- **IRI 形式**は、今は採用しません。旧 API の仕様に IRI 形式の ID は無く、採用すると API の形が変わるためです（AGENTS.md の「互換性」）。IRI を使った新しい API の形は、フェーズ4（新しい API の設計）で改めて決めます。
- **SQL**は、フェーズ3で CockroachDB でも `GENERATED BY DEFAULT AS IDENTITY` が使えるかを確かめます（`docs/conventions.md` の SQL の禁止事項も同時に確認します）。

### 既定の選択

- 仕様と旧実装の型が違うときは、仕様に合わせ、差分を [差分台帳](spec-deviations.md) に書きます。確かめた例：旧 `user/info` の `pxrId` の検索パラメータは仕様で `number` と書かれており、PXR-ID（文字列）と合っていません。どちらに合わせるかは要確認です。`openapi/` を作るとき（MIGRATION.md の 6.1）に、差分台帳へ要確認の項目として書きます。

### 見本

```go
// bootstrap/internal/operator/id.go（見本）
type OperatorID int64 // 旧 operator の id（bigint）。JSON では integer
type PxrID string     // 旧 pxr_id（varchar(255)）。JSON では string
```

```sql
-- bootstrap/internal/operator/store/query.sql（sqlc の入力。テーブル名は旧スキーマに合わせる）
-- name: FindOperatorByID :one
SELECT id, pxr_id FROM operator WHERE id = $1;
```

見本の置き場：`bootstrap/internal/<単位>/id.go`、`bootstrap/internal/<単位>/store/query.sql`（P0-4 で作ります）。

## 8. エラーの形

### 決定（旧実装と同じ形にする）

- 応答の本文は JSON で、`status`（HTTP のステータスと同じ数値）を必ず入れます。
- 旧実装（[pxr-operator-service の GlobalErrorHandler](https://github.com/Personal-Data-Linkage-Module/pxr-operator-service/blob/453e6385a7545be45d6dccc284004457cb8d28f6/src/resources/handler/GlobalErrorHandler.ts)）の場面ごとの形は次のとおりです。

| 場面 | HTTP | 本文 |
| --- | --- | --- |
| 業務のエラー | 業務のエラーごとのステータス | `{"status": 同じ数値, "message": "文"}` |
| 入力の検証の失敗 | 400 | `{"status": 400, "reasons": [{"property": "項目", "value": 値（無ければ null）, "message": "文"}]}` |
| JSON の形式が不正 | 400 | `{"status": 400, "message": "リクエストボディが、JSON形式ではありません"}` |
| CSRF トークンの検証の失敗 | 403 | `{"status": 403, "message": "不正な CSRF トークンまたは CSRF トークンがありません"}` |
| 上のどれでもない想定外のエラー | 503 | `{"status": 503, "message": "未定義のエラーが発生しました"}` |

- 文言は、旧実装の `config/message.json` の日本語を、同じ文で使います。
- 想定外のエラーは 500 ではなく **503** にします（旧実装どおり）。
- **内部の原因（エラーの中身、スタックトレース、SQL の文など）は本文に出しません**。ログに出すだけです（9 節）。
- 旧実装の `httpCode` の分岐（本文の `message` に例外の名前を入れる）は写しません。仕様に無い形のため、写す必要があるかは、差分台帳で要確認として扱います。

### 既定の選択

- 業務のエラーは、`Status`・`Message`・`Err`（内部の原因）を持つ型で表し、応答の書き出しは `platform/` の一箇所だけで行います。ハンドラごとに書きません。

### 見本

```go
// bootstrap/internal/platform/errors.go（見本）
// 旧実装の GlobalErrorHandler と同じ形で返す。
type AppError struct {
	Status  int    // HTTP のステータス（業務ごと）
	Message string // 利用者に返す文（旧 message.json と同じ文）
	Err     error  // 内部の原因。ログに出すだけで、本文には出さない
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

type errorBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// WriteError は、err を旧実装と同じ形の JSON で返す。
func WriteError(w http.ResponseWriter, err error) {
	status, msg := http.StatusServiceUnavailable, "未定義のエラーが発生しました"
	var appErr *AppError
	if errors.As(err, &appErr) {
		status, msg = appErr.Status, appErr.Message
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorBody{Status: status, Message: msg})
}
```

見本の置き場：`bootstrap/internal/platform/errors.go`、エンドポイントのテストは `bootstrap/internal/<単位>/handler_test.go`（5 節）。

## 9. ログの形

### 決定

- **出力先**は標準出力だけにします。1行に1件の JSON（`log/slog` の `slog.NewJSONHandler`）で出します。ログファイルの日付管理は、アプリでは行いません（保管は、コンテナやホストの仕組みに任せます）。
- **1行の項目**は次のとおりです。
  - `time`、`level`、`msg`（英語の短い固定の文。検索しやすくするため）
  - `kind`：`access`（リクエストの記録）、`application`（業務の記録）、`system`（起動や停止など）のどれか
  - `request_id`：リクエストごとに採番する 32 桁の 16 進数。ミドルウェアで採番し、`context` で渡します
  - `err`：エラーがあるときだけ、その文字列
- **レベル**：利用者の入力の誤りや業務のエラーは `WARN`、想定外のエラーは `ERROR`、通常の記録は `INFO` にします。
- **書かないもの**：パスワード、セッションの ID、JWT、ワンタイムコード、SMS の検証コード、利用者の個人情報（氏名、電話番号など）。認証に使える値は、どのレベルでも書きません。
- **旧実装の 5 種類のログ**（system、http、access、application、performance）は、`kind` の 3 種類にまとめます。performance は、フェーズ2で計測の方法を決めるときに改めて決めます。

### 既定の選択

- 旧実装のログの文言は英語と日本語が混ざっていたため、`msg` は英語に統一します。利用者向けの文（応答の `message`）は日本語のままです。

### 出典

- 旧実装のログの設定：[pxr-operator-service の log4js.config.json](https://github.com/Personal-Data-Linkage-Module/pxr-operator-service/blob/453e6385a7545be45d6dccc284004457cb8d28f6/config/log4js.config.json)、[Logging.ts](https://github.com/Personal-Data-Linkage-Module/pxr-operator-service/blob/453e6385a7545be45d6dccc284004457cb8d28f6/src/resources/config/Logging.ts)（リクエストごとの採番の考え方を参考にしました）

### 見本

```go
// bootstrap/internal/platform/log.go（見本）
func NewLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

// リクエストごとの ID を context に入れる（ミドルウェア）
type requestIDKey struct{}

func WithRequestID(ctx context.Context) context.Context {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return context.WithValue(ctx, requestIDKey{}, hex.EncodeToString(b[:]))
}
```

```go
// 呼び出し例（bootstrap/internal/operator/service.go）
logger.WarnContext(ctx, "operator not found", "kind", "application", "err", err)
```

出力例（1行）：

```json
{"time":"2026-10-09T08:00:00Z","level":"WARN","msg":"operator not found","kind":"application","request_id":"4f1c…","err":"operator not found"}
```

見本の置き場：`bootstrap/internal/platform/log.go`（P0-4 で作ります）。
