package operator

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
)

func TestFindOperator(t *testing.T) {
	db := testDB(t)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	existing, existingPxrID := insertOperator(t, db)
	deleted, _ := insertOperator(t, db)
	if _, err := db.ExecContext(context.Background(), `DELETE FROM pxr_operator.operator WHERE id = $1`, deleted); err != nil {
		t.Fatalf("delete: %v", err)
	}

	tests := []struct {
		name    string
		id      OperatorID
		want    Operator
		wantErr error
	}{
		{name: "存在する", id: existing, want: Operator{ID: existing, PxrID: existingPxrID}},
		{name: "存在しない", id: deleted, wantErr: ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(NewStore(db), logger)
			got, err := svc.FindOperator(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got = %+v, want %+v", got, tt.want)
			}
		})
	}
}
