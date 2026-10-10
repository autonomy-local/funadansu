package platform

import (
	"encoding/json"
	"errors"
	"net/http"
)

// 旧実装の GlobalErrorHandler と同じ形で応答します（docs/conventions.md の 8 節）。
// 想定外のエラーは 500 ではなく 503 にします。

// ErrUnexpectedMessage は、想定外のエラーの利用者向けの文です。
const ErrUnexpectedMessage = "未定義のエラーが発生しました"

// AppError は、業務のエラーです。Err は内部の原因で、ログに出すだけで本文には出しません。
type AppError struct {
	Status  int
	Message string
	Err     error
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

// Reason は、入力の検証で1つの項目の誤りを表します。
type Reason struct {
	Property string `json:"property"`
	Value    any    `json:"value"`
	Message  string `json:"message"`
}

// ValidationError は、入力の検証の失敗です（400）。
type ValidationError struct {
	Reasons []Reason
}

func (e *ValidationError) Error() string { return "validation failed" }

type messageBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type reasonsBody struct {
	Status  int      `json:"status"`
	Reasons []Reason `json:"reasons"`
}

// WriteError は err を、旧実装と同じ形の JSON で書き出します。
// 業務のエラーと入力の検証以外は、想定外のエラー（503）として扱います。
func WriteError(w http.ResponseWriter, err error) {
	var validation *ValidationError
	if errors.As(err, &validation) {
		writeJSON(w, http.StatusBadRequest, reasonsBody{Status: http.StatusBadRequest, Reasons: validation.Reasons})
		return
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		writeJSON(w, appErr.Status, messageBody{Status: appErr.Status, Message: appErr.Message})
		return
	}
	writeJSON(w, http.StatusServiceUnavailable, messageBody{Status: http.StatusServiceUnavailable, Message: ErrUnexpectedMessage})
}

// WriteJSON は、成功の応答を JSON で書き出します。
func WriteJSON(w http.ResponseWriter, status int, body any) {
	writeJSON(w, status, body)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
