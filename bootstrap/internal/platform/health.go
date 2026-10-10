package platform

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// Pinger は、DB に ping できることを確かめます（*sql.DB が満たします）。
type Pinger interface {
	PingContext(ctx context.Context) error
}

// Healthz は生存の確認です。DB には触れません（docs/conventions.md の 11 節）。
func Healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		WriteJSON(w, http.StatusOK, messageBody{Status: http.StatusOK, Message: "ok"})
	})
}

// Readyz は DB の ping で準備の確認をします。時限は短く、超えたら 503 です。
func Readyz(db Pinger, timeout time.Duration, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			logger.WarnContext(r.Context(), "readiness check failed", "kind", KindSystem, "err", err.Error())
			WriteError(w, err)
			return
		}
		WriteJSON(w, http.StatusOK, messageBody{Status: http.StatusOK, Message: "ok"})
	})
}
