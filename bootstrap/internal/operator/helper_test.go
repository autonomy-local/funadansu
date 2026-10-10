package operator

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/autonomy-local/funadansu/bootstrap/internal/operator/store"
	"github.com/autonomy-local/funadansu/bootstrap/internal/platform"
)

// testDB は、db-start で起動するローカルの PostgreSQL につなぎます（モックは使わない）。
// 接続先は FUNADANSU_TEST_DATABASE_URL で変えられます。
func testDB(t *testing.T) *sql.DB {
	t.Helper()
	url := os.Getenv("FUNADANSU_TEST_DATABASE_URL")
	if url == "" {
		url = "postgres://postgres@localhost:5432/funadansu?sslmode=disable"
	}
	db, err := platform.OpenDB(platform.Config{DatabaseURL: url, DBMaxConns: 2, DBStatementTimeout: 5 * time.Second})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if err := db.PingContext(context.Background()); err != nil {
		t.Fatalf("PostgreSQL に接続できません（先に db-start を実行してください）: %v", err)
	}
	return db
}

// insertOperator は、テスト用のオペレーターを1件入れ、テストの後で消します。
// 投入と削除も sqlc のクエリを使います（store/query.sql の InsertOperatorForTest など）。
func insertOperator(t *testing.T, db *sql.DB) (OperatorID, PxrID) {
	t.Helper()
	q := store.New(db)
	suffix := randomHex(t)
	pxrID := PxrID("pxr-" + suffix)
	id, err := q.InsertOperatorForTest(context.Background(), store.InsertOperatorForTestParams{
		LoginID:            "login-" + suffix,
		PxrID:              string(pxrID),
		UserID:             "user-" + suffix,
		UniqueCheckLoginID: "login-" + suffix,
	})
	if err != nil {
		t.Fatalf("insert operator: %v", err)
	}
	t.Cleanup(func() {
		_ = q.DeleteOperatorForTest(context.Background(), id)
	})
	return OperatorID(id), pxrID
}

// deleteOperator は、テスト用のオペレーターを消す（存在しない ID の確認に使う）。
func deleteOperator(t *testing.T, db *sql.DB, id OperatorID) {
	t.Helper()
	if err := store.New(db).DeleteOperatorForTest(context.Background(), int64(id)); err != nil {
		t.Fatalf("delete operator: %v", err)
	}
}

func randomHex(t *testing.T) string {
	t.Helper()
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		t.Fatalf("rand: %v", err)
	}
	return hex.EncodeToString(b[:])
}
