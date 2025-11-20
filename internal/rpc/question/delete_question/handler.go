package delete_question

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	entQ "test-question/internal/entity/question"
	"test-question/internal/pkg/rpc"
	"test-question/internal/pkg/rpc/rpc_auth"
)

//go:generate mockery --name=useCase --output=mocks --outpkg=mocks --exported
type (
	useCase interface {
		DeleteQuestion(ctx context.Context, questionID int, userID string) error
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
	questionID, err := strconv.Atoi(idStr)
	if err != nil {
		rpc.WriteBadRequest(w, "invalid_question_id")
		return
	}

	userID := rpc_auth.GetUserID(r.Context())
	if userID == "" {
		rpc.WriteUnauthorized(w)
		return
	}

	err = h.uc.DeleteQuestion(r.Context(), questionID, userID)
	if err != nil {
		switch {
		case errors.Is(err, entQ.ErrQuestionNotFound):
			rpc.WriteNotFound(w, "question_not_found")
			return

		case errors.Is(err, entQ.ErrAccessDenied):
			rpc.WriteForbidden(w)
			return

		default:
			rpc.WriteUnexpectedError(w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
