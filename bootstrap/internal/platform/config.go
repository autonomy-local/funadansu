package platform

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Config は起動時の設定です。main で一度だけ読み、各部品にはコンストラクタの引数で渡します。
type Config struct {
	Addr               string        // 待ち受け先。既定は localhost（ADR 0007）
	DatabaseURL        string        // PostgreSQL の接続先（必須）
	DBMaxConns         int           // 接続プールの上限
	DBStatementTimeout time.Duration // 接続ごとの statement_timeout
	ReadHeaderTimeout  time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	ShutdownTimeout    time.Duration // 終了の時限（Cloud Run の猶予 10 秒に収める）
	ReadyTimeout       time.Duration // /readyz の DB の ping の時限
	MaxBodyBytes       int64         // 本文の上限
	MemoryLimitRatio   float64       // GOMEMLIMIT の比率（cgroup の上限に掛ける）
	MemoryLimitBytes   int64         // GOMEMLIMIT の明示の値。0 は未指定
}

// LoadConfig は getenv（通常は os.Getenv）から設定を読み、検査します。
// エラーには値を出さず、キーの名前だけを出します。
func LoadConfig(getenv func(string) string) (Config, error) {
	r := &envReader{getenv: getenv}
	cfg := Config{
		Addr:               r.str("FUNADANSU_ADDR", "127.0.0.1:8080"),
		DatabaseURL:        r.str("FUNADANSU_DATABASE_URL", ""),
		DBMaxConns:         int(r.integer("FUNADANSU_DB_MAX_CONNS", 5)),
		DBStatementTimeout: r.duration("FUNADANSU_DB_STATEMENT_TIMEOUT", 5*time.Second),
		ReadHeaderTimeout:  r.duration("FUNADANSU_READ_HEADER_TIMEOUT", 5*time.Second),
		ReadTimeout:        r.duration("FUNADANSU_READ_TIMEOUT", 15*time.Second),
		WriteTimeout:       r.duration("FUNADANSU_WRITE_TIMEOUT", 30*time.Second),
		IdleTimeout:        r.duration("FUNADANSU_IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout:    r.duration("FUNADANSU_SHUTDOWN_TIMEOUT", 8*time.Second),
		ReadyTimeout:       r.duration("FUNADANSU_READY_TIMEOUT", 2*time.Second),
		MaxBodyBytes:       r.integer("FUNADANSU_MAX_BODY_BYTES", 1<<20),
		MemoryLimitRatio:   r.ratio("FUNADANSU_MEMORY_LIMIT_RATIO", 0.9),
		MemoryLimitBytes:   r.integer("FUNADANSU_MEMORY_LIMIT_BYTES", 0),
	}
	if r.err != nil {
		return Config{}, r.err
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("FUNADANSU_DATABASE_URL が必要です")
	}
	if cfg.DBMaxConns < 1 {
		return Config{}, errors.New("FUNADANSU_DB_MAX_CONNS は 1 以上にしてください")
	}
	return cfg, nil
}

// envReader は、最初の不正な値のエラーを覚えておき、読み込みを続けます。
type envReader struct {
	getenv func(string) string
	err    error
}

func (r *envReader) str(key, def string) string {
	if v := r.getenv(key); v != "" {
		return v
	}
	return def
}

func (r *envReader) integer(key string, def int64) int64 {
	v := r.getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil || n < 0 {
		r.fail(key)
		return def
	}
	return n
}

func (r *envReader) duration(key string, def time.Duration) time.Duration {
	v := r.getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		r.fail(key)
		return def
	}
	return d
}

func (r *envReader) ratio(key string, def float64) float64 {
	v := r.getenv(key)
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || f <= 0 || f > 1 {
		r.fail(key)
		return def
	}
	return f
}

func (r *envReader) fail(key string) {
	if r.err == nil {
		r.err = fmt.Errorf("%s の形が不正です", key)
	}
}
