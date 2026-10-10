package operator

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"os"
	"testing"
	"time"

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
// 旧スキーマは空の表なので、必要な列だけをここで埋めます。
func insertOperator(t *testing.T, db *sql.DB) (OperatorID, PxrID) {
	t.Helper()
	suffix := randomHex(t)
	pxrID := PxrID("pxr-" + suffix)
	var id int64
	err := db.QueryRowContext(context.Background(), `
		INSERT INTO pxr_operator.operator (
			type, login_id, hpassword, pxr_id, user_information, name, mobile_phone, mail,
			auth, attributes, lock_flg, user_id, region_catalog_code, app_catalog_code,
			wf_catalog_code, client_id, created_by, updated_by, unique_check_login_id
		) VALUES (
			0, $1, 'test', $2, '{}', 'test', '', '', '{}', '{}', false, $3, 0, 0, 0,
			'test', 'test', 'test', $4
		) RETURNING id`, "login-"+suffix, string(pxrID), "user-"+suffix, "login-"+suffix).Scan(&id)
	if err != nil {
		t.Fatalf("insert operator: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), `DELETE FROM pxr_operator.operator WHERE id = $1`, id)
	})
	return OperatorID(id), pxrID
}

func randomHex(t *testing.T) string {
	t.Helper()
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		t.Fatalf("rand: %v", err)
	}
	return hex.EncodeToString(b[:])
}
