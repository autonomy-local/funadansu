package operator

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/autonomy-local/funadansu/bootstrap/internal/platform"
)

// Handler は、この単位の HTTP の入口です。
type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// Register は、見本のエンドポイントを mux につなぎます。
// パスは bootstrap だけの見本です。旧 API の形は、この単位の本物のエンドポイントで写します。
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /bootstrap/operators/{id}", h.find)
}

type operatorBody struct {
	ID    OperatorID `json:"id"`
	PxrID PxrID      `json:"pxrId"`
}

func (h *Handler) find(w http.ResponseWriter, r *http.Request) {
	raw := r.PathValue("id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id < 1 {
		platform.WriteError(w, &platform.ValidationError{Reasons: []platform.Reason{
			{Property: "id", Value: raw, Message: "1 以上の整数で指定してください"},
		}})
		return
	}

	op, err := h.service.FindOperator(r.Context(), OperatorID(id))
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			h.logger.ErrorContext(r.Context(), "find operator failed", "kind", platform.KindApplication, "err", err.Error())
		}
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, operatorBody{ID: op.ID, PxrID: op.PxrID})
}
