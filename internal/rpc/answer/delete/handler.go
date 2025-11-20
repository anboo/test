package delete //nolint:predeclared

import (
	"context"
	"net/http"
	"strconv"

	"test-question/internal/pkg/rpc"
	"test-question/internal/pkg/rpc/rpc_auth"

	entA "test-question/internal/entity/answer"

	"github.com/pkg/errors"
)

//go:generate mockery --name=useCase --output=mocks --outpkg=mocks --exported
type (
	useCase interface {
		DeleteAnswer(ctx context.Context, answerID int, userID string) error
	}
)

type Handler struct {
	uc useCase
}

func NewHandler(uc useCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	answerID, err := strconv.Atoi(idStr)
	if err != nil {
		rpc.WriteBadRequest(w, "invalid answer id")
		return
	}

	userID := rpc_auth.GetUserID(r.Context())
	if userID == "" {
		rpc.WriteUnauthorized(w)
		return
	}

	err = h.uc.DeleteAnswer(r.Context(), answerID, userID)
	if err != nil {
		switch {
		case errors.Is(err, entA.ErrAnswerNotFound):
			rpc.WriteNotFound(w, "answer_not_found")
			return

		case errors.Is(err, entA.ErrAccessDenied):
			rpc.WriteForbidden(w)
			return

		default:
			rpc.WriteUnexpectedError(w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
