// Funadansu の入口です。配線（依存の組み立て）は、この main だけで行います（docs/conventions.md の 3 節）。
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/autonomy-local/funadansu/bootstrap/internal/operator"
	"github.com/autonomy-local/funadansu/bootstrap/internal/platform"
)

func main() {
	if err := run(); err != nil {
		// 起動の失敗は、ログの形になる前に標準エラーへ出します。値は出しません（キーの名前だけ）。
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := platform.LoadConfig(os.Getenv)
	if err != nil {
		return err
	}
	logger := platform.NewLogger(os.Stdout)
	limit, source := platform.ConfigureMemoryLimit(cfg, os.Getenv, platform.ReadFile)
	platform.LogStartup(logger, limit, source)

	db, err := platform.OpenDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	operators := operator.NewService(operator.NewStore(db), logger)

	mux := http.NewServeMux()
	mux.Handle("GET /healthz", platform.Healthz())
	mux.Handle("GET /readyz", platform.Readyz(db, cfg.ReadyTimeout, logger))
	operator.NewHandler(operators, logger).Register(mux)

	handler := platform.Chain(mux,
		platform.AccessLog(logger),
		platform.Recover(logger),
		platform.LimitBody(cfg.MaxBodyBytes),
	)
	srv := platform.NewServer(cfg, handler)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errc := make(chan error, 1)
	go func() {
		logger.Info("listening", "kind", platform.KindSystem, "addr", cfg.Addr)
		errc <- srv.ListenAndServe()
	}()

	select {
	case err := <-errc:
		// 待ち受けが始まらなかったときなど、サーバーが止まった。
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		logger.Info("shutting down", "kind", platform.KindSystem)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
	}
	logger.Info("stopped", "kind", platform.KindSystem)
	return nil
}
