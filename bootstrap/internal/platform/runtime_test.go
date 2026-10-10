package platform

import (
	"errors"
	"math"
	"runtime/debug"
	"testing"
)

// ConfigureMemoryLimit は、プロセス全体の GOMEMLIMIT を変えるため、テストの後で戻す。
func restoreMemoryLimit(t *testing.T) {
	orig := debug.SetMemoryLimit(-1)
	t.Cleanup(func() { debug.SetMemoryLimit(orig) })
}

func TestConfigureMemoryLimit(t *testing.T) {
	noFile := func(string) ([]byte, error) { return nil, errors.New("no file") }
	cgroup := func(v string) func(string) ([]byte, error) {
		return func(string) ([]byte, error) { return []byte(v + "\n"), nil }
	}
	tests := []struct {
		name       string
		cfg        Config
		getenv     func(string) string
		readFile   func(string) ([]byte, error)
		wantLimit  int64
		wantSource string
	}{
		{
			name:       "GOMEMLIMIT が既にあれば上書きしない",
			cfg:        Config{MemoryLimitBytes: 1000},
			getenv:     env(map[string]string{"GOMEMLIMIT": "2GiB"}),
			readFile:   noFile,
			wantLimit:  math.MaxInt64,
			wantSource: MemorySourceEnv,
		},
		{
			name:       "明示のバイト数",
			cfg:        Config{MemoryLimitBytes: 1000},
			getenv:     env(nil),
			readFile:   noFile,
			wantLimit:  1000,
			wantSource: MemorySourceConfig,
		},
		{
			name:       "cgroup の上限に比率を掛ける",
			cfg:        Config{MemoryLimitRatio: 0.5},
			getenv:     env(nil),
			readFile:   cgroup("1000"),
			wantLimit:  500,
			wantSource: MemorySourceCgroup,
		},
		{
			name:       "cgroup の上限が無い（max）",
			cfg:        Config{MemoryLimitRatio: 0.9},
			getenv:     env(nil),
			readFile:   cgroup("max"),
			wantLimit:  math.MaxInt64,
			wantSource: MemorySourceNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restoreMemoryLimit(t)
			limit, source := ConfigureMemoryLimit(tt.cfg, tt.getenv, tt.readFile)
			if limit != tt.wantLimit || source != tt.wantSource {
				t.Errorf("got (%d, %s), want (%d, %s)", limit, source, tt.wantLimit, tt.wantSource)
			}
		})
	}
}
