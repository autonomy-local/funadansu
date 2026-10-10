package platform

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log/slog"
)

// ログの kind（docs/conventions.md の 9 節）。
const (
	KindAccess      = "access"      // リクエストの記録
	KindApplication = "application" // 業務の記録
	KindSystem      = "system"      // 起動や停止など
)

type requestIDKey struct{}

// NewLogger は、標準出力へ1行1件の JSON で書くロガーを返します。
// request_id は、context に入っていれば自動で付きます。
func NewLogger(w io.Writer) *slog.Logger {
	return slog.New(&requestIDHandler{
		Handler: slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}),
	})
}

// WithRequestID は、32 桁の 16 進数の request_id を採番して context に入れます。
func WithRequestID(ctx context.Context) context.Context {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return context.WithValue(ctx, requestIDKey{}, hex.EncodeToString(b[:]))
}

// RequestID は、context の request_id を返します。無ければ空の文字列です。
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// requestIDHandler は、context の request_id を各行に足します。
type requestIDHandler struct {
	slog.Handler
}

func (h *requestIDHandler) Handle(ctx context.Context, r slog.Record) error {
	if id := RequestID(ctx); id != "" {
		r.AddAttrs(slog.String("request_id", id))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *requestIDHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &requestIDHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *requestIDHandler) WithGroup(name string) slog.Handler {
	return &requestIDHandler{Handler: h.Handler.WithGroup(name)}
}
