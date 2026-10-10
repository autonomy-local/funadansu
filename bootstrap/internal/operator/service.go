package operator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/autonomy-local/funadansu/bootstrap/internal/operator/store"
	"github.com/autonomy-local/funadansu/bootstrap/internal/platform"
)

// ErrNotFound は、オペレーターが見つからないときのエラーです（404）。
var ErrNotFound = &platform.AppError{Status: http.StatusNotFound, Message: "オペレーターが見つかりません"}

// Store は、このパッケージが使う DB の操作です（使う側で定義する）。
type Store interface {
	FindOperatorByID(ctx context.Context, id int64) (store.FindOperatorByIDRow, error)
}

// NewStore は、DB の接続から Store を作ります。store/ は、このパッケージの中からだけ使います。
func NewStore(db *sql.DB) Store {
	return store.New(db)
}

type Service struct {
	store  Store
	logger *slog.Logger
}

func NewService(store Store, logger *slog.Logger) *Service {
	return &Service{store: store, logger: logger}
}

// FindOperator は、オペレーターを内部の ID で探します。
func (s *Service) FindOperator(ctx context.Context, id OperatorID) (Operator, error) {
	row, err := s.store.FindOperatorByID(ctx, int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		s.logger.WarnContext(ctx, "operator not found", "kind", platform.KindApplication, "id", id)
		return Operator{}, ErrNotFound
	}
	if err != nil {
		return Operator{}, fmt.Errorf("find operator: %w", err)
	}
	return Operator{ID: OperatorID(row.ID), PxrID: PxrID(row.PxrID)}, nil
}
