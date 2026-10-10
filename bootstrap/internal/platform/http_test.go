package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChainOrderAndRecover(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /boom", func(http.ResponseWriter, *http.Request) { panic("secret detail") })
	mux.HandleFunc("GET /ok", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusTeapot) })
	h := Chain(mux, AccessLog(logger), Recover(logger))

	t.Run("panic は 503 の形で返り、system の記録に出る", func(t *testing.T) {
		buf.Reset()
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/boom", nil))
		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("status = %d, want 503", rec.Code)
		}
		if strings.Contains(rec.Body.String(), "secret detail") {
			t.Errorf("本文に内部の原因が出ている: %s", rec.Body)
		}
		logs := decodeLines(t, buf.String())
		if len(logs) != 2 {
			t.Fatalf("log lines = %d, want 2 (panic と access)", len(logs))
		}
		if logs[0]["kind"] != KindSystem || logs[0]["level"] != "ERROR" {
			t.Errorf("panic の記録 = %v", logs[0])
		}
	})

	t.Run("access の記録は status と request_id を持つ", func(t *testing.T) {
		buf.Reset()
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ok", nil))
		logs := decodeLines(t, buf.String())
		if len(logs) != 1 {
			t.Fatalf("log lines = %d, want 1", len(logs))
		}
		rec2 := logs[0]
		if rec2["kind"] != KindAccess || rec2["path"] != "/ok" || rec2["status"] != float64(http.StatusTeapot) {
			t.Errorf("access の記録 = %v", rec2)
		}
		if id, _ := rec2["request_id"].(string); len(id) != 32 {
			t.Errorf("request_id = %q, want 32 桁", id)
		}
	})
}

func TestLimitBody(t *testing.T) {
	h := Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		if err != nil {
			WriteError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}), LimitBody(4))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", strings.NewReader("12345678")))
	if rec.Code == http.StatusNoContent {
		t.Errorf("上限を超えた本文が通った")
	}
}

func TestHealthAndReady(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	t.Run("healthz は DB に触れずに 200", func(t *testing.T) {
		rec := httptest.NewRecorder()
		Healthz().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want 200", rec.Code)
		}
	})
	t.Run("readyz は ping が失敗すると 503", func(t *testing.T) {
		rec := httptest.NewRecorder()
		Readyz(failingPinger{}, 0, logger).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("status = %d, want 503", rec.Code)
		}
	})
	t.Run("readyz は ping が通れば 200", func(t *testing.T) {
		rec := httptest.NewRecorder()
		Readyz(okPinger{}, 0, logger).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want 200", rec.Code)
		}
	})
}

type failingPinger struct{}

func (failingPinger) PingContext(context.Context) error { return errors.New("down") }

type okPinger struct{}

func (okPinger) PingContext(context.Context) error { return nil }

func decodeLines(t *testing.T, s string) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		if line == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("log line is not JSON: %v: %s", err, line)
		}
		out = append(out, m)
	}
	return out
}
