package platform

import (
	"log/slog"
	"math"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

// cgroupMemoryMax は cgroup v2 のメモリの上限のファイルです。
const cgroupMemoryMax = "/sys/fs/cgroup/memory.max"

// 割り当てのバイト数の出所（docs/conventions.md の 11 節）。
const (
	MemorySourceEnv    = "env"    // 環境変数 GOMEMLIMIT が既にあり、上書きしない
	MemorySourceConfig = "config" // FUNADANSU_MEMORY_LIMIT_BYTES
	MemorySourceCgroup = "cgroup" // cgroup の上限に比率を掛けた値
	MemorySourceNone   = "none"   // 上限を設定しない
)

// ConfigureMemoryLimit は GOMEMLIMIT を決めて設定し、その値と出所を返します。
// 返す値が math.MaxInt64 のときは上限なし（off）です。
// readFile は cgroup の上限を読むための関数です（テストで差し替える）。
func ConfigureMemoryLimit(cfg Config, getenv func(string) string, readFile func(string) ([]byte, error)) (int64, string) {
	if getenv("GOMEMLIMIT") != "" {
		// Go が起動時に読んだ値をそのまま使う。SetMemoryLimit(-1) は値を変えずに現在の値を返す。
		return debug.SetMemoryLimit(-1), MemorySourceEnv
	}
	if cfg.MemoryLimitBytes > 0 {
		debug.SetMemoryLimit(cfg.MemoryLimitBytes)
		return cfg.MemoryLimitBytes, MemorySourceConfig
	}
	if b, err := readFile(cgroupMemoryMax); err == nil {
		if max, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 64); err == nil {
			limit := int64(math.Floor(max * cfg.MemoryLimitRatio))
			debug.SetMemoryLimit(limit)
			return limit, MemorySourceCgroup
		}
	}
	return debug.SetMemoryLimit(-1), MemorySourceNone
}

// LogStartup は、起動時の版と実行の設定を system の記録として1行で出します。
func LogStartup(logger *slog.Logger, memoryLimit int64, memorySource string) {
	gomemlimit := "off"
	if memoryLimit != math.MaxInt64 {
		gomemlimit = strconv.FormatInt(memoryLimit, 10)
	}
	logger.Info("startup",
		"kind", KindSystem,
		"go_version", runtime.Version(),
		"gomaxprocs", runtime.GOMAXPROCS(0),
		"gomemlimit", gomemlimit,
		"memory_source", memorySource,
	)
}

// ReadFile は、通常のファイル読み込みです（ConfigureMemoryLimit に渡す）。
func ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
