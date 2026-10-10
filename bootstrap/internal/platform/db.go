package platform

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

// OpenDB は PostgreSQL の接続プールを開きます。接続はまだ作りません（準備の確認は /readyz）。
// 接続ごとに statement_timeout を設定し、プールの上限は cfg で決めます。
func OpenDB(cfg Config) (*sql.DB, error) {
	conn, err := pgx.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		// 接続先の文字列（パスワードを含む）は、エラーにもログにも出しません。
		return nil, errors.New("FUNADANSU_DATABASE_URL の形が不正です")
	}
	conn.RuntimeParams["statement_timeout"] = strconv.FormatInt(cfg.DBStatementTimeout.Milliseconds(), 10)

	db := stdlib.OpenDB(*conn)
	db.SetMaxOpenConns(cfg.DBMaxConns)
	db.SetMaxIdleConns(cfg.DBMaxConns)
	return db, nil
}
