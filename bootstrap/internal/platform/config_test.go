package platform

import (
	"strings"
	"testing"
	"time"
)

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr string // 空なら成功
		check   func(t *testing.T, cfg Config)
	}{
		{
			name: "既定の値",
			env:  map[string]string{"FUNADANSU_DATABASE_URL": "postgres://localhost/x"},
			check: func(t *testing.T, cfg Config) {
				if cfg.Addr != "127.0.0.1:8080" {
					t.Errorf("Addr = %q, want localhost の既定", cfg.Addr)
				}
				if cfg.ShutdownTimeout != 8*time.Second {
					t.Errorf("ShutdownTimeout = %v, want 8s", cfg.ShutdownTimeout)
				}
				if cfg.MemoryLimitRatio != 0.9 {
					t.Errorf("MemoryLimitRatio = %v, want 0.9", cfg.MemoryLimitRatio)
				}
			},
		},
		{
			name: "値を上書きする",
			env: map[string]string{
				"FUNADANSU_DATABASE_URL":     "postgres://localhost/x",
				"FUNADANSU_ADDR":             "0.0.0.0:9000",
				"FUNADANSU_SHUTDOWN_TIMEOUT": "3s",
				"FUNADANSU_MAX_BODY_BYTES":   "2048",
			},
			check: func(t *testing.T, cfg Config) {
				if cfg.Addr != "0.0.0.0:9000" || cfg.ShutdownTimeout != 3*time.Second || cfg.MaxBodyBytes != 2048 {
					t.Errorf("上書きが反映されない: %+v", cfg)
				}
			},
		},
		{
			name:    "接続先が無い",
			env:     map[string]string{},
			wantErr: "FUNADANSU_DATABASE_URL",
		},
		{
			name:    "時間の形が不正。値は出さない",
			env:     map[string]string{"FUNADANSU_DATABASE_URL": "postgres://localhost/x", "FUNADANSU_READ_TIMEOUT": "secret-value"},
			wantErr: "FUNADANSU_READ_TIMEOUT",
		},
		{
			name:    "比率が範囲外",
			env:     map[string]string{"FUNADANSU_DATABASE_URL": "postgres://localhost/x", "FUNADANSU_MEMORY_LIMIT_RATIO": "1.5"},
			wantErr: "FUNADANSU_MEMORY_LIMIT_RATIO",
		},
		{
			name:    "接続の上限が 0",
			env:     map[string]string{"FUNADANSU_DATABASE_URL": "postgres://localhost/x", "FUNADANSU_DB_MAX_CONNS": "0"},
			wantErr: "FUNADANSU_DB_MAX_CONNS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadConfig(env(tt.env))
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("err = nil, want %q を含むエラー", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("err = %q, want %q を含む", err, tt.wantErr)
				}
				if strings.Contains(err.Error(), "secret-value") {
					t.Errorf("エラーに値が出ている: %q", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			tt.check(t, cfg)
		})
	}
}
