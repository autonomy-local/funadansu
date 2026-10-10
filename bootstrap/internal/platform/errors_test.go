package platform

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "業務のエラー",
			err:        &AppError{Status: 404, Message: "見つからない", Err: errors.New("内部の原因")},
			wantStatus: 404,
			wantBody:   `{"status":404,"message":"見つからない"}`,
		},
		{
			name:       "業務のエラーを包んでも同じ",
			err:        fmt.Errorf("wrap: %w", &AppError{Status: 409, Message: "重複"}),
			wantStatus: 409,
			wantBody:   `{"status":409,"message":"重複"}`,
		},
		{
			name:       "入力の検証の失敗",
			err:        &ValidationError{Reasons: []Reason{{Property: "id", Value: "abc", Message: "整数"}}},
			wantStatus: 400,
			wantBody:   `{"status":400,"reasons":[{"property":"id","value":"abc","message":"整数"}]}`,
		},
		{
			name:       "想定外のエラーは 503 で、内部の原因を出さない",
			err:        errors.New("pq: password authentication failed for user x"),
			wantStatus: 503,
			wantBody:   `{"status":503,"message":"未定義のエラーが発生しました"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			WriteError(rec, tt.err)
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
				t.Errorf("Content-Type = %q", got)
			}
			assertJSONEqual(t, rec.Body.String(), tt.wantBody)
		})
	}
}

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
