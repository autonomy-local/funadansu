package operator

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestMux(db *sql.DB) http.Handler {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	mux := http.NewServeMux()
	NewHandler(NewService(NewStore(db), logger), logger).Register(mux)
	return mux
}

func TestFindOperatorEndpoint(t *testing.T) {
	db := testDB(t)
	id, pxrID := insertOperator(t, db)
	h := newTestMux(db)

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "存在する ID は 200 と JSON の形",
			path:       fmt.Sprintf("/bootstrap/operators/%d", id),
			wantStatus: http.StatusOK,
			wantBody:   fmt.Sprintf(`{"id":%d,"pxrId":%q}`, id, pxrID),
		},
		{
			name:       "存在しない ID は 404 の形",
			path:       "/bootstrap/operators/9999999999",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"status":404,"message":"オペレーターが見つかりません"}`,
		},
		{
			name:       "整数でない ID は 400 の形",
			path:       "/bootstrap/operators/abc",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"status":400,"reasons":[{"property":"id","value":"abc","message":"1 以上の整数で指定してください"}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d（本文: %s）", rec.Code, tt.wantStatus, rec.Body)
			}
			assertJSONEqual(t, rec.Body.String(), tt.wantBody)
		})
	}
}

func TestFindOperatorEndpointUnexpectedError(t *testing.T) {
	db := testDB(t)
	id, _ := insertOperator(t, db)
	h := newTestMux(db)
	// 接続を閉じると、DB の操作は想定外のエラーになる。応答は 503 の形で、内部の原因は出ない。
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/bootstrap/operators/%d", id), nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("body is not JSON: %v", err)
	}
	if body["message"] != "未定義のエラーが発生しました" || body["status"] != float64(503) {
		t.Errorf("body = %v", body)
	}
}

// assertJSONEqual は、キーの順を無視して JSON を比べる。
func assertJSONEqual(t *testing.T, got, want string) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal([]byte(got), &g); err != nil {
		t.Fatalf("got is not JSON: %v: %s", err, got)
	}
	if err := json.Unmarshal([]byte(want), &w); err != nil {
		t.Fatalf("want is not JSON: %v", err)
	}
	gb, _ := json.Marshal(g)
	wb, _ := json.Marshal(w)
	if string(gb) != string(wb) {
		t.Errorf("body = %s, want %s", gb, wb)
	}
}
