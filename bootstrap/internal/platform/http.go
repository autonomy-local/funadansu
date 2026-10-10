package platform

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// Middleware は http.Handler を包みます。
type Middleware func(http.Handler) http.Handler

// Chain は mws を、先頭が一番外側になるように h に包みます。
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// NewServer は、タイムアウトを必ず設定した http.Server を返します（docs/conventions.md の 10 節）。
func NewServer(cfg Config, h http.Handler) *http.Server {
	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           h,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
}

// AccessLog は、リクエストごとに request_id を採番し、access の記録を1行出します。
func AccessLog(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := WithRequestID(r.Context())
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()
			next.ServeHTTP(rec, r.WithContext(ctx))
			logger.InfoContext(ctx, "request",
				"kind", KindAccess,
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}

// Recover は、panic を system の記録に出し、想定外のエラーの形（503）で返します。
func Recover(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if v := recover(); v != nil {
					logger.ErrorContext(r.Context(), "panic recovered",
						"kind", KindSystem,
						"err", fmt.Sprint(v),
						"stack", string(debug.Stack()),
					)
					WriteError(w, errors.New("panic"))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// LimitBody は、本文の上限を決めます。proxy の制限に頼らず、サービスでも上限を持ちます。
func LimitBody(n int64) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
